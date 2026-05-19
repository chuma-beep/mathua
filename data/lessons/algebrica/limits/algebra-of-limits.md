> Content sourced from [Algebrica](https://algebrica.org/algebra-of-limits/) — CC BY-NC 4.0

## Introduction

The definition of a [limit](<../limits/>) offers a framework for describing how a function \\(f(x)\\) approaches a specific value near a given point \\(x_0\\). This definition by itself is not enough for practical calculations. In most cases, we deal with limits when [functions](<../functions>) are added, multiplied, composed, or divided.

The algebra of limits consists of operational rules derived directly from the formal definition of a limit. These rules illustrate the structural compatibility between limits and standard algebraic operations. To illustrate these rules, it is assumed that \\( L \\) and \\( M \\) are [real numbers](<../properties-of-real-numbers/>) such that:

\\[\lim_{x \to x_0} f(x) = L \quad \text{and} \quad \lim_{x \to x_0} g(x) = M \\]

## Limit of a sum

When two functions approach finite values near a point, the sum of the functions approaches the sum of those values. Specifically, if \\( f(x) \\) remains close to \\( L \\) and \\( g(x) \\) remains close to \\( M \\), then their combined variation remains close to \\( L + M \\). Formally we have:

\\[\lim_{x \to x_0} \big(f(x) + g(x)\big) = L + M \\]

This property follows directly from the definition of a limit. For any tolerance around \\( L + M \\), the deviations of \\( f(x) \\) and \\( g(x) \\) can be controlled independently to ensure that their combined deviation remains within the prescribed bound. Consider, for instance, the following expressions that involve [sine](<../sine-function/>) and [cosine functions](<../cosine-function/>):

\\[f(x) = \frac{\sin x}{x} \quad g(x) = \frac{1 - \cos x}{x^2} \\]

and suppose we wish to compute the sum:

\\[\lim_{x \to 0} \big(f(x) + g(x)\big) \\]

Neither function is defined at \\( x = 0 \\), which prevents evaluating the limit by direct substitution. Nevertheless, two established results from the analysis provide:

\\[\lim_{x \to 0} \frac{\sin x}{x} = 1 \quad \lim_{x \to 0} \frac{1 - \cos x}{x^2} = \frac{1}{2} \\]

By the sum rule for limits, we conclude: \\[\lim_{x \to 0} \big(f(x) + g(x)\big) = 1 + \frac{1}{2} = \frac{3}{2} \\]

###### This example demonstrates the utility of the sum rule. Instead of analysing the combined expression directly, which would require substantial algebraic manipulation, the rule permits decomposition of the problem into two independent limits, each of which can be addressed separately.

## Limit of a difference

A similar argument applies to subtraction. If two functions approach \\(L\\) and \\(M\\), then their difference approaches \\(L - M\\). The algebraic properties of the real numbers ensure that subtraction is consistent with the limit process.

\\[\lim_{x \to x_0} \big(f(x) - g(x)\big) = L - M \\]

The proof parallels that of the sum because subtraction can be interpreted as the addition of the additive inverse. For instance consider the functions: \\[f(x) = \frac{1 - \cos x}{x^2} \quad g(x) = \frac{\sin^2 x}{2x^2} \\] Suppose we wish to compute the difference: \\[\lim_{x \to 0} \big(f(x) - g(x)\big) \\] Neither function is defined at \\( x = 0 \\), so direct substitution is not available. However, two known results give:

\\[\lim_{x \to 0} \frac{1 - \cos x}{x^2} = \frac{1}{2} \quad \lim_{x \to 0} \frac{\sin^2 x}{2x^2} = \frac{1}{2} \\]

The second result comes from the fact that \\( \lim_{x \to 0} \frac{\sin x}{x} = 1 \\). Using the difference rule for limits, we find:

\\[\lim_{x \to 0} \big(f(x) - g(x)\big) = \frac{1}{2} - \frac{1}{2} = 0 \\] At first, this result is not obvious from the combined expression: \\[\frac{1 - \cos x}{x^2} - \frac{\sin^2 x}{2x^2} \\]

because each term approaches \\( \frac{1}{2} \\), and finding their difference takes some careful work. The difference rule, like the sum rule, lets us break the problem into two simpler limits, which makes the calculation easier.

## Limit of a constant multiple

If a function gets closer to a value \\( L \\), multiplying it by a constant will also multiply the limit by that constant. This shows that limits behave linearly. In general, for any real constant \\( c \\), we have:

\\[\lim_{x \to x_0} c\,f(x) = cL \\]

The constant does not affect the limit but just changes the final value by scaling it. To analyse this case, suppose we want to find:

\\[\lim_{x \to 0} 3 \frac{\ln(1 + x)}{x} \\]

The [logarithmic function](<../logarithmic-function/>) is not defined at \\( x = 0 \\), so we can’t use direct substitution. However, a well-known result from analysis tells us:

\\[\lim_{x \to 0} \frac{\ln(1 + x)}{x} = 1 \\]

Using the constant multiple rule for limits, we get:

\\[\lim_{x \to 0} 3 \cdot \frac{\ln(1 + x)}{x} = 3 \cdot 1 = 3 \\]

## Limit of a product

If two functions get close to certain values near a point, their product also gets close to the product of those values. For example, if \\( f(x) \\) stays near \\( L \\) and \\( g(x) \\) stays near \\( M \\), then their product stays near \\( L \cdot M \\): \\[\lim_{x \to x_0} \big(f(x) \cdot g(x)\big) = L \cdot M \\]

The key idea behind this rule is that both factors can be controlled independently near \\( x_0 \\), and their combined effect remains bounded. To illustrate this, consider the following functions: \\[f(x) = \frac{e^x - 1}{x} \quad g(x) = \frac{\ln(1 + x)}{x} \\]

Neither function is defined at \\( x = 0 \\), so direct substitution is not available. However, both are remarkable limits whose values are known: \\[\lim_{x \to 0} \frac{e^x - 1}{x} = 1 \quad \lim_{x \to 0} \frac{\ln(1 + x)}{x} = 1 \\]

By the product rule, we can now compute the limit of their product: \\[\lim_{x \to 0} \frac{(e^x - 1)\ln(1 + x)}{x^2} = 1 \cdot 1 = 1 \\]

## Limit of a quotient

When dividing two or more limits, there is an important restriction to keep in mind. If \\( g(x) \\) gets close to a nonzero value \\( M \neq 0 \\), the quotient acts normally near \\( x_0 \\). But if the denominator’s limit is zero, the result can become unpredictable. Assuming \\( M \neq 0 \\), we have:

\\[\lim_{x \to x_0} \frac{f(x)}{g(x)} = \frac{L}{M} \\]

The main idea is that because \\( g(x) \\) stays close to a nonzero number near \\( x_0 \\), it does not get close to zero there. This ensures the quotient does not give rise to any division-by-zero issues near \\( x_0 \\), and the limit behaves as expected. Let’s look at the following functions to analyse this case:

\\[f(x) = \frac{e^x - 1}{x} \quad g(x) = \frac{x^2 + 1}{x + 1} \\]

Suppose we want to find the quotient: \\[\lim_{x \to 0} \frac{f(x)}{g(x)} \\]

The function \\( f(x) \\) is not defined at \\( x = 0 \\), so we cannot use direct substitution for the numerator. However we know that: \\[\lim_{x \to 0} \frac{e^x - 1}{x} = 1 \\]

For the denominator, we can use direct substitution because \\( g(x) \\) is defined at \\( x = 0 \\): \\[\lim_{x \to 0} \frac{x^2 + 1}{x + 1} = \frac{0 + 1}{0 + 1} = 1 \\] Since the denominator’s limit is \\( M = 1 \\) and not zero, we can use the quotient rule and get: \\[\lim_{x \to 0} \frac{f(x)}{g(x)} = \frac{1}{1} = 1 \\]

###### This example shows why it is important to check that the denominator’s limit is not zero before using the quotient rule. In this case, \\( g(x) \\) stays close to \\( 1 \\) near \\( x = 0 \\), so there is no risk of dividing by zero.

## Limits of powers and polynomials

When we have repeated multiplication, we encounter limits on powers. If \\(f(x)\\) approaches \\(L\\), then for any positive integer \\(n\\) we have:

\\[\lim_{x \to x_0} (f(x))^n = L^n \\]

This property means you can find the limit of a [polynomial](<../polynomials/>) by plugging the limiting value into the polynomial. Since polynomials use only sums and products, they follow the same rules for limits as these operations. Consider the following function limit: \\[\lim_{x \to 0} \left(\frac{e^x - 1}{x}\right)^4 \\]

The function is not defined at \\( x = 0 \\), so direct substitution is not available. However, as a [remarkable limit](<../remarkable-limits>), we know that:

\\[\lim_{x \to 0} \frac{e^x - 1}{x} = 1 \\]

If we use the power rule with \\( n = 4 \\), we get:

\\[\lim_{x \to 0} \left(\frac{e^x - 1}{x}\right)^4 = 1^4 = 1 \\]

The power rule allows us to avoid expanding the expression directly, which would make the calculation much harder. Instead, once the limit of the base function is known, we just raise that value to the required power in a single step: \\[\left(\frac{e^x - 1}{x}\right)^4 \to 1^4 = 1 \\]

## Limit of a Composition

We now consider the case of function composition. Given two functions \\( \varphi \\) and \\( f \\), their composition \\( \varphi(f(x)) \\) consists of applying \\( f \\) first and then \\( \varphi \\) to the result. Suppose that:

\\[\lim_{x \to x_0} f(x) = L \\]

Now, if \\( \varphi \\) is continuous at \\( L \\), the limit can be taken through the function: \\[\lim_{x \to x_0} \varphi(f(x)) = \varphi(L) \\]

This property connects the algebra of limits with [continuity](<../continuous-functions/>). The continuity of \\( \varphi \\) ensures that small variations in the input near \\( L \\) produce small variations in the output. Without continuity, this rule cannot be guaranteed.

For instance, consider the following function: \\[f(x) = \frac{\ln(1 + x)}{x} \quad \varphi(t) = \sqrt{t} \\]

Suppose we want to find the limit: \\[\lim_{x \to 0} \sqrt{\frac{\ln(1 + x)}{x}} \\]

The function \\( f(x) \\) is not defined at \\( x = 0 \\), so direct substitution is not available. However, this is one of the remarkable limits, and its value is known to be: \\[\lim_{x \to 0} \frac{\ln(1 + x)}{x} = 1 \\]

Since \\( \varphi(t) = \sqrt{t} \\) is continuous at \\( t = 1 \\), the limit can be taken through the outer function: \\[\lim_{x \to 0} \sqrt{\frac{\ln(1 + x)}{x}} = \sqrt{1} = 1 \\]

###### This example shows that the composition rule reduces the problem to two separate steps: first identifying the limit of the inner function, and then evaluating the outer function at that value. This procedure is justified only if \\( \varphi \\) is continuous at \\( L = 1 \\).

## Summary

|   
---|---  
Sum| \\[\lim_{x \to x_0} \big(f(x) + g(x)\big) = L + M \\]  
Difference| \\[\lim_{x \to x_0} \big(f(x) - g(x)\big) = L - M \\]  
Constant multiple| \\[\lim_{x \to x_0} c\,f(x) = cL \\]  
Product| \\[\lim_{x \to x_0} \big(f(x) \cdot g(x)\big) = L \cdot M \\]  
Quotient| \\[\lim_{x \to x_0} \frac{f(x)}{g(x)} = \frac{L}{M} \quad M \neq 0 \\]  
Power| \\[\lim_{x \to x_0} \big(f(x)\big)^n = L^n \\]  
Composition| \\[\lim_{x \to x_0} \varphi(f(x)) = \varphi(L) \quad \varphi \text{ continuous at } L \\]  
  
## Selected references

  * **University of Washington, M. Greenberg**. [Basic Theorems About Limits](https://sites.math.washington.edu//~greenber/MATH327-LimitTheorems.pdf)

  * **University of Wisconsin, J. Robbin**. [Calculus: Lecture Notes](https://people.math.wisc.edu/~angenent/Free-Lecture-Notes/free221.pdf)

  * **University of California Davis, J. K. Hunter**. [Limits of Functions](https://www.math.ucdavis.edu/~hunter/m125a/intro_analysis_ch2.pdf)


Limits

Limits describe a function’s behavior near a point.

5.9k

[Limits](https://algebrica.org/limits/)

1.6k

[Squeeze Theorem](https://algebrica.org/squeeze-theorem/)

4.1k

[Remarkable Limits](https://algebrica.org/remarkable-limits/)

0 comments[](https://github.com/antoniolupetti/algebrica/blob/main/limits/remarkable-limits.md?plain=1)

3.8k

[Asymptotes](https://algebrica.org/asymptotes/)

1.9k

[Indeterminate Forms of Limits](https://algebrica.org/indeterminate-forms/)

5.1k

[Little-o Notation](https://algebrica.org/little-o-notation/)

1.8k

[Big O Notation](https://algebrica.org/big-o-notation/)
