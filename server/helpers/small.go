package helpers

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
func Minf(a, b float32) float32 {
	if a < b {
		return a
	}
	return b
}

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
func Maxf(a, b float32) float32 {
	if a > b {
		return a
	}
	return b
}

func Ceil(a float32) float32 {
	if float32(int(a)) == a {
		return a
	}
	if a > 0 {
		return float32(int(a + 1))
	}
	return float32(int(a))
}

func Abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}
