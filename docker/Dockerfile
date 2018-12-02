FROM golang:1.9

ARG conf

WORKDIR /go/src/github.com/commojun/nyanbot
COPY . .

RUN mkdir /root/nyanbot
COPY $conf/. /root/nyanbot/.

ENV NYAN_CONF_DIR "/root/nyanbot"

RUN go install github.com/commojun/nyanbot/cmd/hello
CMD ["hello"]