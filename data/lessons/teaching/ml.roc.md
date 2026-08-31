# ROC Curves and Evaluation Metrics

**ROC curve:** Plots true positive rate $\text{TPR}=TP/(TP+FN)$ vs false positive rate $\text{FPR}=FP/(FP+TN)$ as classification threshold varies.

## Threshold Trade-Offs

### Area Under Curve
$\text{AUC}\in[0,1]$: $1$ perfect, $0.5$ random diagonal, $<0.5$ worse than random. $\text{AUC}=P(\text{score(pos)}>\text{score(neg)})$.

### Imbalanced Data
ROC can look optimistic when negatives dominate; precision-recall ($\text{Precision}=TP/(TP+FP)$ vs $\text{Recall}=\text{TPR}$) often more informative for rare positives.

## Example

Classifier on 90 negatives +10 positives: threshold sweep yields points $(0,0)→(1,1)$. AUC $0.85$ means 85% chance a random positive scores higher than a random negative.
