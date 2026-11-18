package lessons

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestStripTrailingReferences(t *testing.T) {
	trailing := "# Differentiation\n\nA derivative is a limit.\n\n## Selected references\n\n  * **MIT OpenCourseWare**. [Lecture notes](https://ocw.mit.edu/x)\n  * **OpenStax**. [Calculus](https://openstax.org/y)\n"
	got := stripTrailingReferences(trailing)
	if strings.Contains(strings.ToLower(got), "selected references") {
		t.Fatalf("heading survived:\n%s", got)
	}
	if !strings.Contains(got, "A derivative is a limit.") {
		t.Fatalf("prose was removed with the list:\n%s", got)
	}
	if strings.HasSuffix(got, "\n\n") {
		t.Fatalf("trailing blank lines left behind: %q", got)
	}

	cases := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "no section is unchanged",
			in:   "# Limits\n\nThe limit is the value a function approaches.\n",
			want: "# Limits\n\nThe limit is the value a function approaches.\n",
		},
		{
			name: "a later section means the list is not trailing",
			in:   "# A\n\nProse.\n\n## Selected references\n\n  * [x](https://x)\n\n## Worked example\n\nMore prose.\n",
			want: "# A\n\nProse.\n\n## Selected references\n\n  * [x](https://x)\n\n## Worked example\n\nMore prose.\n",
		},
		{
			name: "a heading inside a fence is not a section",
			in:   "# A\n\n```text\n## Selected references\n```\n\nStill prose.\n",
			want: "# A\n\n```text\n## Selected references\n```\n\nStill prose.\n",
		},
		{
			name: "case and trailing spaces still match",
			in:   "# A\n\nProse.\n\n## selected references   \n\n  * [x](https://x)\n",
			want: "# A\n\nProse.\n",
		},
	}

	for _, tc := range cases {
		if got := stripTrailingReferences(tc.in); got != tc.want {
			t.Errorf("%s:\n got %q\nwant %q", tc.name, got, tc.want)
		}
	}
}

// The corpus keeps the list on disk, because the files are a vendored copy and
// the scraper would re-add it anyway. What must not survive is the served body.
func TestCurriculumBodiesDropReferenceList(t *testing.T) {
	files := 0
	paths, err := filepath.Glob("../../data/lessons/algebrica/*/*.md")
	if err != nil {
		t.Fatal(err)
	}
	for _, path := range paths {
		body, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("read %s: %v", path, err)
		}
		if strings.Contains(strings.ToLower(string(body)), "## selected references") {
			files++
		}
	}
	if files == 0 {
		t.Fatal("no vendored file carries the section, so this sweep checks nothing")
	}

	loader, err := Load("../../data/lessons")
	if err != nil {
		t.Fatalf("load corpus: %v", err)
	}
	for source, lesson := range loader.All() {
		if strings.Contains(strings.ToLower(lesson.Body), "## selected references") {
			t.Errorf("%s: served body still contains the reference list", source)
		}
		if strings.TrimSpace(lesson.Body) == "" {
			t.Errorf("%s: served body is empty after the cut", source)
		}
	}

	// One concept whose file carries the section, read the way `/study` reads it.
	lesson := loader.Lesson("stat.dist.normal")
	if lesson == nil {
		t.Fatal("stat.dist.normal is missing from the corpus")
	}
	if strings.TrimSpace(lesson.Body) == "" {
		t.Fatal("stat.dist.normal has an empty body")
	}
}
