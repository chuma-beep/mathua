package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/chuma-beep/mathua/internal/storage"
)

// The issues workbench.
//
// A report is filed against production, but the person who fixes it works
// against a copy. The boundary between the two is one JSON file: `issues
// export` on the production box, `issues import` here. Nothing is written back.
// A fix still travels as a commit, because a report is evidence about a
// question, not a queue of patches, and a tool that resolved production rows
// from a laptop would be able to mark a bug fixed before the fix existed.
//
// The one state change this tool does make is local: `issues resolve` records a
// decision on the imported copy. That is how a workbench stays usable after the
// twentieth report on the same concept — the alternative is re-reading all of
// them every session.
const issuesUsage = `usage:
  mathua issues list   [--status open] [--limit 50] [--json]
  mathua issues show   <id> [--origin] [--json]
  mathua issues resolve <id> [--origin] [--status resolved] [--note "what changed"]
  mathua issues export  [--status resolved] [--out reports.json]
  mathua issues import  <file>`

func issuesCommand(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("%s", issuesUsage)
	}
	switch args[0] {
	case "list":
		return issuesList(args[1:])
	case "show":
		return issuesShow(args[1:])
	case "resolve":
		return issuesResolve(args[1:])
	case "export":
		return issuesExport(args[1:])
	case "import":
		return issuesImport(args[1:])
	default:
		return fmt.Errorf("unknown issues command %q\n%s", args[0], issuesUsage)
	}
}

// withRepo runs fn against an open store and closes it. Every subcommand needs
// exactly this, and the close is the part that gets forgotten.
func withRepo(fn func(storage.Repository) error) error {
	repo, _, err := openRepo()
	if err != nil {
		return err
	}
	defer repo.Close()
	return fn(repo)
}

// splitID separates the positional report id from flags, so `show 1 --json`
// works as naturally as `show --json 1`. Go's flag package stops at the first
// non-flag argument, which would otherwise make the documented order fail.
func splitID(args []string) (id string, rest []string) {
	for i, a := range args {
		if strings.HasPrefix(a, "-") {
			continue
		}
		return a, append(append([]string{}, args[:i]...), args[i+1:]...)
	}
	return "", args
}

func issuesList(args []string) error {
	fs := flag.NewFlagSet("issues list", flag.ContinueOnError)
	status := fs.String("status", "", "only reports in this state (open, reviewing, resolved, dismissed)")
	limit := fs.Int("limit", 50, "maximum reports to print")
	asJSON := fs.Bool("json", false, "print the reports as JSON")
	if err := fs.Parse(args); err != nil {
		return err
	}
	return withRepo(func(repo storage.Repository) error {
		reports, err := repo.ListReports(*status, *limit, 0)
		if err != nil {
			return err
		}
		if *asJSON {
			return printJSON(reports)
		}
		if len(reports) == 0 {
			fmt.Println("No reports match.")
			return nil
		}
		// Two ids, because they are different numbers. LOCAL names this copy and is what
		// show/resolve take by default; PROD is the id the admin page and a colleague use.
		fmt.Printf("%-5s %-5s %-9s %-26s %-14s %s\n", "LOCAL", "PROD", "STATUS", "CONCEPT", "KIND", "QUESTION")
		for _, r := range reports {
			fmt.Printf("%-5d %-5s %-9s %-26s %-14s %s\n",
				r.ID, orDash(r.OriginReportID), r.Status, orDash(r.ConceptID), r.Kind,
				oneLine(r.Question, 56))
		}
		fmt.Printf("\n%d report(s)\n", len(reports))
		return nil
	})
}

func issuesShow(args []string) error {
	fs := flag.NewFlagSet("issues show", flag.ContinueOnError)
	asJSON := fs.Bool("json", false, "print the report as JSON")
	byOrigin := fs.Bool("origin", false, "treat <id> as the id it had in production")
	positional, rest := splitID(args)
	if err := fs.Parse(rest); err != nil {
		return err
	}
	if positional == "" {
		return fmt.Errorf("usage: mathua issues show <id> [--json]")
	}
	id, err := strconv.ParseInt(positional, 10, 64)
	if err != nil {
		return fmt.Errorf("report id %q is not a number", positional)
	}
	return withRepo(func(repo storage.Repository) error {
		r, err := findReport(repo, id, *byOrigin)
		if err != nil {
			return err
		}
		if r == nil {
			return fmt.Errorf("no report %d", id)
		}
		if *asJSON {
			return printJSON(r)
		}
		fmt.Printf("#%d  %s  %s\n", r.ID, r.Status, r.Kind)
		if r.OriginReportID != "" {
			fmt.Printf("origin:   #%s in the database this came from\n", r.OriginReportID)
		}
		fmt.Printf("concept:  %s\n", orDash(r.ConceptID))
		fmt.Printf("lesson:   %s\n", orDash(r.LessonID))
		fmt.Printf("source:   %s\n", orDash(r.Source))
		fmt.Printf("filed:    %s\n", r.CreatedAt.Format(time.RFC3339))
		if r.Reason != "" {
			fmt.Printf("reason:   %s\n", r.Reason)
		}
		if r.Resolution != "" {
			fmt.Printf("resolved: %s (%s)\n", r.Resolution, r.ResolvedBy)
		}
		block("question", r.Question)
		block("expected", r.Expected)
		block("answered", r.AnswerGiven)
		block("explanation", r.Explanation)
		block("detail", r.Detail)
		return nil
	})
}

func issuesResolve(args []string) error {
	fs := flag.NewFlagSet("issues resolve", flag.ContinueOnError)
	status := fs.String("status", "resolved", "new state: resolved or dismissed")
	note := fs.String("note", "", "what was decided, recorded on the report")
	byOrigin := fs.Bool("origin", false, "treat <id> as the id it had in production")
	positional, rest := splitID(args)
	if err := fs.Parse(rest); err != nil {
		return err
	}
	if positional == "" {
		return fmt.Errorf("usage: mathua issues resolve <id> [--status resolved] [--note ...]")
	}
	id, err := strconv.ParseInt(positional, 10, 64)
	if err != nil {
		return fmt.Errorf("report id %q is not a number", positional)
	}
	if *status != "resolved" && *status != "dismissed" && *status != "reviewing" && *status != "open" {
		return fmt.Errorf("status %q is not one of open, reviewing, resolved, dismissed", *status)
	}
	return withRepo(func(repo storage.Repository) error {
		target := id
		if *byOrigin {
			found, err := findReport(repo, id, true)
			if err != nil {
				return err
			}
			if found == nil {
				return fmt.Errorf("no report %d", id)
			}
			target = found.ID
		}
		updated, err := repo.UpdateReportStatusAudited(adminActorCLI, target, *status, strings.TrimSpace(*note))
		if err != nil {
			return err
		}
		fmt.Printf("#%d is now %s\n", updated.ID, updated.Status)
		return nil
	})
}

// findReport resolves the id the operator typed. The default is the local row
// id; --origin switches to the id the report had in production, which is the
// number on the admin page and the one a colleague would quote.
func findReport(repo storage.Repository, id int64, byOrigin bool) (*storage.QuestionReport, error) {
	if !byOrigin {
		return repo.GetReport(id)
	}
	return repo.GetReportByOrigin(strconv.FormatInt(id, 10))
}

// issuesExportFile is the file both ends agree on. The version is checked on
// import so a future format change fails loudly instead of importing garbage.
type issuesExportFile struct {
	Version    int                      `json:"version"`
	ExportedAt time.Time                `json:"exported_at"`
	Reports    []storage.QuestionReport `json:"reports"`
}

const issuesExportVersion = 1

func issuesExport(args []string) error {
	fs := flag.NewFlagSet("issues export", flag.ContinueOnError)
	status := fs.String("status", "", "only reports in this state")
	out := fs.String("out", "", "write to this file instead of stdout")
	if err := fs.Parse(args); err != nil {
		return err
	}
	return withRepo(func(repo storage.Repository) error {
		reports, err := exportAll(repo, *status)
		if err != nil {
			return err
		}
		file := issuesExportFile{
			Version:    issuesExportVersion,
			ExportedAt: time.Now().UTC(),
			Reports:    reports,
		}
		data, err := json.MarshalIndent(file, "", "  ")
		if err != nil {
			return err
		}
		data = append(data, '\n')
		if *out == "" {
			_, err = os.Stdout.Write(data)
			return err
		}
		if err := os.WriteFile(*out, data, 0o600); err != nil {
			return err
		}
		fmt.Printf("wrote %d report(s) to %s\n", len(reports), *out)
		return nil
	})
}

// exportAll pages through the queue. ListReports caps a page at 200, so a
// single call would silently export the first page of a larger queue — the
// kind of truncation that only shows up when someone needs the missing rows.
func exportAll(repo storage.Repository, status string) ([]storage.QuestionReport, error) {
	const page = 200
	all := []storage.QuestionReport{}
	for offset := 0; ; offset += page {
		batch, err := repo.ListReports(status, page, offset)
		if err != nil {
			return nil, err
		}
		all = append(all, batch...)
		if len(batch) < page {
			return all, nil
		}
	}
}

func issuesImport(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("usage: mathua issues import <file>")
	}
	data, err := os.ReadFile(args[0])
	if err != nil {
		return err
	}
	var file issuesExportFile
	if err := json.Unmarshal(data, &file); err != nil {
		return fmt.Errorf("read %s: %w", args[0], err)
	}
	if file.Version != issuesExportVersion {
		return fmt.Errorf("%s is export version %d, this binary understands %d", args[0], file.Version, issuesExportVersion)
	}
	return withRepo(func(repo storage.Repository) error {
		inserted, skipped := 0, 0
		for _, r := range file.Reports {
			origin := strconv.FormatInt(r.ID, 10)
			ok, err := repo.ImportReport(origin, r)
			if err != nil {
				return fmt.Errorf("report %s: %w", origin, err)
			}
			if ok {
				inserted++
			} else {
				skipped++
			}
		}
		fmt.Printf("imported %d report(s), %d already present\n", inserted, skipped)
		return nil
	})
}

func printJSON(v interface{}) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	return nil
}

func block(label, text string) {
	if strings.TrimSpace(text) == "" {
		return
	}
	fmt.Printf("\n%s:\n%s\n", label, text)
}

func orDash(s string) string {
	if strings.TrimSpace(s) == "" {
		return "-"
	}
	return s
}

// oneLine keeps a table row to one line and one terminal width.
func oneLine(s string, max int) string {
	s = strings.Join(strings.Fields(s), " ")
	if len(s) <= max {
		return s
	}
	if max <= 1 {
		return s[:max]
	}
	return s[:max-1] + "…"
}
