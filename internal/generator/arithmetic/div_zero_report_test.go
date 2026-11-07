package arithmetic

import (
	"strconv"
	"testing"
)

func TestDivisionByZeroProblemHasAnswer(t *testing.T) {
	for i := 0; i < 20; i++ {
		p := divisionByZeroProblem()
		if p.Question == "" || p.Answer == "" || p.Explanation == "" {
			t.Fatalf("incomplete division-by-zero problem: %+v", p)
		}
		if _, err := strconv.Atoi(p.Answer); err != nil {
			t.Fatalf("answer %q is not an integer: %v", p.Answer, err)
		}
	}
}
