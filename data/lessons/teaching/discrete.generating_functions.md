# Generating Functions

**Ordinary generating function:** $A(x)=\sum_{n\ge0} a_n x^n$. Convolution $c_n=\sum_{k} a_k b_{n-k}$ ↔ $C(x)=A(x)B(x)$.

## Solving Recurrences

### Closed Forms
Fibonacci $a_n=a_{n-1}+a_{n-2}$ gives $A(x)=x/(1-x-x^2)$; extract coefficients via partial fractions.

### Exponential GFs
$E(x)=\sum a_n x^n/n!$ for labelled structures; ordinary for unlabelled.

## Example

Fibonacci: $A(x)=x/(1-x-x^2)$ expands to $0,x, x^2+..., $ coefficients are Fibonacci numbers.
