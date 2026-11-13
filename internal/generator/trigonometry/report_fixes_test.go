package trigonometry

import (
	"regexp"
	"strconv"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
)

// Report 456: a triangle whose sides meet at 180 degrees has zero height, but
// the generator returned r^2/2 for every angle in its list.
func TestSegmentAreaTriangleRespectsAngle(t *testing.T) {
	g := &trigBridgeGen{kind: "segment_area"}
	re := regexp.MustCompile(`two sides of length (\d+) meeting at (\d+) degrees`)
	for i := 0; i < 400; i++ {
		p := g.Generate(generator.GeneratorContext{KPOrdinal: 2})
		m := re.FindStringSubmatch(p.Question)
		if m == nil {
			t.Fatalf("unexpected segment area question: %q", p.Question)
		}
		r, _ := strconv.Atoi(m[1])
		angle, _ := strconv.Atoi(m[2])
		want := r * r / 2
		if angle == 180 {
			want = 0
		}
		if p.Answer != strconv.Itoa(want) {
			t.Fatalf("%q answered %s, want %d", p.Question, p.Answer, want)
		}
	}
}
