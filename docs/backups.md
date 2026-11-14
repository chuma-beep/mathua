# Database backups

Status: production procedure.

The Fly app keeps SQLite on the `mathua_data` volume. Litestream copies the
write-ahead log to Cloudflare R2 every ten seconds and keeps thirty days of
history. The image contains Litestream. The entrypoint starts it only when
`LITESTREAM_BUCKET` is set, so a deployment without backups behaves as before.

## Create the bucket and token

1. Open Cloudflare and go to R2. Create a private bucket named `mathua-backups`.
2. Create an R2 API token with Object Read and Write access, limited to that
   bucket.
3. Note the account ID. The endpoint is
   `https://<account-id>.r2.cloudflarestorage.com`.

## Set the secrets

Run these commands on the machine that holds the Fly credentials. Replace the
placeholders with the bucket, endpoint and token values.

```bash
flyctl secrets set \
  LITESTREAM_BUCKET=mathua-backups \
  LITESTREAM_ENDPOINT=https://<account-id>.r2.cloudflarestorage.com \
  LITESTREAM_ACCESS_KEY_ID=<access-key> \
  LITESTREAM_SECRET_ACCESS_KEY=<secret-key> \
  -a mathua
```

The machine restarts. The entrypoint restores the newest generation when the
database file is missing. It then starts the API under `litestream replicate`.

## Verify a backup exists

```bash
flyctl ssh console -a mathua -C \
  "litestream snapshots -config /etc/litestream.yml /data/mathua.db"
```

The command lists generations with timestamps and sizes. A new deployment can
take up to one snapshot interval to show the first entry.

Check the replica directly:

```bash
flyctl ssh console -a mathua -C \
  "litestream generations -config /etc/litestream.yml /data/mathua.db"
```

## Restore drill

Run the drill once before launch and once a quarter. It restores to `/tmp` so
production keeps serving.

```bash
flyctl ssh console -a mathua -C \
  "litestream restore -config /etc/litestream.yml -o /tmp/mathua-restore.db /data/mathua.db"
flyctl ssh console -a mathua -C \
  "python3 -c \"import sqlite3; db=sqlite3.connect('/tmp/mathua-restore.db'); print(db.execute('pragma integrity_check').fetchone()); print(db.execute('select count(*) from students').fetchone())\""
```

The integrity check must return `ok`. The student count must be plausible for
the restore timestamp.

## Recover a lost volume

1. Stop the machine or scale it to zero.
2. Move the damaged database aside.
3. Restore to the live path.
4. Start the machine.
5. Check `/api/health` and the admin overview.

```bash
flyctl ssh console -a mathua -C \
  "mv /data/mathua.db /data/mathua.db.damaged"
flyctl ssh console -a mathua -C \
  "litestream restore -config /etc/litestream.yml -o /data/mathua.db /data/mathua.db"
```

The entrypoint performs the same restore automatically when the volume is
empty, so a replacement machine recovers without manual steps.

## Failure modes

- A missing endpoint or token stops the machine at boot. This is deliberate:
  a backup configuration that cannot work must not look healthy.
- An empty bucket on the first boot logs a warning and starts with a new
  database.
- Two writers must never share one database. Run one Fly machine.
- Litestream keeps a lock on the database. Do not run a second replication
  process against the same file.
- Retention is thirty days. A restore older than that is not possible.

## Cost

The database is about 743 MB. R2 includes 10 GB of storage, one million Class A
operations and ten million Class B operations each month. WAL frames are small
and frequent, so Class A use stays low. The backup should remain inside the free
tier at the current size.

## Rollback

Unset the Litestream secrets to stop replication. The entrypoint then starts the
API directly. The R2 bucket keeps the existing generations until you delete
them.
