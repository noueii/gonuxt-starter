#1. Build stage
FROM golang:1.24.2-alpine3.21 AS builder
WORKDIR /app
COPY . .
RUN go build -o gonuxt ./cmd/gonuxt-starter
#2. Run stage
FROM alpine:3.21
WORKDIR /app
#RUN apk add --no-cache curl
#RUN curl -fsSL \
#	https://raw.githubusercontent.com/pressly/goose/master/install.sh |\
#	sh


COPY --from=builder /app/gonuxt .
COPY prod.env .
COPY test.env .
COPY app.env .
COPY dev.env .
#COPY scripts/start.sh .
#COPY db/schema ./db/schema

EXPOSE 8080 7777
CMD [ "/app/gonuxt" ]
#ENTRYPOINT [ "/app/start.sh" ]

