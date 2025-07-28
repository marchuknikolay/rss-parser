FROM golang:1.23

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN chmod +x /app/start.sh

RUN go build -v -o ./bin/rss-parser ./cmd/cli \
    && go build -v -o ./bin/migrate ./cmd/migrate