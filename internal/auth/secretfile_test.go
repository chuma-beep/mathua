package auth

import (
	"encoding/hex"
	"os"
	"path/filepath"
	"testing"
)

func TestLoadSecretFile_RoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), ".jwt_secret")
	first, ok := loadSecretFile(path)
	if !ok || len(first) != 32 {
		t.Fatalf("expected fresh 32-byte secret, ok=%v len=%d", ok, len(first))
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat: %v", err)
	}
	if info.Mode().Perm() != 0o600 {
		t.Errorf("expected 0600, got %o", info.Mode().Perm())
	}
	second, ok := loadSecretFile(path)
	if !ok {
		t.Fatal("expected reload to succeed")
	}
	if string(first) != string(second) {
		t.Error("expected persisted secret to be stable across loads (restart survival)")
	}
}

func TestLoadSecretFile_CreatesDirs(t *testing.T) {
	path := filepath.Join(t.TempDir(), "sub", "dir", ".jwt_secret")
	if _, ok := loadSecretFile(path); !ok {
		t.Fatal("expected nested dirs to be created")
	}
}

func TestLoadSecretFile_CorruptQuarantined(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".jwt_secret")
	if err := os.WriteFile(path, []byte("not-a-secret"), 0o600); err != nil {
		t.Fatal(err)
	}
	fresh, ok := loadSecretFile(path)
	if !ok || len(fresh) != 32 {
		t.Fatalf("expected fresh secret after corruption, ok=%v", ok)
	}
	entries, _ := os.ReadDir(dir)
	quarantined := false
	for _, e := range entries {
		if len(e.Name()) > len(".jwt_secret.bad.") && e.Name()[:len(".jwt_secret.bad.")] == ".jwt_secret.bad." {
			quarantined = true
		}
	}
	if !quarantined {
		t.Error("expected corrupt file quarantined aside")
	}
	raw, _ := os.ReadFile(path)
	if _, err := hex.DecodeString(string(raw[:64])); err != nil {
		t.Errorf("expected valid hex secret rewritten, got %q", raw)
	}
}
