# Runge-Kutta 4th Order Method

**RK4:** For $y'=f(x,y)$, $y(x_0)=y_0$, step $h$,

$$k_1=f(x_n,y_n),\;k_2=f(x_n+h/2,y_n+hk_1/2),\;k_3=f(x_n+h/2,y_n+hk_2/2),\;k_4=f(x_n+h,y_n+hk_3)$$

$$y_{n+1}=y_n+h(k_1+2k_2+2k_3+k_4)/6$$

## Accuracy

### Error Orders
Local truncation $O(h^{5})$, global $O(h^{4})$ — halving $h$ reduces global error by $16×$, vs Euler's $O(h^{2})$ local / $O(h)$ global.

### Weighted Average
Increment is weighted average of four slopes, capturing curvature better than Euler's single tangent.

## Example

$y'=y$, $y(0)=1$, $h=0.1$: $k_1=1$, $k_2=1.05$, $k_3=1.0525$, $k_4=1.10525$, $y_1=1+0.1(1+2.1+2.105+1.10525)/6\approx1.10517$ vs exact $e^{0.1}=1.10517$ (error $10^{-7}$).
