# syntax=docker/dockerfile:1
FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY . .

RUN go mod download
RUN go build -o railway-go ./cmd/web/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/railway-go .
COPY config.json /app/config.json
COPY --from=builder /app/api /app/api
EXPOSE 3000

CMD ["./railway-go"]