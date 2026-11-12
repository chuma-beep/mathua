# The Hahn-Banach theorem

Hahn-Banach extends a bounded functional from a subspace to the whole space without increasing its bound. It is the main existence tool for linear functionals.

## Extension of functionals

Let \(f\) be linear on a subspace \(M\) of a real vector space and let \(p\) be a seminorm with \(f(x) \le p(x)\) on \(M\). Then \(f\) extends to \(F\) on the whole space with \(F \le p\) everywhere. On \(M = \{(x,0)\}\) in \(\mathbb{R}^2\) take \(f(x,0) = 3x\) and \(p(x,y) = 4\sqrt{x^2+y^2}\). The bound holds since \(3|x| \le 4|x|\), and \(F(x,y) = 3x\) extends it.

## The dominated form

The domination is what keeps the extension bounded. A functional bounded by \(L\|x\|\) on \(M\) extends with the same bound \(L\). The extension is generally not unique: on \(\mathbb{R}^2\) the functionals \(3x + cy\) all extend \(f\) and all respect the bound when \(|c| \le \sqrt{7}\).

## Separation

A point outside a closed subspace can be separated from it. If \(x_0 \notin M\) and \(M\) is closed, there is a functional vanishing on \(M\) with \(F(x_0) = 1\). This is the geometric form, and it is what proves that a nonzero vector has a nonzero functional detecting it.

Hahn-Banach builds a functional from local data, and the seminorm bound is the only constraint the extension must respect.
