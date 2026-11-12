# Symplectic integrators

A symplectic integrator preserves the phase-space area of a Hamiltonian system. It does not conserve the energy exactly, but it conserves a nearby modified Hamiltonian and therefore keeps the energy bounded over long times.

## Symplectic map

A step is symplectic when its Jacobian satisfies the same identity as a Hamiltonian flow. The identity says that the map preserves the two-form that measures oriented area in phase space.

## Stormer-Verlet

For a separable Hamiltonian the method alternates a half kick from the potential with a full drift from the kinetic term and another half kick. It is explicit, second order and symplectic.

## Long-time behaviour

The energy error of a symplectic method oscillates within a bound instead of growing linearly. A nonsymplectic Runge-Kutta method of the same order drifts, which is why symplectic methods dominate long molecular and orbital integrations.

To check your understanding, write the preserving identity, split the Hamiltonian and compare energy errors.

Preserving the geometric structure matters more than a high order when the integration runs for millions of steps.
