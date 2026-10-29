// Package fire holds the runtime mode for FIRe, the implicit-repetition
// ledger. Shadow records observations only. Active applies banked credit to
// the scheduler for an explicitly acknowledged, server-side cohort.
package fire

import (
	"fmt"
	"strings"
)

// Mode is the resolved startup mode for FIRe.
type Mode string

const (
	// Shadow records ledger rows only. It never writes ConceptProgress, FSRS
	// state, due dates or review timestamps, and the scheduler ignores its
	// banked ordinals.
	Shadow Mode = "shadow"
	// Active may apply banked ordinals to eligible targets. It requires an
	// explicit acknowledgement and a non-empty server-side scope, either a
	// student allowlist or the "all" scope. Applying credit never grants XP or
	// changes mastery evidence.
	Active Mode = "active"
	// Off does not allocate credit and does not apply FIRe outbox effects.
	Off Mode = "off"
)

// ActivationToken is the exact value the server requires before it will run
// FIRe in active mode. A client-provided flag is never consulted.
const ActivationToken = "I acknowledge that FIRe may apply schedule changes"

// Config is the fully resolved FIRe startup configuration. Mode has already
// failed closed to Shadow if Active was requested without its guards.
type Config struct {
	Mode     Mode
	Cohort   string
	Allow    map[string]struct{}
	AllowAll bool
	Degraded string
}

// Parse resolves the configured mode. It is total and fail-closed: empty,
// whitespace, malformed and unknown values all become Shadow.
func Parse(raw string) Mode {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case string(Off):
		return Off
	case string(Active):
		return Active
	default:
		return Shadow
	}
}

// ParseConfig resolves the mode and the activation guards. Active is only
// returned when the acknowledgement is exact, the cohort name is stable and
// non-empty and the scope is either the "all" token or at least one student ID.
// Every other combination degrades to Shadow and records why, so a typo cannot
// activate scheduling effects.
func ParseConfig(raw, acknowledgement, cohort, allowlist string) Config {
	cfg := Config{
		Mode:   Parse(raw),
		Cohort: strings.TrimSpace(cohort),
	}
	allow, allowAll, scopeErr := parseAllowlist(allowlist)
	cfg.Allow = allow
	cfg.AllowAll = allowAll
	if cfg.Mode != Active {
		return cfg
	}
	if strings.TrimSpace(acknowledgement) != ActivationToken {
		cfg.Mode = Shadow
		cfg.Degraded = "activation acknowledgement is missing or invalid"
		return cfg
	}
	if cfg.Cohort == "" {
		cfg.Mode = Shadow
		cfg.Degraded = "activation cohort is empty"
		return cfg
	}
	if scopeErr != "" {
		cfg.Mode = Shadow
		cfg.Degraded = scopeErr
		return cfg
	}
	if !cfg.AllowAll && len(cfg.Allow) == 0 {
		cfg.Mode = Shadow
		cfg.Degraded = "activation scope is empty"
		return cfg
	}
	return cfg
}

// parseAllowlist reads the scope. The bare token "all" or "*" selects every
// learner, including accounts created later. Mixing that token with explicit
// IDs is ambiguous, so it degrades instead of silently choosing one meaning.
func parseAllowlist(raw string) (map[string]struct{}, bool, string) {
	trimmed := strings.TrimSpace(raw)
	if strings.EqualFold(trimmed, "all") || trimmed == "*" {
		return map[string]struct{}{}, true, ""
	}
	out := map[string]struct{}{}
	for _, part := range strings.Split(trimmed, ",") {
		id := strings.TrimSpace(part)
		if id == "" {
			continue
		}
		if strings.EqualFold(id, "all") || id == "*" {
			return map[string]struct{}{}, false, "activation scope mixes all with explicit student IDs"
		}
		out[id] = struct{}{}
	}
	return out, false, ""
}

// ActiveFor reports whether FIRe may apply credit for this learner. Both the
// mode and the server-side scope must agree. An empty learner ID is never
// active, even under the all-learners scope.
func (c Config) ActiveFor(studentID string) bool {
	if c.Mode != Active {
		return false
	}
	id := strings.TrimSpace(studentID)
	if id == "" {
		return false
	}
	if c.AllowAll {
		return true
	}
	_, ok := c.Allow[id]
	return ok
}

// Degrade returns a Shadow copy carrying the reason the requested mode was
// refused. Process edges use it when a durable activation row is missing or
// does not match the deployed configuration.
func (c Config) Degrade(reason string) Config {
	c.Mode = Shadow
	c.Degraded = reason
	return c
}

// CollectsLedger reports whether FIRe should record observations. Shadow and
// Active both do. The method is the single question the engine and outbox ask,
// so a new mode cannot accidentally inherit ledger writes.
func (m Mode) CollectsLedger() bool { return m == Shadow || m == Active }

// String keeps log lines readable and stable.
func (m Mode) String() string { return string(m) }

// Describe returns a short operational summary, including a degradation
// reason when active configuration was rejected.
func (c Config) Describe() string {
	if c.Degraded != "" {
		return fmt.Sprintf("%s (degraded: %s)", c.Mode, c.Degraded)
	}
	if c.Mode == Active {
		if c.AllowAll {
			return fmt.Sprintf("%s (cohort %s, all learners)", c.Mode, c.Cohort)
		}
		return fmt.Sprintf("%s (cohort %s, %d learners)", c.Mode, c.Cohort, len(c.Allow))
	}
	return string(c.Mode)
}
