FROM golang:1.21-alpine

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

RUN go build -v -o ./bin/rss-parser ./cmd/cli

ENTRYPOINT ["./bin/rss-parser"]
CMD ["https://feeds.simplecast.com/qm_9xx0g"]