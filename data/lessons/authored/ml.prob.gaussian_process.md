# Gaussian process priors

A Gaussian process places a joint Gaussian distribution over function values at every finite set of inputs. The kernel fixes the covariance, so conditioning on observations gives a posterior over functions rather than a single fit.

## The covariance matrix

With \(n\) observed inputs the covariance matrix has \(n^2\) entries. Symmetry leaves \(n(n+1)/2\) independent entries before noise is added.

## Posterior prediction

The posterior mean at a new input is a weighted combination of the observed values. The weights come from the kernel between the new input and each observation.

## The noise term

Adding noise variance \(\sigma^2\) to the diagonal makes the matrix positive definite. Without that term the matrix can be singular when two inputs coincide.

To check your understanding, count the covariance entries, count the weights and add the diagonal term.

A Gaussian process is a prior over functions, and conditioning is the update.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller numbers and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
