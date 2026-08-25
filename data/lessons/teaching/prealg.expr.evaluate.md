> Content sourced from [ORCCA](https://pcc.edu/orcca) (Open Resources for Community College Algebra) — CC BY 4.0
# Variables and Evaluating Expressions

Variables and expressions are the basic building blocks for writing algebra. In this section, we explore how to use them.

## Introduction to Variables

When we want to represent an unknown quantity or a quantity whose value can change, we use a *variable*. For example, if you'd like to write about automobile gas mileage, you could use the symbol "g" as a variable to represent a car's gas mileage. The gas mileage $g$ might be $25\,\text{milepergallon}$ (miles per gallon) for one car, $30\,\text{milepergallon}$ for some other car, or other values for other cars. It might be one thing for your car when it was new, and something else ten years later.

Since we are using a variable, we can discuss gas mileage for Honda Civics, Ford Explorers, and all other makes and models at the same time, even though these makes and models each have their own gas mileage.

When variables stand for physical quantities, it's good to use letters that clearly represent those quantities. For example, it is wise to use $g$ for g as mileage. This helps people who read your mathematical writing understand it better. It is common to use $x$, $y$, and $z$ for variables when there is no context to suggest something more meaningful like $g$. You may see the variable $x$ a lot.

It is important to be clear about what unit of measure goes with a variable. With gas mileage $g$, if we all agree to use $\text{milepergallon}$ for its units, then $g$ might be a placeholder for $25$, $30$, etc. On the other hand if we decide to use $\text{kilometerpergallon}$ (kilometers per gallon) for units, those quantities would be $40$, $48$, etc. So it's important to tell readers that $g$ represents gas mileage in *miles per gallon* or *kilometers per gallon* or whatever the case may be.

Sometimes the units we should use for a variable are suggested indirectly. For example if we're told that a car has used so many *gallons* of gas after traveling so many *miles*, then we should measure gas mileage in $\text{milepergallon}$, not $\text{kilometerpergallon}$.

**Naming Variables**
Remi is studying college student demographics, including their ages. Let be the age of a student, measured in.

*Solution*
The unknown quantity is age, which we generally measure in years. So we could say "Let a be the age of a student, measured in years."

Luka needs to drive from Portland, OR to Boise, ID. Let be the amount of time passed since Luka left Portland, measured in.

*Solution*
The amount of time passed is the unknown quantity. Since this is a drive from Portland to Boise, it makes sense to measure this in hours, not minutes or weeks. So we could say "Let t be the amount of time passed since Luka left Portland, measured in hours."

There is a number that you will triple, then add five to, and then take the square root. Let be this number.

*Solution*
The number we are discussing doesn't have any physical context, so we choose to use the generic variable $x$ to represent it.

## Algebraic Expressions

Any combination of variables and numbers using arithmetic operations (like addition, multiplication, etc.) is an *algebraic expression*. The following are examples of algebraic expressions: 
\[x+1\qquad 2\ell+2w\qquad\frac{\sqrt{x}}{y+1}\qquad nRT\]
 Note that this definition of "algebraic expression" does *not* include anything with an equal sign ( $=$) in it. The idea of an *equation* (which has an equals sign) is discussed in.

**Example**
The expression: 
\[\frac{5}{9}(F - 32)\]
 converts a temperature in degrees Fahrenheit to degrees Celsius. To do this, we need a Fahrenheit temperature, $F$. Then we can *evaluate* the expression. This means replacing its variable(s) with specific numbers and finding the result as a single, simplified number. Let's convert the temperature $89^{\circ}F$ to the Celsius scale by evaluating the expression. To do this, we *substitute* the number $89$ in place of the variable $F$. 
\[ \begin{aligned} \frac{5}{9}(89 - 32)&=\frac{5}{9}(57) &&\text{Review order of operations in }\text{.} \\ &=\frac{285}{9} &&\text{Review fraction multiplication in }\text{.} \\ &\approx 31.67 \end{aligned} \]
 This shows us that $89^{\circ}F$ is approximately the same as $31.67^{\circ}C$.

> **Vocabulary**
> The steps in are not "solving", as far as algebra vocabulary is concerned. "Solving" is a word you might want to use because in everyday English you are "finding an answer". However in algebra, there is a special meaning for the term "solving" that is discussed in. Here, when we substitute in values for variables and then compute the result, we are "evaluating an expression", not solving anything.

**Convert Temperature**
Try evaluating the temperature expression for yourself.
If a temperature is $50^{\circ}F$, what is that temperature in Celsius?

*Solution*
$\begin{aligned}[t] \frac{5}{9}(50 - 32)\amp=\frac{5}{9}(18)\\ \amp=\frac{5}{1}(2)=10 \end{aligned}$ So $50^{\circ}F$ is equivalent to $10^{\circ}C$.

If a temperature is $-20^{\circ}F$, what is that temperature in Celsius?

*Solution*
$\begin{aligned}[t] \frac{5}{9}(-20 - 32)\amp=\frac{5}{9}(-52)\\ \amp=-\frac{260}{9}\approx-28.89 \end{aligned}$ So $-20^{\circ}F$ is equivalent to about $-28.89^{\circ}C$.

**Stair Rise and Run**

**Stair Rise and Run**

**Rising Rents**
According to this model, what was the median rent for a living unit in Oregon in 2010?

*Solution*
This model uses $x$ as the number of years since 2000. So for the year 2010, $x$ is $10$: $\begin{aligned} 620+32.35(10)& = 943.50 \end{aligned}$ According to this model, the median monthly rent for a living unit in Oregon in 2010 was $\$943.50$.

According to this model, what was the median rent for a living unit in Oregon in 2020?

*Solution*
For the year 2020, $x$ is $20$: $\begin{aligned} 620+32.35(20)& = 1267 \end{aligned}$ According to this model, the median monthly rent for a living unit in Oregon in 2020 was $\$1267$.

According to this model, what will be the median rent for a living unit in Oregon in 2030?

*Solution*
For the year 2030, $x$ is $30$: $\begin{aligned} 620+32.35(30)& = 1590.50 \end{aligned}$ According to this model, the median monthly rent for a living unit in Oregon in 2030 will be $\$1590.50$.

## Evaluating Expressions with Exponents, Absolute Value, and Radicals

Algebraic expressions might have exponents, absolute value bars, and radicals. This does not change the basic approach to evaluating them.

**Tsunami Speed**
The speed of a tsunami (in meters per second) can be modeled by $\sqrt{9.8d}$, where $d$ is the depth of the tsunami (in meters). Determine the speed of a tsunami that has a depth of $30\,\text{m}$ to four significant digits.

*Solution*
Using $d=30$, we find: 
\[ \begin{aligned} \sqrt{9.8(30)}&=\sqrt{294}&&\text{Review order of operations in }\text{.} \\ &\approx \overbrace{17.14}^{\text{four}}6428\ldots&&\text{Review square root in }\text{.} \end{aligned} \]
 The speed of tsunami with a depth of $30\,\text{m}$ is about $17.15\,\text{m}/$.

We have been evaluating expressions, but we can evaluate formulas in the same way. A *formula* has an equal sign with an expression to the right. On the left of the equal sign, there is a variable that represents the result. For example, we could write the formula $s=\sqrt{9.8d}$ for the speed of a tsunami from.

**Tent Height**
When you are $5$ $\text{ft}$ from the west side, the height is.

*Solution*
When $d=5$, we have: $\begin{aligned} h&= -2\abs{5-3}+6&&\text{Review order of operations in }\text{.} \\ &= -2\abs{2}+6&&\text{Review absolute value in }\text{.} \\ &= -2(2)+6 \\ &= -4+6=2 \end{aligned}$ So when you are $5$ $\text{ft}$ from the west side, the height of the tent is $2$ $\text{ft}$.

When you are $2.5$ $\text{ft}$ from the west side, the height is.

*Solution*
When $d=2.5$, we have: $\begin{aligned} h&= -2\abs{2.5-3}+6 \\ &= -2\abs{-0.5}+6 \\ &= -2(0.5)+6 \\ &=-1+6=5 \end{aligned}$ So when you are $2.5$ $\text{ft}$ from the west side, the height of the tent is $5$ $\text{ft}$.

**Mortgage Payments**

> **Rounding Too Much**
> You might have noticed in the explanation to that during the computations, many decimal places were recorded at each step. Tracking lots of decimal places might be important, depending on what you are working toward. If you round in the middle of your work, you have changed the numbers a little bit from what they *really* should be. As computation continues, this little error can become larger and larger, leaving you with a final result that is too far off from correct. So the best practice is to keep lots of decimal places in all your computations, and then at the very end you may round more if that is appropriate.

## Evaluating Expressions with Negative Numbers

When we substitute negative numbers into an expression, it's important to use parentheses around them or else it's easy to forget that a *negative* number is being raised to a power.

**Example**

**Multivariable Expressions**
$x^3y^2$

*Solution*
$\begin{aligned}[t] (-2)^3(-5)^2\amp=(-8)(25)\\ \amp=-200 \end{aligned}$

$(-2x)^3$

*Solution*
$\begin{aligned}[t] (-2(-2))^3\amp=(4)^3\\ \amp=64 \end{aligned}$

$-3x^2y$

*Solution*
$\begin{aligned}[t] -3(-2)^2(-5)\amp=-3(4)(-5)\\ \amp=60 \end{aligned}$

