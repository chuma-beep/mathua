package notifications

import (
	"fmt"
	"os"
	"strings"

	webpush "github.com/SherClockHolmes/webpush-go"
	"github.com/chuma-beep/mathua/internal/storage"
)

const (
	// PushTTL is how long a push service may hold a message for an offline
	// device. Twelve hours covers an overnight phone without keeping a stale
	// study reminder alive for days.
	PushTTL = 12 * 60 * 60

	defaultVAPIDSubject = "https://mathua.com"
)

// WebPushOptions are the VAPID credentials and the request tunables. The keys
// come from the deployment environment; a missing pair means push is simply
// unavailable, not that the server fails to start.
type WebPushOptions struct {
	PublicKey  string
	PrivateKey string
	Subject    string
	TTL        int
	Client     webpush.HTTPClient
}

// WebPushFromEnv reads the VAPID configuration. It returns false when either
// key is missing, which is the normal state of a development machine.
func WebPushFromEnv() (WebPushOptions, bool) {
	opts := WebPushOptions{
		PublicKey:  strings.TrimSpace(os.Getenv("MATHUA_VAPID_PUBLIC_KEY")),
		PrivateKey: strings.TrimSpace(os.Getenv("MATHUA_VAPID_PRIVATE_KEY")),
		Subject:    strings.TrimSpace(os.Getenv("MATHUA_VAPID_SUBJECT")),
		TTL:        PushTTL,
	}
	if opts.PublicKey == "" || opts.PrivateKey == "" {
		return WebPushOptions{}, false
	}
	if opts.Subject == "" {
		opts.Subject = defaultVAPIDSubject
	}
	return opts, true
}

// WebPushSender sends an encrypted message through a push service.
type WebPushSender struct {
	publicKey  string
	privateKey string
	subject    string
	ttl        int
	client     webpush.HTTPClient
}

// NewWebPushSender validates the credentials and builds the sender. A nil
// sender from WebPushFromEnv is expected when push is not configured.
func NewWebPushSender(opts WebPushOptions) (*WebPushSender, error) {
	if strings.TrimSpace(opts.PublicKey) == "" || strings.TrimSpace(opts.PrivateKey) == "" {
		return nil, fmt.Errorf("web push: VAPID public and private keys are required")
	}
	subject := strings.TrimSpace(opts.Subject)
	if subject == "" {
		subject = defaultVAPIDSubject
	}
	ttl := opts.TTL
	if ttl <= 0 {
		ttl = PushTTL
	}
	return &WebPushSender{
		publicKey:  strings.TrimSpace(opts.PublicKey),
		privateKey: strings.TrimSpace(opts.PrivateKey),
		subject:    subject,
		ttl:        ttl,
		client:     opts.Client,
	}, nil
}

// Send posts one payload to one endpoint. A 404 or 410 means the subscription
// will never work again, so it is reported as PushGoneError and the worker
// deletes it. Any other non-2xx response is a normal error worth a retry.
func (s *WebPushSender) Send(sub storage.NotificationPushSubscription, payload []byte) error {
	if s == nil {
		return fmt.Errorf("web push: sender is nil")
	}
	if strings.TrimSpace(sub.Endpoint) == "" {
		return &PushGoneError{Status: 410, Endpoint: sub.Endpoint}
	}
	resp, err := webpush.SendNotification(payload, &webpush.Subscription{
		Endpoint: sub.Endpoint,
		Keys: webpush.Keys{
			Auth:   sub.Auth,
			P256dh: sub.P256DH,
		},
	}, &webpush.Options{
		HTTPClient:      s.client,
		Subscriber:      s.subject,
		TTL:             s.ttl,
		VAPIDPublicKey:  s.publicKey,
		VAPIDPrivateKey: s.privateKey,
	})
	if err != nil {
		return fmt.Errorf("web push: %w", err)
	}
	if resp == nil {
		return fmt.Errorf("web push: empty response")
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}
	if resp.StatusCode == 404 || resp.StatusCode == 410 {
		return &PushGoneError{Status: resp.StatusCode, Endpoint: sub.Endpoint}
	}
	return fmt.Errorf("web push: endpoint returned %s", resp.Status)
}
