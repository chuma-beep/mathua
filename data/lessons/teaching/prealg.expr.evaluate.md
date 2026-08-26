> Content sourced from [ORCCA](https://pcc.edu/orcca) (Open Resources for Community College Algebra) — CC BY 4.0

# Variables and Evaluating Expressions

Variables and expressions are the basic building blocks for writing algebra. In this section, we explore how to use them.

## Introduction to Variables

When we want to represent an unknown quantity or a quantity whose value can change, we use a *variable*. For example, if you'd like to write about automobile gas mileage, you  could use the symbol "$g$" as a variable to represent a car's gas mileage. The gas mileage $g$ might be $25\,\text{milepergallon}$ (miles per gallon) for one car, $30\,\text{milepergallon}$ for some other car, or other values for other cars. It might be one thing for your car when it was new, and something else ten years later.

Since we are using a variable, we can discuss gas mileage for Honda Civics, Ford Explorers, and all other makes and models at the same time, even though these makes and models each have their own gas mileage.

When variables stand for physical quantities, it's good to use letters that clearly represent those quantities. For example, it is wise to use $g$ for **g**as mileage. This helps people who read your mathematical writing understand it better. It is common to use $x$, $y$, and $z$ for variables when there is no context to suggest something more meaningful like $g$. You may see the variable $x$ a lot.

It is important to be clear about what unit of measure goes with a variable. With gas mileage $g$, if we all agree to use $\text{milepergallon}$ for its units, then $g$ might be a placeholder for $25$, $30$, etc. On the other hand if we decide to use $\text{kilometerpergallon}$ (kilometers per gallon) for units, those quantities would be $40$, $48$, etc. So it's important to tell readers that $g$ represents gas mileage in *miles per gallon* or *kilometers per gallon* or whatever the case may be.

Sometimes the units we should use for a variable are suggested indirectly. For example if we're told that a car has used so many *gallons* of gas after traveling so many *miles*, then we should measure gas mileage in $\text{milepergallon}$, not $\text{kilometerpergallon}$.

\*\*Exercise\*\*

Naming Variables

Remi is studying college student demographics, including their ages. Let  be the age of a student, measured in .

Solution

The unknown quantity is age, which we generally measure in years. So we could say "Let $a$ be the age of a student, measured in years."

Luka needs to drive from Portland, OR to Boise, ID. Let  be the amount of time passed since Luka left Portland, measured in .

Solution

The amount of time passed is the unknown quantity. Since this is a drive from Portland to Boise, it makes sense to measure this in hours, not minutes or weeks. So we could say "Let $t$ be the amount of time passed since Luka left Portland, measured in hours."

There is a number that you will triple, then add five to, and then take the square root. Let  be this number.

Solution

The number we are discussing doesn't have any physical context, so we choose to use the generic variable $x$ to represent it.

## Algebraic Expressions

Any combination of variables and numbers using arithmetic operations (like addition,  multiplication, etc.) is an *algebraic expression*. The following are examples of  algebraic expressions:

$$
 x+1\qquad 2\ell+2w\qquad\frac{\sqrt{x}}{y+1}\qquad nRT 
$$

Note that this definition of "algebraic expression" does *not* include anything with an equal sign ($=$) in it. The idea of an *equation* (which has an equals sign) is discussed in  the section .

\*\*Example\*\*

The expression:

$$
 \frac{5}{9}(F - 32) 
$$

converts a temperature in degrees Fahrenheit to degrees Celsius. To do this, we need a Fahrenheit temperature, $F$. Then we can *evaluate* the expression. This  means replacing its variable(s) with specific numbers and finding the result as a single,  simplified number.

Let's convert the temperature $89^{\circ}F$ to the Celsius scale by evaluating the expression. To do this, we *substitute* the number $89$ in place of the  variable $F$.

$$
\begin{aligned}\frac{5}{9}(\substitute{89} - 32)&=\frac{5}{9}(57) &&\text{Review order of operations in }\text{.} \\ &=\frac{285}{9} &&\text{Review fraction multiplication in }\text{.} \\ &\approx 31.67\end{aligned}
$$

This shows us that $89^{\circ}F$ is approximately the same as $31.67^{\circ}C$.

\*\*Warning\*\*

Vocabulary

The steps in  the example  are not "solving", as far as algebra vocabulary is concerned. "Solving" is a word you might want to use because in everyday English you are "finding an answer". However in algebra, there is a special meaning for the term "solving" that is discussed in  the section . Here, when we substitute in values for variables and then compute the result, we are "evaluating an expression", not solving anything.

\*\*Exercise\*\*

Convert Temperature

Try evaluating the temperature expression for yourself.

If a temperature is $50^{\circ}F$, what is that temperature in Celsius?

Solution

$$
\begin{aligned}[t] \frac{5}{9}(\substitute{50} - 32)&=\frac{5}{9}(18)\\ &=\frac{5}{1}(2)=10 \end{aligned}
$$

So $50^{\circ}F$ is equivalent to $10^{\circ}C$.

If a temperature is $-20^{\circ}F$, what is that temperature in Celsius?

Solution

$$
\begin{aligned}[t] \frac{5}{9}(\substitute{-20} - 32)&=\frac{5}{9}(-52)\\ &=-\frac{260}{9}\approx-28.89 \end{aligned}
$$

So $-20^{\circ}F$ is equivalent to about $-28.89^{\circ}C$.

\*\*Example\*\*

Stair Rise and Run

When building a staircase, you want the rise and run to be consistent from one step to the next.

A convention among contractors is that a staircase run, in inches, is given by $17.5-h$ where $h$ is the rise in inches.

Determine the run for each step of a staircase where the rise is $7\,\text{in}$.

We substitute $7$ for $h$:

$$
\begin{aligned}17.5-\substitute{7}&=10.5\end{aligned}
$$

. So the run is $10.5\,\text{in}$.

A staircase needs to span a total height of $108\,\text{in}$. What is a reasonable number of steps for it to have? What will that mean for the rise of each step? What will the run be for each step?

There is more than one good answer, but if there are $12$ steps then the height of each step will be $\frac{108}{12}$ inches, or $9\,\text{in}$. And that mean the run of each step is found by substituting $9$ for $h$:

$$
\begin{aligned}17.5-\substitute{9}&=8.5\end{aligned}
$$

. So each step would have a run of $8.5\,\text{in}$.

\*\*Exercise\*\*

Stair Rise and Run

\*\*Exercise\*\*

Rising Rents

According to this model, what was the median rent for a living unit in Oregon in 2010?

Solution

According to this model, what was the median rent for a living unit in Oregon in 2020?

Solution

According to this model, what will be the median rent for a living unit in Oregon in 2030?

Solution

## Evaluating Expressions with Exponents, Absolute Value, and Radicals

Algebraic expressions might have exponents, absolute value bars, and radicals. This does not change the basic approach to evaluating them.

\*\*Example\*\*

Tsunami Speed

The speed of a tsunami (in meters per second) can be modeled by $\sqrt{9.8d}$, where $d$ is the depth of the tsunami (in meters). Determine the speed of a tsunami that has a depth of $30\,\text{m}$ to four significant digits.

Using $d=30$, we find:

$$
\begin{aligned}\sqrt{9.8(\substitute{30})}&=\sqrt{294}&&\text{Review order of operations in }\text{.} \\ &\approx \overbrace{17.14}^{\text{four}}6428\ldots&&\text{Review square root in }\text{.}\end{aligned}
$$

The speed of tsunami with a depth of $30\,\text{m}$ is about $17.15\,\text{m}/$.

We have been evaluating expressions, but we can evaluate formulas in the same way. A *formula* has an equal sign with an expression to the right. On the left of the  equal sign, there is a variable that represents the result. For example, we could write the formula $s=\sqrt{9.8d}$ for the speed of a tsunami from  the example .

\*\*Exercise\*\*

Tent Height

When you are $5$ $\text{ft}$ from the west side, the height is .

Solution

When you are $2.5$ $\text{ft}$ from the west side, the height is .

Solution

\*\*Exercise\*\*

Mortgage Payments

\*\*Warning\*\*

Rounding Too Much

You might have noticed in the explanation to  the exercise  that during the computations, many decimal places were recorded at each step. Tracking lots of decimal places might be important, depending on what you are working toward. If you round in the middle of your work, you have changed the numbers a little bit from what they *really* should be. As computation continues, this little error can become larger and larger, leaving you with a final result that is too far off from correct. So the best practice is to keep lots of decimal places in all your computations, and then at the very end you may round more if that is appropriate.

## Evaluating Expressions with Negative Numbers

When we substitute negative numbers into an expression, it's important to use parentheses  around them or else it's easy to forget that a *negative* number is being raised to a power.

\*\*Example\*\*

Evaluate $x^2$ for $x=-2$.

We substitute:

$$
\begin{aligned}x^2&=(\substitute{-2})^2 \\ &=4 \\ x^2&=-2^2\qquad\text{incorrect!} \\ &=-4\end{aligned}
$$

The original expression $x^2$ takes $x$ and squares it, so we want to do the same thing to the number $-2$. But with the incorrect expression $-2^2$, the number $-2$ is not being squared. An exponent has higher priority than negation in the order of operations, so $-2^2$ is the same as $-\left(2^2\right)$, and the wrong number is being squared. With $(-2)^2$ the number $-2$ is being squared, which is what we want.

So it's wise to always use parentheses when substituting in a negative number.

\*\*Exercise\*\*

Multivariable Expressions

$x^3y^2$

Solution

$$
\begin{aligned}[t] (\substitute{-2})^3(\substitute{-5})^2&=(-8)(25)\\ &=-200 \end{aligned}
$$

$(-2x)^3$

Solution

$$
\begin{aligned}[t] (-2(\substitute{-2}))^3&=(4)^3\\ &=64 \end{aligned}
$$

$-3x^2y$

Solution

$$
\begin{aligned}[t] -3(\substitute{-2})^2(\substitute{-5})&=-3(4)(-5)\\ &=60 \end{aligned}
$$

##

\*\*Exercise\*\*

What is a reason for wanting to use a letter other than $x$, $y$, or $z$ as a variable?

\*\*Exercise\*\*

What is the difference between an "algebraic expression" and a "formula", as these things were described in this section? (Other math resources may define these terms differently.)

\*\*Exercise\*\*

What should you watch out for when substituting in a negative number for a variable?

##

