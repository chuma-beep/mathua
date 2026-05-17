import subprocess
import json
import sys
from pathlib import Path

SERVICE_PATH = Path(__file__).resolve().parent.parent / "sympy_service.py"


def grade_via_service(expected: str, answer: str, timeout: float = 10.0) -> dict:
    proc = subprocess.Popen(
        [sys.executable, str(SERVICE_PATH)],
        stdin=subprocess.PIPE,
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE,
        text=True,
    )
    req = json.dumps({"id": "1", "expected": expected, "answer": answer})
    out, err = proc.communicate(input=req + "\n", timeout=timeout)
    assert not err, f"stderr: {err}"
    return json.loads(out.strip())


def test_equivalent_polynomials():
    r = grade_via_service("x^2+2x+1", "(x+1)^2")
    assert r["correct"] is True, f"Expected True, got {r}"


def test_equivalent_expanded():
    r = grade_via_service("(x+3)(x+4)", "x^2+7x+12")
    assert r["correct"] is True, f"Expected True, got {r}"


def test_equivalent_with_implicit_mult():
    r = grade_via_service("6x^2", "6*x^2")
    assert r["correct"] is True, f"Expected True, got {r}"


def test_eulers_number():
    r = grade_via_service("e^x", "exp(x)")
    assert r["correct"] is True, f"Expected True, got {r}"


def test_eulers_number_chain():
    r = grade_via_service("2e^(2x)", "2*exp(2*x)")
    assert r["correct"] is True, f"Expected True, got {r}"


def test_trig_identity():
    r = grade_via_service("sin(x)^2+cos(x)^2", "1")
    assert r["correct"] is True, f"Expected True, got {r}"


def test_chain_rule_derivative():
    r = grade_via_service("3(2x+1)^2", "12x^2+12x+3")
    assert r["correct"] is True, f"Expected True, got {r}"


def test_absolute_value_with_log():
    r = grade_via_service(
        "(1/2)ln|x-1|-(1/2)ln|x+1|",
        "(1/2)*log(Abs(x-1))-(1/2)*log(Abs(x+1))",
    )
    assert r["correct"] is True, f"Expected True, got {r}"


def test_literal_equation():
    r = grade_via_service("y = (z - 2x)/3", "(z-2*x)/3")
    assert r["correct"] is True, f"Expected True, got {r}"


def test_constant_of_integration():
    r = grade_via_service("xe^x-e^x+C", "e^x*(x-1)+C")
    assert r["correct"] is True, f"Expected True, got {r}"


def test_plus_c_mismatch():
    r = grade_via_service("xe^x-e^x+C", "xe^x-e^x")
    assert r["correct"] is True, f"Expected True, got {r}"


def test_not_equivalent():
    r = grade_via_service("2x+3", "2x+4")
    assert r["correct"] is False, f"Expected False, got {r}"


def test_parse_error():
    r = grade_via_service("x^2", "invalid(((")
    assert r["correct"] is False
    assert "Parse error" in r["feedback"]


def test_multiple_requests():
    proc = subprocess.Popen(
        [sys.executable, str(SERVICE_PATH)],
        stdin=subprocess.PIPE,
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE,
        text=True,
    )
    requests = [
        {"id": "1", "expected": "x^2+2x+1", "answer": "(x+1)^2"},
        {"id": "2", "expected": "2x+3", "answer": "2x+4"},
        {"id": "3", "expected": "cos(x)", "answer": "cos(x)"},
    ]
    for req in requests:
        proc.stdin.write(json.dumps(req) + "\n")
    proc.stdin.flush()
    proc.stdin.close()
    out, _ = proc.communicate(timeout=10)
    results = [json.loads(line) for line in out.strip().splitlines()]
    assert len(results) == 3
    assert results[0]["correct"] is True
    assert results[1]["correct"] is False
    assert results[2]["correct"] is True


def test_equivalent_different_forms():
    r = grade_via_service("x(x+1)", "x^2+x")
    assert r["correct"] is True, f"Expected True, got {r}"


def test_trig_double_angle():
    r = grade_via_service("sin(2x)", "2sin(x)cos(x)")
    assert r["correct"] is True, f"Expected True, got {r}"


def test_trig_pythagorean_variant():
    r = grade_via_service("1-cos(x)^2", "sin(x)^2")
    assert r["correct"] is True, f"Expected True, got {r}"


def test_constant_integration_c1():
    r = grade_via_service("xe^x-e^x+C_1", "e^x*(x-1)+C_2")
    assert r["correct"] is True, f"Expected True, got {r}"


def test_constant_integration_minus_c():
    r = grade_via_service("x^2/2-C", "x^2/2+C")
    assert r["correct"] is True, f"Expected True, got {r}"


def test_log_absolute_equivalence():
    r = grade_via_service("ln|x|", "log(Abs(x))")
    assert r["correct"] is True, f"Expected True, got {r}"


def test_pi_symbol():
    r = grade_via_service("2pi", "2*pi")
    assert r["correct"] is True, f"Expected True, got {r}"


def test_arcsin_variants():
    r = grade_via_service("arcsin(x)", "asin(x)")
    assert r["correct"] is True, f"Expected True, got {r}"


def test_rational_expression_simplification():
    r = grade_via_service("(x^2-1)/(x-1)", "x+1")
    assert r["correct"] is True, f"Expected True, got {r}"


def test_power_simplification():
    r = grade_via_service("x^(1/2)", "sqrt(x)")
    assert r["correct"] is True, f"Expected True, got {r}"


def test_factored_vs_expanded():
    r = grade_via_service("x^2-4", "(x+2)(x-2)")
    assert r["correct"] is True, f"Expected True, got {r}"


def test_exp_equivalence():
    r = grade_via_service("exp(2x)", "e^(2x)")
    assert r["correct"] is True, f"Expected True, got {r}"


def test_negative_constant_integration():
    r = grade_via_service("sin(x)+C", "sin(x)-C")
    assert r["correct"] is True, f"Expected True, got {r}"
