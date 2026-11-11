# Sensitivity ranges

A sensitivity range is the interval over which one data change leaves the optimal basis unchanged. Inside the range the shadow prices stay valid and the objective changes linearly. Outside it a different basis becomes optimal and the model must be re-solved.

## Finding an allowable range

Vary one coefficient while holding the basis fixed and require every basic variable to remain non-negative. The inequalities on the right-hand side give the range. The same computation applies to an objective coefficient, where the reduced costs must stay non-negative. The range is therefore derived, not guessed.

## Testing a coefficient change

Substitute the new coefficient and recompute the affected values. If the basis is unchanged, the shadow prices are unchanged and the new objective follows from the linear formula. If a basic variable would become negative, the change is outside the range and the test stops.

## Updating the objective

Inside the range the objective changes by the price times the change. For an objective coefficient the change is the value of the corresponding variable times the coefficient change. Both formulas are exact for a linear program and require no re-solve.

Sensitivity ranges say how far a conclusion can be pushed. They are the difference between a number that is valid and a number that is valid nearby.
