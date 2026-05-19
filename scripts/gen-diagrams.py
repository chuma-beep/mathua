#!/usr/bin/env python3
"""Generate SVG math diagrams for lesson content.

Usage: python3 scripts/gen-diagrams.py
Output: web/next-app/public/diagrams/algebrica/*.svg
"""

import math
import os

OUT = os.path.join(os.path.dirname(__file__), '..', 'web', 'next-app', 'public', 'diagrams', 'algebrica')
DATA = os.path.join(os.path.dirname(__file__), '..', 'data', 'lessons', 'algebrica')

# Color palette matching existing diagrams
C_MAIN = '#333333'
C_AXIS = '#BABABA'
C_FILL = '#E6E6E6'
C_TEXT = '#636363'
C_SINE = '#2980b9'
C_COS = '#e74c3c'
C_TAN = '#27ae60'

# ─── SVG helpers ────────────────────────────────────────────────────────────

def head(w, h):
    return f'<svg width="{w}px" height="{h}px" viewBox="0 0 {w} {h}" xmlns="http://www.w3.org/2000/svg">'

def bg(w, h):
    return f'<rect width="{w}" height="{h}" fill="white"/>'

def text(x, y, t, size=13, color=C_MAIN, anchor='start', family='serif', style='normal'):
    a = f' text-anchor="{anchor}"' if anchor != 'start' else ''
    s = f' font-style="{style}"' if style != 'normal' else ''
    return f'<text x="{x}" y="{y}" font-size="{size}" font-family="{family}" fill="{color}"{a}{s}>{t}</text>'

def line(x1, y1, x2, y2, color=C_AXIS, width=1, dash=''):
    d = f' stroke-dasharray="{dash}"' if dash else ''
    return f'<line x1="{x1}" y1="{y1}" x2="{x2}" y2="{y2}" stroke="{color}" stroke-width="{width}" stroke-linecap="round"{d}/>'

def circle(cx, cy, r, color=C_MAIN, fill='none', width=2):
    f = f' fill="{fill}"' if fill else ''
    return f'<circle cx="{cx}" cy="{cy}" r="{r}"{f} stroke="{color}" stroke-width="{width}"/>'

def dot(cx, cy, r=3.5, color=C_MAIN):
    return f'<circle cx="{cx}" cy="{cy}" r="{r}" fill="{color}"/>'

def polyline(pts, color=C_MAIN, width=2, fill='none'):
    p = ' '.join(f'{x},{y}' for x, y in pts)
    f = f' fill="{fill}"' if fill else ''
    return f'<polyline points="{p}"{f} stroke="{color}" stroke-width="{width}"/>'

def path(d, color=C_MAIN, width=1.5, fill='none'):
    f = f' fill="{fill}"' if fill else ''
    return f'<path d="{d}"{f} stroke="{color}" stroke-width="{width}"/>'

def arrow_tip(x, y, angle, color=C_AXIS):
    hl = 8
    return (f'<polygon points="{x},{y} '
            f'{x-hl*math.cos(angle-0.4)},{y-hl*math.sin(angle-0.4)} '
            f'{x-hl*math.cos(angle+0.4)},{y-hl*math.sin(angle+0.4)}" fill="{color}"/>')

def h_arrow(x1, x2, y, color=C_AXIS, width=1):
    parts = [line(x1, y, x2, y, color, width)]
    if x2 > x1:
        parts.append(arrow_tip(x2, y, 0, color))
    return '\n'.join(parts)

def v_arrow(y1, y2, x, color=C_AXIS, width=1):
    parts = [line(x, y1, x, y2, color, width)]
    if y2 < y1:
        parts.append(arrow_tip(x, y2, -math.pi/2, color))
    return '\n'.join(parts)

def arc(cx, cy, r, a1, a2, color=C_AXIS, width=1):
    start_x = cx + r * math.cos(a1)
    start_y = cy - r * math.sin(a1)
    end_x = cx + r * math.cos(a2)
    end_y = cy - r * math.sin(a2)
    large = 1 if abs(a2 - a1) > math.pi else 0
    sweep = 1 if a2 > a1 else 0
    return f'<path d="M {start_x} {start_y} A {r} {r} 0 {large} {sweep} {end_x} {end_y}" fill="none" stroke="{color}" stroke-width="{width}"/>'

# ─── Diagram generators ─────────────────────────────────────────────────────

def unit_circle_labeled():
    """Labeled unit circle with quadrant angles and (cos, sin) coordinates."""
    cx, cy, r = 250, 250, 160
    lines = [head(500, 500), bg(500, 500)]
    lines.append(h_arrow(40, 460, cy, C_AXIS, 1.5))
    lines.append(v_arrow(cy+10, 40, cx, C_AXIS, 1.5))
    lines.append(circle(cx, cy, r, C_MAIN, width=2))
    lines.append(text(cx, cy+8, 'O', 15, C_TEXT))
    lines.append(text(458, cy+15, 'x', 14, C_TEXT))
    lines.append(text(cx+10, 40, 'y', 14, C_TEXT))

    angles = [0, 30, 45, 60, 90, 120, 135, 150, 180, 210, 225, 240, 270, 300, 315, 330]
    labels = ['0', 'π/6', 'π/4', 'π/3', 'π/2', '2π/3', '3π/4', '5π/6',
              'π', '7π/6', '5π/4', '4π/3', '3π/2', '5π/3', '7π/4', '11π/6']
    cos_vals = ['1', '√3/2', '√2/2', '1/2', '0', '-1/2', '-√2/2', '-√3/2',
                '-1', '-√3/2', '-√2/2', '-1/2', '0', '1/2', '√2/2', '√3/2']
    sin_vals = ['0', '1/2', '√2/2', '√3/2', '1', '√3/2', '√2/2', '1/2',
                '0', '-1/2', '-√2/2', '-√3/2', '-1', '-√3/2', '-√2/2', '-1/2']

    for i, deg in enumerate(angles):
        rad = math.radians(deg)
        px = cx + r * math.cos(rad)
        py = cy - r * math.sin(rad)
        lines.append(line(cx, cy, px, py, C_AXIS, 0.5))
        lines.append(dot(px, py, 3, C_MAIN))

        lx = cx + (r + 22) * math.cos(rad)
        ly = cy - (r + 22) * math.sin(rad)
        align = 'middle'
        ox = 0 if abs(math.cos(rad)) < 0.15 else (-8 if math.cos(rad) > 0 else 8)
        if abs(math.cos(rad)) < 0.15:
            ox = 0
        elif math.cos(rad) > 0:
            ox = 10
        else:
            ox = -10
        oy = 0 if abs(math.sin(rad)) < 0.15 else (4 if math.sin(rad) > 0 else -4)
        lines.append(text(lx+ox, ly+oy, labels[i], 10, C_TEXT, 'middle'))

    lines.append(text(cx, cy-r-30, 'Unit Circle with Notable Angles', 14, C_MAIN, 'middle'))
    lines.append('</svg>')
    return '\n'.join(lines)


def unit_circle_sincos():
    """Unit circle with sine/cosine projections at angle θ (existing)."""
    cx, cy, r = 200, 200, 140
    theta = math.radians(50)
    px = cx + r * math.cos(theta)
    py = cy - r * math.sin(theta)

    lines = [head(400, 400), bg(400, 400)]
    lines.append(h_arrow(30, 370, cy, C_AXIS, 1))
    lines.append(v_arrow(cy+10, 30, cx, C_AXIS, 1))
    lines.append(circle(cx, cy, r, C_MAIN, width=2))
    lines.append(line(cx, cy, px, py, C_MAIN, 2))
    lines.append(line(px, cy, px, py, C_COS, 1.5, '6,4'))
    lines.append(line(cx, py, px, py, C_SINE, 1.5, '6,4'))
    lines.append(dot(px, py, 5))
    lines.append(arc(cx, cy, 30, 0, theta, C_AXIS, 1))
    lines.append(text(cx+5, cy-5, 'O', 14, C_TEXT))
    lines.append(text(368, cy+15, 'x', 14, C_TEXT))
    lines.append(text(cx+8, 32, 'y', 14, C_TEXT))
    lines.append(text(px+8, py-8, 'P', 15, C_MAIN))
    lines.append(text(cx+30+10, cy-30*math.sin(theta)-5, 'θ', 13, C_TEXT))
    lines.append(text(390, cy-8, 'cos(θ)', 13, C_COS, 'end'))
    lines.append(text(cx+10, 20, 'sin(θ)', 13, C_SINE))
    lines.append('</svg>')
    return '\n'.join(lines)


def unit_circle_tangent():
    """Unit circle showing tangent as length ST on tangent line x=1."""
    cx, cy, r = 200, 200, 140
    theta = math.radians(40)
    px = cx + r * math.cos(theta)
    py = cy - r * math.sin(theta)

    # Tangent line x=1 at point S(1,0)
    # S is at x=1 relative to unit circle, so at (cx+r, cy)
    Sx = cx + r
    Sy = cy
    # Line through O and P extended to x = cx+r
    if px != cx:
        t = (Sx - cx) / (px - cx)
        Ty = cy + t * (py - cy)
    else:
        Ty = cy - 100

    lines = [head(450, 400), bg(450, 400)]
    lines.append(h_arrow(30, 420, cy, C_AXIS, 1))
    lines.append(v_arrow(cy+10, 30, cx, C_AXIS, 1))
    lines.append(circle(cx, cy, r, C_MAIN, width=2))
    # Tangent vertical line at x = cx+r
    lines.append(line(Sx, 50, Sx, 350, C_AXIS, 1, '4,4'))
    # Line O->P extended to tangent
    lines.append(line(cx, cy, Sx, Ty, C_MAIN, 1.5))
    # Tangent segment ST
    lines.append(line(Sx, Sy, Sx, Ty, C_TAN, 2.5))
    # Points
    lines.append(dot(px, py, 4))
    lines.append(dot(Sx, Sy, 4))
    lines.append(dot(Sx, Ty, 4))
    # Labels
    lines.append(text(cx+5, cy-5, 'O', 14, C_TEXT))
    lines.append(text(418, cy+15, 'x', 14, C_TEXT))
    lines.append(text(cx+8, 32, 'y', 14, C_TEXT))
    lines.append(text(px+6, py-6, 'P', 13, C_MAIN))
    lines.append(text(Sx+6, Sy-6, 'S', 13, C_MAIN))
    lines.append(text(Sx+6, Ty+4, 'T', 13, C_MAIN))
    lines.append(text(Sx+10, (Sy+Ty)//2, 'tan(θ)', 13, C_TAN))
    lines.append(text(410, 60, f'x = 1', 11, C_TEXT, 'end'))
    lines.append('</svg>')
    return '\n'.join(lines)


def right_triangle_unit_circle():
    """Right triangle inscribed in unit circle for Pythagorean identity."""
    cx, cy, r = 220, 220, 150
    theta = math.radians(35)
    px = cx + r * math.cos(theta)
    py = cy - r * math.sin(theta)
    rx = px
    ry = cy

    lines = [head(440, 440), bg(440, 440)]
    lines.append(h_arrow(40, 410, cy, C_AXIS, 1))
    lines.append(v_arrow(cy+10, 30, cx, C_AXIS, 1))
    lines.append(circle(cx, cy, r, C_MAIN, width=2))
    # Right triangle O-R-P
    lines.append(line(cx, cy, px, py, C_MAIN, 2))
    lines.append(line(cx, cy, rx, ry, C_SINE, 2))
    lines.append(line(rx, ry, px, py, C_COS, 2))
    # Right angle marker
    s = 10
    lines.append(polyline([(rx, ry+s), (rx+s, ry+s), (rx+s, ry)], C_AXIS, 1))
    # Points
    lines.append(dot(cx, cy, 4))
    lines.append(dot(px, py, 4))
    lines.append(dot(rx, ry, 4))
    # Labels
    lines.append(text(cx+5, cy-5, 'O', 14, C_TEXT))
    lines.append(text(px+6, py-6, 'P', 13, C_MAIN))
    lines.append(text(rx+6, ry+5, 'R', 13, C_MAIN))
    lines.append(text(408, cy+15, 'x', 14, C_TEXT))
    lines.append(text(cx+8, 32, 'y', 14, C_TEXT))
    # Side labels
    mx_o = (cx + px) / 2
    my_o = (cy + py) / 2
    lines.append(text(mx_o+8, my_o, '1', 14, C_MAIN))
    lines.append(text((cx+rx)/2, cy+18, 'cos(θ)', 13, C_SINE, 'middle'))
    lines.append(text(px+10, (py+cy)/2, 'sin(θ)', 13, C_COS))
    lines.append(text(cx, cy-r-20, 'Right triangle inscribed in unit circle', 13, C_MAIN, 'middle'))
    lines.append('</svg>')
    return '\n'.join(lines)


def reference_angles():
    """Reference angles shown on unit circle in four quadrants."""
    cx, cy, r = 250, 250, 170
    lines = [head(500, 500), bg(500, 500)]
    lines.append(h_arrow(30, 470, cy, C_AXIS, 1.5))
    lines.append(v_arrow(cy+10, 30, cx, C_AXIS, 1.5))
    lines.append(circle(cx, cy, r, C_MAIN, width=2))

    angles = [
        (50, 'θ', C_SINE),      # Q1
        (130, 'θ', C_COS),       # Q2
        (230, 'θ', C_TAN),       # Q3
        (310, 'θ', C_SINE),      # Q4
    ]
    ref_angles = [50, 50, 50, 50]

    for (deg, label, color), ref_deg in zip(angles, ref_angles):
        rad = math.radians(deg)
        px = cx + r * math.cos(rad)
        py = cy - r * math.sin(rad)
        lines.append(line(cx, cy, px, py, C_MAIN, 1.5))
        lines.append(dot(px, py, 4, color))

        # Reference angle arc
        if deg <= 90:
            a1, a2 = 0, rad
        elif deg <= 180:
            ref_rad = math.radians(180 - deg)
            a1, a2 = rad, rad + ref_rad
        elif deg <= 270:
            ref_rad = math.radians(deg - 180)
            a1, a2 = rad - ref_rad, rad
        else:
            ref_rad = math.radians(360 - deg)
            a1, a2 = rad, rad - ref_rad

        lines.append(arc(cx, cy, 25, a1, a2, color, 1.5))

        # Dashed line to x-axis
        lines.append(line(px, py, px, cy, color, 1, '4,4'))

        # Label
        lx = px + (8 if px > cx else -8)
        ly = py + (4 if py < cy else -4)
        lines.append(text(lx, ly, label, 12, color))

    # Quadrant labels
    lines.append(text(cx+r+15, cy-r+10, 'Quadrant I', 11, C_TEXT))
    lines.append(text(cx-r-15, cy-r+10, 'Quadrant II', 11, C_TEXT, 'end'))
    lines.append(text(cx-r-15, cy+r-10, 'Quadrant III', 11, C_TEXT, 'end'))
    lines.append(text(cx+r+15, cy+r-10, 'Quadrant IV', 11, C_TEXT))

    lines.append(text(cx, cy-r-30, 'Reference Angles on the Unit Circle', 14, C_MAIN, 'middle'))
    lines.append('</svg>')
    return '\n'.join(lines)


def right_triangle_trig():
    """Labeled right triangle with sides opposite/adjacent/hypotenuse."""
    w, h = 350, 300
    # Triangle points: right angle at bottom-left
    ax, ay = 50, 250  # right angle vertex
    bx, by = 250, 250  # adjacent vertex
    cx, cy = 50, 80    # opposite vertex
    theta = math.atan2(ay - cy, bx - ax)  # at the right angle

    lines = [head(w, h), bg(w, h)]
    # Triangle
    lines.append(line(ax, ay, bx, by, C_MAIN, 2.5))
    lines.append(line(ax, ay, cx, cy, C_MAIN, 2.5))
    lines.append(line(bx, by, cx, cy, C_MAIN, 2.5))
    # Right angle marker
    s = 15
    lines.append(polyline([(ax+s, ay), (ax+s, ay-s), (ax, ay-s)], C_AXIS, 1.5))
    # Points
    lines.append(dot(ax, ay, 4))
    lines.append(dot(bx, by, 4))
    lines.append(dot(cx, cy, 4))
    # Labels
    lines.append(text(ax-15, ay+20, 'A (θ)', 13, C_MAIN))
    lines.append(text(bx+10, ay+20, 'B', 13, C_MAIN))
    lines.append(text(cx-15, cy-10, 'C', 13, C_MAIN))
    # Side labels
    lines.append(text((ax+bx)//2, ay+18, 'adjacent (cos)', 12, C_COS, 'middle'))
    lines.append(text(ax-15, (ay+cy)//2, 'opposite (sin)', 12, C_SINE, 'middle'))
    lines.append(text((bx+cx)//2 + 10, (by+cy)//2, 'hypotenuse (1)', 12, C_MAIN))
    # Angle arc
    lines.append(arc(ax, ay, 30, 0, theta, C_AXIS, 1.5))
    lines.append(text(ax+30+8, ay-30*math.sin(theta)+5, 'θ', 13, C_TEXT))
    lines.append('</svg>')
    return '\n'.join(lines)


def pythagorean_theorem():
    """Right triangle with squares on each side."""
    w, h = 400, 400
    # Triangle: right angle at bottom-left
    ax, ay = 40, 280  # right angle
    bx, by = 280, 280  # horizontal leg end
    cx, cy = 40, 100   # vertical leg end
    a_len = bx - ax  # adjacent length
    o_len = ay - cy  # opposite length

    lines = [head(w, h), bg(w, h)]
    # Triangle
    lines.append(line(ax, ay, bx, by, C_MAIN, 2.5))
    lines.append(line(ax, ay, cx, cy, C_MAIN, 2.5))
    lines.append(line(bx, by, cx, cy, C_MAIN, 2.5))
    # Right angle marker
    s = 15
    lines.append(polyline([(ax+s, ay), (ax+s, ay-s), (ax, ay-s)], C_AXIS, 1.5))
    # Square on adjacent (bottom)
    adj = bx - ax
    lines.append(polyline([(ax, ay), (bx, ay), (bx, ay-adj), (ax, ay-adj)], C_COS, 1.5))
    lines.append(text((ax+bx)//2, ay-adj//2, 'a²', 14, C_COS, 'middle'))
    # Square on opposite (left)
    opp = ay - cy
    lines.append(polyline([(ax, ay), (ax, cy), (ax+opp, cy), (ax+opp, ay)], C_SINE, 1.5))
    lines.append(text(ax+opp//2, (ay+cy)//2, 'b²', 14, C_SINE, 'middle'))
    # Square on hypotenuse
    # We need a rotated square on the hypotenuse
    hx, hy = bx - cx, by - cy
    hlen = math.hypot(hx, hy)
    nx, ny = -hy/hlen * hlen, hx/hlen * hlen  # perpendicular
    sq = [(cx, cy), (bx, by), (bx+nx, by+ny), (cx+nx, cy+ny)]
    lines.append(polyline(sq, C_TAN, 1.5))
    lines.append(text((cx+bx+nx)/2 - 5, (cy+by+ny)/2 + 5, 'c²', 14, C_TAN, 'middle'))

    lines.append(text(w//2, h-15, 'Pythagorean Theorem: a² + b² = c²', 14, C_MAIN, 'middle'))
    lines.append('</svg>')
    return '\n'.join(lines)


def law_of_sines():
    """Triangle with circumcircle for Law of Sines."""
    w, h = 400, 350
    cx, cy = 200, 180
    r = 140
    # Triangle inscribed in circle
    angles = [math.radians(25), math.radians(75), math.radians(80)]
    pts = []
    for i, a in enumerate(angles):
        angle = sum(angles[:i+1])
        pts.append((cx + r * math.cos(angle), cy - r * math.sin(angle)))

    pts = [(cx + r, cy), (cx - r * 0.5, cy - r * 0.87), (cx - r * 0.5, cy + r * 0.87)]

    lines = [head(w, h), bg(w, h)]
    lines.append(circle(cx, cy, r, C_AXIS, width=1, fill='none'))
    lines.append(line(pts[0][0], pts[0][1], pts[1][0], pts[1][1], C_MAIN, 2))
    lines.append(line(pts[1][0], pts[1][1], pts[2][0], pts[2][1], C_MAIN, 2))
    lines.append(line(pts[2][0], pts[2][1], pts[0][0], pts[0][1], C_MAIN, 2))
    for p in pts:
        lines.append(dot(p[0], p[1], 4))
    lines.append(text(pts[0][0]+10, pts[0][1]-10, 'A', 14, C_MAIN))
    lines.append(text(pts[1][0]-15, pts[1][1]-10, 'B', 14, C_MAIN))
    lines.append(text(pts[2][0]-15, pts[2][1]+15, 'C', 14, C_MAIN))
    # Side labels
    lines.append(text((pts[0][0]+pts[1][0])//2, (pts[0][1]+pts[1][1])//2-10, 'a', 13, C_COS, 'middle'))
    lines.append(text((pts[1][0]+pts[2][0])//2-15, (pts[1][1]+pts[2][1])//2, 'b', 13, C_SINE, 'middle'))
    lines.append(text((pts[2][0]+pts[0][0])//2+10, (pts[2][1]+pts[0][1])//2, 'c', 13, C_TAN, 'middle'))
    lines.append(text(cx, cy+r+25, 'a/sin(A) = b/sin(B) = c/sin(C) = 2R', 13, C_MAIN, 'middle'))
    lines.append('</svg>')
    return '\n'.join(lines)


def law_of_cosines():
    """Triangle with labeled sides for Law of Cosines."""
    w, h = 380, 300
    # Triangle
    ax, ay = 50, 240
    bx, by = 330, 240
    cx, cy = 130, 50

    lines = [head(w, h), bg(w, h)]
    lines.append(line(ax, ay, bx, by, C_MAIN, 2))
    lines.append(line(bx, by, cx, cy, C_MAIN, 2))
    lines.append(line(cx, cy, ax, ay, C_MAIN, 2))
    for p in [(ax, ay), (bx, by), (cx, cy)]:
        lines.append(dot(p[0], p[1], 4))
    # Labels
    lines.append(text(ax-15, ay+15, 'A', 14, C_MAIN))
    lines.append(text(bx+10, ay+15, 'B', 14, C_MAIN))
    lines.append(text(cx-10, cy-10, 'C', 14, C_MAIN))
    # Side labels
    lines.append(text((ax+bx)//2, ay+18, 'c', 14, C_MAIN, 'middle'))
    lines.append(text((bx+cx)//2+5, (by+cy)//2, 'a', 14, C_COS))
    lines.append(text((cx+ax)//2-5, (cy+ay)//2, 'b', 14, C_SINE, 'end'))
    # Angle label
    theta = math.atan2(cy-ay, cx-ax)
    lines.append(arc(ax, ay, 25, 0, theta, C_AXIS, 1.5))
    lines.append(text(ax+35, ay-30*math.sin(theta)+5, 'γ', 13, C_TEXT))
    lines.append(text(w//2, h-15, 'Law of Cosines: c² = a² + b² - 2ab·cos(γ)', 12, C_MAIN, 'middle'))
    lines.append('</svg>')
    return '\n'.join(lines)


def hyperbolic_functions():
    """Hyperbola x²-y²=1 with cosh/sinh point."""
    w, h = 450, 350
    cx, cy = 225, 175
    scale = 120

    lines = [head(w, h), bg(w, h)]
    lines.append(h_arrow(30, 430, cy, C_AXIS, 1))
    lines.append(v_arrow(cy+10, 20, cx, C_AXIS, 1))
    lines.append(text(428, cy+15, 'x', 14, C_TEXT))
    lines.append(text(cx+10, 22, 'y', 14, C_TEXT))

    # Hyperbola x² - y² = 1: right branch
    pts = []
    for i in range(-100, 101):
        t = i * 3 / 100
        x_val = math.cosh(t)
        y_val = math.sinh(t)
        if x_val * scale > 300:
            break
        pts.append((cx + x_val * scale, cy - y_val * scale))
    if len(pts) > 2:
        lines.append(polyline(pts, C_MAIN, 2))
    # Left branch
    pts = []
    for i in range(-100, 101):
        t = i * 3 / 100
        x_val = -math.cosh(t)
        y_val = math.sinh(t)
        if abs(x_val * scale) > 210:
            continue
        pts.append((cx + x_val * scale, cy - y_val * scale))
    if len(pts) > 2:
        lines.append(polyline(pts, C_MAIN, 2))

    # Asymptotes y = ±x
    asym = 210
    lines.append(line(cx-asym, cy+asym, cx+asym, cy-asym, C_AXIS, 1, '4,4'))
    lines.append(line(cx-asym, cy-asym, cx+asym, cy+asym, C_AXIS, 1, '4,4'))

    # Point P at t = 0.8
    t0 = 0.8
    px = cx + math.cosh(t0) * scale
    py = cy - math.sinh(t0) * scale
    lines.append(line(cx, cy, px, py, C_COS, 1.5))
    lines.append(dot(px, py, 4, C_TAN))
    lines.append(text(px+8, py-6, 'P', 13, C_TAN))

    # Labels
    lines.append(text(cx+5, cy-5, 'O', 14, C_TEXT))
    lines.append(text(cx+120, cy-90, 'x² - y² = 1', 12, C_MAIN))
    lines.append(text(cx+170, cy-20, 'y = x', 11, C_TEXT))
    lines.append(text(cx+120, cy+40, 'y = -x', 11, C_TEXT))
    lines.append(text(w//2, h-10, 'Hyperbolic functions: cosh(x), sinh(x)', 13, C_MAIN, 'middle'))
    lines.append('</svg>')
    return '\n'.join(lines)


def inverse_trig_graphs():
    """Arcsin and Arccos function graphs."""
    w, h = 500, 240
    margin = 45
    gx, gy = margin, h//2
    gw = w - 2*margin
    gh = 80

    lines = [head(w, h), bg(w, h)]
    lines.append(h_arrow(gx, w-margin+5, gy, C_AXIS, 1))
    lines.append(v_arrow(gy+10, margin-5, gx, C_AXIS, 1))
    lines.append(text(w-margin+5, gy+15, 'x', 13, C_TEXT))
    lines.append(text(gx-10, margin-5, 'y', 13, C_TEXT))

    # Domain labels -π/2 to π/2 for arcsin
    for val, label in [(-math.pi/2, '-π/2'), (-math.pi/4, '-π/4'), (math.pi/4, 'π/4'), (math.pi/2, 'π/2')]:
        x = gx + (val + math.pi/2) / math.pi * gw
        lines.append(text(x, gy+15, label, 9, C_TEXT, 'middle'))
    lines.append(text(gx + gw//2, gy+28, 'Domain', 9, C_TEXT, 'middle'))

    # Y-axis labels
    for val, label in [(1, 'π/2'), (-1, '-π/2')]:
        y = gy - val * gh
        lines.append(text(gx-8, y+4, label, 9, C_TEXT, 'end'))
        lines.append(line(gx-4, y, gx+4, y, C_AXIS, 0.5))

    # arcsin(x) on [-1, 1] -> [-π/2, π/2]
    pts = []
    for i in range(201):
        x_val = -1 + i * 2 / 200
        y_val = math.asin(x_val)
        px = gx + (x_val + 1) / 2 * gw  # Actually need to map -1..1 to the drawing
        # Wait, the domain here should map x from -1 to 1 onto the full width
        # But our x-axis is labeled -π/2 to π/2 for arcsin...
        # Let me re-think: arcsin takes x in [-1, 1] and outputs y in [-π/2, π/2]
        # We should plot on x-axis as input, y-axis as output
        px = gx + (x_val + 1) / 2 * gw  # x in [-1,1] maps to full width
        py = gy - y_val / (math.pi/2) * gh  # y in [-π/2, π/2] maps to full height
        pts.append((px, py))
    lines.append(polyline(pts, C_SINE, 2))
    lines.append(text(w-margin-5, gy-gh+12, 'arcsin(x)', 11, C_SINE, 'end'))

    # arccos(x) on [-1, 1] -> [0, π]
    pts = []
    for i in range(201):
        x_val = -1 + i * 2 / 200
        y_val = math.acos(x_val)
        px = gx + (x_val + 1) / 2 * gw
        py = gy - (y_val - math.pi/2) / (math.pi/2) * gh  # center at π/2
        pts.append((px, py))
    lines.append(polyline(pts, C_COS, 2))
    lines.append(text(w-margin-5, gy+gh-10, 'arccos(x)', 11, C_COS, 'end'))

    # Y-axis labels for π
    lines.append(text(gx-8, gy+4, '0', 9, C_TEXT, 'end'))
    for val, label in [(math.pi/2, 'π/2'), (math.pi, 'π'), (-math.pi/2, '-π/2')]:
        y = gy - (val - math.pi/2) / (math.pi/2) * gh
        lines.append(text(gx-8, y+4, label, 9, C_TEXT, 'end'))

    lines.append('</svg>')
    return '\n'.join(lines)


def sec_csc_cot_graphs():
    """Secant, Cosecant, Cotangent function graphs."""
    # Simple: plot sec(x) and csc(x) over [0, 2π]
    w, h = 500, 220
    margin = 40
    gx, gy = margin, h//2
    gw = w - 2*margin
    gh = 80

    lines = [head(w, h), bg(w, h)]
    lines.append(h_arrow(gx, w-margin+5, gy, C_AXIS, 1))
    lines.append(v_arrow(gy+10, margin-5, gx, C_AXIS, 1))

    for i, label in enumerate(['0', 'π/2', 'π', '3π/2', '2π']):
        x = gx + gw * i / 4
        lines.append(text(x, gy+15, label, 9, C_TEXT, 'middle'))
    lines.append(text(w-margin+5, gy+15, 'x', 12, C_TEXT))
    lines.append(text(gx-8, margin-5, 'y', 12, C_TEXT))

    for val, label in [(1, '1'), (-1, '-1')]:
        y = gy - val * gh
        lines.append(text(gx-8, y+4, label, 9, C_TEXT, 'end'))
        lines.append(line(gx-3, y, gx+3, y, C_AXIS, 0.5))

    # Asymptotes at π/2, 3π/2
    for i in range(5):
        x = gx + gw * i / 4
        if i % 2 == 1:  # π/2, 3π/2
            lines.append(line(x, margin-5, x, h-margin+5, C_AXIS, 0.5, '3,3'))

    # sec(x) = 1/cos(x)
    pts = []
    for i in range(200):
        t = (i / 200) * 2 * math.pi
        cos_t = math.cos(t)
        if abs(cos_t) > 0.05:
            y_val = 1 / cos_t
            if abs(y_val) < 4:
                pts.append((gx + t/(2*math.pi)*gw, gy - y_val*gh))
    if pts:
        lines.append(polyline(pts, C_COS, 1.5))

    # csc(x) = 1/sin(x)
    pts = []
    for i in range(200):
        t = (i / 200) * 2 * math.pi
        sin_t = math.sin(t)
        if abs(sin_t) > 0.05:
            y_val = 1 / sin_t
            if abs(y_val) < 4:
                pts.append((gx + t/(2*math.pi)*gw, gy - y_val*gh))
    if pts:
        lines.append(polyline(pts, C_SINE, 1.5))

    lines.append(text(w-margin-5, gy-gh+12, 'sec(x)', 10, C_COS, 'end'))
    lines.append(text(w-margin-5, gy+gh-10, 'csc(x)', 10, C_SINE, 'end'))
    lines.append('</svg>')
    return '\n'.join(lines)


def sine_cosine_graph_svg():
    """Sine and cosine function graphs."""
    w, h = 500, 220
    margin = 40
    gx, gy = margin, h // 2
    gw = w - 2 * margin
    gh = 90

    lines = [head(w, h), bg(w, h)]
    lines.append(h_arrow(gx, w - margin + 10, gy, C_AXIS, 1))
    lines.append(v_arrow(gy+10, margin - 10, gx, C_AXIS, 1))
    lines.append(text(w-margin+5, gy+15, 'x', 11, C_TEXT))
    lines.append(text(gx+5, margin-5, 'y', 11, C_TEXT))

    for i, label in enumerate(['0', 'π/2', 'π', '3π/2', '2π']):
        x = gx + gw * i / 4
        lines.append(text(x, gy+15, label, 10, C_TEXT, 'middle'))
    for val, label in [(1, '1'), (-1, '-1')]:
        y = gy - val * gh
        lines.append(text(gx-8, y+4, label, 10, C_TEXT, 'end'))
        lines.append(line(gx-4, y, gx+4, y, C_AXIS, 0.5))

    pts = []
    for i in range(201):
        t = (i / 200) * 2 * math.pi
        pts.append((gx + t/(2*math.pi)*gw, gy - math.sin(t)*gh))
    lines.append(polyline(pts, C_SINE, 2))
    lines.append(text(w-margin, gy-gh+15, 'sin(x)', 11, C_SINE, 'end'))

    pts = []
    for i in range(201):
        t = (i / 200) * 2 * math.pi
        pts.append((gx + t/(2*math.pi)*gw, gy - math.cos(t)*gh))
    lines.append(polyline(pts, C_COS, 2))
    lines.append(text(w-margin, gy+gh-10, 'cos(x)', 11, C_COS, 'end'))

    lines.append('</svg>')
    return '\n'.join(lines)


def complex_plane():
    """Complex plane with point a+bi."""
    w, h = 400, 400
    cx, cy = 200, 200
    scale = 120
    a, b = 0.7, 0.5  # Real and imaginary parts

    px = cx + a * scale
    py = cy - b * scale

    lines = [head(w, h), bg(w, h)]
    lines.append(h_arrow(30, 380, cy, C_AXIS, 1.5))
    lines.append(v_arrow(cy+10, 30, cx, C_AXIS, 1.5))
    lines.append(text(378, cy+18, 'Re', 14, C_TEXT))
    lines.append(text(cx+10, 30, 'Im', 14, C_TEXT))
    lines.append(text(cx+5, cy-5, 'O', 14, C_TEXT))

    # Point
    lines.append(line(cx, cy, px, py, C_MAIN, 2))
    lines.append(dot(px, py, 5))
    lines.append(text(px+8, py-8, 'z = a + bi', 13, C_MAIN))
    # Dashed projections
    lines.append(line(px, cy, px, py, C_AXIS, 1, '4,4'))
    lines.append(line(cx, py, px, py, C_AXIS, 1, '4,4'))
    # Labels
    lines.append(text(px, cy+18, 'a', 13, C_COS, 'middle'))
    lines.append(text(cx-15, py+4, 'b', 13, C_SINE, 'end'))
    # Angle arc
    theta = math.atan2(b, a)
    lines.append(arc(cx, cy, 30, 0, theta, C_AXIS, 1.5))
    lines.append(text(cx+40, cy-15, 'θ', 13, C_TEXT))
    # Radius label
    mx = (cx + px) / 2
    my = (cy + py) / 2
    lines.append(text(mx+10, my, 'r = |z|', 12, C_MAIN))
    lines.append('</svg>')
    return '\n'.join(lines)


def vector_arrow():
    """Vector arrow in plane with components."""
    w, h = 400, 350
    cx, cy = 80, 280
    vx, vy = 220, -150

    lines = [head(w, h), bg(w, h)]
    lines.append(h_arrow(30, 380, cy, C_AXIS, 1))
    lines.append(v_arrow(cy+10, 30, cx, C_AXIS, 1))
    lines.append(text(378, cy+15, 'x', 13, C_TEXT))
    lines.append(text(cx+8, 30, 'y', 13, C_TEXT))
    lines.append(text(cx-10, cy+5, 'O', 13, C_TEXT))

    # Vector
    dx = cx + vx
    dy = cy + vy
    lines.append(line(cx, cy, dx, dy, C_MAIN, 2.5))
    theta = math.atan2(cy - dy, dx - cx)
    lines.append(arrow_tip(dx, dy, theta, C_MAIN))
    lines.append(dot(dx, dy, 4))
    lines.append(text(dx+10, dy-5, 'v', 14, C_MAIN))

    # Components
    lines.append(line(cx, cy, cx+vx, cy, C_COS, 1.5, '6,4'))
    lines.append(line(cx+vx, cy, dx, dy, C_SINE, 1.5, '6,4'))
    # Component labels
    lines.append(text(cx+vx//2, cy+18, 'v_x', 12, C_COS, 'middle'))
    lines.append(text(cx+vx+12, (cy+dy)//2, 'v_y', 12, C_SINE))
    lines.append('</svg>')
    return '\n'.join(lines)


def vector_addition():
    """Parallelogram law of vector addition."""
    w, h = 400, 350
    ox, oy = 80, 260
    ux, uy = 140, -180
    vx, vy = 180, -60

    lines = [head(w, h), bg(w, h)]
    lines.append(h_arrow(20, 380, oy, C_AXIS, 1))
    lines.append(v_arrow(oy+10, 20, ox, C_AXIS, 1))
    lines.append(text(378, oy+15, 'x', 13, C_TEXT))
    lines.append(text(ox+8, 22, 'y', 13, C_TEXT))
    lines.append(text(ox-10, oy+5, 'O', 13, C_TEXT))

    # Vector u
    u_end = (ox+ux, oy+uy)
    lines.append(line(ox, oy, u_end[0], u_end[1], C_SINE, 2.5))
    theta = math.atan2(oy - u_end[1], u_end[0] - ox)
    lines.append(arrow_tip(u_end[0], u_end[1], theta, C_SINE))
    lines.append(text(u_end[0]+10, u_end[1]-5, 'u', 14, C_SINE))

    # Vector v
    v_end = (ox+vx, oy+vy)
    lines.append(line(ox, oy, v_end[0], v_end[1], C_COS, 2.5))
    theta = math.atan2(oy - v_end[1], v_end[0] - ox)
    lines.append(arrow_tip(v_end[0], v_end[1], theta, C_COS))
    lines.append(text(v_end[0]+10, v_end[1]+5, 'v', 14, C_COS))

    # Parallelogram completion
    uv = (ox+ux+vx, oy+uy+vy)
    lines.append(line(u_end[0], u_end[1], uv[0], uv[1], C_AXIS, 1.5, '4,4'))
    lines.append(line(v_end[0], v_end[1], uv[0], uv[1], C_AXIS, 1.5, '4,4'))

    # Resultant u+v
    lines.append(line(ox, oy, uv[0], uv[1], C_MAIN, 2.5))
    theta = math.atan2(oy - uv[1], uv[0] - ox)
    lines.append(arrow_tip(uv[0], uv[1], theta, C_MAIN))
    lines.append(text(uv[0]+10, uv[1]-5, 'u+v', 14, C_MAIN))
    lines.append('</svg>')
    return '\n'.join(lines)


def pascals_triangle():
    """Pascal's triangle diagram."""
    rows = 6
    cell_w = 50
    cell_h = 35
    cx = 250
    top_y = 40

    lines = [head(500, 300), bg(500, 300)]
    lines.append(text(cx, 20, "Pascal's Triangle", 15, C_MAIN, 'middle'))

    def nCr(n, r):
        return math.comb(n, r)

    for n in range(rows):
        y = top_y + n * cell_h
        for r in range(n + 1):
            x = cx + (r - n / 2) * cell_w
            val = str(nCr(n, r))
            lines.append(text(x, y, val, 12, C_MAIN, 'middle'))

    lines.append('</svg>')
    return '\n'.join(lines)


def completing_square():
    """Visual representation of completing the square."""
    w, h = 350, 300
    s = 80  # Side of the x² square
    bx = 20  # Coefficient for x term (rectangle width)
    gap = 5

    x, y = 40, 40  # Top-left of x² square

    lines = [head(w, h), bg(w, h)]
    # x² square
    lines.append(polyline([(x, y), (x+s, y), (x+s, y+s), (x, y+s)], C_MAIN, 2))
    lines.append(text(x+s//2, y+s//2, 'x²', 14, C_MAIN, 'middle'))

    # bx rectangle (for 2 * (b/2) * x)
    rx = x
    ry = y + s + gap
    rw = s
    rh = bx
    lines.append(polyline([(rx, ry), (rx+rw, ry), (rx+rw, ry+rh), (rx, ry+rh)], C_COS, 2))
    lines.append(text(rx+rw//2, ry+rh//2, 'bx', 13, C_COS, 'middle'))

    # (b/2)² small square
    sx2 = x + s + gap
    sy2 = y
    ss = bx
    lines.append(polyline([(sx2, sy2), (sx2+ss, sy2), (sx2+ss, sy2+ss), (sx2, sy2+ss)], C_COS, 2))
    lines.append(text(sx2+ss//2, sy2+ss//2, '(b/2)²', 11, C_COS, 'middle'))

    # Completed square outline
    big = s + bx
    lines.append(polyline([(x, y), (x+big, y), (x+big, y+big), (x, y+big)], C_TAN, 1.5, 'none'))
    lines.append(text(x+big//2, y+big+20, '(x + b/2)²', 13, C_TAN, 'middle'))

    lines.append(text(w//2, h-15, 'Completing the Square', 13, C_MAIN, 'middle'))
    lines.append('</svg>')
    return '\n'.join(lines)


def number_line_absolute_value():
    """Number line showing absolute value as distance from zero."""
    w, h = 450, 120
    cx, cy = 225, 65
    extent = 180

    lines = [head(w, h), bg(w, h)]
    lines.append(h_arrow(cx-extent-10, cx+extent+10, cy, C_AXIS, 1.5))
    lines.append(text(cx+extent+18, cy+5, 'x', 14, C_TEXT))

    # Tick marks
    for val in range(-3, 4):
        x = cx + val * 50
        lines.append(line(x, cy-5, x, cy+5, C_MAIN, 1))
        lines.append(text(x, cy+20, str(val), 12, C_TEXT, 'middle'))

    # Highlight |3| = 3
    x2 = cx + 3 * 50
    xneg2 = cx - 3 * 50
    lines.append(line(cx, cy, x2, cy, C_SINE, 3))
    lines.append(line(xneg2, cy, cx, cy, C_COS, 3))

    # Distance brackets
    lines.append(text((cx+x2)//2, cy-12, '|3| = 3', 12, C_SINE, 'middle'))
    lines.append(text((xneg2+cx)//2, cy-12, '|-3| = 3', 12, C_COS, 'middle'))

    lines.append(text(cx, cy+45, 'Absolute value as distance from zero', 12, C_MAIN, 'middle'))
    lines.append('</svg>')
    return '\n'.join(lines)


def number_line_intervals():
    """Number line showing different interval types."""
    w, h = 500, 160
    cx, cy = 250, 55
    extent = 220

    lines = [head(w, h), bg(w, h)]
    lines.append(h_arrow(cx-extent-10, cx+extent+10, cy, C_AXIS, 1.5))

    ticks = range(-4, 5)
    for val in ticks:
        x = cx + val * 45
        lines.append(line(x, cy-5, x, cy+5, C_MAIN, 1))
        lines.append(text(x, cy+18, str(val), 11, C_TEXT, 'middle'))

    # Closed interval [-2, 2]
    xl = cx - 2 * 45
    xr = cx + 2 * 45
    lines.append(line(xl, cy-10, xr, cy-10, C_SINE, 2.5))
    lines.append(dot(xl, cy-10, 4, C_SINE))
    lines.append(dot(xr, cy-10, 4, C_SINE))
    lines.append(text((xl+xr)//2, cy-22, '[-2, 2] closed', 10, C_SINE, 'middle'))

    # Open interval (1, 4)
    xl = cx + 1 * 45
    xr = cx + 4 * 45
    lines.append(line(xl, cy-35, xr, cy-35, C_COS, 2))
    lines.append(circle(xl, cy-35, 4, C_COS, 'white'))
    lines.append(circle(xr, cy-35, 4, C_COS, 'white'))
    lines.append(text((xl+xr)//2, cy-48, '(1, 4) open', 10, C_COS, 'middle'))

    lines.append('</svg>')
    return '\n'.join(lines)


def number_line_real():
    """Number line showing real numbers."""
    w, h = 450, 100
    cx, cy = 225, 45
    extent = 200

    lines = [head(w, h), bg(w, h)]
    lines.append(h_arrow(cx-extent-10, cx+extent+10, cy, C_AXIS, 2))

    for val, label in [(-2, '-2'), (-1, '-1'), (0, '0'), (1, '1'), (2, '2')]:
        x = cx + val * 70
        lines.append(line(x, cy-6, x, cy+6, C_MAIN, 1.5))
        lines.append(text(x, cy+20, label, 12, C_TEXT, 'middle'))

    # Fill the line
    lines.append(text(cx, cy-15, 'Each point on the line corresponds to a real number', 11, C_MAIN, 'middle'))
    lines.append('</svg>')
    return '\n'.join(lines)


def number_types_venn():
    """Venn diagram of number systems: N ⊂ Z ⊂ Q ⊂ R ⊂ C."""
    w, h = 420, 310
    cx, cy = 210, 155

    lines = [head(w, h), bg(w, h)]

    # Outer rectangle for complex
    lines.append(polyline([(30, 20), (390, 20), (390, 290), (30, 290)], C_COS, 1.5))
    lines.append(text(380, 285, 'C (Complex)', 11, C_COS, 'end'))

    # Real numbers (large ellipse)
    lines.append(path('M 60 155 C 60 60, 360 60, 360 155 C 360 250, 60 250, 60 155', C_MAIN, 1.5))
    lines.append(text(350, 160, 'R (Real)', 11, C_MAIN, 'end'))

    # Rational numbers (ellipse)
    lines.append(path('M 100 155 C 100 95, 320 95, 320 155 C 320 215, 100 215, 100 155', C_SINE, 1.5))
    lines.append(text(310, 150, 'Q (Rational)', 10, C_SINE, 'end'))

    # Integers (ellipse)
    lines.append(path('M 140 155 C 140 115, 280 115, 280 155 C 280 195, 140 195, 140 155', C_TAN, 1.5))
    lines.append(text(270, 150, 'Z (Integer)', 10, C_TAN, 'end'))

    # Natural numbers
    lines.append(path('M 175 155 C 175 130, 245 130, 245 155 C 245 180, 175 180, 175 155', '#2980b9', 1.5))
    lines.append(text(215, 155, 'N', 12, '#2980b9', 'middle', 'bold'))

    # Irrational annotation
    lines.append(text(80, 100, 'Irrational', 10, '#e74c3c', 'middle'))
    lines.append(line(80, 108, 80, 130, '#e74c3c', 0.5))

    lines.append(text(cx, h-12, 'Hierarchy of Number Systems', 13, C_MAIN, 'middle'))
    lines.append('</svg>')
    return '\n'.join(lines)


def polynomial_roots_graph():
    """Graph of polynomial showing roots on x-axis."""
    w, h = 450, 300
    cx, cy = 225, 160
    scale_x = 60
    scale_y = 120

    lines = [head(w, h), bg(w, h)]
    lines.append(h_arrow(30, 430, cy, C_AXIS, 1.5))
    lines.append(v_arrow(cy+10, 20, cx, C_AXIS, 1.5))
    lines.append(text(428, cy+15, 'x', 13, C_TEXT))
    lines.append(text(cx+8, 22, 'y', 13, C_TEXT))

    for val in range(-3, 4):
        if val == 0:
            continue
        x = cx + val * scale_x
        lines.append(line(x, cy-3, x, cy+3, C_AXIS, 0.5))
        lines.append(text(x, cy+18, str(val), 10, C_TEXT, 'middle'))

    # Polynomial: (x+2)(x)(x-1.5) ≈ x³ + 0.5x² - 3x
    pts = []
    for i in range(300):
        t = (i / 299) * 4 - 2  # [-2, 2]
        # f(t) = 0.25*(t+2)(t)(t-1.5) scaled to fit
        y_val = 1.0 * (t+2) * t * (t-1.5)
        px = cx + t * scale_x
        py = cy - y_val * 30
        if 30 <= px <= 430 and 20 <= py <= cy+10:
            pts.append((px, py))
    if len(pts) > 2:
        lines.append(polyline(pts, C_SINE, 2.5))

    # Root markers
    for root in [-2, 0, 1.5]:
        x = cx + root * scale_x
        lines.append(dot(x, cy, 5, C_COS))
        # Dash down
        lines.append(line(x, cy+5, x, cy+10, C_COS, 1))
        lines.append(text(x, cy+25, f'x={root}', 9, C_COS, 'middle'))

    lines.append(text(w-50, 40, 'f(x) = (x+2)(x)(x-1.5)', 11, C_MAIN))
    lines.append(text(cx, h-15, 'Roots of a Polynomial', 12, C_MAIN, 'middle'))
    lines.append('</svg>')
    return '\n'.join(lines)


def linear_equation_graph():
    """Graph of a linear equation y = mx + b."""
    w, h = 400, 300
    cx, cy = 200, 160
    scale = 40
    m, b = 1.5, -1  # y = 1.5x - 1

    lines = [head(w, h), bg(w, h)]
    lines.append(h_arrow(30, 380, cy, C_AXIS, 1.5))
    lines.append(v_arrow(cy+10, 20, cx, C_AXIS, 1.5))
    lines.append(text(378, cy+15, 'x', 13, C_TEXT))
    lines.append(text(cx+8, 22, 'y', 13, C_TEXT))
    lines.append(text(cx-10, cy+15, 'O', 13, C_TEXT))

    for val in range(-4, 5):
        if val == 0:
            continue
        x = cx + val * scale
        y = cy - val * scale
        lines.append(line(x, cy-3, x, cy+3, C_AXIS, 0.5))
        lines.append(text(x, cy+18, str(val), 10, C_TEXT, 'middle'))
        lines.append(line(cx-3, y, cx+3, y, C_AXIS, 0.5))
        lines.append(text(cx-12, y+4, str(val), 10, C_TEXT, 'end'))

    # Graph line from x=-4 to x=4
    x1, y1 = -4, m*(-4) + b
    x2, y2 = 4, m*4 + b
    px1 = cx + x1 * scale
    py1 = cy - y1 * scale
    px2 = cx + x2 * scale
    py2 = cy - y2 * scale
    lines.append(line(px1, py1, px2, py2, C_SINE, 2.5))

    # Label
    mx_line = (px1 + px2) / 2
    my_line = (py1 + py2) / 2
    lines.append(text(mx_line+15, my_line-10, 'y = 1.5x - 1', 12, C_SINE))

    # Slope indicator
    dx = 40
    dy = -dx * m
    x_mid = cx + 2 * scale
    y_mid = cy - (m*2 + b) * scale
    lines.append(line(x_mid, y_mid, x_mid+dx, y_mid, C_COS, 1, '4,4'))
    lines.append(line(x_mid+dx, y_mid, x_mid+dx, y_mid+dy, C_COS, 1, '4,4'))
    lines.append(text(x_mid+dx//2, y_mid+12, 'Δx', 9, C_COS, 'middle'))
    lines.append(text(x_mid+dx+10, y_mid+dy//2, 'Δy', 9, C_COS))

    lines.append(text(cx, h-15, 'Linear Equation y = mx + b', 12, C_MAIN, 'middle'))
    lines.append('</svg>')
    return '\n'.join(lines)


# ─── Main ───────────────────────────────────────────────────────────────────

def copy_to_data_svgs(name, domain):
    """Copy SVG from public dir to data/lessons svg dir."""
    src = os.path.join(OUT, name)
    dst_dir = os.path.join(DATA, domain, 'svg')
    os.makedirs(dst_dir, exist_ok=True)
    dst = os.path.join(dst_dir, name)
    import shutil
    shutil.copy2(src, dst)


if __name__ == '__main__':
    os.makedirs(OUT, exist_ok=True)

    svgs = {
        # Trigonometry
        'unit-circle-labeled.svg': (unit_circle_labeled, 'trigonometry'),
        'unit-circle-sine-cosine.svg': (unit_circle_sincos, 'trigonometry'),
        'unit-circle-tangent.svg': (unit_circle_tangent, 'trigonometry'),
        'right-triangle-unit-circle.svg': (right_triangle_unit_circle, 'trigonometry'),
        'reference-angles.svg': (reference_angles, 'trigonometry'),
        'right-triangle-trig.svg': (right_triangle_trig, 'trigonometry'),
        'pythagorean-theorem.svg': (pythagorean_theorem, 'trigonometry'),
        'sine-cosine-graph.svg': (sine_cosine_graph_svg, 'trigonometry'),
        'hyperbolic-functions.svg': (hyperbolic_functions, 'trigonometry'),
        'inverse-trig-graphs.svg': (inverse_trig_graphs, 'trigonometry'),
        'sec-csc-cot-graphs.svg': (sec_csc_cot_graphs, 'trigonometry'),
        # Law of sines/cosines
        'law-of-sines.svg': (law_of_sines, 'trigonometry'),
        'law-of-cosines.svg': (law_of_cosines, 'trigonometry'),
        # Vectors
        'vector-arrow.svg': (vector_arrow, 'vectors-and-matrices'),
        'vector-addition.svg': (vector_addition, 'vectors-and-matrices'),
        # Complex numbers
        'complex-plane.svg': (complex_plane, 'complex-numbers'),
        # Polynomials
        'pascals-triangle.svg': (pascals_triangle, 'polynomials'),
        'completing-square.svg': (completing_square, 'polynomials'),
        'polynomial-roots-graph.svg': (polynomial_roots_graph, 'polynomials'),
        # Equations
        'linear-equation-graph.svg': (linear_equation_graph, 'equations'),
        # Sets and numbers
        'number-line-absolute-value.svg': (number_line_absolute_value, 'sets-and-numbers'),
        'number-line-intervals.svg': (number_line_intervals, 'sets-and-numbers'),
        'number-line-real.svg': (number_line_real, 'sets-and-numbers'),
        'number-types-venn.svg': (number_types_venn, 'sets-and-numbers'),
    }

    for name, (gen_func, domain) in svgs.items():
        content = gen_func()
        out_path = os.path.join(OUT, name)
        with open(out_path, 'w') as f:
            f.write(content)
        print(f'Generated {out_path}  [{domain}]')
        copy_to_data_svgs(name, domain)
        print(f'  → data/lessons/algebrica/{domain}/svg/{name}')
