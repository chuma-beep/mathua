> Content sourced from [Algebrica](https://algebrica.org/lines/) — CC BY-NC 4.0

## What are lines

A line is a fundamental geometric object made up of infinitely many points aligned in a perfectly straight path. It has no thickness, no endpoints, and it extends endlessly in both directions. A line is completely determined by any two of its points, and represents the most basic, unchanging connection in space.


The **implicit form** of a line is written as:

\\[ax + by + c = 0\\]

where \\(a\\), \\(b\\), and \\(c\\) are real numbers, and at least one of \\(a\\) or \\(b\\) is non-zero. This equation is called implicit because both variables appear on the same side of the equation, and the relationship between \\(x\\) and \\(y\\) is not isolated.

## Parallel and perpendicular lines

A line is **parallel** to the y-axis when it runs vertically and all of its points share the same x-coordinate. This type of line does not move left or right as you go up or down. It remains perfectly vertical. Its equation is:

\\[x = k\\]

![](/diagrams/algebrica/line-1.png)

A line is parallel to the x-axis when it runs horizontally and all of its points share the same y-coordinate. This type of line does not move up or down as you go left or right. It remains perfectly horizontal. Its equation is:

\\[y = k \\]

![](/diagrams/algebrica/line-3.png)

##### Horizontal lines have a slope of zero and extend infinitely in both directions.

## General equation of a line

In general terms, the **equation of a line** can be written in explicit form as:

\\[y = mx + q \\]

![](/diagrams/algebrica/lines-4.png)

\\(m\\) is the slope of the line and \\(q\\) is the y-intercept. The x-intercept of a line is the value of x for which y = 0 (in other words, it is the root of the line’s equation).


A point lies on a line if and only if its coordinates satisfy the equation of the line. This means that when you substitute the x- and y-values of the point into the equation, both sides of the equation remain equal. For example, the point \\((2, 5)\\) lies on the line \\( y = 2x + 1 \\) because:

\\[y = 2(2) + 1 = 5 \\]

So the equation holds true, and the point belongs to the line.


Two lines \\( r \\) and \\( s \\) are **parallel** if they have the same slope \\( m_r = m_s \\). They are **perpendicular** if their slopes are negative reciprocals \\( m_r = -\dfrac{1}{m_s} \\)

##### The slope ( m ) is undefined for lines parallel to the y-axis, and equal to 0 for lines parallel to the x-axis.

A [linear equation](<../linear-equation>) in two variables, written as \\( y = mx + q \\), expresses a direct relationship between \\( x \\) and \\( y \\). This equation corresponds to a straight line in the coordinate plane, with its structure determining the line’s position and inclination.

## Distance from a point to a line

The distance from a point \\( P(x_P, y_P) \\) to a line \\( r \\) given by the equation \\(ax + by + c = 0 \\) is the length of the segment connecting the point \\( P \\) to the foot of the perpendicular dropped from \\( P \\) onto the line.

![](/diagrams/algebrica/lines-6.png)

This distance is calculated using the formula:

\\[d = \frac{\left| ax_P + by_P + c \right|}{\sqrt{a^2 + b^2}} \\]

##### This expression gives the shortest distance from the point to the line, that is, the perpendicular distance, not the length of any random segment.

## Line passing through two points

Consider the line passing through two points \\( P(x_P, y_P) \\) and \\( Q(x_Q, y_Q) \\). If \\( x_P = x_Q \\), the line is parallel to the y-axis and its equation is: \\[x = x_P \\]


If \\( x_P \ne x_Q \\), the line has a slope \\( m \\) given by: \\[m = \frac{y_Q - y_P}{x_Q - x_P} \\] Its equation is:

\\[y - y_P = m(x - x_P) \\]

This can also be written in the symmetric form:

\\[\frac{y - y_P}{y_Q - y_P} = \frac{x - x_P}{x_Q - x_P} \\]

## Example 1

Let’s find the equation of the line that passes through the points:

\\[P(1, 2) \quad \text{and} \quad Q(3, 6) \\]

To begin, we calculate the slope of the line. The slope \\( m \\) is the ratio between the difference in the y-values and the difference in the x-values of the two points:

\\[m = \frac{y_Q - y_P}{x_Q - x_P} = \frac{6 - 2}{3 - 1} = \frac{4}{2} = 2 \\]

Now that we know the slope is 2, we can write the equation of the line using the point-slope form. We choose point \\( P(1, 2) \\) and plug the values into the formula:

\\[y - 2 = 2(x - 1) \\]

We can leave the equation in this form, or we can expand it into slope-intercept form:

\\[y = 2x - 2 + 2 = 2x \\]

So the line passing through \\( P(1, 2) \\) and \\( Q(3, 6) \\) has the equation:

\\(y = 2x\\)

## Intersection of two lines

If the lines are parallel, they never intersect, because they have the same slope and never cross each other. Two lines that are not parallel intersect at a single point. The coordinates of the **point of intersection** are the solutions of the [system](<../systems-of-linear-equations/>) formed by the equations of the two lines.

## Example 2

For example, let’s consider the following line in slope-intercept form:

\\[y = 2x + 1 \\]

Now let’s take another line with a different slope, so that they are not parallel and will intersect: \\[y = -x + 4 \\]


To find the point of intersection, we solve the system formed by the two equations:

\\[\begin{cases} y = 2x + 1 \\\\[0.5em] y = -x + 4 \end{cases} \\]

By setting the right-hand sides equal to each other, we obtain:

\\[\begin{align} &2x + 1 = -x + 4 \\\\[0.5em] &3x = 3 \\\\[0.5em] &x = 1 \end{align} \\]


Substituting \\( x = 1 \\) into one of the original equations, we find:

\\[y = 2(1) + 1 = 3 \\]

![](/diagrams/algebrica/lines-5.png)

So, the two lines intersect at the point:

\\((x=1, y=3)\\)
