# Generating Functions

**Ordinary generating function:** For a sequence $(a_n)$, $A(x)=\sum_{n\ge0}a_nx^n$. Convolution $c_n=\sum_k a_kb_{n-k}$ is exactly coefficient extraction from the product $C(x)=A(x)B(x)$. The **exponential** generating function $E(x)=\sum a_nx^n/n!$ plays the same role for labelled structures, where products give binomial convolution.

## Worked: $1/(1-x)$ holds all ones

Multiply out $(1-x)(1+x+x^2+\cdots)$:
1. $x^k-x^{k+1}$ telescopes: every power cancels except the leading 1.
2. So $(1-x)\sum_{n\ge0}x^n=1$, i.e. $\sum_{n\ge0}x^n=1/(1-x)$.
3. Every coefficient is 1, matching the constant sequence $a_n=1$.

So $1/(1-x)$ is the ordinary generating function of $(1,1,1,\dots)$, and shifting gives $x^k/(1-x)$ for the sequence starting with $k$ zeros.

## Worked: Fibonacci closed form

Let $A(x)=\sum_{n\ge0}a_nx^n$ with $a_0=0$, $a_1=1$, $a_n=a_{n-1}+a_{n-2}$:
1. $A-xA-x^2A$ kills every recurrence relation, leaving only the $a_1$ term: $A(x)-xA(x)-x^2A(x)=x$.
2. Solve: $A(x)=x/(1-x-x^2)$.
3. Expand: $x+x^2+2x^3+3x^4+5x^5+\cdots$, and the coefficients $0,1,1,2,3,5,\dots$ are the Fibonacci numbers.

So a linear recurrence becomes one algebraic equation for $A(x)$, and series expansion reads the sequence back off.

## Ordinary versus exponential

For unlabelled classes use ordinary products; for labelled ones divide by factorials first. The product $E(x)=e^x\cdot e^x$ has coefficients $\sum_k\binom{n}{k}=2^n$: choosing which $k$ of $n$ labels go left is binomial convolution. In general, ordinary generating functions multiply by convolution, exponential ones by binomial convolution.
