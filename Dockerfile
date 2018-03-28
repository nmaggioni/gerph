FROM golang:1

WORKDIR /app

COPY . .
RUN go get -d -v ./... && mkdir dist && go build -o dist/gerph

VOLUME ["/db"]

ENV PORT 3000

CMD ["sh", "-c", "/app/dist/gerph -dbpath /db/gerph.db -port ${PORT}"]