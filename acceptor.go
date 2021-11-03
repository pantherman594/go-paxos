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
  a.NodeBase.Init()
  go a.readMessages()
}

func (a *Acceptor) readMessages() {
  for m := range a.in {
    var message Message
    err := message.Unmarshal(m)
    if err != nil {
      fmt.Println(err)
      continue
    }

    switch(message.MessageType) {
    case PREPARE:
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

  for _, p := range a.network.Proposers {
    if p.name != recipient {
      continue
    }

    fmt.Println("a to p")
    p.Send(message)
    break
  }
}

func (a *Acceptor) accept(id int32, value int32) {
  for _, l := range a.network.Learners {
    fmt.Println("a to l")
    l.Send(Message{
      MessageType: ACCEPT,
      Id: id,
      Value: value,
    })
  }
}
