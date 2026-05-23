> Content sourced from [ORCCA](https://pcc.edu/orcca) (Open Resources for Community College Algebra) — CC BY 4.0
# Comparison Symbols and Notation for Intervals

Here is a true fact: $8$ is larger than $3$ . That is a comparison between two *specific* numbers. We can also make comparisons using an *unspecified* number,
      like if we say that average rent for an apartment in Portland, OR is more than $1700. We are
      not saying what the average rent is, just that it's larger than $1700. In the first half of
      this section, we examine mathematical notation for making these kinds of comparisons.

In Oregon, only citizens $18$ and older can vote in statewide elections. That is saying
      something about a large group of citizens, not just those who are $18$ . It's saying that
      people who are $37$ and $62$ may vote; and people who are $12$ may not. So it's
      a statement about a large collection of numbers. In the second half of this section, we
      examine the mathematical notation for large collections of numbers like this.

## Comparison Symbols

In everyday language you can say something like "8 is larger than 3" . In
      mathematical writing, we have a shorthand notation for this: "\gt" . It's used as
      follows: \[
        8\gt3
      \] That short expression is read aloud as "8 is greater than 3" . The symbol "\gt" is called the *greater-than symbol.*


We have to be careful when negative numbers are used in a comparison. Is $-8$ greater or
      less than $-3$ ? In one sense $-8$ is larger, because if you owe someone $8$ dollars, that's "more" than owing them $3$ dollars. But the "\gt" symbol
      does not work that way. This symbol tells you which number is farther to the right on a number
      line. With that understanding, $-3$ is greater than $-8$ .

- $$ For example, your answer might look like `4 > 3 > 2 > 1 > 0` .
              If you need to enter $\pi$ , type `pi` .
  *Solution*: We can order these numbers by placing these numbers on a number line. \begin{tikzpicture}
                  \begin{axis}
                    [
                      numberline,
                      xmin=-12,
                      xmax=12,
                      xlabel={},
                      xtick={-10,-5,...,10},
                      minor xtick={-10,-9,...,10},
                      width=\orccaprintwidth,
                    ]
                    \draw[] (axis cs: -7.6,1.5) to[out=-90, in=90] node[pos=0,above] {\(-7.6\)} (axis cs: -7.6,0);
                    \draw[] (axis cs: -6,1)     to[out=-90, in=90] node[pos=0,above] {\(-6\)}   (axis cs: -6,0);
                    \draw[] (axis cs: 6,1.5)    to[out=-90, in=90] node[pos=0,above] {\(6\)}    (axis cs: 6,0);
                    \draw[] (axis cs: 8,1)      to[out=-90, in=90] node[pos=0,above] {\(8\)}    (axis cs: 8,0);
                    \draw[] (axis cs: 9.5,1.5)  to[out=-90, in=90] node[pos=0,above] {\(9.5\)}  (axis cs: 9.5,0);
                  \end{axis}
                \end{tikzpicture} And so we see the answer is $$ .
- \[
                
              \] For example, your answer might look like `4 > 3 > 2 > 1 > 0` . If you
              need to enter $\pi$ , type `pi` .
  *Solution*: We can order these numbers by placing these numbers on a number line. Knowing or
              computing their decimals helps with this: $\pi\approx3.141\ldots$ and $\frac{10}{3}\approx3.333\ldots$ . \begin{tikzpicture}
                  \begin{axis}
                    [
                      numberline,
                      xmin=-12,
                      xmax=12,
                      xlabel={},
                      xtick={-10,-5,...,10},
                      minor xtick={-10,-9,...,10},
                      width=\orccaprintwidth,
                    ]
                    \draw[] (axis cs: -5.2,1.5)  to[out=-90, in=90] node[pos=0,above] {\(-5.2\)}                    (axis cs: -5.2,0);
                    \draw[] (axis cs: 8,1.5)     to[out=-90, in=90] node[pos=0,above] {\(8\)}                       (axis cs: 8,0);
                    \draw[] (axis cs: 3.333,1.5) to[out=-90, in=90] node[pos=0,above] {\(\frac{10}{3}\approx3.33\)} (axis cs: 3.333,0);
                    \draw[] (axis cs: 5.666,1)   to[out=-90, in=90] node[pos=0,above] {\(4.6\)}                       (axis cs: 4.6,0);
                    \draw[] (axis cs: -1,1)      to[out=-90, in=90] node[pos=0,above] {\(\pi\approx3.14\)}            (axis cs: 3.14159,0);
                  \end{axis}
                \end{tikzpicture} And so we see the answer is $$ .

The greater-than symbol has a close relative: the *greater-than-or-equal-to symbol* "\geq" . It means just like it sounds; the left number is either greater than or
      equal to the right number. Consider these examples, five of which are true and one of which is
      false. 8\amp\geq3\amp3\amp\geq-8\amp3\amp\geq3 8\amp\gt3\amp3\amp\gt-8\amp3\amp\reject{\gt}3

While it may seem unhelpful to write $3\geq3$ when you could write $3=3$ , the "\geq" symbol is useful when at least one of the numbers in a comparison is not
      specific, like in these examples: (\text{hourly pay rate})\amp\geq(\text{minimum wage})\amp(\text{age of a voter})\amp\geq18

Sometimes you want to emphasize that one number is *less than* another number. For
      this, we have symbols that are reversed from $\gt$ and $\geq$ . The symbol "\lt" is the *less-than symbol* and it's used like this: \[
        3\lt8
      \] and read aloud as "3 is less than 8" .

gives the complete list of all six comparison symbols.
      We've only discussed three of them so far in this section, but you already know the equals
      symbol. The other two are the "less than or equal to" symbol, "\leq" , and the "not equal to" symbol, "\neq" .

## Set-Builder and Interval Notation

If you write \[
        (\text{age of a voter})\geq18
      \] and have a particular voter in mind, what is that person's age? Maybe they are $18$ , but
      maybe they are older. It's helpful to use a variable $a$ to represent age (in years) and
      then to visualize the possibilities with a number line.

The shaded portion of the number line in is a mathematical *interval* . That means a collection of certain numbers with a "starting point" and a "ending point" . The interval above doesn't really ever end, but we can say $\infty$ (infinity) is the "ending point" in this situation. So this interval starts
      at $18$ and "ends" at $\infty$ .

The number line in is a *visual* representation of a
      collection of certain numbers. We have notations we can use to write down such collections of
      numbers.

> **Definition**
> Set-builder notation attempts to say directly what condition needs to be met by numbers in
          the interval. We write set-builder notation like: \[
            \left\{x\mid\text{condition on }x\right\}
          \] and read it aloud as "the set of all x such that" .

For example, $\left\{x\mid x\geq18\right\}$ is read aloud as "the set of all x
      such that x is greater than or equal to 18" . The breakdown is as follows.

| $\highlight{\{}x\mid x\geq18\highlight{\}}$ | the set of |
| --- | --- |
| $\{\highlight{x}\mid x\geq18\}$ | all $x$ |
| $\{x\highlight{{}\mid{}}x\geq18\}$ | such that |
| $\{x\mid\highlight{x\geq18}\}$ | $x$ is greater than or equal to $18$ |

**Example**

- \begin{tikzpicture}
                  \begin{axis}
                    [
                      numberline,
                      xmin=-6,
                      xmax=6,
                      xlabel={\(x\)},
                      width=\orccaprintwidth,
                    ]
                    \addplot[infiniteclosedinterval] coordinates {(-6,0) (2,0)};
                  \end{axis}
                \end{tikzpicture} If needed, type `>=` for $\geq$ , and `<=` for $\leq$ .
  *Solution*: Since all numbers less than or equal to $2$ are shaded, the set-builder notation
              is $$ .
- \begin{tikzpicture}
                  \begin{axis}
                    [
                      numberline,
                      xmin=-6,
                      xmax=6,
                      xlabel={\(x\)},
                      width=\orccaprintwidth,
                    ]
                    \addplot[infiniteopeninterval] coordinates {(-6,0) (2,0)};
                  \end{axis}
                \end{tikzpicture} If needed, type `>=` for $\geq$ , and `<=` for $\leq$ .
  *Solution*: Since all numbers less than to $2$ are shaded, the set-builder notation is $$ .
- \begin{tikzpicture}
                  \begin{axis}
                    [
                      numberline,
                      xmin=-6,
                      xmax=6,
                      xlabel={\(x\)},
                      width=\orccaprintwidth,
                    ]
                    \addplot[closedinfiniteinterval] coordinates {(2,0) (6,0)};
                  \end{axis}
                \end{tikzpicture} If needed, type `>=` for $\geq$ , and `<=` for $\leq$ .
  *Solution*: Since all numbers greater than or equal to $2$ are shaded, the set-builder
              notation is $$ .

Set-builder notation is useful, but there is an alternative for intervals that is less
      cumbersome.

> **Definition**
> Interval notation describes an interval by saying where it "starts" and "ends" .
          It can look like any of these four options: \[
            (\text{start},\text{end})\qquad
            (\text{start},\text{end}]\qquad
            [\text{start},\text{end})\qquad
            [\text{start},\text{end}]
          \]

In , the interval starts at $18$ . Then it extends forever
      and has no end, so we use the $\infty$ symbol for where this interval "ends" . And we
      write $[18,\infty)$ . There is a subtlety about using the bracket "[" on one
      side and the parenthesis ")" on the other side. The bracket tells us that $18$ *is* part of the interval and the parenthesis tells us that $\infty$ is *not* part of the interval.

Imagine if we wanted to describe all the numbers greater than $18$ , including numbers
      like $18.01$ but not including $18$ itself. Then we would write $(18,\infty)$ .

So there are four types of infinite intervals. Take note of the different uses of round
      parentheses and square brackets.

**Interval Notation from Number Lines**
- \begin{tikzpicture}
                  \begin{axis}
                    [
                      numberline,
                      xmin=-6,
                      xmax=6,
                      xlabel={\(x\)},
                      width=\orccaprintwidth,
                    ]
                    \addplot[infiniteclosedinterval] coordinates {(-6,0) (2,0)};
                  \end{axis}
                \end{tikzpicture} Type `inf` or `infinity` for $\infty$ .
  *Solution*: The shaded interval "starts" at $-\infty$ and ends at $2$ (including $2$ ) so the interval notation is $$ .
- \begin{tikzpicture}
                  \begin{axis}
                    [
                      numberline,
                      xmin=-6,
                      xmax=6,
                      xlabel={\(x\)},
                      width=\orccaprintwidth,
                    ]
                    \addplot[infiniteopeninterval] coordinates {(-6,0) (2,0)};
                  \end{axis}
                \end{tikzpicture} Type `inf` or `infinity` for $\infty$ .
  *Solution*: The shaded interval "starts" at $-\infty$ and ends at $2$ (excluding $2$ ) so the interval notation is $$ .
- \begin{tikzpicture}
                  \begin{axis}
                    [
                      numberline,
                      xmin=-6,
                      xmax=6,
                      xlabel={\(x\)},
                      width=\orccaprintwidth,
                    ]
                    \addplot[closedinfiniteinterval] coordinates {(2,0) (6,0)};
                  \end{axis}
                \end{tikzpicture} Type `inf` or `infinity` for $\infty$ .
  *Solution*: The shaded interval starts at $2$ (including $2$ ) and "ends" at $\infty$ , so the interval notation is $$ .
