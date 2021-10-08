#!/bin/sh

# tcpflow -c -i any -e http
tshark -i any -d tcp.port==80,http -T fields -e tcp.analysis.ack_rtt -e http.request.host   -e http.request.uri -e http.file_data -e http.request.line -e http.response.line

# tshark -i any  -d tcp.port==80,http

# tcpdump -A -s 0 'tcp port 80 and (((ip[2:2] - ((ip[0]&0xf)<<2)) - ((tcp[12]&0xf0)>>2)) != 0)