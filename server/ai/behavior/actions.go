package behavior

type Walk struct {
	Direction string
	behaviorBase
}

func (w *Walk) firstTick(a *AI) (contF, bool) {
	return w.tick(a)
}
func (w *Walk) tick(a *AI) (contF, bool) {
	if a.me.Actions < 1 {
		return running(w)
	}
	if a.world.MoveCharacter(a.me, w.Direction) {
		return w.tickUp(w, success)
	}
	return w.tickUp(w, failure)
}
