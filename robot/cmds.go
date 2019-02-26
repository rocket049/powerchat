package main

import (
	"bytes"
	"encoding/binary"
	"io"
)

type ChatCommand uint8

const (
	CmdChat              ChatCommand = iota //客户端之间的文字信息
	CmdHttpRequest                          //http请求
	CmdHttpReqContinued                     //http post body标志
	CmdHttpReqClose                         //http request 结束标志
	CmdHttpRespContinued                    //不完整数据
	CmdHttpRespClose                        //无权限，关闭Http连接
	CmdFileHeader                           //包括图片
	CmdFileAccept                           //接受传送请求
	CmdFileBlock                            //有其他文件传输过程，请求被阻止
	CmdFileContinued                        //不完整数据
	CmdFileClose                            //最后一块数据
	CmdFileStop                             //接收方中途取消
	CmdFileCancel                           //取消文件传输或下载
	CmdLogin                                //登录
	CmdLogResult                            //登录结果
	CmdRegister                             //注册
	CmdRegResult                            //注册结果
	CmdGetNames                             //通过ＩＤ取得名字
	CmdRetNames                             //返回名字／ＩＤ列表
	CmdAddFriend                            //加好友
	CmdGetFriends                           //取得联系人列表
	CmdRetFriends                           //返回名字／ＩＤ列表
	CmdPing                                 //客户端心跳包
	CmdPong                                 //服务器返回心跳包
	CmdReady                                //服务器第一个信号
	CmdSearchPersons                        //搜索人名
	CmdReturnPersons                        //返回搜索结果
	CmdSysReturn                            //服务器返回信息
	CmdGetStrangers                         //请求陌生人留言
	CmdReturnStrangers                      //返回陌生人留言
	CmdMoveStranger                         //把留言的陌生人加入好友
	CmdRemoveFriend                         //删除好友
	CmdQueryID                              //查询基本身份信息
	CmdReturnQueryID                        //返回身份信息
	CmdUpdatePasswd                         //更新密码
	CmdUpdateDesc                           //更新自述文字
	CmdDeleteMe                             //删除当前登录用户
	CmdUserStatus                           //查询、返回用户状态
)

const HeadSize uint16 = 19
const MsgLimit int = 65000 - 19

func MsgEncode(cmd ChatCommand, from, to int64, msg []byte) ([]byte, error) {
	var header = make([]byte, HeadSize)
	var body []byte
	if len(msg) > MsgLimit {
		body = msg[0:MsgLimit]
	} else {
		body = msg
	}
	size1 := HeadSize + uint16(len(body)) - 2
	binary.BigEndian.PutUint16(header, size1)
	header[2] = byte(cmd)
	binary.BigEndian.PutUint64(header[3:11], uint64(from))
	binary.BigEndian.PutUint64(header[11:19], uint64(to))
	buf := bytes.NewBuffer(header[:])
	_, err := buf.Write(body)
	return buf.Bytes(), err
}

type MsgType struct {
	Cmd  ChatCommand
	From int64
	To   int64
	Msg  []byte
}

func MsgDecode(msg []byte) *MsgType {
	res := new(MsgType)
	res.Cmd = ChatCommand(msg[0])
	res.From = int64(binary.BigEndian.Uint64(msg[1:9]))
	res.To = int64(binary.BigEndian.Uint64(msg[9:17]))
	res.Msg = msg[17:]
	return res
}

func ReadMsg(r io.Reader) ([]byte, error) {
	var lbuf = make([]byte, 2)
	_, err := io.ReadFull(r, lbuf)
	if err != nil {
		return nil, err
	}
	var size1 = binary.BigEndian.Uint16(lbuf)
	var data = make([]byte, int(size1))
	_, err = io.ReadFull(r, data)
	if err != nil {
		return nil, err
	} else {
		return data, nil
	}
}
