#!/usr/bin/env python3
"""
Post-process all Wikipedia lesson files to:
1. Convert <sup>/<sub> tags to LaTeX
2. Convert {{tmath|...}} templates to LaTeX
3. Strip irrelevant sections (See also, Notes, References, etc.)
4. Remove template residue (1=, =}}, |2}}, etc.)
5. Flag locations where diagrams would be useful
6. Clean up missing variable names and broken template output

Usage:
  python3 scripts/clean-lessons.py              # clean all Wikipedia lessons
  python3 scripts/clean-lessons.py --report     # only print issues, don't modify
  python3 scripts/clean-lessons.py --file <id>  # clean a single lesson
"""

import html
import json
import os
import re
import sys

WIKI_DIR = "data/lessons/wikipedia"
LESSONS_JSON = "data/lessons/lessons.json"

DIAGRAM_KEYWORDS = [
    r'\bas shown\b',
    r'\bdiagram\b',
    r'\bfigure\b',
    r'\billustrat(?:ion|ed|rates?)\b',
    r'\bpictur(?:e|es)\b',
    r'\bsketch\b',
    r'\bgraph\b',
    r'\bplot\b',
]

IRRELEVANT_SECTIONS = re.compile(
    r'^#{2,3}\s*(See also|Notes|References|Further reading|'
    r'External links|Bibliography|Sources|Footnotes)\s*$',
    re.MULTILINE | re.IGNORECASE,
)


def inside_math(text: str, pos: int) -> bool:
    before = text[:pos]
    return before.rfind('\\(') > before.rfind('\\)') or before.rfind('\\[') > before.rfind('\\]')


def convert_sup_sub(text: str) -> str:
    def replace_tag(tag, content, pos):
        if inside_math(text, pos):
            return '^{' + content + '}' if tag == 'sup' else '_{' + content + '}'
        return '\\(^{' + content + '}\\)' if tag == 'sup' else '\\(_{' + content + '}\\)'

    # Process <sub> first (inner), then <sup> (may contain <sub>)
    for tag_type in ['sub', 'sup']:
        open_tag = '<' + tag_type + '>'
        close_tag = '</' + tag_type + '>'
        result = []
        i = 0
        while i < len(text):
            pos = text.find(open_tag, i)
            if pos < 0:
                result.append(text[i:])
                break
            result.append(text[i:pos])
            content_start = pos + len(open_tag)
            close_pos = text.find(close_tag, content_start)
            if close_pos >= 0:
                content = text[content_start:close_pos]
                result.append(replace_tag(tag_type, content, pos))
                i = close_pos + len(close_tag)
            else:
                result.append(text[pos])
                i = pos + 1
        text = ''.join(result)
    return text


def convert_tmath(text: str) -> str:
    def replacer(m):
        raw = m.group(0)
        inner = raw[len('{{tmath|'):-len('}}')]
        if inner.startswith('1='):
            inner = inner[2:]
        inner = inner.strip()
        if '\n' in inner or len(inner) > 80:
            return "\n\\[\n" + inner + "\n\\]\n"
        return "\\(" + inner + "\\)"

    result = []
    i = 0
    while i < len(text):
        m = re.search(r'\{\{tmath\|', text[i:])
        if not m:
            result.append(text[i:])
            break
        result.append(text[i:i + m.start()])
        start = i + m.end()
        depth = 1
        j = start
        while j < len(text) and depth > 0:
            if text[j:j+2] == '{{':
                depth += 1
                j += 2
            elif text[j:j+2] == '}}':
                depth -= 1
                j += 2
            else:
                j += 1
        inner = text[start:j-2].strip()
        if inner.startswith('1='):
            inner = inner[2:]
        inner = inner.strip()
        if '\n' in inner or len(inner) > 80:
            result.append("\n\\[\n" + inner + "\n\\]\n")
        else:
            result.append("\\(" + inner + "\\)")
        i = j
    return ''.join(result)


def convert_mvar(text: str) -> str:
    def replacer(m):
        content = m.group(1)
        return "\\(" + content + "\\)"

    text = re.sub(r'\{\{mvar\|([^}]+)\}\}', replacer, text)
    return text


def strip_irrelevant_sections(text: str) -> str:
    match = IRRELEVANT_SECTIONS.search(text)
    if match:
        text = text[:match.start()]
    text = re.sub(r'\n{3,}', '\n\n', text)
    return text.strip()


def clean_template_residue(text: str) -> str:
    text = re.sub(r'\b1=(?:sin|cos|tan|log|ln)\([^)]*\)\s*=\s*', '', text)
    text = re.sub(r'\b1=(\w+)\([^)]*\)\s*=\s*', r'\1', text)
    text = re.sub(r'\b1=∠[^=]*=\s*', '', text)
    text = re.sub(r'\|\d+\}\}', '', text)
    text = re.sub(r'= \}\}', '', text)
    text = re.sub(r'\b\d+=\}\}', '', text)
    text = re.sub(r'\b1=\}', '', text)
    text = re.sub(r'\}\} ?', '', text)
    # Strip remaining "1=" template parameter prefixes (param 1 = value)
    text = re.sub(r'\b1=\s*', '', text)
    text = re.sub(r' +', ' ', text)
    text = re.sub(r'^:\s*', '', text, flags=re.MULTILINE)
    text = re.sub(r'\n{3,}', '\n\n', text)
    return text


def fix_missing_vars(text: str) -> str:
    text = re.sub(r'\b(?:is|the|a|an)\s+[,.?:;]', '', text)
    text = re.sub(r'\bthe\s+of\b', '', text)
    text = re.sub(r'\bof the\s+[,.?:;]', '', text)
    text = re.sub(r',\s*,\s*', ', ', text)
    text = re.sub(r'\s+\.\s*', '. ', text)
    text = re.sub(r'\bsince\s+and\b', 'since', text)
    text = re.sub(r',\s+since\b', '; since', text)
    return text


def clean_math_blocks(text: str) -> str:
    def replacer(m):
        block = m.group(1).strip()
        if not block or block in ('.', ',', ';', ':'):
            return ''
        return m.group(0)

    text = re.sub(r'\\\(\s*([^)]*?)\s*\\\)', replacer, text)
    text = re.sub(r'\\\[\s*([^]]*?)\s*\\\]', replacer, text)
    return text


def convert_span_frac(text: str) -> str:
    text = re.sub(
        r'<span class="frac" role="math">.*?</span>',
        '',
        text,
        flags=re.DOTALL,
    )
    return text


def find_diagram_references(text: str, filepath: str) -> list:
    refs = []
    for pattern in DIAGRAM_KEYWORDS:
        for m in re.finditer(pattern, text, re.IGNORECASE):
            content_before = text[max(0, m.start() - 80):m.start()]
            if re.search(r'\\\(|\\\[', content_before):
                continue
            line_num = text[:m.start()].count('\n') + 1
            start = max(0, m.start() - 40)
            end = min(len(text), m.end() + 60)
            context = text[start:end].replace('\n', ' ')
            refs.append({
                'file': filepath,
                'line': line_num,
                'keyword': m.group(),
                'context': '...' + context + '...',
            })
    return refs


def find_template_residue(text: str, filepath: str) -> list:
    issues = []
    for m in re.finditer(r'\b1=(?:sin|cos|tan|log|ln)\(', text):
        line_num = text[:m.start()].count('\n') + 1
        issues.append({
            'file': filepath,
            'line': line_num,
            'type': 'template_residue',
            'context': text[max(0, m.start() - 20):m.end() + 30].replace('\n', ' '),
        })
    return issues


def find_missing_vars(text: str, filepath: str) -> list:
    issues = []
    for m in re.finditer(r'\b(is|the)\s+[,.?:;]', text):
        line_num = text[:m.start()].count('\n') + 1
        issues.append({
            'file': filepath,
            'line': line_num,
            'type': 'missing_variable',
            'context': text[max(0, m.start() - 30):m.end() + 10].replace('\n', ' '),
        })
    return issues


def process_file(filepath: str, dry_run: bool = False) -> dict:
    filename = os.path.basename(filepath)
    with open(filepath) as f:
        original = f.read()

    text = original
    issues = []

    text = convert_sup_sub(text)
    text = convert_tmath(text)
    text = convert_mvar(text)
    text = convert_span_frac(text)
    text = clean_template_residue(text)
    text = fix_missing_vars(text)
    text = clean_math_blocks(text)
    # Final pass: strip any remaining <sup>/<sub> tags (keep content)
    text = re.sub(r'</?(?:sup|sub)>', '', text)
    # Strip remaining {{...}} templates with brace matching
    def strip_braces(text):
        result = []
        i = 0
        while i < len(text):
            if text[i:i+2] == '{{':
                depth = 1
                j = i + 2
                while j < len(text) and depth > 0:
                    if text[j:j+2] == '{{':
                        depth += 1
                        j += 2
                    elif text[j:j+2] == '}}':
                        depth -= 1
                        j += 2
                    else:
                        j += 1
                # Extract last content param (the most meaningful part)
                inner = text[i+2:j-2]
                if '|' in inner:
                    parts = inner.split('|')
                    content = [p for p in parts if '=' not in p]
                    if content:
                        result.append(content[-1])
                i = j
            else:
                result.append(text[i])
                i += 1
        return ''.join(result)
    text = strip_braces(text)

    diag_refs = find_diagram_references(text, filename)
    tmpl_issues = find_template_residue(text, filename)
    missing_issues = find_missing_vars(text, filename)

    issues.extend(diag_refs)
    issues.extend(tmpl_issues)
    issues.extend(missing_issues)

    text = strip_irrelevant_sections(text)

    if text and not text.endswith('\n'):
        text += '\n'

    changed = text != original

    if changed and not dry_run:
        with open(filepath, 'w') as f:
            f.write(text)

    return {
        'file': filename,
        'changed': changed,
        'issues': issues,
    }


def main():
    dry_run = '--report' in sys.argv

    if not os.path.isdir(WIKI_DIR):
        print(f"Error: {WIKI_DIR} not found", file=sys.stderr)
        sys.exit(1)

    files = sorted(os.listdir(WIKI_DIR))
    md_files = [f for f in files if f.endswith('.md')]

    single_file = None
    for i, arg in enumerate(sys.argv):
        if arg == '--file' and i + 1 < len(sys.argv):
            single_file = sys.argv[i + 1]
            break

    if single_file:
        filepath = os.path.join(WIKI_DIR, single_file)
        if not os.path.exists(filepath):
            filepath = os.path.join(WIKI_DIR, single_file + '.md')
        if not os.path.exists(filepath):
            print(f"Error: lesson '{single_file}' not found", file=sys.stderr)
            sys.exit(1)
        md_files = [os.path.basename(filepath)]

    total_changed = 0
    total_issues = []
    total_sup_sub = 0
    total_tmath = 0

    for fname in md_files:
        filepath = os.path.join(WIKI_DIR, fname)
        result = process_file(filepath, dry_run=dry_run)
        if result['changed']:
            total_changed += 1
        total_issues.extend(result['issues'])

        with open(filepath) as f:
            content = f.read()
        total_sup_sub += len(re.findall(r'<sup[^>]*>|<sub[^>]*>', content))
        total_tmath += len(re.findall(r'\{\{tmath', content))

    if dry_run:
        print(f"\n=== REPORT: {len(md_files)} files inspected ===")
    else:
        print(f"\n=== Cleaned {total_changed}/{len(md_files)} files ===")

    if total_sup_sub > 0:
        print(f"\n⚠ Remaining <sup>/<sub> tags: {total_sup_sub}")
    if total_tmath > 0:
        print(f"⚠ Remaining {{tmath}} templates: {total_tmath}")

    if total_issues:
        diag_issues = [i for i in total_issues if i.get('keyword')]
        tmpl_issues = [i for i in total_issues if i.get('type') == 'template_residue']
        missing_issues = [i for i in total_issues if i.get('type') == 'missing_variable']

        if diag_issues:
            print(f"\n🔷 Diagram references found: {len(diag_issues)}")
            for iss in diag_issues[:15]:
                print(f"  {iss['file']}:{iss['line']} — \"{iss['keyword']}\"")
                print(f"    context: {iss['context'][:100]}")
            if len(diag_issues) > 15:
                print(f"  ... and {len(diag_issues) - 15} more")

        if tmpl_issues:
            print(f"\n⚠ Template residue: {len(tmpl_issues)}")
            for iss in tmpl_issues[:10]:
                print(f"  {iss['file']}:{iss['line']} — {iss['context'][:100]}")

        if missing_issues:
            print(f"\n⚠ Missing variables: {len(missing_issues)}")
            for iss in missing_issues[:10]:
                print(f"  {iss['file']}:{iss['line']} — {iss['context'][:100]}")

    if not dry_run:
        print("\nDone.")


if __name__ == "__main__":
    main()
