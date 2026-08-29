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
    # ── Fractions ─────────────────────────────────────────
    "frac.basics.concept": [
        {"label": "Meaning of fractions", "section": "Understand the Meaning of Fractions",
         "subgoals": ["Identify the numerator and the denominator", "Interpret a fraction as parts of a whole", "Relate fractions to division"]},
        {"label": "Improper fractions and mixed numbers", "section": "Model Improper Fractions and Mixed Numbers",
         "subgoals": ["Recognize an improper fraction", "Model a mixed number", "Distinguish proper from improper"]},
        {"label": "Equivalent fractions", "section": "Model Equivalent Fractions",
         "subgoals": ["Model two fractions that name the same amount", "Multiply numerator and denominator together", "Recognize equivalence visually"]},
    ],
    "frac.basics.parts": [
        {"label": "Meaning of fractions", "section": "Understand the Meaning of Fractions",
         "subgoals": ["Identify the numerator and the denominator", "Describe what each part names", "Relate the denominator to equal parts"]},
        {"label": "Improper fractions and mixed numbers", "section": "Model Improper Fractions and Mixed Numbers",
         "subgoals": ["Recognize an improper fraction", "Model parts beyond a whole", "Name the parts correctly"]},
        {"label": "Property of one", "section": "Property of One",
         "subgoals": ["Recognize a fraction equal to one", "Explain why the parts equal the whole", "Use the property to simplify"]},
    ],
    "frac.basics.number_line": [
        {"label": "Locate fractions on the number line", "section": "Locate Fractions and Mixed Numbers on the Number Line",
         "subgoals": ["Divide the interval into equal parts", "Count the parts from zero", "Place the fraction correctly"]},
        {"label": "Order fractions", "section": "Order Fractions and Mixed Numbers",
         "subgoals": ["Compare positions on the number line", "Order a set of fractions", "Use the number line to verify"]},
        {"label": "Equivalent fractions", "section": "Model Equivalent Fractions",
         "subgoals": ["Model two names for the same point", "Recognize equal spacing", "Link equivalence to the number line"]},
    ],
    "frac.basics.equivalent": [
        {"label": "Model equivalent fractions", "section": "Model Equivalent Fractions",
         "subgoals": ["Model two fractions that name the same amount", "Split parts into smaller equal parts", "See the equivalence visually"]},
        {"label": "Find equivalent fractions", "section": "Find Equivalent Fractions",
         "subgoals": ["Multiply numerator and denominator by the same number", "Keep the value unchanged", "Generate a chain of equivalents"]},
        {"label": "Equivalent fractions property", "section": "Equivalent Fractions Property",
         "subgoals": ["State the property in words", "Apply it to build equivalents", "Use it to compare fractions"]},
    ],
    "frac.ops.simplify": [
        {"label": "Model equivalent fractions", "section": "Model Equivalent Fractions",
         "subgoals": ["Recognize a fraction in higher terms", "Find the common factor", "Reduce to lowest terms"]},
        {"label": "Find equivalent fractions", "section": "Find Equivalent Fractions",
         "subgoals": ["Divide numerator and denominator by a common factor", "Keep the value unchanged", "Repeat until lowest terms"]},
        {"label": "Equivalent fractions property", "section": "Equivalent Fractions Property",
         "subgoals": ["State the property in words", "Apply it in reverse to simplify", "Verify the simplified value"]},
    ],
    "frac.ops.compare": [
        {"label": "Model equivalent fractions", "section": "Model Equivalent Fractions",
         "subgoals": ["Model the two fractions", "Compare the shaded amounts", "Decide which is larger"]},
        {"label": "Find equivalent fractions", "section": "Find Equivalent Fractions",
         "subgoals": ["Rewrite with a common denominator", "Compare the numerators", "Use the comparison symbols"]},
        {"label": "Order fractions", "section": "Order Fractions and Mixed Numbers",
         "subgoals": ["Line the fractions up on a number line", "Order from least to greatest", "Verify with equivalent forms"]},
    ],
    "frac.ops.benchmark": [
        {"label": "Model equivalent fractions", "section": "Model Equivalent Fractions",
         "subgoals": ["Model a fraction close to 0, 1/2, or 1", "Compare shaded parts to half", "Classify the benchmark"]},
        {"label": "Locate fractions on the number line", "section": "Locate Fractions and Mixed Numbers on the Number Line",
         "subgoals": ["Place the fraction between 0 and 1", "Mark the halfway point", "Decide the closer benchmark"]},
        {"label": "Order fractions", "section": "Order Fractions and Mixed Numbers",
         "subgoals": ["Rank fractions by benchmark", "Order mixed numbers too", "Verify on the number line"]},
    ],
    "frac.ops.to_decimal": [
        {"label": "Convert fractions to decimals", "section": "Convert Fractions to Decimals",
         "subgoals": ["Divide the numerator by the denominator", "Write the quotient as a decimal", "Handle terminating results"]},
        {"label": "Convert a fraction to a decimal", "section": "Convert a Fraction to a Decimal",
         "subgoals": ["Set up the long division", "Place the decimal point", "State the decimal value"]},
        {"label": "Repeating decimals", "section": "Repeating Decimals",
         "subgoals": ["Recognize a repeating quotient", "Use the bar notation", "Round when needed"]},
    ],
    "frac.add.same": [
        {"label": "Model fraction addition", "section": "Model Fraction Addition",
         "subgoals": ["Model the first addend", "Model the second addend", "Count the shaded parts"]},
        {"label": "Add fractions with a common denominator", "section": "Add Fractions with a Common Denominator",
         "subgoals": ["Keep the denominator", "Add the numerators", "Simplify the result"]},
        {"label": "Fraction addition rule", "section": "Fraction Addition",
         "subgoals": ["State the rule for same denominators", "Apply it step by step", "Check with a model"]},
    ],
    "frac.sub.same": [
        {"label": "Model fraction subtraction", "section": "Model Fraction Subtraction",
         "subgoals": ["Model the first fraction", "Remove the second amount", "Count what remains"]},
        {"label": "Subtract fractions with a common denominator", "section": "Subtract Fractions with a Common Denominator",
         "subgoals": ["Keep the denominator", "Subtract the numerators", "Simplify the result"]},
        {"label": "Fraction subtraction rule", "section": "Fraction Subtraction",
         "subgoals": ["State the rule for same denominators", "Apply it step by step", "Check with a model"]},
    ],
    "frac.add.diff": [
        {"label": "Find the least common denominator", "section": "Find the Least Common Denominator",
         "subgoals": ["List the multiples of each denominator", "Pick the smallest common multiple", "Use it as the LCD"]},
        {"label": "Find the LCD of two fractions", "section": "Find the least common denominator (LCD) of two fractions.",
         "subgoals": ["Factor each denominator", "Combine the common factors", "Compute the LCD"]},
        {"label": "Convert to equivalent fractions with the LCD", "section": "Convert Fractions to Equivalent Fractions with the LCD",
         "subgoals": ["Scale each fraction to the LCD", "Add the numerators", "Simplify the result"]},
    ],
    "frac.sub.diff": [
        {"label": "Find the least common denominator", "section": "Find the Least Common Denominator",
         "subgoals": ["List the multiples of each denominator", "Pick the smallest common multiple", "Use it as the LCD"]},
        {"label": "Find the LCD of two fractions", "section": "Find the least common denominator (LCD) of two fractions.",
         "subgoals": ["Factor each denominator", "Combine the common factors", "Compute the LCD"]},
        {"label": "Convert to equivalent fractions with the LCD", "section": "Convert Fractions to Equivalent Fractions with the LCD",
         "subgoals": ["Scale each fraction to the LCD", "Subtract the numerators", "Simplify the result"]},
    ],
    "frac.add.word": [
        {"label": "Find the least common denominator", "section": "Find the Least Common Denominator",
         "subgoals": ["Spot the different denominators in the story", "Find their common multiple", "Prepare to add"]},
        {"label": "Convert to equivalent fractions with the LCD", "section": "Convert Fractions to Equivalent Fractions with the LCD",
         "subgoals": ["Rewrite each amount with the LCD", "Add the fractions", "Interpret the total in context"]},
        {"label": "Equivalent fractions property", "section": "Equivalent Fractions Property",
         "subgoals": ["Scale fractions without changing value", "Combine the story amounts", "Simplify the answer"]},
    ],
    "frac.ops.mult": [
        {"label": "Multiply fractions", "section": "Multiply Fractions",
         "subgoals": ["Multiply the numerators", "Multiply the denominators", "Simplify the product"]},
        {"label": "Fraction multiplication rule", "section": "Fraction Multiplication",
         "subgoals": ["State the rule for products of fractions", "Apply it step by step", "Cancel common factors first when possible"]},
        {"label": "Simplify fractions", "section": "Simplify Fractions",
         "subgoals": ["Divide by a common factor", "Reduce to lowest terms", "Verify the product"]},
    ],
    "frac.mult.whole": [
        {"label": "Multiply fractions", "section": "Multiply Fractions",
         "subgoals": ["Rewrite the whole number as a fraction", "Multiply the numerators", "Multiply the denominators"]},
        {"label": "Fraction multiplication rule", "section": "Fraction Multiplication",
         "subgoals": ["Apply the product rule", "Simplify the result", "Interpret the product"]},
        {"label": "Simplify a fraction", "section": "Simplify a fraction.",
         "subgoals": ["Divide by a common factor", "Reduce to lowest terms", "Verify with a model"]},
    ],
    "frac.ops.div": [
        {"label": "Find reciprocals", "section": "Find Reciprocals",
         "subgoals": ["Flip the numerator and denominator", "Identify the reciprocal", "Verify the product is 1"]},
        {"label": "Reciprocal definition", "section": "Reciprocal",
         "subgoals": ["State the definition", "Find reciprocals of fractions and wholes", "Use the reciprocal to divide"]},
        {"label": "Divide fractions", "section": "Divide Fractions",
         "subgoals": ["Keep the first fraction", "Multiply by the reciprocal", "Simplify the result"]},
    ],
    "frac.div.whole": [
        {"label": "Divide fractions", "section": "Divide Fractions",
         "subgoals": ["Rewrite the whole number as a fraction", "Multiply by its reciprocal", "Simplify the result"]},
        {"label": "Fraction division rule", "section": "Fraction Division",
         "subgoals": ["State the division rule", "Apply it with whole divisors", "Check the quotient"]},
        {"label": "Find reciprocals", "section": "Find Reciprocals",
         "subgoals": ["Find the reciprocal of the whole number", "Verify the product is 1", "Use it in the division"]},
    ],
    "frac.mixed.convert": [
        {"label": "Multiply and divide mixed numbers", "section": "Multiply and Divide Mixed Numbers",
         "subgoals": ["Convert a mixed number to an improper fraction", "Multiply the whole by the denominator and add", "Convert back when done"]},
        {"label": "Mixed-number operations", "section": "Multiply or divide mixed numbers.",
         "subgoals": ["Convert both mixed numbers first", "Perform the operation", "Simplify the result"]},
        {"label": "Translate phrases to expressions", "section": "Translate Phrases to Expressions with Fractions",
         "subgoals": ["Read the mixed-number phrase", "Write it as an expression", "Convert to improper form"]},
    ],
    "frac.mixed.add": [
        {"label": "Model addition of mixed numbers", "section": "Model Addition of Mixed Numbers with a Common Denominator",
         "subgoals": ["Model the whole parts", "Model the fraction parts", "Combine the total"]},
        {"label": "Add mixed numbers", "section": "Add Mixed Numbers",
         "subgoals": ["Add the whole parts", "Add the fraction parts", "Carry a whole when the fraction is improper"]},
        {"label": "Add mixed numbers with a common denominator", "section": "Add mixed numbers with a common denominator.",
         "subgoals": ["Add the fractions first", "Add the whole numbers", "Simplify and regroup"]},
    ],
    "frac.mixed.sub": [
        {"label": "Model subtraction of mixed numbers", "section": "Model Subtraction of Mixed Numbers",
         "subgoals": ["Model the first mixed number", "Remove the second amount", "Count what remains"]},
        {"label": "Subtract mixed numbers with a common denominator", "section": "Subtract Mixed Numbers with a Common Denominator",
         "subgoals": ["Subtract the fractions", "Borrow from the whole when needed", "Subtract the whole numbers"]},
        {"label": "Subtract as improper fractions", "section": "Subtract mixed numbers with common denominators as improper fractions.",
         "subgoals": ["Convert to improper fractions", "Subtract the numerators", "Convert the result back"]},
    ],
    "frac.mixed.mult": [
        {"label": "Multiply and divide mixed numbers", "section": "Multiply and Divide Mixed Numbers",
         "subgoals": ["Convert mixed numbers to improper fractions", "Multiply the fractions", "Convert the product back"]},
        {"label": "Mixed-number operations", "section": "Multiply or divide mixed numbers.",
         "subgoals": ["Convert both factors", "Apply the fraction rule", "Simplify the result"]},
        {"label": "Translate phrases to expressions", "section": "Translate Phrases to Expressions with Fractions",
         "subgoals": ["Read the mixed-number phrase", "Write the multiplication", "Solve and simplify"]},
    ],
    # ── Prealgebra · Exponents / negatives / sci notation ──
    "arith.exp.neg": [
        {"label": "Exponent basics", "section": "Exponent Basics",
         "subgoals": ["Identify the base and the exponent", "Expand the power", "Extend to negative exponents"]},
        {"label": "Exponent rules", "section": "Exponent Rules",
         "subgoals": ["Apply the negative-exponent rule", "Rewrite as a reciprocal", "Evaluate the result"]},
        {"label": "Product rule", "section": "Product Rule",
         "subgoals": ["Add exponents with the same base", "Handle negative exponents in products", "Simplify the expression"]},
    ],
    "arith.exp.zero": [
        {"label": "Exponent basics", "section": "Exponent Basics",
         "subgoals": ["Identify the base and the exponent", "Expand the power", "Consider the case of a zero exponent"]},
        {"label": "Exponent rules", "section": "Exponent Rules",
         "subgoals": ["Apply the zero-exponent rule", "Explain why any nonzero base gives 1", "Evaluate expressions with zero exponents"]},
        {"label": "Product to a power rule", "section": "Product to a Power Rule",
         "subgoals": ["Distribute the exponent over a product", "Keep the base rules consistent", "Simplify the expression"]},
    ],
    "arith.exp.sci_notation": [
        {"label": "Negative exponents", "section": "Use the Definition of a Negative Exponent",
         "subgoals": ["Rewrite a negative exponent as a reciprocal", "Convert between forms", "Prepare for scientific notation"]},
        {"label": "Convert to scientific notation", "section": "Convert from Decimal Notation to Scientific Notation",
         "subgoals": ["Move the decimal point after the first digit", "Count the places moved", "Choose a positive or negative power of 10"]},
        {"label": "Convert to decimal form", "section": "Convert Scientific Notation to Decimal Form",
         "subgoals": ["Move the decimal point by the exponent", "Fill zeros as needed", "State the decimal value"]},
    ],
    "arith.sci_notation.ops": [
        {"label": "Basics of scientific notation", "section": "The Basics of Scientific Notation",
         "subgoals": ["Write a number in scientific notation", "Identify the mantissa and the exponent", "Convert both directions"]},
        {"label": "Large and small numbers", "section": "Scientific Notation for Large Numbers",
         "subgoals": ["Use positive exponents for large numbers", "Use negative exponents for small numbers", "Write the results compactly"]},
        {"label": "Operations in scientific notation", "section": "Multiplying and Dividing Using Scientific Notation",
         "subgoals": ["Multiply the mantissas and add exponents", "Divide mantissas and subtract exponents", "Normalize the result"]},
    ],
    "arith.neg.abs_value": [
        {"label": "Definition", "section": "Definition",
         "subgoals": ["Define the absolute value as distance from zero", "Evaluate absolute values of negatives", "Recognize that the result is never negative"]},
        {"label": "Properties", "section": "Properties",
         "subgoals": ["Apply the product property", "Use the triangle-inequality intuition", "Simplify expressions with absolute value"]},
        {"label": "Triangle inequality", "section": "Triangle inequality",
         "subgoals": ["State the inequality in words", "Verify it on examples", "Use it to bound sums"]},
    ],
    "arith.neg.order_ops": [
        {"label": "Signed numbers", "section": "Signed Numbers",
         "subgoals": ["Identify the signs of each term", "Apply the sign rules", "Predict the sign of the result"]},
        {"label": "Adding signed numbers", "section": "Adding",
         "subgoals": ["Add the absolute values for like signs", "Subtract for opposite signs", "Write the result with the right sign"]},
        {"label": "Multiplying signed numbers", "section": "Multiplying",
         "subgoals": ["Multiply the absolute values", "Determine the sign from the count of negatives", "Apply the order of operations"]},
    ],
    # ── Prealgebra · Decimals ─────────────────────────────
    "dec.basics.compare": [
        {"label": "Locate decimals on the number line", "section": "Locate Decimals on the Number Line",
         "subgoals": ["Divide the interval into tenths and hundredths", "Place the decimal exactly", "Read its position"]},
        {"label": "Order decimals", "section": "Order Decimals",
         "subgoals": ["Compare digit by digit from left to right", "Align the decimal points", "Order the set"]},
        {"label": "Equivalent decimals", "section": "Equivalent Decimals",
         "subgoals": ["Recognize trailing zeros as equal", "Use equivalent forms to compare", "Write the comparison"]},
    ],
    "dec.convert.from_frac": [
        {"label": "Convert fractions to decimals", "section": "Convert Fractions to Decimals",
         "subgoals": ["Divide the numerator by the denominator", "Write the quotient as a decimal", "Handle terminating results"]},
        {"label": "Convert a fraction to a decimal", "section": "Convert a Fraction to a Decimal",
         "subgoals": ["Set up the long division", "Place the decimal point", "State the decimal value"]},
        {"label": "Repeating decimals", "section": "Repeating Decimals",
         "subgoals": ["Recognize a repeating quotient", "Use the bar notation", "Compare repeating decimals"]},
    ],
    "dec.convert.to_frac": [
        {"label": "Convert fractions to decimals", "section": "Convert Fractions to Decimals",
         "subgoals": ["Read the decimal places", "Write the decimal over a power of 10", "Simplify the fraction"]},
        {"label": "Repeating decimals", "section": "Repeating Decimals",
         "subgoals": ["Distinguish terminating from repeating", "Convert a terminating decimal exactly", "Round repeating values when needed"]},
        {"label": "Order decimals and fractions", "section": "Order Decimals and Fractions",
         "subgoals": ["Convert both to the same form", "Compare the values", "Order the set"]},
    ],
    "dec.ops.add": [
        {"label": "Add and subtract decimals", "section": "Add and Subtract Decimals",
         "subgoals": ["Align the decimal points", "Fill missing places with zeros", "Add the columns"]},
        {"label": "Add or subtract decimals", "section": "Add or subtract decimals.",
         "subgoals": ["Write the numbers vertically", "Add right to left", "Carry when needed"]},
        {"label": "Multiply decimals", "section": "Multiply Decimals",
         "subgoals": ["Multiply as whole numbers", "Count the decimal places", "Place the decimal point"]},
    ],
    "dec.ops.sub": [
        {"label": "Add and subtract decimals", "section": "Add and Subtract Decimals",
         "subgoals": ["Align the decimal points", "Fill missing places with zeros", "Subtract the columns"]},
        {"label": "Add or subtract decimals", "section": "Add or subtract decimals.",
         "subgoals": ["Write the numbers vertically", "Borrow across columns", "Keep the decimal point aligned"]},
        {"label": "Multiply decimals", "section": "Multiply Decimals",
         "subgoals": ["Multiply as whole numbers", "Count the decimal places", "Place the decimal point"]},
    ],
    "dec.ops.mult": [
        {"label": "Multiply decimals", "section": "Multiply Decimals",
         "subgoals": ["Ignore the decimal points and multiply", "Count the total decimal places", "Place the point in the product"]},
        {"label": "Multiply decimal numbers", "section": "Multiply decimal numbers.",
         "subgoals": ["Multiply digit by digit", "Count places in both factors", "Write the product correctly"]},
        {"label": "Multiply by powers of 10", "section": "Multiply by Powers of $10$",
         "subgoals": ["Move the decimal point right", "Shift once per power of 10", "Fill zeros when needed"]},
    ],
    "dec.ops.div": [
        {"label": "Divide decimals", "section": "Divide Decimals",
         "subgoals": ["Make the divisor a whole number", "Move the decimal point in both", "Divide as with whole numbers"]},
        {"label": "Divide by a whole number", "section": "Divide a decimal by a whole number.",
         "subgoals": ["Place the decimal point in the quotient", "Divide digit by digit", "Bring down digits"]},
        {"label": "Divide by another decimal", "section": "Divide a Decimal by Another Decimal",
         "subgoals": ["Shift both decimal points", "Divide normally", "Check by multiplying back"]},
    ],
    "dec.ops.round": [
        {"label": "Locate decimals on the number line", "section": "Locate Decimals on the Number Line",
         "subgoals": ["Place the decimal between two marks", "Identify the nearest place", "Decide up or down"]},
        {"label": "Order decimals", "section": "Order Decimals",
         "subgoals": ["Compare the digits in order", "Find the rounding place", "Look at the digit to the right"]},
        {"label": "Round decimals", "section": "Round Decimals",
         "subgoals": ["Apply the rounding rule", "Drop the trailing digits", "Write the rounded value"]},
    ],
    # ── Prealgebra · Percent ──────────────────────────────
    "pct.basics.concept": [
        {"label": "Definition of percent", "section": "Use the Definition of Percent",
         "subgoals": ["Read percent as parts per hundred", "Write a percent as a fraction over 100", "Interpret a percent visually"]},
        {"label": "Percent", "section": "Percent",
         "subgoals": ["State the meaning of the word percent", "Relate percents to decimals and fractions", "Compare percents"]},
        {"label": "Convert percents", "section": "Convert Percents to Fractions and Decimals",
         "subgoals": ["Write a percent as a fraction", "Write a percent as a decimal", "Convert in both directions"]},
    ],
    "pct.convert.from_dec": [
        {"label": "Convert decimals and fractions to percents", "section": "Convert Decimals and Fractions to Percents",
         "subgoals": ["Multiply a decimal by 100", "Attach the percent sign", "Convert a fraction through the decimal"]},
        {"label": "Convert a decimal to a percent", "section": "Convert a decimal to a percent.",
         "subgoals": ["Shift the decimal point two places", "Add the percent sign", "Verify against the fraction"]},
        {"label": "Convert a fraction to a percent", "section": "Convert a fraction to a percent.",
         "subgoals": ["Divide to get a decimal", "Convert the decimal to a percent", "Round when needed"]},
    ],
    "pct.convert.to_dec": [
        {"label": "Convert percents to fractions and decimals", "section": "Convert Percents to Fractions and Decimals",
         "subgoals": ["Drop the percent sign", "Divide by 100", "Simplify the fraction"]},
        {"label": "Convert a percent to a fraction", "section": "Convert a percent to a fraction.",
         "subgoals": ["Write over 100", "Reduce to lowest terms", "Handle mixed percents"]},
        {"label": "Convert a percent to a decimal", "section": "Convert a percent to a decimal.",
         "subgoals": ["Shift the decimal point two places left", "Drop the percent sign", "Verify against the fraction"]},
    ],
    "pct.ops.of_number": [
        {"label": "Basic percent equations", "section": "Translate and Solve Basic Percent Equations",
         "subgoals": ["Write percent as a decimal", "Set up amount = percent × base", "Solve for the amount"]},
        {"label": "Applications of percent", "section": "Solve Applications of Percent",
         "subgoals": ["Identify the percent and the base", "Translate the application", "Compute the part"]},
        {"label": "Solve an application", "section": "Solve an application",
         "subgoals": ["Read the problem carefully", "Write the percent equation", "Answer with units"]},
    ],
    "pct.ops.find_rate": [
        {"label": "Basic percent equations", "section": "Translate and Solve Basic Percent Equations",
         "subgoals": ["Identify the part and the base", "Write the ratio part over base", "Convert to a percent"]},
        {"label": "Applications of percent", "section": "Solve Applications of Percent",
         "subgoals": ["Ask what percent A is of B", "Divide A by B", "Convert the quotient to a percent"]},
        {"label": "Solve an application", "section": "Solve an application",
         "subgoals": ["Extract the numbers from the story", "Set up the rate", "State the percentage"]},
    ],
    "pct.ops.increase": [
        {"label": "Basic percent equations", "section": "Translate and Solve Basic Percent Equations",
         "subgoals": ["Find the amount of change", "Divide the change by the original", "Convert to a percent"]},
        {"label": "Applications of percent", "section": "Solve Applications of Percent",
         "subgoals": ["Distinguish increase from decrease", "Apply the change to the base", "Interpret the result"]},
        {"label": "Solve an application", "section": "Solve an application",
         "subgoals": ["Compute the new value", "Report the percent change", "Check the direction of change"]},
    ],
    "pct.ops.discount": [
        {"label": "Solve discount applications", "section": "Solve Discount Applications",
         "subgoals": ["Find the discount amount", "Subtract it from the original price", "State the sale price"]},
        {"label": "Discount", "section": "Discount",
         "subgoals": ["Write the discount as a percent equation", "Solve for the amount off", "Verify the sale price"]},
        {"label": "Solve mark-up applications", "section": "Solve Mark-up Applications",
         "subgoals": ["Find the mark-up amount", "Add it to the cost", "State the selling price"]},
    ],
    "pct.ops.tax_tip": [
        {"label": "Solve sales tax applications", "section": "Solve Sales Tax Applications",
         "subgoals": ["Find the tax amount", "Add it to the price", "State the total cost"]},
        {"label": "Sales tax", "section": "Sales Tax",
         "subgoals": ["Write the tax as a percent equation", "Compute the tax", "Round to cents"]},
        {"label": "Solve commission applications", "section": "Solve Commission Applications",
         "subgoals": ["Find the commission amount", "Add it to any base pay", "State the total"]},
    ],
    # ── Prealgebra · Equations ────────────────────────────
    "prealg.eq.one_step_add": [
        {"label": "Basic principle of algebra", "section": "The Basic Principle of Algebra",
         "subgoals": ["Do the same to both sides", "Keep the equation balanced", "Undo addition with subtraction"]},
        {"label": "Solve one-step equations", "section": "Solving One-Step Equations and Stating Solution Sets",
         "subgoals": ["Identify the operation on the variable", "Apply the inverse operation", "State the solution set"]},
        {"label": "Equations with fractions", "section": "Equations with Fractions",
         "subgoals": ["Undo a fraction with multiplication", "Clear the denominator", "Solve the resulting equation"]},
    ],
    "prealg.eq.one_step_mult": [
        {"label": "Imagine filling in the blanks", "section": "Imagine Filling in the Blanks",
         "subgoals": ["Guess a candidate value", "Check it in the equation", "Adjust until it works"]},
        {"label": "Basic principle of algebra", "section": "The Basic Principle of Algebra",
         "subgoals": ["Do the same to both sides", "Undo multiplication with division", "Keep the equation balanced"]},
        {"label": "Solve one-step equations", "section": "Solving One-Step Equations and Stating Solution Sets",
         "subgoals": ["Identify the coefficient", "Divide both sides by it", "State the solution set"]},
    ],
    "prealg.eq.two_step": [
        {"label": "Solve two-step equations", "section": "Solving Two-Step Equations",
         "subgoals": ["Undo addition or subtraction first", "Then undo multiplication or division", "Check the solution"]},
        {"label": "Solve multistep equations", "section": "Solving Multistep Linear Equations",
         "subgoals": ["Distribute when needed", "Combine like terms", "Isolate the variable"]},
        {"label": "Revisiting applications", "section": "Revisiting Applications",
         "subgoals": ["Translate the story into an equation", "Solve step by step", "Answer with units"]},
    ],
    "prealg.eq.word": [
        {"label": "Translate phrases", "section": "Translating Phrases into Algebraic Expressions and Equations/Inequalities",
         "subgoals": ["Map words to operations", "Define the variable", "Write the equation"]},
        {"label": "Rate models", "section": "Rate Models",
         "subgoals": ["Identify the rate in the story", "Write rate × time = amount", "Solve for the unknown"]},
        {"label": "Percent applications", "section": "Percent Applications",
         "subgoals": ["Convert the percent", "Write the percent equation", "Solve and interpret"]},
    ],
    # ── Prealgebra · Expressions ──────────────────────────
    "prealg.expr.distribute": [
        {"label": "Distributive property", "section": "Distributive Property",
         "subgoals": ["Multiply the term by each term inside", "Keep the signs", "Combine the products"]},
        {"label": "Applying the properties", "section": "Applying the Commutative, Associative, and Distributive Properties",
         "subgoals": ["Reorder terms with commutativity", "Regroup with associativity", "Distribute to simplify"]},
        {"label": "Role of order of operations", "section": "The Role of the Order of Operations",
         "subgoals": ["Simplify inside grouping first", "Multiply and divide before adding", "Compare with and without distribution"]},
    ],
    "prealg.expr.evaluate": [
        {"label": "Introduction to variables", "section": "Introduction to Variables",
         "subgoals": ["Recognize a variable as a placeholder", "Substitute the given value", "Prepare to evaluate"]},
        {"label": "Algebraic expressions", "section": "Algebraic Expressions",
         "subgoals": ["Replace each variable with its value", "Apply the order of operations", "Simplify to a number"]},
        {"label": "Evaluating with negative numbers", "section": "Evaluating Expressions with Negative Numbers",
         "subgoals": ["Substitute negative values carefully", "Track signs in multiplication", "Simplify the result"]},
    ],
    "prealg.expr.like_terms": [
        {"label": "Identifying terms", "section": "Identifying Terms",
         "subgoals": ["Separate an expression into terms", "Recognize the variable parts", "Classify like terms"]},
        {"label": "Combining like terms", "section": "Combining Like Terms",
         "subgoals": ["Add the coefficients of like terms", "Keep the variable part unchanged", "Simplify the expression"]},
        {"label": "Applications", "section": "Applications",
         "subgoals": ["Apply combining to word contexts", "Simplify the model", "Evaluate when needed"]},
    ],
    "prealg.ineq.one_step": [
        {"label": "Solve linear inequalities", "section": "Solving Linear Inequalities",
         "subgoals": ["Isolate the variable", "Keep the inequality direction", "Graph the solution set"]},
        {"label": "Negation", "section": "Negation",
         "subgoals": ["Negate both sides carefully", "Reverse the inequality sign", "Check the direction"]},
        {"label": "Check the solution", "section": "",
         "subgoals": ["Pick a value from the solution set", "Verify it satisfies the inequality", "Pick one outside to confirm"]},
    ],
    "prealg.ineq.two_step": [
        {"label": "Solve multistep inequalities", "section": "Solving Multistep Inequalities",
         "subgoals": ["Undo addition or subtraction first", "Then undo multiplication or division", "Reverse the sign when dividing by a negative"]},
        {"label": "Applications", "section": "Applications",
         "subgoals": ["Translate the constraint into an inequality", "Solve for the variable", "Interpret the range of answers"]},
        {"label": "Graph the solution set", "section": "",
         "subgoals": ["Mark the boundary on a number line", "Shade the valid side", "Use an open or closed circle"]},
    ],
    # ── Prealgebra · Real numbers / types ─────────────────
    "prealg.real.concept": [
        {"label": "Field and order structure", "section": "Field and order structure",
         "subgoals": ["Recognize the arithmetic operations on reals", "Apply the ordering of reals", "Combine structure with intuition"]},
        {"label": "The real line", "section": "The real line",
         "subgoals": ["Place every real number on a line", "Locate negatives and fractions", "Compare positions"]},
        {"label": "The completeness axiom", "section": "The completeness axiom",
         "subgoals": ["State the axiom in words", "See why fractions have limits", "Use it to place irrationals"]},
    ],
    "prealg.real.properties": [
        {"label": "Closure property", "section": "Closure property",
         "subgoals": ["Add two reals and stay in the reals", "Multiply two reals and stay in the reals", "Test closure on examples"]},
        {"label": "Commutative property", "section": "Commutative property",
         "subgoals": ["Swap the order of addition", "Swap the order of multiplication", "Verify the results match"]},
        {"label": "Distributive property", "section": "Distributive property",
         "subgoals": ["Expand a(b + c)", "Factor the expanded form", "Use it to simplify"]},
    ],
    "prealg.types": [
        {"label": "Natural numbers", "section": "Natural numbers",
         "subgoals": ["Identify the counting numbers", "Recognize where they sit in the hierarchy", "Contrast with other types"]},
        {"label": "Integer numbers", "section": "Integer numbers",
         "subgoals": ["Include zero and the negatives", "Locate integers on the number line", "Classify examples"]},
        {"label": "Rational numbers", "section": "Rational numbers",
         "subgoals": ["Write a rational as a ratio of integers", "Recognize repeating decimals as rational", "Distinguish from irrationals"]},
    ],
    "prealg.var.concept": [
        {"label": "Introduction to variables", "section": "Introduction to Variables",
         "subgoals": ["Recognize a variable as a placeholder", "Read expressions with variables", "Choose a letter for an unknown"]},
        {"label": "Algebraic expressions", "section": "Algebraic Expressions",
         "subgoals": ["Build an expression from a phrase", "Identify terms and coefficients", "Evaluate when values are given"]},
        {"label": "Evaluating with exponents and radicals", "section": "Evaluating Expressions with Exponents, Absolute Value, and Radicals",
         "subgoals": ["Apply exponents to the variable", "Handle absolute values", "Evaluate radicals"]},
    ],
    # ── Prealgebra · Ratio ────────────────────────────────
    "ratio.basics.concept": [
        {"label": "Write a ratio as a fraction", "section": "Write a Ratio as a Fraction",
         "subgoals": ["Write the comparison as a fraction", "Order the quantities correctly", "Simplify when possible"]},
        {"label": "Ratios", "section": "Ratios",
         "subgoals": ["Define a ratio as a comparison", "Use colon and fraction forms", "Read ratios in context"]},
        {"label": "Ratios involving decimals", "section": "Ratios Involving Decimals",
         "subgoals": ["Clear the decimals", "Scale to whole numbers", "Simplify the ratio"]},
    ],
    "ratio.ops.proportion": [
        {"label": "Definition of proportion", "section": "Use the Definition of Proportion",
         "subgoals": ["Set two ratios equal", "Identify the four terms", "State the proportion"]},
        {"label": "Proportion", "section": "Proportion",
         "subgoals": ["Recognize a proportion in a problem", "Write the two equal ratios", "Prepare to solve"]},
        {"label": "Cross products", "section": "Cross Products of a Proportion",
         "subgoals": ["Multiply across the equals sign", "Set the products equal", "Solve for the unknown"]},
    ],
    "ratio.ops.rate": [
        {"label": "Write a rate as a fraction", "section": "Write a Rate as a Fraction",
         "subgoals": ["Put the two units in ratio form", "Keep the units in order", "Simplify the fraction"]},
        {"label": "Rate", "section": "Rate",
         "subgoals": ["Define a rate with two different units", "Recognize rates in problems", "Compare rates"]},
        {"label": "Find unit rates", "section": "Find Unit Rates",
         "subgoals": ["Divide to get a denominator of 1", "State the per-unit value", "Use the unit rate to compare"]},
    ],
    "ratio.ops.scale": [
        {"label": "Applications of ratios", "section": "Applications of Ratios",
         "subgoals": ["Read the scale ratio", "Set up the proportion", "Solve for the real size"]},
        {"label": "Ratios of measurements", "section": "Ratios of Two Measurements in Different Units",
         "subgoals": ["Convert the units first", "Write the ratio", "Simplify"]},
        {"label": "Write a rate as a fraction", "section": "Write a Rate as a Fraction",
         "subgoals": ["Write map distance over real distance", "Keep units consistent", "Interpret the scale"]},
    ],
    "ratio.ops.simplify": [
        {"label": "Write a ratio as a fraction", "section": "Write a Ratio as a Fraction",
         "subgoals": ["Write the comparison as a fraction", "Cancel common units", "Reduce the numbers"]},
        {"label": "Ratios", "section": "Ratios",
         "subgoals": ["Express the ratio in lowest terms", "Use equivalent forms", "Keep the order of quantities"]},
        {"label": "Applications of ratios", "section": "Applications of Ratios",
         "subgoals": ["Simplify ratios in context", "Scale quantities by the same factor", "Verify the comparison"]},
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
            if kp["section"] and kp["section"] not in heads:
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
