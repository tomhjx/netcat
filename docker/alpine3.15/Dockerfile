FROM alpine:3.15
RUN apk add --update tzdata alpine-sdk autoconf automake libressl-dev libpcap-dev libnet-dev json-c-dev fts-dev scons \
    wget \
    tcpflow \
    tshark \
    tcpdump && \
    cd /tmp/ && \ 
    wget -O ssldump.tar.gz https://github.com/adulau/ssldump/archive/refs/tags/v1.4.tar.gz && \
    tar -zxvf ssldump.tar.gz && \
    cd ssldump-1.4 && \
    ./autogen.sh && \
    ./configure --prefix=/usr/local && \
    make && make install && \
    cd /tmp/ && \ 
    wget -O mongo.gz  https://github.com/mongodb/mongo/archive/refs/tags/r5.3.2.tar.gz && \
    tar -zxvf mongo.gz && \
    cd mongo-r5.3.2 && \
    scons mongosniff && \
    rm -rf /var/cache/apk/* /tmp/* /sbin/halt /sbin/poweroff /sbin/reboot

COPY mybin /mybin
COPY entrypoint.sh /

RUN cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    echo "Asia/Shanghai" > /etc/timezone \
    chmod +x -R /mybin; \
    chmod +x /entrypoint.sh 
ENTRYPOINT ["/entrypoint.sh"]
CMD ["--help"]