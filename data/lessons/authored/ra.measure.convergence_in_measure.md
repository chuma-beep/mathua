# Convergence in measure

Convergence in measure says the set where a sequence differs from its limit by more than a fixed amount shrinks to nothing. It is weaker than uniform convergence and it does not require a limit almost everywhere.

## The definition

A sequence converges in measure to \(f\) when \(\mu(\{|f_n - f| > \varepsilon\}) \to 0\) for every \(\varepsilon > 0\). On \([0,1]\) the functions \(t^n\) tend to 0 in measure, since the set where \(t^n > 0.01\) is \((0.01^{1/n}, 1]\) and its measure tends to 0.

## Almost everywhere versus in measure

Almost everywhere convergence implies convergence in measure on a finite measure space. The converse fails. The moving bump is an indicator of an interval of length \(2^{-k}\) that slides across \([0,1]\), and it converges to 0 in measure. Every point is covered infinitely often, so there is no pointwise convergence.

## Cauchy in measure

A sequence that is Cauchy in measure has a subsequence converging almost everywhere. The moving bump is Cauchy in measure, and a subsequence whose intervals shrink toward a fixed point converges there. Completeness of \(L^1\) is proved through this route.

Convergence in measure measures the size of the disagreement rather than its location, and that size is what the integral can control.
