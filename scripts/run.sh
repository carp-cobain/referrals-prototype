#!/bin/bash
set -e

# Restore the database if it doesn't already exist.
if [ -f "${DB_PATH}" ]; then
	echo "Database already exists, skipping restore"
else
	echo "No database found, restoring from replica if exists"
	litestream restore -v -if-replica-exists -o "${DB_PATH}" "${REPLICA_URL}"
fi

# Run litestream with the app as the subprocess.
exec litestream replicate -exec "/usr/local/bin/referrals"
