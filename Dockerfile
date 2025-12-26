FROM golang:1.25-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git ca-certificates binutils

COPY go.mod go.sum ./
RUN go mod download

RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.17.1

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o app ./cmd/api && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o seed ./cmd/seed && \
    strip /app/app /app/seed


FROM alpine:3.20

WORKDIR /app

RUN apk add --no-cache ca-certificates curl

COPY --from=builder /go/bin/migrate /usr/local/bin/migrate
COPY --from=builder /app/app /app/app
COPY --from=builder /app/internal/platform/database/migrations /app/migrations
COPY --from=builder /app/seed /app/seed

RUN addgroup -S appgroup && adduser -S appuser -G appgroup && chown -R appuser:appgroup /app

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=3s --start-period=10s CMD ["sh","-c","curl -f http://localhost:8080/health || exit 1"]

USER appuser

ENTRYPOINT ["/app/app"]
CMD []
