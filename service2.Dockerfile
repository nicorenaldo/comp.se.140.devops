FROM golang:1.23

WORKDIR /app

COPY service2-go/ ./

RUN go mod download
RUN go build -o main .

CMD ["./main"]