> Content sourced from [Wikipedia](https://en.wikipedia.org/wiki/Division_%28mathematics%29) — CC BY-SA 4.0

# Basic division without remainder

**Division** is one of the four basic operations of arithmetic. The other operations are addition, subtraction, and multiplication. What is being divided is called the ***dividend***, which is divided by the ***divisor***, and the result is called the *quotient*.

At an elementary level the division of two natural numbers is, among other possible interpretations, the process of calculating the number of times one number is contained within another. For example, if 20 apples are divided evenly between 4 people, everyone receives 5 apples (see picture). However, this number of times or the number contained (divisor) need not be integers.

The division with remainder or Euclidean division of two natural numbers provides an integer *quotient*, which is the number of times the second number is completely contained in the first number, and a *remainder*, which is the part of the first number that remains, when in the course of computing the quotient, no further full chunk of the size of the second number can be allocated. For example, if 21 apples are divided between 4 people, everyone receives 5 apples again, and 1 apple remains.

For division to always yield one number rather than an integer quotient plus a remainder, the natural numbers must be extended to rational numbers or real numbers. In these enlarged number systems, division is the inverse operation to multiplication, that is *a* means *a* × *b* , as long as *b* is not zero. If 1=*b* = 0, then this is a division by zero, which is not defined. 1.In the 21-apples example, everyone would receive 5 apple and a quarter of an apple, thus avoiding any leftover.

Both forms of division appear in various algebraic structures, different ways of defining mathematical structure. Those in which a Euclidean division (with remainder) is defined are called Euclidean domains and include polynomial rings in one indeterminate (which define multiplication and addition over single-variabled formulas). Those in which a division (with a single result) by all nonzero elements is defined are called fields and division rings. In a ring the elements by which division is always possible are called the units (for example, 1 and −1 in the ring of integers). Another generalization of division to algebraic structures is the quotient group, in which the result of "division" is a group rather than a number.

## Introduction
The simplest way of viewing division is in terms of quotition and partition: from the quotition perspective, 20 / 5 means the number of 5s that must be added to get 20. In terms of partition, 20 / 5 means the size of each of 5 parts into which a set of size 20 is divided. For example, 20 apples divide into five groups of four apples, meaning that "twenty divided by five is equal to four". This is denoted as 1=20 / 5 = 4, or 1=. In the example, 20 is the dividend, 5 is the divisor, and 4 is the quotient.

Unlike the other basic operations, when dividing natural numbers there is sometimes a remainder that will not go evenly into the dividend. For example, 10 / 3 leaves a remainder of 1, as 10 is not a multiple of 3. Sometimes this remainder is added to the quotient as a fractional part, so 10 / 3 is equal to or 3.33..., but in the context of integer division, where numbers have no fractional part, the remainder is kept separately (or exceptionally, discarded or rounded). When the remainder is kept as a fraction, it leads to a rational number. The set of all rational numbers is created by extending the integers with all possible results of divisions of integers.

Unlike multiplication and addition, division is not commutative, meaning that *a* / *b* is not always equal to *b* / *a*. Division is also not, in general, associative, meaning that when dividing multiple times, the order of division can change the result. For example, (24 / 6) / 2 , but 24 / (6 / 2) (where the use of parentheses indicates that the operations inside parentheses are performed before the operations outside parentheses).

Division is traditionally considered as left-associative. That is, if there are multiple divisions in a row, the order of calculation goes from left to right:
\(a / b / c = (a / b) / c = a / (b \times c) \;\ne\; a/(b/c)= (a\times c)/b.\)

Division is right-distributive over addition and subtraction, in the sense that
\(\frac{a \pm b}{c} = (a \pm b) / c = (a/c)\pm (b/c) =\frac{a}{c} \pm \frac{b}{c}.\)

This is the same for multiplication, as \((a + b) \times c = a \times c + b \times c\). However, division is *not* left-distributive, as
\(\frac{a}{b + c} = a / (b + c) \;\ne\; (a/b) + (a/c) = \frac{ac+ab}{bc}.\)   For example \(\frac{12}{2+4} = \frac{12}{6} = 2 ,\) but \(\frac{12}{2} + \frac{12}{4} = 6+3 = 9. \)
This is unlike the case in multiplication, which is both left-distributive and right-distributive, and thus distributive.

## Notation

Division is often shown in algebra and science by placing the *dividend* over the *divisor* with a horizontal line, also called a fraction bar, between them. For example, "*a* divided by *b*" can be written as:
\(\frac ab\)

which can also be read out loud as "divide *a* by *b*" or "*a* over *b*". A way to express division all on one line is to write the *dividend* (or numerator), then a slash, then the *divisor* (or denominator), as follows:
\(a/b\)

This is the usual way of specifying division in most computer programming languages; since it can easily be typed as a simple sequence of ASCII characters. (It is also the only notation used for quotient objects in abstract algebra.) Some mathematical software, such as MATLAB and GNU Octave, allows the operands to be written in the reverse order by using the backslash as the division operator:
\(b\backslash a\)

A typographical variation halfway between these two forms uses a solidus (fraction slash), but elevates the dividend and lowers the divisor:
\({}^{a}\!/{}_{b}\)

Any of these forms can be used to display a fraction. A fraction is a division expression where both dividend and divisor are integers (typically called the *numerator* and *denominator*), and there is no implication that the division must be evaluated further. A second way to show division is to use the division sign (÷, also known as obelus though the term has additional meanings), common in arithmetic, in this manner:
\(a \div b\)

This form is infrequent except in elementary arithmetic. ISO 80000-2-10.6 states it should not be used. This division sign is also used alone to represent the division operation itself, as for instance as a label on a key of a calculator. The obelus was introduced by Swiss mathematician Johann Rahn in 1659 in *Teutsche Algebra*. The ÷ symbol is used to indicate subtraction in some European countries, so its use may be misunderstood.

In some non-English-speaking countries, a colon is used to denote division:

\( b\)

This notation was introduced by Gottfried Wilhelm Leibniz in his 1684 *Acta eruditorum*. Leibniz disliked having separate symbols for ratio and division. However, in English usage the colon is restricted to expressing the related concept of ratios.

Since the 19th century, US textbooks have used \(b)a\) or \(b \overline{)a}\) to denote *a* divided by *b*, especially when discussing long division. The history of this notation is not entirely clear because it evolved over time.

## Computing

### Manual methods
Division is often introduced through the notion of "sharing out" a set of objects, for example a pile of lollies, into a number of equal portions. Distributing the objects several at a time in each round of sharing to each portion leads to the idea of 'chunking' a form of division where one repeatedly subtracts multiples of the divisor from the dividend itself.

By allowing one to subtract more multiples than what the partial remainder allows at a given stage, more flexible methods, such as the bidirectional variant of chunking, can be developed as well.

More systematically and more efficiently, two integers can be divided with pencil and paper with the method of short division, if the divisor is small, or long division, if the divisor is larger. If the dividend has a fractional part (expressed as a decimal fraction), one can continue the procedure past the ones place as far as desired. If the divisor has a fractional part, one can restate the problem by moving the decimal to the right in both numbers until the divisor has no fraction, which can make the problem easier to solve (e.g., 10/2.5 = 100/25 = 4).

Division can be calculated with an abacus.

Logarithm tables can be used to divide two numbers, by subtracting the two numbers' logarithms, then looking up the antilogarithm of the result.

Division can be calculated with a slide rule by aligning the divisor on the C scale with the dividend on the D scale. The quotient can be found on the D scale where it is aligned with the left index on the C scale. The user is responsible, however, for mentally keeping track of the decimal point.

### By computer
Modern calculators and computers compute division either by methods similar to long division, or by faster methods; see Division algorithm.

In modular arithmetic (modulo a prime number) and for real numbers, nonzero numbers have a multiplicative inverse. In these cases, a division by may be computed as the product by the multiplicative inverse of. This approach is often associated with the faster methods in computer arithmetic.

## Division in different contexts
### Euclidean division

Euclidean division is the mathematical formulation of the outcome of the usual process of division of integers. It asserts that, given two integers, *a*, the *dividend*, and *b*, the *divisor*, such that *b* ≠ 0, there are unique integers *q*, the *quotient*, and *r*, the remainder, such that *a* = *bq* + *r* and 0 ≤ *r* < |*b*|, where |*b*| denotes the absolute value of *b*.

### Of integers
Integers are not closed under division. Apart from division by zero being undefined, the quotient is not an integer unless the dividend is an integer multiple of the divisor. For example, 26 cannot be divided by 11 to give an integer. Such a case uses one of five approaches:
# Say that 26 cannot be divided by 11; division becomes a partial function.
# Give an approximate answer as a floating-point number. This is the approach usually taken in numerical computation.
# Give the answer as a fraction representing a rational number, so the result of the division of 26 by 11 is \(\tfrac{26}{11}\) (or as a mixed number, so \(\tfrac{26}{11} = 2 \tfrac 4{11}.\)) Usually the resulting fraction should be simplified: the result of the division of 52 by 22 is also \(\tfrac{26}{11}\). This simplification may be done by factoring out the greatest common divisor.
# Give the answer as an integer *quotient* and a *remainder*, so \(\tfrac{26}{11} = 2 \mbox{ remainder } 4.\) To make the distinction with the previous case, this division, with two integers as result, is sometimes called *Euclidean division*, because it is the basis of the Euclidean algorithm.
# Give the integer quotient as the answer, so \(\tfrac{26}{11} = 2.\) This is the *floor function* applied to case 2 or 3. It is sometimes called **integer division**, and denoted by "//".

Dividing integers in a computer program requires special care. Some programming languages treat integer division as in case 5 above, so the answer is an integer. Other languages, such as MATLAB and every computer algebra system return a rational number as the answer, as in case 3 above. These languages also provide functions to get the results of the other cases, either directly or from the result of case 3.

Names and symbols used for integer division include , and. Definitions vary regarding integer division when the dividend or the divisor is negative: rounding may be toward zero (so called T-division) or toward −∞ (F-division); rarer styles can occur – see modulo operation for the details.

Divisibility rules can sometimes be used to quickly determine whether one integer divides exactly into another.

### Of rational numbers
The result of dividing two rational numbers is another rational number when the divisor is not 0. The division of two rational numbers *p*/*q* and *r*/*s* can be computed as

\[
{p/q \over r/s} = {p \over q} \times {s \over r} = {ps \over qr}.
\]

All four quantities are integers, and only *p* may be 0. This definition ensures that division is the inverse operation of multiplication.

### Of real numbers
Division of two real numbers results in another real number (when the divisor is nonzero). It is defined such that *a*/*b* = *c* if and only if *a* = *cb* and *b* ≠ 0.

### Of complex numbers
Dividing two complex numbers (when the divisor is nonzero) results in another complex number, which is found using the conjugate of the denominator:

\[
{p+iq \over r+is} = {(p+iq)(r-is) \over (r+is)(r-is)} = {pr+qs + i(qr-ps) \over r^2+s^2} = {pr+qs \over r^2+s^2} + i{qr-ps \over r^2+s^2}.
\]

This process of multiplying and dividing by \(r-is\) is called 'realisation' or (by analogy) rationalisation. All four quantities *p*, *q*, *r*, *s* are real numbers, and *r* and *s* may not both be 0.

Division for complex numbers expressed in polar form is simpler than the definition above:

\[
{p e^{iq} \over r e^{is= {p e^{iq} e^{-is} \over r e^{is} e^{-is= {p \over r}e^{i(q - s)}.
\]

Again all four quantities *p*, *q*, *r*, *s* are real numbers, and *r* may not be 0.

### Of polynomials
One can define the division operation for polynomials in one variable over a field. Then, as in the case of integers, one has a remainder. See Euclidean division of polynomials, and, for hand-written computation, polynomial long division or synthetic division.

### Of matrices
One can define a division operation for matrices. The usual way to do this is to define 1=*A* / *B* = *AB*\(^{−1}\), where *B*\(^{−1}\) denotes the inverse of *B*, but it is far more common to write out *AB*\(^{−1}\) explicitly to avoid confusion. An elementwise division can also be defined in terms of the Hadamard product.

#### Left and right division
Because matrix multiplication is not commutative, one can also define a left division or so-called *backslash-division* as 1=*A* \ *B* = *A*\(^{−1}\)*B*. For this to be well defined, *B*\(^{−1}\) need not exist, however *A*\(^{−1}\) does need to exist. To avoid confusion, division as defined by 1=*A* / *B* = *AB*\(^{−1}\) is sometimes called *right division* or *slash-division* in this context.

With left and right division defined this way, *A* / (*BC*) is in general not the same as (*A* / *B*) / *C*, nor is (*AB*) \ *C* the same as *A* \ (*B* \ *C*). However, it holds that 1=*A* / (*BC*) = (*A* / *C*) / *B* and 1=(*AB*) \ *C* = *B* \ (*A* \ *C*).

#### Pseudoinverse
To avoid problems when *A*\(^{−1}\) and/or *B*\(^{−1}\) do not exist, division can also be defined as multiplication by the pseudoinverse. That is, 1=*A* / *B* = *AB*\(^{+}\) and 1=*A* \ *B* = *A*\(^{+}\)*B*, where *A*\(^{+}\) and *B*\(^{+}\) denote the pseudoinverses of and. ### Abstract algebra
In abstract algebra, given a magma with binary operation ∗ (which could nominally be termed multiplication), **left division** of *b* by *a* (written *a* \ *b*) is typically defined as the solution *x* to the equation 1=*a* ∗ *x* = *b*, if this exists and is unique. Similarly, **right division** of *b* by *a* (written *b* / *a*) is the solution *y* to the equation 1=*y* ∗ *a* = *b*. Division in this sense does not require ∗ to have any particular properties (such as commutativity, associativity, or an identity element). A magma for which both *a* \ *b* and *b* / *a* exist and are unique for all *a* and all *b* (the Latin square property) is a quasigroup. In a quasigroup, division in this sense is always possible, even without an identity element and hence without inverses.

"Division" in the sense of "cancellation" can be done in any magma by an element with the cancellation property. Examples include matrix algebras, quaternion algebras, and quasigroups. In an integral domain, where not every element need have an inverse, *division* by a cancellative element *a* can still be performed on elements of the form *ab* or *ca* by left or right cancellation, respectively. If a ring is finite and every nonzero element is cancellative, then by an application of the pigeonhole principle, every nonzero element of the ring is invertible, and *division* by any nonzero element is possible. To learn about when *algebras* (in the technical sense) have a division operation, refer to the page on division algebras. In particular Bott periodicity can be used to show that any real normed division algebra must be isomorphic to either the real numbers **R**, the complex numbers **C**, the quaternions **H**, or the octonions **O**.

### Calculus
The derivative of the quotient of two functions is given by the quotient rule:

\[
{\left(\frac fg\right)}' = \frac{f'g - fg'}{g^2}.
\]

## Division by zero

Division of any number by zero in most mathematical systems is undefined, because zero multiplied by any finite number always results in the product being zero. Entry of such an expression into most calculators produces an error message. However, in certain mathematical structures, division by zero is possible, such as in the zero ring and in algebraic structures such as wheels. In these structures, the meaning of division is different from that of traditional definitions.
