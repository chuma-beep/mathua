> Content sourced from [Wikipedia](https://en.wikipedia.org/wiki/Separation_of_variables) — CC BY-SA 4.0

# Separable differential equations

In mathematics, **separation of variables** (also known as the **Fourier method**) is any of several methods for solving ordinary and partial differential equations, in which algebra allows one to rewrite an equation so that each of two variables occurs on a different side of the equation.

## Ordinary differential equations (ODE)
A differential equation for the unknown \(f(x)\) is separable if it can be written in the form

\(\frac{d}{dx} f(x) = g(x)h(f(x))\)

where \(g\) and \(h\) are given functions. This is perhaps more transparent when written using \(y = f(x)\) as:

\(\frac{dy}{dx}=g(x)h(y).\)

So now as long as *h*(*y*) ≠ 0, we can rearrange terms to obtain:

\({dy \over h(y)} = g(x) \, dx,\)

where the two variables *x* and *y* have been separated. Note *dx* (and *dy*) can be viewed, at a simple level, as just a convenient notation, which provides a handy mnemonic aid for assisting with manipulations. A formal definition of *dx* as a differential (infinitesimal) is somewhat advanced.

### Alternative notation
Those who dislike differentials as separate entities may prefer to write this as

\(\frac{1}{h(y)} \frac{dy}{dx} = g(x),\)

but that fails to make it quite as obvious why this is called "separation of variables". Integrating both sides of the equation with respect to \(x\), we have
