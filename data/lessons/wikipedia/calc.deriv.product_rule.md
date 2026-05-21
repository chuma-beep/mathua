> Content sourced from [Wikipedia](https://en.wikipedia.org/wiki/Product_rule) — CC BY-SA 4.0

# Product rule

In calculus, the **product rule** (or **Leibniz rule** or **Leibniz product rule**) is a formula used to find the derivatives of products of two or more functions. For two functions, it may be stated in Lagrange's notation as
\[
(u \cdot v)' = u ' \cdot v + u \cdot v'
\]
 or in Leibniz's notation as
\[
\frac{d}{dx} (u\cdot v) = \frac{du}{dx} \cdot v + u \cdot \frac{dv}{dx}.
\]

The rule may be extended or generalized to products of three or more functions, to a rule for higher-order derivatives of a product, and to other contexts.

## Discovery
Discovery of this rule is credited to Gottfried Leibniz, who demonstrated it using "infinitesimals" (a precursor to the modern differential). (However, J. M. Child, a translator of Leibniz's papers, argues that it is due to Isaac Barrow.) Here is Leibniz's argument: Let *u* and *v* be functions. Then *d(uv)* is the same thing as the difference between two successive *uv*'s; let one of these be *uv*, and the other *u+du* times *v+dv*; then:

\[
\begin{align}
d(u\cdot v) & {} = (u + du)\cdot (v + dv) - u\cdot v \\
& {} = u\cdot dv + v\cdot du + du\cdot dv.
\end{align}
\]

Since the term *du*·*dv* is "negligible" (compared to *du* and *dv*), Leibniz concluded that

\[
d(u\cdot v) = v\cdot du + u\cdot dv
\]

and this is indeed the differential form of the product rule. If we divide through by the differential *dx*, we obtain

\[
\frac{d}{dx} (u\cdot v) = v \cdot \frac{du}{dx} + u \cdot \frac{dv}{dx}
\]

which can also be written in Lagrange's notation as

\[
(u\cdot v)' = v\cdot u' + u\cdot v'.
\]

### First proofs
Both Leibniz and Newton gave proofs that are not rigorous by modern standards. Leibniz reasoned with "infinitely smaller quantities", interpreting products as areas of rectangles, while Newton reasoned with "flowing quantities".

## Examples
*Suppose we want to differentiate \(f(x)=x^2\text {sin}(x).\) By using the product rule, one gets the derivative \(f'(x)=2x\cdot \text {sin}(x)+x^2\text {cos} (x)\) (since the derivative of \(x^2\) is \(2x,\) and the derivative of the sine function is the cosine function).
*One special case of the product rule is the constant multiple rule, which states: if is a number, and \(f(x)\) is a differentiable function, then \(c\cdot f(x)\) is also differentiable, and its derivative is \((cf)'(x)=c \cdot f'(x).\) This follows from the product rule since the derivative of any constant is zero. This, combined with the sum rule for derivatives, shows that differentiation is linear.
*The rule for integration by parts is derived from the product rule, as is (a weak version of) the quotient rule. (It is a "weak" version in that it does not prove that the quotient is differentiable but only says what its derivative is it is differentiable.)

## Proofs
### Limit definition of derivative
Let *h*(*x*) and suppose that and are each differentiable at. We want to prove that is differentiable at and that its derivative, *, is given by **(*x''). To do this, \(f(x)g(x+\Delta x)-f(x)g(x+\Delta x)\) (which is zero, and thus does not change the value) is added to the numerator to permit its factoring, and then properties of limits are used.

\[
\begin{align}
 h'(x) &= \lim_{\Delta x\to 0} \frac{h(x+\Delta x)-h(x)}{\Delta x} \\[5pt]
 &= \lim_{\Delta x\to 0} \frac{f(x+\Delta x)g(x+\Delta x)-f(x)g(x)}{\Delta x} \\[5pt]
 &= \lim_{\Delta x\to 0} \frac{f(x+\Delta x)g(x+\Delta x)-f(x)g(x+\Delta x)+f(x)g(x+\Delta x)-f(x)g(x)}{\Delta x} \\[5pt]
 &= \lim_{\Delta x\to 0} \frac{\big[f(x+\Delta x)-f(x)\big] \cdot g(x+\Delta x) + f(x) \cdot \big[g(x+\Delta x)-g(x)\big]}{\Delta x} \\[5pt]
 &= \lim_{\Delta x\to 0} \frac{f(x+\Delta x)-f(x)}{\Delta x} \cdot \lim_{\Delta x\to 0} g(x+\Delta x)
 + \lim_{\Delta x\to 0} f(x) \cdot \lim_{\Delta x\to 0} \frac{g(x+\Delta x)-g(x)}{\Delta x} \\[5pt]
 &= f'(x)g(x)+f(x)g'(x).
\end{align}
\]

The fact that \(\lim_{\Delta x\to0} g(x+\Delta x) = g(x)\) follows from the fact that differentiable functions are continuous.

### Linear approximations
By definition, if \(f, g: \mathbb{R} \to \mathbb{R}\) are differentiable at \(x\), then we can write linear approximations:

\[
f(x+h) = f(x) + f'(x)h + \varepsilon_1(h)
\]
 and
\[
g(x+h) = g(x) + g'(x)h + \varepsilon_2(h),
\]

where the error terms are small with respect to *h*: that is, \(\lim_{h \to 0} \frac{\varepsilon_1(h)}{h} = \lim_{h \to 0} \frac{\varepsilon_2(h)}{h} = 0,\) also written \(\varepsilon_1, \varepsilon_2 \sim o(h)\). Then:

\[
\begin{align}
 f(x+h)g(x+h) - f(x)g(x) &= (f(x) + f'(x)h +\varepsilon_1(h))(g(x) + g'(x)h + \varepsilon_2(h)) - f(x)g(x) \\[.5em]
 &= f(x)g(x) + f'(x)g(x)h + f(x)g'(x)h -f(x)g(x) + \text{error terms} \\[.5em]
 &= f'(x)g(x)h + f(x)g'(x)h + o(h). \end{align}
\]

The "error terms" consist of items such as \(f(x)\varepsilon_2(h), f'(x)g'(x)h^2\) and \(hf'(x)\varepsilon_1(h)\) which are easily seen to have magnitude \(o(h).\) Dividing by \(h\) and taking the limit \(h\to 0\) gives the result.

### Quarter squares
This proof uses the chain rule and the quarter square function \(q(x)=\tfrac14x^2\) with derivative \(q'(x) = \tfrac12 x\). We have:

\[
uv=q(u+v)-q(u-v),
\]

and differentiating both sides gives:

\[
\begin{align}
f' &= q'(u+v)(u'+v') - q'(u-v)(u'-v') \\[4pt]
&= \left(\tfrac12(u+v)(u'+v')\right) - \left(\tfrac12(u-v)(u'-v')\right) \\[4pt]
&= \tfrac12(uu' + vu' + uv' + vv') - \tfrac12(uu' - vu' - uv' + vv') \\[4pt]
&= vu'+uv'. \end{align}
\]

### Multivariable chain rule
The product rule can be considered a special case of the chain rule for several variables, applied to the multiplication function \(m(u,v) = uv\):

\[
{d (uv) \over dx} = \frac{\partial(uv)}{\partial u}\frac{du}{dx}+\frac{\partial (uv)}{\partial v}\frac{dv}{dx} = v \frac{du}{dx} + u \frac{dv}{dx}.
\]

### Non-standard analysis
Let *u* and *v* be continuous functions in *x*, and let *dx*, *du* and *dv* be infinitesimals within the framework of non-standard analysis, specifically the hyperreal numbers. Using st to denote the standard part function that associates to a finite hyperreal number the real infinitely close to it, this gives

\[
\begin{align}
 \frac{d(uv)}{dx} &= \operatorname{st}\left(\frac{(u + du)(v + dv) - uv}{dx}\right) \\
 &= \operatorname{st}\left(\frac{uv + u \cdot dv + v \cdot du + du \cdot dv -uv}{dx}\right) \\
 &= \operatorname{st}\left(\frac{u \cdot dv + v \cdot du + du \cdot dv}{dx}\right) \\
 &= \operatorname{st}\left(u \frac{dv}{dx} + (v + dv) \frac{du}{dx}\right) \\
 &= u \frac{dv}{dx} + v \frac{du}{dx}.
 \end{align}
\]

This was essentially Leibniz's proof exploiting the transcendental law of homogeneity (in place of the standard part above).

### Smooth infinitesimal analysis
In the context of Lawvere's approach to infinitesimals, let \(dx\) be a nilsquare infinitesimal. Then \(du = u'\ dx\) and \(dv = v'\ dx\), so that

\[
\begin{align}
d(uv) & = (u + du)(v + dv) -uv \\
 & = uv + u \cdot dv + v \cdot du + du \cdot dv - uv \\
 & = u \cdot dv + v \cdot du + du \cdot dv \\
 & = u \cdot dv + v \cdot du
\end{align}
\]

since \(du \, dv = u' v' (dx)^2 = 0.\) Dividing by \(dx\) then gives \(\frac{d(uv)}{dx} = u \frac{dv}{dx} + v \frac{du}{dx}\) or \((uv)' = u \cdot v' + v \cdot u'\).

### Logarithmic differentiation
Let \(h(x) = f(x) g(x)\). Taking the absolute value of each function and the natural log of both sides of the equation,

\[
\ln|h(x)| = \ln|f(x) g(x)|
\]

Applying properties of the absolute value and logarithms,

\[
\ln|h(x)| = \ln|f(x)| + \ln|g(x)|
\]

Taking the logarithmic derivative of both sides and then solving for \(h'(x)\):

\[
\frac{h'(x)}{h(x)} = \frac{f'(x)}{f(x)} + \frac{g'(x)}{g(x)}
\]

Solving for \(h'(x)\) and substituting back \(f(x) g(x)\) for \(h(x)\) gives:

\[
\begin{align}
h'(x) &= h(x)\left(\frac{f'(x)}{f(x)} + \frac{g'(x)}{g(x)}\right) \\
&= f(x) g(x)\left(\frac{f'(x)}{f(x)} + \frac{g'(x)}{g(x)}\right) \\
&= f'(x) g(x) + f(x) g'(x).
\end{align}
\]

Note: Taking the absolute value of the functions is necessary for the logarithmic differentiation of functions that may have negative values, as logarithms are only real-valued for positive arguments. This works because \(\tfrac{d}{dx}(\ln |u|) = \tfrac{u'}{u}\), which justifies taking the absolute value of the functions for logarithmic differentiation.

## Generalizations
### Product of more than two factors
The product rule can be generalized to products of more than two factors. For example, for three factors we have

\[
\frac{d(uvw)}{dx} = \frac{du}{dx}vw + u\frac{dv}{dx}w + uv\frac{dw}{dx}.
\]

For a collection of functions \(f_1, \dots, f_k\), we have

\[
\frac{d}{dx} \left [ \prod_{i=1}^k f_i(x) \right ]
 = \sum_{i=1}^k \left(\left(\frac{d}{dx} f_i(x) \right) \prod_{j=1,j\ne i}^k f_j(x) \right)
= \left( \prod_{i=1}^k f_i(x) \right) \left( \sum_{i=1}^k \frac{f'_i(x)}{f_i(x)} \right).
\]

The logarithmic derivative provides a simpler expression of the last form, as well as a direct proof that does not involve any recursion. The *logarithmic derivative* of a function , denoted here Logder(*f*), is the derivative of the logarithm of the function. It follows that

\[
\operatorname{Logder}(f)=\frac {f'}f.
\]

Using that the logarithm of a product is the sum of the logarithms of the factors, the sum rule for derivatives gives immediately

\[
\operatorname{Logder}(f_1\cdots f_k)= \sum_{i=1}^k\operatorname{Logder}(f_i).
\]

The last above expression of the derivative of a product is obtained by multiplying both members of this equation by the product of the \(f_i.\)

### Higher derivatives

It can also be generalized to the general Leibniz rule for the *n*th derivative of a product of two factors, by symbolically expanding according to the binomial theorem:

\[
d^n(uv) = \sum_{k=0}^n {n \choose k} \cdot d^{(n-k)}(u)\cdot d^{(k)}(v).
\]

Applied at a specific point *x*, the above formula gives:

\[
(uv)^{(n)}(x) = \sum_{k=0}^n {n \choose k} \cdot u^{(n-k)}(x)\cdot v^{(k)}(x).
\]

Furthermore, for the *n*th derivative of an arbitrary number of factors, one has a similar formula with multinomial coefficients:

\[
\left(\prod_{i=1}^kf_i\right)^{\!\!(n)}=\sum_{j_1+j_2+\cdots+j_k=n}{n\choose j_1,j_2,\ldots,j_k}\prod_{i=1}^kf_i^{(j_i)}.
\]

### Higher partial derivatives
For partial derivatives, we have

\[
{\partial^n \over \partial x_1\,\cdots\,\partial x_n} (uv)
= \sum_S {\partial^{|S|} u \over \prod_{i\in S} \partial x_i} \cdot {\partial^{n-|S|} v \over \prod_{i\not\in S} \partial x_i}
\]

where the index runs through all 2\(^{*n*}\) subsets of , and is the cardinality of. For example, when 1=*n* = 3,

\[
\begin{align} & {\partial^3 \over \partial x_1\,\partial x_2\,\partial x_3} (uv) \\[1ex]
= {} & u \cdot{\partial^3 v \over \partial x_1\,\partial x_2\,\partial x_3} + {\partial u \over \partial x_1}\cdot{\partial^2 v \over \partial x_2\,\partial x_3} + {\partial u \over \partial x_2}\cdot{\partial^2 v \over \partial x_1\,\partial x_3} + {\partial u \over \partial x_3}\cdot{\partial^2 v \over \partial x_1\,\partial x_2} \\[1ex]
& + {\partial^2 u \over \partial x_1\,\partial x_2}\cdot{\partial v \over \partial x_3}
+ {\partial^2 u \over \partial x_1\,\partial x_3}\cdot{\partial v \over \partial x_2}
+ {\partial^2 u \over \partial x_2\,\partial x_3}\cdot{\partial v \over \partial x_1}
+ {\partial^3 u \over \partial x_1\,\partial x_2\,\partial x_3}\cdot v. \\[-3ex]&\end{align}
\]

### Banach space
Suppose *X*, *Y*, and *Z* are Banach spaces (which includes Euclidean space) and *B* : *X* × *Y* → *Z* is a continuous bilinear operator. Then *B* is differentiable, and its derivative at the point (*x*,*y*) in *X* × *Y* is the linear map *D*\(_{(*x*,*y*)}\)*B* : *X* × *Y* → *Z* given by

\[
(D_\left( x,y \right)\,B)\left( u,v \right) = B\left( u,y \right) + B\left( x,v \right)\qquad\forall (u,v)\in X \times Y.
\]

This result can be extended to more general topological vector spaces.

### In vector calculus

The product rule extends to various product operations of vector functions on \(\mathbb{R}^n\):

* For scalar multiplication:
\[
(f \cdot \mathbf g)' = f'\cdot \mathbf g + f \cdot \mathbf g'
\]

* For dot product:
\[
(\mathbf f \cdot \mathbf g)' = \mathbf f' \cdot \mathbf g + \mathbf f \cdot \mathbf g'
\]

* For cross product of vector functions on \(\mathbb{R}^3\):
\[
(\mathbf f \times \mathbf g)' = \mathbf f' \times \mathbf g + \mathbf f \times \mathbf g'
\]

There are also analogues for other analogs of the derivative: if *f* and *g* are scalar fields then there is a product rule with the gradient:

\[
\nabla (f \cdot g) = \nabla f \cdot g + f \cdot \nabla g
\]

Such a rule will hold for any continuous bilinear product operation. Let *B* : *X* × *Y* → *Z* be a continuous bilinear map between vector spaces, and let *f* and *g* be differentiable functions into *X* and *Y*, respectively. The only properties of multiplication used in the proof using the limit definition of derivative is that multiplication is continuous and bilinear. So for any continuous bilinear operation,

\[
H(f, g)' = H(f', g) + H(f, g').
\]

This is also a special case of the product rule for bilinear maps in Banach space.

### Derivations in abstract algebra and differential geometry
In abstract algebra, the product rule is the defining property of a derivation. In this terminology, the product rule states that the derivative operator is a derivation on functions.

In differential geometry, a tangent vector to a manifold *M* at a point *p* may be defined abstractly as an operator on real-valued functions which behaves like a directional derivative at *p*: that is, a linear functional *v* which is a derivation,
\[
v(fg) = v(f)\,g(p) + f(p) \, v(g).
\]

Generalizing (and dualizing) the formulas of vector calculus to an *n*-dimensional manifold *M,* one may take differential forms of degrees *k* and *l*, denoted \(\alpha\in \Omega^k(M), \beta\in \Omega^\ell(M)\), with the wedge or exterior product operation \(\alpha\wedge\beta\in \Omega^{k+\ell}(M)\), as well as the exterior derivative \(d:\Omega^m(M)\to\Omega^{m+1}(M)\). Then one has the graded Leibniz rule:
\[
d(\alpha\wedge\beta)= d\alpha \wedge \beta + (-1)^{k} \alpha\wedge d\beta.
\]

## Applications
Among the applications of the product rule is a proof that

\[
{d \over dx} x^n = nx^{n-1}
\]

when *n* is a positive integer (this rule is true even if *n* is not positive or is not an integer, but the proof of that must rely on other methods). The proof is by mathematical induction on the exponent *n*. If *n* = 0 then *x*\(^{*n*}\) is constant and *nx*\(^{*n* − 1}\) = 0. The rule holds in that case because the derivative of a constant function is 0. If the rule holds for any particular exponent *n*, then for the next value, *n* + 1, we have

\[
\begin{align}
\frac{d x^{n+1{dx}
&{}= \frac{d}{dx} \left( x^n\cdot x\right) \\[1ex]
&{}= x \frac{d}{dx} x^n + x^n \frac{d}{dx} x & \text{(the product rule is used here)} \\[1ex]
&{}= x\left(n x^{n-1}\right) + x^n\cdot 1 & \text{(the induction hypothesis is used here)} \\[1ex]
&{}= \left(n + 1\right) x^n.
\end{align}
\]

Therefore, if the proposition is true for *n*, it is true also for *n* + 1, and therefore for all natural *n*.
