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
    if elapsed < 1.5:
        time.sleep(1.5 - elapsed)
    url = API + "?" + urllib.parse.urlencode(params)
    req = urllib.request.Request(url, headers=HEADERS)
    for attempt in range(3):
        try:
            LAST_CALL = time.time()
            with urllib.request.urlopen(req, timeout=15) as resp:
                return json.loads(resp.read())
        except urllib.error.HTTPError as e:
            if e.code == 429 and attempt < 3:
                wait = 5 * (2 ** attempt)
                print(f"  [rate-limit] retrying in {wait}s...", file=sys.stderr)
                time.sleep(wait)
                continue
            raise


def check_redirect(title):
    """Check if a page is a redirect before fetching content."""
    params = {
        "action": "query",
        "titles": title,
        "redirects": "1",
        "format": "json",
        "formatversion": "2",
    }
    try:
        data = api_call(params)
        pages = data.get("query", {}).get("pages", [])
        for p in pages:
            if "redirect" in p.get("title", "").lower() or p.get("pageid", 0) == 0:
                return True
        # Check if any redirect happened
        if data.get("query", {}).get("redirects"):
            return True
    except Exception:
        pass
    return False


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
        wt = data["parse"]["wikitext"]
        if wt.strip().startswith("#REDIRECT"):
            print(f"  [skip] page/section is a redirect: {title}", file=sys.stderr)
            return ""
        return wt
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
    """Remove remaining [[Category:...]] and any unconverted [[File:...]]/[[Image:...]] refs."""
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


def _find_matching_end(text, start):
    """Find the matching }} for {{ at start, handling nested {{...}}."""
    if text[start:start+2] != '{{':
        return -1
    depth = 1
    i = start + 2
    while i < len(text) - 1:
        if text[i:i+2] == '{{':
            depth += 1
            i += 2
        elif text[i:i+2] == '}}':
            depth -= 1
            i += 2
            if depth == 0:
                return i
        else:
            i += 1
    return -1


def _split_params(param_str):
    """Split template parameter string on |, respecting nested {{...}}."""
    parts = []
    depth = 0
    cur = []
    i = 0
    while i < len(param_str):
        ch = param_str[i]
        if ch == '{' and param_str[i:i+2] == '{{':
            depth += 1
            cur.append(ch)
            i += 1
        elif ch == '}' and param_str[i:i+2] == '}}':
            depth -= 1
            cur.append(ch)
            i += 1
        elif ch == '|' and depth == 0:
            parts.append(''.join(cur).strip())
            cur = []
        else:
            cur.append(ch)
        i += 1
    if cur:
        parts.append(''.join(cur).strip())
    return parts


def _strip_named_prefix(param):
    """Strip '1=', '2=', etc. prefix from a template parameter."""
    m = re.match(r'^\d+=(.*)', param)
    if m:
        return m.group(1)
    return param


def _process_known_template(name, params):
    """Process a known template by name, returning replacement text."""
    name = name.strip().lower()

    if name == 'tmath':
        if not params:
            return ''
        inner = _strip_named_prefix(params[0])
        return '\\(%s\\)' % inner.strip()
    if name in ('mvar', 'var'):
        if not params:
            return ''
        return '\\(%s\\)' % _strip_named_prefix(params[0]).strip()
    if name == 'math':
        if not params:
            return ''
        return _strip_named_prefix(params[0]).strip()
    if name == 'frac':
        if len(params) >= 2:
            return '%s/%s' % (_strip_named_prefix(params[0]), _strip_named_prefix(params[1]))
        return ''
    if name in ('sfrac', 'cfrac', 'xfrac'):
        if len(params) >= 2:
            return '%s/%s' % (_strip_named_prefix(params[0]), _strip_named_prefix(params[1]))
        return ''
    if name == 'sqrt':
        if not params:
            return ''
        return 'sqrt(%s)' % _strip_named_prefix(params[0])
    if name == 'abs':
        if not params:
            return ''
        return '|%s|' % _strip_named_prefix(params[0])
    if name in ('nobreak', 'nowrap'):
        if not params:
            return ''
        return _strip_named_prefix(params[0])
    if name in ('big', 'small', 'center', 'left', 'right'):
        if not params:
            return ''
        return _strip_named_prefix(params[0])
    if name in ('abbr', 'tooltip'):
        if not params:
            return ''
        return _strip_named_prefix(params[0])
    if name == 'chem':
        if not params:
            return ''
        return '\\ce{%s}' % _strip_named_prefix(params[0])
    if name == 'lang':
        return _strip_named_prefix(params[-1]) if params else ''
    if name in ('ill', 'interlanguage link'):
        return _strip_named_prefix(params[0]) if params else ''
    if name == 'val':
        if not params:
            return ''
        return _strip_named_prefix(params[0])
    if name in ('su', 'sup', 'sub'):
        return _strip_named_prefix(params[0]) if params else ''
    if name in ('clarify', 'dubious', 'definition', 'example'):
        return _strip_named_prefix(params[0]) if params else ''
    if name == 'not a typo':
        return _strip_named_prefix(params[0]) if params else ''
    if name == 'proper name':
        return _strip_named_prefix(params[0]) if params else ''
    if name in ('citation needed', 'cn', 'fact', 'who', 'when', 'where'):
        return ''
    if name == 'pi':
        return 'π'
    if name in ('e (mathematical constant)', 'munich'):
        return 'e'
    if name in ('i (number)', 'imaginary unit'):
        return 'i'
    if name in ('sqrt', 'radical'):
        return '√'
    if name == 'p/sigma':
        return 'σ'
    if name == 'pi/phi':
        return 'φ'
    if name == "euler's number":
        return 'e'
    if name == 'prime':
        return '′'
    if name == 'prime/prime':
        return '″'
    if name == '1/2':
        return '½'
    if name == '1/3':
        return '⅓'
    if name == '2/3':
        return '⅔'
    if name == '1/4':
        return '¼'
    if name == '3/4':
        return '¾'
    if name == '1/8':
        return '⅛'
    if name in ('ndash', 'spaced ndash'):
        return ' – '
    if name in ('mdash', 'spaced mdash'):
        return ' — '
    if name in ('clear', '-', 'clear left', 'clear right'):
        return ''
    if name in ('sfn', 'harvnb', 'sfnp', 'harv', 'harvc', 'sfnref'):
        return ''
    if name in ('reflist', 'notelist', 'notelist-ua', 'notelist-lr'):
        return ''
    if name in ('colbegin', 'colend', 'col-begin', 'col-end', 'div col', 'div col end'):
        return ''
    if name in ('plainlist', 'plain list', 'ubl', 'unbulleted list', 'ordered list'):
        return ''
    if name in ('short description', 'use dmy dates', 'use mdy dates', 'use british english',
                'use american english', 'good article', 'featured article'):
        return ''
    if name in ('infobox', 'infobox website', 'infobox book', 'infobox person',
                'infobox scientist', 'infobox mathematician'):
        return ''
    if name == '':
        return ''
    if name == '=':
        return '='
    if name == '!':
        return '|'
    if name == '(':
        return '{{'
    if name == ')':
        return '}}'
    if name == "'":
        return "'"
    if name == "'s" or name == "'s":
        return "'s"
    if name in ('displaystyle', 'dspl', 'd'):
        return params[0] if params else ''
    if name == 'strlen':
        return ''
    if name in ('main', 'see', 'further', 'details', 'category see also'):
        return ''
    if name in ('anchor', 'visible anchor', 'vanchor'):
        return ''
    if name == 'as of':
        return _strip_named_prefix(params[0]) if params else ''

    # Unknown template: try to extract text from positional parameters
    if params:
        # Skip named params, collect positional ones
        text_parts = []
        for p in params:
            if '=' not in p:
                text_parts.append(_strip_named_prefix(p))
        if text_parts:
            return ''.join(text_parts)
    return ''


def remove_templates(text):
    """Remove {{...}} templates, converting simple formatting ones.
    Handles nested braces (including LaTeX) inside templates."""
    result = []
    i = 0
    while i < len(text):
        if text[i:i+2] == '{{':
            end = _find_matching_end(text, i)
            if end == -1:
                result.append(text[i])
                i += 1
                continue
            inner = text[i+2:end-2]
            pipe_idx = -1
            depth = 0
            for j, ch in enumerate(inner):
                if ch == '{' and inner[j:j+2] == '{{':
                    depth += 1
                elif ch == '}' and inner[j:j+2] == '}}':
                    depth -= 1
                elif ch == '|' and depth == 0:
                    pipe_idx = j
                    break
            if pipe_idx == -1:
                name = inner.strip()
                params = []
            else:
                name = inner[:pipe_idx]
                param_str = inner[pipe_idx+1:]
                params = _split_params(param_str)
            replacement = _process_known_template(name, params)
            result.append(replacement)
            i = end
        else:
            result.append(text[i])
            i += 1
    # Clean up leftover stray braces
    text = ''.join(result)

    # Protect <math>...</math> blocks from the blanket brace cleanup below.
    # The final regexes strip all {{ and }} from the text, but those patterns
    # appear legitimately inside LaTeX math content (e.g. \frac{\text{rise}}{\text{run}}).
    # Without protection, the }} between numerator and denominator gets removed.
    math_blocks = []
    def save_math(m):
        math_blocks.append(m.group(0))
        return f"\x00MATH{len(math_blocks)-1}\x00"
    text = re.sub(r'<math[^>]*>.*?</math>', save_math, text, flags=re.DOTALL)

    text = re.sub(r'^[\}\{]+', '', text)
    text = re.sub(r'[\}\{]+$', '', text)
    text = re.sub(r'\}\}', '', text)
    text = re.sub(r'\{\{', '', text)

    # Restore math blocks
    text = re.sub(r'\x00MATH(\d+)\x00', lambda m: math_blocks[int(m.group(1))], text)

    return text


def convert_file_refs(text):
    """Convert standalone File: and Image: references to Markdown images.
    Handles: File:Diagram.svg|Caption → ![Caption](/diagrams/wikipedia/Diagram.svg)"""
    def file_replacer(m):
        prefix = m.group(1)  # File or Image
        filename = m.group(2).strip()
        caption = m.group(3).strip() if m.group(3) else filename
        # Remove any size specifiers like |400px
        filename = re.sub(r'\|[0-9]+px', '', filename)
        caption = re.sub(r'\|[0-9]+px', '', caption)
        return '![' + caption + '](/diagrams/wikipedia/' + filename + ')'

    text = re.sub(
        r'^\[\[(File|Image):([^\|\]]+)(?:\|([^\]]*))?\]\]',
        file_replacer,
        text,
        flags=re.MULTILINE,
    )
    # Also handle bare File:/Image: without [[ ]]
    text = re.sub(
        r'^(File|Image):([^\|\n]+)(?:\|([^\n]*))?',
        file_replacer,
        text,
        flags=re.MULTILINE,
    )
    return text


def convert_wiki_tables(text):
    """Convert {| |} wiki table syntax to Markdown tables where possible.
    Complex tables (with colspan, rowspan, etc.) are stripped entirely."""
    # Strip complex tables that have colspan/rowspan or are clearly layout-only
    if re.search(r'colspan|rowspan|style=|[|]{|}\|', text):
        text = re.sub(
            r'\{\|.*?\|\}',
            '',
            text,
            flags=re.DOTALL,
        )
    # Attempt simple table conversion for basic {| ... |} syntax
    # For now, strip all remaining wiki table markup that wasn't caught
    text = re.sub(r'^\{\|.*', '', text, flags=re.MULTILINE)
    text = re.sub(r'^\|\}.*', '', text, flags=re.MULTILINE)
    text = re.sub(r'^\|[\+\!].*', '', text, flags=re.MULTILINE)
    # Convert simple |- row separators to blank lines
    text = re.sub(r'^\|-\s*$', '', text, flags=re.MULTILINE)
    # Convert | cell delimiter to | for simple rows
    text = re.sub(r'^\|(.*?)\|$', r'| \1 |', text, flags=re.MULTILINE)
    # Clean up remaining wiki table artifacts
    text = re.sub(r'\|\}', '', text)
    return text


def convert_links(text):
    """Convert [[Target|display]] → display and [[Target]] → Target."""
    # Process file refs BEFORE regular links so image syntax takes priority
    text = convert_file_refs(text)
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
    # Convert wiki tables before math to avoid mangling LaTeX inside cells
    text = convert_wiki_tables(text)
    # Process math (convert <math>...</math> to LaTeX)
    text = convert_math(text)
    # Convert <sup>/<sub> to LaTeX notation (remaining after template removal)
    text = re.sub(r'<sup>([^<]+)</sup>', r'\\(^{\1}\\)', text)
    text = re.sub(r'<sub>([^<]+)</sub>', r'\\(_{\1}\\)', text)
    # Remove entire <table> and <dl> blocks (too complex for clean conversion)
    text = re.sub(r'<table[^>]*>.*?</table>', '', text, flags=re.DOTALL)
    text = re.sub(r'<dl[^>]*>.*?</dl>', '', text, flags=re.DOTALL)
    # Remove any remaining HTML tags except <b>, <i>, <em>, <strong>
    text = re.sub(r'</(?!b>|i>|em>|strong>)[a-z]+\s*>', '', text)
    text = re.sub(r'<(?![/]?b>|[/]?i>|[/]?em>|[/]?strong>)[a-z]+[^>]*>', '', text, flags=re.DOTALL)
    # Convert links and formatting
    text = convert_links(text)
    text = convert_formatting(text)
    # Strip irrelevant sections
    text = re.sub(
        r'^#{2,3}\s*(See also|Notes|References|Further reading|'
        r'External links|Bibliography|Sources|Footnotes)\s*$.*',
        '', text, flags=re.MULTILINE | re.DOTALL | re.IGNORECASE,
    )
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
        # Fallback: if sections returned nothing, try full-page fetch
        if not sections_text.strip():
            print(f"  sections returned empty, falling back to full page...")
            wt = fetch_wikitext(wiki_title)
            if not wt:
                print(f"  [error] could not fetch page '{wiki_title}'", file=sys.stderr)
                return None
            sections_text = clean_wikitext(wt)
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

    for i, item in enumerate(batch):
        cid = item["concept_id"]
        wp = item["wikipedia"]
        secs = item.get("sections", [])
        print(f"\n[{i+1}/{len(batch)}] {cid} ← Wikipedia: {wp}" + (f" sections {secs}" if secs else ""))
        out = fetch_and_generate(cid, wp, secs)
        if out:
            rel_path = "wikipedia/" + cid + ".md"
            register_lesson(cid, rel_path)
        # Delay between lessons to avoid rate limiting
        if i < len(batch) - 1:
            time.sleep(2.0)


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
