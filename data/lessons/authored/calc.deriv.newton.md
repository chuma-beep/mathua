# Newton's Method and Error Propagation

## Newton's method

Newton's method approximates a root of a differentiable function \\(f\\) by following tangent lines to the \\(x\\)-axis. From a guess \\(x_{n}\\), the tangent line at \\((x_{n}, f(x_{n}))\\) meets the axis at

\\[x_{n+1} = x_{n} - \frac{f(x_{n})}{f'(x_{n})}\\]

provided \\(f'(x_{n}) \neq 0\\). Each step is the root of the linearization of \\(f\\) at \\(x_{n}\\), so near a simple root the approximations converge quadratically: the number of correct digits roughly doubles per step.

### Example

For \\(f(x) = x^{2} - 2\\) starting at \\(x_{0} = 1\\): \\(x_{1} = 1 - (1-2)/2 = 1.5\\), then \\(x_{2} = 1.5 - (2.25-2)/3 \approx 1.4167\\), already close to \\(\sqrt{2}\\).

The method fails when \\(f'(x_{n}) = 0\\), when the guess is far from any root, or when iterates cycle instead of settling.

## Error propagation with differentials

A measurement error \\(\Delta x\\) in the input propagates through \\(y = f(x)\\) as

\\[\Delta y \approx f'(x)\,\Delta x\\]

which is the differential \\(dy = f'(x)\,dx\\) evaluated at the error. The derivative acts as an amplification factor: steep functions magnify input errors, flat ones shrink them.

### Example

For \\(y = x^{2}\\) measured at \\(x = 3 \pm 0.1\\): \\(\Delta y \approx 2\cdot 3 \cdot 0.1 = 0.6\\), so \\(y = 9 \pm 0.6\\). For \\(y = \sqrt{x}\\) at \\(x = 100 \pm 1\\): \\(\Delta y \approx 1/(2\sqrt{100}) = 0.05\\), so \\(y = 10 \pm 0.05\\).
