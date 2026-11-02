# Half-Life and Exponential Decay

## The Half-Life Model

Half-life is the time a decaying quantity needs to fall to half its initial size. If the half-life is h and the initial amount is A0, then after time t the amount is A0 times (1/2)^(t/h). The exponent counts how many half-lives have passed, which is why the model is exponential. After one half-life half remains, after two half-lives a quarter remains, and after three an eighth remains.

## Computing Remaining Amount

To compute a remaining amount, divide the elapsed time by the half-life, raise one half to that power and multiply by the initial amount. When the elapsed time is a whole number of half-lives the arithmetic is exact and no logarithm is needed. When it is not, the exponent is fractional and a calculator is appropriate. The same model describes radioactive decay, medicine clearing from the body and any process that loses a fixed fraction in each equal time interval.

## Solving for Time

Solving for time reverses the calculation. Set the model equal to the target amount, divide by the initial amount, and take a logarithm of both sides. Using base one half gives t = h log_(1/2)(A/A0), and using the natural logarithm gives t = h ln(A/A0) / ln(1/2). A negative elapsed time means the target is larger than the initial amount, which is outside the decay story. Comparing two materials means comparing their half-lives, not their starting amounts, because the starting amount only scales the curve.
