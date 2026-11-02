package statistics

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/chuma-beep/mathua/internal/generator"
)

// The descriptive and regression concepts added in the bridge pilot change
// question shape at each knowledge point, so every generator reads KPOrdinal:
// the first point asks for the definition or a direct computation, the second
// applies the formula, and the third asks the learner to interpret or choose.

type stdDevGen struct{}

func (g *stdDevGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := bridgeOrdinal(ctx)
	d := 1 + rand.Intn(5)
	k := 5 + rand.Intn(15)
	switch o {
	case 1:
		return generator.Problem{
			Question: fmt.Sprintf("A data set is %d, %d, %d, %d. What is its population standard deviation? (enter a number)", k-d, k-d, k+d, k+d),
			Answer:   strconv.Itoa(d),
			Explanation: fmt.Sprintf("The mean is %d. Each value is %d away from the mean, so the population variance is %d and the standard deviation is %d.",
				k, d, d*d, d),
			Facts: map[string]string{"kind": "population_sd", "answer": strconv.Itoa(d)},
		}
	case 2:
		return generator.Problem{
			Question: fmt.Sprintf("A data set is %d, %d, %d, %d, %d. What is its population variance? (enter a number)", k-2*d, k-d, k, k+d, k+2*d),
			Answer:   strconv.Itoa(2 * d * d),
			Explanation: fmt.Sprintf("The mean is %d. The squared deviations are %d, %d, 0, %d and %d, whose average is %d.",
				k, 4*d*d, d*d, d*d, 4*d*d, 2*d*d),
			Facts: map[string]string{"kind": "population_variance", "answer": strconv.Itoa(2 * d * d)},
		}
	default:
		d2 := d + 1 + rand.Intn(3)
		return generator.Problem{
			Question: fmt.Sprintf("Data set A is %d, %d, %d, %d. Data set B is %d, %d, %d, %d. Which set has the larger standard deviation? (enter 1 for A or 2 for B)",
				k-d, k-d, k+d, k+d, k-d2, k-d2, k+d2, k+d2),
			Answer: "2",
			Explanation: fmt.Sprintf("Both sets are centred at %d. Set A has spread %d and set B has spread %d, so B has the larger standard deviation.",
				k, d, d2),
			Facts: map[string]string{"kind": "compare", "answer": "2"},
		}
	}
}

type boxPlotGen struct{}

func (g *boxPlotGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := bridgeOrdinal(ctx)
	a := 2 + rand.Intn(6)
	step := 2 + rand.Intn(5)
	switch o {
	case 1:
		vals := []int{a, a + step, a + 2*step, a + 3*step, a + 4*step}
		return generator.Problem{
			Question:    fmt.Sprintf("A sorted data set is %d, %d, %d, %d, %d. What is the median? (enter a number)", vals[0], vals[1], vals[2], vals[3], vals[4]),
			Answer:      strconv.Itoa(vals[2]),
			Explanation: fmt.Sprintf("With five ordered values the median is the third value, %d.", vals[2]),
			Facts:       map[string]string{"kind": "median", "answer": strconv.Itoa(vals[2])},
		}
	case 2:
		vals := []int{a, a + step, a + 2*step, a + 3*step, a + 4*step, a + 5*step, a + 6*step}
		q1 := vals[1]
		q3 := vals[5]
		return generator.Problem{
			Question: fmt.Sprintf("A sorted data set is %d, %d, %d, %d, %d, %d, %d. What is the interquartile range? (enter a number)",
				vals[0], vals[1], vals[2], vals[3], vals[4], vals[5], vals[6]),
			Answer:      strconv.Itoa(q3 - q1),
			Explanation: fmt.Sprintf("The lower quartile is %d and the upper quartile is %d, so the interquartile range is %d.", q1, q3, q3-q1),
			Facts:       map[string]string{"kind": "iqr", "answer": strconv.Itoa(q3 - q1)},
		}
	default:
		vals := []int{a, a + step, a + 2*step, a + 3*step, a + 4*step, a + 10*step, a + 11*step}
		q1 := vals[1]
		q3 := vals[5]
		iqr := q3 - q1
		low := q1 - 3*iqr/2
		high := q3 + 3*iqr/2
		outliers := 0
		for _, v := range vals {
			if v < low || v > high {
				outliers++
			}
		}
		return generator.Problem{
			Question: fmt.Sprintf("A sorted data set is %d, %d, %d, %d, %d, %d, %d. Using the 1.5 IQR rule, how many values are outliers? (enter a number)",
				vals[0], vals[1], vals[2], vals[3], vals[4], vals[5], vals[6]),
			Answer: strconv.Itoa(outliers),
			Explanation: fmt.Sprintf("The quartiles are %d and %d, so IQR = %d and the fences are %d and %d. Exactly %d values lie outside the fences.",
				q1, q3, iqr, low, high, outliers),
			Facts: map[string]string{"kind": "outliers", "answer": strconv.Itoa(outliers)},
		}
	}
}

type correlationGen struct{}

func (g *correlationGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := bridgeOrdinal(ctx)
	switch o {
	case 1:
		type scenario struct {
			text string
			sign string
			word string
		}
		scenarios := []scenario{
			{"Outside temperature and ice-cream sales tend to rise together.", "1", "positive"},
			{"A child's age and height tend to rise together.", "1", "positive"},
			{"Hours practised and the number of errors tend to move in opposite directions.", "-1", "negative"},
			{"The price of a product and the quantity demanded tend to move in opposite directions.", "-1", "negative"},
			{"Shoe size and score on a vocabulary test show no straight-line pattern.", "0", "no linear"},
			{"Height and the first digit of a phone number show no straight-line pattern.", "0", "no linear"},
			{"Hours studied and exam score tend to rise together.", "1", "positive"},
			{"Speed and travel time for a fixed distance tend to move in opposite directions.", "-1", "negative"},
		}
		s := scenarios[rand.Intn(len(scenarios))]
		return generator.Problem{
			Question:    fmt.Sprintf("%s What is the sign of the correlation coefficient? (enter 1 for positive, -1 for negative, or 0 for no linear association)", s.text),
			Answer:      s.sign,
			Explanation: fmt.Sprintf("This is a %s association, so the correlation sign is %s.", s.word, s.sign),
			Facts:       map[string]string{"kind": "sign", "answer": s.sign},
		}
	case 2:
		sign := 1
		if rand.Intn(2) == 0 {
			sign = -1
		}
		return generator.Problem{
			Question:    fmt.Sprintf("Every point lies exactly on the line y = %dx + %d. What is the correlation coefficient? (enter a number)", sign*(1+rand.Intn(4)), 1+rand.Intn(9)),
			Answer:      strconv.Itoa(sign),
			Explanation: fmt.Sprintf("A perfect straight-line fit gives correlation %d, with the sign of the slope.", sign),
			Facts:       map[string]string{"kind": "perfect", "answer": strconv.Itoa(sign)},
		}
	default:
		type triple struct {
			r      string
			sx, sy int
			cov    int
		}
		options := []triple{
			{"0.5", 4, 5, 10},
			{"0.5", 6, 10, 30},
			{"0.5", 8, 15, 60},
			{"-0.5", 4, 5, -10},
			{"-0.5", 6, 10, -30},
			{"-0.5", 8, 15, -60},
			{"0.8", 5, 5, 20},
			{"0.8", 10, 5, 40},
			{"0.8", 5, 10, 40},
			{"-0.8", 5, 5, -20},
			{"-0.8", 10, 5, -40},
			{"-0.8", 5, 10, -40},
		}
		t := options[rand.Intn(len(options))]
		return generator.Problem{
			Question:    fmt.Sprintf("A data set has correlation r = %s, standard deviation of x equal to %d and standard deviation of y equal to %d. What is the covariance? (enter a number)", t.r, t.sx, t.sy),
			Answer:      strconv.Itoa(t.cov),
			Explanation: fmt.Sprintf("Covariance is r times sx times sy: %s x %d x %d = %d.", t.r, t.sx, t.sy, t.cov),
			Facts:       map[string]string{"kind": "covariance", "answer": strconv.Itoa(t.cov)},
		}
	}
}

type leastSquaresGen struct{}

func (g *leastSquaresGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := bridgeOrdinal(ctx)
	x1 := 1 + rand.Intn(6)
	x2 := x1 + 2 + rand.Intn(4)
	m := 1 + rand.Intn(4)
	b := 1 + rand.Intn(10)
	y1 := m*x1 + b
	y2 := m*x2 + b
	switch o {
	case 1:
		return generator.Problem{
			Question:    fmt.Sprintf("The least-squares line passes through (%d, %d) and (%d, %d). What is its slope? (enter a number)", x1, y1, x2, y2),
			Answer:      strconv.Itoa(m),
			Explanation: fmt.Sprintf("Slope = (%d - %d) / (%d - %d) = %d.", y2, y1, x2, x1, m),
			Facts:       map[string]string{"kind": "slope", "answer": strconv.Itoa(m)},
		}
	case 2:
		return generator.Problem{
			Question:    fmt.Sprintf("A least-squares line has slope %d and passes through (%d, %d). What is its y-intercept? (enter a number)", m, x1, y1),
			Answer:      strconv.Itoa(b),
			Explanation: fmt.Sprintf("Using y = mx + b: %d = %d(%d) + b, so b = %d.", y1, m, x1, b),
			Facts:       map[string]string{"kind": "intercept", "answer": strconv.Itoa(b)},
		}
	default:
		x0 := x2 + 1 + rand.Intn(4)
		y0 := m*x0 + b
		return generator.Problem{
			Question:    fmt.Sprintf("A fitted line is y = %dx + %d. What is the predicted y when x = %d? (enter a number)", m, b, x0),
			Answer:      strconv.Itoa(y0),
			Explanation: fmt.Sprintf("Substitute x = %d: y = %d(%d) + %d = %d.", x0, m, x0, b, y0),
			Facts:       map[string]string{"kind": "predict", "answer": strconv.Itoa(y0)},
		}
	}
}

type chiSquareIndependenceGen struct{}

func (g *chiSquareIndependenceGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := bridgeOrdinal(ctx)
	switch o {
	case 1:
		r := 2 + rand.Intn(3)
		c := 2 + rand.Intn(3)
		return generator.Problem{
			Question:    fmt.Sprintf("A contingency table has %d rows and %d columns. How many degrees of freedom does the chi-square test of independence use? (enter a number)", r, c),
			Answer:      strconv.Itoa((r - 1) * (c - 1)),
			Explanation: fmt.Sprintf("Degrees of freedom = (r - 1)(c - 1) = (%d - 1)(%d - 1) = %d.", r, c, (r-1)*(c-1)),
			Facts:       map[string]string{"kind": "df", "answer": strconv.Itoa((r - 1) * (c - 1))},
		}
	case 2:
		n := 20 * (2 + rand.Intn(4))
		row := 4 * (2 + rand.Intn(5))
		col := 5 * (2 + rand.Intn(5))
		expected := row * col / n
		return generator.Problem{
			Question:    fmt.Sprintf("A table has total count %d. One row totals %d and one column totals %d. What is the expected count in that cell? (enter a number)", n, row, col),
			Answer:      strconv.Itoa(expected),
			Explanation: fmt.Sprintf("Expected count = (row total x column total) / total = (%d x %d) / %d = %d.", row, col, n, expected),
			Facts:       map[string]string{"kind": "expected", "answer": strconv.Itoa(expected)},
		}
	default:
		e := 2 + rand.Intn(6)
		k := 1 + rand.Intn(3)
		observed := e + k*e
		contrib := k * k * e
		return generator.Problem{
			Question:    fmt.Sprintf("A cell has observed count %d and expected count %d. What is its chi-square contribution? (enter a number)", observed, e),
			Answer:      strconv.Itoa(contrib),
			Explanation: fmt.Sprintf("Contribution = (O - E)^2 / E = (%d - %d)^2 / %d = %d.", observed, e, e, contrib),
			Facts:       map[string]string{"kind": "contribution", "answer": strconv.Itoa(contrib)},
		}
	}
}

type studyDesignGen struct{}

func (g *studyDesignGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := bridgeOrdinal(ctx)
	switch o {
	case 1:
		groups := 2 + rand.Intn(4)
		per := 5 + rand.Intn(20)
		return generator.Problem{
			Question:    fmt.Sprintf("An experiment has %d treatment groups with %d participants in each group. How many participants are in the study? (enter a number)", groups, per),
			Answer:      strconv.Itoa(groups * per),
			Explanation: fmt.Sprintf("Total participants = %d groups x %d per group = %d.", groups, per, groups*per),
			Facts:       map[string]string{"kind": "total", "answer": strconv.Itoa(groups * per)},
		}
	case 2:
		groups := 2 + rand.Intn(4)
		blocks := 2 + rand.Intn(4)
		return generator.Problem{
			Question:    fmt.Sprintf("A randomized block design has %d treatments and %d blocks. How many treatment-block cells are there? (enter a number)", groups, blocks),
			Answer:      strconv.Itoa(groups * blocks),
			Explanation: fmt.Sprintf("Each treatment appears once in each block, so there are %d x %d = %d cells.", groups, blocks, groups*blocks),
			Facts:       map[string]string{"kind": "cells", "answer": strconv.Itoa(groups * blocks)},
		}
	default:
		stratum := 40 * (2 + rand.Intn(4))
		population := 4 * stratum
		sample := 4 * (5 + rand.Intn(10))
		alloc := sample / 4
		return generator.Problem{
			Question:    fmt.Sprintf("A population of %d is split into four equal strata. A sample of %d is allocated proportionally. How many people are sampled from one stratum? (enter a number)", population, sample),
			Answer:      strconv.Itoa(alloc),
			Explanation: fmt.Sprintf("Each stratum is one quarter of the population, so it receives one quarter of the sample: %d / 4 = %d.", sample, alloc),
			Facts:       map[string]string{"kind": "allocation", "answer": strconv.Itoa(alloc)},
		}
	}
}

type samplingMethodsGen struct{}

func (g *samplingMethodsGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := bridgeOrdinal(ctx)
	switch o {
	case 1:
		n := 100 * (2 + rand.Intn(4))
		k := 5 * (2 + rand.Intn(4))
		return generator.Problem{
			Question:    fmt.Sprintf("A list of %d people is sampled systematically by taking every %dth person. How many people are selected? (enter a number)", n, k),
			Answer:      strconv.Itoa(n / k),
			Explanation: fmt.Sprintf("The sampling interval is %d, so the number selected is %d / %d = %d.", k, n, k, n/k),
			Facts:       map[string]string{"kind": "systematic", "answer": strconv.Itoa(n / k)},
		}
	case 2:
		n := 100 * (2 + rand.Intn(4))
		k := 5 * (2 + rand.Intn(4))
		return generator.Problem{
			Question:    fmt.Sprintf("A population of %d is divided into clusters of %d. A cluster sample selects 3 clusters. How many people are selected? (enter a number)", n, k),
			Answer:      strconv.Itoa(3 * k),
			Explanation: fmt.Sprintf("Each selected cluster contributes %d people, so 3 clusters contribute 3 x %d = %d.", k, k, 3*k),
			Facts:       map[string]string{"kind": "cluster", "answer": strconv.Itoa(3 * k)},
		}
	default:
		total := 200 * (2 + rand.Intn(4))
		part := total / 4
		sample := 40 * (2 + rand.Intn(4))
		alloc := part * sample / total
		return generator.Problem{
			Question:    fmt.Sprintf("A population of %d has one stratum of %d. A proportional sample of %d is drawn. How many come from that stratum? (enter a number)", total, part, sample),
			Answer:      strconv.Itoa(alloc),
			Explanation: fmt.Sprintf("The stratum is one quarter of the population, so it supplies one quarter of the sample: %d / 4 = %d.", sample, alloc),
			Facts:       map[string]string{"kind": "stratified", "answer": strconv.Itoa(alloc)},
		}
	}
}

func bridgeOrdinal(ctx generator.GeneratorContext) int {
	if ctx.KPOrdinal >= 1 && ctx.KPOrdinal <= 3 {
		return ctx.KPOrdinal
	}
	return 1 + rand.Intn(3)
}
