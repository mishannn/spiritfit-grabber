#!/bin/bash
set -e

# Env params
: "${S3_BUCKET:?Need to set S3_BUCKET}"
: "${BACKUP_INTERVAL_SECONDS:?Need to set BACKUP_INTERVAL_SECONDS}"

echo "Starting cron job for hourly VM dump to S3..."

while true; do
    /vmbackup-prod \
        -storageDataPath=/victoria-metrics-data \
        -snapshot.createURL=http://victoriametrics:8428/snapshot/create \
        -customS3Endpoint=${AWS_ENDPOINT_URL} \
        -dst=s3://${S3_BUCKET}/vm-backup/

    echo "Dump completed and uploaded. Next dump in $BACKUP_INTERVAL_SECONDS seconds."
    sleep "$BACKUP_INTERVAL_SECONDS"
done