> Content sourced from [Wikipedia](https://en.wikipedia.org/wiki/Remainder) — CC BY-SA 4.0

# Division with remainders

In mathematics, the **remainder** is the amount "left over" after performing some computation. In arithmetic, the remainder is the integer "left over" after dividing one integer by another to produce an integer quotient (integer division). In algebra of polynomials, the remainder is the polynomial "left over" after dividing one polynomial by another. The *modulo operation* is the operation that produces such a remainder when given a dividend and divisor.

Alternatively, a remainder is also what is left after subtracting one number from another, although this is more precisely called the *difference*. This usage can be found in some elementary textbooks; colloquially it is replaced by the expression "the rest" as in "Give me two dollars back and keep the rest." However, the term "remainder" is still used in this sense when a function is approximated by a series expansion, where the error expression ("the rest") is referred to as the remainder term.

## Integer division
Given an integer *a* and a non-zero integer *d*, it can be shown that there exist unique integers *q* and *r*, such that 1=*a* = *qd* + *r* and 1=0 ≤ *r* < |*d*|. The number *q* is called the *quotient*, while *r* is called the *remainder*.

(For a proof of this result, see *Euclidean division*. For algorithms describing how to calculate the remainder, see *Division algorithm*.)

The remainder, as defined above, is called the *least positive remainder* or simply the *remainder*.

In some occasions, it is convenient to carry out the division so that *a* is as close to an integral multiple of *d* as possible, that is, we can write
*a* = *kd* + *s*, with |*s*| ≤ |*d*/2| for some integer *k*.
In this case, *s* is called the *least absolute remainder*. As with the quotient and remainder, *k* and *s* are uniquely determined, except in the case where 1=*d* = 2*n* and 1=*s* = ±*n*. For this exception, we have:
*a* = *kd* + *n* = (*k* + 1)*d* − *n*.
A unique remainder can be obtained in this case by some convention—such as always taking the positive value of *s*.

### Examples
In the division of 43 by 5, we have:
43 = 8 × 5 + 3,
so 3 is the least positive remainder. We also have that:
43 = 9 × 5 − 2,
and −2 is the least absolute remainder.

These definitions are also valid if *d* is negative, for example, in the division of 43 by −5,
43 = (−8) × (−5) + 3,
and 3 is the least positive remainder, while,
43 = (−9) × (−5) + (−2)
and −2 is the least absolute remainder.

In the division of 42 by 5, we have:
42 = 8 × 5 + 2,
and since 2 < 5/2, 2 is both the least positive remainder and the least absolute remainder.

In these examples, the (negative) least absolute remainder is obtained from the least positive remainder by subtracting 5, which is *d*. This holds in general. When dividing by *d*, either both remainders are positive and therefore equal, or they have opposite signs. If the positive remainder is *r*\(_{1}\), and the negative one is *r*\(_{2}\), then
*r*\(_{1}\) = *r*\(_{2}\) + *d*.

## For floating-point numbers
When *a* and *d* are floating-point numbers, with *d* non-zero, *a* can be divided by *d* without remainder, with the quotient being another floating-point number. If the quotient is constrained to being an integer, however, the concept of remainder is still necessary. It can be proved that there exists a unique integer quotient *q* and a unique floating-point remainder *r* such that 1=*a* = *qd* + *r* with 0 ≤ *r* < |*d*|.

Extending the definition of remainder for floating-point numbers, as described above, is not of theoretical importance in mathematics; however, many programming languages implement this definition (see *Modulo operation*).

## Methods for Finding Remainders
* **Long Division:** A traditional method for dividing numbers manually, which involves repeated subtraction and digit placement.
*# Write the division problem in the long division format, with the divisor outside the division symbol and the dividend inside.
*# Determine how many times the divisor can go into the first digit or set of digits of the dividend.
*# Divide the selected digits by the divisor and write the quotient above the division symbol.
*# Multiply the quotient by the divisor and write the result underneath the selected digits.
*# Subtract the result from step 4 from the selected digits to find the remainder.
*# If there are more digits in the dividend, bring down the next digit.
*# Repeat steps 2 through 5 until there are no more digits in the dividend to bring down.
*# If the remainder is less than the divisor and there are no more digits to bring down, stop.
*# If there are still digits to bring down but the remainder is not less than the divisor, continue.
*# The final remainder is the result when all digits have been processed and no more digits can be brought down.
* **Modular Arithmetic:** Utilizing modular arithmetic concepts to calculate remainders efficiently, particularly useful in computer science and cryptography.
*# Start with a dividend and a divisor.
*# Divide the dividend by the divisor.
*# Take the remainder obtained from the division.
*# The remainder is the result of the modular arithmetic operation.
* **Remainder Theorem:** A mathematical theorem that provides a systematic approach to finding remainders when dividing polynomials.
*# Identify the polynomial P(x) and the linear divisor x-a.
*# Evaluate P(a), where a is the root of the linear divisor (the value that makes x-a equal to zero).
*# The value obtained in step 2 is the remainder when P(x) is divided by x-a.

## In programming languages

While there are no difficulties inherent in the definitions, there are implementation issues that arise when negative numbers are involved in calculating remainders. Different programming languages have adopted different conventions. For example:
* Pascal chooses the result of the *mod* operation positive, but does not allow *d* to be negative or zero (so, 1=*a* = (*a* div *d* ) × *d* + *a* mod *d* is not always valid).
* C (since C99) chooses the remainder with the same sign as the dividend *a*. (Earlier versions of the C language allowed other choices.)
* Perl, Python (only modern versions) choose the remainder with the same sign as the divisor *d*.
* Scheme offers two functions, *remainder* and *modulo* – Ada, PL/I and CSS have *mod* and *rem*, while Fortran has *mod* and *modulo*; in each case, the former agrees in sign with the dividend *a*, and the latter with the divisor *d*. Common Lisp and Haskell also have *mod* and *rem*, but *mod* uses the sign of the divisor and *rem* uses the sign of the dividend.

## Polynomial division

Euclidean division of polynomials is very similar to Euclidean division of integers and leads to polynomial remainders. Its existence is based on the following theorem: Given two univariate polynomials *a*(*x*) and *b*(*x*) (where *b*(*x*) is a non-zero polynomial) defined over a field (in particular, the reals or complex numbers), there exist two polynomials *q*(*x*) (the *quotient*) and *r*(*x*) (the *remainder*) which satisfy:
*a*(*x*) = *b*(*x*)*q*(*x*) + *r*(*x*)
where
deg(*r*(*x*)) < deg(*b*(*x*)),
where "deg(...)" denotes the degree of the polynomial (the degree of the constant polynomial whose value is always 0 can be defined to be negative, so that this degree condition will always be valid when this is the remainder). Moreover, *q*(*x*) and *r*(*x*) are uniquely determined by these relations.

This differs from the Euclidean division of integers in that, for the integers, the degree condition is replaced by the bounds on the remainder *r* (non-negative and less than the divisor, which insures that *r* is unique.) The similarity between Euclidean division for integers and that for polynomials motivates the search for the most general algebraic setting in which Euclidean division is valid. The rings for which such a theorem exists are called Euclidean domains, but in this generality, uniqueness of the quotient and remainder is not guaranteed.

Polynomial division leads to a result known as the polynomial remainder theorem: If a polynomial *f*(*x*) is divided by *x* − *k*, the remainder is the constant 1=*r* = *f*(*k*).
