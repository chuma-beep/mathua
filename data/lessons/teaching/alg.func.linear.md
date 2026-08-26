> Content sourced from [ORCCA](https://pcc.edu/orcca) (Open Resources for Community College Algebra) — CC BY 4.0

# Technical Definition of a Function

In  Section , we discussed a conceptual understanding of functions and  Definition . In this section we'll start with a more technical definition of what is a function, consistent with the ideas from  Section .

## Formally Defining a Function

\*\*Definition\*\*

Function (Technical Definition)

A *function* is a collection of ordered pairs $(x,y)$ such that any particular value of $x$ is paired with at most one value for $y$.

How is this definition consistent with the informal  Definition , which describes a function as a *process*? Well, if you have a collection of ordered pairs $(x,y)$, you can choose to view the left number as an input, and the right value as the output. If the function's name is $f$ and you want to find $f(x)$ for a particular number $x$, look in the collection of ordered pairs to see if $x$ appears among the first coordinates. If it does, then $f(x)$ is the (unique) $y$-value it was paired with. If it does not, then that $x$ is just not in the domain of $f$, because you have no way to determine what $f(x)$ would be.

\*\*Example\*\*

Using  Definition , a function $f$ could be given by $\{(1,4), (2,3), (5,3), (6,1)\}$.

1. What is $f(1)?$ Since the ordered pair $(1,4)$ appears in the collection of ordered pairs, $f(1)=4$.
2. What is $f(2)?$ Since the ordered pair $(2,3)$ appears in the collection of ordered pairs, $f(2)=3$.
3. What is $f(3)?$ None of the ordered pairs in the collection start with $3$, so $f(3)$ is undefined, and we would say that $3$ is not in the domain of $f$.

\*\*Example\*\*

A Function Given as a Table

Consider the function $g$ expressed by  Figure . How is this "a collection of ordered pairs?" With tables the connection is most easily apparent. Pair off each $x$-value with its $y$-value.

In this case, we can view this function as:

$$
 \left\{(12,0.16), (15,3.2), (18,3.2), (21,1.4), (24,0.98)\right\} 
$$

.

\*\*Example\*\*

A Function Given as a Formula

Consider the function $h$ expressed by the formula $h(x)=x^2$. How is this "a collection of ordered pairs?"

This time, the collection is *really big*. Imagine an $x$-value, like $x=2$. We can calculate that $f(2)=2^2=4$. So the input $2$ pairs with the output $4$ and the ordered pair $(2,4)$ is part of the collection.

You could move on to *any*$x$-value, like say $x=2.1$. We can calculate that $f(2.1)=2.1^2=4.41$. So the input $2.1$ pairs with the output $4.41$ and the ordered pair $(2.1,4.41)$ is part of the collection.

The collection is so large that we cannot literally list all the ordered pairs as was done in  Example  and  Example . We just have to imagine this giant collection of ordered pairs. And if it helps to conceptualize it, we know that the ordered pairs $(2,4)$ and $(2.1,4.41)$ are included.

\*\*Example\*\*

A Function Given as a Graph

Consider the functions $p$ and $q$ expressed in  Figure  and  Figure . How is each of these "a collection of ordered pairs?"

In  Figure , we see that $p(1)=4$, $p(2)=3$, $p(5)=3$, and $p(6)=1$. The graph *literally is* the collection

$$
 \left\{(1,4), (2,3), (5,3), (6,1)\right\} 
$$

.

In  Figure , we can see a few whole number function values, like $q(0)=0$ and $q(1)=2$. But the entire curve has infinitely many points on it and we'd never be able to list them all. We just have to imagine the giant collection of ordered pairs. And if it helps to conceptualize it, we know that the ordered pairs $(0,0)$ and $(1,2)$ are included.

\*\*Exercise\*\*

## Identifying What is *Not* a Function

Just because you have a set of order pairs, a table, a graph, or an equation, it does not necessarily mean that you have a function. Conceptually, whatever you have needs to give consistent outputs if you feed it the same input. More technically, the set of ordered pairs is not allowed to have two ordered pairs that have the same $x$-value but different $y$-values.

\*\*Example\*\*

Consider each set of ordered pairs. Does it make a function?

1. $\left\{\left(5,9\right),\left(3,2\right),\left(\frac{1}{2},0.6\right),\left(5,1\right)\right\}$
2. $\left\{\left(-5,12\right),\left(3,7\right),\left(\sqrt{2},1\right),\left(-0.9,4\right)\right\}$
3. $\left\{\left(5,9\right),\left(3,9\right),\left(4.2,\sqrt{2}\right),\left(\frac{4}{3},\frac{1}{2}\right)\right\}$
4. $\left\{\left(5,9\right),\left(0.7,2\right),\left(\sqrt{25},3\right),\left(\frac{2}{3},\frac{3}{2}\right)\right\}$

1) This set of ordered pairs is *not* a function. The problem is that it has both $(5,9)$ and $(5,1)$. It uses the same $x$-value paired with two different $y$-values. We have no clear way to turn the input $5$ into an output.
2) This set of ordered pairs *is* a function. It is a collection of ordered pairs, and the $x$-values are never reused.
3) This set of ordered pairs *is* a function. It is a collection of ordered pairs, and the $x$-values are never reused. You might note that the *output* value $9$ appears twice, but that doesn't matter. That just tells us that the function turns $5$ into $9$ and it also turns $3$ into $9$.
4) This set of ordered pairs is *not* a function, but it's a little tricky. One of the ordered pairs uses $\sqrt{25}$ as an input value. But that is the same as $5$, which is also used as an input value.

Now that we understand how some sets of ordered pairs might not be functions, what about tables, graphs, and equations? If we are handed one of these things, can we tell whether or not it is giving us a function?

\*\*Exercise\*\*

Does This Table Make a Function?

\*\*Example\*\*

Does This Graph Make a Function?

Which of these graphs make $y$ a function of $x$?

The graph in  Figure  does *not* make $y$ a function of $x$. Two ordered pairs on that graph are $(-3,1)$ and $(-3,-2)$, so an input value is used twice with different output values.

The graph in  Figure  does *not* make $y$ a function of $x$. There are many ordered pairs with the same input value but different output values. For example, $(2,-2)$ and $(2,4)$.

The graph in  Figure *does* make $y$ a function of $x$. It appears that no matter what $x$-value you choose on the $x$-axis, there is exactly one $y$-value paired up with it on the graph.

The graph in  Figure *does* make $y$ a function of $x$, but we should discuss. The hollow dots on the line indicate that the line goes right up to that point, but never reaches it. We say there is a "hole" in the graph at these places. For two of these holes, there is a separate ordered pair immediately above or below the hole. The graph has the ordered pair $(-4,4)$. It *also* has ordered pairs like $(\text{very close to }{-4},\text{very close to }0)$, but it does not have $(-4,0)$. Overall, there is no $x$-value that is used twice with different $y$-values, so this graph does make $y$ a function of $x$

The graph in  Figure  does *not* make $y$ a function of $x$. There are many ordered pairs with the same input value but different output values. For example, $(0,1)$, $(0,3)$, $(0,-1)$, $(0,5)$, and $(0,-6)$ all use $x=0$.

The graph in  Figure  does *not* make $y$ a function of $x$. There are many ordered pairs with the same input value but different output values. For example at $x=2$, there is both a positive and a negative associated $y$-value. It's hard to say exactly what these $y$-values are, but we don't have to.

This last set of examples might reveal something to you. For instance in  Figure , the issue is that there are places on the graph with the same $x$-value, but different $y$-values. Visually, what that means is there are places on the graph that are directly above/below each other. Thinking about this leads to a quick visual "test" to determine if a graph gives $y$ as a function of $x$.

\*\*Fact\*\*

Vertical Line Test

Given a graph in the $xy$-plane, if a vertical line ever touches it in more than one place, the graph does *not* give $y$ as a function of $x$. If vertical lines only ever touch the graph once or never at all, then the graph *does* give $y$ as a function of $x$.

\*\*Example\*\*

In each graph from  Example , we can apply the   .

*A vertical line touching the graph twice makes this graph not give $y$ as a function of $x$.*

*A vertical line touching the graph twice makes this graph not give $y$ as a function of $x$.*

*All vertical lines only touch the graph once, so this graph does give $y$ as a function of $x$.*

*All vertical lines only touch the graph once, or not at all, so this graph does give $y$ as a function of $x$.*

*A vertical line touching the graph more than once makes this graph not give $y$ as a function of $x$.*

*A vertical line touching the graph more than once makes this graph not give $y$ as a function of $x$.*

Lastly, we come to equations. Certain equations with variables $x$ and $y$ clearly make $y$ a function of $x$. For example, $y=x^2+1$ says that if you have an $x$-value, all you have to do is substitute it into that equation and you will have determined an output $y$-value. You could then name the function $f$ and give a formula for it: $f(x)=x^2+1$.

With other equations, it may not be immediately clear whether or not they make $y$ a function of $x$.

\*\*Example\*\*

Do each of these equations make $y$ a function of $x$?

1. $2x+3y=5$
2. $y=\pm\sqrt{x+4}$
3. $x^2+y^2=9$

1) The equation $2x+3y=5$*does* make $y$ a function of $x$. Here are three possible explanations.
   1. You recognize that the graph of this equation would be a non-vertical line, and so it would pass the   .
   2. Imagine that you have a specific value for $x$ and you substitute it in to $2x+3y=5$. Will you be able to use algebra to solve for $y$? All you will need is to simplify, subtract from both sides, and divide on both sides, so you will be able to determine $y$.
   3. Can you just isolate $y$ in terms of $x$? Yes, a few steps of algebra can turn $2x+3y=5$ into $y=\frac{5-2x}{3}$. Now you have an explicit formula for $y$ in terms of $x$, so $y$ is a function of $x$.
2) The equation $y=\pm\sqrt{x+4}$ does *not* make $y$ a function of $x$. Just having the $\pm$ (plus *or* minus) in the equation immediately tells you that for almost any valid $x$-value, there would be *two* associated $y$-values.
3) The equation $x^2+y^2=9$ does *not* make $y$ a function of $x$. Here are three possible explanations.
   1. Imagine that you have a specific value for $x$ and you substitute it in to $x^2+y^2=9$. Will you be able to use algebra to solve for $y$? For example, if you substitute in $x=1$, then you have $1+y^2=9$, which simplifies to $y^2=8$. Can you really determine what $y$ is? No, because it could be $\sqrt{8}$ or it could be $-\sqrt{8}$. So this equation does not provide you with a way to turn $x$-values into $y$-values.
   2. Can you just isolate $y$ in terms of $x$? You might get started and use algebra to convert $x^2+y^2=9$ into $y^2=9-x^2$. But what now? The best you can do is acknowledge that $y$ is either the positive or the negative square root of $9 - x^2$. You might write $y=\pm\sqrt{9-x^2}$. But now for almost any valid $x$-value, there are *two* associated $y$-values.
   3. You recognize that the graph of this equation would be a circle with radius $3$, and so it would not pass the   .

\*\*Exercise\*\*

##

\*\*Exercise\*\*

Suppose you have a "relation". That is, a set of order pairs, a table of $x$- and $y$-values, a graph, or an equation in $x$ and $y$. What is the one thing that could happen that would make the relation *not* be a function?

\*\*Exercise\*\*

Explain how to use the vertical line test.

##

