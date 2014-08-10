package behavior

import (
	"github.com/elegios/topdown/server/ai/pathfind"
)

type PathToTarget struct {
	behaviorBase
}
type FollowPath struct {
	behaviorBase
}

func (p *PathToTarget) firstTick(a *AI) (contF, bool) {
	if a.target == nil || a.me.Pos.Mapid != a.target.Pos.Mapid {
		return p.tickUp(p, failure)
	}
	mPos := a.me.Pos
	tPos := a.target.Pos
	m := a.world.Maps[mPos.Mapid]
	f := func(x, y int) bool {
		return x < 0 || y < 0 || y >= len(m) || x >= len(m[y]) || m[y][x].Collides()
	}
	path, ok := pathfind.Find(f, mPos.X, mPos.Y, tPos.X, tPos.Y)
	if !ok {
		return p.tickUp(p, failure)
	}
	a.path = path
	return p.tickUp(p, success)
}

func (f *FollowPath) firstTick(a *AI) (contF, bool) {
	return f.tick(a)
}
func (f *FollowPath) tick(a *AI) (contF, bool) {
	for a.me.Actions >= 1 {
		if len(a.path) == 0 {
			return f.tickUp(f, success)
		}

		var direction string
		switch {
		case a.me.Pos.X < a.path[0].X:
			direction = "right"

		case a.me.Pos.X > a.path[0].X:
			direction = "left"

		case a.me.Pos.Y < a.path[0].Y:
			direction = "down"

		case a.me.Pos.Y > a.path[0].Y:
			direction = "up"
		}
		if !a.world.MoveCharacter(a.me, direction) {
			return f.tickUp(f, failure)
		}
		if a.me.Pos.X == a.path[0].X && a.me.Pos.Y == a.path[0].Y {
			a.path = a.path[1:]
		}
	}
	return running(f)
}
