# Momentum and Adam Optimization

**Plain gradient descent:** $w_{t+1}=w_t-\eta g_t$. It stalls in flat regions and zigzags in narrow valleys.

**Momentum:** Keep a velocity $v$ that accumulates past gradients,

$$v_{t+1}=\beta v_t + g_t, \quad w_{t+1}=w_t-\eta v_{t+1}$$

**Adam:** Combines momentum with adaptive per-parameter scaling from squared gradients,

$$m_{t+1}=\beta_1 m_t+(1-\beta_1)g_t,\; v_{t+1}=\beta_2 v_t+(1-\beta_2)g_t^2,\; w_{t+1}=w_t-\eta \frac{m_{t+1}}{\sqrt{v_{t+1}}+\epsilon}$$

## Accelerating Through Terrain

### Momentum's Memory
$\beta\in[0,1)$ controls memory: $\beta=0.9$ gives horizon $\approx10$ steps, carrying the optimizer through flat spots and damping oscillations.

### Adam's Adaptation
Each weight gets its own effective learning rate $\eta/\sqrt{v}$, so sparse or steep directions are normalized, speeding convergence without manual scheduling.

## Example

Loss valley: gradient alternates $+5,-5$ across walls with $+0.3$ downhill. Plain SGD zigzags; momentum averages $\pm5$ to near zero on walls while $0.3$ accumulates, moving steadily downhill. Adam rescales walls further for faster descent.
