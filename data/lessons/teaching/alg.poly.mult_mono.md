> Content sourced from [ORCCA](https://pcc.edu/orcca) (Open Resources for Community College Algebra) — CC BY 4.0
# Multiplying Polynomials

Previously in, we learned how to multiply two monomials together (such as $4xy\cdot3x^2$). And in, we learned how to add and subtract polynomials even when there is more than one term (such as $(4x^2-3x)+(5x^2+x-2)$). In this section, we will learn how to *multiply* polynomials with *more than one* term.

**Revenue**

## Review of the Distributive Property

Polynomial multiplication relies on the, and may also rely on the properties of exponents. When we multiply a monomial with a binomial, we can distribute the monomial to each term in the binomial. For example, 
\[ \begin{aligned} \highlight{-4x}(3x^2+5) &= \multiplyleft{(-4x)}\left(3x^2\right)+\multiplyleft{(-4x)}(5) \\ &=-12x^3-20x \end{aligned} \]

**Example**
A rectangle's length is $4$ meters longer than its width. Assume its width is $w$ meters. Use a simplified polynomial to model the rectangle's area in terms of $w$ as the only variable.

*Solution*
In the second line of work above, we recognize that $(w+4)w$ is the same as $w(w+4)$. Whether the $w$ is written before or after the binomial, we are still able to distribute the $w$ to simplify the product.

The distributive property can be understood visually with a *generic rectangle*.

The big rectangle consists of two smaller rectangles. The big rectangle's area is $2x(3x+4)$, and the sum of those two smaller rectangles is $2x\cdot3x+2x\cdot4$. Since the sum of the areas of those two smaller rectangles is the same as the bigger rectangle's area, we have: 
\[ \begin{aligned} 2x(3x+4) &= 2x\cdot3x+2x\cdot4 \\ &= 6x^2+8x \end{aligned} \]
 Generic rectangles can be used to visualize multiplying polynomials.

## Multiplying Binomials

### Multiplying Binomials Using Distribution

Whether we're multiplying a monomial with a polynomial or two larger polynomials together, the first step is still based on the. We'll start with multiplying two binomials and then move on to examples with larger polynomials.

We know we can distribute the $3$ in $(x+2)3$ to obtain $(x+2)\multiplyright{3}=x\multiplyright{3}+2\multiplyright{3}$. We can actually distribute *anything* across $(x+2)$ if it is multiplied. For example: 
\[(x+2)\cat=x\cdot \cat + 2\cdot \cat\]

Keeping this in mind, we can multiply $(x+2)(x+3)$ by distributing the $(x+3)$ across $(x+2)$. Think of $(x+3)$ as the cat. 
\[(x+2)\highlight{(x+3)} = x\highlight{(x+3)} + 2\highlight{(x+3)}\]
 To continue, distribure more: 
\[ \begin{aligned} (x+2)\highlight{(x+3)} &= x\highlight{(x+3)} + 2\highlight{(x+3)} \\ &= x \cdot \highlight{x} + x \cdot \highlight{3} + 2 \cdot \highlight{x} + 2 \cdot \highlight{3} \\ &=x^2+3x+2x+6 \\ &=x^2+5x+6 \end{aligned} \]

To multiply a binomial by another binomial, we simply used distribution twice and then simplified the resulting terms. (Multiplying any two polynomials can be done using these same steps. We will return to that idea later in the section.)

### Multiplying Binomials Using FOIL

Instead of using two applications of the distributive property, people often memorize a shortcut using the acronym FOIL. The letters refer to the pairs of terms from each binomial that end up multiplied together.

If we take a closer look at the example we just completed, $(x+2)(x+3)$, we can highlight how the FOIL process works. FOIL is the acronym for "First, Outer, Inner, Last". 
\[ \begin{aligned} (x+2)(x+3)&= (\overbrace{{x} \stackrel{}{\cdot} {x}}^{\text{F}}) + (\overbrace{{3} \stackrel{}{\cdot} {x}}^{\text{O}}) + (\overbrace{{2} \stackrel{}{\cdot} {x}}^{\text{I}}) + (\overbrace{{2} \stackrel{}{\cdot} {3}}^{\text{L}}) \\ &=x^2+3x+2x+6 \\ &=x^2+5x+6 \end{aligned} \]
 F: $x^2$ The $x^2$ term comes from the product of *first* terms from each binomial. O: $3x$ The $3x$ term comes from the product of the *outer* terms from each binomial. This was $x$ in the front of the first binomial and $3$ in the back of the second binomial. I: $2x$ The $2x$ term comes from the product of the *inner* terms from each binomial. This was $2$ in the back of the first binomial and $x$ in the front of the second binomial. L: $6$ The constant term $6$ comes from the product of the *last* terms of each binomial.

### Multiplying Binomials Using Generic Rectangles

We could approach this same example using the generic rectangle method. To use generic rectangles, we treat $x+2$ as the base of a rectangle, broken up into $x$ and $2$. Similarly we treat $x+3$ as the height, broken up into $x$ and $3$. Their product, $(x+2)(x+3)$, represents the large rectangle's area.

The big rectangle consists of four smaller rectangles. We can find each small rectangle's area in the next diagram with the area formula for a rectangle: $\text{area}=\text{base}\cdot\text{height}$.

To finish finding this product, we add together the areas of the four smaller rectangles: 
\[ \begin{aligned} (x+2)(x+3)&=x^2+3x+2x+6 \\ &=x^2+5x+6 \end{aligned} \]

Notice that the areas of the four smaller rectangles are exactly the same as the four terms we obtained using distribution, which are also the same four terms that came from the FOIL method. The FOIL method and generic rectangles approach are just different ways to represent the distributing.

**Example**
Multiply $(2x-3y)(4x-5y)$ using distribution.

*Solution*
We first distribute the second binomial across $(2x-3y)$. Then we'll distribute again, and simplify the terms that are left. 
\[ \begin{aligned} (2x-3y)\highlight{(4x-5y)}&=2x\highlight{(4x-5y)}-3y\highlight{(4x-5y)} \\ &=8x^2-10xy-12xy+15y^2 \\ &=8x^2-22xy+15y^2 \end{aligned} \]

**Example**
Multiply $(2x-3y)(4x-5y)$ using FOIL.

*Solution*
First, Outer, Inner, Last: Either with arrows on paper or mentally in our heads, we'll pair up the terms from each factor and multiply those pairs together. 
\[ \begin{aligned} (2x-3y)(4x-5y)&= (\overbrace{{\stackrel{}{2x}}\cdot{4x}}^{\large\text{F}})+ (\overbrace{{\stackrel{}{2x}}\cdot{(-5y)}}^{\large\text{O}})+ (\overbrace{{\stackrel{}{-3y}}\cdot{4x}}^{\large\text{I}})+ (\overbrace{{\stackrel{}{-3y}}\cdot{(-5y}}^{\large\text{L}}) \\ &=8x^2-10xy-12xy+15y^2 \\ &=8x^2-22xy+15y^2 \end{aligned} \]

**Example**
Multiply $(2x-3y)(4x-5y)$ using generic rectangles.

*Solution*
We begin by drawing four rectangles and marking their bases and heights with terms in the given binomials:

Next, we calculate each inner rectangle's area by multiplying its base with its height:

Finally, we add up all the inner rectangle areas to find the product: 
\[ \begin{aligned} (2x-3y)(4x-5y)&=8x^2-10xy-12xy+15y^2 \\ &=8x^2-22xy+15y^2 \end{aligned} \]

**Example**
Multiply and simplify the formula for Avery's organic jam revenue, $R$ (in dollars), from. In that example $R=(13+0.25x)(1500-50x)$ and $x$ represents the number of times they raised the price by 25 cents.

*Solution*
To multiply this, we'll use FOIL: 
\[ \begin{aligned} R &= \left(13+0.25x\right)\left(1500-50x\right) \\ &= \left(13\cdot1500\right) - \left(13 \cdot 50x \right) + \left( 0.25x \cdot 1500 \right) - \left( 0.25x \cdot 50x \right) \\ &= 19500 - 650x + 375x - 12.5x^2 \\ &= -12.5x^2 - 275x + 19500 \end{aligned} \]
 Now we have a formula for Avery's revenue based on how many times they raise the price. If we wanted to, we could study this formula more to find ways to maximize that revenue.

**Example**
Tyrone is an artist and he sells each of his paintings for $\$200$. Currently, he can sell $100$ paintings per year. So his annual revenue from selling paintings is $\$200\cdot100=\$20000$. He plans to raise the price. However, for each $20 price increase, his customers will buy $5$ fewer paintings each year. Assume Tyrone would raise the price of his paintings $x$ times, each time by $20. Use an expanded polynomial to represent his new revenue per year.

*Solution*
Currently, each painting costs $200. After raising the price $x$ times, each time by $20, each painting's new price would be $200+20x$ dollars.

Currently, Tyrone sells $100$ paintings per year. After raising the price $x$ times, each time selling $5$ fewer paintings, he would end up selling $100-5x$ paintings per year.

His annual revenue can be calculated by multiplying each painting's price by the number of paintings he would sell: 
\[ \begin{aligned} \text{annual revenue}&=(\text{price})(\text{number of sales}) \\ &=(200+20x)(100-5x) \\ &=200(100)+200(-5x)+20x(100)+20x(-5x) \\ &=20000-1000x+2000x-100x^2 \\ &=-100x^2+1000x+20000 \end{aligned} \]
 After raising the price $x$ times, each time by $20, Tyrone's annual revenue from paintings would be $-100x^2+1000x+20000\$ dollars.

**Example**
What would happen if we needed to multiply two binomials, but also there is a monomial out front? What is the result for $3(x + 2)(2x + 5)$?

*Solution*
We recommend focusing first on the two binomials. We can multiply them using FOIL even though they are written last in the product. 
\[ \begin{aligned} 3(x + 2)(2x + 5) &= 3\left(2x^2 + 5x + 4x + 10\right) \\ &= 3\left(2x^2 + 9x + 10\right) \end{aligned} \]
 Now we can complete the multiplication just by distributing that $3$. 
\[ \begin{aligned} 3(x + 2)(2x + 5) &= 6x^2 + 27x + 30 \end{aligned} \]

## Multiplying Polynomials Larger Than Binomials

To multiply polynomials that have more than two terms, we can use repeated distribution and monomial multiplication. Whether we are working with binomials, trinomials, or even larger polynomials, the process is fundamentally the same.

**Example**

