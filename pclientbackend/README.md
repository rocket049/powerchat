## 本目录为移动平台提供的动态库

可以用 gomobile 编译成安卓SDK库

调用说明：

```
//包名：pclientbackend

//UserDataRet 用户信息数据结构
type UserDataRet struct {
	Id        int64
	Name      string
	Sex       int
	Age       int
	Desc      string
	Timestamp string
	Msg       string
}
//UserDataArray 包装[]UserDataRet的对象，后面有用法说明
type UserDataArray struct {
	Users []UserDataRet
}
//IdArray 包装[]int64的对象，后面有用法说明
type IdArray struct {
	Ids []int64
	pos int
}

//MsgType 通用信息数据包
type MsgType struct {
	Cmd  ChatCommand
	From int64
	To   int64
	Msg  []byte
}

//GetChatClient 初始化，参数：数据目录路径,这个函数必须第一个调用，获得已经初始化好了了对象指针
//func GetChatClient(dataDir,cfgSrc string) *ChatClient
client = GetChatClient( 数据存储目录path1, config.json的内容 )

//NewUser 注册新用户，参数： 名字，密码，性别(1-男，2-女)，出生年份（四位数：1985,2005,...），自述信息。返回：bool
//func (c *ChatClient) NewUser(name, pwd string, sex, birth int, desc string) bool
res = client.NewUser(name, pwd, sex, birth, desc)

//Login 参数：name ,password，登录并返回用户信息
//阻塞函数，最好在线程中运行或者用异步函数包装
//func (c *ChatClient) Login(name, pwd string) *UserDataRet

//GetFriends 返回联系人列表和离线信息，登录成功后应立即调用,阻塞函数，最好在线程中运行或者用异步函数包装
//func (c *ChatClient) GetFriends() *UserDataArray
friends = client.GetFriends()

//GetStrangerMsgs 读取全部陌生人的留言，登录成功后应立即调用，阻塞函数，最好在线程中运行或者用异步函数包装
//func (c *ChatClient) GetStrangerMsgs() *UserDataArray
strangers = client.GetStrangerMsgs()

//SearchPersons 搜索名字包含关键字 key 的用户，阻塞函数，最好在线程中运行或者用异步函数包装
func (c *ChatClient) SearchPersons(key string) *UserDataArray

//NewPasswd 修改密码，参数： Name,OldMd5,NewMd5，返回：bool
//func (c *ChatClient) NewPasswd(name, pwdOld, pwdNew string) bool
res = client.NewPasswd(name,pwdOld,pwdNew)

//CheckPwd 验证密码是否正确，返回值：1 正确，2 密码错误，3 网络错误
//func (c *ChatClient) CheckPwd(name, pwd string) int
res = client.CheckPwd("username", "password")

//SetServeId 设置要访问的联系人 id，打开联系人的WEB页面前调用
//func (c *ChatClient) SetServeId(ident int64)
client.SetServeId(id)

//GetUrl 返回访问访问联系人WEB页面的URL，不会变化，只需调用一次，保存结果
//先用上面的SetServeId设置目标联系人，然后打开本URL访问其网页。
//func (c *ChatClient) GetUrl() string
url = client.GetUrl()

//GetHost 返回服务器的IP地址和端口，格式 IP:PORT
//func (c *ChatClient) GetHost() string
addr = client.GetHost()

//ChatTo 向指定ID发送文字聊天信息
//func (c *ChatClient) ChatTo(to int64, msg string)
client.ChatTo(id,message)

//MultiSend 向id列表中的所有人发送同一个文字信息
func (c *ChatClient) MultiSend(msg string, ids *IdArray)

//Tell 向指定ID发送上线通知
func (c *ChatClient) Tell(uid int64)

//TellAll 向id列表中的所有人发上线通知
func (c *ChatClient) TellAll(uids *IdArray)

//GetMsg 读取信息,网络断开或发生错误时返回null。阻塞函数，需要在线程中循环运行或者用异步函数包装
//func (c *ChatClient) GetMsg() *MsgType

//SendFile 发送文件，参数：id,pathname。 阻塞函数，需要在线程中循环运行或者用异步函数包装
func (c *ChatClient) SendFile(to int64, pathName string)

//AddFriend 加入联系人
func (c *ChatClient) AddFriend(fid int64)

//RemoveFriend 删除联系人
func (c *ChatClient) RemoveFriend(uid int64)

//GetProxyPort 返回代理端口
func (c *ChatClient) GetProxyPort() int 

//SetProxyPort 设置代理端口, port=0则恢复默认值
func (c *ChatClient) SetProxyPort(port int) int

//Quit 退出
func (c *ChatClient) Quit() 

//UpdateDesc 更新自述信息
func (c *ChatClient) UpdateDesc(param string) 

//DeleteMe 删除本用户，参数：用户名和密码
func (c *ChatClient) DeleteMe(name, pwd string) bool

//GetPgPath 返回程序所在路径
func (c *ChatClient) GetPgPath() string

//MoveStrangerToFriend 把留言的陌生人加入联系人名单，附带删除留言信息
func (c *ChatClient) MoveStrangerToFriend(fid int64)

//QueryID 接收到陌生人信息后调用，查询陌生人信息，返回 UserDataRet 结构体，参数：ID int64， msg stirng。 msg - 接受到的陌生人信息
//阻塞函数，最好在线程中运行或者用异步函数包装
func (c *ChatClient) QueryID(uid int64, msg string) *UserDataRet 

//UserStatus 参数：id int64，返回值：0-offline，1-online
//阻塞函数，最好在线程中运行或者用异步函数包装
func (c *ChatClient) UserStatus(uid int64) int
```

简单流程示例：

```
import XXXX/pchatclient

client := pchatclient.GetChatClient()  //初始化
client.Login(name,pwd)  //登录
friends := client.GetFriends()  //返回朋友列表
...  //show friends
stranger_msgs := client.GetStrangerMsgs()  //返回陌生人留言列表
...  //show stranger messages
...  //信息收发
```

## 关于 `GetMsg` 函数读取信息
返回类型为`MsgType`，成员Cmd类型为int8，只会得到2个值:0 - CmdChat, 27 - CmdSysReturn

1. 当 Cmd=0 时，信息来自其他用户用户，此时 From > 0
2. 当 Cmd=27 时，收到的是系统返回信息，此时 From=0

### Cmd=27 时全部系统返回信息，此时 From=0
1. "Offline ID"，其中ID是数字，表示该 ID 下线了。
2. "Version:NUM"，其中NUM是数字，表示github上最新客户端版本号。
3. "DELETE 1\n"，用户删除了自己。
4. "DELETE 0\n"，用户删除自己的操作失败了。
5. "ConnDown\n"，连接中断。

### Cmd=0 时的信息来自其他用户用户，此时 From > 0
这是如果Msg成员的前4个字符有3中情况：

1. "TEXT"，后面的的是文字聊天信息；
2. "JSON"，收到的是文件或图片，格式是"{Name:'该文件的保存路径',Mime:'mime-type',Size:size,Session:0}"，Session字段可以忽略，现阶段都是0；
3. "LOGI"，上线通知。

### UserDataArray 用法
方法 Next()
返回 UserDataRet 先调用 res = client.GetFriends()，返回值就是UserDataArray类型，
接着循环调用 res.Next()，直到返回值为 null

### IdArray 用于 TellAll 和 MultiSend
用 ids = NewIdArray() 初始化，
然后用 ids.Append(id) 添加数据。
