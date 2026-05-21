> Content sourced from [Wikipedia](https://en.wikipedia.org/wiki/Expression_%28mathematics%29) — CC BY-SA 4.0

# Evaluate expression at given value of x

In mathematics, an **expression** is an arrangement of symbols following the context-dependent, syntactic conventions of mathematical notation. Symbols can denote numbers, variables, operations, and functions. Other symbols include punctuation marks and brackets, used for grouping where there is not a well-defined order of operations.

Expressions are commonly distinguished from *formulas*: expressions usually denote mathematical objects, whereas formulas are statements *about* mathematical objects. This is analogous to natural language, where a noun phrase refers to an object, and a whole sentence refers to a fact. For example, \(8x-5\) and \(3\) are both expressions, while the inequality \(8x-5 \geq 3\) is a formula. However, formulas are often considered as expressions that can be evaluated to the Boolean values *true* or *false*.

To *evaluate* an expression means to find a numerical value equivalent to the expression. Expressions can be *evaluated* or *simplified* by replacing operations that appear in them with their result. For example, the expression \(8\times 2-5\) simplifies to \(16-5\), and evaluates to \(11.\)

An expression is often used to define a function, by taking the variables to be arguments, or inputs, of the function, and assigning the output to be the evaluation of the resulting expression. For example, \(x\mapsto x^2+1\) and \(f(x) = x^2 + 1\) define the function that associates to each number its square plus one. An expression with no variables would define a constant function. Usually, two expressions are considered equal or *equivalent* if they define the same function. Such an equality is called a "semantic equality", that is, both expressions "mean the same thing."

## Elementary mathematics
### Variables and evaluation
 In elementary algebra, a *variable* in an expression is a letter that represents a number whose value may change. To *evaluate an expression* with a variable means to find the value of the expression when the variable is substituted with given number. Expressions can be *simplified* by replacing operations that appear in them with their result, or by combining like-terms. The *evaluation* of an expression consists of repeating simplification steps until getting eventually a single number.

For example, take the expression \(4 x^2 + 8\); it can be evaluated for in the following steps:

\(4\cdot3^2 + 8\), (replace x with 3)

\(4 \cdot 9 + 8\) (evaluate the square)

\(36 + 8\) (evaluate the multiplication)

\(44\) (evaluate the addition)

A *term* is a constant or the product of a constant and one or more variables. Some examples include \(7, \; 5x, \; 13x^2y, \; 4b\) The constant of the product is called the coefficient. Terms that are either constants or have the same variables raised to the same powers are called *like terms*. If there are like terms in an expression, one can simplify the expression by combining the like terms. One adds the coefficients and keeps the same variable.

\(4x + 7x +2x = 13x\)

Any variable can be classified as being either a free variable or a bound variable. For a given combination of values for the free variables, an expression may be evaluated, although for some combinations of values of the free variables, the value of the expression may be undefined. Thus an expression represents an operation over constants and free variables and whose output is the resulting value of the expression.

For a non-formalized language, that is, in most mathematical texts outside of mathematical logic, for an individual expression it is not always possible to identify which variables are free and bound. For example, in \(\sum_{i < k} a_{ik}\), depending on the context, the variable \(i\) can be free and \(k\) bound, or vice-versa, but they cannot both be free. Determining which value is assumed to be free depends on context and semantics.

### Equivalence

An expression is often used to define a function, or denote compositions of functions, by taking the variables to be arguments, or inputs, of the function, and assigning the output to be the evaluation of the resulting expression. For example, \(x\mapsto x^2+1\) and \(f(x) = x^2 + 1\) define the function that associates to each number its square plus one. An expression with no variables would define a constant function. In this way, two expressions are said to be equivalent if, for each combination of values for the free variables, they have the same output, i.e., they represent the same function. The equivalence between two expressions is called an identity and is sometimes denoted with \(\equiv.\)

For example, in the expression \(\sum_{n=1}^{3} (2nx),\) the variable *n* is bound, and the variable *x* is free. This expression is equivalent to the simpler expression 12 *x*; that is
\[
\sum_{n=1}^{3} (2nx)\equiv 12x.
\]
 The value for *x* = 3 is 36, which can be denoted
\[
\sum_{n=1}^{3} (2nx)\Big|_{x=3}= 36.
\]

## Well-defined expressions

The language of mathematics exhibits a kind of grammar (called formal grammar) about how expressions may be written. There are two considerations for well-definedness of mathematical expressions, syntax and semantics. Syntax is concerned with the rules used for constructing, or transforming the symbols of an expression without regard to any interpretation or meaning given to them. Expressions that are syntactically correct are called well-formed. Semantics is concerned with the meaning of these well-formed expressions. Expressions that are semantically correct are called well-defined.

### Well-formed
The syntax of mathematical expressions can be described somewhat informally as follows: the allowed operators must have the correct number of inputs in the correct places (usually written with infix notation), the sub-expressions that make up these inputs must be well-formed themselves, have a clear order of operations, etc. Strings of symbols that conform to the rules of syntax are called *well-formed*, and those that are not well-formed are called, *ill-formed*, and do not constitute mathematical expressions.

For example, in arithmetic, the expression *1 + 2 × 3* is well-formed, but
\(\times4)x+,/y\).
is not.

However, being well-formed is not enough to be considered well-defined. For example in arithmetic, the expression \(\frac{1}{0}\) is well-formed, but it is not well-defined (see Division by zero). Such expressions are called undefined.

### Well-defined
Semantics is the study of meaning. Formal semantics is about attaching meaning to expressions. An expression that defines a unique value or meaning is said to be well-defined. Otherwise, the expression is said to be ill defined or ambiguous.**** In general the meaning of expressions is not limited to designating values; for instance, an expression might designate a condition, or an equation that is to be solved, or it can be viewed as an object in its own right that can be manipulated according to certain rules. Certain expressions that designate a value simultaneously express a condition that is assumed to hold, for instance those involving the operator \(\oplus\) to designate an internal direct sum.

In algebra, an expression may be used to designate a value, which might depend on values assigned to variables occurring in the expression. The determination of this value depends on the semantics attached to the symbols of the expression. The choice of semantics depends on the context of the expression. The same syntactic expression *1 + 2 × 3* can have different values (mathematically 7, but also 9), depending on the order of operations implied by the context (See also Operations § Calculators).

For real numbers, the product \(a \times b \times c\) is unambiguous because \((a \times b)\times c = a \times (b \times c)\); hence the notation is said to be *well defined*. This property, also known as associativity of multiplication, guarantees the result does not depend on the sequence of multiplications; therefore, a specification of the sequence can be omitted. The subtraction operation is non-associative; despite that, there is a convention that \(a-b-c\) is shorthand for \((a-b)-c\), thus it is considered "well-defined". On the other hand, Division is non-associative, and in the case of \(a/b/c\), parenthesization conventions are not well established; therefore, this expression is often considered ill-defined.

Unlike with functions, notational ambiguities can be overcome by means of additional definitions (e.g., rules of precedence, associativity of the operator). For example, in the programming language C, the operator - for subtraction is *left-to-right-associative*, which means that a-b-c is defined as (a-b)-c, and the operator = for assignment is *right-to-left-associative*, which means that a=b=c is defined as a=(b=c). In the programming language APL there is only one rule: from right to left – but parentheses first.

## Formal definition
The term 'expression' is part of the language of mathematics, that is to say, it is not defined *within* mathematics, but taken as a primitive part of the language. To attempt to define the term would not be doing mathematics, but rather, one would be engaging in a kind of metamathematics (the metalanguage of mathematics), usually mathematical logic. Within mathematical logic, mathematics is usually described as a kind of formal language, and a well-formed expression can be defined recursively as follows:

The alphabet consists of:

* A set of individual constants: Symbols representing fixed objects in the domain of discourse, such as numerals (1, 2.5, 1/7,... ), sets (\(\varnothing, \{1,2,3\}\),... ), truth values (T or F), etc.
* A set of individual variables: A countably infinite amount of symbols representing variables used for representing an unspecified object in the domain. (Usually letters like , or )
* A set of operations: Function symbols representing operations that can be performed on elements over the domain, like addition (+), multiplication (×), or set operations like union (∪), or intersection (∩). (Functions can be understood as unary operations)
* Brackets ( )

With this alphabet, the recursive rules for forming a well-formed expression (WFE) are as follows:

* Any constant or variable as defined are the atomic expressions, the simplest well-formed expressions (WFE's). For instance, the constant \(2\) or the variable \(x\) are syntactically correct expressions.
* Let \(F\) be a metavariable for any n-ary operation over the domain, and let \(\phi_1, \phi_2,... \phi_n\) be metavariables for any WFE's.
Then \(F(\phi_1, \phi_2,... \phi_n)\) is also well-formed. Here represented with prefix notation, but other notations may be used such as Infix notation like \(3+4\), or possibly non-linear notations such as with matrices or summation notation if allowed.
For instance, if the domain of discourse is the real numbers, \(F\) can denote the binary operation +, then \(\phi_1 + \phi_2\) is well-formed. Or \(F\) can be the unary operation \(\surd\) so \(\sqrt{\phi_1}\) is well-formed.
Brackets are initially around each non-atomic expression, but they can be deleted in cases where there is a defined order of operations, or where order doesn't matter (i.e. where operations are associative).

A well-formed expression can be thought as a syntax tree. The leaf nodes are always atomic expressions. Operations \(+\) and \(\cup\) have exactly two child nodes, while operations \(\sqrt{x}\), \(\text{ln}(x)\) and \(\frac{d}{dx}\) have exactly one. There are countably infinitely many WFE's, however, each WFE has a finite number of nodes.

## Computer science

In computer science, an *expression* is a syntactic entity in a programming language that may be evaluated to determine its value or fail to terminate, in which case the expression is undefined. It is a combination of one or more constants, variables, functions, and operators that the programming language interprets (according to its particular rules of precedence and of association) and computes to produce ("to return", in a stateful environment) another value. This process, for mathematical expressions, is called *evaluation*.
In simple settings, the resulting value is usually one of various primitive types, such as string, Boolean, or numerical (such as integer, floating-point, or complex).

In computer algebra, formulas are viewed as expressions that can be evaluated as a Boolean, depending on the values that are given to the variables occurring in the expressions. For example \(8x-5 \geq 3\) takes the value *false* if is given a value less than 1, and the value *true* otherwise.

Expressions are often contrasted with statements—syntactic entities that have no value (an instruction).

Except for numbers and variables, every mathematical expression may be viewed as the symbol of an operator followed by a sequence of operands. In computer algebra software, the expressions are usually represented in this way. This representation is very flexible, and many things that seem not to be mathematical expressions at first glance, may be represented and manipulated as such. For example, an equation is an expression with "=" as an operator, a matrix may be represented as an expression with "matrix" as an operator and its rows as operands.

See: Computer algebra expression

### Computation

A computation is any type of arithmetic or non-arithmetic calculation that is "well-defined". The notion that mathematical statements should be 'well-defined' had been argued by mathematicians since at least the 1600s, but agreement on a suitable definition proved elusive. A candidate definition was proposed independently by several mathematicians in the 1930s. The best-known variant was formalised by the mathematician Alan Turing, who defined a well-defined statement or calculation as any statement that could be expressed in terms of the initialisation parameters of a Turing machine. Turing's definition apportioned "well-definedness" to a very large class of mathematical statements, including all well-formed algebraic statements, and all statements written in modern computer programming languages.

Despite the widespread uptake of this definition, there are some mathematical concepts that have no well-defined characterisation under this definition. This includes the halting problem and the busy beaver game. It remains an open question as to whether there exists a more powerful definition of 'well-defined' that is able to capture both computable and 'non-computable' statements. All statements characterised in modern programming languages are well-defined, including C++, Python, and Java.

Common examples of computation are basic arithmetic and the execution of computer algorithms. A calculation is a deliberate mathematical process that transforms one or more inputs into one or more outputs or *results*. For example, multiplying 7 by 6 is a simple algorithmic calculation. Extracting the square root or the cube root of a number using mathematical models is a more complex algorithmic calculation.

#### Rewriting
Expressions can be computed by means of an *evaluation strategy.* To illustrate, executing a function call f(a,b) may first evaluate the arguments a and b, store the results in references or memory locations ref_a and ref_b, then evaluate the function's body with those references passed in. This gives the function the ability to look up the original argument values passed in through dereferencing the parameters (some languages use specific operators to perform this), to modify them via assignment as if they were local variables, and to return values via the references. This is the call-by-reference evaluation strategy. Evaluation strategy is part of the semantics of the programming language definition. Some languages, such as PureScript, have variants with different evaluation strategies. Some declarative languages, such as Datalog, support multiple evaluation strategies. Some languages define a calling convention.

In rewriting, a reduction strategy or rewriting strategy is a relation specifying a rewrite for each object or term, compatible with a given reduction relation. A rewriting strategy specifies, out of all the reducible subterms (redexes), which one should be reduced (*contracted*) within a term. One of the most common systems involves lambda calculus.

### Polynomial evaluation

A polynomial consists of variables and coefficients, that involve only the operations of addition, subtraction, multiplication and exponentiation to nonnegative integer powers, and has a finite number of terms. The problem of polynomial evaluation arises frequently in practice. In computational geometry, polynomials are used to compute function approximations using Taylor polynomials. In cryptography and hash tables, polynomials are used to compute *k*-independent hashing.

In the former case, polynomials are evaluated using floating-point arithmetic, which is not exact. Thus different schemes for the evaluation will, in general, give slightly different answers. In the latter case, the polynomials are usually evaluated in a finite field, in which case the answers are always exact.

For evaluating the univariate polynomial \(a_nx^n+a_{n-1}x^{n-1}+\cdots +a_0,\) the most naive method would use \(n\) multiplications to compute \(a_nx^n\), use \(n-1\) multiplications to compute \(a_{n-1} x^{n-1}\) and so on for a total of \(\frac{n(n+1)}{2}\) multiplications and \(n\) additions. Using better methods, such as Horner's rule, this can be reduced to \(n\) multiplications and \(n\) additions. If some preprocessing is allowed, even more savings are possible

## Types of expressions
### Algebraic expression
An *algebraic expression* is an expression built up from algebraic constants, variables, and the algebraic operations (addition, subtraction, multiplication, division and exponentiation by a rational number). For example, 3*x*\(^{2}\) − 2*xy* + *c* is an algebraic expression. Since taking the square root is the same as raising to the power , the following is also an algebraic expression:
\(\sqrt{\frac{1-x^2}{1+x^2\)
See also: Algebraic equation and Algebraic closure

### Polynomial expression
A polynomial is an expression built with scalars (numbers of elements of some field), variables, and the operators of addition, multiplication, and exponentiation to nonnegative integer powers; for example \(3(x+1)^2 - xy.\)

Using associativity, commutativity and distributivity, every polynomial expression is equivalent to a polynomial, that is an expression that is a linear combination of products of integer powers of the indeterminates. For example the above polynomial expression is equivalent (denote the same polynomial as \(3x^2-xy+6x+3.\)

Many author do not distinguish polynomials and polynomial expressions. In this case the expression of a polynomial expression as a linear combination is called the *canonical form*, *normal form*, or *expanded form* of the polynomial.

### Formal expression

A **formal expression** is a kind of string of symbols, created by the same production rules as standard expressions, however, they are used without regard to the meaning of the expression. In this way, two *formal expressions* are considered equal only if they are syntactically equal, that is, if they are the exact same expression. For instance, the formal expressions "2" and "1+1" are not equal.

### Lambda calculus

Formal languages allow formalizing the concept of well-formed expressions.

In the 1930s, a new type of expression, the lambda expression, was introduced by Alonzo Church and Stephen Kleene for formalizing functions and their evaluation. The lambda operators (lambda abstraction and function application) form the basis for lambda calculus, a formal system used in mathematical logic and programming language theory.

The equivalence of two lambda expressions is undecidable (but see unification (computer science)). This is also the case for the expressions representing real numbers, which are built from the integers by using the arithmetical operations, the logarithm and the exponential (Richardson's theorem).

## History

### Early written mathematics

The earliest written mathematics likely began with tally marks, where each mark represented one unit, carved into wood or stone. An example of early counting is the Ishango bone, found near the Nile and dating back over 20,000 years ago, which is thought to show a six-month lunar calendar. Ancient Egypt developed a symbolic system using hieroglyphics, assigning symbols for powers of ten and using addition and subtraction symbols resembling legs in motion. This system, recorded in texts like the Rhind Mathematical Papyrus (c. 2000–1800 BC), influenced other Mediterranean cultures. In Mesopotamia, a similar system evolved, with numbers written in a base-60 (sexagesimal) format on clay tablets written in Cuneiform, a technique originating with the Sumerians around 3000 BC. This base-60 system persists today in measuring time and angles.

### Syncopated stage
The "syncopated" stage of mathematics introduced symbolic abbreviations for commonly used operations and quantities, marking a shift from purely geometric reasoning. Ancient Greek mathematics, largely geometric in nature, drew on Egyptian numerical systems (especially Attic numerals), with little interest in algebraic symbols, until the arrival of Diophantus of Alexandria, who pioneered a form of syncopated algebra in his *Arithmetica,* which introduced symbolic manipulation of expressions. His notation represented unknowns and powers symbolically, but without modern symbols for relations (such as equality or inequality) or exponents. An unknown number was called \(\zeta\). The square of \(\zeta\) was \(\Delta^v\); the cube was \(K^v\); the fourth power was \(\Delta^v\Delta\); the fifth power was \(\Delta K^v\); and \(\pitchfork\) meant to subtract everything on the right from the left. So for example, what would be written in modern notation as:

\[
x^3 - 2x^2 + 10x -1,
\]

Would be written in Diophantus's syncopated notation as:
\(\Kappa^{\upsilon} \overline{\alpha} \; \zeta \overline{\iota} \;\, \pitchfork \;\, \Delta^{\upsilon} \overline{\beta} \; \Mu \overline{\alpha} \,\;\)
In the 7th century, Brahmagupta used different colours to represent the unknowns in algebraic equations in the *Brāhmasphuṭasiddhānta*. Greek and other ancient mathematical advances, were often trapped in cycles of bursts of creativity, followed by long periods of stagnation, but this began to change as knowledge spread in the early modern period.

### Symbolic stage and early arithmetic

The transition to fully symbolic algebra began with Ibn al-Banna' al-Marrakushi (1256–1321) and Abū al-Ḥasan ibn ʿAlī al-Qalaṣādī, (1412–1482) who introduced symbols for operations using Arabic characters. The plus sign (+) appeared around 1351 with Nicole Oresme, likely derived from the Latin *et* (meaning "and"), while the minus sign (−) was first used in 1489 by Johannes Widmann. Luca Pacioli included these symbols in his works, though much was based on earlier contributions by Piero della Francesca. The radical symbol (√) for square root was introduced by Christoph Rudolff in the 1500s, and parentheses for precedence by Niccolò Tartaglia in 1556. François Viète’s *New Algebra* (1591) formalized modern symbolic manipulation. The multiplication sign (×) was first used by William Oughtred and the division sign (÷) by Johann Rahn.

René Descartes further advanced algebraic symbolism in *La Géométrie* (1637), where he introduced the use of letters at the end of the alphabet (x, y, z) for variables, along with the Cartesian coordinate system, which bridged algebra and geometry. Isaac Newton and Gottfried Wilhelm Leibniz independently developed calculus in the late 17th century, with Leibniz's notation becoming the standard.
