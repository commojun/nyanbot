FROM golang:1.9
WORKDIR /go/src/github.com/commojun/nyanbot
COPY . .
RUN go install github.com/commojun/nyanbot/cmd/hello
CMD ["hello"]