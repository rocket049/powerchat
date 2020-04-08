package pclientbackend

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/user"
	"path/filepath"
	"sync"
	"time"
)

var (
	token      []byte
	id         int64
	httpId     int64
	proxyPort  int
	cSrv       *pChatClient
	serverAddr string
	servePort  int = 7890
	pgPath     string
	dataHome   string
	retry      bool = false
	nameRetry  string
	pwdRetry   string
	lServe     *localServer
)

func init() {
	cSrv = new(pChatClient)
}

var (
	//for httpProxy
	httpChan  chan MsgType
	serveChan chan MsgType
	//for export
	cmdChan    chan MsgType
	notifyChan chan MsgType
)
var onceClient sync.Once

func client() {
	var cfg tls.Config
	cfg.InsecureSkipVerify = true
	conn1, err := tls.Dial("tcp", serverAddr, &cfg)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn1.Close()

	httpChan = make(chan MsgType, 10)

	onceClient.Do(func() {
		serveChan = make(chan MsgType, 1)
		cmdChan = make(chan MsgType, 1)
		notifyChan = make(chan MsgType, 1)
	})

	go httpProxy2(conn1)

	var ok bool = true
	if !retry {
		lServe = newLocalServer(conn1)
		res1 := make(chan bool, 1)
		go lServe.Serve(res1)
		ok = <-res1
		close(res1)
	} else {
		lServe.SetConn(conn1)
	}
	if ok {
		readConn(conn1)
	}

	close(httpChan)
}

type UserInfo struct {
	Id       int64
	Name     string
	Sex      int
	Birthday string
	Desc     string
	Pwdmd5   string
}

func OnReady(msg *MsgType) {
	token = msg.Msg
	cSrv.token = msg.Msg
	if retry {
		go retry_login()
	}
}

//goroutine replace httpServe and startcSrv4Glib
type localServer struct {
	conn net.Conn
}

func newLocalServer(c net.Conn) *localServer {
	cSrv.setConn(c)
	return &localServer{conn: c}
}

func (s *localServer) SetConn(c net.Conn) {
	s.conn = c
	cSrv.setConn(c)
}

func (s *localServer) Serve(res1 chan bool) {
	var l net.Listener
	var err error
	for i := 0; i < 8; i++ {
		l, err = net.Listen("tcp", fmt.Sprintf("localhost:%d", servePort))
		if err == nil {
			break
		} else {
			servePort++
			proxyPort = servePort + 2000
		}
	}
	if err != nil {
		panic(err)
	}
	defer l.Close()

	u1, _ := user.Current()
	go startMyHttpServe(filepath.Join(u1.HomeDir, "ChatShare"), fmt.Sprintf("localhost:%d", proxyPort))
	res1 <- true
	go httpRespRouter()
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println("2.accept", err)
			return
		}
		go httpResponse2(s.conn, conn, httpId)

	}
}

func pushHttpChan(msg *MsgType) {
	httpChan <- *msg
}
func pushServeChan(msg *MsgType) {
	serveChan <- *msg
	//fmt.Printf("%v\n", p)
}

//goroutine
func readConn(conn1 net.Conn) {
	for {
		conn1.SetReadDeadline(time.Now().Add(time.Minute * 3))
		msgb, err := ReadMsg(conn1)
		if err != nil {
			log.Printf("ReadMsg:%v\n", err)
			notifyMsg(&MsgType{Cmd: CmdSysReturn, From: 0, To: 0, Msg: []byte("ConnDown\n")})
			break
		}
		msg := MsgDecode(msgb)
		switch msg.Cmd {
		case CmdReady:
			OnReady(msg)
		case CmdChat:
			notifyMsg(msg)
		case CmdHttpRequest:
			pushHttpChan(msg)
		case CmdHttpReqContinued:
			pushHttpChan(msg)
		case CmdHttpReqClose:
			pushHttpChan(msg)

		case CmdHttpRespContinued:
			pushServeChan(msg)
		case CmdHttpRespClose:
			pushServeChan(msg)

		case CmdFileHeader:
			pushFileMsg2(conn1, msg)
		case CmdFileContinued:
			//file
			pushFileMsg2(conn1, msg)
		case CmdFileClose:
			//file
			pushFileMsg2(conn1, msg)
		case CmdFileCancel:
			//file
			pushFileMsg2(conn1, msg)
		case CmdFileAccept:
			//begin send file
			fileResp <- *msg
		case CmdFileBlock:
			//cancel send file
			fileResp <- *msg
		case CmdFileStop:
			//stop send file by session
		case CmdLogResult:
			cmdChan <- *msg
		case CmdRegResult:
			cmdChan <- *msg
		case CmdRetFriends:
			cmdChan <- *msg
		case CmdReturnPersons:
			cmdChan <- *msg
		case CmdSysReturn:
			notifyMsg(msg)
		case CmdReturnStrangers:
			cmdChan <- *msg
		case CmdReturnQueryID:
			cmdChan <- *msg
		case CmdUserStatus:
			cmdChan <- *msg
		default:
			log.Printf("Cmd:%d From:%d To:%d Msg:%s\n", msg.Cmd, msg.From, msg.To, string(msg.Msg))
		}
	}
}

func main_init(dir, cfgSrc string) {
	dataHome = dir

	servePort = 7880
	proxyPort = servePort + 2000

	var cfg1 = make(map[string]string)
	err := json.Unmarshal([]byte(cfgSrc), &cfg1)
	if err != nil {
		log.Fatal(err)
	}
	var ok bool
	serverAddr, ok = cfg1["Host"]
	if ok == false {
		log.Fatal("config file parse error\n")
	}

	go func() {
		for {
			client()
			log.Println("disconnect ,re-connect 10s later.")
			time.Sleep(time.Second * 10)
		}

	}()
}

func getRelatePath(name1 string) string {
	res := filepath.Join(dataHome, name1)
	os.MkdirAll(res, os.ModePerm)
	return res
}

func readManual(uid int64) map[string]string {
	dir1 := getRelatePath("manual")
	pathname := filepath.Join(dir1, fmt.Sprintf("%d", uid))
	p, err := ioutil.ReadFile(pathname)
	if err != nil {
		log.Println("readManual error:", err)
		return nil
	}
	res := make(map[string]string)
	err = json.Unmarshal(p, &res)
	if err != nil {
		log.Println("readManual error:", err)
		return nil
	}
	return res
}

func saveManual(uid int64, data map[string]string) error {
	dir1 := getRelatePath("manual")
	pathname := filepath.Join(dir1, fmt.Sprintf("%d", uid))
	p, err := json.Marshal(data)
	if err != nil {
		return err
	}
	fp, err := os.Create(pathname)
	if err != nil {
		return err
	}
	defer fp.Close()
	_, err = fp.Write(p)
	return err
}
