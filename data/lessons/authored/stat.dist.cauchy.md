# The Cauchy distribution

The Cauchy distribution has density \(1/(\pi(1+x^2))\) centred at its median and no finite mean, so its tails are heavy enough to defeat averaging.

## Heavy tails

The density decays like \(1/x^2\), far slower than the normal's \(\exp(-x^2/2)\). The probability of a value beyond 100 is about 0.0064, so about 6 observations in 1000 land more than 100 units from the centre. No finite mean exists because the integral of \(x f(x)\) diverges.

## The quartiles

The distribution function is \(F(x) = 1/2 + \arctan(x)/\pi\), so the median is 0 and the quartiles solve \(\arctan(x) = \pm \pi/4\), giving \(x = \pm 1\). The interquartile range is therefore 2 for the standard Cauchy, however the distribution is shifted, while the variance does not exist.

## A worked spread

Take 9 ordered observations: minus 8, minus 1, minus 0.4, minus 0.1, 0, 0.2, 0.5, 4 and 1000. The median is the fifth value, 0, whatever the extremes do, while the mean is about 111. A single draw of 1000 moves the mean by more than 100 and leaves the median unchanged.

Heavy tails are not a defect of the model, they are the model. A summary that assumes finite variance is the wrong tool; the mean is not a number worth reporting for such data.
