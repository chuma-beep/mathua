package generator

type Problem struct {
	Question    string
	Answer      string
	Explanation string
}

type Generator interface {
	Generate(difficulty float64) Problem
}
