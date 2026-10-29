package fire

import "testing"

func TestParseSupportsThreeModes(t *testing.T) {
	cases := []struct {
		raw  string
		want Mode
	}{
		{"", Shadow},
		{"   ", Shadow},
		{"bogus", Shadow},
		{"active", Active},
		{"ACTIVE", Active},
		{"shadow", Shadow},
		{"SHADOW", Shadow},
		{" off ", Off},
		{"OFF", Off},
	}
	for _, tc := range cases {
		if got := Parse(tc.raw); got != tc.want {
			t.Errorf("Parse(%q) = %q, want %q", tc.raw, got, tc.want)
		}
	}
}

func TestParseConfigFailsClosed(t *testing.T) {
	cases := []struct {
		name   string
		raw    string
		ack    string
		cohort string
		allow  string
		want   Mode
	}{
		{"shadow needs no acknowledgement", "shadow", "", "", "", Shadow},
		{"off needs no acknowledgement", "off", "", "", "", Off},
		{"active needs acknowledgement", "active", "", "pilot", "s1", Shadow},
		{"active needs cohort", "active", ActivationToken, "", "s1", Shadow},
		{"active needs scope", "active", ActivationToken, "pilot", "", Shadow},
		{"active accepted", "active", ActivationToken, "pilot", "s1,s2", Active},
		{"active all accepted", "active", ActivationToken, "pilot", "all", Active},
		{"active wildcard accepted", "active", ActivationToken, "pilot", "*", Active},
		{"active mixed all rejected", "active", ActivationToken, "pilot", "all,s1", Shadow},
		{"unknown falls back", "nonsense", ActivationToken, "pilot", "s1", Shadow},
	}
	for _, tc := range cases {
		got := ParseConfig(tc.raw, tc.ack, tc.cohort, tc.allow)
		if got.Mode != tc.want {
			t.Errorf("%s: mode = %q, want %q (degraded %q)", tc.name, got.Mode, tc.want, got.Degraded)
		}
	}
}

func TestActiveForRequiresTheAllowlist(t *testing.T) {
	cfg := ParseConfig("active", ActivationToken, "pilot", " s1 , s2 ")
	if !cfg.ActiveFor("s1") || !cfg.ActiveFor("s2") {
		t.Fatal("allowlisted learners must be active")
	}
	if cfg.ActiveFor("s3") {
		t.Fatal("a learner outside the allowlist must not be active")
	}
	shadow := ParseConfig("shadow", "", "", "")
	if shadow.ActiveFor("s1") {
		t.Fatal("shadow must never be active")
	}
}

func TestAllScopeIncludesEveryLearner(t *testing.T) {
	cfg := ParseConfig("active", ActivationToken, "all-users", "all")
	if !cfg.AllowAll {
		t.Fatal("the all token must select every learner")
	}
	if !cfg.ActiveFor("s1") || !cfg.ActiveFor("a-new-account") {
		t.Fatal("the all scope must include existing and future learners")
	}
	if cfg.ActiveFor("") {
		t.Fatal("the all scope must not include an empty learner ID")
	}
	shadow := ParseConfig("shadow", "", "", "all")
	if shadow.ActiveFor("s1") {
		t.Fatal("shadow must never apply credit")
	}
}

func TestAllScopeDescribeIsExplicit(t *testing.T) {
	cfg := ParseConfig("active", ActivationToken, "all-users", "all")
	if got := cfg.Describe(); got != "active (cohort all-users, all learners)" {
		t.Fatalf("Describe() = %q", got)
	}
}

func TestOnlyShadowAndActiveCollectLedger(t *testing.T) {
	if !Shadow.CollectsLedger() {
		t.Fatal("shadow must collect ledger observations")
	}
	if !Active.CollectsLedger() {
		t.Fatal("active must collect ledger observations")
	}
	if Off.CollectsLedger() {
		t.Fatal("off must not collect ledger observations")
	}
}
