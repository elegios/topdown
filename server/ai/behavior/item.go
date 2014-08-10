package behavior

type StandingOnItem struct {
	behaviorBase
}
type PickupItem struct {
	behaviorBase
}

func (s *StandingOnItem) firstTick(a *AI) (contF, bool) {
	if _, ok := a.world.MapItems[a.me.Pos]; ok {
		return s.tickUp(s, success)
	}
	return s.tickUp(s, failure)
}

func (p *PickupItem) firstTick(a *AI) (contF, bool) {
	return p.tick(a)
}
func (p *PickupItem) tick(a *AI) (contF, bool) {
	if a.me.Actions < 1 {
		return running(p)
	}
	if a.world.PickupItem(a.me) {
		return p.tickUp(p, success)
	}
	return p.tickUp(p, failure)
}
