# Finite differences for numerical derivatives

Finite differences approximate derivatives from function values at nearby points. The choice of points controls the accuracy and the error order.

## Forward difference

The forward difference is \((f(x+h)-f(x))/h\). For \(f(x)=x^2\), it equals \(2x+h\), so it carries an error of order \(h\).

A smaller step reduces the truncation error but can increase roundoff error.

The forward difference needs only one value ahead of the point.

## Central difference

The central difference is \((f(x+h)-f(x-h))/(2h)\). For \(x^2\), it equals \(2x\) exactly, so the leading error cancels.

The central difference is second-order accurate for smooth functions. It needs values on both sides of the point.

At a domain boundary the central difference may not be available.

## Second differences

The second difference \(f(x+h)-2f(x)+f(x-h)\) approximates \(h^2f''(x)\). For \(x^2\), it equals \(2h^2\).

Divide by \(h^2\) to approximate the second derivative. The same differencing pattern extends to higher derivatives.

More points usually improve the order but also amplify roundoff.
