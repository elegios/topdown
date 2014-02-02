package main

import (
	"code.google.com/p/go.net/websocket"
	"log"
	"sync"
)

//this data is shared, message relating to a
//character may be accepted even though the
//character has been killed, as knowledge of
//death will not be propagated until a tick.
//It means the server has to
//check that the id really corresponds to a
//character.
type clientData struct {
	charIds map[string]struct{}
	sync.Mutex
}

type mapPos struct {
	X    int  `json:"x"`
	Y    int  `json:"y"`
	Data bits `json:"data"`
}

type tickUpdate struct {
	Maps         map[string][]mapPos `json:"maps"`
	Controllable []string            `json:"controllable"`
	Characters   []*Character        `json:"characters"`
	Props        []Prop              `json:"props"`
	Items        []Item              `json:"items"`
}

type request struct {
	message map[string]string
	ch      chan map[string]string
}

func clientHandler(ch chan<- request) func(*websocket.Conn) {
	return func(conn *websocket.Conn) {
		//TODO: load owned characters
		data := &clientData{
			charIds: make(map[string]struct{}),
		}
		go sendTicks(conn, data)
		locChan := make(chan map[string]string, 1)
		for {
			var mess map[string]string
			err := websocket.JSON.Receive(conn, &mess)
			if err != nil {
				warn.Println("A websocket connection died:", err)
				return
			}

			if !ok(mess, data) {
				warn.Print("A message was not ok:", mess)
				continue
			}

			ch <- request{mess, locChan}
			switch mess["command"] {
			case "create":
				res := <-locChan
				log.Printf("got return from create")
				data.Lock()
				data.charIds[res["id"]] = struct{}{}
				data.Unlock()

			case "itemrequest":
				websocket.JSON.Send(conn, <-locChan)
			}
		}
	}
}

func ok(mess map[string]string, data *clientData) bool {
	if _, ok := mess["character"]; ok {
		data.Lock()
		defer data.Unlock()
		_, ok = data.charIds[mess["character"]]
		return ok
	}
	return true
}
