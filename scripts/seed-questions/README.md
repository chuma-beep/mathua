# Question Seed Script

Populates the mathua question bank with practice questions for each concept.

## Usage

### Generate questions using built-in generators

```bash
# From project root — seeds 8 questions per concept into mathua.db
go run scripts/seed-questions/main.go -db mathua.db -n 8
```

### Import questions from a JSON file

```bash
go run scripts/seed-questions/main.go -db mathua.db -import questions.json
```

## JSON Import Format

```json
[
  {
    "concept_id": "alg.quad.formula",
    "question": "Solve x^2 - 5x + 6 = 0",
    "answer": "x = 2, x = 3",
    "explanation": "Factor: (x-2)(x-3) = 0, so x = 2 or x = 3",
    "source": "math_dataset",
    "difficulty": 0.5
  }
]
```

Fields:
- `concept_id` (required) — must match a mathua concept ID (e.g. `alg.quad.formula`)
- `question` (required) — the problem statement (supports LaTeX)
- `answer` (required) — the final answer
- `explanation` (optional) — step-by-step solution
- `source` (optional) — name of the source dataset
- `difficulty` (optional, default 0.5) — 0.0 to 1.0

## Suggested External Datasets

| Dataset | License | Size | How to convert |
|---------|---------|------|----------------|
| [MATH (hendrycks)](https://github.com/hendrycks/math) | MIT | 12,500 | Use `scripts/convert-math.py` |
| [StackMathQA](https://huggingface.co/datasets/math-ai/StackMathQA) | CC BY-SA | 2M | Use HuggingFace `datasets` library |
| [OpenStax](https://openstax.org/) | CC BY 4.0 | 1000s per book | Scrape exercises from book JSON |

Write a Python script that loads any of these, maps the subject tags to mathua concept IDs, and exports in the JSON format above.
