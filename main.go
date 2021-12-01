package main

import (
  "fmt"
  "log"
  "math/rand"
  "net"
  "os"
  "time"
)

func main() {
  rand.Seed(time.Now().UnixNano())

  if len(os.Args[1:]) < 1 {
    log.Fatalf("Missing role argument. Please run %s [network|client|proposer|acceptor|learner]\n", os.Args[0])
  }

  if os.Args[1] == "network" {
    runNetwork(os.Args[2:])
    return
  }

  var nodeType Message_NodeType

  switch(os.Args[1]) {
  case "client":
    nodeType = CLIENT
  case "proposer":
    nodeType = PROPOSER
  case "acceptor":
    nodeType = ACCEPTOR
  case "learner":
    nodeType = LEARNER
  default:
    log.Fatalf("Invalid role argument. Please run %s [network|client|proposer|acceptor|learner]\n", os.Args[0])
    return
  }

  runNode(nodeType, os.Args[2:])
}

func runNetwork(args []string) {
  address := ":0"
  if len(args) >= 1 {
    address = args[0]
  }

  listener, err := net.Listen("tcp", address)
  if err != nil {
    log.Fatalln(err)
  }

  fmt.Println("Network running on", listener.Addr().String())

  network := Network{
    Nodes: make(map[int32]NodeInfo),
    Clients: make(map[int32]struct{}),
    Proposers: make(map[int32]struct{}),
    Acceptors: make(map[int32]struct{}),
    Learners: make(map[int32]struct{}),
  }
  node := NodeBase{network: &network}
  node.Init()
  network.current = &node

  network.Run(listener)
}

func runNode(nodeType Message_NodeType, args []string) {
  var command string
  var name string

  switch (nodeType) {
  case CLIENT:
    command = "client"
    name = "Client"
  case PROPOSER:
    command = "proposer"
    name = "Proposer"
  case ACCEPTOR:
    command = "acceptor"
    name = "Acceptor"
  case LEARNER:
    command = "learner"
    name = "Learner"
  }

  if len(args) < 1 {
    log.Fatalf("Missing network address. Please run %s %s [address]\n", os.Args[0], command)
  }

  netConn, err := net.Dial("tcp", args[0])
  if err != nil {
    log.Fatalf("Could not connect to network: %v\n", err)
  }

  fmt.Printf("Connected to network at %s\n", netConn.RemoteAddr())

  listener, err := net.Listen("tcp", ":0")
  if err != nil {
    log.Fatalln(err)
  }

  fmt.Printf("%s running on %s\n", name, listener.Addr().String())

  network := &Network{
    id: int32(rand.Int()),
    Nodes: make(map[int32]NodeInfo),
    Clients: make(map[int32]struct{}),
    Proposers: make(map[int32]struct{}),
    Acceptors: make(map[int32]struct{}),
    Learners: make(map[int32]struct{}),
  }

  var node Node

  switch nodeType {
  case CLIENT:
    node = &Client{
      NodeBase: NodeBase{
        network: network,
      },
    }
  case PROPOSER:
    node = &Proposer{
      NodeBase: NodeBase{
        network: network,
      },
      name: network.id,
    }
  case ACCEPTOR:
    node = &Acceptor{
      NodeBase: NodeBase{
        network: network,
      },
    }
  case LEARNER:
    node = &Learner{
      NodeBase: NodeBase{
        network: network,
      },
    }
  default:
    panic("Invalid node type")
  }

  network.current = node

  node.Init()

  node.SendExternal(netConn.RemoteAddr().String(), Message{
    Sender: network.id,
    MessageType: ADD_NODE,
    NodeType: nodeType,
    Address: listener.Addr().String(),
  })

  for {
    conn, err := listener.Accept()
    if err != nil {
      continue
    }

    fmt.Println(conn.RemoteAddr().String())

    go network.handleNode(conn)
  }
}
