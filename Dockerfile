FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o /barista ./cmd/barista

FROM alpine:3.21

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /barista /barista
COPY config/ config/

RUN mkdir -p output/expectations

CMD ["/barista"]
