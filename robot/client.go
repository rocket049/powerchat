package main

import (
	"crypto/tls"
	"encoding/json"
	"flag"
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
	proxyPort  int
	robotSrv   *PClient
	serverAddr string
	servePort  int = 8890
	user1      string
	pwd1       string
)

func init() {
	robotSrv = new(PClient)
}

func client() {
	var cfg tls.Config
	//	roots := x509.NewCertPool()
	//	pem, _ := ioutil.ReadFile("pems/a-cert.pem")
	//	roots.AppendCertsFromPEM(pem)
	//	cfg.RootCAs = roots
	cfg.InsecureSkipVerify = true
	conn1, err := tls.Dial("tcp", serverAddr, &cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer conn1.Close()
	robotSrv.setConn(conn1)
	go httpProxy2(conn1)
	//replace httpServe and rpcSrv.startRpcSrv4Glib
	res1 := make(chan bool, 1)
	go localServe(conn1, res1)
	//go httpServe(conn1)
	go readConn(conn1)
	ok := <-res1
	close(res1)
	if ok {
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
	robotSrv.token = msg.Msg
	go robotSrv.Login(LogParam{Name: user1, Pwd: pwd1})
	go robotSrv.Ping()
}

//for httpProxy
var httpChan chan MsgType
var httpReqChan chan MsgType

var serveChan chan MsgType = make(chan MsgType, 1)

//goroutine replace httpServe and startRpcSrv4Glib
func localServe(conn1 net.Conn, res1 chan bool) {
	go startMyHttpServe(getRelatePath("ChatShare"), fmt.Sprintf("localhost:%d", proxyPort))
	res1 <- true
	httpRespRouter()
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
			robotSrv.notifyMsg(msg)
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
			continue
		case CmdFileBlock:
			//cancel send file
			continue
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
			robotSrv.notifyMsg(msg)
		case CmdReturnStrangers:
			cmdChan <- *msg
		case CmdReturnQueryID:
			cmdChan <- *msg
		case CmdUserStatus:
			cmdChan <- *msg
		default:
			//log.Printf("Cmd:%d From:%d To:%d Msg:%s\n", msg.Cmd, msg.From, msg.To, string(msg.Msg))
		}
	}
}
func main() {
	var u = flag.String("u", "", "username")
	var p = flag.String("p", "", "password")
	flag.Parse()
	user1 = *u
	pwd1 = *p
	servePort = 8890
	proxyPort = 10890
	filepath1, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	path1 := filepath.Dir(filepath1)
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
	client()
}

func getRelatePath(name1 string) string {
	filepath1, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	path1 := filepath.Dir(filepath1)
	res := filepath.Join(path1, "data", name1)
	os.MkdirAll(res, os.ModePerm)
	return res
}
