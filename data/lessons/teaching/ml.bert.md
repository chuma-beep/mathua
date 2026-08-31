# BERT and Pretraining

**BERT:** Bidirectional encoder, masked language modeling (MLM) + next sentence prediction (NSP) pretraining, then fine-tuned; base $12$ layers, large $24$; RoBERTa removes NSP.

## Pretraining

### MLM
Mask $15\%$ tokens, predict; NSP predicts if next sentence follows.

### Fine-Tuning
Pretrained then task-specific fine-tuning.

## Example

BERT base: MLM on Wikipedia/BookCorpus, then fine-tune on SQuAD.
