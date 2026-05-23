> Content sourced from [ORCCA](https://pcc.edu/orcca) (Open Resources for Community College Algebra) — CC BY 4.0
# Solving One-Step Equations

In we learned how to check whether
      or not a give number is a solution to an equation. But the number to check was always given to
      us. In this section, we begin learning how to find solution(s) ourselves.

## Imagine Filling in the Blanks

Let's start with a simple situation so simple that you might not need algebra, but the
      example serves as a warm-up.

**Example**

This approach ( "imagining" what number works in the equation) might work for you to solve *very basic* equations. It's not going to work in general with more complicated
      equations. And so we move on to a more systematic approach that you can use all the time.

## The Basic Principle of Algebra

Let's revisit , but think it through differently.

**Example**

Let's try this strategy with another example.

**Example**

Does this strategy work with multiplication and division?

**Example**

**Example**

These examples explore an important principle for solving an equation applying an
      opposite arithmetic operation. We can revisit and apply this strategy with more care and intentionality.
      If a number plus $2$ is $6$ , what is the number?
      We will use $x$ to represent the unknown number.
      The question translates into the math equation \[
        x+2=6
      \] .
      Picture the equal sign as the middle of a balanced scale.
      The left side has a brick labeled "x lb" together with $2$ one-pound bricks.
      So the total weight on the left is $x+2$ . The right side has $6$ one-pound bricks. (left side) shows the scale.

*Abbreviation for pound*
Why is "lb" the abbreviation for "pound" ?
        It has a connection to a balance scale (⚖),
        which is the symbol for the Zodiac sign L i b ra.

To find the weight of the unknown brick, we can take away $2$ one-pound bricks from *each* side of the scale and that will keep the scale balanced. (right side) shows the solution.

An equation is like a balanced scale: the two sides of the equation are equal,
      and the two sides of the balanced scale have equal weight.
      Just like we can take away $2$ lb from each side of a balanced scale,
      we can subtract $2$ from each side of an equation.
      Instead of drawing two pictures of balance scales,
      we can use algebra symbols and solve the equation $x+2=6$ in the following way: x+2\amp=6\amp\amp\text{like a balanced scale} x+2\subtractright{2}\amp=6\subtractright{2}\amp\amp\text{remove the same quantity from each side} x\amp=4\amp\amp\text{still balanced; now it straight up tells you the solution}

Each line of the algebra above shows what is called an *equivalent equation* . Each of those equations is "algebraically equivalent" to the one that came before it, meaning it has exactly the same solution(s). The final equivalent equation $x=4$ tells us directly that the solution to the equation is $4$ .

In theory, there could have been more than one solution to the equation (although that is
      not the case with this equation). So conceptually, there is a *collection* of solutions
      to any given equation. We call this collection a *solution set* . Any set of numbers that only has one or a few numbers in it is written using curly braces. In this case, the solution set is $\{4\}$ .
      Using braces to surround a collection of numbers listed out is called *set notation* , not to be confused with set-builder notation from .

We have learned we can add or subtract the same number on each side of the equal sign,
      just like we can add or remove the same amount of weight on a balanced scale.
      Can we multiply and divide the same number on each side of the equal sign?
      Let's look at again: If a number times $2$ is $6$ , what is the number?
      Another balance scale can help visualize this.

At first, the scale is balanced.
      If we cut the weight in half on both sides, it should still be balanced.
      We can see from the scale that $x=3$ is correct.

Removing half of the weight from each side of the scale is like dividing both sides of an equation by $2$ : 2x\amp=6 \divideunder{2x}{2}\amp=\divideunder{6}{2} x\amp=3 The equivalent equation in this example is $x=3$ ,
      which tells us that the solution to the equation is $3$ and the solution set is $\{3\}$ .

Similarly, we can multiply each side of an equation by $2$ if that is helpful,
      and it will keep a scale in balance. We can summarize these properties.

**Properties of Equivalent Equations**
If there is an equation $\text{Left}=\text{Right}$ , we can do the following to obtain
          an equivalent equation. \[\text{Left}\addright{c}=\text{Right}\addright{c}\] (add the same number to each side) \[\text{Left}\subtractright{c}=\text{Right}\subtractright{c}\] (subtract the same number from each side) \[\text{(Left)}\multiplyright{c}=\text{(Right)}\multiplyright{c}\] (multiply each side of the equation by the same *nonzero* number) \[\displaystyle\divideunder{\text{Left}}{c}=\divideunder{\text{Right}}{c}\] (divide each side of the equation by the same *nonzero* number)

## Solving One-Step Equations and Stating Solution Sets

Notice when we solved equations in ,
      the final equation looked like $x=\text{number}$ ,
      where the variable $x$ stands alone on one side of the equal sign.
      The goal of solving any equation is to *isolate the variable* in this same manner.

Putting together both strategies
      (applying the opposite operation and balancing equations like a scale)
      that we just explored,
      we summarize how to solve a one-step linear equation.

Let's look at a few examples.

**Example**
Solve for $y$ in the equation $7+y=3$ .

*Solution*



Note that when solving the equation in we found $-3=a$ , and did not bother to write it the other way round as $a=-3$ .
      All that really matters is that we ended with a clear statement of the solution set,
      which was $\{-3\}$ .

**Example**
The formula for a circle's circumference is $c=\pi d$ , where $c$ represents
          circumference, $d$ represents diameter, and $\pi$ is a constant with the value
          of $3.1415926\ldots$ . If a circle's circumference is 12\pi ,
          find the circle's diameter. a circle with the diameter labeled d and the circumference labeled c \begin{tikzpicture}
              \node at (-3,0) {};
              \draw (-2,0) -- (2,0) node[pos=0.5, above] {diameter $d$};
              \def\myshift#1{\raisebox{1ex}}
              \draw [
                fill = firstcolor!20,
                postaction = {
                  decorate,
                    decoration = {
                      text along path,
                      text align = center,
                      text={|\myshift|The circumference $c$ is the distance all the way around.}
                    }
                  }
                ] (0,-2) arc [
                  start angle = 270,
                  end angle = -90,
                  x radius = 2,
                  y radius = 2
                ] ;
              \draw (-2,0) -- (2,0) node[pos = 0.5, above] {diameter $d$};
            \end{tikzpicture}

*Solution*


Examples so far have solved an equation by undoing addition, subtraction, multiplication,
      or division. There is one last arithmetic action that we will look into undoing: *negation* . Negation is when you apply the negative sign to a number. Undoing negation is simple though: just negate again. For example, $-(-42)=42$ .

**Example**
Solve the equation $-b=2$ for $b$ .

*Solution*
Our variable $b$ is not yet isolated because of the negative sign in front.
          To attack that negative sign, we can "negate" each side: -b\amp=2 \negate{-b}\amp=\negate{2} b\amp=-2 We removed the negative sign from $-b$ by negating both sides,
          and we can see that the solution set is $\{-2\}$ .
An alternative is to think of the original negative sign as multiplication by $-1$ .
          In that case, dividing each side by $-1$ will successfully isolate $b$ . -b\amp=2 -1\cdot b\amp=2 \divideunder{-1\cdot b}{-1}\amp=\divideunder{2}{-1} b\amp=-2
Another alternative is to recognize that multiplying on each side by $-1$ will
          also cancel that unwanted negative sign. -b\amp=2 \multiplyleft{-1}(-b)\amp=\multiplyleft{-1}(2) b\amp=-2
It is recommended that you review these three approaches for undoing negation,
          settle on the method that you like, and use it consistently.
We should check the solution
          by substituting $-2$ in for $b$ in the original equation: -b\amp=2 -(\substitute{-2})\amp\confirm{=}2 The solution $-2$ is checked, and the solution set is $\{-2\}$ .

## Equations with Fractions

When an equation has fractions, solving it uses the same principles.
      Of course you may need to use fraction arithmetic.
      Also, you might make good use of the reciprocal of a fraction
      as described in .

**Example**
Solve the equation $\frac{2}{3}+g=\frac{1}{2}$ for $g$ .

*Solution*


When the variable in an equation is multiplied by a fraction,
      you can use the reciprocal of that fraction to help solve the equation.
      The *reciprocal* of a fraction is the fraction you get from swapping the numerator and denominator.
      For example, the reciprocal of $\frac{2}{3}$ is $\frac{3}{2}$ .

A reciprocal is useful because when a fraction is multiplied by its reciprocal,
      the result is $1$ . For example, $\frac{2}{3}\cdot\frac{3}{2}=1$ .
      This helps us remove a fraction when it is multiplied by the variable.

**Example**
Solve the equation $\frac{5}{8}d=7$ for $d$ .

*Solution*


Sometimes the variable is in the numerator of a fraction, like in $\frac{3x}{4}$ .
      This is actually the same as $\frac{3}{4}x$ .
      Either way, $x$ is multiplied by $3$ and divided by $4$ .
      So this is another situation where the reciprocal of a fraction can help.

**Example**
Solve the equation $\frac{3x}{4}=10$ for $x$ .

*Solution*

