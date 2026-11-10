# AdaGrad optimization

## Accumulating squared gradients

AdaGrad adapts the learning rate to each parameter. It accumulates the squares of past gradients and divides the current gradient by the square root of that accumulator. At step t, add the square of the current gradient to the accumulator: G_t = G_{t-1} + g_t^2. A first gradient of 3 gives G_1 = 9. A second gradient of 4 gives G_2 = 9 + 16 = 25.

## The first update

The first update is w_1 = w_0 - eta g_1 / sqrt(G_1). With g_1 = 3, G_1 = 9 and eta = 10, the denominator is 3 and the update is 10. Starting at w_0 = 20 gives w_1 = 10. The sign of the gradient controls the direction. A positive gradient decreases the parameter.

## Later updates and diminishing steps

At the second step, use G_2 in the denominator. With g_2 = 4 and G_2 = 25, the second update is 10 x 4 / 5 = 8. The second step is smaller than the first because the accumulator grew. This diminishing step is the defining behaviour of AdaGrad. Parameters with large past gradients receive smaller steps.
