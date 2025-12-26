FROM alpine:3.20

# Устанавливаем зависимости
RUN apk add --no-cache \
    postgresql-client \
    python3 \
    py3-pip \
    bash \
    curl \
    jq \
    tzdata \
    busybox-suid \
    && pip install --no-cache-dir --break-system-packages awscli \
    && apk add --no-cache bash \
    && mkdir -p /dumps /scripts

# Копируем скрипт entrypoint
COPY entrypoint.sh /scripts/entrypoint.sh
RUN chmod +x /scripts/entrypoint.sh

# Используем bash в качестве entrypoint
ENTRYPOINT ["/scripts/entrypoint.sh"]