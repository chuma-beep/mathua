# Logical Equivalence, De Morgan and Normal Forms

**Logical equivalence:** $P\equiv Q$ means $P$ and $Q$ agree in every row of their joint truth table. De Morgan's laws, $\lnot(p\land q)\equiv\lnot p\lor\lnot q$ and $\lnot(p\lor q)\equiv\lnot p\land\lnot q$, push negations onto literals. Implication and biconditional reduce the same way: $p\to q\equiv\lnot p\lor q$ and $p\leftrightarrow q\equiv(p\to q)\land(q\to p)$.

## Worked: De Morgan for AND

Compare the two sides row by row over $p,q$:
1. $\lnot(p\land q)$ is false in exactly one row: $p$ true and $q$ true.
2. $\lnot p\lor\lnot q$ is false exactly when both disjuncts fail, i.e. $p$ true and $q$ true — the same row.
3. True in the other three rows each, so both columns match everywhere.

So negating a conjunction flips it to a disjunction of negations, and dually for disjunctions.

## Worked: implication in CNF

Check where each side is false:
1. $p\to q$ is false in exactly one row: $p$ true, $q$ false.
2. $\lnot p\lor q$ fails exactly when $\lnot p$ and $q$ both fail: $p$ true, $q$ false — the same row.
3. $\lnot p\lor q$ is a disjunction of literals, hence already a CNF clause.

So $p\to q\equiv\lnot p\lor q$ is the CNF of the implication.

## Worked: biconditional as DNF

Find the rows where the biconditional is true:
1. $p\leftrightarrow q$ is true in two rows: both true, and both false.
2. Those rows give minterms $(p\land q)$ and $(\lnot p\land\lnot q)$, so $p\leftrightarrow q\equiv(p\land q)\lor(\lnot p\land\lnot q)$.
3. Grouping differently, the same table is $(p\to q)\land(q\to p)$: each implication kills one false row.

So the biconditional is both a conjunction of implications and a disjunction of agreement minterms.
