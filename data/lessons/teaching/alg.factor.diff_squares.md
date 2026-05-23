> Content sourced from [ORCCA](https://pcc.edu/orcca) (Open Resources for Community College Algebra) — CC BY 4.0
# Factoring Special Polynomials

Certain polynomials have patterns that you can train yourself to recognize.
      And when they have these patterns,
      there are formulas you can use to factor them,
      much more quickly than using the techniques from Section and Section .

## Difference of Squares

If $b$ is some positive integer,
      then when you multiply $(x-b)(x+b)$ : (x-b)(x+b)\amp=x^2-bx+bx-b^2 \amp=x^2-b^2 .
      The $-bx$ and the $+bx$ cancel each other out.
      So this is telling us that \[
        x^2-b^2=(x-b)(x+b)
      \] .
      And so if we ever encounter a polynomial of the form $x^2-b^2$ (a "difference of squares" ) then we have a quick formula for factoring it.
      Just identify what "b" is,
      and use that in $(x-b)(x+b)$ .

To use this formula,
      it's important to recognize which numbers are perfect squares,
      as in Figure .

**Example**
Factor $x^2-16$ .

*Solution*
The "16" being subtracted here is a perfect square.
          It is the same as $4^2$ .
          So we can take $b=4$ and write: x^2-16\amp=(x-b)(x+b) \amp=(x-4)(x+4)

Try to factor one yourself:

We can do a little better.
      There is nothing special about starting with "x^2" in these examples.
      In full generality:

**The Difference of Squares Formula**
If $A$ and $B$ are any algebraic expressions, then: \[
            A^2 - B^2 = (A-B)(A+B)
          \] .

**Example**
Factor $1-p^2$ .

*Solution*
The "1" at the beginning of this expression is a perfect square;
          it's the same as $1^2$ .
          The "p^2" being subtracted here is also perfect square.
          We can take $A=1$ and $B=p$ ,
          and use : 1-p^2\amp=(A-B)(A+B) \amp=(1-p)(1+p)

**Example**
Factor $m^2n^2-4$ .

*Solution*
Is the "m^2n^2" at the beginning of this expression a perfect square?
          By the properties of exponents, it is the same as $(mn)^2$ , so yes,
          it is a perfect square and we may take $A=mn$ .
          The "4" being subtracted here is also perfect square.
          We can take $B=2$ . tells us: m^2n^2-4\amp=(A-B)(A+B) \amp=(mn-2)(mn+2)

Try to factor one yourself:

**Example**
Factor $x^6-9$ .

*Solution*
Is the "x^6" at the beginning of this expression is a perfect square?
          It may appear to be a *sixth* power,
          but it is *also* a perfect square because we can write $x^6=\left(x^3\right)^2$ .
          So we may take $A=x^3$ .
          The "9" being subtracted here is also perfect square.
          We can take $B=3$ . tells us: x^6-9\amp=(A-B)(A+B) \amp=(x^3-3)(x^3+3)

> It's a common mistake to write something like $x^2+16=(x+4)(x-4)$ .
        This is not what allows you to do,
        and this is in fact incorrect.
        The issue is that $x^2+16$ is a *sum* of squares,
        not a *difference* .
        And it happens that $x^2+16$ is actually prime.
        In fact, any sum of squares without a common factor will always be prime.

## Perfect Square Trinomials

If we expand $(A+B)^2$ : (A+B)^2\amp=(A+B)(A+B) \amp=A^2+BA+AB+B^2 \amp=A^2+2AB+B^2 .
      The $BA$ and the $AB$ equal each other and double up when added together.
      So this is telling us that \[
        A^2+2AB+B^2=(A+B)^2
      \] .
      And so if we ever encounter a polynomial of the form $A^2+2AB+B^2$ (a "perfect square trinomial" ) then we have a quick formula for factoring it.

The tricky part is recognizing when a trinomial you have encountered is in this special form.
      Ask yourself: Are the first and last terms perfect square?
            If so, jot down what $A$ and $B$ would be. When you multiply $2$ with what you wrote down for $A$ and $B$ ,
            i.e. $2AB$ , do you have the middle term?
            If you have this middle term exactly,
            then your polynomial factors as $(A+B)^2$ .
            If the middle term is the negative of $2AB$ ,
            then the sign on your $B$ can be reversed,
            and your polynomial factors as $(A-B)^2$ .

**The Perfect Square Trinomial Formula**
If $A$ and $B$ are any algebraic expressions, then: A^2+2AB+B^2 \amp= (A+B)^2\amp\amp\text{and}\amp A^2-2AB+B^2 \amp= (A-B)^2

**Example**
Factor $x^2+6x+9$ .

*Solution*
The first term, $x^2$ , is clearly a perfect square.
          So we could take $A=x$ .
          The last term, $9$ ,
          is also a perfect square since it is equal to $3^2$ .
          So we could take $B=3$ .
          Now we multiply $2AB=2\cdot x\cdot3$ ,
          and the result is $6x$ .
          This is the middle term, which is what we hope to see.
So we can use : x^2+6x+9\amp=(A+B)^2 \amp=(x+3)^2

**Example**
Factor $4x^2-20xy+25y^2$ .

*Solution*
The first term, $4x^2$ ,
          is a perfect square because it equals $(2x)^2$ .
          So we could take $A=2x$ .
          The last term, $25y^2$ ,
          is also a perfect square since it is equal to $(5y)^2$ .
          So we could take $B=5y$ .
          Now we multiply $2AB=2\cdot (2x)\cdot(5y)$ ,
          and the result is $20xy$ .
          This is the *negative* of the middle term, which we can work with.
          The factored form will be $(A-B)^2$ instead of $(A+B)^2$ .
So we can use : 4x^2-20xy+25y^2\amp=(A-B)^2 \amp=(2x-5y)^2

Try to factor one yourself:

> It is not enough to just see that the first and last terms are perfect squares.
        For example, $9x^2+10x+25$ has its first term equal to $(3x)^2$ and its last term equal to $5^2$ .
        But when you examine $2\cdot(3x)\cdot5$ the result is $30x$ , *not* equal to the middle term.
        So doesn't apply here.
        In fact, this polynomial doesn't factor at all.

## Factoring in Stages

Sometimes factoring a polynomial will take two or more "stages." You might use one of the special patters to factor something into two factors,
      and then those factors might factor even more.
      When the task is to *factor* a polynomial,
      the intention is that you *fully* factor it,
      breaking down the pieces into even smaller pieces when that is possible.

**Factor out any greatest common factor**
Factor $12z^3-27z$ .

*Solution*
The two terms of this polynomial have greatest common factor $3z$ ,
          so the first step in factoring should be to factor this out: \[
            3z\left(4z^2-9\right)
          \] .
          Now we have two factors.
          There is nothing for us to do with $3z$ ,
          but we should ask if $\left(4z^2-9\right)$ can factor further.
          And in fact, that is a difference of squares.
          So we can apply .
          The full process would be: 12z^3-27z\amp=3z\left(4z^2-9\right) \amp=3z(2z-3)(2z+3)

**Recognize a second special pattern**
Factor $p^4-1$ .

*Solution*
Since $p^4$ is the same as $\left(p^2\right)^2$ ,
          we have a difference of squares here.
          We can apply : p^4-1\amp=\left(p^2-1\right)\left(p^2+1\right) It doesn't end here. Of the two factors we found, $\left(p^2+1\right)$ cannot be factored further. But the other one, $\left(p^2-1\right)$ is *also* a difference of squares. So we should apply again: \phantom{p^4-1}\amp=(p-1)(p+1)\left(p^2+1\right)

**Example**
Factor $32x^6y^2-48x^5y+18x^4$ .

*Solution*
The first step of factoring any polynomial is to factor out the common factor if possible.
          For this trinomial, the common factor is $2x^4$ , so we write \[
            32x^6y^2-48x^5y+18x^4=2x^4(16x^2y^2-24xy+9)
          \] .
The square numbers $16$ and $9$ in $16x^2y^2-24xy+9$ hint that maybe we could use .
          Taking $A=4xy$ and $B=3$ ,
          we multiply $2AB=2\cdot(4xy)\cdot 3$ .
          The result is $24xy$ , which is the negative of our middle term.
          So the whole process is: 32x^6y^2-48x^5y+18x^4\amp=2x^4(16x^2y^2-24xy+9) \amp=2x^4(4xy-3)^2
