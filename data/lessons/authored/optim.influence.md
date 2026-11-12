# Influence functions

An influence function estimates how a fitted model would change if one training point were reweighted. It uses the inverse Hessian of the empirical loss to turn a small data perturbation into a parameter perturbation.

## The estimate

The change in parameters is approximately the negative inverse Hessian applied to the gradient of the perturbed point. The Hessian is evaluated at the fitted parameters and the gradient at the point being removed.

## The linear solve

The inverse Hessian is never formed. A conjugate gradient solve computes the product needed for one influence value, and the same factorisation serves many points.

## The caveats

The estimate assumes the loss is twice differentiable and the perturbation is small. It is unreliable for nonconvex models, for points that change the active set and for heavily regularised fits.

To check your understanding, perturb one point, use conjugate gradient and require differentiability.

An influence function is a first-order answer to a counterfactual question about the training set, and its reliability is bounded by the local quadratic model.
