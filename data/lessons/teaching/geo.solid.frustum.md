# Frustums

Frustums is part of the geometry curriculum. This lesson develops the idea in three sections, from the defining picture through the main computation to the interpretation.

## Slicing a Cone

A frustum is the part of a cone or pyramid that remains after a parallel slice removes the top. The two circular faces are parallel, and their radii are usually written R for the larger base and r for the smaller one. The height h is the perpendicular distance between the faces, not the slanted distance along the side. Keeping those three quantities separate avoids most errors in frustum problems.

A cone cut by a plane parallel to its base produces a smaller similar cone above the cut. The removed cone and the original cone share the same shape, so their linear dimensions are proportional. That similarity is the source of every frustum formula. It also explains why the volume is not the average of the two cylinder volumes: the cross-sectional area changes quadratically with height.

Real objects such as buckets, lampshades, cups and truncated towers are frustums. The slant height is the distance from a point on the rim of one base to the nearest point on the rim of the other. It forms the hypotenuse of a right triangle whose legs are the height and the difference of the radii.

## Volume of a Frustum

The volume is V = pi h (R^2 + R r + r^2) / 3. The three terms in the bracket come from integrating the square of the linearly changing radius. When R equals r the formula reduces to pi r^2 h, the volume of a cylinder. When r equals zero it reduces to pi R^2 h / 3, the volume of a cone. Both checks are quick ways to catch a copied formula with a missing factor.

A useful way to remember the bracket is to think of three areas: the large base, the small base and the geometric mean of the two. The formula is the average of those contributions multiplied by the height, with the factor one third. The geometric mean term R r is the one most often forgotten, and omitting it underestimates the volume whenever the two radii differ.

To compute a numerical answer, square both radii, add their product and add the squares, then multiply by the height and by pi, and finally divide by three. Keeping the symbol pi until the last step makes cancellation easier. If a problem asks for the volume divided by pi, the arithmetic becomes an integer or a simple fraction and the final multiplication is unnecessary.

## Surface Area

The lateral surface area of a circular frustum is pi (R + r) l, where l is the slant height. The two circular bases contribute pi R^2 and pi r^2. The total surface area is the sum of all three pieces. For an open object such as a lampshade, one or both bases may be absent, so read the question carefully before adding them.

The slant height satisfies l^2 = h^2 + (R - r)^2. This is the Pythagorean theorem applied to the right triangle formed by the height, the radial difference and the slant. Many textbook problems choose radii and heights that form a Pythagorean triple, which keeps the slant height an integer. Recognising triples such as 3-4-5 or 5-12-13 saves time and reduces rounding.

Surface area scales differently from volume. Doubling every length multiplies areas by four and volumes by eight. That distinction matters when comparing a model with a full-size object, or when estimating how much paint and how much material a container needs. The frustum formulas make the scaling visible because each term contains either two lengths or three.
