> Content sourced from [Wikipedia](https://en.wikipedia.org/wiki/Rounding) — CC BY-SA 4.0

# Round to nearest 1000

Rounding (disambiguation)


**Rounding** or **rounding off** is the process of adjusting a number to an approximate, more convenient value, often with a shorter or simpler representation. For example, replacing $23.4476 with $23.45, the fraction 312/937 with 1/3, or the expression √2 with 1.414.

Rounding is often done to obtain a value that is easier to report and communicate than the original. Rounding can also be important to avoid misleadingly precise reporting of a computed number, measurement, or estimate; for example, a quantity that was computed as 123456 but is known to be accurate only to within a few hundred units is usually better stated as "about 123500".

On the other hand, rounding of exact numbers will introduce some round-off error in the reported result. Rounding is almost unavoidable when reporting many computations – especially when dividing two numbers in integer or fixed-point arithmetic; when computing mathematical functions such as square roots, logarithms, and sines; or when using a floating-point representation with a fixed number of significant digits. In a sequence of calculations, these rounding errors generally accumulate, and in certain ill-conditioned cases they may make the result meaningless.

Accurate rounding of transcendental mathematical functions is difficult because the number of extra digits that need to be calculated to resolve whether to round up or down cannot be known in advance. This problem is known as "the table-maker's dilemma".

Rounding has many similarities to the quantization that occurs when physical quantities must be encoded by numbers or digital signals.

A wavy equals sign (≈) is sometimes used to indicate rounding of exact numbers, e.g. 9.98 ≈ 10. This sign was introduced by Alfred George Greenhill in 1892.

Ideal characteristics of rounding methods include:

# Rounding should be done by a function. This way, when the same input is rounded in different instances, the output is unchanged.
# Calculations done with rounding should be close to those done without rounding.
#* As a result of (1) and (2), the output from rounding should be close to its input, often as close as possible.
# To be considered rounding, the range will be a subset of the domain, often discrete. A classical range is the integers.
# Rounding should preserve symmetries that already exist between the domain and range. With finite precision (or a discrete domain), this translates to removing bias.
# A rounding method should have utility in computer science or human arithmetic where finite precision is used, and speed is a consideration.

Because it is not usually possible for a method to satisfy all ideal characteristics, many different rounding methods exist.

As a general rule, rounding is idempotent; i.e., once a number has been rounded, rounding it again to the same precision will not change its value. Rounding functions are also monotonic; i.e., rounding two numbers to the same absolute precision will not exchange their order (but may give the same value). In the general case of a discrete range, they are piecewise constant functions.

## Types of rounding
Typical rounding problems include:


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


## Rounding to other values
### Rounding to a specified multiple
The most common type of rounding is to round to an integer; or, more generally, to an integer multiple of some increment – such as rounding to whole tenths of seconds, hundredths of a dollar, to whole multiples of 1/2 or 1/8 inch, to whole dozens or thousands, etc.

In general, rounding a number \(x\) to a multiple of some specified positive value \(m\) entails the following steps:
\(\mathrm{roundToMultiple}(x, m) = \mathrm{round}(x/m) \times m\)

For example, rounding *x* = 2.1784 dollars to whole cents (i.e., to a multiple of 0.01) entails computing 2.1784 / 0.01 = 217.84, then rounding that to 218, and finally computing 218 × 0.01 = 2.18.

When rounding to a predetermined number of significant digits, the increment \(m\) depends on the magnitude of the number to be rounded (or of the rounded result).

The increment \(m\) is normally a finite fraction in whatever numeral system is used to represent the numbers. For display to humans, that usually means the decimal numeral system (that is, \(m\) is an integer times a power of 10, like 1/1000 or 25/100). For intermediate values stored in digital computers, it often means the binary numeral system (\(m\) is an integer times a power of 2).

The abstract single-argument "round()" function that returns an integer from an arbitrary real value has at least a dozen distinct concrete definitions presented in the rounding to integer section. The abstract two-argument "roundToMultiple()" function is formally defined here, but in many cases it is used with the implicit value *m* = 1 for the increment and then reduces to the equivalent abstract single-argument function, with also the same dozen distinct concrete definitions.

### Logarithmic rounding
#### Rounding to a specified power
Rounding to a specified *power* is very different from rounding to a specified *multiple*; for example, it is common in computing to need to round a number to a whole power of 2. The steps, in general, to round a positive number \(x\) to a power of some positive number \(b\) other than 1, are:
\(\mathrm{roundToPower}(x, b) = b^{\mathrm{round}(\log_b x)}, x>0, b>0, b \ne 1\)

Many of the caveats applicable to rounding to a multiple are applicable to rounding to a power.

In the chromatic "twelve-tone" scale of music, 3/2 is rounded to 2{sup|7/12} (a fifth), 4/3 is rounded to 2{sup|5/12} (a fourth), 5/4 is rounded to 2{sup|4/12} (a major third), 6/5 is rounded to 2{sup|3/12} (a minor third), and 9/8 is rounded to 2{sup|2/12} (a diminished third).

#### Scaled rounding
This type of rounding, which is also named **rounding to a logarithmic scale**, is a variant of rounding to a specified power. Rounding on a logarithmic scale is accomplished by taking the log of the amount and doing normal rounding to the nearest value on the log scale.

For example, resistors are supplied with preferred numbers on a logarithmic scale. In particular, for resistors with a 10% accuracy, they are supplied with nominal values 100, 120, 150, 180, 220, etc. rounded to multiples of 10 (E12 series). If a calculation indicates a resistor of 165 ohms is required then log(150) = 2.176, log(165) = 2.217 and log(180) = 2.255. The logarithm of 165 is closer to the logarithm of 180 therefore a 180 ohm resistor would be the first choice if there are no other considerations.

Whether a value *x* ∈ (*a*, *b*) rounds to \(a\) or \(b\) depends upon whether the squared value *x*{sup|2} is greater than or less than the product *ab*. The value 165 rounds to 180 in the resistors example because 165{sup|2} = 27225 is greater than 150 × 180 = 27000.

### Floating-point rounding
In floating-point arithmetic, rounding aims to turn a given value \(x\) into a value \(y\) with a specified number of significant digits. In other words, \(y\) should be a multiple of a number \(m\) that depends on the magnitude of \(x\). The number \(m\) is a power of the base (usually 2 or 10) of the floating-point representation.

Apart from this detail, all the variants of rounding discussed above apply to the rounding of floating-point numbers as well. The algorithm for such rounding is presented in the Scaled rounding section above, but with a constant scaling factor *s* = 1, and an integer base *b* > 1.

Where the rounded result would overflow the result for a directed rounding is either the appropriate signed infinity when "rounding away from zero", or the highest representable positive finite number (or the lowest representable negative finite number if \(x\) is negative), when "rounding toward zero". The result of an overflow for the usual case of *round to nearest* is always the appropriate infinity.

### Rounding to a simple fraction
In some contexts it is desirable to round a given number \(x\) to a "neat" fraction – that is, the nearest fraction *y* = *m*/*n* whose numerator \(m\) and denominator \(n\) do not exceed a given maximum. This problem is fairly distinct from that of rounding a value to a fixed number of decimal or binary digits, or to a multiple of a given unit \(m\). This problem is related to Farey sequences, the Stern–Brocot tree, and continued fractions.

### Rounding to an available value
Finished lumber, writing paper, electronic components, and many other products are usually sold in only a few standard values.

Many design procedures describe how to calculate an approximate value, and then "round" to some standard size using phrases such as "round down to nearest standard value", "round up to nearest standard value", or "round to nearest standard value".

When a set of preferred values is equally spaced on a logarithmic scale, choosing the closest preferred value to any given value can be seen as a form of scaled rounding. Such rounded values can be directly calculated.

### Arbitrary bins
More general rounding rules can separate values at arbitrary break points, used for example in data binning. A related mathematically formalized tool is signpost sequences, which use notions of distance other than the simple difference – for example, a sequence may round to the integer with the smallest *relative* (percent) error.

## Rounding in other contexts
### Dithering and error diffusion
When digitizing continuous signals, such as sound waves, the overall effect of a number of measurements is more important than the accuracy of each individual measurement. In these circumstances, dithering, and a related technique, error diffusion, are normally used. A related technique called pulse-width modulation is used to achieve analog type output from an inertial device by rapidly pulsing the power with a variable duty cycle. Delta-sigma modulation is commonly used for converting between real-world signals and digital signals, which allows control of the frequency statistics of quantization.

Error diffusion tries to ensure the error, on average, is minimized. When dealing with a gentle slope from one to zero, the output would be zero for the first few terms until the sum of the error and the current value becomes greater than 0.5, in which case a 1 is output and the difference subtracted from the error so far. Floyd–Steinberg dithering is a popular error diffusion procedure when digitizing images.

As a one-dimensional example, suppose the numbers 0.9677, 0.9204, 0.7451, and 0.3091 occur in order and each is to be rounded to a multiple of 0.01. In this case the cumulative sums, 0.9677, 1.8881 = 0.9677 + 0.9204, 2.6332 = 0.9677 + 0.9204 + 0.7451, and 2.9423 = 0.9677 + 0.9204 + 0.7451 + 0.3091, are each rounded to a multiple of 0.01: 0.97, 1.89, 2.63, and 2.94. The first of these and the differences of adjacent values give the desired rounded values: 0.97, 0.92 = 1.89 − 0.97, 0.74 = 2.63 − 1.89, and 0.31 = 2.94 − 2.63.

### Monte Carlo arithmetic
Monte Carlo arithmetic is a technique in Monte Carlo methods where the rounding is randomly up or down. Stochastic rounding can be used for Monte Carlo arithmetic, but in general, just rounding up or down with equal probability is more often used. Repeated runs will give a random distribution of results which can indicate the stability of the computation.

### Exact computation with rounded arithmetic
It is possible to use rounded arithmetic to evaluate the exact value of a function with integer domain and range. For example, if an integer \(n\) is known to be a perfect square, its square root can be computed by converting \(n\) to a floating-point value \(z\), computing the approximate square root \(x\) of \(z\) with floating point, and then rounding \(x\) to the nearest integer \(y\). If \(n\) is not too big, the floating-point round-off error in \(x\) will be less than 0.5, so the rounded value \(y\) will be the exact square root of \(n\). This is essentially why slide rules could be used for exact arithmetic.

### Double rounding

Rounding a number twice in succession to different levels of precision, with the latter precision being coarser, is not guaranteed to give the same result as rounding once to the final precision except in the case of directed rounding. For instance, rounding 9.46 to the nearest tenth gives 9.5, and then 10 when rounding to the nearest integer using rounding half to even, but would give 9 when rounded directly using the same method. Borman and Chatfield discuss the implications of double rounding when comparing data rounded to one decimal place to specification limits expressed using integers.

In *Martinez v. Allstate* and *Sendejo v. Farmers*, litigated between 1995 and 1997, the insurance companies argued that double rounding premiums was permissible and in fact required. The US courts ruled against the insurance companies and ordered them to adopt rules to ensure single rounding.

Some computer languages and the IEEE 754-2008 standard dictate that in straightforward calculations the result should not be rounded twice. This has been a particular problem with Java, as it is designed to be run identically on different machines; special programming tricks have had to be used to achieve this with x87 floating point. The Java language was changed to allow different results where the difference does not matter and require a strictfp qualifier to be used when the results have to conform accurately; strict floating point was restored in Java 17.

In some algorithms, an intermediate result is computed in a larger precision, then must be rounded to the final precision. Double rounding can be avoided by choosing an adequate rounding for the intermediate computation. This consists in avoiding to round to midpoints for the final rounding (except when the midpoint is exact). In binary arithmetic, the idea is to round the result toward zero, and set the least significant bit to 1 if the rounded result is inexact; this rounding is called *sticky rounding*. Equivalently, it consists in returning the intermediate result when it is exactly representable, and the nearest floating-point number with an odd significand otherwise; this is why it is also known as *rounding to odd*. A concrete implementation of this approach, for binary and decimal arithmetic, is implemented as Rounding to prepare for shorter precision.

### Rounding to prepare for shorter precision
This rounding mode is used to avoid getting a potentially wrong result after multiple roundings. This can be achieved if all roundings except the final one are done using rounding to prepare for shorter precision ("RPSP"), and only the final rounding uses the externally requested mode.

With decimal arithmetic, final digits of 0 and 5 are avoided when the input is not representable exactly; if there is a choice between numbers with the least significant digit 0 or 1, 4 or 5, 5 or 6, 9 or 0, then the digit different from 0 or 5 shall be selected; otherwise, the choice is arbitrary. IBM defines that, in the latter case, a digit with the smaller magnitude shall be selected. RPSP can be applied with the step between two consequent roundings as small as a single digit (for example, rounding to 1/10 can be applied after rounding to 1/100).
For example, when rounding to integer,
* 20.0 is rounded to 20;
* 20.01, 20.1, 20.9, 20.99, 21, 21.01, 21.9, 21.99 are rounded to 21 (avoiding a final 0);
* 22.0, 22.1, 22.9, 22.99 are rounded to 22;
* 24.0, 24.1, 24.9, 24.99 are rounded to 24 (avoiding a final 5);
* 25.0 is rounded to 25;
* 25.01, 25.1 are rounded to 26 (avoiding a final 5).

In the example from "Double rounding" section, rounding 9.46 to one decimal gives 9.4, which rounding to integer in turn gives 9.

With binary arithmetic, this rounding is also called "round to odd" (not to be confused with "round half to odd"). For example, when rounding to 1/4 (0.01 in binary),
* *x* = 2.0 ⇒ result is 2 (10.00 in binary)
* 2.0 < *x* < 2.5 ⇒ result is 2.25 (10.01 in binary)
* *x* = 2.5 ⇒ result is 2.5 (10.10 in binary)
* 2.5 < *x* < 3.0 ⇒ result is 2.75 (10.11 in binary)
* *x* = 3.0 ⇒ result is 3 (11.00 in binary)

For correct results with binary arithmetic, each rounding step must remove at least 2 binary digits, otherwise, wrong results may appear. For example,
* 3.125 RPSP to 1/4 ⇒ result is 3.25
* 3.25 RPSP to 1/2 ⇒ result is 3.5
* 3.5 round-half-to-even to 1 ⇒ result is 4 (wrong)

If the erroneous middle step is removed, the final rounding to integer rounds 3.25 to the correct value of 3.

RPSP is implemented in hardware in IBM zSeries and pSeries. In Python module "Decimal", Tcl module "math", Haskell package "decimal-arithmetic", and possibly others, this mode is called ROUND_05UP or round05up.

### Table-maker's dilemma
William M. Kahan coined the term "The Table-Maker's Dilemma" for the unknown cost of rounding transcendental functions:
Nobody knows how much it would cost to compute {math|*y*{sup|*w* correctly rounded for {em|every} two floating-point arguments at which it does not over/underflow. Instead, reputable math libraries compute elementary transcendental functions mostly within slightly more than half an unit in the last place

The [[IEEE 754 floating-point standard guarantees that add, subtract, multiply, divide, fused multiply–add, square root, and floating-point remainder will give the correctly rounded result of the infinite-precision operation. No such guarantee was given in the 1985 standard for more complex functions and they are typically only accurate to within the last bit at best. However, the 2008 standard guarantees that conforming implementations will give correctly rounded results which respect the active rounding mode; implementation of the functions, however, is optional.

Using the Gelfond–Schneider theorem and Lindemann–Weierstrass theorem, many of the standard elementary functions can be proved to return transcendental results, except on some well-known arguments; therefore, from a theoretical point of view, it is always possible to correctly round such functions. However, for an implementation of such a function, determining a limit for a given precision on how accurate results need to be computed, before a correctly rounded result can be guaranteed, may demand a lot of computation time or may be out of reach. In practice, when this limit is not known (or only a very large bound is known), some decision has to be made in the implementation (see below); but according to a probabilistic model, correct rounding can be satisfied with a very high probability when using an intermediate accuracy of up to twice the number of digits of the target format plus some small constant (after taking special cases into account).

Some programming packages offer correct rounding. The GNU MPFR package gives correctly rounded arbitrary precision results. Some other libraries implement elementary functions with correct rounding in IEEE 754 double precision (binary64):
* IBM's *ml4j*, which stands for *Mathematical Library for Java*, written by Abraham Ziv and Moshe Olshansky in 1999, correctly rounded to nearest only. This library was claimed to be portable, but only binaries for PowerPC/AIX, SPARC/Solaris and x86/Windows NT were provided. According to its documentation, this library uses a first step with an accuracy a bit larger than double precision, a second step based on double-double arithmetic, and a third step with a 768-bit precision based on arrays of IEEE 754 double-precision floating-point numbers.
* IBM's *Accurate portable mathematical library* (abbreviated as APMathLib or just MathLib), also called libultim, in rounding to nearest only. This library uses up to 768 bits of working precision. It was included in the GNU C Library in 2001, but the "slow paths" (providing correct rounding) were removed from 2018 to 2021.
* CRlibm, written in the old Arénaire team (LIP, ENS Lyon), first distributed in 2003. It supports the 4 rounding modes and is proved, using the knowledge of the hardest-to-round cases. More efficient than IBM MathLib. Succeeded by Metalibm (2014), which automates the formal proofs.
* Sun Microsystems's libmcr of 2004, in the 4 rounding modes. For the difficult cases, this library also uses multiple precision, and the number of words is increased by 2 each time the Table-maker's dilemma occurs (with undefined behavior in the very unlikely event that some limit of the machine is reached).
* The CORE-MATH project (2022) provides some correctly rounded functions in the 4 rounding modes for x86-64 processors. Proved using the knowledge of the hardest-to-round cases.
* LLVM libc provides some correctly rounded functions in the 4 rounding modes.

There exist computable numbers for which a rounded value can never be determined no matter how many digits are calculated. Specific instances cannot be given but this follows from the undecidability of the halting problem. For instance, if Goldbach's conjecture is true but unprovable, then the result of rounding the following value, \(n\), up to the next integer cannot be determined: either \(n\)=1+10\(^{−\(k\)}\) where \(k\) is the first even number greater than 4 which is not the sum of two primes, or \(n\)=1 if there is no such number. The rounded result is 2 if such a number \(k\) exists and 1 otherwise. The value before rounding can however be approximated to any given precision even if the conjecture is unprovable.

### Interaction with string searches
Rounding can adversely affect a string search for a number. For example, π rounded to four digits is "3.1416" but a simple search for this string will not discover "3.14159" or any other value of π rounded to more than four digits. In contrast, truncation does not suffer from this problem; for example, a simple string search for "3.1415", which is π truncated to four digits, will discover values of π truncated to more than four digits.

## History
The concept of rounding is very old, perhaps older than the concept of division itself. Some ancient clay tablets found in Mesopotamia contain tables with rounded values of reciprocals and square roots in base 60.
Rounded approximations to π, the length of the year, and the length of the month are also ancient—see base 60 examples.

The *round-half-to-even* method has served as American Standard Z25.1 and ASTM standard E-29 since 1940. The origin of the terms *unbiased rounding* and *statistician's rounding* are fairly self-explanatory. In the 1906 fourth edition of *Probability and Theory of Errors* Robert Simpson Woodward called this "the computer's rule", indicating that it was then in common use by human computers who calculated mathematical tables. For example, it was recommended in Simon Newcomb's c. 1882 book *Logarithmic and Other Mathematical Tables*. Lucius Tuttle's 1916 *Theory of Measurements* called it a "universally adopted rule" for recording physical measurements. Churchill Eisenhart indicated the practice was already "well established" in data analysis by the 1940s.

The origin of the term *bankers' rounding* remains more obscure. If this rounding method was ever a standard in banking, the evidence has proved extremely difficult to find. To the contrary, section 2 of the European Commission report *The Introduction of the Euro and the Rounding of Currency Amounts* suggests that there had previously been no standard approach to rounding in banking; and it specifies that "half-way" amounts should be rounded up.

Until the 1980s, the rounding method used in floating-point computer arithmetic was usually fixed by the hardware, poorly documented, inconsistent, and different for each brand and model of computer. This situation changed after the IEEE 754 floating-point standard was adopted by most computer manufacturers. The standard allows the user to choose among several rounding modes, and in each case specifies precisely how the results should be rounded. These features made numerical computations more predictable and machine-independent, and made possible the efficient and consistent implementation of interval arithmetic.

Currently, much research tends to round to multiples of 5 or 2. For example, Jörg Baten used age heaping in many studies, to evaluate the numeracy level of ancient populations. He came up with the ABCC Index, which enables the comparison of the numeracy among regions possible without any historical sources where the population literacy was measured.

## Rounding functions in programming languages
Most programming languages provide functions or special syntax to round fractional numbers in various ways. The earliest numeric languages, such as Fortran and C, would provide only one method, usually truncation (toward zero). This default method could be implied in certain contexts, such as when assigning a fractional number to an integer variable, or using a fractional number as an index of an array. Other kinds of rounding had to be programmed explicitly; for example, rounding a positive number to the nearest integer could be implemented by adding 0.5 and truncating.

In the last decades, however, the syntax and the standard libraries of most languages have commonly provided at least the four basic rounding functions (up, down, to nearest, and toward zero). The tie-breaking method can vary depending on the language and version or might be selectable by the programmer. Several languages follow the lead of the IEEE 754 floating-point standard, and define these functions as taking a double-precision float argument and returning the result of the same type, which then may be converted to an integer if necessary. This approach may avoid spurious overflows because floating-point types have a larger range than integer types. Some languages, such as PHP, provide functions that round a value to a specified number of decimal digits (e.g., from 4321.5678 to 4321.57 or 4300). In addition, many languages provide a printf or similar string formatting function, which allows one to convert a fractional number to a string, rounded to a user-specified number of decimal places (the *precision*). On the other hand, truncation (round to zero) is still the default rounding method used by many languages, especially for the division of two integer values.

In contrast, CSS and SVG do not define any specific maximum precision for numbers and measurements, which they treat and expose in their DOM and in their IDL interface as strings as if they had infinite precision, and do not discriminate between integers and floating-point values; however, the implementations of these languages will typically convert these numbers into IEEE 754 double-precision floating-point values before exposing the computed digits with a limited precision (notably within standard JavaScript or ECMAScript interface bindings).

## Other rounding standards
Some disciplines or institutions have issued standards or directives for rounding.

### US weather observations
In a guideline issued in mid-1966, the U.S. Office of the Federal Coordinator for Meteorology determined that weather data should be rounded to the nearest round number, with the "round half up" tie-breaking rule. For example, 1.5 rounded to integer should become 2, and −1.5 should become −1. Prior to that date, the tie-breaking rule was "round half away from zero".

### Negative zero in meteorology


Some meteorologists may write "−0" to indicate a temperature between 0.0 and −0.5 degrees (exclusive) that was rounded to an integer. This notation is used when the negative sign is considered important, no matter how small is the magnitude; for example, when rounding temperatures in the Celsius scale, where below zero indicates freezing.

### Negative zero in computation
Some computer representations of numbers, including the floating-point representation, distinguish between 0 (aka +0) and −0.  For these, rounding that would have been to 0 from a negative number can be rounded to −0 rather than +0.

