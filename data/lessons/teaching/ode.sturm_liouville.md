# Sturm-Liouville Theory

**Sturm-Liouville:** $(py')' + (\lambda w - q)y = 0$ with separated boundary conditions. The operator is self-adjoint under the weight $w$: real eigenvalues $\lambda_n$, $w$-orthogonal eigenfunctions, discrete spectrum, complete basis.

## Worked: modes of y'' + λy = 0

Solve $y'' + \lambda y = 0$ with $y(0) = y(\pi) = 0$:
1. General solution: $y = A\cos(\sqrt{\lambda}\,x) + B\sin(\sqrt{\lambda}\,x)$.
2. $y(0) = 0$ kills the cosine: $A = 0$.
3. $y(\pi) = B\sin(\sqrt{\lambda}\,\pi) = 0$ with $B \neq 0$ forces $\sqrt{\lambda} = n$, so $\lambda_n = n^2$ with eigenfunction $\sin(nx)$.

So the boundary pins quantize $\lambda$: only squares $1, 4, 9, \dots$ resonate.

## Worked: orthogonality of sin x and sin 2x

Check $\int_0^{\pi} \sin x \sin 2x \, dx$:
1. Product to sum: $\sin x \sin 2x = (\cos x - \cos 3x)/2$.
2. Integrate: $[\sin x/1 - \sin 3x/3]_0^{\pi} / 2 = 0$.
3. So distinct modes are orthogonal — the $w = 1$ inner product kills cross terms.

So eigenfunction expansion works exactly like a Fourier series: project, then sum.

## Self-adjoint spectrum

Self-adjointness makes every eigenvalue real and eigenfunctions for distinct eigenvalues orthogonal. On a finite interval with separated conditions the spectrum is an infinite discrete ladder $\lambda_1 < \lambda_2 < \cdots$, complete enough to expand any reasonable function — Legendre's equation ($p = 1 - x^2$) is the classic member.
