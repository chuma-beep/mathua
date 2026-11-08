# Backward and forward error

An approximation can be wrong in two directions. Forward error asks how far the computed answer is from the true answer. Backward error asks how much the input would have to change for the computed answer to be exact. A method with small backward error has solved a nearby problem accurately.

## Two ways to be wrong

Forward error is the familiar one: \(|\hat x-x|\). It measures the distance between the answer you got and the answer you wanted. A small forward error means the digits are close.

Backward error measures the size of the perturbation to the data. If a computed root \(\hat x\) gives a small residual \(f(\hat x)\), then \(\hat x\) is an exact root of a slightly different function. The residual is the natural backward-error measure for root finding.

## Residuals and backward error

The residual is \(f(\hat x)\), not \(f(x)\). It can be computed without knowing the true root, which is what makes it practical. A tiny residual means the computed point sits almost on the graph, even if the true root is elsewhere.

Residual and forward error are not the same. A very flat graph can have a small residual far from the root, because a small height change corresponds to a large horizontal change. A steep graph does the opposite.

## Turning backward into forward

Near a simple root, the derivative converts one error into the other: \(|\hat x-x|\approx |f(\hat x)|/|f'(x)|\). A small residual divided by a small derivative can still be a large forward error. This ratio is the condition number of the root problem.

The lesson is to report both quantities. A method that drives the residual to zero has done its job in the backward sense, but the answer may still be far from the root if the problem is ill-conditioned. Conditioning belongs to the problem; stability belongs to the algorithm.
