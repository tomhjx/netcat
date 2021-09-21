FROM alpine:3.14.2
RUN apk add --update tcpflow tcpdump perl && \
    rm -rf /var/cache/apk/* /tmp/* /sbin/halt /sbin/poweroff /sbin/reboot
COPY mybin /mybin
COPY entrypoint.sh /
RUN chmod +x -R /mybin; \
    chmod +x /entrypoint.sh 
ENTRYPOINT ["/entrypoint.sh"]
CMD ["http"]