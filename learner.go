package main

import (
	"fmt"
)

type Learner struct {
	NodeBase
}

var _ = Node(&Learner{})

func (l *Learner) Init() {
	l.NodeBase.Init()
	go l.readMessages()
}

func (l *Learner) readMessages() {
	for message := range l.in {
		switch message.MessageType {
		case ACCEPT:
			fmt.Println("Accepted", message.Value)
		}
	}
}
