#!/usr/bin/env bash
set -euo pipefail

if [ $# -lt 1 ]; then
  echo "Usage: $0 <concept-id> [package-name]"
  echo ""
  echo "Scaffolds a new problem generator for a concept."
  echo "  concept-id    e.g. 'alg.linear.slope'"
  echo "  package-name  Go package name (default: derived from domain)"
  echo ""
  echo "Example: $0 mydomain.myconcept"
  exit 1
fi

CONCEPT_ID="$1"
PACKAGE="${2:-$(echo "$CONCEPT_ID" | cut -d. -f1)}"

DIR="internal/generator/$PACKAGE"
FILE="$DIR/generators.go"
TEST="$DIR/generators_test.go"

mkdir -p "$DIR"

# Convert concept ID to a valid Go type name
STRUCT=$(echo "$CONCEPT_ID" | tr '.' '_' | sed 's/_\([a-z]\)/\U\1/g' | sed 's/^\([a-z]\)/\U\1/')Gen

if [ -f "$FILE" ]; then
  echo "Appending to existing $FILE"
else
  cat > "$FILE" <<GOEOF
package $PACKAGE

import (
	"fmt"
	"math/rand"

	"github.com/chuma-beep/mathua/internal/generator"
)

func Register(reg *generator.Registry) {
	reg.Register("$CONCEPT_ID", &$STRUCT{})
}

type $STRUCT struct{}

func (g *$STRUCT) Generate(difficulty float64) generator.Problem {
	// TODO: use difficulty to scale problem complexity (0.0=easy, 1.0=hard)
	a := rand.Intn(int(5+difficulty*10)) + 1
	b := rand.Intn(int(5+difficulty*10)) + 1
	return generator.Problem{
		Question:    fmt.Sprintf("Example question for %s using a=%d, b=%d.", "$CONCEPT_ID", a, b),
		Answer:      fmt.Sprintf("%d", a+b),
		Explanation: fmt.Sprintf("%d + %d = %d.", a, b, a+b),
	}
}
GOEOF
  echo "Created $FILE"
fi

if [ ! -f "$TEST" ]; then
  cat > "$TEST" <<GOEOF
package $PACKAGE

import (
	"math/rand"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
)

func fuzzGen(t *testing.T, gen generator.Generator) {
	t.Helper()
	for i := 0; i < 50; i++ {
		d := rand.Float64()
		p := gen.Generate(d)
		if p.Question == "" || p.Answer == "" || p.Explanation == "" {
			t.Errorf("empty field at difficulty=%.2f", d)
		}
	}
}

func Test$(echo "$STRUCT" | sed 's/Gen$//')(t *testing.T) { fuzzGen(t, &$STRUCT{}) }
GOEOF
  echo "Created $TEST"
fi

# Register the package in the root generator init if main.go pattern exists
ROOT_REG="internal/generator/import.go"
if [ -f "$ROOT_REG" ]; then
  echo "NOTE: You may need to add an import and register call in $ROOT_REG"
else
  echo "NOTE: Register the generator package in your main.go or equivalent init function:"
  echo "  import \"github.com/chuma-beep/mathua/internal/generator/$PACKAGE\""
  echo "  $PACKAGE.Register(reg)"
fi

echo ""
echo "Done! Scaffolded generator for '$CONCEPT_ID' in $DIR"
echo "Next: implement the Generate method logic in $FILE"
