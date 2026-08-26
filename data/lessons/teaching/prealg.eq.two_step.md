> Content sourced from [ORCCA](https://pcc.edu/orcca) (Open Resources for Community College Algebra) — CC BY 4.0

# Solving Multistep Linear Equations

We solved equations in  the section  where only one step was needed to isolate the variable. Now we will work with equations that need more than one step.

*Alternative Video Lessons*

## Solving Two-Step Equations

**Example**

A water tank can hold up to $140$ gallons of water, but it starts with only $5$ gallons. A tap is turned on, pouring $15$ gallons of water into the tank every minute. After how many minutes will the tank be full?

You might recognize that this is a "rate" scenario like we explored in  the section . If you didn't notice that, you can still explore the given information with a table. You may find a pattern that helps you understand the question better. (And even if you do not find a pattern, spending some time thinking about these quantities stills helps to understand the question better.)

Each additional minute of time gives us $15$ more gallons of water. So after $t$ minutes, we've added "$15$ times $t$" gallons of water to the $5$ gallons that we started with. So after $t$ minutes, we have $15t+5$ gallons. To find when the tank will be full with $140$ gallons, we can write the equation

$$
 15t+5=140 
$$

This is the same equation we would get if we apply  the rate model formula .

To solve this, first we isolate the variable term, $15t$. We need to "remove" the $5$ from the left side of the equation. We can do this in a legal way using the  addition property of equality  by subtracting $5$ from each side of the equation. Once the variable term is isolated, we can eliminate its coefficient and solve for $t$.

The full process is:

$$
\begin{aligned}15t+5&=140 \\ 15t+5\subtractright{5}&=140\subtractright{5} \\ 15t&=135 \\ \divideunder{15t}{15}&=\divideunder{135}{15} \\ t&=9\end{aligned}
$$

We should check this possible solution by substituting $9$ in for $t$ in the original equation:

$$
\begin{aligned}15t+5&=140 \\ 15(\substitute{9})+5&\wonder{=}140 \\ 135+5&\confirm{=}140\end{aligned}
$$

And the solution $9$ is verified.

This problem had *context*. It was not simply solving an equation. It came with a story about a tank filling with water. So we should report a conclusion that uses that context. Something like "The tank will be full after $9$ minutes."

In solving the two-step equation in  the example , we first isolated the variable expression $15t$ and then eliminated the coefficient $15$ by dividing each side of the equation by $15$. These two steps are the heart of our approach to solving linear equations. Try these two steps in the following exercise.

**Exercise**

In  the section , there was  the example . In that example, some background information let us set up an equation, but we didn't try to solve it. Now we can try solving it.

**Exercise**

## Solving Multistep Linear Equations

More complicated equations might need a few setup steps before we can do the two important steps of isolating the variable term and eliminating the coefficient. Here is a general guide for what the full process can be like.

Steps to Solve Linear Equations



**Example**

Ahmed has $$\$2500$$ in his savings account and is going to start saving $$\$550$$ per month. Julia has $$\$4600$$ in her savings account and is going to start saving $$\$250$$ per month. If this situation continues, how long will it take for Ahmed to catch up with Julia in savings?

Ahmed saves $$\$550$$ per month, so he can save $550t$ dollars in $t$ months. With the $$\$2500$$ he started with, after $t$ months he has $550t + 2500$ dollars. Similarly, after $t$ months, Julia has $250t + 4600$ dollars. To find when those two accounts will have the same amount of money, we write the equation

$$
 550t + 2500 = 250t + 4600 
$$

Each side is simplified, but unlike earlier examples, we have the variable $t$ on *both* side of the equation. But we can still use  properties of equality  to have only one $t$-term. We can start by subtracting $250t$ from each side.

$$
\begin{aligned}550t+2500&=250t+4600 \\ 550t+2500\subtractright{250t}&=250t+4600\subtractright{250t} \\ 300t+2500&=4600 \\ 300t+2500\subtractright{2500}&=4600\subtractright{2500} \\ 300t&=2100 \\ \divideunder{300t}{300}&=\divideunder{2100}{300} \\ t&=7\end{aligned}
$$

Checking the solution $7$:

$$
\begin{aligned}550t+2500&=250t+4600 \\ 550(\substitute{7})+2500&\wonder{=}250(\substitute{7})+4600 \\ 3850+2500&\wonder{=}1750+4600 \\ 6350&\confirm{=}6350\end{aligned}
$$

Ahmed will catch up to Julia after $7$ months.

**Exercise**

In  the exercise , we could have moved variable terms to the *right* side of the equal sign and number terms to the *left* side. We chose not to, but there's no reason why we couldn't have done that. Let's explore:

$$
\begin{aligned}5-2x&=5x-9 \\ 5-2x\addright{2x}&=5x-9\addright{2x} \\ 5&=7x-9 \\ 5\addright{9}&=7x-9\addright{9} \\ 14&=7x \\ \divideunder{14}{7}&=\divideunder{7x}{7} \\ 2&=x\end{aligned}
$$

The solution is the same either way.

Also, we could save a step by moving variable terms and constant terms in one step:

$$
\begin{aligned}5-2x&=5x-9 \\ 5-2x\addright{2x+9}&=5x-9\addright{2x+9} \\ 14&=7x \\ \divideunder{14}{7}&=\divideunder{7x}{7} \\ 2&=x\end{aligned}
$$

For the sake of a slow and careful explanation, the examples in this chapter will move variable terms and number terms in separate steps.

The next example requires combining like terms.

**Example**

Solve for $n$ in $n-9+3n=n-3n$.

We start by combining like terms. After this, we can separate the $n$-terms and constant terms, proceeding as before.

$$
\begin{aligned}n-9+3n&=n-3n \\ 4n-9&=-2n \\ 4n-9\subtractright{4n}&=-2n\subtractright{4n} \\ -9&=-6n \\ \divideunder{-9}{-6}&=\divideunder{-6n}{-6} \\ \frac{3}{2}&=n\end{aligned}
$$

Checking the solution $\frac{3}{2}$:

$$
\begin{aligned}n-9+3n&=n-3n \\ \substitute{\frac{3}{2}}-9+3\left(\substitute{\frac{3}{2}}\right)&\wonder{=}\substitute{\frac{3}{2}}-3\left(\substitute{\frac{3}{2}}\right) \\ \frac{3}{2}-9+\frac{9}{2}&\wonder{=}\frac{3}{2}-\frac{9}{2} \\ \frac{12}{2}-9&\wonder{=}-\frac{6}{2} \\ 6-9&\confirm{=}-3\end{aligned}
$$

The solution to the equation $n-9+3n=n-3n$ is $\frac{3}{2}$ and the solution set is $\left\{\frac{3}{2}\right\}$.

**Example**

Solve for $a$ in $4 - (3 - a)=-2 - 2(2a + 1)$.

This time we start by simplifying each side of the equation, which involves distributing as well as combining like terms. Here is a reminder to be careful when distributing a negative sign over a group of terms.

$$
\begin{aligned}4-(3-a)&=-2-2(2a+1) \\ 4-3+a&=-2-4a-2 \\ 1+a&=-4-4a \\ 1+a\addright{4a}&=-4-4a\addright{4a} \\ 1+5a&=-4 \\ 1+5a\subtractright{1}&=-4\subtractright{1} \\ 5a&=-5 \\ \divideunder{5a}{5}&=\divideunder{-5}{5} \\ a&=-1\end{aligned}
$$

Checking the solution $-1$:

$$
\begin{aligned}4-(3-a)&=-2-2(2a+1) \\ 4-(3-(\substitute{-1}))&\wonder{=}-2-2(2(\substitute{-1})+1) \\ 4-(4)&\wonder{=}-2-2(-1) \\ 0&\wonder{=}-2+2 \\ 0&\confirm{=}0\end{aligned}
$$

Therefore the solution to the equation is $-1$ and the solution set is $\{-1\}$.

**Exercise**

## Revisiting Applications

In  the section , we explored several "word problem" scenarios that led to equations, but we did not try to solve those equations. Let's revisit some of those applications and try to solve them.

Here we revisit  the example .

**Example**

A bathtub contains $2.5\,\text{ft}^{3}$ of water. More water is being poured in at a rate of $1.75\,\text{ft}^{3}$ per minute. How long will it be until the amount of water in the bathtub reaches $6.25\,\text{ft}^{3}$?

We have an initial amount of water ($2.5\,\text{ft}^{3}$), a rate at which the amount of water is changing ($1.75\,\text{ft}^{3}$), and a final amount of water we are going to reach ($6.25\,\text{ft}^{3}$). So we can use the pattern for rate modeling from   .

We should clearly identify the variable first though. The solution is supposed to represent an amount of time. So a reasonable variable to use is $t$. Let $t$ be the amount of time, in minutes, that it takes for the tub to reach $6.25\,\text{ft}^{3}$. And we have the equation:

$$
 2.5 + 1.75t = 6.25 
$$

Now we have the skills to solve this equation.

$$
\begin{aligned}2.5 + 1.75t &= 6.25 \\ 2.5 + 1.75t\subtractright{2.5} &= 6.25\subtractright{2.5} \\ 1.75t&= 3.75 \\ \divideunder{1.75t}{1.75}&=\divideunder{3.75}{1.75} \\ t&\approx2.14\end{aligned}
$$

So it will take about $2.14$ minutes for the tub to have $6.25\,\text{ft}^{3}$ of water.

Here we revisit  the example .

**Example**

Jakobi's annual salary as a nurse this year is $$\$73{,}290$$. That's following a $4\%$ raise over last year's salary. What was his salary the previous year?

As soon as you understand that the solution will be Jakobi's salary from last year, that is when you should clearly define a variable. Since it will represent a salary, we choose to use $S$ as the variable. Let $S$ represent Jakobi's salary from last year.

To set up the equation, we need to think about how he arrived at this year's salary. His employer took last year's salary and added $4\%$ to that. In words, we have:

$$
 (\text{last year's salary})+(4\%\text{ of last year's salary}) = (\text{this year's salary}) 
$$

We represent "$4\%$ of last year's salary" with $0.04S$ since $0.04$ is the decimal equivalent to $4\%$. So out equation is:

$$
 S + 0.04S = 73290 
$$

Now we have the skills to solve this equation.

$$
\begin{aligned}S + 0.04S &= 73290 \\ 1.04S &= 73290 \\ \divideunder{1.04S}{1.04}&=\divideunder{73290}{1.04} \\ S&\approx70471\end{aligned}
$$

So last year, Jakobi's salary was about $$\$70{,}471$$.

Here we revisit  the exercise .

**Exercise**

## Differentiating between Simplifying Expressions, Evaluating Expressions and Solving Equations

Consider the following three examples, which have similarities on the surface, but are fundamentally different from each other.

**Example**

Simplify the expression $10-3(x+2)$.

Distribute the $-3$ and then combine like terms.

$$
\begin{aligned}10-3(x+2)&=10-3x-6 \\ &=-3x+4\end{aligned}
$$

Note that our final result is an *expression*.

**Example**

Evaluate the expression $10-3(x+2)$ when $x=2$.

Substitute $2$ in for $x$ in the expression:

$$
\begin{aligned}10-3(x+2)&=10-3(\substitute{2}+2) \\ &=10-3(4) \\ &=10-12 \\ &=-2\end{aligned}
$$

Note that our final result here is a *number*.

**Example**

Solve the equation $10-3(x+2)=x-16$.

Follow the process we have been using in this section to solve a linear equation.

$$
\begin{aligned}10-3(x+2)&=x-16 \\ 10-3x-6&=x-16 \\ -3x+4&=x-16 \\ -3x+4\subtractright{x}&=x-16\subtractright{x} \\ -4x+4&=-16 \\ -4x+4\subtractright{4}&=-16\subtractright{4} \\ -4x&=-20 \\ \divideunder{-4x}{-4}&=\divideunder{-20}{-4} \\ x&=5\end{aligned}
$$

So the solution set is $\{5\}$.

Note that our final result here is a *solution set*.

*Simplifying*, *evaluating*, and *solving* are three different algebra tasks. Students often use these vocabulary terms incorrectly, using one when they meant another. Here is a summary collection of the differences that you should understand between these algebra tasks.

**Exercise**

Matching Vocabulary

Match each piece of mathematics with the task that is most likely to be something you would do with it.

##

**Exercise**

Describe the five steps you might need to take when solving a linear equation.

**Exercise**

In this section there is a reminder to take care with negative numbers when doing what?

**Exercise**

Explain what is wrong with saying "I need to solve $3x+x-8$."

##

