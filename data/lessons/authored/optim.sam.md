# Sharpness-aware minimization

Sharpness-aware minimisation looks for a point whose loss stays low in a small neighbourhood, not just at the point itself. The objective is a maximum over a perturbation, so the method is a minimax problem with an inner ascent step.

## The inner ascent

The method computes the gradient, moves to the point that maximises the loss within a ball, and evaluates the gradient there. The perturbation size is the radius of the ball and is usually small.

## The outer step

The update uses the gradient at the perturbed point with a normal step size. The perturbed point is discarded, so the method pays one extra gradient per step.

## Why sharpness matters

A flat minimum tolerates a shift between training and deployment distributions. Empirical results show better generalisation, though part of the effect may come from the implicit regularisation of the perturbed step.

To check your understanding, compute the gradient, use the perturbed gradient and prefer a flat minimum.

Sharpness-aware minimisation is a minimax method whose adversary is a small ball around the parameters.
