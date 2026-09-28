> Content sourced from [Algebrica](https://algebrica.org/propositional-logic/) — CC BY-NC 4.0

## Propositional language

1.Propositional logic studies arguments whose validity depends on how propositional connectives combine propositions. It treats each atomic proposition as either true or false and determines the truth value of a compound proposition from the truth values of its components.

2.The subject matter of an atomic proposition is irrelevant to this calculation. If an argument has the form $p,$ $p \rightarrow q,$ therefore $q,$ its validity depends on that form, regardless of the propositions represented by $p$ and $q.$ This restriction makes propositional logic precise, but it also limits what the language can express.

3.A propositional language $\mathrm{Prop}[P]$ has atomic symbols, logical connectives, and parentheses. The set $P$ contains the atomic proposition symbols, for example $P = \\{p, q, r\\}.$ Parentheses record how the connectives group the formulas.

4.The language used here has the following connectives:

5.

  * Negation $\neg$
  * Conjunction $\wedge$
  * Disjunction $\lor$
  * Material conditional $\rightarrow$
  * Material biconditional $\leftrightarrow$
  * Exclusive disjunction $\oplus$


6.The choice of primitive connectives is conventional. The biconditional and exclusive disjunction can be defined from other connectives, but separate symbols make common formulas shorter.

7.The language also has two propositional constants, $\top$ and $\bot.$ The first is true under every interpretation, the second is false under every interpretation. They take no arguments, so they are connectives of arity zero, and each of them has the truth value of a compound formula, $\top$ that of $p \lor \neg p$ and $\bot$ that of $p \wedge \neg p.$

8.A formula built from the symbols in $P$ according to the formation rules is a well-formed formula (WFF). Atomic propositions are the simplest WFFs, while every other WFF has one or more WFFs as immediate components.

## Logical connectives

9.A connective is truth-functional when the truth values of its immediate components uniquely determine the truth value of the compound formula. Every connective in $\mathrm{Prop}[P]$ is truth-functional.

10.

  * The negation $\neg p$ is true exactly when $p$ is false.
  * The conjunction $p \wedge q$ is true exactly when $p$ and $q$ are both true.
  * The disjunction $p \lor q$ is true exactly when at least one of $p$ and $q$ is true. The connective $\lor$ has the inclusive meaning, so $p \lor q$ is also true when both disjuncts are true.
  * The material conditional $p \rightarrow q$ is false exactly when its antecedent $p$ is true and its consequent $q$ is false.
  * The material biconditional $p \leftrightarrow q$ is true exactly when $p$ and $q$ have the same truth value.
  * The exclusive disjunction $p \oplus q$ is true exactly when $p$ and $q$ have different truth values.


11.The direction of a conditional requires care. The proposition "$p$ only if $q$" is $p \rightarrow q,$ while "$p$ if $q$" is $q \rightarrow p.$ In the first formula, $p$ is a sufficient condition for $q,$ and $q$ is a necessary condition for $p.$

12.Ordinary language has connectives that are not truth-functional. The truth of "it is necessary that $p$" is not determined by the truth of $p$ alone. Words such as "but" and "although" can express a contrast that conjunction does not retain, while a counterfactual conditional is not, in general, truth-functional. A symbolization in propositional logic has only the truth-functional structure of the original proposition.

## Symbolization example

13.Consider a propositional language $\mathrm{Prop}[P]$ over $P = \\{p, q\\},$ with the following symbolization key:

14.

  * $p =$ the door is open.
  * $q =$ the window is open.


15.The two atomic propositions and the connectives are enough to symbolize several compound propositions.

16.

  * The door is not open is written $\neg p.$
  * The door and the window are both open is written $p \wedge q.$
  * The door or the window is open, possibly both, is written $p \lor q.$
  * If the door is open, then the window is open is written $p \rightarrow q.$
  * The door is open if and only if the window is open is written $p \leftrightarrow q.$
  * Either the door or the window is open, but not both, is written $p \oplus q.$
  * Neither the door nor the window is open is written $\neg(p \lor q).$
  * The negation of the conditional "if the door is open, then the window is open" is written $\neg(p \rightarrow q).$


17.Parentheses are necessary in the last two formulas because the negation has the entire compound formula as its scope. Thus $\neg(p \lor q)$ differs from $\neg p \lor q.$

## Formation rules

18.The formation rules are an inductive definition of the WFFs of $\mathrm{Prop}[P].$

19.

  * Every atomic proposition $p \in P$ is a WFF, and the constants $\top$ and $\bot$ are WFFs.
  * If $\varphi$ is a WFF, then $\neg\varphi$ is a WFF.
  * If $\varphi$ and $\psi$ are WFFs, then each of $(\varphi \wedge \psi),$ $(\varphi \lor \psi),$ $(\varphi \rightarrow \psi),$ $(\varphi \leftrightarrow \psi),$ and $(\varphi \oplus \psi)$ is a WFF.
  * No other expression is a WFF.


20.The last clause excludes every string that cannot be obtained by finitely many applications of the preceding clauses. These clauses also determine the structure of each formula. Every compound WFF has a unique main connective, the connective applied at the last step of its construction.

21.We omit the outermost pair of parentheses when no ambiguity results. Under this convention, the main connective of $\neg(p \wedge q)$ is $\neg,$ while the main connective of $\neg p \wedge q$ is $\wedge.$ The scope of an occurrence of a connective is the subformula to which that occurrence applies. Parentheses therefore determine both the main connective and the scope of the inner connectives.

22.Two further conventions remove the remaining parentheses. The connectives have a precedence order, from the one that binds most strongly to the one that binds least:

$$ \neg \quad \wedge \quad \lor \quad \rightarrow \quad \leftrightarrow $$

23.The formula $\neg p \wedge q \rightarrow r$ is therefore an abbreviation of $((\neg p) \wedge q) \rightarrow r.$ Connectives of equal precedence associate to the left, so $p \wedge q \wedge r$ abbreviates $(p \wedge q) \wedge r.$ The grouping is immaterial for $\wedge$ and $\lor,$ which are associative, and it is not immaterial for $\rightarrow.$ Under $M(p) = M(q) = M(r) = F$ the formula $(p \rightarrow q) \rightarrow r$ is false while $p \rightarrow (q \rightarrow r)$ is true.

> 24.
> 
> Several texts read a chain of conditionals to the right, so that $p \rightarrow q \rightarrow r$ abbreviates $p \rightarrow (q \rightarrow r).$ The two conventions disagree, and explicit parentheses on nested conditionals avoid the ambiguity.

## Semantics

25.The semantics of propositional logic has two truth values:

$$ \mathrm{Bool} := \\{\ T, F\ \\} $$

26.The constants receive the same value in every case, $T$ for $\top$ and $F$ for $\bot.$ The characteristic truth tables define the connectives. A single table for the six connectives is the following:

$$ \begin{array}{cc|cccccc} p & q & \neg p & p \wedge q & p \lor q & p \rightarrow q & p \leftrightarrow q & p \oplus q \\\\[6pt] \hline T & T & F & T & T & T & T & F \\\\[6pt] T & F & F & F & T & F & F & T \\\\[6pt] F & T & T & F & T & T & F & T \\\\[6pt] F & F & T & F & F & T & T & F \end{array} $$

27.A complete truth table has one row for each assignment of truth values to the distinct atomic propositions in a formula. If a formula contains $n$ distinct atomic propositions, its complete truth table has $2^n$ rows. The column under the main connective contains the truth value of the whole formula on each assignment.

28.The table also gives the following logical equivalences:

29.

  * $p \rightarrow q \equiv \neg p \lor q$
  * $p \leftrightarrow q \equiv (p \rightarrow q) \wedge (q \rightarrow p)$
  * $p \oplus q \equiv (p \lor q) \wedge \neg(p \wedge q)$


30.The symbol $\equiv$ is a metalanguage relation between formulas. It means that the formulas have the same truth value under every assignment. The symbol is not another connective of $\mathrm{Prop}[P].$

## Interpretations

31.An interpretation, also called a valuation in propositional logic, is a function that assigns a truth value to each atomic proposition in $P$:

$$ M : P \rightarrow \\{T, F\\} $$

32.The value of a compound formula under $M$ is then determined recursively by its formation and the truth tables of the connectives.

33.

  * If $\varphi$ is true under $M,$ we write $M \models \varphi$ and call $M$ a model of $\varphi.$
  * If $\varphi$ is false under $M,$ we write $M \not\models \varphi$ and call $M$ a counter-model of $\varphi.$


34.Consider the interpretation in which $p$ is true and $q$ is false. The formula $p \rightarrow \neg q$ has the following row:

$$ \begin{array}{cc|cc} p & q & \neg q & p \rightarrow \neg q \\\\[6pt] \hline T & F & T & T \end{array} $$

35.Under this interpretation, the antecedent $p$ and the consequent $\neg q$ are true. Hence $p \rightarrow \neg q$ is true, and $M \models p \rightarrow \neg q.$

36.The following notions are defined by quantifying over interpretations.

37.

  * A formula $\varphi$ is satisfiable if some interpretation satisfies it.
  * A formula $\varphi$ is a tautology if every interpretation satisfies it.
  * A formula $\varphi$ is a contradiction if no interpretation satisfies it.
  * A formula $\varphi$ is contingent if some interpretation satisfies it and some interpretation does not satisfy it.


38.Thus every tautology and every contingent formula is satisfiable, while every contradiction is unsatisfiable. A set of formulas $S$ is jointly satisfiable if some interpretation satisfies every formula in $S.$ If no such interpretation exists, $S$ is jointly unsatisfiable, or inconsistent.

39.Two formulas $\varphi$ and $\psi$ are logically equivalent when they agree under every interpretation:

$$ \varphi \equiv \psi \Longleftrightarrow \forall M \ (M \models \varphi \Longleftrightarrow M \models \psi) $$

## Logical consequence

40.A formula $\varphi$ is a logical consequence of a set of formulas $S$ if every interpretation that satisfies all formulas in $S$ also satisfies $\varphi.$ This relation is written as follows:

$$ S \models \varphi $$

41.Equivalently, no interpretation makes every formula in $S$ true and $\varphi$ false. The symbol $\models$ is a relation in the metalanguage, whereas $\rightarrow$ is a connective that forms a new formula. For two formulas $\varphi$ and $\psi,$ the connection between them is the following:

$$ \varphi \models \psi \Longleftrightarrow \models \varphi \rightarrow \psi $$

42.Consider the set $S = \\{p, p \rightarrow q\\}$ and the proposed consequence $q.$ The relevant truth table is the following:

$$ \begin{array}{cc|c} p & q & p \rightarrow q \\\\[6pt] \hline T & T & T \\\\[6pt] T & F & F \\\\[6pt] F & T & T \\\\[6pt] F & F & T \end{array} $$

43.Only the first row makes both members of $S$ true, and that row also makes $q$ true. Therefore $S \models q.$ The corresponding inference rule is modus ponens.

44.A third formulation of the relation replaces the inspection of the models of $S$ by a question about a single set of formulas:

$$ S \models \varphi \Longleftrightarrow S \cup \\{\neg\varphi\\} \ \text{is unsatisfiable} $$

45.Assume $S \models \varphi$ and let $M$ satisfy every formula of $S \cup \\{\neg\varphi\\}.$ From $M \models S$ we get $M \models \varphi,$ while $M \models \neg\varphi$ gives $M \not\models \varphi,$ and the two conclusions are incompatible. Conversely, assume that $S \cup \\{\neg\varphi\\}$ is unsatisfiable and let $M \models S.$ If $\varphi$ were false under $M,$ then $M$ would satisfy $\neg\varphi$ and hence the whole set. So every model of $S$ is a model of $\varphi.$ The [procedures of automated deduction](<../automated-deduction-in-propositional-logic/>) test the right-hand side, since a proof of unsatisfiability can be searched for mechanically.

## Inference rules

46.An inference rule is a schematic pattern for deriving a conclusion from one or more premises. If $S \vdash \varphi,$ then $\varphi$ has a derivation from premises in $S$ within the chosen proof system. The symbol $\vdash$ concerns derivations, while $\models$ concerns interpretations.

47.A proof system is sound when $S \vdash \varphi$ implies $S \models \varphi,$ and it is complete when $S \models \varphi$ implies $S \vdash \varphi.$ Standard proof systems for propositional logic have both properties. The deductive closure of $S$ is the set of its logical consequences:

$$ \mathrm{Cn}(S) := \\{\ \varphi \mid S \models \varphi \ \\} $$

48.For a sound and complete system, $\mathrm{Cn}(S)$ is also the set of formulas derivable from $S.$ It is infinite for every $S,$ since it contains every tautology of the language.

49.Modus ponens derives $q$ from $p$ and $p \rightarrow q$:

$$ \frac{p \qquad p \rightarrow q}{q} $$

50.Modus tollens derives $\neg p$ from $\neg q$ and $p \rightarrow q$:

$$ \frac{\neg q \qquad p \rightarrow q}{\neg p} $$

51.The hypothetical syllogism derives $p \rightarrow r$ from $p \rightarrow q$ and $q \rightarrow r$:

$$ \frac{p \rightarrow q \qquad q \rightarrow r}{p \rightarrow r} $$

52.In each schema, the formulas above the line are the premises and the formula below the line is the conclusion.

53.For example, let $p$ mean that it is raining, let $q$ mean that the ground is wet, and let $r$ mean that the match is cancelled. From $p \rightarrow q$ and $q \rightarrow r,$ the hypothetical syllogism gives $p \rightarrow r.$ If $p$ is also a premise, modus ponens gives $r.$

54.Further rules govern the remaining connectives. Conjunction elimination and conjunction introduction relate a conjunction to its conjuncts, disjunction introduction weakens a formula to a disjunction, and the disjunctive syllogism removes a disjunct:

$$ \frac{\varphi \wedge \psi}{\varphi} \qquad \frac{\varphi \qquad \psi}{\varphi \wedge \psi} \qquad \frac{\varphi}{\varphi \lor \psi} \qquad \frac{\varphi \lor \psi \qquad \neg\varphi}{\psi} $$

55.A single rule, [resolution](<../automated-deduction-in-propositional-logic/>), covers modus ponens, modus tollens and the disjunctive syllogism at once, and it is the rule on which mechanical proof search is built.

## Normal forms

56.A literal is an atomic proposition or the negation of an atomic proposition. A clause is a disjunction of literals, and a term is a conjunction of literals.

57.A formula is in conjunctive normal form (CNF) if it is a conjunction of clauses:

$$ (l_{1,1} \lor \cdots \lor l_{1,k}) \wedge (l_{2,1} \lor \cdots \lor l_{2,m}) \wedge \cdots $$

58.A formula is in disjunctive normal form (DNF) if it is a disjunction of terms:

$$ (l_{1,1} \wedge \cdots \wedge l_{1,k}) \lor (l_{2,1} \wedge \cdots \wedge l_{2,m}) \lor \cdots $$

59.In either normal form, only $\neg,$ $\wedge,$ and $\lor$ occur, and every negation has an atomic proposition as its scope. A single literal is both a clause and a term, so it is both a CNF formula and a DNF formula.

60.Every propositional formula is logically equivalent to a formula in CNF and to a formula in DNF. One conversion method first removes $\rightarrow,$ $\leftrightarrow,$ and $\oplus,$ then moves each negation inward by double negation and De Morgan's laws, and finally applies the distributive laws.

61.For example, consider $\neg(p \lor q) \rightarrow r.$ The equivalence $\varphi \rightarrow \psi \equiv \neg\varphi \lor \psi$ gives the following calculation:

$$ \neg(p \lor q) \rightarrow r \equiv \neg\neg(p \lor q) \lor r \equiv (p \lor q) \lor r $$

62.The resulting formula is equivalent to $p \lor q \lor r.$ It is a single clause and is therefore in CNF. Since each disjunct is also a one-literal term, the same formula is in DNF.

63.A complete truth table gives another proof of the normal-form theorems. For DNF, take each row on which the original formula is true, form a term that is true only on that row, and disjoin those terms. For CNF, take each row on which the formula is false, form a clause that is false only on that row, and conjoin those clauses. When the formula is a contradiction, $p \wedge \neg p$ is an equivalent normal form. When it is a tautology, $p \lor \neg p$ is an equivalent normal form.

64.The DPLL procedure and several related satisfiability methods take CNF formulas as input, as does the [resolution procedure](<../automated-deduction-in-propositional-logic/>).

The Graph

Concept

The structure of the entry is shown in the conceptual map, where each branch represents a core component and the sub-nodes highlight the specific notions discussed.

Intermediate

2

Requires

2

Enables

The following concepts, [Functions](https://algebrica.org/functions/), [Sets](https://algebrica.org/sets/), are required as prerequisites for this entry.
