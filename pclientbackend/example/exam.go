package main

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/rocket049/powerchat/pclientbackend"
)

func main() {
	cfg, err := ioutil.ReadFile("./config.json")
	if err != nil {
		panic(err)
	}
	client := pclientbackend.GetChatClient(".", string(cfg))
	me := client.Login("test1", "1234")
	fmt.Println("Login:", me)
	for i := 0; i < 5; i++ {
		msg := client.GetMsg()
		if msg == nil {
			break
		}
		fmt.Println(msg.From, string(msg.Msg))
	}

	client = pclientbackend.GetChatClient(".", string(cfg))
	me = client.Login("test1", "1234")
	fmt.Println("Login Again:", me)
	go exportsTest(client)
	limitRecv(client, 10)
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
}
