FROM golang:1.21-alpine

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

RUN go build -v -o rss-parser ./cmd

ENTRYPOINT ["./rss-parser"]
CMD ["https://feeds.simplecast.com/qm_9xx0g"]