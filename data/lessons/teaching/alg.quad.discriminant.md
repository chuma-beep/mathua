> Content sourced from [ORCCA](https://pcc.edu/orcca) (Open Resources for Community College Algebra) — CC BY 4.0
# The Quadratic Formula

In, we solved certain quadratic equations using the square root property. This worked when the equations were in the right format where it was possible to use a square root radical. Now we learn another method to solve quadratic equations in general: the quadratic formula.

## Solving Quadratic Equations with the Quadratic Formula

The standard form for a quadratic equation is 
\[ax^2+bx+c=0\]
 where $a$ is some nonzero number.

When $b=0$ and the equation's form is $ax^2+c=0$, then we can simply use the square root property to solve it as we did with equations like this in. But can we solve equations where $b\neq0$? A general method for doing this is to use the *quadratic formula*.

**The Quadratic Formula**

As we have seen from solving quadratic equations in, there can be two solutions to a quadratic equation. Both solutions are represented in the quadratic formula because of the $\pm$ symbol. We could write the two solutions separately: 
\[ \begin{aligned} x=\frac{-b-\sqrt{b^2-4ac}}{2a}&&\text{ or }&& x=\frac{-b+\sqrt{b^2-4ac}}{2a} \end{aligned} \]

This method for solving quadratic equations will work to solve *every* quadratic equation. It is most helpful when $b\ne0$.

**Example**
Linh is in a physics class that launches a tennis ball from a rooftop that is $90.2$ feet above the ground. They fire it directly upward at a speed of $14.4$ feet per second and measure the time it takes for the ball to hit the ground below. We can model the height of the tennis ball, $h$, in feet, with the quadratic equation $h=-16t^2+14.4t+90.2$, where $x$ represents the time in seconds after the launch. According to the model, when should the ball hit the ground? The ground has height $0$ feet, so we should substitute $0$ in for $h$. This gives us the quadratic equation: 
\[0=-16t^2+14.4t+90.2\]
 We cannot solve this equation with the square root property, so we will use the quadratic formula. First we will identify $\highlight{a=-16}$, $\highlight{b=14.4}$, and $\highlight{c=90.2}$, and substitute these into the quadratic formula: 
\[ \begin{aligned} t&=\frac{-b\pm\sqrt{b^2-4ac}}{2a} \\ t&=\frac{-(14.4)\pm\sqrt{(14.4)^2-4(-16)(90.2)}}{2(-16)} \\ t&=\frac{-14.4\pm\sqrt{5980.16}}{-32} \end{aligned} \]
 These are the exact solutions, but because we have a physical context we want to approximate the solutions with decimals. 
\[t\approx-1.966\qquad\text{or}\qquad t\approx2.866\]
 We don't use the negative solution because a negative time does not make sense in this context. The ball will hit the ground approximately $2.866$ seconds after it is launched.

The quadratic formula can be used to solve any quadratic equation, but it requires that you remembering the formula correctly and that you correctly identify $a$, $b$, and $c$. Also, that you don't make any arithmetic mistakes when you simplify. We recommend that you always check if you could use the square root property before using the quadratic formula.

**Example**
Solve for $x$ in $2x^2-9x+5=0$.

*Solution*
First, we check and see that we cannot use the square root property (because $b\neq0$) so we will use the quadratic formula. We identify $a=2$, $b=-9$, and $c=5$. Substituting these into the quadratic formula: 
\[ \begin{aligned} x&=\frac{-b\pm\sqrt{b^2-4ac}}{2a} \\ x&=\frac{-(-9)\pm\sqrt{(-9)^2-4(2)(5)}}{2(2)} \\ x&=\frac{9\pm\sqrt{81-40}}{4} \\ x&=\frac{9\pm\sqrt{41}}{4} \end{aligned} \]
 This is fully simplified because we cannot simplify $\sqrt{41}$ or reduce the fraction. The solution set is $\left\{\frac{9-\sqrt{41}}{4},\frac{9+\sqrt{41}}{4}\right\}$. In this exercise, there is no reason to approximate these with decimals.

When a quadratic equation does not start out in standard form we must convert it to standard form before we can clearly identify the values of $a$, $b$, and $c$.

**Example**
Solve for $x$ in $x^2=-10x-3$.

*Solution*
First, we convert the equation into standard form by adding $10x$ and $3$ to each side of the equation: 
\[x^2+10x+3=0\]
 Next, we check that we cannot use the square root property so we will use the quadratic formula. We identify $a=1$, $b=10$ and $c=3$. Substituting them into the quadratic formula: 
\[ \begin{aligned} x&=\frac{-b\pm\sqrt{b^2-4ac}}{2a} \\ x&=\frac{-10\pm\sqrt{(10)^2-4(1)(3)}}{2(1)} \\ x&=\frac{-10\pm\sqrt{100-12}}{2} \\ x&=\frac{-10\pm\sqrt{88}}{2} \end{aligned} \]
 We notice that the radical can be simplified: 
\[ \begin{aligned} x&=\frac{-10\pm2\sqrt{22}}{2} \\ x&=\frac{-10}{2}\pm\frac{2\sqrt{22}}{2} \\ x&=-5\pm\sqrt{22} \end{aligned} \]
 The solution set is $\{-5-\sqrt{22}, -5+\sqrt{22}\}$.

The radicand from the quadratic formula, $b^2-4ac$, is called the *discriminant*. When it is a negative number, the quadratic equation has no real solution.

**Example**
Solve for $y$ in $y^2-4y+8=0$.

*Solution*
Identify $a=1$, $b=-4$ and $c=8$. Substitute them into the quadratic formula: 
\[ \begin{aligned} y&=\frac{-b\pm\sqrt{b^2-4ac}}{2a} \\ &=\frac{-(-4)\pm\sqrt{(-4)^2-4(1)(8)}}{2(1)} \\ &=\frac{4\pm\sqrt{16-32}}{2} \\ &=\frac{4\pm\sqrt{-16}}{2} \end{aligned} \]
 The discriminant worked out to be $-16$, which is negative. The square root of a negative number is not a real number, so we will conclude that this equation has no real solutions.

## Applications

Certain "word problems" lead to a quadratic equation where the quadratic formula may be useful.

**Example**
A rectangle is $5$ inches longer than it is wide. The total area of the rectangle is 60 square inches. How wide is the rectangle? (This is asking for the shorter dimension.)

*Solution*
The area of a rectangle is given by multipying its two dimensions. If we let $w$ stand for the width of the rectangle (the shorter dimension), then $w+5$ is the length, and the total area is $w(w+5)$. A diagram can help make more clear what needs to be done.

So we must solve the equation $w(w+5)=60$. After multiplying and moving to standard form, we have $w^2+5w-60=0$. Using the quadratic formula: 
\[ \begin{aligned} w&=\frac{-b\pm\sqrt{b^2-4ac}}{2a} \\ &=\frac{-5\pm\sqrt{5^2-4(1)(-60)}}{2(1)} \\ &=\frac{-5\pm\sqrt{265}}{2} \\ w&\approx-10.64\quad\text{or}\quad w\approx5.64 \end{aligned} \]
 Only the solution $w\approx5.64$ makes sense as the width of a rectangle. So the rectangle is about $5.64$ inches wide.

**Example**
Amita has a food stand where she sells momo (Nepali dumplings). If she charges $p$ dollars for each momo, she estimates that she'd sell $80{,}000-20{,}000p$ over the course of a year. (This reflects how raising the price can lead to fewer people making purchases.) What is the largest price that Amita could set to end up with a revenue of $\$70{,}000$ for the year? (Note that revenue is not the same as profit, and we will not be accounting for Amita's expenses.)

*Solution*
If Amita sets the price at $p$ dollars per momo, her total revenue will be $(\text{number sold})\cdot\text{price}$, which is $(80000-20000p)p$, which simplifies to $80000p-20000p^2$. We have been asked to set that equal to $70000$. 
\[ \begin{aligned} 80000p-20000p^2 & = 70000 \\ -20000p^2+80000p-70000 & = 0 \\ -2p^2+8p-7 & = 0 \end{aligned} \]
 
\[ \begin{aligned} p&=\frac{-b\pm\sqrt{b^2-4ac}}{2a} \\ &=\frac{-8\pm\sqrt{8^2-4(-2)(-7)}}{2(-2)} \\ &=\frac{-8\pm\sqrt{8}}{-4} \\ p\approx& 2.71 \quad\text{or}\quad x\approx1.29 \end{aligned} \]
 The possible solutions are about $2.71$ and $1.29$. They are both valid solutions. Amita can reach the target revenue of $\$70{,}000$ by setting the price per momo to either $\$1.29$ or $\$2.71$.

## Radical Equations

Sometimes a radical equation gives rise to a quadratic equation, and the quadratic formula can then be useful.

**Example**
Solve for $z$ in $\sqrt{z}+2=z$.

*Solution*
We will isolate the radical first, and then square both sides. 
\[ \begin{aligned} \sqrt{z}+2&=z \\ \sqrt{z}&=z-2 \\ \left(\sqrt{z}\right)^{\highlight{2}}&=(z-2)^{\highlight{2}} \\ z&=z^2-4z+4 \\ 0&=z^2-5z+4 \\ z&=\frac{5\pm\sqrt{(-5)^2-4(1)(4)}}{2} \\ &=\frac{5\pm\sqrt{25-16}}{2} \\ &=\frac{5\pm\sqrt{9}}{2} \\ &=\frac{5\pm3}{2} \end{aligned} \]
 
\[ \begin{aligned} z&=\frac{5-3}{2}&\text{ or }&& z&=\frac{5+3}{2} \\ z&=1&\text{ or }&& z&=4 \end{aligned} \]
 There was a step in this process where we squared both sides of an equation. Squaring both sides of an equation can cause false equations to become true. For example $-2=2$ is false, but if you square both sides then $4=4$ is true. This means that when an equation solving process involves squaring both sides, and we end up with "solutions" like the two $z$-values above, they are only *potential* solutions. Logically we've whittled down the list of possible solutions to just these two numbers. But we need to actually check whether or not they actually work. 
\[ \begin{aligned} \sqrt{1}+2&\wonder{=}1&\sqrt{4}+2&\wonder{=}4 \\ 1+2&\reject{=}1&2+2&\confirm{=}4 \end{aligned} \]
 It turned out that $1$ is an *extraneous solution*, but $4$ is a valid solution. So the equation has one solution, $4$, and the solution set is $\{4\}$.

**Example**
Solve the equation $\sqrt{2n-6}=1+\sqrt{n-2}$ for $n$.

*Solution*
We cannot isolate two radicals at the same time. One of the radicals is already isolated, so we will simply square both sides and then try to clean up what remains. 
\[ \begin{aligned} \sqrt{2n-6}&=1+\sqrt{n-2} \\ \left(\sqrt{2n-6}\right)^{\highlight{2}}&=\left(1+\sqrt{n-2}\right)^{\highlight{2}} \\ 2n-6&=1^2+2\sqrt{n-2}+\left(\sqrt{n-2}\right)^2 \\ 2n-6&=1+2\sqrt{n-2}+n-2 \\ 2n-6&=2\sqrt{n-2}+n-1 \\ n-5&=2\sqrt{n-2} \\ (n-5)^{\highlight{2}}&=\left(2\sqrt{n-2}\right)^{\highlight{2}} \\ n^2-10n+25&=4(n-2) \\ n^2-10n+25&=4n-8 \\ n^2-14n+33&=0 \end{aligned} \]
 
\[ \begin{aligned} n&=\frac{14\pm\sqrt{14^2-4(1)(33)}}{2} \\ &=\frac{14\pm\sqrt{196-132}}{2} \\ &=\frac{14\pm\sqrt{64}}{2} \\ &=\frac{14\pm8}{2} \end{aligned} \]
 
\[ \begin{aligned} n&=\frac{14-8}{2}&\text{ or }&& n&=\frac{14+8}{2} \\ n&=3&\text{ or }&& n&=11 \end{aligned} \]
 So our two potential solutions are $3$ and $11$. We should now verify that they truly are solutions. 
\[ \begin{aligned} \sqrt{2(3)-6}&\wonder{=}1+\sqrt{3-2}&\sqrt{2(11)-6}&\wonder{=}1+\sqrt{11-2} \\ \sqrt{6-6}&\wonder{=}1+\sqrt{1}&\sqrt{22-6}&\wonder{=}1+\sqrt{9} \\ \sqrt{0}&\wonder{=}1+1&\sqrt{16}&\wonder{=}1+3 \\ 0&\reject{=}2&4&\confirm{=}4 \end{aligned} \]
 So, $11$ is the only solution. The solution set is $\{11\}$.

