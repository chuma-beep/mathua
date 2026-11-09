# Power regression

Power regression fits a model \(y = a x^{b}\) to data where one quantity scales as a power of another. The two parameters are the leading coefficient and the exponent, and a logarithmic transformation turns the fit into a straight line.

## The exponent

Taking logarithms of both sides gives \(\ln y = \ln a + b \ln x\). In the logged variables the model is linear, and the exponent \(b\) is the slope.

The exponent describes how fast the output grows with the input. An exponent of 2 is a quadratic scaling, an exponent of 3 is cubic, and an exponent of \(1/2\) is a square-root scaling.

For the points \((1, 2)\), \((2, 8)\) and \((3, 18)\), the ratios \(8/2 = 4\) and \(18/2 = 9\) are \(2^2\) and \(3^2\). The exponent is 2.

## The leading coefficient

The coefficient \(a\) is the value at \(x = 1\), because \(1^{b} = 1\). It sets the overall scale of the curve without changing its shape.

When the data do not include \(x = 1\), the coefficient can still be recovered from any point once the exponent is known. Dividing the value by \(x^{b}\) gives \(a\).

In the example above, the point \((1, 2)\) gives \(a = 2\) directly. The full model is \(y = 2x^2\), which reproduces all three points.

## A predicted value

Prediction substitutes the new input into \(y = a x^{b}\). The exponent is applied before the multiplication by \(a\), and the order matters.

For \(y = 2x^2\), the prediction at \(x = 4\) is \(2 \times 16 = 32\). The same model gives 8 at \(x = 2\) and 18 at \(x = 3\).

Power models appear in geometry and physics whenever area or volume scales with a length. The exponent is often a small integer or a simple fraction, and the coefficient carries the units and the constant factor.
