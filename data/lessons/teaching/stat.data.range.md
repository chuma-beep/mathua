> Content sourced from [ORCCA](https://pcc.edu/orcca) (Open Resources for Community College Algebra) — CC BY 4.0
# Exploring Two-Variable Data and Rate of Change

This section is about making observations about data and the patterns that you might see. In some cases, we are able to turn those observations into useful mathematical calculations.

## Modeling data with two variables

Math helps us understand data from the world around us. We can use what we discover to understand the world better and make better decisions. Here's an example with economic data from the US, plotted in a Cartesian plane.

For the years from 2000 to 2021, consider what percent of American wealth was held by the wealthiest 1% of Americans. The table in gives the numbers (source: https://fee.org/articles/the-top-1-hold-a-record-amount-of-wealth-in-the-us-here-s-how-much-and-why/), but any pattern there might not be apparent when looking at the data organized this way. Plotting the data in a Cartesian coordinates system can make an overall pattern or trend become visible.

What observations do you see now that you couldn't easily see from the numbers in the table? Do you see evidence of the COVID pandemic? Evidence of the Great Recession of 2008?

Overall, do you see a larger pattern with wealth distribution? Assuming that you see the rising pattern, is it easier to see that with the graph than with the table?

## Patterns in Tables

**Example**
Find a pattern in each table, using only the table itself. What is the missing entry in each table? Can you describe each pattern in words and/or mathematics?

*Solution*
First table Each word on the right has the opposite meaning of the word to its left. Second table Each city on the right is the capital of the country to its left. Third table Each number on the right is double the number to its left.

Generally in a table with two columns of data, we can think of the table as *assigning* value on the right to each value on the left. The first table *assigns* "white" to "black", as its opposite. The second table *assigns* "Paris" to "France", as its capital city. The third table *assigns* $10$ to $5$, as its double.

Only the third table in is a table of numbers. Let's examine that data graphically.

With the data plotted, and the question being what should happen when $x$ is $5$, our eyes can converge to the point $(5,10)$ and we conclude the missing value will be $10$. Graphically, we didn't have to use the observation that the $y$-values were twice the $x$-values.

For each of the following tables, find an equation that describes the pattern you see. Numerical pattern recognition may or may not come naturally for you and you may want to use a graph to help visually process the numbers. Either way, pattern recognition is an important mathematical skill that anyone can develop. The solutions for these exercises offer some hints about what patterns you might look for.

## Rate of Change

For an hourly wage-earner, the amount of money they earn depends on how many hours they work. If a worker earns $\$15$ per hour, then $10$ hours of work corresponds to $\$150$ of pay. Working *one* additional hour will change $10$ hours to $11$ hours; and this will cause the $\$150$ in pay to rise by *fifteen* dollars to $\$165$ in pay. Any time we compare how one amount changes (dollars earned) as a consequence of another amount changing (hours worked), we are talking about a *rate of change*.

Given a table of two-variable data, between any two rows we can compute a *rate of change*.

**Example**

Use the data in to find the rate of change in Oregon invasive cancer diagnoses between 2000 and 2003.

*Solution*
To find the rate of change between 2000 and 2003, calculate $\frac{17590 - 17458}{2003 - 2000} = 44$. So the rate of change was .

And what was the rate of change between 2015 and 2020?

*Solution*
To find the rate of change between 2005 and 2020, calculate $\frac{20151 - 22154}{2020 - 2015} = -400.6$. So the rate of change was .

We are ready to give a formal definition for "rate of change". Considering our work from and, we settle on:

> **Definition**
> If $\left(x_1,y_1\right)$ and $\left(x_2,y_2\right)$ are two data points from a set of two-variable data, then the *rate of change* between them is 
\[\frac{\text{change in \$y\$}}{\text{change in \$x\$}}=\frac{\Delta y}{\Delta x}=\frac{y_2-y_1}{x_2-x_1}\]. The Greek letter delta, $\Delta$, is used to represent "change in" since it is the first letter of the Greek word for "difference".

In and we found three rates of change. highlights the three pairs of points that were used to make these calculations.

Note how the larger the numerical rate of change between two points, the steeper the line is that connects them. Also when the $y$-values went down as you read the graph left-to-right, the rate of change was negative. This is such an important observation, we'll put it in an official remark.

Whenever the rate of change is constant no matter which two $(x,y)$-pairs (or data pairs) are chosen from a data set, then you can conclude the graph will be a straight line *even without making the graph*. We call this kind of relationship a *linear* relationship. We'll study linear relationships in more detail throughout this chapter. Right now in this section, we feel it is important to simply identify if data has a linear relationship or not.

Let's return to the data that we opened the section with, in. Is that data linear? Well, yes and no. To be completely honest, it's not linear. It's easy to pick out pairs of points where the steepness changes from one pair to the next. In other words, the points do not line up into a single straight line.

However if we step back, there does seem to be an overall upward trend that is captured by the line overlaying the data in. Points on the overlaid line *do* have a linear pattern. Let's estimate the rate of change between some pair of points on this line. We are free to use any pair of points to do this, so let's make this calculation easier by choosing points we can clearly identify on the graph: $(2005,27.7)$ and $(2018,31)$.

The rate of change between those two points is 
\[\frac{(31-27.7)\,\text{pct. points}}{(2018-2005)\,\text{years}}=\frac{3.4\,\text{pct. points}}{13\,\text{years}}\approx0.2615\,\frac{\text{pct. points}}{\text{year}}\]
 So we might say that *on average*, the rate of change expressed by this data is $0.2615$ percentage points per year.
