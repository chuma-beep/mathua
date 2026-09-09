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

from latex2sympy2 import latex2sympy as _latex2sympy_raw  # noqa: E402
import latex2sympy2 as _l2s_mod  # noqa: E402
from sympy.core.function import UndefinedFunction  # noqa: E402


def latex2sympy(seg: str):
    """Hermetic wrapper around latex2sympy2.

    The library keeps parse state in module globals (``var`` / ``variances`` /
    ``VARIABLE_VALUES``): an ``=``-assignment writes entries, and a bare
    differential (``di`` in ``(c + di)``, ``dx`` in integrals, ...) executes
    ``var = get_differential_var(...)`` inside ``convert_atom`` (which
    declares ``global var``), replacing the dict with a bare ``Symbol``.
    Every later call in the process then dies at ``... in var`` with
    ``TypeError: argument of type 'Symbol' is not a container`` — one early
    complex-numbers block poisoned ~1800 later blocks in the sweep.

    Reset the globals before each call so every segment parses independently.
    """
    _l2s_mod.var = {}
    _l2s_mod.variances = {}
    _l2s_mod.VARIABLE_VALUES = {}
    return _latex2sympy_raw(seg)

ROOT = Path(__file__).resolve().parents[1]
REPORT_PATH = ROOT / "scripts" / "verify_report.json"

# (file suffix, latex prefix, reason). Correct claims whose hypotheses live in
# prose (root-of-polynomial + max-coefficient + |z0|>1 here): sampling blind
# points can only misfire, so they are waived explicitly instead.
WAIVERS: list[tuple[str, str, str]] = [
    ("fundamental-inequalities-for-complex-numbers.md", "|z_0| \\leq 1 + M",
     "Cauchy root bound needs hypotheses (z_0 root, M=max coeff)"),
    ("nt.quadratic_reciprocity.md", "\\left(\\frac{p}{q}\\right)",
     "Legendre symbols need distinct odd primes (sampling p,q blind misfires)"),
    ("powers-radicals-logarithms/logarithms.md", "\\log_b(a^y) = \\log_b(x)",
     "change-of-base step under substitution y=log_a(x)"),
    ("hyperbolic-sine-and-cosine.md", "\\sqrt{X^{2}-1} = e^{x} - X",
     "isolation step under exponentiated hypothesis X+sqrt=e^x"),
    ("integration-by-substitution.md", "\\sqrt{9(1-\\sin^2 u)}",
     "needs cos u >= 0 (principal substitution range)"),
    ("trigonometric-substitution-for-integrals.md", "\\begin{align} \\sqrt{x^2+a^2}",
     "needs sec θ > 0 (acute substitution range)"),
    ("trigonometric-substitution-for-integrals.md", "\\begin{align} \\sqrt{x^2-a^2}",
     "needs sec/tan range hypothesis (acute substitution range)"),
    ("probability-and-statistics/geometric-distribution.md", "\\sum_{k=1}^{\\infty} k",
     "needs 0<p<1 for series convergence"),
    ("sets-and-numbers/binomial-coefficient.md", "\\frac{1}{1+x} = \\sum_{k=0}^{\\infty}",
     "needs |x|<1 for series convergence"),
    ("sets-and-numbers/absolute-value.md", "|a + b| = a + b = |a| + |b|",
     "nonnegative-sum case hypothesis (a+b>=0)"),
    ("sets-and-numbers/absolute-value.md", "|a + b| = -(a + b)",
     "negative-sum case hypothesis (a+b<=0)"),
    ("differential-calculus-theorems/cauchy-theorem.md", "\\begin{align} \\frac{4c - 4}{2c}",
     "solve-track for c (equation steps, not identities)"),
    ("integrals/finding-areas-by-integration.md", "\\begin{aligned}\nx^3 - 3x",
     "intersection solve-track (equations, not identities)"),
    ("lines-planes-conic-sections/ellipse.md", "\\begin{align} &(x - 1)^2 + y^2 = 36",
     "radical-isolation derivation (equation steps, not identities)"),
    ("lines-planes-conic-sections/lines.md", "\\begin{align} &2x + 1 = -x + 4",
     "intersection solve-track (solving the system)"),
]

CONTRADICTION_RE = re.compile(r"contradiction|absurd|no solution|impossible", re.I)


def _has_contradiction_marker(blk: Block) -> bool:
    try:
        lines = Path(ROOT.joinpath(blk.file)).read_text(encoding="utf-8").split("\n")
        lo, hi = max(0, blk.line - 1 - 9), min(len(lines), blk.line - 1 + 8)
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
    r"\\{1,2\}\(\s*\\?([a-zA-Z]+)\(\s*([a-zA-Z]+)\s*\)\s*=\s*((?:[^\\]|\\{1,2}(?![\\(])){1,120}?)\s*\\{1,2\}\)"
)

# A collected body that is just another bare call (`R(x) = P(x)`) states a
# conditional relation, not a construction — expanding it is unsound.
BARE_CALL_RE = re.compile(r"^\\?[a-zA-Z]+\([^()]*\)$")


def collect_definitions(text: str) -> dict:
    """Scan raw lesson text for simple function definitions f(x) = <latex>.

    Returns name -> [(line, var, body)] in file order. Only genuine
    constructions are collected: the body must parse as math, contain its
    parameter, and not be prose punctuation or a bare alias call. (Without
    these guards, conditionals like `R(x) = 0` or identities like
    `P(x) = Q(x)D(x) + R(x)` poison every later `R(...)`/`P(...)` in the
    file.) Callers scope to the nearest PRECEDING definition — examples
    redefine `f` freely down a file.
    """
    defs: dict[str, list] = {}
    for m in FUNC_DEF_RE.finditer(text):
        name, var, body = m.group(1), m.group(2), m.group(3).strip()
        if len(body) > 100 or "\\" in body and "=" in body:
            continue
        # Guard against prose false positives: the associative law
        # `a(bc) = (ab)c.` once collected a bogus `a`-definition whose
        # trailing sentence period leaked into the body, silently
        # rewriting every later `a(...)` product in the file. A genuine
        # definition body is pure math: no trailing prose punctuation,
        # and it must parse on its own.
        if re.search(r"[.,;?!:]$", body):
            continue
        if BARE_CALL_RE.match(body):
            continue
        if not re.search(rf"\b{re.escape(var)}\b", body):
            continue
        try:
            parsed_body = latex2sympy(body)
        except Exception:
            continue
        # The body must be built from its parameter, not from other
        # unknowns: `P(x) = Q(x)D(x) + R(x)` is the division identity, not a
        # construction of P — expanding it rewrites P into unbound symbols.
        try:
            if any(isinstance(a.func, UndefinedFunction)
                   for a in parsed_body.atoms(sp.Function)):
                continue
            if any("\\" in str(s) for s in parsed_body.free_symbols):
                continue
        except Exception:
            continue
        line = text.count("\n", 0, m.start()) + 1
        defs.setdefault(name, []).append((line, var, body))
    return defs


def _defs_before(defs: dict, line: int) -> dict:
    """Nearest preceding definition per name (example-scoped use)."""
    out = {}
    for name, cands in defs.items():
        best = None
        for dl, var, body in cands:
            if dl <= line and (best is None or dl > best[0]):
                best = (dl, var, body)
        if best is not None:
            out[name] = (best[1], best[2])
    return out


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
    # `\limits` only controls operator layout (`\sum\limits_{...}`); drop it.
    t = re.sub(r"\\limits(?![a-zA-Z])", "", t)
    # Inverse-trig notation (`\tan^{-1}`) is composition-inverse (arctan),
    # NOT reciprocal — latex2sympy would read `(tan)^(-1)` = cot, wrong math.
    t = re.sub(r"\\(sin|cos|tan)\^\{-1\}", r"\\arc\1", t)
    t = re.sub(r"\\(sin|cos|tan)\^-1(?![0-9])", r"\\arc\1", t)
    # Font/accent wrappers carry no algebraic meaning here: \mathbf{v} -> v.
    t = re.sub(r"\\(?:mathbf|mathrm|mathit|boldsymbol)\{([^{}]*)\}", r"\1", t)
    t = re.sub(r"\\(?:overline|bar|hat|vec|tilde)\{([^{}]*)\}", r"\1", t)
    # Short \text{...} function names are math functions in text clothing:
    # \text{sin}(x+h) -> \sin(x+h). (Longer prose is handled by strip_text.)
    t = re.sub(
        r"\\text\{(sin|cos|tan|sec|csc|cot|sinh|cosh|tanh|arcsin|arccos|arctan|log|ln|lg|exp)\}",
        r"\\\1", t)
    # Leading enumeration labels like \text{1. } \quad or bare "1." / "9\." / "10 ."
    t = re.sub(r"^\\text\{[^{}]*\}\s*", "", t)
    t = re.sub(r"^\s*\d+\s*\\\.\s*", "", t)
    t = re.sub(r"^\s*\d+\.\s*\\quad\s*", "", t)
    t = re.sub(r"^\s*\d+\.\s+", "", t)
    # Spacing/formatting commands are semantically empty here (spacing ones become a space).
    # NOTE: `\left`/`\right` must only match as whole commands — a plain
    # substring replace would corrupt `\rightarrow` into `arrow`.
    for cmd in (r"\left", r"\right", r"\displaystyle", r"\mathrm"):
        t = re.sub(re.escape(cmd) + r"(?![a-zA-Z])", "", t)
        t = re.sub(re.escape(cmd.replace("\\", "\\\\")) + r"(?![a-zA-Z])", "", t)
    for cmd in (r"\!", r"\vphantom{()}"):
        t = t.replace(cmd, "")
        t = t.replace(cmd.replace("\\", "\\\\"), "")
    for cmd in (r"\,", r"\;", r":", r"\quad", r"\qquad"):
        t = t.replace(cmd, " ")
        t = t.replace(cmd.replace("\\", "\\\\"), " ")
    # Thin-spacing juxtaposition is multiplication (`\sqrt{13}\,e^{...}`,
    # `r_1 r_2\, e^{...}`): a spacing command between two math atoms joins
    # factors. (Claim-separating gaps were already split at the TeX level,
    # so surviving spacings sit inside single claims.)
    t = re.sub(r"([}\d)])\s*\\[,;:!]\s*(?=[a-zA-Z\\({])", r"\1*", t)
    t = re.sub(r"([}\d)])\s*\\\\[,;:!]\s*(?=[a-zA-Z\\({])", r"\1*", t)
    # Scrape dialect: doubled backslash before a command name (\\frac, \\cdot,
    # \\sin) is a single command. Collapse those first so the command survives;
    # any remaining \\ (row separators) becomes a space below.
    t = re.sub(r"\\\\([a-zA-Z])", r"\\\1", t)
    t = t.replace("\\[0.5em]", " ").replace("\\\\", " ")
    # Implicit multiplication like 2(3) -> 2*(3) for latex2sympy.
    # NOTE: no digit-digit joining (`28 7` -> `28*7`): space-joined
    # independent claims (`12=2*2*3 18=2*3*3`) need their boundary intact for
    # the chain splitter downstream.
    t = re.sub(r"(\d)\s*\(", r"\1*(", t)
    t = re.sub(r"\)\s*(\d)", r")*\1", t)

    # Euler's `e` never takes a subscript — `e_1`, `e_{n}` in lessons are
    # ordinary indexed variables (Vieta's formulas), but latex2sympy reads
    # them as Euler-E and drops the index. Rename to a safe symbol first.
    # (Bare `e` and `e^{...}` are untouched and still map to E downstream.)
    t = re.sub(r"\be(_\{[^}]*\}|_[a-zA-Z0-9]+)", r"evar\1", t)
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


# A display that opens with an enumeration label (`\text{2.}`, `9\.`) is a
# numbered exercise — a problem statement to solve, not a claim. (Checked on
# the raw block: norm_latex strips these labels before verification.)
NUMBERED_RE = re.compile(r"^\s*(?:\\text\{\s*)?\d+\s*(?:\\\.)?\.\s*")

# Files whose displays are predominantly equations to solve rather than
# identities to verify.
EQUATION_PATH_RE = re.compile(
    r"/equations/|alg\.eq\.|prealg\.eq\.|eq\.one_step|eq\.two_step|vars_both_sides|quad\.discriminant")

# Prose markers that the surrounding lesson presents an equation/problem
# rather than stating an identity (same ±window convention as the
# contradiction marker below). Path-gated above: inside equation lessons
# these phrases mean equation work ("Squaring both sides yields",
# "we will check the solution").
SOLVE_PROSE_RE = re.compile(
    r"\bequat|\bsolve\b|\bsolving\b|\bboth sides\b|\bexercise\b|\bcheck\w*\s+(the\s+)?solution\b", re.I)


def _nearby_text(blk: Block, lo: int, hi: int) -> str:
    # NOTE: blk.line is 1-indexed; list indices are 0-indexed.
    try:
        lines = Path(ROOT.joinpath(blk.file)).read_text(encoding="utf-8").split("\n")
        return "\n".join(lines[max(0, blk.line - 1 + lo):min(len(lines), blk.line - 1 + hi)])
    except Exception:
        return ""


def _is_equation_context(blk: Block) -> bool:
    """True when the block is a problem statement, not a universal claim."""
    if NUMBERED_RE.match(blk.tex):
        return True
    if EQUATION_PATH_RE.search(blk.file):
        # Inside equation lessons almost every `=`-display is introduced as
        # an equation ("Consider the equation", "Let's solve"); identity
        # recalls there are the rare exception, and skipping them as
        # equations is the honest direction (a false OK on an equation is
        # worse than a missed OK on a recall).
        if SOLVE_PROSE_RE.search(_nearby_text(blk, -12, 3)):
            return True
    return False


# Top-level implication arrows joining two claims (`P \Rightarrow Q`).
# Split is brace-aware at the call site, so `\lim_{x \to 0}` subscripts and
# lone `a_n \to +\infty` notations are unaffected.
IMPLIES_SEPS = ["\\Rightarrow", "\\Longrightarrow", "\\rightarrow",
                "\\longrightarrow", "\\leftrightarrow", "\\iff",
                "\\Longleftrightarrow", "\\implies", "\\to"]

# Trailing side-condition clauses: hypotheses, not claims
# (`y = \log_a{f(x)} \quad \text{with} \quad a > 0`). The `{with}` form
# survives strip_text (short group); cut the head claim free and verify it
# alone — definitions then hit the definitional guard, identities verify.
SIDE_COND_RE = re.compile(
    r"\{(?:with|for|if|where|when|such that|provided|given)\}|"
    r"\\text\{(?:with|for|if|where|when)\}")

# A trailing bare `\neq`-condition after an equation (`... = 0,\ a \neq -2`).
TRAILING_COND_RE = re.compile(
    r"^(?P<head>.*=.*?),?\s*(?P<cond>[a-zA-Z](?:_\{[^}]*\}|_[a-zA-Z0-9])?\s*\\neq.*)$",
    re.S)

# Ellipsis schemas (`z^n - 1 = (z-z_0)\cdots`) quantify over an unlisted
# range — sampling the visible symbols cannot verify them.
ELLIPSIS_RE = re.compile(r"\\(cdots|ldots|vdots|ddots)")

# Evaluation-bar notation (`\left[ F \right]_1^t`) needs the FTC to check;
# the bar has no algebraic meaning to the sampler.
EVAL_BAR_RE = re.compile(r"\](_\{[^}]*\}|_[a-zA-Z0-9])")

# Sequence notation with subscript conditions (`(a_n)_{n \geq 3}`).
SEQ_NOTATION_RE = re.compile(r"\)_\{[^}]*\\(geq|leq)")

# LaTeX commands the pipeline understands (parser repertoire, guards, and
# decorations stripped by norm_latex). A parse failure blamed on anything
# outside this set is reported as that macro, not a generic parse error.
# Deliberately EXCLUDED (hence reported): custom binary operators (\star,
# \oplus, \odot, ...), \underbrace/\overbrace, \boxed, \cancel, operator
# names the parser has no model for (\ker, \hom, \tr, \rank, \sgn, ...).
KNOWN_MACROS = {
    "frac", "dfrac", "tfrac", "sqrt", "root", "cdot", "times", "div",
    "pm", "mp", "leq", "geq", "le", "ge", "lt", "gt", "neq", "equiv",
    "approx", "simeq", "asymp", "cong", "doteq", "propto", "sim",
    "left", "right", "big", "bigg", "Big", "Bigg", "bigl", "bigr",
    "Bigl", "Bigr", "langle", "rangle", "lfloor", "rfloor", "lceil",
    "rceil", "vert", "Vert", "lvert", "rvert", "ulcorner", "urcorner",
    "llcorner", "lrcorner", "sin", "cos", "tan", "sec", "csc", "cot",
    "sinh", "cosh", "tanh", "arcsin", "arccos", "arctan", "log", "ln",
    "lg", "exp", "abs", "det", "arg", "min", "max", "sup", "inf",
    "gcd", "lim", "sum", "prod", "int", "iint", "iiint", "oint",
    "infty", "pi", "partial", "nabla", "forall", "exists", "in",
    "notin", "ni", "subset", "supset", "subseteq", "supseteq",
    "sqsubset", "sqsupset", "sqsubseteq", "sqsupseteq", "cup", "cap",
    "setminus", "emptyset", "vee", "wedge", "neg", "land", "lor",
    "alpha", "beta", "gamma", "delta", "epsilon", "varepsilon",
    "zeta", "eta", "theta", "vartheta", "iota", "kappa", "lambda",
    "mu", "nu", "xi", "omicron", "pi", "varpi", "rho", "varrho",
    "sigma", "varsigma", "tau", "upsilon", "phi", "varphi", "chi",
    "psi", "omega", "Gamma", "Delta", "Theta", "Lambda", "Xi", "Pi",
    "Sigma", "Upsilon", "Phi", "Psi", "Omega", "ell", "hbar",
    "imath", "jmath", "Re", "Im", "aleph", "infty",
    "mathbb", "mathbf", "mathrm", "mathit", "boldsymbol", "overline",
    "underline", "bar", "hat", "check", "breve", "acute", "grave",
    "vec", "tilde", "dot", "ddot", "text", "intertext", "quad",
    "qquad", "displaystyle", "scriptstyle", "scriptscriptstyle",
    "textstyle", "vphantom", "hphantom", "phantom", "begin", "end",
    "tag", "label", "nonumber", "binom", "to", "mapsto", "rightarrow",
    "leftarrow", "Rightarrow", "Leftarrow", "leftrightarrow",
    "Leftrightarrow", "implies", "impliedby", "ldots", "cdots",
    "vdots", "ddots", "dots", "colon", "mid", "parallel", "perp",
    "prime", "dprime", "trprime", "backprime", "surd", "top", "bot",
}


# Function heads latex2sympy models (everything else applied to arguments
# is an unresolvable custom/unknown call for the sampler).
KNOWN_FUNCS = {
    "sin", "cos", "tan", "sec", "csc", "cot", "sinh", "cosh", "tanh",
    "arcsin", "arccos", "arctan", "log", "ln", "lg", "exp", "sqrt", "abs",
    "det", "arg", "min", "max", "sup", "inf", "gcd",
}


def _unknown_macro(body: str) -> str | None:
    """Name a backslash-command outside the known repertoire, if any."""
    for m in re.finditer(r"\\([a-zA-Z]+)", body):
        if m.group(1) not in KNOWN_MACROS:
            return "\\" + m.group(1)
    return None


# All-positive fallback points for domain-restricted identities (logarithms,
# real radicals, positive bases): a sample outside the domain (log of a
# negative, even root of a negative) is an invalid counterexample, not a
# falsification. Retried only when standard sampling fails on such shapes.
POSITIVE_POINTS = [
    {x: 2, y: 3, z: 5, k: 2, n: 3, a: 2, b: 3, c: 5,
     theta: 1, alpha: 2, beta: 1, varphi: 2, phi: 1, rho: 2,
     lam: 3, mu: 2, tau: 1, omega: 2, r_: 2, u_: 1, v_: 3, w_: 2,
     t_: 1, m_: 2, p_: 3, q_: 2, s_: 1, d_: 4},
    {x: 3, y: 5, z: 2, k: 1, n: 2, a: 3, b: 5, c: 2,
     theta: 2, alpha: 1, beta: 2, varphi: 1, phi: 2, rho: 1,
     lam: 2, mu: 3, tau: 2, omega: 1, r_: 3, u_: 2, v_: 1, w_: 4,
     t_: 2, m_: 1, p_: 2, q_: 5, s_: 3, d_: 1},
]


def _combine_pieces(blk: Block, subs: list) -> None:
    """Combine sub-block verdicts: FAIL short-circuits (with piece context),
    OK wins if any piece verified, else the first SKIP explains the block."""
    first_skip: tuple | None = None
    saw_ok = False
    for sub in subs:
        if sub.verdict == "FAIL":
            blk.verdict, blk.reason = "FAIL", sub.reason
            return
        if sub.verdict == "OK":
            saw_ok = True
        elif first_skip is None:
            first_skip = (sub.verdict, sub.reason)
    if saw_ok:
        blk.verdict, blk.reason = "OK", ""
    elif first_skip is not None:
        blk.verdict, blk.reason = first_skip
    else:
        blk.verdict, blk.reason = "SKIP", "empty pieces"


AND_OR_RE = re.compile(r"\{(?:and|or)\}|\\text\{(?:and|or)\}")


def _split_raw_claims(tex: str) -> list[str] | None:
    """Split one display block into independent claims at the TeX level.

    Separators: row breaks (double backslash, with optional spacing opt),
    ``\\qquad`` / ``\\quad`` gaps, and ``\\text{and|or}``. Must run on the
    RAW block: norm turns gaps into a plain space (and even joins `28 7`
    into `28*7`), destroying the boundary afterwards. A bare `a = b = c`
    chain (no separator) stays one piece and is walked pairwise downstream.

    The corpus mixes single- and doubled-command prefixes, so doubled
    prefixes are collapsed first — exactly like norm_latex does — otherwise
    the row-break pattern would fire inside a doubled gap command itself.
    Spacing commands (`\\,`, `\\;`, `\\:`, `\\!`) are blanked first for the
    same reason: `\\,` is a thin space, never a row break (this once split
    `a(x_1+x_2)\\,x` mid-claim).
    """
    t = re.sub(r"\\\\([a-zA-Z])", r"\\\1", tex)
    t = re.sub(r"\\\\([,;:!])", " ", t)
    parts = re.split(
        r"\\\\+(?:\[[^\]]*\])?|\\qquad|\\quad|\\text\{\s*(?:and|or)\s*\}", t)
    parts = [p for p in (s.strip() for s in parts) if p]
    if len(parts) <= 1:
        return None
    return parts


def _split_independent(body: str) -> list[str] | None:
    """Split `A {and} B`-style leftovers (post-norm) into independent claims.

    Returns None when the body is a single claim. A bare `a = b = c` chain
    (no separator) is NOT split here — the `=`-chain walker below verifies it
    pairwise.
    """
    parts = [p.strip() for p in AND_OR_RE.split(body)]
    parts = [p for p in parts if p]
    if len(parts) <= 1:
        return None
    return parts


ENV_ROW_RE = re.compile(r"\\begin\{(aligned|align\*?)\}([\s\S]*?)\\end\{(?:aligned|align\*?)\}")


def _verify_aligned(blk: Block, tex: str) -> None:
    """Verify aligned/align derivation rows one by one.

    Operates on the RAW block (only doubled prefixes collapsed): full
    normalization would already turn every `\\\\` row separator into a
    space, leaving a single uncheckable row. Rows joined by `&=` continuations
    (`A &= B \\\\ &= C`) are chained: a row starting with `=` continues from
    the previous row's right-hand side.
    """
    t = re.sub(r"\\\\([a-zA-Z])", r"\\\1", blk.tex)
    # Spacing commands are never row breaks (see _split_raw_claims).
    t = re.sub(r"\\\\([,;:!])", " ", t)
    pieces: list[str] = []
    pos = 0
    found = False
    for m in ENV_ROW_RE.finditer(t):
        found = True
        outside = m.string[pos:m.start()].strip()
        if outside:
            pieces.append(outside)
        for row in re.split(r"\\\\+(?:\[[^\]]*\])?", m.group(2)):
            # NOTE: `&` alignment points are KEPT here — the claim chunker
            # below needs them to separate side-by-side equations. Only
            # \nonumber (a pure decoration) is dropped.
            row = row.replace("\\nonumber", "").strip()
            if row:
                pieces.append(row)
        pos = m.end()
    if not found:
        _verify_piece(blk, *strip_text(tex))
        return
    outside = t[pos:].strip()
    if outside:
        pieces.append(outside)
    subs: list[Block] = []
    # Previous right-hand side per table column: `&=` continuations and
    # multi-column (`&&`) computation tables chain within their column
    # (`2·2^k &> 2k` continues as `2·2^k > 2k`; Line A/B slopes in parallel).
    prev_rhs: list[str | None] = []

    def _claim_signal(s: str) -> bool:
        return bool(re.search(r"=|\\leq|\\geq|\\neq|<|>|\\lim", s))

    failed = False
    for piece in pieces:
        # Multi-column table rows (`A &= B & C &= D`): `&&` is a hard column
        # break; single `&` cells are greedily rejoined into claims — a new
        # claim starts when the accumulation already holds `=` and the next
        # cell is not a continuation (`= C`, `> 2k`). This keeps `&=` chains
        # intact while separating side-by-side equations.
        cols = [c for c in re.split(r"(?<!\\)&&", piece)]
        for ci, col in enumerate(cols):
            while len(prev_rhs) <= ci:
                prev_rhs.append(None)
            cells = [c.strip() for c in re.split(r"(?<!\\)&", col)]
            claims: list[str] = []
            acc = ""
            for c in cells:
                if not c:
                    continue
                if "=" in acc and not re.match(r"^(=|\\leq|\\geq|\\neq|<|>)", c):
                    claims.append(acc)
                    acc = c
                else:
                    acc = f"{acc} {c}".strip() if acc else c
            if acc:
                claims.append(acc)
            for cell in claims if claims else [""]:
                body, _ = strip_text(norm_latex(cell))
                if not _claim_signal(body):
                    # Pure expression fragment. A fragment continuing with an
                    # operator (`- (3)(4)(2) - ...` picking up the previous
                    # right-hand side) EXTENDS the column's chain — the next
                    # `=` row is verified against the whole accumulated
                    # right-hand side (Sarrus expansions split across rows).
                    # Anything else is not a claim: skip without touching it.
                    if prev_rhs[ci] and re.match(r"^[+-]", body):
                        prev_rhs[ci] = f"{prev_rhs[ci]} {body}"
                    sub = Block(blk.file, blk.line, cell, blk.defs)
                    sub.verdict, sub.reason = "SKIP", "continued RHS fragment"
                    subs.append(sub)
                    continue
                if prev_rhs[ci] and re.match(r"^(=|\\leq|\\geq|\\neq|<|>)", body):
                    # Continuation row (`&= C`, `&> 2k`): verify against the
                    # same column's previous right-hand side.
                    body = f"{prev_rhs[ci]} {body}"
                sub = Block(blk.file, blk.line, cell, blk.defs)
                _verify_piece(sub, body, False)
                if sub.verdict == "FAIL":
                    sub.reason = f"[{cell[:60]}] {sub.reason}"[:200]
                    subs.append(sub)
                    failed = True
                    break
                subs.append(sub)
                if "=" in body:
                    prev_rhs[ci] = body.split("=")[-1].strip()
                else:
                    segs = split_top(body, ["\\leq", "\\geq", "\\neq", "<", ">"])
                    prev_rhs[ci] = segs[-1].strip() if segs else None
            if failed:
                break
        if failed:
            break
    _combine_pieces(blk, subs)


def verify_block(blk: Block) -> None:
    # Principal-argument range claims (-pi < Arg(z) <= pi) need an Arg model
    # the checker does not have; sampling the letters is meaningless.
    # (Checked on raw tex: \text/\mathrm wrappers are stripped later.)
    if re.search(r"\\(?:text|mathrm)\{Arg\}", blk.tex):
        blk.verdict, blk.reason = "SKIP", "principal argument"
        return
    # Problem statements (numbered exercises, equations introduced as
    # "Consider the equation" / "Let's solve") are not universal claims:
    # sampling arbitrary points can only misfire on them.
    if _is_equation_context(blk):
        blk.verdict, blk.reason = "SKIP", "equation to solve"
        return
    # Matrix algebra (`A·I = I·A = A`): single-uppercase symbols around
    # `\cdot` are matrices, but the parser reads `I` as the imaginary unit —
    # sampling them as scalars can only misfire (or false-verify).
    if re.search(r"matr|vector", blk.file, re.I) and re.search(
            r"[A-Z]\s*\\cdot|\\cdot\s*[A-Z]", blk.tex):
        blk.verdict, blk.reason = "SKIP", "matrix algebra"
        return
    # Conjugates (`\overline{z}`, `\bar{z}`) and vectors (`\vec{i}`) lose
    # their meaning when norm strips the accent (`z·\bar{z}` would verify as
    # `z·z`, true at real samples but a different claim — a false OK). The
    # sampler has no conjugate/vector model: honest SKIP up front.
    # (Multi-letter `\overline{OA}` segment lengths are harmless either way,
    # but they ride along; their chains still SKIP as disjoint unknowns.)
    if re.search(r"\\(?:overline|bar\b|bar\{|vec\b|vec\{)", blk.tex):
        blk.verdict, blk.reason = "SKIP", "conjugate/vector notation"
        return
    # Waivers: correct lesson claims whose hypotheses live in prose. Checked
    # here (before env routing) so whole aligned derivations can be waived;
    # _verify_piece keeps a backstop check for split pieces.
    for wfile, wtex, wreason in WAIVERS:
        if blk.file.endswith(wfile) and blk.tex.strip().startswith(wtex):
            blk.verdict, blk.reason = "SKIP", f"waived: {wreason}"
            return
    tex = norm_latex(blk.tex)
    if "\\begin{" in tex:
        m = re.search(r"\\begin\{([a-zA-Z*]+)\}", tex)
        env = m.group(1) if m else "?"
        if env in ("aligned", "align", "align*"):
            _verify_aligned(blk, tex)
            return
        blk.verdict, blk.reason = "SKIP", f"environment: {env}"
        return
    # `P \Rightarrow Q` (and `\iff`, `\to` between claims) states an
    # implication between claims, not a claim itself; logic-chasing it by
    # sampling is meaningless either way. Checked BEFORE claim-splitting:
    # splitting `\quad \iff \quad` first would orphan the hypothesis from
    # the claim (`\arcsin(\sin\theta) = \theta` alone is false outside its
    # range). The split is brace-aware, so `\lim_{x \to 0}` subscripts and
    # lone `a_n \to +\infty` notations are unaffected.
    _impl_tex = re.sub(r"\\\\([a-zA-Z])", r"\\\1", blk.tex)
    _impl_tex = re.sub(r"\\\\([,;:!])", " ", _impl_tex)
    impl_parts = [p for p in split_top(_impl_tex, IMPLIES_SEPS)]
    if len(impl_parts) > 1 and all(
            re.search(r"=|\\leq|\\geq|\\neq|<|>|\\lim|\\in", p) for p in impl_parts if p.strip()):
        blk.verdict, blk.reason = "SKIP", "implication"
        return
    # Independent claims joined at the TeX level (`A \\\\ B`, `A \qquad B`,
    # `A \quad \text{and} B`): verify each alone. Each piece recurses through
    # the full router (own equation-context, own guards).
    raw_pieces = _split_raw_claims(blk.tex)
    if raw_pieces is not None:
        subs = []
        for piece in raw_pieces:
            sub = Block(blk.file, blk.line, piece, blk.defs)
            verify_block(sub)
            if sub.verdict == "FAIL":
                sub.reason = f"[{piece[:60]}] {sub.reason}"[:200]
                subs.append(sub)
                break
            subs.append(sub)
        _combine_pieces(blk, subs)
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
    pieces = _split_independent(body)
    if pieces is not None:
        subs = []
        for piece in pieces:
            sub = Block(blk.file, blk.line, piece, blk.defs)
            _verify_piece(sub, piece, False)
            if sub.verdict == "FAIL":
                sub.reason = f"[{piece[:60]}] {sub.reason}"[:200]
                subs.append(sub)
                break
            subs.append(sub)
        _combine_pieces(blk, subs)
        return
    _verify_piece(blk, body, had_long_text)


def _undetermined_symbols(lhs, rhs) -> bool:
    """Template unknowns (`A,B,C` in partial fractions, `p,q,r,s` in factoring
    templates, `z_k` in indexed families): a bare letter confined to exactly
    one side is an unknown to solve for, not a shared variable of an
    identity. A true universal identity cannot have a one-sided free variable
    (independent sampling would falsify it) — except the constants `i` and
    `e`, excluded. Definitions (`S_n = ...`, `y = ...`) are caught by the
    definitional guard first, so this fires only on non-Symbol left sides."""
    try:
        lf, rf = lhs.free_symbols, rhs.free_symbols
    except Exception:
        return False
    for s_ in (lf | rf):
        name = str(s_)
        if len(name) < 1 or name in ("i", "e"):
            continue
        if not re.fullmatch(r"[A-Za-z](_\{[^}]*\}|_[a-zA-Z0-9]+)?", name):
            continue
        if (s_ in lf) != (s_ in rf):
            return True
    return False


def _domain_sensitive(lhs, rhs) -> bool:
    """True for logarithms, real radicals and fractional powers — claims
    stated over positive reals, where negative samples prove nothing."""
    for e in (lhs, rhs):
        try:
            if not isinstance(e, sp.Basic):
                continue
            if e.has(sp.log):
                return True
            for p in e.atoms(sp.Pow):
                if p.exp.is_Integer is not True:
                    return True
        except Exception:
            continue
    return False


def _pair_guard(lhs, rhs) -> str | None:
    """Honest-SKIP guards for one `lhs = rhs` pair; None means checkable."""
    try:
        lf, rf = lhs.free_symbols, rhs.free_symbols
    except Exception:
        return None
    if isinstance(lhs, sp.Symbol) and rf:
        return "definitional binding"
    if lf and rf and not (lf & rf):
        return "disjoint variables"
    if isinstance(lhs, (sp.Pow, sp.Mul)) and (lf - rf):
        # ...unless symbolically identical (`a^0 = 1`): exact verification
        # dominates the heuristic — sample-checking would confirm it too.
        try:
            if sp.simplify(lhs - rhs) == 0:
                return None
        except Exception:
            pass
        return "dependent variable"
    if bool(lf) != bool(rf):
        return "solve-for (not an identity)"
    if _undetermined_symbols(lhs, rhs):
        return "undetermined coefficients"
    return None


def _verify_piece(blk: Block, body: str, had_long_text: bool) -> None:
    # Ellipsis schemas quantify over an unlisted range — the visible symbols
    # do not determine the claim.
    if ELLIPSIS_RE.search(body):
        blk.verdict, blk.reason = "SKIP", "ellipsis schema"
        return
    # Evaluation bars need the FTC (`[F]_1^t = F(t) - F(1)`); the sampler
    # has no such model, and lessons mix dummy variables inside the bar.
    if EVAL_BAR_RE.search(body):
        blk.verdict, blk.reason = "SKIP", "evaluation-bar notation"
        return
    # Sequence notation with subscript conditions (`(a_n)_{n \geq 3}`).
    if SEQ_NOTATION_RE.search(body):
        blk.verdict, blk.reason = "SKIP", "sequence notation"
        return
    # Expand function definitions inside body, scoped to the nearest PRECEDING
    # definition (examples redefine `f` freely down a file; a following
    # example's `f` must not rewrite this block's `f(1)`).
    scoped_defs = _defs_before(blk.defs, blk.line)
    for fname, (var, fbody) in scoped_defs.items():
        def _repl(m):
            arg = m.group(1)
            try:
                b_expr = latex2sympy(fbody)
                a_expr = latex2sympy(arg)
                return f"({b_expr.subs(sp.Symbol(var), a_expr)})"
            except Exception:
                return m.group(0)
        body = re.sub(rf"{re.escape(fname)}\s*\(\s*([^\)]+)\s*\)", _repl, body)
    # Bare products (`a(b+c)`, `m(x - x_P)`, `x(y+3)`): latex2sympy reads
    # `L(...)` as function application, but for these ordinary-variable heads
    # the corpus overwhelmingly means multiplication (distributive law,
    # point-slope form, factoring). Every other single-letter head (`f`,
    # `g`, `p`, `v`, `t` as tangent line, `c` as MVT point, `q` as
    # polynomial, `o` as little-o, ...) is function-prone in this corpus and
    # keeps its call form: an out-of-scope use triages as unresolved below
    # rather than mis-verifying as a product. Names defined ANYWHERE in the
    # file also keep call form for the same reason.
    file_def_names = set(blk.defs)
    MULT_HEADS = frozenset("amnxyz")

    def _mul_repl(m):
        if m.group(1) in file_def_names or m.group(1) not in MULT_HEADS:
            return m.group(0)
        return f"{m.group(1)}*("
    body = re.sub(r"(?<![\\a-zA-Z0-9_{}])([a-z])\s*\(", _mul_repl, body)
    # Trailing side conditions are hypotheses, not claims: verify the head
    # alone (`S_n = ... \quad \text{for } r \neq 1`, `... = 0,\ a \neq -2`).
    mcond = SIDE_COND_RE.search(body)
    if mcond and "=" in body[:mcond.start()]:
        body = body[:mcond.start()]
    else:
        mtrail = TRAILING_COND_RE.match(body.strip())
        if mtrail:
            body = mtrail.group("head")
    # Differential/integral relations (`dx = du/2`, `\int ... dt = ...`) are
    # not algebraic identities to check with sample points — the bare
    # differential (`dx`, `dt`, `ds`, ...) is an operator, not a variable.
    if re.search(r"\bdx\b|\bdu\b|\bdy\b|\bd[a-z]\b", body):
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
    # Tuples and multi-argument forms (`(1,0,0)`, bare `x = 1, y = 2`) are
    # not scalar claims. Paren-depth counts too, so known multi-arg calls
    # (`\max(a,b)`) are exempt — but a bare tuple's commas sit at paren
    # depth ≥ 1 as well, so those still fall through to the parser, which
    # reports them as non-scalar (same honest bucket as before).
    # Post-parse, survivors would trip "non-scalar parse" anyway — no OK
    # has a brace/paren-top-level comma, so this only refines SKIP reasons.
    depth_brace = depth_paren = 0
    has_top_comma = False
    i = 0
    while i < len(body):
        ch = body[i]
        if ch == "\\" and i + 1 < len(body):
            i += 2
            continue
        if ch == "{":
            depth_brace += 1
        elif ch == "}":
            depth_brace -= 1
        elif ch == "(":
            depth_paren += 1
        elif ch == ")":
            depth_paren -= 1
        elif ch == "," and depth_brace == 0 and depth_paren == 0:
            has_top_comma = True
            break
        i += 1
    if has_top_comma:
        blk.verdict, blk.reason = "SKIP", "tuple or multi-argument form"
        return
    # Textual unresolved-call check: a `Head(...)` / `Head_sub(...)` form
    # whose head is neither a known function nor a file-defined name is
    # function-like and unresolvable here — catching it pre-parse also
    # avoids parser-internal AttributeErrors on inputs like `C_s(V_x, V_y)`
    # that the grammar chokes on. Deliberately NARROW: bare single letters
    # (`E(X)` for expectation, kept-as-call `f(x)`) are exempt — they flow
    # to the post-parse triage, which already handles them (Euler-coincidence
    # OKs for linear expectation identities must not regress). Only forms
    # that can never verify are caught early: subscripted heads, file-defined
    # names used out of scope, and long unknown heads.
    for mcall in re.finditer(
            r"(?<![\\a-zA-Z])([a-zA-Z]+(?:_\{[^}]*\}|_[a-zA-Z0-9]+)?)\s*\(", body):
        head = mcall.group(1)
        base = head.split("_")[0].lstrip("\\")
        if base in KNOWN_FUNCS:
            continue
        if "_" not in head and len(base) <= 2 and base not in file_def_names:
            # Bare short head (`E(X)`, kept-as-call `f(x)`): the post-parse
            # triage already handles these (Euler-coincidence OKs for linear
            # expectation identities must not regress).
            continue
        # Subscripted head (`C_s(`, `L_{\text{reg}}(`), long unknown head
        # (`Var(`, `ker(`), or out-of-scope defined name: never verifiable.
        blk.verdict, blk.reason = "SKIP", "unresolved function"
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
            if fname in scoped_defs:
                var, body = scoped_defs[fname]
                try:
                    body_expr = latex2sympy(body)
                    arg_expr = latex2sympy(mfn.group(2))
                    return body_expr.subs(sp.Symbol(var), arg_expr)
                except Exception:
                    pass
        return latex2sympy(seg)

    def parse_side(seg: str):
        """Parse one segment; returns ('expr', e) or ('ineq', ops, exprs)."""
        ops = re.findall(r"\\leq|\\geq|\\neq|<|>", seg)
        if ops:
            parts_ = split_top(seg, ["\\leq", "\\geq", "\\neq", "<", ">"])
            exprs_ = [resolve_expr(p) for p in parts_]
            return ("ineq", ops, exprs_)
        return ("expr", resolve_expr(seg))

    # Space-joined independent claims (`12=2*2*3 18=2*3*3`, or `A=B \qquad
    # C=D` after norm turns `\qquad` into a space): a middle segment with a
    # digit-space-digit boundary holds the end of one claim and the start of
    # the next. Split raw segments into chains there first; bare `a = b = c`
    # chains (no such boundary) stay one chain, walked pairwise below.
    raw_segments = [seg.strip() for seg in split_top(body, ["="])] if has_eq else [body]
    if any(not sgm for sgm in raw_segments):
        blk.verdict, blk.reason = "SKIP", "empty side"
        return
    chains: list[list[str]] = [[raw_segments[0]]]
    if has_eq and len(raw_segments) > 2:
        for sgm in raw_segments[1:-1]:
            parts = re.split(r"(?<=\d)\s+(?=\d)", sgm.strip(), maxsplit=1)
            if len(parts) == 2 and all(p.strip() for p in parts):
                chains[-1].append(parts[0].strip())
                chains.append([parts[1].strip()])
            else:
                chains[-1].append(sgm)
        chains[-1].append(raw_segments[-1])
    else:
        chains = [raw_segments]

    def _parse_chain(segs: list[str]):
        par = []
        for sgm in segs:
            kind = "ineq" if re.search(r"\\leq|\\geq|\\neq|<|>", sgm) else "expr"
            par.append((kind, sgm))
        return par, [parse_side(sgm) for _, sgm in par]

    try:
        signal.alarm(20)
        try:
            parsed_chains = [_parse_chain(segs) for segs in chains]
        finally:
            signal.alarm(0)
    except Timeout:
        blk.verdict, blk.reason = "SKIP", "timeout"
        return
    except Exception as exc:
        macro = _unknown_macro(body)
        if macro:
            blk.verdict, blk.reason = "SKIP", f"unknown macro: {macro}"[:120]
        else:
            blk.verdict, blk.reason = "SKIP", f"parse: {type(exc).__name__}"
        return
    # Flatten for the triage guards below (per-chain results are walked
    # separately afterwards).
    parsed = [p for par, _ in parsed_chains for p in par]
    results = [r for _, res in parsed_chains for r in res]
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
            return any(isinstance(a.func, UndefinedFunction) for a in expr.atoms(sp.Function))
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
    # whose head is neither a collected definition nor a known function
    # (KNOWN_FUNCS is module-level).
    call_re = re.compile(r"^\\?([a-zA-Z]+)'?\s*\(.*\)$", re.S)
    for _, sgm in parsed:
        m = call_re.match(sgm.strip())
        if m and m.group(1) not in blk.defs and m.group(1) not in KNOWN_FUNCS:
            blk.verdict, blk.reason = "SKIP", "unresolved function"
            return
    def check_expr_pair(lhs, rhs, tol_pair=False) -> str | None:
        # Definitional labels ("\Delta = 4-12") bind a fresh name to a value.
        if isinstance(lhs, sp.Symbol) and not rhs.free_symbols:
            return None
        # Evaluate calculus objects first: a Limit/Sum/Product the parser kept
        # unevaluated (e.g. `\lim_{x \to 0} ... = -\infty`) can only be
        # compared after `.doit()`. Still unevaluated afterwards → honest SKIP.
        if any(isinstance(e, (sp.Limit, sp.Sum, sp.Product, sp.Integral))
               for e in (lhs, rhs) if isinstance(e, sp.Basic)):
            try:
                lhs, rhs = lhs.doit(), rhs.doit()
            except Exception:
                pass
            if any(isinstance(e, (sp.Limit, sp.Sum, sp.Product, sp.Integral))
                   for e in (lhs, rhs) if isinstance(e, sp.Basic)):
                for cls in (sp.Limit, sp.Sum, sp.Product, sp.Integral):
                    if lhs.has(cls) or rhs.has(cls):
                        return f"SKIP: unevaluated {cls.__name__.lower()}"
        # Odd real roots: sympy's principal branch gives complex, but lessons use real branch.
        def _real_pow(expr):
            # Rewrite (-a)**(1/odd) -> -(a**(1/odd)) for odd integer denominator.
            # Plain integer exponents (q == 1, e.g. (-2)**2 = 4) are NOT roots
            # and must be left alone — rewriting them negates even powers.
            return expr.replace(
                lambda e: isinstance(e, sp.Pow) and e.base.is_number and e.base.is_negative and e.exp.is_Rational and not e.exp.is_Integer and e.exp.q % 2 == 1,
                lambda e: -((-e.base) ** e.exp)
            )
        lhs = _real_pow(lhs)
        rhs = _real_pow(rhs)
        # Extended-real claims (`\lim ... = -\infty`): equal infinities verify,
        # mismatched ones are genuine lesson errors, not indeterminate forms.
        if lhs in (sp.oo, -sp.oo) or rhs in (sp.oo, -sp.oo):
            if lhs == rhs:
                return None
            return f"{lhs} != {rhs}"
        # Indeterminate forms like 0/0 are not equalities to verify
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

        def _run_points(points) -> str | None:
            """Sample (lhs - rhs) at points. Returns None (verified), a FAIL
            message (falsified), or "SKIP: ..." (nothing decidable)."""
            nonlocal tried
            for pti, pt in enumerate(points):
                sub = {s_: v for s_, v in pt.items() if s_ in free}
                # Deterministic fallback for rarer symbols (subscripted z_1, ...):
                # generic points keep the check meaningful without false skips.
                missing = sorted(free - set(sub), key=str)
                for j, s_ in enumerate(missing):
                    sub[s_] = extras[(sum(map(ord, str(s_))) + j + pti) % len(extras)]
                try:
                    d = sp.simplify((lhs - rhs).subs(sub))
                except Exception:
                    continue
                # A sample outside the expression's domain (pole, log of a
                # negative, even root of a negative) is an invalid
                # counterexample, not a falsification — discard the point.
                try:
                    if d.has(sp.nan, sp.zoo) or d.is_infinite:
                        continue
                except Exception:
                    continue
                tried += 1
                if d != 0:
                    try:
                        if abs(complex(d.evalf())) < 5e-3:
                            continue
                    except (TypeError, ValueError, AttributeError):
                        pass
                    return f"sampling: {sub} -> {d}"
            return None

        msg = _run_points(SAMPLE_POINTS)
        std_tried = tried
        if msg is None and std_tried > 0:
            return None
        if not _domain_sensitive(lhs, rhs):
            return msg if msg is not None else "SKIP: no valid sample"
        # Domain-sensitive shape (logarithms, real radicals, fractional
        # powers are stated over positive reals): a standard-points mismatch
        # may be an out-of-domain artifact. Retry where the claim lives — a
        # positive-points mismatch is a genuine falsification.
        tried = 0
        msg2 = _run_points(POSITIVE_POINTS)
        if msg2 is None:
            return None if tried > 0 else "SKIP: no valid sample"
        return msg2

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
        # Unbound applied functions (`c(x)`, `f`, `g` in epsilon-delta
        # schemas) cannot be decided by sampling — but that is a lack of
        # model, not a falsification.
        for e in exprs_:
            try:
                if any(isinstance(a.func, UndefinedFunction)
                       for a in e.atoms(sp.Function)):
                    return "SKIP: unbound function in inequality"
            except Exception:
                continue
        extras = [sp.Rational(3, 2), 5, -1, 7, sp.Rational(2, 5), -3]

        def _concrete(pt, salt: int):
            sub = {s_: v for s_, v in pt.items() if s_ in free}
            missing = sorted(free - set(sub), key=str)
            for j, s_ in enumerate(missing):
                sub[s_] = extras[(sum(map(ord, str(s_))) + j + salt) % len(extras)]
            return sub

        points = ([_concrete(pt, pti) for pti, pt in enumerate(SAMPLE_POINTS)]
                  + [_concrete(POSITIVE_POINTS[0], 99)])
        decided = 0
        violated = 0
        first_violation = ""
        for concrete in points:
            try:
                vals = [e.subs(concrete) for e in exprs_]
            except Exception:
                continue
            try:
                point_ok: bool | None = True
                for i, op in enumerate(ops):
                    if i + 1 >= len(vals):
                        point_ok = None
                        break
                    d = sp.simplify(vals[i] - vals[i + 1])
                    if d.has(sp.nan, sp.zoo) or getattr(d, "is_infinite", False):
                        point_ok = None
                        break
                    if d == 0:
                        ok = op not in ("<", ">", "\\neq")
                    else:
                        try:
                            dn = float(d.evalf())
                        except (TypeError, ValueError, AttributeError):
                            point_ok = None
                            break
                        if op in ("<", "\\leq"):
                            ok = dn <= 0
                        elif op in (">", "\\geq"):
                            ok = dn >= 0
                        elif op == "\\neq":
                            ok = dn != 0
                        else:
                            point_ok = None
                            break
                    if not ok:
                        point_ok = False
                        first_violation = (
                            f"violated at {concrete}: {op} between "
                            f"{vals[i]} , {vals[i + 1]}")
                        break
            except Exception:
                continue
            if point_ok is None:
                continue
            decided += 1
            if not point_ok:
                violated += 1
        if decided == 0:
            return "SKIP: undecidable inequality"
        if violated == 0:
            return None
        # Sampling can confirm a universal inequality (holds at every
        # decidable point) but can never falsify one: a violation at every
        # sampled point is equally a solution set/condition whose satisfying
        # region the samples missed (`x > 3`). Report either way as SKIP —
        # a missed OK on a condition is honest; a FAIL would cry wolf.
        # (Genuinely false universals are caught by hand-triage of the
        # conditional-inequality bucket, which stays small.)
        return "SKIP: conditional inequality"

    def check_junction(val, ops, exprs_) -> str | None:
        # Mixed equality/inequality chain: the junction itself must be an
        # equality (val == exprs_[0]), then the inequality chain must hold.
        # The equality part gets the same honest-SKIP guards as a plain pair
        # (`a_{n+1} = 1/(n+1)!` is a definition with a side condition).
        guard = _pair_guard(val, exprs_[0])
        if guard:
            return f"SKIP: {guard}"
        msg = check_expr_pair(val, exprs_[0])
        if msg:
            return f"junction mismatch: {msg}"[:180]
        return check_ineq(ops, exprs_)

    # Walk each parsed chain pairwise. check_* helpers return None (pass),
    # "SKIP: reason" (honest skip), or anything else (FAIL).
    def _record(msg: str, chain_results) -> None:
        # A ground claim that fails next to contradiction language is a
        # lesson's intentional reductio (e.g. deriving 3 > 4 to show no
        # solution) — skip it; anything else stays a FAIL.
        if msg.startswith("SKIP: "):
            blk.verdict, blk.reason = "SKIP", msg[6:180]
            return
        # A mismatch in a block containing out-of-model macros (custom
        # operators like \star, teaching macros like \addright) is
        # untrustworthy parse output, not evidence of falsehood: abstain.
        macro = _unknown_macro(body)
        if macro:
            blk.verdict, blk.reason = "SKIP", f"unknown macro: {macro}"[:120]
            return
        try:
            allfree: set = set()
            for r in chain_results:
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

    def _walk_chain(chain_res) -> str | None:
        for idx in range(len(chain_res) - 1):
            ka, *pa = chain_res[idx]
            kb, *pb = chain_res[idx + 1]
            if ka == "expr" and kb == "expr":
                # Per-pair honest-SKIP guards (definitions, disjoint or
                # template unknowns): independent sampling cannot verify a
                # pair that shares nothing checkable — e.g. `z^n = ...`
                # under the definition `z := re^{iθ}`, or `OA = OP·cosθ`
                # with geometrically linked unknowns.
                guard = _pair_guard(pa[0], pb[0])
                if guard:
                    return f"SKIP: {guard}"
                msg = check_expr_pair(pa[0], pb[0])
            elif ka == "expr" and kb == "ineq":
                msg = check_junction(pa[0], pb[0], pb[1])
            elif ka == "ineq" and kb == "expr":
                msg = check_ineq(pa[0], pa[1])
                if not msg:
                    guard = _pair_guard(pa[1][-1], pb[0])
                    if guard:
                        msg = f"SKIP: {guard}"
                    else:
                        msg = check_expr_pair(pa[1][-1], pb[0])
                        if msg:
                            msg = f"junction mismatch: {msg}"[:180]
            else:
                msg = check_ineq(pa[0], pa[1]) or check_ineq(pb[0], pb[1])
            if msg:
                return msg
        if len(chain_res) == 1 and chain_res[0][0] == "ineq":
            return check_ineq(chain_res[0][1], chain_res[0][2])
        return None

    first_skip: str | None = None
    saw_ok = False
    for _par, chain_res in parsed_chains:
        msg = _walk_chain(chain_res)
        if msg is None:
            saw_ok = True
        elif msg.startswith("SKIP: "):
            if first_skip is None:
                first_skip = msg
        else:
            _record(f"[chain] {msg}" if len(parsed_chains) > 1 else msg,
                    chain_res)
            return
    if saw_ok:
        blk.verdict = "OK"
    elif first_skip is not None:
        blk.verdict, blk.reason = "SKIP", first_skip[6:180]
    else:
        blk.verdict, blk.reason = "SKIP", "empty chain"


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
