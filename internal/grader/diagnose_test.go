package grader

import "testing"

// A diagnosis is only worth showing when it is certainly true. Silence is the
// right answer to "what went wrong?" far more often than any guess — a learner
// told the wrong reason for a mistake learns the wrong lesson, which is the
// opposite of what an explanation is for.
func TestDiagnose_Numeric(t *testing.T) {
	cases := map[string]struct {
		gradingType GradingType
		expected    string
		answer      string
		want        string
		why         string
	}{
		"sign flipped": {
			GradingNumeric, "-7", "7",
			"You have the right number with the sign flipped — the answer is 7.",
			"the magnitude is right and only the sign is wrong",
		},
		"off by one": {
			GradingNumeric, "12", "13",
			"You are exactly one away, so a count or a boundary is off by one.",
			"an adjacent integer is a distinct, common error",
		},
		"rounded too early": {
			GradingNumeric, "3.14", "3.1",
			"That is the answer rounded to fewer decimal places than it needs — 3.14 keeps them.",
			"the answer is the right value with precision lost",
		},
		"unrelated numbers": {
			GradingNumeric, "42", "17", "",
			"nothing distinguishes these, so say nothing",
		},
		"equal within tolerance": {
			GradingNumeric, "0.5", "1/2", "",
			"the numeric grader calls this correct; there is no mistake to name",
		},
		"blank answer": {
			GradingNumeric, "42", "", "",
			"an empty submission has no mistake to characterise",
		},
		"unparseable answer": {
			GradingNumeric, "42", "banana", "",
			"not a number, so no numeric relationship applies",
		},
		"choice type is not diagnosed numerically": {
			GradingMultipleChoice, "3", "3", "",
			"a multiple-choice type says nothing about the shape of an answer",
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			if got := Diagnose(tc.gradingType, tc.expected, tc.answer); got != tc.want {
				t.Errorf("Diagnose = %q, want %q — %s", got, tc.want, tc.why)
			}
		})
	}
}

func TestDiagnose_ComparisonReversed(t *testing.T) {
	cases := map[string]struct {
		expected, answer, want string
	}{
		"reversed":       {">", "<", "The comparison goes the other way round — it should be >, not <."},
		"other way":      {"<", ">", "The comparison goes the other way round — it should be <, not >."},
		"not equal":      {"=", "<", ""},
		"correct":        {">", ">", ""},
		"nonsense input": {">", "banana", ""},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			if got := Diagnose(GradingComparison, tc.expected, tc.answer); got != tc.want {
				t.Errorf("Diagnose = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestDiagnose_NamesTheOneWrongValue(t *testing.T) {
	// Knowing which one is off is the difference between a fixable mistake and
	// a shrug — but only when exactly one is wrong. Two wrong values mean the
	// method, not a slip, and saying "both" names nothing.
	cases := []struct {
		name        string
		gradingType GradingType
		expected    string
		answer      string
		want        string
	}{
		{"second value wrong", GradingTuple, "(2, 7)", "(2, 8)", "Only the 2nd value is wrong; the rest match."},
		{"third value wrong", GradingTuple, "(2, 7, 9)", "(2, 7, 4)", "Only the 3rd value is wrong; the rest match."},
		{
			"eleventh value uses the right suffix", GradingOrdering,
			"1,2,3,4,5,6,7,8,9,10,11", "1,2,3,4,5,6,7,8,9,10,99",
			"Only the 11th value is wrong; the rest match.",
		},
		{"two values wrong says nothing", GradingTuple, "(2, 7, 9)", "(3, 8, 9)", ""},
		{"length mismatch says nothing", GradingTuple, "(2, 7)", "(2, 7, 9)", ""},
		{"bare commas are accepted", GradingTuple, "2,7", "2,8", "Only the 2nd value is wrong; the rest match."},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := Diagnose(tc.gradingType, tc.expected, tc.answer); got != tc.want {
				t.Errorf("Diagnose = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestOrdinal(t *testing.T) {
	for n, want := range map[int]string{1: "1st", 2: "2nd", 3: "3rd", 4: "4th", 11: "11th", 12: "12th", 13: "13th", 21: "21st", 22: "22nd", 23: "23rd", 101: "101st"} {
		if got := ordinal(n); got != want {
			t.Errorf("ordinal(%d) = %q, want %q", n, got, want)
		}
	}
}
