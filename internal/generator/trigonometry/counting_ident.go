package trigonometry

import "github.com/chuma-beep/mathua/internal/generator"

func unitCircleFactsProblem() generator.Problem {
	pool := []trigVerdict{
		{"sin 30 degrees = 1/2", true, "The y-coordinate at 30 degrees is 1/2."},
		{"cos 60 degrees = 1/2", true, "The x-coordinate at 60 degrees is 1/2."},
		{"sin 90 degrees = 1", true, "The point at 90 degrees is (0, 1)."},
		{"cos 180 degrees = -1", true, "The point at 180 degrees is (-1, 0)."},
		{"tan 45 degrees = 1", true, "sin 45 and cos 45 are equal, so the ratio is 1."},
		{"sin 270 degrees = -1", true, "The point at 270 degrees is (0, -1)."},
		{"cos 0 degrees = 1", true, "The point at 0 degrees is (1, 0)."},
		{"tan 180 degrees = 0", true, "sin 180 is 0 and cos 180 is -1, so the ratio is 0."},
		{"sin 60 degrees = 1/2", false, "sin 60 is sqrt(3)/2."},
		{"cos 30 degrees = 1/2", false, "cos 30 is sqrt(3)/2."},
		{"tan 90 degrees = 0", false, "tan 90 is undefined because cos 90 is 0."},
		{"sin 180 degrees = 1", false, "sin 180 is 0."},
		{"cos 270 degrees = 1", false, "cos 270 is 0."},
		{"tan 30 degrees = 1/2", false, "tan 30 is 1/sqrt(3), about 0.577."},
		{"sin 45 degrees = sqrt(3)/2", false, "sin 45 is sqrt(2)/2."},
		{"cos 90 degrees = 1", false, "cos 90 is 0."},
	}
	return countTrigVerdicts("How many of these unit-circle statements are true?", pool)
}

func specialAngleProblem() generator.Problem {
	pool := []trigVerdict{
		{"sin 30 degrees = 1/2", true, "The 30-60-90 triangle gives a short side half the hypotenuse."},
		{"cos 30 degrees = sqrt(3)/2", true, "The adjacent side is the long leg of the 30-60-90 triangle."},
		{"tan 60 degrees = sqrt(3)", true, "tan 60 is sin 60 over cos 60, which is sqrt(3)."},
		{"sin 45 degrees = sqrt(2)/2", true, "The legs are equal, so the ratio is 1/sqrt(2)."},
		{"cos 45 degrees = sqrt(2)/2", true, "The legs are equal, so the ratio is 1/sqrt(2)."},
		{"tan 45 degrees = 1", true, "sin 45 and cos 45 are equal."},
		{"sin 60 degrees = sqrt(3)/2", true, "The opposite side is the long leg of the 30-60-90 triangle."},
		{"cos 60 degrees = 1/2", true, "The adjacent side is half the hypotenuse."},
		{"sin 30 degrees = sqrt(3)/2", false, "sin 30 is 1/2."},
		{"cos 30 degrees = 1/2", false, "cos 30 is sqrt(3)/2."},
		{"tan 60 degrees = 1/sqrt(3)", false, "tan 60 is sqrt(3); 1/sqrt(3) is tan 30."},
		{"sin 45 degrees = 1/2", false, "sin 45 is sqrt(2)/2."},
		{"cos 60 degrees = sqrt(3)/2", false, "cos 60 is 1/2."},
		{"tan 30 degrees = sqrt(3)/2", false, "tan 30 is 1/sqrt(3)."},
		{"sin 60 degrees = 1/2", false, "sin 60 is sqrt(3)/2."},
		{"tan 45 degrees = sqrt(2)/2", false, "tan 45 is 1."},
	}
	return countTrigVerdicts("How many of these exact-value statements are true?", pool)
}

func productToSumProblem() generator.Problem {
	pool := []trigVerdict{
		{"sin A cos B = 1/2[sin(A+B) + sin(A-B)]", true, "This is the standard product-to-sum identity."},
		{"cos A sin B = 1/2[sin(A+B) - sin(A-B)]", true, "Cosine first reverses the sign of the difference term."},
		{"cos A cos B = 1/2[cos(A+B) + cos(A-B)]", true, "The product of cosines gives a sum of cosines."},
		{"sin A sin B = 1/2[cos(A-B) - cos(A+B)]", true, "The product of sines gives a difference of cosines."},
		{"2 sin A cos B = sin(A+B) + sin(A-B)", true, "Multiply the first identity by 2."},
		{"2 cos A sin B = sin(A+B) - sin(A-B)", true, "Multiply the second identity by 2."},
		{"2 cos A cos B = cos(A+B) + cos(A-B)", true, "Multiplication by 2 clears the half factor."},
		{"2 sin A sin B = cos(A-B) - cos(A+B)", true, "The difference is ordered smaller first in this form."},
		{"sin A cos B = 1/2[sin(A+B) - sin(A-B)]", false, "The sine product has a plus sign."},
		{"cos A sin B = 1/2[sin(A+B) + sin(A-B)]", false, "The cosine-first product has a minus sign."},
		{"cos A cos B = 1/2[cos(A+B) - cos(A-B)]", false, "The cosine product has a plus sign."},
		{"sin A sin B = 1/2[cos(A+B) - cos(A-B)]", false, "The order of the cosine terms is reversed."},
		{"2 sin A cos B = sin(A+B) - sin(A-B)", false, "The difference form belongs to cos A sin B."},
		{"2 cos A sin B = sin(A+B) + sin(A-B)", false, "The cosine-first form has a minus sign."},
		{"2 cos A cos B = cos(A-B) - cos(A+B)", false, "The cosine product adds the two cosines."},
		{"2 sin A sin B = cos(A+B) - cos(A-B)", false, "The order of the cosine terms is reversed."},
	}
	return countTrigVerdicts("How many of these product-to-sum identities are stated correctly?", pool)
}

func sumToProductProblem() generator.Problem {
	pool := []trigVerdict{
		{"sin A + sin B = 2 sin((A+B)/2) cos((A-B)/2)", true, "The sum of sines uses the half-sum and half-difference."},
		{"sin A - sin B = 2 cos((A+B)/2) sin((A-B)/2)", true, "The difference of sines swaps sine and cosine."},
		{"cos A + cos B = 2 cos((A+B)/2) cos((A-B)/2)", true, "The sum of cosines is a product of cosines."},
		{"cos A - cos B = -2 sin((A+B)/2) sin((A-B)/2)", true, "The difference of cosines carries a minus sign."},
		{"sin A + sin B has factor sin((A+B)/2)", true, "The half-sum sine is one factor."},
		{"cos A + cos B has factor cos((A+B)/2)", true, "The half-sum cosine is one factor."},
		{"cos A - cos B has a leading minus sign", true, "The standard form is negative."},
		{"2 sin((A+B)/2) cos((A-B)/2) = sin A + sin B", true, "This is the same identity read right to left."},
		{"sin A + sin B = 2 cos((A+B)/2) cos((A-B)/2)", false, "The sum of sines uses a sine and a cosine factor."},
		{"sin A - sin B = 2 sin((A+B)/2) cos((A-B)/2)", false, "That is the sum, not the difference."},
		{"cos A + cos B = 2 sin((A+B)/2) sin((A-B)/2)", false, "The sum of cosines uses cosine factors."},
		{"cos A - cos B = 2 sin((A+B)/2) sin((A-B)/2)", false, "The minus sign is missing."},
		{"sin A + sin B has factor cos((A+B)/2)", false, "The half-sum factor is sine."},
		{"cos A - cos B has factor cos((A-B)/2)", false, "It uses two sine factors."},
		{"sin A - sin B has a leading minus sign", false, "Only the cosine difference carries the minus sign."},
		{"2 cos((A+B)/2) sin((A-B)/2) = sin A + sin B", false, "That product is sin A - sin B."},
	}
	return countTrigVerdicts("How many of these sum-to-product identities are stated correctly?", pool)
}

func pythagoreanVariantProblem() generator.Problem {
	pool := []trigVerdict{
		{"sin^2 x + cos^2 x = 1", true, "This is the Pythagorean identity."},
		{"1 + tan^2 x = sec^2 x", true, "Divide the Pythagorean identity by cos^2 x."},
		{"1 + cot^2 x = csc^2 x", true, "Divide the Pythagorean identity by sin^2 x."},
		{"sec^2 x - tan^2 x = 1", true, "Rearrange the tangent form."},
		{"csc^2 x - cot^2 x = 1", true, "Rearrange the cotangent form."},
		{"tan^2 x = sec^2 x - 1", true, "Subtract 1 from the tangent form."},
		{"cot^2 x = csc^2 x - 1", true, "Subtract 1 from the cotangent form."},
		{"sin^2 x = 1 - cos^2 x", true, "Rearrange the Pythagorean identity."},
		{"1 + tan^2 x = csc^2 x", false, "The tangent form gives sec^2 x."},
		{"1 + cot^2 x = sec^2 x", false, "The cotangent form gives csc^2 x."},
		{"sec^2 x + tan^2 x = 1", false, "The correct relation subtracts the tangent term."},
		{"sin^2 x - cos^2 x = 1", false, "The squares add to 1, not subtract."},
		{"sec^2 x - csc^2 x = 1", false, "The secant and cosecant forms are separate."},
		{"tan^2 x + cot^2 x = 1", false, "There is no such identity."},
		{"1 - tan^2 x = sec^2 x", false, "The tangent term is added."},
		{"cos^2 x - sin^2 x = 1", false, "The squares add to 1; their difference is cos 2x."},
	}
	return countTrigVerdicts("How many of these Pythagorean variants are correct?", pool)
}

func generalSolutionProblem() generator.Problem {
	pool := []trigVerdict{
		{"sin x = 0 gives x = 180k degrees", true, "The zeros of sine are multiples of 180 degrees."},
		{"cos x = 0 gives x = 90 + 180k degrees", true, "Cosine vanishes at the odd multiples of 90 degrees."},
		{"tan x = 0 gives x = 180k degrees", true, "Tangent shares the zeros of sine."},
		{"sin x = 1 gives x = 90 + 360k degrees", true, "The maximum occurs once per full turn."},
		{"cos x = 1 gives x = 360k degrees", true, "The maximum of cosine is at full turns."},
		{"sin x = -1 gives x = 270 + 360k degrees", true, "The minimum occurs once per full turn."},
		{"cos x = -1 gives x = 180 + 360k degrees", true, "The minimum of cosine is at half turns."},
		{"tan x = 1 gives x = 45 + 180k degrees", true, "Tangent repeats every 180 degrees."},
		{"sin x = 0 gives x = 90 + 180k degrees", false, "Those are the zeros of cosine."},
		{"cos x = 0 gives x = 180k degrees", false, "Those are the zeros of sine."},
		{"tan x = 0 gives x = 90 + 180k degrees", false, "Tangent is undefined there."},
		{"sin x = 1 gives x = 180 + 360k degrees", false, "That is cos x = -1."},
		{"cos x = 1 gives x = 180 + 360k degrees", false, "That is cos x = -1."},
		{"tan x = 1 gives x = 45 + 360k degrees", false, "Tangent repeats every 180 degrees, so half the solutions are missing."},
		{"sin x = -1 gives x = 90 + 360k degrees", false, "That is sin x = 1."},
		{"cos x = -1 gives x = 360k degrees", false, "That is cos x = 1."},
	}
	return countTrigVerdicts("How many of these general solutions are correct, with k an integer?", pool)
}

func advancedInequalityProblem() generator.Problem {
	pool := []trigVerdict{
		{"sin x > 0 on (0, 180) degrees, modulo 360", true, "Sine is positive in quadrants I and II."},
		{"cos x > 0 on (-90, 90) degrees, modulo 360", true, "Cosine is positive in quadrants I and IV."},
		{"tan x > 0 in quadrants I and III", true, "Sine and cosine have the same sign there."},
		{"sin x < 0 on (180, 360) degrees, modulo 360", true, "Sine is negative in quadrants III and IV."},
		{"cos x < 0 on (90, 270) degrees, modulo 360", true, "Cosine is negative in quadrants II and III."},
		{"sin x >= 1 only at 90 + 360k degrees", true, "The maximum value of sine is 1."},
		{"cos x <= -1 only at 180 + 360k degrees", true, "The minimum value of cosine is -1."},
		{"tan x < 0 in quadrants II and IV", true, "Sine and cosine have opposite signs there."},
		{"sin x > 0 on (180, 360) degrees", false, "Sine is negative there."},
		{"cos x > 0 on (90, 270) degrees", false, "Cosine is negative there."},
		{"tan x > 0 in quadrants II and IV", false, "Tangent is negative there."},
		{"sin x < 0 on (0, 180) degrees", false, "Sine is positive there."},
		{"cos x < 0 on (-90, 90) degrees", false, "Cosine is positive there."},
		{"sin x >= 1 for every x", false, "Sine never exceeds 1."},
		{"cos x <= -1 for every x", false, "Cosine is never below -1."},
		{"tan x < 0 in quadrants I and III", false, "Tangent is positive there."},
	}
	return countTrigVerdicts("How many of these inequality statements are correct, with k an integer?", pool)
}

func hyperbolicIdentityProblem() generator.Problem {
	pool := []trigVerdict{
		{"cosh^2 x - sinh^2 x = 1", true, "This is the fundamental hyperbolic identity."},
		{"sinh(x+y) = sinh x cosh y + cosh x sinh y", true, "This matches the sine addition formula with signs removed."},
		{"cosh(x+y) = cosh x cosh y + sinh x sinh y", true, "The hyperbolic cosine addition formula adds both products."},
		{"sinh 2x = 2 sinh x cosh x", true, "This is the double-angle formula for sinh."},
		{"cosh 2x = cosh^2 x + sinh^2 x", true, "This is one double-angle form for cosh."},
		{"cosh 2x = 2 cosh^2 x - 1", true, "This is the equivalent form using cosh only."},
		{"tanh x = sinh x / cosh x", true, "Tangent is defined as sine over cosine, and the hyperbolic case is analogous."},
		{"sech^2 x = 1 - tanh^2 x", true, "Divide the fundamental identity by cosh^2 x."},
		{"cosh^2 x + sinh^2 x = 1", false, "The fundamental identity subtracts."},
		{"sinh(x+y) = sinh x cosh y - cosh x sinh y", false, "That is the circular sine formula; the hyperbolic form adds."},
		{"cosh(x+y) = cosh x cosh y - sinh x sinh y", false, "That is the circular cosine formula; the hyperbolic form adds."},
		{"sinh 2x = sinh x cosh x", false, "The factor 2 is missing."},
		{"cosh 2x = cosh^2 x - sinh^2 x", false, "That difference is 1, not cosh 2x."},
		{"tanh x = cosh x / sinh x", false, "That is coth x."},
		{"sech^2 x = 1 + tanh^2 x", false, "The correct relation subtracts tanh^2 x."},
		{"cosh^2 x - sinh^2 x = 0", false, "The difference is 1."},
	}
	return countTrigVerdicts("How many of these hyperbolic identities are correct?", pool)
}
