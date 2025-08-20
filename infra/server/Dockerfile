# ---- builder ----
ARG GO_VERSION=1.22
FROM golang:${GO_VERSION}-alpine AS builder
WORKDIR /app
RUN apk add --no-cache git ca-certificates && update-ca-certificates

# module files
COPY src/go.mod ./
COPY src/go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download

# copy source
COPY src/ ./
RUN go mod tidy

# install swag CLI
RUN GOBIN=/usr/local/bin go install github.com/swaggo/swag/cmd/swag@v1.16.3

# generate swagger docs
RUN swag init -g cmd/api/main.go -o docs

# build binaries
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o /app/api ./cmd/api
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o /app/migrate ./cmd/migrate

# ---- runner ----
FROM alpine:3.20
RUN apk add --no-cache ca-certificates tzdata curl mysql-client && update-ca-certificates
WORKDIR /app
COPY --from=builder /app/api /app/api
COPY --from=builder /app/migrate /app/migrate
COPY --from=builder /app/docs /app/docs
COPY src/.env /app/.env
EXPOSE 8080
ENTRYPOINT ["/app/api"]
