> Content sourced from [Algebrica](https://algebrica.org/logarithmic-function/) — CC BY-NC 4.0

## Introduction

The [logarithmic](<../logarithms>) function is the inverse of the [exponential function](<../exponential-function>). Therefore, its [domain](<../determining-the-domain-of-a-function/>) and range are inverted compared to the exponential function. A logarithmic function is defined as a function of the form:

\\[y = \log{_a}x \quad \text{with} \quad a \gt 0 \quad a \neq 1 \quad \forall x \in \mathbb{R}^+ \\]

There are two cases to consider when examining the logarithmic function \\( \log_a(x) \\).

  * If the base \\( a > 1 \\), the function is increasing, meaning it grows as \\( x \\) increases.

  * If \\( 0 < a < 1 \\), the function is decreasing, meaning it diminishes as \\( x \\) increases.


![](/diagrams/algebrica/logharithm-6.png)

In both cases, the only point where the logarithmic function takes the value \\( 0 \\) is at \\( x = 1 \\). This is because: \\(\log_a(1) = 0\\). This result is true because by definition, by the property of [powers](<../powers>), we have \\( a^0 = 1 \\).

##### The value of the base, whether it is greater than 1 or between 0 and 1, plays a crucial role in many applications. In particular, when dealing with [logarithmic inequalities](<../logarithmic-inequalities/>), a base between 0 and 1 requires reversing the direction of the inequality, due to the decreasing nature of the logarithmic function.

## Properties

  * Domain: \\( (0, +\infty) \\).
  * Range: \\( \mathbb{R} \\).
  * Roots: \\(x=1\\).
  * The logarithmic function is strictly [increasing](<../increasing-and-decreasing-functions/>) on \\((0,+\infty)\\) when the base satisfies \\(a>1\\). If \\(0<a<1\\), the function is strictly decreasing on the same interval.
  * The function is neither even nor odd, because it is not defined for negative values of \\( x \\).
  * The function is [continuous](<../continuous-functions/>) on \\( (0, +\infty) \\).
  * The function is differentiable on its entire domain, with derivative \\( f’(x) = 1/x \\).
  * The function has no [maximum or minimum](<../maximum-minimum-and-inflection-points/>) on its domain.


## Limits, derivatives, and integrals

A fundamental [limit](<../limits/>) involving the logarithmic function describes its behavior near zero and at infinity and plays an important role in the study of logarithmic functions. For a logarithmic function with base \\(a>1\\), as the variable approaches zero from the right, the value of the logarithm decreases without bound, while it grows without bound as the variable tends to infinity. This behavior is formally expressed by the following limits:

\\[1. \quad \lim_{x \to 0^+} \log_a x = -\infty \\] \\[2. \quad \lim_{x \to +\infty} \log_a x = +\infty \\]

When the logarithm base is \\( 0 < a < 1 \\), we have:

\\[3. \quad \lim_{x \to 0^+} \log_a x = +\infty \\] \\[4. \quad \lim_{x \to +\infty} \log_a x = -\infty\\]


The [derivative](<../derivatives>) of the logarithmic function with base \\(a\\) can be obtained by applying the standard differentiation rules for logarithms. In particular, for \\(x>0\\) and \\(a>0\\) with \\(a\neq 1\\), the derivative of \\(\log_a(x)\\) with respect to \\(x\\) is given by: \\[5\. \quad \frac{d}{dx}\,\log_a(x) = \frac{1}{x\,\ln(a)} \\]


The [indefinite integral](<../indefinite-integrals/>) of the logarithmic function with base \\(a\\) can be derived by recalling the relationship between logarithms with different bases. Since: \\[\log_a(x) = \frac{\log(x)}{\log(a)} \\]

the integration can be reduced to the integral of the natural logarithm. Using standard techniques from calculus, in particular [integration by parts](<../integration-by-parts/>), we obtain the following result: \\[6\. \quad \int \frac{\log(x)}{\log(a)} \, dx = \frac{x\,(\log(x)-1)}{\log(a)} + c \\]

## Natural logarithm

More generally, considering the function \\( y = \ln(x) \\), the natural logarithm, we have: \\[7. \quad \ln(x) = \log_e(x)\\]

Using the change of base formula, this can be rewritten as: \\[8. \quad \ln(x) = \frac{\log_a(x)}{\log_a(e)}\\]


For the natural logarithm \\(\ln(x)\\), which is defined for all \\(x>0\\), the rate of change with respect to \\(x\\) is particularly simple. Indeed, the derivative of \\(\ln(x)\\) is given by: \\[9\. \quad \frac{d}{dx}\,\ln(x) = \frac{1}{x} \\]

##### This result reflects the fact that the natural logarithm is the inverse function of the [exponential function](<../exponential-function>) \\(e^x\\), and it explains why \\(\ln(x)\\) grows slowly as \\(x\\) increases.


An important integral representation of the natural logarithm involves the [definite integral](<../definite-integrals>) of \\(\ln(t)\\) over an interval starting at the origin. Specifically, consider the integral \\[10 . \quad \int_{0}^{x} \ln(t) \, dt \\] which is defined for \\(x>0\\). Although the integrand \\(\ln(t)\\) is not defined at \\(t=0\\), the integral is interpreted as an improper integral. Its value can be computed using integration by parts, yielding: \\[\int \ln(t)\,dt = t\ln(t) - t + c \\] Applying this antiderivative and taking the limit as the lower bound approaches zero, we obtain: \\[\int_{0}^{x} \ln(t)\,dt = \bigl[t\ln(t) - t\bigr]_{0}^{x} = x\ln(x) - x + 1 \\]


A fundamental definite integral involving the natural logarithm is \\[11\. \quad \int_{0}^{1} \ln(x) \, dx = -1 \\] This integral is interpreted as an improper integral, since the natural logarithm is not defined at \\(x=0\\). Its value can be obtained by using the antiderivative \\(x\ln(x) - x\\) and evaluating the limit at the lower bound. The result shows that, on the interval \\((0,1)\\), the negative values of \\(\ln(x)\\) dominate, leading to a net negative area equal to \\(-1\\).

The logarithmic function models the growth of information in binary decisions. Each step in binary search answers one yes/no question, which corresponds to one binary bit. This is a deep reason why logarithms are fundamental in computer science and information theory.


## Binary search and the logarithmic function

Imagine we want to find an item in a sorted list of \\( N \\) elements. If we check the elements one by one, and the target is in the \\( n \\)-th position, we may need up to \\( n \\) steps to find it. This approach is known as linear search, and its worst-case time complexity is:

\\[\mathcal{O}(n) \\]

However, linear search is inefficient for large lists. Can we do better? Yes — by using an algorithm that reduces the search space at each step. This leads us to **binary search** , which has logarithmic complexity:

\\[\mathcal{O}(\log n) \\]

Binary search is an efficient algorithm used to find a target element in a sorted list. The core idea is simple: at each step, the algorithm divides the list in half and keeps only the half that may contain the target. This process continues until the element is found or the list can no longer be divided.

At each step, the search space is divided by 2. This means that if the initial list has \\( n \\) elements, the number of steps required to complete the search is proportional to the number of times we can divide \\( n \\) by 2 before reaching 1. This is exactly the definition of the base-2 logarithm:

\\[\text{Number of steps} \approx \log_2 n \\]

Thus, the time complexity of binary search is:

\\[\mathcal{O}(\log n) \\]

This connection shows how logarithmic functions naturally describe processes that involve repeated halving, such as decision trees, divide-and-conquer algorithms, and exponential decay.


If a list contains 1,000,000 elements, binary search can find any element in at most:

\\[\log_2(1,000,000) \approx 19.9 \\]

So, fewer than 20 steps are needed—an enormous efficiency gain compared to linear search, which may require up to 1,000,000 steps.

##### The elements must be sorted according to an ordering criterion. Binary search only works on sequences that are ordered—typically in ascending order. Without a defined ordering, the algorithm cannot determine which half of the list to discard at each step.
