# Anatomy of an optimization problem

An optimization problem names three things before it asks for anything: the quantity to be made small, the variables that may move, and the rules those variables must obey. Writing the three down separately is the whole skill, because most modelling errors are a variable that was never declared or a constraint that was silently dropped.

## Objective, variables and constraints

The objective is a single number computed from the variables, and the solver either minimizes or maximizes it. Variables carry a declared domain, real, non-negative, integer or binary, and that domain is a constraint in its own right. Constraints are equalities or inequalities that every accepted point must satisfy. A model is complete when each symbol in the objective appears in the variable list and each restriction in the sentence appears as a constraint.

## Classifying a model

Classification follows the shape of the objective and the constraints. A linear objective with linear constraints is a linear program. A quadratic objective with linear constraints is a quadratic program. Integer or binary variables make the model discrete. A model with a product of two variables, or a variable in a denominator, is nonlinear. The class decides which algorithm is legitimate, so the label is a claim about structure, not about difficulty.

## Writing an objective from a sentence

Translate the sentence in stages. Name the decision, assign a symbol and a unit. Write the quantity the sentence wants to reduce or increase as a sum of per-unit contributions. Then write each restriction as an inequality with a direction that matches the wording, at most, at least or exactly. Check the units of every term in the objective and in each constraint, because a mismatch there usually means a missing conversion factor rather than a wrong formula.

A model that states the objective, the variable domains and the constraints can be checked by someone else. A model that states only a formula cannot, because the reader has to guess which symbols were free and which rules were intended.
