package pathfind

type container struct {
	first *containerLink
}
type containerLink struct {
	f    int
	ns   []*node
	next *containerLink
}

func h(x1, y1, tx, ty int) int {
	return x1 - tx + y1 - ty
}

func (c *container) add(n *node, tx, ty int) {
	if n == nil {
		return
	}
	f := n.dist + h(n.x, n.y, tx, ty)

	if c.first == nil {
		c.first = &containerLink{
			f:  f,
			ns: []*node{n},
		}
		return
	}
	prev := &c.first
	link := c.first

	for ; link != nil; prev, link = &link.next, link.next {
		switch {
		case link.f < f:
			continue

		case link.f == f:
			link.ns = append(link.ns, n)
			return

		default:
			*prev = &containerLink{
				f:    f,
				ns:   []*node{n},
				next: link,
			}
			return
		}
	}

	*prev = &containerLink{
		f:  f,
		ns: []*node{n},
	}
}

func (c *container) popSmallest() (n *node) {
	if c.first == nil {
		panic("tried to pop smallest when there is nothing")
	}

	n = c.first.ns[0]
	if len(c.first.ns) == 1 {
		c.first = c.first.next
	} else {
		c.first.ns = c.first.ns[1:]
	}
	return
}
