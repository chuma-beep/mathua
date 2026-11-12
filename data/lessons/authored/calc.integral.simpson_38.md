# Simpson's three-eighths rule

## The three-eighths weights

The three-eighths rule fits a cubic through four equally spaced points. On one panel of width \(3h\) the weights are 1, 3, 3 and 1, giving \((3h/8)(f_0 + 3f_1 + 3f_2 + f_3)\). The weights sum to 8, which is what makes a constant function integrate exactly.

## The composite rule

The composite rule repeats that panel. With \(n\) a multiple of three and \(h = (b-a)/n\), interior nodes at multiples of three carry weight 2 and all other interior nodes carry weight 3, while the two ends carry weight 1. Each block of three subintervals is one panel of the basic rule.

## One panel by hand

One panel by hand: \(\int_0^3 x^3 dx\) with \(h = 1\) uses the values 0, 1, 8 and 27. The rule gives \((3/8)(0 + 3 + 24 + 27) = (3/8)(54) = 20.25\). The exact value is \(81/4 = 20.25\), so the rule is exact for a cubic, as the error term predicts.

A complete answer names the relation, substitutes the data and states the result in the language of the problem. When a second method is available, compare the two results before moving on.
