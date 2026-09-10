# Continuous maps between spaces

Continuity is the bridge that lets one metric space talk to another.

## Epsilon-Delta Definition

A function \\(f: X \to Y\\) between metric spaces is **continuous at a point** \\(x_0 \in X\\) when small moves in \\(X\\) produce small moves in \\(Y\\): for every tolerance \\(\varepsilon > 0\\) there is a \\(\delta > 0\\) such that whenever \\(d_X(x, x_0) < \delta\\), we get \\(d_Y(f(x), f(x_0)) < \varepsilon\\). The function is **continuous** when this holds at every point.

## Open Set Definition

This is the same \\(\varepsilon\\)-\\(\delta\\) idea from real analysis — the output can be pinned down to any accuracy by choosing the input accurately enough. What makes topology powerful is a purely set-theoretic reformulation: \\(f\\) is continuous exactly when the *preimage* of every open set in \\(Y\\) is open in \\(X\\). No distances are needed to check it, only which regions count as open.

## Equivalent Conditions

From that single characterization the classic theorems fall out cheaply. Compositions of continuous maps are continuous. Continuous images of connected sets stay connected, and of compact sets stay compact. And on compact domains, continuous real-valued functions are bounded and attain their extrema — the engine behind every existence theorem you will meet in optimization.
