# Consensus optimization

Consensus optimisation gives each agent a local copy of a shared variable and penalises disagreement.

## Local copies

Each agent holds its own copy and has its own objective. A consensus constraint forces every copy to equal a common value.

## Averaging

The shared value is the average of the local copies. Each agent then takes a step on its own objective plus a pull toward the average.

## Agreement

The disagreement is the spread of the copies. A small spread means the agents have reached consensus.

To check your understanding, write the consensus constraint, compute the average and measure the spread.

Consensus is decomposition where the coupling is an average rather than a price.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
