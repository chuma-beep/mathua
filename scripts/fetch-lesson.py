#!/usr/bin/env python3
"""
Fetch a Wikipedia page/section and convert it into a mathua lesson markdown file.

Usage:
  python3 scripts/fetch-lesson.py <concept_id> <wikipedia_title> [section_numbers...]

Examples:
  python3 scripts/fetch-lesson.py alg.linear.slope Slope 2 3
  python3 scripts/fetch-lesson.py trig.sin_cos_def "Sine and cosine" 1 2
  python3 scripts/fetch-lesson.py --batch scripts/wiki-map.json

Output:
  Writes to data/lessons/wikipedia/<concept_id>.md
  Appends entry to data/lessons/lessons.json
"""

import html
import json
import os
import re
import sys
import time
import urllib.request
import urllib.parse
import urllib.error

API = "https://en.wikipedia.org/w/api.php"
LESSONS_DIR = "data/lessons"
WIKI_DIR = os.path.join(LESSONS_DIR, "wikipedia")
LESSONS_JSON = os.path.join(LESSONS_DIR, "lessons.json")
CONCEPTS_JSON = "web/next-app/data/concepts.json"


def load_concepts():
    with open(CONCEPTS_JSON) as f:
        return {c["id"]: c for c in json.load(f)}


def load_lessons_json():
    if not os.path.exists(LESSONS_JSON):
        return []
    with open(LESSONS_JSON) as f:
        return json.load(f)


def save_lessons_json(entries):
    os.makedirs(LESSONS_DIR, exist_ok=True)
    with open(LESSONS_JSON, "w") as f:
        json.dump(entries, f, indent=2)
        f.write("\n")


HEADERS = {"User-Agent": "Mathua/1.0 (lesson-fetcher; mathua-project)"}
LAST_CALL = 0


def api_call(params):
    global LAST_CALL
    elapsed = time.time() - LAST_CALL
    if elapsed < 0.5:
        time.sleep(0.5 - elapsed)
    url = API + "?" + urllib.parse.urlencode(params)
    req = urllib.request.Request(url, headers=HEADERS)
    for attempt in range(3):
        try:
            LAST_CALL = time.time()
            with urllib.request.urlopen(req, timeout=15) as resp:
                return json.loads(resp.read())
        except urllib.error.HTTPError as e:
            if e.code == 429 and attempt < 2:
                wait = 2 ** (attempt + 1)
                print(f"  [rate-limit] retrying in {wait}s...", file=sys.stderr)
                time.sleep(wait)
                continue
            raise


def fetch_wikitext(title, section=None):
    params = {
        "action": "parse",
        "page": title,
        "prop": "wikitext",
        "format": "json",
        "formatversion": "2",
    }
    if section is not None:
        params["section"] = str(section)
    try:
        data = api_call(params)
        return data["parse"]["wikitext"]
    except Exception as e:
        print(f"  [error] fetching section {section}: {e}", file=sys.stderr)
        return ""


def fetch_sections(title):
    params = {
        "action": "parse",
        "page": title,
        "prop": "tocdata",
        "format": "json",
        "formatversion": "2",
    }
    try:
        data = api_call(params)
        tocdata = data["parse"].get("tocdata", {})
        if isinstance(tocdata, dict):
            return tocdata.get("sections", [])
        return []
    except Exception as e:
        print(f"  [warn] fetching TOC for '{title}': {e}", file=sys.stderr)
        return []


def remove_file_refs(text):
    """Remove [[File:...]], [[Image:...]], [[Category:...]] entirely."""
    # Remove [[File:...]], [[Image:...]], [[Category:...]] handling nested [[...]]
    result = []
    i = 0
    while i < len(text):
        m = re.match(r'\[\[(File|Image|Category):', text[i:])
        if m:
            depth = 1
            j = i + 2  # past the opening [[
            while j < len(text) and depth > 0:
                if text[j:j+2] == '[[':
                    depth += 1
                    j += 2
                elif text[j:j+2] == ']]':
                    depth -= 1
                    j += 2
                else:
                    j += 1
            i = j
        else:
            result.append(text[i])
            i += 1
    return ''.join(result)


def remove_gallery(text):
    """Remove <gallery>...</gallery> blocks."""
    return re.sub(r'<gallery>.*?</gallery>', '', text, flags=re.DOTALL)


def remove_refs(text):
    """Remove <ref>...</ref> and <ref .../>."""
    text = re.sub(r'<ref[^>]*/>', '', text)
    text = re.sub(r'<ref[^>]*>.*?</ref>', '', text, flags=re.DOTALL)
    return text


def remove_templates(text):
    """Remove {{...}} templates, converting simple formatting ones."""
    # Convert {{frac|a|b}} → a/b
    text = re.sub(r'\{\{frac\|([^|}]+)\|([^|}]+)\}\}', r'\1/\2', text)
    # Convert {{sqrt|a}} → sqrt(a)
    text = re.sub(r'\{\{sqrt\|([^}]+)\}\}', r'sqrt(\1)', text)
    # Convert {{abs|a}} → |a|
    text = re.sub(r'\{\{abs\|([^}]+)\}\}', r'|\1|', text)
    # Convert {{nobreak|...}}, {{nowrap|...}} → just content
    text = re.sub(r'\{\{(?:nobreak|nowrap)\|([^}]+)\}\}', r'\1', text)
    # Convert {{math|...}} → just content  
    text = re.sub(r'\{\{math\|([^}]+)\}\}', r'\1', text)
    # Remove all remaining {{...}} templates
    text = re.sub(r'\{\{[^}]*\}\}', '', text)
    # Clean up any leftover single braces from unmatched templates
    text = re.sub(r'^[\}\{]+', '', text)
    text = re.sub(r'[\}\{]+$', '', text)
    return text


def convert_links(text):
    """Convert [[Target|display]] → display and [[Target]] → Target."""
    text = re.sub(r'\[\[([^\]]*\|([^\]]*))\]\]', r'\2', text)
    text = re.sub(r'\[\[([^\]]+)\]\]', r'\1', text)
    text = re.sub(r'\[(https?://[^\s\]]+)\s+([^\]]+)\]', r'\2', text)
    text = re.sub(r'\[(https?://[^\]]+)\]', r'\1', text)
    return text


def convert_formatting(text):
    """Convert wiki formatting to markdown."""
    text = re.sub(r"'''(.*?)'''", r'**\1**', text)
    text = re.sub(r"''(.*?)''", r'*\1*', text)
    return text


def convert_math(text):
    """
    Convert <math>...</math> to LaTeX.
    - Inline: wrap in \\( ... \\)
    - Display (on its own line, possibly : indented): wrap in \\[ ... \\]
    """
    def replacer(match):
        content = match.group(2).strip()
        attrs = (match.group(1) or '').lower()
        is_display = 'display="block"' in attrs or 'display="displaystyle"' in attrs
        if is_display:
            return '\n\\[\n' + content + '\n\\]\n'
        else:
            return '\\(' + content + '\\)'

    text = re.sub(r'<math\s*([^>]*)>(.*?)</math>', replacer, text, flags=re.DOTALL)
    return text


def clean_wikitext(wikitext):
    """Convert raw wikitext to clean lesson content."""
    text = wikitext
    # Convert wiki headings (== X == → ## X, === X === → ### X, etc.)
    def heading_replacer(m):
        level = len(m.group(1))
        return '#' * level + ' ' + m.group(2)
    text = re.sub(r'^(==+)\s*(.*?)\s*\1\s*$', heading_replacer, text, flags=re.MULTILINE)
    # Strip : indent markers (used for display math, already handled)
    text = re.sub(r'^:\s*', '', text, flags=re.MULTILINE)
    # Remove comments
    text = re.sub(r'<!--.*?-->', '', text, flags=re.DOTALL)
    # Remove __TOC__, __NOTOC__, etc.
    text = re.sub(r'__(?:TOC|NOTOC|FORCETOC|NOEDITSECTION)__', '', text)
    # Remove <includeonly>, <noinclude>, <onlyinclude>
    text = re.sub(r'<includeonly>.*?</includeonly>', '', text, flags=re.DOTALL)
    text = re.sub(r'</?noinclude>', '', text, flags=re.DOTALL)
    text = re.sub(r'</?onlyinclude>', '', text, flags=re.DOTALL)
    # Remove templates, refs, files FIRST (before math conversion to avoid mangling LaTeX)
    text = remove_templates(text)
    text = remove_refs(text)
    text = remove_file_refs(text)
    text = remove_gallery(text)
    # Process math (convert <math>...</math> to LaTeX)
    text = convert_math(text)
    # Remove any remaining HTML tags except <b>, <i>, <em>, <strong>, <sub>, <sup>
    text = re.sub(r'</(?!b>|i>|em>|strong>|sub>|sup>)[a-z]+\s*>', '', text)
    text = re.sub(r'<(?![/]?b>|[/]?i>|[/]?em>|[/]?strong>|[/]?sub>|[/]?sup>)[a-z]+[^>]*>', '', text, flags=re.DOTALL)
    # Convert links and formatting
    text = convert_links(text)
    text = convert_formatting(text)
    # Clean up excessive whitespace
    text = re.sub(r'\n{4,}', '\n\n\n', text)
    text = re.sub(r' +\n', '\n', text)
    # Remove :indent after processing
    text = re.sub(r'^:\s*', '', text, flags=re.MULTILINE)
    # Remove empty sections (heading followed by nothing or just whitespace)
    text = re.sub(r'^#{2,3}\s+.*?\n(?=\n#{2,3}|\Z)', '', text, flags=re.MULTILINE)
    # Decode HTML entities
    text = html.unescape(text)
    return text.strip()


def generate_lesson(concept_id, title, sections_text):
    """Assemble the final lesson markdown."""
    # Load concept label for the H1 title
    concepts = load_concepts()
    concept = concepts.get(concept_id, {})
    lesson_title = concept.get("label", title)

    lines = [
        f'> Content sourced from [Wikipedia](https://en.wikipedia.org/wiki/{urllib.parse.quote(title.replace(" ", "_"))}) — CC BY-SA 4.0',
        '',
        f'# {lesson_title}',
        '',
        sections_text,
        '',
    ]
    return '\n'.join(lines)


def fetch_and_generate(concept_id, wiki_title, sections):
    """Fetch a Wikipedia page and generate a lesson file."""
    os.makedirs(WIKI_DIR, exist_ok=True)
    out_path = os.path.join(WIKI_DIR, f"{concept_id}.md")

    if sections:
        toc = fetch_sections(wiki_title) or []
        # Build parent mapping: for each section index, find its top-level parent
        parent_top = {}
        current_top = None
        for s in toc:
            if s.get("tocLevel") == 1:
                current_top = s["index"]
            parent_top[s["index"]] = current_top
        requested_top = {str(sec) for sec in sections}
        parts = []
        seen_heading_texts = set()
        for sec in sections:
            sec_str = str(sec)
            # Skip sub-sections whose parent is also being fetched (content already included)
            if parent_top.get(sec_str) and parent_top[sec_str] in requested_top and parent_top[sec_str] != sec_str:
                continue
            print(f"  fetching section {sec}...")
            wt = fetch_wikitext(wiki_title, sec)
            if not wt:
                continue
            cleaned = clean_wikitext(wt)
            if not cleaned.strip():
                continue
            # Deduplicate: if this sub-section's heading was already included by a parent, skip the heading
            lines = cleaned.split('\n')
            deduped = []
            for line in lines:
                if re.match(r'^#{2,4}\s+\S', line):
                    heading_text = line.strip()
                    if heading_text in seen_heading_texts:
                        continue
                    seen_heading_texts.add(heading_text)
                deduped.append(line)
            parts.append('\n'.join(deduped))
        sections_text = '\n\n'.join(parts)
    else:
        print(f"  fetching full page...")
        wt = fetch_wikitext(wiki_title)
        if not wt:
            print(f"  [error] could not fetch page '{wiki_title}'", file=sys.stderr)
            return None
        sections_text = clean_wikitext(wt)

    lesson = generate_lesson(concept_id, wiki_title, sections_text)

    with open(out_path, "w") as f:
        f.write(lesson)
        f.write("\n")

    print(f"  wrote {out_path}")
    return out_path


def register_lesson(concept_id, relative_path):
    """Register the lesson in lessons.json."""
    entries = load_lessons_json()
    # Check if concept already has an entry
    already = any(e.get("concept_id") == concept_id for e in entries)
    if already:
        print(f"  [warn] concept '{concept_id}' already in lessons.json, skipping", file=sys.stderr)
        return

    entries.append({
        "concept_id": concept_id,
        "source": relative_path,
    })
    save_lessons_json(entries)
    print(f"  registered in {LESSONS_JSON}")


def process_batch(batch_path):
    """Process a batch mapping file (JSON array of {concept_id, wikipedia, sections})."""
    with open(batch_path) as f:
        batch = json.load(f)

    for item in batch:
        cid = item["concept_id"]
        wp = item["wikipedia"]
        secs = item.get("sections", [])
        print(f"\n{cid} ← Wikipedia: {wp}" + (f" sections {secs}" if secs else ""))
        out = fetch_and_generate(cid, wp, secs)
        if out:
            rel_path = "wikipedia/" + cid + ".md"
            register_lesson(cid, rel_path)


def main():
    if len(sys.argv) < 3:
        print(__doc__)
        sys.exit(1)

    if sys.argv[1] == "--batch":
        process_batch(sys.argv[2])
        return

    concept_id = sys.argv[1]
    wiki_title = sys.argv[2]
    sections = [int(s) for s in sys.argv[3:]] if len(sys.argv) > 3 else []

    print(f"\n{concept_id} ← Wikipedia: {wiki_title}" + (f" sections {sections}" if sections else ""))
    out = fetch_and_generate(concept_id, wiki_title, sections)
    if out:
        rel_path = "wikipedia/" + concept_id + ".md"
        register_lesson(concept_id, rel_path)


if __name__ == "__main__":
    main()
