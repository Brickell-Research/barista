FROM golang:1.23-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /explore ./cmd/explore

FROM alpine:3.21
RUN apk add --no-cache git ca-certificates
COPY --from=build /explore /usr/local/bin/explore
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh
WORKDIR /data
COPY config/services.yml config/services.yml
ENTRYPOINT ["/entrypoint.sh"]
