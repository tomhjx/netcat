#!/bin/sh

# tshark -i any -n -f 'tcp port 27017'
tshark -i any -r /data/resources/mongodb0.pcap