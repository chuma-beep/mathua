# Richardson extrapolation

Richardson extrapolation combines two approximations with different step sizes to cancel the leading error term.

## The error expansion

If an approximation has error \(O(h^2)\), halving the step reduces the error by a factor of four. The leading errors at \(h\) and \(h/2\) are proportional.

The two approximations share the same limit, so a weighted difference can cancel the \(h^2\) term.

The cancellation works term by term in the error expansion.

## The extrapolation formula

For second-order error, the extrapolated value is \((4D(h/2)-D(h))/3\). The weights 4/3 and -1/3 cancel the leading term.

The result is fourth-order accurate when the original method is second-order and the error expansion is smooth.

The formula uses the ratio of the error factors, not the values themselves.

## Higher order and cautions

For an error of order \(h^p\), the correction factor is \(2^p\). Higher-order extrapolation needs more approximations.

Extrapolation assumes a regular error expansion. Near a singularity or with noisy data, it can make the estimate worse.

Two levels of extrapolation are common, but each additional level needs another step size.
