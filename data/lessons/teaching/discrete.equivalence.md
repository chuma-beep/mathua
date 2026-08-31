# Logical Equivalence, De Morgan and Normal Forms

**Logical equivalence:** $P\equiv Q$ means $P$ and $Q$ have the same truth value in every row. De Morgan: $\lnot(p\land q)\equiv\lnot p\lor\lnot q$, $\lnot(p\lor q)\equiv\lnot p\land\lnot q$. Normal forms (CNF/DNF) are canonical conjunctions/disjunctions of literals.

## Transforming Formulas

### De Morgan's Laws
Negating a conjunction flips to a disjunction of negations. Example: $\lnot(p\land q)$ is false only when both $p,q$ true; $\lnot p\lor\lnot q$ is true when at least one fails — same table.

### Implication and Biconditional
$p\to q\equiv\lnot p\lor q$, $p\leftrightarrow q\equiv(p\to q)\land(q\to p)\equiv(p\land q)\lor(\lnot p\land\lnot q)$.

## Example

Put $p\to q$ in CNF: $p\to q\equiv\lnot p\lor q$ is already a disjunction of literals, hence CNF.
