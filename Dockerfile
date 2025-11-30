# Build stage
FROM golang:1.22 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o lcp-server ./cmd/server

# Runtime stage
FROM gcr.io/distroless/base-debian12
WORKDIR /srv/lcp

COPY --from=builder /app/lcp-server /usr/local/bin/lcp-server

ENV SERVER_PORT=:8080
EXPOSE 8080

ENTRYPOINT ["/usr/local/bin/lcp-server"]
