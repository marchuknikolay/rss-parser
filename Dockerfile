FROM golang:1.21-alpine

WORKDIR /app

COPY . .

RUN go build -o rss-parser ./cmd

ENTRYPOINT ["./rss-parser"]
CMD ["https://feeds.simplecast.com/qm_9xx0g"]