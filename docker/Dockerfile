FROM debian:bullseye AS builder

ARG TARGETPLATFORM
ARG VERSION=0.0.2

RUN apt-get update && apt-get install wget -y
RUN wget -o nyanbot.tar.gz https://github.com/commojun/nyanbot/releases/download/v${VERSION}/nyanbot_${VERSION}_Linux_armv6.tar.gz && \
    tar -xvf ./nyanbot_${VERSION}_Linux_armv6.tar.gz

FROM debian:bullseye-slim

RUN apt-get update && apt-get install -y locales && rm -rf /var/lib/apt/lists/* \
    && localedef -i ja_JP -c -f UTF-8 -A /usr/share/locale/locale.alias ja_JP.UTF-8
ENV LANG ja_JP.UTF-8

COPY --from=builder /nyan-* /usr/local/bin/

RUN useradd -r app
USER app
WORKDIR /home/app/

CMD ["nyan-server"]
