package arithmetic

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/chuma-beep/mathua/internal/generator"
)

// arithBridgeGen serves the concepts added by the 2026 content expansion.
// Every branch returns a number, and every knowledge point asks a different
// question rather than repeating the same shape with new digits.
type arithBridgeGen struct{ kind string }

func arithOrdinal(ctx generator.GeneratorContext) int {
	if ctx.KPOrdinal >= 1 && ctx.KPOrdinal <= 3 {
		return ctx.KPOrdinal
	}
	return 1 + rand.Intn(3)
}

func arithProb(o int, q, a, e string) generator.Problem {
	return generator.Problem{
		Question:    q,
		Answer:      a,
		Explanation: e,
		Facts:       map[string]string{"answer": a},
		KPOrdinal:   o,
	}
}

func (g *arithBridgeGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := arithOrdinal(ctx)
	switch g.kind {
	case "mult_two_by_two":
		switch o {
		case 1:
			a := 12 + rand.Intn(87)
			b := 12 + rand.Intn(87)
			return arithProb(o, fmt.Sprintf("What is %d times %d? (enter a number)", a, b), strconv.Itoa(a*b), fmt.Sprintf("Split both factors: %d x %d = %d.", a, b, a*b))
		case 2:
			a := 12 + rand.Intn(40)
			b := 12 + rand.Intn(40)
			return arithProb(o, fmt.Sprintf("A number times %d equals %d. What is the number? (enter a number)", a, a*b), strconv.Itoa(b), fmt.Sprintf("Divide the product by the known factor: %d / %d = %d.", a*b, a, b))
		default:
			a := 12 + rand.Intn(30)
			b := 12 + rand.Intn(30)
			return arithProb(o, fmt.Sprintf("A warehouse has %d crates with %d items in each crate. How many items are there in total? (enter a number)", a, b), strconv.Itoa(a*b), fmt.Sprintf("Multiply crates by items per crate: %d x %d = %d.", a, b, a*b))
		}
	case "div_two_digit":
		switch o {
		case 1:
			b := 12 + rand.Intn(29)
			q := 3 + rand.Intn(10)
			return arithProb(o, fmt.Sprintf("What is %d divided by %d? (enter a number)", b*q, b), strconv.Itoa(q), fmt.Sprintf("The divisor times the quotient gives the dividend: %d x %d = %d.", b, q, b*q))
		case 2:
			b := 12 + rand.Intn(29)
			q := 4 + rand.Intn(10)
			return arithProb(o, fmt.Sprintf("A number divided by %d equals %d. What is the number? (enter a number)", b, q), strconv.Itoa(b*q), fmt.Sprintf("Multiply the divisor by the quotient: %d x %d = %d.", b, q, b*q))
		default:
			b := 12 + rand.Intn(29)
			q := 3 + rand.Intn(8)
			r := 1 + rand.Intn(b-1)
			return arithProb(o, fmt.Sprintf("What is the remainder when %d is divided by %d? (enter a number)", b*q+r, b), strconv.Itoa(r), fmt.Sprintf("%d = %d x %d + %d, so the remainder is %d.", b*q+r, b, q, r, r))
		}
	case "divisibility":
		switch o {
		case 1:
			m := []int{2, 3, 4, 5, 6, 9, 10}[rand.Intn(7)]
			nums := make([]int, 4)
			count := 0
			for i := range nums {
				nums[i] = 20 + rand.Intn(280)
				if nums[i]%m == 0 {
					count++
				}
			}
			return arithProb(o, fmt.Sprintf("How many of these numbers are divisible by %d: %d, %d, %d and %d? (enter a number)", m, nums[0], nums[1], nums[2], nums[3]), strconv.Itoa(count), fmt.Sprintf("Check each number against the divisibility rule for %d. Exactly %d of the four are divisible.", m, count))
		case 2:
			a := 1 + rand.Intn(8)
			c := rand.Intn(10)
			d := (9 - (100*a+c)%9) % 9
			return arithProb(o, fmt.Sprintf("The three-digit number %d_%d is divisible by 9. What is the missing digit? (enter a number)", a, c), strconv.Itoa(d), fmt.Sprintf("The digit sum must be a multiple of 9. With digit %d the sum is %d.", d, a+d+c))
		default:
			nums := make([]int, 4)
			count := 0
			for i := range nums {
				nums[i] = 30 + rand.Intn(270)
				if nums[i]%6 == 0 {
					count++
				}
			}
			return arithProb(o, fmt.Sprintf("How many of these numbers are divisible by 6: %d, %d, %d and %d? (enter a number)", nums[0], nums[1], nums[2], nums[3]), strconv.Itoa(count), fmt.Sprintf("A number is divisible by 6 when it is divisible by both 2 and 3. Exactly %d of the four qualify.", count))
		}
	case "sub_zeros":
		switch o {
		case 1:
			a := 100 * (2 + rand.Intn(8))
			b := 11 + rand.Intn(a-20)
			return arithProb(o, fmt.Sprintf("What is %d minus %d? (enter a number)", a, b), strconv.Itoa(a-b), fmt.Sprintf("Borrow through the zeros, then subtract: %d - %d = %d.", a, b, a-b))
		case 2:
			a := 100 * (2 + rand.Intn(8))
			c := 11 + rand.Intn(a-20)
			return arithProb(o, fmt.Sprintf("%d minus a number equals %d. What is the number? (enter a number)", a, c), strconv.Itoa(a-c), fmt.Sprintf("The missing subtrahend is %d - %d = %d.", a, c, a-c))
		default:
			a := 100 * (3 + rand.Intn(7))
			b := 100 + rand.Intn(700)
			return arithProb(o, fmt.Sprintf("A fundraiser has a goal of %d dollars and has raised %d dollars. How many dollars are still needed? (enter a number)", a, b), strconv.Itoa(a-b), fmt.Sprintf("Subtract the amount raised from the goal: %d - %d = %d.", a, b, a-b))
		}
	case "fact_families":
		switch o {
		case 1:
			a := 2 + rand.Intn(11)
			b := 2 + rand.Intn(11)
			return arithProb(o, fmt.Sprintf("What is %d times %d? (enter a number)", a, b), strconv.Itoa(a*b), fmt.Sprintf("The fact family for %d and %d has product %d.", a, b, a*b))
		case 2:
			a := 2 + rand.Intn(11)
			b := 2 + rand.Intn(11)
			return arithProb(o, fmt.Sprintf("In the fact family with product %d and one factor %d, what is the other factor? (enter a number)", a*b, a), strconv.Itoa(b), fmt.Sprintf("Divide the product by the known factor: %d / %d = %d.", a*b, a, b))
		default:
			a := 2 + rand.Intn(11)
			b := 2 + rand.Intn(11)
			return arithProb(o, fmt.Sprintf("Using the fact family of %d and %d, what is %d divided by %d? (enter a number)", a, b, a*b, b), strconv.Itoa(a), fmt.Sprintf("The related division is %d / %d = %d.", a*b, b, a))
		}
	case "place_millions":
		switch o {
		case 1:
			n := 1000000 + rand.Intn(8999999)
			places := []int{1000000, 100000, 10000, 1000, 100, 10, 1}
			names := []string{"millions", "hundred-thousands", "ten-thousands", "thousands", "hundreds", "tens", "ones"}
			i := rand.Intn(7)
			d := (n / places[i]) % 10
			return arithProb(o, fmt.Sprintf("In the number %d, what is the value of the digit in the %s place? (enter a number)", n, names[i]), strconv.Itoa(d*places[i]), fmt.Sprintf("The digit is %d and its place is worth %d, so its value is %d.", d, places[i], d*places[i]))
		case 2:
			n := 1000000 + rand.Intn(8999999)
			return arithProb(o, fmt.Sprintf("How many whole millions are in %d? (enter a number)", n), strconv.Itoa(n/1000000), fmt.Sprintf("The millions period is the leading digit, which is %d.", n/1000000))
		default:
			n := 1000000 + rand.Intn(8999999)
			th := (n / 1000) % 1000
			return arithProb(o, fmt.Sprintf("In the number %d, what is the three-digit number formed by the thousands period? (enter a number)", n), strconv.Itoa(th), fmt.Sprintf("Group the digits in threes from the right; the thousands period is %03d.", th))
		}
	case "missing_addend":
		switch o {
		case 1:
			a := 2 + rand.Intn(48)
			b := 2 + rand.Intn(48)
			return arithProb(o, fmt.Sprintf("%d plus what number equals %d? (enter a number)", a, a+b), strconv.Itoa(b), fmt.Sprintf("Subtract the known addend from the total: %d - %d = %d.", a+b, a, b))
		case 2:
			a := 100 + rand.Intn(800)
			b := 20 + rand.Intn(80)
			return arithProb(o, fmt.Sprintf("%d plus what number equals %d? (enter a number)", a, a+b), strconv.Itoa(b), fmt.Sprintf("Use the inverse operation: %d - %d = %d.", a+b, a, b))
		default:
			need := 20 + rand.Intn(70)
			have := 3 + rand.Intn(need-3)
			return arithProb(o, fmt.Sprintf("A team needs %d points and has scored %d. How many more points does it need? (enter a number)", need, have), strconv.Itoa(need-have), fmt.Sprintf("The missing addend is %d - %d = %d.", need, have, need-have))
		}
	case "estimate_quotient":
		switch o {
		case 1:
			b := 12 + rand.Intn(18)
			q := 4 + rand.Intn(20)
			a := b*q + rand.Intn(b)
			b0 := 10 * ((b + 5) / 10)
			a0 := 10 * ((a + 5) / 10)
			if b0 == 0 {
				b0 = 10
			}
			return arithProb(o, fmt.Sprintf("Estimate %d divided by %d by rounding the dividend to the nearest hundred and the divisor to the nearest ten. What quotient do you get? (enter a number)", a, b), strconv.Itoa(a0/b0), fmt.Sprintf("%d rounds to %d and %d rounds to %d, so the estimate is %d divided by %d, which is %d.", a, a0, b, b0, a0, b0, a0/b0))
		case 2:
			b := 11 + rand.Intn(20)
			q := 3 + rand.Intn(15)
			a := b * q
			b0 := 10 * ((b + 5) / 10)
			if b0 == 0 {
				b0 = 10
			}
			a0 := 10 * ((a + 5) / 10)
			return arithProb(o, fmt.Sprintf("Estimate %d divided by %d by rounding both numbers to the nearest ten. What is the estimated quotient? (enter a number)", a, b), strconv.Itoa(a0/b0), fmt.Sprintf("Rounding gives %d and %d. Their quotient is about %d.", a0, b0, a0/b0))
		default:
			b := 12 + rand.Intn(20)
			q := 5 + rand.Intn(20)
			return arithProb(o, fmt.Sprintf("What is the exact quotient when %d is divided by %d? (enter a number)", b*q, b), strconv.Itoa(q), fmt.Sprintf("The product %d x %d is %d, so the quotient is %d.", b, q, b*q, q))
		}
	case "factor_multiple":
		switch o {
		case 1:
			m := 3 + rand.Intn(9)
			n := 10 + rand.Intn(80)
			next := m * ((n + m - 1) / m)
			return arithProb(o, fmt.Sprintf("What is the smallest multiple of %d that is greater than or equal to %d? (enter a number)", m, n), strconv.Itoa(next), fmt.Sprintf("Divide %d by %d and round up, then multiply back: the next multiple is %d.", n, m, next))
		case 2:
			m := 3 + rand.Intn(8)
			lo := 10 + rand.Intn(40)
			hi := lo + 20 + rand.Intn(60)
			count := hi/m - (lo-1)/m
			return arithProb(o, fmt.Sprintf("How many multiples of %d lie between %d and %d inclusive? (enter a number)", m, lo, hi), strconv.Itoa(count), fmt.Sprintf("The count is floor(%d/%d) minus floor(%d/%d), which is %d.", hi, m, lo-1, m, count))
		default:
			a := 2 + rand.Intn(8)
			b := 3 + rand.Intn(9)
			lcm := a * b
			for d := 2; d <= a && d <= b; d++ {
				if a%d == 0 && b%d == 0 {
					lcm = a / d * b
					break
				}
			}
			return arithProb(o, fmt.Sprintf("What is the least common multiple of %d and %d? (enter a number)", a, b), strconv.Itoa(lcm), fmt.Sprintf("The smallest positive common multiple of %d and %d is %d.", a, b, lcm))
		}
	case "neg_compare":
		switch o {
		case 1:
			a := -(1 + rand.Intn(30))
			b := -(1 + rand.Intn(30))
			if a == b {
				b--
			}
			larger := a
			if b > a {
				larger = b
			}
			return arithProb(o, fmt.Sprintf("Which is larger, %d or %d? (enter the larger number)", a, b), strconv.Itoa(larger), fmt.Sprintf("On the number line %d lies to the right of %d, so %d is larger.", larger, a+b-larger, larger))
		case 2:
			k := -(1 + rand.Intn(20))
			nums := make([]int, 4)
			count := 0
			for i := range nums {
				nums[i] = -30 + rand.Intn(60)
				if nums[i] < k {
					count++
				}
			}
			return arithProb(o, fmt.Sprintf("How many of these numbers are less than %d: %d, %d, %d and %d? (enter a number)", k, nums[0], nums[1], nums[2], nums[3]), strconv.Itoa(count), fmt.Sprintf("A number is less than %d when it lies further left. Exactly %d of the four qualify.", k, count))
		default:
			a := -(1 + rand.Intn(20))
			return arithProb(o, fmt.Sprintf("What is the largest integer that is less than %d? (enter a number)", a), strconv.Itoa(a-1), fmt.Sprintf("The integers decrease to the left, so the largest integer below %d is %d.", a, a-1))
		}
	case "even_odd":
		switch o {
		case 1:
			nums := make([]int, 4)
			count := 0
			for i := range nums {
				nums[i] = 10 + rand.Intn(90)
				if nums[i]%2 == 0 {
					count++
				}
			}
			return arithProb(o, fmt.Sprintf("How many of these numbers are even: %d, %d, %d and %d? (enter a number)", nums[0], nums[1], nums[2], nums[3]), strconv.Itoa(count), fmt.Sprintf("A number is even when its last digit is 0, 2, 4, 6 or 8. Exactly %d of the four are even.", count))
		case 2:
			a := 10 + rand.Intn(90)
			b := 10 + rand.Intn(90)
			parity := (a + b) % 2
			return arithProb(o, fmt.Sprintf("Is the sum %d + %d even or odd? Enter 0 for even and 1 for odd. (enter a number)", a, b), strconv.Itoa(parity), fmt.Sprintf("The parities are %d and %d, so the sum has parity %d.", a%2, b%2, parity))
		default:
			a := 10 + rand.Intn(90)
			b := 10 + rand.Intn(90)
			parity := (a * b) % 2
			return arithProb(o, fmt.Sprintf("Is the product %d times %d even or odd? Enter 0 for even and 1 for odd. (enter a number)", a, b), strconv.Itoa(parity), fmt.Sprintf("A product is odd only when both factors are odd, so the parity is %d.", parity))
		}
	case "standard_form":
		switch o {
		case 1:
			a := 2 + rand.Intn(8)
			b := 1 + rand.Intn(9)
			n := 3 + rand.Intn(5)
			val := (a*10 + b) * intPow10(n-1)
			return arithProb(o, fmt.Sprintf("What is %d.%d times 10 to the power %d in ordinary form? (enter a number)", a, b, n), strconv.Itoa(val), fmt.Sprintf("Move the decimal point %d places right: %d.%d becomes %d.", n, a, b, val))
		case 2:
			a := 2 + rand.Intn(8)
			b := 1 + rand.Intn(9)
			n := 3 + rand.Intn(5)
			val := (a*10 + b) * intPow10(n-1)
			return arithProb(o, fmt.Sprintf("What power of ten belongs in the standard form of %d? Write the number as %d.%d times ten to that power. (enter a number)", val, a, b), strconv.Itoa(n), fmt.Sprintf("The decimal point moves %d places left, so the exponent is %d.", n, n))
		default:
			n := 2 + rand.Intn(6)
			return arithProb(o, fmt.Sprintf("Write 0.%0*d as a power of ten. What is the exponent? (enter a number)", n, 1), strconv.Itoa(-n), fmt.Sprintf("The decimal point moves %d places right, so the exponent is minus %d.", n, n))
		}
	}
	return generator.Problem{}
}

func intPow10(n int) int {
	v := 1
	for i := 0; i < n; i++ {
		v *= 10
	}
	return v
}
