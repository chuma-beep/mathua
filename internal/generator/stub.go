package generator

type Stub struct {
	ConceptID string
}

func (s *Stub) Generate(ctx GeneratorContext) Problem {
	return Problem{
		Question:    "This concept is not yet available for practice. Please check back later.",
		Answer:      "unavailable",
		Explanation: "A practice generator for " + s.ConceptID + " hasn't been built yet.",
	}
}
