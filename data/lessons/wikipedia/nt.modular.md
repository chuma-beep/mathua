> Content sourced from [Wikipedia](https://en.wikipedia.org/wiki/Modular_arithmetic) — CC BY-SA 4.0

# Modular arithmetic

(mod )*" notation|the binary operation *mod()'' |Modulo

In mathematics, **modular arithmetic** is a system of arithmetic operations for integers, differing from the usual ones in that numbers "wrap around" when reaching or exceeding a certain value, called the **modulus**. The modern approach to number theory using modular arithmetic was developed by Carl Friedrich Gauss in his book *Disquisitiones Arithmeticae*, published in 1801.

*Modular arithmetic modulo* consists of systematically replacing the results of additions, multiplications, and subtractions by the remainder of the division by. A remarkable property of modular arithmetic is that the result of a computation does not depend on whether the division by is performed after each operation, only once at the end of the computation, or at the end of the computation and after some intermediate resultstypically when an intermediate result becomes too large.

## Motivating example
A familiar setting exhibiting modular arithmetic is the hour hand on a 12-hour clock. If the hour hand points to 7 now, then 8 hours later it will point to 3. Ordinary addition would result in 7 + 8 , but 15 reads as 3 on the clock face. This is because the hour hand makes one rotation every 12 hours and the hour number starts over when the hour hand passes 12. We say that 15 is *congruent* to 3 modulo 12, and we write 15 ≡ 3 (mod 12), so 7 + 8 ≡ 3 (mod 12).

Similarly, if one waits 8 hours and then 8 more hours (thus 16 hours in total), the clock will show the same time change as if one waited 4 hours. This is reflected by the identity 2 × 8 ≡ 4 (mod 12). After a wait of exactly 12 hours, the hour hand will be right where it started, so 12 acts as 0; one writes 12 ≡ 0 (mod 12).

## Congruence
Given an integer *m* ≥ 1, called a **modulus**, two integers and are said to be **congruent** modulo , if their difference *a* − *b* is an integer multiple of ; that is, if there is an integer *k* such that
*a* − *b* = *km*.
Congruence modulo is a congruence relation, meaning that it is an equivalence relation compatible with addition, subtraction, and multiplication. Congruence modulo is denoted by

\(a \equiv b \pmod m.\)

The parentheses mean that (mod *m*) applies to the entire equation, not just to the right-hand side (here, ).

This notation is not to be confused with the notation *b* mod *m* or (*b* mod *m*) (without parentheses immediately before "mod"), which refers to the remainder of *b* when divided by *m*, known as the modulo operation; that is, *b* mod *m* denotes the unique integer such that 0 ≤ *r* < *m* and *r* ≡ *b* (mod *m*). So, the relation \(a \equiv b \pmod m\) must be read \((a \equiv b) \bmod m,\) and is equivalent with \(a\bmod m = b\bmod m.\)

The congruence relation *a* ≡ *b* (mod *m*) may be rewritten as
\(\exists k\in\Z\quad a=km+b,\)
explicitly showing its relationship with Euclidean division. However, the *b* here need not be the remainder in the division of *a* by *m*. Rather, *a* ≡ *b* (mod *m*) asserts that *a* and *b* have the same remainder when divided by *m*. That is,
*a* = *p m* + *r*,
*b* = *q m* + *r*,
where 0 ≤ *r* < *m* is the common remainder. We recover the previous relation (*a* − *b* = *k m*) by subtracting these two expressions and setting *k* = *p* − *q*.

Because the congruence modulo is defined by the divisibility by and because −1 is a unit in the ring of integers, a number is divisible by −*m* exactly if it is divisible by. This means that every non-zero integer may be taken as a modulus.
### Examples
In modulus 12, one can assert that:
38 ≡ 14 (mod 12)
because the difference is 38 − 14 = 24 = 2 × 12, a multiple of 12. Equivalently, 38 and 14 have the same remainder 2 when divided by 12.

The definition of congruence also applies to negative values. For example:
\(\begin{aligned}
2 &\equiv -3 \pmod 5\\
-8 &\equiv \phantom{+}7 \pmod 5\\
-3 &\equiv -8 \pmod 5.
\end{aligned}\)

## Basic properties

The congruence relation satisfies all the conditions of an equivalence relation:
* Reflexivity: *a* ≡ *a* (mod *m*)
* Symmetry: *a* ≡ *b* (mod *m*) if and only if *b* ≡ *a* (mod *m*).
* Transitivity: If *a* ≡ *b* (mod *m*) and *b* ≡ *c* (mod *m*), then *a* ≡ *c* (mod *m*)

If *a*\(_{1}\) ≡ *b*\(_{1}\) (mod *m*) and *a*\(_{2}\) ≡ *b*\(_{2}\) (mod *m*), or if *a* ≡ *b* (mod *m*), then:
* *a* + *k* ≡ *b* + *k* (mod *m*) for any integer *k* (compatibility with translation)
* *k a* ≡ *k b* (mod *m*) for any integer *k* (compatibility with scaling)
* *k a* ≡ *k b* (mod *k m*) for any integer *k*
* *a*\(_{1}\) + *a*\(_{2}\) ≡ *b*\(_{1}\) + *b*\(_{2}\) (mod *m*) (compatibility with addition)
* *a*\(_{1}\) − *a*\(_{2}\) ≡ *b*\(_{1}\) − *b*\(_{2}\) (mod *m*) (compatibility with subtraction)
* *a*\(_{1}\) *a*\(_{2}\) ≡ *b*\(_{1}\) *b*\(_{2}\) (mod *m*) (compatibility with multiplication)
* *a*\(^{*k*}\) ≡ *b*\(^{*k*}\) (mod *m*) for any non-negative integer *k* (compatibility with exponentiation)
* *p*(*a*) ≡ *p*(*b*) (mod *m*), for any polynomial *p*(*x*) with integer coefficients (compatibility with polynomial evaluation)

If *a* ≡ *b* (mod *m*), then it is generally false that *k\(^{a}\)* ≡ *k\(^{b}\)* (mod *m*). However, the following is true:
* If *c* ≡ *d* (mod *φ*(*m*)), where *φ* is Euler's totient function, then *a*\(^{*c*}\) ≡ *a*\(^{*d*}\) (mod *m*)—provided that *a* is coprime with *m*.

If *a* ≡ *b* (mod *mn*), then *a* ≡ *b* (mod *m*) and *a* ≡ *b* (mod *n*).

For cancellation of common terms, we have the following rules:
* If *a* + *k* ≡ *b* + *k* (mod *m*), where *k* is any integer, then *a* ≡ *b* (mod *m*).
* If *k a* ≡ *k b* (mod *m*) and *k* is coprime with *m*, then *a* ≡ *b* (mod *m*).
* If *k a* ≡ *k b* (mod *k m*) and *k* ≠ 0, then *a* ≡ *b* (mod *m*).

The last rule can be used to move modular arithmetic into division. If *b* divides *a*, then (*a*/*b*) mod *m* = (*a* mod (*b m*)) / *b*.

The modular multiplicative inverse is defined by the following rules:
* Existence: There exists an integer denoted *a*\(^{−1}\) such that *aa*\(^{−1}\) ≡ 1 (mod *m*) if and only if *a* is coprime with *m*. This integer *a*\(^{−1}\) is called a *modular multiplicative inverse* of modulo *m*.
* If *a* ≡ *b* (mod *m*) and *a*\(^{−1}\) exists, then *a*\(^{−1}\) ≡ *b*\(^{−1}\) (mod *m*) (compatibility with multiplicative inverse, and, if *a* = *b*, uniqueness modulo *m*).
* If *ax* ≡ *b* (mod *m*) and *a* is coprime to *m*, then the solution to this linear congruence is given by *x* ≡ *a*\(^{−1}\)*b* (mod *m*).

The multiplicative inverse *x* ≡ *a*\(^{−1}\) (mod *m*) may be efficiently computed by solving Bézout's equation *a x* + *m y* = 1 for *x*, *y*, by using the Extended Euclidean algorithm.

In particular, if *p* is a prime number, then *a* is coprime with *p* for every *a* such that 0 < *a* < *p*; thus a multiplicative inverse exists for all *a* that is not congruent to zero modulo *p*.

## Advanced properties
Some of the more advanced properties of congruence relations are the following:
* Fermat's little theorem: If *p* is prime and does not divide *a*, then *a*.
* Euler's theorem: If *a* and *m* are coprime, then *a*, where *φ* is Euler's totient function.
* A simple consequence of Fermat's little theorem is that if *p* is prime, then *a*\(^{−1}\) ≡ *a* is the multiplicative inverse of 0 < *a* < *p*. More generally, from Euler's theorem, if *a* and *m* are coprime, then *a* (mod *m*). Hence, if *ax* ≡ *1* (mod *m*), then *x* ≡ *a*.
* Another simple consequence is that if *a* ≡ *b* (mod *φ*(*m*)), where *φ* is Euler's totient function, then *k*\(^{*a*}\) ≡ *k*\(^{*b*}\) (mod *m*) provided *k* is coprime with *m*.
* Wilson's theorem: *p* is prime if and only if (*p* − 1)! ≡ −1 (mod *p*).
* Chinese remainder theorem: For any *a*, *b* and coprime *m*, *n*, there exists a unique *x* (mod *mn*) such that *x* ≡ *a* (mod *m*) and *x* ≡ *b* (mod *n*). In fact, *x* ≡ *b m*\(_{*n*}\)\(^{−1}\) *m* + *a n*\(_{*m*}\)\(^{−1}\) *n* (mod *mn*) where *m*\(_{*n*}\)\(^{−1}\) is the inverse of *m* modulo *n* and *n*\(_{*m*}\)\(^{−1}\) is the inverse of *n* modulo *m*.
* Lagrange's theorem: If *p* is prime and *f* (*x*) is a polynomial with integer coefficients such that is not a divisor of *a*\(_{0}\), then the congruence *f* (*x*) ≡ 0 (mod *p*) has at most *d* non-congruent solutions.
* Primitive root modulo *m*: A number *g* is a primitive root modulo *m* if, for every integer *a* coprime to *m*, there is an integer *k* such that *g*\(^{*k*}\) ≡ *a* (mod *m*). A primitive root modulo *m* exists if and only if *m* is equal to 2, 4, *p*\(^{*k*}\) or 2*p*\(^{*k*}\), where *p* is an odd prime number and *k* is a positive integer. If a primitive root modulo *m* exists, then there are exactly *φ*(*φ*(*m*)) such primitive roots, where *φ* is the Euler's totient function.
* Quadratic residue: An integer *a* is a quadratic residue modulo *m*, if there exists an integer *x* such that *x*\(^{2}\) ≡ *a* (mod *m*). Euler's criterion asserts that, if *p* is an odd prime, and is not a multiple of , then *a* is a quadratic residue modulo *p* if and only if
*: *a*.

## Congruence classes
The congruence relation is an equivalence relation. The equivalence class modulo of an integer *a* is the set of all integers of the form *a* + *k m*, where is any integer. It is called the **congruence class** or **residue class** of *a* modulo *m*, and may be denoted (*a* mod *m*), or as or [*a*] when the modulus *m* is known from the context.

Each residue class modulo *m* contains exactly one integer in the range \(0,... , |m| - 1\). Thus, these \(|m|\) integers are representatives of their respective residue classes.

It is generally easier to work with integers than sets of integers; that is, the representatives most often considered, rather than their residue classes.

Consequently, (*a* mod *m*) denotes generally the unique integer such that 0 ≤ *r* < *m* and *r* ≡ *a* (mod *m*); it is called the **residue** of *a* modulo *m*.

In particular, (*a* mod *m*) = (*b* mod *m*) is equivalent to *a* ≡ *b* (mod *m*), and this explains why "=" is often used instead of "≡" in this context.

## Residue systems
Each residue class modulo *m* may be represented by any one of its members, although we usually represent each residue class by the smallest nonnegative integer which belongs to that class (since this is the proper remainder which results from division). Any two members of different residue classes modulo *m* are incongruent modulo *m*. Furthermore, every integer belongs to one and only one residue class modulo *m*.

The set of integers is called the **least residue system modulo *m***. Any set of *m* integers, no two of which are congruent modulo *m*, is called a **complete residue system modulo *m***.

The least residue system is a complete residue system, and a complete residue system is simply a set containing precisely one representative of each residue class modulo *m*. For example, the least residue system modulo 4 Some other complete residue systems modulo 4 include:
*
*
*
*
*
*

Some sets that are *not* complete residue systems modulo 4 are:
* ; since 6 is congruent to 22 modulo 4.
* ; since a complete residue system modulo 4 must have exactly 4 incongruent residue classes.

### Reduced residue systems

Given the Euler's totient function *φ*(*m*), any set of *φ*(*m*) integers that are relatively prime to *m* and mutually incongruent under modulus *m* is called a **reduced residue system modulo *m***. The set from above, for example, is an instance of a reduced residue system modulo 4.

### Covering systems

Covering systems represent yet another type of residue system that may contain residues with varying moduli.

## Integers modulo *m*
In the context of this paragraph, the modulus *m* is almost always taken as positive.

The set of all congruence classes modulo *m* is a ring called the **ring of integers modulo *m***, and is denoted \(\mathbb{Z}/m\mathbb{Z}\), \(\mathbb{Z}/(m)\), \(\mathbb{Z}/m\), or \(\mathbb{Z}_m\). The ring \(\mathbb{Z}/m\mathbb{Z}\) is fundamental to various branches of mathematics (see ** below).
(In some parts of number theory the notation \(\mathbb{Z}_m\) is avoided because it can be confused with the set of *m*-adic integers.)

For *m* > 0 one has
\(\mathbb{Z}/m\mathbb{Z} = \left\{ \overline{a}_m \mid a \in \mathbb{Z}\right\} = \left\{ \overline{0}_m, \overline{1}_m, \overline{2}_m,\ldots, \overline{m{-}1}_m \right\}.\)

When *m* = 1, \(\mathbb{Z}/m\mathbb{Z}\) is the zero ring; when *m* = 0, \(\mathbb{Z}/m\mathbb{Z}\) is not an empty set; rather, it is isomorphic to \(\mathbb{Z}\); since. Addition, subtraction, and multiplication are defined on \(\mathbb{Z}/m\mathbb{Z}\) by the following rules:
* \(\overline{a}_m + \overline{b}_m = \overline{(a + b)}_m\)
* \(\overline{a}_m - \overline{b}_m = \overline{(a - b)}_m\)
* \(\overline{a}_m \overline{b}_m = \overline{(a b)}_m.\)

The properties given before imply that, with these operations, \(\mathbb{Z}/m\mathbb{Z}\) is a commutative ring. For example, in the ring \(\mathbb{Z}/24\mathbb{Z}\), one has
\(\overline{12}_{24} + \overline{21}_{24} = \overline{33}_{24}= \overline{9}_{24}\)
as in the arithmetic for the 24-hour clock.

The notation \(\mathbb{Z}/m\mathbb{Z}\) is used because this ring is the quotient ring of \(\mathbb{Z}\) by the ideal \(m\mathbb{Z}\), the set formed by all multiples of *m*, that is, all numbers *k m* with \(k\in\mathbb{Z}.\)

Under addition, \(\mathbb Z/m\Z\) is a cyclic group. All finite cyclic groups are isomorphic with \(\mathbb Z/m\mathbb Z\) for some. The ring of integers modulo *m* is a field; that is, every nonzero element has a multiplicative inverse, if and only if *m* is prime. If *m* = *p* is a prime power with *k* > 1, there exists a unique (up to isomorphism) finite field \(\mathrm{GF}(m) =\mathbb F_m\) with *m* elements, which is *not* isomorphic to \(\mathbb Z/m\mathbb Z\), which fails to be a field because it has zero-divisors.

If *m* > 1, \((\mathbb Z/m\mathbb Z)^\times\) denotes the multiplicative group of the integers modulo *m* that are invertible. It consists of the congruence classes , where *a* is coprime to *m*; these are precisely the classes possessing a multiplicative inverse. They form an abelian group under multiplication; its order is *φ*(*m*), where is Euler's totient function.

## Applications
In pure mathematics, modular arithmetic is one of the foundations of number theory, touching on almost every aspect of its study, and it is also used extensively in group theory, ring theory, knot theory, and abstract algebra. In applied mathematics, it is used in computer algebra, cryptography, computer science, chemistry and the visual and musical arts.

A very practical application is to calculate checksums within serial number identifiers. For example, International Standard Book Number (ISBN) uses modulo 11 (for 10-digit ISBN) or modulo 10 (for 13-digit ISBN) arithmetic for error detection. Likewise, International Bank Account Numbers (IBANs) use modulo 97 arithmetic to spot user input errors in bank account numbers. In chemistry, the last digit of the CAS registry number (a unique identifying number for each chemical compound) is a check digit, which is calculated by taking the last digit of the first two parts of the CAS registry number times 1, the previous digit times 2, the previous digit times 3 etc., adding all these up and computing the sum modulo 10.

In cryptography, modular arithmetic directly underpins public key systems such as RSA and Diffie–Hellman, and provides finite fields which underlie elliptic curves, and is used in a variety of symmetric key algorithms including Advanced Encryption Standard (AES), International Data Encryption Algorithm (IDEA), and RC4. RSA and Diffie–Hellman use modular exponentiation.

In computer algebra, modular arithmetic is commonly used to limit the size of integer coefficients in intermediate calculations and data. It is used in polynomial factorization, a problem for which all known efficient algorithms use modular arithmetic. It is used by the most efficient implementations of polynomial greatest common divisor, exact linear algebra and Gröbner basis algorithms over the integers and the rational numbers. As posted on Fidonet in the 1980s and archived at Rosetta Code, modular arithmetic was used to disprove Euler's sum of powers conjecture on a Sinclair QL microcomputer using just one-fourth of the integer precision used by a CDC 6600 supercomputer to disprove it two decades earlier via a brute force search.

In computer science, modular arithmetic is often applied in bitwise operations and other operations involving fixed-width, cyclic data structures. The modulo operation, as implemented in many programming languages and calculators, is an application of modular arithmetic that is often used in this context. The logical operator XOR sums 2 bits, modulo 2.

The use of long division to turn a fraction into a repeating decimal in any base *b* is equivalent to modular multiplication of *b* modulo the denominator. For example, for decimal, *b* = 10.

In music, arithmetic modulo 12 is used in the consideration of the system of twelve-tone equal temperament, where octave and enharmonic equivalency occurs (that is, pitches in a 1:2 or 2:1 ratio are equivalent, and C-sharp is considered the same as D-flat).

The method of casting out nines offers a quick check of decimal arithmetic computations performed by hand. It is based on modular arithmetic modulo 9, and specifically on the crucial property that 10 ≡ 1 (mod 9).

Arithmetic modulo 7 is used in algorithms that determine the day of the week for a given date. In particular, Zeller's congruence and the Doomsday algorithm make heavy use of modulo-7 arithmetic.

More generally, modular arithmetic also has application in disciplines such as politics (for example, apportionment), economics (for example, game theory) and other areas of the social sciences, where proportional division and allocation of resources plays a central part of the analysis.

## Computational complexity
Since modular arithmetic has such a wide range of applications, it is important to know how hard it is to solve a system of congruences. A linear system of congruences can be solved in polynomial time with a form of Gaussian elimination, for details see linear congruence theorem. Algorithms, such as Montgomery reduction, also exist to allow simple arithmetic operations, such as multiplication and exponentiation modulo *m*, to be performed efficiently on large numbers.

Some operations, like finding a discrete logarithm or a quadratic congruence appear to be as hard as integer factorization and thus are a starting point for cryptographic algorithms and encryption. These problems might be NP-intermediate.

Solving a system of non-linear modular arithmetic equations is NP-complete.
