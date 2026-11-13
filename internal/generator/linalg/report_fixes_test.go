package linalg

import (
	"strings"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
)

// Report 461: the vectors were printed as "[-3 -3 4]", which reads as a matrix
// row and made the trace question look undefined. Print them with commas.
func TestOuterProductPrintsSeparatedVectors(t *testing.T) {
	g := &w3LinalgGen{kind: "outer_product"}
	for i := 0; i < 200; i++ {
		p := g.Generate(generator.GeneratorContext{KPOrdinal: 2})
		if strings.Contains(p.Question, "u=[") || strings.Contains(p.Question, "v=[") {
			t.Fatalf("vector is still printed as a bracket list: %q", p.Question)
		}
		if !strings.Contains(p.Question, "u=(") || !strings.Contains(p.Question, ", ") {
			t.Fatalf("vector is not comma-separated: %q", p.Question)
		}
	}
}
