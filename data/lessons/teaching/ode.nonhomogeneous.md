# Nonhomogeneous Linear ODEs

A second-order linear nonhomogeneous ODE has the form:

$$y'' + p(x)y' + q(x)y = r(x)$$

where $r(x) \neq 0$. The **general solution** is:

$$y = y_h + y_p$$

where $y_h$ solves the homogeneous equation ($y'' + py' + qy = 0$) and $y_p$ is any particular solution.

## Method of Undetermined Coefficients

For constant-coefficient equations with $r(x)$ of special form:

**Case 1:** $r(x) = e^{\alpha x}$ — Try $y_p = Ae^{\alpha x}$.
**Case 2:** $r(x) = \cos(\beta x)$ or $\sin(\beta x)$ — Try $y_p = A\cos(\beta x) + B\sin(\beta x)$.
**Case 3:** $r(x) = polynomial$ — Try $y_p$ as a general polynomial of the same degree.

**Modification Rule:** If the trial form is already a solution of the homogeneous equation, multiply by $x$ (or $x^2$ if needed).

**Example 1:** Solve $y'' - 3y' + 2y = e^{3x}$.

Homogeneous: $r^2 - 3r + 2 = 0$, roots $r = 1, 2$.
$y_h = C_1e^x + C_2e^{2x}$

Particular: Try $y_p = Ae^{3x}$, substitute:
$9Ae^{3x} - 9Ae^{3x} + 2Ae^{3x} = 2Ae^{3x} = e^{3x}$, so $A = 1/2$.

General solution: $y = C_1e^x + C_2e^{2x} + \frac{1}{2}e^{3x}$

**Example 2:** Solve $y'' + y = \sin x$.

Homogeneous: $r^2 + 1 = 0$, roots $r = \pm i$.
$y_h = C_1\cos x + C_2\sin x$

Particular: Since $\sin x$ is already in $y_h$, try $y_p = Ax\cos x + Bx\sin x$.
After substitution: $A = 0$, $B = -\frac{1}{2}$.
$y_p = -\frac{1}{2}x\cos x$

General solution: $y = C_1\cos x + C_2\sin x - \frac{1}{2}x\cos x$

## Variation of Parameters

For any $r(x)$, given two linearly independent solutions $y_1, y_2$ of the homogeneous equation:

$$y_p = -y_1 \int \frac{y_2 r}{W} dx + y_2 \int \frac{y_1 r}{W} dx$$

where $W = y_1y_2' - y_2y_1'$ is the Wronskian.
