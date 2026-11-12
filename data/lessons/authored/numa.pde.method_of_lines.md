# The method of lines

The method of lines discretises space first and leaves time continuous. The partial differential equation becomes a system of ordinary differential equations that any ODE solver can advance.

## Semi-discrete system

Each interior node contributes one equation for its unknown. With n interior nodes the system has n components and a banded Jacobian.

## Stiffness

Diffusion discretisation produces eigenvalues spread over a wide range, so the system is stiff. Explicit solvers then need very small steps.

## Choice of solver

Implicit solvers handle the stiffness but need a linear solve per step. The spatial and temporal errors must be balanced.

To check your understanding, discretise space, inspect the eigenvalues and balance both errors.

The method of lines turns a PDE into a stiff ODE system with structure.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
