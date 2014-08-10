package behavior

type HasTarget struct {
	behaviorBase
}
type TargetRandomCharacter struct {
	behaviorBase
}

func (h *HasTarget) firstTick(a *AI) (contF, bool) {
	if a.target != nil {
		return h.tickUp(h, success)
	}
	return h.tickUp(h, failure)
}

func (t *TargetRandomCharacter) firstTick(a *AI) (contF, bool) {
	for p, c := range a.world.MapCharacters {
		if p.Mapid == a.me.Pos.Mapid && c.Id != a.me.Id {
			a.target = c
			return t.tickUp(t, success)
		}
	}
	return t.tickUp(t, failure)
}
