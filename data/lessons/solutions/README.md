# Solution schemas

One file per concept: `<concept-id>.json`. Each renders the explanation for
**every** problem a concept generates, so the learner sees the same reasoning
shape with different numbers each time.

## Why a schema instead of prose

The arithmetic stays in the generator, which is the single source of truth for
the answer. A schema names values the generator already computed; it never
recomputes them, so an explanation cannot contradict the answer it explains.

- **Numbers** come from the generator (`Problem.Facts`).
- **Prose** comes from here.

## Format

```json
{
  "concept": "frac.add.diff",
  "setup": "You are adding $[[a_num]]/[[a_den]]$ and $[[b_num]]/[[b_den]]$.",
  "steps": [
    {
      "fact": "lcm",
      "say": "The bottoms differ, so put both over the smallest common multiple. LCM([[a_den]], [[b_den]]) = [[lcm]]."
    },
    {
      "fact": "b_scaled",
      "say": "That scales the numerators too: $[[b_num]]/[[b_den]]$ becomes $[[b_scaled]]/[[lcm]]$."
    }
  ],
  "answer": "Adding the tops gives [[num]]/[[lcm]]."
}
```

| Field | Required | Meaning |
|---|---|---|
| `concept` | no | Concept id. Defaults to the filename. |
| `setup` | no | Frames the problem, with values inlined. |
| `steps[].fact` | yes | The named fact the step rests on. If the generator does not publish it, the step is skipped. |
| `steps[].say` | yes | The authored prose for that step. |
| `steps[].optional` | no | `true` when the step belongs to only one of the question shapes a multi-shape generator asks. Suppresses the drift report for a fact that is legitimately absent on the other shape. |
| `answer` | **yes** | Closes on the final result. A schema with no `answer` is rejected at load. |

## Placeholders

`[[name]]` — double square brackets. Single braces are LaTeX and the corpus is
full of them, so they cannot be the delimiter.

Available in every schema:

| Name | Value |
|---|---|
| `[[answer]]` | The graded answer. Present so prose never has to restate it. |

Everything else is a fact the concept's generator publishes. See each
generator's `Problem{...}` literal for the names it sets.

## One generator can ask two questions

Some generators draw between a *recognition* question and a *production* one —
classify this angle, or find its complement. One schema covers both, with each
question's steps marked `"optional": true` so the step that does not apply is
skipped instead of dragging the whole schema down. Leave the flag off a step
whose absence would mean you and the generator have drifted apart: that case
drops the schema, which is the point.

## Failure is always a fallback, never a guess

A schema is used only if it renders completely. If the concept has no schema,
or any placeholder has no value behind it, the whole schema is skipped and the
generator's own `Explanation` is served instead. The skip is logged. A learner
is never shown a half-filled explanation.

## Checking a schema

`go test ./internal/solutions/...` renders every shipped schema against a
real generated problem and fails on any unresolved placeholder, naming the
concept and the fact. Add the schema and the generator's `Facts` in the same
commit, or that test is the thing that will tell you they disagree.
