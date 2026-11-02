package geometry

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/chuma-beep/mathua/internal/generator"
)

func intSqrt(n int) int {
	r := 0
	for r*r < n {
		r++
	}
	return r
}

type frustumGen struct{}

func (g *frustumGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := geoOrdinal(ctx)
	switch o {
	case 1:
		k := 1 + rand.Intn(8)
		h := 3 * k
		r := 1 + rand.Intn(9)
		R := r + 1 + rand.Intn(10-r)
		ans := k * (R*R + R*r + r*r)
		return generator.Problem{
			Question:    fmt.Sprintf("A frustum has height %d and base radii %d and %d. What is its volume divided by pi? (enter a number)", h, R, r),
			Answer:      strconv.Itoa(ans),
			Explanation: fmt.Sprintf("The volume is pi h (R^2 + Rr + r^2)/3. Here h/3 is %d and the bracket is %d, so the quotient by pi is %d.", k, R*R+R*r+r*r, ans),
			KPOrdinal:   o,
			Facts:       map[string]string{"answer": strconv.Itoa(ans)},
		}
	case 2:
		triples := [][2]int{{3, 4}, {6, 8}, {5, 12}, {8, 15}, {7, 24}, {9, 12}}
		t := triples[rand.Intn(len(triples))]
		l := t[0]*t[0] + t[1]*t[1]
		return generator.Problem{
			Question:    fmt.Sprintf("A frustum has height %d and the difference of its radii is %d. What is its slant height? (enter a number)", t[1], t[0]),
			Answer:      strconv.Itoa(intSqrt(l)),
			Explanation: fmt.Sprintf("The slant height is the hypotenuse of a right triangle with legs %d and %d, so it is %d.", t[0], t[1], intSqrt(l)),
			KPOrdinal:   o,
			Facts:       map[string]string{"answer": strconv.Itoa(intSqrt(l))},
		}
	default:
		triples := [][2]int{{3, 4}, {6, 8}, {5, 12}, {8, 15}, {7, 24}, {9, 12}}
		t := triples[rand.Intn(len(triples))]
		l := intSqrt(t[0]*t[0] + t[1]*t[1])
		r := 1 + rand.Intn(5)
		R := r + t[0]
		ans := (R + r) * l
		return generator.Problem{
			Question:    fmt.Sprintf("A frustum has base radii %d and %d and slant height %d. What is its lateral surface area divided by pi? (enter a number)", R, r, l),
			Answer:      strconv.Itoa(ans),
			Explanation: fmt.Sprintf("The lateral area is pi (R+r)l, and (%d + %d) times %d equals %d.", R, r, l, ans),
			KPOrdinal:   o,
			Facts:       map[string]string{"answer": strconv.Itoa(ans)},
		}
	}
}
