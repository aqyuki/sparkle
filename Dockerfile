#==================== Builder ====================
FROM golang:1.22.5-alpine3.19 AS builder

WORKDIR /app
COPY . .
RUN go mod download
RUN go mod verify
RUN go build -ldflags="-s -w" -trimpath -o sparkle cmd/sparkle/main.go

#==================== Runner ====================
FROM gcr.io/distroless/cc-debian12 AS runner

WORKDIR /app
COPY --from=builder --chown=root:root /app/sparkle /app/sparkle
STOPSIGNAL SIGINT
ENTRYPOINT ["./sparkle"]
