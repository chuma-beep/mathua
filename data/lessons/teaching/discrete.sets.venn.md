> Content sourced from [Discrete Mathematics: An Open Introduction, 3rd ed](https://discrete.openmathbooks.org/dmoi3/sec_intro-sets.html) by Oscar Levin — CC BY-SA 4.0

# Venn Diagrams

## Section 0.3 Sets

### Subsection  Notation

#### Example 0.3.1.

1. \(\{x \st x + 3 \in \N\}\text{.}\)
ð
1. \(\{x \in \N \st x + 3 \in \N\}\text{.}\)
ð
1. \(\{x \st x \in \N \vee -x \in \N\}\text{.}\)
ð
1. \(\{x \st x \in \N \wedge -x \in \N\}\text{.}\)
ð
1. This is the set of all numbers which are 3 less than a natural number (i.e., that if you add 3 to them, you get a natural number). The set could also be written as \(\{-3, -2, -1, 0, 1, 2, \ldots\}\) (note that 0 is a natural number, so \(-3\) is in this set because \(-3 + 3 = 0\)).ð

ð
1. This is the set of all natural numbers which are 3 less than a natural number. So here we just have \(\{0, 1, 2,3 \ldots\}\text{.}\)ð

ð
1. This is the set of all integers  (positive and negative whole numbers, written \(\Z\)). In other words, \(\{\ldots, -2, -1, 0, 1, 2, \ldots\}\text{.}\)ð

ð
1. Here we want all numbers \(x\) such that \(x\) and \(-x\) are natural numbers. There is only one: 0. So we have the set \(\{0\}\text{.}\)ð

ð
#### Example 0.3.2.

1. \(\displaystyle A = \{x \in \Z \st x^2 \in \N\}\)ð

ð
1. \(\displaystyle B = \{x^2 \st x \in \N\}\)ð

ð
1. The set of integers that pass the condition that their square is a natural number.  Well, every integer, when you square it, gives you a non-negative integer, so a natural number.  Thus \(A = \Z = \{\ldots, -2, -1, 0, 1, 2, 3, \ldots\}\text{.}\)ð

ð
1. Here we are looking for the set of all \(x^2\)s where \(x\) is a natural number.  So this set is simply the set of perfect squares.  \(B = \{0, 1, 4, 9, 16, \ldots\}\text{.}\)ð

Another way we could have written this set, using more strict set builder notation, would be as \(B = \{x \in \N \st x = n^2 \text{ for some } n \in \N\}\text{.}\)ð

ð
#### Special sets.

#### Set Theory Notation.

###### Investigate!

1. Find the cardinality of each set below.


\(A = \{3,4,\ldots, 15\}\text{.}\)
ð


\(B = \{n \in \N \st 2 \lt n \le 200\}\text{.}\)
ð


\(C = \{n \le 100 \st n \in \N \wedge \exists m \in \N (n = 2m+1)\}\text{.}\)
ð


ð

ð
1. \(A = \{3,4,\ldots, 15\}\text{.}\)
ð
1. \(B = \{n \in \N \st 2 \lt n \le 200\}\text{.}\)
ð
1. \(C = \{n \le 100 \st n \in \N \wedge \exists m \in \N (n = 2m+1)\}\text{.}\)
ð
1. Find two sets \(A\) and \(B\) for which \(|A| = 5\text{,}\) \(|B| = 6\text{,}\) and \(|A\cup B| = 9\text{.}\) What is \(|A \cap B|\text{?}\)ð

ð
1. Find sets \(A\) and \(B\) with \(|A| = |B|\) such that \(|A\cup B| = 7\) and \(|A \cap B| = 3\text{.}\) What is \(|A|\text{?}\)ð

ð
1. Let \(A = \{1,2,\ldots, 10\}\text{.}\) Define \(\mathcal{B}_2 = \{B \subseteq A \st |B| = 2\}\text{.}\) Find \(|\mathcal{B}_2|\text{.}\)ð

ð
1. For any sets \(A\) and \(B\text{,}\) define \(AB = \{ab \st a\in A \wedge b \in B\}\text{.}\) If \(A = \{1,2\}\) and \(B = \{2,3,4\}\text{,}\) what is \(|AB|\text{?}\) What is \(|A \times B|\text{?}\)

ð
### Subsection  Relationships Between Sets

#### Example 0.3.3.

1. \(A \subset B\text{.}\)
ð
1. \(B \subset A\text{.}\)
ð
1. \(B \in C\text{.}\)
ð
1. \(\emptyset \in A\text{.}\)
ð
1. \(\emptyset \subset A\text{.}\)
ð
1. \(A \lt D\text{.}\)
ð
1. \(3 \in C\text{.}\)
ð
1. \(3 \subset C\text{.}\)
ð
1. \(\{3\} \subset C\text{.}\)
ð
1. False. For example, \(1\in A\) but \(1 \notin B\text{.}\)ð

ð
1. True. Every element in \(B\) is an element in \(A\text{.}\)ð

ð
1. False. The elements in \(C\) are 1, 2, and 3. The set \(B\) is not equal to 1, 2, or 3.ð

ð
1. False. \(A\) has exactly 6 elements, and none of them are the empty set.ð

ð
1. True. Everything in the empty set (nothing) is also an element of \(A\text{.}\) Notice that the empty set is a subset of every set.ð

ð
1. Meaningless. A set cannot be less than another set.ð

ð
1. True. \(3\) is one of the elements of the set \(C\text{.}\)ð

ð
1. Meaningless. \(3\) is not a set, so it cannot be a subset of another set.ð

ð
1. True. \(3\) is the only element of the set \(\{3\}\text{,}\) and is an element of \(C\text{,}\) so every element in \(\{3\}\) is an element of \(C\text{.}\)ð

ð
#### Example 0.3.4.

#### Example 0.3.5.

1. Find the cardinality of \(A = \{23, 24, \ldots, 37, 38\}\text{.}\)ð

ð
1. Find the cardinality of \(B = \{1, \{2, 3, 4\}, \emptyset\}\text{.}\)ð

ð
1. If \(C = \{1,2,3\}\text{,}\) what is the cardinality of \(\pow(C)\text{?}\)ð

ð
1. Since \(38 - 23 = 15\text{,}\) we can conclude that the cardinality of the set is \(|A| = 16\) (you need to add one since 23 is included).ð

ð
1. Here \(|B| = 3\text{.}\) The three elements are the number 1, the set \(\{2,3,4\}\text{,}\) and the empty set.ð

ð
1. We wrote out the elements of the power set \(\pow(C)\) above, and there are 8 elements (each of which is a set). So \(\card{\pow(C)} = 8\text{.}\) (You might wonder if there is a relationship between \(\card{A}\) and \(\card{\pow(A)}\) for all sets \(A\text{.}\) This is a good question which we will return to in ChapterÂ 1.)ð

ð
### Subsection  Operations On Sets

#### Example 0.3.6.

1. \(A \cup B\text{.}\)
ð
1. \(A \cap B\text{.}\)
ð
1. \(B \cap C\text{.}\)
ð
1. \(A \cap D\text{.}\)
ð
1. \(\bar{B \cup C}\text{.}\)
ð
1. \(A \setminus B\text{.}\)
ð
1. \((D \cap \bar C) \cup \bar{A \cap B}\text{.}\)
ð
1. \(\emptyset \cup C\text{.}\)
ð
1. \(\emptyset \cap C\text{.}\)
ð
1. \(A \cup B = \{1, 2, 3, 4, 5, 6\} = A\) since everything in \(B\) is already in \(A\text{.}\)

ð
1. \(A \cap B = \{2, 4, 6\} = B\) since everything in \(B\) is in \(A\text{.}\)

ð
1. \(B \cap C = \{2\}\) as the only element of both \(B\) and \(C\) is 2.
ð
1. \(A \cap D = \emptyset\) since \(A\) and \(D\) have no common elements.
ð
1. \(\bar{B \cup C} = \{5, 7, 8, 9, 10\}\text{.}\) First we find that \(B \cup C = \{1, 2, 3, 4, 6\}\text{,}\) then we take everything not in that set.
ð
1. \(A \setminus B = \{1, 3, 5\}\) since the elements 1, 3, and 5 are in \(A\) but not in \(B\text{.}\) This is the same as \(A \cap \bar B\text{.}\)

ð
1. \((D \cap \bar C) \cup \bar{A \cap B} = \{1, 3, 5, 7, 8, 9, 10\}\text{.}\) The set contains all elements that are either in \(D\) but not in \(C\) (i.e., \(\{7,8,9\}\)), or not in both \(A\) and \(B\) (i.e., \(\{1,3,5,7,8,9,10\}\)).
ð
1. \(\emptyset \cup C = C\) since nothing is added by the empty set.
ð
1. \(\emptyset \cap C = \emptyset\) since nothing can be both in a set and in the empty set.
ð
#### Example 0.3.7.

#### Example 0.3.8.

### Subsection  Venn Diagrams

### Exercises  Exercises

#### 1.

1. \(\displaystyle A \cup B\text{.}\)ð

ð
1. \(\displaystyle A \cap B\text{.}\)ð

ð
1. \(\displaystyle A \setminus B\text{.}\)ð

ð
1. \(\displaystyle B \setminus A\text{.}\)ð

ð
1. \(A \cup B = {\left\{3,5,6,7,8,9,12\right\}}\text{.}\) It includes everything that is in \(A\) or \(B\) or both.ð

ð
1. \(A \cap B = {\left\{3,6,7\right\}}\text{.}\)  It contains everything that is in both \(A\) and \(B\text{.}\)ð

ð
1. \(A \setminus B = {\left\{5,8\right\}}\text{.}\)  It contains everything that is in \(A\) except anything that is also in \(B\text{.}\)  We could also have written this set as \(A \cap \bar{B}\text{.}\)ð

ð
1. \(B \setminus A = {\left\{9,12\right\}}\text{.}\) It contains everything in \(B\) except anything that is also in \(A\text{.}\)  Another way to write this is \(B \cap \bar{A}\text{.}\)  Note that \(A \setminus B \ne B \setminus A\text{.}\)ð

ð
#### 2.

1. \(\renewcommand{\d}{\displaystyle}\newcommand{\N}{\mathbb N}\newcommand{\st}{:}\{n \in \N \st n^2 - 4 \ge 5\}\text{.}\) .ð

ð
1. \(\renewcommand{\d}{\displaystyle}\newcommand{\N}{\mathbb N}\newcommand{\st}{:}\{n \in \N \st n^2 - 9 \in \N\}\text{.}\) ð

ð
1. \(\renewcommand{\d}{\displaystyle}\newcommand{\N}{\mathbb N}\newcommand{\st}{:}\{n^2+1 \st n \in \N\}\text{.}\) ð

ð
1. \(\renewcommand{\d}{\displaystyle}\newcommand{\N}{\mathbb N}\newcommand{\st}{:}\{n \in \N \st n = k^2 + 1 \text{ for some } k \in \N\}\text{.}\) ð

ð
1. This is the set \(\{3, 4, 5 \ldots \}\) since we need each element to be a natural number whose square is at least 4 more than 5. Since \(3^2 - 4 = 5\) but \(2^2 - 4 = 0\text{,}\) we see that the first such natural number is 3.ð

ð
1. We get the same set as we did in the previous part, and the smallest non-negative number for which \(n^2 - 9\) is a natural number is 3.ð

Note that if we didnât specify \(n \ge 0\) by saying that \(n \in \mathbb N\text{,}\) then any integer less than \(-3\) would also be in the set, so there would not be a least element.ð

ð
1. This is the set \(\{1, 2, 5, 10, \ldots\}\text{,}\) namely the set of numbers that are the result of squaring and adding 1 to a natural number.  (\(0^2 + 1 = 1\text{,}\) \(1^2 + 1 = 2\text{,}\) \(2^2 + 1 = 5\) and so on.)  Thus the least element of the set is 1.ð

ð
1. Now we are looking for natural numbers that are equal to taking some natural number, squaring it, and adding 1.  That is, \(\{1, 2, 5, 10, \ldots\}\text{,}\) the same set as the previous part.  So again, the least element is 1.ð

ð
#### 3.

1. \(|A|\) when \(A = \{5, 6,7,8,\ldots,35\}\text{.}\) ð

ð
1. \(|A|\) when \(\newcommand{\Z}{\mathbb Z}\newcommand{\st}{:}A = \{x \in \Z \st -2 \le x \le 97\}\text{.}\) ð

ð
1. \(|A \cap B|\) when \(\renewcommand{\d}{\displaystyle}\newcommand{\N}{\mathbb N}\newcommand{\st}{:}A = \{x \in \N \st x \le 37\}\) and \(\renewcommand{\d}{\displaystyle}\newcommand{\N}{\mathbb N}\newcommand{\st}{:}B = \{x \in \N \st x \mbox{ is prime} \}\) ð

ð
1. \(31 = 35 - 5 + 1\text{.}\)  We need to add one in order to count 5 itself.ð

ð
1. \(100 = 97 + 2 + 1\text{.}\)ð

ð
1. There are 12 primes not greater than 37: \({\left\{2,3,5,7,11,13,17,19,23,29,31,37\right\}}\text{.}\)ð

ð
#### 4.

#### 5.

#### 6.

#### 7.

#### 8.

1. Find \(A \cap B\text{.}\)ð

ð
1. Find \(A \cup B\text{.}\)ð

ð
1. Find \(A \setminus B\text{.}\)ð

ð
1. Find \(A \cap \overline{(B \cup C)}\text{.}\)ð

ð
1. \(\displaystyle A \cap B = {\left\{5,6,7\right\}}\text{.}\)ð

ð
1. \(\displaystyle A \cup B = {\left\{3,4,5,6,7,8,9\right\}}\text{.}\)ð

ð
1. \(\displaystyle A \setminus B ={\left\{3,4\right\}}\text{.}\)ð

ð
1. \(\displaystyle A \cap \overline{(B \cup C)} = {\left\{3\right\}}\text{.}\)ð

ð
#### 9.

1. Find \(A \cap B\text{.}\)ð

ð
1. Find \(A \setminus B\text{.}\)ð

ð
1. \(A \cap B\) will be the set of natural numbers that are both at least 3 and less than 16, and even.  That is, \(\renewcommand{\d}{\displaystyle}\newcommand{\N}{\mathbb N}\newcommand{\st}{:}\renewcommand{\v}{\vtx{above}{}}A \cap B = \{x \in \N \st 3 \le x \lt  16 \wedge x \text{ is even}\} = {\left\{4,6,8,10,12,14\right\}}\text{.}\)ð

ð
1. \(A \setminus B\) is the set of all elements that are in \(A\) but not \(B\text{.}\)  So this is \(\renewcommand{\d}{\displaystyle}\newcommand{\N}{\mathbb N}\newcommand{\st}{:}\renewcommand{\v}{\vtx{above}{}}\{x \in \N \st 3 \le x \lt  16 \wedge x \text{ is odd}\} = {\left\{3,5,7,9,11,13,15\right\}}\text{.}\)ð

Note this is the same set as \(A \cap \overline{B}\text{.}\)ð

ð
#### 10.

1. Find \(A \cap B\text{.}\)ð

ð
1. Find \(A
\cup B\text{.}\)ð

ð
1. Find \(B \cap C\text{.}\)ð

ð
1. Find \(B \cup C\text{.}\)ð

ð
#### 11.

#### 12.

#### 13.

1. Is \(\Z^+ \subseteq 2\Z\text{?}\) Explain.ð

ð
1. Is \(2\Z \subseteq \Z^+\text{?}\) Explain.ð

ð
1. Find \(2\Z \cap 3\Z\text{.}\) Describe the set in words, and using set notation.ð

ð
1. Express \(\{x \in \Z \st \exists y\in \Z (x = 2y \vee x = 3y)\}\) as a union or intersection of two sets already described in this problem.ð

ð
1. No.ð

ð
1. No.ð

ð
1. \(2\Z \cap 3\Z\) is the set of all integers which are multiples of both 2 and 3 (so multiples of 6). Therefore \(2\Z \cap 3\Z
= \{x \in \Z \st \exists y\in \Z(x = 6y)\}\text{.}\)

ð
1. \(2\Z \cup 3\Z\text{.}\)
ð
#### 14.

#### 15.

1. \(\displaystyle A \cup \bar B\)
ð
1. \(\displaystyle \bar{(A \cup B)}\)
ð
1. \(\displaystyle A
\cap (B \cup C)\)
ð
1. \(\displaystyle (A \cap B) \cup C\)
ð
1. \(\displaystyle \bar A \cap B \cap
\bar C\)
ð
1. \(\displaystyle (A \cup B) \setminus C\)
ð
1. \(A
\cup \bar B\text{:}\)ð


ð
1. \(\bar{(A
\cup B)}\text{:}\)ð


ð
1. \(A
\cap (B \cup C)\text{:}\)ð


ð
1. \((A
\cap B) \cup C\text{:}\)ð


ð
1. \(\bar
A \cap B \cap \bar C\text{:}\)ð


ð
1. \((A
\cup B) \setminus C\text{:}\)ð


ð
#### 16.

#### 17.

#### 18.

#### 19.

#### 20.

#### 21.

#### 22.

#### 23.

1. What are the smallest and largest possible values of \(|A \cup
B|\text{?}\) Explain.ð

ð
1. What are the smallest and largest possible values of \(|A \cap B|\text{?}\) Explain.ð

ð
1. What are the smallest and largest possible values of \(|A \times B|\text{?}\) Explain.ð

ð
#### 24.

1. A set \(A \subseteq \N\) with \(|A| =
10\) such that \(X \setminus A = \{10, 12, 14\}\text{.}\)ð

ð
1. A set \(B \in \pow(X)\) with \(|B| = 5\text{.}\)ð

ð
1. A set \(C \subseteq \pow(X)\) with \(|C| = 5\text{.}\)ð

ð
1. A set \(D \subseteq X \times X\) with \(|D| = 5\)ð

ð
1. A set \(E
\subseteq X\) such that \(|E| \in E\text{.}\)ð

ð
#### 25.

1. Suppose that \(A \subseteq B\) and \(B \subseteq C\text{.}\) Does this mean that \(A \subseteq C\text{?}\) Prove your answer. Hint: to prove that \(A \subseteq C\) you must prove the implication, âfor all \(x\text{,}\) if \(x \in A\) then \(x \in C\text{.}\)âð

ð
1. Suppose that \(A \in B\) and \(B \in C\text{.}\) Does this mean that \(A \in C\text{?}\) Give an example to prove that this does NOT always happen (and explain why your example works). You should be able to give an example where \(|A| = |B| = |C| = 2\text{.}\)ð

ð
#### 26.

#### 27.

#### 28.

#### 29.

#### 30.