# Variation of Parameters

**Variation of parameters:** For $y''+p(x)y'+q(x)y=g(x)$ with homogeneous basis $y_1, y_2$ and Wronskian $W=y_1y_2'-y_1'y_2$, a particular solution is

$$y_p = u_1 y_1 + u_2 y_2, \quad u_1'=-\frac{y_2 g}{W},\; u_2'=\frac{y_1 g}{W}$$

## Building the Particular Solution

### The Wronskian
$W$ measures linear independence of $y_1, y_2$. For constant-coefficient equations it is often constant or an exponential and never zero when $y_1, y_2$ are independent.

### Integrating for $u_1, u_2$
Integrate $u_1'$ and $u_2'$ to get $u_1, u_2$, then form $y_p$. Any integration constant can be dropped (it reproduces the homogeneous part).

## Example

Solve $y''+y=\sec x$. Homogeneous $y_1=\cos x$, $y_2=\sin x$, $W=1$. Then $u_1'=-\sin x \sec x=-\tan x$, so $u_1=\ln|\cos x|$. $u_2'=\cos x \sec x=1$, so $u_2=x$. Hence $y_p=\ln|\cos x|\cos x + x\sin x$, and $y=C_1\cos x+C_2\sin x+y_p$.
