#!/bin/sh
# Start the API under Litestream when an R2 replica is configured.
#
# Litestream is opt-in. Without LITESTREAM_BUCKET this script is a plain exec,
# so local Docker and Fly deployments without backups behave exactly as before.
# With it, the database is restored only when the volume is empty, then the
# process runs under `litestream replicate`, which streams WAL frames to R2.
set -eu

if [ -z "${LITESTREAM_BUCKET:-}" ]; then
  exec "$@"
fi

# Postgres deployments do not use Litestream. The replica is SQLite only.
case "${DATABASE_URL:-/data/mathua.db}" in
  postgres://*|postgresql://*)
    exec "$@"
    ;;
esac

: "${LITESTREAM_ENDPOINT:?LITESTREAM_ENDPOINT is required when LITESTREAM_BUCKET is set}"
: "${LITESTREAM_ACCESS_KEY_ID:?LITESTREAM_ACCESS_KEY_ID is required when LITESTREAM_BUCKET is set}"
: "${LITESTREAM_SECRET_ACCESS_KEY:?LITESTREAM_SECRET_ACCESS_KEY is required when LITESTREAM_BUCKET is set}"

db="${DATABASE_URL:-/data/mathua.db}"
case "$db" in
  sqlite:*) db="${db#sqlite:}" ;;
esac

# A fresh volume has no database. Restore the newest generation first. An empty
# bucket is normal on the first boot, so a failed restore is a warning and the
# app starts with a new database that Litestream will seed.
if [ ! -f "$db" ]; then
  if ! litestream restore -if-db-not-exists -config /etc/litestream.yml "$db"; then
    echo "litestream: no restore generation, starting with an empty database" >&2
  fi
fi

exec litestream replicate -config /etc/litestream.yml -exec "$*"
