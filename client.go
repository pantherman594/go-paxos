package main

import (
  "fmt"
  "time"
)

type Client struct {
  NodeBase
}
var _ = Node(&Client{})

func (c *Client) Init() {
  c.NodeBase.Init()
  go c.readMessages()
  go func() {
  time.Sleep(3  * time.Second)
  fmt.Println("PROPOSE")
  c.Propose(123)
  }()
}

func (c *Client) readMessages() {
  for message := range c.in {
    fmt.Println(message)
  }
}

func (c *Client) Propose(value int32) {
  message := Message{
    MessageType: PROPOSE,
    Value: value,
  }

  for id := range c.network.Proposers {
    fmt.Println("c to p")
    c.SendExternal(c.network.Nodes[id].address, message)
  }
}
