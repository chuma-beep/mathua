package optimization

import (
	"fmt"

	"github.com/chuma-beep/mathua/internal/generator"
)

// Wave 13 extends the optimization domain to one hundred and fifty concepts.
// The families reuse the Wave 12 adapter: three integer-valued closures per
// concept, selected by the ordinal, so the numeric grader can check every
// answer.
var wave13Families = map[string][3]func() trio{
	// -----------------------------------------------------------------------
	// Variance-reduced methods
	// -----------------------------------------------------------------------
	"optim.svrg": {
		func() trio {
			g, h, mean := limit(2, 20), limit(2, 20), limit(2, 20)
			return mk(fmt.Sprintf("An SVRG correction uses the snapshot gradient %d for a point whose true gradient is %d, and the snapshot mean is %d. Enter g - h + m.", g, h, mean), g-h+mean,
				fmt.Sprintf("%d - %d + %d = %d, which keeps the estimator unbiased while shrinking its variance.", g, h, mean, g-h+mean))
		},
		func() trio {
			inner, full := limit(5, 40), limit(2, 12)
			return mk(fmt.Sprintf("An SVRG epoch runs %d inner steps and one full pass costing %d gradient evaluations. Enter the evaluations per epoch.", inner, full), inner+full,
				fmt.Sprintf("The epoch costs %d + %d = %d evaluations.", inner, full, inner+full))
		},
		func() trio {
			variance, steps := limit(6, 80), limit(1, 9)
			return mk(fmt.Sprintf("After %d inner steps an SVRG variance of %d is divided by %d. Enter the reduced variance using integer division.", steps, variance, steps+1), variance/(steps+1),
				fmt.Sprintf("%d divided by %d is %d.", variance, steps+1, variance/(steps+1)))
		},
	},
	"optim.saga": {
		func() trio {
			table, epochs := limit(4, 16), limit(2, 6)
			return mk(fmt.Sprintf("SAGA keeps %d stored gradients and performs %d epochs. Enter the total table reads.", table, epochs), table*epochs,
				fmt.Sprintf("%d reads per epoch over %d epochs is %d.", table, epochs, table*epochs))
		},
		func() trio {
			stored, delta := limit(1, 30), limit(-6, 6)
			return mk(fmt.Sprintf("A SAGA table stores the gradient %d and replaces it with %d added. Enter the stored value afterwards.", stored, delta), stored+delta,
				fmt.Sprintf("%d + (%d) = %d.", stored, delta, stored+delta))
		},
		func() trio {
			stored, reads := limit(3, 12), limit(2, 7)
			return mk(fmt.Sprintf("Each of %d stored gradients is read %d times per sweep. Enter the reads per sweep.", stored, reads), stored*reads,
				fmt.Sprintf("%d times %d is %d.", stored, reads, stored*reads))
		},
	},
	"optim.sag": {
		func() trio {
			base, count, extra := limit(1, 12), limit(2, 9), limit(0, 5)
			sum := base*count + extra
			return mk(fmt.Sprintf("SAG averages %d stored gradients whose sum is %d. Enter the average using integer division.", count, sum), sum/count,
				fmt.Sprintf("%d divided by %d is %d, which is the average the step uses.", sum, count, sum/count))
		},
		func() trio {
			total, updated := limit(6, 24), limit(1, 5)
			return mk(fmt.Sprintf("A SAG table has %d entries and a short run updates %d of them. Enter the number never touched.", total, updated), total-updated,
				fmt.Sprintf("%d - %d = %d entries stay at their initial value.", total, updated, total-updated))
		},
		func() trio {
			sum, old, delta := limit(20, 80), limit(1, 9), limit(-4, 4)
			return mk(fmt.Sprintf("A SAG stored sum is %d and the entry %d is replaced by %d. Enter the new sum.", sum, old, old+delta), sum+delta,
				fmt.Sprintf("The sum changes by (%d) - %d = %d, so it becomes %d.", old+delta, old, delta, sum+delta))
		},
	},
	"optim.control_variate": {
		func() trio {
			g, c, h, mean := limit(2, 20), limit(1, 5), limit(2, 20), limit(1, 10)
			return mk(fmt.Sprintf("A control-variate estimator is g - c(h - m) with g = %d, c = %d, h = %d and m = %d. Enter its value.", g, c, h, mean), g-c*(h-mean),
				fmt.Sprintf("%d - %d·(%d - %d) = %d.", g, c, h, mean, g-c*(h-mean)))
		},
		func() trio {
			base, corr := limit(1, 9), limit(1, 9)
			return mk(fmt.Sprintf("A correlated estimator keeps (1 - ρ²) of a variance of %d, with ρ = 0.%d. Enter the remaining variance in hundredths.", base*100, corr), base*(100-corr*corr),
				fmt.Sprintf("(1 - 0.%d²)·%d = %d hundredths.", corr, base*100, base*(100-corr*corr)))
		},
		func() trio {
			ratio, variance := limit(1, 9), limit(1, 9)
			return mk(fmt.Sprintf("The optimal control coefficient is c / v. With c = %d and v = %d it is a whole number. Enter it.", ratio*variance, variance), ratio,
				fmt.Sprintf("%d divided by %d is %d.", ratio*variance, variance, ratio))
		},
	},
	"optim.importance_sampling": {
		func() trio {
			weight, total := limit(1, 9), limit(10, 60)
			return mk(fmt.Sprintf("Importance weights of %d are normalised by a total of %d. Enter the numerator of one normalised weight over that denominator.", weight, total), weight,
				fmt.Sprintf("The normalised weight is %d/%d, so the numerator is %d.", weight, total, weight))
		},
		func() trio {
			ratio, base := limit(1, 9), limit(1, 9)
			return mk(fmt.Sprintf("A likelihood ratio is p/q with p = %d and q = %d. Enter p/q when it is a whole number.", ratio*base, base), ratio,
				fmt.Sprintf("%d divided by %d is %d.", ratio*base, base, ratio))
		},
		func() trio {
			weight, count := limit(2, 15), limit(2, 9)
			return mk(fmt.Sprintf("An importance weight of %d is summed over %d samples. Enter the total.", weight, count), weight*count,
				fmt.Sprintf("%d times %d is %d.", weight, count, weight*count))
		},
	},
	"optim.stratified_sampling": {
		func() trio {
			strata, each := limit(2, 12), limit(2, 14)
			return mk(fmt.Sprintf("A stratified sample draws %d items from each of %d strata. Enter the sample size.", each, strata), strata*each,
				fmt.Sprintf("%d strata times %d items is %d.", strata, each, strata*each))
		},
		func() trio {
			strata, within, between := limit(2, 5), limit(1, 9), limit(1, 20)
			return mk(fmt.Sprintf("A stratified variance is v/k + b with within-variance %d, %d strata and between-variance %d. Enter it.", within*strata, strata, between), within+between,
				fmt.Sprintf("%d/%d + %d = %d.", within*strata, strata, between, within+between))
		},
		func() trio {
			per, strata, extra := limit(2, 12), limit(2, 8), limit(0, 5)
			return mk(fmt.Sprintf("%d samples are split into %d strata with %d left over. Enter the size of the smallest stratum.", per*strata+extra, strata, extra), per,
				fmt.Sprintf("%d divided by %d is %d with remainder %d, so the smallest stratum holds %d.", per*strata+extra, strata, per*strata+extra-per*strata, extra, per))
		},
	},
	"optim.rank_one_sampling": {
		func() trio {
			u, v := limit(2, 15), limit(2, 9)
			return mk(fmt.Sprintf("A rank-one update is u vᵀ with u = %d and v = %d in one dimension. Enter its single entry.", u, v), u*v,
				fmt.Sprintf("%d times %d is %d.", u, v, u*v))
		},
		func() trio {
			kept, of := limit(1, 9), 10*limit(1, 10)
			return mk(fmt.Sprintf("A rank-one sketch keeps %d of %d coordinates. Enter the kept percentage.", kept, of), kept*100/of,
				fmt.Sprintf("%d·100/%d = %d percent.", kept, of, kept*100/of))
		},
		func() trio {
			rows, cols := limit(2, 20), limit(2, 20)
			return mk(fmt.Sprintf("A rank-one factor stores a vector of length %d and a vector of length %d. Enter the stored entries.", rows, cols), rows+cols,
				fmt.Sprintf("%d + %d = %d entries.", rows, cols, rows+cols))
		},
	},
	"optim.svrg_epoch": {
		func() trio {
			inner := limit(5, 40)
			return mk(fmt.Sprintf("An SVRG epoch runs %d inner steps and one snapshot. Enter its gradient evaluations.", inner), inner+1,
				fmt.Sprintf("%d inner steps plus the snapshot give %d evaluations.", inner, inner+1))
		},
		func() trio {
			epochs, cost := limit(2, 20), limit(2, 12)
			return mk(fmt.Sprintf("A run takes %d epochs at %d evaluations each. Enter the total.", epochs, cost), epochs*cost,
				fmt.Sprintf("%d times %d is %d.", epochs, cost, epochs*cost))
		},
		func() trio {
			epochs, period := limit(6, 40), limit(2, 5)
			return mk(fmt.Sprintf("A snapshot is refreshed every %d epochs. Over %d epochs, enter the number of refreshes.", period, epochs), epochs/period,
				fmt.Sprintf("%d divided by %d is %d refreshes.", epochs, period, epochs/period))
		},
	},
	"optim.gradient_aggregation": {
		func() trio {
			workers, each := limit(2, 12), limit(2, 10)
			return mk(fmt.Sprintf("%d workers each send %d gradients. Enter the aggregated count.", workers, each), workers*each,
				fmt.Sprintf("%d times %d is %d.", workers, each, workers*each))
		},
		func() trio {
			nodes := limit(2, 40)
			return mk(fmt.Sprintf("A ring all-reduce over %d nodes sends one message per ordered pair. Enter the messages.", nodes), nodes*(nodes-1),
				fmt.Sprintf("%d·(%d - 1) = %d messages.", nodes, nodes, nodes*(nodes-1)))
		},
		func() trio {
			a, b, w := limit(2, 20), limit(2, 20), limit(2, 6)
			return mk(fmt.Sprintf("Two partial sums %d and %d are combined with weights 1 and %d. Enter the aggregate.", a, b, w), a+w*b,
				fmt.Sprintf("%d + %d·%d = %d.", a, w, b, a+w*b))
		},
	},
	"optim.variance_bound": {
		func() trio {
			sigma := limit(2, 40)
			return mk(fmt.Sprintf("A gradient has standard deviation %d. Enter its variance.", sigma), sigma*sigma,
				fmt.Sprintf("%d squared is %d.", sigma, sigma*sigma))
		},
		func() trio {
			batch, reduced := limit(2, 5), limit(1, 9)
			return mk(fmt.Sprintf("A batch of %d averages gradients whose standard deviation is %d. Enter the squared reduced deviation.", batch, batch*reduced), reduced*reduced,
				fmt.Sprintf("The averaged deviation is %d, whose square is %d.", reduced, reduced*reduced))
		},
		func() trio {
			batch, bound := limit(2, 9), limit(1, 9)
			return mk(fmt.Sprintf("A variance bound is C/n with C = %d and n = %d. Enter the bound.", bound*batch, batch), bound,
				fmt.Sprintf("%d divided by %d is %d.", bound*batch, batch, bound))
		},
	},

	// -----------------------------------------------------------------------
	// Adaptive methods
	// -----------------------------------------------------------------------
	"optim.adadelta": {
		func() trio {
			g := 2 * limit(1, 25)
			return mk(fmt.Sprintf("Adadelta with decay one half adds half of the squared gradient %d². Enter the added amount.", g), g*g/2,
				fmt.Sprintf("Half of %d is %d.", g*g, g*g/2))
		},
		func() trio {
			value := limit(1, 30)
			return mk(fmt.Sprintf("A root mean square averages two equal squared values %d² and %d². Enter the root mean square.", value, value), value,
				fmt.Sprintf("The mean of %d is %d, whose square root is %d.", value*value, value*value, value))
		},
		func() trio {
			first, second := 2*limit(1, 10), 2*limit(1, 10)
			return mk(fmt.Sprintf("Two accumulated squares %d and %d are averaged. Enter the average.", first, second), (first+second)/2,
				fmt.Sprintf("(%d + %d)/2 = %d.", first, second, (first+second)/2))
		},
	},
	"optim.rmsprop": {
		func() trio {
			old, grad := 2*limit(1, 10), 2*limit(1, 8)
			return mk(fmt.Sprintf("RMSProp with decay one half updates the square average %d using the new square %d². Enter the new average.", old, grad), old/2+grad*grad/2,
				fmt.Sprintf("%d/2 + %d/2 = %d.", old, grad*grad, old/2+grad*grad/2))
		},
		func() trio {
			root, rate := limit(1, 20), limit(1, 9)
			return mk(fmt.Sprintf("A step divides a rate of %d by the root mean square %d. Enter the step.", rate*root, root), rate,
				fmt.Sprintf("%d divided by %d is %d.", rate*root, root, rate))
		},
		func() trio {
			a, b := limit(2, 20), limit(2, 20)
			return mk(fmt.Sprintf("Two running averages %d and %d are combined as (a + 3b)/4. Enter the numerator a + 3b.", a, b), a+3*b,
				fmt.Sprintf("%d + 3·%d = %d.", a, b, a+3*b))
		},
	},
	"optim.adamw": {
		func() trio {
			decay, scale := limit(1, 9), limit(1, 10)
			return mk(fmt.Sprintf("AdamW subtracts a decay of one tenth times %d from a weight of %d tens. Enter the subtracted amount.", decay, scale), decay*scale,
				fmt.Sprintf("%d tenths of %d tens is %d.", decay, scale, decay*scale))
		},
		func() trio {
			decay, weight := limit(1, 9), limit(2, 20)
			return mk(fmt.Sprintf("A decoupled penalty multiplies the decay %d by the weight %d. Enter the product.", decay, weight), decay*weight,
				fmt.Sprintf("%d times %d is %d.", decay, weight, decay*weight))
		},
		func() trio {
			first, second := limit(1, 9), limit(1, 9)
			return mk(fmt.Sprintf("Two regularisation strengths %d and %d are combined. Enter their sum.", first, second), first+second,
				fmt.Sprintf("%d + %d = %d.", first, second, first+second))
		},
	},
	"optim.amsgrad": {
		func() trio {
			a, b := limit(2, 30), limit(2, 30)
			return mk(fmt.Sprintf("AMSGrad keeps the maximum of the squared averages %d and %d. Enter the maximum.", a, b), maxInt(a, b),
				fmt.Sprintf("The larger of %d and %d is %d.", a, b, maxInt(a, b)))
		},
		func() trio {
			a, b, c := limit(2, 20), limit(2, 20), limit(2, 20)
			return mk(fmt.Sprintf("Three squared averages %d, %d and %d are reduced to their maximum. Enter it.", a, b, c), maxInt(a, maxInt(b, c)),
				fmt.Sprintf("The largest of the three is %d.", maxInt(a, maxInt(b, c))))
		},
		func() trio {
			root, scale := limit(1, 15), limit(1, 9)
			return mk(fmt.Sprintf("A step divides a gradient of %d tens by the square root %d. Enter the step.", scale*root, root), scale,
				fmt.Sprintf("%d divided by %d is %d.", scale*root, root, scale))
		},
	},
	"optim.bias_correction": {
		func() trio {
			moment := limit(1, 30)
			return mk(fmt.Sprintf("The first moment %d is corrected by 1/(1 - 0.9), which is a factor of ten. Enter the corrected moment.", moment), 10*moment,
				fmt.Sprintf("10·%d = %d.", moment, 10*moment))
		},
		func() trio {
			moment := limit(1, 30)
			return mk(fmt.Sprintf("The second moment %d is corrected by 1/(1 - 0.99), which is a factor of one hundred. Enter the corrected moment.", moment), 100*moment,
				fmt.Sprintf("100·%d = %d.", moment, 100*moment))
		},
		func() trio {
			steps := limit(1, 25)
			return mk(fmt.Sprintf("A bias correction divides by 1 - (1/2)^t and 2^t is %d. Enter t.", 1<<steps), steps,
				fmt.Sprintf("2^%d = %d.", steps, 1<<steps))
		},
	},
	"optim.learning_rate_decay": {
		func() trio {
			rate, step := limit(1, 10), limit(1, 9)
			return mk(fmt.Sprintf("A rate is η/(1 + k t) with η = %d and k t = %d, chosen so the result is whole. Enter the decayed rate.", rate*(1+step), step), rate,
				fmt.Sprintf("%d divided by %d is %d.", rate*(1+step), 1+step, rate))
		},
		func() trio {
			base, halvings := limit(1, 9), limit(1, 6)
			return mk(fmt.Sprintf("A rate of %d is halved %d times. Enter the result.", base<<halvings, halvings), base,
				fmt.Sprintf("%d divided by %d is %d.", base<<halvings, 1<<halvings, base))
		},
		func() trio {
			base := limit(1, 20)
			return mk(fmt.Sprintf("A linear decay from %d reaches three quarters of its value at the given time. Enter that value.", 4*base), 3*base,
				fmt.Sprintf("Three quarters of %d is %d.", 4*base, 3*base))
		},
	},
	"optim.warmup": {
		func() trio {
			rate := 2 * limit(1, 20)
			return mk(fmt.Sprintf("A warmup reaches the rate %d at its end. At the halfway step it is half. Enter that rate.", rate), rate/2,
				fmt.Sprintf("Half of %d is %d.", rate, rate/2))
		},
		func() trio {
			warm, decay := limit(2, 20), limit(2, 30)
			return mk(fmt.Sprintf("A schedule warms up for %d steps and decays for %d. Enter its total length.", warm, decay), warm+decay,
				fmt.Sprintf("%d + %d = %d.", warm, decay, warm+decay))
		},
		func() trio {
			base, step := limit(1, 10), limit(1, 6)
			return mk(fmt.Sprintf("A warmup scales a base rate of %d by step %d of a period of %d. Enter the scaled rate when the period is four times the step.", base*4, step, 4*step), base,
				fmt.Sprintf("The scale is %d/%d = 1/4, so %d becomes %d.", step, 4*step, base*4, base))
		},
	},
	"optim.adaptive_regret": {
		func() trio {
			root := limit(1, 20)
			return mk(fmt.Sprintf("An adaptive regret bound is the square root of T. With T = %d, enter the bound.", root*root), root,
				fmt.Sprintf("The square root of %d is %d.", root*root, root))
		},
		func() trio {
			exponent := limit(1, 25)
			return mk(fmt.Sprintf("A regret bound is log T with T = 2^t. At t = %d, enter the bound.", exponent), exponent,
				fmt.Sprintf("log of 2^%d is %d.", exponent, exponent))
		},
		func() trio {
			constant, root := limit(2, 12), limit(2, 12)
			return mk(fmt.Sprintf("A bound is C·sqrt(T) with C = %d and sqrt(T) = %d. Enter the bound.", constant, root), constant*root,
				fmt.Sprintf("%d times %d is %d.", constant, root, constant*root))
		},
	},
	"optim.per_coordinate_rate": {
		func() trio {
			first, second := limit(1, 20), limit(1, 20)
			return mk(fmt.Sprintf("Two coordinates use rates %d and %d. Enter their sum.", first, second), first+second,
				fmt.Sprintf("%d + %d = %d.", first, second, first+second))
		},
		func() trio {
			entries, coords := limit(2, 9), limit(2, 20)
			return mk(fmt.Sprintf("A diagonal preconditioner stores %d entries for each of %d coordinates. Enter the stored entries.", entries, coords), entries*coords,
				fmt.Sprintf("%d times %d is %d.", entries, coords, entries*coords))
		},
		func() trio {
			root, scale := limit(1, 15), limit(1, 9)
			return mk(fmt.Sprintf("A coordinate rate is η/sqrt(G) with η = %d and sqrt(G) = %d. Enter the rate.", scale*root, root), scale,
				fmt.Sprintf("%d divided by %d is %d.", scale*root, root, scale))
		},
	},
	"optim.weight_decay": {
		func() trio {
			decay, tenths := limit(1, 9), limit(1, 10)
			return mk(fmt.Sprintf("An L2 penalty of %d tenths times a weight of %d tens adds to the gradient. Enter the added amount.", decay, tenths), decay*tenths,
				fmt.Sprintf("%d tenths of %d tens is %d.", decay, tenths, decay*tenths))
		},
		func() trio {
			total, each := limit(2, 12), limit(1, 9)
			return mk(fmt.Sprintf("A penalty sums %d weights of %d each. Enter the total.", total, each), total*each,
				fmt.Sprintf("%d times %d is %d.", total, each, total*each))
		},
		func() trio {
			// The answer space is half the weight range. Thirty values rather
			// than twenty keeps 200 draws above the twenty-distinct design
			// target with certainty instead of leaving it to sampling luck.
			weight := 2 * limit(1, 30)
			return mk(fmt.Sprintf("Decoupled decay multiplies a weight of %d by one half. Enter the result.", weight), weight/2,
				fmt.Sprintf("Half of %d is %d.", weight, weight/2))
		},
	},
}

func registerWave13(reg *generator.Registry) {
	for id, families := range wave13Families {
		reg.Register(id, &w12Gen{families: families})
	}
	for id, families := range wave13FamiliesB {
		reg.Register(id, &w12Gen{families: families})
	}
	for id, families := range wave13FamiliesC {
		reg.Register(id, &w12Gen{families: families})
	}
}
