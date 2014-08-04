package helpers

func Visible(blocksVision func(int, int) bool, x1, y1, x2, y2 int) bool {
	cx := x1
	cy := y1

	dx := x2 - cx
	dy := y2 - cy
	if dx < 0 {
		dx = 0 - dx
	}
	if dy < 0 {
		dy = 0 - dy
	}

	var sx int
	var sy int
	if cx < x2 {
		sx = 1
	} else {
		sx = -1
	}
	if cy < y2 {
		sy = 1
	} else {
		sy = -1
	}
	err := dx - dy

	for {
		if blocksVision(cx, cy) {
			return false
		}
		if (cx == x2) && (cy == y2) {
			return true
		}
		e2 := 2 * err
		if e2 > (0 - dy) {
			err -= dy
			cx += sx
		}
		if e2 < dx {
			err += dx
			cy += sy
		}
	}
}
