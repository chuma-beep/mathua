package server

import "net/http"

// handleAdminAccountability exposes the resolved XP policy and process-local
// counters. xp_events remains the durable source of truth; this view answers
// what the policy decided and what storage applied since startup.
func (s *Server) handleAdminAccountability(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, map[string]interface{}{
		"config":  s.eng.AccountabilityConfig(),
		"metrics": s.eng.AccountabilityMetrics(),
	})
}
