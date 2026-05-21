> Content sourced from [Algebrica](https://algebrica.org/root-mean-square/) — CC BY-NC 4.0

## What is the quadratic mean?

The **quadratic mean** , also called the root mean square, belongs to the general family of [power means](<../introduction-to-the-mean>). It is obtained by taking the square root of the [arithmetic mean](<../arithmetic-mean>) of the squared values in a dataset. This measure is particularly useful when the direction of the data is irrelevant, but the magnitude of each value matters.

Because each value is squared before being averaged, larger numbers have a stronger influence on the result. For this reason, the quadratic mean effectively represents quantities that combine according to quadratic relationships, where variations in magnitude must be preserved rather than canceled by sign.

In simple terms, it describes the equilibrium point of the squared distribution, providing a realistic measure of the average intensity or effective value of a set of data.


In general form, the quadratic mean is expressed as:

\\[M_2 = \sqrt{\frac{1}{n} \sum_{i=1}^{n} x_i^2} \\]

where \\(x_1, x_2, \ldots, x_n\\) are the observed values and \\(n\\) is the total number of elements.

  * The quadratic mean can be applied to any set of real numbers, positive or negative.
  * Since the calculation involves squaring each term, it is always greater than or equal to the arithmetic and geometric means.
  * It provides an accurate description of data representing magnitudes, such as voltage, acceleration, power, or standard deviation, where the overall strength of variation is more important than its direction.


## Example 1

Let’s consider the average temperature variation in a small mountain town during five consecutive autumn days. Because temperatures can fluctuate above and below zero, we will use the quadratic mean to capture the overall magnitude of the variation, regardless of sign.

Day| Temperature (°C)  
---|---  
Monday| –3.5  
Tuesday| 0.0  
Wednesday| 2.8  
Thursday| –1.6  
Friday| 3.2  
  
Substituting the observed values to the quadratic mean formula, we get:

\\[\begin{align} M_2 &= \sqrt{\frac{(-3.5)^2 + 0.0^2 + (2.8)^2 + (-1.6)^2 + (3.2)^2}{5}} \\\\[3pt] &= \sqrt{\frac{12.25 + 0.00 + 7.84 + 2.56 + 10.24}{5}} \\\\[3pt] & = \sqrt{\frac{32.89}{5}} \approx 2.56 \end{align} \\]


  * If we consider the arithmetic mean (0.18 °C), it is noticeably lower than the quadratic mean (2.56 °C).
  * This happens because the arithmetic mean takes into account the sign of each value, so negative temperatures offset the positive ones.
  * The quadratic mean, on the other hand, measures the overall magnitude of the variations, providing a more realistic picture of the actual thermal intensity during the period.


Hence, the quadratic mean temperature is approximately:

\\[M_2 = 2.56 \text{ °C} \\]

##### This result shows that, even though the temperature fluctuates above and below zero, the quadratic mean captures the overall intensity of these variations.
