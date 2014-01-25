package main

func visible(bs [][]bits, x1, y1, x2, y2 int) bool {
	if x1 == x2 && y1 == y2 {
		return true
	}
	if abs(x1-x2) > abs(y1-y2) {
		dir := 1
		if x1 > x2 {
			dir = -1
		}
		k := k(x1, y1, x2, y2)
		m := m(x1, y1, k)
		for i := x1 + dir; i != x2; i += dir {
			if bs[f(k, i, m, dir)][i].blocksVision() {
				return false
			}
		}
	} else {
		dir := 1
		if y1 > y2 {
			dir = -1
		}
		k := k(y1, x1, y2, x2)
		m := m(y1, x1, k)
		for j := y1 + dir; j != y2; j += dir {
			if bs[j][f(k, j, m, dir)].blocksVision() {
				return false
			}
		}
	}

	return true
}

func k(x1, y1, x2, y2 int) float32 {
	if x1 == x2 {
		return 0
	}
	return float32(y1-y2) / float32(x1-x2)
}

func m(x, y int, k float32) float32 {
	return float32(y) - k*float32(x)
}

func f(k float32, x int, m float32, dir int) int {
	y := k*(float32(x)) + m + 0.5
	//if y-float32(int(y)) > 0.5 {
	//	return int(y) + 1
	//}
	return int(y)
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}
