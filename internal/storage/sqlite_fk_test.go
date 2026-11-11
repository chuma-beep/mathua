package storage

import (
	"context"
	"path/filepath"
	"testing"
)

// A file-backed store pools eight connections. foreign_keys is a per-connection
// pragma, so a one-shot Exec only reaches whichever connection it lands on and
// ON DELETE CASCADE silently stops working on the rest. Holding all eight
// connections at once proves the DSN reached every one of them.
func TestSQLiteFileStoreEnablesForeignKeysOnEveryConnection(t *testing.T) {
	store, err := NewSQLiteStore(filepath.Join(t.TempDir(), "fk.db"))
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	t.Cleanup(func() { store.Close() })

	ctx := context.Background()
	conns := make([]interface{ Close() error }, 0, 8)
	defer func() {
		for _, c := range conns {
			c.Close()
		}
	}()
	for i := 0; i < 8; i++ {
		conn, err := store.db.Conn(ctx)
		if err != nil {
			t.Fatalf("take connection %d: %v", i, err)
		}
		conns = append(conns, conn)
		var foreignKeys, busyTimeout int
		if err := conn.QueryRowContext(ctx, "PRAGMA foreign_keys").Scan(&foreignKeys); err != nil {
			t.Fatalf("read foreign_keys on connection %d: %v", i, err)
		}
		if err := conn.QueryRowContext(ctx, "PRAGMA busy_timeout").Scan(&busyTimeout); err != nil {
			t.Fatalf("read busy_timeout on connection %d: %v", i, err)
		}
		if foreignKeys != 1 {
			t.Fatalf("connection %d has foreign keys off", i)
		}
		if busyTimeout != 5000 {
			t.Fatalf("connection %d busy timeout = %d, want 5000", i, busyTimeout)
		}
	}
}
