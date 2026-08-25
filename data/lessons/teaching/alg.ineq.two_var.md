> Content sourced from [ORCCA](https://pcc.edu/orcca) (Open Resources for Community College Algebra) — CC BY 4.0
# Solving Inequalities Graphically

In this text, we have mostly focused on solving inequalities algebraically. While we have had some practice solving inequalities graphically with technology, we want to solidify those skills. Solving using graphing is special because the graphing utility we use can do much of the heavy lifting and all that is left is to analyze the graph that is shown to us. So let's let our favorite graphing program make some graphs for us and then we can interpret the results.

## Solving Inequalities Graphically

**Example**

Let's turn to an example involving a linear equation.

**Example**
Solve the inequality $3x-2\lt7$ graphically.

*Solution*

## Solving Absolute Value and Quadratic Inequalities Graphically

Recall in Section that we learned that graphs of absolute value function are in general shaped like "V" s. We can now solve some absolute value inequalities graphically.

**Example**

**Example**
Solve the inequality $\abs{\frac{2}{3}x+1} \lt 3$ graphically.

*Solution*
To solve the inequality $\abs{\frac{2}{3}x+1} \lt 3$, we will start by making a graph with both $y=\abs{\frac{2}{3}x+1}$ and $y=3$.

The last examples had absolute value expressions being *less than* some value. We now need to investigate what happens when we have an absolute value expression that is *greater than* a value.

**Example**

**Example**
Solve the inequality $\abs{\frac{1}{3}x+2} \ge 6$ graphically.

*Solution*
To solve the inequality $\abs{\frac{1}{3}x+2} \ge 6$, we will start by making a graph with both $y=\abs{\frac{1}{3}x+2}$ and $y=6$.

Solving inequalities with quadratic expressions graphically is very similar to solving absolute value inequalities graphically.

**Example**
Graphically solve the following quadratic inequalities. 
1. $42(x-2)^2-60 \geq 21x-39$
2. $42(x-2)^2-60 \lt 21x-39$

*Solution*
1. To solve $42(x-2)^2-60 \ge 21x-39$, we need to determine where the $y$-values of the parabola are higher than (or equal to) those of the line. This region is highlighted in Figure. We can see that $42(x-2)^2-60 \ge 21x-39$ for all values of $x$ where $x\le 1$ or $x\ge 3.5$. We can write this solution set in interval notation as $(-\infty,1]\cup[3.5,\infty)$ or in set-builder notation as $\{x\mid x\le 1 \text{ or } x\ge 3.5\}$.
2. To now solve $42(x-2)^2-60 \lt 21x-39$, we will need to determine where the $y$-values of the parabola are *less* than those of the line. This region is highlighted in Figure. So the solutions to this inequality include all values of $x$ for which $1\lt x \lt 3.5$. We can write this solution set in interval notation as $(1,3.5)$ or in set-builder notation as $\{x \mid 1 \lt x \lt 3.5\}$.

## Solving Compound Inequalities Graphically

**Example**
Figure shows a graph of $y=f(x)$. Use the graph to solve the inequality $2\le f(x) \lt 6$.

*Solution*

**Example**
Figure shows a graph of $y=g(x)$. Use the graph to solve the inequality $-4\lt g(x) \le 3$.

*Solution*
To solve $-4\lt g(x) \le 3$, we first draw the horizontal lines $y=-4$ and $y=3$. To solve this inequality we notice that there are two pieces of the function $g$ that are trapped between the $y$-values $-4$ and $3$.

The solution set is the compound inequality $(-2.1,0.7)\cup(2.4,3.2]$.

**Example**
Phuong is taking the standard climbing route on Mount Hood from Timberline Lodge up the Southside Hogsback to the summit and back down the same way. Her altitude can be very closely modeled by an absolute value function since the angle of ascent is nearly constant. Let $x$ represent the number of miles walked from Timberline Lodge, and let $f(x)$ represent the altitude, in miles, after walking for a distance $x$. The altitude can be modeled by $f\left(x\right)=2.1-0.3077\cdot\abs{x-3.25}$. Note that below Timberline Lodge this model fails to be accurate. 
1. Solve the equation $f(x)=1.1$ graphically and interpret the results in the context of the problem.
2. Altitude sickness can occur at or above altitudes $1.5$ miles. Set up and solve an inequality graphically to find out how far Phuong can walk the trail and still be under $1.5$ miles of elevation.

*Solution*
1. First, we substitute the formula for $f(x)$ and simplify the equation. $\begin{aligned} f(x)&=1.1 \\ 2.1-0.3077\cdot\abs{x-3.25} &=1.1 \end{aligned}$ At this point, we should make a graph of both $y=2.1-0.3077\cdot\abs{x-3.25}$ and $y=1.1$ and find their intersections. Next, we should note that we are looking for the $x$-values of the intersections. These solutions are $0$ and $6.5$. According to the model, Phuong will be at $1.1$ miles of elevation after walking about $0$ miles as well as about $6.5$ miles along the trail. This implies that Timberline Lodge is at $1.1$ miles of elevation. In addition, it implies that the entire hike is $6.5$ miles round trip, ending at Timberline Lodge again.
2. The inequality we are looking for will describe when the altitude is below $1.5$ miles, but also above $1.1$ miles based on the reality of the situation (since the model only works above Timberline lodge at $1.1$ miles of altitude). Since $f(x)$ is the altitude, the inequality we need is $1.1\le f(x)\lt 1.5$, which becomes $1.1\le2.1-0.3077\cdot\abs{x-3.25}\lt 1.5$. Let's examine the graph again to solve this inequality: We are looking for places on the graph where the $y$-value is above $1.1$, but also where the graph is below $1.5$. To find this, we will draw in lines at both of those $y$-values and find intersections with $f$. The highlighted portions of the graph have $x$-values that satisfy the inequalities $0\le x\lt1.3 \text{ or } 5.2\lt x\le 6.5$. In conclusion, based both on our math and the reality of the situation, regions of the trail that are below $1.5$ miles are those that are from Timberline Lodge (at $0$ miles on the trail), to $1.3$ miles along the trail and then also from $5.2$ miles along the trail (and by now we are on our way back down) to $6.5$ miles along the trail (back at Timberline Lodge). If we wanted to write this in interval notation, we might write $[0,1.3)\cup(5.2,6.5]$. There is a big portion along the trail (from $1.3$ miles to $5.2$ miles) that Phuong will be above the $1.5$ mile altitude and should watch for signs of altitude sickness.

