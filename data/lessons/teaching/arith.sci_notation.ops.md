> Content sourced from [ORCCA](https://pcc.edu/orcca) (Open Resources for Community College Algebra) — CC BY 4.0

# Scientific Notation

Very large and very small numbers can be awkward to write and calculate with. These kinds of numbers can show in the sciences. For example in biology, a human hair might be as thick as $0.000181$ meters. And the closest that Mars gets to the sun is $206620000$ meters. Keeping track of the decimal places and extra zeros raises the potential for mistakes to be made. In this section, we discuss a format used for very large and very small numbers called *scientific notation* that helps alleviate the issues with these numbers.

## The Basics of Scientific Notation

An October 3, 2016 CBS News headline read:

> Federal Debt in FY 2016 Jumped $$\$1{,}422{,}827{,}047{,}452.46$$---that's $$\$12{,}036$$ Per Household.

The article also later states:

> By the close of business on Sept. 30, 2016, the last day of fiscal 2016, it had climbed to $$\$19{,}573{,}444{,}713{,}936.79$$.

When presented in this format, trying to comprehend the value of these numbers can be overwhelming. More commonly, such numbers would be presented in a descriptive manner:

- The federal debt climbed by $1.42$ trillion dollars in 2016.
- The federal debt was $19.6$ trillion dollars at the close of business on Sept. 30, 2016.

In science, government, business, and many other disciplines, it's not uncommon to deal with very large numbers like these. When numbers get this large, it can be hard to discern when a number has eleven digits and when it has twelve.

We have descriptive language for all numbers based on the place value of the different digits: ones, tens, thousands, ten thousands, etc. We tend to rely upon this language more when we start dealing with larger numbers. Here's a chart for some of the most common numbers we see and use in the world around us:

*Whole Number Powers of $10$*

Each number above has a corresponding power of ten and this power of ten will be important as we start to work with the content in this section. This descriptive language also covers even larger numbers: trillion, quadrillion, quintillion, sextillion, septillion, and so on. There's also corresponding language to describe very small numbers, such as thousandth, millionth, billionth, trillionth, etc.

Through centuries of scientific progress, humanity became increasingly aware of very large numbers and very small measurements. As one example, the star that is nearest to our sun is . Proxima Centauri is about $25{,}000{,}000{,}000{,}000$ miles from our sun. Again, many will find the descriptive language easier to read: Proxima Centauri is about $25$ trillion miles from our sun.

To make computations involving such numbers more manageable, a standardized notation called "scientific notation" was established. The foundation of scientific notation is the fact that multiplying or dividing by a power of $10$ will move the decimal point of a number so many places to the right or left, respectively. So first, let's take a moment to review that level of basic arithmetic.

**Exercise**

Multiplying a number by $10^n$ where $n$ is a positive integer had the effect of moving the decimal point $n$ places to the right.

Every number can be written as a product of a number between $1$ and $10$ and a power of $10$. For example, $650 = 6.5 \times 100$. Since $100 = 10^2$, we can also write

$$
 650 = 6.5 \times 10^{2} 
$$

and this is our first example of writing a number in scientific notation.

**Definition**

A positive number is written in *scientific notation* when it has the form $a \times 10^n$ where $n$ is an integer and $1 \le a \lt 10 $. In other words, $a$ has precisely one nonzero digit to the left of the decimal place. The exponent $n$ used here is called the number's *order of magnitude*. The number $a$ is sometimes called the *significand* or the *mantissa*.

Some conventions do not require $a$ to be between $1$ and $10$, excluding both values, but that is the convention used in this book.

Some calculators and computer readouts cannot display exponents in superscript. In some cases, these devices will display scientific notation in the form `6.5E2` instead of $6.5\times10^2$.

## Scientific Notation for Large Numbers

To write a number larger than $10$ in scientific notation, like $89412$, first write the number with the decimal point right after its first digit, like $8.9412$. Now count how many places there are between where the decimal point originally was and where it is now.

$$
8.\overbrace{9412}^{4}
$$

Use that count as the power of $10$. In this example, we have

$$
89412=8.9412\times10^4
$$

Scientific notation communicates the "essence" of the number ($8.9412$) and then its size, or order of magnitude ($10^4$).

**Example**

To get a sense of how scientific notation works, let's consider familiar lengths of time converted to seconds.

Note that roughly $2.6$*million* seconds is one month, while roughly $2.5$*billion* seconds is an entire lifetime.

**Exercise**

**Exercise**

## Scientific Notation for Small Numbers

Scientific notation can also be useful when working with numbers smaller than $1$. As we saw in  Figure , we can represent thousands, millions, billions, trillions, etc., with positive integer exponents on $10$. We can similarly represent numbers smaller than $1$ (which are written as tenths, hundredths, thousandths, millionths, billionths, trillionths, etc.), with *negative* integer exponents on $10$. This relationship is outlined in  Figure .

*Negative Integer Powers of $10$*

To see how this works with a digit other than $1$, let's look at $0.005$. When we state $0.005$ as a number, we say "5 thousandths." Thus $0.005=5\times \frac{1}{1000}$. The fraction $\frac{1}{1000}$ can be written as $\frac{1}{10^3}$, which we know is equivalent to $10^{-3}$. Using negative exponents, we can then rewrite $0.005$ as $5\times10^{-3}$. This is the scientific notation for $0.005$.

In practice, we won't generally do that much computation. To write a small number in scientific notation we start as we did before and place the decimal point behind the first nonzero digit. We then count the number of decimal places between where the decimal had originally been and where it now is. Keep in mind that negative powers of ten are used to help represent very small numbers (smaller than $1$) and positive powers of ten are used to represent very large numbers (larger than $1$). So to convert $0.005$ to scientific notation, we have:

$$
 0\overbrace{.\highlight{005}}^{3}=5\times 10^{-3} 
$$

**Example**

In quantum mechanics, there is an important value called . Written as a decimal, the value of Planck's constant (rounded to six significant digits) is

$$
 0.000\,000\,000\,000\,000\,000\,000\,000\,000\,000\,000\,662\,607 
$$

.

In scientific notation, this number will be $6.62607\times 10^{\mathord{?}}$. To determine the exponent, we need to count the number of places from where the decimal originally is to where we will move it (following the first "6"):

$$
 0\overbrace{.\highlight{000\,000\,000\,000\,000\,000\,000\,000\,000\,000\,000\,6}}^{34\text{ places}}62\,607 
$$

So in scientific notation, Planck's Constant is $6.62607 \times 10^{-34}$. It will be much easier to use $6.62607 \times 10^{-34}$ in a calculation, and an added benefit is that scientific notation quickly communicates both the value and the order of magnitude of Planck's Constant.

**Exercise**

**Exercise**

**Exercise**

## Multiplying and Dividing Using Scientific Notation

One main reason for having scientific notation is to make calculations involving immensely large or small numbers easier to perform. By having the order of magnitude separated out in scientific notation, we can separate any calculation into two components.

**Example**

On Sept. 30th, 2016, the  was about $$\$19{,}600{,}000{,}000{,}000$$ and the US population was about $323{,}000{,}000$. What was the average debt per person that day?

1. Calculate the answer using the numbers provided, which are not in scientific notation.
2. First, confirm that the given values in scientific notation are $1.96 \times 10^{13}$ and $3.23 \times 10^8$. Then calculate the answer using scientific notation.

We've been asked to answer the same question, but to perform the calculation using two different approaches. In both cases, we'll need to divide the debt by the population.

1. We may need to use a calculator to handle such large numbers and we have to be careful that we type the correct number of `0`s.
   $$
   \frac{19600000000000}{323000000}\approx 60681.11
   $$
2. To perform this calculation using scientific notation, our work would begin by setting up the quotient as $\frac{1.96 \times 10^{13}}{3.23 \times 10^8}$. Dividing this quotient follows the same process we did with variable expressions of the same format, such as $\frac{1.96 w^{13}}{3.23 w^8}$. In both situations, we'll divide the coefficients and then use exponent properties to simplify the powers.
   $$
   \begin{aligned}\frac{1.96 \times 10^{13}}{3.23 \times 10^8} &= \frac{1.96 }{3.23} \times\frac{10^{13}}{ 10^8} \\ &\approx 0.6068111 \times 10^5 \\ &\approx 60681.11\end{aligned}
   $$

The federal debt per capita in the US on September 30th, 2016 was about $$\$60{,}681.11$$ per person. Both calculations give us the same answer, but the calculation relying upon scientific notation has less room for error and allows us to perform the calculation as two smaller steps.

Whenever we multiply or divide numbers that are written in scientific notation, we must separate the calculation for the coefficients from the calculation for the powers of ten, just as we simplified earlier expressions using variables and the exponent properties.

**Example**

1. Multiply $\left( 2\times 10^5 \right)\left( 3\times10^4 \right)$.
2. Divide $\dfrac{8\times 10^{17}}{4\times 10^2}$.

We will simplify the significand/mantissa parts as one step and then simplify the powers of $10$ as a separate step.

1. $$
\begin{aligned}[t] \left( 2\times 10^5 \right)\left( 3\times10^4 \right) &= \left( 2\times 3 \right)\times \left(10^5 \times 10^4 \right)\\ &= 6 \times 10^{9} \end{aligned}
$$
2. $$
\begin{aligned}[t] \frac{8 \times 10^{17}}{4\times 10^2} &= \frac{8}{4} \times \frac{10^{17}}{10^2}\\ &= 2 \times 10^{15} \end{aligned}
$$

Often when we multiply or divide numbers in scientific notation, the resulting value will not be in scientific notation. Suppose we were multiplying $\left( 9.3\times 10^{17} \right)\left( 8.2 \times 10^{-6} \right)$ and need to state our answer using scientific notation. We would start as we have previously:

$$
\begin{aligned}\left( 9.3\times 10^{17} \right)\left( 8.2 \times 10^{-6} \right) &=\left( 9.3\times 8.2 \right)\times \left( 10^{17} \times 10^{-6} \right) \\ &= 76.26 \times 10^{11}\end{aligned}
$$

While this is a correct value, it is not written using scientific notation. One way to convert this answer into scientific notation is to turn just the coefficient into scientific notation and momentarily ignore the power of ten:

$$
\begin{aligned}&=\highlight{76.26} \times 10^{11} \\ &= \highlight{7.626 \times 10^1} \times 10^{11}\end{aligned}
$$

Now that the coefficient fits into the proper format, we can combine the powers of ten and have our answer written using scientific notation.

$$
\begin{aligned}&=7.626 \times \highlight{10^1 \times 10^{11}} \\ &= 7.626 \times 10^{12}\end{aligned}
$$

**Example**

Multiply or divide as indicated. Write your answer using scientific notation.

1.
2.

Again, we'll separate out the work for the significand/mantissa from the work for the powers of ten. If the resulting coefficient is not between $1$ and $10$, we'll need to adjust that coefficient to put it into scientific notation.

1. $$
\begin{aligned}[t] \left( 8 \times 10^{21} \right)\left( 2 \times 10^{-7} \right) &= \left( 8 \times 2 \right)\times\left( 10^{21} \times 10^{-7} \right)\\ &= \highlight{16} \times 10^{14}\\ &= \highlight{1.6\times 10^1} \times 10^{14}\\ &= 1.6 \times 10^{15} \end{aligned}
$$

   We need to remember to apply the  product property  for exponents to the powers of ten.
2. $$
\begin{aligned}[t] \frac{ 2 \times 10^{-6} }{ 8 \times 10^{-19} } &= \frac{ 2 }{ 8 }\times\frac{ 10^{-6} }{ 10^{-19} }\\ &= \highlight{0.25} \times 10^{13}\\ &= \highlight{2.5\times 10^{-1}} \times 10^{13}\\ &= 2.5 \times 10^{12} \end{aligned}
$$

There are times where we will have to raise numbers written in scientific notation to a power. For example, suppose we have to find the area of a square whose radius is $3\times 10^7$ feet. To perform this calculation, we first remember the formula for the area of a square, $A=s^2$ and then substitute $3\times 10^7$ for $s$: $A = \left( 3\times 10^7 \right)^2$. To perform this calculation, we'll need to remember to use the  product to a power property  and the  power to a power property :

$$
\begin{aligned}A &= \left( 3\times 10^7 \right)^2 \\ &= \left( 3\right)^2 \times \left(10^7 \right)^2 \\ &= 9 \times 10^{14}\end{aligned}
$$

##

**Exercise**

Which number is very large and which number is very small?

$$
 9.99\times10^{-47}\qquad1.01\times10^{23} 
$$

**Exercise**

Since some computer/calculator screens can't display an exponent, how might a computer/calculator display the number $2.318\times10^{13}$?

**Exercise**

Why do we bother having scientific notation for numbers?

##

