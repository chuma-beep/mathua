# Wallis integrals

The Wallis integrals are the moments of powers of sine on the quarter period. They satisfy a two-term recurrence and have separate closed forms for odd and even powers.

## Odd powers

Define
\[
W_n=\int_0^{\pi/2}\sin^n x\,dx.
\]
For odd \(n=2m+1\),
\[
W_{2m+1}=\frac{2\cdot4\cdots 2m}{1\cdot3\cdots(2m+1)}.
\]
The product contains the even numbers up to \(2m\) and the odd numbers up to \(2m+1\). Every factor is a rational number, so the answer is rational.

The first values are \(W_1=1\), \(W_3=2/3\), \(W_5=8/15\) and \(W_7=16/35\). The sequence decreases toward zero.

## Even powers

For even \(n=2m\),
\[
W_{2m}=\frac{\pi}{2}\cdot\frac{1\cdot3\cdots(2m-1)}{2\cdot4\cdots 2m}.
\]
The answer is a rational multiple of \(\pi\). The coefficient is the product of the odd numbers divided by the product of the even numbers, with an extra factor \(1/2\).

The first coefficients are \(1/2\), \(3/8\), \(5/16\) and \(35/128\). They decrease toward zero, and the products are related to the central binomial coefficients.

## The recurrence

Integration by parts gives
\[
W_n=\frac{n-1}{n}W_{n-2}
\]
for \(n\ge2\). The recurrence connects powers of the same parity, so the odd and even families are computed separately.

The recurrence also proves that \(W_n\) decreases with \(n\) and that the ratio \(W_{2m}/W_{2m+1}\) tends to 1. That limit is the Wallis product for \(\pi\).
