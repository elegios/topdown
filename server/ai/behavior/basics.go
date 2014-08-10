package behavior

import (
	"math/rand"
)

type Selector struct {
	Bs []Behavior
	behaviorBase
}
type Sequence struct {
	Bs []Behavior
	behaviorBase
}
type Inverter struct {
	B Behavior
	behaviorBase
}
type Random struct {
	Bs []Behavior
	behaviorBase
}
type Parallel struct {
	Bs           []Behavior
	QuickSuccess bool
	QuickFailure bool
	behaviorBase
}
type Repeat struct {
	PropagateSuccess bool
	PropagateFailure bool
	B                Behavior
	behaviorBase
}
type Wait struct {
	behaviorBase
}

func (s *Selector) firstTick(a *AI) (contF, bool) {
	return s.Bs[0].firstTick(a)
}
func (s *Selector) tickFromBelow(a *AI, b Behavior, r behaviorResult) (contF, bool) {
	if r == success {
		return s.tickUp(s, r)
	}
	for i := range s.Bs {
		if b == s.Bs[i] && i+1 < len(s.Bs) {
			return s.Bs[i+1].firstTick(a)
		}
	}
	return s.tickUp(s, failure)
}
func (s *Selector) setParent(p behaviorParent) {
	s.parent = p
	for _, b := range s.Bs {
		b.setParent(s)
	}
}

func (s *Sequence) firstTick(a *AI) (contF, bool) {
	return s.Bs[0].firstTick(a)
}
func (s *Sequence) tickFromBelow(a *AI, b Behavior, r behaviorResult) (contF, bool) {
	if r == failure {
		return s.tickUp(s, failure)
	}
	for i := range s.Bs {
		if b == s.Bs[i] && i+1 < len(s.Bs) {
			return s.Bs[i+1].firstTick(a)
		}
	}
	return s.tickUp(s, success)
}
func (s *Sequence) setParent(p behaviorParent) {
	s.parent = p
	for _, b := range s.Bs {
		b.setParent(s)
	}
}

func (i *Inverter) firstTick(a *AI) (contF, bool) {
	return i.B.firstTick(a)
}
func (i *Inverter) tickFromBelow(a *AI, b Behavior, r behaviorResult) (contF, bool) {
	return i.tickUp(i, !r)
}
func (i *Inverter) setParent(p behaviorParent) {
	i.parent = p
	i.B.setParent(i)
}

func (r *Random) firstTick(a *AI) (contF, bool) {
	return r.Bs[rand.Intn(len(r.Bs))].firstTick(a)
}
func (r *Random) tickFromBelow(a *AI, _ Behavior, res behaviorResult) (contF, bool) {
	return r.tickUp(r, res)
}
func (r *Random) setParent(p behaviorParent) {
	r.parent = p
	for _, b := range r.Bs {
		b.setParent(r)
	}
}

func (p *Parallel) firstTick(a *AI) (contF, bool) {
	ns := make([]contF, len(p.Bs))
	for i, b := range p.Bs {
		cont := true
		ns[i] = b.firstTick
		for cont && ns[i] != nil {
			ns[i], cont = ns[i](a)
		}
	}
	return p.tick(ns), true
}
func (p *Parallel) tick(ns []contF) contF {
	return func(a *AI) (contF, bool) {
		cont := make([]bool, len(ns))
		for i, n := range ns {
			if n == nil {
				continue
			}
			ns[i], cont[i] = n(a)
			for cont[i] && ns[i] != nil {
				ns[i], cont[i] = ns[i](a)
			}
		}

		// everything has now run until completion or running
		allSuccess := !p.QuickSuccess
		allFailure := !p.QuickFailure
		for i, n := range ns {
			if n == nil {
				switch {
				case p.QuickFailure && !cont[i]:
					return p.tickUp(p, failure)

				case p.QuickSuccess && cont[i]:
					return p.tickUp(p, success)
				}
				allSuccess = allSuccess && cont[i]
				allFailure = allFailure && !cont[i]

			} else {
				allFailure = false
				allSuccess = false
			}
		}
		if allSuccess {
			return p.tickUp(p, success)
		}
		if allFailure {
			return p.tickUp(p, failure)
		}
		return p.tick(ns), false
	}
}
func (p *Parallel) tickFromBelow(_ *AI, _ Behavior, r behaviorResult) (contF, bool) {
	return nil, bool(r)
}
func (p *Parallel) setParent(parent behaviorParent) {
	p.parent = parent
	for _, b := range p.Bs {
		b.setParent(p)
	}
}

func (w *Wait) firstTick(_ *AI) (contF, bool) {
	return running(w)
}
func (w *Wait) tick(_ *AI) (contF, bool) {
	return w.tickUp(w, success)
}

func (r *Repeat) firstTick(a *AI) (contF, bool) {
	return r.B.firstTick(a)
}
func (r *Repeat) tickFromBelow(a *AI, _ Behavior, res behaviorResult) (contF, bool) {
	if r.PropagateFailure && res == failure {
		return r.tickUp(r, failure)
	}
	if r.PropagateSuccess && res == success {
		return r.tickUp(r, success)
	}
	return r.B.firstTick(a)
}
func (r *Repeat) setParent(p behaviorParent) {
	r.parent = p
	r.B.setParent(r)
}
