FROM golang:1.20.4-alpine3.18

RUN apk add --no-cache git && \
    go install github.com/cespare/reflex@latest

WORKDIR /app

CMD [ "reflex", "--start-service", "--", "go", "run", "/app/cmd/cli", "serve", "--web", "--with-token-cleanup" ]
