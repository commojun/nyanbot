FROM golang:1.9

ENV NYAN_USER_DIR "/root/nyanbot"

WORKDIR /go/src/github.com/commojun/nyanbot
COPY . .

RUN make deps
RUN go install github.com/commojun/nyanbot/cmd/hello
