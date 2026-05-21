> Content sourced from [Algebrica](https://algebrica.org/arithmetic-sequence/) — CC BY-NC 4.0

## What is an arithmetic sequence

A [sequence](<../sequences>) \\( a_n \\) is called an **arithmetic sequence** (or arithmetic progression) if it consists of numbers arranged in such a way that the difference between any term and the one before it is constant. It is characterized by terms of the form:

\\[a_1, a_2, \dots, a_n \quad \text{with} \quad a_n - a_{n-1} = d\\]

  * By convention, the first term of an arithmetic progression is typically indexed with \\( n = 1.\\)
  * \\( d \\) represents the difference between two consecutive terms in an arithmetic progression, and it is known as the common difference.
  * If \\( d > 0 \\), the progression is increasing.
  * If \\( d < 0 \\), the progression is decreasing.
  * If \\( d = 0 \\), the progression is constant.


Let’s consider, for example, the sequence of non-negative even numbers:

![](/diagrams/algebrica/arithmetic-sequence.png)


An arithmetic sequence can also be defined using a recursive formula:

\\[a_n = a_1 + n \cdot d \quad \text{where } a_1, d \in \mathbb{R} \\]

![](/diagrams/algebrica/arithmetic-sequence-2.png)

An arithmetic progression exhibits a characteristic stepwise pattern, where the height of each step corresponds to the common difference between consecutive terms in the sequence.


In an arithmetic progression, each term \\( a_n \\) is obtained by adding the first term \\( a_1 \\) to the product of the common difference \\( d \\) and \\( (n - 1) \\). This gives the general formula for the \\( n \\)-th term:

\\[a_n = a_1 + (n - 1) \cdot d \quad \text{for } n \geq 1 \\]

This formula allows you to compute any term in the sequence directly, without listing all the previous ones.

## Example

Let’s define an arithmetic sequence with first term \\( a_1 = 2 \\) and common difference \\( d = 3 \\). We use the formula:

\\[a_n = a_1 + (n - 1) \cdot d \\]

Plug in the values:

\\[a_n = 2 + (n - 1) \cdot 3 \\]

Now calculate the first few terms:

  * \\( a_1 = 2 \\)
  * \\( a_2 = 2 + 1 \cdot 3 = 5 \\)
  * \\( a_3 = 2 + 2 \cdot 3 = 8 \\)
  * \\( a_4 = 2 + 3 \cdot 3 = 11 \\)
  * \\( a_5 = 2 + 4 \cdot 3 = 14 \\)


The resulting sequence is:  
\\[2,\ 5,\ 8,\ 11,\ 14,\ \dots \\]

## Sum of \\(n\\) terms of an arithmetic progression

The sum \\( S_n \\) of the first \\( n \\) terms \\( a_1, a_2, \dots, a_n \\) of an arithmetic progression is equal to the product of \\( n \\) and the average of the first and last term:

\\[S_n = n \cdot \frac{a_1 + a_n}{2} \\]

This formula allows you to quickly compute the total sum of a finite number of terms in an arithmetic progression. For example, consider the arithmetic progression of non-negative even numbers:  
\\[2,\ 4,\ 6,\ 8,\ 10 \\]

We want to calculate the sum of the first 5 terms \\( (n = 5).\\) Using the formula, we have:  
\\[S_5 = 5 \cdot \frac{2 + 10}{2} = 5 \cdot 6 = 30 \\]

##### This illustrates the same reasoning behind Gauss’s trick: by pairing the first and last terms, you can quickly compute the total sum of an arithmetic progression.
