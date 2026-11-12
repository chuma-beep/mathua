# Hamiltonian Monte Carlo

Hamiltonian Monte Carlo augments the target with a momentum variable and simulates Hamiltonian dynamics. Gradients steer the proposal along level sets, so proposals are far less random than a random walk.

## Leapfrog steps

Each leapfrog step evaluates the gradient once. A trajectory of \(L\) steps therefore costs \(L\) gradient evaluations.

## The acceptance ratio

The Metropolis correction compares the initial and final Hamiltonian. Energy conservation makes that ratio close to one for a stable step size.

## Step size and trajectory length

A small step size reduces discretisation error but needs more steps. The no-U-turn criterion chooses the trajectory length adaptively.

To check your understanding, count the gradient calls, compare energies and balance step size.

Hamiltonian dynamics turns gradient information into distant, well-accepted proposals.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller numbers and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
