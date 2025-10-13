# =========================
# 1) Builder
# =========================
FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /app

ARG VERSION=dev
ARG COMMIT=none
ARG BUILD_TIME

# Cache de dependencias
COPY go.mod go.sum ./
RUN go mod download

# Código
COPY . .

# Build binario estático
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -ldflags="-s -w \
      -X 'main.version=${VERSION}' \
      -X 'main.commit=${COMMIT}' \
      -X 'main.buildTime=${BUILD_TIME}'" \
    -o /app/bin/server ./cmd/server

# =========================
# 2) Runtime
# =========================
FROM alpine:3.20

RUN addgroup -S app && adduser -S app -G app \
  && apk add --no-cache ca-certificates tzdata

ENV GIN_MODE=release \
    PORT=8080 \
    TZ=America/Santiago

COPY --from=builder /app/bin/server /usr/local/bin/server

RUN mkdir -p /app && chown -R app:app /app
WORKDIR /app

EXPOSE 8080
USER app

CMD ["/usr/local/bin/server"]
