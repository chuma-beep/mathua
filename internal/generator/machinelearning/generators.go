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
			"output (or prediction)",
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
			"In a neural network with sigmoid activation, what is the derivative of \\(\\sigma(x) = 1/(1+e^{-x})\\) in terms of \\(\\sigma(x)\\)?",
			"σ(x)(1-σ(x))",
			"The derivative of the sigmoid function is \\(\\sigma(x)(1-\\sigma(x))\\), which makes backpropagation efficient.",
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
			"During backpropagation, the gradient of the loss with respect to a weight depends on the gradient from the ____ layer.",
			"next (or subsequent)",
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
	type entry struct{question, answer, exp string}
	entries:=[]entry{{"Cross-entropy loss between true label 1 and prediction 0.8 is (-log(0.8)≈?) (0.22/1.20)","0.22","CE = -log(0.8) ≈ 0.22."},{"Cross-entropy is appropriate for ____ tasks.","classification","Cross-entropy for classification."},{"Binary cross-entropy for target 0 and prediction 0.1 is (-log(0.9)≈?) (0.10/2.30)","0.10","CE = -log(0.9) ≈ 0.10."},{"Does cross-entropy penalize confident wrong predictions more than MSE? (yes/no)","yes","CE grows ∞ as p→0."}}
	e:=entries[rand.Intn(len(entries))]
	return generator.Problem{Question:e.question,Answer:e.answer,Explanation:e.exp}
}
type regularizationGen struct{}
func (g *regularizationGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct{question, answer, exp string}
	entries:=[]entry{{"L2 regularization adds λ||w||² to the loss. What does λ control? (regularization strength)","regularization strength","λ trades off."},{"Does L2 regularization shrink weights toward zero? (yes/no)","yes","Gradient 2λw pulls."},{"What is the gradient of λ||w||² w.r.t w? (2λw)","2λw","d/dw λw² = 2λw."},{"Is regularization used to reduce overfitting? (yes/no)","yes","Penalizing large weights."}}
	e:=entries[rand.Intn(len(entries))]
	return generator.Problem{Question:e.question,Answer:e.answer,Explanation:e.exp}
}
type overfittingGen struct{}
func (g *overfittingGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct{question, answer, exp string}
	entries:=[]entry{{"Overfitting is when training error is ____ than test error. (lower/higher)","lower","Training lower."},{"What curve shape indicates overfitting when plotting error vs model capacity?","U-shaped test error","Test drops then rises."},{"Does more training data generally reduce overfitting? (yes/no)","yes","More samples."},{"Which error measures generalization: training or test error?","test","Test error."}}
	e:=entries[rand.Intn(len(entries))]
	return generator.Problem{Question:e.question,Answer:e.answer,Explanation:e.exp}
}
type validationGen struct{}
func (g *validationGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct{question, answer, exp string}
	entries:=[]entry{{"What is the purpose of a validation set?","tune hyperparameters (or model selection)","Validation for hyperparams."},{"k-fold cross-validation with k=5 splits data into how many folds?","5","k folds."},{"Should the test set be used for hyperparameter tuning? (yes/no)","no","Leakage."},{"Does cross-validation give a more reliable estimate than a single split? (yes/no)","yes","Averaging."}}
	e:=entries[rand.Intn(len(entries))]
	return generator.Problem{Question:e.question,Answer:e.answer,Explanation:e.exp}
}
type momentumGen struct{}
func (g *momentumGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct{question, answer, exp string}
	entries:=[]entry{{"Momentum update: v_{t+1}=β v_t + g_t. What does β control? (momentum / decay)","momentum","β controls persistence."},{"Adam combines momentum and ____ scaling.","adaptive (or RMSProp)","Adaptive."},{"Does momentum help escape shallow local minima? (yes/no)","yes","Velocity carries."},{"With β=0.9, the effective memory horizon is roughly how many steps? (10)","10","Horizon ≈1/(1-β)."}}
	e:=entries[rand.Intn(len(entries))]
	return generator.Problem{Question:e.question,Answer:e.answer,Explanation:e.exp}
}
type dropoutGen struct{}
func (g *dropoutGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale:=int(1+ctx.Difficulty*4)
	type entry struct{question, answer, exp string}
	easy:=[]entry{{"Dropout randomly zeroes activations with probability p during training. At test time, are activations scaled by (1-p) or kept as-is? (1-p)","1-p","Inverted dropout."},{"Does dropout reduce co-adaptation of neurons? (yes/no)","yes","Robustness."},{"Is dropout only applied during training? (yes/no)","yes","Disabled at inference."}}
	hard:=[]entry{{"With p=0.5, what is the expected number of active neurons out of 100? (50)","50","Expected (1-p)*100."},{"Does dropout approximate an ensemble of 2^n subnetworks? (yes/no)","yes","Each mask subnetwork."},{"Is dropout a form of regularization? (yes/no)","yes","Generalization."}}
	pool:=easy
	if scale>3{pool=append(easy,hard...)}
	e:=pool[rand.Intn(len(pool))]
	return generator.Problem{Question:e.question,Answer:e.answer,Explanation:e.exp}
}
type cnnGen struct{}
func (g *cnnGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale:=int(1+ctx.Difficulty*4)
	type entry struct{question, answer, exp string}
	easy:=[]entry{{"CNNs use ____ layers to detect local patterns.","convolutional","Convolutions."},{"Does weight sharing in CNNs reduce parameters vs fully connected? (yes/no)","yes","Same filter reused."},{"What operation follows convolution in a CNN block? (activation + pooling)","activation","After convolution activation."}}
	hard:=[]entry{{"For a 5x5 input, 3x3 filter, stride 1, no padding, what is output size? (3x3)","3x3","(5-3)/1+1=3."},{"Does max-pooling provide translation invariance? (yes/no)","yes","Robust to shifts."},{"Are CNNs especially suited for image data? (yes/no)","yes","Local connectivity."}}
	pool:=easy
	if scale>3{pool=append(easy,hard...)}
	e:=pool[rand.Intn(len(pool))]
	return generator.Problem{Question:e.question,Answer:e.answer,Explanation:e.exp}
}
type rnnGen struct{}
func (g *rnnGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale:=int(1+ctx.Difficulty*4)
	type entry struct{question, answer, exp string}
	easy:=[]entry{{"RNNs maintain a ____ state across time steps.","hidden","Hidden state."},{"Does an RNN share weights across time steps? (yes/no)","yes","Same recurrent weights."},{"What problem plagues vanilla RNNs on long sequences? (vanishing gradient)","vanishing gradient","Backprop multiplies small gradients."}}
	hard:=[]entry{{"Do LSTMs use gates (input, forget, output) to combat vanishing gradients? (yes/no)","yes","Gating."},{"Is an RNN with hidden size 0 able to model sequences? (yes/no)","no","No capacity."},{"Does unfolding an RNN yield a deep feedforward network? (yes/no)","yes","Unrolling creates depth."}}
	pool:=easy
	if scale>3{pool=append(easy,hard...)}
	e:=pool[rand.Intn(len(pool))]
	return generator.Problem{Question:e.question,Answer:e.answer,Explanation:e.exp}
}
type batchNormGen struct{}
func (g *batchNormGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale:=int(1+ctx.Difficulty*4)
	type entry struct{question, answer, exp string}
	easy:=[]entry{{"Batch norm normalizes activations to have mean 0 and variance ____ (1)","1","Normalized."},{"Does batch norm have learnable scale γ and shift β? (yes/no)","yes","γ·x̂+β."},{"Is batch norm applied before or after activation? (before)","before","Typically before."}}
	hard:=[]entry{{"Does batch norm allow higher learning rates? (yes/no)","yes","Stabilized."},{"At test time, does batch norm use batch statistics or running averages? (running)","running","Population stats."},{"Does batch norm reduce internal covariate shift? (yes/no)","yes","Stated goal."}}
	pool:=easy
	if scale>3{pool=append(easy,hard...)}
	e:=pool[rand.Intn(len(pool))]
	return generator.Problem{Question:e.question,Answer:e.answer,Explanation:e.exp}
}
type rocGen struct{}
func (g *rocGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale:=int(1+ctx.Difficulty*4)
	type entry struct{question, answer, exp string}
	easy:=[]entry{{"ROC plots TPR vs ____ (FPR)","FPR","ROC."},{"Does AUC = 1 mean perfect classifier? (yes/no)","yes","Area 1."},{"Is AUC = 0.5 equivalent to random guessing? (yes/no)","yes","Diagonal."}}
	hard:=[]entry{{"For imbalanced data, is precision-recall often more informative than ROC? (yes/no)","yes","PR emphasizes."},{"Does higher AUC mean better ranking of positives above negatives? (yes/no)","yes","AUC = P(score(pos)>score(neg))."},{"If classifier thresholds from 0 to 1, does ROC go from (1,1) to (0,0) or opposite? (0,0→1,1)","0,0 to 1,1","Same curve."}}
	pool:=easy
	if scale>3{pool=append(easy,hard...)}
	e:=pool[rand.Intn(len(pool))]
	return generator.Problem{Question:e.question,Answer:e.answer,Explanation:e.exp}
}
type transformerGen struct{}
func (g *transformerGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale:=int(1+ctx.Difficulty*4)
	type entry struct{question, answer, exp string}
	easy:=[]entry{{"Does attention compute weighted sum of values by query-key similarity? (yes/no)","yes","Attention = softmax(QK^T/√d)V."},{"Is self-attention permutation equivariant without positional encoding? (yes/no)","yes","Positional encoding adds order."},{"Does transformer use multi-head attention? (yes/no)","yes","Multiple heads."}}
	hard:=[]entry{{"Is scaled dot-product attention O(n²) in sequence length? (yes/no)","yes","n×n matrix."},{"Does positional encoding use sin/cos of different frequencies? (yes/no)","yes","Sinusoidal."},{"Is transformer decoder autoregressive with causal mask? (yes/no)","yes","Mask prevents future."}}
	pool:=easy
	if scale>3{pool=append(easy,hard...)}
	e:=pool[rand.Intn(len(pool))]
	return generator.Problem{Question:e.question,Answer:e.answer,Explanation:e.exp}
}
type ganGen struct{}
func (g *ganGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale:=int(1+ctx.Difficulty*4)
	type entry struct{question, answer, exp string}
	easy:=[]entry{{"GAN has generator and discriminator playing minimax game? (yes/no)","yes","G vs D."},{"Does GAN training aim for Nash equilibrium? (yes/no)","yes","Equilibrium."},{"Is GAN generator trained to minimize log(1-D(G(z)))? (yes/no)","yes","Original loss."}}
	hard:=[]entry{{"Does mode collapse mean G outputs limited diversity? (yes/no)","yes","Few modes."},{"Is Wasserstein GAN using Earth mover distance? (yes/no)","yes","WGAN."},{"Does D provide gradient for G? (yes/no)","yes","Via backprop."}}
	pool:=easy
	if scale>3{pool=append(easy,hard...)}
	e:=pool[rand.Intn(len(pool))]
	return generator.Problem{Question:e.question,Answer:e.answer,Explanation:e.exp}
}
type adamDetailsGen struct{}
func (g *adamDetailsGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale:=int(1+ctx.Difficulty*4)
	type entry struct{question, answer, exp string}
	easy:=[]entry{{"Does Adam maintain first and second moment estimates? (yes/no)","yes","m and v."},{"Is bias correction important early in Adam? (yes/no)","yes","m̂=m/(1-β1^t)."},{"Does Adam default β1=0.9, β2=0.999? (yes/no)","yes","Defaults."}}
	hard:=[]entry{{"Does Adam update = -η m̂/(√v̂+ε)? (yes/no)","yes","Rule."},{"Is Adam adaptive per-parameter? (yes/no)","yes","Per-weight."},{"Does Adam combine momentum and RMSProp? (yes/no)","yes","Hybrid."}}
	pool:=easy
	if scale>3{pool=append(easy,hard...)}
	e:=pool[rand.Intn(len(pool))]
	return generator.Problem{Question:e.question,Answer:e.answer,Explanation:e.exp}
}
type calibrationGen struct{}
func (g *calibrationGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale:=int(1+ctx.Difficulty*4)
	type entry struct{question, answer, exp string}
	easy:=[]entry{{"Is calibrated classifier's confidence ≈ true accuracy? (yes/no)","yes","Reliability."},{"Does reliability diagram plot accuracy vs confidence? (yes/no)","yes","Binned."},{"Is ECE expected calibration error? (yes/no)","yes","Weighted avg."}}
	hard:=[]entry{{"Does overconfidence mean ECE high? (yes/no)","yes","Gap."},{"Is temperature scaling a calibration method? (yes/no)","yes","Softmax temp."},{"Does well-calibrated with 80% conf mean 80% correct? (yes/no)","yes","Definition."}}
	pool:=easy
	if scale>3{pool=append(easy,hard...)}
	e:=pool[rand.Intn(len(pool))]
	return generator.Problem{Question:e.question,Answer:e.answer,Explanation:e.exp}
}
type augmentationGen struct{}
func (g *augmentationGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale:=int(1+ctx.Difficulty*4)
	type entry struct{question, answer, exp string}
	easy:=[]entry{{"Is data augmentation used to increase training data artificially? (yes/no)","yes","Transforms."},{"Does augmentation reduce overfitting? (yes/no)","yes","More variety."},{"Is augmentation applied at test time? (no)","no","Test clean."}}
	hard:=[]entry{{"Is cutout/mixup a form of augmentation? (yes/no)","yes","Masks."},{"Does augmentation approximate invariance? (yes/no)","yes","Ignore transforms."},{"Is augmentation a form of regularization? (yes/no)","yes","Effective data."}}
	pool:=easy
	if scale>3{pool=append(easy,hard...)}
	e:=pool[rand.Intn(len(pool))]
	return generator.Problem{Question:e.question,Answer:e.answer,Explanation:e.exp}
}
