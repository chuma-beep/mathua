> Content sourced from [Wikipedia](https://en.wikipedia.org/wiki/Logical_connective) — CC BY-SA 4.0

# Logical connectives: AND, OR, NOT, IMPLIES

In logic, a **logical connective** (also called a **logical operator**, **sentential connective**, or **sentential operator**) is an operator that combines or modifies one or more logical variables or formulas, similarly to how arithmetic connectives like \(+\) and \(-\) combine or negate arithmetic expressions. For instance, in the syntax of propositional logic, the binary connective \(\lor\) (meaning "or") can be used to join the two logical formulas \(P\) and \(Q\), producing the complex formula \(P \lor Q\).

Unlike in algebra, there are many symbols in use for each logical connective. The table "Logical connectives" shows examples.

Common connectives include negation, disjunction, conjunction, implication, and equivalence. In standard systems of classical logic, these connectives are interpreted as truth functions, though they receive a variety of alternative interpretations in nonclassical logics. Their classical interpretations are similar to the meanings of natural language expressions such as English "not", "or", "and", and "if", but not identical. Discrepancies between natural language connectives and those of classical logic have motivated nonclassical approaches to natural language meaning.

## Overview
In formal languages, truth functions are denoted by fixed symbols, ensuring that well-formed statements have a single interpretation. These symbols are called *logical connectives*, *logical operators*, *propositional operators*, or, in classical logic, *truth-functional connectives*. For the rules which allow new well-formed formulas to be constructed by joining other well-formed formulas using truth-functional connectives, see well-formed formula.

Logical connectives can be used to link zero or more statements, so one can speak about *-ary logical connectives*. The boolean constants *True* and *False* can be thought of as nullary operators. Negation is a unary connective, and so on.

<table>
  <tr>
    <th>colspan=2 | Symbol, name</th>
    <th>colspan=4 | Truthtable</th>
    <th>| Venndiagram</th>
    <td>-</td>
    <th>colspan=7 | Zeroary connectives (constants)</th>
    <td>-</td>
    <td>\(\top\)</td>
    <td>style="text-align:left; | Truth/tautology</td>
    <td>colspan=4 | 1</td>
    <td>-</td>
    <td>\(\bot\)</td>
    <td>style="text-align:left; | Falsity/contradiction</td>
    <td>colspan=4 | 0</td>
    <td>-</td>
    <th>colspan=7 | Unary connectives</th>
    <td>-</td>
    <td>colspan=2 | \(p\) =</td>
    <td>colspan=2 | 0</td>
    <td>colspan=2 | 1</td>
    <td>-</td>
    <td></td>
    <td>style="text-align:left; | Proposition \(p\)</td>
    <td>colspan=2 | 0</td>
    <td>colspan=2 | 1</td>
    <td>-</td>
    <td>\(\neg\)</td>
    <td>style="text-align:left; | Negation</td>
    <td>colspan=2 | 1</td>
    <td>colspan=2 | 0</td>
    <td>-</td>
    <th>colspan=9 | Binary connectives</th>
    <td>-</td>
    <td>colspan=2 | \(p\) =</td>
    <td>0</td>
    <td>0</td>
    <td>1</td>
    <td>1</td>
    <td>-</td>
    <td>colspan=2 | \(q\) =</td>
    <td>0</td>
    <td>1</td>
    <td>0</td>
    <td>1</td>
    <td>-</td>
    <td>\(\and\)</td>
    <td>| Conjunction</td>
    <td>0</td>
    <td>0</td>
    <td>0</td>
    <td>1</td>
    <td>-</td>
    <td>\(\uparrow\)</td>
    <td>| Alternative denial</td>
    <td>1</td>
    <td>1</td>
    <td>1</td>
    <td>0</td>
    <td>-</td>
    <td>\(\vee\)</td>
    <td>style="text-align:left; | Disjunction</td>
    <td>0</td>
    <td>1</td>
    <td>1</td>
    <td>1</td>
    <td>-</td>
    <td>\(\downarrow\)</td>
    <td>style="text-align:left; | Joint denial</td>
    <td>1</td>
    <td>0</td>
    <td>0</td>
    <td>0</td>
    <td>-</td>
    <td>\(\nleftrightarrow\)</td>
    <td>style="text-align:left; | Exclusive or</td>
    <td>0</td>
    <td>1</td>
    <td>1</td>
    <td>0</td>
    <td>-</td>
    <td>\(\leftrightarrow\)</td>
    <td>style="text-align:left; | Biconditional</td>
    <td>1</td>
    <td>0</td>
    <td>0</td>
    <td>1</td>
    <td>-</td>
    <td>\(\rightarrow\)</td>
    <td>style="text-align:left; | Material conditional</td>
    <td>1</td>
    <td>1</td>
    <td>0</td>
    <td>1</td>
    <td>-</td>
    <td>\(\nrightarrow\)</td>
    <td>style="text-align:left; | Material nonimplication</td>
    <td>0</td>
    <td>0</td>
    <td>1</td>
    <td>0</td>
    <td>-</td>
    <td>\(\leftarrow\)</td>
    <td>style="text-align:left; | Converse implication</td>
    <td>1</td>
    <td>0</td>
    <td>1</td>
    <td>1</td>
    <td>-</td>
    <td>\(\nleftarrow\)</td>
    <td>style="text-align:left; | Converse nonimplication</td>
    <td>0</td>
    <td>1</td>
    <td>0</td>
    <td>0</td>
    <td>-</td>
    <td>colspan=7" | More information</td>
    <td>
### List of common logical connectives Commonly used logical connectives include the following ones. * Negation (not): \(\neg\), \(\sim\), \(N\) (prefix) in which \(\neg\) is the most modern and widely used, and \(\sim\) is also common; * Conjunction (and): \(\wedge\), \(\&\), \(K\) (prefix) in which \(\wedge\) is the most modern and widely used; * Disjunction (or): \(\vee\), \(A\) (prefix) in which \(\vee\) is the most modern and widely used; * Implication (if...then): \(\to\), \(\supset\), \(\Rightarrow\), \(C\) (prefix) in which \(\to\) is the most modern and widely used, and \(\supset\) is also common; * Equivalence (if and only if): \(\leftrightarrow\), \(\subset\!\!\!\supset\), \(\Leftrightarrow\), \(\equiv\), \(E\) (prefix) in which \(\leftrightarrow\) is the most modern and widely used, and \(\subset\!\!\!\supset\) is commonly used where \(\supset\) is also used. For example, the meaning of the statements *it is raining* (denoted by \(p\)) and *I am indoors* (denoted by \(q\)) is transformed, when the two are combined with logical connectives: * It is ***not*** raining (\(\neg p\)); * It is raining ***and*** I am indoors (\(p \wedge q\)); * It is raining ***or*** I am indoors (\(p \lor q\)); * ***If*** it is raining, ***then*** I am indoors (\(p \rightarrow q\)); * ***If*** I am indoors, ***then*** it is raining (\(q \rightarrow p\)); * I am indoors ***if and only if*** it is raining (\(p \leftrightarrow q\)). It is also common to consider the *always true* formula and the *always false* formula to be connective (in which case they are nullary). * True formula: \(\top\), \(1\), \(V\) (prefix), or \(\mathrm{T}\); * False formula: \(\bot\), \(0\), \(O\) (prefix), or \(\mathrm{F}\). This table summarizes the terminology: {|</td>
  </tr>
  <tr>
    <th>Connective</th>
    <th>In English</th>
    <th>Noun for parts</th>
    <th>Verb phrase</th>
  </tr>
  <tr>
    <th>Conjunction</th>
    <td>Both A and B</td>
    <td>conjunct</td>
    <td>A and B are conjoined</td>
  </tr>
  <tr>
    <th>Disjunction</th>
    <td>Either A or B, or both</td>
    <td>disjunct</td>
    <td>A and B are disjoined</td>
  </tr>
  <tr>
    <th>Negation</th>
    <td>It is not the case that A</td>
    <td>negatum/negand</td>
    <td>A is negated</td>
  </tr>
  <tr>
    <th>Conditional</th>
    <td>If A, then B</td>
    <td>antecedent, consequent</td>
    <td>B is implied by A</td>
  </tr>
  <tr>
    <th>Biconditional</th>
    <td>A if, and only if, B</td>
    <td>equivalents</td>
    <td>A and B are equivalent</td>
  </tr>
</table>

| Connective | In English | Noun for parts | Verb phrase |
| --- | --- | --- | --- |
| Conjunction | Both A and B | conjunct | A and B are conjoined |
| Disjunction | Either A or B, or both | disjunct | A and B are disjoined |
| Negation | It is not the case that A | negatum/negand | A is negated |
| Conditional | If A, then B | antecedent, consequent | B is implied by A |
| Biconditional | A if, and only if, B | equivalents | A and B are equivalent |

### History of notations
* Negation: the symbol \(\neg\) appeared in Heyting in 1930 (compare to Frege's symbol ⫟ in his Begriffsschrift); the symbol \(\sim\) appeared in Russell in 1908; an alternative notation is to add a horizontal line on top of the formula, as in \(\overline{p}\); another alternative notation is to use a prime symbol as in \(p'\).
* Conjunction: the symbol \(\wedge\) appeared in Heyting in 1930 (compare to Peano's use of the set-theoretic notation of intersection \(\cap\)); the symbol \(\&\) appeared at least in Schönfinkel in 1924; the symbol \(\cdot\) comes from Boole's interpretation of logic as an elementary algebra.
* Disjunction: the symbol \(\vee\) appeared in Russell in 1908 (compare to Peano's use of the set-theoretic notation of union \(\cup\)); the symbol \(+\) is also used, in spite of the ambiguity coming from the fact that the \(+\) of ordinary elementary algebra is an exclusive or when interpreted logically in a two-element ring; punctually in the history a \(+\) together with a dot in the lower right corner has been used by Peirce.
* Implication: the symbol \(\to\) appeared in Hilbert in 1918; \(\supset\) was used by Russell in 1908 (compare to Peano's Ɔ the inverted C); \(\Rightarrow\) appeared in Bourbaki in 1954.
* Equivalence: the symbol \(\equiv\) in Frege in 1879; \(\leftrightarrow\) in Becker in 1933 (not the first time and for this see the following); \(\Leftrightarrow\) appeared in Bourbaki in 1954; other symbols appeared punctually in the history, such as \(\supset\subset\) in Gentzen, \(\sim\) in Schönfinkel or \(\subset\supset\) in Chazal,
* True: the symbol \(1\) comes from Boole's interpretation of logic as an elementary algebra over the two-element Boolean algebra; other notations include \(\mathrm{V}\) (abbreviation for the Latin word "verum") to be found in Peano in 1889.
* False: the symbol \(0\) comes also from Boole's interpretation of logic as a ring; other notations include \(\Lambda\) (rotated \(\mathrm{V}\)) to be found in Peano in 1889.

Some authors used letters for connectives: \(\operatorname{u.}\) for conjunction (German's "und" for "and") and \(\operatorname{o.}\) for disjunction (German's "oder" for "or") in early works by Hilbert (1904); \(Np\) for negation, \(Kpq\) for conjunction, \(Dpq\) for alternative denial, \(Apq\) for disjunction, \(Cpq\) for implication, \(Epq\) for biconditional in Łukasiewicz in 1929.

### Redundancy
Such a logical connective as converse implication "\(\leftarrow\)" is actually the same as material conditional with swapped arguments; thus, the symbol for converse implication is redundant. In some logical calculi (notably, in classical logic), certain essentially different compound statements are logically equivalent. A less trivial example of a redundancy is the classical equivalence between \(\neg p\vee q\) and \(p\to q\). Therefore, a classical-based logical system does not need the conditional operator "\(\to\)" if "\(\neg\)" (not) and "\(\vee\)" (or) are already in use, or may use the "\(\to\)" only as a syntactic sugar for a compound having one negation and one disjunction.

There are sixteen Boolean functions associating the input truth values \(p\) and \(q\) with four-digit binary outputs. These correspond to possible choices of binary logical connectives for classical logic. Different implementations of classical logic can choose different functionally complete subsets of connectives.

One approach is to choose a *minimal* set, and define other connectives by some logical form, as in the example with the material conditional above. The following are the minimal functionally complete sets of operators in classical logic whose arities do not exceed 2:
;One element: \(\{\uparrow\}\), \(\{\downarrow\}\).
;Two elements: \(\{\vee, \neg\}\), \(\{\wedge, \neg\}\), \(\{\to, \neg\}\), \(\{\gets, \neg\}\), \(\{\to, \bot\}\), \(\{\gets, \bot\}\), \(\{\to, \nleftrightarrow\}\), \(\{\gets, \nleftrightarrow\}\), \(\{\to, \nrightarrow\}\), \(\{\to, \nleftarrow\}\), \(\{\gets, \nrightarrow\}\), \(\{\gets, \nleftarrow\}\), \(\{\nrightarrow, \neg\}\), \(\{\nleftarrow, \neg\}\), \(\{\nrightarrow, \top\}\), \(\{\nleftarrow, \top\}\), \(\{\nrightarrow, \leftrightarrow\}\), \(\{\nleftarrow, \leftrightarrow\}\).
;Three elements: \(\{\lor, \leftrightarrow, \bot\}\), \(\{\lor, \leftrightarrow, \nleftrightarrow\}\), \(\{\lor, \nleftrightarrow, \top\}\), \(\{\land, \leftrightarrow, \bot\}\), \(\{\land, \leftrightarrow, \nleftrightarrow\}\), \(\{\land, \nleftrightarrow, \top\}\).

Another approach is to use with equal rights connectives of a certain convenient and functionally complete, but *not minimal* set. This approach requires more propositional axioms, and each equivalence between logical forms must be either an axiom or provable as a theorem.

The situation, however, is more complicated in intuitionistic logic. Of its five connectives, {∧, ∨, →, ¬, ⊥}, only negation "¬" can be reduced to other connectives (see for more). Neither conjunction, disjunction, nor material conditional has an equivalent form constructed from the other four logical connectives.

## Natural language
The standard logical connectives of classical logic have rough equivalents in the grammars of natural languages. In English, as in many languages, such expressions are typically grammatical conjunctions. However, they can also take the form of complementizers, verb suffixes, and particles. The denotations of natural language connectives is a major topic of research in formal semantics, a field that studies the logical structure of natural languages.

The meanings of natural language connectives are not precisely identical to their nearest equivalents in classical logic. In particular, disjunction can receive an exclusive interpretation in many languages. Some researchers have taken this fact as evidence that natural language semantics is nonclassical. However, others maintain classical semantics by positing pragmatic accounts of exclusivity which create the illusion of nonclassicality. In such accounts, exclusivity is typically treated as a scalar implicature. Related puzzles involving disjunction include free choice inferences, Hurford's Constraint, and the contribution of disjunction in alternative questions.

Other apparent discrepancies between natural language and classical logic include the paradoxes of material implication, donkey anaphora and the problem of counterfactual conditionals. These phenomena have been taken as motivation for identifying the denotations of natural language conditionals with logical operators including the strict conditional, the variably strict conditional, as well as various dynamic operators.

The following table shows the standard classically definable approximations for the English connectives.

| English word | Connective | Symbol | Logical gate |
| --- | --- | --- | --- |
| not | negation | \(\neg\) | NOT |
| and | conjunction | \(\and\) | AND |
| or | disjunction | \(\vee\) | OR |
| if...then | material implication | \(\rightarrow\) | IMPLY |
| ... if | converse implication | \(\leftarrow\) |  |
| either...or | exclusive disjunction | \(\nleftrightarrow\) | XOR |
| if and only if | biconditional | \(\leftrightarrow\) | XNOR |
| not both | alternative denial | \(\uparrow\) | NAND |
| neither...nor | joint denial | \(\downarrow\) | NOR |
| but not | material nonimplication | \(\nrightarrow\) | NIMPLY |
| not...but | converse nonimplication | \(\nleftarrow\) |  |

## Properties
Some logical connectives possess properties that may be expressed in the theorems containing the connective. Some of those properties that a logical connective may have are:

; Associativity: Within an expression containing two or more of the same associative connectives in a row, the order of the operations does not matter as long as the sequence of the operands is not changed.
; Commutativity:The operands of the connective may be swapped, preserving logical equivalence to the original expression.
; Distributivity: A connective denoted by · distributes over another connective denoted by +, if *a* · (*b* + *c*) = (*a* · *b*) + (*a* · *c*) for all operands ,. ; Idempotence: Whenever the operands of the operation are the same, the compound is logically equivalent to the operand.
; Absorption: A pair of connectives ∧, ∨ satisfies the absorption law if \(a\land(a\lor b)=a\) for all operands ,. ; Monotonicity: If *f*(*a*\(_{1}\),... , *a*\(_{*n*}\)) ≤ *f*(*b*\(_{1}\),... , *b*\(_{*n*}\)) for all *a*\(_{1}\),... , *a*\(_{*n*}\), *b*\(_{1}\),... , *b*\(_{*n*}\) ∈ {0,1} such that *a*\(_{1}\) ≤ *b*\(_{1}\), *a*\(_{2}\) ≤ *b*\(_{2}\),... , *a*\(_{*n*}\) ≤ *b*\(_{*n*}\). E.g., ∨, ∧, ⊤, ⊥.
; Affinity: Each variable always makes a difference in the truth-value of the operation or it never makes a difference. E.g., ¬, ↔, \(\nleftrightarrow\), ⊤, ⊥.
; Duality: To read the truth-value assignments for the operation from top to bottom on its truth table is the same as taking the complement of reading the table of the same or another connective from bottom to top. Without resorting to truth tables it may be formulated as *g̃*(¬*a*\(_{1}\),... , ¬*a*\(_{*n*}\)) = ¬*g*(*a*\(_{1}\),... , *a*\(_{*n*}\)). E.g., ¬.
; Truth-preserving: The compound all those arguments are tautologies is a tautology itself. E.g., ∨, ∧, ⊤, →, ↔, ⊂ (see validity).
; Falsehood-preserving: The compound all those argument are contradictions is a contradiction itself. E.g., ∨, ∧, \(\nleftrightarrow\), ⊥, ⊄, ⊅ (see validity).
; Involutivity (for unary connectives): *f*(*f*(*a*)) = *a*. E.g. negation in classical logic.

For classical and intuitionistic logic, the "=" symbol means that corresponding implications "...→..." and "...←..." for logical compounds can be both proved as theorems, and the "≤" symbol means that "...→..." for logical compounds is a consequence of corresponding "...→..." connectives for propositional variables. Some many-valued logics may have incompatible definitions of equivalence and order (entailment).

Both conjunction and disjunction are associative, commutative and idempotent in classical logic, most varieties of many-valued logic and intuitionistic logic. The same is true about distributivity of conjunction over disjunction and disjunction over conjunction, as well as for the absorption law.

In classical logic and some varieties of many-valued logic, conjunction and disjunction are dual, and negation is self-dual, the latter is also self-dual in intuitionistic logic.

## Order of precedence
As a way of reducing the number of necessary parentheses, one may introduce precedence rules: ¬ has higher precedence than ∧, ∧ higher than ∨, and ∨ higher than →. So for example, \(P \vee Q \and{\neg R} \rightarrow S\) is short for \((P \vee (Q \and (\neg R))) \rightarrow S\).

Here is a table that shows a commonly used precedence of logical operators.
| Operator | Precedence |
| --- | --- |
| \(\neg\) | 1 |
| \(\and\) | 2 |
| \(\vee\) | 3 |
| \(\rightarrow\) | 4 |
| \(\leftrightarrow\) | 5 |

However, not all compilers use the same order; for instance, an ordering in which disjunction is lower precedence than implication or bi-implication has also been used. Sometimes precedence between conjunction and disjunction is unspecified requiring to provide it explicitly in given formula with parentheses. The order of precedence determines which connective is the "main connective" when interpreting a non-atomic formula.

## Table and Hasse diagram
The 16 logical connectives can be partially ordered to produce the following Hasse diagram. The partial order is defined by declaring that \(x \leq y\) if and only if whenever \(x\) holds then so does \(y.\)

## Applications
Logical connectives are used in computer science and in set theory.

### Computer science

A truth-functional approach to logical operators is implemented as logic gates in digital circuits. Practically all digital circuits (the major exception is DRAM) are built up from NAND, NOR, NOT, and transmission gates; see more details in Truth function in computer science. Logical operators over bit vectors (corresponding to finite Boolean algebras) are bitwise operations.

But not every usage of a logical connective in computer programming has a Boolean semantic. For example, lazy evaluation is sometimes implemented for *P* ∧ *Q* and *P* ∨ *Q*, so these connectives are not commutative if either or both of the expressions , have side effects. Also, a conditional, which in some sense corresponds to the material conditional connective, is essentially non-Boolean because for if (P) then Q;, the consequent Q is not executed if the antecedent P is false (although a compound as a whole is successful ≈ "true" in such case). This is closer to intuitionist and constructivist views on the material conditional— rather than to classical logic's views.

### Set theory

Logical connectives are used to define the fundamental operations of set theory, as follows:

| Set operation | Connective | Definition |
| --- | --- | --- |
| Intersection | Conjunction | \(A \cap B = \{x : x \in A \land x \in B \}\) |
| Union | Disjunction | \(A \cup B = \{x : x \in A \lor x \in B \}\) |
| Complement | Negation | \(\overline{A} = \{x : x \notin A \}\) |
| Subset | Implication | \(A \subseteq B \leftrightarrow (x \in A \rightarrow x \in B)\) |
| Equality | Biconditional | \(A = B \leftrightarrow (\forall X)[A \in X \leftrightarrow B \in X]\) |

This definition of set equality is equivalent to the axiom of extensionality.
