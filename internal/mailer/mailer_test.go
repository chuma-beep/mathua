package mailer

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSend_PostsTheDocumentedSendByteRequest(t *testing.T) {
	var gotMethod, gotPath, gotAuth, gotType string
	var body map[string]string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		gotAuth = r.Header.Get("Authorization")
		gotType = r.Header.Get("Content-Type")
		raw, _ := io.ReadAll(r.Body)
		if err := json.Unmarshal(raw, &body); err != nil {
			t.Errorf("request body is not JSON: %v", err)
		}
		w.WriteHeader(http.StatusAccepted)
	}))
	defer srv.Close()

	t.Setenv("SENDBYTE_API_KEY", "sb_live_test")
	t.Setenv("SENDBYTE_FROM", "Mathua <hello@mathua.com>")
	t.Setenv("SENDBYTE_API_URL", srv.URL)
	t.Setenv("SMTP_HOST", "")

	err := New().Send(Message{
		To:      "learner@example.com",
		Subject: "Mathua password reset",
		Text:    "Line one\r\nLine two <unsafe>",
	})
	if err != nil {
		t.Fatalf("send: %v", err)
	}
	if gotMethod != http.MethodPost || gotPath != "/v1/emails" {
		t.Errorf("request = %s %s, want POST /v1/emails", gotMethod, gotPath)
	}
	if gotAuth != "Bearer sb_live_test" {
		t.Errorf("authorization = %q", gotAuth)
	}
	if !strings.HasPrefix(gotType, "application/json") {
		t.Errorf("content type = %q", gotType)
	}
	if body["from"] != "Mathua <hello@mathua.com>" || body["to"] != "learner@example.com" {
		t.Errorf("addresses = %+v", body)
	}
	if body["subject"] != "Mathua password reset" {
		t.Errorf("subject = %q", body["subject"])
	}
	if !strings.Contains(body["html"], "<p>Line one</p>") || !strings.Contains(body["html"], "<p>Line two &lt;unsafe&gt;</p>") {
		t.Errorf("html = %q, want one escaped paragraph per line", body["html"])
	}
}

func TestSend_RefusesWithoutKeyOrSender(t *testing.T) {
	t.Setenv("SENDBYTE_API_KEY", "")
	t.Setenv("SMTP_HOST", "")
	t.Setenv("SENDBYTE_FROM", "hello@mathua.com")
	if err := New().Send(Message{To: "a@b.c", Subject: "s", Text: "b"}); err == nil {
		t.Error("sent without an API key or SMTP host")
	}

	t.Setenv("SENDBYTE_API_KEY", "sb_live_test")
	t.Setenv("SENDBYTE_FROM", "")
	t.Setenv("SMTP_FROM", "")
	if err := New().Send(Message{To: "a@b.c", Subject: "s", Text: "b"}); err == nil {
		t.Error("sent without a From address")
	}
}

func TestSend_ReportsTheProviderRefusal(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		_, _ = w.Write([]byte(`{"error":"sender domain is not verified"}`))
	}))
	defer srv.Close()

	t.Setenv("SENDBYTE_API_KEY", "sb_test")
	t.Setenv("SENDBYTE_FROM", "hello@example.com")
	t.Setenv("SENDBYTE_API_URL", srv.URL)

	err := New().Send(Message{To: "learner@example.com", Subject: "s", Text: "b"})
	if err == nil {
		t.Fatal("expected the provider refusal to surface")
	}
	if !strings.Contains(err.Error(), "422") || !strings.Contains(err.Error(), "not verified") {
		t.Errorf("error = %v, want the status and provider message", err)
	}
}

func TestConfiguredAndSendBytePreference(t *testing.T) {
	hit := false
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hit = true
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	t.Setenv("SENDBYTE_API_KEY", "sb_test")
	t.Setenv("SENDBYTE_FROM", "hello@mathua.com")
	t.Setenv("SENDBYTE_API_URL", srv.URL)
	t.Setenv("SMTP_HOST", "")

	if !Configured() {
		t.Fatal("Configured() is false with a SendByte key")
	}
	if err := New().Send(Message{To: "learner@example.com", Subject: "Mathua password reset", Text: "body"}); err != nil {
		t.Fatalf("Send: %v", err)
	}
	if !hit {
		t.Error("Send did not reach SendByte")
	}
}

func TestPlainTextToHTML(t *testing.T) {
	cases := []struct {
		input string
		want  string
	}{
		{"", ""},
		{"   \n   \n   ", ""},
		{
			"Hello world",
			"<p>Hello world</p>",
		},
		{
			"First paragraph\r\n\r\nSecond paragraph",
			"<p>First paragraph</p><p>Second paragraph</p>",
		},
		{
			"Special chars: <script>alert('test')</script> & \"quotes\"",
			"<p>Special chars: &lt;script&gt;alert(&#39;test&#39;)&lt;/script&gt; &amp; &#34;quotes&#34;</p>",
		},
	}

	for _, c := range cases {
		got := plainTextToHTML(c.input)
		if got != c.want {
			t.Errorf("plainTextToHTML(%q): got %q, want %q", c.input, got, c.want)
		}
	}
}
