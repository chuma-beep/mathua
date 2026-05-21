> Content sourced from [Wikipedia](https://en.wikipedia.org/wiki/Truth_table) — CC BY-SA 4.0

# Truth tables

A **truth table** is a mathematical table used in logic—specifically in connection with Boolean algebra, Boolean functions, and propositional calculus—which sets out the functional values of logical expressions on each of their functional arguments, that is, for each combination of values taken by their logical variables. In particular, truth tables can be used to show whether a propositional expression is true for all legitimate input values, that is, logically valid.

A truth table has one column for each input variable (for example, A and B), and one final column showing the result of the logical operation that the table represents (for example, A XOR B). Each row of the truth table contains one possible configuration of the input variables (for instance, A=true, B=false), and the result of the operation for those values.

A proposition's truth table is a graphical representation of its truth function. The truth function can be more useful for mathematical purposes, although the same information is encoded in both.

Ludwig Wittgenstein is generally credited with inventing and popularizing the truth table in his *Tractatus Logico-Philosophicus*, which was completed in 1918 and published in 1921. Such a system was also independently proposed in 1921 by Emil Leon Post.

## History
Irving Anellis's research shows that C.S. Peirce appears to be the earliest logician (in 1883) to devise a truth table matrix.

From the summary of Anellis's paper:
 In 1997, John Shosky discovered, on the verso of a page of the typed transcript of Bertrand Russell's 1912 lecture on "The Philosophy of Logical Atomism" truth table matrices. The matrix for negation is Russell's, alongside of which is the matrix for material implication in the hand of Ludwig Wittgenstein. It is shown that an unpublished manuscript identified as composed by Peirce in 1893 includes a truth table matrix that is equivalent to the matrix for material implication discovered by John Shosky. An unpublished manuscript by Peirce identified as having been composed in 1883–84 in connection with the composition of Peirce's "On the Algebra of Logic: A Contribution to the Philosophy of Notation" that appeared in the *American Journal of Mathematics* in 1885 includes an example of an indirect truth table for the conditional.

## Applications
Truth tables can be used to prove many other logical equivalences. For example, consider the following truth table:

| | \(p\) | | \(q\) | | \(\neg p\) | | \(\neg p\vee q\) | | \(p\rightarrow q\) |
| --- | --- | --- | --- | --- |
|  |  |  |  |  |
|  |  |  |  |  |
|  |  |  |  |  |
|  |  |  |  |  |

This demonstrates the fact that \(p\rightarrow q\) is logically equivalent to \(\neg p\vee q\).

### Truth table for logic gates
Here is a truth table that gives definitions of each of the 6 possible 2-input logic gate functions of two Boolean variables P and Q:

<table>
  <tr>
    <th>\(P\) || \(Q\) || \(P\and Q\) || \(P\vee Q\) || \(P\uparrow Q\) || \(P\downarrow Q\) || \(P\nleftrightarrow Q\) || \(P\leftrightarrow Q\)</th>
  </tr>
  <tr>
    <td></td>
    <td></td>
    <td></td>
    <td></td>
    <td></td>
    <td></td>
    <td></td>
  </tr>
  <tr>
    <td></td>
    <td></td>
    <td></td>
    <td></td>
    <td></td>
    <td></td>
    <td></td>
  </tr>
  <tr>
    <td></td>
    <td></td>
    <td></td>
    <td></td>
    <td></td>
    <td></td>
    <td></td>
  </tr>
  <tr>
    <td></td>
    <td></td>
    <td></td>
    <td></td>
    <td></td>
    <td></td>
    <td></td>
  </tr>
  <tr>
    <td>|Name(function)</td>
    <td>AND(conjunction)</td>
    <td>OR(disjunction)</td>
    <td>NAND(non-conjunction)</td>
    <td>NOR(non-disjunction)</td>
    <td>XOR(non-equivalence)</td>
    <td>XNOR(equivalence)</td>
  </tr>
  <tr>
    <td>| *where* *means* **true** *and* *means* **false**</td>
  </tr>
</table>

### Condensed truth tables for binary operators
For binary operators, a condensed form of truth table is also used, where the row headings and the column headings specify the operands and the table cells specify the result. For example, Boolean logic uses this condensed truth table notation:

<table>
  <tr>
    <td>| {|</td>
  </tr>
  <tr>
    <th>∧</th>
    <th>T</th>
    <th>F</th>
  </tr>
  <tr>
    <th>T</th>
  </tr>
  <tr>
    <th>F</th>
  </tr>
</table>

| ∧ | T | F |
| --- | --- | --- |
| T |  |  |
| F |  |  |

| style="width:80px;"|
|
| ∨ | T | F |
| --- | --- | --- |
| T |  |  |
| F |  |  |

|}

This notation is useful especially if the operations are commutative, although one can additionally specify that the rows are the first operand and the columns are the second operand. This condensed notation is particularly useful in discussing multi-valued extensions of logic, as it significantly cuts down on combinatoric explosion of the number of rows otherwise needed. It also provides for quickly recognizable characteristic "shape" of the distribution of the values in the table which can assist the reader in grasping the rules more quickly.

### Truth tables in digital logic
Truth tables are also used to specify the function of hardware look-up tables (LUTs) in digital logic circuitry. For an n-input LUT, the truth table will have values (or rows in the above tabular format), completely specifying a Boolean function for the LUT. By representing each Boolean value as a bit in a binary number, truth table values can be efficiently encoded as integer values in electronic design automation (EDA) software. For example, a 32-bit integer can encode the truth table for a LUT with up to 5 inputs.

When using an integer representation of a truth table, the output value of the LUT can be obtained by calculating a bit index *k* based on the input values of the LUT, in which case the LUT's output value is the *k*th bit of the integer. For example, to evaluate the output value of a LUT given an array of *n* Boolean input values, the bit index of the truth table's output value can be computed as follows: if the *i*th input is true, let \(V_i = 1\), else let \(V_i = 0\). Then the *k*th bit of the binary representation of the truth table is the LUT's output value, where

\[
k = V_0 \times 2^0 + V_1 \times 2^1 + V_2 \times 2^2 + \dots + V_{n-1} \times 2^{n-1}.
\]

Truth tables are a simple and straightforward way to encode Boolean functions, however given the exponential growth in size as the number of inputs increase, they are not suitable for functions with a large number of inputs. Other representations which are more memory efficient are text equations and binary decision diagrams.

### Applications of truth tables in digital electronics
In digital electronics and computer science (fields of applied logic engineering and mathematics), truth tables can be used to reduce basic Boolean operations to simple correlations of inputs to outputs, without the use of logic gates or code. For example, a binary addition can be represented with the truth table:

| | | | | | | | |
| --- | --- | --- | --- |

where A is the first operand, B is the second operand, C is the carry digit, and R is the result.

This truth table is read left to right:
* Value pair (A, B) equals value pair (C, R).
* Or for this example, A plus B equal result R, with the Carry C.

This table does not describe the logic operations necessary to implement this operation, rather it simply specifies the function of inputs to output values.

With respect to the result, this example may be arithmetically viewed as modulo 2 binary addition, and as logically equivalent to the exclusive-or (exclusive disjunction) binary logic operation.

In this case it can be used for only very simple inputs and outputs, such as 1s and 0s. However, if the number of types of values one can have on the inputs increases, the size of the truth table will increase.

For instance, in an addition operation, one needs two operands, A and B. Each can have one of two values, zero or one. The number of combinations of these two values is 22, or four. So the result is four possible outputs of C and R. If one were to use base 3, the size would increase to 33, or nine possible outputs.

The first "addition" example above is called a half-adder. A full-adder is when the carry from the previous operation is provided as input to the next adder. Thus, a truth table of eight rows would be needed to describe a full adder's logic:

A B C* | C R
0 0 0 | 0 0
0 1 0 | 0 1
1 0 0 | 0 1
1 1 0 | 1 0
0 0 1 | 0 1
0 1 1 | 1 0
1 0 1 | 1 0
1 1 1 | 1 1

Same as previous, but..
C* = Carry from previous adder

## Methods of writing truth tables
Regarding the *guide columns* to the left of a table, which represent propositional variables, different authors have different recommendations about how to fill them in, although this is of no logical significance.

### Alternating method
Lee Archie, a professor at Lander University, recommends this procedure, which is commonly followed in published truth-tables:

# Write out the number of variables (corresponding to the number of statements) in alphabetical order.
# The number of lines needed is 2\(^{*n*}\) where n is the number of variables. (E. g., with three variables, 2\(^{3}\) = 8).
# Start in the right-hand column and alternate 's and 's until you run out of lines.
# Then move left to the next column and alternate pairs of 's and 's until you run out of lines.
# Then continue to the next left-hand column and double the numbers of 's and 's until completed.

This method results in truth-tables such as the following table for *P* → (*Q* ∨ *R* → (*R* → ¬*P*)), produced by Stephen Cole Kleene:
| \(P\) | \(Q\) | \(R\) | \(P\rightarrow (Q\vee R\rightarrow(R\rightarrow \neg P))\) |
| --- | --- | --- | --- |
|  |  |  |  |
|  |  |  |  |
|  |  |  |  |
|  |  |  |  |
|  |  |  |  |
|  |  |  |  |
|  |  |  |  |
|  |  |  |  |

### Combinatorial method
Colin Howson, on the other hand, believes that "it is a good practical rule" to do the following:to start with all Ts, then all the ways (three) two Ts can be combined with one F, then all the ways (three) one T can be combined with two Fs, and then finish with all Fs. If a compound is built up from n distinct sentence letters, its truth table will have 2\(^{n}\) rows; since there are two ways of assigning T or F to the first letter, and for each of these there will be two ways of assigning T or F to the second, and for each of these there will be two ways of assigning T or F to the third, and so on, giving 2.2.2. …, n times, which is equal to 2\(^{n}\).

This results in truth tables like this table "showing that (*A*→*C*)∧(*B*→*C*) and (*A*∨*B*)→*C* are truth-functionally equivalent", modeled after a table produced by Howson:
| \(A\) | \(B\) | \(C\) | \((A\rightarrow C)\and(B\rightarrow C)\) | \((A\vee B)\rightarrow C\) |
| --- | --- | --- | --- | --- |
|  |  |  |  |  |
|  |  |  |  |  |
|  |  |  |  |  |
|  |  |  |  |  |
|  |  |  |  |  |
|  |  |  |  |  |
|  |  |  |  |  |
|  |  |  |  |  |

## Size of truth tables
If there are *n* input variables then there are 2\(^{*n*}\) possible combinations of their truth values. A given function may produce true or false for each combination so the number of different functions of *n* variables is the double exponential 2\(^{2*n*}\).

<table>
  <tr>
    <th>*n*</th>
    <th>2\(^{*n*}\)</th>
    <th>| 2\(^{2*n*}\)</th>
  </tr>
  <tr>
    <td>0</td>
    <td>1</td>
    <td>| 2</td>
    <td>|</td>
  </tr>
  <tr>
    <td>1</td>
    <td>2</td>
    <td>| 4</td>
    <td>|</td>
  </tr>
  <tr>
    <td>2</td>
    <td>4</td>
    <td>| 16</td>
    <td>|</td>
  </tr>
  <tr>
    <td>3</td>
    <td>8</td>
    <td>| 256</td>
    <td>|</td>
  </tr>
  <tr>
    <td>4</td>
    <td>16</td>
    <td>| 65,536</td>
    <td>|</td>
  </tr>
  <tr>
    <td>5</td>
    <td>32</td>
    <td>| 4,294,967,296</td>
    <td>| ≈ 4.3</td>
  </tr>
  <tr>
    <td>6</td>
    <td>64</td>
    <td>| 18,446,744,073,709,551,616</td>
    <td>| ≈ 1.8</td>
  </tr>
  <tr>
    <td>7</td>
    <td>128</td>
    <td>|</td>
    <td>| ≈ 3.4</td>
  </tr>
  <tr>
    <td>8</td>
    <td>256</td>
    <td>|</td>
    <td>| ≈ 1.2</td>
  </tr>
</table>

Truth tables for functions of three or more variables are rarely given.

## Function Tables
It can be useful to have the output of a truth table expressed as a function of some variable values, instead of just a literal truth or false value. These may be called "function tables" to differentiate them from the more general "truth tables". For example, one value, may be used with an XOR gate to conditionally invert another value,. In other words, when is false, the output and when is true, the output is \(\neg X\). The function table for this would look like:
| \(G\) | \(G\nleftrightarrow X\) |
| --- | --- |
| F | \(X\) |
| T | \(\neg X\) |

Similarly, a 4-to-1 multiplexer with select inputs \(S_0\) and \(S_1\), data inputs , and , and output (as displayed in the image) would have this function table:

| \(S_1\) | \(S_0\) |  |
| --- | --- | --- |
| F | F |  |
| F | T |  |
| T | F |  |
| T | T |  |

## Sentential operator truth tables
### Overview table
Here is an extended truth table giving definitions of all sixteen possible truth functions of two Boolean variables ***p*** and ***q***:

<table>
  <tr>
    <th>\(p\) || \(q\)</th>
    <th>|</th>
    <th>\(\bot\)||\(p\downarrow q\)||\(p\nleftarrow q\)||\(\neg p\)||\(p\nrightarrow q\)||\(\neg q\)||\(p\nleftrightarrow q\)||\(p\uparrow q\)||\(p\and q\)||\(p\leftrightarrow q\)||\(q\)||\(p\rightarrow q\)||\(p\)||\(p\leftarrow q\)||\(p\vee q\)||\(\top\)</th>
  </tr>
  <tr>
    <td></td>
    <td>|</td>
    <td></td>
    <td></td>
    <td></td>
    <td></td>
    <td></td>
    <td></td>
    <td></td>
    <td></td>
    <td></td>
    <td></td>
    <td></td>
    <td></td>
    <td></td>
    <td></td>
    <td></td>
  </tr>
  <tr>
    <td></td>
    <td>|</td>
    <td></td>
    <td></td>
    <td></td>
    <td></td>
    <td></td>
    <td></td>
    <td></td>
    <td></td>
    <td></td>
    <td></td>
    <td></td>
    <td></td>
    <td></td>
    <td></td>
    <td></td>
  </tr>
  <tr>
    <td></td>
    <td>|</td>
    <td></td>
    <td></td>
    <td></td>
    <td></td>
    <td></td>
    <td></td>
    <td></td>
    <td></td>
    <td></td>
    <td></td>
    <td></td>
    <td></td>
    <td></td>
    <td></td>
    <td></td>
  </tr>
  <tr>
    <td></td>
    <td>|</td>
    <td></td>
    <td></td>
    <td></td>
    <td></td>
    <td></td>
    <td></td>
    <td></td>
    <td></td>
    <td></td>
    <td></td>
    <td></td>
    <td></td>
    <td></td>
    <td></td>
    <td></td>
  </tr>
  <tr>
    <td>|</td>
  </tr>
  <tr>
    <td>|</td>
    <td>|</td>
    <td></td>
    <td></td>
    <td></td>
    <td></td>
    <td></td>
    <td></td>
    <td></td>
    <td></td>
    <td></td>
    <td></td>
    <td></td>
    <td></td>
    <td></td>
    <td></td>
    <td></td>
  </tr>
  <tr>
    <td>|</td>
    <td>|</td>
    <td></td>
    <td></td>
    <td></td>
    <td></td>
    <td></td>
    <td></td>
    <td></td>
    <td></td>
    <td></td>
    <td></td>
    <td></td>
    <td></td>
    <td></td>
    <td></td>
    <td></td>
  </tr>
  <tr>
    <td>|</td>
    <td>|</td>
    <td>\(\bot\)</td>
    <td>\(p\downarrow q\)</td>
    <td>\(p\nrightarrow q\)</td>
    <td>\(\neg q\)</td>
    <td>\(p\nleftarrow q\)</td>
    <td>\(\neg p\)</td>
    <td>\(p\nleftrightarrow q\)</td>
    <td>\(p\uparrow q\)</td>
    <td>\(p\and q\)</td>
    <td>\(p\leftrightarrow q\)</td>
    <td>\(p\)</td>
    <td>\(p\leftarrow q\)</td>
    <td>\(q\)</td>
    <td>\(p\rightarrow q\)</td>
    <td>\(p\vee q\)</td>
    <td>\(\top\)</td>
  </tr>
  <tr>
    <td>|</td>
    <td>|</td>
    <td>\(\top\)</td>
    <td>\(p\vee q\)</td>
    <td>\(p\leftarrow q\)</td>
    <td>\(p\)</td>
    <td>\(p\rightarrow q\)</td>
    <td>\(q\)</td>
    <td>\(p\leftrightarrow q\)</td>
    <td>\(p\and q\)</td>
    <td>\(p\uparrow q\)</td>
    <td>\(p\nleftrightarrow q\)</td>
    <td>\(\neg q\)</td>
    <td>\(p\nrightarrow q\)</td>
    <td>\(\neg p\)</td>
    <td>\(p\nleftarrow q\)</td>
    <td>\(p\downarrow q\)</td>
    <td>\(\bot\)</td>
  </tr>
  <tr>
    <td>|</td>
    <td>|</td>
    <td>\(\top\)</td>
    <td>\(p\uparrow q\)</td>
    <td>\(p\rightarrow q\)</td>
    <td>\(\neg p\)</td>
    <td>\(p\leftarrow q\)</td>
    <td>\(\neg q\)</td>
    <td>\(p\leftrightarrow q\)</td>
    <td>\(p\downarrow q\)</td>
    <td>\(p\vee q\)</td>
    <td>\(p\nleftrightarrow q\)</td>
    <td>\(q\)</td>
    <td>\(p\nleftarrow q\)</td>
    <td>\(p\)</td>
    <td>\(p\nrightarrow q\)</td>
    <td>\(p\and q\)</td>
    <td>\(\bot\)</td>
  </tr>
  <tr>
    <td>|</td>
    <td>|</td>
    <td></td>
    <td></td>
    <td>F</td>
    <td></td>
    <td></td>
    <td></td>
    <td>F</td>
    <td></td>
    <td>T</td>
    <td>T</td>
    <td>T, F</td>
    <td>T</td>
    <td></td>
    <td></td>
    <td>F</td>
  </tr>
  <tr>
    <td>|</td>
    <td>|</td>
    <td></td>
    <td></td>
    <td></td>
    <td></td>
    <td>F</td>
    <td></td>
    <td>F</td>
    <td></td>
    <td>T</td>
    <td>T</td>
    <td></td>
    <td></td>
    <td>T, F</td>
    <td>T</td>
    <td>F</td>
  </tr>
</table>

where

T = true.
F = false.
The **Com** row indicates whether an operator, **op**, is commutative – *P* op *Q* = *Q* op *P*.
The **Assoc** row indicates whether an operator, **op**, is associative – (*P* op *Q*) op *R* = *P* op (*Q* op *R*).
The **Adj** row shows the operator **op2** such that *P* op *Q* = *Q* op2 *P*.
The **Neg** row shows the operator **op2** such that *P* op *Q* = ¬(*P* op2 *Q*).
The **Dual** row shows the dual operation obtained by interchanging T with F, and AND with OR.
The **L id** row shows the operator's left identities if it has any values such that *I* op *Q* = *Q*.
The **R id** row shows the operator's right identities if it has any values such that *P* op *I* = *P*.

### Wittgenstein table
In proposition 5.101 of the *Tractatus Logico-Philosophicus*, Wittgenstein listed the table above as follows:
<table>
  <tr>
    <th>scope=col |</th>
    <th>scope=col | Truthvalues</th>
    <th>scope=col |</th>
    <th>scope=col | Operator</th>
    <th>scope=col | Operation name</th>
    <th>scope=col | Tractatus</th>
  </tr>
  <tr>
    <td>0</td>
    <td>(F F F F)(p, q)</td>
    <td>⊥</td>
    <td>false</td>
    <td>**Opq**</td>
    <td>Contradiction</td>
    <td>p and not p; and q and not q</td>
  </tr>
  <tr>
    <td>1</td>
    <td>(F F F T)(p, q)</td>
    <td>NOR</td>
    <td>**p** ↓ **q**</td>
    <td>**Xpq**</td>
    <td>Logical NOR</td>
    <td>neither *p* nor *q*</td>
  </tr>
  <tr>
    <td>2</td>
    <td>(F F T F)(p, q)</td>
    <td>↚</td>
    <td>**p** ↚ **q**</td>
    <td>**Mpq**</td>
    <td>Converse nonimplication</td>
    <td>*q* and not *p*</td>
  </tr>
  <tr>
    <td>3</td>
    <td>(F F T T)(p, q)</td>
    <td>**¬p**, **~p**</td>
    <td>**¬p**</td>
    <td>**Np**, **Fpq**</td>
    <td>Negation</td>
    <td>not *p*</td>
  </tr>
  <tr>
    <td>4</td>
    <td>(F T F F)(p, q)</td>
    <td>↛</td>
    <td>**p** ↛ **q**</td>
    <td>| **Lpq**</td>
    <td>Material nonimplication</td>
    <td>*p* and not *q*</td>
  </tr>
  <tr>
    <td>5</td>
    <td>(F T F T)(p, q)</td>
    <td>**¬q**, **~q**</td>
    <td>**¬q**</td>
    <td>**Nq**, **Gpq**</td>
    <td>Negation</td>
    <td>not *q*</td>
  </tr>
  <tr>
    <td>6</td>
    <td>(F T T F)(p, q)</td>
    <td>XOR</td>
    <td>**p** ⊕ **q**</td>
    <td>**Jpq**</td>
    <td>Exclusive disjunction</td>
    <td>*p* or *q*, but not both</td>
  </tr>
  <tr>
    <td>7</td>
    <td>(F T T T)(p, q)</td>
    <td>NAND</td>
    <td>**p** ↑ **q**</td>
    <td>**Dpq**</td>
    <td>Logical NAND</td>
    <td>not both *p* and *q*</td>
  </tr>
  <tr>
    <td>8</td>
    <td>(T F F F)(p, q)</td>
    <td>AND</td>
    <td>**p** ∧ **q**</td>
    <td>**Kpq**</td>
    <td>Logical conjunction</td>
    <td>*p* and *q*</td>
  </tr>
  <tr>
    <td>9</td>
    <td>(T F F T)(p, q)</td>
    <td>XNOR</td>
    <td>**p** iff **q**</td>
    <td>**Epq**</td>
    <td>Logical biconditional</td>
    <td>if *p* then *q*; and if *q* then *p*</td>
  </tr>
  <tr>
    <td>10</td>
    <td>(T F T F)(p, q)</td>
    <td>**q**</td>
    <td>**q**</td>
    <td>**Hpq**</td>
    <td>Projection function</td>
    <td>*q*</td>
  </tr>
  <tr>
    <td>11</td>
    <td>(T F T T)(p, q)</td>
    <td>**p** → **q**</td>
    <td>if **p** then **q**</td>
    <td>**Cpq**</td>
    <td>Material implication</td>
    <td>if *p* then *q*</td>
  </tr>
  <tr>
    <td>12</td>
    <td>(T T F F)(p, q)</td>
    <td>**p**</td>
    <td>**p**</td>
    <td>**Ipq**</td>
    <td>Projection function</td>
    <td>*p*</td>
  </tr>
  <tr>
    <td>13</td>
    <td>(T T F T)(p, q)</td>
    <td>**p** ← **q**</td>
    <td>if **q** then **p**</td>
    <td>**Bpq**</td>
    <td>Converse implication</td>
    <td>if *q* then *p*</td>
  </tr>
  <tr>
    <td>14</td>
    <td>(T T T F)(p, q)</td>
    <td>OR</td>
    <td>**p** ∨ **q**</td>
    <td>**Apq**</td>
    <td>Logical disjunction</td>
    <td>*p* or *q*</td>
  </tr>
  <tr>
    <td>15</td>
    <td>(T T T T)(p, q)</td>
    <td>⊤</td>
    <td>true</td>
    <td>**Vpq**</td>
    <td>Tautology</td>
    <td>if p then p; and if q then q</td>
  </tr>
</table>

The truth table represented by each row is obtained by appending the sequence given in **Truthvalues**\(_{row}\) to the table
| scope=row | *p* | T | T | F | F |
| --- | --- | --- | --- | --- |
| scope=row | *q* | T | F | T | F |

For example, the table
| scope=row | *p* | T | T | F | F |
| --- | --- | --- | --- | --- |
| scope=row | *q* | T | F | T | F |
| scope=row | *11* | T | F | T | T |

represents the truth table for Material implication. Logical operators can also be visualized using Venn diagrams.

### Nullary operations
There are 2 nullary operations:
*Always true
*Never true, unary *falsum*

#### Logical true
The output value is always true, because this operator has zero operands and therefore no input values
| | *p* | | *T* |
| --- | --- |
| T | T |
| F | T |

#### Logical false
The output value is never true: that is, always false, because this operator has zero operands and therefore no input values
| | *p* | | *F* |
| --- | --- |
| T | F |
| F | F |

### Unary operations
There are 2 unary operations:
*Unary *identity*
*Unary *negation*

#### Logical identity
Logical identity is an operation on one logical value p, for which the output value remains p.

The truth table for the logical identity operator is as follows:

| | *p* | | *p* |
| --- | --- |
| T | T |
| F | F |

#### Logical negation
Logical negation is an operation on one logical value, typically the value of a proposition, that produces a value of *true* if its operand is false and a value of *false* if its operand is true.

The truth table for **NOT p** (also written as **¬p**, **Np**, **Fpq**, or **~p**) is as follows:

| | *p* | | *¬p* |
| --- | --- |
| T | F |
| F | T |

### Binary operations
There are 16 possible truth functions of two binary variables, each operator has its own name.

#### Logical conjunction (AND)
Logical conjunction is an operation on two logical values, typically the values of two propositions, that produces a value of *true* if both of its operands are true.

The truth table for **p AND q** (also written as **p ∧ q**, **Kpq**, **p & q**, or **p** \(\cdot\) **q**) is as follows:

| | *p* | | *q* | | *p* ∧ *q* |
| --- | --- | --- |
| T | T | T |
| T | F | F |
| F | T | F |
| F | F | F |

In ordinary language terms, if both *p* and *q* are true, then the conjunction *p* ∧ *q* is true. For all other assignments of logical values to *p* and to *q* the conjunction *p* ∧ *q* is false.

It can also be said that if *p*, then *p* ∧ *q* is *q*, otherwise *p* ∧ *q* is *p*.

#### Logical disjunction (OR)
Logical disjunction is an operation on two logical values, typically the values of two propositions, that produces a value of *true* if at least one of its operands is true.

The truth table for **p OR q** (also written as **p ∨ q**, **Apq**, **p || q**, or **p + q**) is as follows:

| | *p* | | *q* | | *p* ∨ *q* |
| --- | --- | --- |
| T | T | T |
| T | F | T |
| F | T | T |
| F | F | F |

Stated in English, if *p*, then *p* ∨ *q* is *p*, otherwise *p* ∨ *q* is *q*.

#### Logical implication
Logical implication and the material conditional are both associated with an disabled on two deleteMevalue]]s, typically the values of two propositions, which produces a value of *false* if the first operand is true and the second operand is false, and a value of *true* otherwise.

The truth table associated with the logical implication **p implies q** (symbolized as **p ⇒ q**, or more rarely **Cpq**) is as follows:

| | *p* | | *q* | | *p* ⇒ *q* |
| --- | --- | --- |
| T | T | T |
| T | F | F |
| F | T | T |
| F | F | T |

The truth table associated with the material conditional **if p then q** (symbolized as **p → q**) is as follows:

| | *p* | | *q* | | *p* → *q* |
| --- | --- | --- |
| T | T | T |
| T | F | F |
| F | T | T |
| F | F | T |

**p ⇒ q** and **p → q** are equivalent to **¬p ∨ q**.

#### Logical equality
Logical equality (also known as biconditional or exclusive nor) is an operation on two logical values, typically the values of two propositions, that produces a value of *true* if both operands are false or both operands are true.

The truth table for **p XNOR q** (also written as **p ↔ q**, **Epq**, **p = q**, or **p ≡ q**) is as follows:

| | *p* | | *q* | | *p* ↔ *q* |
| --- | --- | --- |
| T | T | T |
| T | F | F |
| F | T | F |
| F | F | T |

So p EQ q is true if p and q have the same truth value (both true or both false), and false if they have different truth values.

#### Exclusive disjunction
Exclusive disjunction is an operation on two logical values, typically the values of two propositions, that produces a value of *true* if one but not both of its operands is true.

The truth table for **p XOR q** (also written as **Jpq**, or **p ⊕ q**) is as follows:

| | *p* | | *q* | | **p** ⊕ **q** |
| --- | --- | --- |
| T | T | F |
| T | F | T |
| F | T | T |
| F | F | F |

For two propositions, **XOR** can also be written as (p ∧ ¬q) ∨ (¬p ∧ q).

#### Logical NAND
The logical NAND is an operation on two logical values, typically the values of two propositions, that produces a value of *false* if both of its operands are true. In other words, it produces a value of *true* if at least one of its operands is false.

The truth table for **p NAND q** (also written as **p ↑ q**, **Dpq**, or **p | q**) is as follows:

| | *p* | | *q* | | *p* ↑ *q* |
| --- | --- | --- |
| T | T | F |
| T | F | T |
| F | T | T |
| F | F | T |

It is frequently useful to express a logical operation as a compound operation, that is, as an operation that is built up or composed from other operations. Many such compositions are possible, depending on the operations that are taken as basic or "primitive" and the operations that are taken as composite or "derivative".

In the case of logical NAND, it is clearly expressible as a compound of NOT and AND.

The negation of a conjunction: ¬(*p* ∧ *q*), and the disjunction of negations: (¬*p*) ∨ (¬*q*) can be tabulated as follows:

| | *p* | | *q* | | *p* ∧ *q* | | ¬(*p* ∧ *q*) | | ¬*p* | | ¬*q* | | (¬*p*) ∨ (¬*q*) |
| --- | --- | --- | --- | --- | --- | --- |
| T | T | T | F | F | F | F |
| T | F | F | T | F | T | T |
| F | T | F | T | T | F | T |
| F | F | F | T | T | T | T |

#### Logical NOR
The logical NOR is an operation on two logical values, typically the values of two propositions, that produces a value of *true* if both of its operands are false. In other words, it produces a value of *false* if at least one of its operands is true. ↓ is also known as the Peirce arrow after its inventor, Charles Sanders Peirce, and is a Sole sufficient operator.

The truth table for **p NOR q** (also written as **p ↓ q**, or **Xpq**) is as follows:

| | *p* | | *q* | | *p* ↓ *q* |
| --- | --- | --- |
| T | T | F |
| T | F | F |
| F | T | F |
| F | F | T |

The negation of a disjunction ¬(*p* ∨ *q*), and the conjunction of negations (¬*p*) ∧ (¬*q*) can be tabulated as follows:

| | *p* | | *q* | | *p* ∨ *q* | | ¬(*p* ∨ *q*) | | ¬*p* | | ¬*q* | | (¬*p*) ∧ (¬*q*) |
| --- | --- | --- | --- | --- | --- | --- |
| T | T | T | F | F | F | F |
| T | F | T | F | F | T | F |
| F | T | T | F | T | F | F |
| F | F | F | T | T | T | T |

Inspection of the tabular derivations for NAND and NOR, under each assignment of logical values to the functional arguments *p* and *q*, produces the identical patterns of functional values for ¬(*p* ∧ *q*) as for (¬*p*) ∨ (¬*q*), and for ¬(*p* ∨ *q*) as for (¬*p*) ∧ (¬*q*). Thus the first and second expressions in each pair are logically equivalent, and may be substituted for each other in all contexts that pertain solely to their logical values.

This equivalence is one of De Morgan's laws.
