package probability

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
)

// probVerdict is one statement in a count question. The answer is how many of
// the drawn statements are true, so the learner has to judge each claim rather
// than recognise a single phrase.
type probVerdict struct {
	label string
	yes   bool
	note  string
}

const probDraw = 8

func countProbVerdicts(lead string, pool []probVerdict) generator.Problem {
	var yes, no []probVerdict
	for _, v := range pool {
		if v.yes {
			yes = append(yes, v)
		} else {
			no = append(no, v)
		}
	}
	countChoice := rand.Intn(probDraw + 1)
	if countChoice > len(yes) {
		countChoice = len(yes)
	}
	if probDraw-countChoice > len(no) {
		countChoice = probDraw - len(no)
	}
	yesShuffled := append([]probVerdict(nil), yes...)
	noShuffled := append([]probVerdict(nil), no...)
	rand.Shuffle(len(yesShuffled), func(i, j int) { yesShuffled[i], yesShuffled[j] = yesShuffled[j], yesShuffled[i] })
	rand.Shuffle(len(noShuffled), func(i, j int) { noShuffled[i], noShuffled[j] = noShuffled[j], noShuffled[i] })

	var question string
	var chosen []probVerdict
	for attempt := 0; attempt < 40; attempt++ {
		chosen = append([]probVerdict(nil), yesShuffled[:countChoice]...)
		chosen = append(chosen, noShuffled[:probDraw-countChoice]...)
		rand.Shuffle(len(chosen), func(i, j int) { chosen[i], chosen[j] = chosen[j], chosen[i] })
		parts := make([]string, len(chosen))
		for i, v := range chosen {
			parts[i] = fmt.Sprintf("(%c) %s", 'a'+i, v.label)
		}
		question = fmt.Sprintf("%s %s (enter a number)", lead, strings.Join(parts, " "))
		if !strings.Contains(question, strconv.Itoa(countChoice)) {
			break
		}
	}

	notes := make([]string, len(chosen))
	for i, v := range chosen {
		verdict := "false"
		if v.yes {
			verdict = "true"
		}
		notes[i] = fmt.Sprintf("%s: %s, %s", v.label, verdict, v.note)
	}
	return generator.Problem{
		Question:    question,
		Answer:      strconv.Itoa(countChoice),
		Explanation: fmt.Sprintf("Exactly %d of the 8 statements are true. %s", countChoice, strings.Join(notes, " ")),
		Facts:       map[string]string{"count": strconv.Itoa(countChoice)},
	}
}

func probabilityAxiomsProblem() generator.Problem {
	pool := []probVerdict{
		{"P(A) >= 0 for every event A", true, "Non-negativity is the first axiom."},
		{"P(Omega) = 1", true, "The sample space has probability one."},
		{"P(A union B) = P(A) + P(B) for disjoint A and B", true, "Finite additivity follows from countable additivity."},
		{"P(empty set) = 0", true, "The empty event has no outcomes."},
		{"P(A complement) = 1 - P(A)", true, "An event and its complement partition Omega."},
		{"A subset B implies P(A) <= P(B)", true, "Monotonicity follows from additivity."},
		{"P(A union B) = P(A) + P(B) - P(A intersect B)", true, "This is inclusion-exclusion for two events."},
		{"Countable additivity over disjoint events is required", true, "It is the third Kolmogorov axiom."},
		{"P(A union B) = P(A) + P(B) for every pair of events", false, "The overlap is counted twice unless the events are disjoint."},
		{"P(Omega) = 0", false, "The whole sample space has probability one."},
		{"P(A complement) = P(A)", false, "They add to one, so they are equal only when both are 1/2."},
		{"Probabilities may exceed 1", false, "Every probability lies between 0 and 1."},
		{"Finite additivity alone defines a probability measure on every sigma-algebra", false, "Countable additivity is needed for infinite sample spaces."},
		{"P(empty set) = 1", false, "The empty event has probability zero."},
		{"P(A minus B) = P(A) + P(B)", false, "It is P(A) - P(A intersect B)."},
		{"Countable additivity applies to overlapping events without correction", false, "It applies to pairwise disjoint events."},
	}
	return countProbVerdicts("How many of these statements about probability axioms are true?", pool)
}

func sigmaAlgebraProblem() generator.Problem {
	pool := []probVerdict{
		{"A sigma-algebra contains Omega", true, "The whole space belongs to it."},
		{"A sigma-algebra is closed under complement", true, "Closure under complement is required."},
		{"A sigma-algebra is closed under countable unions", true, "Countable unions stay in the family."},
		{"A sigma-algebra is closed under countable intersections", true, "De Morgan follows from unions and complements."},
		{"A sigma-algebra contains the empty set", true, "It is the complement of Omega."},
		{"A sigma-algebra is closed under set difference", true, "A minus B is A intersect B complement."},
		{"The power set of Omega is a sigma-algebra", true, "It contains every subset and is closed under all operations."},
		{"The intersection of two sigma-algebras is a sigma-algebra", true, "Every axiom is inherited by the intersection."},
		{"Every subset of Omega must belong to a sigma-algebra", false, "Only the power set contains every subset."},
		{"The union of two sigma-algebras is always a sigma-algebra", false, "It need not be closed under unions."},
		{"A sigma-algebra is closed under unions but not complements", false, "Complement closure is part of the definition."},
		{"A sigma-algebra must be finite", false, "Borel sets form an uncountable sigma-algebra."},
		{"The intersection of sigma-algebras may fail to be a sigma-algebra", false, "The intersection always is one."},
		{"A sigma-algebra must contain every singleton", false, "Only when the singletons are measurable."},
		{"Complement closure is optional once countable unions are assumed", false, "Both are required."},
		{"Countable additivity alone forces a family to be a sigma-algebra", false, "Set closure axioms are separate."},
	}
	return countProbVerdicts("How many of these statements about sigma-algebras are true?", pool)
}

func continuityProblem() generator.Problem {
	pool := []probVerdict{
		{"A_n increasing to A implies P(A_n) tends to P(A)", true, "This is continuity from below."},
		{"A_n decreasing to A implies P(A_n) tends to P(A)", true, "This is continuity from above."},
		{"Continuity from below follows from countable additivity", true, "Write the limit event as a disjoint union."},
		{"Continuity from above needs P(A_1) finite", true, "Otherwise both sides can be infinite."},
		{"P(liminf A_n) <= liminf P(A_n)", true, "This is Fatou's lemma for probabilities."},
		{"P(union of A_n) = limit of P(union up to n)", true, "The partial unions increase to the full union."},
		{"Finite additivity plus continuity gives countable additivity", true, "The two formulations are equivalent."},
		{"P(intersection of A_n) = limit of P(intersection up to n)", true, "The partial intersections decrease to the full intersection."},
		{"Continuity from above holds without any finiteness hypothesis", false, "The sets need finite measure at some stage."},
		{"P(A_n) tending to P(A) forces A_n to increase to A", false, "The probabilities can converge without set convergence."},
		{"P(limsup A_n) = limsup P(A_n) always", false, "Only an inequality holds in general."},
		{"P(liminf A_n) >= liminf P(A_n) always", false, "Fatou gives the opposite inequality."},
		{"Continuity from below requires a finite measure space", false, "It holds for any measure."},
		{"P(union of A_n) = limsup P(A_n) always", false, "The union probability is the limit of partial unions."},
		{"Continuity of a probability measure implies independence", false, "Continuity concerns limits of events, not independence."},
		{"A decreasing sequence always has P(A_n) tending to P(empty set)", false, "It tends to the probability of the intersection."},
	}
	return countProbVerdicts("How many of these statements about continuity of probability are true?", pool)
}

func pairwiseIndependenceProblem() generator.Problem {
	pool := []probVerdict{
		{"Pairwise independence means P(A_i intersect A_j) = P(A_i)P(A_j) for every i < j", true, "Every pair must factor."},
		{"Pairwise independence does not imply mutual independence", true, "Triple intersections can still fail to factor."},
		{"Mutual independence implies pairwise independence", true, "Every subcollection must factor."},
		{"Mutual independence of three events includes the triple product", true, "All subcollection products are required."},
		{"Pairwise independent events can be jointly dependent", true, "The classic two-coin example shows this."},
		{"Complements of independent events remain independent", true, "Complementing preserves the product."},
		{"Independent events with positive probability have positive intersection", true, "The product of positive probabilities is positive."},
		{"Pairwise independence is a symmetric relation among events", true, "The condition treats every pair alike."},
		{"Pairwise independence implies mutual independence", false, "It does not control triple intersections."},
		{"Mutual independence only requires pairwise products", false, "Every subcollection must factor."},
		{"Pairwise independence requires the full intersection product", false, "That is mutual independence."},
		{"Independent events must be disjoint", false, "Disjoint events with positive probability are dependent."},
		{"Disjoint events with positive probability are independent", false, "Their intersection has probability zero while the product is positive."},
		{"Pairwise independence forces P(A intersect B intersect C) to factor", false, "That is an extra condition."},
		{"Taking complements destroys independence", false, "Complements stay independent."},
		{"Three pairwise independent events must have a triple product", false, "This is the standard gap between the notions."},
	}
	return countProbVerdicts("How many of these statements about pairwise independence are true?", pool)
}

func randomVariableProblem() generator.Problem {
	pool := []probVerdict{
		{"A random variable is a measurable function from Omega to R", true, "Measurability is with respect to the sigma-algebras."},
		{"The preimage of every Borel set must be measurable", true, "This is the measurability condition."},
		{"The indicator of an event is a random variable", true, "It takes values 0 and 1."},
		{"A continuous function of a random variable is a random variable", true, "Composition of measurable maps is measurable."},
		{"The sum of two random variables is a random variable", true, "Addition is measurable."},
		{"A constant function is a random variable", true, "Constants are measurable."},
		{"A random variable induces a distribution on R", true, "This is the pushforward measure."},
		{"Composing a random variable with a measurable map gives a random variable", true, "Measurability is preserved."},
		{"Any function from Omega to R is a random variable", false, "It must be measurable."},
		{"A random variable must be continuous", false, "It can be discrete, singular or mixed."},
		{"A random variable must be bounded", false, "Unbounded random variables are common."},
		{"A random variable maps R to Omega", false, "It maps the sample space to the real line."},
		{"The preimage of every subset of R must be measurable", false, "Only Borel sets are required."},
		{"The product of two random variables may fail to be a random variable", false, "The product is measurable."},
		{"An indicator function is not a random variable", false, "It is the simplest kind."},
		{"Every random variable has a density", false, "Only absolutely continuous ones do."},
	}
	return countProbVerdicts("How many of these statements about random variables are true?", pool)
}

func cdfProblem() generator.Problem {
	pool := []probVerdict{
		{"F(x) = P(X <= x) defines the cdf", true, "This is the standard definition."},
		{"A cdf is nondecreasing", true, "Larger thresholds include more outcomes."},
		{"A cdf is right-continuous", true, "This follows from continuity from above."},
		{"A cdf tends to 0 at negative infinity", true, "The event X <= x shrinks to the empty set."},
		{"A cdf tends to 1 at positive infinity", true, "The event X <= x grows to the whole space."},
		{"A cdf takes values in [0, 1]", true, "It is a probability."},
		{"P(a < X <= b) = F(b) - F(a)", true, "The interval probability is the increment."},
		{"Jumps of a cdf correspond to atoms", true, "The jump size is the point mass."},
		{"A cdf is left-continuous", false, "It is right-continuous."},
		{"F(x) = P(X < x) defines the cdf", false, "The standard cdf uses X <= x."},
		{"A cdf is nonincreasing", false, "It is nondecreasing."},
		{"The limit at negative infinity is 1", false, "It is 0."},
		{"P(a <= X <= b) = F(b) - F(a) in general", false, "The left endpoint mass is excluded from the increment."},
		{"A cdf must be continuous", false, "Discrete variables have jumps."},
		{"A cdf determines a density", false, "It determines the distribution, which need not have a density."},
		{"A cdf can decrease between two points", false, "Monotonicity forbids that."},
	}
	return countProbVerdicts("How many of these statements about cumulative distribution functions are true?", pool)
}

func memorylessProblem() generator.Problem {
	pool := []probVerdict{
		{"The geometric distribution is memoryless on the nonnegative integers", true, "It is the discrete memoryless law."},
		{"The exponential distribution is memoryless on [0, infinity)", true, "It is the continuous memoryless law."},
		{"P(X > s + t | X > s) = P(X > t) is the memoryless property", true, "The past does not change the remaining wait."},
		{"Every memoryless continuous distribution is exponential", true, "The functional equation forces the exponential form."},
		{"Every memoryless distribution on the nonnegative integers is geometric", true, "The discrete solution is geometric."},
		{"Exponential survival multiplies over disjoint waits", true, "P(X > s + t) = P(X > s)P(X > t)."},
		{"Geometric survival after n trials is (1 - p)^n", true, "Each trial must fail."},
		{"The memoryless property characterises these two families", true, "One discrete and one continuous."},
		{"The Poisson distribution is memoryless", false, "Poisson counts do not have this property."},
		{"The uniform distribution is memoryless", false, "A completed wait changes the remaining distribution."},
		{"The normal distribution is memoryless", false, "It has no such survival factorisation."},
		{"P(X > s + t) = P(X > s)P(X > t) for every distribution", false, "This characterises the memoryless families."},
		{"Memoryless continuous distributions include the gamma family", false, "Only the exponential case is memoryless."},
		{"The binomial distribution is memoryless", false, "It has a fixed number of trials."},
		{"The exponential distribution lacks the memoryless property", false, "It is the defining continuous example."},
		{"The geometric distribution lacks the memoryless property", false, "It is the defining discrete example."},
	}
	return countProbVerdicts("How many of these statements about the memoryless property are true?", pool)
}

func exponentialFamilyProblem() generator.Problem {
	pool := []probVerdict{
		{"An exponential family density has the form exp(eta T(x) - A(eta)) h(x)", true, "This is the canonical form."},
		{"eta is called the natural parameter", true, "It indexes the family."},
		{"T(x) is a sufficient statistic", true, "It captures the parameter information."},
		{"The Bernoulli distribution belongs to an exponential family", true, "Its log-odds is the natural parameter."},
		{"The Poisson distribution belongs to an exponential family", true, "The log mean is natural."},
		{"The normal distribution belongs to an exponential family", true, "Its mean and variance give two parameters."},
		{"The log-partition function is convex", true, "Its derivatives give moments."},
		{"Exponential family members have conjugate priors", true, "Conjugacy is a family property."},
		{"Every distribution is an exponential family", false, "Many distributions are not."},
		{"The beta distribution does not belong to an exponential family", false, "It does."},
		{"The gamma distribution does not belong to an exponential family", false, "It does."},
		{"The sufficient statistic is always x itself", false, "It can be a transform of x."},
		{"The log-partition function is concave", false, "It is convex."},
		{"The natural parameter is always the mean", false, "It is a transformed parameter."},
		{"Exponential family members cannot have conjugate priors", false, "They typically do."},
		{"A regular exponential family may have support depending on the parameter", false, "That violates regularity."},
	}
	return countProbVerdicts("How many of these statements about exponential families are true?", pool)
}
