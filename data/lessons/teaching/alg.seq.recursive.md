# Recursively Defined Sequences

## Recursive Rules

A recursively defined sequence gives one or more starting values and a rule that produces each later term from earlier ones. For example, a1 = 3 and a(n+1) = a(n) + 4 defines an arithmetic sequence by repeated addition. The rule is not a formula for the nth term directly; it is a procedure, and generating terms means applying it in order. Writing the first several terms is the fastest way to see the pattern.

## Arithmetic and Geometric Recursion

When each term is the previous term plus a constant, the sequence is arithmetic and the constant is the common difference. When each term is the previous term multiplied by a constant, the sequence is geometric and the constant is the common ratio. Recursion can also combine terms, as in the Fibonacci rule a(n+2) = a(n+1) + a(n), which is neither arithmetic nor geometric. Recognising the type tells you which closed form to look for.

## Closed Forms

A closed form expresses a(n) directly in terms of n. For arithmetic recursion it is a1 + (n - 1)d, and for geometric recursion it is a1 r^(n-1). Substituting n gives a later term without listing every earlier one, which is the main reason to convert. For other recurrences, computing several terms and guessing a pattern is a reasonable start, but the guess should be checked against the rule before it is used.
