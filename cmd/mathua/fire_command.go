package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/chuma-beep/mathua/internal/concepts"
	"github.com/chuma-beep/mathua/internal/engine"
	"github.com/chuma-beep/mathua/internal/fire"
	"github.com/chuma-beep/mathua/internal/generator/all"
	"github.com/chuma-beep/mathua/internal/storage"
)

// resolveFIReConfig reads the process configuration and then requires the
// durable activation row to agree with it. A missing or mismatched row degrades
// active to shadow. This is the second key: an environment variable alone can
// never switch scheduling effects on.
func resolveFIReConfig(repo storage.Repository) fire.Config {
	cfg := fire.ParseConfig(
		os.Getenv("MATHUA_FIRE_MODE"),
		os.Getenv("MATHUA_FIRE_ACK"),
		os.Getenv("MATHUA_FIRE_COHORT"),
		os.Getenv("MATHUA_FIRE_ALLOWLIST"),
	)
	if cfg.Mode != fire.Active {
		return cfg
	}
	store, ok := repo.(storage.FireActivationStore)
	if !ok {
		return cfg.Degrade("storage does not support durable FIRe activation")
	}
	act, err := store.GetFireActivation()
	if err != nil {
		return cfg.Degrade("read durable activation: " + err.Error())
	}
	if act == nil {
		return cfg.Degrade("durable activation row is missing")
	}
	if act.Mode != string(fire.Active) {
		return cfg.Degrade("durable activation is not active")
	}
	if strings.TrimSpace(act.Cohort) != cfg.Cohort {
		return cfg.Degrade("durable activation cohort does not match the deployed cohort")
	}
	return cfg
}

func fireCommand(args []string) error {
	if len(args) == 0 {
		return errors.New("usage: mathua fire status|report|activate|deactivate")
	}
	repo, _, err := openRepo()
	if err != nil {
		return err
	}
	defer repo.Close()

	switch args[0] {
	case "status":
		cfg := resolveFIReConfig(repo)
		store, ok := repo.(storage.FireActivationStore)
		var act *storage.FireActivation
		if ok {
			act, err = store.GetFireActivation()
			if err != nil {
				return err
			}
		}
		out := map[string]any{
			"configured_mode": cfg.Mode,
			"cohort":          cfg.Cohort,
			"allowlist_size":  len(cfg.Allow),
			"allow_all":       cfg.AllowAll,
			"degraded":        cfg.Degraded,
			"description":     cfg.Describe(),
			"activation":      act,
		}
		return writeJSONStdout(out)
	case "report":
		dag, err := concepts.LoadDir(dataPath("concepts"))
		if err != nil {
			return err
		}
		cfg := resolveFIReConfig(repo)
		eng := engine.New(repo, dag, all.Registry(), nil, nil, engine.WithFIReConfig(cfg))
		defer eng.Close()
		report, err := eng.FIReDryRun(time.Now().UTC())
		if err != nil {
			return err
		}
		return writeJSONStdout(report)
	case "activate":
		return activateFIRe(repo, args[1:])
	case "deactivate":
		return deactivateFIRe(repo, args[1:])
	default:
		return fmt.Errorf("unknown fire command %q", args[0])
	}
}

func activateFIRe(repo storage.Repository, args []string) error {
	opts, err := parseFireFlags(args)
	if err != nil {
		return err
	}
	if opts.ack != fire.ActivationToken {
		return fmt.Errorf("activation requires --ack %s", fire.ActivationToken)
	}
	if strings.TrimSpace(opts.cohort) == "" {
		return errors.New("activation requires --cohort")
	}
	if strings.TrimSpace(opts.allowlist) == "" {
		return errors.New("activation requires --allowlist with at least one student ID or all")
	}
	// Validate the full configuration before recording the row, so a typo is
	// rejected here instead of silently degrading at the next startup.
	if cfg := fire.ParseConfig("active", opts.ack, opts.cohort, opts.allowlist); cfg.Mode != fire.Active {
		return fmt.Errorf("activation rejected: %s", cfg.Degraded)
	}
	store, ok := repo.(storage.FireActivationStore)
	if !ok {
		return errors.New("storage does not support durable FIRe activation")
	}
	now := time.Now().UTC()
	act, err := store.SetFireActivation(storage.FireActivation{
		Mode:           string(fire.Active),
		Cohort:         strings.TrimSpace(opts.cohort),
		AcknowledgedBy: opts.actor,
		AcknowledgedAt: now,
		UpdatedAt:      now,
	})
	if err != nil {
		return err
	}
	return writeJSONStdout(act)
}

func deactivateFIRe(repo storage.Repository, args []string) error {
	opts, err := parseFireFlags(args)
	if err != nil {
		return err
	}
	store, ok := repo.(storage.FireActivationStore)
	if !ok {
		return errors.New("storage does not support durable FIRe activation")
	}
	existing, err := store.GetFireActivation()
	if err != nil {
		return err
	}
	act, err := store.SetFireActivation(storage.FireActivation{
		Mode:           string(fire.Shadow),
		Cohort:         opts.cohort,
		AcknowledgedBy: opts.actor,
		UpdatedAt:      time.Now().UTC(),
	})
	if err != nil {
		return err
	}
	if existing != nil {
		fmt.Fprintln(os.Stderr, "FIRe disabled; existing application rows remain auditable and are not reversed automatically")
	}
	return writeJSONStdout(act)
}

type fireFlags struct {
	ack       string
	cohort    string
	allowlist string
	actor     string
}

func parseFireFlags(args []string) (fireFlags, error) {
	out := fireFlags{actor: "cli"}
	for i := 0; i < len(args); i++ {
		key := args[i]
		value := ""
		if eq := strings.IndexByte(key, '='); eq >= 0 {
			value = key[eq+1:]
			key = key[:eq]
		} else {
			i++
			if i >= len(args) {
				return out, fmt.Errorf("%s needs a value", key)
			}
			value = args[i]
		}
		switch key {
		case "--ack":
			out.ack = value
		case "--cohort":
			out.cohort = value
		case "--allowlist":
			out.allowlist = value
		case "--actor":
			out.actor = value
		default:
			return out, fmt.Errorf("unknown fire flag %q", key)
		}
	}
	return out, nil
}

func writeJSONStdout(v any) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}
