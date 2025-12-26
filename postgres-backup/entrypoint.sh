#!/bin/bash
set -e

# Параметры из env
: "${PG_HOST:?Need to set PG_HOST}"
: "${PG_PORT:=5432}"
: "${PG_USER:?Need to set PG_USER}"
: "${PG_PASSWORD:?Need to set PG_PASSWORD}"
: "${PG_DATABASE:?Need to set PG_DATABASE}"
: "${S3_BUCKET:?Need to set S3_BUCKET}"
: "${DUMP_DIR:=/dumps}"
: "${BACKUP_INTERVAL_SECONDS:?Need to set BACKUP_INTERVAL_SECONDS}"

export PGPASSWORD="$PG_PASSWORD"

mkdir -p "$DUMP_DIR"

FILENAME="dump_${PG_DATABASE}.sql"
FILEPATH="$DUMP_DIR/$FILENAME"

echo "Starting cron job for hourly PostgreSQL dump to S3..."

while true; do
    echo "Dumping database $PG_DATABASE..."
    pg_dump -h "$PG_HOST" -p "$PG_PORT" -U "$PG_USER" "$PG_DATABASE" > "$FILEPATH"

    echo "Uploading $FILENAME to S3 bucket $S3_BUCKET..."
    aws s3 cp "$FILEPATH" "s3://$S3_BUCKET/$FILENAME"

    echo "Dump completed and uploaded. Next dump in $BACKUP_INTERVAL_SECONDS seconds."
    sleep "$BACKUP_INTERVAL_SECONDS"
done