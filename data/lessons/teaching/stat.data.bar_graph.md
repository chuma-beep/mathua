> Content sourced from [ORCCA](https://pcc.edu/orcca) (Open Resources for Community College Algebra) — CC BY 4.0

# Overview of Graphing

In this section, we will review how to graph lines and general functions which will be useful when we graph parabolas in the next section.

## Graphing Lines by Plotting Points

Sometimes, the easiest way to make a graph of an equation is by making a table and plotting points. (This was the approach in  Section .) Let's refresh ourselves on how this works.

**Example**

A bathtub is holding $12$ gallons of water. The drain starts to leak water at a constant rate of $0.6$ gallons per second. A linear function with formula $W(x)=-0.6x+12$ can be used to model the amount of water, in gallons, in the tub $x$ seconds after it started draining. Let's make a graph of this function. The most straightforward method to graph any function is to build a table of $x$- and $y$-values, and then plot the points.

*A table of values for $W(x)=-0.6x+12$*

*A graph of $W(x)=-0.6x+12$*

**Exercise**

## Graphing Lines in Slope-Intercept Form

Recall that the  slope-intercept form  of a line equation is $y=mx+b$ where $m$ is the slope and $(0,b)$ is the vertical intercept.

**Example**

An efficient method to graph $y=-0.6x+12$ is to use the fact that it is in slope-intercept form. To quickly make a graph, examine the equation and pick out the slope (in this case $-0.6$) and vertical intercept (in this case $(0,12)$), and then plot slope-triangles from the intercept to locate more points on the line. One key point here is that it helps to have the  slope  written as a fraction. In this case,

$$
-0.6=-\frac{6}{10}=-\frac{3}{5}
$$

. So start our graph at $(0,12)$ and go forward $5$ units and then down $3$ units to reach more points.

Since we know that we will go forward $5$ units and then down $3$ units, and that we will start our graph at $(0,12)$, we can choose to orient and scale our axes to see a more complete picture of $W$ than we achieved by plotting convenient points in  Example .

*A graph of $W(x)=-0.6x+12$*

**Example**

Find the slope and vertical intercept of $y=h(x)$, where $h(x)=\frac{5}{3}x-4$. Then use slope triangles to find two more points on the line and sketch it.

The slope is $\frac{5}{3}$ and the vertical intercept is $(0,-4)$. Starting at $(0,-4)$, we go forward $3$ units and up $5$ units to reach more points: $(3,1)$ and $(6,6)$.

*A graph of $h(x)=\frac{5}{3}x-4$*

## Graphing Lines in Point-Slope Form

Recall that the  point-slope form  of a line equation is $y=m(x-x_0)+y_0$ where $m$ is the slope and $(x_0,y_0)$ is a point on the line. The reason that $(x_0,y_0)$ is a point on the line is because you can substitute in $x_0$ for $x$ and then $y_0$ is the result for $y$.

$$
 \underset{\overset{\downarrow}{y_0}}{\strut y \strut}=m(\overset{\overset{x_0}{\downarrow}}{\strut x \strut}-x_0)+y_0 
$$

**Example**

The population of Monarch butterflies has been on  since the 1980s, as have populations of many migratory animals. Efforts to restore the population haven't had great success yet. There are several distinct populations of Monarchs that probably never meet each other: the Hawaii population, the Florida Keys population, the Western population, and the Eastern population. Of these, the Eastern population is by far the largest and we can model this population of Monarch butterflies with a simple linear function.

$$
M(x)=-(x-2006)+15
$$

approximates the total number of acres of Mexican forest that the Eastern population of Monarchs hibernates in during winter in year $x$. This formula is only valid from 1995 to 2018, the years that the population has been well studied.

Let's make graph of this equation given the information provided, but only between 1995 and 2018.

*A graph of $M(x)=-(x-2006)+15$*

**Example**

Find the slope and a point on the graph of $y=m(x)$, where $m(x)=-\frac{9}{5}(x+1)-3$. Then use slope triangles to find two more points on the line and sketch it.

The slope of the line is $-\frac{9}{5}$, and the point given by the equation is $(-1,-3)$. So to graph $h$, start at $(-1,-3)$, and the go forward $5$ units and down $9$ units to reach more points: $(4,-12)$ and $(9,-21)$.

*A graph of $m(x)=-\frac{9}{5}(x+1)-3$*

## Graphing Lines Using Intercepts

Recall that the  standard form  of a line equation is $Ax+By=C$ where where $A$, $B$, and $C$ are three numbers (each of which might be $0$, although at least one of $A$ and $B$ must be nonzero).

**Example**

Recall our bathtub draining problem from  Example , where $W(x)=-0.6x+12$ modeled the amount of water, in gallons, in the tub $x$ seconds after it started draining. Let's write the line equation $y=-0.6x+12$ in standard form.

To find the standard form of the equation, we do as in  Subsection . First, we will replace the variable $W(x)$ with $y$ because standard form relates $x$ and $y$ and does not use function notation. So $W(x)=-0.6x+12$ becomes $y=-0.6x+12$. Now to convert to standard form, move both $x$ and $y$ to the left-hand side.

$$
\begin{aligned}y&=-0.6x+12 \\ 0.6x+y&=12\end{aligned}
$$

The equation is in standard form written as $0.6x+y=12$.

If a linear function is given in standard form, we can relative easily find the equation's $x$- and $y$-intercepts by substituting in $y=0$ and $x=0$, respectively.

**Example**

Let's find the intercepts of $0.6x+y=12$, still relating back to  Example . Then we may graph the equation using those intercepts.

To find the $x$-intercept, set $y=0$ and solve for $x$.

$$
\begin{aligned}0.6x+\substitute{(0)}&=12 \\ 0.6x&=12 \\ x&=20\end{aligned}
$$

So the $x$-intercept is the point $(\firsthighlight{20},\secondhighlight{0})$. In context, this means that $\firsthighlight{20}$ minutes after the tub started to drain, $\secondhighlight{0}$ gallons of water remained. This is telling us that the tub is empty!

To find the $y$-intercept, set $x=0$ and solve for $y$.

$$
\begin{aligned}0.6\substitute{(0)}+y&=12 \\ y&=12\end{aligned}
$$

So, the $y$-intercept is the point $(\firsthighlight{0},\secondhighlight{12})$. In context, this means that $\firsthighlight{0}$ minutes after the tub started to drain, $\secondhighlight{12}$ gallons of water remained. This is telling us how much water was initially in the tub.

When you sketch the graph of a straight line on paper using a straight edge, even having the straight edge off by a very small angle can have a large effect on where the line is drawn, and make your sketch too inaccurate for many purposes. To protect against this, you can find a third point on the line. To that end, we choose $x=5$, then solve for $y$ in the given equation.

$$
\begin{aligned}0.6(\substitute{5})+y&=12 \\ 3+y&=12 \\ y&=9\end{aligned}
$$

Note we chose $x=5$ because it would end up with an integer for $y$. So another point on the line must be $(\firsthighlight{5},\secondhighlight{9})$. In context, this means that $\firsthighlight{5}$ minutes after the tub started to drain, $\secondhighlight{9}$ gallons of water remained.

Now with the $x$- and $y$-intercepts known along with an additional point, we may plot these points and draw the line that runs through them.

*A graph of $3x+5y=60$*

**Exercise**

## Graphing Functions by Plotting Points

Any function, linear or not, can be graphed by building a table of $x$- and $y$-values and plotting points. Let's look at a few more examples.

**Example**

Imagine a company called Corduroy's-Я-Us that makes pants. Their profit from their Royal Blue Corduroys, in thousands of dollars, can be modeled by the function $P(x)=-0.5x^2+33x-200$ where $x$ is the price of each pair of Royal Blue pants that they sell. Let's build a table of values and plot the function's graph.

In this context, the value of $x$ must be positive. Furthermore, we shouldn't really consider $x$-values like $1$, $2$, etc., because it is not realistic that the price of a pair of new pants would be so low. Instead we try multiples of $10$: $10$, $20$, etc.

*A table of values for $P(x)=-0.5x^2+33x-200$*

With the values in  Table , we can sketch the graph. Note that we have to estimate the how the graph curves which is a limitation of graphing a function by plotting points compared with using algebraic techniques.

**Exercise**

**Example**

Human-initiated global warming has been the subject of some debate. However, one aspect of the debate is undeniable fact: the amount of atmospheric carbon dioxide ($\mathrm{CO}_2$: ) is being  and is increasing faster and faster. The measured yearly average atmospheric carbon dioxide levels in parts per million (ppm) since 1958 can be very closely approximated by the function $C(x)=244+29\cdot1.0148^x$ where $x$ represents the number of years since the year 1900. Before 1958, the greenhouse gases weren't regularly measured. Create a table of values rounded to the nearest whole number for the carbon dioxide levels since 1958.

Since 1958 is $58$ years since 1900, we will start our table at $x=58$ and go by $10$s up through $x=118$, which would stand for the year 2018.

*A table of values for $C(x)=244+29\cdot1.0148^x$*

*A graph of $C(x)=244+29\cdot1.0148^x$*

*A graph of $C$ with
&#x20;          &#x20;
&#x20;           overlaid and the function extrapolated beyond known dates.*

##

**Exercise**

What are the four methods we recalled to graph lines in this section?

**Exercise**

Why might it be better to represent a line in point-slope form than slope intercept form?

**Exercise**

Explain how an equation for a line given in slope-intercept or point-slope form can be graphed without creating a table of values.

**Exercise**

Describe one or more possible issues you might encounter after creating a table of points for a function and trying to use those points to make a graph.

##

