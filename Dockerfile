ARG GO_VERSION=1.25
FROM golang:${GO_VERSION}-alpine AS builder

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GOPATH=/go

WORKDIR /src

COPY go.mod go.sum ./
RUN go env -w GOPROXY=https://proxy.golang.org,direct \
    && go mod download

COPY . .

RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    go build -trimpath -ldflags="-s -w" -o /app ./

FROM gcr.io/distroless/static-debian11 AS runtime

ARG APP_UID=10001
ARG APP_GID=10001

COPY --from=builder /app /app

USER ${APP_UID}:${APP_GID}

WORKDIR /

ENTRYPOINT ["/app"]
CMD ["-c", "/config.yaml"]