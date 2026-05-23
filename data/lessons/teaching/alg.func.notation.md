> Content sourced from [ORCCA](https://pcc.edu/orcca) (Open Resources for Community College Algebra) — CC BY 4.0
# Function Basics

In this section, we will introduce a topic that will be essential for continued mathematical learning: functions.
      Functions should be thought of as machines that turn one number into another number,
      much like a cash register can turn a number of pounds of fruit into a price.

## Informal Definition of a Function

The $\sqrt{\phantom{x}}$ symbol represents a *process* ;
      it's a way for us to turn numbers into other numbers.
      This idea of having a process for turning numbers into other numbers is the fundamental topic of this chapter.

> **Definition**
> A *function* is a process for turning numbers into (potentially) different numbers.
          The process must be *consistent* ,
          in that whenever you apply it to some particular number,
          you always get the same result.

Section covers a more technical definition for functions,
      and covers topics that are more appropriate when using that definition. Definition is so broad that you probably use functions all the time.

**Example**

The parentheses in $\operatorname{sqrt}\highlight{(\phantom{x})}$ are very important.
      To see why, try to put yourself in the "mind" of a computer.
      The computer will recognize `sqrt` and know that it needs to compute a square root but without parentheses it will think that it needs to compute `sqrt4` and then put a `9` on the end,
      which would produce a final result of $29$ .
      This is probably not what was intended.
      And so the purpose of the parentheses in `sqrt(49)` is to  be deliberately clear.

Functions have their own names.
      We've seen a function named $\operatorname{sqrt}$ ,
      but any name you can imagine is allowable.
      In the sciences, it is common to name functions with whole words,
      like $\operatorname{weight}$ or $\operatorname{health\_index}$ .
      In math,
      we often abbreviate such function names to $w$ or $h$ .
      And of course, since the word "function" itself starts with "f," we will often name a function $f$ .

> **Notation Ambiguity**
> In some contexts,
        the symbol $t$ might represent a variable (a number that is represented by a letter) and in other contexts, $t$ might represent a function (a process for changing numbers into other numbers).
        By staying conscious of the *context* of an investigation,
        we avoid confusion.

Next we need to discuss how we go about using a function's name.

> **Definition**
> The standard notation for referring to functions involves giving the function itself a name,
          and then writing: \[
            \begin{matrix}
            \text{name}\\
            \text{of}\\
            \text{function}
            \end{matrix}
            \left(
            \begin{matrix}
            \\
            \text{input}\\
            \\
            \end{matrix}
            \right)
          \]

**Example**


In the following examples, a function is given using a formula,
      and we will evaluate the function at specific values.
      See Section for a review on evaluating expressions.

**Example**
Let $V$ be the function defined by $V(t)=-5t+1$ . Find $V(4)$ . Find $V(0)$ . Find $V\mathopen{}\left(-\frac{1}{3}\right)\mathclose{}$ .

*Solution*
$\begin{aligned}
                    V(\substitute{4})\amp = -5(\substitute{4})+1\\
                    \amp = -20+1 = -19
                    \end{aligned}$ $\begin{aligned}
                    V(\substitute{0})\amp = -5(\substitute{0})+1\\
                    \amp = 0+1 = 1
                    \end{aligned}$ $\begin{aligned}
                    V\mathopen{}\left(\substitute{-\frac{1}{3}}\right)\mathclose{}\amp = -5\left(\substitute{-\frac{1}{3}}\right)+1\\
                    \amp = \frac{5}{3}+1\\
                    \amp = \frac{5}{3}+\frac{3}{3} = \frac{8}{3}
                    \end{aligned}$

**Example**
Let $L$ be the function defined by $L(z)=2z^2-z+3$ . Find $L(1)$ . Find $L(-9)$ . Find $L\mathopen{}\left(\frac{3}{2}\right)\mathclose{}$ .

*Solution*
$\begin{aligned}
                    L(\substitute{1}) \amp =2(\substitute{1})^2-(\substitute{1})+3\\
                    \amp = 2(1)-1+3\\
                    \amp = 2-1+3=4
                    \end{aligned}$ $\begin{aligned}
                    L(\substitute{-9}) \amp =2(\substitute{-9})^2-(\substitute{-9})+3\\
                    \amp = 2(81)+9+3\\
                    \amp = 162+9+3=174
                    \end{aligned}$ $\begin{aligned}
                    L\mathopen{}\left(\substitute{\frac{3}{2}}\right)\mahtclose{}\amp=2\left(\substitute{\frac{3}{2}}\right)^2-\left(\substitute{\frac{3}{2}}\right)+3\\
                    \amp = 2\left(\frac{9}{4}\right)-\frac{3}{2}+3\\
                    \amp = \frac{9}{2}-\frac{3}{2}+\frac{6}{2}\\
                    \amp = \frac{12}{2}=6
                    \end{aligned}$



> **More Notation Ambiguity**
> As mentioned in Warning ,
        we need to remain conscious of the context of any symbol we are using.
        
        Consider the expression $a(b)$ .
        This could easily mean the output of a function $a$ with input $b$ .
        It could also mean that two numbers $a$ and $b$ need to be multiplied.
        It all depends on the context in which these symbols are being used.

## Tables and Graphs

Since functions are potentially complicated,
      we want ways to understand them more easily.
      Two basic tools for understanding a function better are tables and graphs.

**A Table for the Budget Deficit Function**


**A Table for the Square Root Function**


Another powerful tool for understanding a function better is a graph.
      Given a function $f$ ,
      one way to make its graph is to take a table of input and output values,
      and read each row as the coordinates of a point in the $xy$ -plane.

**A Graph for the Budget Deficit Function**

**A Graph for the Square Root Function**

While it is important to be able to make a graph of a function $f$ ,
      we also need to be capable of looking at a graph and reading it well.
      A graph of $f$ provides us with helpful specific information about $f$ ;
      it tells us what $f$ does to its input values.
      When we were making graphs, we plotted points of the form \[
        (\text{input},\text{output})
      \] Now given a graph of $f$ ,
      we interpret coordinates in the same way.

**Example**


**Unemployment Rates**


## Translating Between Four Descriptions of the Same Function

We have noted that functions are complicated,
      and we want ways to make them easier to understand.
      It's common to find a problem involving a function and not know how to find a solution to that problem.
      Most functions have at least four standard ways to think about them,
      and if we learn how to translate between these four perspectives,
      we often find that one of them makes a given problem easier to solve.

**Example**
Consider a function $f$ that squares its input and then adds $1$ .
          Translate this verbal description of $f$ into a table, a graph,
          and a formula.

*Solution*
Lastly, we must find a formula for $f$ .
          This means we need to write an algebraic expression that says the same thing about $f$ as the verbal description,
          the table,
          and the graph.
          For this example, we can focus on the verbal description.
          Since $f$ takes its input, squares it,
          and adds $1$ , we have that \[
            f(x)=x^2+1
          \] .

**Example**
Let $F$ be the function that takes a Celsius temperature as input and outputs the corresponding Fahrenheit temperature.
          Translate this verbal description of $F$ into a table, a graph,
          and a formula.

*Solution*
To make a table for $F$ ,
          we will need to rely on what we know about Celsius and Fahrenheit temperatures.
          It is a fact that the freezing temperature of water at sea level is 0 , which equals 32 .
          Also, the boiling temperature of water at sea level is 100 , which is the same as 212 .
          One more piece of information we might have is that standard human body temperature is 37 , or 98.6 .
          All of this is compiled in Figure .
          Note that we tabulated inputs and outputs by working with the context of the function,
          not with any computations.
To find a formula for $F$ ,
          the verbal definition is not of much direct help.
          But $F$ 's graph does seem to be a straight line.
          And linear equations are familiar to us.
          This line has a $y$ -intercept at $(0,32)$ and a slope we can calculate: $\frac{212-32}{100-0}=\frac{180}{100}=\frac{9}{5}$ .
          So the equation of this line is $y=\frac{9}{5}C+32$ .
          On the other hand, the equation of this graph is $y=F(C)$ ,
          since it is a graph of the function $F$ .
          So evidently, \[
            F(C)=\frac{9}{5}C+32
          \] .
