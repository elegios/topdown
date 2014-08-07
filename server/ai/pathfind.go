package ai

import (
	"github.com/elegios/topdown/server/helpers"
)

type Point struct {
	X, Y int
}

type node struct {
	prev *node
	dist int
	x, y int
}

// Will (probably, I haven't proven it) find a fastest path
// if there is one, but loop infinitely otherwise
func Find(bs func(int, int) bool, x, y, tx, ty int) []Point {
	open := new(container)
	start := &node{dist: 0, x: x, y: y}

	n := nextHori(1, bs, start, tx, ty)
	if isEnd(n, tx, ty) {
		return buildPath(n)
	}
	open.add(n, tx, ty)

	n = nextHori(-1, bs, start, tx, ty)
	if isEnd(n, tx, ty) {
		return buildPath(n)
	}
	open.add(n, tx, ty)

	n = nextVert(1, bs, start, tx, ty)
	if isEnd(n, tx, ty) {
		return buildPath(n)
	}
	open.add(n, tx, ty)

	n = nextVert(-1, bs, start, tx, ty)
	if isEnd(n, tx, ty) {
		return buildPath(n)
	}
	open.add(n, tx, ty)

	for {
		current := open.popSmallest()

		var n1, n2, n3 *node
		switch {
		case current.x < current.prev.x:
			n1 = nextHori(-1, bs, current, tx, ty)
			n2 = nextVert(-1, bs, current, tx, ty)
			n3 = nextVert(1, bs, current, tx, ty)

		case current.x > current.prev.x:
			n1 = nextHori(1, bs, current, tx, ty)
			n2 = nextVert(-1, bs, current, tx, ty)
			n3 = nextVert(1, bs, current, tx, ty)

		case current.y < current.prev.y:
			n1 = nextVert(-1, bs, current, tx, ty)
			n2 = nextHori(-1, bs, current, tx, ty)
			n3 = nextHori(1, bs, current, tx, ty)

		case current.y > current.prev.y:
			n1 = nextVert(1, bs, current, tx, ty)
			n2 = nextHori(-1, bs, current, tx, ty)
			n3 = nextHori(1, bs, current, tx, ty)
		}

		if isEnd(n1, tx, ty) {
			return buildPath(n1)
		}
		if isEnd(n2, tx, ty) {
			return buildPath(n2)
		}
		if isEnd(n3, tx, ty) {
			return buildPath(n3)
		}

		open.add(n1, tx, ty)
		open.add(n2, tx, ty)
		open.add(n3, tx, ty)
	}
}

func isEnd(n *node, tx, ty int) bool {
	return n != nil && n.x == tx && n.y == ty
}

func buildPath(end *node) []Point {
	steps := make([]Point, 0)
	for end.prev != nil {
		steps = append(steps, Point{X: end.x, Y: end.y})
		end = end.prev
	}
	for i := 0; i < len(steps)/2; i++ {
		steps[i], steps[len(steps)-i-1] = steps[len(steps)-i-1], steps[i]
	}
	return steps
}

func nextHori(d int, bs func(int, int) bool, n *node, tx, ty int) *node {
	x, y := n.x+d, n.y
	for ; !bs(x, y); x += d {
		if (x == tx && y == ty) ||
			(bs(x-d, y-1) && !bs(x, y-1)) ||
			(bs(x-d, y+1) && !bs(x, y+1)) {

			return &node{
				prev: n,
				dist: n.dist + helpers.Abs(x-n.x),
				x:    x,
				y:    y,
			}
		}
	}

	return nil
}

func nextVert(d int, bs func(int, int) bool, n *node, tx, ty int) *node {
	movingNode := &node{
		dist: n.dist + 1,
		prev: n,
		x:    n.x,
		y:    n.y + d,
	}
	for !bs(movingNode.x, movingNode.y) {

		if (movingNode.x == tx && movingNode.y == ty) ||
			nextHori(1, bs, movingNode, tx, ty) != nil ||
			nextHori(-1, bs, movingNode, tx, ty) != nil {

			return movingNode
		}

		movingNode.y += d
		movingNode.dist++
	}

	return nil
}

func h(x1, y1, tx, ty int) int {
	return x1 - tx + y1 - ty
}

type container struct {
	first *containerLink
}
type containerLink struct {
	f    int
	ns   []*node
	next *containerLink
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
