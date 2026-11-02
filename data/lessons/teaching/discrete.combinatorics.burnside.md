# Burnside's Lemma

## Symmetry and Orbits

When configurations are considered the same up to rotation or reflection, the objects of interest are orbits of a group action. Burnside's lemma counts orbits by averaging the number of configurations fixed by each group element.

The group is usually a cyclic group of rotations or a dihedral group of rotations and reflections. Each element permutes the positions, and a configuration is fixed when every position receives the same colour as its image.

The average is over group elements, not over configurations, which is what makes the count tractable.

## Fixed Configurations

A rotation of n positions by d steps splits the positions into cycles of length n divided by the greatest common divisor of n and d. A colouring is fixed by that rotation exactly when every cycle is constant.

With k colours, a rotation with c cycles therefore fixes k to the power c configurations. Reflections require separate cycle counts because their cycle structure differs from rotations.

Listing the cycle count for every group element is the main computational step.

## Counting the Orbits

Add the fixed counts over all group elements and divide by the size of the group. For rotations the group size is n; for the full dihedral group it is 2n.

The result is always an integer, which is a useful check on the arithmetic. Small cases such as a square with two colours can be checked by drawing the orbits.

Burnside's lemma is the bridge between group actions and ordinary counting problems.
