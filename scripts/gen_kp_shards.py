#!/usr/bin/env python3
"""Generate KP shard files for a domain from an inline spec.

Usage: python3 scripts/gen_kp_shards.py <domain>
Each KP: {label, section, subgoals[3]} — the section must resolve to a real
heading in the concept's lesson source (teaching/* wins per loader).
"""
import json
import os
import re
import sys

ROOT = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
KP_DIR = os.path.join(ROOT, "data", "lessons", "kp")
LESSONS = os.path.join(ROOT, "data", "lessons")

# concept -> 3 KPs, each {label, section, subgoals}
SPEC = {
    # ── Arithmetic · Addition ─────────────────────────────
    "arith.add.single": [
        {"label": "Use addition notation", "section": "Use Addition Notation",
         "subgoals": ["Identify the addends and the sum", "Read 3 + 4 as three plus four", "Translate word phrases into addition notation"]},
        {"label": "Model addition with blocks", "section": "Model Addition of Whole Numbers",
         "subgoals": ["Model the first addend with ones blocks", "Model the second addend", "Count the total"]},
        {"label": "Add whole numbers", "section": "Add Whole Numbers Without Models",
         "subgoals": ["Line up digits by place value", "Add ones, then tens, then hundreds", "Carry when a column sums to 10 or more"]},
    ],
    "arith.add.double": [
        {"label": "Use addition notation", "section": "Use Addition Notation",
         "subgoals": ["Recognize the addends and the sum", "Read multi-digit expressions", "Estimate before adding"]},
        {"label": "Model addition with blocks", "section": "Model Addition of Whole Numbers",
         "subgoals": ["Model 17 + 26 with tens rods and ones", "Exchange 10 ones for 1 ten", "Read the final model"]},
        {"label": "Add whole numbers", "section": "Add Whole Numbers Without Models",
         "subgoals": ["Write addends with aligned place values", "Add each column right to left", "Carry when needed"]},
    ],
    "arith.add.triple": [
        {"label": "Use addition notation", "section": "Use Addition Notation",
         "subgoals": ["Identify the addends and the sum", "Read three-digit expressions", "Relate vertical form to notation"]},
        {"label": "Model addition with blocks", "section": "Model Addition of Whole Numbers",
         "subgoals": ["Model three-digit addends with blocks", "Combine hundreds, tens, and ones", "Exchange across columns"]},
        {"label": "Add whole numbers", "section": "Add Whole Numbers Without Models",
         "subgoals": ["Align hundreds, tens, and ones", "Add right to left, carrying as needed", "Check with estimation"]},
    ],
    "arith.add.carry": [
        {"label": "Use addition notation", "section": "Use Addition Notation",
         "subgoals": ["Identify when a column sum exceeds 9", "Prepare to regroup 10 ones into 1 ten", "Recognize carrying notation"]},
        {"label": "Model addition with blocks", "section": "Model Addition of Whole Numbers",
         "subgoals": ["Model the addends with base-10 blocks", "Exchange 10 ones for 1 ten", "Count the result"]},
        {"label": "Add whole numbers", "section": "Add Whole Numbers Without Models",
         "subgoals": ["Write the small 1 above the next column", "Add each place value right to left", "Carry repeatedly when needed"]},
    ],
    "arith.add.word": [
        {"label": "Use addition notation", "section": "Use Addition Notation",
         "subgoals": ["Find the action words that mean add", "Translate a word phrase into math", "Choose the addends"]},
        {"label": "Translate word phrases", "section": "Translate Word Phrases to Math Notation",
         "subgoals": ["Spot keywords like sum and total", "Write the expression", "Add to simplify"]},
        {"label": "Add in applications", "section": "Add Whole Numbers in Applications",
         "subgoals": ["Read the problem and underline the question", "Write a word phrase for what to find", "Answer with the correct units"]},
    ],
    # ── Arithmetic · Subtraction ──────────────────────────
    "arith.sub.single": [
        {"label": "Use subtraction notation", "section": "Use Subtraction Notation",
         "subgoals": ["Identify the minuend and the subtrahend", "Read 7 - 2 as seven minus two", "Relate subtraction to the missing addend"]},
        {"label": "Model subtraction", "section": "Model Subtraction of Whole Numbers",
         "subgoals": ["Model the minuend", "Take away the subtrahend", "Count what remains"]},
        {"label": "Subtract whole numbers", "section": "Subtract Whole Numbers",
         "subgoals": ["Line up digits by place value", "Subtract right to left", "Borrow when needed"]},
    ],
    "arith.sub.double": [
        {"label": "Use subtraction notation", "section": "Use Subtraction Notation",
         "subgoals": ["Recognize the minuend, subtrahend, and difference", "Write subtraction vertically", "Estimate the difference"]},
        {"label": "Model subtraction", "section": "Model Subtraction of Whole Numbers",
         "subgoals": ["Model tens and ones", "Take away ones first", "Exchange 1 ten for 10 ones when needed"]},
        {"label": "Subtract whole numbers", "section": "Subtract Whole Numbers",
         "subgoals": ["Align place values", "Subtract ones, then tens", "Borrow across columns"]},
    ],
    "arith.sub.borrow": [
        {"label": "Use subtraction notation", "section": "Use Subtraction Notation",
         "subgoals": ["Spot when borrowing is required", "Read the vertical form", "Prepare to regroup a ten"]},
        {"label": "Model subtraction", "section": "Model Subtraction of Whole Numbers",
         "subgoals": ["Model the larger number", "Exchange 1 ten for 10 ones", "Subtract each column"]},
        {"label": "Subtract whole numbers", "section": "Subtract Whole Numbers",
         "subgoals": ["Write the regrouped digits", "Subtract right to left", "Verify by adding back"]},
    ],
    "arith.sub.word": [
        {"label": "Use subtraction notation", "section": "Use Subtraction Notation",
         "subgoals": ["Find the action words that mean subtract", "Translate the phrase into math", "Identify the larger number first"]},
        {"label": "Model subtraction", "section": "Model Subtraction of Whole Numbers",
         "subgoals": ["Model the starting amount", "Remove the amount taken away", "Count what is left"]},
        {"label": "Subtract whole numbers", "section": "Subtract Whole Numbers",
         "subgoals": ["Set up the vertical subtraction", "Subtract and simplify", "Check with addition"]},
    ],
    # ── Arithmetic · Multiplication ───────────────────────
    "arith.mult.concept": [
        {"label": "Use multiplication notation", "section": "Use Multiplication Notation",
         "subgoals": ["Recognize the factors and the product", "Read 2 × 6 as two times six", "Rewrite repeated addition as multiplication"]},
        {"label": "Model multiplication", "section": "Model Multiplication of Whole Numbers",
         "subgoals": ["Model rows and columns", "Count the total objects", "Relate the array to the product"]},
        {"label": "Multiply whole numbers", "section": "Multiply Whole Numbers",
         "subgoals": ["Use known times facts", "Line up by place value", "Multiply and add partial products"]},
    ],
    "arith.mult.tables": [
        {"label": "Use multiplication notation", "section": "Use Multiplication Notation",
         "subgoals": ["Recognize the factors in a fact", "Say the fact aloud", "Relate 2 × 3 to 3 × 2"]},
        {"label": "Model multiplication", "section": "Model Multiplication of Whole Numbers",
         "subgoals": ["Build the array for the fact", "Count in groups", "State the product"]},
        {"label": "Multiply whole numbers", "section": "Multiply Whole Numbers",
         "subgoals": ["Recall the 1–12 table fact", "Practice for speed and accuracy", "Use the table to estimate"]},
    ],
    "arith.mult.2_5_10": [
        {"label": "Use multiplication notation", "section": "Use Multiplication Notation",
         "subgoals": ["Recognize the patterns in the 2s, 5s, and 10s", "Double for the 2s table", "Count by 5s and 10s"]},
        {"label": "Model multiplication", "section": "Model Multiplication of Whole Numbers",
         "subgoals": ["Group objects in 2s, 5s, and 10s", "Count the groups", "State the fact"]},
        {"label": "Multiply whole numbers", "section": "Multiply Whole Numbers",
         "subgoals": ["Apply the pattern to larger numbers", "Multiply by 10 by appending a zero", "Verify with repeated addition"]},
    ],
    "arith.mult.double": [
        {"label": "Use multiplication notation", "section": "Use Multiplication Notation",
         "subgoals": ["Identify the factors", "Rewrite as repeated addition", "Estimate the product"]},
        {"label": "Model multiplication", "section": "Model Multiplication of Whole Numbers",
         "subgoals": ["Break the multiplier into tens and ones", "Model with arrays or blocks", "Combine the parts"]},
        {"label": "Multiply whole numbers", "section": "Multiply Whole Numbers",
         "subgoals": ["Multiply the ones", "Multiply the tens", "Add the partial products"]},
    ],
    "arith.mult.triple": [
        {"label": "Use multiplication notation", "section": "Use Multiplication Notation",
         "subgoals": ["Identify the factors", "Estimate the product", "Prepare the vertical form"]},
        {"label": "Model multiplication", "section": "Model Multiplication of Whole Numbers",
         "subgoals": ["Model three-digit by one-digit", "Multiply by the ones first", "Regroup across columns"]},
        {"label": "Multiply whole numbers", "section": "Multiply Whole Numbers",
         "subgoals": ["Multiply ones, then tens, then hundreds", "Carry between columns", "Check with estimation"]},
    ],
    "arith.mult.word": [
        {"label": "Use multiplication notation", "section": "Use Multiplication Notation",
         "subgoals": ["Find the action words that mean multiply", "Translate the phrase into an expression", "Choose the factors"]},
        {"label": "Model multiplication", "section": "Model Multiplication of Whole Numbers",
         "subgoals": ["Draw the equal groups", "Count the total", "Write the multiplication"]},
        {"label": "Multiply whole numbers", "section": "Multiply Whole Numbers",
         "subgoals": ["Compute the product", "Answer with the correct units", "Check the reasonableness"]},
    ],
    # ── Arithmetic · Division ─────────────────────────────
    "arith.div.concept": [
        {"label": "Use division notation", "section": "Use Division Notation",
         "subgoals": ["Identify the dividend, divisor, and quotient", "Read 12 ÷ 3 as twelve divided by three", "Relate division to multiplication"]},
        {"label": "Model division", "section": "Model Division of Whole Numbers",
         "subgoals": ["Share the dividend into equal groups", "Count how many per group", "State the quotient"]},
        {"label": "Divide whole numbers", "section": "Divide Whole Numbers",
         "subgoals": ["Use related multiplication facts", "Divide step by step", "Check with multiplication"]},
    ],
    "arith.div.basic": [
        {"label": "Use division notation", "section": "Use Division Notation",
         "subgoals": ["Recognize the parts of a division", "Read the notation", "Relate to a fact family"]},
        {"label": "Model division", "section": "Model Division of Whole Numbers",
         "subgoals": ["Partition into equal groups", "Count groups or per-group", "Verify with multiplication"]},
        {"label": "Divide whole numbers", "section": "Divide Whole Numbers",
         "subgoals": ["Divide without remainder", "Use times tables", "Check the quotient"]},
    ],
    "arith.div.long": [
        {"label": "Use division notation", "section": "Use Division Notation",
         "subgoals": ["Set up the long-division bracket", "Divide the leading digit", "Bring down the next digit"]},
        {"label": "Model division", "section": "Model Division of Whole Numbers",
         "subgoals": ["Model with base-10 blocks", "Share the tens first", "Share the remainder as ones"]},
        {"label": "Divide whole numbers", "section": "Divide Whole Numbers",
         "subgoals": ["Divide, multiply, subtract, bring down", "Repeat for each digit", "Check by multiplying back"]},
    ],
    "arith.div.remainder": [
        {"label": "Use division notation", "section": "Use Division Notation",
         "subgoals": ["Recognize when division is not exact", "Read the remainder", "Relate to the fact family"]},
        {"label": "Model division", "section": "Model Division of Whole Numbers",
         "subgoals": ["Share the dividend into equal groups", "Count the leftover", "Write quotient with remainder"]},
        {"label": "Divide whole numbers", "section": "Divide Whole Numbers",
         "subgoals": ["Divide as usual", "Record the remainder", "Check: divisor × quotient + remainder"]},
    ],
    "arith.div.word": [
        {"label": "Use division notation", "section": "Use Division Notation",
         "subgoals": ["Find the action words that mean divide", "Translate the phrase into math", "Choose the dividend and divisor"]},
        {"label": "Model division", "section": "Model Division of Whole Numbers",
         "subgoals": ["Share the items into equal groups", "Count per group", "Handle leftovers"]},
        {"label": "Divide whole numbers", "section": "Divide Whole Numbers",
         "subgoals": ["Compute the quotient", "Answer with the correct units", "Check the reasonableness"]},
    ],
    # ── Arithmetic · Exponents ────────────────────────────
    "arith.exp.concept": [
        {"label": "Introduction to powers", "section": "Introduction to powers",
         "subgoals": ["Identify the base and the exponent", "Read a power aloud", "Expand as repeated multiplication"]},
        {"label": "Powers with real exponents", "section": "Powers with real exponents",
         "subgoals": ["Evaluate powers with integer exponents", "Apply signs to the base", "Interpret the value"]},
        {"label": "Fundamental rules of powers", "section": "Fundamental rules of powers",
         "subgoals": ["Multiply powers with the same base", "Divide powers with the same base", "Raise a power to a power"]},
    ],
    "arith.exp.evaluate": [
        {"label": "Introduction to powers", "section": "Introduction to powers",
         "subgoals": ["Identify the base and the exponent", "Expand the power", "Simplify the product"]},
        {"label": "Powers with real exponents", "section": "Powers with real exponents",
         "subgoals": ["Evaluate positive exponents", "Evaluate with a negative base", "Estimate by comparison"]},
        {"label": "Fundamental rules of powers", "section": "Fundamental rules of powers",
         "subgoals": ["Apply the product rule", "Apply the quotient rule", "Combine the rules"]},
    ],
    "arith.exp.product_rule": [
        {"label": "Introduction to powers", "section": "Introduction to powers",
         "subgoals": ["Recognize a product of powers", "Identify the common base", "Count the factors"]},
        {"label": "Fundamental rules of powers", "section": "Fundamental rules of powers",
         "subgoals": ["Add the exponents of the same base", "Keep the base unchanged", "Apply to a product of three powers"]},
        {"label": "Powers with real exponents", "section": "Powers with real exponents",
         "subgoals": ["Extend the rule to larger exponents", "Verify by expanding", "Apply the rule in expressions"]},
    ],
    "arith.exp.quotient_rule": [
        {"label": "Introduction to powers", "section": "Introduction to powers",
         "subgoals": ["Recognize a quotient of powers", "Identify the common base", "Cancel common factors"]},
        {"label": "Fundamental rules of powers", "section": "Fundamental rules of powers",
         "subgoals": ["Subtract the exponents of the same base", "Keep the base unchanged", "Handle a zero exponent"]},
        {"label": "Powers with real exponents", "section": "Powers with real exponents",
         "subgoals": ["Extend the rule", "Verify by expanding", "Apply the rule in expressions"]},
    ],
    "arith.exp.power_rule": [
        {"label": "Introduction to powers", "section": "Introduction to powers",
         "subgoals": ["Recognize a power raised to a power", "Identify the base and both exponents", "Count nested factors"]},
        {"label": "Fundamental rules of powers", "section": "Fundamental rules of powers",
         "subgoals": ["Multiply the exponents", "Keep the base unchanged", "Apply the rule in steps"]},
        {"label": "Powers with real exponents", "section": "Powers with real exponents",
         "subgoals": ["Extend the rule to three exponents", "Verify by expanding", "Apply the rule in expressions"]},
    ],
    # ── Arithmetic · Factoring ────────────────────────────
    "arith.factor.find": [
        {"label": "Identify multiples", "section": "Identify Multiples of Numbers",
         "subgoals": ["Recognize a multiple of a number", "List multiples by skip counting", "Identify multiples in a list"]},
        {"label": "Divisibility tests", "section": "Use Common Divisibility Tests",
         "subgoals": ["Apply the 2, 3, 5, 9, and 10 tests", "Check divisibility quickly", "Use the tests to narrow factors"]},
        {"label": "Find all factors", "section": "Find All the Factors of a Number",
         "subgoals": ["Start with 1 and the number itself", "Test each integer in order", "List the factor pairs"]},
    ],
    "arith.factor.composite": [
        {"label": "Identify multiples", "section": "Identify Multiples of Numbers",
         "subgoals": ["Define composite by its factors", "Recognize more than two factors", "Contrast with primes"]},
        {"label": "Divisibility tests", "section": "Use Common Divisibility Tests",
         "subgoals": ["Test divisibility to find factors", "Verify the factor count", "Use the tests to classify"]},
        {"label": "Find all factors", "section": "Find All the Factors of a Number",
         "subgoals": ["List all factors of the number", "Count the factors", "Decide composite or prime"]},
    ],
    "arith.factor.prime": [
        {"label": "Identify multiples", "section": "Identify Multiples of Numbers",
         "subgoals": ["Define prime by exactly two factors", "Check 1 and the number itself", "List the small primes"]},
        {"label": "Divisibility tests", "section": "Use Common Divisibility Tests",
         "subgoals": ["Test divisibility by small primes", "Rule out composite candidates", "Classify each number"]},
        {"label": "Find all factors", "section": "Find All the Factors of a Number",
         "subgoals": ["List the factor pairs", "Confirm exactly two factors", "Conclude prime"]},
    ],
    "arith.factor.prime_fact": [
        {"label": "Find all factors", "section": "Find All the Factors of a Number",
         "subgoals": ["Find the factors of the number", "Check which factors are prime", "Separate the prime factors"]},
        {"label": "Divisibility tests", "section": "Use Common Divisibility Tests",
         "subgoals": ["Divide by small primes", "Rewrite as a product of primes", "Continue until all factors are prime"]},
        {"label": "Identify multiples", "section": "Identify Multiples of Numbers",
         "subgoals": ["Express the number as a product", "Verify by multiplication", "Write the prime factorization"]},
    ],
    "arith.factor.gcf": [
        {"label": "Motivation for factoring", "section": "Motivation for Factoring",
         "subgoals": ["List the factors of each number", "Find the common factors", "Explain why factoring matters"]},
        {"label": "Identifying the GCF", "section": "Identifying the Greatest Common Factor",
         "subgoals": ["Compare the factor lists", "Pick the largest common factor", "Verify it divides both numbers"]},
        {"label": "Factoring out the GCF", "section": "Factoring Out the Greatest Common Factor",
         "subgoals": ["Identify the GCF", "Divide each term by the GCF", "Rewrite as a product"]},
    ],
    "arith.factor.lcm": [
        {"label": "Prime factorization", "section": "Find the Prime Factorization of a Composite Number",
         "subgoals": ["Break the number into any factor pair", "Factor each branch", "Collect the prime factors"]},
        {"label": "Factor-tree method", "section": "Prime Factorization Using the Factor Tree Method",
         "subgoals": ["Draw the factor tree", "Factor until every leaf is prime", "Write the product of primes"]},
        {"label": "Using prime factorization", "section": "Prime Factorization",
         "subgoals": ["Verify the prime product", "Compare with the ladder method", "Apply prime factorization to the LCM"]},
    ],
    # ── Arithmetic · Place value ──────────────────────────
    "arith.place.tens": [
        {"label": "Counting and whole numbers", "section": "Identify Counting Numbers and Whole Numbers",
         "subgoals": ["Distinguish counting numbers from whole numbers", "Recognize zero as a whole number", "Order numbers by size"]},
        {"label": "Model whole numbers", "section": "Model Whole Numbers",
         "subgoals": ["Model numbers with base-10 blocks", "Group 10 ones into 1 ten", "Read the model"]},
        {"label": "Place value of a digit", "section": "Identify the Place Value of a Digit",
         "subgoals": ["Label the ones and tens places", "Name the value of each digit", "Write in expanded form"]},
    ],
    "arith.place.hundreds": [
        {"label": "Counting and whole numbers", "section": "Identify Counting Numbers and Whole Numbers",
         "subgoals": ["Recognize whole numbers up to hundreds", "Count by hundreds", "Compare sizes"]},
        {"label": "Model whole numbers", "section": "Model Whole Numbers",
         "subgoals": ["Model with ones, tens, and hundreds blocks", "Regroup across columns", "Read the model"]},
        {"label": "Place value of a digit", "section": "Identify the Place Value of a Digit",
         "subgoals": ["Label the ones, tens, and hundreds places", "Name each digit's value", "Write in expanded form"]},
    ],
    "arith.place.thousands": [
        {"label": "Counting and whole numbers", "section": "Identify Counting Numbers and Whole Numbers",
         "subgoals": ["Recognize numbers into the thousands", "Count in thousands", "Compare magnitudes"]},
        {"label": "Model whole numbers", "section": "Model Whole Numbers",
         "subgoals": ["Model thousands with blocks", "Regroup 10 hundreds into 1 thousand", "Read the model"]},
        {"label": "Place value of a digit", "section": "Identify the Place Value of a Digit",
         "subgoals": ["Label up to the thousands place", "Name each digit's value", "Write in expanded form"]},
    ],
    # ── Arithmetic · Rounding ─────────────────────────────
    "arith.round.tens": [
        {"label": "Place value of a digit", "section": "Identify the Place Value of a Digit",
         "subgoals": ["Find the rounding digit (tens)", "Look at the digit to its right", "Decide to round up or down"]},
        {"label": "Model whole numbers", "section": "Model Whole Numbers",
         "subgoals": ["Locate the number on a number line", "Mark the nearest tens", "Pick the closer bound"]},
        {"label": "Round whole numbers", "section": "Round Whole Numbers",
         "subgoals": ["Apply the rounding rule", "Write the rounded number with zeros", "Check against the number line"]},
    ],
    "arith.round.hundreds": [
        {"label": "Place value of a digit", "section": "Identify the Place Value of a Digit",
         "subgoals": ["Find the rounding digit (hundreds)", "Look at the digit to its right", "Decide to round up or down"]},
        {"label": "Model whole numbers", "section": "Model Whole Numbers",
         "subgoals": ["Locate the number on a number line", "Mark the nearest hundreds", "Pick the closer bound"]},
        {"label": "Round whole numbers", "section": "Round Whole Numbers",
         "subgoals": ["Apply the rounding rule", "Write the rounded number with zeros", "Check against the number line"]},
    ],
    "arith.round.thousands": [
        {"label": "Place value of a digit", "section": "Identify the Place Value of a Digit",
         "subgoals": ["Find the rounding digit (thousands)", "Look at the digit to its right", "Decide to round up or down"]},
        {"label": "Model whole numbers", "section": "Model Whole Numbers",
         "subgoals": ["Locate the number on a number line", "Mark the nearest thousands", "Pick the closer bound"]},
        {"label": "Round whole numbers", "section": "Round Whole Numbers",
         "subgoals": ["Apply the rounding rule", "Write the rounded number with zeros", "Check against the number line"]},
    ],
    # ── Arithmetic · Square roots ─────────────────────────
    "arith.sqrt.perfect": [
        {"label": "Square root facts", "section": "Square Root Facts",
         "subgoals": ["Recognize the perfect squares", "Recall common square roots", "Relate squaring and taking roots"]},
        {"label": "Calculating square roots", "section": "Calculating Square Roots with a Calculator",
         "subgoals": ["Enter the radical", "Read the calculator result", "Identify exact versus rounded values"]},
        {"label": "Square roots of fractions", "section": "Square Roots of Fractions",
         "subgoals": ["Take the root of the numerator", "Take the root of the denominator", "Simplify the fraction"]},
    ],
    "arith.sqrt.simplify": [
        {"label": "Definition of radicals", "section": "Definition of radicals",
         "subgoals": ["Identify the radicand and the index", "Rewrite square roots as powers", "Relate roots to squaring"]},
        {"label": "Properties of radicals", "section": "Properties",
         "subgoals": ["Apply the product rule for radicals", "Extract perfect-square factors", "Simplify the radical"]},
        {"label": "Rationalizing the denominator", "section": "Rationalizing the denominator",
         "subgoals": ["Remove a radical from the denominator", "Multiply the numerator and denominator", "Simplify the result"]},
    ],
    # ── Arithmetic · Negatives ────────────────────────────
    "arith.neg.number_line": [
        {"label": "Signed numbers", "section": "Signed Numbers",
         "subgoals": ["Place negative numbers left of zero", "Read positions on the number line", "Compare negative values"]},
        {"label": "Adding signed numbers", "section": "Adding",
         "subgoals": ["Add by moving on the number line", "Move left for negatives", "Find the sum"]},
        {"label": "Subtracting signed numbers", "section": "Subtracting",
         "subgoals": ["Subtract by moving right", "Rewrite subtraction as addition of the opposite", "Find the difference"]},
    ],
    "arith.neg.add_sub": [
        {"label": "Signed numbers", "section": "Signed Numbers",
         "subgoals": ["Identify the signs of the addends", "Apply the sign rules", "Predict the sign of the result"]},
        {"label": "Adding signed numbers", "section": "Adding",
         "subgoals": ["Add numbers with the same sign", "Add numbers with opposite signs", "Combine absolute values"]},
        {"label": "Subtracting signed numbers", "section": "Subtracting",
         "subgoals": ["Keep the first number", "Change subtraction to addition of the opposite", "Apply the addition rules"]},
    ],
    "arith.neg.mult_div": [
        {"label": "Signed numbers", "section": "Signed Numbers",
         "subgoals": ["Identify the signs of the factors", "Apply the sign rule for products", "Predict the result sign"]},
        {"label": "Multiplying signed numbers", "section": "Multiplying",
         "subgoals": ["Multiply the absolute values", "Determine the sign", "State the product"]},
        {"label": "Powers of signed numbers", "section": "Powers",
         "subgoals": ["Multiply powers of negative bases", "Apply the sign rules to exponents", "Handle even and odd powers"]},
    ],
    # ── Arithmetic · Order of operations ──────────────────
    "arith.order_ops.basic": [
        {"label": "Grouping symbols", "section": "Grouping Symbols",
         "subgoals": ["Recognize parentheses and brackets", "Evaluate inside grouping first", "Handle nested groups"]},
        {"label": "Order of operations", "section": "Order of Operations",
         "subgoals": ["Multiply before adding", "Divide before subtracting", "Work left to right on ties"]},
        {"label": "Absolute value and grouping", "section": "Absolute Value and Implied Grouping",
         "subgoals": ["Treat absolute value as a grouping symbol", "Evaluate the inside first", "Simplify the result"]},
    ],
    "arith.order_ops.full": [
        {"label": "Grouping symbols", "section": "Grouping Symbols",
         "subgoals": ["Recognize all grouping symbols", "Evaluate the inner groups first"]},
        {"label": "Order of operations", "section": "Order of Operations",
         "subgoals": ["Apply BODMAS in full", "Handle exponents before multiplying", "Work left to right"]},
        {"label": "Absolute value and grouping", "section": "Absolute Value and Implied Grouping",
         "subgoals": ["Evaluate absolute values first", "Apply implied grouping", "Simplify fully"]},
    ],
    "arith.order_ops.nested": [
        {"label": "Grouping symbols", "section": "Grouping Symbols",
         "subgoals": ["Spot nested parentheses", "Start with the innermost group", "Work outward"]},
        {"label": "Order of operations", "section": "Order of Operations",
         "subgoals": ["Apply the full order inside each group", "Track operations carefully", "Simplify step by step"]},
        {"label": "Absolute value and grouping", "section": "Absolute Value and Implied Grouping",
         "subgoals": ["Recognize absolute value inside brackets", "Evaluate the innermost part first", "Simplify the expression"]},
    ],
    # ── Arithmetic · Decimals ─────────────────────────────
    "arith.dec.intro": [
        {"label": "Name decimals", "section": "Name Decimals",
         "subgoals": ["Read the whole-number part", "Read the decimal part as tenths or hundredths", "Use and for the decimal point"]},
        {"label": "Write decimals", "section": "Write Decimals",
         "subgoals": ["Write the whole-number part", "Place the decimal point", "Fill zeros for place value"]},
        {"label": "Convert decimals to fractions", "section": "Convert Decimals to Fractions or Mixed Numbers",
         "subgoals": ["Write the decimal as a fraction", "Simplify the fraction", "Convert to a mixed number"]},
    ],
}


def lesson_source(concept_id):
    with open(os.path.join(LESSONS, "lessons.json")) as f:
        entries = json.load(f)
    srcs = [m["source"] for m in entries if m["concept_id"] == concept_id]
    if not srcs:
        return None
    return sorted(srcs)[-1]  # teaching/* wins per loader


def headings(path):
    body = open(path, encoding="utf-8").read()
    heads = []
    for line in body.split("\n"):
        t = line.strip()
        m = re.match(r"^(#{2,4})\s+(.+)$", t)
        if m:
            heads.append(m.group(2).strip())
    return heads


def main():
    domain = sys.argv[1] if len(sys.argv) > 1 else "arith"
    os.makedirs(KP_DIR, exist_ok=True)
    problems = []
    written = 0
    for cid, kps in SPEC.items():
        if not cid.startswith(domain):
            continue
        src = lesson_source(cid)
        if src is None:
            problems.append(f"{cid}: no lesson source")
            continue
        path = os.path.join(LESSONS, src)
        if not os.path.exists(path):
            problems.append(f"{cid}: source missing {src}")
            continue
        heads = set(headings(path))
        for kp in kps:
            if kp["section"] not in heads:
                problems.append(f"{cid}: section {kp['section']!r} not in {src}")
        with open(os.path.join(KP_DIR, f"{cid}.json"), "w", encoding="utf-8") as f:
            json.dump(kps, f, indent=2, ensure_ascii=False)
            f.write("\n")
        written += 1
    print(f"wrote {written} kp files for domain '{domain}'")
    if problems:
        print("PROBLEMS:")
        for p in problems:
            print("  ", p)
        sys.exit(1)


if __name__ == "__main__":
    main()
