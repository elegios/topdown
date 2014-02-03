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
	Props        []types.Prop              `json:"props"`
	Items        []types.Item              `json:"items"`
	positions    map[types.Position]bool
}

func sendTicks(conn *websocket.Conn, data *clientData) {
	ch := make(chan struct{}, 1)
	otm.AddP <- helpers.Pair{conn, ch}
	defer func() { otm.Rem <- conn }()
	for {
		<-ch
		slog.Println("Preparing to send tick to", conn.LocalAddr())
		update := tickUpdate{
			Maps:         make(map[string][]types.MapPos),
			Controllable: make([]string, 0),
			Characters:   make([]*types.Character, 0),
			Props:        make([]types.Prop, 0),
			Items:        make([]types.Item, 0),
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
		j := helpers.Max(c.Y-c.ViewDist, 0)
		maxJ := helpers.Min(c.Y+c.ViewDist, len(world.Maps[c.Mapname])-1)
		maxI := helpers.Min(c.X+c.ViewDist, len(world.Maps[c.Mapname][0])-1)
		pos := types.Position{
			Mapid: c.Mapname,
		}
		for ; j <= maxJ; j++ {
			pos.Y = j
			i := helpers.Max(c.X-c.ViewDist, 0)
			for ; i <= maxI; i++ {
				pos.X = i
				if update.positions[pos] ||
					(j-c.Y)*(j-c.Y)+(i-c.X)*(i-c.X) > c.ViewDist*c.ViewDist+1 ||
					!helpers.Visible(world.Maps[c.Mapname], c.X, c.Y, i, j) {
					continue
				}
				update.positions[types.Position{c.Mapname, i, j}] = true
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
			update.Props = append(update.Props, p)
		}
		if i, ok := world.MapItems[pos]; ok {
			update.Items = append(update.Items, i)
		}
	}
}
