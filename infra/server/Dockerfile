# ---- builder ----
ARG GO_VERSION=1.22
FROM golang:${GO_VERSION}-alpine AS builder
WORKDIR /app
RUN apk add --no-cache git ca-certificates && update-ca-certificates

# copy module files first (cache-friendly)
COPY src/go.mod ./
COPY src/go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download

# copy the rest of the source
COPY src/ ./
RUN go mod tidy

# build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o /app/api ./cmd/api

# ---- runner ----
FROM alpine:3.20
RUN apk add --no-cache ca-certificates tzdata curl && update-ca-certificates
WORKDIR /app
COPY --from=builder /app/api /app/api
EXPOSE 8080
ENTRYPOINT ["/app/api"]
