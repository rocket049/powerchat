package main

import (
	"fmt"
	"io/ioutil"
	"time"

	"pclientbackend"
)

func main() {
	cfg, err := ioutil.ReadFile("./config.json")
	if err != nil {
		panic(err)
	}
	client := pclientbackend.GetChatClient(".", string(cfg))
	time.Sleep(time.Second * 3)
	me := client.Login("test1", "1234")
	fmt.Println("Login:", me)
	//fmt.Println(client.SendFile(1, "/home/fuhz/pclientbackend-sources.jar"))
	//limitRecv(client, 10)

	//client = pclientbackend.GetChatClient(".", string(cfg))
	//me = client.Login("test1", "1234")
	//fmt.Println("Login Again:", me)
	//go exportsTest(client)
	limitRecv(client, 20)
	client.Quit()
}

func limitRecv(client *pclientbackend.ChatClient, n int) {
	for i := 0; i < n; i++ {
		msg := client.GetMsg()
		if msg == nil {
			break
		}
		fmt.Println(msg.From, string(msg.Msg))
	}
}

func exportsTest(client *pclientbackend.ChatClient) {
	var n time.Duration = 1
	client.AddFriend(20)
	time.Sleep(time.Second * n)
	client.ChatTo(1, "hello")
	time.Sleep(time.Second * n)
	client.Tell(1)
	time.Sleep(time.Second * n)
	client.UpdateDesc("test user " + time.Now().Format("2006-01-02 15:04:05"))
	time.Sleep(time.Second * n)
	client.NewPasswd("test1", "1234", "123456")
	fmt.Println("NewPwd Ret:", client.CheckPwd("test1", "1234"))
	fmt.Println("NewPwd Ret:", client.CheckPwd("test1", "123456"))
	time.Sleep(time.Second * n)
	client.NewPasswd("test1", "123456", "1234")
	fmt.Println("NewPwd Ret:", client.CheckPwd("test1", "1234"))
}
