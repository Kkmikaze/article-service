# Stage 1: Build the Go application
FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY . .

ARG SERVICE_VERSION

ENV SERVICE_VERSION=${SERVICE_VERSION} \
    TZ=Asia/Jakarta

RUN apk add --no-cache tzdata \
 && ln -sf /usr/share/zoneinfo/Asia/Jakarta /etc/localtime \
 && echo "Asia/Jakarta" > /etc/timezone \
 && go mod tidy \
 && CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w -X main.version=${SERVICE_VERSION}" -o engine ./cmd/server

# Stage 2: Minimal final image
FROM scratch

COPY --from=builder /app/.env /.env
COPY --from=builder /app/engine /engine
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/localtime /etc/localtime
COPY --from=builder /etc/timezone /etc/timezone

ENV TZ=Asia/Jakarta

ENTRYPOINT ["/engine"]
CMD ["engine", "--rpc", "8080", "--gateway", "80"]

