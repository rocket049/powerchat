package main

import (
	"encoding/binary"
	"encoding/json"
	"io"
	"log"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

var (
	fileChan chan MsgType
)

type fileHeaderType struct {
	Name    string
	Mime    string
	Size    int64
	Session uint32
}

//gorutine
func pushFileMsg(msg *MsgType) {
	fileChan <- *msg
}

//goroutine
func startFileServ(conn1 io.Writer) {
	fileChan = make(chan MsgType, 10)
	for {
		h1, ok := <-fileChan
		if ok == false {
			log.Println("close file Header Chan")
			return
		}
		if h1.Cmd != CmdFileHeader {
			//log.Println("error not a fileHeader")
			continue
		}
		h2 := new(fileHeaderType)
		err := json.Unmarshal(h1.Msg, h2)
		if err != nil {
			//log.Println("error parse fileHeader")
			continue
		}

		//block file transfer
		var from = h1.From
		var b2 [4]byte
		binary.BigEndian.PutUint32(b2[:], h2.Session)
		msg, _ := MsgEncode(CmdFileBlock, 0, from, b2[:])
		conn1.Write(msg)
	}
}
