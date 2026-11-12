# Uniform limits of continuous functions

A uniform limit of continuous functions is continuous. Pointwise limits need not be, so uniformity is exactly the strength this theorem needs.

## The uniform metric

The distance between bounded functions is \(d(f,g) = \sup_x |f(x) - g(x)|\). The sequence \(f_n(t) = \sqrt{t^2 + 1/n}\) converges to \(|t|\) with \(d(f_n, |\cdot|) \le 1/\sqrt{n}\), which is 0.1 at \(n = 100\). Convergence in this metric is uniform convergence.

## Continuity is preserved

If each \(f_n\) is continuous at \(x\) and \(f_n \to f\) uniformly, then \(f\) is continuous at \(x\). Take \(\varepsilon = 0.01\) and choose \(n\) with \(d(f_n,f) < 0.003\). Then choose \(\delta\) from the continuity of that \(f_n\). The three errors add to at most 0.01. The limit \(|t|\) is continuous even though each approximant is smooth at the origin.

## Interchanging limits

Uniform convergence lets two limits be exchanged: \(\lim_{t \to x} \lim_n f_n(t) = \lim_n \lim_{t \to x} f_n(t)\). Pointwise convergence does not. The functions \(t^n\) on \([0,1]\) converge pointwise to the function that is 0 below 1 and 1 at 1, and the two orders disagree at the endpoint.

Uniformity controls the function as a whole, and that control is what continuity of the limit requires.
