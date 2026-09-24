package machinelearning

import (
	"math/rand"

	"github.com/chuma-beep/mathua/internal/generator"
)

func Register(reg *generator.Registry) {
	reg.Register("ml.perceptron", &perceptronGen{})
	reg.Register("ml.loss.mse", &mseGen{})
	reg.Register("ml.forward_pass", &forwardPassGen{})
	reg.Register("ml.gradient_descent", &gradientDescentGen{})
	reg.Register("ml.backpropagation", &backpropGen{})
	reg.Register("ml.loss.cross_entropy", &crossEntropyGen{})
	reg.Register("ml.regularization", &regularizationGen{})
	reg.Register("ml.overfitting", &overfittingGen{})
	reg.Register("ml.validation", &validationGen{})
	reg.Register("ml.optim.momentum", &momentumGen{})
	reg.Register("ml.regularization.dropout", &dropoutGen{})
	reg.Register("ml.arch.cnn", &cnnGen{})
	reg.Register("ml.arch.rnn", &rnnGen{})
	reg.Register("ml.optim.batchnorm", &batchNormGen{})
	reg.Register("ml.eval.roc", &rocGen{})
	reg.Register("ml.arch.transformer", &transformerGen{})
	reg.Register("ml.arch.gan", &ganGen{})
	reg.Register("ml.optim.adam", &adamDetailsGen{})
	reg.Register("ml.eval.calibration", &calibrationGen{})
	reg.Register("ml.data.augmentation", &augmentationGen{})
	reg.Register("ml.nlp.transformer", &nlpTransformerGen{})
	reg.Register("ml.reinforcement.qlearning", &qlearningGen{})
	reg.Register("ml.reinforcement.policy_gradient", &policyGradientGen{})
	reg.Register("ml.nlp.bert", &bertGen{})
	reg.Register("ml.reinforcement.actor_critic", &actorCriticGen{})
}

type backpropGen struct{}

func (g *backpropGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		question string
		answer   string
		exp      string
	}
	entries := []entry{
		{
			"In backpropagation, which rule from calculus is used to compute gradients through multiple layers?",
			"chain rule",
			"Backpropagation applies the chain rule to compute the gradient of the loss with respect to each weight.",
		},
		{
			"The forward pass computes the ____ given input data.",
			"output",
			"The forward pass propagates input through each layer to compute the network's output or prediction.",
		},
		{
			"The backward pass computes the ____ of the loss with respect to each weight.",
			"gradient",
			"The backward pass uses the chain rule to compute how much each weight contributed to the error.",
		},
		{
			"What optimization algorithm uses gradients computed by backpropagation to update weights?",
			"gradient descent",
			"Gradient descent updates weights in the opposite direction of the gradient to minimize the loss.",
		},
		{
			"In a neural network with sigmoid activation, the derivative at x=0 is σ(0)(1-σ(0)). What is it? (enter a number)",
			"0.25",
			"The derivative of the sigmoid function is \\(\\sigma(x)(1-\\sigma(x))\\); at x=0: 0.5·0.5=0.25, which makes backpropagation efficient.",
		},
		{
			"What does the learning rate control in gradient descent with backpropagation?",
			"step size",
			"The learning rate determines how large a step is taken in the direction of the negative gradient.",
		},
		{
			"Backpropagation is an application of the ____ rule from multivariable calculus.",
			"chain",
			"The multivariable chain rule allows computing how changes in each weight affect the final loss.",
		},
		{
			"What is the name of the problem where gradients become very small in deep networks, making learning slow?",
			"vanishing gradient",
			"Vanishing gradients occur when gradients get progressively smaller as they are backpropagated through many layers.",
		},
		{
			"Which activation function helps mitigate the vanishing gradient problem compared to sigmoid?",
			"ReLU",
			"ReLU (Rectified Linear Unit) has derivative 1 for positive inputs, helping gradients flow through deep networks.",
		},
		{
			"The chain rule means each layer's gradient depends on the gradient from the layer above it. Which layer: above or below? (type one word)",
			"above",
			"The chain rule means each layer's gradient depends on the gradient from the layer above it.",
		},
	}
	e := entries[rand.Intn(len(entries))]

	return generator.Problem{
		Question:    e.question,
		Answer:      e.answer,
		Explanation: e.exp,
	}
}

type perceptronGen struct{}

func (g *perceptronGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		question string
		answer   string
		exp      string
	}
	entries := []entry{
		{"What does a perceptron compute? (weighted sum + activation / convolution)", "weighted sum + activation", "A perceptron computes a weighted sum of inputs plus bias, then applies an activation function."},
		{"A perceptron with weights [0.5, -0.5] and bias 0, inputs [1,1] gives weighted sum?", "0", "0.5*1 + (-0.5)*1 + 0 = 0."},
		{"Perceptron activation is typically a ____ function.", "step (or sign)", "Original perceptron uses a step function to output 0 or 1."},
		{"Is a single perceptron able to solve XOR? (yes/no)", "no", "XOR is not linearly separable, so a single perceptron cannot solve it; needs a multilayer network."},
	}
	e := entries[rand.Intn(len(entries))]
	return generator.Problem{Question: e.question, Answer: e.answer, Explanation: e.exp}
}

type mseGen struct{}

func (g *mseGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		question string
		answer   string
		exp      string
	}
	entries := []entry{
		{"Mean squared error between prediction 2 and target 3 is?", "1", "MSE = (2-3)^2 = 1. Mean over one sample is 1."},
		{"MSE between predictions [0,1] and targets [1,1] is?", "0.5", "((0-1)^2 + (1-1)^2)/2 = (1+0)/2 = 0.5."},
		{"MSE is appropriate for ____ tasks.", "regression", "Mean squared error measures distance for continuous regression targets."},
		{"Derivative of MSE (1/2)(y-t)^2 w.r.t y is?", "y - t", "d/dy 1/2(y-t)^2 = y-t."},
	}
	e := entries[rand.Intn(len(entries))]
	return generator.Problem{Question: e.question, Answer: e.answer, Explanation: e.exp}
}

type forwardPassGen struct{}

func (g *forwardPassGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		question string
		answer   string
		exp      string
	}
	entries := []entry{
		{"In a forward pass, what is computed layer by layer?", "activations (or outputs)", "Forward pass propagates input through each layer's weights and activation to compute the output."},
		{"Forward pass of a linear layer with W=[[1,0],[0,1]] and x=[2,3] gives?", "[2,3]", "Identity matrix leaves the input unchanged: Wx = x."},
		{"What operation follows the weighted sum in each neuron?", "activation", "After computing the weighted sum plus bias, an activation function is applied."},
		{"Is forward pass done before or after backpropagation? (before/after)", "before", "Forward pass computes predictions; backpropagation then computes gradients."},
	}
	e := entries[rand.Intn(len(entries))]
	return generator.Problem{Question: e.question, Answer: e.answer, Explanation: e.exp}
}

type gradientDescentGen struct{}

func (g *gradientDescentGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		question string
		answer   string
		exp      string
	}
	entries := []entry{
		{"Gradient descent updates weights in the direction of the ____ gradient.", "negative (or opposite)", "Weights are moved opposite to the gradient to minimize loss."},
		{"What does learning rate control in gradient descent?", "step size", "Learning rate scales the gradient step; too large diverges, too small is slow."},
		{"With gradient 2 and learning rate 0.1, what is the weight update magnitude?", "0.2", "Update = learning_rate * gradient = 0.1*2 = 0.2."},
		{"What variant of gradient descent uses a single sample per update?", "stochastic (or SGD)", "Stochastic gradient descent updates after each sample, while batch uses the whole dataset."},
	}
	e := entries[rand.Intn(len(entries))]
	return generator.Problem{Question: e.question, Answer: e.answer, Explanation: e.exp}
}

type crossEntropyGen struct{}

func (g *crossEntropyGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct{ question, answer, exp string }
	entries := []entry{{"Cross-entropy loss between true label 1 and prediction 0.8 is (-log(0.8)≈?) (0.22/1.20)", "0.22", "CE = -log(0.8) ≈ 0.22."}, {"Cross-entropy is appropriate for ____ tasks.", "classification", "Cross-entropy for classification."}, {"Binary cross-entropy for target 0 and prediction 0.1 is (-log(0.9)≈?) (0.10/2.30)", "0.10", "CE = -log(0.9) ≈ 0.10."}, {"Does cross-entropy penalize confident wrong predictions more than MSE? (yes/no)", "yes", "CE grows ∞ as p→0."}}
	e := entries[rand.Intn(len(entries))]
	return generator.Problem{Question: e.question, Answer: e.answer, Explanation: e.exp}
}

type regularizationGen struct{}

func (g *regularizationGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct{ question, answer, exp string }
	entries := []entry{{"L2 regularization adds λ||w||² to the loss. What does λ control? (regularization strength)", "regularization strength", "λ trades off."}, {"Does L2 regularization shrink weights toward zero? (yes/no)", "yes", "Gradient 2λw pulls."}, {"What is the gradient of λ||w||² w.r.t w? (2λw)", "2λw", "d/dw λw² = 2λw."}, {"Is regularization used to reduce overfitting? (yes/no)", "yes", "Penalizing large weights."}}
	e := entries[rand.Intn(len(entries))]
	return generator.Problem{Question: e.question, Answer: e.answer, Explanation: e.exp}
}

type overfittingGen struct{}

func (g *overfittingGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct{ question, answer, exp string }
	entries := []entry{{"Overfitting is when training error is ____ than test error. (lower/higher)", "lower", "Training lower."}, {"What curve shape indicates overfitting when plotting error vs model capacity?", "U-shaped test error", "Test drops then rises."}, {"Does more training data generally reduce overfitting? (yes/no)", "yes", "More samples."}, {"Which error measures generalization: training or test error?", "test", "Test error."}}
	e := entries[rand.Intn(len(entries))]
	return generator.Problem{Question: e.question, Answer: e.answer, Explanation: e.exp}
}

type validationGen struct{}

func (g *validationGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct{ question, answer, exp string }
	entries := []entry{{"What is the purpose of a validation set?", "tune hyperparameters (or model selection)", "Validation for hyperparams."}, {"k-fold cross-validation with k=5 splits data into how many folds?", "5", "k folds."}, {"Should the test set be used for hyperparameter tuning? (yes/no)", "no", "Leakage."}, {"Does cross-validation give a more reliable estimate than a single split? (yes/no)", "yes", "Averaging."}}
	e := entries[rand.Intn(len(entries))]
	return generator.Problem{Question: e.question, Answer: e.answer, Explanation: e.exp}
}

type momentumGen struct{}

func (g *momentumGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct{ question, answer, exp string }
	entries := []entry{{"Momentum update: v_{t+1}=β v_t + g_t. What does β control? (momentum / decay)", "momentum", "β controls persistence."}, {"Adam combines momentum and ____ scaling.", "adaptive (or RMSProp)", "Adaptive."}, {"Does momentum help escape shallow local minima? (yes/no)", "yes", "Velocity carries."}, {"With β=0.9, the effective memory horizon is roughly how many steps? (10)", "10", "Horizon ≈1/(1-β)."}}
	e := entries[rand.Intn(len(entries))]
	return generator.Problem{Question: e.question, Answer: e.answer, Explanation: e.exp}
}

type dropoutGen struct{}

func (g *dropoutGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*4)
	type entry struct{ question, answer, exp string }
	easy := []entry{{"Dropout randomly zeroes activations with probability p during training. At test time, are activations scaled by (1-p) or kept as-is? (1-p)", "1-p", "Inverted dropout."}, {"Does dropout reduce co-adaptation of neurons? (yes/no)", "yes", "Robustness."}, {"Is dropout only applied during training? (yes/no)", "yes", "Disabled at inference."}}
	hard := []entry{{"With p=0.5, what is the expected number of active neurons out of 100? (50)", "50", "Expected (1-p)*100."}, {"Does dropout approximate an ensemble of 2^n subnetworks? (yes/no)", "yes", "Each mask subnetwork."}, {"Is dropout a form of regularization? (yes/no)", "yes", "Generalization."}}
	pool := easy
	if scale > 3 {
		pool = append(easy, hard...)
	}
	e := pool[rand.Intn(len(pool))]
	return generator.Problem{Question: e.question, Answer: e.answer, Explanation: e.exp}
}

type cnnGen struct{}

func (g *cnnGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*4)
	type entry struct{ question, answer, exp string }
	easy := []entry{{"CNNs use ____ layers to detect local patterns.", "convolutional", "Convolutions."}, {"Does weight sharing in CNNs reduce parameters vs fully connected? (yes/no)", "yes", "Same filter reused."}, {"What operation follows convolution in a CNN block? (activation + pooling)", "activation", "After convolution activation."}}
	hard := []entry{{"For a 5x5 input, 3x3 filter, stride 1, no padding, what is output size? (3x3)", "3x3", "(5-3)/1+1=3."}, {"Does max-pooling provide translation invariance? (yes/no)", "yes", "Robust to shifts."}, {"Are CNNs especially suited for image data? (yes/no)", "yes", "Local connectivity."}}
	pool := easy
	if scale > 3 {
		pool = append(easy, hard...)
	}
	e := pool[rand.Intn(len(pool))]
	return generator.Problem{Question: e.question, Answer: e.answer, Explanation: e.exp}
}

type rnnGen struct{}

func (g *rnnGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*4)
	type entry struct{ question, answer, exp string }
	easy := []entry{{"RNNs maintain a ____ state across time steps.", "hidden", "Hidden state."}, {"Does an RNN share weights across time steps? (yes/no)", "yes", "Same recurrent weights."}, {"What problem plagues vanilla RNNs on long sequences? (vanishing gradient)", "vanishing gradient", "Backprop multiplies small gradients."}}
	hard := []entry{{"Do LSTMs use gates (input, forget, output) to combat vanishing gradients? (yes/no)", "yes", "Gating."}, {"Is an RNN with hidden size 0 able to model sequences? (yes/no)", "no", "No capacity."}, {"Does unfolding an RNN yield a deep feedforward network? (yes/no)", "yes", "Unrolling creates depth."}}
	pool := easy
	if scale > 3 {
		pool = append(easy, hard...)
	}
	e := pool[rand.Intn(len(pool))]
	return generator.Problem{Question: e.question, Answer: e.answer, Explanation: e.exp}
}

type batchNormGen struct{}

func (g *batchNormGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*4)
	type entry struct{ question, answer, exp string }
	easy := []entry{{"Batch norm normalizes activations to have mean 0 and variance ____ (1)", "1", "Normalized."}, {"Does batch norm have learnable scale γ and shift β? (yes/no)", "yes", "γ·x̂+β."}, {"Is batch norm applied before or after activation? (before)", "before", "Typically before."}}
	hard := []entry{{"Does batch norm allow higher learning rates? (yes/no)", "yes", "Stabilized."}, {"At test time, does batch norm use batch statistics or running averages? (running)", "running", "Population stats."}, {"Does batch norm reduce internal covariate shift? (yes/no)", "yes", "Stated goal."}}
	pool := easy
	if scale > 3 {
		pool = append(easy, hard...)
	}
	e := pool[rand.Intn(len(pool))]
	return generator.Problem{Question: e.question, Answer: e.answer, Explanation: e.exp}
}

type rocGen struct{}

func (g *rocGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*4)
	type entry struct{ question, answer, exp string }
	easy := []entry{{"ROC plots TPR vs ____ (FPR)", "FPR", "ROC."}, {"Does AUC = 1 mean perfect classifier? (yes/no)", "yes", "Area 1."}, {"Is AUC = 0.5 equivalent to random guessing? (yes/no)", "yes", "Diagonal."}}
	hard := []entry{{"For imbalanced data, is precision-recall often more informative than ROC? (yes/no)", "yes", "PR emphasizes."}, {"Does higher AUC mean better ranking of positives above negatives? (yes/no)", "yes", "AUC = P(score(pos)>score(neg))."}, {"If classifier thresholds from 0 to 1, does ROC go from (1,1) to (0,0) or opposite? (0,0→1,1)", "0,0 to 1,1", "Same curve."}}
	pool := easy
	if scale > 3 {
		pool = append(easy, hard...)
	}
	e := pool[rand.Intn(len(pool))]
	return generator.Problem{Question: e.question, Answer: e.answer, Explanation: e.exp}
}

type transformerGen struct{}

func (g *transformerGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*4)
	type entry struct{ question, answer, exp string }
	easy := []entry{{"With n=4 tokens, what shape is the QK^T attention matrix?", "4x4", "Every query meets every key: 4 queries by 4 keys gives a 4x4 score matrix."}, {"With d_model=512 and h=8 heads, what is the per-head dimension?", "64", "Width splits across heads: 512/8 = 64 dimensions per head."}, {"Attention weights [0.7,0.2,0.1] sum to what?", "1", "Softmax normalizes scores so the weights sum to 1."}, {"Values [10,20,30] with weights [1,0,0] give what output?", "10", "All weight on the first value returns it: 1*10+0*20+0*30 = 10."}, {"Is self-attention without positional encoding permutation equivariant? (yes/no)", "yes", "Reordering inputs reorders outputs identically; encodings add order."}, {"Is scaled dot-product attention linear O(n) in sequence length? (yes/no)", "no", "The n-by-n score matrix makes it quadratic O(n^2)."}}
	hard := []entry{{"Is scaled dot-product attention O(n²) in sequence length? (yes/no)", "yes", "n×n matrix."}, {"Does positional encoding use sin/cos of different frequencies? (yes/no)", "yes", "Sinusoidal."}, {"Is transformer decoder autoregressive with causal mask? (yes/no)", "yes", "Mask prevents future."}}
	pool := easy
	if scale > 3 {
		pool = append(easy, hard...)
	}
	e := pool[rand.Intn(len(pool))]
	return generator.Problem{Question: e.question, Answer: e.answer, Explanation: e.exp}
}

type ganGen struct{}

func (g *ganGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*4)
	type entry struct{ question, answer, exp string }
	easy := []entry{{"At the Nash equilibrium, what value does D output everywhere?", "0.5", "Generator matches data so real and fake are indistinguishable: D = 0.5."}, {"D scores 4 fakes [0.1,0.4,0.5,0.9]; how many fool D (score above 0.5)?", "1", "Only 0.9 exceeds 0.5, so 1 of 4 fakes fools the discriminator."}, {"G covers 2 of 10 data modes; what fraction is covered?", "0.2", "Coverage 2/10 = 0.2, a sign of mode collapse."}, {"A GAN trains how many networks: generator and discriminator?", "2", "Two players: G generates fakes, D classifies real vs fake."}, {"Does GAN training aim for a Nash equilibrium? (yes/no)", "yes", "Neither player can improve alone at equilibrium."}, {"Is mode collapse a sign of healthy GAN training? (yes/no)", "no", "Collapse means G emits few modes; healthy G covers the data."}}
	hard := []entry{{"Does mode collapse mean G outputs limited diversity? (yes/no)", "yes", "Few modes."}, {"Is Wasserstein GAN using Earth mover distance? (yes/no)", "yes", "WGAN."}, {"Does D provide gradient for G? (yes/no)", "yes", "Via backprop."}}
	pool := easy
	if scale > 3 {
		pool = append(easy, hard...)
	}
	e := pool[rand.Intn(len(pool))]
	return generator.Problem{Question: e.question, Answer: e.answer, Explanation: e.exp}
}

type adamDetailsGen struct{}

func (g *adamDetailsGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*4)
	type entry struct{ question, answer, exp string }
	easy := []entry{{"First Adam step with gradient 2, beta1=0.9, beta2=0.999, eta=0.1: what is the update magnitude?", "0.1", "Bias-corrected m-hat=2, v-hat=4, so 0.1*2/sqrt(4) = 0.1."}, {"After step 1 with beta1=0.9, m=0.18. What is bias-corrected m-hat?", "1.8", "Divide out the initialization bias: 0.18/0.1 = 1.8."}, {"With beta=0.9, the effective memory horizon is roughly how many steps?", "10", "Horizon is about 1/(1-beta) = 1/0.1 = 10."}, {"Adam default beta1 equals what?", "0.9", "Standard defaults: beta1=0.9, beta2=0.999."}, {"Does Adam maintain first and second moment estimates? (yes/no)", "yes", "m tracks mean gradient, v tracks mean squared gradient."}, {"Is Adam non-adaptive with one global step size? (yes/no)", "no", "Each weight gets its own effective rate eta/sqrt(v)."}}
	hard := []entry{{"Does Adam update = -η m̂/(√v̂+ε)? (yes/no)", "yes", "Rule."}, {"Is Adam adaptive per-parameter? (yes/no)", "yes", "Per-weight."}, {"Does Adam combine momentum and RMSProp? (yes/no)", "yes", "Hybrid."}}
	pool := easy
	if scale > 3 {
		pool = append(easy, hard...)
	}
	e := pool[rand.Intn(len(pool))]
	return generator.Problem{Question: e.question, Answer: e.answer, Explanation: e.exp}
}

type calibrationGen struct{}

func (g *calibrationGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*4)
	type entry struct{ question, answer, exp string }
	easy := []entry{{"A bin holds 100 of 200 samples with |acc-conf|=0.2. What is its ECE contribution?", "0.1", "Weight 100/200 = 0.5 times gap 0.2 gives 0.1."}, {"100 predictions at 80 percent confidence with 80 correct are miscalibrated by how much?", "0", "Accuracy 0.8 matches confidence 0.8: gap 0, calibrated."}, {"60 correct of 100 at 80 percent confidence means overconfidence by what fraction?", "0.2", "Confidence 0.8 minus accuracy 0.6 gives 0.2 overconfidence."}, {"Reliability diagrams bin predictions by what: confidence or loss?", "confidence", "Predictions group by confidence; accuracy plots per bin."}, {"Is a calibrated classifier's confidence approximately its accuracy? (yes/no)", "yes", "Calibration means confidence matches empirical accuracy."}, {"Does temperature scaling change the predicted class? (yes/no)", "no", "Scaling logits keeps argmax fixed; only confidence changes."}}
	hard := []entry{{"Does overconfidence mean ECE high? (yes/no)", "yes", "Gap."}, {"Is temperature scaling a calibration method? (yes/no)", "yes", "Softmax temp."}, {"Does well-calibrated with 80% conf mean 80% correct? (yes/no)", "yes", "Definition."}}
	pool := easy
	if scale > 3 {
		pool = append(easy, hard...)
	}
	e := pool[rand.Intn(len(pool))]
	return generator.Problem{Question: e.question, Answer: e.answer, Explanation: e.exp}
}

type augmentationGen struct{}

func (g *augmentationGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*4)
	type entry struct{ question, answer, exp string }
	easy := []entry{{"Is data augmentation used to increase training data artificially? (yes/no)", "yes", "Transforms."}, {"Does augmentation reduce overfitting? (yes/no)", "yes", "More variety."}, {"Is augmentation applied at test time? (no)", "no", "Test clean."}}
	hard := []entry{{"Is cutout/mixup a form of augmentation? (yes/no)", "yes", "Masks."}, {"Does augmentation approximate invariance? (yes/no)", "yes", "Ignore transforms."}, {"Is augmentation a form of regularization? (yes/no)", "yes", "Effective data."}}
	pool := easy
	if scale > 3 {
		pool = append(easy, hard...)
	}
	e := pool[rand.Intn(len(pool))]
	return generator.Problem{Question: e.question, Answer: e.answer, Explanation: e.exp}
}

type nlpTransformerGen struct{}

func (g *nlpTransformerGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*4)
	type entry struct{ q, a, e string }
	easy := []entry{{"Does transformer use self-attention for NLP? (yes/no)", "yes", "Attention."}, {"Is BERT bidirectional transformer? (no)", "no", "BERT is, but this asks transformer generally."}, {"Does positional encoding add order to transformer? (yes/no)", "yes", "Sinusoidal."}}
	hard := []entry{{"Is attention O(n²) in sequence length? (yes/no)", "yes", "Quadratic."}, {"Does transformer decoder use causal mask? (yes/no)", "yes", "Autoregressive."}, {"Is multi-head attention multiple subspaces? (yes/no)", "yes", "Heads."}}
	pool := easy
	if scale > 3 {
		pool = append(easy, hard...)
	}
	e := pool[rand.Intn(len(pool))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

type qlearningGen struct{}

func (g *qlearningGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*4)
	type entry struct{ q, a, e string }
	easy := []entry{{"Q=0, alpha=0.5, r=2, gamma=0.9, max next Q=4. What is the updated Q?", "2.8", "Target 2+0.9*4 = 5.6; Q moves halfway: 0+0.5*5.6 = 2.8."}, {"With epsilon=0.1, an epsilon-greedy policy exploits with what probability?", "0.9", "Exploit with 1-epsilon = 0.9, explore randomly with 0.1."}, {"With discount gamma=0, the Bellman target equals what: r or r+Q?", "r", "Zero discount drops the future term, leaving reward r."}, {"Q(s,a) estimates expected return or immediate reward: return or reward?", "return", "Q is the expected discounted return from s taking a."}, {"Is Q-learning off-policy? (yes/no)", "yes", "It learns greedy values from any exploratory behavior."}, {"Does Q-learning need a model of the environment? (yes/no)", "no", "Model-free: updates come from sampled transitions."}}
	hard := []entry{{"Does Q-learning converge with infinite exploration? (yes/no)", "yes", "Convergence."}, {"Is Q-learning model-free? (yes/no)", "yes", "No model."}, {"Does overestimation bias exist in Q-learning? (yes/no)", "yes", "Max bias."}}
	pool := easy
	if scale > 3 {
		pool = append(easy, hard...)
	}
	e := pool[rand.Intn(len(pool))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

type policyGradientGen struct{}

func (g *policyGradientGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*4)
	type entry struct{ q, a, e string }
	easy := []entry{{"REINFORCE with alpha=0.1, return G=4, score grad 0.5: what is the update?", "0.2", "Update alpha*G*grad = 0.1*4*0.5 = 0.2."}, {"Q=7, V=5. What is the advantage A=Q-V?", "2", "The action beats its state average by 7-5 = 2."}, {"A variance-cutting signal that adds no bias is an example of what: advantage or entropy?", "advantage", "Baselines like advantage center updates without bias."}, {"Policy gradient samples from the current or an old policy: current or old?", "current", "REINFORCE is on-policy: trajectories come from the current policy."}, {"Does REINFORCE use Monte Carlo returns? (yes/no)", "yes", "Full-episode returns weight each score gradient."}, {"Is the REINFORCE gradient estimator biased? (yes/no)", "no", "Unbiased but high-variance; baselines cut the variance."}}
	hard := []entry{{"Is actor-critic both value and policy? (yes/no)", "yes", "Actor-critic."}, {"Does entropy regularization encourage exploration? (yes/no)", "yes", "Entropy bonus."}, {"Is PPO clipped surrogate objective? (yes/no)", "yes", "PPO."}}
	pool := easy
	if scale > 3 {
		pool = append(easy, hard...)
	}
	e := pool[rand.Intn(len(pool))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

type bertGen struct{}

func (g *bertGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*4)
	type entry struct{ q, a, e string }
	easy := []entry{{"BERT base has 12 layers and large has 24. How many times deeper is large?", "2", "24/12 = 2: large doubles base depth."}, {"Masking 15 percent of 100 tokens masks how many?", "15", "0.15*100 = 15 masked tokens per batch."}, {"BERT base width 768 with 12 heads gives what head dimension?", "64", "768/12 = 64 dimensions per attention head."}, {"BERT pretraining pairs MLM with what second task: NSP or causal LM?", "NSP", "Next sentence prediction: does sentence B follow A."}, {"Is BERT a bidirectional encoder? (yes/no)", "yes", "Full left and right context informs every token."}, {"Is BERT trained from scratch for every downstream task? (yes/no)", "no", "Pretrain once, then fine-tune the same base per task."}}
	hard := []entry{{"Does BERT base have 12 layers? (yes/no)", "yes", "Base 12, large 24."}, {"Is next sentence prediction part of BERT pretraining? (yes/no)", "yes", "NSP."}, {"Does RoBERTa remove NSP? (yes/no)", "yes", "RoBERTa."}}
	pool := easy
	if scale > 3 {
		pool = append(easy, hard...)
	}
	e := pool[rand.Intn(len(pool))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

type actorCriticGen struct{}

func (g *actorCriticGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*4)
	type entry struct{ q, a, e string }
	easy := []entry{{"r=1, gamma=0.9, next value 10, current value 8: what is the TD error?", "2", "Target 1+0.9*10 = 10; error 10-8 = 2."}, {"Q=7, V=5. What is the advantage?", "2", "Advantage Q-V = 7-5 = 2."}, {"The critic fits values while the actor updates what: policy or value?", "policy", "Actor holds the policy, critic the value baseline."}, {"A2C runs actors synchronously or asynchronously: synchronously or asynchronously?", "synchronously", "A2C is synchronous; A3C is the asynchronous variant."}, {"Does the critic provide a value baseline? (yes/no)", "yes", "V(s) centers updates as the advantage baseline."}, {"Is the actor updated with raw returns and no baseline? (yes/no)", "no", "Actor-critic weights updates by advantage, not raw return."}}
	hard := []entry{{"Is advantage = Q - V? (yes/no)", "yes", "Advantage."}, {"Does GAE generalize advantage? (yes/no)", "yes", "Generalized."}, {"Is DDPG actor-critic for continuous actions? (yes/no)", "yes", "DDPG."}}
	pool := easy
	if scale > 3 {
		pool = append(easy, hard...)
	}
	e := pool[rand.Intn(len(pool))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}
