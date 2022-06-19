#!/bin/sh

# https://www.wireshark.org/docs/dfref/a/amqp.html
# tshark -i any -r /data/resources/rabbitmq.pcap -Y amqp
tshark -i any -Y amqp