# Arithmetico-geometric series

## The standard form

An arithmetico-geometric series multiplies an arithmetic progression by a geometric one: \(\sum_{n \ge 1} (a + (n-1)d) r^{n-1}\) with \(|r| < 1\). Taking \(a = 1\), \(d = 1\) and \(r = 1/2\) gives \(1 + 1 + 0.75 + 0.5 + \cdots\). The first four terms total 3.25. The ratio \(r\) controls the decay while the arithmetic factor grows, so convergence needs \(|r| < 1\).

## The summation formula

The sum is \(a/(1-r) + dr/(1-r)^2\). With \(a = 1\), \(d = 1\) and \(r = 1/2\) the two parts are 2 and 2, so the total is 4. The partial sums confirm it: after four terms the total is 3.25 and the fifth term 0.3125 lifts it to 3.5625.

## A worked total

A worked total: \(\sum_{n\ge1} n/2^n\). The series is \(\sum n r^n\) with \(r = 1/2\), whose sum is \(r/(1-r)^2 = 0.5/0.25 = 2\). Adding terms gives 1.625 after four and 1.875 after six, so the limit 2 is consistent with both.

A complete answer names the relation, substitutes the data and states the result in the language of the problem. When a second method is available, compare the two results before moving on.
