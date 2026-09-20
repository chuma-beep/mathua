> Content sourced from [Discrete Mathematics: An Open Introduction, 3rd ed](https://discrete.openmathbooks.org/dmoi3/sec_propositional.html) by Oscar Levin — CC BY-SA 4.0

# Propositional Logic

## Section 3.1 Propositional Logic

###### Investigate!

> Troll 1: If we are cousins, then we are both knaves.

Troll 2: We are cousins or we are both knaves.


### Subsection  Truth Tables

> If you get more doubles than any other player then you will lose, or if you lose then you must have bought the most properties.


| \(P\) | \(Q\) | \(P\wedge Q\) |
| T | T | T |
| T | F | F |
| F | T | F |
| F | F | F |

| \(P\) | \(Q\) | \(P\vee Q\) |
| T | T | T |
| T | F | T |
| F | T | T |
| F | F | F |

| \(P\) | \(Q\) | \(P\imp Q\) |
| T | T | T |
| T | F | F |
| F | T | T |
| F | F | T |

| \(P\) | \(Q\) | \(P\iff Q\) |
| T | T | T |
| T | F | F |
| F | T | F |
| F | F | T |

| \(P\) | \(\neg P\) |
| T | F |
| F | T |
|  |

#### Example 3.1.1.

| \(P\) | \(Q\) | \(\neg P\) | \(\neg P \vee Q\) |
| T | T | F | T |
| T | F | F | F |
| F | T | T | T |
| F | F | T | T |

#### Example 3.1.2.

| \(P\) | \(Q\) | \(R\) | \(P \imp Q\) | \(Q \imp R\) | \((P \imp Q) \vee (Q \imp R)\) |
| T | T | T | T | T | T |
| T | T | F | T | F | T |
| T | F | T | F | T | T |
| T | F | F | F | T | T |
| F | T | T | T | T | T |
| F | T | F | T | F | T |
| F | F | T | T | T | T |
| F | F | F | T | T | T |

### Subsection  Logical Equivalence

| \(P\) | \(Q\) | \(P \imp Q\) | \(\neg P \vee Q\) |
| T | T | T | T |
| T | F | F | F |
| F | T | T | T |
| F | F | T | T |

#### Logical Equivalence.

#### Example 3.1.3.

| \(P\) | \(Q\) | \(\neg(P \vee Q)\) | \(\neg P \wedge \neg Q\) |
| T | T | F | F |
| T | F | F | F |
| F | T | F | F |
| F | F | T | T |

#### De Morgan’s Laws.

#### Implications are Disjunctions.

#### Double Negation.

#### Example 3.1.4.

#### Negation of an Implication.

#### Example 3.1.5.

| \(P\) | \(Q\) | \(R\) | \((P\vee Q) \imp R\) | \((P\imp R) \vee (Q \imp R)\) |
| T | T | T | T | T |
| T | T | F | F | F |
| T | F | T | T | T |
| T | F | F | F | T |
| F | T | T | T | T |
| F | T | F | F | T |
| F | F | T | T | T |
| F | F | F | T | T |
|  |

### Subsection  Deductions

###### Investigate!

> If Edith eats her vegetables, then she can have a cookie. Edith ate her vegetables. Therefore Edith gets a cookie.


|  | \(P \imp Q\) |
|  | \(P\) |
| \(\therefore\) | \(Q\) |

| \(P\) | \(Q\) | \(P\imp Q\) |
| T | T | T |
| T | F | F |
| F | T | T |
| F | F | T |

#### Example 3.1.6.

|  | \(P \imp Q\) |
|  | \(\neg P \imp Q\) |
| \(\therefore\) | \(Q\) |

| \(P\) | \(Q\) | \(P\imp Q\) | \(\neg P\) | \(\neg P \imp Q\) |
| T | T | T | F | T |
| T | F | F | F | T |
| F | T | T | T | T |
| F | F | T | T | F |

#### Example 3.1.7.

|  | \(P \imp R\) |
|  | \(Q \imp R\) |
|  | \(R\) |
| \(\therefore\) | \(P \vee Q\) |

| \(P\) | \(Q\) | \(R\) | \(P \imp R\) | \(Q \imp R\) | \(P \vee Q\) |
| T | T | T | T | T | T |
| T | T | F | F | F | T |
| T | F | T | T | T | T |
| T | F | F | F | T | T |
| F | T | T | T | T | T |
| F | T | F | T | F | T |
| F | F | T | T | T | F |
| F | F | F | T | T | F |

|  | \(P \imp R\) |
|  | \(Q \imp R\) |
|  | \(P \vee Q\) |
| \(\therefore\) | \(R\) |

### Subsection  Beyond Propositions

> All primes greater than 2 are odd.


#### Example 3.1.8.

#### Example 3.1.9.

|  | \(\exists y \forall x P(x,y)\) |
| \(\therefore\) | \(\forall x \exists y P(x,y)\) |

### Exercises  Exercises

#### 1.

1. Translate the above statement into symbols. Clearly state which statement is \(P\) and which is \(Q\text{.}\)

1. Make a truth table for the statement.

1. Assuming the statement is true, what (if anything) can you conclude if there will be cake?

1. Assuming the statement is true, what (if anything) can you conclude if there will not be cake?

1. Suppose you found out that the statement was a lie. What can you conclude?

1. \(P\text{:}\) it’s your birthday; \(Q\text{:}\) there will be cake. \((P \vee Q) \imp Q\)

1. Hint: you should get three T’s and one F.

1. Only that there will be cake.

1. It’s NOT your birthday!

1. It’s your birthday, but the cake is a lie.

#### 2.

| \(P\) | \(Q\) | \(P \wedge Q\) | \(P \vee Q\) | \((P \wedge Q) \rightarrow (P \vee Q))\) |
| T | T |  |  |  |
| T | F |  |  |  |
| F | T |  |  |  |
| F | F |  |  |  |

| \(P\) | \(Q\) | \(P \wedge Q\) | \(P \vee Q\) | \((P \wedge Q) \rightarrow (P \vee Q))\) |
| T | T | T | T | T |
| T | F | F | T | T |
| F | T | F | T | T |
| F | F | F | F | T |

#### 3.

| \(P\) | \(Q\) | \(\neg Q\) | \(Q \rightarrow P\) | \(\neg Q \vee (Q \rightarrow P))\) |
| T | T |  |  |  |
| T | F |  |  |  |
| F | T |  |  |  |
| F | F |  |  |  |

- That \(P\) and \(Q\) are both true.

- That \(P\) is true and \(Q\) is false.

- That \(P\) is false and \(Q\) is true.

- That \(P\) and \(Q\) are both false.

- None of the above.

| \(P\) | \(Q\) | \(\neg Q\) | \(Q \rightarrow P\) | \(\neg Q \vee (Q \rightarrow P))\) |
| T | T | F | T | T |
| T | F | T | T | T |
| F | T | F | F | F |
| F | F | T | T | T |

#### 4.

| \(P\) | \(Q\) | \(R\) | \(\neg Q\) | \(R \rightarrow \neg Q\) | \(P \vee (R \rightarrow \neg Q)\) |
| T | T | T |  |  |  |
| T | T | F |  |  |  |
| T | F | T |  |  |  |
| T | F | F |  |  |  |
| F | T | T |  |  |  |
| F | T | F |  |  |  |
| F | F | T |  |  |  |
| F | F | F |  |  |  |

| \(P\) | \(Q\) | \(R\) | \(\neg Q\) | \(R \rightarrow \neg Q\) | \(P \vee (R \rightarrow \neg Q)\) |
| T | T | T | F | F | T |
| T | T | F | F | T | T |
| T | F | T | T | T | T |
| T | F | F | T | T | T |
| F | T | T | F | F | F |
| F | T | F | F | T | T |
| F | F | T | T | T | T |
| F | F | F | T | T | T |

#### 5.

1. Translate Geoff’s order into logical symbols.

1. The waiter knows that Geoff is either a liar or a truth-teller (so either everything he says is false, or everything is true). Which is it?

1. What, if anything, can the waiter conclude about the ingredients in Geoff’s desired calzone?

#### 6.

#### 7.

| \(P\) | \(Q\) | \(R\) | \(P \rightarrow (Q \vee R)\) | \((P \rightarrow Q) \vee (P\rightarrow R)\) |
| T | T | T |  |  |
| T | T | F |  |  |
| T | F | T |  |  |
| T | F | F |  |  |
| F | T | T |  |  |
| F | T | F |  |  |
| F | F | T |  |  |
| F | F | F |  |  |

- Yes, because the columns for the two statements are identical.

- Yes, because even though the columns are not identical, there are some rows in which they are identical.

- No, because the statements are not always true.

- No, because the columns for the two statements are not identical.

- Impossible to determine without more information.

| \(P\) | \(Q\) | \(R\) | \(P \rightarrow (Q \vee R)\) | \((P \rightarrow Q) \vee (P\rightarrow R)\) |
| T | T | T | T | T |
| T | T | F | T | T |
| T | F | T | T | T |
| T | F | F | F | F |
| F | T | T | T | T |
| F | T | F | T | T |
| F | F | T | T | T |
| F | F | F | T | T |

#### 8.

1. \(\neg(P
\imp \neg Q)\text{.}\)
1. \((\neg P \vee
\neg Q) \imp \neg (\neg Q \wedge R)\text{.}\)
1. \(\neg((P \imp \neg Q) \vee \neg (R \wedge
\neg
R))\text{.}\)
1. It is false that if Sam is not a man then Chris is a woman, and that Chris is not a woman.

1. \(P
\wedge Q\text{.}\)
1. \((\neg P \vee \neg R) \imp (Q \vee \neg R)\) or, replacing the implication with a disjunction first: \((P \wedge Q) \vee (Q \vee \neg R)\text{.}\)

1. \((P
\wedge Q) \wedge (R \wedge \neg R)\text{.}\) This is necessarily false, so it is also equivalent to \(P \wedge \neg P\text{.}\)

1. Either Sam is a woman and Chris is a man, or Chris is a woman.

#### 9.

1. \(\neg((\neg P \wedge Q) \vee \neg(R \vee \neg S))\text{.}\)
1. \(\neg((\neg P \imp \neg
Q) \wedge (\neg Q \imp R))\) (careful with the implications).
1. For both parts above, verify your answers are correct using truth tables. That is, use a truth table to check that the given statement and your proposed simplification are actually logically equivalent.

#### 10.

1. Make a truth table for the statement \((T \vee S) \imp \neg P\text{.}\)

1. If you believed the statement was false, what properties would a counterexample need to possess? Explain by referencing your truth table.

1. If the statement were true, what could you conclude about the number 5657, which is definitely prime? Again, explain using the truth table.

1. There will be three rows in which the statement is false.

1. Consider the three rows that evaluate to false and say what the truth values of \(T\text{,}\) \(S\text{,}\) and \(P\) are there.

1. You are looking for a row in which \(P\) is true, and the whole statement is true.

#### 11.

#### 12.

|  | \(P \rightarrow Q\) |
|  | \(\neg Q\) |
| \(\therefore\) | \(\neg P\) |

| \(P\) | \(Q\) | \(P \rightarrow Q\) | \(\neg Q\) | \(\neg P\) |
| T | T |  |  |  |
| T | F |  |  |  |
| F | T |  |  |  |
| F | F |  |  |  |

- No, because the conclusion is not always true.

- Yes, because there is a row in which both premises are true.

- No, because the columns for the two premises are not identical.

- Yes, in every row where both premises are true, the conclusion is also true.

- Impossible to determine without more information.

| \(P\) | \(Q\) | \(P \rightarrow Q\) | \(\neg Q\) | \(\neg P\) |
| T | T | T | F | F |
| T | F | F | T | F |
| F | T | T | F | T |
| F | F | T | T | T |

#### 13.

|  | \(P \rightarrow (Q \vee R)\) |
|  | \(\neg(P \rightarrow Q)\) |
| \(\therefore\) | \(R\) |

| \(P\) | \(Q\) | \(R\) | \(P \rightarrow (Q \vee R)\) | \(\neg(P \rightarrow Q)\) |
| T | T | T |  |  |
| T | T | F |  |  |
| T | F | T |  |  |
| T | F | F |  |  |
| F | T | T |  |  |
| F | T | F |  |  |
| F | F | T |  |  |
| F | F | F |  |  |

- Yes, because there is a row in which both premises are true.

- No, because the statements are not always true.

- Yes, in every row where both premises are true, the conclusion is also true.

- No, because the columns for the two premises are not identical.

- Impossible to determine without more information.

| \(P\) | \(Q\) | \(R\) | \(P \rightarrow (Q \vee R)\) | \(\neg(P \rightarrow Q)\) |
| T | T | T | T | F |
| T | T | F | T | F |
| T | F | T | T | T |
| T | F | F | F | T |
| F | T | T | T | F |
| F | T | F | T | F |
| F | F | T | T | F |
| F | F | F | T | F |

#### 14.

|  | \((P \wedge Q) \rightarrow R\) |
|  | \(\neg P \vee \neg Q\) |
| \(\therefore\) | \(\neg R\) |

| \(P\) | \(Q\) | \(R\) | \((P \wedge Q) \rightarrow R\) | \(\neg P \vee \neg Q\) | \(\neg R\) |
| T | T | T |  |  |  |
| T | T | F |  |  |  |
| T | F | T |  |  |  |
| T | F | F |  |  |  |
| F | T | T |  |  |  |
| F | T | F |  |  |  |
| F | F | T |  |  |  |
| F | F | F |  |  |  |

- Yes, because there is a row in which the conclusion and both premises are true.

- Yes, because in every row that the conclusion is true, one of the premises is true.

- No, because the columns for the two premises are not identical.

- No, because there is a row in which both premises are true but the conclusion is false.

- Impossible to determine without more information.

| \(P\) | \(Q\) | \(R\) | \((P \wedge Q) \rightarrow R\) | \(\neg P \vee \neg Q\) | \(\neg R\) |
| T | T | T | T | F | F |
| T | T | F | F | F | T |
| T | F | T | T | T | F |
| T | F | F | T | T | T |
| F | T | T | T | T | F |
| F | T | F | T | T | T |
| F | F | T | T | T | F |
| F | F | F | T | T | T |

#### 15.

1. Prove that the following is a valid deduction rule:




\(P
\imp Q\)



\(Q \imp R\)


\(\therefore\)
\(P \imp R\)


|  | \(P
\imp Q\) |
|  | \(Q \imp R\) |
| \(\therefore\) | \(P \imp R\) |

1. Prove that the following is a valid deduction rule for any \(n \ge 2\text{:}\)




\(P_1 \imp P_2\)



\(P_2
\imp P_3\)



\(\vdots\)



\(P_{n-1}
\imp P_n\)


\(\therefore\)
\(P_1
\imp P_n\text{.}\)


I suggest you don’t go through the trouble of writing out a \(2^n\) row truth table. Instead, you should use part (a) and mathematical induction.

|  | \(P_1 \imp P_2\) |
|  | \(P_2
\imp P_3\) |
|  | \(\vdots\) |
|  | \(P_{n-1}
\imp P_n\) |
| \(\therefore\) | \(P_1
\imp P_n\text{.}\) |

#### 16.

1. \(\neg \exists x \forall y (\neg O(x) \vee E(y))\text{.}\)
1. \(\neg \forall x \neg
\forall y \neg(x \lt y \wedge \exists z (x \lt z \vee y \lt z))\text{.}\)
1. There is a number \(n\) for which no other number is either less \(n\) than or equal to \(n\text{.}\)

1. It is false that for every number \(n\) there are two other numbers which \(n\) is between.

1. \(\forall
x \exists y (O(x) \wedge \neg E(y))\text{.}\)
1. \(\exists x \forall y (x \ge y \vee \forall
z (x \ge z \wedge y \ge z))\text{.}\)
1. There is a number \(n\) for which every other number is strictly greater than \(n\text{.}\)

1. There is a number \(n\) which is not between any other two numbers.

#### 17.

1. \(\neg \forall x \forall y (x \lt y \vee y \lt x)\text{.}\)

1. \(\neg(\exists
x P(x) \imp \forall y P(y))\text{.}\)

#### 18.

1. Every number is either even or odd.

1. There is a sequence that is both arithmetic and geometric.

1. For all numbers \(n\text{,}\) if \(n\) is prime, then \(n+3\) is not prime.

#### 19.

#### 20.

|  | \(P_1\) |
|  | \(P_2\) |
|  | \(\vdots\) |
|  | \(P_n\) |
| \(\therefore\) | \(Q\) |

#### 21.

|  | \(P \rightarrow Q\) |
|  | \(P \wedge \neg Q\) |
| \(\therefore\) | \(R\) |

| \(P\) | \(Q\) | \(R\) | \(P \rightarrow Q\) | \(P \wedge \neg Q\) |
| T | T | T |  |  |
| T | T | F |  |  |
| T | F | T |  |  |
| T | F | F |  |  |
| F | T | T |  |  |
| F | T | F |  |  |
| F | F | T |  |  |
| F | F | F |  |  |

- Yes, in every row where both premises are true, the conclusion is also true.

- No, because the premises are never both true in the same row.

- Yes, because there is a row in which both premises are true.

- No, because the columns for the two premises are not identical.

- Impossible to determine without more information.

| \(P\) | \(Q\) | \(R\) | \(P \rightarrow Q\) | \(P \wedge \neg Q\) |
| T | T | T | T | F |
| T | T | F | T | F |
| T | F | T | F | T |
| T | F | F | F | T |
| F | T | T | T | F |
| F | T | F | T | F |
| F | F | T | T | F |
| F | F | F | T | F |