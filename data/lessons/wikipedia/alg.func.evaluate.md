> Content sourced from [Wikipedia](https://en.wikipedia.org/wiki/Function_%28mathematics%29) — CC BY-SA 4.0

# Evaluate a function at a given value

### Partial functions


Partial functions are defined similarly to ordinary functions, with the "total" condition removed. That is, a *partial function* from  to  is a binary relation  between  and  such that, for every \(x\in X,\) there is *at most one*  in  such that \((x,y) \in R.\)

Using functional notation, this means that, given \(x\in X,\) either \(f(x)\) is in , or it is undefined.

The set of the elements of  such that \(f(x)\) is defined and belongs to  is called the *domain of definition* of the function. A partial function from  to  is thus an ordinary function that has as its domain a subset of  called the domain of definition of the function. If the domain of definition equals , one often says that the partial function is a *total function*.

In several areas of mathematics, the term "function" refers to partial functions rather than to ordinary (total) functions. This is typically the case when functions may be specified in a way that makes difficult or even impossible to determine their domain.

In calculus, a *real-valued function of a real variable* or *real function* is a partial function from the set \(\R\) of the real numbers to itself. Given a real function \(f:x\mapsto f(x)\) its multiplicative inverse \(x\mapsto 1/f(x)\) is also a real function. The determination of the domain of definition of a multiplicative inverse of a (partial) function amounts to compute the zeros of the function, the values where the function is defined but not its multiplicative inverse.

Similarly, a *function of a complex variable* is generally a partial function whose domain of definition is a subset of the complex numbers \(\Complex\). The difficulty of determining the domain of definition of a complex function is illustrated by the multiplicative inverse of the Riemann zeta function: the determination of the domain of definition of the function \(z\mapsto 1/\zeta(z)\) is more or less equivalent to the proof or disproof of one of the major open problems in mathematics, the Riemann hypothesis.

In computability theory, a general recursive function is a partial function from the integers to the integers whose values can be computed by an algorithm (roughly speaking). The domain of definition of such a function is the set of inputs for which the algorithm does not run forever. A fundamental theorem of computability theory is that there cannot exist an algorithm that takes an arbitrary general recursive function as input and tests whether 0 belongs to its domain of definition (see Halting problem).

