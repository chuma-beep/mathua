# Hyperparameter gradients

A hypergradient differentiates a validation loss with respect to a regularisation strength or another hyperparameter. The derivative exists even though the hyperparameter never appears in the inner gradient directly.

## The stationarity route

At the inner optimum the gradient of the training loss vanishes. Differentiating that identity with respect to the hyperparameter gives a linear system for the derivative of the weights, and the chain rule then gives the hypergradient.

## The cost

The linear system is the Hessian of the training loss, so one hypergradient needs one Hessian-vector product solve. Conjugate gradient handles the solve without forming the Hessian.

## Why it beats grid search

A hypergradient gives a direction, so the search follows the validation surface instead of sampling it. That matters when each inner training run is expensive and the hyperparameter count is large.

To check your understanding, differentiate the stationarity identity, form a Hessian-vector product and follow the validation surface.

Hypergradients replace a grid over expensive retrainings with a derivative that reuses the inner solution.
