package mathutil

import "testing"

func TestMatrix2x2InverseReportedEntry(t *testing.T) {
	m := Matrix2x2{A: -3, B: -1, C: 1, D: 0}
	if got := m.Det(); got != 1 {
		t.Fatalf("Det() = %d, want 1", got)
	}
	inv, ok := m.Inverse()
	if !ok {
		t.Fatal("Inverse() reported a singular matrix")
	}
	if inv.C != -1 {
		t.Fatalf("inverse (2,1) = %d, want -1 for %+v", inv.C, m)
	}
}
