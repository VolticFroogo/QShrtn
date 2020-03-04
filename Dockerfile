FROM golang:1.13.8

WORKDIR /go/src/github.com/VolticFroogo/QShrtn
COPY . .
RUN go build -o main .

CMD ["./main"]
