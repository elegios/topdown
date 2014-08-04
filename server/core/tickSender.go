package core

import (
	"code.google.com/p/go.net/websocket"
	"github.com/elegios/topdown/server/helpers"
	"github.com/elegios/topdown/server/types"
)

type tickUpdate struct {
	Maps         map[string][]types.MapPos `json:"maps"`
	Controllable []string                  `json:"controllable"`
	Characters   []*types.Character        `json:"characters"`
	Props        map[string][]types.Prop   `json:"props"`
	Items        map[string][]types.Item   `json:"items"`
	Updates      []interface{}             `json:"updates"`
	positions    map[types.Position]bool
}

func sendTicks(conn *websocket.Conn, data *clientData) {
	ch := make(chan struct{}, 1)
	otm.AddP <- helpers.Pair{conn, ch}
	defer func() { otm.Rem <- conn }()
	for {
		<-ch
		slog.Println("Preparing to send tick to", conn.Request().RemoteAddr)
		update := tickUpdate{
			Maps:         make(map[string][]types.MapPos),
			Controllable: make([]string, 0),
			Characters:   make([]*types.Character, 0),
			Props:        make(map[string][]types.Prop),
			Items:        make(map[string][]types.Item),
			Updates:      make([]interface{}, 0),
			positions:    make(map[types.Position]bool),
		}
		updateControllable(data, &update)
		collectMapPos(&update)
		collectVisible(&update)
		//done collecting data
		otm.Done()
		err := websocket.JSON.Send(conn, &update)
		if err != nil {
			wlog.Println("ticksender died:", err)
			return
		}
	}
}

func updateControllable(data *clientData, update *tickUpdate) {
	data.Lock()
	defer data.Unlock()
	for charId := range data.charIds {
		_, ok := world.Charids[charId]
		if !ok {
			delete(data.charIds, charId)
			continue
		}
		update.Controllable = append(update.Controllable, charId)
	}
}

func collectMapPos(update *tickUpdate) {
	for _, cid := range update.Controllable {
		c := world.Charids[cid]
		p := c.Pos
		visionFunc := func(x, y int) bool {
			return world.Maps[p.Mapid][y][x].BlocksVision()
		}
		j := helpers.Max(p.Y-c.ViewDist, 0)
		maxJ := helpers.Min(p.Y+c.ViewDist, len(world.Maps[p.Mapid])-1)
		maxI := helpers.Min(p.X+c.ViewDist, len(world.Maps[p.Mapid][0])-1)
		pos := types.Position{
			Mapid: p.Mapid,
		}
		for ; j <= maxJ; j++ {
			pos.Y = j
			i := helpers.Max(p.X-c.ViewDist, 0)
			for ; i <= maxI; i++ {
				pos.X = i
				if update.positions[pos] ||
					(j-p.Y)*(j-p.Y)+(i-p.X)*(i-p.X) > c.ViewDist*c.ViewDist+1 ||
					!helpers.Visible(visionFunc, p.X, p.Y, i, j) {
					continue
				}
				update.positions[types.Position{p.Mapid, i, j}] = true
			}
		}
	}
	for pos := range update.positions {
		update.Maps[pos.Mapid] = append(update.Maps[pos.Mapid], pos.GetMapPos(world.Maps[pos.Mapid]))
	}
}

func collectVisible(update *tickUpdate) {
	for pos := range update.positions {
		if c, ok := world.MapCharacters[pos]; ok {
			update.Characters = append(update.Characters, c)
		}
		if p, ok := world.MapProps[pos]; ok {
			update.Props[pos.Mapid] = append(update.Props[pos.Mapid], p)
		}
		if i, ok := world.MapItems[pos]; ok {
			update.Items[pos.Mapid] = append(update.Items[pos.Mapid], i)
		}
	}
	for _, u := range world.Updates {
		if update.positions[u.Pos] {
			update.Updates = append(update.Updates, u.Content)
		}
	}
}
