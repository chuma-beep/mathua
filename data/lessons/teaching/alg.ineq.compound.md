> Content sourced from [ORCCA](https://pcc.edu/orcca) (Open Resources for Community College Algebra) — CC BY 4.0
# Compound Inequalities

On the newest version of the SAT (an exam that often qualifies students for colleges) the minimum score that you can earn is $400$ and the maximum score that you can earn is $1600$. This means that only numbers between $400$ and $1600$, including these endpoints, are possible scores. To plot all of these values on a number line would look something like:

Going back to the original statement, "the minimum score that you can earn is 400 and the maximum score that you can earn is 1600," this really says two things. First, it says that $(\text{a SAT score})\geq400$, and second, that $(\text{a SAT score})\le1600$. When we combine two inequalities like this into a single problem, it becomes a *compound inequality*.

Our lives are often constrained by the compound inequalities of reality: you need to buy enough materials to complete your project, but you can only fit so much into your vehicle; you would like to finish your degree early, but only have so much money and time to put toward your courses; you would like a vegetable garden big enough to supply you with veggies all summer long, but your yard or balcony only gets so much sun. In the rest of the section we hope to illuminate how to think mathematically about problems like these.

Before continuing, a review on how notation for intervals works may be useful, and you may benefit from revisiting Section. Then a refresher on solving linear inequalities may also benefit you, which you can revisit in Section and Section.

## Unions of Intervals

> **Definition**
> The *union* of two sets, $A$ and $B$, is the set of all elements contained in either $A$ or $B$ (or both). We write $A\cup B$ to indicate the union of the two sets. In other words, the union of two sets is what you get if you toss every number in both sets into a bigger set.

**Example**

**Example**
Visualize the union of the sets $(-\infty,4)$ and $[7,\infty)$.

*Solution*
First we make a number line with both intervals drawn to understand what both sets mean.

The two intervals should be viewed as a single object when stating the union, so here is the picture of the union. It looks the same, but now it is a graph of a single set.

> **Definition**
> The *intersection* of two sets, $A$ and $B$, is the set of all elements that are in $A$ *and* $B$. We write $A\cap B$ to indicate the intersection of the two sets. In other words, the intersection of two sets is where the two sets overlap.

**Example**

**Example**
Find the intersection of the sets $(-\infty,5)$ and $[3,\infty)$.

*Solution*
To find the intersection of the sets $(-\infty,5)$ and $[3,\infty)$, first we draw a number line with both intervals drawn to visualize where the sets overlap.

Recall that the intersection of two sets is the set of the numbers in common to both sets. In English, we might say that the lines overlap at every number between $-3$ and $5$. This description is the same as the interval $[-3,5)$.

In conclusion, 
\[[-3,\infty)\cap(-\infty,5)=[-3,5)\].

**Example**
Simplify the intersections and unions. 
1. $(-\infty,12)\cup[-3,\infty)$
2. $(-\infty,12)\cap[-3,\infty)$
3. $(-\infty,-2]\cup[4,\infty)$
4. $(-\infty,-2]\cap[4,\infty)$

*Solution*
1. $(-\infty,12)\cup[-3,\infty)=\mathbb{R}$
2. $(-\infty,12)\cap[-3,\infty)=[-3,12)$
3. $(-\infty,-2]\cup[4,\infty)=(-\infty,-2]\cup[4,\infty)$ This union cannot be simplified because the two sets have nothing in common.
4. $(-\infty,-2]\cap[4,\infty)=\emptyset$ Since the two sets have nothing in common, their intersection is empty.

## Or Compound Inequalities

> **Definition**
> A *compound inequality* is a grouping of two or more inequalities into a larger inequality statement. These usually come in two flavors: "or" and "and" inequalities. For an example of an "or" compound inequality, you might get a discount at the movie theater if your age is less than $13$ *or* greater than $64$. For an example of an "and" compound inequality, to purchase a drink at a bar in Oregon, you need to be over $21$ years old *and* be have money for your drink. You need to fulfill *both* requirements.

In math, the technical term *or* means "either or both." So, mathematically, if we asked if you would like "chocolate cake or apple pie" for dessert, your choices are either "chocolate cake," "apple pie," or "both chocolate cake and apple pie." This is slightly different than the English "or" which usually means "one or the other but not both."

"Or" shows up in math between equations (as in when solving a quadratic equation, you might end up with "x=2 or x=-3") or between inequalities (which is what we're about to discuss).

**Example**
Solve the compound inequality. 
\[x\le 1\quad\text{or}\quad x\gt 4\]

*Solution*
An "or" statement becomes a union of solution sets, so the solution set to the compound inequality must be: 
\[(-\infty,1]\cup(4,\infty)\].

**Example**
Solve the compound inequality. 
\[3-5x\gt-7\quad\text{or}\quad2-x\le-3\]

*Solution*
First we need to do some algebra to isolate $x$ in each piece. Note that we are going to do algebra on both pieces simultaneously. Also note that the mathematical symbol "or" should be written on each line. 
\[ \begin{aligned} 3-5x&\gt-7&\text{or}&& 2-x&\le-3 \\ -5x&\gt-10&\text{or}&& -x&\le-5 \\ \frac{-5x}{-5}&\mathbin{\highlight{\lt}}\frac{-10}{-5}&\text{or}&& \frac{-x}{-1}&\mathbin{\highlight{\geq}}\frac{-5}{-1} \\ x&\lt2&\text{or}&& x&\geq5 \end{aligned} \]

The solution set for the compound inequality $x\lt2$ is $(-\infty,2)$ and the solution set to $x\ge5$ is $[5,\infty)$. To do the "or" portion of the problem, we need to take the union of these two sets. Let's first make a graph of the solution sets to visualize the problem.

The union combines both solution sets into one, and so 
\[(-\infty,2)\cup[5,\infty)\]

We have finished the problem, but for the sake of completeness, let's try to verify that our answer is reasonable.

1. First, let's choose a number that is *not* in our proposed solution set. We will arbitrarily choose $\highlight{3}$. $\begin{aligned} 3-5x&\gt-7&\text{or}&& 2-x&\le-3 \\ 3-5(\highlight{3})&\wonder{\gt}-7&\text{or}&& 2-(\highlight{3})&\wonder{\le}-3 \\ -9&\reject{\gt}-7&\text{or}&& -1&\reject{\le}-3 \end{aligned}$ This value made *both* inequalities false which is why $3$ isn't in our solution set.
2. Next, let's choose a number that *is* in our solution region. We will arbitrarily choose $\highlight{1}$. $\begin{aligned} 3-5x&\gt-7&\text{or}&& 2-x&\le-3 \\ 3-5(\highlight{1})&\wonder{\gt}-7&\text{or}&& 2-(\highlight{1})&\wonder{\le}-3 \\ -12&\confirm{\lt}-7&\text{or}&& -1&\reject{\leq}-3 \end{aligned}$ This value made *one* of the inequalities true. Since this is an "or" statement, only one *or* the other piece has to be true to make the compound inequality true.
3. Last, what will happen if we choose a value that was in the other solution region in Figure, like the number $\highlight{6}$? $\begin{aligned} 3-5x&\gt-7&\text{or}&& 2-x&\le-3 \\ 3-5(\highlight{6})&\wonder{\gt}-7&\text{or}&& 2-(\highlight{6})&\wonder{\le}-3 \\ -27&\reject{\gt}-7&\text{or}&& -4&\confirm{\le}-3 \end{aligned}$ This solution made the *other* inequality piece true.

This completes the check. Numbers from within the solution region make the compound inequality true and numbers outside the solution region make the compound inequality false.

**Example**
Solve the compound inequality. 
\[\frac{3}{4}t+2\le \frac{5}{2}\quad\text{or}\quad -\frac{1}{2}(t-3)\lt -2\]

*Solution*
First we will solve each inequality for $t$. Recall that we usually try to clear denominators by multiplying both sides by the least common denominator. 
\[ \begin{aligned} \frac{3}{4}t+2&\le \frac{5}{2}&\text{or}&& -\frac{1}{2}(t-3)&\lt -2 \\ \multiplyleft{4}\left(\frac{3}{4}t+2\right)&\le \multiplyleft{4}\frac{5}{2}&\text{or}&& \multiplyleft{2}\left(-\frac{1}{2}(t-3)\right)&\lt \multiplyleft{2}(-2) \\ 3t+8&\le 10&\text{or}&& -t+3&\lt -4 \\ 3t&\le 2&\text{or}&& -t&\lt -7 \\ \frac{3t}{3}&\le \frac{2}{3}&\text{or}&& \frac{-t}{-1}&\mathbin{\highlight{\gt}} \frac{-7}{-1} \\ t&\le\frac{2}{3}&\text{or}&& t&\gt 7 \end{aligned} \]

The solution set to $t\le \frac{2}{3}$ is $\left(-\infty,\frac{2}{3}\right]$ and the solution set to $t\gt 7$ is $(7,\infty)$. Figure shows these two sets.

Note that the two sets do not overlap so there will be no way to simplify the union. Thus the solution set to the compound inequality is: 
\[\left(-\infty,\frac{2}{3}\right]\cup(7,\infty)\]

**Example**
Solve the compound inequality. 
\[3y-15\gt 6\quad\text{or}\quad7-4y\ge y-3\]

*Solution*
First we solve each inequality for $y$. 
\[ \begin{aligned} 3y-15&\gt 6&\text{or}&& 7-4y&\ge y-3 \\ 3y&\gt 21&\text{or}&& -5y&\ge -10 \\ \frac{3y}{3}&\gt \frac{21}{3}&\text{or}&& \frac{-5y}{-5}&\mathbin{\highlight{\le}} \frac{-10}{-5} \\ y&\gt 7&\text{or}&& y&\le 2 \end{aligned} \]

The solution set to $y\gt 7$ is $(7,\infty)$ and the solution set to $y\le 2$ is $(-\infty,2]$. Figure shows these two sets.

So the solution set to the compound inequality is: 
\[(-\infty,2]\cup(7,\infty)\]

## Three-Part Inequalities

There are two different kinds of "and" compound inequalities. One type has an expression that is "between" two values, like $A \lt B \le C$, that we will call "three-part inequalities". The other type has two inequalities joined by the word "and," as in $A \lt B \text{ and } C\geq D$. We will start with the three-part inequalities.

The inequality $1\leq 2\lt 3$ says a lot more than you might think. It actually says four different single inequalities which are highlighted for you to see. 
\[\highlight{1}\mathbin{\highlight{\leq}}\highlight{2}\mathbin{\lowlight{\lt}}\lowlight{3}\quad \highlight{1}\mathbin{\highlight{\leq}}\lowlight{2}\mathbin{\lowlight{\lt}}\highlight{3}\quad \highlight{1}\mathbin{\lowlight{\leq}}\lowlight{2}\mathbin{\highlight{\lt}}\highlight{3}\quad \lowlight{1}\mathbin{\lowlight{\leq}}\highlight{2}\mathbin{\highlight{\lt}}\highlight{3}\]

This might seem trivial at first, but if you are presented with an inequality like $-1\lt3\geq2$, at first it might look sensible; however, in reality, you need to check that *all four* linear inequalities make sense. Those are highlighted here. 
\[\highlight{-1}\mathbin{\highlight{\lt}}\highlight{3}\mathbin{\lowlight{\geq}}\lowlight{2}\quad \highlight{-1}\mathbin{\highlight{\lt}}\lowlight{3}\mathbin{\lowlight{\geq}}\highlight{2}\quad \highlight{-1}\mathbin{\lowlight{\lt}}\lowlight{3}\mathbin{\highlight{\geq}}\highlight{2}\quad \lowlight{-1}\mathbin{\lowlight{\lt}}\highlight{3}\mathbin{\highlight{\geq}}\highlight{2}\]

One of these inequalities is false: $-1\ngeq2$. This implies that the entire original inequality, $-1\lt3\geq2$, is nonsense.

**Example**
Decide whether or not the following inequalities are true or false. 
1. True or False: $-5\lt7\le12$?
2. True or False: $-7\le-10\lt4$?
3. True or False: $-2\le0\geq1$?
4. True or False: $5\gt-3\geq-9$?
5. True or False: $3\lt3\le5$?
6. True or False: $9\gt1\lt5$?
7. True or False: $3\lt8\le-2$?
8. True or False: $-9\lt-4\le-2$?

*Solution*
We need to go through all four single inequalities for each. If the inequality is false, for simplicity's sake, we will only highlight the one single inequality that makes the inequality false. 
1. True: $-5\lt7\le12$.
2. False: $-7\reject{\le}-10$.
3. False: $-2\reject{\ge}1$.
4. True: $5\gt-3\geq-9$.
5. False: $3\reject{\lt}3$.
6. False: $9\reject{\lt}5$.
7. False: $8\reject{\le}-2$.
8. True: $-9\lt-4\le-2$.

As a general hint, no (nontrivial) three-part inequality can ever be true if the inequality signs are not pointing in the same direction. So no matter what numbers $a$, $b$, and $c$ are, both $a\lt b\geq c$ and $a\geq b \lt c$ cannot be true! Soon you will be writing inequalities like $2\lt x \le 4$ and you need to be sure to check that your answer is feasible. You will know that if you get $2\gt x \le 4$ or $2\lt x \geq 4$ that something went wrong in the solving process. The only exception is that something like $1\le1\geq1$ is true because $1=1=1$, although this shouldn't come up very often!

**Example**
Write the solution set to the compound inequality. 
\[-7\lt x\le 5\]

*Solution*
The solutions to the three-part inequality $-7\lt x\le 5$ are those numbers that are trapped between $-7$ and $5$, including $5$ but not $-7$. Keep in mind that there are infinitely many decimal numbers and irrational numbers that satisfy this inequality like $-2.781828$ and $\pi$. We will write these numbers in interval notation as $(-7,5]$ or in set builder notation as $\{x\mid -7\lt x\le 5\}$.

**Example**
Solve the compound inequality. 
\[4\le 9x+13\lt 20\]

*Solution*

**Example**
Solve the compound inequality. 
\[-13\lt 7-\frac{4}{3}x\le 15\]

*Solution*

## Solving And Inequalities

Here we will deal with the other kind of compound inequality: the "and" variety.

**Example**
Solve the compound inequality. 
\[4-2t\gt-2\quad\text{and}\quad3t+1\geq-2\]

*Solution*
\[ \begin{aligned} 4-2t&\gt -2&\text{and}&& 3t+1&\geq -2 \\ 4-2t\subtractright{4}&\gt -2\subtractright{4} &\text{and}&& 3t+1\subtractright{1}&\geq -2\subtractright{1} \\ -2t&\gt -6&\text{and}&& 3t&\geq -3 \\ \frac{-2t}{-2}& \mathbin{\highlight{\lt}} \frac{-6}{-2} &\text{and}&& \frac{3t}{3}&\geq \frac{-3}{3} \\ t&\lt 3&\text{and}&& t&\geq -1 \end{aligned} \]

The solution set to $t\lt3$ is $(-\infty,3)$ and the solution set to $t\geq-1$ is $[-1,\infty)$. Shown is a graph of these solution sets.

Recall that an "and" problem finds the intersection of the solution sets. Intersection finds the $t$-values where the two lines overlap, so the solution to the compound inequality must be 
\[(-\infty,3)\cap[-1,\infty)=[-1,3)\].

We have finished the problem, but for the sake of completeness, let's try to "verify" that our answer is reasonable.

1. First, choose a number within our solution region and test that it makes both original inequalities true. We will arbitrarily choose $\highlight{1}$. $\begin{aligned} 4-2t&\gt -2&\text{and}&& 3t+1&\geq -2 \\ 4-2(\highlight{1})&\wonder{\gt} -2&\text{and}&& 3(\highlight{1})+1&\wonder{\geq} -2 \\ 2&\confirm{\gt} -2&\text{and}&& 4&\confirm{\geq} -2 \end{aligned}$
2. Next, choose a value outside the solution set and test that it makes *at least* one of the inequalities false. We will arbitrarily choose $\highlight{4}$. $\begin{aligned} 4-2t&\gt -2&\text{and}&& 3t+1&\geq -2 \\ 4-2(\highlight{4})&\wonder{\gt} -2&\text{and}&& 3(\highlight{4})+1&\wonder{\geq} -2 \\ -4&\reject{\gt} -2&\text{and}&& 13&\confirm{\geq} -2 \end{aligned}$ Since one of the inequalities is false and this is an "and" question, the compound inequality is false for this value which is what expected by picking a number outside the solution set.
3. Last, we should choose a number that is not a solution that is on the "other side" of the solution set. We will arbitrarily choose $\highlight{-2}$. $\begin{aligned} 4-2t&\gt -2&\text{and}&& 3t+1&\geq -2 \\ 4-2(\highlight{-2})&\wonder{\gt} -2&\text{and}&& 3(\highlight{-2})+1&\wonder{\geq} -2 \\ 8&\confirm{\gt} -2&\text{and}&& -5&\reject{\geq} -2 \end{aligned}$ Again, since one of the inequalities is false and this is an "and" question, the compound inequality is false for $-2$.

So, numbers outside the proposed solution region make the compound inequality false, and numbers inside the region make the compound inequality true. We have verified our solution set.

## Applications of Compound inequalities

**Example**
Raphael's friend is getting married and he's decided to give them some dishes from their registry. Raphael doesn't want to seem cheap but isn't a wealthy man either, so he wants to buy "enough" but not "too many." He's decided that he definitely wants to spend at least $\$150$ on his friend, but less than $\$250$. Each dish is $\$21.70$ and shipping on an order of any size is going to be $\$19.99$. Given his budget, set up and algebraically solve a compound inequality to find out what his different options are for the number of dishes that he can buy.

*Solution*
First, we should define our variable. Let $x$ represent the number of dishes that Raphael can afford. Next we should write a compound inequality that describes this situation. In this case, Raphael wants to spend between $\$150$ and $\$250$ and, since he's buying $x$ dishes, the price that he will pay is $21.70x+19.99$. All of this translates to a triple inequality 
\[150 \lt 21.70x+19.99 \lt 250\]

Now we have to solve this inequality in the usual way. 
\[ \begin{aligned} 150 &\lt 21.70x+19.99 \lt 250 \\ 150\subtractright{19.99} &\lt 21.70x+19.99\subtractright{19.99} \lt 250\subtractright{19.99} \\ 130.01 &\lt 21.70x \lt 230.01 \\ \frac{130.01}{21.70} &\lt \frac{21.70x}{21.70} \lt \frac{230.01}{21.70} \\ 5.991 &\lt x \lt 10.6&\text{(note: these values are approximate)} \end{aligned} \]

The interpretation of this inequality is a little tricky. Remember that $x$ represents the number of dishes Raphael can afford. Since you cannot buy $5.991$ dishes (manufacturers will typically only ship whole number amounts of tableware) his minimum purchase must be $6$ dishes. We have a similar problem with his maximum purchase: clearly he cannot buy $10.6$ dishes. So, should we round up or down? If we rounded up, that would be $11$ dishes and that would cost $\$21.70\cdot\highlight{11}+\$19.99=\$258.69$, which is outside his price range. Therefore, we should actually round *down* in this case.

In conclusion, Raphael should buy somewhere between 6 and 10 dishes for his friend to stay within his budget.

**Example**
Oak Ridge National Laboratory, a renowned scientific research facility, compiled some https://tedb.ornl.gov/data/ on fuel efficiency of a mid-size hybrid car versus the speed that the car was driven. A model for the fuel efficiency $e(x)$ (in miles per gallon, $\text{milepergallon}$) at a speed $x$ (in miles per hour, $\text{mileperhour}$) is $e(x)=88-0.7x$. 
1. Evaluate and interpret $e(60)$ in the context of the problem.
2. Note that this model only applies between certain speeds. The maximum fuel efficiency for which this formula applies is $55\,\text{milepergallon}$ and the minimum fuel efficiency for which it applies is $33\,\text{milepergallon}$. Set up and algebraically solve a compound inequality to find the range of speeds for which this model applies.

*Solution*
1. Let's evaluate $e(60)$ first. $\begin{aligned} e(x)&=88-0.7x \\ e(\highlight{60})&=88-0.7(\highlight{60}) \\ &=46 \end{aligned}$ So, when the hybrid car travels at a speed of $60\,\text{mileperhour}$, it has a fuel efficiency of $46\,\text{milepergallon}$.
2. In this case, the minimum efficiency is $33\,\text{milepergallon}$ and the maximum efficiency is $55\,\text{milepergallon}$. We need to trap our formula between these two values to solve for the respective speeds. $\begin{aligned} 33 &\lt 88-0.7x \lt 55 \\ 33\subtractright{88} &\lt 88-0.7x\subtractright{88} \lt 55\subtractright{88} \\ -55 &\lt -0.7x \lt -33 \\ \frac{-55}{-0.7} &\mathbin{\highlight{\gt}} \frac{-0.7x}{-0.7} \mathbin{\highlight{\gt}} \frac{-33}{-0.7} \\ 78.57 &\gt x \gt 47.14&\text{(note: these values are approximate)} \end{aligned}$

This inequality says that our model is applicable when the car's speed is between about $47\,\text{mileperhour}$ and about $79\,\text{mileperhour}$.

