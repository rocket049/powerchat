package main

/*
#cgo pkg-config: glib-2.0
#include "glib.h"

struct NUParam {
	char* Name;
	int Sex;
	int Birth;
	char* Desc;
	char* Pwd;
};
struct UserData{
	gint64 Id;
	char* Name;
	int Sex;
	int Age;
	char* Desc;
};

//AppendFriendFn(id,name,sex,age,desc,msg,msg_time)
typedef void (*AppendUserFn)(gint64,char*,int,int,char*,char*,char*);
static void callAppendUser(void *f,gint64 id,char *name,int sex,int age,char *desc,char *msg,char *msg_time){
	AppendUserFn fn = (AppendUserFn)f;
	fn(id,name,sex,age,desc,msg,msg_time);
}
//NotifyFn(Cmd,From,To,Msg)
typedef void (*NotifyFn)(char,gint64,gint64,char*);
static void callNotify(void *f,char Cmd,gint64 From,gint64 To,char* Msg){
	NotifyFn fn = (NotifyFn)f;
	fn(Cmd,From,To,Msg);
}
static gint64 Get(gint64 *v,int i){
	return v[i];
}

static void FillUserData(struct UserData *v,gint64 id,char *name,int sex,int age,char *desc){
	v->Id = id;
	v->Name = name;
	v->Sex = sex;
	v->Age = age;
	v->Desc = desc;
}
*/
import "C"

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
	"unsafe"
)

var (
	cmdChan chan MsgType = make(chan MsgType, 1)
	noticer *Noticer
)

func init() {
	noticer, _ = NewNoticer()
}

type pChatClient struct {
	conn      net.Conn
	token     []byte
	id        int64
	httpId    int64
	proxyPort int
	fileSend  *FileSender
	notifyFn  unsafe.Pointer
}

func (c *pChatClient) setToken(tk []byte) {
	c.token = make([]byte, len(tk))
	copy(c.token, tk)
}

func (c *pChatClient) setConn(connection net.Conn) {
	c.conn = connection
}

func (c *pChatClient) setID(ident int64) {
	c.id = ident
}

//export Client_SetHttpId
func Client_SetHttpId(ident C.gint64) {
	cSrv.httpId = int64(ident)
	httpId = cSrv.httpId
	clearLocRouter()
}

//export Client_NewUser
func Client_NewUser(p *C.struct_NUParam) C.int {
	name1 := strings.TrimSpace(C.GoString(p.Name))
	if len(C.GoString(p.Name)) != len(name1) {
		return -1
	}
	if len(C.GoString(p.Pwd)) == 0 {
		return -1
	}
	var user1 = &UserInfo{0, name1, int(p.Sex),
		fmt.Sprintf("%d-01-01", int(p.Birth)), C.GoString(p.Desc), string(newuserMd5(name1, C.GoString(p.Pwd)))}
	b, err := json.Marshal(user1)
	if err != nil {
		log.Println(err)
		return -1
	}
	msg, _ := MsgEncode(CmdRegister, 0, 0, b)
	cSrv.conn.Write(msg)
	ret := <-cmdChan
	if ret.Cmd == CmdRegResult && string(ret.Msg[0:2]) == "OK" {
		return 1
	} else {
		return -1
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

//NewPasswd params Name,OldMd5-login中的算法,NewMd5-NewUser的算法
//export Client_NewPasswd
func Client_NewPasswd(name, pwdOld, pwdNew *C.char) C.int {
	gName := C.GoString(name)
	gPwdold := C.GoString(pwdOld)
	gPwdnew := C.GoString(pwdNew)
	msg1 := make(map[string][]byte)
	msg1["old"] = loginMd5(gName, gPwdold, cSrv.token)
	msg1["new"] = newuserMd5(gName, gPwdnew)
	msg1["name"] = []byte(gName)

	msg2, err := json.Marshal(msg1)
	if err != nil {
		return 0
	}
	msg, _ := MsgEncode(CmdUpdatePasswd, 0, 0, msg2)
	cSrv.conn.Write(msg)
	res := checkPwd(gName, gPwdnew)
	return C.int(res)
}

//checkPwd 测试密码是否正确，同步方法
func checkPwd(name, pwd string) int {
	dgam := &LogDgam{Name: name, Pwdmd5: loginMd5(name, pwd, cSrv.token)}
	bmsg, _ := json.Marshal(dgam)
	msg, _ := MsgEncode(CmdLogin, 0, 0, bmsg)
	cSrv.conn.Write(msg)
	resp, ok := <-cmdChan
	if ok == false {
		return 3
	}
	s := string(resp.Msg[0:4])
	if strings.HasPrefix(s, "FAIL") {
		return 2
	}
	return 1
}

//export Client_Login
func Client_Login(name, pwd *C.char, p *C.struct_UserData) C.int {
	dgam := &LogDgam{Name: C.GoString(name), Pwdmd5: loginMd5(C.GoString(name), C.GoString(pwd), cSrv.token)}
	bmsg, _ := json.Marshal(dgam)
	msg, _ := MsgEncode(CmdLogin, 0, 0, bmsg)
	cSrv.conn.Write(msg)
	resp, ok := <-cmdChan
	if ok == false {
		return C.int(0)
	}
	if resp.Cmd != CmdLogResult {
		return C.int(0)
	}
	s := string(resp.Msg[0:4])
	if strings.HasPrefix(s, "FAIL") {
		return C.int(0)
	}
	var u UserBaseInfo
	err := json.Unmarshal(resp.Msg, &u)
	if err != nil {
		return C.int(0)
	}
	cSrv.id = u.Id
	id = u.Id

	//load manual config
	cfg1 := readManual(cSrv.id)
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
	C.FillUserData(p, C.gint64(u.Id), C.CString(u.Name), C.int(u.Sex),
		C.int(time.Now().Year()-u.Birthday.Year()), C.CString(u.Desc))
	return C.int(1)
}

type msgOfflineData struct {
	Timestamp string
	Msg       string
}
type FriendData struct {
	Id         int64
	Name       string
	Sex        int
	Age        int
	Desc       string
	MsgOffline msgOfflineData
}
type UserBaseInfo struct {
	Id         int64
	Name       string
	Sex        int
	Birthday   time.Time
	Desc       string
	MsgOffline string
}

//export Client_GetFriends
func Client_GetFriends(callback unsafe.Pointer) {
	req, _ := MsgEncode(CmdGetFriends, cSrv.id, 0, []byte("\n"))
	cSrv.conn.Write(req)
	var resp MsgType
	var ok bool
	for {
		resp, ok = <-cmdChan
		if ok == false {
			return
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
		return
	}

	//ret := []FriendData{}
	for _, v := range frds {
		var offmsg msgOfflineData
		if len(v.MsgOffline) > 17 {
			json.Unmarshal([]byte(v.MsgOffline), &offmsg)
		}
		//ret = append(ret, FriendData{Id: v.Id, Name: v.Name, Sex: v.Sex,
		//Age: time.Now().Year() - v.Birthday.Year(), Desc: v.Desc, MsgOffline: offmsg})
		C.callAppendUser(callback, C.gint64(v.Id), C.CString(v.Name), C.int(v.Sex),
			C.int(time.Now().Year()-v.Birthday.Year()), C.CString(v.Desc),
			C.CString(offmsg.Msg), C.CString(offmsg.Timestamp))
	}
	//log.Println("GetFriends")
	go notifyVersion()
	return
}

//export Client_UserStatus
func Client_UserStatus(uid C.gint64) C.int {
	req, _ := MsgEncode(CmdUserStatus, 0, int64(uid), []byte("\n"))
	cSrv.conn.Write(req)
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

//export Client_QueryID
func Client_QueryID(uid C.gint64, msg *C.char, callback unsafe.Pointer) {
	req, _ := MsgEncode(CmdQueryID, cSrv.id, int64(uid), []byte("\n"))
	cSrv.conn.Write(req)
	var resp MsgType
	var ok bool
	for {
		resp, ok = <-cmdChan
		if ok == false {
			return
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
		return
	}
	C.callAppendUser(callback, C.gint64(v.Id), C.CString(v.Name), C.int(v.Sex),
		C.int(time.Now().Year()-v.Birthday.Year()), C.CString(v.Desc),
		msg, C.CString(time.Now().Format("2006-01-02 15:04:05")))
}

//export Client_MoveStrangerToFriend
func Client_MoveStrangerToFriend(fid C.gint64) {
	req, _ := MsgEncode(CmdMoveStranger, cSrv.id, int64(fid), []byte("\n"))
	cSrv.conn.Write(req)
}

//export Client_GetStrangerMsgs
func Client_GetStrangerMsgs(callback unsafe.Pointer) {
	req, _ := MsgEncode(CmdGetStrangers, cSrv.id, 0, []byte("\n"))
	cSrv.conn.Write(req)
	var resp MsgType
	var ok bool
	for {
		resp, ok = <-cmdChan
		if ok == false {
			return
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
		return
	}

	//ret := []FriendData{}
	for _, v := range frds {
		var offmsg msgOfflineData
		if len(v.MsgOffline) > 17 {
			json.Unmarshal([]byte(v.MsgOffline), &offmsg)
		}
		//ret = append(ret, FriendData{Id: v.Id, Name: v.Name, Sex: v.Sex,
		//Age: time.Now().Year() - v.Birthday.Year(), Desc: v.Desc, MsgOffline: offmsg})
		C.callAppendUser(callback, C.gint64(v.Id), C.CString(v.Name), C.int(v.Sex),
			C.int(time.Now().Year()-v.Birthday.Year()), C.CString(v.Desc),
			C.CString(offmsg.Msg), C.CString(offmsg.Timestamp))
	}
	//*res = ret
	//log.Println("GetStrangerMsgs")
}

//export Client_SearchPersons
func Client_SearchPersons(key *C.char, callback unsafe.Pointer) {
	req, _ := MsgEncode(CmdSearchPersons, cSrv.id, 0, []byte(C.GoString(key)))
	cSrv.conn.Write(req)
	var resp MsgType
	var ok bool
	for {
		resp, ok = <-cmdChan
		if ok == false {
			return
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
		return
	}
	//ret := []FriendData{}
	for _, v := range frds {
		//ret = append(ret, FriendData{Id: v.Id, Name: v.Name, Sex: v.Sex,
		//Age: time.Now().Year() - v.Birthday.Year(), Desc: v.Desc})
		if v.Id == cSrv.id {
			continue
		}
		C.callAppendUser(callback, C.gint64(v.Id), C.CString(v.Name), C.int(v.Sex),
			C.int(time.Now().Year()-v.Birthday.Year()), C.CString(v.Desc),
			C.CString(""), C.CString(""))
	}
	//*res = ret
}

type ChatMessage struct {
	To  int64
	Msg string
}

//export Client_ChatTo
func Client_ChatTo(to C.gint64, msg *C.char) {
	//log.Println(p.Msg)
	p := &ChatMessage{To: int64(to), Msg: C.GoString(msg)}
	buf := bytes.NewBufferString("TEXT")
	buf.WriteString(p.Msg)
	msg1, _ := MsgEncode(CmdChat, 0, p.To, buf.Bytes())
	cSrv.conn.Write(msg1)
}

//export Client_Tell
func Client_Tell(uid C.gint64) {
	msg, _ := MsgEncode(CmdChat, 0, int64(uid), []byte("LOGI"))
	cSrv.conn.Write(msg)
}

//export Client_TellAll
func Client_TellAll(uid *C.gint64, num C.int) {
	uids := make([]int64, int(num))
	var i C.int
	for i = 0; i < num; i++ {
		uids[i] = int64(C.Get(uid, i))
	}
	mMsg := &MultiSendMsg{Ids: uids, Msg: "LOGI"}
	multiSendGo(mMsg)
}

//export Client_MultiSend
func Client_MultiSend(msg *C.char, ids *C.gint64, num C.int) {
	uids := make([]int64, int(num))
	var i C.int
	for i = 0; i < num; i++ {
		uids[i] = int64(C.Get(ids, i))
	}
	mMsg := &MultiSendMsg{Ids: uids, Msg: fmt.Sprintf("TEXT%s", C.GoString(msg))}
	multiSendGo(mMsg)
}

func multiSendGo(param *MultiSendMsg) error {
	bmsg, err := json.Marshal(param)
	if err != nil {
		return err
	}
	msg, _ := MsgEncode(CmdMultiSend, 0, 0, bmsg)
	cSrv.conn.Write(msg)
	return nil
}

//export Client_Ping
func Client_Ping() {
	msg, _ := MsgEncode(CmdPing, 0, 0, []byte("\n"))
	_, err := cSrv.conn.Write(msg)
	if err != nil {
		log.Println(err)
		cSrv.conn.Close()
	}
}

//export Client_HttpConnect
func Client_HttpConnect(uid C.gint64) {
	httpId = int64(uid)
	cSrv.httpId = httpId
}

//export Client_ProxyPort
func Client_ProxyPort(port C.int) C.int {
	var port1 int = servePort + 2000
	if int(port) != 0 {
		port1 = int(port)
	}
	cSrv.proxyPort = port1
	proxyPort = port1
	cfg1 := make(map[string]string)
	cfg1["ProxyPort"] = fmt.Sprintf("%d", port1)
	saveManual(cSrv.id, cfg1)
	return C.int(port1)
}

var notifyFn unsafe.Pointer

//export Client_SetNotifyFn
func Client_SetNotifyFn(callback unsafe.Pointer) {
	notifyFn = callback
	cSrv.notifyFn = callback
}

func notifyMsg(msg *MsgType) {
	if noticer != nil && msg.Cmd == CmdChat {
		go noticer.Play()
	}
	C.callNotify(cSrv.notifyFn, C.char(msg.Cmd), C.gint64(msg.From), C.gint64(msg.To), C.CString(string(msg.Msg)))
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

//export Client_SendFile
func Client_SendFile(to C.gint64, pathName *C.char) C.int {
	param := &SFParam{To: int64(to), PathName: C.GoString(pathName)}
	if cSrv.fileSend != nil {
		if cSrv.fileSend.Status() {
			notifyMsg(&MsgType{Cmd: CmdChat, From: 0, To: 0, Msg: []byte("OneOnly!\n")})
			return 1
		}
	}
	//log.Println(param.PathName);
	var sender = new(FileSender)
	sender.Prepare(param.PathName, param.To, cSrv.conn)
	ret := sender.SendFileHeader()
	if ret > 0 {
		return C.int(ret)
	}
	cSrv.fileSend = sender
	go sender.SendFileBody()
	return 0
}

//export Client_AddFriend
func Client_AddFriend(fid C.gint64) {
	if int64(fid) == cSrv.id {
		return
	}
	req, _ := MsgEncode(CmdAddFriend, 0, int64(fid), []byte("\n"))
	cSrv.conn.Write(req)
}

//export Client_RemoveFriend
func Client_RemoveFriend(uid C.gint64) {
	req, _ := MsgEncode(CmdRemoveFriend, cSrv.id, int64(uid), []byte("\n"))
	cSrv.conn.Write(req)
}

//export Client_GetProxyPort
func Client_GetProxyPort() C.int {
	return C.int(proxyPort)
}

//export Client_Quit
func Client_Quit() {
	log.Println("rpc serve normal exit")
	noticer.Close()
	cSrv.conn.Close()
}

//export Client_GetHost
func Client_GetHost(p **C.char) {
	*p = C.CString(serverAddr)
}

//export Client_OpenPath
func Client_OpenPath(param *C.char) {
	//log.Println("Open:", C.GoString(param))
	myOpen(C.GoString(param))
}

//export Client_UpdateDesc
func Client_UpdateDesc(param *C.char) {
	req, _ := MsgEncode(CmdUpdateDesc, 0, 0, []byte(C.GoString(param)))
	cSrv.conn.Write(req)
}

type CheckDelData struct {
	Md5   []byte
	Token []byte
}

//rpc service block
//DeleteMe param=[]string{name,pwd}
//export Client_DeleteMe
func Client_DeleteMe(name, pwd *C.char) C.int {
	var token [8]byte
	io.ReadFull(rand.Reader, token[:])
	md5v := loginMd5(C.GoString(name), C.GoString(pwd), token[:])
	var checkd = &CheckDelData{Md5: md5v, Token: token[:]}
	jsond, err := json.Marshal(checkd)
	if err != nil {
		return 0
	}
	req, _ := MsgEncode(CmdDeleteMe, 0, 0, jsond)
	_, err = cSrv.conn.Write(req)
	if err == nil {
		return 1
	} else {
		return 0
	}
}

//export Client_GetPgPath
func Client_GetPgPath(p **C.char) {
	filepath1, _ := os.Executable()
	*p = C.CString(fmt.Sprintf("%s", filepath.Dir(filepath1)))
}

//export Client_MakeLauncher
func Client_MakeLauncher() {
	makeLauncher()
}

//main
func main() {}
