> Content sourced from [ORCCA](https://pcc.edu/orcca) (Open Resources for Community College Algebra) — CC BY 4.0
# Solving One-Step Inequalities

In this section, we learn that solving small inequalities is not all that different from solving small equations.

## Solving Linear Inequalities

To "solve an inequality" in algebra context means to identify all the solutions for that inequality.For the most part, the properties from apply to inequalities too, not just equations. Here are some numerical examples to consider.

**Example**
Solve the inequality $t+7\lt5$.

*Solution*
There is not much difference between the steps to solve this inequality and the steps to solve the *equation* $t+7=5$. We can subtract $7$ from each side. 
\[ \begin{aligned} t+7&\lt5 \\ t+7\subtractright{7}&\lt5\subtractright{7} \\ t&\lt-2 \end{aligned} \]

When we solve a linear inequality, there are usually infinitely many solutions. The solution set has infinitely many numbers in it. This is unlike when we solve a linear equation, where there is usually only one solution and one number in the solution set. For this example, any number less than $-2$ is a solution.

There are at least three ways to represent an inequality's solution set: graphically, with set-builder notation, and with interval notation. (Interval notation and set-builder notation are discussed in.) Graphically, the solution set is part of a number line:

Using interval notation, we write the solution set by reading the number line from left to right. The solution set is $(-\infty,-2)$.

Using set-builder notation, we write the solution with generic set braces, declaring $t$ to be the variable, and writing the condition that $t$ needs to meet: $\{t\mid t\lt-2\}$.

As with equations, we should check solutions to catch human mistakes. Since there are infinitely many solutions, it's impossible to literally check them all. So we settle for something that gives us confidence our solution set is correct, but does not take forever to do.

According to our solution, all values of $t$ for which $t\lt-2$ are solutions and all values of $t$ for which $t\geq2$ are not solutions. So an approach we can use is to check if one number less than $-2$ (any number, your choice) satisfies the inequality. And *also* that $-2$ itself does *not* satisfy the inequality. And *also* that one number greater than $-2$ (any number, your choice) does *not* satisfy the inequality.

Here we will test $-3$, $-2$, and $0$ in the original inequality. 
\[ \begin{aligned} && t+7&\lt5&& \\ -3+7&\wonder{\lt}5& -2+7&\wonder{\lt}5& 0+7&\wonder{\lt}5 \\ 4&\confirm{\lt}5& 5&\reject{\lt}5& 7&\reject{\lt}5 \end{aligned} \]
 It worked! The number $-3$ is a solution, and both $-2$ and $0$ are *not*. This is what we expected. This is evidence that our solution set is correct, and we can feel more secure that we did not make a human mistake when we were solving. While it takes time and space to make three checks, it's worth it.

## Negation

Something interesting happens when we multiply or divide by a *negative* number on each side of an inequality: the direction reverses! To understand why, consider, where the numbers $2$ and $4$ are each multiplied by $-1$.

Starting with $2\lt4$, if we multipled each sides by $-1$ and left the inequality sign alone, we would get the *false* inequality $-2\reject{\lt}-4$. We should change the direction so we have the *true* inequality $-2\gt-4$.

**Changing the Direction of the Inequality Sign**
When multiplying or dividing each side of an inequality by a *negative* number, the inequality sign must change direction.

Do not change the inequality direction when multiplying/dividing by a *positive* number, or when *adding/subtracting* by any number.

**Example**
Solve the inequality $-2x\geq12$. State the solution set graphically, using interval notation, and using set-builder notation.

*Solution*
To solve this inequality, we will divide each side by $-2$: 
\[ \begin{aligned} -2x&\geq12 \\ \frac{-2x}{-2}&\secondhighlight{\leq}\frac{12}{-2}&&\text{Note the change in direction.} \\ x&\leq-6 \end{aligned} \]

The inequality sign changed direction in the same step where we divided by a negative number. Graphically, the solution set is part of a number line:

Using interval notation, we write the solution set as $(-\infty,-6]$. Using set-builder notation, we write the solution set as $\{x\mid x\leq-6\}$.

We should check that some number less than $-6$ is a solution, that $-6$ itself is also a solution, and that some number greater than $-6$ is not a solution. 
\[ \begin{aligned} && -2x&\ge12&& \\ -2(-7)&\wonder{\geq}12& -2(-6)&\wonder{\ge}12& -2(-5)&\wonder{\ge}12 \\ 14&\confirm{\geq}12& 12&\confirm{\ge}12& 10&\reject{\ge}12 \end{aligned} \]
 Everything came out as expected, so our solution is reasonably checked.

