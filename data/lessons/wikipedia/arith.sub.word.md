> Content sourced from [Wikipedia](https://en.wikipedia.org/wiki/Subtraction) — CC BY-SA 4.0

# Subtraction word problems

Subtractother usesSubtraction (disambiguation)


**Subtraction** (which is signified by the minus sign, −) is one of the four arithmetic operations along with addition, multiplication and division. Subtraction is an operation that represents removal of objects from a collection. For example, in the adjacent picture, there are 5 − 2 peaches—meaning 5 peaches with 2 taken away, resulting in a total of 3 peaches. Therefore, the *difference* of 5 and 2 is 3; that is, 5 − 2 = 3. While primarily associated with natural numbers in arithmetic, subtraction can also represent removing or decreasing physical and abstract quantities using different kinds of objects including negative numbers, fractions, irrational numbers, vectors, decimals, functions, and matrices.

In a sense, subtraction is the inverse of addition. That is, *c* = *a* − *b* if and only if *c* + *b* = *a*. In words: the difference of two numbers is the number that gives the first one when added to the second one.

Subtraction follows several important patterns. Subtraction is anti-commutative, meaning that changing the order changes the sign of the answer. It is also not associative, meaning that when one subtracts more than two numbers, the order in which subtraction is performed matters. Because 0 is the additive identity, subtraction of it does not change a number. Subtraction also obeys predictable rules concerning related operations, such as addition and multiplication. All of these rules can be proven, starting with the subtraction of integers and generalizing up through the real numbers and beyond. General binary operations that follow these patterns are studied in abstract algebra.

In computability theory, considering subtraction is not well-defined over natural numbers, operations between numbers are actually defined using "truncated subtraction" or monus.

## Notation and terminology
Subtraction is usually written using the minus sign "−" between the terms; that is, in infix notation. The result is expressed with an equals sign. For example, \(2 - 1 = 1\) (pronounced as "two minus one equals one") and \(4 - 6 = -2\) (pronounced as "four minus six equals negative two"). Nonetheless, there are some situations where subtraction is understood, even though no symbol appears; in accounting, a column of two numbers, with the lower number in red, usually indicates that the lower number in the column is to be subtracted, with the difference written below, under a line.

The number being subtracted is the **subtrahend**, while the number it is subtracted from is the **minuend**. The result is the **difference**, that is:

\[
{\rm minuend} - {\rm subtrahend} = {\rm difference}.
\]


All of this terminology derives from Latin. "Subtraction" is an English word derived from the Latin verb *subtrahere*, which in turn is a compound of *sub* "from under" and *trahere* "to pull". Thus, to subtract is to *draw from below*, or to *take away*. Using the gerundive suffix *-nd* results in "subtrahend", "thing to be subtracted". Likewise, from *minuere* "to reduce or diminish", one gets "minuend", which means "thing to be diminished".

## Of integers and real numbers
### Integers

Imagine a line segment of length *b* with the left end labeled *a* and the right end labeled *c*.
Starting from *a*, it takes *b* steps to the right to reach *c*. This movement to the right is modeled mathematically by addition:
*a* + *b* = *c*.

From *c*, it takes *b* steps to the *left* to get back to *a*. This movement to the left is modeled by subtraction:
*c* − *b* = *a*.


Now, a line segment labeled with the numbers 1, 2, and 3. From position 3, it takes no steps to the left to stay at 3, so 3 − 0 = 3. It takes 2 steps to the left to get to position 1, so 3 − 2 = 1. This picture is inadequate to describe what would happen after going 3 steps to the left of position 3. To represent such an operation, the line must be extended.

To subtract arbitrary natural numbers, one begins with a line containing every natural number (0, 1, 2, 3, 4, 5, 6, ...). From 3, it takes 3 steps to the left to get to 0, so 3 − 3 = 0. But 3 − 4 is still invalid, since it again leaves the line. The natural numbers are not a useful context for subtraction.

The solution is to consider the integer number line (..., −3, −2, −1, 0, 1, 2, 3, ...). This way, it takes 4 steps to the left from 3 to get to −1:
3 − 4 = −1.

### Natural numbers
Subtraction of natural numbers is not closed: the difference is not a natural number unless the minuend is greater than or equal to the subtrahend. For example, 26 cannot be subtracted from 11 to give a natural number. Such a case uses one of two approaches:
# Conclude that 26 cannot be subtracted from 11; subtraction becomes a partial function.
# Give the answer as an integer representing a negative number, so the result of subtracting 26 from 11 is −15.

### Real numbers
The field of real numbers can be defined specifying only two binary operations, addition and multiplication, together with unary operations yielding additive and multiplicative inverses. The subtraction of a real number (the subtrahend) from another (the minuend) can then be defined as the addition of the minuend and the additive inverse of the subtrahend. For example, 3 − *π* = 3 + (−*π*). Alternatively, instead of requiring these unary operations, the binary operations of subtraction and division can be taken as basic.

## Properties
### Anti-commutativity
Subtraction is anti-commutative, meaning that if one reverses the terms in a difference left-to-right, the result is the negative of the original result. Symbolically, if *a* and *b* are any two numbers, then
*a* − *b* = −(*b* − *a)*.

### Non-associativity
Subtraction is non-associative, which comes up when one tries to define repeated subtraction. In general, the expression
"*a* − *b* − *c*"
can be defined to mean either (*a* − *b*) − *c* or *a* − (*b* − *c*), but these two possibilities lead to different answers. To resolve this issue, one must establish an order of operations, with different orders yielding different results.

### Predecessor
In the context of integers, subtraction of one also plays a special role: for any integer *a*, the integer (*a* − 1) is the largest integer less than *a*, also known as the predecessor of *a*.

## Units of measurement
When subtracting two numbers with units of measurement such as kilograms or pounds, they must have the same unit. In most cases, the difference will have the same unit as the original numbers.

### Percentages
Changes in percentages can be reported in at least two forms, percentage change and percentage point change. Percentage change represents the relative change between the two quantities as a percentage, while percentage point change is simply the number obtained by subtracting the two percentages.

As an example, suppose that 30% of widgets made in a factory are defective. Six months later, 20% of widgets are defective. The percentage change is 20% − 30%/30% = −1/3 = −33/1%, while the percentage point change is −10 percentage points.

## In computing
The method of complements is a technique used to subtract one number from another using only the addition of positive numbers. This method was commonly used in mechanical calculators, and is still used in modern computers.


! Binarydigit
! Ones' complement
|-align="center"
| 0
| 1
|-align="center"
| 1
| 0

To subtract a binary number *y* (the subtrahend) from another number *x* (the minuend), the ones' complement of *y* is added to *x* and one is added to the sum. The leading digit "1" of the result is then discarded.

The method of complements is especially useful in binary (radix 2) since the ones' complement is very easily obtained by inverting each bit (changing "0" to "1" and vice versa). And adding 1 to get the two's complement can be done by simulating a carry into the least significant bit. For example:

   01100100  (x, equals decimal 100)
 - 00010110  (y, equals decimal 22)

becomes the sum:

   01100100  (x)
 + 11101001  (ones' complement of y)
 +        1  (to get the two's complement)
 ——————————
  101001110

Dropping the initial "1" gives the answer: 01001110 (equals decimal 78)

## The teaching of subtraction in schools
Methods used to teach subtraction to elementary school vary from country to country, and within a country, different methods are adopted at different times. In what is known in the United States as traditional mathematics, a specific process is taught to students at the end of the 1st year (or during the 2nd year) for use with multi-digit whole numbers, and is extended in either the fourth or fifth grade to include decimal representations of fractional numbers.

### In America
Almost all American schools currently teach a method of subtraction using borrowing or regrouping (the decomposition algorithm) and a system of markings called crutches. Although a method of borrowing had been known and published in textbooks previously, the use of crutches in American schools spread after William A. Brownell published a study—claiming that crutches were beneficial to students using this method. This system caught on rapidly, displacing the other methods of subtraction in use in America at that time.

### In Europe
Some European schools employ a method of subtraction called the Austrian method, also known as the additions method. There is no borrowing in this method. There are also crutches (markings to aid memory), which vary by country.

### Comparing the two main methods
Both these methods break up the subtraction as a process of one digit subtractions by place value. Starting with a least significant digit, a subtraction of the subtrahend:
*s*\(_{*j*}\) *s*\(_{*j*−1}\) ... *s*\(_{1}\)
from the minuend
*m*\(_{*k*}\) *m*\(_{*k*−1}\) ... *m*\(_{1}\),
where each *s*\(_{*i*}\) and *m*\(_{*i*}\) is a digit, proceeds by writing down *m*\(_{1}\) − *s*\(_{1}\), *m*\(_{2}\) − *s*\(_{2}\), and so forth, as long as *s*\(_{*i*}\) does not exceed *m*\(_{*i*}\). Otherwise, *m*\(_{*i*}\) is increased by 10 and some other digit is modified to correct for this increase. The American method corrects by attempting to decrease the minuend digit *m*\(_{*i*+1}\) by one (or continuing the borrow leftwards until there is a non-zero digit from which to borrow). The European method corrects by increasing the subtrahend digit *s*\(_{*i*+1}\) by one.

**Example:** 704 − 512.

\begin{array}{rrrr}
        & \color{Red}-1 \\
        & C & D & U  \\
        & 7 & 0 & 4 \\
        & 5 & 1 & 2 \\
      \hline
        & 1 & 9 & 2 \\
   \end{array}
   \begin{array}{l}
       { \color{Red}\longleftarrow \rm carry }\\
       \\
       \longleftarrow \; \rm Minuend\\
       \longleftarrow \; \rm Subtrahend\\
       \longleftarrow \rm{Rest \; or \; Difference}\\
   \end{array}

The minuend is 704, the subtrahend is 512. The minuend digits are *m*\(_{3}\) = 7, *m*\(_{2}\) = 0 and *m*\(_{1}\) = 4. The subtrahend digits are *s*\(_{3}\) = 5, *s*\(_{2}\) = 1 and *s*\(_{1}\) = 2. Beginning at the one's place, 4 is not less than 2 so the difference 2 is written down in the result's one's place. In the ten's place, 0 is less than 1, so the 0 is increased by 10, and the difference with 1, which is 9, is written down in the ten's place. The American method corrects for the increase of ten by reducing the digit in the minuend's hundreds place by one. That is, the 7 is struck through and replaced by a 6. The subtraction then proceeds in the hundreds place, where 6 is not less than 5, so the difference is written down in the result's hundred's place. We are now done, the result is 192.

The Austrian method does not reduce the 7 to 6. Rather it increases the subtrahend hundreds digit by one. A small mark is made near or below this digit (depending on the school). Then the subtraction proceeds by asking what number when increased by 1, and 5 is added to it, makes 7. The answer is 1, and is written down in the result's hundreds place.

There is an additional subtlety in that the student always employs a mental subtraction table in the American method. The Austrian method often encourages the student to mentally use the addition table in reverse. In the example above, rather than adding 1 to 5, getting 6, and subtracting that from 7, the student is asked to consider what number, when increased by 1, and 5 is added to it, makes 7.

## Subtraction by hand
### Austrian method
Example:


### Subtraction from left to right
Example:


### American method
In this method, each digit of the subtrahend is subtracted from the digit above it starting from right to left. If the top number is too small to subtract the bottom number from it, we add 10 to it; this 10 is "borrowed" from the top digit to the left, which we subtract 1 from. Then we move on to subtracting the next digit and borrowing as needed, until every digit has been subtracted. Example:


### Trade first
A variant of the American method where all borrowing is done before all subtraction.

Example:


### Partial differences
The partial differences method is different from other vertical subtraction methods because no borrowing or carrying takes place. In their place, one places plus or minus signs depending on whether the minuend is greater or smaller than the subtrahend. The sum of the partial differences is the total difference.

Example:


### Nonvertical methods
#### Counting up
Instead of finding the difference digit by digit, one can count up the numbers between the subtrahend and the minuend.

Example:
1234 − 567 = can be found by the following steps:
* 567 + **3** = 570
* 570 + **30** = 600
* 600 + **400** = 1000
* 1000 + **234** = 1234
Add up the value from each step to get the total difference: 3 + 30 + 400 + 234 = 667.

#### Breaking up the subtraction
Another method that is useful for mental arithmetic is to split up the subtraction into small steps.

Example:
1234 − 567 = can be solved in the following way:
* 1234 − **500** = 734
* 734 − **60** = 674
* 674 − **7** = 667

#### Same change
The same change method uses the fact that adding or subtracting the same number from the minuend and subtrahend does not change the answer. One simply adds the amount needed to get zeros in the subtrahend.

Example:

"1234 − 567 =" can be solved as follows:
* 1234 − 567 = 1237 − 570 = 1267 − 600 = 667

