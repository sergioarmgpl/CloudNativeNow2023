FROM golang:1.20

WORKDIR /app

#COPY go.mod ./
#COPY go.sum ./
COPY . ./
WORKDIR /app/src

RUN go mod download

RUN go build -o consumer consumer.go
#EXPOSE 3000
CMD ["./consumer"]