package main

import (
  "fmt"
  "log"
  "net"
  "os"
)

func main() {
  if len(os.Args[1:]) < 1 {
    log.Fatalf("Missing role argument. Please run %s [network|client|proposer|acceptor|learner]\n", os.Args[0])
  }

  switch(os.Args[1]) {
  case "network":
    runNetwork(os.Args[2:])
  case "client":
    runClient(os.Args[2:])
  case "proposer":
    runProposer(os.Args[2:])
  case "acceptor":
    runAcceptor(os.Args[2:])
  case "learner":
    runLearner(os.Args[2:])
  default:
    log.Fatalf("Invalid role argument. Please run %s [network|client|proposer|acceptor|learner]\n", os.Args[0])
  }

  // network := Network{}
  // network.Init()
  // network.Clients[0].Propose(3)
  // network.Clients[0].Propose(8)
  // network.Clients[0].Propose(4)

  // network.Wait()
}

func runNetwork(args []string) {
  // network := Network{}
  address := ":0"
  if len(args) >= 1 {
    address = args[0]
  }

  listener, err := net.Listen("tcp", address)
  if err != nil {
    log.Fatalln(err)
  }

  fmt.Println("Network running on", listener.Addr().String())

  for {
    conn, err := listener.Accept()
    if err != nil {
      continue
    }

    fmt.Println(conn.RemoteAddr().String())
  }
}

func runClient(args []string) {
  if len(args) < 1 {
    log.Fatalf("Missing network address. Please run %s client [address]\n", os.Args[0])
  }

  conn, err := net.Dial("tcp", args[0])
  if err != nil {
    log.Fatalf("Could not connect to network: %v\n", err)
  }

  fmt.Printf("Connected to network at %s\n", conn.RemoteAddr())

  listener, err := net.Listen("tcp", ":0")
  if err != nil {
    log.Fatalln(err)
  }

  fmt.Println("Client running on", listener.Addr().String())

  for {
    conn, err := listener.Accept()
    if err != nil {
      continue
    }

    fmt.Println(conn.RemoteAddr().String())
  }
}

func runProposer(args []string) {
}

func runAcceptor(args []string) {
}

func runLearner(args []string) {
}
