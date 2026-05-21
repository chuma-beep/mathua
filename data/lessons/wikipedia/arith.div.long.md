> Content sourced from [Wikipedia](https://en.wikipedia.org/wiki/Long_division) — CC BY-SA 4.0

# Long division

In arithmetic, **long division** is a standard division algorithm suitable for dividing multi-digit Hindu-Arabic numerals (positional notation) that is simple enough to perform by hand. It breaks down a division problem into a series of easier steps.

As in all division problems, one number, called the dividend, is divided by another, called the divisor, producing a result called the quotient. It enables computations involving arbitrarily large numbers to be performed by following a series of simple steps. The abbreviated form of long division is called short division, which is almost always used instead of long division when the divisor has only one digit.

## History
Related algorithms have existed since the 12th century.
Al-Samawal al-Maghribi (1125–1174) performed calculations with decimal numbers that essentially require long division, leading to infinite decimal results, but without formalizing the algorithm.
Caldrini (1491) is the earliest printed example of long division, known as the *Danda* method in medieval Italy, and it became more practical with the introduction of decimal notation for fractions by Pitiscus (1608).
The specific algorithm in modern use was introduced by Henry Briggs 1600.

## Education
Inexpensive calculators and computers have become the most common tools for performing division in educational and professional contexts worldwide, reducing reliance on traditional paper-and-pencil techniques. Internally, these devices implement various division algorithms, many of which rely on iterative approximations and multiplication to improve computational efficiency.

Educational approaches to teaching division vary across countries and regions, reflecting differing curricular priorities. In North America, long division has been de-emphasized or, in some cases, removed from portions of the curriculum as part of reform mathematics, which emphasizes conceptual understanding and the use of technology.

In contrast, many education systems in Europe and Asia continue to emphasize mastery of standard algorithms, including long division, as a foundational arithmetic skill. For example, curricula in countries such as Japan and Germany typically introduce and reinforce long division during primary education, often alongside mental arithmetic strategies and problem-solving techniques.

International assessments such as the Trends in International Mathematics and Science Study (TIMSS) highlight these differences, showing variation in how procedural fluency and conceptual understanding are balanced across educational systems.

These differing approaches reflect broader educational philosophies regarding the balance between procedural fluency, conceptual understanding, and the role of technology in mathematics education.

## Method
In English-speaking countries, long division does not use the division slash or division sign symbols but instead constructs a **tableau**. The divisor is separated from the dividend by a right parenthesis or vertical bar ]]; the dividend is separated from the quotient by a vinculum (i.e., an overbar). The combination of these two symbols is sometimes known as a **long division symbol,** **division bracket**, or even a **bus stop**. It developed in the 18th century from an earlier single-line notation separating the dividend from the quotient by a left parenthesis.

The process is begun by dividing the left-most digit of the dividend by the divisor. The quotient (rounded down to an integer) becomes the first digit of the result, and the remainder is calculated (this step is notated as a subtraction). This remainder carries forward when the process is repeated on the following digit of the dividend (notated as 'bringing down' the next digit to the remainder). When all digits have been processed and no remainder is left, the process is complete.

An example is shown below, representing the division of 500 by 4 (with a result of 125).
 125 (Explanations)
 4)500
 4 ( 4 × 1 = 4)
 10 ( 5 - 4 = 1)
 8 ( 4 × 2 = 8)
 20 (10 - 8 = 2)
 20 ( 4 × 5 = 20)
 0 (20 - 20 = 0)

A more detailed breakdown of the steps goes as follows:

# Find the shortest sequence of digits starting from the left end of the dividend, 500, that the divisor 4 goes into at least once. In this case, this is simply the first digit, 5. The largest number that the divisor 4 can be multiplied by without exceeding 5 is 1, so the digit 1 is put above the 5 to start constructing the quotient.
# Next, the 1 is multiplied by the divisor 4, to obtain the largest whole number that is a multiple of the divisor 4 without exceeding the 5 (4 in this case). This 4 is then placed under and subtracted from the 5 to get the remainder, 1, which is placed under the 4 under the 5.
# Afterwards, the first as-yet unused digit in the dividend, in this case the first digit 0 after the 5, is copied directly underneath itself and next to the remainder 1, to form the number 10.
# At this point the process is repeated enough times to reach a stopping point: The largest number by which the divisor 4 can be multiplied without exceeding 10 is 2, so 2 is written above as the second leftmost quotient digit. This 2 is then multiplied by the divisor 4 to get 8, which is the largest multiple of 4 that does not exceed 10; so 8 is written below 10, and the subtraction 10 minus 8 is performed to get the remainder 2, which is placed below the 8.
# The next digit of the dividend (the last 0 in 500) is copied directly below itself and next to the remainder 2 to form 20. Then the largest number by which the divisor 4 can be multiplied without exceeding 20, which is 5, is placed above as the third leftmost quotient digit. This 5 is multiplied by the divisor 4 to get 20, which is written below and subtracted from the existing 20 to yield the remainder 0, which is then written below the second 20.
# At this point; since there are no more digits to bring down from the dividend and the last subtraction result was 0, we can be assured that the process finished.

If the last remainder when we ran out of dividend digits had been something other than 0, there would have been two possible courses of action:

# We could just stop there and say that the dividend divided by the divisor is the quotient written at the top with the remainder written at the bottom, and write the answer as the quotient followed by a fraction that is the remainder divided by the divisor.
# We could extend the dividend by writing it as, say, 500.000... and continue the process (using a decimal point in the quotient directly above the decimal point in the dividend), in order to get a decimal answer, as in the following example.

 31.75
 4)127.00
 12 (12 ÷ 4 = 3)
 07 (0 remainder, bring down next figure)
 4 (7 ÷ 4 = 1 r 3)
 3.0 (bring down 0 and the decimal point)
 2.8 (7 × 4 = 28, 30 ÷ 4 = 7 r 2)
 20 (an additional zero is brought down)
 20 (5 × 4 = 20)
 0

In this example, the decimal part of the result is calculated by continuing the process beyond the units digit, "bringing down" zeros as being the decimal part of the dividend.

This example also illustrates that, at the beginning of the process, a step that produces a zero can be omitted. Since the first digit 1 is less than the divisor 4, the first step is instead performed on the first two digits 12. Similarly, if the divisor were 13, one would perform the first step on 127 rather than 12 or 1.

### Basic procedure for long division of
# Find the location of all decimal points in the dividend and divisor. # If necessary, simplify the long division problem by moving the decimals of the divisor and dividend by the same number of decimal places, to the right (or to the left), so that the decimal of the divisor is to the right of the last digit.
# When doing long division, keep the numbers lined up straight from top to bottom under the tableau.
# After each step, be sure the remainder for that step is less than the divisor. If it is not, there are three possible problems: the multiplication is wrong, the subtraction is wrong, or a greater quotient is needed.
# In the end, the remainder, is added to the growing quotient as a fraction,. ### Invariant property and correctness
The basic presentation of the steps of the process (above) **focuses on what** steps are to be performed,
rather than the *properties of those steps* that ensure the result will be correct
(specifically, that *q × m + r = n*, where *q* is the final quotient and *r* the final remainder).
A slight variation of presentation requires more writing,
and requires that we change, rather than just update, digits of the quotient,
but can shed more light on *why* these steps actually produce the right answer
by allowing evaluation of *q × m + r* at intermediate points in the process.
This illustrates the key property used in the derivation of the algorithm
(below).

Specifically, we amend the above basic procedure so that
we fill the space after the digits of the *quotient* under construction with 0's, to at least the 1's place,
and include those 0's in the numbers we write below the division bracket.

This lets us maintain an invariant relation at every step:
*q × m + r = n*, where *q* is the partially-constructed quotient (above the division bracket)
and *r* the partially-constructed remainder (bottom number below the division bracket).
Note that, initially *q=0* and *r=n*, so this property holds initially;
the process reduces r and increases q with each step,
eventually stopping when *r 125 (*q'', changes from 000 to 100 to 120 to 125 as per notes below)
 4)500
 400 ( 4 × 100 = 400)
 100 (500 - 400 = 100; now *q=100*, *r=100*; note *q×4+r = 500*.)
 80 ( 4 × 20 = 80)
 20 (100 - 80 = 20; now *q=120*, *r= 20*; note *q×4+r = 500*.)
 20 ( 4 × 5 = 20)
 0 ( 20 - 20 = 0; now *q=125*, *r= 0*; note *q×4+r = 500*.)

### Example with multi-digit divisor

A divisor of any number of digits can be used. In this example, 1260257 is to be divided by 37. First the problem is set up as follows:

 37)1260257

Digits of the number 1260257 are taken until a number greater than or equal to 37 occurs. So 1 and 12 are less than 37, but 126 is greater. Next, the greatest multiple of 37 less than or equal to 126 is computed. So 3 × 37 = 111 < 126, but 4 × 37 > 126. The multiple 111 is written underneath the 126 and the 3 is written on the top where the solution will appear:

 3
 37)1260257
 111

Note carefully which place-value column these digits are written into. The 3 in the quotient goes in the same column (ten-thousands place) as the 6 in the dividend 1260257, which is the same column as the last digit of 111.

The 111 is then subtracted from the line above, ignoring all digits to the right:

 3
 37)1260257
 111
 15

Now the digit from the next smaller place value of the dividend is copied down and appended to the result 15:

 3
 37)1260257
 111
 150

The process repeats: the greatest multiple of 37 less than or equal to 150 is subtracted. This is 148 = 4 × 37, so a 4 is added to the top as the next quotient digit. Then the result of the subtraction is extended by another digit taken from the dividend:

 34
 37)1260257
 111
 150
 148
 22

The greatest multiple of 37 less than or equal to 22 is 0 × 37 = 0. Subtracting 0 from 22 gives 22, we often don't write the subtraction step. Instead, we simply take another digit from the dividend:

 340
 37)1260257
 111
 150
 148
 225

The process is repeated until 37 divides the last line exactly:

 34061
 37)1260257
 111
 150
 148
 225
 222
 37

### Mixed mode long division
For non-decimal currencies (such as the British £sd system before 1971) and measures (such as avoirdupois) **mixed mode** division must be used. Consider dividing 50 miles 600 yards into 37 pieces:

 mi - yd - ft - in
 1 - 634 1 9 r. 15"
 37) 50 - 600 - 0 - 0
 37 22880 66 348
 13 23480 66 348
 1760 222 37 333
 22880 128 29 15
 ===== 111 348 ==
 170 ===
 148
 22
 66
 ==

Each of the four columns is worked in turn. Starting with the miles: 50/37 = 1 remainder 13. No further division is
possible, so perform a long multiplication by 1,760 to convert miles to yards, the result is 22,880 yards. Carry this to the top of the yards column and add it to the 600 yards in the dividend giving 23,480. Long division of 23,480 / 37 now proceeds as normal yielding 634 with remainder 22. The remainder is multiplied by 3 to get feet and carried up to the feet column. Long division of the feet gives 1 remainder 29 which is then multiplied by twelve to get 348 inches. Long division continues with the final remainder of 15 inches being shown on the result line.

### Interpretation of decimal results
When the quotient is not an integer and the division process is extended beyond the decimal point, one of two things can happen:

# The process can terminate, which means that a remainder of 0 is reached; or
# A remainder could be reached that is identical to a previous remainder that occurred after the decimal points were written. In the latter case, continuing the process would be pointless, because from that point onward the same sequence of digits would appear in the quotient over and over. So a bar is drawn over the repeating sequence to indicate that it repeats forever (i.e., every rational number is either a terminating or repeating decimal).

## Notation in non-English-speaking countries

China, Japan, Korea use the same notation as English-speaking nations including India. Elsewhere, the same general principles are used, but the figures are often arranged differently.

### Latin America
In Latin America (except Argentina, Bolivia, Mexico, Colombia, Paraguay, Venezuela, Uruguay and Brazil), the calculation is almost exactly the same, but is written down differently as shown below with the same two examples used above. Usually the quotient is written under a bar drawn under the divisor. A long vertical line is sometimes drawn to the right of the calculations.

 500 ÷ 4 = 125 (Explanations)
 4 ( 4 × 1 = 4)
 10 ( 5 - 4 = 1)
 8 ( 4 × 2 = 8)
 20 (10 - 8 = 2)
 20 ( 4 × 5 = 20)
 0 (20 - 20 = 0)

and

 127 ÷ 4 = 31.75
 124
 30 (bring down 0; decimal to quotient)
 28 (7 × 4 = 28)
 20 (an additional zero is added)
 20 (5 × 4 = 20)
 0
In Mexico, the English-speaking world notation is used, except that only the result of the subtraction is annotated and the calculation is done mentally, as shown below:

 125 (Explanations)
 4)500
 10 ( 5 - 4 = 1)
 20 (10 - 8 = 2)
 0 (20 - 20 = 0)

In Bolivia, Brazil, Paraguay, Venezuela, French-speaking Canada, Colombia, and Peru, the European notation (see below) is used, except that the quotient is not separated by a vertical line, as shown below:

 127|4
 −124 31,75
 30
 −28
 20
 −20
 0

Same procedure applies in Mexico, Uruguay and Argentina, only the result of the subtraction is annotated and the calculation is done mentally.

### Eurasia
In Spain, Italy, France, Portugal, Lithuania, Romania, Turkey, Greece, Belgium, Belarus, Ukraine, and Russia, the divisor is to the right of the dividend, and separated by a vertical bar. The division also occurs in the column, but the quotient (result) is written below the divider, and separated by the horizontal line. The same method is used in Iran, Vietnam, and Mongolia.

 127|4
 −12 |31,75
 7
 −4
 30
 −28
 20
 -20
 0

In Cyprus, as well as in France, a long vertical bar separates the dividend and subsequent subtractions from the quotient and divisor, as in the example below of 6359 divided by 17, which is 374 with a remainder of 1.

 6359|17
 −51 |374
 125 |
 −119 |
 69|
 −68|
 1|

Decimal numbers are not divided directly, the dividend and divisor are multiplied by a power of ten so that the division involves two whole numbers. Therefore, if one were dividing 12,7 by 0,4 (commas being used instead of decimal points), the dividend and divisor would first be changed to 127 and 4, and then the division would proceed as above.

In Austria, Germany and Switzerland, the notational form of a normal equation is used. <dividend> : <divisor> = <quotient>, with the colon ":" denoting a binary infix symbol for the division operator (analogous to "/" or "÷"). In these regions the decimal separator is written as a comma. (cf. first section of Latin American countries above, where it's done virtually the same way):

 127 : 4 = 31,75
 −12
 07
 −4
 30
 −28
 20
 −20
 0

The same notation is adopted in Denmark, Norway, Bulgaria, North Macedonia, Poland, Croatia, Slovenia, Hungary, Czech Republic, Slovakia, Vietnam and in Serbia.

In the Netherlands, the following notation is used:

 12 / 135 \ 11,25
 12
 15
 12
 30
 24
 60
 60
 0

In Finland, the Italian method detailed above was replaced by the Anglo-American one in the 1970s. In the early 2000s, however, some textbooks have adopted the German method as it retains the order between the divisor and the dividend.

## Algorithm for arbitrary base

Every natural number \(n\) can be uniquely represented in an arbitrary number base \(b>1\) as a sequence of digits \(n=\alpha_{0}\alpha_{1}\alpha_{2}...\alpha_{k-1}\) where \(0\leq\alpha_{i}1\), \(r_{i-1} \geq 0\), \(\alpha_{i+l-1} \geq 0\), this is always true. For the right side of the inequality we assume there exists a smallest \(\beta_{i}^{\prime}\) such that
\(br_{i-1} + \alpha_{i+l-1} < m(\beta_{i}^{\prime} + 1)\)
Since this is the smallest \(\beta_{i}^{\prime}\) that the inequality holds true, this must mean that for \(\beta_{i}^{\prime}-1\)
\(br_{i-1} + \alpha_{i+l-1} \geq m\beta_{i}^{\prime}\)
which is exactly the same as the left side of the inequality. Thus, \(\beta_{i}=\beta_{i}^{\prime}\). As \(\beta_{i}\) will always exist, so will \(\beta_{i}^{\prime}\) equal to \(\beta_{i}\), and there is only one unique \(\beta_{i}\) that is valid for the inequality. Thus we have proven the existence and uniqueness of \(\beta_{i}\).

The final quotient is \(q = q_{k-l}\) and the final remainder is \(r = r_{k-l}\)

### Examples
In base 10, using the example above with \(n = 1260257\) and \(m = 37\), the initial values \(q_{-1} = 0\) and \(r_{-1} = 1\).

{| class="wikitable"
|-
! \(0 \leq i \leq k - l\)
! \(\alpha_{i+l-1}\)
! \(d_{i} = b r_{i-1} + \alpha_{i+l-1}\)
! \(\beta_{i}\)
! \(r_{i} = d_{i} - m \beta_{i}\)
! \(q_{i} = b q_{i-1} + \beta_{i}\)
|-
| 0 || 2 || \(10 \cdot 1 + 2 = 12\) || 0 || \(12 - 37 \cdot 0 = 12\) || \(10 \cdot 0 + 0 = 0\)
|-
| 1 || 6 || \(10 \cdot 12 + 6 = 126\) || 3 || \(126 - 37 \cdot 3 = 15\) || \(10 \cdot 0 + 3 = 3\)
|-
| 2 || 0 || \(10 \cdot 15 + 0 = 150\) || 4 || \(150 - 37 \cdot 4 = 2\) || \(10 \cdot 3 + 4 = 34\)
|-
| 3 || 2 || \(10 \cdot 2 + 2 = 22\) || 0 || \(22 - 37 \cdot 0 = 22\) || \(10 \cdot 34 + 0 = 340\)
|-
| 4 || 5 || \(10 \cdot 22 + 5 = 225\) || 6 || \(225 - 37 \cdot 6 = 3\) || \(10 \cdot 340 + 6 = 3406\)
|-
| 5 || 7 || \(10 \cdot 3 + 7 = 37\) || 1 || \(37 - 37 \cdot 1 = 0\) || \(10 \cdot 3406 + 1 = 34061\)
|}

Thus, \(q = 34061\) and \(r = 0\).

In base 16, with \(n = \text{f412df}\) and \(m = 12\), the initial values are \(q_{-1} = 0\) and \(r_{-1} = \text{f}\).

{| class="wikitable"
|-
! \(0 \leq i \leq k - l\)
! \(\alpha_{i+l-1}\)
! \(d_{i} = b r_{i-1} + \alpha_{i+l-1}\)
! \(\beta_{i}\)
! \(r_{i} = d_{i} - m \beta_{i}\)
! \(q_{i} = b q_{i-1} + \beta_{i}\)
|-
| 0 || 4 || \(10 \cdot \text{f} + 4 = \text{f4}\) || \(\text{d}\) || \(\text{f4} - 12 \cdot \text{d} = \text{a}\) || \(10 \cdot 0 + \text{d} = \text{d}\)
|-
| 1 || 1 || \(10 \cdot \text{a} + 1 = \text{a1}\) || 8 || \(\text{a1} - 12 \cdot 8 = 11\) || \(10 \cdot \text{d} + 8 = \text{d8}\)
|-
| 2 || 2 || \(10 \cdot 11 + 2 = 112\) || \(\text{f}\) || \(112 - 12 \cdot \text{f} = 4\) || \(10 \cdot \text{d8} + \text{f} = \text{d8f}\)
|-
| 3 || \(\text{d} = 13\) || \(10 \cdot 4 + \text{d} = \text{4d}\) || 4 || \(\text{4d} - 12 \cdot 4 = 5\) || \(10 \cdot \text{d8f} + 4 = \text{d8f4}\)
|-
| 4 || \(\text{f} = 15\) || \(10 \cdot 5 + \text{f} = \text{5f}\) || 5 || \(\text{5f} - 12 \cdot 5 = 5\) || \(10 \cdot \text{d8f4} + 5 = \text{d8f45}\)
|}

Thus, \(q = \text{d8f45}\) and \(r = \text{5}\).

If one doesn't have the addition, subtraction, or multiplication tables for base memorised, then this algorithm still works if the numbers are converted to decimal and at the end are converted back to base. For example, with the above example,
\(n = \text{f412df}_{16} = 15 \cdot 16^5 + 4 \cdot 16^4 + 1 \cdot 16^3 + 2 \cdot 16^2 + 13 \cdot 16^1 + 15 \cdot 16^0\)
and
\(m = \text{12}_{16} = 1 \cdot 16^1 + 2 \cdot 16^0 = 18\)
with \(b = 16\). The initial values are \(q_{-1} = 0\) and \(r_{-1} = 15\).

{| class="wikitable"
|-
! \(0 \leq i \leq k - l\)
! \(\alpha_{i+l-1}\)
! \(d_{i} = b r_{i-1} + \alpha_{i+l-1}\)
! \(\beta_{i}\)
! \(r_{i} = d_{i} - m \beta_{i}\)
! \(q_{i} = b q_{i-1} + \beta_{i}\)
|-
| 0 || 4 || \(16 \cdot 15 + 4 = 244\) || \(13 = \text{d}\) || \(244 - 18 \cdot 13 = 10\) || \(16 \cdot 0 + 13 = 13\)
|-
| 1 || 1 || \(16 \cdot 10 + 1 = 161\) || 8 || \(161 - 18 \cdot 8 = 17\) || \(16 \cdot 13 + 8\)
|-
| 2 || 2 || \(16 \cdot 17 + 2 = 274\) || \(15 = \text{f}\) || \(274 - 18 \cdot 15 = 4\) || \(16 \cdot (16 \cdot 13 + 8) + 15 = 16^2 \cdot 13 + 16 \cdot 8 + 15\)
|-
| 3 || \(\text{d} = 13\) || \(16 \cdot 4 + 13 = 77\) || 4 || \(77 - 18 \cdot 4 = 5\) || \(16 \cdot (16^2 \cdot 13 + 16 \cdot 8 + 15) + 4 = 16^3 \cdot 13 + 16^2 \cdot 8 + 16 \cdot 15 + 4\)
|-
| 4 || \(\text{f} = 15\) || \(16 \cdot 5 + 15 = 95\) || 5 || \(95 - 18 \cdot 5 = 5\) || \(16 \cdot (16^3 \cdot 13 + 16^2 \cdot 8 + 16 \cdot 15 + 4 = 16^4 \cdot 13 + 16^3 \cdot 8 + 16^2 \cdot 15 + 16^1 \cdot 4 + 5\)
|}
Thus, \(q = 16^4 \cdot 13 + 16^3 \cdot 8 + 16^2 \cdot 15 + 16^1 \cdot 4 + 5 = \text{d8f45}_{16}\) and \(r = 5 = \text{5}_{16}\).

This algorithm can be done using the same kind of pencil-and-paper notations as shown in above sections.

 d8f45 r. 5
 12 ) f412df
 ea
 a1
 90
 112
 10e
 4d
 48
 5f
 5a
 5

### Rational quotients
If the quotient is not constrained to be an integer, then the algorithm does not terminate for \(i>k-l\). Instead, if \(i>k-l\) then \(\alpha_{i}=0\) by definition. If the remainder \(r_{i}\) is equal to zero at any iteration, then the quotient is a \(b\)-adic fraction, and is represented as a finite decimal expansion in base \(b\) positional notation. Otherwise, it is still a rational number but not a \(b\)-adic rational, and is instead represented as an infinite repeating decimal expansion in base \(b\) positional notation.

### Binary division

### Performance
On each iteration, the most time-consuming task is to select \(\beta_{i}\). We know that there are \(b\) possible values, so we can find \(\beta_{i}\) using \(O(\log(b))\) comparisons. Each comparison will require evaluating \(d_{i}-m\beta_{i}\). Let \(k\) be the number of digits in the dividend \(n\) and \(l\) be the number of digits in the divisor \(m\). The number of digits in \(d_{i} \leq l + 1\). The multiplication of \(m\beta_{i}\) is therefore \(O(l)\), and likewise the subtraction of \(d_{i}-m\beta_{i}\). Thus it takes \(O(l\log(b))\) to select \(\beta_{i}\). The remainder of the algorithm are addition and the digit-shifting of \(q_{i}\) and \(r_{i}\) to the left one digit, and so takes time \(O(k)\) and \(O(l)\) in base \(b\), so each iteration takes \(O(l\log(b) + k + l)\), or just \(O(l\log(b) + k)\). For all \(k - l + 1\) digits, the algorithm takes time \(O((k - l + 1)(l\log(b) + k))\), or \(O(kl\log(b) + k^2)\) in base \(b\).

## Generalizations
### Rational numbers
Long division of integers can easily be extended to include non-integer dividends, as long as they are rational. This is because every rational number has a repeating decimal expansion. The procedure can also be extended to include divisors which have a finite or terminating decimal expansion (i.e. decimal fractions). In this case the procedure involves multiplying the divisor and dividend by the appropriate power of ten so that the new divisor is an integer – taking advantage of the fact that *a* ÷ *b* = (*ca*) ÷ (*cb*) – and then proceeding as above.

### Polynomials
A generalised version of this method called polynomial long division is also used for dividing polynomials (sometimes using a shorthand version called synthetic division).
