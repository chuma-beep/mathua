> Content sourced from [Algebrica](https://algebrica.org/homogeneous-trigonometric-equations/) — CC BY-NC 4.0

## What are homogeneous trigonometric equations

Homogeneous trigonometric equations are [equations](<../equations>) in which all terms involve trigonometric functions, such as [sine and cosine](<../sine-and-cosine>), [raised](<../powers>) to the same degree. A first-degree homogeneous trigonometric equation is generally represented by the form:

\\[a \sin x + b \cos x = 0 \\]

Each term in the equation is of degree \\(1\\) and \\( a \\) and \\( b \\) are [real coefficients](<../types-of-numbers>). A general form of a first-degree trigonometric equation is:

\\[a \sin x + b \cos x + c = 0 \\]

This equation is considered homogeneous when the constant term \\( c = 0 \\). Similar considerations apply to second-degree homogeneous trigonometric equations, which include only terms of second degree. These equations are generally written in the form:

\\[a \sin^2 x + b \sin x \cos x + c \cos^2 x = 0 \\]

##### The definition can be generalized to equations of the third, fourth, fifth degree, and so on. In these cases, all terms must involve trigonometric [functions](<../functions>) raised to the same degree, ensuring the equation remains homogeneous.

Homogeneous equations of degree \\( n \\) can be solved by dividing each term of the equation by \\( \cos^n x \\), thereby transforming the equation into one expressed in terms of \\( \tan x \\).

## Example 1

Let’s solve the first-degree homogeneous trigonometric equation: \\[\sin x - \sqrt{3} \cos x = 0 \\]


Dividing the equation by \\( \cos x \\), we obtain:

\\[\frac{\sin x}{\cos x} - \sqrt{3} = 0\\]

Since the ratio of sine to cosine equals the [tangent](<../tangent-and-cotangent>), we can rewrite the equation in terms of \\( \tan x. \\)

\\[\tan x = \sqrt{3}\\]

Solving the equation we get the value:

\\[x = \arctan(\sqrt{3}) = \frac{\pi}{3} \\]

Since the tangent function is periodic with period \\( \pi \\), the general solution is:

\\[x = \frac{\pi}{3} + k\pi, \quad k \in \mathbb{Z} \\]

Whenever [trigonometric identities](<../trigonometric-identities/>) allow us to rewrite a trigonometric equation in homogeneous form, it is generally preferable to do so, as it can simplify the calculations and make the solution process more manageable.

## Example 2

Let us now consider the second-degree trigonometric equation:

\\[\sin^2 x + \sin 2x - \cos^2 x = 0 \\]

At first glance, this is not a homogeneous equation, because not all terms are of the same degree. In particular, \\( \sin 2x \\) is a first-degree term, while \\( \sin^2 x \\) and \\( \cos^2 x \\) involve second-degree expressions. In an equation like this, we cannot proceed by dividing all terms by \\( \cos^2 x \\), because the first-degree term does not transform into a useful expression for solving the equation. In fact, this method works only when all terms are of the same degree, as in homogeneous equations.


However, using [trigonometric identities](<../trigonometric-identities/>), specifically the double angle formula, we can rewrite \\( \sin(2x) \\) as \\( 2\sin(x)\cos(x) \\). The equation then becomes:

\\[\sin^2 x + 2\sin x \cos x - \cos^2 x = 0 \\]

In this way, we have transformed the original equation into a second-degree homogeneous trigonometric equation, which can now be solved by dividing all terms by \\( \cos^2 x \\). The equation becomes:

\\[\frac{\sin^2 x}{\cos^2 x} + \frac{2\sin x \cos x}{\cos^2 x} - \frac{\cos^2 x}{\cos^2 x} = 0 \\]


Simplifying each term, we obtain:

\\[\tan^2 x + 2\tan x - 1 = 0 \\]

This is a [quadratic equation](<../quadratic-equations>) in \\( \tan x \\), and can now be solved using the [quadratic formula](<../quadratic-formula>). Let’s substitute \\( t = \tan x \\). The equation becomes:

\\[t^2 + 2t - 1 = 0 \\]

Using the quadratic formula, we have:

\\[t = \frac{-2 \pm \sqrt{4 + 4}}{2} = \frac{-2 \pm 2\sqrt{2}}{2} = -1 \pm \sqrt{2} \\]

The two solutions are:

\\[t_1 = -1 + \sqrt{2}, \quad t_2 = -1 - \sqrt{2} \\]


Substituting back \\( t = \tan x \\), we obtain the solutions:

\\[\tan x = -1 + \sqrt{2} \quad \text{and} \quad \tan x = -1 - \sqrt{2} \\]

The general solutions are:

\\[x_1 = \arctan(-1 + \sqrt{2}) + k\pi \quad k \in \mathbb{Z} \\] \\[x_2 = \arctan(-1 - \sqrt{2}) + k\pi \quad k \in \mathbb{Z} \\]
