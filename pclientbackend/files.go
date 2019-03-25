package pclientbackend

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"io"
	"log"
	"math/rand"
	"mime"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type fileReceiver struct {
	From      int64
	Header    *fileHeaderType
	File      *os.File
	Lock      sync.Mutex
	timeStamp time.Time
}

func (s *fileReceiver) UpdateTime() {
	s.timeStamp = time.Now()
}

func (s *fileReceiver) IsRunning() bool {
	if s.From == 0 {
		return false
	} else {
		if s.timeStamp.Add(time.Second * 60).After(time.Now()) {
			return true
		} else {
			name1 := s.File.Name()
			s.File.Close()
			os.Remove(name1)
			return false
		}
	}
}

var (
	receiver *fileReceiver
)

func init() {
	rand.Seed(time.Now().Unix())
	receiver = &fileReceiver{From: 0, Header: nil, Lock: sync.Mutex{}}
}

func getFileDir() string {
	return getRelatePath("RecvFiles")
}

type fileHeaderType struct {
	Name    string
	Mime    string
	Size    int64
	Session uint32
}

//gorutine
func pushFileMsg2(conn1 io.Writer, msg *MsgType) {
	receiver.Lock.Lock()
	defer receiver.Lock.Unlock()
	if receiver.IsRunning() == false && msg.Cmd == CmdFileHeader {
		receiver.From = msg.From
		receiver.Header = new(fileHeaderType)
		err := json.Unmarshal(msg.Msg, receiver.Header)
		if err != nil {
			receiver.From = 0
			receiver.Header = nil
			return
		}
		fileDir := getFileDir()
		receiver.File, err = os.Create(filepath.Join(fileDir, receiver.Header.Name))
		if err != nil {
			log.Println("error create file:", receiver.Header.Name)
			receiver.From = 0
			receiver.Header = nil
			return
		}
		receiver.UpdateTime()
		notifyMsg(&MsgType{Cmd: CmdChat, From: receiver.From, To: 0,
			Msg: []byte("TEXTSending:" + receiver.Header.Name)})
		//request Accept
		bs := make([]byte, 4)
		binary.BigEndian.PutUint32(bs, receiver.Header.Session)
		resp, _ := MsgEncode(CmdFileAccept, 0, receiver.From, bs)
		conn1.Write(resp)
		log.Printf("Accept %s\n", receiver.Header.Name)
	} else if receiver.From == msg.From && msg.Cmd == CmdFileContinued && len(msg.Msg) >= 4 {
		if dataSession(msg.Msg) != receiver.Header.Session {
			return
		}
		receiver.File.Write(msg.Msg[4:])
		receiver.UpdateTime()
	} else if receiver.From == msg.From && msg.Cmd == CmdFileClose && len(msg.Msg) >= 4 {
		if dataSession(msg.Msg) != receiver.Header.Session {
			return
		}
		name1 := receiver.File.Name()
		receiver.File.Close()
		notifyFile(name1, receiver.From)
		//From=0 释放 receiver
		//log.Println("complete:", receiver.Header.Name)
		receiver.From = 0
	} else if receiver.From == msg.From && msg.Cmd == CmdFileCancel && len(msg.Msg) >= 4 {
		if dataSession(msg.Msg) != receiver.Header.Session {
			return
		}
		name1 := receiver.File.Name()
		receiver.File.Close()
		os.Remove(name1)
		notifyMsg(&MsgType{Cmd: CmdChat, From: receiver.From, To: 0,
			Msg: []byte("TEXTCancel:" + receiver.Header.Name)})
		receiver.From = 0
	} else if receiver.IsRunning() == true && msg.Cmd == CmdFileHeader {
		h1 := new(fileHeaderType)
		err := json.Unmarshal(msg.Msg, h1)
		if err != nil {
			return
		}
		var bs [4]byte
		binary.BigEndian.PutUint32(bs[:], h1.Session)
		resp, _ := MsgEncode(CmdFileBlock, 0, msg.From, bs[:])
		conn1.Write(resp)
	}
}

func dataSession(data []byte) uint32 {
	return binary.BigEndian.Uint32(data[:4])
}

func notifyFile(pathname string, from int64) {
	fh1, err := os.Stat(pathname)
	if err != nil {
		log.Println(err)
		return
	}
	secs := strings.Split(pathname, ".")
	var typ1 string
	if len(secs) > 1 {
		typ1 = secs[len(secs)-1]
		log.Println("Type:", typ1)
	}
	//wd1, _ := os.Getwd()
	//path1 := filepath.Join(wd1, pathname)
	fh2 := &fileHeaderType{Name: pathname,
		Mime:    mime.TypeByExtension("." + typ1),
		Size:    int64(fh1.Size()),
		Session: 0}
	b1, err := json.Marshal(fh2)
	if err != nil {
		log.Println(err)
		return
	}
	msgbuf := bytes.NewBufferString("JSON")
	msgbuf.Write(b1)
	msg := MsgType{Cmd: CmdChat, From: from, To: 0, Msg: msgbuf.Bytes()}
	//msgChan <- msg
	notifyMsg(&msg)
}

type FileSender struct {
	mutex1   sync.Mutex
	session  uint32
	running  bool
	pathname string
	conn     io.Writer
	to       int64
	size     int64
	sendSize int64
}

func (s *FileSender) prepare(pathname string, to int64, conn1 io.Writer) {
	for {
		s.session = rand.Uint32()
		if s.session != 0 {
			break
		}
	}

	s.pathname = pathname
	s.conn = conn1
	s.to = to
	s.running = true
}

var fileResp = make(chan MsgType, 1)

func (s *FileSender) sendFileHeader() (bool, uint32) {
	fh1, err := os.Stat(s.pathname)
	if err != nil {
		log.Println(err)
		return false, 0
	}
	secs := strings.Split(s.pathname, ".")
	var typ1 string
	if len(secs) > 1 {
		typ1 = secs[len(secs)-1]
	}
	fh2 := &fileHeaderType{Name: filepath.Base(s.pathname),
		Mime:    mime.TypeByExtension("." + typ1),
		Size:    int64(fh1.Size()),
		Session: s.session}
	b1, err := json.Marshal(fh2)
	if err != nil {
		log.Println(err)
		return false, 0
	}
	msg, _ := MsgEncode(CmdFileHeader, 0, s.to, b1)
	s.conn.Write(msg)
	log.Printf("Send header: %s\n", fh2.Name)
	res, ok := <-fileResp
	if ok == false {
		log.Println("internal error")
		return false, 0
	}
	if binary.BigEndian.Uint32(res.Msg) != s.session {
		return false, 0
	}
	s.size = fh1.Size()
	return true, s.session
}

func (s *FileSender) cancelTrans() {
	s.mutex1.Lock()
	s.running = false
	s.mutex1.Unlock()
	b1 := make([]byte, 4)
	binary.BigEndian.PutUint32(b1, s.session)
	msg, _ := MsgEncode(CmdFileCancel, 0, s.to, b1)
	s.conn.Write(msg)
}

//goroutine
func (s *FileSender) sendFileBody() {
	f1, err := os.Open(s.pathname)
	if err != nil {
		s.cancelTrans()
		notifyMsg(&MsgType{Cmd: CmdChat, From: 0, To: 0, Msg: []byte("Error:  " + err.Error())})
		return
	}
	defer f1.Close()
	b1 := make([]byte, 4000)
	binary.BigEndian.PutUint32(b1[:4], s.session)
	var running bool = true
	for running {
		s.mutex1.Lock()
		running = s.running
		s.mutex1.Unlock()
		n, _ := f1.Read(b1[4:])
		if n <= 0 {
			break
		}
		s.sendSize += int64(n)
		msg, _ := MsgEncode(CmdFileContinued, 0, s.to, b1[:n+4])
		s.conn.Write(msg)
	}
	msg, _ := MsgEncode(CmdFileClose, 0, s.to, b1[:4])
	s.conn.Write(msg)
	s.mutex1.Lock()
	s.running = false
	s.mutex1.Unlock()
}

func (s *FileSender) getSent() (sent, size int64) {
	return s.sendSize, s.size
}

func (s *FileSender) status() bool {
	s.mutex1.Lock()
	res := s.running
	s.mutex1.Unlock()
	return res
}
