package main

import (
	"code.google.com/p/go.net/websocket"
	"sync"
)

type Pair struct {
	Conn *websocket.Conn
	Out  chan struct{}
}

type OneToMany struct {
	In   chan struct{}
	AddP chan Pair
	Rem  chan *websocket.Conn
	Out  map[*websocket.Conn]chan struct{}
	sync.WaitGroup
}

func newOneToMany() *OneToMany {
	return &OneToMany{
		In:   make(chan struct{}),
		AddP: make(chan Pair),
		Rem:  make(chan *websocket.Conn),
		Out:  make(map[*websocket.Conn]chan struct{}),
	}
}

func (o *OneToMany) Listen() {
	for {
		select {
		case e := <-o.In:
			o.Add(len(o.Out))
			for _, c := range o.Out {
				c <- e
			}
			o.Done()

		case p := <-o.AddP:
			o.Out[p.Conn] = p.Out

		case c := <-o.Rem:
			delete(o.Out, c)
		}
	}
}
