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
  for m := range l.in {
    var message Message
    err := message.Unmarshal(m)
    if err != nil {
      fmt.Println(err)
      continue
    }

    switch(message.MessageType) {
    case ACCEPT:
      fmt.Println("Accepted", message.Value)
      l.network.Close()
    }
  }
}
