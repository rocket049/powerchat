package main

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	cmdChan chan MsgType = make(chan MsgType, 1)
)

type PClient struct {
	conn      net.Conn
	token     []byte
	id        int64
	httpId    int64
	proxyPort int
}

func (c *PClient) setToken(tk []byte) {
	c.token = make([]byte, len(tk))
	copy(c.token, tk)
}

func (c *PClient) setConn(connection net.Conn) {
	c.conn = connection
}

func (c *PClient) setID(ident int64) {
	c.id = ident
}

type LogParam struct {
	Name string
	Pwd  string
}
type LogDgam struct {
	Name   string
	Pwdmd5 []byte
}

func loginMd5(name, pwd string, token []byte) []byte {
	buf := bytes.NewBufferString(name)
	buf.WriteString(pwd)
	var pwdmd5 = md5.Sum(buf.Bytes())
	e := make([]byte, 32)
	hex.Encode(e, pwdmd5[:])
	bbuf := bytes.NewBufferString("")
	bbuf.Write(token)
	bbuf.Write(e)
	b := md5.Sum(bbuf.Bytes())
	return b[:]
}

//service
func (c *PClient) Login(p LogParam) {
	dgam := &LogDgam{Name: p.Name, Pwdmd5: loginMd5(p.Name, p.Pwd, c.token)}
	bmsg, _ := json.Marshal(dgam)
	msg, _ := MsgEncode(CmdLogin, 0, 0, bmsg)
	c.conn.Write(msg)
	resp, ok := <-cmdChan
	if ok == false {
		panic("internal error 1")
	}
	if resp.Cmd != CmdLogResult {
		panic("internal error 2")
	}
	s := string(resp.Msg[0:4])
	if strings.HasPrefix(s, "FAIL") {
		panic("login fail 3")
	}

	//log.Println("logined")
}

type ChatMessage struct {
	To  int64
	Msg string
}

//rpc service
func (c *PClient) ChatTo(p ChatMessage) error {
	//log.Println(p.Msg)
	buf := bytes.NewBufferString("TEXT")
	buf.WriteString(p.Msg)
	msg, _ := MsgEncode(CmdChat, 0, p.To, buf.Bytes())
	c.conn.Write(msg)
	return nil
}

//Ping goroutine service
func (c *PClient) Ping() error {
	msg, _ := MsgEncode(CmdPing, 0, 0, []byte("\n"))
	var timer1 = time.NewTimer(time.Second * 50)
	for {
		//t1 := <-timer1.C
		<-timer1.C
		_, err := c.conn.Write(msg)
		//log.Printf("ping %v\n", t1.Second())
		if err != nil {
			return err
		}
		timer1.Reset(time.Second * 50)
	}
}

type MsgReturn struct {
	T    uint8
	From int64
	To   int64
	Msg  string
}

//rpc notify the client
func (c *PClient) notifyMsg(msg *MsgType) error {
	//to deal msg
	//log.Println("cmd:", msg.Cmd)
	if msg.Cmd == CmdChat {
		c.replyChat(msg)
		//log.Println("chat:", msg.From, string(msg.Msg))
	}
	return nil
}

func (c *PClient) replyChat(msg *MsgType) {
	if len(msg.Msg) < 4 {
		return
	}
	pre := string(msg.Msg[:4])
	switch pre {
	case "LOGI":
		c.ChatTo(ChatMessage{To: msg.From, Msg: "Welcome!\n这是一个服务号，可以回答你的问题。"})
		c.ChatTo(ChatMessage{To: msg.From, Msg: getPage(0)})
	case "TEXT":
		if c.showTopic(msg) {
			return
		}
		if c.showIndex(msg) {
			return
		}
		fallback_msg := fmt.Sprintf("不明指令：%s\n请输入'p0'显示首页！", string(msg.Msg[4:]))
		c.ChatTo(ChatMessage{To: msg.From, Msg: fallback_msg})
	}
}

func (c *PClient) showTopic(msg *MsgType) bool {
	ok, err := regexp.Match("^\\d+$", msg.Msg[4:])
	if err != nil {
		//log.Println(err)
		return false
	}
	if ok {
		path1 := filepath.Join(getRelatePath("topics"), string(msg.Msg[4:]))
		text1, err1 := ioutil.ReadFile(path1)
		var res = "No such topic."
		if err1 == nil {
			res = string(text1)
		}
		c.ChatTo(ChatMessage{To: msg.From, Msg: res})
	}
	return ok
}

func (c *PClient) showIndex(msg *MsgType) bool {
	ok, err := regexp.Match("^p\\d+$", msg.Msg[4:])
	if err != nil {
		//log.Println(err)
		return false
	}
	if ok {
		n, err := strconv.ParseInt(string(msg.Msg[5:]), 10, 32)
		if err != nil {
			return false
		}
		c.ChatTo(ChatMessage{To: msg.From, Msg: getPage(int(n))})
	}
	return ok
}
