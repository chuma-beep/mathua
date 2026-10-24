package linalg

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
)

type kernelGen struct{}

// kernelGen builds a matrix whose rank is known by construction: a basis of
// independent rows plus extra rows that are integer combinations of them. The
// nullity is then n minus that rank, and the printed matrix makes the rank
// something the learner can find by elimination.
func (g *kernelGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	n := 2 + rand.Intn(2)
	rank := 1 + rand.Intn(n-1)

	// Each basis row owns a pivot in the first rank columns and may vary only in
	// the free columns, so the rows are independent by construction.
	basis := make([][]int, rank)
	for i := range basis {
		basis[i] = make([]int, n)
		for j := 0; j < n; j++ {
			switch {
			case j == i:
				basis[i][j] = 1
			case j < rank:
				basis[i][j] = 0
			default:
				basis[i][j] = rand.Intn(3) - 1
			}
		}
	}

	rows := append([][]int{}, basis...)
	for len(rows) < n {
		combo := make([]int, n)
		for _, b := range basis {
			mult := rand.Intn(3) - 1
			if mult == 0 {
				mult = 1
			}
			for j := 0; j < n; j++ {
				combo[j] += mult * b[j]
			}
		}
		rows = append(rows, combo)
	}
	for i := len(rows) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		rows[i], rows[j] = rows[j], rows[i]
	}

	rowText := make([]string, len(rows))
	rowFacts := make([]string, len(rows))
	for i, r := range rows {
		cells := make([]string, len(r))
		for j, v := range r {
			cells[j] = strconv.Itoa(v)
		}
		rowText[i] = "[" + strings.Join(cells, ", ") + "]"
		rowFacts[i] = strings.Join(cells, ",")
	}

	return generator.Problem{
		Question: fmt.Sprintf("The linear map T(x) = Ax has matrix A = %s. What is the dimension of the kernel of T? (enter a number)", strings.Join(rowText, ", ")),
		Answer:   strconv.Itoa(n - rank),
		Explanation: fmt.Sprintf(
			"Row reduction leaves %d independent rows, so the rank is %d. Rank-nullity gives nullity %d - %d = %d.",
			rank, rank, n, rank, n-rank),
		Facts: map[string]string{"rows": strings.Join(rowFacts, ";"), "n": strconv.Itoa(n), "answer": strconv.Itoa(n - rank)},
	}
}
