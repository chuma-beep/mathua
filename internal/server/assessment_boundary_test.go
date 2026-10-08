package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/chuma-beep/mathua/internal/auth"
	"github.com/chuma-beep/mathua/internal/concepts"
	"github.com/chuma-beep/mathua/internal/engine"
	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/planning"
	"github.com/chuma-beep/mathua/internal/storage"
)

// The assessment boundary, as executable code.
//
// Diagnostics and tests measure knowledge. They do not teach. So no verdict —
// including the final one — may carry the served worked solution, and no verdict
// may carry a diagnosis that lets the learner derive the answer.
//
// This file exists because the opposite was once a deliberate feature:
// TestQuizAnswer_ExplanationOnBothVerdicts asserted the solution reached the
// client on every verdict, and TestStudyAnswer_AnchoredNotClientSupplied asserts
// the same for Learn. The first of those was wrong for an assessment and is
// inverted below; the second is right and stays, which is the whole point of the
// boundary being drawn per handler rather than per response shape.
//
// The served problem is testGen's: question "2+2=?", answer "4", explanation
// "2+2=4". The explanation is a distinctive string, so "does it appear anywhere
// in the response body" is a real question and not a restatement of "is the key
// called explanation". A future handler that smuggles the solution out under a
// different key fails these tests, which key-name assertions would not catch.

const (
	servedAnswer       = "4"
	servedExplanation  = "2+2=4"
	derivedDiagnosis   = "You are exactly one away"
	answerTooFarOff    = "999" // no diagnosis is derivable from this
	answerOneAwayWrong = "3"   // |4-3| == 1, the diagnosis that leaks
)

// assessmentServer is a two-concept server with auth, so a quiz covers two
// questions: the first answer is mid-assessment and the second is final. The
// single-concept guestServer collapses every quiz to one question and cannot
// express the mid-flight/final distinction, which is exactly where a "strip it
// unless done" rule would leak.
func assessmentServer(t *testing.T) (*http.ServeMux, string) {
	t.Helper()
	d, _ := concepts.Build([]concepts.Concept{
		{ID: "a", Label: "A", Domain: "d", GradingType: "numeric", Prerequisites: []string{},
			MasteryThreshold: concepts.MasteryThreshold{Streak: 3, AvgTimeSeconds: 60}},
		{ID: "b", Label: "B", Domain: "d", GradingType: "numeric", Prerequisites: []string{"a"},
			MasteryThreshold: concepts.MasteryThreshold{Streak: 3, AvgTimeSeconds: 60}},
	})
	store, _ := storage.NewSQLiteStore(":memory:")
	t.Cleanup(func() { store.Close() })
	reg := generator.NewRegistry()
	reg.Register("a", &testGen{})
	reg.Register("b", &testGen{})
	s := New(engine.New(store, d, reg, nil, planning.New(d)), store, auth.New(store))
	mux := http.NewServeMux()
	s.Register(mux)

	rec := postGuest(t, mux, `{}`, "")
	if rec.Code != 200 {
		t.Fatalf("create guest: %d %s", rec.Code, rec.Body.String())
	}
	var res map[string]interface{}
	_ = json.Unmarshal(rec.Body.Bytes(), &res)
	token, _ := res["token"].(string)
	if token == "" {
		t.Fatal("expected a guest token")
	}
	return mux, token
}

// postJSONRaw returns the raw body so a test can ask what is in the response
// rather than only what the decoded map exposes. An explanation tucked into an
// unexpected key is still an explanation.
func postJSONRaw(t *testing.T, mux *http.ServeMux, path, body, token string) (int, string) {
	t.Helper()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("POST", path, bytes.NewReader([]byte(body)))
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	mux.ServeHTTP(rec, req)
	return rec.Code, rec.Body.String()
}

// assertNoSolution fails if the response body carries the served solution or a
// solution-bearing diagnosis. The two keys are checked explicitly so a failure
// names the leak, and the body is grepped so a leak under a new key is caught.
func assertNoSolution(t *testing.T, what, body string) {
	t.Helper()
	if strings.Contains(body, servedExplanation) {
		t.Errorf("%s: response body contains the served explanation %q: %s", what, servedExplanation, body)
	}
	if strings.Contains(body, derivedDiagnosis) {
		t.Errorf("%s: response body contains an answer-revealing diagnosis: %s", what, body)
	}
	var decoded map[string]interface{}
	if err := json.Unmarshal([]byte(body), &decoded); err == nil {
		if _, ok := decoded["explanation"]; ok {
			t.Errorf("%s: response has an explanation key: %s", what, body)
		}
		if _, ok := decoded["diagnosis"]; ok {
			t.Errorf("%s: response has a diagnosis key: %s", what, body)
		}
	}
}

// startDiagnostic opens a diagnostic over both concepts and returns the session
// and first concept id.
func startDiagnostic(t *testing.T, mux *http.ServeMux) (sessionID, conceptID string) {
	t.Helper()
	code, start := postJSON(t, mux, "/api/goal/diagnostic", `{"name":"tester","concept_ids":["a","b"]}`)
	if code != 200 {
		t.Fatalf("start diagnostic: %d %v", code, start)
	}
	sessionID, _ = start["session_id"].(string)
	conceptID, _ = start["concept_id"].(string)
	if sessionID == "" || conceptID == "" {
		t.Fatalf("expected session and concept, got %v", start)
	}
	return sessionID, conceptID
}

// 1. wrong diagnostic answer → no explanation, no derivable diagnosis.
func TestDiagnosticBoundary_WrongAnswerCarriesNoSolution(t *testing.T) {
	mux, _ := assessmentServer(t)
	sid, cid := startDiagnostic(t, mux)

	code, body := postJSONRaw(t, mux, "/api/goal/diagnostic/answer",
		`{"session_id":"`+sid+`","concept_id":"`+cid+`","answer":"`+answerOneAwayWrong+`","elapsed":5.0}`, "")
	if code != 200 {
		t.Fatalf("answer: %d %s", code, body)
	}
	assertNoSolution(t, "diagnostic wrong answer", body)

	// The verdict itself must still be there: the boundary strips teaching, not measurement.
	var decoded map[string]interface{}
	_ = json.Unmarshal([]byte(body), &decoded)
	if _, ok := decoded["correct"]; !ok {
		t.Errorf("the response must still report correctness: %s", body)
	}
}

// 2. correct diagnostic answer → no explanation.
func TestDiagnosticBoundary_CorrectAnswerCarriesNoSolution(t *testing.T) {
	mux, _ := assessmentServer(t)
	sid, cid := startDiagnostic(t, mux)

	code, body := postJSONRaw(t, mux, "/api/goal/diagnostic/answer",
		`{"session_id":"`+sid+`","concept_id":"`+cid+`","answer":"`+servedAnswer+`","elapsed":5.0}`, "")
	if code != 200 {
		t.Fatalf("answer: %d %s", code, body)
	}
	assertNoSolution(t, "diagnostic correct answer", body)
}

// 3. wrong quiz answer, mid-assessment (done:false) → no explanation.
// The two-concept quiz serves "a" then "b", so the first answer is not the last.
func TestQuizBoundary_WrongMidFlightCarriesNoSolution(t *testing.T) {
	mux, token := assessmentServer(t)
	sess := quizPost(t, mux, "/api/quiz/session", `{}`, token)

	code, body := postJSONRaw(t, mux, "/api/quiz/answer",
		`{"session_id":"`+sess["session_id"].(string)+`","concept_id":"a","answer":"`+answerOneAwayWrong+`","elapsed":5}`, token)
	if code != 200 {
		t.Fatalf("answer: %d %s", code, body)
	}
	var decoded map[string]interface{}
	_ = json.Unmarshal([]byte(body), &decoded)
	if done, _ := decoded["done"].(bool); done {
		t.Fatalf("expected a mid-flight quiz verdict, got done:true: %s", body)
	}
	assertNoSolution(t, "quiz wrong mid-flight", body)
}

// 4. correct quiz answer, mid-assessment → no explanation.
func TestQuizBoundary_CorrectMidFlightCarriesNoSolution(t *testing.T) {
	mux, token := assessmentServer(t)
	sess := quizPost(t, mux, "/api/quiz/session", `{}`, token)

	code, body := postJSONRaw(t, mux, "/api/quiz/answer",
		`{"session_id":"`+sess["session_id"].(string)+`","concept_id":"a","answer":"`+servedAnswer+`","elapsed":5}`, token)
	if code != 200 {
		t.Fatalf("answer: %d %s", code, body)
	}
	var decoded map[string]interface{}
	_ = json.Unmarshal([]byte(body), &decoded)
	if done, _ := decoded["done"].(bool); done {
		t.Fatalf("expected a mid-flight quiz verdict, got done:true: %s", body)
	}
	assertNoSolution(t, "quiz correct mid-flight", body)
}

// 5 and 6. The final quiz answer, right and wrong. `done:true` ends the
// assessment; it does not begin the lesson. This is the case a rule phrased as
// "reveal once finished" would get wrong, and it is why the boundary is stated
// per handler rather than tied to the flag.
func TestQuizBoundary_FinalVerdictCarriesNoSolution(t *testing.T) {
	for _, tc := range []struct {
		name   string
		answer string
		wrong  bool
	}{
		{"wrong final", answerOneAwayWrong, true},
		{"correct final", servedAnswer, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			mux, token := assessmentServer(t)
			sess := quizPost(t, mux, "/api/quiz/session", `{}`, token)
			sid := sess["session_id"].(string)

			// First concept: answer and stay mid-flight.
			code, first := postJSONRaw(t, mux, "/api/quiz/answer",
				`{"session_id":"`+sid+`","concept_id":"a","answer":"`+answerTooFarOff+`","elapsed":5}`, token)
			if code != 200 {
				t.Fatalf("first answer: %d %s", code, first)
			}
			var mid map[string]interface{}
			_ = json.Unmarshal([]byte(first), &mid)
			if done, _ := mid["done"].(bool); done {
				t.Fatalf("expected the first of two questions to be mid-flight: %s", first)
			}

			// Second concept: the final verdict.
			code, body := postJSONRaw(t, mux, "/api/quiz/answer",
				`{"session_id":"`+sid+`","concept_id":"b","answer":"`+tc.answer+`","elapsed":5}`, token)
			if code != 200 {
				t.Fatalf("final answer: %d %s", code, body)
			}
			var decoded map[string]interface{}
			_ = json.Unmarshal([]byte(body), &decoded)
			if done, _ := decoded["done"].(bool); !done {
				t.Fatalf("expected the final verdict to be done:true: %s", body)
			}
			if got, _ := decoded["correct"].(bool); got == tc.wrong {
				t.Fatalf("correct = %v, want wrong = %v: %s", got, tc.wrong, body)
			}
			assertNoSolution(t, "quiz final verdict", body)

			// And the solution is still on the session, because the boundary strips it on the
			// way out and not on the way in: the attempts transcript and grading are untouched.
			if !strings.Contains(body, `"done":true`) {
				t.Errorf("the flag's own encoding should be unchanged: %s", body)
			}
		})
	}
}

// 7. The granular feedback a learner is allowed: a verdict, and nothing that
// reconstructs the answer. This asserts the allowed half of the contract so the
// boundary cannot be "satisfied" by blanking every response.
func TestAssessmentBoundary_KeepsVerdictDropsDerivation(t *testing.T) {
	mux, token := assessmentServer(t)
	sess := quizPost(t, mux, "/api/quiz/session", `{}`, token)

	_, body := postJSONRaw(t, mux, "/api/quiz/answer",
		`{"session_id":"`+sess["session_id"].(string)+`","concept_id":"a","answer":"`+answerOneAwayWrong+`","elapsed":5}`, token)

	// The wrong answer is reported as wrong — that is assessment state and it must survive.
	if !strings.Contains(body, `"correct":false`) {
		t.Errorf("the verdict must still cross the boundary: %s", body)
	}
	// The diagnosis that would follow from it must not.
	if strings.Contains(body, "one away") || strings.Contains(body, "off by one") {
		t.Errorf("derivation-based diagnosis crossed the boundary: %s", body)
	}
}

// 8. Learn is unchanged. /api/study/answer is the teaching seam: it is
// *supposed* to return the worked solution and the expected answer, because
// Learn teaches and teaching is what it is for. If the assessment helper were
// ever applied to it, this fails — which is the point of pinning it.
func TestStudyAnswer_LearnSeamStillTeaches(t *testing.T) {
	_, mux, _ := guestServer(t)
	token, studentID := quizGuest(t, mux)

	// Serve a practice set, which writes the anchor the answer is graded against.
	rec := httptest.NewRecorder()
	get := httptest.NewRequest("GET", "/api/lessons/a/practice?student_id="+studentID+"&count=1", nil)
	get.Header.Set("Authorization", "Bearer "+token)
	mux.ServeHTTP(rec, get)
	if rec.Code != 200 {
		t.Fatalf("practice: %d %s", rec.Code, rec.Body.String())
	}
	var set struct {
		Questions []struct{ Question, Answer, Explanation string } `json:"questions"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &set); err != nil {
		t.Fatalf("decode practice: %v", err)
	}
	if len(set.Questions) == 0 {
		t.Fatal("expected a served question")
	}

	code, body := postJSONRaw(t, mux, "/api/study/answer",
		`{"concept_id":"a","answer":"`+answerOneAwayWrong+`","elapsed":5.0,"question":`+
			quoteJSON(set.Questions[0].Question)+`,"student_id":"`+studentID+`"}`, token)
	if code != 200 {
		t.Fatalf("study answer: %d %s", code, body)
	}
	var res struct {
		Explanation    string `json:"explanation"`
		ExpectedAnswer string `json:"expected_answer"`
	}
	if err := json.Unmarshal([]byte(body), &res); err != nil {
		t.Fatalf("decode study answer: %v", err)
	}
	if res.Explanation == "" {
		t.Errorf("Learn must still return the worked solution: %s", body)
	}
	if res.ExpectedAnswer == "" {
		t.Errorf("Learn must still return the expected answer for its corrective feedback: %s", body)
	}
}

func quoteJSON(s string) string {
	b, _ := json.Marshal(s)
	return string(b)
}
