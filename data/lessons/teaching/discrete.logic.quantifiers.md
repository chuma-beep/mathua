> Content sourced from [Discrete Mathematics: An Open Introduction, 3rd ed](https://discrete.openmathbooks.org/dmoi3/sec_intro-statements.html) by Oscar Levin — CC BY-SA 4.0

# Predicates and Quantifiers

## Section 0.2 Mathematical Statements

###### Investigate!

> Troll 1: If I am a knave, then there are exactly two knights here.

Troll 2: Troll 1 is lying.

Troll 3: Either we are all knaves or at least one of us is a knight.


### Subsection  Atomic and Molecular Statements

#### Example 0.2.1.

- Telephone numbers in the USA have 10 digits.

- The moon is made of cheese.

- 42 is a perfect square.

- Every even number greater than 2 can be expressed as the sum of two primes.

- \(\displaystyle 3+7 = 12\)

- Would you like some cake?

- The sum of two squares.

- \(1+3+5+7+\cdots+2n+1\text{.}\)
- Go to your room!

- \(\displaystyle 3+x = 12\)

> Telephone numbers in the USA have 10 digits and 42 is a perfect square.


#### Logical Connectives.

- \(P \wedge Q\) is read “\(P\) and \(Q\text{,}\)” and called a conjunction. 
- \(P \vee Q\) is read “\(P\) or \(Q\text{,}\)” and called a disjunction. 
- \(P \imp Q\) is read “if \(P\) then \(Q\text{,}\)” and called an implication or conditional. 
- \(P \iff Q\) is read “\(P\) if and only if \(Q\text{,}\)” and called a biconditional. 
- \(\neg P\) is read “not \(P\text{,}\)” and called a negation. 
#### Truth Conditions for Connectives.

- \(P \wedge Q\) is true when both \(P\) and \(Q\) are true. 
- \(P \vee Q\) is true when \(P\) or \(Q\) or both are true. 
- \(P \imp Q\) is true when \(P\) is false or \(Q\) is true or both. 
- \(P \iff Q\) is true when \(P\) and \(Q\) are both true, or both false. 
- \(\neg P\) is true when \(P\) is false. 
### Subsection  Implications

#### Implications.

- \(P\) is the hypothesis (or antecedent).
- \(Q\) is the conclusion (or consequent).
> If \(a\) and \(b\) are the legs of a right triangle with hypotenuse \(c\text{,}\) then \(a^2 + b^2 = c^2\text{.}\)


#### Example 0.2.2.

> If Bob gets a 90 on the final, then Bob will pass the class.


#### Example 0.2.3.

1. If \(1=1\text{,}\) then most horses have 4 legs.

1. If \(0=1\text{,}\) then \(1=1\text{.}\)

1. If 8 is a prime number, then the 7624th digit of \(\pi\) is an 8.

1. If the 7624th digit of \(\pi\) is an 8, then \(2+2 = 4\text{.}\)

1. Here both the hypothesis and the conclusion are true, so the implication is true. It does not matter that there is no meaningful connection between the true mathematical fact and the fact about horses.

1. Here the hypothesis is false and the conclusion is true, so the implication is true.

1. I have no idea what the 7624th digit of \(\pi\) is, but this does not matter. Since the hypothesis is false, the implication is automatically true.

1. Similarly here, regardless of the truth value of the hypothesis, the conclusion is true, making the implication true.

#### Direct Proofs of Implications.

#### Example 0.2.4.

##### Proof.

#### Converse and Contrapositive.

- The converse  of an implication \(P \imp Q\) is the implication \(Q \imp P\text{.}\) The converse is NOT logically equivalent to the original implication. That is, whether the converse of an implication is true is independent of the truth of the implication.

- The contrapositive  of an implication \(P \imp Q\) is the statement \(\neg Q \imp \neg P\text{.}\) An implication and its contrapositive are logically equivalent (they are either both true or both false).

#### Example 0.2.5.

#### Example 0.2.6.

1. Sue gets a 93% on her final.

1. Sue gets an A in the class.

1. Sue does not get a 93% on her final.

1. Sue does not get an A in the class.

1. We have \(P \imp Q\) and \(P\text{,}\) so \(Q\) follows. Sue gets an A.

1. You cannot conclude anything. Sue could have gotten the A because she did extra credit for example. Notice that we do not know that if Sue gets an \(A\text{,}\) then she gets a 93% on her final. That is the converse of the original implication, so it might or might not be true.

1. The contrapositive of the converse of \(P \imp Q\) is \(\neg P \imp \neg Q\text{,}\) which states that if Sue does not get a 93% on the final, then she will not get an A in the class. But this does not follow from the original implication. Again, we can conclude nothing. Sue could have done extra credit.

1. What would happen if Sue does not get an A but did get a 93% on the final? Then \(P\) would be true and \(Q\) would be false. This makes the implication \(P \imp Q\) false! It must be that Sue did not get a 93% on the final. Notice now we have the implication \(\neg Q \imp \neg P\) which is the contrapositive of \(P \imp Q\text{.}\) Since \(P \imp Q\) is assumed to be true, we know \(\neg Q \imp \neg P\) is true as well.

#### If and only if.

> \(P \iff Q\) is logically equivalent to \((P \imp Q) \wedge (Q \imp P)\text{.}\)


#### Example 0.2.7.

#### Example 0.2.8.

1. I am asleep if I dream.

1. I dream only if I am asleep.

1. In order to dream, I must be asleep.

1. To dream, it is necessary that I am asleep.

1. To be asleep, it is sufficient to dream.

1. I am not dreaming unless I am asleep.

1. I dream if I am asleep.

1. I am asleep only if I dream.

1. It is necessary that I dream in order to be asleep.

1. It is sufficient that I be asleep in order to dream.

1. If I don’t dream, then I’m not asleep.

#### Necessary and Sufficient.

- “\(P\) is necessary for \(Q\)” means \(Q \imp P\text{.}\)

- “\(P\) is sufficient for \(Q\)” means \(P \imp Q\text{.}\)

- If \(P\) is necessary and sufficient for \(Q\text{,}\) then \(P \iff Q\text{.}\)

#### Example 0.2.9.

### Subsection  Predicates and Quantifiers

###### Investigate!

1. You can fool some people all of the time.

1. You can fool everyone some of the time.

1. You can always fool some people.

1. Sometimes you can fool everyone.

#### Universal and Existential Quantifiers.

#### Quantifiers and Negation.

> \(\neg \forall x P(x)\) is equivalent to \(\exists x \neg P(x)\text{.}\)


\(\neg \exists x P(x)\) is equivalent to \(\forall x \neg P(x)\text{.}\)


#### Implicit Quantifiers.

### Exercises  Exercises

#### 1.

1. Eat your vegetables!

1. Sally ate her vegetables.

1. Sally ate her vegetables and got a cookie.

#### 2.

1. Everybody can be fooled sometimes.

1. The Broncos will win the Super Bowl, or I’ll eat my hat.

1. If a set contains three elements, then the sum of those elements is at least 6.

1. Every even number is divisible by 2

1. Every natural number greater than 1 is either prime or composite.

#### 3.

1. Translate “Jack and Jill both passed math” into symbols.

1. Translate “If Jack passed math, then Jill did not” into symbols.

1. Translate “\(P \vee Q\)” into English.

1. Translate “\(\neg(P \wedge Q) \imp Q\)” into English.

1. Suppose you know that if Jack passed math, then so did Jill. What can you conclude if you know that:


Jill passed math?


Jill did not pass math?



1. Jill passed math?
1. Jill did not pass math?
1. \(P
\wedge Q\text{.}\)
1. \(P \imp \neg Q\text{.}\)
1. Jack passed math or Jill passed math (or both).

1. If Jack and Jill did not both pass math, then Jill did.

1. Nothing else.


Jack did not pass math either.



1. Nothing else.
1. Jack did not pass math either.
#### 4.

1. 7 is prime, and 16 is not prime.

1. If 7 is not prime, then 7 is my favorite number.

1. If 16 is my favorite number, then \(16+1\) is my favorite number.

1. 16 is prime, or 7 is prime.

1. 17 is my favorite number, and 16 is not prime.

1. If 16 is not prime, then 16 is my favorite number.

#### 5.

1. If the circle is purple, then the square is yellow.

1. The square is not yellow, or the circle is purple.

1. The square and the circle are both purple.

1. The square and the circle are both yellow.

1. If the circle is not purple, then the square is not yellow.

#### 6.

1. The diamond is purple if and only if the circle is blue.

1. The diamond is purple.

1. The circle is blue.

1. The diamond is purple if and only if the circle is not blue.

#### 7.

1. If you will not give me a cow, then I will not give you magic beans.

1. If you will give me a cow, then I will not give you magic beans.

1. If I will give you magic beans, then you will not give me a cow.

1. If I will give you magic beans, then you will give me a cow.

1. You will give me a cow, and I will not give you magic beans.

1. If I will not give you magic beans, then you will not give me a cow.

#### 8.

1. Write the converse of the statement.

1. Write the contrapositive of the statement.

1. Is it possible for the contrapositive to be false? If it was, what would that tell you?

1. Suppose the original statement is true, and that Oscar drinks milk. Can you conclude anything (about his eating Chinese food)? Explain.

1. Suppose the original statement is true, and that Oscar does not drink milk. Can you conclude anything (about his eating Chinese food)? Explain.

#### 9.

1. Only viscous graphs satisfy condition (V).

1. Satisfying condition (V) is a sufficient condition for a graph to be viscous.

1. A graph is viscous only if it satisfies condition (V).

1. Every viscous graph satisfies condition (V).

1. Satisfying condition (V) is a necessary condition for a graph to be viscous.

#### 10.

1. To lose weight, you must exercise.

1. To lose weight, all you need to do is exercise.

1. Every American is patriotic.

1. You are patriotic only if you are American.

1. The set of rational numbers is a subset of the real numbers.

1. A number is prime if it is not even.

1. Either the Broncos will win the Super Bowl, or they won’t play in the Super Bowl.

1. If you have lost weight, then you exercised.

1. If you exercise, then you will lose weight.

1. If you are American, then you are patriotic.

1. If you are patriotic, then you are American.

1. If a number is rational, then it is real.

1. If a number is not even, then it is prime. (Or the contrapositive: if a number is not prime, then it is even.)

1. If the Broncos don’t win the Super Bowl, then they didn’t play in the Super Bowl. Alternatively, if the Broncos play in the Super Bowl, then they will win the Super Bowl.

#### 11.

1. You will be rich only if you win the lottery.

1. You will be rich if you win the lottery.

1. Either you win the lottery, or else you are not rich.

1. If you are rich, you must have won the lottery.

1. You will win the lottery if you are rich.

#### 12.

1. Is \(P(15)\) true or false?


True



False



Neither (not a statement)




1. True

1. False

1. Neither (not a statement)

1. What, if anything, can you conclude about \(\exists x P(x)\) from the truth value of \(P(15)\text{?}\)





\(\exists x P(x)\) must be true.




\(\exists x P(x)\) must be false.




\(\exists x P(x)\) could be true or could be false.




1. \(\exists x P(x)\) must be true.

1. \(\exists x P(x)\) must be false.

1. \(\exists x P(x)\) could be true or could be false.

1. What, if anything, can you conclude about \(\forall x P(x)\) from the truth value of \(P(15)\text{?}\)





\(\forall x P(x)\) must be true.




\(\forall x P(x)\) must be false.




\(\forall x P(x)\) could be true or could be false.




1. \(\forall x P(x)\) must be true.

1. \(\forall x P(x)\) must be false.

1. \(\forall x P(x)\) could be true or could be false.

#### 13.

1. Is \(P(15)\) true or false?


True



False



Neither (not a statement)




1. True

1. False

1. Neither (not a statement)

1. What, if anything, can you conclude about \(\exists x P(x)\) from the truth value of \(P(15)\text{?}\)





\(\exists x P(x)\) must be true.




\(\exists x P(x)\) must be false.




\(\exists x P(x)\) could be true or could be false.




1. \(\exists x P(x)\) must be true.

1. \(\exists x P(x)\) must be false.

1. \(\exists x P(x)\) could be true or could be false.

1. What, if anything, can you conclude about \(\forall x P(x)\) from the truth value of \(P(15)\text{?}\)





\(\forall x P(x)\) must be true.




\(\forall x P(x)\) must be false.




\(\forall x P(x)\) could be true or could be false.




1. \(\forall x P(x)\) must be true.

1. \(\forall x P(x)\) must be false.

1. \(\forall x P(x)\) could be true or could be false.

#### 14.

1. What would you need to do to prove \(\forall
x P(x)\) is true?

1. What would you need to do to prove \(\forall x P(x)\) is false?

1. What would you need to do to prove \(\exists x P(x)\) is true?

1. What would you need to do to prove \(\exists x P(x)\) is false?

1. The claim that \(\forall x P(x)\) means that \(P(n)\) is true no matter what \(n\) you consider in the domain of discourse. Thus the only way to prove that \(\forall x
P(x)\) is true is to check or otherwise argue that \(P(n)\) is true for all \(n\) in the domain.

1. To prove \(\forall x P(x)\) is false all you need is one example of an element in the domain for which \(P(n)\) is false. This is often called a counterexample.

1. We are simply claiming that there is some element \(n\) in the domain of discourse for which \(P(n)\) is true. If you can find one such element, you have verified the claim.

1. Here we are claiming that no element we find will make \(P(n)\) true. The only way to be sure of this is to verify that every element of the domain makes \(P(n)\) false. Note that the level of proof needed for this statement is the same as to prove that \(\forall x P(x)\) is true.

#### 15.

|  | 1 | 2 | 3 | 4 |
| --- | --- | --- | --- | --- |
| 1 | T | F | F | F |
| 2 | F | T | T | F |
| 3 | T | T | T | T |
| 4 | F | F | F | F |

1. \(\displaystyle \exists x \forall y P(x,y)\text{.}\)

1. \(\displaystyle \exists y \forall x P(x,y)\text{.}\)

1. \(\displaystyle \forall x \exists y P(x,y)\text{.}\)

1. \(\displaystyle \forall y \exists x P(x,y)\text{.}\)

#### 16.

1. No number is both even and odd.

1. One more than any even number is an odd number.

1. There is prime number that is even.

1. Between any two numbers there is a third number.

1. There is no number between a number and one more than that number.

1. \(\neg
\exists x (E(x) \wedge O(x))\text{.}\)
1. \(\forall x (E(x) \imp O(x+1))\text{.}\)
1. \(\exists
x(P(x) \wedge E(x))\) (where \(P(x)\) means “\(x\) is prime”).
1. \(\forall x
\forall y \exists z(x \lt z \lt y \vee y \lt z \lt x)\text{.}\)
1. \(\forall x \neg \exists
y
(x \lt y \lt x+1)\text{.}\)
#### 17.

1. \(\forall x (E(x) \imp E(x +2))\text{.}\)
1. \(\forall x \exists y (\sin(x)
= y)\text{.}\)
1. \(\forall
y \exists x (\sin(x) = y)\text{.}\)
1. \(\forall x \forall y (x^3 = y^3 \imp x = y)\text{.}\)
1. Any even number plus 2 is an even number.

1. For any \(x\) there is a \(y\) such that \(\sin(x) = y\text{.}\) In other words, every number \(x\) is in the domain of sine.

1. For every \(y\) there is an \(x\) such that \(\sin(x) = y\text{.}\) In other words, every number \(y\) is in the range of sine (which is false).

1. For any numbers, if the cubes of two numbers are equal, then the numbers are equal.

#### 18.

#### 19.

1. \(\forall x \exists y (y^2 = x)\text{.}\)
1. \(\forall x \forall y (x \lt y \imp \exists z
(x \lt z \lt y))\text{.}\)
1. \(\exists x \forall y \forall z (y \lt z \imp y \le x \le z)\text{.}\)
#### 20.

1. Write the converse and the contrapositive of the statement, saying which is which. Note: the original statement claims that an implication is true for all \(n\text{,}\) and it is that implication that we are taking the converse and contrapositive of.

1. Write the negation of the original statement. What would you need to show to prove that the statement is false?

1. Even though you don’t know whether 10 is solitary (in fact, nobody knows this), is the statement “if 10 is prime, then 10 is solitary” true or false? Explain.

1. It turns out that 8 is solitary. Does this tell you anything about the truth or falsity of the original statement, its converse or its contrapositive? Explain.

1. Assuming that the original statement is true, what can you say about the relationship between the set \(P\) of prime numbers and the set \(S\) of solitary numbers. Explain.
