package main

import (
  "fmt"
)

type Client struct {
  NodeBase
}
var _ = Node(&Client{})

func (c *Client) Init() {
  c.NodeBase.Init()
  go c.readMessages()
}

func (c *Client) readMessages() {
  for message := range c.in {
    fmt.Println(message)
  }
}

func (c *Client) Propose(value int32) {
  for _, p := range c.network.Proposers {
    fmt.Println("c to p")
    p.Send(Message{
      MessageType: PROPOSE,
      Value: value,
    })
  }
}
