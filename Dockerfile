FROM golang:1.13.8

WORKDIR /go/src/github.com/VolticFroogo/QShrtn
COPY . .
RUN go get -d -v ./...
RUN go install -v ./...
RUN go build -o main .

CMD ["./main"]
