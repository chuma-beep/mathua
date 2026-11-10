# Composition of transformations

## Applying transformations in order

A composition applies one transformation and then another. The order matters for rotations and reflections, and it can matter for translations combined with other maps.

Write the image after the first transformation before applying the second. Keeping an intermediate point prevents the two steps from being merged incorrectly.

## Inverse compositions

To undo a composition, undo the last transformation first and then the first one. The inverse of a composition is the composition of the inverses in reverse order.

Check an inverse by composing it with the original transformation. The result should return every point to itself.

## Coordinates and fixed points

A fixed point is unchanged by the transformation. Translations have no fixed points, while rotations and reflections can have fixed points or fixed lines. Coordinates make these properties easy to test.

A composition of two translations is another translation. A composition of two reflections can be a translation or a rotation, depending on whether the mirror lines are parallel or intersecting.
