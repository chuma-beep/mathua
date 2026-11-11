// Package mailer owns outbound email. Phase 3 moved it out of internal/auth so
// notification delivery can send through the same provider path as password
// resets and email verification, without importing authentication.
package mailer

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	stdhtml "html"
	"io"
	"net/http"
	"net/smtp"
	"os"
	"strings"
	"time"
)

// Message is one plain-text email. The transport converts it to the provider's
// HTML shape when the provider requires one.
type Message struct {
	To      string
	Subject string
	Text    string
}

// Sender is the seam notification delivery uses in tests.
type Sender interface {
	Send(Message) error
}

// Mailer sends through SendByte when its API key is present, otherwise through
// SMTP. SendByte is preferred because it is the provider this deployment uses;
// SMTP remains supported for existing installations.
type Mailer struct {
	HTTPClient *http.Client
}

// New returns a mailer with the production timeout. A nil client would use
// http.DefaultClient, which has no timeout, so the field is always populated.
func New() *Mailer {
	return &Mailer{HTTPClient: &http.Client{Timeout: 10 * time.Second}}
}

// Configured reports whether either provider has enough configuration to be
// attempted. It is deliberately a cheap environment check, not a connection
// test; a provider outage still surfaces as a send error.
func Configured() bool {
	return sendByteConfigured() || strings.TrimSpace(os.Getenv("SMTP_HOST")) != ""
}

// Send validates and dispatches one message.
func (m *Mailer) Send(msg Message) error {
	if strings.TrimSpace(msg.To) == "" {
		return errors.New("mailer: recipient is required")
	}
	if strings.TrimSpace(msg.Subject) == "" {
		return errors.New("mailer: subject is required")
	}
	if sendByteConfigured() {
		return m.sendSendByte(msg)
	}
	return m.sendSMTP(msg)
}

func sendByteConfigured() bool {
	return strings.TrimSpace(os.Getenv("SENDBYTE_API_KEY")) != ""
}

func sendByteBaseURL() string {
	if v := strings.TrimSpace(os.Getenv("SENDBYTE_API_URL")); v != "" {
		return strings.TrimRight(v, "/")
	}
	return "https://api.sendbyte.africa"
}

func sendByteFrom() string {
	if v := strings.TrimSpace(os.Getenv("SENDBYTE_FROM")); v != "" {
		return v
	}
	return strings.TrimSpace(os.Getenv("SMTP_FROM"))
}

func (m *Mailer) sendSendByte(msg Message) error {
	key := strings.TrimSpace(os.Getenv("SENDBYTE_API_KEY"))
	if key == "" {
		return errors.New("sendbyte: SENDBYTE_API_KEY is not set")
	}
	from := sendByteFrom()
	if from == "" {
		return errors.New("sendbyte: SENDBYTE_FROM is not set")
	}
	payload, err := json.Marshal(map[string]string{
		"from":    from,
		"to":      msg.To,
		"subject": msg.Subject,
		"html":    plainTextToHTML(msg.Text),
	})
	if err != nil {
		return fmt.Errorf("sendbyte: encode request: %w", err)
	}
	req, err := http.NewRequest(http.MethodPost, sendByteBaseURL()+"/v1/emails", bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("sendbyte: build request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+key)
	req.Header.Set("Content-Type", "application/json")
	client := m.HTTPClient
	if client == nil {
		client = New().HTTPClient
	}
	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("sendbyte: %w", err)
	}
	defer res.Body.Close()
	if res.StatusCode < 200 || res.StatusCode > 299 {
		// The provider explains refusals (unverified sender, sandbox key, bad
		// address) in the body; surface it so the next attempt is diagnosable.
		snippet, _ := io.ReadAll(io.LimitReader(res.Body, 512))
		return fmt.Errorf("sendbyte: %s: %s", res.Status, strings.TrimSpace(string(snippet)))
	}
	return nil
}

func (m *Mailer) sendSMTP(msg Message) error {
	host := strings.TrimSpace(os.Getenv("SMTP_HOST"))
	if host == "" {
		return errors.New("smtp: SMTP_HOST is not set")
	}
	port := strings.TrimSpace(os.Getenv("SMTP_PORT"))
	if port == "" {
		port = "587"
	}
	user := os.Getenv("SMTP_USER")
	pass := os.Getenv("SMTP_PASS")
	from := os.Getenv("SMTP_FROM")
	if from == "" {
		from = user
	}
	if strings.TrimSpace(from) == "" {
		return errors.New("smtp: SMTP_FROM is not set")
	}
	var auth smtp.Auth
	if user != "" {
		auth = smtp.PlainAuth("", user, pass, host)
	}
	body := []byte("From: " + from + "\r\nTo: " + msg.To + "\r\nSubject: " + msg.Subject + "\r\n\r\n" + msg.Text)
	return smtp.SendMail(host+":"+port, auth, from, []string{msg.To}, body)
}

// plainTextToHTML keeps the mail templates readable as paragraphs. The bodies
// are authored as plain text with CRLF line breaks, and turning them into one
// <p> per non-empty line preserves the shape without an HTML template.
func plainTextToHTML(body string) string {
	var out strings.Builder
	for _, line := range strings.Split(strings.ReplaceAll(body, "\r\n", "\n"), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		out.WriteString("<p>")
		out.WriteString(stdhtml.EscapeString(line))
		out.WriteString("</p>")
	}
	return out.String()
}
