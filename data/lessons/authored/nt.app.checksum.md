# Modular checksums

A checksum appends a digit that makes a weighted sum vanish modulo a small number. The modulus catches single-digit errors and many transpositions.

## ISBN-10

The ten digits are weighted from ten down to one, and the total is zero modulo eleven. The final digit is chosen to satisfy that congruence.

## ISBN-13

The thirteen digits alternate weights one and three, and the total is zero modulo ten. The alternating weights catch adjacent transpositions.

## The Luhn check

The Luhn algorithm doubles every second digit from the right and subtracts nine when the result exceeds nine. The check digit makes the total zero modulo ten.

To check your understanding, weight the digits, alternate one and three and double alternate digits.

A checksum is a congruence, and its error-detection power comes from choosing weights that separate the errors it must catch.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
