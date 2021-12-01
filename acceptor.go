package main

import (
  "fmt"
)

type Acceptor struct {
  NodeBase

  max_id int32
}
var _ = Node(&Acceptor{})

func (a *Acceptor) Init() {
  a.max_id = -1
  a.NodeBase.Init()
  go a.readMessages()
}

func (a *Acceptor) readMessages() {
  for message := range a.in {
    switch(message.MessageType) {
    case PREPARE:
      fmt.Println("RECEIVE PREPARE")
      if message.Id > a.max_id {
        a.max_id = message.Id
        go a.promise(message)
      }
    case PROPOSE:
      if a.max_id == message.Id {
        go a.accept(message.Id, message.Value)
      }
    }
  }
}

func (a *Acceptor) promise(message Message) {
  recipient := message.Sender
  message.MessageType = PROMISE
  message.Sender = 0

  fmt.Println("a to p promise")
  a.SendExternal(a.network.Nodes[recipient].address, message)
}

func (a *Acceptor) accept(id int32, value int32) {
  message := Message{
    MessageType: ACCEPT,
    Id: id,
    Value: value,
  }

  for id := range a.network.Learners {
    fmt.Println("a to l accept")
    a.SendExternal(a.network.Nodes[id].address, message)
  }
}
