# Exponential regression

Exponential regression fits a model of the form \(y = a b^{x}\) to data that grow or decay by a constant factor. The two parameters have direct meanings: \(a\) is the value at \(x = 0\), and \(b\) is the factor by which the value changes when \(x\) increases by one.

## The growth factor

If the model is \(y = a b^{x}\), then consecutive values have the ratio \(y_{k+1}/y_k = b\). The ratio is constant, and that is the signature of exponential data.

For the points \((0, 4)\), \((1, 8)\) and \((2, 16)\), the ratios are \(8/4 = 2\) and \(16/8 = 2\). The growth factor is 2. A ratio above 1 means growth, and a ratio between 0 and 1 means decay.

Taking logarithms turns the model into a line: \(\ln y = \ln a + x \ln b\). A linear fit on the logged values recovers \(\ln b\) as the slope, and exponentiating gives \(b\).

## The initial value

The parameter \(a\) is the model's value at \(x = 0\), because \(b^{0} = 1\). If the data include a point at \(x = 0\), the initial value can be read directly.

When the table starts at \(x = 1\), the value at \(x = 0\) is recovered by dividing by the growth factor. A point \((1, 8)\) with growth factor 2 implies \(a = 8/2 = 4\).

The initial value is not the same as the first number in the table unless that number happens to sit at \(x = 0\). Checking the input before naming \(a\) avoids that error.

## Predicting with the fitted model

Prediction substitutes a new input into \(y = a b^{x}\). Once \(a\) and \(b\) are known, the arithmetic is a single power and a single product.

For \(y = 3 \cdot 2^{x}\), the prediction at \(x = 5\) is \(3 \times 32 = 96\). The same model gives \(3 \times 2^{0} = 3\) at the start, so the value has grown by a factor of 32 over five steps.

Predictions outside the data range are extrapolations. The arithmetic still works, but the claim that the pattern continues is an assumption about the situation, not a result of the regression.
