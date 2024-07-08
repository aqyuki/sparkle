#==================== Builder ====================
FROM golang:1.22.5-alpine3.19 AS builder

ARG BOT_VERSION=unknown
WORKDIR /app
COPY . .
RUN go mod download
RUN go mod verify
RUN go build -ldflags="-s -w -X 'github.com/aqyuki/sparkle/internal/info.Version=${BOT_VERSION}'" -trimpath -o sparkle cmd/sparkle/main.go

#==================== Runner ====================
FROM gcr.io/distroless/cc-debian12 AS runner

WORKDIR /app
COPY --from=builder --chown=root:root /app/sparkle /app/sparkle
STOPSIGNAL SIGINT
ENTRYPOINT ["./sparkle"]
