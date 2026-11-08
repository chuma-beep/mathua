# Cramer's rule for two equations

Cramer's rule solves a square linear system with determinants. For two equations in two unknowns, each unknown is a ratio of two determinants: the denominator is the determinant of the coefficient matrix, and each numerator replaces one column with the constants.

## Determinants as denominators

For the system \(ax+by=e\) and \(cx+dy=f\), the coefficient determinant is \(D=ad-bc\). It measures how much the two equations disagree about direction. If \(D\ne0\), the lines cross once and the system has a unique solution.

The same \(D\) sits under every unknown. Compute it once, carefully, and check its sign. A sign error in \(D\) changes both answers, which is why the determinant deserves its own line of work.

## Building the two numerators

Replace the \(x\)-column with the constants and keep the \(y\)-column: \(D_x=ed-bf\). Replace the \(y\)-column and keep the \(x\)-column: \(D_y=af-ec\). Then \(x=D_x/D\) and \(y=D_y/D\).

The pattern is column replacement, not row replacement. Mixing the two produces a plausible number that fails substitution. After computing, substitute both values into the original equations. One quick check catches most arithmetic slips.

## Singular systems

If \(D=0\), the denominator vanishes and Cramer's rule cannot produce a unique answer. The lines are parallel or identical. A zero determinant therefore does not always mean no solution; it means the system needs a different test.

Substitute to distinguish the two cases. A contradiction such as \(0=5\) means no solution. An identity such as \(0=0\) means infinitely many solutions. Cramer's rule is a tool for the unique case, and recognising when it does not apply is part of using it well.
