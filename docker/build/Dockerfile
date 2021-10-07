FROM alpine:3.14.2
RUN apk add --update tzdata alpine-sdk autoconf automake libressl-dev libpcap-dev libnet-dev json-c-dev fts-dev \
    wget \
    tcpflow \
    tshark \
    tcpdump && \
    rm -rf /var/cache/apk/* /tmp/* /sbin/halt /sbin/poweroff /sbin/reboot && \
    wget -O ssldump.tar.gz https://github.com/adulau/ssldump/archive/refs/tags/v1.4.tar.gz && \
    tar -zxvf ssldump.tar.gz && \
    cd ssldump-1.4 && \
    ./autogen.sh && \
    ./configure --prefix=/usr/local && \
    make && make install

COPY mybin /mybin
COPY entrypoint.sh /

RUN cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    echo "Asia/Shanghai" > /etc/timezone \
    chmod +x -R /mybin; \
    chmod +x /entrypoint.sh 
ENTRYPOINT ["/entrypoint.sh"]
CMD ["--help"]