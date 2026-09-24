# Laurent Series

**Laurent series:** For $f$ analytic on an annulus $r<|z-z_0|<R$, the expansion $f(z)=\sum a_n(z-z_0)^n$ runs over all integers $n$, negative included. The negative-power part is the principal part; its $a_{-1}$ coefficient is the residue.

## Worked: one over z is its own Laurent series

Expand $f(z)=1/z$ at 0:
1. Write $f(z)=z^{-1}$: a single power, already of Laurent form.
2. Read the coefficients: $a_{-1}=1$ and every other $a_n=0$.
3. Split the parts: principal part $1/z$, regular part 0.

So a pole shows finite principal part: finitely many negative powers means a pole, and the pole order is the most negative exponent present.

## Worked: one over z squared has residue zero

Expand $f(z)=1/z^2$ at 0:
1. Write $f(z)=z^{-2}+0 \cdot z^{-1}$: the $z^{-1}$ slot is empty.
2. Read the residue: $a_{-1}=0$.
3. Conclude $\oint_{|z|=1} dz/z^2=0$: no $z^{-1}$ term, no integral.

So the residue is exactly $a_{-1}$: integrating term by term kills every power except $z^{-1}$.

## Worked: exponential of one over z is essential

Expand $f(z)=e^{1/z}$ at 0:
1. Substitute $w=1/z$ into $e^w=1+w+w^2/2+\cdots$.
2. Get $1+z^{-1}+z^{-2}/2+z^{-3}/6+\cdots$: negative powers never stop.
3. Conclude the singularity is essential, with $a_{-1}=1$.

So infinite principal part means essential singularity: near one, $f$ attains nearly every value infinitely often, unlike poles which simply blow up.
