package machinelearning

import "github.com/chuma-beep/mathua/internal/generator"

func Register(reg *generator.Registry) {
	reg.Register("ml.backpropagation", &generator.Stub{ConceptID: "ml.backpropagation"})
}
