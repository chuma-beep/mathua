# BERT and Pretraining

**BERT:** A bidirectional transformer encoder pretrained with masked language modeling (MLM) plus next-sentence prediction (NSP), then fine-tuned per task. Base has 12 layers, large has 24; RoBERTa later dropped NSP.

## Worked: masking one batch

Take 100 tokens with the standard 15 percent mask rate:
1. Choose $0.15\cdot 100 = 15$ tokens to mask.
2. Of the 15, roughly 12 become [MASK], 2 become random tokens, 1 stays unchanged — the mix stops the model gaming the mask token.
3. The loss predicts only those 15 originals: bidirectional context (both sides) informs every prediction.

So MLM turns raw text into supervision: hide tokens, predict them from full surroundings — no labels needed.

## Worked: pretrain then fine-tune

Take BERT base pretrained on Wikipedia and BookCorpus:
1. Pretraining learns general language: syntax, facts, word relations — expensive, done once.
2. Fine-tuning on SQuAD swaps in one question-answering head and trains all weights a few more epochs — cheap, per task.
3. The same base adapts to classification, NER, or entailment by changing only the head and the data.

So pretraining amortizes language knowledge across tasks: one big unsupervised run, many small supervised adaptations.

## When NSP goes away

Ablations showed NSP barely helped: RoBERTa removed it, trained longer on more data with dynamic masking, and scored higher. Pretrain with MLM at scale; treat auxiliary objectives as guilty until proven useful.
