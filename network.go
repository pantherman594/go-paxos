package main

import (
	"fmt"
	"net"
	"sync"
)

type NodeInfo struct {
	id       int32
	address  string
	nodeType Message_NodeType
}

type Network struct {
	id int32

	Nodes     map[int32]NodeInfo
	Clients   map[int32]struct{}
	Proposers map[int32]struct{}
	Acceptors map[int32]struct{}
	Learners  map[int32]struct{}

	conn    *net.Conn
	current Node

	close_lock sync.Mutex
	closed     bool
	done       chan struct{}
}

func (n *Network) Run(listener net.Listener) {
	n.closed = false
	n.done = make(chan struct{})

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}

		fmt.Println(conn.RemoteAddr().String())
		go n.handleNode(conn)
	}
}

func (net *Network) handleNode(conn net.Conn) {
	defer conn.Close()

	var buf [512]byte
	n, err := conn.Read(buf[:])
	if err != nil {
		fmt.Println(err)
		return
	}

	net.current.Receive(buf[:n])
}

func (n *Network) AddNode(id int32, nodeType Message_NodeType, addr string) {
	if id == n.id {
		return
	}

	if _, ok := n.Nodes[id]; ok {
		return
	}

	n.Nodes[id] = NodeInfo{id, addr, nodeType}

	switch nodeType {
	case CLIENT:
		n.Clients[id] = struct{}{}
	case PROPOSER:
		n.Proposers[id] = struct{}{}
	case ACCEPTOR:
		n.Acceptors[id] = struct{}{}
	case LEARNER:
		n.Learners[id] = struct{}{}
	default:
		panic("Invalid node type: " + nodeType.String())
	}

	addMessage := Message{
		Sender:      id,
		MessageType: ADD_NODE,
		NodeType:    nodeType,
		Address:     addr,
	}

	for nodeId, nodeInfo := range n.Nodes {
		n.current.SendExternal(nodeInfo.address, addMessage)

		if n.id == 0 {
			addNodeMessage := Message{
				Sender:      nodeId,
				MessageType: ADD_NODE,
				NodeType:    nodeInfo.nodeType,
				Address:     nodeInfo.address,
			}

			n.current.SendExternal(addr, addNodeMessage)
		}
	}
}

type Node interface {
	Init()
	Receive(buffer []byte) error
	SendExternal(recipient string, message Message) error
	readMessages()
}

type NodeBase struct {
	network *Network
	in      chan Message
}

func (n *NodeBase) Init() {
	n.in = make(chan Message, 10)
}

func (n *NodeBase) Receive(buf []byte) error {
	var message Message
	err := message.Unmarshal(buf)
	if err != nil {
		return err
	}

	switch message.MessageType {
	case ADD_NODE:
		n.network.AddNode(message.Sender, message.NodeType, message.Address)
	default:
		fmt.Println(message)
		n.in <- message
	}

	return nil
}

func (n *NodeBase) SendExternal(recipient string, message Message) error {
	fmt.Println("send", message, "to", recipient)
	buf, err := message.Marshal()
	if err != nil {
		return err
	}

	fmt.Println(buf)

	conn, err := net.Dial("tcp", recipient)
	if err != nil {
		return err
	}

	fmt.Println("conn")

	_, err = conn.Write(buf)
	if err != nil {
		return err
	}

	fmt.Println("write")

	return nil
}

func (n *NodeBase) readMessages() {}

var _ = Node(&NodeBase{})
