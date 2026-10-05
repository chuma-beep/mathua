// Package all registers every generator onto one registry.
//
// This exists so that the set of generators a learner can meet and the set a diagnostic
// measures are the same list. `cmd/mathua/main.go` used to wire the seventeen
// sub-packages inline, which meant nothing outside `main` could enumerate "every concept
// that has a generator" — and that is precisely the question you need answered to find
// out whether a concept can supply enough distinct questions to keep a practice feed
// going. A diagnostic built from its own hand-written list would drift silently from
// production the first time a domain was added.
//
// It lives outside `internal/generator` because every sub-package imports
// `internal/generator` for `Problem` and `GeneratorContext`, so a registrar inside it
// would be an import cycle.
package all

import (
	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/abstract"
	"github.com/chuma-beep/mathua/internal/generator/algebra"
	"github.com/chuma-beep/mathua/internal/generator/arithmetic"
	"github.com/chuma-beep/mathua/internal/generator/calculus"
	"github.com/chuma-beep/mathua/internal/generator/complex"
	"github.com/chuma-beep/mathua/internal/generator/discrete"
	"github.com/chuma-beep/mathua/internal/generator/fractions"
	"github.com/chuma-beep/mathua/internal/generator/geometry"
	"github.com/chuma-beep/mathua/internal/generator/linalg"
	"github.com/chuma-beep/mathua/internal/generator/machinelearning"
	"github.com/chuma-beep/mathua/internal/generator/numtheory"
	"github.com/chuma-beep/mathua/internal/generator/odes"
	"github.com/chuma-beep/mathua/internal/generator/prealgebra"
	"github.com/chuma-beep/mathua/internal/generator/precalculus"
	"github.com/chuma-beep/mathua/internal/generator/statistics"
	"github.com/chuma-beep/mathua/internal/generator/topology"
	"github.com/chuma-beep/mathua/internal/generator/trigonometry"
)

// Register adds every generator in the corpus to reg.
//
// Kept in the same order as the registrations this replaced, so a diff against main.go
// shows only the move.
func Register(reg *generator.Registry) {
	arithmetic.Register(reg)
	fractions.Register(reg)
	geometry.Register(reg)
	prealgebra.Register(reg)
	algebra.Register(reg)
	trigonometry.Register(reg)
	precalculus.Register(reg)
	statistics.Register(reg)
	numtheory.Register(reg)
	complex.Register(reg)
	linalg.Register(reg)
	machinelearning.Register(reg)
	discrete.Register(reg)
	calculus.Register(reg)
	odes.Register(reg)
	abstract.Register(reg)
	topology.Register(reg)
}

// Registry returns a registry with every generator registered.
func Registry() *generator.Registry {
	reg := generator.NewRegistry()
	Register(reg)
	return reg
}
