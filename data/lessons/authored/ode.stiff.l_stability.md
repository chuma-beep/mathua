# L-stability

L-stability adds a vanishing stability function at infinity to A-stability. A very stiff mode is then annihilated in one step instead of being answered with a slowly decaying numerical ringing.

## The extra condition

The stability function must tend to zero as the real part of the step parameter tends to minus infinity. A-stability only requires the modulus to stay below one, which allows a limit of minus one.

## Backward Euler

Backward Euler has a stability function that decays to zero at infinity, so it is L-stable. The trapezoidal rule has a limit of minus one, so it is A-stable but not L-stable.

## Practical effect

An L-stable method damps a fast transient without oscillation. A merely A-stable method can produce a decaying alternation that a learner sees as numerical noise on a solution that should be smooth.

To check your understanding, state the limit at infinity, write its stability function and apply to a stiff mode.

L-stability is the condition that separates a usable stiff solver from one that rings on the fastest mode.
