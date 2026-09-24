# Stability of Equilibria

**Equilibrium:** $x^*$ with $f(x^*) = 0$ for $x' = f(x)$. **Stable** means nearby solutions stay nearby; **asymptotically stable** adds convergence to $x^*$; **unstable** means some nearby solution escapes.

## Worked: decay in x' = -2x

Solve $x' = -2x$ from $x(0) = x_0$:
1. Separate: $x(t) = x_0 e^{-2t}$.
2. Ratio after one unit of time: $x(1)/x(0) = e^{-2} \approx 0.135$.
3. Every start shrinks toward 0, so 0 is asymptotically stable.

So a negative rate constant is a contraction: the equilibrium pulls all nearby starts in.

## Worked: counting equilibria of x' = 4 - x^2

Find and classify equilibria for $x' = 4 - x^2$:
1. Set $4 - x^2 = 0$: equilibria at $x = 2$ and $x = -2$.
2. Between them $x' > 0$ (rightward flow); outside, $x' < 0$ (leftward flow).
3. So $x = 2$ attracts from the left and $x = -2$ repels to the right: both half-stable, neither asymptotically stable.

So counting (zero, one, or two equilibria as $\mu$ varies) comes before classifying.

## Linearization decides hyperbolic cases

For $x' = Ax$, the origin is asymptotically stable when every eigenvalue has negative real part, unstable when any has positive real part. Lyapunov's indirect method extends this: the Jacobian at an equilibrium classifies it whenever no eigenvalue sits on the imaginary axis. Centres (pure imaginary) stay stable but never asymptotic.
