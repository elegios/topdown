package core

import (
	"code.google.com/p/go.net/websocket"
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

type request struct {
	message map[string]string
	ch      chan map[string]string
}

func clientHandler(ch chan<- request) func(*websocket.Conn) {
	return func(conn *websocket.Conn) {
		slog.Println("Got a websocket connection")
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
				wlog.Println("A websocket connection died:", err)
				return
			}

			if !ok(mess, data) {
				wlog.Print("A message was not ok:", mess)
				continue
			}

			ch <- request{mess, locChan}
			switch mess["command"] {
			case "create":
				res := <-locChan
				slog.Printf("got return from create")
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
