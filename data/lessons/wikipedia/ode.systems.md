> Content sourced from [Wikipedia](https://en.wikipedia.org/wiki/System_of_differential_equations) — CC BY-SA 4.0

# Systems of differential equations

In mathematics, a **system of differential equations** is a finite set of differential equations. Such a system can be either linear or non-linear. Also, such a system can be either a system of ordinary differential equations or a system of partial differential equations.

Examples of systems of differential equations often emerge in the real world from systems that have interconnected components that exchange quantities between each other. This could be water cascading down a set of tanks, bouncing masses connected with springs, or heat transferring between different parts of a building using Newton's law of cooling.

## Linear systems of differential equations

A first-order linear system of ODEs is a system in which every equation is first order and depends on the unknown functions linearly. Here, we consider systems with an equal number of unknown functions and equations. These may be written as

\(\frac{dx_j}{dt} = a_{j1}(t) x_1 + \ldots + a_{jn}(t)x_n + g_{j}(t), \qquad j=1,\ldots,n\)

where \(n\) is a positive integer and \(a_{ji}(t),g_{j}(t)\) are arbitrary functions of the independent variable \(t\).

A first-order linear system of ODEs may be written in matrix form:

\(\frac{d}{dt} \begin{bmatrix} x_1 \\ x_2 \\ \vdots \\ x_n \end{bmatrix} = \begin{bmatrix} a_{11} & \ldots & a_{1n} \\ a_{21} & \ldots & a_{2 n} \\ \vdots & \ldots & \vdots \\ a_{n1} & & a_{n n} \end{bmatrix} \begin{bmatrix} x_1 \\ x_2 \\ \vdots \\ x_n \end{bmatrix} + \begin{bmatrix} g_1 \\ g_2 \\ \vdots \\ g_n \end{bmatrix} ,\)

or simply

\(\mathbf{\dot{x(t) = \mathbf{A}(t)\mathbf{x}(t) + \mathbf{g}(t)\).

### Homogeneous systems of differential equations
A linear system is said to be homogeneous if \(g_j(t)=0\) for each \(j\) and for all values of \(t\); otherwise, it is referred to as non-homogeneous. Homogeneous systems have the property that if \(\mathbf{x_1},\ldots ,\mathbf{x_p}\) are linearly independent solutions to the system, then any linear combination of these, \(C_1 \mathbf{x _1}+ \ldots + C_p \mathbf{x _p}\), is also a solution to the linear system where \(C_1, \ldots, C_p\) are constant.

The case where the coefficients \(a_{ji}(t)\) are all constant has a general solution: \(\mathbf{x} = C_1 \mathbf{v_1}e^{\lambda_1 t } + \ldots + C_n \mathbf{v_n}e^{\lambda_n t }\), where \(\lambda_i\) is an eigenvalue of the matrix \(\mathbf{A}\) with corresponding eigenvectors \(\mathbf{v}_i\) for \(1 \leq i \leq n\). This general solution only applies in cases where \(\mathbf{A}\) has n distinct eigenvalues, cases with fewer distinct eigenvalues must be treated differently.

## Linear independence of solutions
For an arbitrary system of ODEs, a set of solutions \(\mathbf{x_1}(t), \ldots ,\mathbf{x_n}(t)\) are said to be linearly independent if:

\(C_1\mathbf{x_1}(t) + \ldots + C_n \mathbf{x_n}(t) = 0 \quad \forall t\) is satisfied only for \(C_1 = \ldots = C_n=0\).

A second-order differential equation \(\ddot{x} = f(t,x,\dot{x})\) may be converted into a system of first order linear differential equations by defining \(y=\dot{x}\), which gives us the first-order system:

\(\begin{cases} \dot{x} & = & y \\ \dot{y} & = & f(t,x,y) \end{cases}\)

Just as with any linear system of two equations, two solutions may be called linearly independent if \(C_1 \mathbf{x}_1 + C_2 \mathbf{x}_2=\mathbf{0 }\) implies \(C_1 = C_2 = 0\), or equivalently that their Wronskian \(\begin{vmatrix} x_1 & x_ 2 \\ \dot{x}_ 1 & \dot{x}_ 2 \end{vmatrix}\) is non-zero. For second-order systems, any two solutions to a second-order ODE are called linearly independent if their Wronskian is non-zero as well.

## Overdetermination of systems of differential equations
Like any system of equations, a system of linear differential equations is said to be overdetermined if there are more equations than the unknowns. For an overdetermined system to have a solution, it needs to satisfy certain compatibility conditions. For example, consider the system:
\(\frac{\partial u}{\partial x_i} = f_i, 1 \le i \le m.\)
Then the necessary conditions for the system to have a solution are:
\(\frac{\partial f_i}{\partial x_k} - \frac{\partial f_k}{\partial x_i} = 0, 1 \le i, k \le m.\)

See also: Cauchy problem and Ehrenpreis's fundamental principle.

## Nonlinear system of differential equations
A system is considered nonlinear when one or more of the governing equations in the system is nonlinear, often leading to complex behavior, and are occasionally harder to predict. Fewer established techniques to find generalized solutions and prove smoothness for nonlinear systems of differential equations. Their computation often consists of linearization and/or numerical methods such as the Runge-Kutta methods.

Perhaps the most famous example of a nonlinear system of differential equations is the Navier–Stokes equations, used to model fluid dynamics and air flow. Unlike the linear case, the existence of a solution of a nonlinear system is a difficult problem and is considered one of the greatest unsolved problems in physics. (cf. Navier–Stokes existence and smoothness.)

Other examples of nonlinear systems of differential equations include the Lotka–Volterra equations, which are used to predict predator and prey populations over time, and the Lorenz system, a forecasting model that was used to demonstrate chaos theory.

## Differential system

A **differential system** is a means of studying a system of partial differential equations using geometric ideas such as differential forms and vector fields.

For example, the compatibility conditions of an overdetermined system of differential equations can be succinctly stated in terms of differential forms (i.e., for a form to be exact, it needs to be closed). See integrability conditions for differential systems for more.
