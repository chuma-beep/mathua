package realanalysis

import "github.com/chuma-beep/mathua/internal/generator"

// Register adds the real-analysis generators to reg.
func Register(reg *generator.Registry) {
	// The order and the completeness of the real line.
	reg.Register("ra.order.supremum", &supremumGen{})
	reg.Register("ra.order.completeness", &completenessGen{})
	reg.Register("ra.order.archimedean", &archimedeanGen{})
	reg.Register("ra.order.nested_intervals", &nestedIntervalsGen{})
	reg.Register("ra.order.density", &densityGen{})
	reg.Register("ra.order.countable", &countableGen{})
	reg.Register("ra.order.uncountable", &uncountableGen{})
	reg.Register("ra.order.extended", &extendedGen{})
	// Sequences and series.
	reg.Register("ra.seq.convergence", &sequenceGen{})
	reg.Register("ra.seq.monotone", &monotoneGen{})
	reg.Register("ra.seq.subsequence", &subsequenceGen{})
	reg.Register("ra.seq.bolzano_weierstrass", &bolzanoWeierstrassGen{})
	reg.Register("ra.seq.cauchy", &cauchyGen{})
	reg.Register("ra.seq.series", &seriesGen{})
	reg.Register("ra.seq.tests", &ratioTestGen{})
	reg.Register("ra.seq.absolute", &absoluteGen{})
	reg.Register("ra.seq.rearrangement", &rearrangementGen{})
	reg.Register("ra.seq.cauchy_criterion", &cauchyCriterionGen{})
	// Limits and continuity.
	reg.Register("ra.limit.function", &functionLimitGen{})
	reg.Register("ra.limit.sequential", &sequentialGen{})
	reg.Register("ra.limit.continuity", &continuityGen{})
	reg.Register("ra.limit.ivt", &ivtGen{})
	reg.Register("ra.limit.evt", &evtGen{})
	reg.Register("ra.limit.uniform", &uniformGen{})
	reg.Register("ra.limit.compact", &compactGen{})
	reg.Register("ra.limit.connected", &connectedGen{})
	reg.Register("ra.limit.discontinuity", &discontinuityGen{})
	reg.Register("ra.limit.monotone_fn", &monotoneFnGen{})
	// Differentiation and integration, made rigorous.
	reg.Register("ra.deriv.mvt", &mvtGen{})
	reg.Register("ra.deriv.chain", &chainGen{})
	reg.Register("ra.deriv.taylor", &taylorGen{})
	reg.Register("ra.deriv.lhopital", &lhopitalGen{})
	reg.Register("ra.deriv.riemann", &riemannGen{})
	reg.Register("ra.deriv.darboux", &darbouxGen{})
	reg.Register("ra.deriv.ftc", &ftcGen{})
	reg.Register("ra.deriv.parts", &partsGen{})
	reg.Register("ra.deriv.improper", &improperGen{})
	reg.Register("ra.deriv.uniform_conv", &uniformConvGen{})
	reg.Register("ra.deriv.power_series", &powerSeriesGen{})
	reg.Register("ra.deriv.interchange", &interchangeGen{})
	// Metric spaces.
	reg.Register("ra.metric.open_closed", &openClosedGen{})
	reg.Register("ra.metric.closure", &closureGen{})
	reg.Register("ra.metric.convergence", &metricConvGen{})
	reg.Register("ra.metric.complete", &completeGen{})
	reg.Register("ra.metric.compact", &metricCompactGen{})
	reg.Register("ra.metric.continuity", &metricContinuityGen{})
	reg.Register("ra.metric.contraction", &contractionGen{})
	reg.Register("ra.metric.baire", &baireGen{})
	reg.Register("ra.metric.equivalent", &equivalentGen{})
	registerConversions(reg)
	registerWave5(reg)
	registerWave5Remaining(reg)
	registerWave6A(reg)
	registerWave6B(reg)
	for id, kind := range raBridgeKinds {
		reg.Register(id, &raBridgeGen{kind: kind})
	}
}
