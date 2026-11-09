# Logarithmic regression

Logarithmic regression fits a model \(y = a + b \ln x\) to data that grow quickly at first and then level off. The output is linear in the logarithm of the input, so a plot against \(\ln x\) is a straight line.

## The coefficient of the logarithm

The coefficient \(b\) is the change in \(y\) for each one-unit increase in \(\ln x\). A positive \(b\) means the output still grows, but more and more slowly.

If \(x\) is multiplied by \(e\), then \(\ln x\) increases by one, so the fitted value increases by \(b\). That is the natural unit for this model.

For the points \((1, 2)\), \((e, 5)\) and \((e^2, 8)\), the values rise by 3 for each unit of \(\ln x\). The coefficient is \(b = 3\).

## The intercept

The intercept \(a\) is the fitted value when \(\ln x = 0\), which happens at \(x = 1\). It is not the value at \(x = 0\), because the logarithm is undefined there.

The intercept sets the height of the curve. Changing \(a\) shifts the whole fitted line up or down without changing its slope.

In the example, the point at \(x = 1\) has \(y = 2\), so \(a = 2\). The model is \(y = 2 + 3 \ln x\).

## A predicted value

Prediction substitutes the new input into the logarithm and then applies the linear rule. For an input of the form \(e^{k}\), the logarithm is simply \(k\).

For \(y = 2 + 3 \ln x\), the prediction at \(x = e^{3}\) is \(2 + 3 \times 3 = 11\). At \(x = e\) it is 5, which matches the second data point.

The model is defined only for positive inputs. A logarithmic regression therefore cannot describe a trend through zero or through negative values, and that restriction belongs to the model rather than to the fitting procedure.
