package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// The deployed routing of a Next.js static export, which is where the administration surface
// actually lives in production.
//
// A `next build && next export` writes one HTML file per route at the top level
// (`out/admin.html`) and, for any route that has child routes, a *directory* of the same name
// (`out/admin/users.html`). Those two facts collide in one place. This test exists because the
// collision shipped: `nextStaticFS` checked only that the path existed, so `/admin` matched the
// directory, went to http.FileServer, was redirected to `/admin/` and then rendered a listing of
// `audit.html`, `reports.html`, `users.html`. A 301 to a file index, on the first page of the
// administration workflow, in production.
//
// Every assertion below is written as the request a browser makes, against a real directory
// laid out the way `next export` lays one out. A test that mocked the filesystem would have
// passed against the broken version.

func writeExport(t *testing.T, files map[string]string) string {
	t.Helper()
	root := t.TempDir()
	for name, body := range files {
		full := filepath.Join(root, filepath.FromSlash(name))
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatalf("mkdir %s: %v", name, err)
		}
		if err := os.WriteFile(full, []byte(body), 0o644); err != nil {
			t.Fatalf("write %s: %v", name, err)
		}
	}
	return root
}

// A small export with the shape that matters: a parent route that has children, a childless
// route, and an asset.
func adminExport(t *testing.T) string {
	return writeExport(t, map[string]string{
		"index.html":             "<html>home</html>",
		"learn.html":             "<html>learn</html>",
		"admin.html":             "<html>admin-overview</html>",
		"admin/users.html":       "<html>admin-users</html>",
		"admin/audit.html":       "<html>admin-audit</html>",
		"admin/reports.html":     "<html>admin-reports</html>",
		"docs.html":              "<html>docs</html>",
		"docs/architecture.html": "<html>docs-architecture</html>",
		"domains.html":           "<html>domains</html>",
		"favicon.ico":            "icon",
	})
}

func get(t *testing.T, root, path string) *httptest.ResponseRecorder {
	t.Helper()
	rec := httptest.NewRecorder()
	nextStaticFS(root).ServeHTTP(rec, httptest.NewRequest(http.MethodGet, path, nil))
	return rec
}

// The regression itself: `/admin` has a directory of the same name, and that directory must not
// be what answers.
func TestNextStaticFSServesThePageNotADirectoryListing(t *testing.T) {
	root := adminExport(t)
	for _, path := range []string{"/admin", "/docs"} {
		t.Run(path, func(t *testing.T) {
			rec := get(t, root, path)
			if rec.Code != 200 {
				t.Fatalf("status = %d, want 200 (a redirect means the directory matched)", rec.Code)
			}
			// The failure mode is a 200, not an error status: http.FileServer is happy to list a
			// directory. So the status alone would not have caught it.
			if strings.Contains(rec.Body.String(), "<pre>") || strings.Contains(rec.Body.String(), "Index of") {
				t.Fatalf("served a directory listing instead of the page:\n%s", rec.Body.String())
			}
			if !strings.Contains(rec.Body.String(), "admin-overview") && !strings.Contains(rec.Body.String(), "docs<") {
				t.Errorf("body does not look like the page: %s", rec.Body.String())
			}
		})
	}
}

// A childless route has no directory of its own, so it never hit this. Asserted so the fix
// cannot regress the case it was written for.
func TestNextStaticFSServesChildlessRoutes(t *testing.T) {
	root := adminExport(t)
	for path, want := range map[string]string{
		"/learn":   "learn",
		"/domains": "domains",
		"/docs":    "docs<",
	} {
		rec := get(t, root, path)
		if rec.Code != 200 {
			t.Errorf("%s: status = %d, want 200", path, rec.Code)
		}
		if !strings.Contains(rec.Body.String(), want) {
			t.Errorf("%s: body = %s", path, rec.Body.String())
		}
	}
}

// Child routes were already working and must keep working — they take the direct-file branch.
func TestNextStaticFSServesChildRoutes(t *testing.T) {
	root := adminExport(t)
	for path, want := range map[string]string{
		"/admin/users":       "admin-users",
		"/admin/audit":       "admin-audit",
		"/admin/reports":     "admin-reports",
		"/docs/architecture": "docs-architecture",
	} {
		rec := get(t, root, path)
		if rec.Code != 200 {
			t.Errorf("%s: status = %d, want 200", path, rec.Code)
		}
		if !strings.Contains(rec.Body.String(), want) {
			t.Errorf("%s: body = %s, want it to contain %q", path, rec.Body.String(), want)
		}
	}
}

func TestNextStaticFSServesRootAndAssets(t *testing.T) {
	root := adminExport(t)
	if rec := get(t, root, "/"); rec.Code != 200 || !strings.Contains(rec.Body.String(), "home") {
		t.Errorf("/: status = %d body = %s", rec.Code, rec.Body.String())
	}
	if rec := get(t, root, "/favicon.ico"); rec.Code != 200 {
		t.Errorf("/favicon.ico: status = %d, want 200", rec.Code)
	}
}

func TestNextStaticFSRejectsPathTraversal(t *testing.T) {
	root := adminExport(t)
	for _, path := range []string{"/../etc/passwd", "/admin/../../etc/passwd", "/%2e%2e/secret"} {
		rec := get(t, root, path)
		if rec.Code == 200 {
			t.Errorf("%s: status = 200, want a rejection", path)
		}
		if strings.Contains(rec.Body.String(), "root:") {
			t.Errorf("%s: served /etc/passwd", path)
		}
	}
}

func TestNextStaticFSUnknownRouteIs404(t *testing.T) {
	root := adminExport(t)
	if rec := get(t, root, "/nope"); rec.Code != 404 {
		t.Errorf("status = %d, want 404", rec.Code)
	}
	// A path that only looks like a route must not be rewritten onto a real page.
	if rec := get(t, root, "/admin.html.bak"); rec.Code == 200 {
		t.Errorf("served a page for an unexpected filename")
	}
}
