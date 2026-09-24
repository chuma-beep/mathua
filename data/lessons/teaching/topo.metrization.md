# Urysohn Metrization Theorem

**Urysohn metrization:** Every second-countable regular Hausdorff space is metrizable. Three topological hypotheses buy a compatible metric; in particular compact Hausdorff plus second-countable implies metrizable.

## Worked: the interval is metrizable

Check [0, 1] against the three hypotheses:

1. Hausdorff: distinct points separate by disjoint subintervals.
2. Regular: a point and a disjoint closed set separate by shrinking intervals around each.
3. Second-countable: intervals with rational endpoints form a countable base.

So [0, 1] satisfies all three hypotheses, and the usual absolute-value metric is the compatible one.

## Worked: embedding into Hilbert space

See where the metric comes from for a general $X$:

1. Regularity plus countability yields countably many Urysohn functions $f_n:X\to[0, 1]$.
2. Assemble them into $F(x)=(f_1(x), f_2(x)/2, \dots)$ landing in the Hilbert cube.
3. $F$ is an embedding, and the cube is metric, so $X$ inherits a metric from it.

So second-countability plus regularity builds the embedding, and the metric is pulled back from the ambient cube.

## When metrizability fails

Test the Sorgenfrey line (lower-limit topology):

1. It is regular and Hausdorff: half-open intervals separate points from closed sets.
2. But it has no countable base: each $[x, x+e)$ forces uncountably many distinct basic opens.
3. Hence the three-hypothesis package breaks, and indeed no metric gives this topology.

So countable base is essential: without it, regularity plus Hausdorff cannot produce a metric.
