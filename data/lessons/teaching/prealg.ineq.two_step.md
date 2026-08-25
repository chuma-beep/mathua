> Content sourced from [ORCCA](https://pcc.edu/orcca) (Open Resources for Community College Algebra) — CC BY 4.0
# Solving Multistep Linear Inequalities

We solved inequalities in where only one step was needed to isolate the variable. Now we will work with inequalities that need more than one step.

## Solving Multistep Inequalities

When solving a linear inequality, we almost follow the exact same steps as we do in Process. One difference is that when we multiply or divide by a negative number on both sides of an inequality, the direction of the inequality symbol must switch. The other difference is that checking a solution set takes more effort.

**Example**
Solve for $t$ in the inequality $-3t+5\geq11$. Write the solution set in both set-builder notation and interval notation.

*Solution*
We'll handle this much like we would handle an equation, at least for the first step. 
\[ \begin{aligned} -3t+5&\geq11 \\ -3t+5\subtractright{5}&\geq11\subtractright{5} \\ -3t&\geq6 \\ \frac{-3t}{-3}&\mathbin{\secondhighlight{\le}}\frac{6}{-3} \\ t&\leq-2 \end{aligned} \]

Note that when we divided both sides of the inequality by $-3$, we had to switch the direction of the inequality symbol. At this point we think that the solution set in set-builder notation is $\{t\mid t\leq-2\}$, and the solution set in interval notation is $(-\infty,-2]$.

Since there are infinitely many solutions, it's impossible to literally check them all. We believe that all values of $t$ for which $t\leq-2$ are solutions. We check that one number less than $-2$ (any number, your choice) satisfies the inequality. *And* that $-2$ satisfies the inequality. *And* that one number greater than $-2$ (any number, your choice) does *not* satisfy the inequality. We choose to check the values $-10$, $-2$, and $0$.

\[ \begin{aligned} && -3t+5&\ge11& & \\ -3(-10)+5&\wonder{\geq}11& -3(-2)+5&\wonder{\geq}11& -3(0)+5&\wonder{\geq}11 \\ 30+5&\wonder{\geq}11& 6+5&\wonder{\geq}11& 0+5&\wonder{\geq}11 \\ 35&\confirm{\geq}11& 11&\confirm{\geq}11& 5&\reject{\geq}11 \end{aligned} \]
 So both $-10$ and $-2$ are solutions as expected, while $0$ is not. This is evidence that our solution set is correct. Making these checks would help us catch an error if we had made one. While it certainly does take time and space to make three checks like this, it has its value.

**Example**
Solve for $z$ in the inequality $(6z+5)-(2z-3)\gt-12$. Write the solution set in both set-builder notation and interval notation.

*Solution*
Here, our first step will be simplifying the left side. 
\[ \begin{aligned} (6z+5)-(2z-3)&\gt-12 \\ 6z+5-2z+3&\gt-12 \\ 4z+8&\gt-12 \\ 4z+8\subtractright{8}&\gt-12\subtractright{8} \\ 4z&\gt-20 \\ \frac{4z}{4}&\gt\frac{-20}{4} \\ z&\gt -5 \end{aligned} \]

Note that we divided both sides of the inequality by $4$ and since this is a positive number we did *not* need to switch the direction of the inequality symbol. At this point we think that the solution set in set-builder notation is $\{z\mid z\gt-5\}$, and the solution set in interval notation is $(-5,\infty)$.

Since there are infinitely many solutions, it's impossible to literally check them all. We believe that all values of $z$ for which $z\gt-5$ are solutions. We check that one number less than $-5$ (any number, your choice) does *not* satisfy the inequality. *And* that $-5$ does *not* satisfy the inequality. *And* that one number greater than $-5$ (any number, your choice) *does* satisfy the inequality. We choose to check the values $-10$, $-5$, and $0$.

\[ \begin{aligned} (6(-10)+5)-(2(-10)-3)&\wonder{\gt}-12 \\ (-60+5)-(-20-3)&\wonder{\gt}-12 \\ -55-(-23)&\wonder{\gt}-12 \\ -32&\reject{\gt}-12 \end{aligned} \]
 
\[ \begin{aligned} (6(-5)+5)-(2(-5)-3)&\wonder{\gt}-12 \\ (-30+5)-(-10-3)&\wonder{\gt}-12 \\ -25-(-13)&\wonder{\gt}-12 \\ -12&\reject{\gt}-12 \end{aligned} \]
 
\[ \begin{aligned} (6(0)+5)-(2(0)-3)&\wonder{\gt}-12 \\ (0+5)-(0-3)&\wonder{\gt}-12 \\ 5-(-3)&\wonder{\gt}-12 \\ 8&\confirm{\gt}-12 \end{aligned} \]
 So both $-10$ and $-5$ are not solutions as expected, while $0$ is a solution. This is evidence that our solution set is correct. The solution set in set-builder notation is $\{z\mid z\gt-5\}$. The solution set in interval notation is $(-5,\infty)$.

## Applications

**Rate Problem**
When an experiment started, the pressure inside a gas container was $4.2$ atm (one atm is the standard pressure air at sea level). As the container was heated, the pressure increased by $0.7$ atm per minute. The maximum pressure the container is rated to handle is $21.7$ atm. Heating must be stopped once the pressure reaches $21.7$ atm. Over what time interval was the container in a safe state (meaning the pressure was less than or equal to $21.7$ atm)?

*Solution*
This is a situation where something had an initial value (the pressure starts at $4.2$ atm) and then changed at a constant rate (it increased by $0.7$ atm per minute). So we can use the rate model formula. Except we are not exactly interested in the pressure *equaling* the final value of $21.7$ atm. Instead, we are asked about when the pressure was *less than or equal to* $21.7$ atm. So we have the inequality: 
\[ \begin{aligned} 0.7t+4.2&\leq21.7 \\ 0.7t+4.2\subtractright{4.2}&\leq21.7\subtractright{4.2} \\ 0.7t&\leq17.5 \\ \frac{0.7t}{0.7}&\leq\frac{17.5}{0.7} \\ t&\leq25 \end{aligned} \]

In summary, the container was safe as long as $t\leq25$. Assuming that the time $t$ also must be greater than or equal to zero, this means $0\leq t\leq 25$. We can write this as the time interval as $[0,25]$. Thus the container was safe between $0$ minutes and $25$ minutes.

**Percent Problem**
The population of a certain country grew by $7\%$ over the course of the past decade. One town in this country grew in population too, but not as fast as the country did overall. Its current population is $22{,}341$. What might its population have been ten years ago?

*Solution*
Let $x$ be the town's population from ten years ago. If the town had grown by $7\%$, then it's population would be $x + 0.07x$. But since it actually grew less quickly than $7\%$ per decade, $x + 0.07x$ would work out to more than the town's current population of $22{,}341$. So we have the inequality: 
\[ \begin{aligned} x + 0.07x&\gt22341 \\ 1.07x&\gt22341 \\ \frac{1.07x}{1.07}&\gt\frac{22341}{1.07} \\ x&\gt20879.4\ldots \end{aligned} \]

So the town's population from ten years ago was at least $20880$.

