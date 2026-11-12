# Extraneous solutions in trigonometric equations

Squaring both sides of a trigonometric equation is irreversible, so it can produce candidates that solve the squared equation and not the original.

## Where squaring adds roots

Squaring removes the sign of both sides. That removed sign is what creates the extra roots. Take
\(\sin x + \cos x = 1\). Squaring gives \(1 + 2\sin x\cos x = 1\), so \(\sin x\cos x = 0\) and the
candidates are 0, 90, 180 and 270 degrees. The original holds at 0 and 90 degrees and fails at the
other two, where the sum is -1 instead of 1.

## Testing each candidate

The test is substitution into the equation as written, before any transformation. At 180 degrees,
\(\sin x + \cos x\) is \(0 - 1 = -1\), so the candidate is rejected. At 270 degrees it is
\(-1 + 0 = -1\), also rejected. A useful shortcut is to note the sign of each side before
squaring: the right side is positive, so a candidate making the left side negative cannot work.

## Reporting the valid set

The same discipline applies to every squared equation. The valid set is often smaller than the
candidate list. For \(\cos x - \sin x = 1\), squaring gives \(1 - 2\sin x\cos x = 1\), the same
four candidates. Substitution keeps only 0 and 270 degrees, so this equation has two solutions in
a turn rather than four.

State the valid set alone and say that the rejected values came from squaring.
