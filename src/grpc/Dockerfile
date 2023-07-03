FROM golang:1.20

WORKDIR /app

#COPY go.mod ./
#COPY go.sum ./
COPY . ./
WORKDIR /app/src
RUN go get github.com/redis/go-redis/v9
RUN go mod download

RUN go build -o server server.go
EXPOSE 50051
CMD ["./server"]