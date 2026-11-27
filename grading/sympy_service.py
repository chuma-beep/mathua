#!/usr/bin/env python3

import sys
import json
import re
import random
import signal

from sympy.parsing.sympy_parser import (
    parse_expr,
    standard_transformations,
    implicit_multiplication_application,
    convert_xor,
)
from sympy import Expr, S, simplify, expand, trigsimp, radsimp, powsimp, factor, pi, E, I, Symbol

MAX_INPUT_LENGTH = 500
ALLOWED_CHARS = set("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789+-*/^()[]=., _|π")

TRANSFORMATIONS = standard_transformations + (
    implicit_multiplication_application,
    convert_xor,
)

COMMON_SYMBOLS = {
    'x': Symbol('x'),
    'y': Symbol('y'),
    'z': Symbol('z'),
    't': Symbol('t'),
    'n': Symbol('n'),
    'a': Symbol('a'),
    'b': Symbol('b'),
    'c': Symbol('c'),
    'pi': pi,
    'e': E,
    'E': E,
    'I': I,
}

_FUNC_ALIASES = {
    'arcsinh': 'asinh',
    'arccosh': 'acosh',
    'arctanh': 'atanh',
    'arcsech': 'asech',
    'arccsch': 'acsch',
    'arccoth': 'acoth',
    'arcsin': 'asin',
    'arccos': 'acos',
    'arctan': 'atan',
    'arcsec': 'asec',
    'arccsc': 'acsc',
    'arccot': 'acot',
    'ln': 'log',
}

# Longest names first so inverse and hyperbolic functions (arcsin, asec, sinh)
# are matched atomically and never split into a*sin(...) or a*sec(...).
_POWER_FUNCS = tuple(sorted((
    'arcsinh', 'arccosh', 'arctanh', 'arcsech', 'arccsch', 'arccoth',
    'arcsin', 'arccos', 'arctan', 'arcsec', 'arccsc', 'arccot',
    'asinh', 'acosh', 'atanh', 'asech', 'acsch', 'acoth',
    'sinh', 'cosh', 'tanh', 'sech', 'csch', 'coth',
    'asin', 'acos', 'atan', 'asec', 'acsc', 'acot',
    'sin', 'cos', 'tan', 'sec', 'csc', 'cot',
    'log', 'ln',
), key=len, reverse=True))

_CALL_FUNCS = tuple(sorted(_POWER_FUNCS + ('sqrt', 'exp', 'Abs'), key=len, reverse=True))

_FUNC_CALL_RE = re.compile(r"(" + "|".join(_CALL_FUNCS) + r")\s*\(")

_FUNC_POWER_RE = re.compile(
    r"(" + "|".join(_POWER_FUNCS) + r")\s*(?:\^|\*\*)\s*(?:\((\d+)\)|(\d+))\s*\("
)

_INTEGRAL_CONST_RE = re.compile(r"(?<=\S)\s*[+-]\s*[Cc](?:_\d+)?\s*$")
_EXPECTED_INTEGRAL_RE = re.compile(r"\S\s*[+-]\s*C(?:_\d+)?\s*$")

# Group unparenthesized coefficient-variable exponents on Euler's number
# (e.g. e^2x -> e^(2x), e^-3x -> e^(-3x)) while leaving chained powers
# (e^-x^2), function exponents (e^sin(x)) and spaces intact for SymPy.
_EULER_MONOMIAL_EXP_RE = re.compile(
    r"e(\^|\*\*)([+-]?\d+(?:\.\d+)?[a-zA-Z]+)(?![a-zA-Z\^]|\*\*)"
)


def _find_matching_paren(s: str, open_idx: int) -> int:
    """Return index of matching ')' for '(' at open_idx, or -1 if unbalanced."""
    depth = 1
    i = open_idx + 1
    while i < len(s):
        if s[i] == '(':
            depth += 1
        elif s[i] == ')':
            depth -= 1
            if depth == 0:
                return i
        i += 1
    return -1


def _rewrite_func_powers(s: str) -> str:
    """Rewrite sin^2(x) and sin^(2)(x) into ((sin(x))**2) for unsigned integer powers."""
    result = []
    pos = 0
    while True:
        m = _FUNC_POWER_RE.search(s, pos)
        if not m:
            result.append(s[pos:])
            break
        open_idx = m.end() - 1
        close_idx = _find_matching_paren(s, open_idx)
        if close_idx == -1:
            result.append(s[pos:m.end()])
            pos = m.end()
            continue
        func_name = m.group(1)
        exponent = m.group(2) or m.group(3)
        inner = _rewrite_func_powers(s[open_idx + 1:close_idx])
        prefix = '*' if m.start() > 0 and (s[m.start() - 1].isalnum() or s[m.start() - 1] == ')') else ''
        result.append(s[pos:m.start()])
        result.append(f"{prefix}(({func_name}({inner}))**{exponent})")
        pos = close_idx + 1
    return ''.join(result)


def _normalize_func_calls(s: str) -> str:
    """Normalize function aliases and insert '*' for implicit multiplication before calls."""
    def repl(m: re.Match) -> str:
        name = _FUNC_ALIASES.get(m.group(1), m.group(1))
        start = m.start()
        prefix = '*' if start > 0 and (s[start - 1].isalnum() or s[start - 1] == ')') else ''
        return f"{prefix}{name}("

    return _FUNC_CALL_RE.sub(repl, s)


def preprocess(s: str, is_integral: bool = False) -> str:
    s = s.strip()

    # Strip variable assignments: "y = ..." -> "..."
    s = re.sub(r"^[a-zA-Z_]\w*\s*=\s*", "", s)

    # Strip trailing integration constant (+ C, - C, + C_1, + c) only when
    # expected is an indefinite integral, preventing false positives on
    # non-integral answers and preserving algebraic variables like "a + b + c".
    if is_integral:
        s = _INTEGRAL_CONST_RE.sub("", s)

    # Normalize Unicode pi before function and multiplication passes
    s = s.replace("π", "pi")

    # Rewrite trig/log power notation (sin^2(x) -> ((sin(x))**2))
    s = _rewrite_func_powers(s)

    # Handle ln|x| and |x| before function call normalization
    s = re.sub(r"ln\|([^|]+)\|", r"log(Abs(\1))", s)
    s = re.sub(r"ln\s*\(\s*\|([^|]+)\|\s*\)", r"log(Abs(\1))", s)
    s = re.sub(r"\|([^|]+)\|", r"Abs(\1)", s)

    # Normalize function aliases (arcsin -> asin, ln -> log) and implicit multiplication
    s = _normalize_func_calls(s)

    # Group unparenthesized coefficient-variable exponents on e (e.g. e^-2x -> e^(-2x))
    s = _EULER_MONOMIAL_EXP_RE.sub(r"e\1(\2)", s)

    # Help with implicit multiplication: 2sin(x) -> 2*sin(x), 2x -> 2*x
    s = re.sub(r"(\d)([a-zA-Z(])", r"\1*\2", s)

    return s


def _try_parse(s: str):
    try:
        expr = parse_expr(s, transformations=TRANSFORMATIONS, local_dict=COMMON_SYMBOLS)
    except Exception:
        return None
    if not isinstance(expr, Expr):
        return None
    return expr


def _has_undefined_atom(expr: Expr) -> bool:
    return expr.has(S.NaN, S.ComplexInfinity)


def _has_nonfinite_atom(expr: Expr) -> bool:
    return expr.has(S.NaN, S.ComplexInfinity, S.Infinity, S.NegativeInfinity)


def _symbolic_equal(ee: Expr, ea: Expr) -> tuple[bool, str]:
    if _has_undefined_atom(ee) or _has_undefined_atom(ea):
        return False, "Not equivalent"

    if _has_nonfinite_atom(ee) or _has_nonfinite_atom(ea):
        return (ee == ea, "" if ee == ea else "Not equivalent")

    diff = ee - ea

    if diff == 0:
        return True, ""

    if expand(diff) == 0:
        return True, ""

    if trigsimp(diff) == 0:
        return True, ""

    if radsimp(diff) == 0:
        return True, ""

    if powsimp(diff) == 0:
        return True, ""

    if factor(diff) == 0:
        return True, ""

    if simplify(diff) == 0:
        return True, ""

    try:
        if ee.equals(ea):
            return True, ""
    except (TypeError, ValueError, ZeroDivisionError, ArithmeticError, NotImplementedError, AttributeError):
        return False, "Not equivalent"

    return False, "Not equivalent"


def _numerical_check(ee: Expr, ea: Expr, num_points: int = 10) -> tuple[bool, str]:
    if _has_nonfinite_atom(ee) or _has_nonfinite_atom(ea):
        return False, "Not equivalent"

    free_symbols = sorted(ee.free_symbols | ea.free_symbols, key=lambda sym: sym.name)

    if not free_symbols:
        return False, "Not equivalent"

    diff_expr = ee - ea
    rng = random.Random(42)
    compared = 0
    min_compared = max(1, num_points // 2)

    for _ in range(num_points):
        subs = {}
        for sym in free_symbols:
            val = rng.uniform(-5, 5)
            if abs(val) < 0.1:
                val = 0.5
            subs[sym] = val

        try:
            ee_val = ee.evalf(30, subs=subs)
            ea_val = ea.evalf(30, subs=subs)
            if _has_nonfinite_atom(ee_val) or _has_nonfinite_atom(ea_val):
                continue
            if not (ee_val.is_number and ea_val.is_number):
                continue
            diff_val = diff_expr.evalf(30, subs=subs)
            if _has_nonfinite_atom(diff_val) or not diff_val.is_number:
                continue
            diff_mag = abs(diff_val)
            if _has_nonfinite_atom(diff_mag):
                continue
            if diff_mag > 1e-6:
                return False, "Not equivalent"
        except Exception:
            continue

        compared += 1

    if compared < min_compared:
        return False, "Not equivalent"

    return True, ""


def grade(expected: str, answer: str) -> tuple[bool, str]:
    if len(expected) > MAX_INPUT_LENGTH or len(answer) > MAX_INPUT_LENGTH:
        return False, "Input too long"
    if not set(expected).issubset(ALLOWED_CHARS) or not set(answer).issubset(ALLOWED_CHARS):
        return False, "Input contains disallowed characters"

    is_integral = bool(_EXPECTED_INTEGRAL_RE.search(expected))
    e = preprocess(expected, is_integral=is_integral)
    a = preprocess(answer, is_integral=is_integral)

    ee = _try_parse(e)
    ea = _try_parse(a)

    if ee is None or ea is None:
        return False, "Parse error: could not parse expression"

    correct, feedback = _symbolic_equal(ee, ea)
    if correct:
        return True, ""

    return _numerical_check(ee, ea)


def _timeout_handler(signum, frame):
    raise TimeoutError("Grading timed out")


def main():
    has_sigalrm = hasattr(signal, "SIGALRM")
    for line in sys.stdin:
        line = line.strip()
        if not line:
            continue
        req = None
        try:
            req = json.loads(line)
            if not isinstance(req, dict):
                raise ValueError("Request must be a JSON object")
            if has_sigalrm:
                signal.signal(signal.SIGALRM, _timeout_handler)
                signal.alarm(3)
            try:
                correct, feedback = grade(req["expected"], req["answer"])
            finally:
                if has_sigalrm:
                    signal.alarm(0)
            resp = {"id": req["id"], "correct": correct, "feedback": feedback}
        except TimeoutError:
            req_id = req.get("id", "") if isinstance(req, dict) else ""
            resp = {"id": req_id, "correct": False, "feedback": "Grading timed out"}
        except Exception as ex:
            req_id = req.get("id", "") if isinstance(req, dict) else ""
            resp = {"id": req_id, "correct": False, "feedback": f"Service error: {ex}"}
        sys.stdout.write(json.dumps(resp) + "\n")
        sys.stdout.flush()


if __name__ == "__main__":
    main()
