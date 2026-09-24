# Generative Adversarial Networks

**GAN:** A generator $G(z)$ maps noise $z$ to data space while a discriminator $D(x)$ scores real versus fake; training is minimax, with $D$ maximizing $\log D(x) + \log(1-D(G(z)))$ and $G$ minimizing the second term. The discriminator's gradient teaches the generator.

## Worked: one discriminator update

Take one real sample with $D = 0.9$ and one fake with $D = 0.2$:
1. Real term: $\log 0.9 \approx -0.11$ — nearly perfect, small push to raise $D$ on real data.
2. Fake term: $\log(1-0.2) = \log 0.8 \approx -0.22$ — the fake fools somewhat, push to lower $D$ on fakes.
3. Total objective $\approx -0.33$; ascending it sharpens $D$, and the same backward pass hands $G$ gradients pointing toward more convincing fakes.

So the two networks share one backward pass: $D$ learns to discriminate, $G$ learns through $D$'s eyes.

## Worked: equilibrium

Take training converged with $G$ matching the data distribution exactly:
1. Real and fake samples are indistinguishable, so the best $D$ outputs 0.5 everywhere.
2. At $D = 0.5$ the objective is $\log 0.5 + \log 0.5 \approx -1.39$, and neither player can improve alone — a Nash equilibrium.
3. In practice training hovers near, not at, this point: gradients shrink as $D$ approaches chance.

So the target is precise: generator distribution equal to data, discriminator at chance — any deviation from 0.5 says who is winning.

## When the generator collapses

If $G$ discovers 2 of 10 data modes fool $D$ reliably, it may emit only those: mode collapse, with coverage 0.2 and no pressure to diversify. Wasserstein GANs replace the classifier loss with Earth mover distance for smoother gradients — change the game when the players stall.
