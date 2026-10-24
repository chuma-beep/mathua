package complex

import (
	"fmt"
	"math/rand"

	"github.com/chuma-beep/mathua/internal/generator"
)

// The Euler identity concept held four situations: one question shape at each
// quarter turn. Everything below is computed from a drawn parameter, so the
// learner meets the identity applied rather than recited.

var eulerUnits = []string{"i", "-1", "-i", "1"}

// eulerAngle renders k*pi/2 for k in 1..8 the way a textbook writes it.
func eulerAngle(k int) string {
	switch k {
	case 1:
		return `\pi/2`
	case 2:
		return `\pi`
	case 3:
		return `3\pi/2`
	case 4:
		return `2\pi`
	case 5:
		return `5\pi/2`
	case 6:
		return `3\pi`
	case 7:
		return `7\pi/2`
	case 8:
		return `4\pi`
	}
	return `0`
}

// eulerOther returns the three units that are not the answer, which is exactly
// how many distractors a four-option question needs.
func eulerOther(answer string) []string {
	out := make([]string, 0, 3)
	for _, u := range eulerUnits {
		if u != answer {
			out = append(out, u)
		}
	}
	return out
}

// twiceCosine gives e^{i d} + e^{-i d} for the four angles asked below.
func twiceCosine(deg int) string {
	switch deg {
	case 60:
		return "1"
	case 90:
		return "0"
	case 120:
		return "-1"
	}
	return "-2"
}

func (g *eulerIdentityGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	switch rand.Intn(5) {
	case 0:
		// A full turn in quarter steps, so every value of the unit circle is
		// reachable and the argument never prints the answer.
		k := 1 + rand.Intn(8)
		answer := eulerUnits[(k-1)%4]
		return complexChoice(
			fmt.Sprintf("Use Euler's formula to evaluate \\(e^{i%s}\\).", eulerAngle(k)),
			answer,
			fmt.Sprintf("\\(e^{i\\theta}=\\cos\\theta+i\\sin\\theta\\). At \\(\\theta=%s\\) the cosine and sine give %s.", eulerAngle(k), answer),
			eulerOther(answer)...,
		)
	case 1:
		// Whole multiples of pi: the value alternates rather than cycling.
		k := 1 + rand.Intn(6)
		answer := "1"
		other := "-1"
		parity := "even"
		if k%2 == 1 {
			answer = "-1"
			other = "1"
			parity = "odd"
		}
		return complexChoice(
			fmt.Sprintf("Use Euler's formula to evaluate \\(e^{i%d\\pi}\\).", k),
			answer,
			fmt.Sprintf("An integer multiple of \\(\\pi\\) lands on the real axis. Since %d is %s, the value is %s.", k, parity, answer),
			other, "i", "-i",
		)
	case 2:
		// The modulus is one for every argument, the fact most often confused
		// with the value itself.
		deg := []int{30, 45, 60, 90, 120, 180}[rand.Intn(6)]
		return complexChoice(
			fmt.Sprintf("What is \\(\\left|e^{i%d^\\circ}\\right|\\)?", deg),
			"1",
			"The exponential of a purely imaginary number has modulus one, because its real and imaginary parts are a cosine and a sine of the same argument.",
			"0", "2", "-1",
		)
	case 3:
		// Powers of i, asked as powers so the cycle has to be used.
		m := 2 + rand.Intn(4)
		answer := eulerUnits[(m-1)%4]
		return complexChoice(
			fmt.Sprintf("Use Euler's formula to evaluate \\(i^{%d}\\).", m),
			answer,
			fmt.Sprintf("\\(i=e^{i\\pi/2}\\), so \\(i^{%d}\\) turns the argument to %d times a quarter turn, where the unit circle reads %s.", m, m, answer),
			eulerOther(answer)...,
		)
	default:
		// A sum of conjugate exponentials, which is twice a cosine.
		deg := []int{60, 90, 120, 180}[rand.Intn(4)]
		answer := twiceCosine(deg)
		return complexChoice(
			fmt.Sprintf("What is \\(e^{i%d^\\circ}+e^{-i%d^\\circ}\\)?", deg, deg),
			answer,
			fmt.Sprintf("The two exponentials are conjugates, so their sum is twice the cosine of the angle. Twice the cosine of %d degrees is %s.", deg, answer),
			"1", "0", "-1", "-2",
		)
	}
}
