package main

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/rocket049/go-jsonrpc2glib"
)

var (
	token      []byte
	id         int64
	httpId     int64
	proxyPort  int
	rpcSrv     *PClient
	serverAddr string
	servePort  int = 7890
)

func init() {
	rpcSrv = new(PClient)
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
	go httpProxy(conn1)
	//replace httpServe and rpcSrv.startRpcSrv4Glib
	res1 := make(chan bool, 1)
	go localServe(conn1, res1)
	//go httpServe(conn1)
	u1, _ := user.Current()
	go startMyHttpServe(filepath.Join(u1.HomeDir, "ChatShare"), fmt.Sprintf("localhost:%d", proxyPort))
	go readConn(conn1)
	ok := <-res1
	close(res1)
	if ok {
		startFileServ(conn1)
	} /* else {
		go startFileServ(conn1)
	loop1:
		for i := 0; i < 50; i++ {
			//commonds
			var cmd string
			fmt.Scanf("%s", &cmd)
			switch cmd {
			case "new":
				register(conn1)
			case "login":
				login(conn1)
			case "chat":
				chat(conn1)
			case "proxy":
				httpConnect()
			case "file":
				//sendFile
				sendFile(conn1)
			case "users":
				getUsers(conn1, []int64{5, 6})
			case "friend":
				addFriend(conn1)
			case "friends":
				getFriends(conn1)
			case "proxyport":
				proxyport()
			case "q":
				break loop1
			default:
				fmt.Println("error commond ", cmd)
			}
		}
	} */
	fmt.Println("quit")
}

//func httpConnect() {
//	fmt.Printf("ID:")
//	fmt.Scanf("%d", &httpId)
//}

//func proxyport() {
//	fmt.Printf("Port:")
//	fmt.Scanf("%d", &proxyPort)
//}

type UserInfo struct {
	Id       int64
	Name     string
	Sex      int
	Birthday string
	Desc     string
	Pwdmd5   string
}

//func getFriends(conn1 io.Writer) {
//	req, _ := MsgEncode(CmdGetFriends, 0, 0, []byte("\n"))
//	conn1.Write(req)
//}

//func addFriend(conn1 io.Writer) {
//	var fid int64
//	fmt.Printf("ID:")
//	fmt.Scanf("%d", &fid)
//	req, _ := MsgEncode(CmdAddFriend, 0, fid, []byte("\n"))
//	conn1.Write(req)
//}

//func getUsers(conn1 io.Writer, ids []int64) {
//	b, err := json.Marshal(ids)
//	if err != nil {
//		log.Println(err)
//		return
//	}
//	req, _ := MsgEncode(CmdGetNames, 0, 0, b)
//	log.Println("send:", string(b))
//	conn1.Write(req)
//}

//func sendFile(conn1 io.Writer) {
//	var name string
//	var to int64
//	fmt.Printf("To:")
//	fmt.Scanf("%d", &to)
//	fmt.Printf("File Name:")
//	fmt.Scanf("%s", &name)
//	var sender = new(FileSender)
//	sender.Prepare(name, to, conn1)
//	ok, _ := sender.SendFileHeader()
//	if ok == false {
//		return
//	}
//	sender.SendFileBody()
//	fmt.Println("finish send file.")
//}

//func register(conn1 io.Writer) {
//	var name string
//	var sex int
//	var birthday string
//	var desc string
//	var pwd string
//	fmt.Printf("Name:")
//	fmt.Scanf("%s", &name)
//	fmt.Printf("Sex(1/2):")
//	fmt.Scanf("%d", &sex)
//	fmt.Printf("Birthday:")
//	fmt.Scanf("%s", &birthday)
//	fmt.Printf("Description:")
//	fmt.Scanf("%s", &desc)
//	fmt.Printf("Password:")
//	fmt.Scanf("%s", &pwd)
//	//MD5(name+pwd) 避免重复
//	buf := bytes.NewBufferString(name)
//	buf.WriteString(pwd)
//	var pwdmd5 = md5.Sum(buf.Bytes())
//	e := make([]byte, 32)
//	hex.Encode(e, pwdmd5[:])
//	var user1 = &UserInfo{0, name, sex, birthday, desc, string(e)}
//	p, err := json.Marshal(user1)
//	if err != nil {
//		log.Println(err)
//		return
//	}
//	msg, _ := MsgEncode(CmdRegister, 0, 0, p)
//	conn1.Write(msg)
//}

////type LogDgam struct {
////	Name   string
////	Pwdmd5 []byte
////}

//func login(conn1 io.Writer) {
//	var name string
//	var pwd string
//	fmt.Printf("Name:")
//	fmt.Scanf("%s", &name)
//	fmt.Printf("Password:")
//	fmt.Scanf("%s", &pwd)
//	//MD5(name+pwd) 避免重复
//	buf := bytes.NewBufferString(name)
//	buf.WriteString(pwd)
//	var pwdmd5 = md5.Sum(buf.Bytes())
//	e := make([]byte, 32)
//	hex.Encode(e, pwdmd5[:])
//	bbuf := bytes.NewBufferString("")
//	bbuf.Write(token)
//	bbuf.Write(e)
//	b := md5.Sum(bbuf.Bytes())
//	dgam := &LogDgam{Name: name, Pwdmd5: b[:]}
//	bmsg, _ := json.Marshal(dgam)
//	msg, _ := MsgEncode(CmdLogin, 0, 0, bmsg)
//	conn1.Write(msg)
//}

//func chat(conn1 io.Writer) {
//	var s string
//	fmt.Printf("to ID:")
//	fmt.Scanf("%d", &id)
//	fmt.Printf("Message:")
//	fmt.Scanf("%s", &s)
//	msg, _ := MsgEncode(CmdChat, 0, id, []byte(s))
//	conn1.Write(msg)
//}

func OnReady(msg *MsgType) {
	token = msg.Msg
	rpcSrv.token = msg.Msg
}

//for httpProxy
var httpChan chan MsgType
var httpReqChan chan MsgType

//goroutine
func httpProxy(conn1 io.ReadWriter) {
HEAD1:
	httpReqChan = make(chan MsgType, 1)

	for {
	LOOP2:
		msg, ok := <-httpChan
		if ok == false {
			return
		}
		//send request header
		var httpConn io.ReadWriteCloser
		httpConn, err := net.Dial("tcp", fmt.Sprintf("localhost:%d", proxyPort))
		if err != nil {
			log.Printf("localhost:%d :%v\n", proxyPort, err)
			httpConn = errReader
			//continue
		}
		//httpConn.SetReadDeadline(time.Now().Add(time.Second * 10))
		httpConn.Write(msg.Msg)
		//fmt.Println("sent header")
		//send request body
		for {
			rbody, ok := <-httpReqChan
			if ok == false {
				httpConn.Close()
				goto HEAD1
			}
			if rbody.Cmd == CmdHttpReqClose {
				break
			} else if rbody.Cmd == CmdHttpReqContinued {
				httpConn.Write(rbody.Msg)
			}
		}
		//fmt.Println("sent req body")
		//recv response
		reader1 := bufio.NewReader(httpConn)
		//header
		headbuf := bytes.NewBufferString("")
		var size1 int
		for {
			line1, _, err := reader1.ReadLine()
			if err != nil {
				log.Printf("ReadLine:%v\n", err)
				cr, _ := MsgEncode(CmdHttpRespClose, id, msg.From, []byte("Error"))
				conn1.Write(cr)
				goto LOOP2
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
		r, _ := MsgEncode(CmdHttpRespContinued, id, msg.From, headbuf.Bytes())
		conn1.Write(r)
		//body
		var data []byte
		if size1 > 0 {
			data = make([]byte, size1)
			io.ReadFull(reader1, data)
			n := size1 / 4000
			for i := 0; i < n; i++ {
				p := data[i*4000 : i*4000+4000]
				r, _ := MsgEncode(CmdHttpRespContinued, id, msg.From, p)
				conn1.Write(r)
			}
			p := data[n*4000:]
			r, _ = MsgEncode(CmdHttpRespContinued, id, msg.From, p)
			conn1.Write(r)
		} else {
			var tail bool
			for {
				chunk1 := readChunk(reader1)
				if chunk1 == nil {
					break
				}
				//fmt.Printf("Chunk:%v\r\n", chunk1)
				r, _ := MsgEncode(CmdHttpRespContinued, id, msg.From, chunk1)
				conn1.Write(r)
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
		r, _ = MsgEncode(CmdHttpRespClose, id, msg.From, []byte("\r\n"))
		conn1.Write(r)
		httpConn.Close()
	}
}

func readChunk(reader1 *bufio.Reader) []byte {
	line1, _, err := reader1.ReadLine()
	if err != nil {
		log.Printf("ReadLine:%v\n", err)
		return nil
	}
	if len(line1) == 0 {
		return []byte("\r\n")
	}
	size1, err := strconv.ParseInt(string(line1), 16, 32)
	if err != nil {
		return []byte{}
	}
	//fmt.Printf("chunk size: %d\n", size1)
	if size1 == 0 {
		return []byte("0\r\n")
	}
	data := make([]byte, size1)
	io.ReadFull(reader1, data)
	buf := bytes.NewBuffer(line1)
	buf.WriteString("\r\n")
	buf.Write(data)
	return buf.Bytes()
}

var serveChan chan MsgType = make(chan MsgType, 1)

func httpResponse(conn1 io.ReadWriter, conn2 net.Conn) {
	//conn2.SetReadDeadline(time.Now().Add(time.Second * 5))
	reader1 := bufio.NewReader(conn2)
	reqbuf := bytes.NewBufferString("")
	var size1 int
	//read header
LOOP1:
	for {
		size1 = 0
		for {
			line1, _, err := reader1.ReadLine()
			if err != nil {
				log.Printf("ReadLine:%v\n", err)
				conn2.Close()
				break LOOP1
			}
			s1 := strings.ToLower(string(line1))
			//fmt.Printf("%s\n", s1)
			if strings.HasPrefix(s1, "content-length:") {
				sz1, err := strconv.ParseInt(s1[16:], 10, 32)
				if err == nil {
					size1 = int(sz1)
				}
			}
			reqbuf.Write(line1)
			reqbuf.WriteString("\r\n")
			if len(line1) == 0 {
				log.Println("request complete")
				break
			}
		}
		//read post body
		if size1 > 0 {
			//fmt.Printf("content-length:%d\n", size1)
			body := make([]byte, size1)
			io.ReadFull(reader1, body)
			reqbuf.Write(body)
			//reqbuf.WriteString("\r\n")
		}
		//超过4000就分段
		buf := make([]byte, 4000)
		n, _ := reqbuf.Read(buf)
		req, _ := MsgEncode(CmdHttpRequest, id, httpId, buf[:n])
		conn1.Write(req)
		for {
			n, _ = reqbuf.Read(buf)
			if n == 0 {
				break
			}
			//fmt.Printf("BufferRead:%d\n", n)
			req, _ := MsgEncode(CmdHttpReqContinued, id, httpId, buf[:n])
			conn1.Write(req)
		}
		req, _ = MsgEncode(CmdHttpReqClose, id, httpId, []byte("\r\n"))
		conn1.Write(req)
		//从远端返回
		//fmt.Printf("Wait response\n")
		for {
			res, ok := <-serveChan
			if ok == false {
				conn2.Close()
				serveChan = make(chan MsgType, 1)
				break LOOP1
			}
			if res.Cmd == CmdHttpRespClose {
				if string(res.Msg) == "Offline" {
					conn2.Close()
					break LOOP1
				} else if string(res.Msg) == "Error" {
					conn2.Close()
					break LOOP1
				} else {
					break
				}
			} else if res.Cmd == CmdHttpRespContinued {
				conn2.Write(res.Msg)
			}
		}
		//fmt.Printf("trans finish\n")
	}
}

//goroutine replace httpServe and startRpcSrv4Glib
func localServe(conn1 net.Conn, res1 chan bool) {
	l, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", servePort))
	if err != nil {
		panic(err)
	}
	defer l.Close()
	res1 <- true
	jsonrpc2glib.DebugMode(false)
	serveID := 0
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println("2.accept", err)
			return
		}

		switch serveID {
		case 0:
			rpcSrv.setConn(conn1, conn)
			go rpcSrv.startServe()
			serveID += 1
		case 1:
			go httpResponse(conn1, conn)
		}

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
			rpcSrv.notifyMsg(msg)
		case CmdHttpRequest:
			go pushHttpChan(msg)
		case CmdHttpRespContinued:
			pushServeChan(msg)
		case CmdHttpRespClose:
			pushServeChan(msg)
		case CmdHttpReqContinued:
			httpReqChan <- *msg
		case CmdHttpReqClose:
			httpReqChan <- *msg
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
			rpcSrv.notifyMsg(msg)
		case CmdReturnStrangers:
			cmdChan <- *msg
		case CmdReturnQueryID:
			cmdChan <- *msg
		default:
			log.Printf("Cmd:%d From:%d To:%d Msg:%s\n", msg.Cmd, msg.From, msg.To, string(msg.Msg))
		}
	}
}
func main() {
	var port1 = flag.Int("port", 7890, "local serve port")
	flag.Parse()
	servePort = *port1
	proxyPort = servePort + 2000
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
	exe1, err := os.Executable()
	if err != nil {
		log.Println(err)
		return "."
	}
	dir1 := filepath.Dir(exe1)
	res := filepath.Join(dir1, name1)
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
