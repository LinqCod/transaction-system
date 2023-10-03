FROM golang:1.21-alpine AS builder
WORKDIR /app

COPY ./publisher_service/go.mod ./publisher_service/go.sum ./
RUN go mod download

COPY .env .
COPY ./publisher_service .

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