> Content sourced from [Wikipedia](https://en.wikipedia.org/wiki/Rounding) — CC BY-SA 4.0

# Round to nearest 10

## Rounding to integer
The most basic form of rounding is to replace an arbitrary number by an integer. All the following rounding modes are concrete implementations of an abstract single-argument "round()" procedure. These are true functions (with the exception of those that use randomness).

### Directed rounding to an integer
These four methods are called **directed rounding to an integer**, as the displacements from the original number \(x\) to the rounded value \(y\) are all directed toward or away from the same limiting value (0, +∞, or −∞). Directed rounding is used in interval arithmetic and is often required in financial calculations.

If \(x\) is positive, round-down is the same as round-toward-zero, and round-up is the same as round-away-from-zero. If \(x\) is negative, round-down is the same as round-away-from-zero, and round-up is the same as round-toward-zero. In any case, if \(x\) is an integer, \(y\) is just \(x\).

Where many calculations are done in sequence, the choice of rounding method can have a very significant effect on the result. A famous instance involved a new index set up by the Vancouver Stock Exchange in 1982. It was initially set at 1000.000 (three decimal places of accuracy), and after 22 months had fallen to about 520, although the market appeared to be rising. The problem was caused by the index being recalculated thousands of times daily, and always being truncated (rounded down) to 3 decimal places, in such a way that the rounding errors accumulated. Recalculating the index for the same period using rounding to the nearest thousandth rather than truncation corrected the index value from 524.811 up to 1098.892.

For the examples below, sgn(*x*) refers to the sign function applied to the original number, \(x\).

#### Rounding down
One may **round down** (or take the **floor**, or **round toward negative infinity**): \(y\) is the largest integer that does not exceed \(x\).

\(y = \mathrm{floor}(x) = \left\lfloor x \right\rfloor = -\left\lceil -x \right\rceil\)

For example, 23.7 gets rounded to 23, and −23.2 gets rounded to −24.

#### Rounding up
One may also **round up** (or take the **ceiling**, or **round toward positive infinity**): \(y\) is the smallest integer that is not less than \(x\).

\(y = \operatorname{ceil}(x) = \left\lceil x \right\rceil = -\left\lfloor -x \right\rfloor\)

For example, 23.2 gets rounded to 24, and −23.7 gets rounded to −23.

#### Rounding toward zero
One may also **round toward zero** (or **truncate**, or **round away from infinity**): \(y\) is the integer that is closest to \(x\) such that it is between 0 and \(x\) (included); i.e. \(y\) is the integer part of \(x\), without its fraction digits.

\(y = \operatorname{truncate}(x)
= \sgn(x) \left\lfloor |x| \right\rfloor
= -\sgn(x) \left\lceil -|x| \right\rceil
= \begin{cases}
  \left\lfloor x \right\rfloor & x \ge 0 \\[5mu]
  \left\lceil x \right\rceil   & x < 0
\end{cases}\)

For example, 23.7 gets rounded to 23, and −23.7 gets rounded to −23.

#### Rounding away from zero
One may also **round away from zero** (or **round toward infinity**): \(y\) is the integer that is closest to 0 (or equivalently, to \(x\)) such that \(x\) is between 0 and \(y\) (included).

\(y = \sgn(x) \left\lceil |x| \right\rceil
= -\sgn(x) \left\lfloor -|x| \right\rfloor
= \begin{cases}
  \left\lceil x \right\rceil   & x \ge 0 \\[5mu]
  \left\lfloor x \right\rfloor & x < 0
\end{cases}\)

For example, 23.2 gets rounded to 24, and −23.2 gets rounded to −24.

### Rounding to the nearest integer
These six methods are called **rounding to the nearest integer**. Rounding a number \(x\) to the nearest integer requires some tie-breaking rule for those cases when \(x\) is exactly half-way between two integers – that is, when the fraction part of \(x\) is exactly 0.5.

If it were not for the 0.5 fractional parts, the round-off errors introduced by the round to nearest method would be symmetric: for every fraction that gets rounded down (such as 0.268), there is a complementary fraction (namely, 0.732) that gets rounded up by the same amount.

When rounding a large set of fixed-point numbers with uniformly distributed fractional parts, the rounding errors by all values, with the omission of those having 0.5 fractional part, would statistically compensate each other. This means that the expected (average) value of the rounded numbers is equal to the expected value of the original numbers when numbers with fractional part 0.5 from the set are removed.

In practice, floating-point numbers are typically used, which have even more computational nuances because they are not equally spaced.

#### Rounding half up
One may **round half up** (or **round half toward positive infinity**), a tie-breaking rule that is widely used in many disciplines. That is, half-way values of \(x\) are always rounded up. If the fractional part of \(x\) is exactly 0.5, then *y* = *x* + 0.5

\(y = \left\lfloor x + \tfrac12 \right\rfloor = -\left\lceil -x - \tfrac12 \right\rceil = \left\lceil \tfrac12 \lfloor 2x \rfloor \right\rceil\)

For example, 23.5 gets rounded to 24, and −23.5 gets rounded to −23.

Some programming languages (such as Java and Python) use "half up" to refer to *round half away from zero* rather than *round half toward positive infinity*.

This method only requires checking one digit to determine rounding direction in two's complement and similar representations.

#### Rounding half down
One may also **round half down** (or **round half toward negative infinity**) as opposed to the more common *round half up*. If the fractional part of \(x\) is exactly 0.5, then *y* = *x* − 0.5

\(y = \left\lceil x - \tfrac12 \right\rceil = -\left\lfloor -x + \tfrac12 \right\rfloor = \left\lfloor \tfrac12 \lceil 2x \rceil \right\rfloor\)

For example, 23.5 gets rounded to 23, and −23.5 gets rounded to −24.

Some programming languages (such as Java and Python) use "half down" to refer to *round half toward zero* rather than *round half toward negative infinity*.

#### Rounding half toward zero
One may also **round half toward zero** (or **round half away from infinity**) as opposed to the conventional *round half away from zero*. If the fractional part of \(x\) is exactly 0.5, then *y* = *x* − 0.5 if \(x\) is positive, and *y* = *x* + 0.5 if \(x\) is negative.

\(y = \sgn(x) \left\lceil |x| - \tfrac12 \right\rceil
= -\sgn(x) \left\lfloor -|x| + \tfrac12 \right\rfloor
= \begin{cases}
  \left\lceil x - \tfrac12 \right\rceil = \left\lfloor \tfrac12 \lceil 2x \rceil \right\rfloor & x \ge 0 \\[5mu]
  \left\lfloor x + \tfrac12 \right\rfloor = \left\lceil \tfrac12 \lfloor 2x \rfloor \right\rceil & x < 0
\end{cases}\)

For example, 23.5 gets rounded to 23, and −23.5 gets rounded to −23.

This method treats positive and negative values symmetrically, and therefore is free of overall positive/negative bias if the original numbers are positive or negative with equal probability. It does, however, still have bias toward zero.

#### Rounding half away from zero
One may also **round half away from zero** (or **round half toward infinity**), a tie-breaking rule that is commonly taught and used, namely: If the fractional part of \(x\) is exactly 0.5, then *y* = *x* + 0.5 if \(x\) is positive, and *y* = *x* − 0.5 if \(x\) is negative.

\(y = \sgn(x) \left\lfloor |x| + \tfrac12 \right\rfloor
= -\sgn(x) \left\lceil -|x| - \tfrac12 \right\rceil
= \begin{cases}
  \left\lfloor x + \tfrac12 \right\rfloor = \left\lceil \tfrac12 \lfloor 2x \rfloor \right\rceil & x \ge 0 \\[5mu]
  \left\lceil x - \tfrac12 \right\rceil = \left\lfloor \tfrac12 \lceil 2x \rceil \right\rfloor & x < 0
\end{cases}\)

For example, 23.5 gets rounded to 24, and −23.5 gets rounded to −24.

This can be more efficient on computers that use sign-magnitude representation for the values to be rounded, because only the first omitted digit needs to be considered to determine if it rounds up or down. This is one method used when rounding to significant figures due to its simplicity.

This method, also known as **commercial rounding**, treats positive and negative values symmetrically, and therefore is free of overall positive/negative bias if the original numbers are positive or negative with equal probability. It does, however, still have bias away from zero.

It is often used for currency conversions and price roundings (when the amount is first converted into the smallest significant subdivision of the currency, such as cents of a euro) as it is easy to explain by just considering the first fractional digit, independently of supplementary precision digits or sign of the amount (for strict equivalence between the paying and recipient of the amount).

#### Rounding half to even
One may also **round half to even**, a tie-breaking rule without positive/negative bias *and* without bias toward/away from zero. By this convention, if the fractional part of \(x\) is 0.5, then \(y\) is the even integer nearest to \(x\). Thus, for example, 23.5 becomes 24, as does 24.5; however, −23.5 becomes −24, as does −24.5. This function minimizes the expected error when summing over rounded figures, regardless of the inputs being mostly positive or mostly negative, provided they are neither mostly even nor mostly odd.

This variant of the round-to-nearest method is also called **convergent rounding**, **statistician's rounding**, **Dutch rounding**, **Gaussian rounding**, **odd–even rounding**, or **bankers' rounding**.

This is the default rounding mode used in IEEE 754 operations for results in binary floating-point formats.

By eliminating bias, repeated addition or subtraction of independent numbers, as in a one-dimensional random walk, will give a rounded result with an error that tends to grow in proportion to the square root of the number of operations rather than linearly.

date=October 2025

#### Rounding half to odd
One may also **round half to odd**, a similar tie-breaking rule to round half to even. In this approach, if the fractional part of \(x\) is 0.5, then \(y\) is the odd integer nearest to \(x\). Thus, for example, 23.5 becomes 23, as does 22.5; while −23.5 becomes −23, as does −22.5.

This method is also free from positive/negative bias and bias toward/away from zero, provided the numbers to be rounded are neither mostly even nor mostly odd.  It also shares the round half to even property of distorting the original distribution, as it increases the probability of odds relative to evens. It was the method used for bank balances in the United Kingdom when it decimalized its currencyreason=Needs a complete reference..

This variant is almost never used in computations, except in situations where one wants to avoid increasing the scale of floating-point numbers, which have a limited exponent range. With *round half to even*, a non-infinite number would round to infinity, and a small denormal value would round to a normal non-zero value. Effectively, this mode prefers preserving the existing scale of tie numbers, avoiding out-of-range results when possible for numeral systems of even radix (such as binary and decimal).reason=The *problem* would be avoided only for halfway numbers. So, how can this be useful in practice? Any practical example?.

### Randomized rounding to an integer
#### Alternating tie-breaking
One method, more obscure than most, is to alternate direction when rounding a number with 0.5 fractional part. All others are rounded to the closest integer. Whenever the fractional part is 0.5, alternate rounding up or down: for the first occurrence of a 0.5 fractional part, round up, for the second occurrence, round down, and so on. Alternatively, the first 0.5 fractional part rounding can be determined by a random seed.  "Up" and "down" can be any two rounding methods that oppose each other - toward and away from positive infinity or toward and away from zero.

If occurrences of 0.5 fractional parts occur significantly more than a restart of the occurrence "counting", then it is effectively bias free. With guaranteed zero bias, it is useful if the numbers are to be summed or averaged.

#### Random tie-breaking
If the fractional part of \(x\) is 0.5, choose \(y\) randomly between *x* + 0.5 and *x* − 0.5, with equal probability. All others are rounded to the closest integer.

Like round-half-to-even and round-half-to-odd, this rule is essentially free of overall bias, but it is also fair among even and odd \(y\) values. An advantage over alternate tie-breaking is that the last direction of rounding on the 0.5 fractional part does not have to be "remembered".

#### Stochastic rounding
Rounding as follows to one of the closest integer toward negative infinity and the closest integer toward positive infinity, with a probability dependent on the proximity is called stochastic rounding and will give an unbiased result on average.
\(\operatorname {Round} (x) = \begin{cases}
\lfloor x \rfloor & \text { with probability } 1 - (x - \lfloor x \rfloor) = \lfloor x \rfloor - x + 1 \\[5mu]
\lfloor x \rfloor + 1 & \text { with probability } {x - \lfloor x \rfloor}
\end{cases}\)

For example, 1.6 would be rounded to 1 with probability 0.4 and to 2 with probability 0.6.

Stochastic rounding can be accurate in a way that a rounding function can never be. For example, suppose one started with 0 and added 0.3 to that one hundred times while rounding the running total between every addition. The result would be 0 with regular rounding, but with stochastic rounding, the expected result would be 30, which is the same value obtained without rounding. This can be useful in machine learning where the training may use low precision arithmetic iteratively. Stochastic rounding is also a way to achieve 1-dimensional dithering.

### Comparison of approaches for rounding to an integer


! rowspan="3" | Value
! colspan="11"| Functional methods
! colspan="6" | Randomized methods

! colspan="4" | Directed rounding
! colspan="6" | Round to nearest
! rowspan="2" | Round to prepare for shorter precision
! colspan="2" | Alternating tie-break
! colspan="2" | Random tie-break
! colspan="2" | Stochastic
|- style="line-height:110%;"
! Down(toward −∞)
! Up(toward +∞)
! Toward 0
! Away From 0
! Half Down(toward −∞)
! Half Up(toward +∞)
! Half Toward 0
! Half Away From 0
! Half to Even
! Half to Odd
!Average
!SD
!Average
!SD
!Average
!SD

| +2.8
| rowspan="3" | +2
| rowspan="3" | +3
| rowspan="3" | +2
| rowspan="3" | +3
| +3
| rowspan="2" | +3
| +3
| rowspan="2" | +3
| +3
| rowspan="2" | +3
| rowspan="3" | +2
| +3
|0
| +3
|0
| +2.8
|0.04

| +2.5
| rowspan="3" | +2
| rowspan="3" | +2
| rowspan="4" | +2
| +2.505
|0
| +2.5
|0.05
| +2.5
|0.05

| +2.2
| rowspan="3" | +2
| rowspan="3" | +2
| rowspan="2" | +2
| rowspan="2" | +2
| rowspan="2" |0
| rowspan="2" | +2
| rowspan="2" |0
| +2.2
|0.04

| +1.8
| rowspan="3" | +1
| rowspan="3" | +2
| rowspan="3" | +1
| rowspan="3" | +2
| rowspan="6" | +1
| +1.8
| 0.04

| +1.5
| rowspan="3" | +1
| rowspan="3" | +1
| rowspan="4" | +1
| +1.505
| 0
| +1.5
| 0.05
| +1.5
| 0.05

| +1.2
| rowspan="3" | +1
| rowspan="3" | +1
| rowspan="2" | +1
| rowspan="2" | +1
| rowspan="2" | 0
| rowspan="2" | +1
| rowspan="2" | 0
| +1.2
| 0.04

| +0.8
| rowspan="3" | 0
| rowspan="3" | +1
| rowspan="6" | 0
| rowspan="3" | +1
| +0.8
| 0.04

| +0.5
| rowspan="3" | 0
| rowspan="4" | 0
| rowspan="4" | 0
| +0.505
| 0
| +0.5
| 0.05
| +0.5
| 0.05

| +0.2
| rowspan="3" | 0
| rowspan="2" | 0
| rowspan="2" | 0
| rowspan="2" | 0
| rowspan="2" | 0
| rowspan="2" | 0
| rowspan="2" | 0
| +0.2
| 0.04

| −0.2
| rowspan="3" | −1
| rowspan="3" | 0
| rowspan="3" | −1
| rowspan="6" | −1
|−0.2
| 0.04

| −0.5
| rowspan="3" | −1
| rowspan="3" | −1
| rowspan="4" | −1
|−0.495
| 0
|−0.5
| 0.05
|−0.5
| 0.05

| −0.8
| rowspan="3" | −1
| rowspan="3" | −1
| rowspan="2" | −1
| rowspan="2" | −1
| rowspan="2" | 0
| rowspan="2" | −1
| rowspan="2" | 0
|−0.8
| 0.04

| −1.2
| rowspan="3" | −2
| rowspan="3" | −1
| rowspan="3" | −1
| rowspan="3" | −2
|−1.2
| 0.04

| −1.5
| rowspan="3" | −2
| rowspan="3" | −2
| rowspan="4" | −2
|−1.495
| 0
|−1.5
| 0.05
|−1.5
| 0.05

| −1.8
| rowspan="3" | −2
| rowspan="3" | -2
| rowspan="2" | −2
| rowspan="2" | −2
| rowspan="2" | 0
| rowspan="2" | −2
| rowspan="2" | 0
|−1.8
| 0.04

|−2.2
| rowspan="3" |−3
| rowspan="3" |−2
| rowspan="3" |−2
| rowspan="3" |−3
| rowspan="3" |−2
|−2.2
|0.04

|−2.5
| rowspan="2" |−3
| rowspan="2" |−3
| rowspan="2" |−3
|−2.495
|0
|−2.5
|0.05
|−2.5
|0.05

|−2.8
|−3
|−3
|−3
|−3
|0
|−3
|0
|−2.8
|0.04
|

