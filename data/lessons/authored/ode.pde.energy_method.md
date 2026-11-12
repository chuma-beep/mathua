# Energy methods for parabolic equations

Multiplying a parabolic equation by its own solution and integrating by parts produces an energy that cannot increase. The estimate is the standard proof that the heat equation smooths data and forgets its initial condition.

## Multiply and integrate

For u_t = u_xx, multiply by u and integrate over the domain. The product u u_t is half the time derivative of u squared, and u u_xx is a divergence minus the square of the gradient.

## Boundary terms

A Dirichlet condition makes the boundary term vanish because the solution is fixed there. A Neumann condition makes the normal derivative vanish. Either way the boundary contributes no positive term, which is what turns the identity into an inequality.

## The decay estimate

The energy decreases at a rate equal to the Dirichlet integral of the gradient. A Poincare inequality bounds that integral below by a constant times the energy, so the energy decays exponentially and the solution is unique.

To check your understanding, multiply the equation by the solution, apply the divergence theorem and apply the Poincare inequality.

The energy identity turns a partial differential equation into an inequality about one scalar, and the inequality is usually enough to prove uniqueness.
