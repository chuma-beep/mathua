# Higher derivatives of sine and cosine

Differentiating sine or cosine repeatedly returns to the starting function after four steps, so every higher derivative is one of four functions.

## The four-step cycle

The cycle for \(f(x) = \sin x\) reads \(f' = \cos x\), \(f'' = -\sin x\), \(f''' = -\cos x\) and
\(f'''' = \sin x\). Since 100 is divisible by 4, the hundredth derivative of \(\sin x\) is
\(\sin x\) again. The hundred and first is \(\cos x\). Cosine enters the same cycle one step
ahead, so its hundredth derivative is \(\cos x\) as well. The pattern repeats forever.

## Derivatives of a scaled angle

A scaled angle brings out a factor at every step. For \(y = \cos 2x\) the derivatives are
\(-2\sin 2x\), \(-4\cos 2x\), \(8\sin 2x\) and \(16\cos 2x\). The coefficient after \(n\) steps
has magnitude \(2^{n}\), with the sign from the cycle, so the tenth derivative is
\(-1024\cos 2x\). At \(x = 0\) its value is -1024. A modest scaling compounds quickly, while the
sign alternates as it goes.

## The pattern at zero

Evaluated at zero, the derivatives of sine alternate between 0 and \(\pm 1\): 0, 1, 0, -1 and
around again. The seventh derivative of \(\sin x\) at 0 is therefore -1, because 7 leaves
remainder 3 and the third derivative is \(-\cos x\). Cosine starts at 1 instead, so its sixth
derivative at 0 is -1.

The remainder on division by 4 selects the function; the value at zero settles the sign.
