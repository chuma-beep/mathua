> Content sourced from [ORCCA](https://pcc.edu/orcca) (Open Resources for Community College Algebra) — CC BY 4.0
# Multiplying Polynomials

Previously in , we learned how to
      multiply two monomials together (such as $4xy\cdot3x^2$ ). And in , we learned how to add and subtract
      polynomials even when there is more than one term (such as $(4x^2-3x)+(5x^2+x-2)$ ). In
      this section, we will learn how to *multiply* polynomials with *more than one* term.

## Review of the Distributive Property

Polynomial multiplication relies on the ,
      and may also rely on the properties of exponents . When we
      multiply a monomial with a binomial, we can distribute the monomial to each term in the
      binomial. For example, \highlight{-4x}(3x^2+5) \amp= \multiplyleft{(-4x)}\left(3x^2\right)+\multiplyleft{(-4x)}(5) \amp=-12x^3-20x

**Example**
A rectangle's length is $4$ meters longer than its width. Assume its width is $w$ meters. Use a simplified polynomial to model the rectangle's area in terms of $w$ as the only variable.

*Solution*
In the second line of work above, we recognize that $(w+4)w$ is the same as $w(w+4)$ . Whether the $w$ is written before or after the binomial, we are still
          able to distribute the $w$ to simplify the product.


The distributive property can be understood visually with a *generic rectangle* .

The big rectangle consists of two smaller rectangles. The big rectangle's area is $2x(3x+4)$ , and the sum of those two smaller rectangles is $2x\cdot3x+2x\cdot4$ .
      Since the sum of the areas of those two smaller rectangles is the same as the bigger
      rectangle's area, we have: 2x(3x+4) \amp= 2x\cdot3x+2x\cdot4 \amp= 6x^2+8x Generic rectangles can be used to visualize multiplying polynomials.

## Multiplying Binomials

Notice that the areas of the four smaller rectangles are exactly the same as the four terms we
      obtained using distribution, which are also the same four terms that came from the FOIL method. The FOIL method and generic rectangles approach are
      just different ways to represent the distributing.

**Example**
Multiply $(2x-3y)(4x-5y)$ using distribution.

*Solution*
We first distribute the second binomial across $(2x-3y)$ . Then we'll distribute
          again, and simplify the terms that are left. (2x-3y)\highlight{(4x-5y)}\amp=2x\highlight{(4x-5y)}-3y\highlight{(4x-5y)} \amp=8x^2-10xy-12xy+15y^2 \amp=8x^2-22xy+15y^2

**Example**
Multiply $(2x-3y)(4x-5y)$ using FOIL .

*Solution*
First, Outer, Inner, Last: Either with arrows on paper or mentally in our heads, we'll
          pair up the terms from each factor and multiply those pairs together. (2x-3y)(4x-5y)\amp=
            (\overbrace{{\stackrel{}{2x}}\cdot{4x}}^{\large\text{F}})+
            (\overbrace{{\stackrel{}{2x}}\cdot{(-5y)}}^{\large\text{O}})+
            (\overbrace{{\stackrel{}{-3y}}\cdot{4x}}^{\large\text{I}})+
            (\overbrace{{\stackrel{}{-3y}}\cdot{(-5y}}^{\large\text{L}}) \amp=8x^2-10xy-12xy+15y^2 \amp=8x^2-22xy+15y^2

**Example**
Multiply $(2x-3y)(4x-5y)$ using generic rectangles.

*Solution*
We begin by drawing four rectangles and marking their bases and heights with terms in the given binomials:
Next, we calculate each inner rectangle's area by multiplying its base with its height:
Finally, we add up all the inner rectangle areas to find the product: (2x-3y)(4x-5y)\amp=8x^2-10xy-12xy+15y^2 \amp=8x^2-22xy+15y^2

**Example**
Multiply and simplify the formula for Avery's organic jam revenue, $R$ (in dollars),
          from . In that example $R=(13+0.25x)(1500-50x)$ and $x$ represents the number of times they raised the price by 25 cents.

*Solution*
To multiply this, we'll use FOIL : R \amp= \left(13+0.25x\right)\left(1500-50x\right) \amp= \left(13\cdot1500\right) - \left(13 \cdot 50x  \right) + \left( 0.25x \cdot 1500 \right) - \left( 0.25x \cdot 50x \right) \amp= 19500  - 650x + 375x - 12.5x^2 \amp= -12.5x^2 - 275x + 19500 Now we have a formula for Avery's revenue based on how many times they raise the price. If
          we wanted to, we could study this formula more to find ways to maximize that revenue.

**Example**
Tyrone is an artist and he sells each of his paintings for $\$200$ . Currently, he can
          sell $100$ paintings per year. So his annual revenue from selling paintings is $\$200\cdot100=\$20000$ . He plans to raise the price. However, for each $20 price
          increase, his customers will buy $5$ fewer paintings each year. Assume Tyrone would raise the price of his paintings $x$ times, each time by $20. Use
          an expanded polynomial to represent his new revenue per year.

*Solution*
Currently, each painting costs $200. After raising the price $x$ times, each time by
          $20, each painting's new price would be $200+20x$ dollars.
Currently, Tyrone sells $100$ paintings per year. After raising the price $x$ times, each time selling $5$ fewer paintings, he would end up selling $100-5x$ paintings per year.
His annual revenue can be calculated by multiplying each painting's price by the number of
          paintings he would sell: \text{annual revenue}\amp=(\text{price})(\text{number of sales}) \amp=(200+20x)(100-5x) \amp=200(100)+200(-5x)+20x(100)+20x(-5x) \amp=20000-1000x+2000x-100x^2 \amp=-100x^2+1000x+20000 After raising the price $x$ times, each time by $20, Tyrone's annual revenue from
          paintings would be $-100x^2+1000x+20000$ dollars.

**Example**
What would happen if we needed to multiply two binomials, but also there is a monomial out
          front? What is the result for $3(x + 2)(2x + 5)$ ?

*Solution*
We recommend focusing first on the two binomials. We can multiply them using FOIL even though they are written last in the product. 3(x + 2)(2x + 5) \amp= 3\left(2x^2 + 5x + 4x + 10\right) \amp= 3\left(2x^2 + 9x + 10\right) Now we can complete the multiplication just by distributing that $3$ . 3(x + 2)(2x + 5) \amp= 6x^2 + 27x + 30

## Multiplying Polynomials Larger Than Binomials

To multiply polynomials that have more than two terms, we can use repeated distribution and
      monomial multiplication. Whether we are working with binomials, trinomials, or even larger
      polynomials, the process is fundamentally the same.

**Example**

