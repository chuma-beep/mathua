package optimization

import (
	"strconv"
	"strings"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
)

var wave11IDs = []string{
	"optim.problem.anatomy", "optim.feasible.set", "optim.local_global",
	"optim.convex_set", "optim.strict_convexity", "optim.subgradient",
	"optim.epigraph", "optim.coercivity", "optim.first_order", "optim.second_order",
	"optim.hessian", "optim.gradient_step", "optim.step_size", "optim.newton_optim",
	"optim.bfgs", "optim.coordinate_descent", "optim.conjugate_direction",
	"optim.armijo", "optim.wolfe", "optim.trust_region_ratio",
	"optim.lp_standard_form", "optim.basic_solution", "optim.simplex_pivot",
	"optim.lp_dual", "optim.dual_function", "optim.weak_duality",
	"optim.strong_duality", "optim.complementary_slackness", "optim.shadow_price",
	"optim.lp_sensitivity", "optim.kkt", "optim.slater",
	"optim.equality_lagrange", "optim.active_constraint", "optim.penalty",
	"optim.barrier", "optim.augmented_lagrangian", "optim.projected_gradient",
	"optim.active_set", "optim.trust_region_subproblem", "optim.kkt_matrix",
	"optim.duality_gap", "optim.integer_program", "optim.lp_relaxation",
	"optim.branch_bound", "optim.cutting_plane", "optim.knapsack",
	"optim.transportation", "optim.set_cover", "optim.integrality_gap",
}

var wave12IDs = []string{
	"optim.central_path", "optim.path_following", "optim.primal_dual_interior",
	"optim.newton_system", "optim.conic_program", "optim.second_order_cone",
	"optim.semidefinite", "optim.dual_certificate", "optim.infeasibility_certificate",
	"optim.interior_convergence",
	"optim.sgd", "optim.mini_batch", "optim.variance_reduction", "optim.momentum",
	"optim.adagrad", "optim.adam", "optim.learning_rate_schedule",
	"optim.stochastic_noise", "optim.sample_complexity", "optim.online_gradient",
	"optim.pareto_dominance", "optim.pareto_front", "optim.weighted_sum",
	"optim.epsilon_constraint", "optim.goal_programming", "optim.tradeoff_rate",
	"optim.kkt_multiobjective", "optim.pareto_optimality", "optim.nondominated_sort",
	"optim.crowding_distance",
	"optim.sparsity", "optim.first_order_method", "optim.dual_coordinate",
	"optim.decomposition", "optim.admm", "optim.consensus",
	"optim.stochastic_averaging", "optim.subsampling", "optim.matrix_free",
	"optim.preconditioning",
	"optim.least_squares", "optim.regularized_least_squares", "optim.lasso",
	"optim.proximal_operator", "optim.proximal_gradient", "optim.soft_threshold",
	"optim.orthogonal_projection", "optim.quadratic_program",
	"optim.trust_region_global", "optim.nonsmooth_optimality",
}

var wave13IDs = []string{
	"optim.svrg", "optim.saga", "optim.sag", "optim.control_variate",
	"optim.importance_sampling", "optim.stratified_sampling",
	"optim.rank_one_sampling", "optim.svrg_epoch", "optim.gradient_aggregation",
	"optim.variance_bound", "optim.adadelta", "optim.rmsprop", "optim.adamw",
	"optim.amsgrad", "optim.bias_correction", "optim.learning_rate_decay",
	"optim.warmup", "optim.adaptive_regret", "optim.per_coordinate_rate",
	"optim.weight_decay", "optim.lp_as_socp", "optim.qp_as_socp",
	"optim.robust_lp", "optim.portfolio_socp", "optim.sdp_relaxation",
	"optim.max_cut_sdp", "optim.matrix_completion", "optim.certificate_check",
	"optim.conic_duality", "optim.cone_membership", "optim.dual_decomposition",
	"optim.primal_decomposition", "optim.consensus_admm",
	"optim.alternating_direction", "optim.block_coordinate",
	"optim.jacobi_iteration", "optim.gauss_seidel", "optim.column_generation",
	"optim.dantzig_wolfe", "optim.benders", "optim.proximal_point",
	"optim.moreau_envelope", "optim.dual_smoothing",
	"optim.smoothing_parameter", "optim.huber_loss", "optim.hinge_loss",
	"optim.subgradient_method", "optim.polyak_step", "optim.bundle_method",
	"optim.nonsmooth_kkt",
}

var wave14IDs = []string{
	"optim.minimax", "optim.saddle_point", "optim.zero_sum_game",
	"optim.robust_minimax", "optim.primal_dual_saddle", "optim.extragradient",
	"optim.optimistic_gradient", "optim.mirror_descent", "optim.dual_averaging",
	"optim.saddle_gap",
	"optim.optimal_control", "optim.control_hamiltonian", "optim.hjb",
	"optim.lqr", "optim.riccati", "optim.mpc", "optim.shooting",
	"optim.collocation", "optim.costate", "optim.bang_bang",
	"optim.implicit_diff", "optim.bilevel", "optim.meta_gradient",
	"optim.hypergradient", "optim.diff_layers", "optim.influence",
	"optim.sam", "optim.adversarial_training", "optim.dro",
	"optim.gradient_penalty",
	"optim.branch_cut", "optim.lagrangian_relaxation",
	"optim.subgradient_lagrangian", "optim.gomory_cut",
	"optim.cutting_stock", "optim.facility_location",
	"optim.tsp_relaxation", "optim.quadratic_assignment", "optim.big_m",
	"optim.sos1",
	"optim.kkt_solver", "optim.inertia_correction", "optim.regularized_newton",
	"optim.iterative_refinement", "optim.finite_difference",
	"optim.complex_step", "optim.autodiff", "optim.jfnk", "optim.lbfgs",
	"optim.line_search_interp",
}

// allOptimizationIDs is every concept the package owns, across all four waves.
var allOptimizationIDs = func() []string {
	ids := append([]string{}, wave11IDs...)
	ids = append(ids, wave12IDs...)
	ids = append(ids, wave13IDs...)
	return append(ids, wave14IDs...)
}()

func wave11Registry(t *testing.T) *generator.Registry {
	t.Helper()
	reg := generator.NewRegistry()
	Register(reg)
	return reg
}

// Every answer is a number the numeric grader can parse, and every question
// names the quantity it wants rather than a theorem.
func TestEveryAnswerIsNumeric(t *testing.T) {
	reg := wave11Registry(t)
	for _, id := range allOptimizationIDs {
		for _, kp := range []int{1, 2, 3} {
			for draw := 0; draw < 20; draw++ {
				p, err := reg.GenerateContext(id, generator.GeneratorContext{KPOrdinal: kp})
				if err != nil {
					t.Fatalf("%s kp%d: %v", id, kp, err)
				}
				if _, err := strconv.Atoi(p.Answer); err != nil {
					t.Fatalf("%s kp%d answer %q is not an integer", id, kp, p.Answer)
				}
				if p.Explanation == "" {
					t.Fatalf("%s kp%d has no explanation", id, kp)
				}
				if p.KPOrdinal != kp {
					t.Fatalf("%s reported kp%d for kp%d", id, p.KPOrdinal, kp)
				}
				lower := strings.ToLower(p.Question)
				if strings.Contains(lower, "theorem") {
					t.Fatalf("%s asks for a theorem name: %q", id, p.Question)
				}
			}
		}
	}
}

// The gate asks for five distinct questions per ordinal in twelve draws and a
// design target of twenty per ordinal in two hundred draws.
func TestEveryOrdinalHasEnoughVariety(t *testing.T) {
	reg := wave11Registry(t)
	for _, id := range allOptimizationIDs {
		for kp := 1; kp <= 3; kp++ {
			short := map[string]bool{}
			for draw := 0; draw < 12; draw++ {
				p, err := reg.GenerateContext(id, generator.GeneratorContext{KPOrdinal: kp})
				if err != nil {
					t.Fatalf("%s kp%d: %v", id, kp, err)
				}
				short[p.Question] = true
			}
			if len(short) < 5 {
				t.Errorf("%s kp%d produced %d distinct questions in 12 draws", id, kp, len(short))
			}
			long := map[string]bool{}
			for draw := 0; draw < 200; draw++ {
				p, err := reg.GenerateContext(id, generator.GeneratorContext{KPOrdinal: kp})
				if err != nil {
					t.Fatalf("%s kp%d: %v", id, kp, err)
				}
				long[p.Question] = true
			}
			if len(long) < 20 {
				t.Errorf("%s kp%d produced %d distinct questions in 200 draws", id, kp, len(long))
			}
		}
	}
}
