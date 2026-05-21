> Content sourced from [Algebrica](https://algebrica.org/logarithmic-inequalities/) — CC BY-NC 4.0

## Introduction

Logarithmic inequalities are inequalities that involve one or more [logarithmic](<../logarithms>) expressions, in which the unknown \\(x\\) appears either in the argument of the logarithm or, in some cases, in the base itself. Before tackling logarithmic inequalities, it is essential to have a solid understanding of logarithms, their fundamental properties, and the standard methods used to solve [logarithmic equations](<../logarithmic-inequalities/>), as these tools are crucial for approaching and resolving such inequalities correctly.


Logarithmic inequalities may have the general form:

\\[\log_af(x) \lesseqgtr \log_ag(x) \\]

  * \\(a\\) is the base of the logarithm, with (a > 0) and \\(a \neq 1.\\)
  * \\(f(x)\\) and \\(g(x)\\) are algebraic expressions depending on the variable \\(x.\\)
  * The symbol \\(\lesseqgtr\\) denotes one of the relations \\(\le\\), \\(=\\), or \\(\ge.\\)


Moreover, for the inequality to be well defined, the following conditions must be satisfied:

\\[\begin{cases} f(x) > 0\\\\[0.6em] g(x) > 0\\\\[0.6em] \end{cases} \\]


An inequality that contains a logarithmic expression in which the unknown variable does not appear in either the argument or the base of the logarithm is not a logarithmic inequality. In other words, an inequality such as: \\[\log_3(9) - 3x > 0 \\] is not a logarithmic inequality, since the logarithmic term is a constant. By contrast,  
\\[\log_3(9x) - 3x > 0 \\] is a logarithmic inequality, because the variable \\(x\\) appears inside the argument of the logarithm and directly affects its domain and behavior.

## How to solve logarithmic inequalities

The solution process for logarithmic inequalities is general, however, for explanatory convenience, let us consider the following inequality: \\[\log_a f(x) \ge \log_a g(x) \\]

The procedure can be structured into four fundamental steps:

  * Determine the [domain](<../determining-the-domain-of-a-function/>) of the inequality by imposing the admissibility conditions. The arguments of all logarithmic expressions must be strictly positive, and the base must satisfy: \\[\begin{cases} f(x) > 0 \\\\[0.6em] g(x) > 0 \\\\[0.6em] \end{cases} \\]

  * The base must satisfy: \\[\begin{cases} a > 0 \\\\[0.6em] a \neq 1 \\\\[0.6em] \end{cases} \\]

  * If \\(a > 1\\), the logarithmic function is [increasing](<../increasing-and-decreasing-functions/>) and the inequality preserves its direction when the logarithms are removed.

  * If \\(0 < a < 1\\), the logarithmic function is decreasing and the direction of the inequality must be reversed when eliminating the logarithms.

  * Use the monotonicity of the logarithmic function to remove the logarithms and reduce the problem to an equivalent algebraic inequality involving \\(f(x)\\) and \\(g(x)\\).

  * Solve the resulting algebraic inequality and intersect the solution set with the domain previously determined, discarding any values that do not satisfy the original logarithmic conditions.


To explain the role played by the base of the logarithm, let us recall the behavior of the [logarithmic function](<../logarithmic-function/>) when \\(0 < a < 1\\). From its graph, we observe that the function is strictly decreasing over its entire domain, with a vertical asymptote along the \\(y\\)-axis:

![Graph of the logarithmic function with base between zero and one.](/diagrams/algebrica/logharithm-5-1.png)

##### The dashed curve represents the logarithmic function with base \\(a > 1\\). In this case, the function is strictly increasing. In both cases, when \\(x = 1\\), the value of the logarithmic function is \\(0\\), and the graphs intersect at the point \\((1,0)\\).

## Example 1

Consider the logarithmic inequality: \\[\log_{\frac{1}{2}}(x + 3) \ge \log_{\frac{1}{2}}(2x - 1) \\]

We begin by determining the domain of the inequality. The arguments of the logarithms must be strictly positive: \\[\begin{cases} x + 3 > 0 \\\\[0.6em] 2x - 1 > 0 \end{cases} \\]

Using a graphical representation and considering the solution intervals of the [linear inequalities](<../linear-inequalities/>) in the previous system, we find that their intersection, which determines the domain of the logarithmic inequality, is precisely given by \\(x > 1/2\\).

| \\[-3\\]| \\[\frac{1}{2}\\]|   
---|---|---|---  
| | |   
| | |   
| | |   
| | |   
  
Therefore, the domain \\(D\\) of the original inequality is given by the following interval: \\[\left(-\frac{1}{2}, +\infty\right) \\]


Next, we analyze the base of the logarithm. Since the base satisfies \\[0 < \frac{1}{2} < 1 \\] the logarithmic function is strictly decreasing. As a consequence, when the logarithms are removed, the direction of the inequality must be reversed. Therefore, the given inequality: \\[\log_{\frac{1}{2}}(x + 3) \ge \log_{\frac{1}{2}}(2x - 1) \\] is equivalent to: \\[x + 3 \le 2x - 1 \\]


We now solve the resulting algebraic inequality: \\[x + 3 \le 2x - 1 \quad \rightarrow \quad x \ge 4 \\]

Finally, we intersect this result with the domain previously determined. Since the domain requires \\(x > \frac{1}{2}\\), the condition \\(x \ge 4\\) is admissible.

Hence, the solution set of the logarithmic inequality is: \\[x \ge 4 \\]

## Example 2

Consider the logarithmic inequality:

\\[\log_{\frac{1}{2}}(x+1) > \log_2(2-x) \\]

We begin by determining the [domain](<../determining-the-domain-of-a-function/>). The arguments of the logarithms must be strictly positive, hence:

\\[\begin{cases} x+1 > 0 \\\\[0.6em] 2-x > 0 \\\\[0.6em] \end{cases} \\]

Using a graphical representation and considering the solution intervals of the linear inequalities in the previous system, we find that their intersection is given by \\( -1 < x < 2\\).

| \\[-1\\]| \\[2\\]|   
---|---|---|---  
| | |   
| | |   
| | |   
| | |   
  
Therefore, the domain \\(D\\) of the original inequality is given by the following interval: \\[(-1, 2) \\]


Next, we rewrite the logarithm with base \\(\tfrac{1}{2}\\) in terms of base \\(2\\). Since \\(\tfrac{1}{2} = 2^{-1}\\), we have \\[\log_{\frac{1}{2}}(x+1) = \frac{\log_2(x+1)}{\log_2(\frac{1}{2})} = -\log_2(x+1) \\]

Substituting into the original inequality, we obtain:

\\[\log_2(x+1) < -\log_2(2-x) \\]

Bringing all logarithmic terms to the same side and applying the properties of logarithms, we get: \\[\log_2(x+1) + \log_2(2-x) < 0 \\]

which, by the sum property of logarithms, stating that the [sum of two logarithms](<../logarithms>) equals the logarithm of their product, allows us to rewrite the expression as follows:

\\[\log_2!\bigl((x+1)(2-x)\bigr) < 0 \\]

Since the logarithmic function with base \\(2>1\\) is strictly increasing, this inequality is equivalent to \\[0 < (x+1)(2-x) < 1 \\]

Within the domain \\((-1,2)\\), the product \\((x+1)(2-x)\\) is always positive, so it suffices to solve:

\\[(x+1)(2-x) < 1 \\]

Expanding and simplifying, we obtain \\[x^2 - x - 1 > 0 \\]

The associated [quadratic equation](<../quadratic-equations>) has roots: \\[x = \frac{1 \pm \sqrt{5}}{2} \\]

The inequality is satisfied outside the interval determined by these roots. Intersecting this result with the domain \\((-1,2)\\), we finally obtain the solution set: \\[-1 < x < \frac{1-\sqrt{5}}{2}\\]
