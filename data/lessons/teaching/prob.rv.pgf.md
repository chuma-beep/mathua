# Probability Generating Functions

This lesson develops the definition, the main calculation and the interpretation of the idea.

## Encoding a count distribution

For a nonnegative integer-valued random variable, the probability generating function is

$$G(s)=E[s^X]=\sum_k P(X=k)s^k.$$

The coefficient of $s^k$ is the probability that $X=k$. The function packages the whole distribution into one expression, and evaluating it at $s=0$ returns $P(X=0)$.

The radius of convergence is at least one because the probabilities sum to one. Inside that radius the series can be differentiated term by term.

## Derivatives and moments

Differentiating once and setting $s=1$ gives the mean:

$$G'(1)=E[X].$$

Differentiating twice gives the second factorial moment, so the variance follows from $G''(1)+G'(1)-G'(1)^2$. Higher derivatives give higher factorial moments.

For a Bernoulli variable with success probability $p$, $G(s)=1-p+ps$. For a Poisson variable with rate $\lambda$, $G(s)=e^{\lambda(s-1)}$. These forms are often easier to manipulate than the corresponding sums.

## Sums of independent counts

If $X$ and $Y$ are independent, the generating function of $X+Y$ is the product:

$$G_{X+Y}(s)=G_X(s)G_Y(s).$$

That product rule is the main reason to use generating functions. It turns convolution of distributions into multiplication of ordinary functions, and it proves that independent Poisson variables add by adding their rates.

The same idea applies to a random sum, where the number of terms is itself random, provided the terms and the count are independent.

