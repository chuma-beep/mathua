package machinelearning

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
)

// A count question asks how many of eight statements are true. Each pool holds
// sixteen statements with eight true and eight false cases, so the answer can
// be drawn uniformly from 0 to 8 and every statement is explained afterwards.
type mlVerdict struct {
	label string
	yes   bool
	note  string
}

const mlDraw = 8

func countMLVerdicts(lead string, pool []mlVerdict) generator.Problem {
	var yes, no []mlVerdict
	for _, v := range pool {
		if v.yes {
			yes = append(yes, v)
		} else {
			no = append(no, v)
		}
	}
	count := rand.Intn(mlDraw + 1)
	if count > len(yes) {
		count = len(yes)
	}
	if mlDraw-count > len(no) {
		count = mlDraw - len(no)
	}
	rand.Shuffle(len(yes), func(i, j int) { yes[i], yes[j] = yes[j], yes[i] })
	rand.Shuffle(len(no), func(i, j int) { no[i], no[j] = no[j], no[i] })
	chosen := append([]mlVerdict{}, yes[:count]...)
	chosen = append(chosen, no[:mlDraw-count]...)
	rand.Shuffle(len(chosen), func(i, j int) { chosen[i], chosen[j] = chosen[j], chosen[i] })

	answer := strconv.Itoa(count)
	question := lead
	for attempt := 0; attempt < 40; attempt++ {
		var parts []string
		for i, v := range chosen {
			parts = append(parts, fmt.Sprintf("(%c) %s", 'a'+i, v.label))
		}
		question = fmt.Sprintf("%s %s (enter a number)", lead, strings.Join(parts, " "))
		if !strings.Contains(question, answer) {
			break
		}
		rand.Shuffle(len(chosen), func(i, j int) { chosen[i], chosen[j] = chosen[j], chosen[i] })
	}

	var notes []string
	for _, v := range chosen {
		word := "false"
		if v.yes {
			word = "true"
		}
		notes = append(notes, fmt.Sprintf("%s: %s, %s", v.label, word, v.note))
	}
	return generator.Problem{
		Question:    question,
		Answer:      answer,
		Explanation: fmt.Sprintf("Exactly %d of the 8 statements are true. %s", count, strings.Join(notes, " ")),
		Facts:       map[string]string{"count": answer},
	}
}

func crossEntropyProblem() generator.Problem {
	return countMLVerdicts("How many of these statements about cross-entropy are true?", []mlVerdict{
		{"For a one-hot label the loss is minus the log of the predicted probability of the true class", true, "the surviving term is -log p"},
		{"Cross-entropy is non-negative", true, "probabilities lie between 0 and 1"},
		{"For a Bernoulli label it equals the negative log-likelihood", true, "the two expressions are the same"},
		{"With softmax its gradient in the logits is predicted minus target", true, "the softmax Jacobian cancels the log"},
		{"It penalises a confident wrong prediction very strongly", true, "the loss grows without bound as p tends to 0"},
		{"Binary cross-entropy uses y log p plus (1-y) log(1-p)", true, "the label selects one term"},
		{"For softmax logits the loss is convex", true, "the log-sum-exp form is convex"},
		{"Cross-entropy equals KL divergence plus the entropy of the label distribution", true, "a standard decomposition"},
		{"Cross-entropy is bounded above by 1", false, "it is unbounded as the true-class probability tends to 0"},
		{"It ignores which class is the true one", false, "only the true-class probability enters"},
		{"Mean squared error is always the better classification loss", false, "cross-entropy matches the probabilistic output"},
		{"Its gradient in the logits is predicted plus target", false, "the gradient is predicted minus target"},
		{"It applies only to regression", false, "it is the standard classification loss"},
		{"Softmax cross-entropy is not differentiable", false, "it is differentiable everywhere"},
		{"The loss falls when the true-class probability falls", false, "it rises in that case"},
		{"Cross-entropy is the same quantity as accuracy", false, "accuracy counts correct labels, the loss scores probabilities"},
	})
}

func overfittingProblem() generator.Problem {
	return countMLVerdicts("How many of these statements about overfitting are true?", []mlVerdict{
		{"It shows as low training loss and high validation loss", true, "the gap is the signature"},
		{"More training data usually reduces it", true, "more evidence constrains the fit"},
		{"Regularisation usually reduces it", true, "a penalty restricts the hypothesis"},
		{"A simpler model usually reduces it", true, "lower capacity fits noise less"},
		{"Early stopping can reduce it", true, "training stops before noise is fitted"},
		{"It means the model has memorised noise in the training set", true, "noise cannot generalise"},
		{"A large train-validation gap points to variance", true, "the model is unstable across samples"},
		{"Cross-validation can detect it", true, "held-out folds expose the gap"},
		{"It is always a sign of a healthy model", false, "it is a failure to generalise"},
		{"Validation loss is normally lower than training loss", false, "the reverse gap is the concern"},
		{"Adding parameters always improves generalisation", false, "capacity can fit noise"},
		{"Overfitting is the same as high bias", false, "it is high variance"},
		{"Tuning on the test set is an acceptable remedy", false, "that leaks the test set"},
		{"Dropout makes generalisation worse", false, "dropout is a regulariser"},
		{"High training accuracy alone proves generalisation", false, "held-out performance is the test"},
		{"It cannot occur with a small model", false, "even small models can overfit small data"},
	})
}

func dropoutProblem() generator.Problem {
	return countMLVerdicts("How many of these statements about dropout are true?", []mlVerdict{
		{"It randomly zeroes activations during training", true, "the mask is sampled per pass"},
		{"Inverted dropout scales survivors by 1/(1-p)", true, "the expected activation is preserved"},
		{"It is disabled at inference", true, "prediction uses the full network"},
		{"It acts as a regulariser", true, "it discourages reliance on single units"},
		{"It reduces co-adaptation of neurons", true, "units cannot rely on fixed partners"},
		{"It approximates an ensemble of subnetworks", true, "each mask is a shared-weight subnetwork"},
		{"Here p is the probability of dropping a unit", true, "p is the drop rate"},
		{"It is applied per activation rather than per layer", true, "each unit gets an independent mask"},
		{"It is applied at test time", false, "the mask is off at inference"},
		{"It increases co-adaptation", false, "it does the opposite"},
		{"It is deterministic", false, "masks are random"},
		{"Here p is the probability of keeping a unit", false, "the keep probability is 1-p"},
		{"It replaces batch normalisation", false, "the two mechanisms are different"},
		{"It is used only on the output layer", false, "it is normally used on hidden layers"},
		{"It always lowers training loss", false, "training loss usually rises while validation improves"},
		{"It has no effect on variance", false, "it is a variance-reducing regulariser"},
	})
}

func rnnProblem() generator.Problem {
	return countMLVerdicts("How many of these statements about recurrent networks are true?", []mlVerdict{
		{"An RNN carries a hidden state across time steps", true, "the state is the memory"},
		{"It shares weights across time steps", true, "the same recurrent matrix is reused"},
		{"Vanilla RNNs can suffer from vanishing gradients", true, "repeated multiplication shrinks them"},
		{"LSTMs use input, forget and output gates", true, "gates control the cell state"},
		{"Unfolding an RNN gives a deep feedforward computation", true, "each time step becomes a layer"},
		{"It can process sequences of variable length", true, "the recurrence is applied step by step"},
		{"Training uses backpropagation through time", true, "the unfolded graph is differentiated"},
		{"The order of the input sequence matters", true, "the state depends on the order"},
		{"It has no weight sharing across time", false, "weights are shared"},
		{"It is a purely feedforward network", false, "the recurrence is the defining feature"},
		{"Its gradients are always stable", false, "they can vanish or explode"},
		{"LSTMs have no gates", false, "they have three standard gates"},
		{"It cannot process sequences", false, "sequences are its purpose"},
		{"Its hidden state is independent of the past", false, "the state summarises the past"},
		{"It has no learnable parameters", false, "it has recurrent and input weights"},
		{"Its training order is irrelevant to the loss", false, "the sequence order changes the computation"},
	})
}

func batchNormProblem() generator.Problem {
	return countMLVerdicts("How many of these statements about batch normalisation are true?", []mlVerdict{
		{"It normalises activations to mean 0 and variance 1", true, "the batch statistics standardise each feature"},
		{"It has learnable scale gamma and shift beta", true, "the affine step restores capacity"},
		{"Training uses statistics from the current batch", true, "the batch supplies the moments"},
		{"Inference uses running averages", true, "population statistics replace batch statistics"},
		{"It can allow a higher learning rate", true, "it stabilises the optimisation"},
		{"It is often applied before the activation", true, "a common ordering"},
		{"It reduces sensitivity to initialisation", true, "the normalisation controls scale"},
		{"Its behaviour depends on batch size", true, "small batches give noisy moments"},
		{"Training uses test-set statistics", false, "that would leak the test data"},
		{"It has no learnable parameters", false, "gamma and beta are learned"},
		{"It is applied only after softmax", false, "it is used throughout hidden layers"},
		{"It is always applied after the activation", false, "the usual order is before"},
		{"It is independent of the batch", false, "batch moments enter training"},
		{"It replaces dropout", false, "the two are separate regularisers"},
		{"It cannot be used in convolutional networks", false, "it is standard in CNNs"},
		{"It removes the need for any activation function", false, "normalisation is not an activation"},
	})
}

func rocProblem() generator.Problem {
	return countMLVerdicts("How many of these statements about ROC curves are true?", []mlVerdict{
		{"The axes are true positive rate and false positive rate", true, "TPR against FPR"},
		{"An AUC of 0.5 corresponds to random ranking", true, "the diagonal is chance"},
		{"An AUC of 1 means perfect separation", true, "every positive outranks every negative"},
		{"The curve sweeps the classification threshold", true, "each threshold gives one point"},
		{"AUC is invariant to class balance", true, "it ranks within classes"},
		{"AUC equals the probability that a positive scores above a negative", true, "the ranking interpretation"},
		{"The diagonal is the chance baseline", true, "no discrimination"},
		{"Precision-recall is often more informative under imbalance", true, "PR focuses on the positive class"},
		{"ROC plots precision against recall", false, "that is the PR curve"},
		{"An AUC of 0 is the worst possible classifier", false, "0 is perfect reversed ranking"},
		{"The AUC depends on the chosen threshold", false, "it integrates over all thresholds"},
		{"Higher AUC always means better calibration", false, "AUC measures ranking, not calibration"},
		{"ROC requires balanced classes", false, "it is defined for any class balance"},
		{"AUC is the same as accuracy", false, "accuracy uses one threshold"},
		{"ROC plots the loss against training time", false, "it plots rates"},
		{"The threshold is fixed at 0.5 while drawing the curve", false, "the threshold varies"},
	})
}

func transformerProblem() generator.Problem {
	return countMLVerdicts("How many of these statements about transformers are true?", []mlVerdict{
		{"Self-attention compares every token with every other token", true, "all query-key pairs are scored"},
		{"With n tokens the attention score matrix is n by n", true, "one row per query"},
		{"Multi-head attention splits the model width across heads", true, "each head has width d_model/h"},
		{"Positional encoding supplies order information", true, "attention alone is permutation equivariant"},
		{"Attention costs O(n^2) in sequence length", true, "the score matrix is quadratic"},
		{"Layer normalisation is used in the block", true, "it stabilises training"},
		{"Residual connections are used in the block", true, "they carry information across layers"},
		{"A decoder uses a causal mask", true, "it prevents looking ahead"},
		{"It is built on a recurrent state", false, "there is no recurrence"},
		{"Attention is linear in sequence length", false, "the score matrix is quadratic"},
		{"No positional information is needed", false, "order must be injected"},
		{"It uses only one attention head", false, "multi-head attention is standard"},
		{"Attention weights are left unnormalised", false, "softmax normalises them"},
		{"A decoder may attend to future tokens", false, "the causal mask blocks them"},
		{"It has no residual connections", false, "residuals are standard"},
		{"It only accepts a fixed sequence length", false, "variable lengths are padded or packed"},
	})
}

func ganProblem() generator.Problem {
	return countMLVerdicts("How many of these statements about GANs are true?", []mlVerdict{
		{"A GAN trains a generator and a discriminator", true, "two competing networks"},
		{"Training is a minimax game", true, "each player optimises against the other"},
		{"At the Nash equilibrium the discriminator outputs 0.5", true, "real and fake are indistinguishable"},
		{"Mode collapse means the generator covers few modes", true, "diversity is lost"},
		{"Wasserstein GAN uses an earth-mover distance", true, "the critic estimates it"},
		{"The discriminator supplies a gradient to the generator", true, "the generator learns through D"},
		{"Training is adversarial", true, "the two losses oppose each other"},
		{"The discriminator classifies real against fake", true, "that is its task"},
		{"It trains a single network", false, "there are two networks"},
		{"The generator outputs class labels", false, "it outputs samples"},
		{"Training is always stable", false, "instability is common"},
		{"Mode collapse is a sign of healthy training", false, "it is a failure mode"},
		{"The discriminator is never trained", false, "it is trained alongside G"},
		{"A GAN gives an explicit likelihood", false, "it is implicit"},
		{"The generator has no latent input", false, "it maps noise to samples"},
		{"The discriminator outputs one of ten classes", false, "it outputs a real-fake score"},
	})
}

func adamProblem() generator.Problem {
	return countMLVerdicts("How many of these statements about Adam are true?", []mlVerdict{
		{"It maintains first and second moment estimates", true, "mean and squared gradient"},
		{"It applies bias correction to both moments", true, "the initial zeros are corrected"},
		{"The defaults are beta1 0.9 and beta2 0.999", true, "standard settings"},
		{"It adapts the step size per parameter", true, "each weight has its own rate"},
		{"It combines momentum with RMSProp scaling", true, "a hybrid method"},
		{"The update is eta times m-hat over the square root of v-hat plus epsilon", true, "the Adam rule"},
		{"It is an adaptive optimiser", true, "the effective rate varies"},
		{"With beta 0.9 the memory horizon is about 10 steps", true, "1/(1-beta)"},
		{"It uses one global step size", false, "the rate is per parameter"},
		{"It skips bias correction", false, "bias correction is central"},
		{"It tracks only the first moment", false, "it tracks both moments"},
		{"The default beta1 is 0.999", false, "that is beta2"},
		{"It requires a convex loss", false, "it is used for nonconvex networks"},
		{"It is exactly stochastic gradient descent", false, "it is adaptive"},
		{"The second moment is the raw gradient", false, "it is the squared gradient average"},
		{"It has no epsilon term", false, "epsilon stabilises the division"},
	})
}

func augmentationProblem() generator.Problem {
	return countMLVerdicts("How many of these statements about data augmentation are true?", []mlVerdict{
		{"It increases the variety of the training data", true, "transforms create new views"},
		{"It usually reduces overfitting", true, "more variety constrains the fit"},
		{"It is applied during training", true, "test data stays clean"},
		{"Mixup and cutout are augmentation methods", true, "both create altered samples"},
		{"It can approximate invariance to a transform", true, "the model sees transformed copies"},
		{"It acts as a regulariser", true, "it restricts what the model can rely on"},
		{"It uses label-preserving transforms", true, "the label must remain correct"},
		{"It can improve robustness", true, "more varied inputs are covered"},
		{"It is applied at test time", false, "that would change the evaluation"},
		{"It replaces the model", false, "it changes the data"},
		{"It guarantees higher accuracy", false, "no method guarantees it"},
		{"It increases overfitting", false, "it normally reduces it"},
		{"It only ever uses horizontal flips", false, "many transforms exist"},
		{"It may change labels arbitrarily", false, "transforms must preserve the label"},
		{"It has no regularisation effect", false, "it is a regulariser"},
		{"It makes the dataset smaller", false, "it expands the effective data"},
	})
}
