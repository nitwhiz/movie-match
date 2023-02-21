FROM golang:1.20-alpine3.17 AS builder

WORKDIR /build

COPY ./ /build

RUN go build -o dist/movie-match ./cmd/cli

FROM alpine:3.17

WORKDIR /opt/movie-match

COPY --from=builder /build/dist/movie-match /opt/movie-match/movie-match

RUN mkdir /opt/movie-match/posters

VOLUME /opt/movie-match/posters

ENTRYPOINT [ "/opt/movie-match/movie-match" ]

CMD [ "serve", "--web", "--with-token-cleanup", "--with-media-auto-pull" ]
