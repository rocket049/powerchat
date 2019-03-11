package pclientbackend

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path/filepath"
)

var (
	token      []byte
	id         int64
	httpId     int64
	proxyPort  int
	cSrv       *ChatClient
	serverAddr string
	servePort  int = 7890
	pgPath     string
	dataHome   string
)

func init() {
	cSrv = nil
}

func client(ctl1 chan int) {
	var cfg tls.Config
	cfg.InsecureSkipVerify = true
	conn1, err := tls.Dial("tcp", serverAddr, &cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer conn1.Close()
	go httpProxy2(conn1)
	//replace httpServe and cSrv.startcSrv4Glib
	res1 := make(chan bool, 1)
	go localServe(conn1, res1)
	//go httpServe(conn1)
	go readConn(conn1)
	ok := <-res1
	close(res1)
	if ok {
		ctl1 <- 1
		startFileServ(conn1)
	}
	fmt.Println("quit")
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
}

//for httpProxy
var httpChan chan MsgType
var httpReqChan chan MsgType

var serveChan chan MsgType = make(chan MsgType, 1)

//goroutine replace httpServe and startcSrv4Glib
func localServe(conn1 net.Conn, res1 chan bool) {
	//only on connection
	l, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", servePort))
	if err != nil {
		panic(err)
	}
	defer l.Close()
	cSrv.setConn(conn1)
	go startMyHttpServe(getRelatePath("ChatShare"), fmt.Sprintf("localhost:%d", proxyPort))
	res1 <- true
	go httpRespRouter()
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println("2.accept", err)
			return
		}
		go httpResponse2(conn1, conn, httpId)

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
		msgb, err := ReadMsg(conn1)
		if err != nil {
			log.Printf("ReadMsg:%v\n", err)
			return
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
			pushFileMsg(msg)
		case CmdFileContinued:
			//file
			pushFileMsg(msg)
		case CmdFileClose:
			//file
			pushFileMsg(msg)
		case CmdFileCancel:
			//file
			pushFileMsg(msg)
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
func main_init(dir string) {
	dataHome = dir

	servePort = 7890
	proxyPort = servePort + 2000

	filepath1, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	path1 := filepath.Dir(filepath1)
	pgPath = path1

	var cfg1 = make(map[string]string)
	cfgFile, err := ioutil.ReadFile(filepath.Join(path1, "config.json"))
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(cfgFile, &cfg1)
	if err != nil {
		log.Fatal(err)
	}
	var ok bool
	serverAddr, ok = cfg1["Host"]
	if ok == false {
		log.Fatal("config file parse error\n")
	}
	httpChan = make(chan MsgType, 10)
	ctl1 := make(chan int, 1)
	go client(ctl1)
	<-ctl1
	close(ctl1)
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
