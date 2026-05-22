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
from sympy import simplify, expand, trigsimp, radsimp, powsimp, factor, pi, E, I, Symbol

MAX_INPUT_LENGTH = 500
ALLOWED_CHARS = set("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789+-*/^()[]=., _|pi eEInaxbyczt")

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

# Seed random for deterministic grading
random.seed(42)


def preprocess(s: str) -> str:
    s = s.strip()

    # Strip variable assignments: "y = ..." -> "..."
    s = re.sub(r"^[a-zA-Z_]\w*\s*=\s*", "", s)

    # Strip trailing "+ C", "- C", "+C", "-C" (constant of integration)
    s = re.sub(r"\s*[+-]\s*[Cc]\s*$", "", s)
    s = re.sub(r"\s*[+-]\s*C_\d+\s*$", "", s)

    # Normalize Euler's number: e^x -> exp(x), xe^x -> x*exp(x)
    # Only skip if 'e' is part of a known function name (exp, log, etc.)
    # If followed by ^ or **, it's definitely Euler's number
    def replace_e_power(s):
        result = []
        i = 0
        while i < len(s):
            # Check for e^ or e** pattern (Euler's number)
            is_e_power = False
            if s[i] == 'e':
                has_caret = i + 1 < len(s) and s[i+1] == '^'
                has_stars = i + 2 < len(s) and s[i+1:i+3] == '**'
                if has_caret or has_stars:
                    # If followed by ^ or **, it's Euler's number
                    # (function names like 'exp' are followed by '(', not '^')
                    is_e_power = True

            if is_e_power:
                result.append('exp(')
                # Skip 'e^' or 'e**'
                if s[i+1] == '^':
                    i += 2
                else:
                    i += 3

                # Handle the exponent
                if i < len(s) and s[i] == '(':
                    # Already has parens: e^(...) -> exp(...)
                    depth = 1
                    result.append(s[i])  # opening (
                    i += 1
                    while i < len(s) and depth > 0:
                        if s[i] == '(':
                            depth += 1
                        elif s[i] == ')':
                            depth -= 1
                        result.append(s[i])
                        i += 1
                else:
                    # No parens: e^x -> exp(x)
                    while i < len(s) and (s[i].isalnum() or s[i] in '._'):
                        result.append(s[i])
                        i += 1
                result.append(')')
            else:
                result.append(s[i])
                i += 1
        return ''.join(result)

    s = replace_e_power(s)

    # Normalize inverse trig names BEFORE implicit multiplication (to avoid matching 'n' in 'arcsin')
    s = re.sub(r"\barcsin\b", r"asin", s)
    s = re.sub(r"\barccos\b", r"acos", s)
    s = re.sub(r"\barctan\b", r"atan", s)

    # Handle implicit multiplication: xexp(x) -> x*exp(x), xsin(x) -> x*sin(x)
    s = re.sub(r"([a-zA-Z\)])(exp|log|sin|cos|tan|asin|acos|atan|sqrt|Abs)\(", r"\1*\2(", s)

    # Handle ln|x| -> log(Abs(x))
    s = re.sub(r"ln\|([^|]+)\|", r"log(Abs(\1))", s)
    s = re.sub(r"ln\s*\(\s*\|([^|]+)\|\s*\)", r"log(Abs(\1))", s)

    # Handle |x| -> Abs(x)
    s = re.sub(r"\|([^|]+)\|", r"Abs(\1)", s)

    # Normalize ln -> log
    s = s.replace("ln", "log")

    # Help with implicit multiplication: 2sin(x) -> 2*sin(x)
    s = re.sub(r"(\d)([a-zA-Z(])", r"\1*\2", s)

    # Normalize π -> pi
    s = s.replace("π", "pi")

    return s


def _try_parse(s: str):
    try:
        return parse_expr(s, transformations=TRANSFORMATIONS, local_dict=COMMON_SYMBOLS)
    except Exception:
        return None


def _symbolic_equal(ee, ea) -> tuple[bool, str]:
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
    except Exception as ex:
        pass

    return False, "Not equivalent"


def _numerical_check(ee, ea, num_points: int = 10) -> tuple[bool, str]:
    free_symbols = ee.free_symbols | ea.free_symbols

    if not free_symbols:
        return ee == ea, "Not equivalent"

    for _ in range(num_points):
        subs = {}
        for sym in free_symbols:
            val = random.uniform(-5, 5)
            if abs(val) < 0.1:
                val = 0.5
            subs[sym] = val

        try:
            ee_val = complex(ee.subs(subs))
            ea_val = complex(ea.subs(subs))

            if not (abs(ee_val) < 1e10 and abs(ea_val) < 1e10):
                continue

            if abs(ee_val - ea_val) > 1e-6:
                return False, "Not equivalent"
        except Exception as ex:
            continue

    return True, ""


def grade(expected: str, answer: str) -> tuple[bool, str]:
    if len(expected) > MAX_INPUT_LENGTH or len(answer) > MAX_INPUT_LENGTH:
        return False, "Input too long"
    if not set(expected).issubset(ALLOWED_CHARS) or not set(answer).issubset(ALLOWED_CHARS):
        return False, "Input contains disallowed characters"

    e = preprocess(expected)
    a = preprocess(answer)

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
    for line in sys.stdin:
        line = line.strip()
        if not line:
            continue
        try:
            req = json.loads(line)
            signal.signal(signal.SIGALRM, _timeout_handler)
            signal.alarm(10)
            try:
                correct, feedback = grade(req["expected"], req["answer"])
            finally:
                signal.alarm(0)
            resp = {"id": req["id"], "correct": correct, "feedback": feedback}
        except TimeoutError:
            resp = {"id": req.get("id", ""), "correct": False, "feedback": "Grading timed out"}
        except Exception as ex:
            resp = {"id": req.get("id", ""), "correct": False, "feedback": f"Service error: {ex}"}
        sys.stdout.write(json.dumps(resp) + "\n")
        sys.stdout.flush()


if __name__ == "__main__":
    main()
