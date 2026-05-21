> Content sourced from [Wikipedia](https://en.wikipedia.org/wiki/Least_common_multiple) — CC BY-SA 4.0

# Least common multiple

In arithmetic and number theory, the **least common multiple** (**LCM**), **lowest common multiple**, or **smallest common multiple** (**SCM**) of two integers *a* and *b*, usually denoted by lcm(*a*, *b*), is the smallest positive integer that is divisible by both *a* and *b*. Since division of integers by zero is undefined, this definition has meaning only if *a* and *b* are both different from zero. However, some authors define lcm(*a*, 0) as 0 for all *a*; since 0 is the only common multiple of *a* and 0.

The least common multiple of the denominators of two fractions is the "lowest common denominator" (lcd), and can be used for adding, subtracting or comparing the fractions.

The least common multiple of more than two integers *a*, *b*, *c*,... , usually denoted by lcm(*a*, *b*, *c*,... ), is defined as the smallest positive integer that is divisible by each of *a*, *b*, *c*,... ## Overview
A multiple of a number is the product of that number and an integer. For example, 10 is a multiple of 5 because 5 × 2 = 10, so 10 is divisible by 5 and 2. Because 10 is the smallest positive integer that is divisible by both 5 and 2, it is the least common multiple of 5 and 2. By the same principle, 10 is the least common multiple of −5 and −2 as well.

### Notation
The least common multiple of two integers *a* and *b* is denoted as lcm(*a*, *b*). Some older textbooks use [*a*, *b*].

### Example
\(\operatorname{lcm}(4, 6)\)

Multiples of 4 are:

\(4, 8, 12, 16, 20, 24, 28, 32, 36, 40, 44, 48, 52, 56, 60, 64, 68, 72, 76,... \)

Multiples of 6 are:

\(6, 12, 18, 24, 30, 36, 42, 48, 54, 60, 66, 72,... \)

*Common multiples* of 4 and 6 are the numbers that are in both lists:

\(12, 24, 36, 48, 60, 72,... \)

In this list, the smallest number is 12. Hence, the *least common multiple* is 12.

## Applications
When adding, subtracting, or comparing simple fractions, the least common multiple of the denominators (often called the lowest common denominator) is used, because each of the fractions can be expressed as a fraction with this denominator. For example,

\({2\over21}+{1\over6}={4\over42}+{7\over42}={11\over42}\)

where the denominator 42 was used, because it is the least common multiple of 21 and 6.

### Gears problem
Suppose there are two meshing gears in a machine, having *m* and *n* teeth, respectively, and the gears are marked by a line segment drawn from the center of the first gear to the center of the second gear. When the gears begin rotating, the number of rotations the first gear must complete to realign the line segment can be calculated by using \(\operatorname{lcm}(m, n)\). The first gear must complete \(\frac{\operatorname{lcm}(m, n)}{m}\) rotations for the realignment. By that time, the second gear will have made \(\frac{\operatorname{lcm}(m, n)}{n}\) rotations.

### Planetary alignment

Suppose there are three planets revolving around a star that take *l*, *m*, and *n* units of time, respectively, to complete their orbits. Assume that *l*, *m*, and *n* are integers. Assuming the planets started moving around the star after an initial linear alignment, all the planets attain a linear alignment again after \(\operatorname{lcm}(l, m, n)\) units of time. At this time, the first, second and third planet will have completed \(\frac{\operatorname{lcm}(l, m, n)}{l}\), \(\frac{\operatorname{lcm}(l, m, n)}{m}\) and \(\frac{\operatorname{lcm}(l, m, n)}{n}\) orbits, respectively, around the star.

## Calculation
There are several ways to compute least common multiples.

### Using the greatest common divisor
The least common multiple can be computed from the greatest common divisor (gcd) with the formula
\(\operatorname{lcm}(a,b)=\frac{|ab|}{\gcd(a,b)}.\)

To avoid introducing integers that are larger than the result, it is convenient to use the equivalent formulas
\(\operatorname{lcm}(a,b)=|a|\,\frac{|b|}{\gcd(a,b)} = |b|\,\frac{|a|}{\gcd(a,b)} ,\)
where the result of the division is always an integer.

These formulas are also valid when exactly one of *a* and *b* is 0; since gcd|*a*|. However, if both *a* and *b* are 0, these formulas would cause division by zero; so, lcm0 must be considered as a special case.

To return to the example above,

\(\operatorname{lcm}(21,6)
=6\times\frac {21}{\gcd(21,6)}
=6\times\frac {21} 3
=6\times 7
= 42.\)

There are fast algorithms, such as the Euclidean algorithm for computing the gcd that do not require the numbers to be factored. For very large integers, there are even faster algorithms for the three involved operations (multiplication, gcd, and division); see Fast multiplication. As these algorithms are more efficient with factors of similar size, it is more efficient to divide the largest argument of the lcm by the gcd of the arguments, as in the example above.

### Using prime factorization
The unique factorization theorem indicates that every positive integer greater than 1 can be written in only one way as a product of prime numbers. The prime numbers can be considered as the atomic elements which, when combined, make up a composite number.

For example:

\(90 = 2^1 \cdot 3^2 \cdot 5^1 = 2 \cdot 3 \cdot 3 \cdot 5.\)

Here, the composite number 90 is made up of one atom of the prime number 2, two atoms of the prime number 3, and one atom of the prime number 5.

This fact can be used to find the lcm of a set of numbers.

Example: lcm(8,9,21)

Factor each number and express it as a product of prime number powers.

\(\begin{aligned}
8 & = 2^3 \\
9 & = 3^2 \\
21 & = 3^1 \cdot 7^1
\end{aligned}\)

The lcm will be the product of multiplying the highest power of each prime number together. The highest power of the three prime numbers 2, 3, and 7 is 2\(^{3}\), 3\(^{2}\), and 7\(^{1}\), respectively. Thus,

\(\operatorname{lcm}(8,9,21) = 2^3 \cdot 3^2 \cdot 7^1 = 8 \cdot 9 \cdot 7 = 504.\)

This method is not as efficient as reducing to the greatest common divisor; since there is no known general efficient algorithm for integer factorization.

The same method can also be illustrated with a Venn diagram as follows, with the prime factorization of each of the two numbers demonstrated in each circle and *all* factors they share in common in the intersection. The lcm then can be found by multiplying all of the prime numbers in the diagram.

Here is an example:

48 = 2 × 2 × 2 × 2 × 3,
180 = 2 × 2 × 3 × 3 × 5,

sharing two "2"s and a "3" in common:

Least common multiple = 2 × 2 × 2 × 2 × 3 × 3 × 5 = 720
Greatest common divisor = 2 × 2 × 3 = 12
Product = 2 × 2 × 2 × 2 × 3 × 2 × 2 × 3 × 3 × 5 = 8640

This also works for the greatest common divisor (gcd), except that instead of multiplying all of the numbers in the Venn diagram, one multiplies only the prime factors that are in the intersection. Thus the gcd of 48 and 180 is 2 × 2 × 3 = 12.

## Formulas
### Fundamental theorem of arithmetic
According to the fundamental theorem of arithmetic, every integer greater than 1 can be represented uniquely as a product of prime numbers, up to the order of the factors:

\(n = 2^{n_2} 3^{n_3} 5^{n_5} 7^{n_7} \cdots = \prod_p p^{n_p},\)

where the exponents *n*\(_{2}\), *n*\(_{3}\),... are non-negative integers; for example, 84 = 2\(^{2}\) 3\(^{1}\) 5\(^{0}\) 7\(^{1}\) 11\(^{0}\) 13\(^{0}\)... Given two positive integers \(a = \prod_p p^{a_p}\) and \(b = \prod_p p^{b_p}\), their greatest common divisor and least common multiple are given by the formulas
\(\gcd(a,b) = \prod_p p^{\min(a_p, b_p)}\)

and
\(\operatorname{lcm}(a,b) = \prod_p p^{\max(a_p, b_p)}.\)

Since
\(\min(x,y) + \max(x,y) = x + y,\)
this gives
\(\gcd(a,b) \operatorname{lcm}(a,b) = ab.\)

In fact, every rational number can be written uniquely as the product of primes, if negative exponents are allowed. When this is done, the above formulas remain valid. For example:
\(\begin{aligned}
 4 &= 2^2 3^0, & 6 &= 2^1 3^1, & \gcd(4, 6) &= 2^1 3^0 = 2, & \operatorname{lcm}(4,6) &= 2^2 3^1 = 12. \\[8pt]
 \tfrac{1}{3} &= 2^0 3^{-1} 5^0, & \tfrac{2}{5} &= 2^1 3^0 5^{-1}, & \gcd\left(\tfrac13, \tfrac{2}{5}\right) &= 2^0 3^{-1} 5^{-1} = \tfrac{1}{15}, & \operatorname{lcm}\left(\tfrac{1}{3}, \tfrac{2}{5}\right) &= 2^1 3^0 5^0 = 2, \\[8pt]
 \tfrac{1}{6} &= 2^{-1} 3^{-1}, & \tfrac{3}{4} &= 2^{-2} 3^1, & \gcd\left(\tfrac{1}{6}, \tfrac{3}{4}\right) &= 2^{-2} 3^{-1} = \tfrac{1}{12}, & \operatorname{lcm}\left(\tfrac{1}{6}, \tfrac{3}{4}\right) &= 2^{-1} 3^1 = \tfrac{3}{2}.
\end{aligned}\)

### Lattice-theoretic
The positive integers may be partially ordered by divisibility: if *a* divides *b* (that is, if *b* is an integer multiple of *a*) write *a* ≤ *b* (or equivalently, *b* ≥ *a*). (Note that the usual magnitude-based definition of ≤ is not used here.)

Under this ordering, the positive integers become a lattice, with meet given by the gcd and join given by the lcm. The proof is straightforward, if a bit tedious; it amounts to checking that lcm and gcd satisfy the axioms for meet and join. Putting the lcm and gcd into this more general context establishes a duality between them:

*If a formula involving integer variables, gcd, lcm, ≤ and ≥ is true, then the formula obtained by switching gcd with lcm and switching ≥ with ≤ is also true.* (Remember ≤ is defined as divides).

The following pairs of dual formulas are special cases of general lattice-theoretic identities.

**Commutative laws**
\(\operatorname{lcm}(a, b) = \operatorname{lcm}(b, a),\)
\(\gcd(a, b) =\gcd( b, a).\)

**Associative laws**
\(\operatorname{lcm}(a,\operatorname{lcm}(b, c)) = \operatorname{lcm}(\operatorname{lcm}( b),c),\)
\(\gcd(a, \gcd(b, c)) = \gcd(\gcd(a,b), c).\)

**Absorption laws:**
\(\operatorname{lcm}(a, \gcd(a,b)) = a,\)
\(\gcd(a, \operatorname{lcm}(a, b)) = a.\)

**Idempotent laws**
\(\operatorname{lcm}(a, a) = a,\)
\(\gcd(a, a) = a.\)

**Define divides in terms of lcm and gcd**
\(a \ge b \iff a = \operatorname{lcm}(a,b),\)
\(a \le b \iff a = \gcd(a,b).\)

It can also be shown that this lattice is distributive; that is, lcm distributes over gcd and gcd distributes over lcm:

\(\operatorname{lcm}(a,\gcd(b,c)) = \gcd(\operatorname{lcm}(a,b),\operatorname{lcm}(a,c)),\)
\(\gcd(a,\operatorname{lcm}(b,c)) = \operatorname{lcm}(\gcd(a,b),\gcd(a,c)).\)

This identity is self-dual:
\(\gcd(\operatorname{lcm}(a,b),\operatorname{lcm}(b,c),\operatorname{lcm}(a,c))=\operatorname{lcm}(\gcd(a,b),\gcd(b,c),\gcd(a,c)).\)

### Other
* Let *D* be the product of *ω*(*D*) distinct prime numbers (that is, *D* is squarefree).

Then

\(|\{(x,y) \;:\; \operatorname{lcm}(x,y) = D\}| = 3^{\omega(D)},\)

where the absolute bars || denote the cardinality of a set.

* If none of \(a_1, a_2, \ldots , a_r\) is zero, then

\(\operatorname{lcm}(a_1, a_2, \ldots , a_r) = \operatorname{lcm}(\operatorname{lcm}(a_1, a_2, \ldots , a_{r-1}), a_r).\)

## In commutative rings
The least common multiple can be defined generally over commutative rings as follows:

Let and be elements of a commutative ring. A *common multiple* of and is an element of such that both and divide (that is, there exist elements and of such that *ax* and *by* ). A *least common multiple* of and is a common multiple that is minimal, in the sense that for any other common multiple of and , divides. In general, two elements in a commutative ring can have no least common multiple or more than one. However, any two least common multiples of the same pair of elements are associates. In a unique factorization domain, any two elements have a least common multiple. In a principal ideal domain, the least common multiple of and can be characterised as a generator of the intersection of the ideals generated by and (the intersection of a collection of ideals is always an ideal).
