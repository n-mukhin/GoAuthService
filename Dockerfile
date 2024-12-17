FROM golang:1.20-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN go build -o server cmd/server/main.go

EXPOSE 8080

CMD ["./server"]
