> Content sourced from [ORCCA](https://pcc.edu/orcca) (Open Resources for Community College Algebra) — CC BY 4.0

# Slope-Intercept Form

In this section, we will explore the most common way to write the equation of a line. It's known as "slope-intercept form".

## Slope-Intercept Definition

Recall  the example , where Yara started with $$\$50$$ in her savings account, and from then on deposited $$\$20$$ each week. In that example, we used $x$ to represent how many weeks have passed. After $x$ weeks, Yara has added $20x$ dollars. Since she started with $$\$50$$, she now has

$$
 y=20x+50 
$$

. In this example, there is a constant rate of change of $20$ dollars per week, so we call that the slope. We also saw in  the figure  that plotting Yara's balance over time makes a straight-line graph.

The graph of Yara's savings has some things in common with almost every straight-line graph. There is a slope, and there is a place where the line crosses the $y$-axis.  the figure  illustrates this in the abstract.

*Yara's savings*

*Generic line*

We already have a symbol, $m$, for the slope of a line. That other feature, where the line crosses the $y$-intercept is of interest to us now. The *$y$-intercept* of a line is a *point* where the line crosses the $y$-axis. Since it's on the $y$-axis, the $x$-coordinate of this point is $0$. It is standard to call the point $(0,b)$ the $y$-intercept, and call the number $b$ the "$y$-coordinate of the $y$-intercept". It is almost inevitable that people will find this too wordy, and will call $b$ the $y$-intercept. But technically, the $y$-intercept is $(0,b)$.

**Exercise**

Use  the figure  to answer this question.

One way to write the equation for Yara's savings was

$$
 y=20x+50 
$$

where $m=20$ and $b=50$ are immediately visible in the equation. Now we generalize this.

**Definition**

Slope-Intercept Form

When $x$ and $y$ have a linear relationship where $m$ is the slope and $(0,b)$ is the $y$-intercept, one equation for this relationship is

$$
 y=mx+b 
$$

and this equation is called the *slope-intercept form* of the line. It is called this because the slope and $y$-intercept are immediately discernible from the numbers in the equation.

**Exercise**

**Remark**

The number $b$ is the $y$-value when $x=0$. Therefore it is common to refer to $b$ as the *initial value* or *starting value* of a linear relationship.

## Graphing Slope-Intercept Equations

**Example**

With a simple equation like $y=2x+3$, we can see that this is a line whose slope is $2$ and which has initial value $3$. So starting at $y=3$ on the $y$-axis, each time we increase the $x$-value by $1$, the $y$-value increases by $2$. With these basic observations, we can quickly produce a table and/or a graph.

**Example**

The conversion formula for a Celsius temperature into Fahrenheit is $F=\frac{9}{5}C+32$. This appears to be in slope-intercept form, except that $x$ and $y$ are replaced with $C$ and $F$. Suppose you are asked to graph this equation. How will you proceed? You *could* make a table of values as we did in  the example  but that takes time and effort. Since the equation is in slope-intercept form, there is a *better* way.

Since this equation is for converting a Celsius temperature to a Fahrenheit temperature, it makes sense to let $C$ be the horizontal axis variable and $F$ be the vertical axis variable. Note the slope is $\frac{9}{5}$ and the vertical intercept (here, the $F$-intercept) is $(0,32)$.

1. Set up the axes using an appropriate window and labels. Considering the freezing temperature of water ($0^{\circ}$ Celsius or $32^{\circ}$ Fahrenheit), and the boiling temperature of water ($100^{\circ}$ Celsius or $212^{\circ}$ Fahrenheit), it's reasonable to let $C$ run through at least $0$ to $100$ and $F$ run through at least $32$ to $212$.
2. Plot the $F$-intercept, which is at $(0,32)$.
3. Starting at the $F$-intercept, use slope triangles to reach the next point. Since our slope is $\frac{9}{5}$, that suggests a "run" of $5$ and a "rise" of $9$ might work. But as  the figure  indicates, such slope triangles are too tiny. You can actually use any fraction equivalent to $\frac{9}{5}$ to plot using the slope, as in $\frac{18}{10}$, $\frac{90}{50}$, $\frac{900}{50}$, or $\frac{45}{25}$ which all reduce to $\frac{9}{5}$. Given the size of our graph, we will use $\frac{90}{50}$ to plot points, where we will try a "run" of $50$ and a "rise" of $90$.
4. Connect your points with a straight line, use arrowheads, and label the equation.

*Graphing $F=\frac{9}{5}C+32$*

**Example**

Graph $y=-\frac{2}{3}x+10$.

*Graphing $y=-\frac{2}{3}x+10$*

**Exercise**

## Writing a Slope-Intercept Equation Given a Graph

We can write a linear equation in slope-intercept form based on its graph. We need to be able to calculate the line's slope and see its $y$-intercept.

**Exercise**

**Exercise**

## Writing a Slope-Intercept Equation Given Two Points

Any two points uniquely determine a line. Once you identify two points, there is a process to find the slope-intercept form of the equation of the line that connects them.

**Example**

Find the slope-intercept form of an equation for the line that passes through the points $(0,5)$ and $(8,-5)$.

Our goal is to write $y=mx+b$, with specific numbers for $m$ and $b$. The first step is to find the slope, $m$. To do this, recall the  slope formula  from  the section . It says that if a line passes through the points $(x_1,y_1)$ and $(x_2,y_2)$, then the slope is found by the formula $m=\frac{y_2-y_1}{x_2-x_1}$. Applying this to our two points $(\overset{x_1}{0},\overset{y_1}{5})$ and $(\overset{x_2}{8},\overset{y_2}{-5})$, we see that the slope is:

$$
\begin{aligned}m&=\frac{y_2-y_1}{x_2-x_1} \\ &=\frac{\substitute{-5}-\substitute{5}}{\substitute{8}-\substitute{0}} \\ &=\frac{-10}{8}=-\frac{5}{4}\end{aligned}
$$

We are trying to write $y=mx+b$. Since we already found the slope, we know that we want to write $y=-\frac{5}{4}x+b$ but we need a specific number for $b$. We *happen* to know that one point on this line is $(0,5)$, which is on the $y$-axis because its $x$-value is $0$. So $(0,5)$ is this line's $y$-intercept, and therefore $b=5$. So our equation is $y=-\frac{5}{4}x+5$.

**Example**

Find the slope-intercept form of an equation for the line that passes through the points $(3,-8)$ and $(-6,1)$.

We first find the slope between our two points: $(\overset{x_1}{3},\overset{y_1}{-8})$ and $(\overset{x_2}{-6},\overset{y_2}{1})$. Using the  slope formula  again, we have:

$$
\begin{aligned}m&=\frac{y_2-y_1}{x_2-x_1} \\ &=\frac{\substitute{1}-\substitute{(-8)}}{\substitute{{-}6}-\substitute{3}} \\ &=\frac{9}{-9} \\ &=-1\end{aligned}
$$

Now that we have the slope, we can write $y=-1x+b$, or more simply: $y=-x+b$. Unlike in  the example , we are not given the value of $b$ because neither of our two given points have an $x$-value of $0$. To find $b$, remember that we have two points that we already know should make the equation true! This means we can substitute *either* point into the equation (for the $x$ and the $y$) and solve for $b$. Let's arbitrarily choose $(3,-8)$ to substitute in.

$$
\begin{aligned}y&=-x+b \\ \substitute{-8}&=-(\substitute{3})+b&\text{(Now solve for }b\text{.)} \\ -8&=-3+b \\ -8\addright{3}&=-3+b\addright{3} \\ -5&=b\end{aligned}
$$

We conclude that the slope-intercept line equation is $y=-x-5$.

**Exercise**

**Exercise**

## Modeling with Slope-Intercept Form

We can model many relationships using slope-intercept form, and then solve related questions using algebra. Here are a few examples.

**Example**

Uber is a ride-sharing company. Its pricing in Portland factors in how much time and how many miles a trip takes. But if you assume that rides average out at a speed of $30\,\text{mileperhour}$, then their pricing scheme boils down to a base of $$\$7.35$$ for the trip, plus $$\$3.85$$ per mile. Use a slope-intercept equation and algebra to answer these questions.

1. How much is the fare if a trip is $5.3$ miles long?
2. With $$\$100$$ available to you, how long of a trip can you afford?

The rate of change (slope) is $$\$3.85$$ per mile, and the starting value is $$\$7.35$$. So the slope-intercept equation is

$$
 y=3.85x+7.35 
$$

. In this equation, $x$ stands for the number of miles in a trip, and $y$ stands for the amount of money to be charged.

If a trip is $5.3$ miles long, we substitute $x=5.3$ into the equation and we have:

$$
\begin{aligned}y&=3.85x+7.35 \\ &=3.85(\substitute{5.3})+7.35 \\ &=20.405+7.35 \\ &=27.755\end{aligned}
$$

And the $5.3$-mile ride will cost you about $$\$27.76$$. (We say "about", because this was all assuming you average $30\,\text{mileperhour}$.)

Next, to find how long of a trip would cost $$\$100$$, we substitute $y=100$ into the equation and solve for $x$:

$$
\begin{aligned}y&=3.85x+7.35 \\ \substitute{100}&=3.85x+7.35 \\ 100\subtractright{7.35}&=3.85x \\ 92.65&=3.85x \\ \divideunder{92.65}{3.85}&=x \\ 24.06&\approx x\end{aligned}
$$

So with $$\$100$$ you could afford a little more than a $24$-mile trip.

**Exercise**

##

**Exercise**

How does "slope-intercept form" get its name?

**Exercise**

What are two phrases you can use for "$b$" in a slope-intercept form line equation?

**Exercise**

Explain the two basic steps to graphing a line when you have the equation in slope-intercept form. (Not counting the step where you draw and label the axes and ticks.)

##

