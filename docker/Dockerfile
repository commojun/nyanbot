FROM golang:alpine
MAINTAINER commojun <junpooooow@gmail.com>

RUN apk add make gcc g++ git

ENV NYAN_USER_DIR "/root/nyanbot"
WORKDIR /go/src/github.com/commojun/nyanbot
COPY . .

# install
RUN make deps
RUN go install github.com/commojun/nyanbot/cmd/hello
RUN go install github.com/commojun/nyanbot/cmd/alarm

# crontab
COPY crontab/crontab.txt /etc/crontabs/root
CMD ["crond", "-f"]
