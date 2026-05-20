#!/usr/bin/env bash
set -euo pipefail

if [ $# -lt 2 ]; then
  echo "Usage: $0 <concept-id> <lesson-file-path>"
  echo ""
  echo "Scaffolds a new lesson markdown file and registers it in lessons.json."
  echo "  concept-id       e.g. 'alg.linear.slope'"
  echo "  lesson-file-path relative path under data/lessons/, e.g. 'algebra/slope.md'"
  echo ""
  echo "Example: $0 alg.linear.slope algebra/slope.md"
  exit 1
fi

CONCEPT_ID="$1"
LESSON_PATH="$2"
LESSONS_JSON="data/lessons/lessons.json"
LESSON_FILE="data/lessons/$LESSON_PATH"

mkdir -p "$(dirname "$LESSON_FILE")"

if [ ! -f "$LESSON_FILE" ]; then
  # Derive a title from the concept ID
  TITLE=$(echo "$CONCEPT_ID" | sed 's/\./ /g' | sed 's/\b\(.\)/\u\1/g')
  cat > "$LESSON_FILE" <<MDEOF
# $TITLE

TODO: Write a lesson about $CONCEPT_ID.

## Definition

TODO: Define the concept.

## Examples

TODO: Add examples.

## Practice

TODO: Add practice problems.
MDEOF
  echo "Created $LESSON_FILE"
else
  echo "$LESSON_FILE already exists, skipping"
fi

# Check if already registered
if [ -f "$LESSONS_JSON" ]; then
  if grep -q "\"$CONCEPT_ID\"" "$LESSONS_JSON" 2>/dev/null; then
    echo "Concept '$CONCEPT_ID' already registered in $LESSONS_JSON"
  else
    # Append to lessons.json (remove trailing bracket, add comma, add entry, close)
    tmp=$(mktemp)
    # Use python for reliable JSON manipulation
    python3 -c "
import json
with open('$LESSONS_JSON') as f:
    data = json.load(f)
data.append({'concept_id': '$CONCEPT_ID', 'source': '$LESSON_PATH'})
with open('$LESSONS_JSON', 'w') as f:
    json.dump(data, f, indent=2)
    f.write('\n')
"
    echo "Registered '$CONCEPT_ID' -> '$LESSON_PATH' in $LESSONS_JSON"
  fi
fi

echo ""
echo "Done! Next: write the lesson content in $LESSON_FILE"
