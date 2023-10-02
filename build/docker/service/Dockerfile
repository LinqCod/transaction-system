FROM golang:1.21-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app ./cmd/main.go

FROM alpine:latest

RUN apk update && apk upgrade

RUN rm -rf /var/cache/apk/* && \
    rm -rf /tmp/*

WORKDIR /app

COPY --from=builder /app/app .
COPY --from=builder /app/.env .

EXPOSE 8080

CMD ["./app"]