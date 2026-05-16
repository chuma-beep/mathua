#!/usr/bin/env python3

import sys
import json
import re

from sympy.parsing.sympy_parser import (
    parse_expr,
    standard_transformations,
    implicit_multiplication_application,
    convert_xor,
)
from sympy import simplify

TRANSFORMATIONS = standard_transformations + (
    implicit_multiplication_application,
    convert_xor,
)


def preprocess(s: str) -> str:
    s = re.sub(r"^[a-zA-Z_]\w*\s*=\s*", "", s)
    s = re.sub(r"\s*\+\s*[A-Za-z]\s*$", "", s.strip())
    s = s.replace("e^", "E^").replace("e**", "E**")
    s = re.sub(r"ln\|([^|]+)\|", r"log(Abs(\1))", s)
    s = re.sub(r"\|([^|]+)\|", r"Abs(\1)", s)
    s = s.replace("ln", "log")
    return s


def grade(expected: str, answer: str) -> tuple[bool, str]:
    e = preprocess(expected)
    a = preprocess(answer)
    try:
        ee = parse_expr(e, transformations=TRANSFORMATIONS)
        ea = parse_expr(a, transformations=TRANSFORMATIONS)
        diff = simplify(ee - ea)
        if diff == 0:
            return True, ""
        return False, "Not equivalent"
    except Exception as ex:
        return False, f"Parse error: {ex}"


def main():
    for line in sys.stdin:
        line = line.strip()
        if not line:
            continue
        try:
            req = json.loads(line)
            correct, feedback = grade(req["expected"], req["answer"])
            resp = {"id": req["id"], "correct": correct, "feedback": feedback}
        except Exception as ex:
            resp = {"id": req.get("id", ""), "correct": False, "feedback": f"Service error: {ex}"}
        sys.stdout.write(json.dumps(resp) + "\n")
        sys.stdout.flush()


if __name__ == "__main__":
    main()
