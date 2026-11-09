# Fibonacci and Binet's formula

The Fibonacci sequence begins with two ones, and every later term is the sum of the two before it. The recurrence is simple enough to run by hand, and Binet's formula gives the same terms in closed form.

## The recurrence

The sequence is defined by \(F_1 = 1\), \(F_2 = 1\) and \(F_k = F_{k-1} + F_{k-2}\) for \(k \ge 3\). The first terms are 1, 1, 2, 3, 5, 8, 13 and 21.

Computing a term means stepping forward from the two known values. Each step uses only the previous two terms, so the arithmetic stays small even when the index grows.

The recurrence is the practical way to compute a specific term. It is exact, it needs no rounding and it works for every positive index.

## Neighbouring terms

Two consecutive Fibonacci numbers determine the next one by addition. If \(F_n = 13\) and \(F_{n+1} = 21\), then \(F_{n+2} = 34\).

The same pair also determines the previous term by subtraction, because \(F_{n-1} = F_{n+1} - F_n\). With 13 and 21, the previous term is 8.

This is why a table of Fibonacci values can be extended in either direction. Moving forward adds, and moving backward subtracts.

## The sum of the first n terms

The prefix sums have a compact identity: \(F_1 + F_2 + \cdots + F_n = F_{n+2} - 1\). Adding the first five terms gives \(1+1+2+3+5 = 12\), and \(F_7 - 1 = 13 - 1 = 12\).

The identity can be proved by induction. Adding the next term to a prefix sum shifts the right side by exactly that term, because \(F_{n+3} - F_{n+2} = F_{n+1}\).

Binet's formula expresses the same sequence with powers of the golden ratio. It is useful for estimates and for proving identities, while the recurrence remains the cleanest route to an exact term.
