package main

import (
	"fmt"
	"io/ioutil"

	"github.com/rocket049/powerchat/pclientbackend"
)

func main() {
	cfg, err := ioutil.ReadFile("./config.json")
	if err != nil {
		panic(err)
	}
	client := pclientbackend.GetChatClient(".", string(cfg))
	defer client.Quit()
	me := client.Login("test1", "1234")
	fmt.Println("Login:", me)
	for {
		msg := client.GetMsg()
		if msg == nil {
			break
		}
		fmt.Println(msg.From, string(msg.Msg))
	}
}
