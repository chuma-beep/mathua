> Content sourced from [ORCCA](https://pcc.edu/orcca) (Open Resources for Community College Algebra) — CC BY 4.0
# Slope-Intercept Form

In this section, we will explore the most common way to write the equation of a line. It's
      known as "slope-intercept form" .

## Slope-Intercept Definition

Recall , where Yara started with $\$50$ in her savings
      account, and from then on deposited $\$20$ each week. In that example, we used $x$ to represent how many weeks have passed. After $x$ weeks, Yara has added $20x$ dollars. Since she started with $\$50$ , she now has \[
        y=20x+50
      \] .
      In this example, there is a constant rate of change of $20$ dollars per week, so we call
      that the slope. We also saw in that plotting Yara's balance
      over time makes a straight-line graph.

The graph of Yara's savings has some things in common with almost every straight-line graph.
      There is a slope, and there is a place where the line crosses the $y$ -axis. illustrates this in the abstract.

We already have a symbol, $m$ , for the slope of a line. That other feature, where the
      line crosses the $y$ -intercept is of interest to us now. The *y-intercept* of a line is a *point* where the line crosses the $y$ -axis. Since it's on the $y$ -axis, the $x$ -coordinate of this point is $0$ . It is standard to call the point $(0,b)$ the $y$ -intercept, and call the
      number $b$ the "y-coordinate of the y-intercept" . It is almost
      inevitable that people will find this too wordy, and will call $b$ the $y$ -intercept. But technically, the $y$ -intercept is $(0,b)$ .

Use to answer this question.

One way to write the equation for Yara's savings was \[
        y=20x+50
      \] where $m=20$ and $b=50$ are immediately visible in the equation. Now we generalize
      this.

> **Definition**
> When $x$ and $y$ have a linear relationship where $m$ is the slope and $(0,b)$ is the $y$ -intercept, one equation for this relationship is y=mx+b and this equation is called the *slope-intercept form* of the line. It is
          called this because the slope and $y$ -intercept are immediately discernible from the
          numbers in the equation.


## Graphing Slope-Intercept Equations

**Example**

**Example**

**Example**


## Writing a Slope-Intercept Equation Given a Graph

We can write a linear equation in slope-intercept form based on its graph. We need to be able
      to calculate the line's slope and see its $y$ -intercept.



## Writing a Slope-Intercept Equation Given Two Points

Any two points uniquely determine a line. Once you identify two points, there is a process to
      find the slope-intercept form of the equation of the line that connects them.

**Example**
Find the slope-intercept form of an equation for the line that passes through the points $(0,5)$ and $(8,-5)$ .

*Solution*
Our goal is to write $y=mx+b$ , with specific numbers for $m$ and $b$ . The
          first step is to find the slope, $m$ . To do this, recall the slope formula from . It says that if a line passes through the points $(x_1,y_1)$ and $(x_2,y_2)$ , then the slope is found by the formula $m=\frac{y_2-y_1}{x_2-x_1}$ . Applying this to our two points $(\overset{x_1}{0},\overset{y_1}{5})$ and $(\overset{x_2}{8},\overset{y_2}{-5})$ , we see that the slope is: m\amp=\frac{y_2-y_1}{x_2-x_1} \amp=\frac{\substitute{-5}-\substitute{5}}{\substitute{8}-\substitute{0}} \amp=\frac{-10}{8}=-\frac{5}{4} We are trying to write $y=mx+b$ . Since we already found the slope, we know that we
          want to write $y=-\frac{5}{4}x+b$ but we need a specific number for $b$ . We *happen* to know that one point on this line is $(0,5)$ , which is on the $y$ -axis because its $x$ -value is $0$ . So $(0,5)$ is this line's $y$ -intercept, and therefore $b=5$ . So our equation is $y=-\frac{5}{4}x+5$ .

**Example**
Find the slope-intercept form of an equation for the line that passes through the points $(3,-8)$ and $(-6,1)$ .

*Solution*
We first find the slope between our two points: $(\overset{x_1}{3},\overset{y_1}{-8})$ and $(\overset{x_2}{-6},\overset{y_2}{1})$ . Using the slope formula again, we have: m\amp=\frac{y_2-y_1}{x_2-x_1} \amp=\frac{\substitute{1}-\substitute{(-8)}}{\substitute{{-}6}-\substitute{3}} \amp=\frac{9}{-9} \amp=-1
Now that we have the slope, we can write $y=-1x+b$ , or more simply: $y=-x+b$ .
          Unlike in , we are not given
          the value of $b$ because neither of our two given points have an $x$ -value of $0$ . To find $b$ , remember that we have two points that we already know should
          make the equation true! This means we can substitute *either* point into the
          equation (for the $x$ and the $y$ ) and solve for $b$ . Let's arbitrarily
          choose $(3,-8)$ to substitute in. y\amp=-x+b \substitute{-8}\amp=-(\substitute{3})+b\amp\text{(Now solve for }b\text{.)} -8\amp=-3+b -8\addright{3}\amp=-3+b\addright{3} -5\amp=b
We conclude that the slope-intercept line equation is $y=-x-5$ .



## Modeling with Slope-Intercept Form

We can model many relationships using slope-intercept form, and then solve related questions
      using algebra. Here are a few examples.

**Example**
Uber is a ride-sharing company. Its pricing in Portland factors in how much time and how
          many miles a trip takes. But if you assume that rides average out at a speed of 30 ,
          then their pricing scheme boils down to a base of $\$7.35$ for the trip, plus $\$3.85$ per mile. Use a slope-intercept equation and algebra to answer these
          questions. How much is the fare if a trip is $5.3$ miles long? With $\$100$ available to you, how long of a trip can you afford?

*Solution*
The rate of change (slope) is $\$3.85$ per mile, and the starting value is $\$7.35$ . So the slope-intercept equation is \[
            y=3.85x+7.35
          \] .
          In this equation, $x$ stands for the number of miles in a trip, and $y$ stands
          for the amount of money to be charged.
If a trip is $5.3$ miles long, we substitute $x=5.3$ into the equation and we
          have: y\amp=3.85x+7.35 \amp=3.85(\substitute{5.3})+7.35 \amp=20.405+7.35 \amp=27.755 And the $5.3$ -mile ride will cost you about $\$27.76$ . (We say "about" ,
          because this was all assuming you average 30 .)
Next, to find how long of a trip would cost $\$100$ , we substitute $y=100$ into
          the equation and solve for $x$ : y\amp=3.85x+7.35 \substitute{100}\amp=3.85x+7.35 100\subtractright{7.35}\amp=3.85x 92.65\amp=3.85x \divideunder{92.65}{3.85}\amp=x 24.06\amp\approx x So with $\$100$ you could afford a little more than a $24$ -mile trip.

