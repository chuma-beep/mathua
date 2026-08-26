> Content sourced from [ORCCA](https://pcc.edu/orcca) (Open Resources for Community College Algebra) — CC BY 4.0

# Standard Form

We've seen that a linear relationship can be expressed with an equation in  slope-intercept form  or with an equation in  point-slope form . There is a third form that you can use to write line equations. It's known as "standard form".

## Standard Form Definition

Imagine gathering donations to pay for a $$\$10{,}000$$ medical procedure you cannot afford. Oversimplifying the mathematics a bit, suppose that there were only two types of donors in the world: those who will donate $$\$20$$ and those who will donate $$\$100$$. How many of each, or what combination, do you need to reach the funding goal? That is, if $x$ people donate $$\$20$$ and $y$ people donate $$\$100$$, what numbers could $x$ and $y$ be to meet your goal? The donors of the first type have collectively donated $20x$ dollars, and the donors of the second type have collectively donated $100y$.

To reach $$\$10{,}000$$, altogether you'd need to have

$$
 20x+100y=10000 
$$

This is an example of a line equation in *standard form*.

**Definition**

Standard Form

It is always possible to write an equation for a line in the form

$$
 Ax+By=C 
$$

where $A$, $B$, and $C$ are three numbers (each of which might be $0$, although at least one of $A$ and $B$ must be nonzero). This form of a line equation is called *standard form*. In the context of an application, the meaning of $A$, $B$, and $C$ depends on that context. This equation is called "standard" form perhaps because *any* line can be written this way, even vertical lines, which do not have slope and therefore cannot be written using slope-intercept or point-slope form.

**Exercise**

Returning to the example with donations for the medical procedure, let's examine the equation

$$
 20x+100y=10000 
$$

. What units are attached to each part of this equation? The $10000$ is in dollars. Both $x$ and $y$ are numbers of people. Both the $20$ and the $100$ are in dollars per person. So the terms $20x$ and $100y$ are ultimately in dollars. Note how both sides of the equation are in dollars.

What is the slope of the linear relationship? It's not immediately visible since $m$ is not part of the standard form equation. But we can use algebra to isolate $y$:

$$
\begin{aligned}20x+100y&=10000 \\ 100y&=\highlight{-20x}+10000 \\ y&=\divideunder{-20x+10000}{100} \\ y&=\frac{-20x}{100}+\frac{10000}{100} \\ y&=-\frac{1}{5}x+100\end{aligned}
$$

. And we see that the slope is $-\frac{1}{5}$.

What units are on that slope? As always, the units on slope are $\frac{y\text{-unit}}{x\text{-unit}}$. In this case that's $\frac{\text{person}}{\text{person}}$, which sounds a little weird, but this slope of $-\frac{1}{5}\frac{\text{person}}{\text{person}}$ is saying that for every 5 extra *people* who donate $$\$20$$, you need $1$ fewer *person* donating $$\$100$$ to still reach your goal.

What is the $y$-intercept? Since we've already converted the equation into slope-intercept form, we can see that it is at $(0,100)$. This tells us that if $0$ people donate $$\$20$$, then you will need $100$ people donating $$\$100$$ to meet the goal.

What does a graph for this line look like? We've already converted into slope-intercept form, and we could use that to make the graph. But when given a line in standard form, there is another approach that may be preferable. Returning to

$$
 20x+100y=10000 
$$

, let's calculate the $y$-intercept and the $x$-intercept from scratch. Recall that these are *points* where the line crosses the $y$-axis and $x$-axis. To be on the $y$-axis means that $x=0$, and to be on the $x$-axis means that $y=0$. All these "$0$"s make the resulting algebra easy to finish:

$$
\begin{aligned}20x+100y&=10000&20x+100y&=10000 \\ 20(\substitute{0})+100y&=10000&20x+100(\substitute{0})&=10000 \\ 100y&=10000&20x&=10000 \\ y&=\divideunder{10000}{100}& x&=\divideunder{10000}{20} \\ y&=100& x&=500\end{aligned}
$$

So we have a $y$-intercept at $(0,100)$ and an $x$-intercept at $(500,0)$. If we plot these, we get to mark especially relevant points given the context, and then drawing a straight line between them gives us  the figure .

## The Importance of $x$- and $y$-Intercepts

With a linear relationship (and other types of equations too), we are often interested in the $x$-intercept and $y$-intercept because they have special meaning in context. For example, in  the figure , the $x$-intercept implies that if *no one* donates $$\$100$$, you need $500$ people to donate $$\$20$$ to get us to $$\$10{,}000$$. And the $y$-intercept implies if *no one* donates $$\$20$$, you need $100$ people to donate $$\$100$$. Let's look at another example.

**Example**

James owns a restaurant that uses about $32$ lb of flour every day. He just purchased $1200$ lb of flour. Model the amount of flour that remains $x$ days later with a linear equation, and interpret the meaning of its $x$-intercept and $y$-intercept.

Since the rate of change is constant ($-32$ lb every day), and we know the initial value, we can model the amount of flour at the restaurant with a  slope-intercept form  equation:

$$
 y=-32x+1200 
$$

where $x$ represents the number of days passed since the initial purchase, and $y$ represents the amount of flour left (in lb).

A line's $x$-intercept is on the $x$-axis, so its $y$-value must be $0$.  To find this line's $x$-intercept, we substitute $y$ with $0$, and solve for $x$:

$$
\begin{aligned}y&=-32x+1200 \\ \substitute{0}&=-32x+1200 \\ \highlight{-1200}&=-32x \\ \divideunder{-1200}{-32}&=x \\ 37.5&=x\end{aligned}
$$

So the line's $x$-intercept is at $(37.5,0)$. In context this means the flour would last for $37.5$ days.

A line's $y$-intercept is on the $y$-axis, so its $x$-value must be $0$. This line equation is already in slope-intercept form, so we simply recognize that its $y$-intercept is at $(0,1200)$. In general though, we would substitute $x$ with $0$, and we have:

$$
\begin{aligned}y&=-32x+1200 \\ y&=-32(\substitute{0})+1200 \\ y&=1200\end{aligned}
$$

So yes, the line's $y$-intercept is at $(0,1200)$. This means that when the flour was purchased, there was $1200$ lb of it. In other words, the $y$-intercept tells us one of the original pieces of information: in the beginning, James purchased $1200$ of flour.

The important thing is that both intercepts have relevant meaning in the context of the example. One was the initial amount of flour, and the other was how long until we would run out of flour.

If a line is in standard form, it may be easiest to graph it using its two intercepts.

**Example**

Graph $2x-3y=-6$ using its intercepts. And then use the intercepts to calculate the line's slope.

To graph a line by its $x$-intercept and $y$-intercept, it might help to first set up a table like in  the figure :

*Intercepts of $2x-3y=-6$*

A table like this might help you stay focused on searching for *two* points. As To find an $x$-intercept, $y$ must be $0$. This is why we put $0$ in the $y$-value cell of the $x$-intercept row. Similarly, a line's $y$-intercept has $x=0$, and we put $0$ into the $x$-value cell of the $y$-intercept row.

Next, we calculate the line's $x$-intercept by substituting $y=0$ into the equation

$$
\begin{aligned}2x-3y&=-6 \\ 2x-3(\substitute{0})&=-6 \\ 2x&=-6 \\ x&=-3\end{aligned}
$$

So the line's $x$-intercept is $(-3,0)$.

Similarly, we substitute $x=0$ into the equation to calculate the $y$-intercept:

$$
\begin{aligned}2x-3y&=-6 \\ 2(\substitute{0})-3y&=-6 \\ -3y&=-6 \\ y&=2\end{aligned}
$$

So the line's $y$-intercept is $(0,2)$.

*Intercepts of $2x-3y=-6$*

With both intercepts' coordinates, we can graph the line:

*Graph of $2x-3y=-6$*

There is a slope triangle from the $x$-intercept to the origin up to the $y$-intercept. It tells us that the slope is

$$
 m=\frac{\Delta y}{\Delta x}=\frac{2}{3} 
$$

.

This last example generalizes to a fact worth noting.

**Fact**

If a line's $x$-intercept is at $(r,0)$ and its $y$-intercept is at $(0,b)$, then the slope of the line is $-\frac{b}{r}$. (Unless the line passes through the origin, in which case both $r$ and $b$ equal $0$, and then this fraction is undefined. And the slope of the line could be anything.)

**Exercise**

What is its $x$-intercept?

Solution

What is its $y$-intercept?

Solution

What is its slope?

Solution

**Exercise**

## Transforming between Standard Form and Slope-Intercept Form

Sometimes a linear equation arises in  standard form , but it would be useful to see that equation in  slope-intercept form . Or perhaps, vice versa.

**Example**

Change $2x-3y=-6$ to slope-intercept form, and then graph it. (This is the same equation from  the example ).

Since a line in slope-intercept form looks like $y=\ldots$, we will isolate $y$:

$$
\begin{aligned}2x-3y&=-6 \\ -3y&=-6\subtractright{2x} \\ -3y&=-2x-6 \\ y&=\divideunder{-2x-6}{-3} \\ y&=\frac{-2x}{-3}-\frac{6}{-3} \\ y&=\frac{2}{3}x+2\end{aligned}
$$

In the third line, we wrote $-2x-6$ on the right side, instead of $-6-2x$. The only reason we did this is because we are headed to slope-intercept form, where the $x$-term is traditionally written first.

Now we can see that the slope is $\frac{2}{3}$ and the $y$-intercept is at $(0,2)$. With these things found, we can graph the line using slope triangles.

Compare this graphing method with the Graphing by Intercepts method in  the example . We have more points in this graph, thus we can graph the line more accurately.

*Graphing $2x-3y=-6 $ with Slope Triangles*

**Example**

Graph $2x-3y=0$.

First, we will try (and fail) to graph this line using its $x$- and $y$-intercepts.

Trying to find the $x$-intercept:

$$
\begin{aligned}2x-3y&=0 \\ 2x-3(\substitute{0})&=0 \\ 2x&=0 \\ x&=0\end{aligned}
$$

So the line's $x$-intercept is at $(0,0)$, at the origin. Hmm, the origin is *also* on the $y$-axis So we've also found the $y$-intercept even though that was not the immediate goal.

Since both intercepts are the same point, there is no way to use the intercepts alone to graph this line. So what can be done?

One is to convert the line equation into slope-intercept form:

$$
\begin{aligned}2x-3y&=0 \\ -3y&=0\subtractright{2x} \\ -3y&=-2x \\ y&=\divideunder{-2x}{-3} \\ y&=\frac{2}{3}x\end{aligned}
$$

So the line's slope is $\frac{2}{3}$, and we can graph the line using slope triangles and the intercept at $(0,0)$, as in  the figure .

*Graphing $2x-3y=0 $ with Slope Triangles*

If $C=0$ in a  standard form equation , it's convenient to graph it by first converting the equation to  slope-intercept form .

**Example**

Write the equation $y=\frac{2}{3}x+2$ in standard form.

Once we subtract $\frac{2}{3}x$ on both sides of the equation, we have

$$
 -\frac{2}{3}x+y=2 
$$

Technically, this equation is already in standard form $Ax+By=C$. However, you might like to end up with an equation that has no fractions, so you could multiply each side by $3$:

$$
\begin{aligned}\multiplyleft{3}\left(-\frac{2}{3}x+y\right)&=\multiplyleft{3}2 \\ -2x+3y&=6\end{aligned}
$$

##

**Exercise**

What kind of line has an equation in standard form, but cannot be written in slope-intercept form or point-slope form?

**Exercise**

What are some reasons why you might want to find the $x$- and $y$-intercepts of a line?

**Exercise**

What is not immediately apparent from standard form, that *is* immediately apparent from slope-intercept form and point-slope form?

##

