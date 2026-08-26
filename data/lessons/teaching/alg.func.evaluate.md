> Content sourced from [ORCCA](https://pcc.edu/orcca) (Open Resources for Community College Algebra) — CC BY 4.0

# Function Basics

In this section, we will introduce a topic that will be essential for continued mathematical learning: functions. Functions should be thought of as machines that turn one number into another number, much like a cash register can turn a number of pounds of fruit into a price.

## Informal Definition of a Function

We are familiar with the $\sqrt{\phantom{x}}$ symbol. This symbol is used to turn numbers into their square roots. Sometimes it's simple to do this on paper or in our heads, and sometimes it helps to have a calculator. We can see some calculations in  Figure .

*Values of $\sqrt{x}$*

The $\sqrt{\phantom{x}}$ symbol represents a *process*; it's a way for us to turn numbers into other numbers. This idea of having a process for turning numbers into other numbers is the fundamental topic of this chapter.

\*\*Definition\*\*

Function (Informal Definition)

A *function* is a process for turning numbers into (potentially) different numbers. The process must be *consistent*, in that whenever you apply it to some particular number, you always get the same result.

Section  covers a more technical definition for functions, and covers topics that are more appropriate when using that definition.  Definition  is so broad that you probably use functions all the time.

\*\*Example\*\*

In each of these examples, some process is used for turning one number into another.

- If you convert a person's birth year into their age, you are using a function.
- If you look up the Kelly Blue Book value of a Honda Odyssey based on how old it is, you are using a function.
- If you use the expected guest count for a party to determine how many pizzas you should order, you are using a function.

The $\sqrt{\phantom{x}}$ function is consistent; for example, every time you evaluate $\sqrt{9}$, you always get 3. One interesting fact is that $\sqrt{\phantom{x}}$ is not found on most keyboards, and yet computers can still find square roots. Computer technicians write $\operatorname{sqrt}(\phantom{x})$ when they want to compute a square root, as we see in  Figure .

*Values of $\operatorname{sqrt}(x)$*

The parentheses in $\operatorname{sqrt}\highlight{(\phantom{x})}$ are very important. To see why, try to put yourself in the "mind" of a computer. The computer will recognize `sqrt` and know that it needs to compute a square root but without parentheses it will think that it needs to compute `sqrt4` and then put a `9` on the end, which would produce a final result of $29$. This is probably not what was intended. And so the purpose of the parentheses in `sqrt(49)` is to  be deliberately clear.

Functions have their own names. We've seen a function named $\operatorname{sqrt}$, but any name you can imagine is allowable. In the sciences, it is common to name functions with whole words, like $\operatorname{weight}$ or $\operatorname{health\_index}$. In math, we often abbreviate such function names to $w$ or $h$. And of course, since the word "function" itself starts with "f," we will often name a function $f$.

\*\*Warning\*\*

Notation Ambiguity

In some contexts, the symbol $t$ might represent a variable (a number that is represented by a letter) and in other contexts, $t$ might represent a function (a process for changing numbers into other numbers). By staying conscious of the *context* of an investigation, we avoid confusion.

Next we need to discuss how we go about using a function's name.

\*\*Definition\*\*

Function Notation

The standard notation for referring to functions involves giving the function itself a name, and then writing:

$$
 \begin{matrix} \text{name}\\ \text{of}\\ \text{function} \end{matrix} \left( \begin{matrix} \\ \text{input}\\ \\ \end{matrix} \right) 
$$

\*\*Example\*\*

$f(13)$ is pronounced "f of 13." The word "of" is very important, because it reminds us that $f$ is a process and we are about to apply that process to the input value $13$. So $f$ is the function, $13$ is the input, and $f(13)$ is the output we'd get from using $13$ as input.

$f(x)$ is pronounced "f of x." This is just like the previous example, except that the input is not any specific number. The value of $x$ could be $13$ or any other number. Whatever $x$'s value, $f(x)$ means the corresponding output from the function $f$.

$\operatorname{BudgetDeficit}(2017)$ is pronounced "BudgetDeficit of 2017." This is probably about a function that takes a year as input, and gives that year's federal budget deficit as output. The process here of changing a year into a dollar amount might not involve any mathematical formula, but rather looking up information from the Congressional Budget Office's website.

\*\*Note\*\*

While a function has a name like $f$, and the input to that function often has a variable name like $x$, the expression $f(x)$ represents the output of the function. To be clear, $f(x)$ is *not* a function. Rather, $f$ is a function, and $f(x)$ its output when the number $x$ was used as input.

\*\*Exercise\*\*

In the following examples, a function is given using a formula, and we will evaluate the function at specific values. See  Section  for a review on evaluating expressions.

\*\*Example\*\*

Let $V$ be the function defined by $V(t)=-5t+1$.

1.
2.
3.

1)
2)
3)

\*\*Example\*\*

Let $L$ be the function defined by $L(z)=2z^2-z+3$.

1.
2.
3.

1)
2)
3)

\*\*Exercise\*\*

\*\*Exercise\*\*

\*\*Warning\*\*

More Notation Ambiguity

As mentioned in  Warning , we need to remain conscious of the context of any symbol we are using.  Consider the expression $a(b)$. This could easily mean the output of a function $a$ with input $b$. It could also mean that two numbers $a$ and $b$ need to be multiplied. It all depends on the context in which these symbols are being used.

Sometimes it's helpful to think of a function as a machine, as in  Figure .  A *function* has the capacity to take in all kinds of different numbers into it's hopper (feeding tray) as inputs and transform them into their outputs.

*Imagining a function as a machine. (Image by Duane Nykamp using Mathematica.)*

## Tables and Graphs

Since functions are potentially complicated, we want ways to understand them more easily. Two basic tools for understanding a function better are tables and graphs.

\*\*Example\*\*

A Table for the Budget Deficit Function

Consider the function $\operatorname{BudgetDeficit}$, that takes a year as its input and outputs the US federal budget deficit for that year. For example, the Congressional Budget Office's website tells us that $\operatorname{BudgetDeficit}(2009)$ is $$\$1.41$$ trillion. If we'd like to understand this function better, we might make a table of all the inputs and outputs we can find. Using the CBO's website (), we can put together  Table .

How is this table helpful? There are things about the function that we can see now by looking at the numbers in this table.

- We can see that the budget deficit had a spike between 2008 and 2009.
- And it fell again between 2012 and 2013.
- It appears to stay roughly steady for several years at a time, with occasional big jumps or drops.

These observations help us understand the function $\operatorname{BudgetDeficit}$ a little better.

\*\*Exercise\*\*

\*\*Example\*\*

A Table for the Square Root Function

Let's return to our example of the function $\operatorname{sqrt}$. Tabulating some inputs and outputs reveals  Figure .

How is this table helpful? Here are some observations that we can make now.

- We can see that when input numbers increase, so do output numbers.
- We can see even though outputs are increasing, they increase by less and less with each step forward in $x$.

These observations help us understand $\operatorname{sqrt}$ a little better. For instance, based on these observations which do you think is larger: the difference between $\operatorname{sqrt}(23)$ and $\operatorname{sqrt}(24)$, or the difference between $\operatorname{sqrt}(85)$ and $\operatorname{sqrt}(86)$?

\*\*Exercise\*\*

Another powerful tool for understanding a function better is a graph. Given a function $f$, one way to make its graph is to take a table of input and output values, and read each row as the coordinates of a point in the $xy$-plane.

\*\*Example\*\*

A Graph for the Budget Deficit Function

Returning to the function $\operatorname{BudgetDeficit}$ that we studied in  Example , in order to make a graph of this function we view  Table  as a list of points with $x$ and $y$ coordinates, as in  Figure . We then plot these points on a set of coordinate axes, as in  Figure . The points have been connected with a curve so that we can see the overall pattern given by the progression of points. Since there was not any actual data for inputs in between any two years, the curve is dashed. That is, this curve is dashed because it just represents someone's best guess as to how to connect the plotted points. Only the plotted points themselves are precise.

How has this graph helped us to understand the function better? All of the observations that we made in  Example  are perhaps even more clear now. For instance, the spike in the deficit between 2008 and 2009 is now visually apparent. Seeking an explanation for this spike, we recall that there was a financial crisis in late 2008. Revenue from income taxes dropped at the same time that federal money was spent to prevent further losses.

\*\*Example\*\*

A Graph for the Square Root Function

Let's now construct a graph for $\operatorname{sqrt}$. Tabulating inputs and outputs gives the points in  Figure , which in turn gives us the graph in  Figure .

Just as in the previous example, we've plotted points where we have concrete coordinates, and then we have made our best attempt to connect those points with a curve. Unlike the previous example, here we believe that points will continue to follow the same pattern indefinitely to the right, and so we have added an arrowhead to the graph.

What has this graph done to improve our understanding of $\operatorname{sqrt}$? As inputs ($x$-values) increase, the outputs ($y$-values) increase too, although not at the same rate. In fact we can see that our graph is steep on its left, and less steep as we move to the right. This confirms our earlier observation in  Example  that outputs increase by smaller and smaller amounts as the input increases.

\*\*Remark\*\*

Graph of a Function

Given a function $f$, when we refer to a *graph of $f$* we are *not* referring to an entire picture, like  Figure . A graph of $f$ is only *part* of that picture---the curve and the points that it connects. Everything else (axes, tick marks, the grid, labels, and the surrounding white space) is just useful decoration so that we can read the graph more easily.

\*\*Remark\*\*

A Common Wording Misunderstanding

It is common to refer to the graph of $f$ as the *graph of the equation $y=f(x)$*. However, we should avoid saying "the graph of $f(x)$." That would indicate a  misunderstanding of our notation. Since$f(x)$ is the output for a certain input $x$. That means that $f(x)$ is just a number  and not worthy of a two-dimensional picture.

While it is important to be able to make a graph of a function $f$, we also need to be capable of looking at a graph and reading it well. A graph of $f$ provides us with helpful specific information about $f$; it tells us what $f$ does to its input values. When we were making graphs, we plotted points of the form

$$
 (\text{input},\text{output}) 
$$

Now given a graph of $f$, we interpret coordinates in the same way.

\*\*Example\*\*

In  Figure  we have a graph of a function $f$. If we wish to find $f(1)$, we recognize that $1$ is being used as an input. So we would want to find a point of the form $(1,\phantom{y})$. Seeking out $x$-coordinate $1$ in  Figure , we find that the only such point is $(1,2)$. Therefore the output for $1$ is $2$; in other words $f(1)=2$.

\*\*Exercise\*\*

\*\*Example\*\*

Unemployment Rates

Suppose that $u$ is the unemployment function of time. That is, $u(t)$ is the unemployment rate in the United States in year $t$. The graph of the equation $y=u(t)$ is given in  Figure  ().

*Unemployment in the United States*

What was the unemployment in 2008? It is a straightforward matter to use  Figure  to find that unemployment was almost $6\%$ in 2008. Asking this question is exactly the same thing as asking to find $u(2008)$. That is, we have one question that can either be asked in an everyday-English way or which can be asked in a terse, mathematical notation-heavy way:

"What was unemployment in 2008?"

"Find $u(2008)$."

If we use the table to establish that $u(2009)\approx9.25$, then we should be prepared to translate that into everyday-English using the context of the function: In 2009, unemployment in the  was about $9.25\%$.

If we ask the question "when was unemployment at $5\%$," we can read the graph and see that there were two such times: mid-2007 and about 2016. But there is again a more mathematical notation-heavy way to ask this question. Namely, since we are being told that the output of $u$ is $5$, we are being asked to solve the equation $u(t)=5$. So the following communicate the same thing:

"When was unemployment at $5\%$?"

"Solve the equation $u(t)=5$."

And our answer to this question is:

"Unemployment was at $5\%$ in mid-2007 and about 2016."

"$t\approx2007.5$ or $t\approx2016$."

\*\*Exercise\*\*

## Translating Between Four Descriptions of the Same Function

We have noted that functions are complicated, and we want ways to make them easier to understand. It's common to find a problem involving a function and not know how to find a solution to that problem. Most functions have at least four standard ways to think about them, and if we learn how to translate between these four perspectives, we often find that one of them makes a given problem easier to solve.

The four modes for working with a given function are

- a verbal description
- a table of inputs and outputs
- a graph of the function
- a formula for the function

This has been visualized in  Figure .

*Function Perspectives*

\*\*Example\*\*

Consider a function $f$ that squares its input and then adds $1$. Translate this verbal description of $f$ into a table, a graph, and a formula.

To make a table for $f$, we'll have to select some input $x$-values. These choices are left entirely up to us, so we might as well choose small, easy-to-work-with values. However we shouldn't shy away from negative input values. Given the verbal description, we should be able to compute a column of output values.  Figure  is one possible table that we might end up with.

Once we have a table for $f$, we can make a graph for $f$ as in  Figure , using the table to plot points.

Lastly, we must find a formula for $f$. This means we need to write an algebraic expression that says the same thing about $f$ as the verbal description, the table, and the graph. For this example, we can focus on the verbal description. Since $f$ takes its input, squares it, and adds $1$, we have that

$$
 f(x)=x^2+1 
$$

.

\*\*Example\*\*

Let $F$ be the function that takes a Celsius temperature as input and outputs the corresponding Fahrenheit temperature. Translate this verbal description of $F$ into a table, a graph, and a formula.

To make a table for $F$, we will need to rely on what we know about Celsius and Fahrenheit temperatures. It is a fact that the freezing temperature of water at sea level is $0\,\text{celsius}$, which equals $32\,\text{fahrenheit}$. Also, the boiling temperature of water at sea level is $100\,\text{celsius}$, which is the same as $212\,\text{fahrenheit}$. One more piece of information we might have is that standard human body temperature is $37\,\text{celsius}$, or $98.6\,\text{fahrenheit}$. All of this is compiled in  Figure . Note that we tabulated inputs and outputs by working with the context of the function, not with any computations.

Once a table is established, making a graph by plotting points is a simple matter, as in  Figure . The three plotted points seem to be in a straight line, so we think it is reasonable to connect them in that way.

To find a formula for $F$, the verbal definition is not of much direct help. But $F$'s graph does seem to be a straight line. And linear equations are familiar to us. This line has a $y$-intercept at $(0,32)$ and a slope we can calculate: $\frac{212-32}{100-0}=\frac{180}{100}=\frac{9}{5}$. So the equation of this line is $y=\frac{9}{5}C+32$. On the other hand, the equation of this graph is $y=F(C)$, since it is a graph of the function $F$. So evidently,

$$
 F(C)=\frac{9}{5}C+32 
$$

.

##

\*\*Exercise\*\*

When $g$ is a function, how should you say out loud "$g(x)$?"

\*\*Exercise\*\*

There are four main ways to communicate how a function turns its inputs into its outputs. What are they?

\*\*Exercise\*\*

What is usually an acceptable way to type "the square root of $x$" if you have to type it using a regular keyboard?

##

\*\*Exercise\*\*

\*\*Exercise\*\*

\*\*Exercise\*\*

\*\*Exercise\*\*

