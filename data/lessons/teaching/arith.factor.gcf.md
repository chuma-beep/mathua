> Content sourced from [ORCCA](https://pcc.edu/orcca) (Open Resources for Community College Algebra) — CC BY 4.0

# Factoring Out the Common Factor

In  Chapter , we learned how to multiply polynomials, such as when you start with $(x+2)(x+3)$ and obtain $x^2+5x+6$. This chapter, starting with this section, is about the *opposite* process---factoring. For example, starting with $x^2+5x+6$ and obtaining $(x+2)(x+3)$. We will start with the simplest kind of factoring: for example starting with $x^2+2x$ and obtaining $x(x+2)$.

## Motivation for Factoring

When you write $x^2+2x$, you have an algebraic expression built with two terms---two parts that are *added* together. When you write $x(x+2)$, you have an algebraic expression built with two factors---two parts that are *multiplied* together. Factoring is useful, because sometimes (but not always) having your expression written as parts that are *multiplied* together makes it easy to simplify the expression.

You've seen this with fractions. To simplify $\frac{15}{35}$, breaking down the numerator and denominator into factors is useful: $\frac{3\cdot5}{7\cdot5}$. Now you can see that the factors of $5$ cancel.

There are other reasons to appreciate the value in factoring. One reason is that there is a relationship between a factored polynomial and the horizontal intercepts of its graph. For example in the graph of $y=(x+2)(x-3)$, the horizontal intercepts are $(-2,0)$ and $(3,0)$. Note the $x$-values are $-2$ and $3$, and think about what happens when you subsitutue those numbers in for $x$ in $y=(x+2)(x-3)$. We will explore this more fully in  Section .

*A graph of $y=(x+2)(x-3)$*

## Identifying the Greatest Common Factor

The most basic technique for factoring involves recognizing the *greatest common factor* between two expressions, which is the largest factor that goes in evenly to both expressions. For example, the greatest common factor between $6$ and $8$ is $2$, since $2$ divides nicely into both $6$ and $8$ and no larger number would divide nicely into both $6$ and $8$.

Similarly, the greatest common factor between $4x$ and $3x^2$ is $x$. If you write $4x$ as a product of its factors, you have $2\cdot 2 \cdot x$. And if you fully factor $3x^2$, you have $3\cdot x\cdot x$. The only factor they have in common is $x$, so that is the greatest common factor. No larger expression goes in nicely to both expressions.

\*\*Example\*\*

Finding the Greatest Common Factor

What is the common factor between $6x^2$ and $70x$? Break down each of these into its factors:

$$
\begin{aligned}6x^2 & =2\cdot3\cdot x\cdot x & 70x & =2\cdot5\cdot7\cdot x \\ 6x^2 & =\attention{2}\cdot3\cdot \attention{x}\cdot x & 70x & =\attention{2}\cdot5\cdot7\cdot \attention{x}\end{aligned}
$$

With $2$ and $x$ in common, the greatest common factor is $2x$.

\*\*Exercise\*\*

## Factoring Out the Greatest Common Factor

We have learned the distributive property: $a(b+c)=ab+ac$. Perhaps you have thought of this as a way to "distribute" the number $a$ to each of $b$ and $c$. In this section, we will use the distributive property in the opposite way. If you have an expression $ab+ac$, it is equal to $a(b+c)$. In that example, we factored out $a$, which is the common factor between $ab$ and $ac$.

The following steps use the distributive property to factor out the greatest common factor between two or more terms.

Factoring Out the Greatest Common Factor

1. Identify the greatest common factor in all terms.
2. Write the greatest common factor outside a pair of parentheses with the appropriate addition or subtraction signs inside.
3. For each term from the original expression, what would you multiply the greatest common factor by to result in that term? Write your answer in the parentheses.

\*\*Example\*\*

To factor $12x^2+15x$:

1. The greatest common factor between $12x^2$ and $15x$ is $3x$.
2. $3x(\phantom{4x}+\phantom{5})$
3. $3x(4x+5)$

\*\*Example\*\*

Factor the polynomial $3x^3+3x^2-9$.

1. We identify the greatest common factor as $3$, because $3$ is the only common factor between $3x^3$, $3x^2$ and $9$.
2. We write:
   $$
    3x^3+3x^2-9=3(\phantom{x^2}+\phantom{x^2}-\phantom{3}) 
   $$
   .
3. We ask the question "$3$ times what gives $3x^3$?" The answer is $x^3$. Now we have:
   $$
    3x^3+3x^2-9=3(x^3+\phantom{x^2}-\phantom{3}) 
   $$
   .

   We ask the question "$3$ times what gives $3x^2$?" The answer is $x^2$. Now we have:
   $$
    3x^3+3x^2-9=3(x^3+x^2-\phantom{3}) 
   $$
   .

   We ask the question "$3$ times what gives $9$?" The answer is $3$. Now we have:
   $$
    3x^3+3x^2-9=3(x^3+x^2-3) 
   $$
   .

To check that this is correct, multiplying through $3(x^3+x^2-3)$ should give the original expression $3x^3+3x^2-9$. We check this, and it does.

\*\*Exercise\*\*

## Visualizing With Rectangles

In  Section , we learned one way to multiply polynomials using rectangle diagrams. Similarly, we can factor a polynomial with a rectangle diagram.

Factoring Out the Greatest Common Factor Using Rectangles

1. Put the terms into adjacent rectangles. Think of these as labeling the areas of each rectangle.
2. Identify the greatest common factor, and mark the height of the overall rectangle with it.
3. Mark the width of each rectangle based on each rectangle's area and height.
4. Since the overall rectangle's area equals its width times its height, the height is one factor, and the sum of the widths is another factor.

\*\*Example\*\*

We will factor $12x^2+15x$, the same polynomial from the example in  Algorithm , so that you may compare the two styles.

So $12x^2+15x$ factors as $3x(4x+5)$.

## More Examples of Factoring out the Common Factor

Previous examples did not cover every nuance with factoring out the greatest common factor. Here are a few more factoring examples that attempt to do so.

\*\*Example\*\*

Factor $-35m^5+5m^4-10m^3$.

First, we identify the common factor. The number $5$ is the greatest common factor of the three coefficients (which were $-35$, $5$, and $-10$) and also $m^3$ is the largest expression that divides $m^5$, $m^4$, and $m^3$. Therefore the greatest common factor is $5m^3$.

In this example, the leading term is a negative number. When this happens, we will make it common practice to take that negative as part of the greatest common factor. So we will proceed by factoring out $-5m^3$. Note the signs change inside the parentheses.

$$
\begin{aligned}-35m^5\highlight{{}+{}}5m^4\highlight{{}-{}}10m^3&=-5m^3(\phantom{7m^2}\highlight{{}-{}}\phantom{m}\highlight{{}+{}}\phantom{2}) \\ &=-5m^3(7m^2-\phantom{m}+\phantom{2}) \\ &=-5m^3(7m^2-m+\phantom{2}) \\ &=-5m^3(7m^2-m+2)\end{aligned}
$$

\*\*Example\*\*

Factor $14-7n^2+28n^4-21n$.

Notice that the terms are not in a standard order, with powers of $n$ decreasing as you read left to right. It is usually a best practice to rearrange the terms into the standard order first.

$$
 14-7n^2+28n^4-21n=28n^4-7n^2-21n+14 
$$

.

The number $7$ divides all of the numerical coefficients. Separately, no power of $n$ is part of the greatest common factor because the $14$ term has no $n$ factors. So the greatest common factor is just $7$. We proceed by factoring that out:

$$
\begin{aligned}14-7n^2+28n^4-21n&=28n^4-7n^2-21n+14 \\ &=7\mathopen{}\left(4n^4-n^2-3n+2\right)\mathclose{}\end{aligned}
$$

\*\*Example\*\*

Factor $24ab^2+16a^2b^3-12a^3b^2$.

There are two variables in this polynomial, but that does not change the factoring strategy. The greatest numerical factor between the three terms is $4$. The variable $a$ divides all three terms, and $b^2$ divides all three terms. So we have:

$$
\begin{aligned}24ab^2+16a^2b^3-12a^3b^2 &=4ab^2\mathopen{}\left(\phantom{6}+\phantom{4ab}-\phantom{3a^2}\right)\mathclose{}</mrow> <mrow>&=4ab^2\mathopen{}\left(6+\phantom{4ab}-\phantom{3a^2}\right)\mathclose{}</mrow> <mrow>&=4ab^2\mathopen{}\left(6+4ab-\phantom{3a^2}\right)\mathclose{}</mrow> <mrow> &=4ab^2\mathopen{}\left(6+4ab-3a^2\right)\mathclose{}\end{aligned}
$$

\*\*Example\*\*

Factor $4m^2n-3xy$.

There are no common factors in those two terms (unless you want to count $1$ or $-1$, but we do not count these for the purposes of identifying a greatest common factor). In this situation we can say the polynomial is *prime* or *irreducible*, and leave it as it is.

\*\*Example\*\*

Factor $-x^3+2x+18$.

There are no common factors in those three terms, and it would be correct to state that this polynomial is prime or irreducible. However, since its leading coefficient is negative, it may be wise to factor out a negative sign. So, it could be factored as $-\mathopen{}\left(x^3-2x-18\right)\mathclose{}$. Note that *every* term is negated as the leading negative sign is extracted.

##

\*\*Exercise\*\*

Given two terms, how would you describe their "greatest common factor?"

\*\*Exercise\*\*

If a simplified polynomial has four terms, and you factor out its greatest common factor, how many terms will remain inside a set of parentheses?

##

