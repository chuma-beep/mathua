package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// The header shows one public number, and it must not spend GitHub's
// unauthenticated budget (sixty requests an hour per address) once per reader.
// The count is fetched server side and held for six hours, so every browser
// asks this origin instead of GitHub. That also keeps the browser's network
// surface inside the site, which the static-export e2e checks require.
var githubStarsEndpoint = "https://api.github.com/repos/chuma-beep/mathua"

const githubStarsTTL = 6 * time.Hour

type githubStarsCache struct {
	mu    sync.Mutex
	count int
	at    time.Time
}

type githubStarsResponse struct {
	Count int `json:"count"`
}

func (s *Server) handleGitHubStars(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// One fetch at a time. The lock is held across the request so a burst of
	// first visitors does not start a burst of GitHub calls; the fetch has its
	// own five-second deadline.
	s.stars.mu.Lock()
	defer s.stars.mu.Unlock()

	if !s.stars.at.IsZero() && time.Since(s.stars.at) < githubStarsTTL {
		writeJSON(w, githubStarsResponse{Count: s.stars.count})
		return
	}
	count, err := fetchGitHubStars(r.Context())
	if err != nil {
		// A stale count is still true enough for a label, and better than
		// dropping the number because GitHub is briefly unavailable.
		if !s.stars.at.IsZero() {
			writeJSON(w, githubStarsResponse{Count: s.stars.count})
			return
		}
		writeError(w, "failed to read GitHub stars", http.StatusBadGateway)
		return
	}
	s.stars.count = count
	s.stars.at = time.Now()
	writeJSON(w, githubStarsResponse{Count: count})
}

func fetchGitHubStars(ctx context.Context) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, githubStarsEndpoint, nil)
	if err != nil {
		return 0, err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "mathua")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("github answered %d", res.StatusCode)
	}
	var body struct {
		StargazersCount int `json:"stargazers_count"`
	}
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		return 0, err
	}
	if body.StargazersCount < 0 {
		return 0, fmt.Errorf("github answered a negative count")
	}
	return body.StargazersCount, nil
}
