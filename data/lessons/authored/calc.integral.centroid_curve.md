# Centroid of a curve

The centroid of a curve is its balance point. It is a pair of averages weighted by arc length, so the answer depends on the geometry of the curve and not only on its endpoints.

## The x-coordinate

For a curve \(y=f(x)\) from \(x=a\) to \(x=b\), the arc length element is
\[
ds=\sqrt{1+(f'(x))^2}\,dx.
\]
The centroid coordinates are
\[
\bar x=\frac{1}{L}\int_a^b x\,ds,\qquad \bar y=\frac{1}{L}\int_a^b f(x)\,ds,
\]
where \(L\) is the arc length.

Each coordinate is an average over the curve. The denominator \(L\) is the total weight, and the numerator adds up the coordinate times the weight.

For a straight segment the arc length element is constant up to a factor. The factor cancels in the average, so the centroid is the midpoint of the endpoints.

## The y-coordinate

The \(y\)-coordinate uses the same weighted average with \(f(x)\) in place of \(x\). For a straight segment from \((x_1,y_1)\) to \((x_2,y_2)\),
\[
\bar y=\frac{y_1+y_2}{2}.
\]
The result does not depend on the slope, because the length weighting is uniform along a straight line.

For a curved graph the average is no longer the midpoint. A curve that spends more of its length near one end pulls the centroid toward that end.

## Arc length

The weight in both averages is the arc length. For a straight segment the length is the distance between the endpoints:
\[
L=\sqrt{(x_2-x_1)^2+(y_2-y_1)^2}.
\]
When the displacement is a Pythagorean triple, the length is an integer and the centroid coordinates stay rational.

Computing the length first is the usual order of work. The centroid is then one integral away, and for a straight segment the integral is just the average of the endpoint coordinates.
