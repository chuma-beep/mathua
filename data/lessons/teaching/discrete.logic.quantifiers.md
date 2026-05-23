> Content sourced from [Discrete Mathematics: An Open Introduction, 3rd ed](https://discrete.openmathbooks.org/dmoi3/sec_intro-statements.html) by Oscar Levin — CC BY-SA 4.0

# Predicates and Quantifiers

## Section 0.2 Mathematical Statements

###### Investigate!

> Troll 1: If I am a knave, then there are exactly two knights here.ð

Troll 2: Troll 1 is lying.ð

Troll 3: Either we are all knaves or at least one of us is a knight.ð

ð

### Subsection  Atomic and Molecular Statements

#### Example 0.2.1.

- Telephone numbers in the USA have 10 digits.ð

ð
- The moon is made of cheese.ð

ð
- 42 is a perfect square.ð

ð
- Every even number greater than 2 can be expressed as the sum of two primes.ð

ð
- \(\displaystyle 3+7 = 12\)ð

ð
- Would you like some cake?ð

ð
- The sum of two squares.ð

ð
- \(1+3+5+7+\cdots+2n+1\text{.}\)
ð
- Go to your room!ð

ð
- \(\displaystyle 3+x = 12\)ð

ð
> Telephone numbers in the USA have 10 digits and 42 is a perfect square.ð

ð

#### Logical Connectives.

- \(P \wedge Q\) is read â\(P\) and \(Q\text{,}\)â and called a conjunction. 
ð
- \(P \vee Q\) is read â\(P\) or \(Q\text{,}\)â and called a disjunction. 
ð
- \(P \imp Q\) is read âif \(P\) then \(Q\text{,}\)â and called an implication or conditional. 
ð
- \(P \iff Q\) is read â\(P\) if and only if \(Q\text{,}\)â and called a biconditional. 
ð
- \(\neg P\) is read ânot \(P\text{,}\)â and called a negation. 
ð
#### Truth Conditions for Connectives.

- \(P \wedge Q\) is true when both \(P\) and \(Q\) are true. 
ð
- \(P \vee Q\) is true when \(P\) or \(Q\) or both are true. 
ð
- \(P \imp Q\) is true when \(P\) is false or \(Q\) is true or both. 
ð
- \(P \iff Q\) is true when \(P\) and \(Q\) are both true, or both false. 
ð
- \(\neg P\) is true when \(P\) is false. 
ð
### Subsection  Implications

#### Implications.

- \(P\) is the hypothesis (or antecedent).
ð
- \(Q\) is the conclusion (or consequent).
ð
> If \(a\) and \(b\) are the legs of a right triangle with hypotenuse \(c\text{,}\) then \(a^2 + b^2 = c^2\text{.}\)ð

ð

#### Example 0.2.2.

> If Bob gets a 90 on the final, then Bob will pass the class.ð

ð

#### Example 0.2.3.

1. If \(1=1\text{,}\) then most horses have 4 legs.ð

ð
1. If \(0=1\text{,}\) then \(1=1\text{.}\)ð

ð
1. If 8 is a prime number, then the 7624th digit of \(\pi\) is an 8.ð

ð
1. If the 7624th digit of \(\pi\) is an 8, then \(2+2 = 4\text{.}\)ð

ð
1. Here both the hypothesis and the conclusion are true, so the implication is true. It does not matter that there is no meaningful connection between the true mathematical fact and the fact about horses.ð

ð
1. Here the hypothesis is false and the conclusion is true, so the implication is true.ð

ð
1. I have no idea what the 7624th digit of \(\pi\) is, but this does not matter. Since the hypothesis is false, the implication is automatically true.ð

ð
1. Similarly here, regardless of the truth value of the hypothesis, the conclusion is true, making the implication true.ð

ð
#### Direct Proofs of Implications.

#### Example 0.2.4.

##### Proof.

#### Converse and Contrapositive.

- The converse  of an implication \(P \imp Q\) is the implication \(Q \imp P\text{.}\) The converse is NOT logically equivalent to the original implication. That is, whether the converse of an implication is true is independent of the truth of the implication.ð

ð
- The contrapositive  of an implication \(P \imp Q\) is the statement \(\neg Q \imp \neg P\text{.}\) An implication and its contrapositive are logically equivalent (they are either both true or both false).ð

ð
#### Example 0.2.5.

#### Example 0.2.6.

1. Sue gets a 93% on her final.ð

ð
1. Sue gets an A in the class.ð

ð
1. Sue does not get a 93% on her final.ð

ð
1. Sue does not get an A in the class.ð

ð
1. We have \(P \imp Q\) and \(P\text{,}\) so \(Q\) follows. Sue gets an A.ð

ð
1. You cannot conclude anything. Sue could have gotten the A because she did extra credit for example. Notice that we do not know that if Sue gets an \(A\text{,}\) then she gets a 93% on her final. That is the converse of the original implication, so it might or might not be true.ð

ð
1. The contrapositive of the converse of \(P \imp Q\) is \(\neg P \imp \neg Q\text{,}\) which states that if Sue does not get a 93% on the final, then she will not get an A in the class. But this does not follow from the original implication. Again, we can conclude nothing. Sue could have done extra credit.ð

ð
1. What would happen if Sue does not get an A but did get a 93% on the final? Then \(P\) would be true and \(Q\) would be false. This makes the implication \(P \imp Q\) false! It must be that Sue did not get a 93% on the final. Notice now we have the implication \(\neg Q \imp \neg P\) which is the contrapositive of \(P \imp Q\text{.}\) Since \(P \imp Q\) is assumed to be true, we know \(\neg Q \imp \neg P\) is true as well.ð

ð
#### If and only if.

> \(P \iff Q\) is logically equivalent to \((P \imp Q) \wedge (Q \imp P)\text{.}\)ð

ð

#### Example 0.2.7.

#### Example 0.2.8.

1. I am asleep if I dream.ð

ð
1. I dream only if I am asleep.ð

ð
1. In order to dream, I must be asleep.ð

ð
1. To dream, it is necessary that I am asleep.ð

ð
1. To be asleep, it is sufficient to dream.ð

ð
1. I am not dreaming unless I am asleep.ð

ð
1. I dream if I am asleep.ð

ð
1. I am asleep only if I dream.ð

ð
1. It is necessary that I dream in order to be asleep.ð

ð
1. It is sufficient that I be asleep in order to dream.ð

ð
1. If I donât dream, then Iâm not asleep.ð

ð
#### Necessary and Sufficient.

- â\(P\) is necessary for \(Q\)â means \(Q \imp P\text{.}\)

ð
- â\(P\) is sufficient for \(Q\)â means \(P \imp Q\text{.}\)

ð
- If \(P\) is necessary and sufficient for \(Q\text{,}\) then \(P \iff Q\text{.}\)ð

ð
#### Example 0.2.9.

### Subsection  Predicates and Quantifiers

###### Investigate!

1. You can fool some people all of the time.ð

ð
1. You can fool everyone some of the time.ð

ð
1. You can always fool some people.ð

ð
1. Sometimes you can fool everyone.ð

ð
#### Universal and Existential Quantifiers.

#### Quantifiers and Negation.

> \(\neg \forall x P(x)\) is equivalent to \(\exists x \neg P(x)\text{.}\)ð


\(\neg \exists x P(x)\) is equivalent to \(\forall x \neg P(x)
\text{.}\)ð

ð

#### Implicit Quantifiers.

### Exercises  Exercises

#### 1.

1. Eat your vegetables!ð

ð
1. Sally ate her vegetables.ð

ð
1. Sally ate her vegetables and got a cookie.ð

ð
#### 2.

1. Everybody can be fooled sometimes.ð

ð
1. The Broncos will win the Super Bowl, or Iâll eat my hat.ð

ð
1. If a set contains three elements, then the sum of those elements is at least 6.ð

ð
1. Every even number is divisible by 2ð

ð
1. Every natural number greater than 1 is either prime or composite.ð

ð
#### 3.

1. Translate âJack and Jill both passed mathâ into symbols.ð

ð
1. Translate âIf Jack passed math, then Jill did notâ into symbols.ð

ð
1. Translate â\(P \vee Q\)â into English.ð

ð
1. Translate â\(\neg(P \wedge Q) \imp Q\)â into English.ð

ð
1. Suppose you know that if Jack passed math, then so did Jill. What can you conclude if you know that:


Jill passed math?
ð


Jill did not pass math?
ð


ð

ð
1. Jill passed math?
ð
1. Jill did not pass math?
ð
1. \(P
\wedge Q\text{.}\)
ð
1. \(P \imp \neg Q\text{.}\)
ð
1. Jack passed math or Jill passed math (or both).ð

ð
1. If Jack and Jill did not both pass math, then Jill did.ð

ð
1. Nothing else.
ð


Jack did not pass math either.
ð


ð

ð
1. Nothing else.
ð
1. Jack did not pass math either.
ð
#### 4.

1. 7 is prime, and 16 is not prime.ð

ð
1. If 7 is not prime, then 7 is my favorite number.ð

ð
1. If 16 is my favorite number, then \(16+1\) is my favorite number.ð

ð
1. 16 is prime, or 7 is prime.ð

ð
1. 17 is my favorite number, and 16 is not prime.ð

ð
1. If 16 is not prime, then 16 is my favorite number.ð

ð
#### 5.

1. If the circle is purple, then the square is yellow.ð

ð
1. The square is not yellow, or the circle is purple.ð

ð
1. The square and the circle are both purple.ð

ð
1. The square and the circle are both yellow.ð

ð
1. If the circle is not purple, then the square is not yellow.ð

ð
#### 6.

1. The diamond is purple if and only if the circle is blue.ð

ð
1. The diamond is purple.ð

ð
1. The circle is blue.ð

ð
1. The diamond is purple if and only if the circle is not blue.ð

ð
#### 7.

1. If you will not give me a cow, then I will not give you magic beans.ð

ð
1. If you will give me a cow, then I will not give you magic beans.ð

ð
1. If I will give you magic beans, then you will not give me a cow.ð

ð
1. If I will give you magic beans, then you will give me a cow.ð

ð
1. You will give me a cow, and I will not give you magic beans.ð

ð
1. If I will not give you magic beans, then you will not give me a cow.ð

ð
#### 8.

1. Write the converse of the statement.ð

ð
1. Write the contrapositive of the statement.ð

ð
1. Is it possible for the contrapositive to be false? If it was, what would that tell you?ð

ð
1. Suppose the original statement is true, and that Oscar drinks milk. Can you conclude anything (about his eating Chinese food)? Explain.ð

ð
1. Suppose the original statement is true, and that Oscar does not drink milk. Can you conclude anything (about his eating Chinese food)? Explain.ð

ð
#### 9.

1. Only viscous graphs satisfy condition (V).ð

ð
1. Satisfying condition (V) is a sufficient condition for a graph to be viscous.ð

ð
1. A graph is viscous only if it satisfies condition (V).ð

ð
1. Every viscous graph satisfies condition (V).ð

ð
1. Satisfying condition (V) is a necessary condition for a graph to be viscous.ð

ð
#### 10.

1. To lose weight, you must exercise.ð

ð
1. To lose weight, all you need to do is exercise.ð

ð
1. Every American is patriotic.ð

ð
1. You are patriotic only if you are American.ð

ð
1. The set of rational numbers is a subset of the real numbers.ð

ð
1. A number is prime if it is not even.ð

ð
1. Either the Broncos will win the Super Bowl, or they wonât play in the Super Bowl.ð

ð
1. If you have lost weight, then you exercised.ð

ð
1. If you exercise, then you will lose weight.ð

ð
1. If you are American, then you are patriotic.ð

ð
1. If you are patriotic, then you are American.ð

ð
1. If a number is rational, then it is real.ð

ð
1. If a number is not even, then it is prime. (Or the contrapositive: if a number is not prime, then it is even.)ð

ð
1. If the Broncos donât win the Super Bowl, then they didnât play in the Super Bowl. Alternatively, if the Broncos play in the Super Bowl, then they will win the Super Bowl.ð

ð
#### 11.

1. You will be rich only if you win the lottery.ð

ð
1. You will be rich if you win the lottery.ð

ð
1. Either you win the lottery, or else you are not rich.ð

ð
1. If you are rich, you must have won the lottery.ð

ð
1. You will win the lottery if you are rich.ð

ð
#### 12.

1. Is \(P(15)\) true or false?


Trueð

ð


Falseð

ð


Neither (not a statement)ð

ð


ð

ð
1. Trueð

ð
1. Falseð

ð
1. Neither (not a statement)ð

ð
1. What, if anything, can you conclude about \(\exists x P(x)\) from the truth value of \(P(15)\text{?}\)ð





\(\exists x P(x)\) must be true.ð

ð



\(\exists x P(x)\) must be false.ð

ð



\(\exists x P(x)\) could be true or could be false.ð

ð


ð

ð
1. \(\exists x P(x)\) must be true.ð

ð
1. \(\exists x P(x)\) must be false.ð

ð
1. \(\exists x P(x)\) could be true or could be false.ð

ð
1. What, if anything, can you conclude about \(\forall x P(x)\) from the truth value of \(P(15)\text{?}\)ð





\(\forall x P(x)\) must be true.ð

ð



\(\forall x P(x)\) must be false.ð

ð



\(\forall x P(x)\) could be true or could be false.ð

ð


ð

ð
1. \(\forall x P(x)\) must be true.ð

ð
1. \(\forall x P(x)\) must be false.ð

ð
1. \(\forall x P(x)\) could be true or could be false.ð

ð
#### 13.

1. Is \(P(15)\) true or false?


Trueð

ð


Falseð

ð


Neither (not a statement)ð

ð


ð

ð
1. Trueð

ð
1. Falseð

ð
1. Neither (not a statement)ð

ð
1. What, if anything, can you conclude about \(\exists x P(x)\) from the truth value of \(P(15)\text{?}\)ð





\(\exists x P(x)\) must be true.ð

ð



\(\exists x P(x)\) must be false.ð

ð



\(\exists x P(x)\) could be true or could be false.ð

ð


ð

ð
1. \(\exists x P(x)\) must be true.ð

ð
1. \(\exists x P(x)\) must be false.ð

ð
1. \(\exists x P(x)\) could be true or could be false.ð

ð
1. What, if anything, can you conclude about \(\forall x P(x)\) from the truth value of \(P(15)\text{?}\)ð





\(\forall x P(x)\) must be true.ð

ð



\(\forall x P(x)\) must be false.ð

ð



\(\forall x P(x)\) could be true or could be false.ð

ð


ð

ð
1. \(\forall x P(x)\) must be true.ð

ð
1. \(\forall x P(x)\) must be false.ð

ð
1. \(\forall x P(x)\) could be true or could be false.ð

ð
#### 14.

1. What would you need to do to prove \(\forall
x P(x)\) is true?ð

ð
1. What would you need to do to prove \(\forall x P(x)\) is false?ð

ð
1. What would you need to do to prove \(\exists x P(x)\) is true?ð

ð
1. What would you need to do to prove \(\exists x P(x)\) is false?ð

ð
1. The claim that \(\forall x P(x)\) means that \(P(n)\) is true no matter what \(n\) you consider in the domain of discourse. Thus the only way to prove that \(\forall x
P(x)\) is true is to check or otherwise argue that \(P(n)\) is true for all \(n\) in the domain.ð

ð
1. To prove \(\forall x P(x)\) is false all you need is one example of an element in the domain for which \(P(n)\) is false. This is often called a counterexample.ð

ð
1. We are simply claiming that there is some element \(n\) in the domain of discourse for which \(P(n)\) is true. If you can find one such element, you have verified the claim.ð

ð
1. Here we are claiming that no element we find will make \(P(n)\) true. The only way to be sure of this is to verify that every element of the domain makes \(P(n)\) false. Note that the level of proof needed for this statement is the same as to prove that \(\forall x P(x)\) is true.ð

ð
#### 15.

|  | 1 | 2 | 3 | 4 |
| --- | --- | --- | --- | --- |
| 1 | T | F | F | F |
| 2 | F | T | T | F |
| 3 | T | T | T | T |
| 4 | F | F | F | F |

1. \(\displaystyle \exists x \forall y P(x,y)\text{.}\)ð

ð
1. \(\displaystyle \exists y \forall x P(x,y)\text{.}\)ð

ð
1. \(\displaystyle \forall x \exists y P(x,y)\text{.}\)ð

ð
1. \(\displaystyle \forall y \exists x P(x,y)\text{.}\)ð

ð
#### 16.

1. No number is both even and odd.ð

ð
1. One more than any even number is an odd number.ð

ð
1. There is prime number that is even.ð

ð
1. Between any two numbers there is a third number.ð

ð
1. There is no number between a number and one more than that number.ð

ð
1. \(\neg
\exists x (E(x) \wedge O(x))\text{.}\)
ð
1. \(\forall x (E(x) \imp O(x+1))\text{.}\)
ð
1. \(\exists
x(P(x) \wedge E(x))\) (where \(P(x)\) means â\(x\) is primeâ).
ð
1. \(\forall x
\forall y \exists z(x \lt z \lt y \vee y \lt z \lt x)\text{.}\)
ð
1. \(\forall x \neg \exists
y
(x \lt y \lt x+1)\text{.}\)
ð
#### 17.

1. \(\forall x (E(x) \imp E(x +2))\text{.}\)
ð
1. \(\forall x \exists y (\sin(x)
= y)\text{.}\)
ð
1. \(\forall
y \exists x (\sin(x) = y)\text{.}\)
ð
1. \(\forall x \forall y (x^3 = y^3 \imp x = y)\text{.}\)
ð
1. Any even number plus 2 is an even number.ð

ð
1. For any \(x\) there is a \(y\) such that \(\sin(x) = y\text{.}\) In other words, every number \(x\) is in the domain of sine.ð

ð
1. For every \(y\) there is an \(x\) such that \(\sin(x) = y\text{.}\) In other words, every number \(y\) is in the range of sine (which is false).ð

ð
1. For any numbers, if the cubes of two numbers are equal, then the numbers are equal.ð

ð
#### 18.

#### 19.

1. \(\forall x \exists y (y^2 = x)\text{.}\)
ð
1. \(\forall x \forall y (x \lt y \imp \exists z
(x \lt z \lt y))\text{.}\)
ð
1. \(\exists x \forall y \forall z (y \lt z \imp y \le x \le z)\text{.}\)
ð
#### 20.

1. Write the converse and the contrapositive of the statement, saying which is which. Note: the original statement claims that an implication is true for all \(n\text{,}\) and it is that implication that we are taking the converse and contrapositive of.ð

ð
1. Write the negation of the original statement. What would you need to show to prove that the statement is false?ð

ð
1. Even though you donât know whether 10 is solitary (in fact, nobody knows this), is the statement âif 10 is prime, then 10 is solitaryâ true or false? Explain.ð

ð
1. It turns out that 8 is solitary. Does this tell you anything about the truth or falsity of the original statement, its converse or its contrapositive? Explain.ð

ð
1. Assuming that the original statement is true, what can you say about the relationship between the set \(P\) of prime numbers and the set \(S\) of solitary numbers. Explain.ð

ð