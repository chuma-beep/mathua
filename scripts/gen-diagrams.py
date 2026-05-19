#!/usr/bin/env python3
"""Generate SVG math diagrams for lesson content."""

import math
import os

OUT = os.path.join(os.path.dirname(__file__), '..', 'web', 'next-app', 'public', 'diagrams', 'algebrica')

def svg_head(w, h):
    return f'<svg width="{w}px" height="{h}px" viewBox="0 0 {w} {h}" xmlns="http://www.w3.org/2000/svg">'

def arrow(x1, y1, x2, y2, color='#333', width=1.5):
    """Line with arrowhead."""
    dx = x2 - x1
    dy = y2 - y1
    angle = math.atan2(dy, dx)
    head_len = 8
    return f'''<line x1="{x1}" y1="{y1}" x2="{x2}" y2="{y2}" stroke="{color}" stroke-width="{width}" stroke-linecap="round"/>
<polygon points="{x2},{y2} {x2-head_len*math.cos(angle-0.4)},{y2-head_len*math.sin(angle-0.4)} {x2-head_len*math.cos(angle+0.4)},{y2-head_len*math.sin(angle+0.4)}" fill="{color}"/>'''

def unit_circle_svg():
    """Unit circle with sine/cosine projections at angle θ."""
    cx, cy, r = 200, 200, 140
    theta = math.radians(50)
    px = cx + r * math.cos(theta)
    py = cy - r * math.sin(theta)

    lines = [svg_head(400, 400), '<defs><marker id="dot" viewBox="0 0 10 10" refX="5" refY="5"><circle cx="5" cy="5" r="5" fill="#333"/></marker></defs>']
    lines.append(f'<rect width="400" height="400" fill="white"/>')
    # Axes with arrows
    lines.append(arrow(30, cy, 370, cy, '#999', 1))
    lines.append(arrow(cx, 370, cx, 30, '#999', 1))
    # Unit circle
    lines.append(f'<circle cx="{cx}" cy="{cy}" r="{r}" fill="none" stroke="#333" stroke-width="2"/>')
    # Radius line to P
    lines.append(f'<line x1="{cx}" y1="{cy}" x2="{px}" y2="{py}" stroke="#333" stroke-width="2"/>')
    # Cosine projection (x-axis)
    lines.append(f'<line x1="{px}" y1="{cy}" x2="{px}" y2="{py}" stroke="#e74c3c" stroke-width="1.5" stroke-dasharray="6,4"/>')
    # Sine projection (y-axis)
    lines.append(f'<line x1="{cx}" y1="{py}" x2="{px}" y2="{py}" stroke="#2980b9" stroke-width="1.5" stroke-dasharray="6,4"/>')
    # Point P
    lines.append(f'<circle cx="{px}" cy="{py}" r="5" fill="#333"/>')
    # Angle arc
    large = 1 if theta > math.pi else 0
    end_x = cx + 30 * math.cos(theta)
    end_y = cy - 30 * math.sin(theta)
    lines.append(f'<path d="M {cx+30} {cy} A 30 30 0 {large} 0 {end_x} {end_y}" fill="none" stroke="#999" stroke-width="1"/>')
    # Labels
    lines.append(f'<text x="{cx+5}" y="{cy-5}" font-size="14" font-family="serif" fill="#999">O</text>')
    lines.append(f'<text x="{370}" y="{cy+15}" font-size="14" font-family="serif" fill="#999">x</text>')
    lines.append(f'<text x="{cx+8}" y="{35}" font-size="14" font-family="serif" fill="#999">y</text>')
    lines.append(f'<text x="{px+8}" y="{py-8}" font-size="15" font-family="serif" fill="#333">P</text>')
    lines.append(f'<text x="{end_x+10}" y="{end_y-5}" font-size="13" font-family="serif" fill="#999">θ</text>')
    lines.append(f'<text x="{390}" y="{cy-8}" text-anchor="end" font-size="13" font-family="serif" fill="#e74c3c">cos(θ)</text>')
    lines.append(f'<text x="{cx+10}" y="{20}" font-size="13" font-family="serif" fill="#2980b9">sin(θ)</text>')
    lines.append('</svg>')
    return '\n'.join(lines)


def sine_cosine_graph_svg():
    """Sine and cosine function graphs side by side."""
    w, h = 500, 220
    margin = 40
    gx, gy = margin, h // 2
    gw = w - 2 * margin
    gh = 90

    lines = [svg_head(w, h)]
    lines.append(f'<rect width="{w}" height="{h}" fill="white"/>')

    # Axes
    lines.append(arrow(gx, gy, w - margin + 10, gy, '#999', 1))
    lines.append(arrow(gx, margin - 10, gx, h - margin + 10, '#999', 1))

    # X-axis labels
    for i, label in enumerate(['0', 'π/2', 'π', '3π/2', '2π']):
        x = gx + gw * i / 4
        lines.append(f'<text x="{x}" y="{gy+15}" text-anchor="middle" font-size="10" font-family="serif" fill="#999">{label}</text>')
    lines.append(f'<text x="{w-margin+5}" y="{gy+15}" font-size="11" font-family="serif" fill="#999">x</text>')
    lines.append(f'<text x="{gx+5}" y="{margin-5}" font-size="11" font-family="serif" fill="#999">y</text>')

    # Y-axis labels
    for val, label in [(1, '1'), (-1, '-1')]:
        y = gy - val * gh
        lines.append(f'<text x="{gx-8}" y="{y+4}" text-anchor="end" font-size="10" font-family="serif" fill="#999">{label}</text>')
        lines.append(f'<line x1="{gx-4}" y1="{y}" x2="{gx+4}" y2="{y}" stroke="#999" stroke-width="0.5"/>')

    # Sine wave
    sine_pts = []
    for i in range(201):
        t = (i / 200) * 2 * math.pi
        x = gx + (t / (2 * math.pi)) * gw
        y = gy - math.sin(t) * gh
        sine_pts.append(f'{x},{y}')
    lines.append(f'<polyline points="{" ".join(sine_pts)}" fill="none" stroke="#2980b9" stroke-width="2"/>')
    lines.append(f'<text x="{w-margin}" y="{gy-gh+15}" text-anchor="end" font-size="11" font-family="serif" fill="#2980b9">sin(x)</text>')

    # Cosine wave
    cos_pts = []
    for i in range(201):
        t = (i / 200) * 2 * math.pi
        x = gx + (t / (2 * math.pi)) * gw
        y = gy - math.cos(t) * gh
        cos_pts.append(f'{x},{y}')
    lines.append(f'<polyline points="{" ".join(cos_pts)}" fill="none" stroke="#e74c3c" stroke-width="2"/>')
    lines.append(f'<text x="{w-margin}" y="{gy+gh-10}" text-anchor="end" font-size="11" font-family="serif" fill="#e74c3c">cos(x)</text>')

    lines.append('</svg>')
    return '\n'.join(lines)


if __name__ == '__main__':
    os.makedirs(OUT, exist_ok=True)

    svgs = {
        'unit-circle-sine-cosine.svg': unit_circle_svg(),
        'sine-cosine-graph.svg': sine_cosine_graph_svg(),
    }

    for name, content in svgs.items():
        path = os.path.join(OUT, name)
        with open(path, 'w') as f:
            f.write(content)
        print(f'Generated {path}')
