# Riemann Sums

A Riemann sum turns a continuous function into a finite sum of rectangles. It is the bridge between adding values and integrating a function.

## Partitioning an interval

A partition of \([a,b]\) divides the interval into subintervals. If the subintervals have equal width and there are \(n\) of them, the width is \(\Delta x=(b-a)/n\). Choose one sample point in each subinterval, such as the left endpoint, the right endpoint or the midpoint.

The choice of sample point matters for a finite sum. For \(f(x)=x^2\) on \([0,3]\) with three subintervals, the right endpoints are 1, 2 and 3. The right Riemann sum is \(1^2+2^2+3^2=14\), because each rectangle has width one.

## Building a Riemann sum

For a partition with widths \(\Delta x_i\) and sample points \(x_i^*\), the Riemann sum is \(\sum f(x_i^*)\Delta x_i\). A right-endpoint sum uses the right end of each subinterval. A left-endpoint sum uses the left end.

On \([0,4]\) with four unit subintervals, the left endpoints for \(f(x)=x\) are 0, 1, 2 and 3, so the left sum is 6. The right endpoints are 1, 2, 3 and 4, so the right sum is 10. The two sums differ because the function is increasing.

## The definite integral as a limit

As the number of subintervals grows and the largest width shrinks to zero, the Riemann sums approach the definite integral. For an integrable function, \(\lim_{n\to\infty}\sum f(x_i^*)\Delta x=\int_a^b f(x)\,dx\).

The limit is independent of the sample points for a continuous function. For \(f(x)=x\) on \([0,4]\), both the left and right sums approach the area 8. This limiting value is the definite integral, and it can be positive, negative or zero.
