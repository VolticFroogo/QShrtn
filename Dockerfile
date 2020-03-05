FROM golang:1.14.0

WORKDIR /go/src/github.com/VolticFroogo/QShrtn
COPY . .
RUN go build -o main .

CMD ["./main"]
