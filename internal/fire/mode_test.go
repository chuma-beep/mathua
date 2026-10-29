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
		{"active needs allowlist", "active", ActivationToken, "pilot", "", Shadow},
		{"active accepted", "active", ActivationToken, "pilot", "s1,s2", Active},
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
