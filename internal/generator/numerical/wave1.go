package numerical

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/chuma-beep/mathua/internal/generator"
)

// ---------------------------------------------------------------------------
// Backward error
// ---------------------------------------------------------------------------

type backwardErrorGen struct{}

func (g *backwardErrorGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := ctx.KPOrdinal
	if kp < 1 || kp > 3 {
		kp = rand.Intn(3) + 1
	}

	var question, explanation string
	var answer int
	switch kp {
	case 1:
		residual := rand.Intn(15) + 1
		answer = residual
		question = fmt.Sprintf("A computed root x̂ gives the residual f(x̂) = %d. What is the absolute backward error measured by that residual? (enter a number)", residual)
		explanation = fmt.Sprintf("The absolute backward error is the size of the residual, so it is %d. It measures how far the data had to move to make x̂ an exact root.", residual)

	case 2:
		derivative := rand.Intn(8) + 2
		forward := rand.Intn(9) + 1
		residual := derivative * forward
		answer = forward
		question = fmt.Sprintf("Near a computed root, the residual is %d and |f'(x̂)| = %d. Using the linear estimate, what is the forward error |x - x̂|? (enter a number)", residual, derivative)
		explanation = fmt.Sprintf("The linear estimate gives forward error ≈ residual / |f'(x̂)| = %d / %d = %d.", residual, derivative, forward)

	default:
		backward := rand.Intn(5) + 1 // percent
		condition := []int{100, 200, 400, 500, 1000}[rand.Intn(5)]
		answer = backward * condition / 100
		question = fmt.Sprintf("A problem has condition number %d. The relative backward error is %d%%. What relative forward error does the condition number predict, in percent? (enter a number)", condition, backward)
		explanation = fmt.Sprintf("Forward error ≈ condition number × backward error = %d × %d%% = %d%%.", condition, backward, answer)
	}

	return generator.Problem{
		Question:    question,
		Answer:      strconv.Itoa(answer),
		Explanation: explanation,
		Facts:       map[string]string{"answer": strconv.Itoa(answer)},
		KPOrdinal:   kp,
	}
}

// ---------------------------------------------------------------------------
// Floating-point error
// ---------------------------------------------------------------------------

type floatingErrorGen struct{}

func (g *floatingErrorGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := ctx.KPOrdinal
	if kp < 1 || kp > 3 {
		kp = rand.Intn(3) + 1
	}

	var question, explanation string
	var answer int
	switch kp {
	case 1:
		trueValue := rand.Intn(80) + 20
		err := rand.Intn(6) + 1
		answer = err
		question = fmt.Sprintf("The value %d is approximated by %d. What is the absolute error? (enter a number)", trueValue, trueValue+err)
		explanation = fmt.Sprintf("Absolute error is |%d - %d| = %d.", trueValue+err, trueValue, err)

	case 2:
		percent := []int{1, 2, 4, 5, 10}[rand.Intn(5)]
		trueValue := 100 / percent * (rand.Intn(8) + 2)
		err := trueValue * percent / 100
		answer = percent
		question = fmt.Sprintf("The value %d is approximated by %d. What is the relative error, in percent? (enter a number)", trueValue, trueValue+err)
		explanation = fmt.Sprintf("Relative error is %d/%d = %d%%.", err, trueValue, percent)

	default:
		digits := 4 + rand.Intn(2) // 4..5 digits
		base := rand.Intn(9) + 1
		value := base
		for i := 1; i < digits; i++ {
			value = value*10 + rand.Intn(10)
		}
		sig := rand.Intn(3) + 2 // 2..4
		scale := pow10(digits - sig)
		answer = ((value + scale/2) / scale) * scale
		question = fmt.Sprintf("Round %d to %d significant figures. (enter a number)", value, sig)
		explanation = fmt.Sprintf("Keep the first %d significant digits and round the next digit, giving %d.", sig, answer)
	}

	return generator.Problem{
		Question:    question,
		Answer:      strconv.Itoa(answer),
		Explanation: explanation,
		Facts:       map[string]string{"answer": strconv.Itoa(answer)},
		KPOrdinal:   kp,
	}
}

func pow10(n int) int {
	out := 1
	for i := 0; i < n; i++ {
		out *= 10
	}
	return out
}

// ---------------------------------------------------------------------------
// Brent's method
// ---------------------------------------------------------------------------

type brentGen struct{}

func (g *brentGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := ctx.KPOrdinal
	if kp < 1 || kp > 3 {
		kp = rand.Intn(3) + 1
	}

	var question, explanation string
	var answer int
	switch kp {
	case 1:
		steps := rand.Intn(5) + 1
		width := (rand.Intn(20) + 1) << steps
		answer = width >> steps
		noun := "steps"
		if steps == 1 {
			noun = "step"
		}
		question = fmt.Sprintf("Brent's method begins with a bracket of width %d. If it takes %d bisection %s, what is the bracket width afterwards? (enter a number)", width, steps, noun)
		explanation = fmt.Sprintf("Each bisection halves the bracket, so after %d steps the width is %d / 2^%d = %d.", steps, width, steps, answer)

	case 2:
		steps := rand.Intn(7) + 1
		width := (rand.Intn(10) + 1) << steps
		target := width >> steps
		answer = steps
		question = fmt.Sprintf("A bracket has width %d. How many bisection steps are needed to make the width at most %d? (enter a number)", width, target)
		explanation = fmt.Sprintf("Each step halves the width. After %d steps the width is %d / 2^%d = %d, which is at most the target %d.", steps, width, steps, target, target)

	default:
		steps := rand.Intn(7) + 1
		answer = 1 << steps
		question = fmt.Sprintf("Brent's bisection fallback halves the bracket at every step. By what factor has the bracket width shrunk after %d steps? (enter a number)", steps)
		explanation = fmt.Sprintf("Each bisection halves the width, so after %d steps it has shrunk by a factor of 2^%d = %d.", steps, steps, answer)
	}

	return generator.Problem{
		Question:    question,
		Answer:      strconv.Itoa(answer),
		Explanation: explanation,
		Facts:       map[string]string{"answer": strconv.Itoa(answer)},
		KPOrdinal:   kp,
	}
}
