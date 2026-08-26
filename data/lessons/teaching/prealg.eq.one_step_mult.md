> Content sourced from [ORCCA](https://pcc.edu/orcca) (Open Resources for Community College Algebra) — CC BY 4.0

# Solving One-Step Equations

In  the section  we learned how to check whether or not a give number is a solution to an equation. But the number to check was always given to us. In this section, we begin learning how to find solution(s) ourselves.

## Imagine Filling in the Blanks

Let's start with a simple situation---so simple that you might not need algebra, but the example serves as a warm-up.

\*\*Example\*\*

A number plus $2$ is $6$. What is that number?

You may be so familiar with arithmetic that you know the answer already. The *algebra* approach is to translate "A number plus $2$ is $6$" into an equation:

$$
 x+2=6 
$$

where $x$ is the number we are trying to find. How do we find the value for $x$ that makes the equation true?

One valid option is to *imagine* what number you could put in place of $x$ that would result in a true equation.

- Would $0$ work? No, that would mean $0+2=6$, which is false.
- Would $17$ work? No, that would mean $17+2=6$, which is false.
- Would $4$ work? Yes, because $4+2=6$ is a true equation.

So one solution to the equation is $4$. No other numbers can be a solution, because when you add $2$ to something smaller than $4$, the result is smaller than $6$. And when you add $2$ to something larger than $4$, the result is larger than $6$.

This approach ("imagining" what number works in the equation) might work for you to solve *very basic* equations. It's not going to work in general with more complicated equations. And so we move on to a more systematic approach that you can use all the time.

## The Basic Principle of Algebra

Let's revisit  the example , but think it through differently.

\*\*Example\*\*

If a number plus $2$ is $6$, what is the number?

If a number *plus*$2$ equals $6$, then the number is a little smaller than $6$. We should be able to *subtract*$2$ from $6$ to get that unknown number. (We are using the *opposite* operation from addition, which is subtraction.)

*Add $2$ to mystery number to get $6$, or subtract $2$ from $6$ to get mystery number*

Doing that subtraction: $6-2=4$.

Let's try this strategy with another example.

\*\*Example\*\*

If a number *minus*$2$ equals $6$, what is that number? The mystery number must be a little larger than $6$. The opposite of subtraction is addition, so if we *add*$2$ to $6$ we will find the unknown number. So the unknown number is $6+2=8$.

Does this strategy work with multiplication and division?

\*\*Example\*\*

If a number multiplied by $2$ makes $6$, what is that number? The mystery number is small, since it gets multiplied by $2$ to make $6$. If we *divide*$6$ by $2$, we will find the unknown number. Note that division is the opposite action of multiplication.

*Multiply mystery number by $2$ to get $6$, or divide $6$ by $2$ to get mystery number*

So the unknown number is $\frac{6}{2}=3$.

\*\*Example\*\*

If a number divided by $2$ equals $6$, what is the number? We must be starting with a larger number, since cutting it in half makes $6$. If we *multiply*$6$ by $2$ (because multiplying is the opposite of dividing) then we find the unknown number is $6\cdot2=12$.

These examples explore an important principle for solving an equation---applying an opposite arithmetic operation. We can revisit  the example  and apply this strategy with more care and intentionality. If a number plus $2$ is $6$, what is the number? We will use $x$ to represent the unknown number. The question translates into the math equation

$$
 x+2=6 
$$

. Picture the equal sign as the middle of a balanced scale. The left side has a brick labeled "$x$ lb" together with $2$ one-pound bricks. So the total weight on the left is $x+2$. The right side has $6$ one-pound bricks.  the figure  (left side) shows the scale.

*Balance scale representing $x+2=6$, and the solution after taking away $2$ from each side*

To find the weight of the unknown brick, we can take away $2$ one-pound bricks from *each* side of the scale and that will keep the scale balanced.  the figure  (right side) shows the solution.

An equation is like a balanced scale: the two sides of the equation are equal, and the two sides of the balanced scale have equal weight. Just like we can take away $2$ lb from each side of a balanced scale, we can subtract $2$ from each side of an equation. Instead of drawing two pictures of balance scales, we can use algebra symbols and solve the equation $x+2=6$ in the following way:

$$
\begin{aligned}x+2&=6&&\text{like a balanced scale} \\ x+2\subtractright{2}&=6\subtractright{2}&&\text{remove the same quantity from each side} \\ x&=4&&\text{still balanced; now it straight up tells you the solution}\end{aligned}
$$

Each line of the algebra above shows what is called an *equivalent equation*.  Each of those equations is "algebraically equivalent" to the one that came before it,  meaning it has exactly the same solution(s). The final equivalent equation $x=4$ tells us directly that the solution to the equation is $4$.

In theory, there could have been more than one solution to the equation (although that is not the case with this equation). So conceptually, there is a *collection* of solutions to any given equation. We call this collection a *solution set*.  Any set of numbers that only has one or a few numbers in it  is written using curly braces. In this case, the solution set is $\{4\}$. Using braces to surround a collection of numbers listed out is called *set notation*, not to be confused with set-builder notation from  the section .

We have learned we can add or subtract the same number on each side of the equal sign, just like we can add or remove the same amount of weight on a balanced scale. Can we multiply and divide the same number on each side of the equal sign? Let's look at  the example  again: If a number times $2$ is $6$, what is the number? Another balance scale can help visualize this.

*Balance scale representing $2x=6$, and the solution after taking away half from each side*

At first, the scale is balanced. If we cut the weight in half on both sides, it should still be balanced. We can see from the scale that $x=3$ is correct.

Removing half of the weight from each side of the scale is like dividing both sides of an equation by $2$:

$$
\begin{aligned}2x&=6 \\ \divideunder{2x}{2}&=\divideunder{6}{2} \\ x&=3\end{aligned}
$$

The equivalent equation in this example is $x=3$, which tells us that the solution to the equation is $3$ and the solution set is $\{3\}$.

\*\*Remark\*\*

Note that when we divide each side of an equation by a number, we use a fraction bar, not a division symbol. The equation $\divideunder{2x}{2}=\divideunder{6}{2}$ could be written as $2x\divideright{2}=6\divideright{2}$, but algebra tends to avoid using the $\div$ symbol. In part, this is because when writing by hand, it might be confused with a subtraction sign.

Similarly, we can multiply each side of an equation by $2$ if that is helpful, and it will keep a scale in balance. We can summarize these properties.

\*\*Fact\*\*

Properties of Equivalent Equations

If there is an equation $\text{Left}=\text{Right}$, we can do the following to obtain an equivalent equation.

$$
\text{Left}\addright{c}=\text{Right}\addright{c}
$$

(add the same number to each side)

$$
\text{Left}\subtractright{c}=\text{Right}\subtractright{c}
$$

(subtract the same number from each side)

$$
\text{(Left)}\multiplyright{c}=\text{(Right)}\multiplyright{c}
$$

(multiply each side of the equation by the same *nonzero* number)

$$
\displaystyle\divideunder{\text{Left}}{c}=\divideunder{\text{Right}}{c}
$$

(divide each side of the equation by the same *nonzero* number)

## Solving One-Step Equations and Stating Solution Sets

Notice when we solved equations in  the subsection , the final equation looked like $x=\text{number}$, where the variable $x$ stands alone on one side of the equal sign. The goal of solving any equation is to *isolate the variable* in this same manner.

Putting together both strategies (applying the opposite operation and balancing equations like a scale) that we just explored, we summarize how to solve a one-step linear equation.

Steps to Solving Simple (One-Step) Linear Equations



Let's look at a few examples.

\*\*Example\*\*

Solve for $y$ in the equation $7+y=3$.

To isolate $y$, we need to remove $7$ from the left side. Since $7$ is being *added* to $y$, we need to *subtract*$7$ from each side of the equation.

$$
\begin{aligned}7+y&=3 \\ 7+y\subtractright{7}&=3\subtractright{7} \\ y&=-4\end{aligned}
$$

We should check the solution. To do that, substitute $-4$ in for $y$ in the original equation:

$$
\begin{aligned}7+y&=3 \\ 7+(\substitute{-4})&\wonder{=}3 \\ 3&\confirm{=}3\end{aligned}
$$

The solution $-4$ is checked, and the solution set is $\{-4\}$.

\*\*Exercise\*\*

\*\*Exercise\*\*

Note that when solving the equation in  the exercise  we found $-3=a$, and did not bother to write it the other way round as $a=-3$. All that really matters is that we ended with a clear statement of the solution set, which was $\{-3\}$.

\*\*Example\*\*

The formula for a circle's circumference is $c=\pi d$, where $c$ represents circumference, $d$ represents diameter, and $\pi$ is a constant with the value of $3.1415926\ldots$.

If a circle's circumference is $12\pi\,\text{ft}$, find the circle's diameter.

The circumference is given as $12\pi$ feet, so we will substitute $c$ in the formula with $12\pi$, giving the equation $12\pi=\pi d$. Now we will solve for $d$:

$$
\begin{aligned}12\pi&=\pi d \\ \divideunder{12\pi}{\pi}&=\divideunder{\pi d}{\pi} \\ 12&=d\end{aligned}
$$

We should check the solution by substituting $12$ in for $d$ in the original equation:

$$
\begin{aligned}12\pi&=\pi d \\ 12\pi&\wonder{=}\pi \substitute{(12)} \\ 12\pi&\confirm{=}12\pi\end{aligned}
$$

This checks out, so the circle's diameter is $12\,\text{ft}$.

\*\*Exercise\*\*

Examples so far have solved an equation by undoing addition, subtraction, multiplication, or division. There is one last arithmetic action that we will look into undoing: *negation*. Negation is when you apply the negative sign to a number.  Undoing negation is simple though: just negate again. For example, $-(-42)=42$.

\*\*Example\*\*

Solve the equation $-b=2$ for $b$.

Our variable $b$ is not yet isolated because of the negative sign in front. To attack that negative sign, we can "negate" each side:

$$
\begin{aligned}-b&=2 \\ \negate{-b}&=\negate{2} \\ b&=-2\end{aligned}
$$

We removed the negative sign from $-b$ by negating both sides, and we can see that the solution set is $\{-2\}$.

An alternative is to think of the original negative sign as multiplication by $-1$. In that case, dividing each side by $-1$ will successfully isolate $b$.

$$
\begin{aligned}-b&=2 \\ -1\cdot b&=2 \\ \divideunder{-1\cdot b}{-1}&=\divideunder{2}{-1} \\ b&=-2\end{aligned}
$$

Another alternative is to recognize that multiplying on each side by $-1$ will also cancel that unwanted negative sign.

$$
\begin{aligned}-b&=2 \\ \multiplyleft{-1}(-b)&=\multiplyleft{-1}(2) \\ b&=-2\end{aligned}
$$

It is recommended that you review these three approaches for undoing negation, settle on the method that you like, and use it consistently.

We should check the solution by substituting $-2$ in for $b$ in the original equation:

$$
\begin{aligned}-b&=2 \\ -(\substitute{-2})&\confirm{=}2\end{aligned}
$$

The solution $-2$ is checked, and the solution set is $\{-2\}$.

## Equations with Fractions

When an equation has fractions, solving it uses the same principles. Of course you may need to use fraction arithmetic. Also, you might make good use of the reciprocal of a fraction as described in  the example .

\*\*Example\*\*

Solve the equation $\frac{2}{3}+g=\frac{1}{2}$ for $g$.

Since $\frac23$ is added to $g$, we will subtract $\frac23$ from each side.

$$
\begin{aligned}\frac{2}{3}+g&=\frac{1}{2} \\ \frac{2}{3}+g\subtractright{\frac{2}{3}}&=\frac{1}{2}\subtractright{\frac{2}{3}} \\ g&=\highlight{\frac{3}{6}}-\highlight{\frac{4}{6}} \\ g&=-\frac{1}{6}\end{aligned}
$$

We should check the solution by substituting $-\frac{1}{6}$ in for $g$:

$$
\begin{aligned}\frac{2}{3}+g&=\frac{1}{2} \\ \frac{2}{3}+\left(\substitute{-\frac{1}{6}}\right)&\wonder{=}\frac{1}{2} \\ \frac{4}{6}+\left(-\frac{1}{6}\right)&\wonder{=}\frac{1}{2} \\ \frac{3}{6}&\confirm{=}\frac{1}{2}\end{aligned}
$$

The solution $-\frac{1}{6}$ is checked, and the solution set is $\left\{-\frac{1}{6}\right\}$.

\*\*Exercise\*\*

When the variable in an equation is multiplied by a fraction, you can use the reciprocal of that fraction to help solve the equation. The *reciprocal* of a fraction is the fraction you get  from swapping the numerator and denominator. For example, the reciprocal of $\frac{2}{3}$ is $\frac{3}{2}$.

A reciprocal is useful because when a fraction is multiplied by its reciprocal, the result is $1$. For example, $\frac{2}{3}\cdot\frac{3}{2}=1$. This helps us remove a fraction when it is multiplied by the variable.

\*\*Example\*\*

Solve the equation $\frac{5}{8}d=7$ for $d$.

Our variable $d$ is multiplied by the fraction $\frac{5}{8}$. While we *could* divide on each side by $\frac{5}{8}$, that leads to a messy four-level equation: $\frac{\frac{5}{8}d}{\frac{5}{8}}=\frac{7}{\frac{5}{8}}$. To avoid this, we can just *multiply* on each side by the reciprocal of $\frac{5}{8}$.

$$
\begin{aligned}\frac{5}{8}d&=7 \\ \multiplyleft{\frac{8}{5}}\frac{5}{8}d&=\multiplyleft{\frac{8}{5}}7 \\ 1d&=\frac{8}{5}\cdot\highlight{\frac{7}{1}} \\ d&=\frac{56}{5}\end{aligned}
$$

We should check the solution by substituting $\frac{56}{5}$ in for $d$:

$$
\begin{aligned}\frac{5}{8}\substitute{\left(\frac{56}{5}\right)}&\wonder{=}7 \\ \frac{\cancel{5}}{8}\left(\frac{56}{\cancel{5}}\right)&\wonder{=}7 \\ \frac{1}{\cancel{8}}\left(\frac{\cancelto{7}{56}}{1}\right)&\wonder{=}7 \\ \frac{7}{1}&\confirm{=}7\end{aligned}
$$

The solution $\frac{56}{5}$ is checked, and the solution set is $\left\{\frac{56}{5}\right\}$.

\*\*Exercise\*\*

Sometimes the variable is in the numerator of a fraction, like in $\frac{3x}{4}$. This is actually the same as $\frac{3}{4}x$. Either way, $x$ is multiplied by $3$ and divided by $4$. So this is another situation where the reciprocal of a fraction can help.

\*\*Example\*\*

Solve the equation $\frac{3x}{4}=10$ for $x$.

We can multiply on each side by $\frac{4}{3}$ and that will isolate $x$:

$$
\begin{aligned}\frac{3x}{4}&=10 \\ \multiplyleft{\frac{4}{3}}\left(\frac{3x}{4}\right)&=\multiplyleft{\frac{4}{3}}(10) \\ \frac{4}{\cancel{3}}\left(\frac{\cancel{3}x}{4}\right)&=\frac{4}{3}\left(\highlight{\frac{10}{1}}\right) \\ \frac{\cancel{4}}{1}\left(\frac{x}{\cancel{4}}\right)&=\frac{40}{3} \\ x&=\frac{40}{3}\end{aligned}
$$

We should check the solution by substituting $\frac{40}{3}$ in for $x$ in the original equation:

$$
\begin{aligned}\frac{3x}{4}&=10 \\ \frac{3\left(\substitute{\frac{40}{3}}\right)}{4}&\wonder{=}10 \\ \frac{40}{4}&\wonder{=}10 \\ 10&\confirm{=}10\end{aligned}
$$

The solution $\frac{40}{3}$ is checked, and the solution set is $\left\{\frac{40}{3}\right\}$.

\*\*Exercise\*\*

##

\*\*Exercise\*\*

If you imagine the equation $2x+3=11$ as a balance scale with bricks on each side, how many bricks do you imagine are on the left side? How many *types* of brick do you imagine being on the left side?

\*\*Exercise\*\*

What is the opposite operation of multiplying by a negative number?

\*\*Exercise\*\*

Each time you solve an algebra equation, there is something you should be in the habit of doing at the end. Describe that thing you should do.

##

