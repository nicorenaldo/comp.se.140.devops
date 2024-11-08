FROM golang:1.23

WORKDIR /app

COPY service1-go/ ./

RUN go mod download
RUN go build -o main .

CMD ["./main"]