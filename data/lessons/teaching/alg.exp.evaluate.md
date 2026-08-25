> Content sourced from [ORCCA](https://pcc.edu/orcca) (Open Resources for Community College Algebra) — CC BY 4.0
# Introduction to Exponent Rules

In this section, we look at some properties of exponents that we can use when simplifying expressions that have multiplication and exponents.

## Exponent Basics

Before we discuss any exponent rules, let's remind ourselves about exponent fundamentals. When working with expressions with exponents, we have the following vocabulary: 
\[\text{base}^{\text{exponent}} = \text{power}\]
 For example, when we calculate $8^{2} = 64$, the *base* is $8$, the *exponent* is $2$, and the expression $8^{2}$ is called the 2nd *power* of $8$.

Exponents indicate repeated multiplication. When the exponent is a positive integer, the power can be rewritten as repeated multiplication of the base. For example, the $4$ th power of $3$ can be written as $4$ factors of $3$ like so: 
\[3^{4} = 3 \cdot 3 \cdot 3 \cdot 3\]

## Exponent Rules

### Product Rule

If we write out $3^5\cdot 3^2$ without using exponents, we'd have: 
\[3^5 \cdot 3^2 = \left(3 \cdot 3\cdot 3\cdot 3\cdot 3\right) \cdot \left(3 \cdot 3\right)\]
 If we then count how many $3$ s are being multiplied together, we find we have $5+2=7$, a total of seven $3$ s. So $3^5\cdot 3^2$ simplifies like this: 
\[ \begin{aligned} 3^5\cdot 3^2 &= 3^{5+2} \\ &= 3^7 \end{aligned} \]

**Example**
Simplify $x^2\cdot x^3$. To simplify $x^2\cdot x^3$, we write this out in its expanded form, as a product of $x$ 's, we have 
\[ \begin{aligned} x^2\cdot x^3 &=(x\cdot x)(x \cdot x \cdot x) \\ &=x\cdot x\cdot x \cdot x \cdot x \\ &=x^5 \end{aligned} \]
 Note that we obtained the exponent of $5$ by adding $2$ and $3$.

This demonstrates our first exponent rule, the *Product Rule*: when multiplying two expressions that have the same base, we can simplify the product by adding the exponents. $x^m \cdot x^n = x^{m+n}$

Recall that $x=x^1$. It helps to remember this when multiplying certain expressions together.

**Example**
Multiply $x(x^3+2)$ by using the distributive property. According to the distributive property, 
\[x(x^3+2)=x\cdot x^3 + x\cdot2\]
 How can we simplify that term $x\cdot x^3$? It's really the same as $x^1\cdot x^3$, so according to the Product Rule, it is $x^4$. So we have: 
\[ \begin{aligned} x(x^3+2)&=x\cdot x^3 + x\cdot2 \\ &=x^4+2x \end{aligned} \]

### Power to a Power Rule

If we write out $\left(3^5\right)^2$ without using exponents, we'd have $3^5$ multiplied by itself: 
\[ \begin{aligned} \left(3^5\right)^2 &= \left(3^5\right)\cdot \left(3^5\right) \\ &= \left(3\cdot 3\cdot 3\cdot 3 \cdot 3 \right) \cdot \left(3 \cdot 3\cdot 3\cdot 3\cdot 3\right) \end{aligned} \]
 If we again count how many $3$ s are being multiplied, we have a total of two groups each with five $3$ s. So we'd have $2\cdot 5=10$ instances of a $3$. So $\left(3^5\right)^2$ simplifies like this: 
\[ \begin{aligned} \left(3^5\right)^2 &= 3^{2\cdot 5} \\ &= 3^{10} \end{aligned} \]

**Example**
Simplify $\left(x^2\right)^3$. To simplify $\left(x^2\right)^3$, we write this out in its expanded form, as a product of $x$ 's, we have 
\[ \begin{aligned} \left(x^2\right)^3 &=\left(x^2\right) \cdot \left(x^2\right)\cdot\left(x^2\right) \\ &=(x \cdot x)\cdot (x \cdot x)\cdot (x \cdot x) \\ &=x^6 \end{aligned} \]
 Note that we obtained the exponent of $6$ by multiplying $2$ and $3$.

This demonstrates our second exponent rule, the *Power to a Power Rule*: when a base is raised to an exponent and that expression is raised to another exponent, we multiply the exponents. 
\[\left(x^m\right)^n = x^{m \cdot n}\]

### Product to a Power Rule

The third exponent rule deals with having multiplication inside a set of parentheses and an exponent outside the parentheses. If we write out $\left(3t\right)^5$ without using an exponent, we'd have $3t$ multiplied by itself five times: 
\[(3t)^5= (3t)(3t)(3t)(3t)(3t)\]
 Keeping in mind that there is multiplication between every $3$ and $t$, and multiplication between all of the parentheses pairs, we can reorder and regroup the factors: 
\[ \begin{aligned} \left(3t\right)^5 &= (3\cdot t)\cdot (3\cdot t)\cdot (3\cdot t)\cdot (3\cdot t)\cdot (3\cdot t) \\ &= \left(3\cdot 3\cdot 3\cdot 3\cdot 3 \right) \cdot \left(t \cdot t \cdot t \cdot t \cdot t\right) \\ &= 3^5 t^5 \end{aligned} \]
 We could leave it written this way if $3^5$ feels especially large. But if you are able to evaluate $3^5=243$, then perhaps a better final version of this expression is $243t^5$.

We essentially applied the outer exponent to each factor inside the parentheses. It is important to see how the exponent $5$ applied to *both* the $3$ *and* the $t$, not just to the $t$.

**Example**
Simplify $(xy)^5$. To simplify $(xy)^5$, we write this out in its expanded form, as a product of $x$ 's and $y$ 's, we have 
\[ \begin{aligned} (xy)^5 &=(x \cdot y) \cdot (x \cdot y) \cdot (x \cdot y) \cdot (x \cdot y) \cdot (x \cdot y) \\ &=(x \cdot x \cdot x \cdot x \cdot x) \cdot (y \cdot y \cdot y \cdot y \cdot y) \\ &=x^5 y^5 \end{aligned} \]
 Note that the exponent on $xy$ can simply be applied to both $x$ and $y$.

This demonstrates our third exponent rule, the *Product to a Power Rule*: when a product is raised to an exponent, we can apply the exponent to each factor in the product. 
\[\left(x\cdot y\right)^n = x^{n}\cdot y^{n}\]

**Summary of the Rules of Exponents for Multiplication**

- **Product Rule**: $b^{m} \cdot b^{n} = b^{m+n}$
- **Power to a Power Rule**: $(b^{m})^{n} = b^{m\cdot n}$
- **Product to a Power Rule**: $(bc)^{m} = b^{m} \cdot c^{m}$

Many examples will make use of more than one exponent rule. In deciding which exponent rule to work with first, it's important to remember that the order of operations still applies.

**Example**
Simplify the following expressions. 
1. $\left(3^7r^5\right)^4$
2. $\left(t^3\right)^2\cdot \left(t^4\right)^5$

*Solution*
1. Since we cannot simplify anything inside the parentheses, we'll begin simplifying this expression using the. We'll apply the outer exponent of 4 to each factor inside the parentheses. Then we'll use the to finish the simplification process. $\begin{aligned} \left(3^7r^5\right)^4 &= \left(3^7\right)^4 \cdot \left(r^5\right)^4 \\ &= 3^{7\cdot4} \cdot r^{5\cdot 4} \\ &= 3^{28}r^{20} \end{aligned}$ Note that $3^{28}$ is too large to actually compute, even with a calculator, so we leave it written as $3^{28}$.
2. According to the order of operations, we should first simplify any exponents before carrying out any multiplication. Therefore, we'll begin simplifying this by applying the and then finish using the. $\begin{aligned} \left(t^3\right)^2\cdot \left(t^4\right)^5 &= t^{3\cdot2}\cdot t^{4\cdot5} \\ &= t^6 \cdot t^{20} \\ &= t^{6+20} \\ &= t^{26} \end{aligned}$

