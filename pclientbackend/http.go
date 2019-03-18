package pclientbackend

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"sync"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

var router1 = &sync.Map{}

//goroutine http proxy
func httpProxy2(conn1 io.ReadWriter) {
	for {
		msg, ok := <-httpChan
		if ok == false {
			return
		}
		if msg.Cmd == CmdHttpReqContinued {
			ch1 := make(chan MsgType, 1)
			cid := binary.BigEndian.Uint32(msg.Msg)
			router1.Store(cid, ch1)
			//star serve ch1
			go proxyChan(ch1, conn1, msg.From, cid)
		} else {
			cid := binary.BigEndian.Uint32(msg.Msg[:4])
			ch1, ok := router1.Load(cid)
			if ok {
				ch1.(chan MsgType) <- msg
			}
		}
	}
}

func proxyChan(ch1 chan MsgType, conn1 io.ReadWriter, from int64, cid uint32) {
	httpConn, err := net.Dial("tcp", fmt.Sprintf("localhost:%d", proxyPort))
	if err != nil {
		log.Printf("localhost:%d :%v\n", proxyPort, err)
		return
	}
	defer httpConn.Close()
	timeout_ch := make(chan int, 1)
	defer router1.Delete(cid)
	defer close(ch1)

	go proxyResopnse(conn1, httpConn, from, timeout_ch, cid)
	for {
		select {
		case rbody, ok := <-ch1:
			if ok == false {
				return
			}
			if rbody.Cmd == CmdHttpReqClose {
				log.Println("recv CmdHttpReqClose")
				return
			}
			if rbody.Cmd == CmdHttpRequest {
				httpConn.Write(rbody.Msg[4:])
			}
		case res := <-timeout_ch:
			if res != 1 {
				return
			}
		case <-time.After(time.Second * 30):
			return
		}
	}
}

func proxyResopnse(conn1 io.ReadWriter, httpConn io.ReadWriter, from int64, timeout_ch chan int, cid uint32) {
	//recv response
	defer close(timeout_ch)
	header := make([]byte, 4)
	binary.BigEndian.PutUint32(header, cid)
	buf := make([]byte, 4000)
	buffer := bytes.NewBufferString("")
	for {
		n, err := httpConn.Read(buf)
		if err != nil {
			log.Printf("ProxyResp:%v\n", err)
			cr, _ := MsgEncode(CmdHttpRespClose, id, from, header)
			conn1.Write(cr)
			timeout_ch <- 0
			return
		}
		if n == 0 {
			break
		}
		buffer.Reset()
		buffer.Write(header)
		buffer.Write(buf[:n])
		r, _ := MsgEncode(CmdHttpRespContinued, id, from, buffer.Bytes())
		conn1.Write(r)
		timeout_ch <- 1
	}
}

//local router
var locRouter = &sync.Map{}
var counter uint32 = rand.Uint32()
var lock1 sync.Mutex

func getConnID() uint32 {
	lock1.Lock()
	counter++
	ret := counter
	lock1.Unlock()
	return ret
}

//goroutine http local serve
func httpResponse2(conn1 io.ReadWriter, locConn net.Conn, to int64) {
	cid := getConnID()
	locRouter.Store(cid, locConn)
	defer locConn.Close()
	defer locRouter.Delete(to)
	header := make([]byte, 4)
	binary.BigEndian.PutUint32(header, cid)
	r, _ := MsgEncode(CmdHttpReqContinued, 0, to, header)
	conn1.Write(r)
	//read and redirect
	buf := make([]byte, 4000)
	buffer := bytes.NewBufferString("")
	for {
		n, err := locConn.Read(buf)
		if err != nil {
			log.Printf("Browser:%v\n", err)
			r, _ := MsgEncode(CmdHttpReqClose, 0, to, header)
			conn1.Write(r)
			return
		}
		if n > 0 {
			buffer.Reset()
			buffer.Write(header)
			buffer.Write(buf[:n])
			r, _ = MsgEncode(CmdHttpRequest, 0, to, buffer.Bytes())
			conn1.Write(r)
		}
	}
}

func httpRespRouter() {
	for {
		msg, ok := <-serveChan
		if ok == false {
			return
		}
		cid := binary.BigEndian.Uint32(msg.Msg[:4])
		iConn, ok := locRouter.Load(cid)
		if ok == false {
			continue
		}
		lConn, ok := iConn.(io.WriteCloser)
		if ok == false {
			continue
		}
		if msg.Cmd == CmdHttpRespContinued {
			lConn.Write(msg.Msg[4:])
		} else if msg.Cmd == CmdHttpRespClose {
			lConn.Close()
		}
	}
}

func clearLocRouter() {
	var ka = []interface{}{}
	var va = []interface{}{}
	locRouter.Range(func(k, v interface{}) bool {
		va = append(va, v)
		ka = append(ka, k)
		return true
	})
	for _, k := range ka {
		locRouter.Delete(k)
	}
	for _, v := range va {
		v.(io.Closer).Close()
	}
}
