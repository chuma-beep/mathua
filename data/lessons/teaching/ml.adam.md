# Adam Optimizer Details

**Adam:** Maintains first moment $m_t=\beta_1 m_{t-1}+(1-\beta_1)g_t$ and second moment $v_t=\beta_2 v_{t-1}+(1-\beta_2)g_t^2$ with bias correction $\hat{m}=m/(1-\beta_1^t)$, $\hat{v}=v/(1-\beta_2^t)$, update $w←w-η\hat{m}/(\sqrt{\hat{v}}+ε)$.

## Adaptive Per-Parameter

### Defaults
$\beta_1=0.9$, $\beta_2=0.999$, $ε=10^{-8}$; combines momentum and RMSProp, adaptive per-weight learning rate.

### Behaviour
Early steps bias-corrected; later $v$ stabilizes. Works well without manual scheduling.

## Example

Gradient $g=[2,0.1]$: $m$ smooths direction, $v$ scales $2$ down, $0.1$ up, so both dimensions progress steadily.
