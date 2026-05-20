#!/usr/bin/env python3
"""
Convert external math datasets to mathua's JSON import format.

Usage:
    python scripts/convert-math-dataset.py --source huggingface --dataset imyangyixuan/hendrycks_math
    python scripts/convert-math-dataset.py --source huggingface --dataset math-ai/StackMathQA --limit 10000

Maps dataset subjects to mathua concept IDs using keyword matching.
Output: questions-import.json (ready for `go run scripts/seed-questions/main.go -import`)
"""

import argparse
import json
import re
import sys

# Mapping from external subject labels to mathua concept ID prefixes
SUBJECT_MAP = {
    'algebra': 'alg',
    'intermediate algebra': 'alg',
    'prealgebra': 'prealg',
    'pre-algebra': 'prealg',
    'geometry': 'geo',
    'counting & probability': 'count',
    'counting and probability': 'count',
    'number theory': 'nt',
    'precalculus': 'prealg',
    'calculus': 'calc',
    'trigonometry': 'trig',
    'statistics': 'stat',
    'linear algebra': 'linalg',
    'arithmetic': 'arith',
    'fractions': 'frac',
    'complex numbers': 'complex',
    'discrete math': 'discrete',
}

def map_to_concept_id(subject: str, problem_text: str) -> str | None:
    prefix = None
    subj_lower = subject.lower().strip()
    for key, val in SUBJECT_MAP.items():
        if key in subj_lower:
            prefix = val
            break
    if not prefix:
        return None
    return f"{prefix}.external"

def clean_answer(answer: str) -> str:
    answer = re.sub(r'\\boxed\{([^}]*)\}', r'\1', answer)
    answer = answer.strip()
    return answer

def convert_huggingface(dataset_name: str, limit: int | None = None):
    try:
        from datasets import load_dataset
    except ImportError:
        print("Error: install datasets library: pip install datasets", file=sys.stderr)
        sys.exit(1)

    ds = load_dataset(dataset_name, split='train')
    if limit:
        ds = ds.select(range(min(limit, len(ds))))

    questions = []
    for i, row in enumerate(ds):
        problem = row.get('problem', '')
        solution = row.get('solution', '')
        subject = row.get('type', '')
        level = row.get('level', '')

        if not problem:
            continue

        concept_id = map_to_concept_id(subject, problem)
        if not concept_id:
            continue

        answer = clean_answer(solution.split('\n')[-1] if solution else '')

        questions.append({
            'concept_id': concept_id,
            'question': problem,
            'answer': answer or 'See solution',
            'explanation': solution or '',
            'source': dataset_name,
            'difficulty': float(level[-1]) / 5.0 if level and level[-1].isdigit() else 0.5,
        })

    return questions

def convert_json_file(input_path: str, subject_key: str = 'subject', 
                       problem_key: str = 'question', answer_key: str = 'answer',
                       solution_key: str = 'solution'):
    with open(input_path) as f:
        data = json.load(f)

    questions = []
    for row in data:
        problem = row.get(problem_key, '')
        if not problem:
            continue
        subject = row.get(subject_key, '')
        concept_id = map_to_concept_id(subject, problem)
        if not concept_id:
            continue
        questions.append({
            'concept_id': concept_id,
            'question': problem,
            'answer': row.get(answer_key, ''),
            'explanation': row.get(solution_key, ''),
            'source': input_path,
            'difficulty': 0.5,
        })
    return questions

if __name__ == '__main__':
    parser = argparse.ArgumentParser(description='Convert math datasets to mathua JSON format')
    parser.add_argument('--source', choices=['huggingface', 'json'], required=True)
    parser.add_argument('--dataset', help='HuggingFace dataset name (for --source huggingface)')
    parser.add_argument('--input', help='JSON file path (for --source json)')
    parser.add_argument('--limit', type=int, help='Max number of questions to convert')
    parser.add_argument('--output', default='questions-import.json', help='Output file path')

    args = parser.parse_args()

    if args.source == 'huggingface' and not args.dataset:
        parser.error('--dataset is required for huggingface source')
    if args.source == 'json' and not args.input:
        parser.error('--input is required for json source')

    if args.source == 'huggingface':
        questions = convert_huggingface(args.dataset, args.limit)
    else:
        questions = convert_json_file(args.input)

    with open(args.output, 'w') as f:
        json.dump(questions, f, indent=2)

    print(f"Converted {len(questions)} questions → {args.output}")
    print(f"Run: go run scripts/seed-questions/main.go -import {args.output}")
