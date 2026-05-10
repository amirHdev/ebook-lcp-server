# Build stage
FROM golang:1.24 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
RUN go install github.com/readium/readium-lcp-server/lcpencrypt@v1.13.4 \
	&& go install github.com/readium/readium-lcp-server/lcpserver@v1.13.4 \
	&& go install github.com/readium/readium-lcp-server/lsdserver@v1.13.4

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o lcp-server ./cmd/server

# Runtime stage
FROM docker.io/library/debian:trixie-slim
RUN apt-get update \
	&& apt-get install -y --no-install-recommends ca-certificates \
	&& rm -rf /var/lib/apt/lists/*
WORKDIR /srv/lcp

COPY --from=builder /app/lcp-server /usr/local/bin/lcp-server
COPY --from=builder /go/bin/lcpencrypt /usr/local/bin/lcpencrypt
COPY --from=builder /go/bin/lcpserver /usr/local/bin/lcpserver
COPY --from=builder /go/bin/lsdserver /usr/local/bin/lsdserver
COPY --from=builder /app/docs /srv/lcp/docs

ENV SERVER_PORT=:8080
EXPOSE 8080

ENTRYPOINT ["/usr/local/bin/lcp-server"]
