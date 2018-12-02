FROM golang:alpine
MAINTAINER commojun <junpooooow@gmail.com>

ENV NYAN_USER_DIR "/root/nyanbot"

WORKDIR /go/src/github.com/commojun/nyanbot
COPY . .

RUN make deps
RUN go install github.com/commojun/nyanbot/cmd/hello
RUN go install github.com/commojun/nyanbot/cmd/alarm

# crontab
COPY
