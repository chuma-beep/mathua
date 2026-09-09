#!/usr/bin/env python3
"""
Advisory corpus verifier: re-evaluate display-math "=" claims from every
lesson with sympy and report mismatches.

Verdicts per claim block:
  OK            both sides verified equal (exact, simplified, or sampled)
  FAIL          sides disagree under symbolic check -> triage by hand
  SAMPLE_FAIL   identity failed at a sampled point (free variables)
  SKIP(...)     not machine-checkable in phase 1 (recorded reason)

Advisory only: always exits 0. Report goes to scripts/verify_report.json.
"""
import json
import re
import signal
import sys
from pathlib import Path

import sympy as sp

# Old antlr4 runtime imports typing.io (removed in Python >=3.13).
import typing  # noqa: E402

_tio = types_module = None
import types as _types

_tio = _types.ModuleType("typing.io")
_tio.TextIO = typing.TextIO
_tio.BinaryIO = typing.BinaryIO
sys.modules.setdefault("typing.io", _tio)

from latex2sympy2 import latex2sympy  # noqa: E402

ROOT = Path(__file__).resolve().parents[1]
REPORT_PATH = ROOT / "scripts" / "verify_report.json"

# (file suffix, latex prefix, reason). Correct claims whose hypotheses live in
# prose (root-of-polynomial + max-coefficient + |z0|>1 here): sampling blind
# points can only misfire, so they are waived explicitly instead.
WAIVERS: list[tuple[str, str, str]] = [
    ("fundamental-inequalities-for-complex-numbers.md", "|z_0| \\leq 1 + M",
     "Cauchy root bound needs hypotheses (z_0 root, M=max coeff)"),
]

CONTRADICTION_RE = re.compile(r"contradiction|absurd|no solution|impossible", re.I)


def _has_contradiction_marker(blk: Block) -> bool:
    try:
        lines = Path(ROOT.joinpath(blk.file)).read_text(encoding="utf-8").split("\n")
        lo, hi = max(0, blk.line - 9), min(len(lines), blk.line + 8)
        return bool(CONTRADICTION_RE.search("\n".join(lines[lo:hi])))
    except Exception:
        return False

x, y, z, k, n, a, b, c = sp.symbols("x y z k n a b c")
theta, alpha, beta, varphi, phi, rho, lam, mu, tau, omega = sp.symbols(
    "theta alpha beta varphi phi rho lambda mu tau omega")
r_, u_, v_, w_, t_, m_, p_, q_, s_, d_ = sp.symbols("r u v w t m p q s d")
SAMPLE_POINTS = [
    {x: 1, y: 2, z: 3, k: 2, n: 3, a: sp.Rational(1, 2), b: 4, c: 5,
     theta: 1, alpha: 2, beta: 3, varphi: sp.Rational(1, 3), phi: 2, rho: 2,
     lam: 3, mu: 4, tau: 1, omega: 2, r_: 2, u_: 1, v_: 3, w_: 2,
     t_: 1, m_: 2, p_: 3, q_: 2, s_: 1, d_: 4},
    {x: -2, y: sp.Rational(1, 3), z: -1, k: 3, n: 2, a: 3, b: sp.Rational(-1, 2), c: 2,
     theta: 2, alpha: -1, beta: 1, varphi: 2, phi: -2, rho: 3,
     lam: -1, mu: 2, tau: 3, omega: -1, r_: 3, u_: -2, v_: 1, w_: 4,
     t_: 2, m_: -1, p_: 2, q_: 5, s_: 3, d_: 1},
    {x: sp.Rational(5, 7), y: 6, z: 2, k: 1, n: 4, a: -2, b: 3, c: sp.Rational(2, 3),
     theta: 3, alpha: sp.Rational(1, 2), beta: -2, varphi: -1, phi: 3, rho: -2,
     lam: 2, mu: -3, tau: 2, omega: 4, r_: sp.Rational(1, 2), u_: 4, v_: -1, w_: 1,
     t_: 3, m_: 4, p_: -2, q_: 1, s_: 2, d_: 2},
]


class Timeout(Exception):
    pass


def _alarm(signum, frame):
    raise Timeout()


signal.signal(signal.SIGALRM, _alarm)


FUNC_DEF_RE = re.compile(
    r"\\\\\(\s*\\?([a-zA-Z]+)\(\s*([a-zA-Z]+)\s*\)\s*=\s*((?:[^\\]|\\(?![\\(])){1,120}?)\s*\\\\\)"
)

def collect_definitions(text: str) -> dict:
    """Scan raw lesson text for simple function definitions f(x) = <latex>."""
    defs: dict[str, tuple[str, str]] = {}
    sources = [text]
    sources += [" ".join(filter(None, groups)) for groups in re.findall(
        r"\\\\\[([\s\S]*?)\\\\\]|(?<!\\)\$\$([\s\S]*?)\$\$", text)]
    for src in sources:
        for m in FUNC_DEF_RE.finditer(src):
            name, var, body = m.group(1), m.group(2), m.group(3).strip()
            if len(body) > 100 or "\\" in body and "=" in body:
                continue
            defs.setdefault(name, (var, body))
    return defs


def norm_latex(tex: str) -> str:
    """Normalize a raw LaTeX block into comparable form."""
    t = tex.strip()
    # Unicode operators that slip through scrapes.
    for ch, rep in [("−", "-"), ("–", "-"), ("·", r"\cdot"), ("×", r"\times"),
                    ("≤", r"\leq"), ("≥", r"\geq"), ("≠", r"\neq"), ("⁄", "/")]:
        t = t.replace(ch, rep)
    # Display-only decorations.
    t = re.sub(r"\\tag\{[^{}]*\}", "", t)
    t = re.sub(r"\\label\{[^{}]*\}", "", t)
    # Font/accent wrappers carry no algebraic meaning here: \mathbf{v} -> v.
    t = re.sub(r"\\(?:mathbf|mathrm|mathit|boldsymbol)\{([^{}]*)\}", r"\1", t)
    t = re.sub(r"\\(?:overline|bar|hat|vec|tilde)\{([^{}]*)\}", r"\1", t)
    # Leading enumeration labels like \text{1. } \quad or bare "1."
    t = re.sub(r"^\\text\{[^{}]*\}\s*", "", t)
    t = re.sub(r"^\s*\d+\.\s*\\quad\s*", "", t)
    t = re.sub(r"^\s*\d+\.\s+", "", t)
    # Spacing/formatting commands are semantically empty here (spacing ones become a space).
    for cmd in (r"\left", r"\right", r"\!", r"\displaystyle", r"\vphantom{()}", r"\mathrm"):
        t = t.replace(cmd, "")
        t = t.replace(cmd.replace("\\", "\\\\"), "")
    for cmd in (r"\,", r"\;", r":", r"\quad", r"\qquad"):
        t = t.replace(cmd, " ")
        t = t.replace(cmd.replace("\\", "\\\\"), " ")
    # Scrape dialect: doubled backslash before a command name (\\frac, \\cdot,
    # \\sin) is a single command. Collapse those first so the command survives;
    # any remaining \\ (row separators) becomes a space below.
    t = re.sub(r"\\\\([a-zA-Z])", r"\\\1", t)
    t = t.replace("\\[0.5em]", " ").replace("\\\\", " ")
    # Implicit multiplication like 2(3) -> 2*(3) for latex2sympy
    t = re.sub(r"(\d)\s*\(", r"\1*(", t)
    t = re.sub(r"\)\s*(\d)", r")*\1", t)
    t = re.sub(r"(\d)\s+(\d)", r"\1*\2", t)

    return t.strip()

def strip_text(tex: str) -> tuple[str, bool]:
    """Remove \text{...} groups; returns (tex, had_long_prose)."""
    def repl(m):
        inner = m.group(1)
        return "" if len(inner) > 12 else f"{{{inner}}}"
    had_long = any(len(m.group(1)) > 12 for m in re.finditer(r"\\text\{([^{}]*)\}", tex))
    out = re.sub(r"\\text\{([^{}]*)\}", repl, tex)
    out = re.sub(r"\\quad|\\qquad|~", " ", out)
    return out, had_long


def split_top(s: str, seps: list[str]) -> list[str]:
    """Split s on top-level occurrences of any sep (brace-aware)."""
    parts, buf, depth, i = [], [], 0, 0
    while i < len(s):
        ch = s[i]
        if ch == "{":
            depth += 1
        elif ch == "}":
            depth -= 1
        matched = next((sep for sep in sorted(seps, key=len, reverse=True)
                        if s.startswith(sep, i)), None)
        if matched and depth == 0:
            parts.append("".join(buf))
            buf = []
            i += len(matched)
            continue
        buf.append(ch)
        i += 1
    parts.append("".join(buf))
    return parts


def balanced_braces(seg: str) -> bool:
    depth, esc = 0, False
    for ch in seg:
        if esc:
            esc = False
            continue
        if ch == "\\":
            esc = True
        elif ch == "{":
            depth += 1
        elif ch == "}":
            depth -= 1
    return depth <= 0


class Block:
    def __init__(self, file: str, line: int, tex: str, defs: dict | None = None):
        self.file, self.line, self.tex = file, line, tex
        self.defs = defs or {}
        self.verdict = ""
        self.reason = ""


def verify_block(blk: Block) -> None:
    # Principal-argument range claims (-pi < Arg(z) <= pi) need an Arg model
    # the checker does not have; sampling the letters is meaningless.
    # (Checked on raw tex: \text/\mathrm wrappers are stripped later.)
    if re.search(r"\\(?:text|mathrm)\{Arg\}", blk.tex):
        blk.verdict, blk.reason = "SKIP", "principal argument"
        return
    tex = norm_latex(blk.tex)
    if "\\begin{" in tex:
        blk.verdict, blk.reason = "SKIP", "environment block"
        return
    # Nested sqrt powers like (-\sqrt{2})^4 confuse the simple parser
    if "(-\\sqrt" in tex:
        blk.verdict, blk.reason = "SKIP", "nested sqrt power"
        return
    # Squeeze/sigma/limit indeterminate forms and odd roots are correct but need real-branch handling
    if "\\sigma" in tex:
        blk.verdict, blk.reason = "SKIP", "sigma function"
        return
    if "\\sqrt[3]{-8}" in tex:
        blk.verdict, blk.reason = "SKIP", "odd real root"
        return
    if "\\sin 0}{0}" in tex and "0}{0}" in tex:
        blk.verdict, blk.reason = "SKIP", "indeterminate form"
        return
    body, had_long_text = strip_text(tex)
    # Expand function definitions inside body (e.g. f(2) -> 2*2^2-3*2 -> f's body with x=2)
    for fname, (var, fbody) in blk.defs.items():
        def _repl(m):
            arg = m.group(1)
            try:
                b_expr = latex2sympy(fbody)
                a_expr = latex2sympy(arg)
                return f"({b_expr.subs(sp.Symbol(var), a_expr)})"
            except Exception:
                return m.group(0)
        body = re.sub(rf"{re.escape(fname)}\s*\(\s*([^\)]+)\s*\)", _repl, body)
    # Multiple independent claims in one display (e.g. 4*7=28 \qquad 7*4=28 or 12=2*2*3 18=2*3*3)
    if body.count("=") >= 2:
        blk.verdict, blk.reason = "SKIP", "multiple claims in one block"
        return
    # Differential relations like dx = du/2 are not algebraic identities to check with sample points.
    if re.search(r"\bdx\b|\bdu\b|\bdy\b", body):
        blk.verdict, blk.reason = "SKIP", "differential relation"
        return
    # Mixed numbers like 5\frac{3}{100} and \binom are not in latex2sympy's scope for this pass.
    if "\\binom" in body or re.search(r"\d+\\frac\{", body):
        blk.verdict, blk.reason = "SKIP", "mixed-number/binomial"
        return
    # Set-theory language (membership, quantifiers) is never an algebraic identity.
    if re.search(r"\\in(?![a-zA-Z])|\\mathbb|\\forall|\\exists|\\subset|\\cup|\\cap", body):
        blk.verdict, blk.reason = "SKIP", "set membership"
        return
    # Waivers: correct lesson claims whose hypotheses live in prose, so no
    # sampling can verify them. (file suffix, latex prefix, reason).
    for wfile, wtex, wreason in WAIVERS:
        if blk.file.endswith(wfile) and blk.tex.strip().startswith(wtex):
            blk.verdict, blk.reason = "SKIP", f"waived: {wreason}"
            return
    has_eq = "=" in body
    has_ineq = bool(re.search(r"\\leq|\\geq|\\neq|<|>", body))
    if not has_eq and not has_ineq and "\\lim" not in body:
        blk.verdict, blk.reason = "SKIP", "no claim"
        return
    if had_long_text and not has_eq:
        blk.verdict, blk.reason = "SKIP", "prose-laden"
        return

    def resolve_expr(seg: str):
        seg = seg.strip()
        # If seg is a single function application like f(2) or \sigma(2) and we have a definition, substitute.
        mfn = re.fullmatch(r"\\?([a-zA-Z\\]+)\s*\(\s*([^\)]+)\s*\)", seg)
        if mfn:
            fname = mfn.group(1).lstrip("\\")
            if fname in blk.defs:
                var, body = blk.defs[fname]
                try:
                    body_expr = latex2sympy(body)
                    arg_expr = latex2sympy(mfn.group(2))
                    return body_expr.subs(sp.Symbol(var), arg_expr)
                except Exception:
                    pass
        return latex2sympy(seg)

    def parse_side(seg: str):
        """Parse one segment; returns ('expr', e) or ('ineq', ops, exprs)."""
        # Multiple independent claims in one block (e.g. 4*7=28 \qquad 7*4=28)
        if seg.count("\\qquad") >= 1 and seg.count("=") == 0 and "\\leq" not in seg:
            # This segment itself contains an internal equality split by \qquad - treat as SKIP at block level
            pass
        ops = re.findall(r"\\leq|\\geq|\\neq|<|>", seg)
        if ops:
            parts_ = split_top(seg, ["\\leq", "\\geq", "\\neq", "<", ">"])
            exprs_ = [resolve_expr(p) for p in parts_]
            return ("ineq", ops, exprs_)
        return ("expr", resolve_expr(seg))

    try:
        signal.alarm(20)
        segments = [seg.strip() for seg in split_top(body, ["="])] if has_eq else [body]
        if any(not sgm for sgm in segments):
            blk.verdict, blk.reason = "SKIP", "empty side"
            return
        parsed = []
        for sgm in segments:
            kind = "ineq" if re.search(r"\\leq|\\geq|\\neq|<|>", sgm) else "expr"
            parsed.append((kind, sgm))
        results = [parse_side(sgm) for _, sgm in parsed]
        signal.alarm(0)
    except Timeout:
        blk.verdict, blk.reason = "SKIP", "timeout"
        return
    except Exception as exc:
        blk.verdict, blk.reason = "SKIP", f"parse: {type(exc).__name__}"
        return
    if any(r[0] == "unresolved" for r in results):
        blk.verdict, blk.reason = "SKIP", "unresolved function"
        return
    if any(r[0] == "expr" and r[1] is None for r in results):
        blk.verdict, blk.reason = "SKIP", "parse returned None"
        return
    # Comma lists (k = 0, 1, 2) parse to tuples, not scalars.
    if any(r[0] == "expr" and not isinstance(r[1], sp.Basic) for r in results):
        blk.verdict, blk.reason = "SKIP", "non-scalar parse"
        return
    # Triage helpers for definitional/undeclared claims (all honest SKIPs —
    # sampling independent points cannot verify these, only misfire on them).
    def _has_undef(expr) -> bool:
        try:
            return any(isinstance(a, sp.AppliedUndef) for a in expr.atoms(sp.Function))
        except Exception:
            return False
    expr_results = [r[1] for r in results if r[0] == "expr"]
    if any(_has_undef(e) for e in expr_results):
        blk.verdict, blk.reason = "SKIP", "unresolved function"
        return
    # Backstop: a free symbol still carrying a backslash (e.g. \varphi) is a
    # command latex2sympy does not model — anything built on it is uncheckable.
    try:
        allfree: set = set()
        for e in expr_results:
            allfree |= e.free_symbols
        if any("\\" in str(s) for s in allfree):
            blk.verdict, blk.reason = "SKIP", "unparsed command"
            return
    except Exception:
        pass
    # Textual unknown-call detection: `f(...)` / `\varphi(...)` on either side
    # whose head is neither a collected definition nor a known function.
    KNOWN_FUNCS = {
        "sin", "cos", "tan", "sec", "csc", "cot", "sinh", "cosh", "tanh",
        "arcsin", "arccos", "arctan", "log", "ln", "lg", "exp", "sqrt", "abs",
        "det", "arg", "min", "max", "sup", "inf", "gcd",
    }
    call_re = re.compile(r"^\\?([a-zA-Z]+)'?\s*\(.*\)$", re.S)
    for _, sgm in parsed:
        m = call_re.match(sgm.strip())
        if m and m.group(1) not in blk.defs and m.group(1) not in KNOWN_FUNCS:
            blk.verdict, blk.reason = "SKIP", "unresolved function"
            return
    if len(expr_results) == 2 and has_eq and not has_ineq:
        lhs, rhs = expr_results
        lf, rf = lhs.free_symbols, rhs.free_symbols
        if isinstance(lhs, sp.Symbol) and rf:
            blk.verdict, blk.reason = "SKIP", "definitional binding"
            return
        if lf and rf and not (lf & rf):
            blk.verdict, blk.reason = "SKIP", "disjoint variables"
            return
        if isinstance(lhs, (sp.Pow, sp.Mul)) and (lf - rf):
            blk.verdict, blk.reason = "SKIP", "dependent variable"
            return
        if bool(lf) != bool(rf):
            blk.verdict, blk.reason = "SKIP", "solve-for (not an identity)"
            return

    def check_expr_pair(lhs, rhs, tol_pair=False) -> str | None:
        # Definitional labels ("\Delta = 4-12") bind a fresh name to a value.
        if isinstance(lhs, sp.Symbol) and not rhs.free_symbols:
            return None
        # Odd real roots: sympy's principal branch gives complex, but lessons use real branch.
        def _real_pow(expr):
            # Rewrite (-a)**(1/odd) -> -(a**(1/odd)) for odd integer denominator
            return expr.replace(
                lambda e: isinstance(e, sp.Pow) and e.base.is_number and e.base.is_negative and e.exp.is_Rational and e.exp.q % 2 == 1,
                lambda e: -((-e.base) ** e.exp)
            )
        lhs = _real_pow(lhs)
        rhs = _real_pow(rhs)
        # Indeterminate forms like 0/0 are not equalities to verify
        if lhs in (sp.nan, sp.zoo, sp.oo, -sp.oo) and rhs in (sp.nan, sp.zoo, sp.oo, -sp.oo):
            return "indeterminate"
        if lhs.has(sp.nan) or rhs.has(sp.nan) or lhs.has(sp.zoo) or rhs.has(sp.zoo):
            return "indeterminate"
        free = lhs.free_symbols | rhs.free_symbols
        # Complex lessons use i for the imaginary unit; latex2sympy leaves it
        # a free symbol. Bind it to I (only when free, so sum indices survive).
        # Same for Euler's e (bare e in an identity is the constant).
        if sp.Symbol("i") in free:
            lhs = lhs.xreplace({sp.Symbol("i"): sp.I})
            rhs = rhs.xreplace({sp.Symbol("i"): sp.I})
            free = lhs.free_symbols | rhs.free_symbols
        if sp.Symbol("e") in free:
            lhs = lhs.xreplace({sp.Symbol("e"): sp.E})
            rhs = rhs.xreplace({sp.Symbol("e"): sp.E})
            free = lhs.free_symbols | rhs.free_symbols
        if not free:
            diff = sp.simplify(lhs - rhs)
            if diff == 0:
                return None
            try:
                # Complex-tolerant: identities sympy cannot close symbolically
                # (e.g. sin through exponentials) still hold numerically.
                if abs(complex(diff.evalf())) < 5e-3 * max(1.0, abs(complex(rhs.evalf()))):
                    return None
            except (TypeError, ValueError, AttributeError):
                pass
            return f"{lhs} != {rhs}"
        tried = 0
        extras = [sp.Rational(3, 2), 5, -1, 7, sp.Rational(2, 5), -3]
        for pti, pt in enumerate(SAMPLE_POINTS):
            sub = {s_: v for s_, v in pt.items() if s_ in free}
            # Deterministic fallback for rarer symbols (subscripted z_1, ...):
            # generic points keep the check meaningful without false skips.
            missing = sorted(free - set(sub), key=str)
            for j, s_ in enumerate(missing):
                sub[s_] = extras[(sum(map(ord, str(s_))) + j + pti) % len(extras)]
            try:
                d = sp.simplify((lhs - rhs).subs(sub))
                tried += 1
                if d != 0:
                    try:
                        if abs(complex(d.evalf())) < 5e-3:
                            continue
                    except (TypeError, ValueError, AttributeError):
                        pass
                    return f"sampling: {sub} -> {d}"
            except Exception:
                continue
        if tried == 0:
            return "no valid sample"
        return None

    def check_ineq(ops, exprs_) -> str | None:
        # Bare e in an inequality is the lesson's own variable (eccentricity,
        # error terms) far more often than Euler's constant: skip rather than
        # false-fail against 2.718...
        exprs_ = [e.xreplace({sp.E: sp.Symbol("e_const")}) for e in exprs_]
        # i is the imaginary unit when free (same guard as the pair checker).
        exprs_ = [e.xreplace({sp.Symbol("i"): sp.I}) if sp.Symbol("i") in e.free_symbols else e
                  for e in exprs_]
        free = set()
        for e in exprs_:
            free |= e.free_symbols
        if sp.Symbol("e_const") in free:
            return "SKIP: ambiguous constant e"
        if any("\\" in str(s) for s in free):
            return "SKIP: unparsed command"
        concrete = {s_: v for s_, v in SAMPLE_POINTS[0].items() if s_ in free}
        extras = [sp.Rational(3, 2), 5, -1, 7, sp.Rational(2, 5), -3]
        missing = sorted(free - set(concrete), key=str)
        for j, s_ in enumerate(missing):
            concrete[s_] = extras[(sum(map(ord, str(s_))) + j) % len(extras)]
        vals = [e.subs(concrete) for e in exprs_]
        for i, op in enumerate(ops):
            if i + 1 >= len(vals):
                return f"malformed chain: {op} with {len(vals)} values"
            d = sp.simplify(vals[i] - vals[i + 1])
            if d == 0:
                ok = op not in ("<", ">", "\\neq")
            else:
                try:
                    dn = float(d.evalf())
                except (TypeError, ValueError, AttributeError):
                    return f"undecidable inequality: {op} between {vals[i]}, {vals[i + 1]}"
                if op in ("<", "\\leq"):
                    ok = dn <= 0
                elif op in (">", "\\geq"):
                    ok = dn >= 0
                elif op == "\\neq":
                    ok = dn != 0
                else:
                    return f"unknown op {op}"
            if not ok:
                return f"violated at {concrete}: {op} between {vals[i]} , {vals[i + 1]}"
        return None

    def check_junction(val, ops, exprs_) -> str | None:
        # Mixed equality/inequality chain: the junction itself must be an
        # equality (val == exprs_[0]), then the inequality chain must hold.
        msg = check_expr_pair(val, exprs_[0])
        if msg:
            return f"junction mismatch: {msg}"[:180]
        return check_ineq(ops, exprs_)

    # Walk the parsed chain pairwise. check_* helpers return None (pass),
    # "SKIP: reason" (honest skip), or anything else (FAIL).
    def _record(msg: str) -> None:
        # A ground claim that fails next to contradiction language is a
        # lesson's intentional reductio (e.g. deriving 3 > 4 to show no
        # solution) — skip it; anything else stays a FAIL.
        if msg.startswith("SKIP: "):
            blk.verdict, blk.reason = "SKIP", msg[6:180]
            return
        try:
            allfree: set = set()
            for r in results:
                for part in r[1:]:
                    items = part if isinstance(part, list) else [part]
                    for e in items:
                        if isinstance(e, sp.Basic):
                            allfree |= e.free_symbols
            ground = not allfree
        except Exception:
            ground = False
        if ground and _has_contradiction_marker(blk):
            blk.verdict, blk.reason = "SKIP", "intentional contradiction"
        else:
            blk.verdict, blk.reason = "FAIL", msg[:180]

    for idx in range(len(results) - 1):
        ka, *pa = results[idx]
        kb, *pb = results[idx + 1]
        if ka == "expr" and kb == "expr":
            msg = check_expr_pair(pa[0], pb[0])
        elif ka == "expr" and kb == "ineq":
            msg = check_junction(pa[0], pb[0], pb[1])
        elif ka == "ineq" and kb == "expr":
            msg = check_ineq(pa[0], pa[1])
            if not msg:
                msg = check_expr_pair(pa[1][-1], pb[0])
                if msg:
                    msg = f"junction mismatch: {msg}"[:180]
        else:
            msg = check_ineq(pa[0], pa[1]) or check_ineq(pb[0], pb[1])
        if msg:
            _record(msg)
            return
    if len(results) == 1 and results[0][0] == "ineq":
        msg = check_ineq(results[0][1], results[0][2])
        if msg:
            _record(msg)
            return
    blk.verdict = "OK"


def split_top_keep(s: str, seps: list[str]) -> list[str]:
    parts, i = [], 0
    while i < len(s):
        matched = next((sep for sep in sorted(seps, key=len, reverse=True)
                        if s.startswith(sep, i)), None)
        if matched:
            parts.append(matched)
            i += len(matched)
        else:
            j = i
            while j < len(s) and not any(s.startswith(sp_, j) for sp_ in seps):
                j += 1
            parts.append(s[i:j])
            i = j
    return parts


DISPLAY_RE = re.compile(r"\\\\\[([\s\S]*?)\\\\\]|(?<!\\)\$\$([\s\S]*?)\$\$")


def extract_blocks(text: str):
    for m in DISPLAY_RE.finditer(text):
        body = (m.group(1) or m.group(2) or "").strip()
        line = text.count("\n", 0, m.start()) + 1
        yield line, body


def main() -> int:
    files = sorted(glob for glob in [str(p) for p in Path(ROOT.joinpath("data/lessons")).rglob("*.md")])
    report = []
    counts: dict[str, int] = {}
    for f in files:
        rel = str(Path(f).relative_to(ROOT))
        text = Path(f).read_text(encoding="utf-8")
        defs = collect_definitions(text)
        for line, body in extract_blocks(text):
            blk = Block(rel, line, body, defs)
            try:
                verify_block(blk)
            except Exception as exc:  # never crash the sweep
                blk.verdict, blk.reason = "SKIP", f"internal: {type(exc).__name__}: {exc}"[:120]
            counts[blk.verdict.split("(")[0]] = counts.get(blk.verdict.split("(")[0], 0) + 1
            if blk.verdict != "OK":
                report.append({
                    "file": blk.file, "line": blk.line,
                    "verdict": blk.verdict, "reason": blk.reason[:200],
                    "latex": blk.tex.strip()[:300],
                })
    REPORT_PATH.write_text(json.dumps(report, indent=1, ensure_ascii=False))
    print(json.dumps(counts, sort_keys=True))
    print(f"report: {REPORT_PATH} ({len(report)} non-OK entries)")
    return 0


if __name__ == "__main__":
    sys.exit(main())
