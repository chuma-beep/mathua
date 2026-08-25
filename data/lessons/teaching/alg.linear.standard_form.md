> Content sourced from [ORCCA](https://pcc.edu/orcca) (Open Resources for Community College Algebra) — CC BY 4.0
# Standard Form

We've seen that a linear relationship can be expressed with an equation in slope-intercept form or with an equation in point-slope form. There is a third form that you can use to write line equations. It's known as "standard form".

## Standard Form Definition

Imagine gathering donations to pay for a $\$10{,}000$ medical procedure you cannot afford. Oversimplifying the mathematics a bit, suppose that there were only two types of donors in the world: those who will donate $\$20$ and those who will donate $\$100$. How many of each, or what combination, do you need to reach the funding goal? That is, if $x$ people donate $\$20$ and $y$ people donate $\$100$, what numbers could $x$ and $y$ be to meet your goal? The donors of the first type have collectively donated $20x$ dollars, and the donors of the second type have collectively donated $100y$.

To reach $\$10{,}000$, altogether you'd need to have 
\[20x+100y=10000\]
 This is an example of a line equation in *standard form*.

> **Definition**
> It is always possible to write an equation for a line in the form Ax+By=C where $A$, $B$, and $C$ are three numbers (each of which might be $0$, although at least one of $A$ and $B$ must be nonzero). This form of a line equation is called *standard form*. In the context of an application, the meaning of $A$, $B$, and $C$ depends on that context. This equation is called "standard" form perhaps because *any* line can be written this way, even vertical lines, which do not have slope and therefore cannot be written using slope-intercept or point-slope form.

Returning to the example with donations for the medical procedure, let's examine the equation 
\[20x+100y=10000\]. What units are attached to each part of this equation? The $10000$ is in dollars. Both $x$ and $y$ are numbers of people. Both the $20$ and the $100$ are in dollars per person. So the terms $20x$ and $100y$ are ultimately in dollars. Note how both sides of the equation are in dollars.

What is the slope of the linear relationship? It's not immediately visible since $m$ is not part of the standard form equation. But we can use algebra to isolate $y$: 
\[ \begin{aligned} 20x+100y&=10000 \\ 100y&=\highlight{-20x}+10000 \\ y&=\frac{-20x+10000}{100} \\ y&=\frac{-20x}{100}+\frac{10000}{100} \\ y&=-\frac{1}{5}x+100 \end{aligned} \]. And we see that the slope is $-\frac{1}{5}$.

What units are on that slope? As always, the units on slope are $\frac{y\text{-unit}}{x\text{-unit}}$. In this case that's $\frac{\text{person}}{\text{person}}$, which sounds a little weird, but this slope of $-\frac{1}{5}\frac{\text{person}}{\text{person}}$ is saying that for every 5 extra *people* who donate $\$20$, you need $1$ fewer *person* donating $\$100$ to still reach your goal.

What is the $y$-intercept? Since we've already converted the equation into slope-intercept form, we can see that it is at $(0,100)$. This tells us that if $0$ people donate $\$20$, then you will need $100$ people donating $\$100$ to meet the goal.

What does a graph for this line look like? We've already converted into slope-intercept form, and we could use that to make the graph. But when given a line in standard form, there is another approach that may be preferable. Returning to 
\[20x+100y=10000\], let's calculate the $y$-intercept and the $x$-intercept from scratch. Recall that these are *points* where the line crosses the $y$-axis and $x$-axis. To be on the $y$-axis means that $x=0$, and to be on the $x$-axis means that $y=0$. All these "0" s make the resulting algebra easy to finish: 
\[ \begin{aligned} 20x+100y&=10000&20x+100y&=10000 \\ 20(0)+100y&=10000&20x+100(0)&=10000 \\ 100y&=10000&20x&=10000 \\ y&=\frac{10000}{100}& x&=\frac{10000}{20} \\ y&=100& x&=500 \end{aligned} \]

So we have a $y$-intercept at $(0,100)$ and an $x$-intercept at $(500,0)$. If we plot these, we get to mark especially relevant points given the context, and then drawing a straight line between them gives us.

## The Importance of x- and y-Intercepts

With a linear relationship (and other types of equations too), we are often interested in the $x$-intercept and $y$-intercept because they have special meaning in context. For example, in, the $x$-intercept implies that if *no one* donates $\$100$, you need $500$ people to donate $\$20$ to get us to $\$10{,}000$. And the $y$-intercept implies if *no one* donates $\$20$, you need $100$ people to donate $\$100$. Let's look at another example.

**Example**

If a line is in standard form, it may be easiest to graph it using its two intercepts.

**Example**
Graph $2x-3y=-6$ using its intercepts. And then use the intercepts to calculate the line's slope.

*Solution*
To graph a line by its $x$-intercept and $y$-intercept, it might help to first set up a table like in:

A table like this might help you stay focused on searching for *two* points. As To find an $x$-intercept, $y$ must be $0$. This is why we put $0$ in the $y$-value cell of the $x$-intercept row. Similarly, a line's $y$-intercept has $x=0$, and we put $0$ into the $x$-value cell of the $y$-intercept row.

With both intercepts' coordinates, we can graph the line:

This last example generalizes to a fact worth noting.

If a line's $x$-intercept is at $(r,0)$ and its $y$-intercept is at $(0,b)$, then the slope of the line is $-\frac{b}{r}$. (Unless the line passes through the origin, in which case both $r$ and $b$ equal $0$, and then this fraction is undefined. And the slope of the line could be anything.)

What is its $x$-intercept?

*Solution*
To find the $x$-intercept: $\begin{aligned} 2x+4.3y&=\frac{1000}{99} \\ 2x+4.3(0)&=\frac{1000}{99} \\ 2x&=\frac{1000}{99} \\ x&=\frac{500}{99} \end{aligned}$ So the $x$-intercept is at $\left(\frac{500}{99},0\right)$.

What is its $y$-intercept?

*Solution*
To find the $y$-intercept: $\begin{aligned} 2x+4.3y&=\frac{1000}{99} \\ 2(0)+4.3y&=\frac{1000}{99} \\ 4.3y&=\frac{1000}{99} \\ y&=\multiplyleft{\frac{1}{4.3}}\frac{1000}{99} \\ y&\approx2.349\ldots \end{aligned}$ So the $y$-intercept is at about $(0,2.349)$.

What is its slope?

*Solution*
Since we have the $x$- and $y$-intercepts, we can calculate the slope: $m\approx-\frac{2.349}{\frac{500}{99}}=-\frac{2.349\cdot99}{500}\approx-0.4561$.

## Transforming between Standard Form and Slope-Intercept Form

Sometimes a linear equation arises in standard form, but it would be useful to see that equation in slope-intercept form. Or perhaps, vice versa.

**Example**
Change $2x-3y=-6$ to slope-intercept form, and then graph it. (This is the same equation from).

*Solution*
Since a line in slope-intercept form looks like $y=\ldots$, we will isolate $y$:

Now we can see that the slope is $\frac{2}{3}$ and the $y$-intercept is at $(0,2)$. With these things found, we can graph the line using slope triangles.

**Example**
Graph $2x-3y=0$.

*Solution*
First, we will try (and fail) to graph this line using its $x$- and $y$-intercepts.

Trying to find the $x$-intercept: 
\[ \begin{aligned} 2x-3y&=0 \\ 2x-3(0)&=0 \\ 2x&=0 \\ x&=0 \end{aligned} \]
 So the line's $x$-intercept is at $(0,0)$, at the origin. Hmm, the origin is *also* on the $y$-axis So we've also found the $y$-intercept even though that was not the immediate goal.

Since both intercepts are the same point, there is no way to use the intercepts alone to graph this line. So what can be done?

If $C=0$ in a standard form equation, it's convenient to graph it by first converting the equation to slope-intercept form.

**Example**
Write the equation $y=\frac{2}{3}x+2$ in standard form.

*Solution*
Once we subtract $\frac{2}{3}x$ on both sides of the equation, we have 
\[-\frac{2}{3}x+y=2\]
 Technically, this equation is already in standard form $Ax+By=C$. However, you might like to end up with an equation that has no fractions, so you could multiply each side by $3$: 
\[ \begin{aligned} \multiplyleft{3}\left(-\frac{2}{3}x+y\right)&=\multiplyleft{3}2 \\ -2x+3y&=6 \end{aligned} \]

