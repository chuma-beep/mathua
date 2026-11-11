# Wave 11 candidate ledger

Status: proposed on 2026-11-10.

Wave 11 raises the corpus from 2,000 to 2,050 concepts. It creates the Optimization domain and one university course, `optim`.

Every id and normalized label below was checked against every file in `data/concepts/` on 2026-11-10. There is no exact id collision. There is no normalized-label collision after replacing the first `The assignment problem` candidate with `The transportation problem`, because `graph.assignment` already uses that label.

## Allocation

| Domain | Tier | Now | Target | Add | Course |
|---|---|---:|---:|---:|---|
| optimization | university | 0 | 50 | +50 | `optim` |
| **Total** | | | | **+50** | |

Optimization is a university domain. It needs no authored solution schemas.

## Subdomains

| Subdomain | Concepts |
|---|---:|
| `optimization.foundations` | 10 |
| `optimization.unconstrained` | 10 |
| `optimization.convex` | 12 |
| `optimization.constrained` | 10 |
| `optimization.discrete` | 8 |

## Foundations: 10 candidates

| id | Label | Prerequisite | Families |
|---|---|---|---|
| optim.problem.anatomy | Anatomy of an optimization problem | calc.deriv.optimization | Identify objective, variables and constraints; classify a model; write an objective from a sentence |
| optim.feasible.set | Feasible and infeasible sets | optim.problem.anatomy | Test a point; find the feasible interval; decide whether the set is empty |
| optim.local_global | Local and global minima | optim.problem.anatomy | Compare values; identify a global minimum; distinguish local from global |
| optim.convex_set | Convex sets | optim.problem.anatomy | Test a combination; find a separating line; decide convexity from a picture description |
| optim.strict_convexity | Strict convexity and uniqueness | optim.convex_set, calc.deriv.convexity | Test strict convexity; count minimizers; use the midpoint inequality |
| optim.subgradient | Subgradients | optim.convex_set, calc.deriv.convexity | Test a subgradient; find the subdifferential of an absolute value; use the supporting line |
| optim.epigraph | Epigraphs and sublevel sets | optim.convex_set, calc.deriv.convexity | Test an epigraph point; find a sublevel set; relate convexity to the epigraph |
| optim.coercivity | Coercivity and existence | optim.local_global, calc.deriv.optimization | Test coercivity; decide existence; find a minimizing sequence bound |
| optim.first_order | First-order optimality | optim.local_global, calc.deriv.gradient | Test a stationary point; solve a gradient condition; use the descent direction |
| optim.second_order | Second-order optimality and saddle points | optim.first_order, calc.deriv.higher_order_value | Test a Hessian; classify a stationary point; find a saddle direction |

## Unconstrained optimization: 10 candidates

| id | Label | Prerequisite | Families |
|---|---|---|---|
| optim.hessian | The Hessian matrix | optim.second_order, calc.deriv.higher_order_value | Build a Hessian; evaluate it at a point; test positive definiteness |
| optim.gradient_step | A gradient descent step | optim.first_order | Take one step; predict the next iterate; compare two step directions |
| optim.step_size | Step size and convergence | optim.gradient_step | Test a step size; count iterations to tolerance; choose a stable step |
| optim.newton_optim | Newton's method for minimization | optim.hessian, calc.deriv.newton | Take a Newton step; test convergence; compare with gradient descent |
| optim.bfgs | The BFGS update | optim.gradient_step, optim.hessian | Update an inverse Hessian; test the secant condition; apply a rank-two update |
| optim.coordinate_descent | Coordinate descent | optim.gradient_step | Take a coordinate step; choose the next coordinate; test a coordinate minimum |
| optim.conjugate_direction | Conjugate directions | optim.gradient_step, optim.hessian | Test conjugacy; build a conjugate direction; count exact quadratic steps |
| optim.armijo | The Armijo condition | optim.step_size | Test sufficient decrease; choose a backtracking step; compare two trial steps |
| optim.wolfe | The Wolfe conditions | optim.armijo | Test curvature; test both conditions; choose a valid step |
| optim.trust_region_ratio | The trust-region ratio | optim.second_order, optim.hessian | Compute the ratio; accept or reject a step; update the radius |

## Convex optimization and duality: 12 candidates

| id | Label | Prerequisite | Families |
|---|---|---|---|
| optim.lp_standard_form | Linear programming standard form | alg.ineq.linear_programming, optim.feasible.set | Convert to standard form; read the matrix; count basic variables |
| optim.basic_solution | Basic feasible solutions | optim.lp_standard_form, linalg.systems.row_reduction | Find a basic solution; test feasibility; count possible bases |
| optim.simplex_pivot | The simplex pivot | optim.basic_solution | Choose an entering variable; perform a pivot; test optimality |
| optim.lp_dual | The dual linear program | optim.lp_standard_form | Write the dual; read a dual variable; compare primal and dual dimensions |
| optim.dual_function | The Lagrange dual function | optim.lp_dual, calc.multi.lagrange | Evaluate the dual function; minimize over a primal variable; find a dual bound |
| optim.weak_duality | Weak duality | optim.dual_function | Compare bounds; test a dual-feasible pair; prove a gap sign |
| optim.strong_duality | Strong duality | optim.weak_duality, optim.slater | Test a zero gap; find a primal-dual pair; decide whether strong duality applies |
| optim.complementary_slackness | Complementary slackness | optim.strong_duality | Test a pair; find a complementary solution; use a positive multiplier |
| optim.shadow_price | Shadow prices | optim.complementary_slackness | Read a shadow price; estimate a constraint change; compare two prices |
| optim.lp_sensitivity | Sensitivity ranges | optim.shadow_price | Find an allowable range; test a coefficient change; update the objective |
| optim.kkt | Karush-Kuhn-Tucker conditions | optim.dual_function, optim.equality_lagrange | Test a KKT point; solve a small system; identify a missing condition |
| optim.slater | Slater's condition | optim.convex_set, optim.dual_function | Test strict feasibility; decide whether the condition holds; find an interior point |

## Constrained optimization: 10 candidates

| id | Label | Prerequisite | Families |
|---|---|---|---|
| optim.equality_lagrange | Equality-constrained optimization | calc.multi.lagrange, optim.second_order | Solve a two-variable problem; test a multiplier; classify a constrained point |
| optim.active_constraint | Active and inactive constraints | optim.kkt | Identify active constraints; test a candidate; count active inequalities |
| optim.penalty | Quadratic penalty method | optim.equality_lagrange | Evaluate a penalized objective; take one penalty step; choose a penalty parameter |
| optim.barrier | Logarithmic barrier method | optim.active_constraint | Evaluate a barrier; take a barrier step; test the central path |
| optim.augmented_lagrangian | Augmented Lagrangian | optim.penalty, optim.dual_function | Update a multiplier; evaluate the augmented objective; test a convergence step |
| optim.projected_gradient | Projected gradient descent | optim.gradient_step, optim.convex_set | Project a point; take a projected step; test a bound constraint |
| optim.active_set | Active-set methods | optim.active_constraint, optim.simplex_pivot | Choose the working set; take an active-set step; test optimality on a face |
| optim.trust_region_subproblem | The trust-region subproblem | optim.trust_region_ratio, optim.kkt | Test a boundary solution; solve a two-variable subproblem; find the multiplier |
| optim.kkt_matrix | The KKT matrix | optim.kkt, optim.hessian | Assemble a KKT matrix; test nonsingularity; solve a small system |
| optim.duality_gap | The duality gap | optim.strong_duality, optim.complementary_slackness | Compute a gap; test a stopping rule; compare two candidate pairs |

## Discrete optimization: 8 candidates

| id | Label | Prerequisite | Families |
|---|---|---|---|
| optim.integer_program | Integer programming | optim.lp_standard_form | Classify a model; test integrality; write an integer objective |
| optim.lp_relaxation | LP relaxation | optim.integer_program | Solve a relaxation; compare bounds; test an integrality gap |
| optim.branch_bound | Branch and bound | optim.integer_program, optim.lp_relaxation | Take a branch; compute a bound; prune a node |
| optim.cutting_plane | Cutting planes | optim.integer_program | Test a cut; add a valid inequality; decide whether a cut is violated |
| optim.knapsack | The knapsack problem | optim.integer_program, discrete.algorithms.dynamic_programming | Solve a small instance; read a DP table; test a fractional bound |
| optim.transportation | The transportation problem | optim.lp_standard_form, optim.simplex_pivot | Balance a table; find a basic solution; compute a transportation cost |
| optim.set_cover | Set cover | optim.integer_program, graph.vertex_cover | Test a cover; find a minimum cover; compare with a vertex cover |
| optim.integrality_gap | The integrality gap | optim.lp_relaxation, optim.branch_bound | Compute a gap; compare two formulations; test a rounding bound |

## Prerequisite map

The first ten form a connected foundations block:

```text
calc.deriv.optimization
  -> optim.problem.anatomy
     -> optim.feasible.set
     -> optim.local_global
     -> optim.convex_set
        -> optim.strict_convexity
        -> optim.subgradient
        -> optim.epigraph
  -> optim.first_order
     -> optim.second_order
        -> optim.hessian
```

All later edges point to an existing concept or an earlier same-wave candidate. No edge points to a concept that is not in the corpus.

## Overlap audit

The candidates below have adjacent existing material. Each boundary is deliberate.

| New concept | Existing neighbour | Boundary |
|---|---|---|
| optim.first_order | calc.deriv.optimization | Calculus finds a maximum or minimum of a one-variable function. The new concept states the multivariable first-order condition and tests a gradient. |
| optim.second_order | calc.deriv.convexity | Calculus classifies a one-variable stationary point. The new concept uses the Hessian and separates a saddle from a local minimum. |
| optim.hessian | calc.deriv.higher_order_value | The existing concept evaluates a higher derivative. The new concept assembles the matrix of second partials and tests definiteness. |
| optim.gradient_step | ml.gradient_descent | The ML concept describes the training loop. The new concept computes an optimization step and tests a convergence condition. |
| optim.step_size | numa.optimization.line_search | The numerical concept chooses a line-search step. The new concept tests stability and counts iterations for a quadratic. |
| optim.newton_optim | calc.deriv.newton | The calculus concept finds a root. The new concept minimizes a function with the Hessian and compares the step. |
| optim.conjugate_direction | numa.linear.cg | The numerical concept solves a linear system. The new concept tests conjugacy and the finite termination of a quadratic. |
| optim.lp_standard_form | alg.ineq.linear_programming | The algebra concept solves a two-variable problem graphically. The new concept rewrites a general LP in standard form. |
| optim.lp_dual | alg.ineq.linear_programming | The algebra concept does not build a dual. The new concept writes the dual and compares dimensions. |
| optim.kkt | calc.multi.lagrange | The calculus concept handles equalities. The new concept adds inequalities, multipliers and complementary slackness. |
| optim.penalty | optim.equality_lagrange | The new concept is a numerical method for the same equality-constrained problem. |
| optim.projected_gradient | ml.gradient_descent | The ML concept is unconstrained. The new concept projects each step onto a feasible set. |
| optim.integer_program | discrete.algorithms.dynamic_programming | The discrete concept is one algorithm. The new domain states an integer model and compares it with its relaxation. |
| optim.knapsack | discrete.algorithms.dynamic_programming | The new concept uses the existing DP as a tool and adds the fractional bound and model classification. |
| optim.transportation | graph.assignment | The graph concept is the assignment problem. The transportation problem has supplies, demands and a different basic solution. |

The remaining concepts have no adjacent existing label.

## Family design

Every family computes something. The three families are split as follows:

- **Model family:** identify or write the objective, variables, constraints or dual.
- **Computation family:** evaluate a step, bound, value, gap or certificate.
- **Decision family:** test a condition, choose a direction, accept a step or classify a point.

No family asks for the theorem name. No family is always yes or always no. Finite families report their pool size in the review.

## Review order

1. Author and inspect the first ten.
2. Run the optimization generator tests, the independent recomputation tests and the variety test.
3. Add the lesson, KP, course and enrichment rows.
4. Regenerate the graph and count artifacts.
5. Author the remaining forty only after the first ten pass.
