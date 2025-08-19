# ---- dev ----
ARG GO_VERSION=1.22
FROM golang:${GO_VERSION}-alpine AS dev

WORKDIR /app

RUN apk add --no-cache git ca-certificates tzdata curl build-base && update-ca-certificates

RUN go install github.com/air-verse/air@latest

COPY src/go.mod ./ 
COPY src/go.sum ./ 
RUN go mod download

COPY src/ ./ 

CMD ["air"]
