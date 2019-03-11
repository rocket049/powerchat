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

func init() {
	rand.Seed(time.Now().Unix())
}

var (
	fileChan chan MsgType
	fileResp chan MsgType
)

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
func pushFileMsg(msg *MsgType) {
	fileChan <- *msg
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
	msg := MsgType{Cmd: CmdFileHeader, From: from, To: 0, Msg: msgbuf.Bytes()}
	//msgChan <- msg
	notifyMsg(&msg)
}

//goroutine
func startFileServ(conn1 io.Writer) {
	fileResp = make(chan MsgType, 1)
	fileChan = make(chan MsgType, 10)
	var fileDir = getFileDir()
	for {
		h1, ok := <-fileChan
		if ok == false {
			log.Println("close file Header Chan")
			return
		}
		if h1.Cmd != CmdFileHeader {
			log.Println("error not a fileHeader")
			continue
		}
		h2 := new(fileHeaderType)
		err := json.Unmarshal(h1.Msg, h2)
		if err != nil {
			log.Println("error parse fileHeader")
			continue
		}
		var from = h1.From
		var session = h2.Session
		var bs = make([]byte, 4)
		file1, err := os.Create(filepath.Join(fileDir, h2.Name))
		if err != nil {
			log.Println("error create file:", h2.Name)
			continue
		}
		//request Accept
		binary.BigEndian.PutUint32(bs, session)
		msg, _ := MsgEncode(CmdFileAccept, 0, h1.From, bs)
		conn1.Write(msg)
		log.Printf("Accept %s\n", h2.Name)
		//receive file
		for {
			b1, ok := <-fileChan
			if ok == false {
				log.Println("close file chan")
				file1.Close()
				return
			}
			if b1.Cmd == CmdFileContinued {
				if b1.From != from {
					continue
				}
				if dataSession(b1.Msg) != session {
					continue
				}
				file1.Write(b1.Msg[4:])
			} else if b1.Cmd == CmdFileClose {
				//finish
				if dataSession(b1.Msg) != session {
					continue
				}
				file1.Close()
				//recv notify
				//msgChan <- MsgType{Cmd: h1.Cmd, From: h1.From, To: 0, Msg: []byte("Recv 1 file")}
				notifyFile(file1.Name(), h1.From)
				break
			} else if b1.Cmd == CmdFileCancel {
				//cancel
				if dataSession(b1.Msg) != session {
					continue
				}
				name1 := file1.Name()
				file1.Close()
				os.Remove(name1)
				break
			} else if b1.Cmd == CmdFileHeader {
				var tmph fileHeaderType
				var b2 [4]byte
				err := json.Unmarshal(h1.Msg, &tmph)
				if err != nil {
					log.Println("error fileHeader")
					continue
				}
				binary.BigEndian.PutUint32(b2[:], tmph.Session)
				msg, _ := MsgEncode(CmdFileBlock, 0, b1.From, b2[:])
				conn1.Write(msg)
			}
		}
	}
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

func (s *FileSender) Prepare(pathname string, to int64, conn1 io.Writer) {
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

func (s *FileSender) SendFileHeader() (bool, uint32) {
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

func (s *FileSender) CancelTrans() {
	s.mutex1.Lock()
	s.running = false
	s.mutex1.Unlock()
	b1 := make([]byte, 4)
	binary.BigEndian.PutUint32(b1, s.session)
	msg, _ := MsgEncode(CmdFileCancel, 0, s.to, b1)
	s.conn.Write(msg)
}

//goroutine
func (s *FileSender) SendFileBody() {
	f1, err := os.Open(s.pathname)
	if err != nil {
		s.CancelTrans()
		log.Panicln(err)
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

func (s *FileSender) GetSent() (sent, size int64) {
	return s.sendSize, s.size
}

func (s *FileSender) Status() bool {
	s.mutex1.Lock()
	res := s.running
	s.mutex1.Unlock()
	return res
}
