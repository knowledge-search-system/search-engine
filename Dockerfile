FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux \
    go build \
    -ldflags="-s -w" \
    -o search-engine \
    ./cmd/main.go

FROM alpine:3.22

WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/search-engine .

COPY --from=builder /app/config ./config

EXPOSE 8080

EXPOSE 9090

ENTRYPOINT ["./search-engine"]