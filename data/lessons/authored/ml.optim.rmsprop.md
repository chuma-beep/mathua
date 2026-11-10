# RMSProp optimization

## The decayed average

RMSProp keeps a decayed average of squared gradients. The average gives more weight to recent gradients and less to old ones. The decay rate controls how quickly the memory fades.

The square root of the average scales the gradient. A large average produces a smaller update, which stabilizes training.

## The update

The update is the learning rate times the gradient divided by the square root of the decayed average. The parameter moves opposite the gradient.

A small constant is often added inside the square root to avoid division by zero. The structure of the update is the same with or without it.

## Comparing with AdaGrad

AdaGrad accumulates every squared gradient, so its denominator grows without bound. RMSProp forgets old gradients, so its denominator can shrink and its updates can remain useful.

The decay rate is the main difference. A decay near one remembers longer, while a smaller decay reacts faster to recent gradients.

A final check is to restate the answer in the units and context of the original question. If the result is plausible and the units match, the calculation is ready.
