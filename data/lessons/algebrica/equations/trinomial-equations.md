> Content sourced from [Algebrica](https://algebrica.org/trinomial-equations/) — CC BY-NC 4.0

## What are trinomial equations

Trinomial equations are a specific type of [polynomial equation](<../polynomial-equations/>) that consist of three terms involving constants and [powers](<../powers>) of a variable.

\\[ax^{2n} + bx^{n} + c = 0 \\]

This form is especially useful because it allows for substitution strategies such as setting \\(x^n = y\\) to transform the equation into a standard quadratic.

  * \\( a \\), \\( b \\), and \\( c \\) are numerical coefficients.
  * \\( x \\) is the unknown variable.


It is worth noting that if \\( n = 1 \\), the equation becomes a [quadratic equation](<../quadratic-equations>), often written in the form:

\\[ax^2 + bx + c = 0 \\]

## How to solve trinomial equations

A general set of steps can be followed to solve a trinomial equation of the form \\( ax^{2n} + bx^n + c = 0 \\). The first step is to make a substitution that transforms the equation into a quadratic equation in \\( y \\), which is easier to solve. For example, we let:

\\[x^n = y \\]

and obtain:

\\[ay^2 + by + c = 0 \\]

We can use techniques such as [factoring](<../factoring-quadratic-equations>), [completing the square](<../completing-the-square>), or the [quadratic formula](<../quadratic-formula>) to solve the resulting quadratic equation in \\( y \\). Once we find the solutions for \\( y \\), we substitute them back using the original substitution \\( y = x^n \\) to determine the corresponding values of \\( x \\). In short:

  * Let \\( y = x^n \\), so the equation becomes \\( ay^2 + by + c = 0 \\).

  * Solve the quadratic equation in \\( y \\).

  * For each solution \\( y_i \\), solve \\( x^n = y_i \\) to find the corresponding \\( x \\)-values.

  * Check all solutions in the original equation to ensure they are valid.


##### It is important to verify each solution by substituting it into the original equation to ensure its validity. The process may seem complex, but it’s actually practical and straightforward to implement.

When solving equations through substitution, extra solutions, called extraneous solutions, can sometimes appear. These might satisfy the transformed equation but not the original one. By substituting each solution back into the original equation, we ensure that it truly works and doesn’t violate any implicit conditions.

## Example

Solve the equation \\(3x^4-7x^2 + 2 = 0\\).

Let’s substitute \\(y = x^2\\) to transform the equation into a quadratic equation:

\\[3y^2-7y + 2 = 0\\]


We can use the [quadratic formula](<../quadratic-formula>) to find the value of \\(y\\). We obtain:

\begin{align*} y &= \frac{{-(-7) \pm \sqrt{{(-7)^2-4 \cdot 3 \cdot 2}}}}{{2 \cdot 3}} \\\\[0.6em] &= \frac{{7 \pm \sqrt{{49-24}}}}{6} \\\\[0.6em] &= \frac{{7 \pm \sqrt{25}}}{6} \\\\[0.6em] &= \frac{{7 \pm 5}}{6} \end{align*}

So, we found:

\\[y_1 = \frac{12}{6} = 2\\] \\[y_2 = \frac{2}{6} = \frac{1}{3}\\]


Once the solution for \\(y\\) are found, they can be substituted back into the original equation with \\(x^2 = y\\) to find the corresponding values of \\(x\\).

For \\(y = 2\\) we have: \\[x^2 = 2 \quad \rightarrow \quad x = \pm \sqrt{2} \\]

For \\(y = \frac{1}{3}\\) we have \\[x^2 = \frac{1}{3} \quad \rightarrow \quad \pm \sqrt{\frac{1}{3}} \\]


It’s essential to verify that the solutions are correct by substituting them back into the original equation to ensure accuracy. To check if the given values are solutions to the equation \\(3x^4-7x^2 + 2 = 0\\), we must plug them into the equation and verify if they satisfy it.

For \\( x = \sqrt{2} \\):

\\[3(\sqrt{2})^4-7(\sqrt{2})^2 + 2 = 3 \cdot 2^2-7 \cdot 2 + 2 \\] \\[= 12 - 14 + 2 = 0 \\]


For \\( x = -\sqrt{2} \\): \\[3(-\sqrt{2})^4-7(-\sqrt{2})^2 + 2 = 3 \cdot 2^2-7 \cdot 2 + 2 \\] \\[= 12 - 14 + 2 = 0 \\]


For \\( x = \sqrt{\dfrac{1}{3}} \\): \\[3\left(\sqrt{\frac{1}{3}}\right)^4-7\left(\sqrt{\frac{1}{3}}\right)^2 + 2 = 3 \cdot \left(\frac{1}{3}\right)^2 - 7 \cdot \frac{1}{3} + 2 \\] \\[= \frac{3}{9}-\frac{7}{3} + 2 = 0 \\]


For \\( x = -\sqrt{\dfrac{1}{3}} \\): \\[3\left(-\sqrt{\frac{1}{3}}\right)^4-7\left(-\sqrt{\frac{1}{3}}\right)^2 + 2 = 3 \cdot \left(\frac{1}{3}\right)^2-7 \cdot \frac{1}{3} + 2 \\] \\[= \frac{3}{9}- \frac{7}{3} + 2 = 0 \\]

All solutions satisfy the original equation.

The solution to the equation is: \\[x = \pm \sqrt{2} \quad \text{and} \quad \pm \sqrt{\frac{1}{3}}\\]
