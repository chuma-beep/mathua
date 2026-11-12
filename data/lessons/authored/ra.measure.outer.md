# Outer measure

Outer measure assigns a size to every subset of the line by covering it with intervals. It is the first step of the Lebesgue construction.

## Covering by intervals

The outer measure of \(A\) is the infimum of the total lengths of countable covers by open intervals. A single point has outer measure 0, since an interval of length 0.001 covers it. The rationals in \([0,1]\) also have outer measure 0, because they are countable and each point can be given an interval of length \(\varepsilon/2^k\).

## Monotonicity and subadditivity

If \(A \subseteq B\) then \(m^*(A) \le m^*(B)\). Also \(m^*(A \cup B) \le m^*(A) + m^*(B)\). For \(A = [0,1]\) and \(B = [1,2]\) the two measures are both 1 and the sum is 2, which is the measure of the union. The inequality can be strict for sets that overlap in a complicated way.

## From outer to Lebesgue

Outer measure is not additive on all sets. A Vitali set built from one representative of each coset of the rationals has positive outer measure, and its complement has outer measure 1. The two numbers add to more than 1, which is the measure of the interval. Restricting to the sets that split every test set additively removes the pathology and produces Lebesgue measure.

Outer measure covers first and asks questions later, and the questions are answered by the splitting criterion.
