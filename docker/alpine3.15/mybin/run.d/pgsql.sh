#!/bin/sh


# https://wiki.wireshark.org/PostgresProtocol
# https://www.wireshark.org/docs/dfref/p/pgsql.html

tshark -i any -T fields -e pgsql.query -e pgsql.val.data -e pgsql.parameter_name -e pgsql.parameter_value -r /data/resources/pgsql.cap
