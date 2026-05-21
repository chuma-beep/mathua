> Content sourced from [Wikipedia](https://en.wikipedia.org/wiki/Greatest_common_divisor) — CC BY-SA 4.0

# Greatest common factor

In mathematics, the **greatest common divisor** (**GCD**), also known as **greatest common factor (GCF)**, of two or more integers, which are not all zero, is the largest positive integer that divides each of the integers. For two integers *x*, *y*, the greatest common divisor of *x* and *y* is denoted \(\gcd (x,y)\). For example, the GCD of 8 and 12 is 4, that is, gcd4.

In the name "greatest common divisor", the adjective "greatest" may be replaced by "highest", and the word "divisor" may be replaced by "factor", so that other names include **highest common factor (HCF)**, etc. Historically, other names for the same concept have included **greatest common measure**.

This notion can be extended to polynomials (see *Polynomial greatest common divisor*) and other commutative rings (see ** below).

## Overview
### Definition
The *greatest common divisor* (GCD) of integers and , at least one of which is nonzero, is the greatest positive integer such that is a divisor of both and ; that is, there are integers and such that 1=*a* = *de* and 1=*b* = *df*, and is the largest such integer. The GCD of and is generally denoted gcd(*a*, *b*). explains this as: "Many authors write (*a*, *b*) for g.c.d.(*a*, *b*). We do not, because we shall often use (*a*, *b*) to represent a point in the Euclidean plane."

When one of *a* and *b* is zero, the GCD is the absolute value of the nonzero integer: gcdgcd(0, *a*) = |*a*|. This case is important as the terminating step of the Euclidean algorithm.

The above definition is unsuitable for defining gcd(0, 0); since there is no greatest integer *n* such that 1=0 × *n* = 0. However, zero is its own greatest divisor if *greatest* is understood in the context of the divisibility relation, so gcd(0, 0) is commonly defined as 0. This preserves the usual identities for GCD, and in particular Bézout's identity, namely that gcd(*a*, *b*) generates the same ideal as. This convention is followed by many computer algebra systems. Nonetheless, some authors leave gcd(0, 0) undefined.

The GCD of and is their greatest positive common divisor in the preorder relation of divisibility. This means that the common divisors of and are exactly the divisors of their GCD. This is commonly proved by using either Euclid's lemma, the fundamental theorem of arithmetic, or the Euclidean algorithm. This is the meaning of "greatest" that is used for the generalizations of the concept of GCD.

### Example
The number 54 can be expressed as a product of two integers in several different ways:
\(54 \times 1 = 27 \times 2 = 18 \times 3 = 9 \times 6.\)

Thus the complete list of *divisors* of 54 is 1, 2, 3, 6, 9, 18, 27, 54.
Similarly, the divisors of 24 are 1, 2, 3, 4, 6, 8, 12, 24.
The numbers that these two lists have *in common* are the *common divisors* of 54 and 24, that is,
\(1, 2, 3, 6.\)

Of these, the greatest is 6, so it is the *greatest common divisor*:
\(\gcd(54,24) = 6.\)

Computing all divisors of the two numbers in this way is usually not efficient, especially for large numbers that have many divisors. Much more efficient methods are described in **.

### Coprime numbers

Two numbers are called relatively prime, or coprime, if their greatest common divisor equals 1. For example, 9 and 28 are coprime.

### A geometric view

For example, a 24-by-60 rectangular area can be divided into a grid of: 1-by-1 squares, 2-by-2 squares, 3-by-3 squares, 4-by-4 squares, 6-by-6 squares or 12-by-12 squares. Therefore, 12 is the greatest common divisor of 24 and 60. A 24-by-60 rectangular area can thus be divided into a grid of 12-by-12 squares, with two squares along one edge (1=24/12 = 2) and five squares along the other (1=60/12 = 5).

## Applications
### Reducing fractions

The greatest common divisor is useful for reducing fractions to the lowest terms. For example, gcd14, therefore,
\(\frac{42}{56}=\frac{3 \cdot 14 }{ 4 \cdot 14}=\frac{3 }{ 4}.\)

### Least common multiple

The least common multiple of two integers that are not both zero can be computed from their greatest common divisor, by using the relation
\(\operatorname{lcm}(a,b)=\frac{|a\cdot b|}{\operatorname{gcd}(a,b)}.\)

## Calculation

### Using prime factorizations
Greatest common divisors can be computed by determining the prime factorizations of the two numbers and comparing factors. For example, to compute gcd(48, 180), we find the prime factorizations 48 = 2\(^{4}\) · 3\(^{1}\) and 180 = 2\(^{2}\) · 3\(^{2}\) · 5\(^{1}\); the GCD is then 2\(^{min(4,2)}\) · 3\(^{min(1,2)}\) · 5\(^{min(0,1)}\) = 2\(^{2}\) · 3\(^{1}\) · 5\(^{0}\) = 12 The corresponding LCM is then
2\(^{max(4,2)}\) · 3\(^{max(1,2)}\) · 5\(^{max(0,1)}\) =
2\(^{4}\) · 3\(^{2}\) · 5\(^{1}\) = 720.

In practice, this method is only feasible for small numbers, as computing prime factorizations takes too long.

### Euclid's algorithm

The method introduced by Euclid for computing greatest common divisors is based on the fact that, given two positive integers and such that *a* > *b*, the common divisors of and are the same as the common divisors of *a* – *b* and. So, Euclid's method for computing the greatest common divisor of two positive integers consists of replacing the larger number with the difference of the numbers, and repeating this until the two numbers are equal: that is their greatest common divisor.

For example, to compute gcd(48,18), one proceeds as follows:
\(\begin{align}\gcd(48,18)\quad&\to\quad \gcd(48-18, 18)= \gcd(30,18)\\
&\to \quad \gcd(30-18, 18)= \gcd(12,18)\\
&\to \quad \gcd(12,18-12)= \gcd(12,6)\\
&\to \quad \gcd(12-6,6)= \gcd(6,6).\end{align}\)
So gcd6.

This method can be very slow if one number is much larger than the other. So, the variant that follows is generally preferred.

### Euclidean algorithm

A more efficient method is the *Euclidean algorithm*, a variant in which the difference of the two numbers and is replaced by the *remainder* of the Euclidean division (also called *division with remainder*) of by. Denoting this remainder as *a* mod *b*, the algorithm replaces (*a*, *b*) with (*b*, *a* mod *b*) repeatedly until the pair is (*d*, 0), where is the greatest common divisor.

For example, to compute gcd(48,18), the computation is as follows:
\(\begin{align}\gcd(48,18)\quad&\to\quad \gcd(18, 48\bmod 18)= \gcd(18, 12)\\
&\to \quad \gcd(12, 18\bmod 12)= \gcd(12,6)\\
&\to \quad \gcd(6,12\bmod 6)= \gcd(6,0).\end{align}\)
This again gives gcd6.

### Binary GCD algorithm

The binary GCD algorithm is a variant of Euclid's algorithm that is specially adapted to the binary representation of the numbers, which is used in most computers.

The binary GCD algorithm differs from Euclid's algorithm essentially by dividing by two every even number that is encountered during the computation. Its efficiency results from the fact that, in binary representation, testing parity consists of testing the right-most digit, and dividing by two consists of removing the right-most digit.

The method is as follows, starting with *a* and *b* that are the two positive integers whose GCD is sought.
# If *a* and *b* are both even, then divide both by two until at least one of them becomes odd; let be the number of these paired divisions.
# If *a* is even, then divide it by two until it becomes odd.
# If *b* is even, then divide it by two until it becomes odd.
#: Now, *a* and *b* are both odd and will remain odd until the end of the computation
# While *a* ≠ *b* do
#* If *a* > *b*, then replace with *a* – *b* and divide the result by two until becomes odd (as *a* and *b* are both odd, there is, at least, one division by 2).
#* If *a* < *b*, then replace with *b* – *a* and divide the result by two until becomes odd.
# Now, 1=*a* = *b*, and the greatest common divisor is \(2^d a.\)

Step 1 determines as the highest power of 2 that divides *a* and *b*, and thus their greatest common divisor. None of the steps changes the set of the odd common divisors of *a* and *b*. This shows that when the algorithm stops, the result is correct. The algorithm stops eventually; since each steps divides at least one of the operands by at least 2. Moreover, the number of divisions by 2 and thus the number of subtractions is at most the total number of digits.

Example: (*a*, *b*, *d*) = (48, 18, 0) → (24, 9, 1) → (12, 9, 1) → (6, 9, 1) → (3, 9, 1) → (3, 3, 1) ; the original GCD is thus the product 6 of 1=2\(^{*d*}\) = 2\(^{1}\) and 1=*a* = *b* = 3.

The binary GCD algorithm is particularly easy to implement and particularly efficient on binary computers. Its computational complexity is
\(O((\log a + \log b)^2).\)
The square in this complexity comes from the fact that division by 2 and subtraction take a time that is proportional to the number of bits of the input.

The computational complexity is usually given in terms of the length *n* of the input. Here, this length is 1=*n* = log *a* + log *b*, and the complexity is thus
\(O(n^2)\).

### Lehmer's GCD algorithm

Lehmer's algorithm is based on the observation that the initial quotients produced by Euclid's algorithm can be determined based on only the first few digits; this is useful for numbers that are larger than a computer word. In essence, one extracts initial digits, typically forming one or two computer words, and runs Euclid's algorithms on these smaller numbers, as long as it is guaranteed that the quotients are the same with those that would be obtained with the original numbers. The quotients are collected into a small 2-by-2 transformation matrix (a matrix of single-word integers) to reduce the original numbers. This process is repeated until numbers are small enough that the binary algorithm (see below) is more efficient.

This algorithm improves speed, because it reduces the number of operations on very large numbers, and can use hardware arithmetic for most operations. In fact, most of the quotients are very small, so a fair number of steps of the Euclidean algorithm can be collected in a 2-by-2 matrix of single-word integers. When Lehmer's algorithm encounters a quotient that is too large, it must fall back to one iteration of Euclidean algorithm, with a Euclidean division of large numbers.

### Other methods

If *a* and *b* are both nonzero, the greatest common divisor of *a* and *b* can be computed by using least common multiple (LCM) of *a* and *b*:
\(\gcd(a,b)=\frac{|a\cdot b|}{\operatorname{lcm}(a,b)}\),
but more commonly the LCM is computed from the GCD.

Using Thomae's function *f*,
\(\gcd(a,b) = a f\left(\frac b a\right),\)
which generalizes to *a* and *b* rational numbers or commensurable real numbers.

Keith Slavin has shown that for odd *a* ≥ 1:
\(\gcd(a,b)=\log_2\prod_{k=0}^{a-1} (1+e^{-2i\pi k b/a})\)
which is a function that can be evaluated for complex *b*. Wolfgang Schramm has shown that
\(\gcd(a,b)=\sum\limits_{k=1}^a \exp (2\pi ikb/a) \cdot \sum\limits_{d\left| a\right.} \frac{c_d (k)}{d}\)
is an entire function in the variable *b* for all positive integers *a* where *c*\(_{*d*}\)(*k*) is Ramanujan's sum.

### Complexity
The computational complexity of the computation of greatest common divisors has been widely studied. If one uses the Euclidean algorithm and the elementary algorithms for multiplication and division, the computation of the greatest common divisor of two integers of at most bits is *O*(*n*\(^{2}\)). This means that the computation of greatest common divisor has, up to a constant factor, the same complexity as the multiplication.

However, if a fast multiplication algorithm is used, one may modify the Euclidean algorithm for improving the complexity, but the computation of a greatest common divisor becomes slower than the multiplication. More precisely, if the multiplication of two integers of *n* bits takes a time of *T*(*n*), then the fastest known algorithm for greatest common divisor has a complexity *O*(*T*(*n*) log *n*). This implies that the fastest known algorithm has a complexity of *O*(*n* (log *n*)\(^{2}\)).

Previous complexities are valid for the usual models of computation, specifically multitape Turing machines and random-access machines.

The computation of the greatest common divisors belongs thus to the class of problems solvable in quasilinear time. *A fortiori*, the corresponding decision problem belongs to the class P of problems solvable in polynomial time. The GCD problem is not known to be in NC, and so there is no known way to parallelize it efficiently; nor is it known to be P-complete, which would imply that it is unlikely to be possible to efficiently parallelize GCD computation. Shallcross et al. showed that a related problem (EUGCD, determining the remainder sequence arising during the Euclidean algorithm) is NC-equivalent to the problem of integer linear programming with two variables; if either problem is in **NC** or is **P-complete**, the other is as well. Since **NC** contains NL, it is also unknown whether a space-efficient algorithm for computing the GCD exists, even for nondeterministic Turing machines.

Although the problem is not known to be in **NC**, parallel algorithms asymptotically faster than the Euclidean algorithm exist; the fastest known deterministic algorithm is by Chor and Goldreich, which (in the CRCW-PRAM model) can solve the problem in *O*(*n*/log *n*) time with *n*\(^{1+*ε*}\) processors. Randomized algorithms can solve the problem in *O*((log *n*)\(^{2}\)) time on \(O\left(e^{\sqrt{n \log n\right)\) processors (this is superpolynomial).

## Properties
* For every positive integer *a*, gcd*a*.
* Every common divisor of *a* and *b* is a divisor of gcd(*a*, *b*).
* gcd(*a*, *b*), where *a* and *b* are not both zero, may be defined alternatively and equivalently as the smallest positive integer *d* which can be written in the form 1=*d* = *a*⋅*p* + *b*⋅*q*, where *p* and *q* are integers. This expression is called Bézout's identity. Numbers *p* and *q* like this can be computed with the extended Euclidean algorithm.
* gcd|*a*|, for *a* ≠ 0; since any number is a divisor of 0, and the greatest divisor of *a* is |*a*|. This is usually used as the base case in the Euclidean algorithm.
* If *a* divides the product *b*⋅*c*, and gcd*d*, then *a*/*d* divides *c*.
* If *m* is a positive integer, then gcd*m*⋅gcd(*a*, *b*).
* If *m* is any integer, then gcdgcd(*a*, *b*). Equivalently, gcdgcd(*a*,*b*).
* If *m* is a positive common divisor of *a* and *b*, then gcdgcd(*a*, *b*)/*m*.
* If gcd*d*, then gcd1.
* The GCD is a commutative function: gcdgcd(*b*, *a*).
* The GCD is an associative function: 1=gcd(*a*, gcd(*b*, *c*)) = gcd(gcd(*a*, *b*), *c*). Thus 1=gcd(*a*, *b*, *c*,.. .) can be used to denote the GCD of multiple arguments.
* The GCD is a multiplicative function in the following sense: if *a*\(_{1}\) and *a*\(_{2}\) are relatively prime, then 1=gcd(*a*\(_{1}\)⋅*a*\(_{2}\), *b*) = gcd(*a*\(_{1}\), *b*)⋅gcd(*a*\(_{2}\), *b*).
* gcd(*a*, *b*) is closely related to the least common multiple lcm(*a*, *b*): we have
*: 1=gcd(*a*, *b*)⋅lcm(*a*, *b*) = |*a*⋅*b*|.
This formula is often used to compute least common multiples: one first computes the GCD with Euclid's algorithm and then divides the product of the given numbers by their GCD.
* The following versions of distributivity hold true:
*: 1=gcd(*a*, lcm(*b*, *c*)) = lcm(gcd(*a*, *b*), gcd(*a*, *c*))
*: 1=lcm(*a*, gcd(*b*, *c*)) = gcd(lcm(*a*, *b*), lcm(*a*, *c*)).
* If we have the unique prime factorizations of 1=*a* = *p*\(_{1}\)\(^{*e*\(_{1}\)}\) *p*\(_{2}\)\(^{*e*\(_{2}\)}\) ⋅⋅⋅ *p*\(_{*m*}\)\(^{*e*\(_{*m*}\)}\) and 1=*b* = *p*\(_{1}\)\(^{*f*\(_{1}\)}\) *p*\(_{2}\)\(^{*f*\(_{2}\)}\) ⋅⋅⋅ *p*\(_{*m*}\)\(^{*f*\(_{*m*}\)}\) where 1=*e*\(_{*i*}\) ≥ 0 and 1=*f*\(_{*i*}\) ≥ 0, then the GCD of *a* and *b* is
*: gcd*p*\(_{1}\)\(^{min(*e*\(_{1}\),*f*\(_{1}\))}\) *p*\(_{2}\)\(^{min(*e*\(_{2}\),*f*\(_{2}\))}\) ⋅⋅⋅ *p*\(_{*m*}\)\(^{min(*e*\(_{*m*}\),*f*\(_{*m*}\))}\).
* It is sometimes useful to define gcd0 and lcm0 because then the natural numbers become a complete distributive lattice with GCD as meet and LCM as join operation. This extension of the definition is also compatible with the generalization for commutative rings given below.
* In a Cartesian coordinate system, gcd(*a*, *b*) can be interpreted as the number of segments between points with integral coordinates on the straight line segment joining the points (0, 0) and (*a*, *b*).
* For non-negative integers *a* and *b*, where *a* and *b* are not both zero, provable by considering the Euclidean algorithm in base *n*:
*: 1=gcd(*n*\(^{*a*}\) − 1, *n*\(^{*b*}\) − 1) = *n*\(^{gcd(*a*,*b*)}\) − 1.
* An identity involving Euler's totient function:
*: \(\gcd(a,b) = \sum_{k|a \text{ and }k|b} \varphi(k). \)
* GCD Summatory function (Pillai's arithmetical function):
\(\sum_{k=1}^n \gcd(k,n)
= \sum_{d|n} d \varphi \left( \frac n d \right)
=n\sum_{d|n}\frac{\varphi(d)}{d}
=n\prod_{p|n}\left(1+\nu_p(n)\left(1-\frac{1}{p}\right)\right)\) where \(\nu_p(n)\) is the *p*-adic valuation.

## Probabilities and expected value
In 1972, James E. Nymann showed that *k* integers, chosen independently and uniformly from , are coprime with probability 1/*ζ*(*k*) as *n* goes to infinity, where *ζ* refers to the Riemann zeta function. (See coprime for a derivation.) This result was extended in 1987 to show that the probability that *k* random integers have greatest common divisor *d* is *d*\(^{−*k*}\)/ζ(*k*).

Using this information, the expected value of the greatest common divisor function can be seen (informally) to not exist when 1=*k* = 2. In this case the probability that the GCD equals *d* is *d*\(^{−2}\)/*ζ*(2), so we have
\(\mathrm{E}( \mathrm{2} ) = \sum_{d=1}^\infty d (d^{-2}/\zeta(2)) = \frac{1}{\zeta(2)} \sum_{d=1}^\infty \frac{1}{d}.\)

This last summation is the harmonic series, which diverges. However, when *k* ≥ 3, the expected value is well-defined, and by the above argument, it is
\(\mathrm{E}(k) = \sum_{d=1}^\infty d^{1-k} \zeta(k)^{-1} = \frac{\zeta(k-1)}{\zeta(k)}.\)

For 1=*k* = 3, this is approximately equal to 1.3684. For 1=*k* = 4, it is approximately 1.1106.

## In commutative rings

The notion of greatest common divisor can more generally be defined for elements of an arbitrary commutative ring, although in general there need not exist one for every pair of elements.

* If is a commutative ring, and and are in , then an element of is called a *common divisor* of and if it divides both and (that is, if there are elements and in such that *d*·*x* = *a* and *d*·*y* = *b*).
* If is a common divisor of and , and every common divisor of and divides , then is called a *greatest common divisor* of and *b*.

With this definition, two elements and may very well have several greatest common divisors, or none at all. If is an integral domain, then any two GCDs of and must be associate elements; since by definition either one must divide the other. Indeed, if a GCD exists, any one of its associates is a GCD as well.

Existence of a GCD is not assured in arbitrary integral domains. However, if is a unique factorization domain or any other GCD domain, then any two elements have a GCD. If is a Euclidean domain in which euclidean division is given algorithmically (as is the case for instance when 1=*R* = *F*[*X*] where is a field, or when is the ring of Gaussian integers), then greatest common divisors can be computed using a form of the Euclidean algorithm based on the division procedure.

The following is an example of an integral domain with two elements that do not have a GCD:
\(R = \mathbb{Z}\left[\sqrt{-3}\,\,\right],\quad a = 4 = 2\cdot 2 = \left(1+\sqrt{-3}\,\,\right)\left(1-\sqrt{-3}\,\,\right),\quad b = \left(1+\sqrt{-3}\,\,\right)\cdot 2.\)

The elements 2 and \(1 + \sqrt{-3}\) are two maximal common divisors (that is, any common divisor which is a multiple of 2 is associated to 2, the same holds for \(1 + \sqrt{-3}\), but they are not associated, so there is no greatest common divisor of and *b*.

Corresponding to the Bézout property we may, in any commutative ring, consider the collection of elements of the form *pa* + *qb*, where and range over the ring. This is the ideal generated by and , and is denoted simply (*a*, *b*). In a ring all of whose ideals are principal (a principal ideal domain or PID), this ideal will be identical with the set of multiples of some ring element *d*; then this is a greatest common divisor of and *b*. But the ideal (*a*, *b*) can be useful even when there is no greatest common divisor of and *b*. (Indeed, Ernst Kummer used this ideal as a replacement for a GCD in his treatment of Fermat's Last Theorem, although he envisioned it as the set of multiples of some hypothetical, or *ideal*, ring element , whence the ring-theoretic term.)
