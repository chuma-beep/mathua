# Adversarial examples

An adversarial example is an input perturbed so that a model changes its prediction. The perturbation is usually small in a chosen norm, so a human sees no difference.

## The perturbation budget

An \(\ell_\infty\) budget of \(\epsilon\) allows every coordinate to move by at most \(\epsilon\). The ball has \(2\epsilon\) side length in each dimension.

## The fast gradient sign method

One step moves every coordinate by \(\epsilon\) in the sign of the gradient. The step costs one backward pass.

## The transfer

Perturbations often transfer between models trained on the same task. Transfer is why an attack need not know the target weights.

To check your understanding, set the norm, take the gradient and attack a surrogate.

Adversarial examples exploit the local sensitivity of a trained model.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller numbers and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
