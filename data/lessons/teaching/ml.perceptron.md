# The Perceptron

**Perceptron:** For inputs $x$ and weights $w$ with bias $b$, output $\text{step}(w \cdot x + b)$: 1 when the weighted sum plus bias is positive, 0 otherwise. It is the simplest neural unit — a linear score followed by a threshold.

## Worked: one forward computation

Take weights $w=[1, 2]$, inputs $x=[3, 4]$, bias $-10$:
1. Form the dot product: $1\cdot 3 + 2\cdot 4 = 11$.
2. Add the bias: $11 - 10 = 1$.
3. Apply the step: 1 is positive, so the perceptron outputs 1.

So the bias shifts the decision boundary: without it the sum 11 would still give 1, but a bias of $-12$ would flip the same input to 0.

## Worked: why XOR needs more

A single perceptron draws one straight line, so it separates only linearly separable sets:
1. AND is separable: $w=[1, 1]$, $b=-1.5$ gives outputs 0, 0, 0, 1 on the four corners — correct.
2. XOR labels opposite corners the same: no single line puts $(0,1)$ and $(1,0)$ on one side with $(0,0)$ and $(1,1)$ on the other.
3. Hence no choice of $w$ and $b$ solves XOR — check: any line leaving $(0,0)$ at 0 and $(1,1)$ at 0 traps one of the 1-corners on the wrong side.

So linear threshold units compose into multilayer networks: one hidden layer of two perceptrons solves XOR.

## When the step is not enough

The step activation has zero gradient almost everywhere, so it cannot be trained by gradient descent. Sigmoid or ReLU replacements keep the same weighted-sum front end but pass a usable derivative backward — the perceptron idea survives, the threshold does not.
