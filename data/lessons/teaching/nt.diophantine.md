# Diophantine Equations

A **Diophantine equation** is an equation where only integer solutions are allowed. The simplest type is the linear Diophantine equation:

$$ax + by = c$$

where $a, b, c$ are integers and we seek integer solutions for $x, y$.

**Existence of solutions:**
The equation $ax + by = c$ has integer solutions **if and only if** $\gcd(a, b)$ divides $c$.

**Finding a solution using the Euclidean Algorithm:**

**Example:** Solve $252x + 105y = 21$.

First, compute $\gcd(252, 105)$ using the Euclidean algorithm:
$252 = 105 \times 2 + 42$
$105 = 42 \times 2 + 21$
$42 = 21 \times 2 + 0$

Work backwards to express 21 as a combination:
$21 = 105 - 42 \times 2$
$42 = 252 - 105 \times 2$
Substitute: $21 = 105 - (252 - 105 \times 2) \times 2$
$21 = 105 \times 5 - 252 \times 2$

So $x = -2$, $y = 5$ is one solution: $252(-2) + 105(5) = -504 + 525 = 21$.

**General solution:**
If $(x_0, y_0)$ is one solution, the general solution is:
$$x = x_0 + \frac{b}{d}t, \quad y = y_0 - \frac{a}{d}t$$
where $d = \gcd(a, b)$ and $t$ is any integer.

**Example:** For $252x + 105y = 21$, $d = 21$:
$$x = -2 + \frac{105}{21}t = -2 + 5t$$
$$y = 5 - \frac{252}{21}t = 5 - 12t$$

Check with $t = 1$: $x = 3, y = -7$
$252(3) + 105(-7) = 756 - 735 = 21$ ✓
