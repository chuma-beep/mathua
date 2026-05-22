> Content sourced from [Wikipedia](https://en.wikipedia.org/wiki/Order_of_operations) — CC BY-SA 4.0

# Nested parentheses in order of operations

## Special cases
### Unary minus sign
There are differing conventions concerning the unary operation '−' (usually pronounced "minus"). In written or printed mathematics, the expression −3\(^{2}\) is interpreted to mean −(3\(^{2}\)) {=} −9.BSARGS

In some applications and programming languages, notably Microsoft Excel (and other spreadsheet applications), unary operations have a higher priority than binary operations, that is, the unary minus has higher precedence than exponentiation, so in those languages −3\(^{2}\) will be interpreted as (−3)\(^{2}\) {=} 9.Microsoft 2005Berger This does not apply to the binary minus operation '−'; for example in Microsoft Excel the formulas =-2^2, =(-2)^2 and =0+-2^2 return 4, but the formulas =0-2^2 and =-(2^2) return −4.

### Mixed division and multiplication
There is no universal convention for interpreting an expression containing both division denoted by '÷' and multiplication denoted by '×'. Proposed conventions include assigning the operations equal precedence and evaluating them from left to right, or equivalently treating division as multiplication by the reciprocal and then evaluating in any order;Chrystal evaluating all multiplications first followed by divisions from left to right; or eschewing such expressions and instead always disambiguating them by explicit parentheses.Cajori

Beyond primary education, the symbol '÷' for division is seldom used, but is replaced by the use of algebraic fractions.Wu These are most explicitly and unambiguously written "vertically" with the numerator stacked above the denominator separated by a fraction bar. But they can also be written "horizontally" with the numerator and denominator separated by the slash symbol '/'.iso That is, expressions such as *a* ÷ *b* are avoided in favor of {sfrac|*a*|*b*} or *a* / *b*.

 Multiplication denoted by juxtaposition (also known as implied multiplication) creates a visual unit and is often given higher precedence than most other operations. In academic literature, when inline fractions are combined with implied multiplication without explicit parentheses, the multiplication is conventionally interpreted as having higher precedence than division, so that e.g. 1 / 2*n* is interpreted to mean 1 / (2 · *n*) rather than (1 / 2) · *n*.BSChrystalLennesStrogatz For instance, the manuscript submission instructions for the *Physical Review* journals directly state that multiplication has precedence over division,APS and this is also the convention observed in physics textbooks such as the *Course of Theoretical Physics* by Landau and Lifshitz and mathematics textbooks such as *Concrete Mathematics* by Graham, Knuth, and  Patashnik.GKP However, some authors recommend against expressions such as *a* / *bc*, preferring the explicit use of parenthesis *a* / (*bc*).Peterson

More complicated cases are more ambiguous. For instance, the notation 1 / 2*π*(*a* + *b*) could plausibly mean either
1 / [2*π* · (*a* + *b*)] or [1 / (2*π*)] · (*a* + *b*).FatemanCaspi Sometimes interpretation depends on context. The *Physical Review* submission instructions recommend against expressions of the form *a* / *b* / *c*; more explicit expressions (*a* / *b*) / *c* or *a* / (*b* / *c*) are unambiguous.APS


This ambiguity has been the subject of Internet memes such as "8 ÷ 2(2 + 2)", for which there are two conflicting interpretations: 8 ÷ [2 · (2 + 2)] = 1 and (8 ÷ 2) · (2 + 2) = 16.StrogatzHaelle Mathematics education researcher Hung-Hsi Wu points out that "one never gets a computation of this type in real life", and calls such contrived examples "a kind of Gotcha! parlor game designed to trap an unsuspecting person by phrasing it in terms of a set of unreasonably convoluted rules".Wu

### Serial exponentiation

If exponentiation is indicated by stacked symbols using superscript notation, the usual rule is to work from the top down:BSNIST
*a**b*\(^{*c*}\) {=} *a*(*b*\(^{*c*}\)),
which typically is not equal to (*a*\(^{*b*}\))\(^{*c*}\). This convention is useful because there is a property of exponentiation that (*a*\(^{*b*}\))\(^{*c*}\) = *a*\(^{*bc*}\), so it's unnecessary to use serial exponentiation for this.

However, when exponentiation is represented by an explicit symbol such as a caret (^) or arrow (↑), there is no common standard. For example, Microsoft Excel and computation programming language MATLAB evaluate *a*^*b*^*c* as (*a*\(^{*b*}\))\(^{*c*}\), but Google Search and Wolfram Alpha as *a*(*b*\(^{*c*}\)). Thus 4^3^2 is evaluated to 4,096 in the first case and to 262,144 in the second case.

