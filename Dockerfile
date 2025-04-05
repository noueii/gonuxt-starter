#1. Build stage
FROM golang:1.24.2-alpine3.21 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go
#2. Run stage
FROM alpine:3.21
WORKDIR /app
RUN apk add --no-cache curl
RUN curl -fsSL \
	https://raw.githubusercontent.com/pressly/goose/master/install.sh |\
	sh


COPY --from=builder /app/main .
COPY prod.env .
COPY test.env .
COPY scripts/start.sh .
COPY db/schema ./schema

EXPOSE 8080
CMD [ "/app/main" ]
ENTRYPOINT [ "/app/start.sh" ]

