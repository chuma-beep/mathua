# Weak convergence

Weak convergence is convergence tested against every bounded functional. It is coarser than norm convergence and it is often compact where the norm ball is not.

## The weak topology

A net \(x_\alpha\) converges weakly to \(x\) when \(f(x_\alpha) \to f(x)\) for every bounded functional \(f\). The weak topology is the smallest one making every such functional continuous. In a Hilbert space this reads \(\langle x_\alpha, y\rangle \to \langle x, y\rangle\) for every \(y\).

## Weak versus norm

Norm convergence implies weak convergence, never the reverse. In \(\ell^2\) the standard basis vectors \(e_n\) satisfy \(\langle e_n, y\rangle = y_n \to 0\) for every \(y\), so \(e_n \to 0\) weakly. The norms stay at 1, so there is no norm convergence and no norm-convergent subsequence.

## Compactness

A bounded sequence in a reflexive space has a weakly convergent subsequence, by Banach-Alaoglu applied to the closed unit ball of the bidual. In a Hilbert space the unit ball is weakly compact. The sequence \(e_n\) converges weakly to 0 without any subsequence converging in norm, which shows the two topologies differ even on a bounded set.

Weak convergence keeps the limits that compactness promises, and it gives up the norm to keep them.
