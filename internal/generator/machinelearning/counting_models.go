package machinelearning

import "github.com/chuma-beep/mathua/internal/generator"

func nlpTransformerProblem() generator.Problem {
	return countMLVerdicts("How many of these statements about transformers in NLP are true?", []mlVerdict{
		{"Self-attention lets each token use context from the whole sequence", true, "every position attends to others"},
		{"Positional encoding supplies word order", true, "attention alone has no order"},
		{"Multi-head attention uses several representation subspaces", true, "each head learns a different projection"},
		{"Attention is quadratic in sequence length", true, "the score matrix is n by n"},
		{"A decoder uses a causal mask", true, "future tokens are hidden"},
		{"Pretraining is followed by fine-tuning", true, "a pretrained model is adapted"},
		{"Attention weights sum to one", true, "softmax normalises each row"},
		{"It handles long-range dependencies better than a vanilla RNN", true, "the path length is short"},
		{"It is only a bag of words", false, "order enters through attention and encodings"},
		{"It processes tokens with a recurrence", false, "there is no recurrent state"},
		{"It needs no positional information", false, "order must be injected"},
		{"It has only one attention head", false, "multi-head attention is standard"},
		{"Exact attention is linear in sequence length", false, "it is quadratic"},
		{"A decoder is bidirectional", false, "it is masked and causal"},
		{"It cannot model long-range dependencies", false, "that is a strength"},
		{"It is a fixed set of grammar rules", false, "it is a learned network"},
	})
}

func qlearningProblem() generator.Problem {
	return countMLVerdicts("How many of these statements about Q-learning are true?", []mlVerdict{
		{"It is off-policy", true, "it learns greedy values from exploratory behaviour"},
		{"It is model-free", true, "it needs no transition model"},
		{"It updates Q toward reward plus discounted maximum next Q", true, "the Bellman target"},
		{"Epsilon-greedy explores with probability epsilon", true, "otherwise it exploits"},
		{"Q(s,a) estimates the expected discounted return", true, "not the immediate reward"},
		{"It converges with sufficient exploration", true, "every pair must be visited often enough"},
		{"The target uses the maximum over next actions", true, "the greedy value"},
		{"A discount of zero reduces the target to the reward", true, "the future term vanishes"},
		{"It is on-policy", false, "that is SARSA"},
		{"It requires a model of the environment", false, "it is model-free"},
		{"It updates only a state value", false, "it updates action values"},
		{"It needs no exploration", false, "unvisited actions cannot be valued"},
		{"Q is the immediate reward only", false, "it is a return"},
		{"It converges without visiting a state", false, "the value stays unknown"},
		{"It is a policy-gradient method", false, "it is value-based"},
		{"The target is the reward alone", false, "the bootstrapped future term is included"},
	})
}

func policyGradientProblem() generator.Problem {
	return countMLVerdicts("How many of these statements about policy gradients are true?", []mlVerdict{
		{"REINFORCE is a Monte Carlo policy-gradient method", true, "it uses full-episode returns"},
		{"It is on-policy", true, "trajectories come from the current policy"},
		{"It uses sampled returns to weight the score function", true, "the REINFORCE estimator"},
		{"It is unbiased but has high variance", true, "a baseline reduces the variance"},
		{"A baseline does not add bias", true, "it only centres the estimator"},
		{"The advantage is Q minus V", true, "it measures the action against its state average"},
		{"Entropy regularisation encourages exploration", true, "it rewards uncertain policies"},
		{"It optimises the policy directly", true, "there is no separate value target"},
		{"It is off-policy", false, "it needs current trajectories"},
		{"It is unbiased with low variance by itself", false, "variance is the main problem"},
		{"It requires a model of the environment", false, "it is model-free"},
		{"It can only produce deterministic policies", false, "it usually produces stochastic policies"},
		{"A baseline introduces bias", false, "a state-only baseline leaves the gradient unbiased"},
		{"It needs no sampling", false, "the expectation is estimated from samples"},
		{"The advantage is V minus Q", false, "the sign is reversed"},
		{"It is the same as Q-learning", false, "one is policy-based, the other value-based"},
	})
}

func bertProblem() generator.Problem {
	return countMLVerdicts("How many of these statements about BERT are true?", []mlVerdict{
		{"BERT is a bidirectional encoder", true, "left and right context are used"},
		{"Masked language modelling is a pretraining task", true, "masked tokens are predicted"},
		{"Next sentence prediction is a pretraining task", true, "a second binary objective"},
		{"A pretrained BERT is fine-tuned per downstream task", true, "the same base is adapted"},
		{"BERT base has 12 layers and width 768", true, "the base configuration"},
		{"BERT base uses 12 attention heads", true, "head width is 64"},
		{"About 15 percent of tokens are masked", true, "the standard masking rate"},
		{"Its embeddings are contextual", true, "the representation depends on neighbours"},
		{"It is a decoder-only model", false, "it is an encoder"},
		{"It uses a causal mask", false, "it is bidirectional"},
		{"It only predicts the next token", false, "it fills masked tokens"},
		{"It needs no pretraining", false, "pretraining is the point"},
		{"BERT base has 24 layers", false, "that is BERT large"},
		{"It masks every token", false, "only about 15 percent"},
		{"It is unidirectional", false, "it is bidirectional"},
		{"It is trained from scratch for every task", false, "it is fine-tuned"},
	})
}

func actorCriticProblem() generator.Problem {
	return countMLVerdicts("How many of these statements about actor-critic methods are true?", []mlVerdict{
		{"The actor holds the policy", true, "it selects actions"},
		{"The critic estimates a value function", true, "it supplies the baseline"},
		{"The advantage is Q minus V", true, "the action is compared with its state value"},
		{"The TD error is reward plus discounted next value minus current value", true, "the one-step error"},
		{"A2C runs synchronous actors", true, "A3C is the asynchronous variant"},
		{"GAE generalises the advantage estimate", true, "it trades bias against variance"},
		{"DDPG is an actor-critic method for continuous actions", true, "a deterministic actor with a critic"},
		{"The critic provides a baseline", true, "it centres the policy update"},
		{"The actor estimates the value", false, "the critic does"},
		{"The critic holds the policy", false, "the actor does"},
		{"No baseline is used", false, "the critic is the baseline"},
		{"A2C is asynchronous", false, "A2C is synchronous"},
		{"The advantage is V minus Q", false, "the sign is reversed"},
		{"It only handles discrete actions", false, "continuous variants exist"},
		{"It never bootstraps", false, "the critic bootstraps from the next state"},
		{"The actor is updated with raw returns only", false, "it uses the advantage"},
	})
}

func logisticRegressionProblem() generator.Problem {
	return countMLVerdicts("How many of these statements about logistic regression are true?", []mlVerdict{
		{"It maps a linear score through the sigmoid", true, "the output is a probability"},
		{"Cross-entropy is its natural loss", true, "it is the Bernoulli likelihood"},
		{"Its decision boundary is linear", true, "the score is zero on a hyperplane"},
		{"It outputs a probability", true, "the sigmoid maps to 0 and 1"},
		{"The log-odds are linear in the inputs", true, "log(p/(1-p)) = w*x + b"},
		{"Its loss is convex", true, "no local minima"},
		{"A 0.5 threshold classifies by the sign of the score", true, "the boundary"},
		{"Coefficients encode feature effects on the log-odds", true, "a standard interpretation"},
		{"It outputs an unbounded real score", false, "the sigmoid bounds the output"},
		{"Mean squared error is its natural loss", false, "cross-entropy is"},
		{"Its boundary is always nonlinear", false, "it is linear"},
		{"It predicts a class directly with no probability", false, "it predicts a probability"},
		{"It has no intercept term", false, "a bias is standard"},
		{"It is a nonparametric model", false, "it has fixed coefficients"},
		{"It cannot be regularised", false, "L1 and L2 variants exist"},
		{"It requires balanced classes", false, "class weights can adjust the fit"},
	})
}

func vanishingGradientProblem() generator.Problem {
	return countMLVerdicts("How many of these statements about vanishing and exploding gradients are true?", []mlVerdict{
		{"The sigmoid derivative is at most 0.25", true, "its peak is one quarter"},
		{"Repeated multiplication of small derivatives shrinks the gradient", true, "depth turns it into a product"},
		{"ReLU helps mitigate vanishing gradients", true, "its derivative is 1 on the positive side"},
		{"Residual connections give gradients a short path", true, "the identity path carries them"},
		{"Gradient clipping targets exploding gradients", true, "it caps the update norm"},
		{"The gradient at an early layer is a product of layer derivatives", true, "the chain rule"},
		{"Batch normalisation can help stabilise gradients", true, "it controls activation scale"},
		{"LSTM gates mitigate vanishing gradients", true, "the cell state carries information"},
		{"The sigmoid derivative is above 1", false, "it peaks at 0.25"},
		{"Gradients in deep networks always explode", false, "they can vanish or explode"},
		{"ReLU causes vanishing gradients", false, "it reduces them"},
		{"Clipping fixes vanishing gradients", false, "it caps large gradients"},
		{"Depth has no effect on the gradient", false, "depth multiplies the factors"},
		{"Vanishing gradients are only a data problem", false, "they come from the architecture"},
		{"Residual connections worsen vanishing gradients", false, "they help"},
		{"A learning rate of zero is the remedy", false, "no learning would occur"},
	})
}

func imbalanceProblem() generator.Problem {
	return countMLVerdicts("How many of these statements about class imbalance are true?", []mlVerdict{
		{"Accuracy can be misleading under imbalance", true, "the majority class dominates"},
		{"Precision and recall separate the two error types", true, "they describe the positive class"},
		{"Resampling can address imbalance", true, "it changes class frequencies"},
		{"Class weights can address imbalance", true, "they change the cost"},
		{"A dataset with 990 negatives and 10 positives gives 99 percent accuracy to an always-negative model", true, "990/1000"},
		{"Precision-recall curves are often more informative than ROC under imbalance", true, "PR focuses on positives"},
		{"Undersampling discards majority data", true, "information is lost"},
		{"Stratified splits preserve class proportions", true, "each fold keeps the mix"},
		{"Accuracy alone is sufficient", false, "it hides minority errors"},
		{"The remedy is to remove the minority class", false, "that discards the signal"},
		{"An always-negative model gets 50 percent accuracy", false, "it gets the majority share"},
		{"Balanced accuracy always equals plain accuracy", false, "they differ under imbalance"},
		{"Class weights change the number of examples", false, "they change the loss weights"},
		{"Undersampling adds majority data", false, "it removes data"},
		{"Precision always equals recall", false, "they measure different errors"},
		{"The majority class is always the positive class", false, "it depends on the labelling"},
	})
}

func ensembleProblem() generator.Problem {
	return countMLVerdicts("How many of these statements about ensembles are true?", []mlVerdict{
		{"Bagging trains models on bootstrap samples", true, "resampling with replacement"},
		{"Random forests choose a random feature subset at each split", true, "the trees are decorrelated"},
		{"Bagging mainly reduces variance", true, "averaging stabilises predictions"},
		{"Boosting mainly reduces bias", true, "each learner corrects residuals"},
		{"Averaging partly independent errors cancels them", true, "the variance of the average falls"},
		{"Boosting trains learners sequentially", true, "each depends on the previous fit"},
		{"Stacking uses a meta-learner", true, "it combines base predictions"},
		{"More trees usually reduce variance up to a point", true, "the average stabilises"},
		{"Bagging mainly reduces bias", false, "it reduces variance"},
		{"Boosting only reduces variance", false, "it reduces bias primarily"},
		{"A random forest is a single tree", false, "it is an ensemble"},
		{"Bagging trains on the test set", false, "it trains on resampled training data"},
		{"Boosting trains learners independently in parallel", false, "the stages are sequential"},
		{"Averaging increases variance", false, "it reduces variance"},
		{"An ensemble always overfits", false, "averaging often helps"},
		{"A random forest has no randomness", false, "it uses bootstrap and feature sampling"},
	})
}
