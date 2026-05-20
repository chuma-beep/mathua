> Content sourced from [Wikipedia](https://en.wikipedia.org/wiki/Order_of_operations) — CC BY-SA 4.0

# Order of operations without exponents

In mathematics and computer programming, the **order of operations** is a collection of conventions about which arithmetic operations to perform first in order to evaluate a given mathematical expression.

These conventions are formalized with a ranking of the operations. The rank of an operation is called its **precedence**, and an operation with a *higher* precedence is performed before operations with *lower* precedence. Calculators generally perform operations with the same precedence from left to right, but some programming languages and calculators adopt different conventions.

For example, multiplication is granted a higher precedence than addition, and it has been this way since the introduction of modern algebraic notation. Thus, in the expression 1 + 2 × 3, the multiplication is performed before addition, and the expression has the value 1 + (2 × 3) , and not (1 + 2) × 3 . When exponents were introduced in the 16th and 17th centuries, they were given precedence over both addition and multiplication and placed as a superscript to the right of their base. Thus 3 + 5<sup>2</sup>  and 3 × 5<sup>2</sup> .

These conventions exist to avoid notational ambiguity while allowing notation to remain brief. Where it is desired to override the precedence conventions, or even simply to emphasize them, parentheses ( ) can be used. For example, (2 + 3) × 4  forces addition to precede multiplication, while (3 + 5)<sup>2</sup>  forces addition to precede exponentiation. If multiple pairs of parentheses are required in a mathematical expression (such as in the case of nested parentheses), the parentheses may be replaced by other types of brackets to avoid confusion, as in [2 × (3 + 4)] − 5 .

These conventions are meaningful only when the usual notation (called infix notation) is used. When functional or Polish notation are used for all operations, the order of operations results from the notation itself.

## Conventional order
The order of operations, that is, the order in which the operations in an expression are usually performed, results from a convention adopted  throughout mathematics, science, technology and many computer programming languages. It is summarized as:
# Parentheses
# Exponentiation
# Multiplication and division
# Addition and subtraction
This means that to evaluate an expression, one first evaluates any sub-expression inside parentheses, working inside to outside if there is more than one set. Whether inside parentheses or not, the operation that is higher in the above list should be applied first. Operations of the same precedence are conventionally evaluated from left to right.

If each division is replaced with multiplication by the reciprocal (multiplicative inverse) then the associative and commutative laws of multiplication allow the factors in each term to be multiplied together in any order. Sometimes multiplication and division are given equal precedence, or sometimes multiplication is given higher precedence than division; see  below. If each subtraction is replaced with addition of the opposite (additive inverse), then the associative and commutative laws of addition allow terms to be added in any order.

The radical symbol  which signifies a square root is traditionally extended by a bar (the vinculum) over the radicand; this avoids the need for parentheses around the radicand. Other functions use parentheses around the input to avoid ambiguity. The parentheses can be omitted if the input is a single numerical variable or constant, as in the case of sin *x*  and sin *π* .  Traditionally this convention extends to monomials; thus, sin 3*x*  and even sin  sin(*xy*)}}, but sin *x* + *y* , because *x* + *y* is not a monomial. However, this convention is not universally understood, and some authors prefer explicit parentheses. Some calculators and programming languages require parentheses around function inputs, while others do not.

 Parentheses and alternate symbols of grouping can be used to override the usual order of operations or to make the intended order explicit. Grouped symbols can be treated as a single expression.

### Examples
Multiplication before addition:
\(1 + 2 \times 3 = 1 + 6 = 7 .\)

Parenthetical subexpressions are evaluated first:
\((1 + 2) \times 3 = 3 \times 3 = 9 .\)

Exponentiation before multiplication, multiplication before subtraction:
\(1 - 2 \times 3 ^ 4 = 1 - 2 \times 81 = 1 - 162 = -161 .\)

When an expression is written as a superscript, the superscript is considered to be grouped by its position above its base:
\(1+2^{3+4} = 1+2^7 = 1+128 = 129 .\)

The operand of a root symbol is determined by the overbar:
\(\sqrt{1 + 3} + 5 = \sqrt 4 + 5 = 2 + 5 = 7.\)

A horizontal fractional line forms two grouped subexpressions, one above divided by another below:
\(\frac{1 + 2}{3 + 4} + 5 = \frac{3}{7} + 5.\)

Parentheses can be nested, and should be evaluated from the inside outward. For legibility, outer parentheses can be made larger than inner parentheses. Alternately, other grouping symbols, such as curly braces { } or square brackets [ ], are sometimes used along with parentheses ( ). For example:
\(\bigl[ (1 + 2) \div (3 + 4) \bigr] + 5 = (3 \div 7) + 5\)

## Special cases
### Unary minus sign
There are differing conventions concerning the unary operation '−' (usually pronounced "minus"). In written or printed mathematics, the expression −3<sup>2</sup> is interpreted to mean −(3<sup>2</sup>) .

In some applications and programming languages, notably Microsoft Excel (and other spreadsheet applications), unary operations have a higher priority than binary operations, that is, the unary minus has higher precedence than exponentiation, so in those languages −3<sup>2</sup> will be interpreted as (−3)<sup>2</sup> . This does not apply to the binary minus operation '−'; for example in Microsoft Excel the formulas =-2^2, =(-2)^2 and =0+-2^2 return 4, but the formulas =0-2^2 and =-(2^2) return −4.

### Mixed division and multiplication
There is no universal convention for interpreting an expression containing both division denoted by '÷' and multiplication denoted by '×'. Proposed conventions include assigning the operations equal precedence and evaluating them from left to right, or equivalently treating division as multiplication by the reciprocal and then evaluating in any order; evaluating all multiplications first followed by divisions from left to right; or eschewing such expressions and instead always disambiguating them by explicit parentheses.

Beyond primary education, the symbol '÷' for division is seldom used, but is replaced by the use of algebraic fractions. These are most explicitly and unambiguously written "vertically" with the numerator stacked above the denominator separated by a fraction bar. But they can also be written "horizontally" with the numerator and denominator separated by the slash symbol '/'. That is, expressions such as *a* ÷ *b* are avoided in favor of  or *a* / *b*.

 Multiplication denoted by juxtaposition (also known as implied multiplication) creates a visual unit and is often given higher precedence than most other operations. In academic literature, when inline fractions are combined with implied multiplication without explicit parentheses, the multiplication is conventionally interpreted as having higher precedence than division, so that e.g. 1 / 2*n* is interpreted to mean 1 / (2 · *n*) rather than (1 / 2) · *n*. For instance, the manuscript submission instructions for the *Physical Review* journals directly state that multiplication has precedence over division, and this is also the convention observed in physics textbooks such as the *Course of Theoretical Physics* by Landau and Lifshitz and mathematics textbooks such as *Concrete Mathematics* by Graham, Knuth, and  Patashnik. However, some authors recommend against expressions such as *a* / *bc*, preferring the explicit use of parenthesis *a* / (*bc*).

More complicated cases are more ambiguous. For instance, the notation 1 / 2*π*(*a* + *b*) could plausibly mean either
1 / [2*π* · (*a* + *b*)] or [1 / (2*π*)] · (*a* + *b*). Sometimes interpretation depends on context. The *Physical Review* submission instructions recommend against expressions of the form *a* / *b* / *c*; more explicit expressions (*a* / *b*) / *c* or *a* / (*b* / *c*) are unambiguous.


This ambiguity has been the subject of Internet memes such as "8 ÷ 2(2 + 2)", for which there are two conflicting interpretations: 8 ÷ [2 · (2 + 2)] = 1 and (8 ÷ 2) · (2 + 2) = 16. Mathematics education researcher Hung-Hsi Wu points out that "one never gets a computation of this type in real life", and calls such contrived examples "a kind of Gotcha! parlor game designed to trap an unsuspecting person by phrasing it in terms of a set of unreasonably convoluted rules".

### Serial exponentiation

If exponentiation is indicated by stacked symbols using superscript notation, the usual rule is to work from the top down:
*a*<sup>*b*<sup>*c*</sup></sup> ,
which typically is not equal to (*a*<sup>*b*</sup>)<sup>*c*</sup>. This convention is useful because there is a property of exponentiation that (*a*<sup>*b*</sup>)<sup>*c*</sup> = *a*<sup>*bc*</sup>, so it's unnecessary to use serial exponentiation for this.

However, when exponentiation is represented by an explicit symbol such as a caret (^) or arrow (↑), there is no common standard. For example, Microsoft Excel and computation programming language MATLAB evaluate *a*^*b*^*c* as (*a*<sup>*b*</sup>)<sup>*c*</sup>, but Google Search and Wolfram Alpha as *a*<sup>(*b*<sup>*c*</sup>)</sup>. Thus 4^3^2 is evaluated to 4,096 in the first case and to 262,144 in the second case.

## Mnemonics


Mnemonic acronyms are often taught in primary schools to help students remember the order of operations. The acronym ***PEMDAS***, which stands for **P**arentheses, **E**xponents, **M**ultiplication/**D**ivision, **A**ddition/**S**ubtraction, is common in the United States and France. Sometimes the letters are expanded into words of a mnemonic sentence such as "Please Excuse My Dear Aunt Sally". The United Kingdom and other Commonwealth countries may use ***BODMAS*** (or sometimes ***BOMDAS***), standing for **B**rackets, **O**f, **D**ivision/**M**ultiplication, **A**ddition/**S**ubtraction, with "of" meaning fraction multiplication. Sometimes the **O** is instead expanded as **O**rder, meaning exponent or root, or replaced by **I** for **I**ndices in the alternative mnemonic ***BIDMAS***. In Canada and New Zealand ***BEDMAS*** is common.

These mnemonics may be misleading when written this way. For example, misinterpreting any of the above rules to mean "addition first, subtraction afterward" would incorrectly evaluate the expression \(a - b + c\) as \(a - (b + c)\), while the correct evaluation is \((a - b) + c\). These values are different when \(c\ne 0\).

In Germany, the convention is simply taught as , "dot operations before line operations" referring to the graphical shapes of the symbols (multiplication), (division), (addition), and (subtraction). This avoids the potential for the above misunderstanding.

Mnemonic acronyms have been criticized for not developing a conceptual understanding of the order of operations, and not addressing student questions about its purpose or flexibility. Students learning the order of operations via mnemonic acronyms routinely make mistakes, as do some pre-service teachers. Even when students correctly learn the acronym, a disproportionate focus on memorization of trivia crowds out substantive mathematical content. The acronym's procedural application does not match experts' intuitive understanding of mathematical notation: mathematical notation indicates groupings in ways other than parentheses or brackets and a mathematical expression is a tree-like hierarchy rather than a linearly "ordered" structure; furthermore, there is no single order by which mathematical expressions must be simplified or evaluated and no universal canonical simplification for any particular expression, and experts fluently apply valid transformations and substitutions in whatever order is convenient, so learning a rigid procedure can lead students to a misleading and limiting understanding of mathematical notation.

## Calculators

Different calculators follow different orders of operations. Many simple calculators without a stack implement chain input, working in button-press order without any priority given to different operations, give a different result from that given by more sophisticated calculators. For example, on a simple calculator, typing 1 + 2 × 3 = yields 9, while a more sophisticated calculator will use a more standard priority, so typing 1 + 2 × 3 = yields 7.

Calculators may associate exponents to the left or to the right. For example, the expression *a*^*b*^*c* is interpreted as *a*<sup>(*b*<sup>*c*</sup>)</sup> on the TI-92 and the TI-30XS MultiView in "Mathprint mode", whereas it is interpreted as (*a*<sup>*b*</sup>)<sup>*c*</sup> on the TI-30XII and the TI-30XS MultiView in "Classic mode".

An expression like 1/2*x* is interpreted as 1/(2*x*) by TI-82, as well as many modern Casio calculators (configurable on some like the fx-9750GIII), but as (1/2)*x* by TI-83 and every other TI calculator released since 1996, as well as by all Hewlett-Packard calculators with algebraic notation. While the first interpretation may be expected by some users due to the nature of implied multiplication, the latter is more in line with the rule that multiplication and division are of equal precedence.

When the user is unsure how a calculator will interpret an expression, parentheses can be used to remove the ambiguity.

Order of operations arose due to the adaptation of infix notation in standard mathematical notation, which can be notationally ambiguous without such conventions, as opposed to postfix notation or prefix notation, which do not need orders of operations. Hence, calculators utilizing reverse Polish notation (RPN) using a stack to enter expressions in the correct order of precedence do not need parentheses or any possibly model-specific order of execution.

## Programming languages
Most programming languages use precedence levels that conform to the order commonly used in mathematics, though others, such as APL, Smalltalk, Occam and Mary, have no operator precedence rules (in APL, evaluation is strictly right to left; in Smalltalk, it is strictly left to right).

Furthermore, because many operators are not associative, the order within any single level is usually defined by grouping left to right so that 16/4/4 is interpreted as (16/4)/4  rather than 16/(4/4) ; such operators are referred to as "left associative". Exceptions exist; for example, languages with operators corresponding to the cons operation on lists usually make them group right to left ("right associative"), e.g. in Haskell, 1:2:3:4:[] == 1:(2:(3:(4:[]))) == [1,2,3,4].

Dennis Ritchie, creator of the C language, said of the precedence in C (shared by programming languages that borrow those rules from C, for example, C++, Perl and PHP) that it would have been preferable to move the bitwise operators above the comparison operators. Many programmers have become accustomed to this order, but more recent popular languages like Python and Ruby do have this order reversed. The relative precedence levels of operators found in many C-style languages are as follows:

{| class="wikitable"
|1 || ()   []   ->   .   :: || Function call, scope, array/member access
|-
|2 || !   ~   -   +   *   &   sizeof   *type cast*   ++   --   || (most) unary operators, sizeof and type casts (right to left)
|-
|3 || *   /   % MOD || Multiplication, division, modulo
|-
|4 || +   - || Addition and subtraction
|-
|5 || <<   >> || Bitwise shift left and right
|-
|6 || <   <=   >   >= || Comparisons: less-than and greater-than
|-
|7 || ==   != || Comparisons: equal and not equal
|-
|8 || & || Bitwise AND
|-
|9 || ^ || Bitwise exclusive OR (XOR)
|-
|10 || | || Bitwise inclusive (normal) OR
|-
|11 || && || Logical AND
|-
|12 || || || Logical OR
|-
|13 || ? : || Conditional expression (ternary)
|-
|14 || =   +=   -=   *=   /=   %=   &=   |=   ^=   <<=   >>= || Assignment operators (right to left)
|-
|15 || , || Comma operator
|}
Examples:
* !A + !B is interpreted as (!A) + (!B)
* ++A + !B is interpreted as (++A) + (!B)
* A + B * C is interpreted as A + (B * C)
* A || B && C is interpreted as A || (B && C)
* A && B == C is interpreted as A && (B == C)
* A & B == C is interpreted as A & (B == C)
(In Python, Ruby, PARI/GP and other popular languages, A & B == C is interpreted as (A & B) == C.)

Source-to-source compilers that compile to multiple languages need to explicitly deal with the issue of different order of operations across languages. Haxe for example standardizes the order and enforces it by inserting brackets where it is appropriate.

The accuracy of software developer knowledge about binary operator precedence has been found to closely follow their frequency of occurrence in source code.

## History
The order of operations emerged progressively over centuries. The rule that multiplication has precedence over addition was incorporated into the development of algebraic notation in the 1600s, since the distributive property implies this as a natural hierarchy. As recently as the 1920s, the historian of mathematics Florian Cajori identifies disagreement about whether multiplication should have precedence over division, or whether they should be treated equally. The term "order of operations" and the "PEMDAS/BEDMAS" mnemonics were formalized only in the late 19th or early 20th century, as demand for standardized textbooks grew.
Ambiguity about issues such as whether implicit multiplication takes precedence over explicit multiplication and division in such expressions as *a*/2*b*, which could be interpreted as *a*/(2*b*) or (*a*/2) × *b*, imply that the conventions are not yet completely stable.

## See also
* Common operator notation (for a more formal description)
* Hyperoperation
* Logical connective#Order of precedence
* Operator associativity
* Operator overloading
* Operator precedence in C and C++
* Polish notation
* Reverse Polish notation

## Notes
 (p. 22), and the first volume of the *Feynman Lectures* contains expressions such as 1/2sqrt(*N*) (p. 6–7). In both books, these expressions are written with the convention that the solidus is evaluated last.


}}

## References










## Further reading
*

## External links
*
* Zachary, Joseph L. (1997) "Operator Precedence", supplement to *Introduction to Scientific Programming*. University of Utah. Maple worksheet, Mathematica notebook.

