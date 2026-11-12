# Adaptive step sizes

Adagrad divides each coordinate's step by the square root of its accumulated squared gradients. Frequently used coordinates slow down.

## The accumulator

The accumulator adds the square of each gradient coordinate. The step for a coordinate is the base step times the gradient divided by the square root of the accumulator.

## The effect

A coordinate with a large history takes small steps, which suits sparse features. The accumulator only grows, so the step eventually becomes very small.

## A scalar case

In one dimension the accumulated step after many equal gradients is the base step divided by the square root of the count. That is a diminishing schedule produced automatically.

To check your understanding, update the accumulator, explain a slowed coordinate and take equal gradients.

Adagrad is a per-coordinate step size learned from the gradient history.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
