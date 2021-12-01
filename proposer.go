package main

import (
	"fmt"
	"sync/atomic"
)

type Proposed struct {
	value    int32
	accepted uint32
}

type Proposer struct {
	NodeBase

	name    int32
	next_id int32

	proposed map[int32]*Proposed
}

var _ = Node(&Proposer{})

func (p *Proposer) Init() {
	p.NodeBase.Init()
	p.proposed = make(map[int32]*Proposed)
	go p.readMessages()
}

func (p *Proposer) readMessages() {
	for message := range p.in {
		switch message.MessageType {
		case PROPOSE:
			id := p.next_id
			p.next_id += 1

			p.proposed[id] = &Proposed{
				value:    message.Value,
				accepted: 0,
			}

			go p.prepare(id, message.Value)
		case PROMISE:
			p.receivePromise(message)
		}
	}
}

func (p *Proposer) prepare(id int32, value int32) {
	message := Message{
		MessageType: PREPARE,
		Sender:      p.name,
		Id:          id,
		Value:       value,
	}

	for id := range p.network.Acceptors {
		fmt.Println("p to a prepare")
		p.SendExternal(p.network.Nodes[id].address, message)
	}
}

func (p *Proposer) receivePromise(message Message) {
	fmt.Println("recv promise")
	prop := p.proposed[message.Id]
	accepted := atomic.AddUint32(&prop.accepted, 1)

	fmt.Println(accepted, len(p.network.Acceptors))
	if float64(accepted) > float64(len(p.network.Acceptors))/2 {
		fmt.Println("propose", accepted)
		go p.propose(message.Id, prop.value)
	}
}

func (p *Proposer) propose(id int32, value int32) {
	message := Message{
		MessageType: PROPOSE,
		Sender:      p.name,
		Id:          id,
		Value:       value,
	}

	for id := range p.network.Acceptors {
		fmt.Println("p to a propose")
		p.SendExternal(p.network.Nodes[id].address, message)
	}
}
