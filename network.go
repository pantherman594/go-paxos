package main

import (
  "net"
  "sync"
)

type Network struct {
  Clients []*Client
  Proposers []*Proposer
  Acceptors []*Acceptor
  Learners []*Learner

  conn *net.Conn

  close_lock sync.Mutex
  closed bool
  done chan struct{}
}

func (n *Network) Init() {
  n.closed = false
  n.done = make(chan struct{})
  for i := 0; i < 1; i++ {
    c := &Client{
      NodeBase: NodeBase{ network: n },
    }

    n.Clients = append(n.Clients, c)
    c.Init()
  }

  for i := 0; i < 1; i++ {
    p := &Proposer{
      NodeBase: NodeBase{ network: n },
      name: 0,
      next_id: 0,
    }

    n.Proposers = append(n.Proposers, p)
    p.Init()
  }

  for i := 1; i < 4; i++ {
    a := &Acceptor{
      NodeBase: NodeBase{ network: n },
      max_id: -1,
    }

    n.Acceptors = append(n.Acceptors, a)
    a.Init()
  }

  for i := 0; i < 1; i++ {
    l := &Learner{
      NodeBase: NodeBase{ network: n },
    }

    n.Learners = append(n.Learners, l)
    l.Init()
  }
}

func (n *Network) Wait() {
  _, _ = <- n.done
}

func (n *Network) Close() {
  n.close_lock.Lock()
  defer n.close_lock.Unlock()

  if n.closed {
    return
  }

  for _, p := range n.Proposers {
    close(p.in)
  }
  for _, a := range n.Acceptors {
    close(a.in)
  }
  for _, l := range n.Learners {
    close(l.in)
  }

  close(n.done)
  n.closed = true
}

type Node interface {
  Init()
  Send(message Message) error
  readMessages()
}

type NodeBase struct {
  network *Network
  conn *net.Conn
  in chan []byte
}

func (n *NodeBase) Init() {
  n.in = make(chan []byte, 10)
}

func (n *NodeBase) Send(message Message) error {
  b, err := message.Marshal()
  if err != nil {
    return err
  }

  n.in <- b

  return nil
}

func (n *NodeBase) readMessages() { }

var _ = Node(&NodeBase{})
