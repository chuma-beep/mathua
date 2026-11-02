# Exponential Generating Functions

## Labelled Structures

An ordinary generating function encodes a sequence by placing the coefficient of x to the power n in front of x to the power n. This is natural when the objects are unlabelled and order matters only through size.

When the n positions are labelled, the natural encoding divides the coefficient by n factorial. The resulting series is an exponential generating function, and its coefficients count labelled structures.

The factorial denominator is bookkeeping for the labels, not part of the count itself.

## Standard Series

The exponential series has coefficient one over n factorial, so its exponential generating function is e to the power x. Multiplying the exponent by k multiplies the n-th coefficient by k to the power n.

The geometric series one over one minus x has ordinary coefficients all equal to one, so as an exponential generating function its n-th coefficient is n factorial.

Recognising these standard series lets you read off counts without expanding every term.

## Products and Selections

Multiplying exponential generating functions corresponds to labelling a structure by splitting the label set between two independent parts. This is the labelled analogue of the product rule for ordinary generating functions.

The coefficient of x to the power n over n factorial in a product is a sum of binomial choices, which is exactly the number of ways to allocate labels.

Use exponential generating functions when the objects are permutations, functions or other labelled arrangements, and ordinary generating functions when only sizes matter.
