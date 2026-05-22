> Content sourced from [Wikipedia](https://en.wikipedia.org/wiki/Function_%28mathematics%29) — CC BY-SA 4.0

# Function notation f(x)

### Formal definition
Diagram of a function
Diagram of a relation that is not a function. One reason is that 2 is the first element in more than one ordered pair.  Another reason is that neither 3 nor 4 are the first element (input) of any ordered pair.

The above definition of a function is essentially that of the founders of calculus, Leibniz, Newton and Euler. However, it cannot be formalized, since there is no mathematical definition of an "assignment". It is only at the end of the 19th century that the first formal definition of a function could be provided, in terms of set theory. This set-theoretic definition is based on the fact that a function establishes a *relation* between the elements of the domain and some (possibly all) elements of the codomain. Mathematically, a binary relation between two sets *X* and *Y* is a subset of the set of all ordered pairs \((x, y)\) such that \(x\in X\) and \(y\in Y.\) The set of all these pairs is called the Cartesian product of *X* and *Y* and denoted \(X\times Y.\) Thus, the above definition may be formalized as follows.

A *function* with domain *X* and codomain *Y* is a binary relation \(R\) between *X* and *Y* that satisfies the two following conditions:
* For every \(x\) in \(X\) there exists \(y\) in \(Y\) such that \((x,y)\in R.\)
* If \((x,y)\in R\) and \((x,z)\in R,\) then \(y=z.\)

This definition may be rewritten more formally, without referring explicitly to the concept of a relation, but using more notation (including set-builder notation):

A function is formed by three sets (often as an ordered triple), the *domain* \(X,\) the *codomain* \(Y,\) and the *graph* \(R\) that satisfy the three following conditions.
*\(R \subseteq \{(x,y) \mid x\in X, y\in Y\}\)
*\(\forall x\in X, \exists y\in Y, \left(x, y\right) \in R \qquad\)
*\((x,y)\in R \land (x,z)\in R \implies y=z\qquad\)

A relation satisfying these conditions is called a functional relation.

The more usual terminology and notation can be derived from this formal definition as follows. Let \(f\) be a function defined by a functional relation \(R\). For every \(x\) in the domain of \(f\), the unique element of the codomain that is related to \(x\) is denoted \(f(x)\). If \(y\) is this element, one writes commonly \(y=f(x)\) instead of \((x,y)\in R\,\) or \(xRy\), and one says that "\(f\) maps \(x\) to \(y\)", "\(y\) is the image by \(f\) of \(x\)", or "the application of \(f\) on \(x\) gives \(y\)", etc.

