package mathutil

func GCD(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func IntPow(a, b int) int {
	result := 1
	for i := 0; i < b; i++ {
		result *= a
	}
	return result
}

func Abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}
