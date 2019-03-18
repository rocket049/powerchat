package main

import (
	"fmt"

	"github.com/rocket049/powerchat/pclientbackend"
)

func main() {
	client := pclientbackend.GetChatClient(".", "./config.json")
	defer client.Quit()
	me := client.Login("test1", "1234")
	fmt.Println(me)
	for {
		msg := client.GetMsg()
		if msg == nil {
			break
		}
		fmt.Println(msg.From, string(msg.Msg))
	}
}
