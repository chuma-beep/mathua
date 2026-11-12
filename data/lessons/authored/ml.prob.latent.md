# Latent variable models

A latent variable model explains observations through unobserved factors. Marginalising those factors couples the data and makes the likelihood intractable in most interesting cases.

## The marginal likelihood

The likelihood integrates the joint over the latent variables. The integral has no closed form once the model is nonlinear.

## Parameter counting

A model with latent dimension \(d\) and observation dimension \(m\) has a parameter count that grows with both. That count sets the capacity of the model.

## Inference

Expectation-maximisation or variational inference fills the gap left by the intractable integral. Both alternate between inferring the latents and updating the parameters.

To check your understanding, write the integral, count the parameters and run an E step.

Latent variables buy structure at the price of an integral.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller numbers and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
