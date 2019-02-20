package main

import (
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type ClientData struct {
	Conn    io.WriteCloser
	Counter uint8
}

var clients = &sync.Map{}

func server() {
	filepath1, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	path1 := filepath.Dir(filepath1)
	cert, err := tls.LoadX509KeyPair(filepath.Join(path1, "pems/a-cert.pem"),
		filepath.Join(path1, "pems/a-key.pem"))
	if err != nil {
		log.Fatal(err)
	}
	cfg := &tls.Config{Certificates: []tls.Certificate{cert}}

	var cfg1 = make(map[string]string)
	cfgFile, err := ioutil.ReadFile(filepath.Join(path1, "config.json"))
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(cfgFile, &cfg1)
	if err != nil {
		log.Fatal(err)
	}
	host, ok := cfg1["Host"]
	if ok == false {
		log.Fatal("config file parse error\n")
	}
	listener, err := tls.Listen("tcp", host, cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()
	for {
		conn1, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go serveConn(conn1)
	}
}

//goroutine
func serveConn(conn1 net.Conn) {
	defer conn1.Close()
	var client = new(ClientType)
	client.ServeReady(conn1)
	client.Auth = false

	for {
		conn1.SetReadDeadline(time.Now().Add(time.Minute * 5))
		msg, err := ReadMsg(conn1)
		if err != nil {
			log.Printf("Disconnect:%v\n", client.Id)
			break
		}
		smsg := MsgDecode(msg)
		if smsg == nil {
			log.Printf("MsgDecode error\n")
			break
		}
		switch smsg.Cmd {
		case CmdPing:
			err := client.Pong(conn1)
			if err != nil {
				log.Printf("on CmdPing:%v\n", err)
			}
		case CmdRegister:
			err := client.RegisterUser(conn1, smsg)
			if err != nil {
				log.Printf("on CmdRegister:%v\n", err)
			}
		case CmdLogin:
			err := client.Login(conn1, smsg)
			if err != nil {
				log.Printf("on CmdLogin:%v\n", err)
			} else {
				defer clients.Delete(client.Id)
			}
		case CmdGetNames:
			//[id1,id2,...]
			err := client.GetNames(conn1, smsg)
			if err != nil {
				log.Printf("on CmdGetNames:%v\n", err)
			}
		case CmdAddFriend:
			err := addFriend(client.Id, smsg.To)
			if err != nil {
				log.Printf("on CmdAddFriend:%v\n", err)
			}
		case CmdGetFriends:
			//get friends
			err := client.GetFriends(conn1)
			if err != nil {
				log.Printf("on CmdGetFriends:%v\n", err)
			}
		case CmdGetStrangers:
			err := client.GetStrangers(conn1)
			if err != nil {
				log.Printf("on CmdGetStrangers:%v\n", err)
			}
		case CmdMoveStranger:
			//move stranger to friend
			err := moveStrangerToFriend(smsg.From, smsg.To)
			if err != nil {
				log.Printf("on CmdMoveStranger:%v\n", err)
			}
		case CmdRemoveFriend:
			err := removeFriend(smsg.From, smsg.To)
			if err != nil {
				log.Printf("on CmdRemoveFriend:%v\n", err)
			}

		case CmdSearchPersons:
			err := client.SearchPersons(conn1, string(smsg.Msg))
			if err != nil {
				log.Printf("on CmdSearchPersons:%v\n", err)
			}
		case CmdHttpRequest:
			if client.IsOnline(smsg.To) {
				client.Redirect(conn1, smsg)
			} else {
				req := MsgType{Cmd: CmdHttpRespClose, From: 0, To: client.Id, Msg: []byte("Offline")}
				client.Redirect(conn1, &req)
			}
		case CmdChat:
			//log.Println("Chat:",string(smsg.Msg))
			if client.IsOnline(smsg.To) == false {
				offlineMsg(client.Id, smsg.To, string(smsg.Msg))
			}
			//print("Redirect\n")
			client.Redirect(conn1, smsg)
		case CmdQueryID:
			//query user info by id
			err := client.QueryUserById(conn1, smsg)
			if err != nil {
				log.Printf("on CmdQueryID:%v\n", err)
			}
		case CmdFileHeader:
			//阻止陌生人的文件传输
			if isFriend(client.Id, smsg.To) == false {
				client.SysResp(conn1, CmdFileBlock, "Must Be Friend")
			} else {
				client.Redirect(conn1, smsg)
			}
		case CmdUpdatePasswd:
			err := client.UpdatePass(smsg)
			if err != nil {
				log.Printf("on CmdUpdatePasswd:%v\n", err)
			}
		case CmdUpdateDesc:
			err := client.UpdateDesc(smsg)
			if err != nil {
				log.Printf("on CmdUpdateDesc:%v\n", err)
			}
		case CmdDeleteMe:
			client.DeleteUser(conn1, smsg)

		case CmdUserStatus:
			//return status
			client.UserStatus(conn1, smsg)
		default:
			err := client.Redirect(conn1, smsg)
			if err != nil {
				log.Printf("on Code: %d :%v\n", smsg.Cmd, err)
			}
		}
	}
	if client.Id > 0 {
		clients.Delete(client.Id)
	}
}

type ClientType struct {
	Id    int64
	Token [32]byte
	Auth  bool
}

func (c *ClientType) UserStatus(conn1 io.Writer, msg *MsgType) error {
	if c.IsOnline(msg.To) {
		//online
		c.SysResp(conn1, CmdUserStatus, "Y")
	} else {
		//offline
		c.SysResp(conn1, CmdUserStatus, "N")
	}
	return nil
}

func (c *ClientType) QueryUserById(conn1 io.Writer, msg *MsgType) error {
	list1, err := getUsersByIds([]int64{msg.To})
	if err != nil {
		c.SysResp(conn1, CmdReturnQueryID, "")
		return err
	}
	var res *UserBaseInfo = nil
	for _, v := range list1 {
		res = v
		break
	}
	b, err := json.Marshal(res)
	if err != nil {
		c.SysResp(conn1, CmdReturnQueryID, "")
		return err
	}
	resp, _ := MsgEncode(CmdReturnQueryID, 0, 0, b)
	_, err = conn1.Write(resp)
	return nil
}

func (c *ClientType) IsOnline(uid int64) bool {
	_, ok := clients.Load(uid)
	return ok
}

func (c *ClientType) SearchPersons(conn1 io.Writer, key string) error {
	list1, err := searchUsers(key)
	if err != nil {
		c.SysResp(conn1, CmdReturnPersons, "")
		return err
	}
	b, err := json.Marshal(list1)
	if err != nil {
		c.SysResp(conn1, CmdReturnPersons, "")
		return err
	}
	resp, _ := MsgEncode(CmdReturnPersons, 0, 0, b)
	_, err = conn1.Write(resp)
	return err
}

func (c *ClientType) GetStrangers(conn1 io.Writer) error {
	list1, err := getStrangers(c.Id)
	if err != nil {
		c.SysResp(conn1, CmdReturnStrangers, "")
		return err
	}
	b, err := json.Marshal(list1)
	if err != nil {
		c.SysResp(conn1, CmdReturnStrangers, "")
		return err
	}
	resp, _ := MsgEncode(CmdReturnStrangers, 0, 0, b)
	_, err = conn1.Write(resp)
	return err
}

func (c *ClientType) GetFriends(conn1 io.Writer) error {
	list1, err := getFriends(c.Id)
	if err != nil {
		c.SysResp(conn1, CmdRetFriends, "")
		return err
	}
	b, err := json.Marshal(list1)
	if err != nil {
		c.SysResp(conn1, CmdRetFriends, "")
		return err
	}
	resp, _ := MsgEncode(CmdRetFriends, 0, 0, b)
	_, err = conn1.Write(resp)
	return err
}

func (c *ClientType) GetNames(conn1 io.Writer, msg *MsgType) error {
	var ids []int64
	err := json.Unmarshal(msg.Msg, &ids)
	if err != nil {
		c.SysResp(conn1, CmdRetNames, "")
		return err
	}
	names, err := getUsersByIds(ids)
	if err != nil {
		c.SysResp(conn1, CmdRetNames, "")
		return err
	}
	b, err := json.Marshal(names)
	if err != nil {
		c.SysResp(conn1, CmdRetNames, "")
		return err
	}
	rmsg, _ := MsgEncode(CmdRetNames, 0, c.Id, b)
	_, err = conn1.Write(rmsg)
	return err
}
func (c *ClientType) Redirect(conn1 io.Writer, msg *MsgType) error {
	if c.Auth == false {
		c.SysResp(conn1, CmdSysReturn, "Permission Denied")
		return errors.New("not auth")
	}
	v, ok := clients.Load(msg.To)
	if ok == false {
		c.SysResp(conn1, CmdSysReturn, fmt.Sprintf("Offline %d", msg.To))
		return errors.New("Not logined")
	}
	req, err := MsgEncode(msg.Cmd, c.Id, msg.To, msg.Msg)
	if err != nil {
		log.Printf("MsgEncode:%v\n", err)
		c.SysResp(conn1, CmdSysReturn, err.Error())
		return err
	}
	_, err = io.Copy(v.(*ClientData).Conn, bytes.NewBuffer(req))
	if err != nil {
		log.Printf("io.Copy:%v\n", err)
		c.SysResp(conn1, CmdSysReturn, err.Error())
		return err
	}
	return nil
}

func (c *ClientType) ServeReady(conn1 io.Writer) error {
	var p = c.Token[:]
	io.ReadFull(rand.Reader, p)
	first, _ := MsgEncode(CmdReady, 0, 0, p)
	_, err := io.Copy(conn1, bytes.NewBuffer(first))
	return err
}

func (c *ClientType) Pong(conn1 io.Writer) error {
	v, ok := clients.Load(c.Id)
	if ok == false {
		link1 := v.(*ClientData).Conn
		link1.Close()
		clients.Delete(c.Id)
		return nil
	}
	v.(*ClientData).Counter++
	msg, _ := MsgEncode(CmdPong, 0, c.Id, []byte("\n"))
	_, err := io.Copy(conn1, bytes.NewBuffer(msg))
	return err
}

type UserRegType struct {
	Name     string
	Sex      int
	Birthday string
	Desc     string
	Pwdmd5   string
}

func (c *ClientType) RegisterUser(conn1 io.Writer, msg *MsgType) error {
	var u UserRegType
	err := json.Unmarshal(msg.Msg, &u)
	if err != nil {
		c.SysResp(conn1, CmdSysReturn, err.Error())
		return err
	}
	id, err := insertUser(u.Name, u.Sex, u.Birthday, u.Desc, u.Pwdmd5)
	var rmsg string
	if err != nil {
		rmsg = fmt.Sprintf("ERR:%s", err.Error())
	} else {
		rmsg = fmt.Sprintf("OK:%d", id)
	}
	req, _ := MsgEncode(CmdRegResult, 0, 0, []byte(rmsg))
	_, err = conn1.Write(req)
	return err
}

type LogDgam struct {
	Name   string
	Pwdmd5 []byte
}

func (c *ClientType) Login(conn1 io.WriteCloser, msg *MsgType) error {
	u1 := new(LogDgam)
	err := json.Unmarshal(msg.Msg, u1)
	if err != nil {
		c.SysResp(conn1, CmdLogResult, err.Error())
		return err
	}
	u, err := getUserByName(u1.Name)
	if err != nil {
		c.SysResp(conn1, CmdLogResult, "FAIL NOT EXIST")
		return err
	}
	buf := bytes.NewBufferString("")
	buf.Write(c.Token[:])
	buf.Write([]byte(u.Pwdmd5))
	b := md5.Sum(buf.Bytes())
	if bytes.Compare(b[:], u1.Pwdmd5) != 0 {
		c.SysResp(conn1, CmdLogResult, "FAIL AUTH")
		return errors.New("auth fail")
	}
	c.Id = u.Id
	c.Auth = true
	clients.Store(c.Id, &ClientData{Conn: conn1, Counter: 0})
	u.Pwdmd5 = ""
	reqb, _ := json.Marshal(u)
	c.SysResp(conn1, CmdLogResult, string(reqb))
	return nil
}

func (c *ClientType) UpdatePass(msg *MsgType) error {
	info1 := make(map[string][]byte)
	err := json.Unmarshal(msg.Msg, &info1)
	if err != nil {
		return err
	}
	name1 := string(info1["name"])
	old1 := info1["old"]
	new1 := info1["new"]

	u, err := getUserByName(name1)
	if err != nil {
		return err
	}
	buf := bytes.NewBufferString("")
	buf.Write(c.Token[:])
	buf.Write([]byte(u.Pwdmd5))
	b := md5.Sum(buf.Bytes())
	if bytes.Compare(b[:], old1) != 0 {
		return errors.New("auth fail")
	}

	return updatePasswd(c.Id, string(new1))
}

func (c *ClientType) UpdateDesc(msg *MsgType) error {
	desc := string(msg.Msg)
	return updateDesc(c.Id, desc)
}

type CheckDelData struct {
	Md5   []byte
	Token []byte
}

//DeleteUser smsg.Msg={md5:md5_v,token:token_v}
func (c *ClientType) DeleteUser(conn1 io.Writer, msg *MsgType) error {

	var checkd CheckDelData
	err := json.Unmarshal(msg.Msg, &checkd)
	if err != nil {
		c.SysResp(conn1, CmdSysReturn, "DELETE 0\n")
		return err
	}
	u, err := getUserById(c.Id)
	if err != nil {
		c.SysResp(conn1, CmdSysReturn, "DELETE 0\n")
		return err
	}
	//check pwd
	buf := bytes.NewBufferString("")
	buf.Write(checkd.Token)
	buf.WriteString(u.Pwdmd5)
	b := md5.Sum(buf.Bytes())
	if bytes.Compare(b[:], checkd.Md5) != 0 {
		c.SysResp(conn1, CmdSysReturn, "DELETE 0\n")
		return errors.New("auth fail")
	}
	err = deleteUser(c.Id)
	if err != nil {
		c.SysResp(conn1, CmdSysReturn, "DELETE 0\n")
		return err
	}

	c.SysResp(conn1, CmdSysReturn, "DELETE 1\n")
	v1, ok := clients.Load(c.Id)
	if ok {
		link1 := v1.(*ClientData).Conn
		link1.Close()
		clients.Delete(c.Id)
	}
	return nil
}

func (c *ClientType) SysResp(conn1 io.Writer, cmd ChatCommand, s string) {
	r, _ := MsgEncode(cmd, 0, 0, []byte(s))
	_, err := conn1.Write(r)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	server()
}
