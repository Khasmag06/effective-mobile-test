# Build stage
FROM golang:1.21.0-alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/app ./cmd/app/main.go

# Final stage
FROM alpine:latest

WORKDIR /effective-mobile-test

COPY --from=builder /build/migrations ./migrations
COPY --from=builder /build/bin/app .
COPY --from=builder /build/.env .

CMD ["./app"]
