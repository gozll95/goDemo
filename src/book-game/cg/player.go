package cg

import(
	"fmt"
)

type Player struct {
	Name  string `json:"name"`
	Level int    `json:"level"`
	Exp   int    `json:"exp"`
	Room  int    `json:"room"`

	mq chan *Message `json:"mq"` //等待收取的消息
}

func NewPlayer()*Player{
	m:=make(chan *Message,1024)
	player:=&Player{"",0,0,0,m}
	go func(p *Player){
		for{
			msg:=<-p.mq
			fmt.Println(p.Name,"received message:",msg.Content)
		}
	}(player)

	return player
}