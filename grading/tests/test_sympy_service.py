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
    payload = "".join(json.dumps(req) + "\n" for req in requests)
    out, _ = proc.communicate(input=payload, timeout=10)
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


sys.path.insert(0, str(SERVICE_PATH.parent))
import sympy_service  # noqa: E402


def test_numerical_check_rejects_overflow_and_undefined_answers():
    """Overflowing, undefined or large-offset wrong answers must never pass."""
    for expected, answer in [
        ("x + 1", "10^11"),
        ("x + 1", "10^12"),
        ("2x + 3", "e^100"),
        ("x^2 + 1", "1/0"),
        ("x + 1", "0/0"),
        ("x + 1", "x/(x-x)"),
        ("x + 1", "log(0)"),
        ("x + 1", "oo"),
        ("x + 1", "zoo"),
        ("x + 1", "nan"),
        ("1/0", "log(0)"),
        ("0/0", "nan"),
        ("10^7*x + 1", "10^7*x + 2"),
        ("10^8 + x", "10^8 + 2*x"),
        ("10^8 + x", "10^8"),
        ("10^10*x + 1000*x^2", "10^10*x"),
        ("e^20 + x", "e^20 + 2*x"),
    ]:
        correct, _ = sympy_service.grade(expected, answer)
        assert correct is False, f"Expected False for ({expected!r}, {answer!r})"


def test_non_expr_inputs_return_parse_error_without_crashing():
    """List or tuple syntax must return a clean parse error, not a TypeError."""
    for answer in ["[x + 1]", "(x, 1)"]:
        correct, feedback = sympy_service.grade("x + 1", answer)
        assert correct is False
        assert "Parse error" in feedback


def test_integration_constant_not_stripped_for_non_integrals():
    """Appending + C or - C to a non-integral answer must be rejected."""
    for expected, answer in [
        ("x^2 + 1", "x^2 + 1 + C"),
        ("x^2 + 1", "x^2 + 1 - C"),
        ("x^2 + 1", "x^2 + 1 + C_1"),
    ]:
        correct, _ = sympy_service.grade(expected, answer)
        assert correct is False, f"Expected False for ({expected!r}, {answer!r})"

    correct, feedback = sympy_service.grade("-C", "-C")
    assert correct is True, f"Bare -C should not be stripped to empty: {feedback!r}"


def test_inverse_trig_not_split_into_product():
    """Inverse and hyperbolic trig functions must not split into products."""
    for expected, answer in [
        ("arcsin(x)", "a*sin(x)"),
        ("arccos(x)", "a*cos(x)"),
        ("arctan(x)", "a*tan(x)"),
        ("asec(x)", "a*sec(x)"),
        ("acsc(x)", "a*csc(x)"),
        ("acot(x)", "a*cot(x)"),
        ("arcsec(x)", "a*r*c*sec(x)"),
        ("arcsinh(x)", "a*r*c*sinh(x)"),
    ]:
        correct, _ = sympy_service.grade(expected, answer)
        assert correct is False, f"Expected False for ({expected!r}, {answer!r})"

    for expected, answer in [
        ("arctan(1)", "pi/4"),
        ("arcsin(1/2)", "pi/6"),
        ("acot(1)", "pi/4"),
        ("asec(2)", "pi/3"),
        ("acsc(2)", "pi/6"),
        ("sin(arcsin(x))", "x"),
        ("tan(arctan(x))", "x"),
        ("2*arcsin(x)", "2arcsin(x)"),
        ("x*arcsin(x)", "xarcsin(x)"),
        ("x*asin(x)", "xasin(x)"),
        ("sech(x)^2", "sech^2(x)"),
    ]:
        correct, feedback = sympy_service.grade(expected, answer)
        assert correct is True, f"Expected True for ({expected!r}, {answer!r}), got {feedback!r}"


def test_euler_exponents_preserve_precedence_and_functions():
    """Signed, chained, spaced and functional exponents on e must parse accurately."""
    for expected, answer in [
        ("e^(-x)", "e^-x"),
        ("exp(-x)", "e^-x"),
        ("e^(-2x)", "e^-2x"),
        ("2*e^(-3x)", "2e^-3x"),
        ("e^(-(x+1))", "e^-(x+1)"),
        ("e^(-x) + 1", "e^-x + 1"),
        ("e^(-x)", "e**-x"),
        ("x*e^(-x)", "xe^-x"),
        ("x*e^(2x)", "xe^2x"),
        ("x*e^(-2x)", "xe^-2x"),
        ("e^(2*x)*(x+1)", "e^2x(x+1)"),
        ("e^(-x^2)", "e^-x^2"),
        ("e^(sin(x))", "e^sin(x)"),
        ("x", "e^ln(x)"),
        ("e^x * sin(x)", "e^xsin(x)"),
        ("e^x", "e^ x"),
        ("e^(-x)", "e^- x"),
        ("e^(x+1)", "e^ (x+1)"),
        ("e^Abs(x)", "e^|x|"),
    ]:
        correct, feedback = sympy_service.grade(expected, answer)
        assert correct is True, f"Expected True for ({expected!r}, {answer!r}), got {feedback!r}"

    for expected, answer in [
        ("e^(-2x)", "e^-x^2"),
        ("e^(2x)", "e^x^2"),
        ("x^2*e^2", "xe^2x"),
        ("e^(2*x)*y^2", "e^2xy^2"),
        ("x", "e^sin(x)"),
        ("x", "e^cos(x)"),
        ("x", "e^tan(x)"),
    ]:
        correct, _ = sympy_service.grade(expected, answer)
        assert correct is False, f"Expected False for ({expected!r}, {answer!r})"


def test_trig_and_log_power_notation():
    """Standard sin^n(x) notation emitted by calculus generators and MathLive."""
    for expected, answer in [
        ("sin^3(x)/3+C", "sin(x)^3/3+C"),
        ("sin^3(x)/3+C", "(1/3)*sin(x)^3+C"),
        ("sin^3(x)/3+C", "sin^3(x)/3"),
        ("sin^2(x)+cos^2(x)", "1"),
        ("sin ^2(x)+cos ^2(x)", "1"),
        ("cos(2x+1)^2", "cos^2(2x+1)"),
        ("cos(2*(x+1))^2", "cos^2(2*(x+1))"),
        ("x*sin(x)^2", "xsin^2(x)"),
        ("ln(x)^2", "ln^2(x)"),
        ("sin(x)^6", "sin^2(x)^3"),
    ]:
        correct, feedback = sympy_service.grade(expected, answer)
        assert correct is True, f"Expected True for ({expected!r}, {answer!r}), got {feedback!r}"

    correct, _ = sympy_service.grade("sin(x)^8", "sin^2(x)^3")
    assert correct is False, "sin^2(x)^3 must equal sin(x)^6, not sin(x)^8"


def test_unicode_pi_and_implicit_function_multiplication():
    """Unicode pi and implicit multiplication before ln/sec/csc/cot."""
    for expected, answer in [
        ("2*pi", "2π"),
        ("pi/2", "π/2"),
        ("x*ln(x)", "xln(x)"),
        ("x*ln|x|", "xln|x|"),
        ("x*sec(x)", "xsec(x)"),
        ("x*csc(x)", "xcsc(x)"),
        ("x*cot(x)", "xcot(x)"),
        ("2sec(2x)tan(2x)", "2tan(2x)sec(2x)"),
    ]:
        correct, feedback = sympy_service.grade(expected, answer)
        assert correct is True, f"Expected True for ({expected!r}, {answer!r}), got {feedback!r}"


def test_numerical_check_is_deterministic_across_calls():
    """_numerical_check must not depend on prior global RNG state."""
    import random
    ee = sympy_service._try_parse("x^2 + 1")
    ea_eq = sympy_service._try_parse("(x + 1)^2 - 2*x")
    ea_neq = sympy_service._try_parse("x^2 + 2")
    random.seed(1)
    first_eq = sympy_service._numerical_check(ee, ea_eq)
    first_neq = sympy_service._numerical_check(ee, ea_neq)
    assert first_eq == (True, "")
    assert first_neq[0] is False
    for seed in (7, 999, 123456):
        random.seed(seed)
        assert sympy_service._numerical_check(ee, ea_eq) == first_eq
        assert sympy_service._numerical_check(ee, ea_neq) == first_neq


def test_service_subprocess_handles_malformed_and_valid_requests():
    """Subprocess mode must survive malformed JSON and grade subsequent requests."""
    proc = subprocess.Popen(
        [sys.executable, str(SERVICE_PATH)],
        stdin=subprocess.PIPE,
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE,
        text=True,
    )
    lines = [
        "{not json",
        "[]",
        json.dumps({"id": "1", "expected": "sin^3(x)/3+C", "answer": "(1/3)*sin(x)^3"}),
        json.dumps({"id": "2", "expected": "x + 1", "answer": "10^11"}),
        json.dumps({"id": "3", "expected": "x + 1", "answer": "[x + 1]"}),
    ]
    out, err = proc.communicate(input="\n".join(lines) + "\n", timeout=15)
    assert not err, f"stderr: {err}"
    results = [json.loads(line) for line in out.strip().splitlines()]
    assert len(results) == 5
    assert results[0]["id"] == "" and results[0]["correct"] is False
    assert results[1]["id"] == "" and results[1]["correct"] is False
    assert results[2]["id"] == "1" and results[2]["correct"] is True
    assert results[3]["id"] == "2" and results[3]["correct"] is False
    assert results[4]["id"] == "3" and results[4]["correct"] is False and "Parse error" in results[4]["feedback"]

