> Content sourced from [ORCCA](https://pcc.edu/orcca) (Open Resources for Community College Algebra) — CC BY 4.0
# Order of Operations

When you write something down, it's important that the people who read it will understand what
      you actually meant. But language can be ambiguous. If we say in English, "two times three
      squared" , do we mean that: $2$ is multiplied by $3$ , and then the result is squared? That would mean
            first we have $6$ , and then we square $6$ to end with $36$ . Or that $2$ is multiplied by "three squared" ? That would mean first we square $3$ to get $9$ , and then we multiply by $2$ to end with $18$ . So it makes a difference, and the English phrase "two times three squared" is arguably
      ambiguous.

English is allowed to be ambiguous. But math needs to be unambiguous and mean the same thing
      for everyone who reads it. So for this reason, there are some rules that we've all agreed to
      that control what a math expression really means. These rules are called the "order of
      operations" , which we review here.

## Grouping Symbols

Consider the expression $2\cdot3^2$ . There are two math operations here: at some point
      two things will be multiplied, and at some point something will be raised to a power. The
      result depends on which operation you decide to do first: If you multiply $2\cdot3$ , and
      then square the result, you end with $36$ . If you square $3$ , and then multiply that
      result by $2$ , you end with $18$ . So if we want all people everywhere to interpret $2\cdot3^2$ in the same way, then only *one* of these can be correct.

One tool that we have to clearly tell readers which thing to do first is a pair of grouping
      symbols, like parentheses and brackets. If you *intend* to do the multiplication first,
      then writing $(2\cdot3)^2$ clearly tells your reader to do that. And if you *intend* to execute the power first, then writing $2\cdot\left(3^2\right)$ clearly
      tells your reader to do that.

To visualize the difference between $2\cdot \left(3^2\right)$ or $(2\cdot 3)^2$ ,
      consider these garden plots:

If we find $3^2$ , we have the area of one of the small square garden plots on the left.
      Then if we double that, we have $2\cdot\left(3^2\right)$ , the area of the entire left
      garden plot.

But if we find $(2\cdot3)^2$ , then first we are doubling $3$ . So we are getting the
      area of a large square garden plot whose sides are twice as long. We end up with the area of
      the entire garden plot on the right.

The point is that these amounts are different.


## Order of Operations

If math expressions used grouping symbols for every arithmetic operation, we wouldn't need "order of operations" . Some computer systems work that way, *requiring* the use of
      grouping symbols all the time. But it is more common to allow math expressions that don't have
      grouping symbols everywhere, like $5+3\cdot2$ . Should the addition $5+3$ be done
      first, or should the multiplication $3\cdot2$ be done first? We have a set of rules the
      world has agreed to, known as the "order of operations" . They tell us what to do first.

The *order of operations* is nothing more than an agreement that we all have made
      to prioritize doing arithmetic operations in a certain order.

To help remember the order of operations, consider the acronym PEMDAS . You might
      use mnemonic devices to help remember this such as "Please Excuse My Dear Aunt Sally" , "People Eat More Donuts After School" , etc.

We'll start with a few examples that only invoke a few operations each.

**Example**

**Practice with order of operations**
- $5-3(7-4)^2={}$
  *Solution*: 5-3(\nextoperation{7-4})^2\amp=5-3\nextoperation{(3)^2}
- $\phantom{5-3(7-4)^2}={}$
  *Solution*: \amp=5-\nextoperation{3(9)}
- $\phantom{5-3(7-4)^2}={}$
  *Solution*: \amp=5-27
- $\phantom{5-3(7-4)^2}={}$
  *Solution*: \amp=-22



## Absolute Value and Implied Grouping

Grouping symbols are more than just parentheses and brackets. Each of the following operations *implies* some grouping.

Absolute Value Bars Absolute value bars group the expression inside just like a set of parentheses would. In
            the expression $3\abs{2-5}$ , the first thing to do is subtract $2-5$ since
            that is inside the absolute value bars. Radicals The same is true with a radical symbol. Everything inside the radical is grouped. For
            example with $4+2\sqrt{12-3}$ , the first arithmetic to do is subtract $12-3$ . Fraction Bars A fraction bar can create *two* groups, one in the numerator and one in the
            denominator. With the expression $\frac{2+3}{5-2}\div4$ , the first arithmetic to
            take care of is adding $2+3$ and subtracting $5-2$ . These two groups can be
            worked on separately in any order. Exponents Content that is inside an exponent is treated as one group, as with $2+^{2\cdot3}$ .
            In that example, the first arithmetic to take care of is multiplying $2\cdot3$ .

Each of these implied groupings also ask you to do something once the arithmetic on the inside
      is completed. Actually taking the absolute value or the square root, perhaps. Doing the
      division in the case of a fraction. Raising something to a power. But *before* doing
      those things, all of the arithmetic *inside* the groups should be take care of.

**Example**

**Implied Grouping**


## Understanding (-a)^m versus -a^m

We noted in the order of operations that using the minus sign to negate a number has the same priority as multiplication and
      division.

How would you write a math expression that takes the number $-4$ and squares it? Is it OK
      to write $-4^2$ ? How about $(-4)^2$ ?

These expressions mean very different things. The second option, $(-4)^2$ is squaring the
      number $-4$ . The parentheses make this clear. The result is $16$ .

The first expression $-4^2$ is different. There are two actions here: a negation and
      exponentiation. According to the order of operations, the exponentiation has higher priority,
      so we should do $4^2$ first. -4^2\amp=-\nextoperation{4^2} \amp=-16 and this is not the same as $(-4)^2$ , which is *positive* $16$ .

> **Negative Numbers Raised to Powers**
> You may find yourself needing to raise a negative number to a power, and using a calculator
        to do the work for you. If you do not understand the issue described above, then you may get
        incorrect results. Entering `-4^2` into a calculator or computer will result in $-16$ . Entering `(-4)^2` into a calculator or computer will result in $16$ . Try entering these into your own calculator.

**Negating and Raising to Powers**
- $-3^4={}$ and $(-3)^4={}$
  *Solution*: $-3^4=-81$ and $(-3)^4=81$
- $-4^3={}$ and $(-4)^3={}$
  *Solution*: $-4^3=-64$ and $(-4)^3=-64$
- $-1.1^2={}$ and $(-1.1)^2={}$
  *Solution*: $-1.1^2=-1.21$ and $(-1.1)^2=1.21$

You might notice in that $-4^3$ and $(-4)^3$ each have the same result, $-64$ . It's true that the results are the same,
      but the two expressions say different things. With $-4^3$ , you raise to a power first,
      then negate. With $(-4)^3$ , you negate first, then raise to a power. It's like two
      different roads that happen to lead to the same place, which happens sometimes.
