# Householder reflections

A Householder reflection mirrors a vector across the hyperplane perpendicular to a chosen direction. It is orthogonal, symmetric and has determinant minus 1.

## Reflection formula

The Householder matrix is the identity minus 2 v v transpose divided by v transpose v. Applying it to x reflects x across the hyperplane perpendicular to v. The formula preserves lengths and angles.

## Determinant

A Householder reflection is an orthogonal reflection in a hyperplane. It reverses orientation, so its determinant is minus 1. The result is independent of the chosen vector v.

## Symmetry

The matrix is symmetric because v v transpose is symmetric. It is also its own inverse, so applying the reflection twice returns the original vector. Symmetry is the quickest check on a proposed Householder matrix.

The three sections are one chain. The first names the object, the second performs the computation and the third reads the result in context. Practice asks one family from each section, so the three moves stay connected.
