> Content sourced from [ORCCA](https://pcc.edu/orcca) (Open Resources for Community College Algebra) — CC BY 4.0

# Order of Operations

When you write something down, it's important that the people who read it will understand what you actually meant. But language can be ambiguous. If we say in English, "two times three squared", do we mean that:

- $2$ is multiplied by $3$, and then the result is squared? That would mean first we have $6$, and then we square $6$ to end with $36$.
- Or that $2$ is multiplied by "three squared"? That would mean first we square $3$ to get $9$, and then we multiply by $2$ to end with $18$.

So it makes a difference, and the English phrase "two times three squared" is arguably ambiguous.

English is allowed to be ambiguous. But math needs to be unambiguous and mean the same thing for everyone who reads it. So for this reason, there are some rules that we've all agreed to that control what a math expression really means. These rules are called the "order of operations", which we review here.

## Grouping Symbols

Consider the expression $2\cdot3^2$. There are two math operations here: at some point two things will be multiplied, and at some point something will be raised to a power. The result depends on which operation you decide to do first: If you multiply $2\cdot3$, and then square the result, you end with $36$. If you square $3$, and then multiply that result by $2$, you end with $18$. So if we want all people everywhere to interpret $2\cdot3^2$ in the same way, then only *one* of these can be correct.

One tool that we have to clearly tell readers which thing to do first is a pair of grouping symbols, like parentheses and brackets. If you *intend* to do the multiplication first, then writing $(2\cdot3)^2$clearly tells your reader to do that. And if you *intend* to execute the power first, then writing $2\cdot\left(3^2\right)$ clearly tells your reader to do that.

To visualize the difference between $2\cdot \left(3^2\right)$ or $(2\cdot 3)^2$, consider these garden plots:

*is squared, then doubled: $2\cdot\left(3^2\right)$*

*&#x20;is doubled, then squared: $(2\cdot3)^2$*

If we find $3^2$, we have the area of one of the small square garden plots on the left. Then if we double that, we have $2\cdot\left(3^2\right)$, the area of the entire left garden plot.

But if we find $(2\cdot3)^2$, then first we are doubling $3$. So we are getting the area of a large square garden plot whose sides are twice as long. We end up with the area of the entire garden plot on the right.

The point is that these amounts are different.

\*\*Exercise\*\*

## Order of Operations

If math expressions used grouping symbols for every arithmetic operation, we wouldn't need "order of operations". Some computer systems work that way, *requiring* the use of grouping symbols all the time. But it is more common to allow math expressions that don't have grouping symbols everywhere, like $5+3\cdot2$. Should the addition $5+3$ be done first, or should the multiplication $3\cdot2$ be done first? We have a set of rules the world has agreed to, known as the "order of operations". They tell us what to do first.

The *order of operations* is nothing more than an agreement that we all have made to prioritize doing arithmetic operations in a certain order.

Order of Operations

- Grouping symbols "group" the expression inside them, and the arithmetic within that group needs to be done first.
- After grouping symbols, exponentiation has the highest priority. Raise things to powers before doing any other arithmetic.
- After all exponentiation is done, start doing multiplication, division, and negation. These things all have equal priority. If there are more than one of them in your expression, do these things in order from left to right as you would naturally read the expression.
- After all other arithmetic is done, addition and subtraction are all that is left. They have equal priority. If there are more than one of them in your expression, do these things in order from left to right as you would naturally read the expression.

To help remember the order of operations, consider the acronym . You might use mnemonic devices to help remember this such as "Please Excuse My Dear Aunt Sally", "People Eat More Donuts After School", etc.

We'll start with a few examples that only invoke a few operations each.

\*\*Example\*\*

Use the order of operations to simplify the following expressions.

1. $10+2\cdot 3$. With this expression, we have addition and multiplication. The order of operations says the multiplication has higher priority, so do that first:
   $$
   \begin{aligned}10+2\cdot 3& =10+2\cdot 3 \\ &=10+\highlight{6} \\ &=\highlight{16}\end{aligned}
   $$
2. $4+10\div 2 - 1$. According to the order of operations, the first thing to do is the division. After that, we'll apply the addition and subtraction working left to right:
   $$
   \begin{aligned}4+10\div2-1&=4+10\div2-1 \\ &=4+\highlight{5}-1 \\ &=\highlight{9}-1 \\ &=\highlight{8}\end{aligned}
   $$
3. $7-10+4$. This expression *only* has subtraction and addition. These operations tie for priority (even though "A" comes before "S" in ). So we work left to right to do them:
   $$
   \begin{aligned}7-10+4&=7-10+4 \\ &=\highlight{-3}+4 \\ &=1\end{aligned}
   $$
4. $20\div 4\cdot 7$. This expression has only division and multiplication. Again, we have two operations that tie for priority, so we do them in order from left to right:
   $$
   \begin{aligned}20\div 4\cdot 5&=20\div 4 \cdot 5 \\ &=\highlight{5}\cdot5 \\ &=\highlight{25}\end{aligned}
   $$
5. $(6+7)^2$. Here we have addition inside parentheses, and an exponent of $2$ outside. We must do the arithmetic inside the parentheses first:
   $$
   \begin{aligned}(6+7)^2&= (6+7)^2 \\ &= \highlight{13}^2 \\ &= \highlight{169}\end{aligned}
   $$
6. $4(2)^3$. This expression has multiplication and an exponent. There are parentheses, but no operation inside them. Parentheses used this way are just to make it clear that the $4$ and $2$ are separate numbers, not to be confused with the number $42$. Exponentiation has the higher priority, so we'll do that part first, and then multiply:
   $$
   \begin{aligned}4(2)^3 &= 4(2)^3 \\ &= 4(\highlight{8}) \\ &= \highlight{32}\end{aligned}
   $$

\*\*Remark\*\*

There are several different ways to write multiplication. We can use the symbols $\cdot$, $\times$, and $*$ to mean multiplication. We can also write two things right next to each other with no symbol in between them to mean multiplication. That is what is happening in  Item , where the $4$ is written right next to the $(2)^3$ with no symbol in between.

Using a symbol for multiplication is called "explicit multiplication" and not writing  any symbol at all is called "implicit multiplication". For this textbook, explicit and  implicit multiplication have the same priority in the order of operations. However there are some conventions out in the real world where implicit multiplication has a higher priority in the order of operations than explicit multiplication. You may have seen memes with expressions like $6\div2(3)$ that play on how the real world has more than one convention for the order of operations.

\*\*Exercise\*\*

Practice with order of operations

$5-3(7-4)^2={}$

Solution

$\phantom{5-3(7-4)^2}={}$

Solution

$\phantom{5-3(7-4)^2}={}$

Solution

$\phantom{5-3(7-4)^2}={}$

Solution

\*\*Exercise\*\*

\*\*Exercise\*\*

## Absolute Value and Implied Grouping

Grouping symbols are more than just parentheses and brackets. Each of the following operations *implies* some grouping.



Each of these implied groupings also ask you to do something once the arithmetic on the inside is completed. Actually taking the absolute value or the square root, perhaps. Doing the division in the case of a fraction. Raising something to a power. But *before* doing those things, all of the arithmetic *inside* the groups should be take care of.

\*\*Example\*\*

Use the order of operations to simplify the following expressions.

1. $4-3\abs{5-7}$. The absolute value bars group the $5-7$ so we must do that subtraction first. Then we take the absolute value and continue:
   $$
   \begin{aligned}4-3\abs{5-7} &= 4-3\abs{5-7} \\ &= 4-3\abs{\highlight{-2}} \\ &= 4-3(\highlight{2}) \\ &= 4-\highlight{6} \\ &= \highlight{-2}\end{aligned}
   $$
   It would be a mistake to subtract $4-3$ first, because that $3$ is multiplied by the $\abs{5-7}$. So subtracting $4-3$ would violate the order of operations.
2. $8-\sqrt{5^2-8\cdot 2}$. The radical is grouping $5^2-8\cdot 2$, which must be simplified first. Then we take the square root and continue:
   $$
   \begin{aligned}8-\sqrt{5^2-8\cdot 2} &= 8-\sqrt{5^2-8\cdot 2} \\ &= 8-\sqrt{\highlight{25}-8\cdot 2} \\ &= 8-\sqrt{25-\highlight{16}} \\ &= 8-\sqrt{\highlight{9}} \\ &= 8-\highlight{3} \\ &= \highlight{5}\end{aligned}
   $$
3. $\dfrac{2^4+3\cdot 6}{5-18\div 2}$. We recognize that the fraction bar is creating two groups. We should simplify the numerator and denominator separately according to the order of operations, and proceed from there:
   $$
   \begin{aligned}\frac{2^4+3\cdot 6}{5-18\div 2} &= \frac{2^4+3\cdot 6}{5-18\div 2} \\ &=\frac{\highlight{16}+3\cdot 6}{5-\highlight{9}} \\ &=\frac{16+\highlight{18}}{\highlight{-4}} \\ &=\frac{\highlight{34}}{-4} \\ &=-\frac{17}{2}\end{aligned}
   $$

\*\*Exercise\*\*

Implied Grouping

\*\*Exercise\*\*

## Understanding $(-a)^m$ versus $-a^m$

We noted in the  order of operations  that using the minus sign to negate a number has the same priority as multiplication and division.

How would you write a math expression that takes the number $-4$ and squares it? Is it OK to write $-4^2$? How about $(-4)^2$?

These expressions mean very different things. The second option, $(-4)^2$ is squaring the number $-4$. The parentheses make this clear. The result is $16$.

The first expression $-4^2$ is different. There are two actions here: a negation and exponentiation. According to the order of operations, the exponentiation has higher priority, so we should do $4^2$ first.

$$
\begin{aligned}-4^2&=-4^2 \\ &=-16\end{aligned}
$$

and this is not the same as $(-4)^2$, which is *positive*$16$.

\*\*Warning\*\*

Negative Numbers Raised to Powers

You may find yourself needing to raise a negative number to a power, and using a calculator to do the work for you. If you do not understand the issue described above, then you may get incorrect results.

- Entering `-4^2` into a calculator or computer will result in $-16$.
- Entering `(-4)^2` into a calculator or computer will result in $16$.

Try entering these into your own calculator.

\*\*Exercise\*\*

Negating and Raising to Powers

$-3^4={}$ and $(-3)^4={}$

Solution

$-3^4=-81$ and $(-3)^4=81$

$-4^3={}$ and $(-4)^3={}$

Solution

$-4^3=-64$ and $(-4)^3=-64$

$-1.1^2={}$ and $(-1.1)^2={}$

Solution

$-1.1^2=-1.21$ and $(-1.1)^2=1.21$

You might notice in  the exercise  that $-4^3$ and $(-4)^3$ each have the same result, $-64$. It's true that the results are the same, but the two expressions say different things. With $-4^3$, you raise to a power first, then negate. With $(-4)^3$, you negate first, then raise to a power. It's like two different roads that happen to lead to the same place, which happens sometimes.

##

