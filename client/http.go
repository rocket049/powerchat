package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

var router1 = &sync.Map{}

//goroutine
func httpProxy2(conn1 io.ReadWriter) {
	for {
		msg, ok := <-httpChan
		if ok == false {
			return
		}
		if msg.Cmd == CmdHttpRequest {
			ch1, ok := router1.Load(msg.From)
			if ok {
				ch1.(chan MsgType) <- msg
			} else {
				ch1 := make(chan MsgType, 1)
				router1.Store(msg.From, ch1)
				//star serve ch1
				go proxyChan(ch1, conn1, msg.From)
				ch1 <- msg
			}

		} else {
			ch1, ok := router1.Load(msg.From)
			if ok {
				ch1.(chan MsgType) <- msg
			}
		}
	}
}

func proxyChan(ch1 chan MsgType, conn1 io.ReadWriter, from int64) {
	var httpConn io.ReadWriteCloser
	httpConn, err := net.Dial("tcp", fmt.Sprintf("localhost:%d", proxyPort))
	if err != nil {
		log.Printf("localhost:%d :%v\n", proxyPort, err)
		httpConn = errReader
		//continue
	}
	defer httpConn.Close()
	timeout_ch := make(chan int, 1)
	defer close(timeout_ch)
	defer router1.Delete(from)
	defer close(ch1)

	go proxyResopnse(conn1, httpConn, from, timeout_ch)
	for {
		select {
		case rbody, ok := <-ch1:
			if ok == false {
				return
			}
			if rbody.Cmd != CmdHttpReqClose {
				httpConn.Write(rbody.Msg)
			}
		case res := <-timeout_ch:
			if res == 0 {
				return
			}
		case <-time.After(time.Second * 60):
			return
		}
	}
}

func proxyResopnse(conn1 io.ReadWriter, httpConn io.ReadWriter, from int64, timeout_ch chan int) {
	//recv response
	reader1 := bufio.NewReader(httpConn)
	for {
		//header
		headbuf := bytes.NewBufferString("")
		var size1 int
		for {
			line1, _, err := reader1.ReadLine()
			if err != nil {
				log.Printf("ReadLine:%v\n", err)
				cr, _ := MsgEncode(CmdHttpRespClose, id, from, []byte("Error"))
				conn1.Write(cr)
				timeout_ch <- 0
				return
			}
			headbuf.Write(line1)
			headbuf.WriteString("\r\n")
			s1 := strings.ToLower(string(line1))
			//fmt.Printf("%s\n\r", s1)
			if strings.HasPrefix(s1, "content-length:") {
				sz1, err := strconv.ParseInt(s1[16:], 10, 32)
				if err == nil {
					size1 = int(sz1)
				}
			}
			if len(line1) == 0 {
				break
			}
		}
		r, _ := MsgEncode(CmdHttpRespContinued, id, from, headbuf.Bytes())
		conn1.Write(r)
		timeout_ch <- 1
		//body
		var data []byte
		if size1 > 0 {
			data = make([]byte, size1)
			_, err := io.ReadFull(reader1, data)
			if err != nil {
				log.Printf("ReadFull:%v\n", err)
				cr, _ := MsgEncode(CmdHttpRespClose, id, from, []byte("Error"))
				conn1.Write(cr)
				timeout_ch <- 0
				return
			}
			n := size1 / 4000
			for i := 0; i < n; i++ {
				p := data[i*4000 : i*4000+4000]
				r, _ := MsgEncode(CmdHttpRespContinued, id, from, p)
				conn1.Write(r)
				timeout_ch <- 1
			}
			p := data[n*4000:]
			r, _ = MsgEncode(CmdHttpRespContinued, id, from, p)
			conn1.Write(r)
			timeout_ch <- 1
		} else {
			var tail bool
			for {
				chunk1 := readChunk(reader1)
				if chunk1 == nil {
					break
				}
				//fmt.Printf("Chunk:%v\r\n", chunk1)
				r, _ := MsgEncode(CmdHttpRespContinued, id, from, chunk1)
				conn1.Write(r)
				timeout_ch <- 1
				if string(chunk1) == "0\r\n" {
					tail = true
				}
				if tail {
					if string(chunk1) == "\r\n" {
						break
					}
				}
			}
		}
		r, _ = MsgEncode(CmdHttpRespClose, id, from, []byte("\r\n"))
		conn1.Write(r)
		timeout_ch <- 1
	}
}
