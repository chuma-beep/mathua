# Multiple kernel learning

Multiple kernel learning combines several kernels into one. The combination weights are learned with the classifier instead of being fixed in advance.

## The combination

The combined kernel is a non-negative weighted sum of base kernels. The weights are usually constrained to sum to one.

## The cost

Each base kernel supplies an \(n \times n\) Gram matrix. With \(k\) kernels the storage is \(k n^2\) entries before optimisation.

## The purpose

Different kernels see different similarity notions. Learning the weights selects the useful notions for the task.

To check your understanding, weight the sum, count the Gram matrices and compare kernel views.

Multiple kernel learning chooses the similarity rather than assuming it.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller numbers and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
