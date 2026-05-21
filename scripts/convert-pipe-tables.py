#!/usr/bin/env python3
"""
Convert Wikipedia pipe-table syntax ({| ... |}) to Markdown/HTML.

This script handles three output formats based on table complexity:
  1. Simple grids (no colspan/rowspan) => Markdown pipe tables
  2. Tables with colspan/rowspan      => HTML <table> (valid in Markdown)
  3. Definition-list style tables     => HTML <dl> definitions

Usage:
  python3 scripts/convert-pipe-tables.py                # convert all files
  python3 scripts/convert-pipe-tables.py --report       # only report, don't modify
  python3 scripts/convert-pipe-tables.py --file <id>    # convert single file
"""

import os
import re
import sys

WIKI_DIR = "data/lessons/wikipedia"

# ---------- pipe-table parser ----------

def parse_table_attributes(line: str) -> str:
    """Extract table-level attributes from {| ... line."""
    m = re.match(r'\{\|\s*(.*)', line)
    if not m:
        return ''
    return m.group(1).strip()


def is_header_cell(line: str, pos: int) -> bool:
    """Check if the cell at pos is a header (starts with !)."""
    # Skip past any leading whitespace
    i = pos
    while i < len(line) and line[i] == ' ':
        i += 1
    return i < len(line) and line[i] == '!'


def split_on_separator(text: str, sep: str) -> list:
    """Split on separator, respecting LaTeX math, templates, and wiki links."""
    parts = []
    current = []
    i = 0
    in_math = False
    in_template = 0
    in_link = False

    while i < len(text):
        if text[i:i+2] == '\\(':
            in_math = True
            current.append(text[i])
            i += 1
            continue
        if text[i:i+2] == '\\)' and in_math:
            in_math = False
            current.append(text[i])
            i += 1
            continue
        if in_math:
            current.append(text[i])
            i += 1
            continue

        if text[i:i+2] == '{{':
            in_template += 1
            current.append(text[i])
            i += 1
            continue
        if text[i:i+2] == '}}' and in_template > 0:
            in_template -= 1
            current.append(text[i])
            i += 1
            continue
        if in_template > 0:
            current.append(text[i])
            i += 1
            continue

        if text[i:i+2] == '[[':
            in_link = True
            current.append(text[i])
            i += 1
            continue
        if text[i:i+2] == ']]' and in_link:
            in_link = False
            current.append(text[i])
            i += 1
            continue
        if in_link:
            current.append(text[i])
            i += 1
            continue

        # Check separator
        if text[i:i+len(sep)] == sep:
            parts.append(''.join(current).strip())
            current = []
            i += len(sep)
            continue

        current.append(text[i])
        i += 1

    if current:
        parts.append(''.join(current).strip())
    return parts


def split_cells(row_text: str) -> list:
    """
    Split a wiki-table row line into cell contents.
    Handles !! (header separator) and || (data separator).
    """
    stripped = row_text.lstrip()
    if not stripped:
        return []

    # Determine cell type from first character
    first_char = stripped[0]
    if first_char == '!':
        cell_type = 'h'
        content = stripped[1:]  # strip leading !
        sep = '!!'
    elif first_char == '|':
        cell_type = 'd'
        content = stripped[1:]  # strip leading |
        sep = '||'
    else:
        return []

    raw_cells = split_on_separator(content, sep)
    # Strip Wikipedia attributes from each cell
    cells = [(cell_type, c.strip()) for c in raw_cells]
    return cells


def strip_wiki_attributes(cell_text: str) -> str:
    """Remove Wikipedia-specific HTML attributes from cell text."""
    # Remove style="...", class="...", colspan="N", rowspan="N" from cell start
    cell_text = re.sub(
        r'\s*(?:style|class|align|valign|width|bgcolor|colspan|rowspan)="[^"]*"\s*',
        ' ', cell_text
    )
    cell_text = re.sub(r'\s+', ' ', cell_text).strip()
    return cell_text


def extract_cell_attributes(cell_text: str) -> dict:
    """Extract colspan/rowspan from cell text and return them separately."""
    attrs = {}
    m = re.search(r'colspan="(\d+)"', cell_text)
    if m:
        attrs['colspan'] = int(m.group(1))
    m = re.search(r'rowspan="(\d+)"', cell_text)
    if m:
        attrs['rowspan'] = int(m.group(1))
    return attrs


def has_span_attrs(cell_text: str) -> bool:
    """Check if cell has colspan or rowspan attributes."""
    return bool(re.search(r'(?:colspan|rowspan)="\d+"', cell_text))


def cell_to_html(cell_type: str, content: str) -> str:
    """Convert a single cell to HTML tag + content."""
    tag = 'th' if cell_type == 'h' else 'td'
    content = strip_wiki_attributes(content)
    # Extract and apply colspan/rowspan
    attrs = extract_cell_attributes(content)
    attr_str = ''
    if attrs.get('colspan'):
        attr_str += f' colspan="{attrs["colspan"]}"'
    if attrs.get('rowspan'):
        attr_str += f' rowspan="{attrs["rowspan"]}"'
    # Remove attributes from content
    content = re.sub(r'\s*(?:colspan|rowspan)="\d+"', '', content).strip()
    return f'    <{tag}{attr_str}>{content}</{tag}>'


def table_to_markdown(table_data: dict) -> str:
    """Convert parsed table to Markdown pipe table (simple, no spanning)."""
    parts = []
    header_row = None
    rows = table_data['rows']
    has_header = False

    # Determine if first row is header
    processed = []
    for row in rows:
        cells = row['cells']
        cell_types = [c[0] for c in cells]
        if all(t == 'h' for t in cell_types):
            has_header = True
        processed.append(cells)

    if not processed:
        return ''

    # Calculate column count
    max_cols = max(len(row) for row in processed)

    # Build Markdown table
    lines = []
    for ri, row in enumerate(processed):
        md_cells = []
        for cell in row:
            content = strip_wiki_attributes(cell[1])
            content = content.replace('\n', ' ').strip()
            md_cells.append(content)
        # Pad to max_cols
        while len(md_cells) < max_cols:
            md_cells.append('')
        lines.append('| ' + ' | '.join(md_cells) + ' |')

        if ri == 0 and has_header:
            # Separator row
            sep = '| ' + ' | '.join(['---'] * max_cols) + ' |'
            lines.append(sep)

    # If no header, add separator after first data row
    if not has_header and len(lines) > 0:
        sep = '| ' + ' | '.join(['---'] * max_cols) + ' |'
        lines.insert(1, sep)

    return '\n'.join(lines)


def table_to_html(table_data: dict) -> str:
    """Convert parsed table to HTML <table>."""
    parts = ['<table>']

    caption = table_data.get('caption', '')
    if caption:
        parts.append(f'  <caption>{caption}</caption>')

    for row in table_data['rows']:
        parts.append('  <tr>')
        for cell_type, content in row['cells']:
            parts.append(cell_to_html(cell_type, content))
        parts.append('  </tr>')

    parts.append('</table>')
    return '\n'.join(parts)


def is_definition_table(rows: list) -> bool:
    """
    Heuristic: if table has exactly 2 columns and no header,
    and cells contain LaTeX formulas paired with descriptions, it's a def-list.
    """
    if not rows:
        return False
    # Check all rows have exactly 2 cells
    for row in rows:
        if len(row['cells']) != 2:
            return False
        # First cell should contain LaTeX
    # No headers
    types = [c[0] for row in rows for c in row['cells']]
    if 'h' in types:
        return False
    return True


def table_to_definition_list(table_data: dict) -> str:
    """Convert property-definition table to HTML <dl>."""
    parts = ['<dl>']
    for row in table_data['rows']:
        cells = row['cells']
        if len(cells) >= 2:
            term = strip_wiki_attributes(cells[0][1]).strip()
            desc = strip_wiki_attributes(cells[1][1]).strip()
            parts.append(f'  <dt>{term}</dt>')
            parts.append(f'  <dd>{desc}</dd>')
    parts.append('</dl>')
    return '\n'.join(parts)


def classify_table(table_text: str) -> str:
    """Classify a table and return the desired output format."""
    has_spans = bool(re.search(r'(?:colspan|rowspan)="\d+"', table_text))
    has_wikitable = 'class="wikitable"' in table_text

    # Count rows and columns
    row_lines = [l for l in table_text.split('\n') if l.strip().startswith('|-')]
    data_lines = [l for l in table_text.split('\n') if l.strip().startswith('|') and not l.strip().startswith('|}')]

    # Check if it's a definition-style table (no wikitable, 2-col, pattern)
    if not has_wikitable and not has_spans:
        # Check for 2-column pattern with LaTeX in first col
        cell_count = 0
        for l in table_text.split('\n'):
            s = l.strip()
            if s.startswith('|') and not s.startswith('|}'):
                cell_count += 1
        if cell_count <= 20:  # reasonable number of cells
            pass  # Might be definition list

    if has_spans:
        return 'html'
    return 'markdown'


def parse_table(text: str, start: int) -> tuple:
    """
    Parse a complete pipe table starting at {|.
    Returns (table_data_dict, end_position).
    """
    i = start
    assert text[i:i+2] == '{|', f"Not a table at {start}: {text[start:start+20]}"

    # Parse table-level attributes
    end_of_first = text.find('\n', i)
    if end_of_first < 0:
        return None, i + 2
    table_attrs = text[i:end_of_first].strip()

    i = end_of_first + 1
    rows = []
    current_row = []
    caption = ''
    in_row = True  # Content after {| line IS the first row
    continuation = False
    seen_row_sep = False  # Track if we've seen |- yet

    while i < len(text):
        # End of table
        if text[i:i+2] == '|}':
            if current_row:
                rows.append({'cells': current_row})
            result = {
                'attrs': table_attrs,
                'caption': caption,
                'rows': rows,
            }
            return result, i + 2

        # Newline
        if text[i] == '\n':
            i += 1
            continue

        # Caption
        if text[i:i+2] == '|+':
            end_line = text.find('\n', i)
            caption = text[i+2:end_line].strip()
            i = end_line + 1
            continue

        # Row separator
        if text[i:i+2] == '|-':
            if current_row:
                rows.append({'cells': current_row})
                current_row = []
            in_row = True
            seen_row_sep = True
            continuation = False
            i += 2
            continue

        # Cell content (| or ! at start of line)
        if (text[i] == '|' or text[i] == '!') and not continuation:
            end_line = text.find('\n', i)
            if end_line < 0:
                end_line = len(text)
            line = text[i:end_line]

            # Skip if this is |} (end of table)
            stripped = line.strip()
            if stripped == '|}':
                continue

            # Check if line has inline multi-cell separators
            cells = split_cells(line)

            if not current_row and not seen_row_sep:
                # First row of cells (headers or data before |-)
                current_row = cells
            else:
                current_row.extend(cells)

            i = end_line + 1
            continue

        # Continuation line (continues last cell content)
        if current_row:
            end_line = text.find('\n', i)
            if end_line < 0:
                end_line = len(text)
            line = text[i:end_line].strip()

            if line.startswith('|') or line.startswith('!'):
                cells = split_cells(text[i:end_line])
                current_row.extend(cells)
                i = end_line + 1
                continue
            elif line and not line.startswith('|'):
                # Continuation of last cell
                last = current_row[-1]
                current_row[-1] = (last[0], last[1] + '\n' + line)
                i = end_line + 1
                continue

        i += 1

    return None, len(text)


def convert_table(text: str, start: int) -> tuple:
    """
    Convert a pipe table starting at start to Markdown/HTML.
    Returns (converted_text, end_position).
    """
    table_data, end = parse_table(text, start)
    if not table_data:
        return text[start:end], end

    rows = table_data['rows']
    if not rows:
        return '', end

    # Check for nested tables (if any cell or caption contains another {|)
    full_text = text[start:end]
    # The full_text ends with |} (the current table's closing). Look for {|
    # inside the table body (before the final |}) but NOT inside LaTeX math.
    # Exclude the first {| (the table opener itself).
    inner = full_text[:full_text.rfind('|}')] if '|}' in full_text else full_text
    has_nested = False
    if '{|' in inner:
        # Find all {| after the first one
        first_brace = inner.index('{|')
        after_first = inner[first_brace + 2:]
        if '{|' in after_first:
            for m in re.finditer(re.escape('{|'), after_first):
                abs_pos = start + first_brace + 2 + m.start()
                if not inside_math(text, abs_pos):
                    has_nested = True
                    break
    if has_nested:
        result = table_to_html(table_data)
        return result, end

    # Classify and convert
    classification = classify_table(full_text)

    if classification == 'html' or has_span_in_rows(rows):
        result = table_to_html(table_data)
    elif is_definition_table(rows):
        result = table_to_definition_list(table_data)
    else:
        result = table_to_markdown(table_data)

    return result.strip(), end


def has_span_in_rows(rows: list) -> bool:
    """Check if any cell in rows has colspan/rowspan."""
    for row in rows:
        for cell_type, content in row['cells']:
            if has_span_attrs(content):
                return True
    return False


def caption_index(text: str) -> int:
    """Find where caption ends."""
    # Skip past {| line
    idx = text.find('\n')
    return idx + 1 if idx > 0 else 0


def inside_math(text: str, pos: int) -> bool:
    """Check if pos is inside LaTeX math mode \(...\) or \[...\]."""
    before = text[:pos]
    # Check for \(...\) pair
    last_inline_open = before.rfind('\\(')
    last_inline_close = before.rfind('\\)')
    if last_inline_open > last_inline_close:
        return True
    # Check for \[...\] pair
    last_display_open = before.rfind('\\[')
    last_display_close = before.rfind('\\]')
    if last_display_open > last_display_close:
        return True
    return False


def find_table_starts(text: str) -> list:
    """Find all {| that start actual pipe tables (not inside math or templates)."""
    positions = []
    i = 0
    in_template = 0
    in_math = False
    in_display_math = False

    while i < len(text):
        # Track template depth
        if text[i:i+2] == '{{':
            in_template += 1
            i += 2
            continue
        if text[i:i+2] == '}}' and in_template > 0:
            in_template -= 1
            i += 2
            continue

        # Track math mode
        if text[i:i+2] == '\\(' and not in_template:
            in_math = True
            i += 2
            continue
        if text[i:i+2] == '\\)' and in_math:
            in_math = False
            i += 2
            continue
        if text[i:i+2] == '\\[' and not in_template:
            in_display_math = True
            i += 2
            continue
        if text[i:i+2] == '\\]' and in_display_math:
            in_display_math = False
            i += 2
            continue

        # Check for {| - only count if not inside math or template
        if text[i:i+2] == '{|' and not in_math and not in_display_math and in_template == 0:
            positions.append(i)

        i += 1

    return positions


def convert_pipe_tables(text: str) -> str:
    """Find and convert all pipe tables in text."""
    # Find actual table start positions (not inside math/templates)
    table_starts = find_table_starts(text)
    if not table_starts:
        return text

    result = []
    prev_end = 0

    for pos in table_starts:
        # Add text before this table
        result.append(text[prev_end:pos])

        # Parse and convert this table
        converted, end = convert_table(text, pos)
        result.append(converted.strip())
        result.append('\n\n')
        prev_end = end

    # Add remaining text
    result.append(text[prev_end:])

    return ''.join(result)


def process_file(filepath: str, dry_run: bool = False) -> dict:
    """Process a single file, converting all pipe tables."""
    with open(filepath) as f:
        original = f.read()

    # Quick check for pipe tables
    if '{|' not in original:
        return {'file': os.path.basename(filepath), 'changed': False, 'tables_found': 0}

    # Count only actual table starts (not inside math/templates)
    table_starts = find_table_starts(original)
    count_before = len(table_starts)
    if count_before == 0:
        return {'file': os.path.basename(filepath), 'changed': False, 'tables_found': 0}

    text = convert_pipe_tables(original)

    # Clean up excessive blank lines
    text = re.sub(r'\n{3,}', '\n\n', text)

    changed = text != original

    if changed and not dry_run:
        with open(filepath, 'w') as f:
            f.write(text)

    return {
        'file': os.path.basename(filepath),
        'changed': changed,
        'tables_found': count_before,
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
    total_tables = 0

    for fname in md_files:
        filepath = os.path.join(WIKI_DIR, fname)
        result = process_file(filepath, dry_run=dry_run)
        if result['changed']:
            total_changed += 1
            if not dry_run:
                print(f"  {fname}: converted {result['tables_found']} table(s)")
        total_tables += result['tables_found']

    if dry_run:
        print(f"\n=== REPORT: {len(md_files)} files, {total_tables} pipe tables found ===")
    else:
        print(f"\n=== Converted tables in {total_changed}/{len(md_files)} files ===")
        print(f"    Total pipe tables converted: {total_tables}")
        print("\nDone.")


if __name__ == "__main__":
    main()
