FROM golang:1.22-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o proxy .

FROM alpine:3.20
WORKDIR /app
RUN apk add --no-cache ca-certificates
COPY --from=builder /app/proxy .

ENV KAFKA_BROKERS=kafka:9092 \
    KAFKA_INPUT_TOPIC=bff-requests \
    GRADUAL_MIGRATION=false \
    MOVIES_MIGRATION_PERCENT=0

ENTRYPOINT ["./proxy"]
