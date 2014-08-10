package behavior

import (
	"github.com/elegios/topdown/server/ai/pathfind"
	"github.com/elegios/topdown/server/types"
)

type AI struct {
	nextF  contF
	me     *types.Character
	world  *types.World
	target *types.Character
	path   []pathfind.Point
}

func NewAI(me *types.Character, world *types.World, b Behavior) *AI {
	return &AI{
		nextF: b.firstTick,
		me:    me,
		world: world,
	}
}
func (a *AI) Tick() {
	cont := true
	for cont {
		a.nextF, cont = a.nextF(a)
	}
}

type behaviorResult bool

const (
	success behaviorResult = true
	failure behaviorResult = false
)

type contF func(*AI) (contF, bool)

func Finalize(b Behavior) {
	b.setParent(nil)
}

type Behavior interface {
	firstTick(*AI) (contF, bool)
	setParent(behaviorParent)
}
type behaviorParent interface {
	tickFromBelow(*AI, Behavior, behaviorResult) (contF, bool)
}
type behaviorRunning interface {
	tick(*AI) (contF, bool)
}

func running(b behaviorRunning) (contF, bool) {
	return b.tick, false
}

type behaviorBase struct {
	parent behaviorParent
}

func (b *behaviorBase) setParent(p behaviorParent) {
	b.parent = p
}
func (b *behaviorBase) tickUp(self Behavior, r behaviorResult) (contF, bool) {
	return func(a *AI) (contF, bool) {
		return b.parent.tickFromBelow(a, self, r)
	}, true
}
