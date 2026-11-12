# Zero-sum games

In a zero-sum game one player's gain is the other's loss. A mixed strategy is a probability vector, the payoff is bilinear in the two strategies, and the equilibrium is a saddle point of that bilinear form.

## The payoff matrix

The row player chooses a distribution over rows and the column player a distribution over columns. The expected payoff is the bilinear product of the two distributions with the payoff matrix. A pure strategy is the special case where one coordinate is one.

## The linear program

Computing an equilibrium is a linear program: maximise a scalar value subject to the payoff against every pure column being at least that value. The dual program is the column player's problem, and strong duality makes the two values equal.

## Randomisation as a guarantee

A mixed strategy bounds the opponent's payoff regardless of what the opponent does. That guarantee is the reason randomisation appears in security and in repeated play, where a predictable pure choice is exploitable.

To check your understanding, write the bilinear expectation, maximise the floor value and bound the opponent payoff.

A zero-sum game is a linear program in disguise, and the equilibrium value is a guarantee against any opposing strategy.
