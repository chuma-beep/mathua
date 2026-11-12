# Cesaro summation

Cesaro summation assigns a value to some divergent series by averaging their partial sums. It extends ordinary summation without contradicting it.

## Averaging partial sums

Let \(s_n\) be the partial sums of a series. The Cesaro means are \(\sigma_n = (s_0 + s_1 + \dots + s_n)/(n+1)\). For Grandi's series \(1 - 1 + 1 - 1 + \dots\) the partial sums alternate between 1 and 0, and the averages tend to \(1/2\). The series is given the value \(1/2\).

## Regularity

If a series converges to \(s\), its Cesaro means also converge to \(s\). Take \(1/2 + 1/4 + 1/8 + \dots\), whose partial sums are \(1 - 2^{-n}\). The average of the first \(n\) partial sums is \(1 - (1 - 2^{-n})/n\), which tends to 1. Averaging never changes a convergent sum.

## Divergent series made summable

Some divergent series become summable. Where Abel summation also applies, the two values agree. The Fourier series of a jump function converges in the Cesaro sense at the jump, to the midpoint of the two one-sided limits. Fejer's theorem is exactly this statement for continuous functions.

Cesaro summation replaces a limit of terms with a limit of averages, and the average is often the more stable object.
