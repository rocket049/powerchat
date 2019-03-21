package pclientbackend

//symbols for mobile device anroid/ios

import (
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var (
	cmdChan chan MsgType = make(chan MsgType, 1)
)

type msgOfflineData struct {
	Timestamp string
	Msg       string
}
type UserDataRet struct {
	Id        int64
	Name      string
	Sex       int
	Age       int
	Desc      string
	Timestamp string
	Msg       string
}
type UserDataArray struct {
	Users []UserDataRet
	pos   int
}
type IdArray struct {
	Ids []int64
}
type UserBaseInfo struct {
	Id         int64
	Name       string
	Sex        int
	Birthday   time.Time
	Desc       string
	MsgOffline string
}

type ChatClient struct {
	conn      net.Conn
	token     []byte
	id        int64
	httpId    int64
	proxyPort int
	fileSend  *FileSender
}

func (p *UserDataArray) Next() (res *UserDataRet) {
	if p.pos >= len(p.Users) {
		res = nil
	} else {
		res = &p.Users[p.pos]
	}
	p.pos++
	return res
}

//NewIdArray 用于初始化TellAll和MultiSend的ids/uids参数
//export NewIdArray
func NewIdArray() *IdArray {
	return new(IdArray)
}

//Append 向对象内部增加id
func (p *IdArray) Append(id int64) {
	p.Ids = append(p.Ids, id)
}

func (c *ChatClient) setToken(tk []byte) {
	c.token = make([]byte, len(tk))
	copy(c.token, tk)
}

func (c *ChatClient) setConn(connection net.Conn) {
	c.conn = connection
	go c.startPing()
}

func (c *ChatClient) setID(ident int64) {
	c.id = ident
}

//GetChatClient 初始化，参数：数据目录路径
func GetChatClient(dataDir, cfgSrc string) *ChatClient {
	if cSrv == nil {
		//only initial once
		cSrv = new(ChatClient)
		main_init(dataDir, cfgSrc)
	}
	return cSrv
}

//SetServeId 设置要访问的联系人 id
func (c *ChatClient) SetServeId(ident int64) {
	c.httpId = ident
	httpId = ident
	clearLocRouter()
}

//NewUser 注册新用户，参数： 名字，密码，性别(1-男，2-女)，出生年份（四位数：1985,2005,...），自述信息
func (c *ChatClient) NewUser(name, pwd string, sex, birth int, desc string) bool {
	name1 := strings.TrimSpace(name)
	if len(name) != len(name1) {
		return false
	}
	if len(pwd) == 0 {
		return false
	}
	var user1 = &UserInfo{0, name1, sex,
		fmt.Sprintf("%d-01-01", birth), desc, string(newuserMd5(name1, pwd))}
	b, err := json.Marshal(user1)
	if err != nil {
		log.Println(err)
		return false
	}
	msg, _ := MsgEncode(CmdRegister, 0, 0, b)
	c.conn.Write(msg)
	ret := <-cmdChan
	if ret.Cmd == CmdRegResult && string(ret.Msg[0:2]) == "OK" {
		return true
	} else {
		return false
	}
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

func newuserMd5(name, pwd string) []byte {
	pwdmd5 := md5.Sum([]byte(name + pwd))
	e := make([]byte, 32)
	hex.Encode(e, pwdmd5[:])
	return e
}

//NewPasswd 参数： Name,OldMd5-login中的算法,NewMd5-NewUser的算法
func (c *ChatClient) NewPasswd(name, pwdOld, pwdNew string) bool {
	msg1 := make(map[string][]byte)
	msg1["old"] = loginMd5(name, pwdOld, c.token)
	msg1["new"] = newuserMd5(name, pwdNew)
	msg1["name"] = []byte(name)

	msg2, err := json.Marshal(msg1)
	if err != nil {
		return false
	}
	msg, _ := MsgEncode(CmdUpdatePasswd, 0, 0, msg2)
	c.conn.Write(msg)
	return true
}

//Login 参数：name ,password
//阻塞函数，最好在线程中运行或者用异步函数包装
func (c *ChatClient) Login(name, pwd string) *UserDataRet {
	dgam := &LogDgam{Name: name, Pwdmd5: loginMd5(name, pwd, c.token)}
	bmsg, _ := json.Marshal(dgam)
	msg, _ := MsgEncode(CmdLogin, 0, 0, bmsg)
	c.conn.Write(msg)
	resp, ok := <-cmdChan
	if ok == false {
		return nil
	}
	if resp.Cmd != CmdLogResult {
		return nil
	}
	s := string(resp.Msg[0:4])
	if strings.HasPrefix(s, "FAIL") {
		return nil
	}
	var u UserBaseInfo
	err := json.Unmarshal(resp.Msg, &u)
	if err != nil {
		return nil
	}
	c.id = u.Id
	id = u.Id

	//load manual config
	cfg1 := readManual(c.id)
	if cfg1 != nil {
		sPort, ok := cfg1["ProxyPort"]
		if ok {
			//log.Println(reflect.TypeOf(sPort).Name())
			port, err := strconv.ParseInt(sPort, 10, 32)
			if err == nil {
				proxyPort = int(port)
			} else {
				log.Println(err)
			}
			log.Println("proxy port:", proxyPort, ",manual:", port)
		}
	}

	return &UserDataRet{Id: u.Id, Name: u.Name, Sex: u.Sex,
		Age: time.Now().Year() - u.Birthday.Year(), Desc: u.Desc, Timestamp: "", Msg: ""}
}

//GetFriends 返回联系人列表和离线信息,阻塞函数，最好在线程中运行或者用异步函数包装
func (c *ChatClient) GetFriends() *UserDataArray {
	req, _ := MsgEncode(CmdGetFriends, c.id, 0, []byte("\n"))
	c.conn.Write(req)
	var resp MsgType
	var ok bool
	for {
		resp, ok = <-cmdChan
		if ok == false {
			return nil
		}
		if resp.Cmd == CmdRetFriends {
			break
		} else {
			cmdChan <- resp
			time.Sleep(time.Millisecond * 100)
		}
	}
	//json
	frds := make(map[int64]UserBaseInfo)
	err := json.Unmarshal(resp.Msg, &frds)
	if err != nil {
		return nil
	}

	ret := []UserDataRet{}
	for _, v := range frds {
		var offmsg msgOfflineData
		if len(v.MsgOffline) > 17 {
			json.Unmarshal([]byte(v.MsgOffline), &offmsg)
		}
		ret = append(ret, UserDataRet{Id: v.Id, Name: v.Name, Sex: v.Sex,
			Age: time.Now().Year() - v.Birthday.Year(), Desc: v.Desc, Timestamp: offmsg.Timestamp, Msg: offmsg.Msg})
	}
	//log.Println("GetFriends")
	go notifyVersion()
	return &UserDataArray{Users: ret, pos: 0}
}

//UserStatus 参数：id int64，返回值：0-offline，1-online
//阻塞函数，最好在线程中运行或者用异步函数包装
func (c *ChatClient) UserStatus(uid int64) int {
	req, _ := MsgEncode(CmdUserStatus, 0, uid, []byte("\n"))
	c.conn.Write(req)
	var resp MsgType
	var ok bool
	for {
		resp, ok = <-cmdChan
		if ok == false {
			return 0
		}
		if resp.Cmd == CmdUserStatus {
			break
		} else {
			cmdChan <- resp
			time.Sleep(time.Millisecond * 100)
		}
	}

	if string(resp.Msg) == "Y" {
		return 1
	} else {
		return 0
	}
}

//QueryID 查询陌生人信息，返回 UserDataRet 结构体，参数：ID int64， msg stirng。 msg - 接受到的陌生人信息
//阻塞函数，最好在线程中运行或者用异步函数包装
func (c *ChatClient) QueryID(uid int64, msg string) *UserDataRet {
	req, _ := MsgEncode(CmdQueryID, c.id, uid, []byte("\n"))
	c.conn.Write(req)
	var resp MsgType
	var ok bool
	for {
		resp, ok = <-cmdChan
		if ok == false {
			return nil
		}
		if resp.Cmd == CmdReturnQueryID {
			break
		} else {
			cmdChan <- resp
			time.Sleep(time.Millisecond * 100)
		}
	}
	var v UserBaseInfo
	err := json.Unmarshal(resp.Msg, &v)
	if err != nil {
		return nil
	}
	return &UserDataRet{Id: v.Id, Name: v.Name, Sex: v.Sex, Age: time.Now().Year() - v.Birthday.Year(),
		Desc: v.Desc, Timestamp: time.Now().Format("2006-01-02 15:04:05"), Msg: msg}
}

//MoveStrangerToFriend 把留言的陌生人加入联系人名单
func (c *ChatClient) MoveStrangerToFriend(fid int64) {
	req, _ := MsgEncode(CmdMoveStranger, c.id, fid, []byte("\n"))
	c.conn.Write(req)
}

//GetStrangerMsgs 读取全部陌生人的留言，阻塞函数，最好在线程中运行或者用异步函数包装
func (c *ChatClient) GetStrangerMsgs() *UserDataArray {
	req, _ := MsgEncode(CmdGetStrangers, c.id, 0, []byte("\n"))
	c.conn.Write(req)
	var resp MsgType
	var ok bool
	for {
		resp, ok = <-cmdChan
		if ok == false {
			return nil
		}
		if resp.Cmd == CmdReturnStrangers {
			break
		} else {
			cmdChan <- resp
			time.Sleep(time.Millisecond * 100)
		}
	}
	//json
	frds := make(map[int64]UserBaseInfo)
	err := json.Unmarshal(resp.Msg, &frds)
	if err != nil {
		return nil
	}

	ret := []UserDataRet{}
	for _, v := range frds {
		var offmsg msgOfflineData
		if len(v.MsgOffline) > 17 {
			json.Unmarshal([]byte(v.MsgOffline), &offmsg)
		}
		ret = append(ret, UserDataRet{Id: v.Id, Name: v.Name, Sex: v.Sex,
			Age: time.Now().Year() - v.Birthday.Year(), Desc: v.Desc, Timestamp: offmsg.Timestamp, Msg: offmsg.Msg})
	}
	return &UserDataArray{Users: ret, pos: 0}
}

//SearchPersons 搜索用户，阻塞函数，最好在线程中运行或者用异步函数包装
func (c *ChatClient) SearchPersons(key string) *UserDataArray {
	req, _ := MsgEncode(CmdSearchPersons, c.id, 0, []byte(key))
	c.conn.Write(req)
	var resp MsgType
	var ok bool
	for {
		resp, ok = <-cmdChan
		if ok == false {
			return nil
		}
		if resp.Cmd == CmdReturnPersons {
			break
		} else {
			cmdChan <- resp
			time.Sleep(time.Millisecond * 100)
		}
	}
	//json
	frds := make(map[int64]UserBaseInfo)
	err := json.Unmarshal(resp.Msg, &frds)
	if err != nil {
		return nil
	}
	ret := []UserDataRet{}
	for _, v := range frds {
		ret = append(ret, UserDataRet{Id: v.Id, Name: v.Name, Sex: v.Sex,
			Age: time.Now().Year() - v.Birthday.Year(), Desc: v.Desc, Timestamp: "", Msg: ""})
	}
	return &UserDataArray{Users: ret, pos: 0}
}

type ChatMessage struct {
	To  int64
	Msg string
}

//ChatTo 发送文字聊天信息
func (c *ChatClient) ChatTo(to int64, msg string) {
	//log.Println(p.Msg)
	p := &ChatMessage{To: to, Msg: msg}
	buf := bytes.NewBufferString("TEXT")
	buf.WriteString(p.Msg)
	msg1, _ := MsgEncode(CmdChat, 0, p.To, buf.Bytes())
	c.conn.Write(msg1)
}

//Tell 上线通知
func (c *ChatClient) Tell(uid int64) {
	msg, _ := MsgEncode(CmdChat, 0, uid, []byte("LOGI"))
	c.conn.Write(msg)
}

//TellAll 向id列表中的所有人发上线通知
func (c *ChatClient) TellAll(uids *IdArray) {
	mMsg := &MultiSendMsg{Ids: uids.Ids, Msg: "LOGI"}
	c.multiSendGo(mMsg)
}

//MultiSend 向id列表中的所有人发送同一个文字信息
func (c *ChatClient) MultiSend(msg string, ids *IdArray) {
	mMsg := &MultiSendMsg{Ids: ids.Ids, Msg: fmt.Sprintf("TEXT%s", msg)}
	c.multiSendGo(mMsg)
}

func (c *ChatClient) multiSendGo(param *MultiSendMsg) error {
	bmsg, err := json.Marshal(param)
	if err != nil {
		return err
	}
	msg, _ := MsgEncode(CmdMultiSend, 0, 0, bmsg)
	c.conn.Write(msg)
	return nil
}

//startPing 心跳数据
func (c *ChatClient) startPing() {
	msg, _ := MsgEncode(CmdPing, 0, 0, []byte("\n"))
	for {
		time.Sleep(time.Second * 60)
		_, err := c.conn.Write(msg)
		if err != nil {
			c.conn.Close()
			break
		}
	}
}

//SetProxyPort 设置代理端口
func (c *ChatClient) SetProxyPort(port int) int {
	var port1 int = servePort + 2000
	if port > 0 {
		port1 = port
	}
	c.proxyPort = port1
	proxyPort = port1
	cfg1 := make(map[string]string)
	cfg1["ProxyPort"] = fmt.Sprintf("%d", port1)
	saveManual(c.id, cfg1)
	return port1
}

var notifyChan = make(chan MsgType, 1)

//GetMsg 读取信息，阻塞函数，需要在线程中运行或者用异步函数包装
func (c *ChatClient) GetMsg() *MsgType {
	msg, ok := <-notifyChan
	if ok {
		return &msg
	} else {
		return nil
	}
}

//notifyMsg
func notifyMsg(msg *MsgType) {
	notifyChan <- *msg
}

//notifyVersion:rpc notify the client, go routine
type Version struct {
	Version uint `json:"version"`
}

func notifyVersion() {
	r, err := http.Get("https://gitee.com/rocket049/powerchat/raw/master/release.json")
	if err != nil {
		log.Println(err)
		return
	}
	defer r.Body.Close()
	buf := make([]byte, 50)
	n, err := r.Body.Read(buf)
	if n == 0 {
		log.Println(err)
		return
	}
	//fmt.Println(string(buf[:n]))
	var v1 Version
	json.Unmarshal(buf[:n], &v1)
	res := &MsgType{Cmd: CmdSysReturn,
		From: 0,
		To:   0,
		Msg:  []byte(fmt.Sprintf("Version:%v", v1.Version))}
	notifyMsg(res)
}

type SFParam struct {
	To       int64
	PathName string
}

//SendFile 发送文件，参数：id,pathname
func (c *ChatClient) SendFile(to int64, pathName string) {
	param := &SFParam{To: to, PathName: pathName}
	if c.fileSend != nil {
		if c.fileSend.status() {
			return
		}
	}
	//log.Println(param.PathName);
	var sender = new(FileSender)
	sender.prepare(param.PathName, param.To, c.conn)
	ok, _ := sender.sendFileHeader()
	if ok == false {
		return
	}
	c.fileSend = sender
	go sender.sendFileBody()
}

//AddFriend 加入联系人
func (c *ChatClient) AddFriend(fid int64) {
	req, _ := MsgEncode(CmdAddFriend, 0, fid, []byte("\n"))
	c.conn.Write(req)
}

//RemoveFriend 删除联系人
func (c *ChatClient) RemoveFriend(uid int64) {
	req, _ := MsgEncode(CmdRemoveFriend, c.id, uid, []byte("\n"))
	c.conn.Write(req)
}

//GetProxyPort 返回代理端口
func (c *ChatClient) GetProxyPort() int {
	return proxyPort
}

//Quit 退出
func (c *ChatClient) Quit() {
	c.conn.Close()
	close(fileChan)
	close(notifyChan)
}

//GetHost 返回服务器 IP:PORT
func (c *ChatClient) GetHost() string {
	return serverAddr
}

//GetUrl 返回访问BaseURL，不会变化
func (c *ChatClient) GetUrl() string {
	return fmt.Sprintf("http://localhost:%d", servePort)
}

//UpdateDesc 更新自述信息
func (c *ChatClient) UpdateDesc(param string) {
	req, _ := MsgEncode(CmdUpdateDesc, 0, 0, []byte(param))
	c.conn.Write(req)
}

type CheckDelData struct {
	Md5   []byte
	Token []byte
}

//DeleteMe 删除本用户，参数：用户名和密码
func (c *ChatClient) DeleteMe(name, pwd string) bool {
	var token [8]byte
	io.ReadFull(rand.Reader, token[:])
	md5v := loginMd5(name, pwd, token[:])
	var checkd = &CheckDelData{Md5: md5v, Token: token[:]}
	jsond, err := json.Marshal(checkd)
	if err != nil {
		return false
	}
	req, _ := MsgEncode(CmdDeleteMe, 0, 0, jsond)
	_, err = c.conn.Write(req)
	if err == nil {
		return true
	} else {
		return false
	}
}

//GetPgPath 返回程序所在路径
func (c *ChatClient) GetPgPath() string {
	filepath1, _ := os.Executable()
	return filepath.Dir(filepath1)
}

//main
func main() {}
