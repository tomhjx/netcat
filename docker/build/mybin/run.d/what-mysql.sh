#!/bin/sh

# tcpflow -c -p -i any port 3306| grep -i -E "TRANSACTION|REPLACE|SELECT|UPDATE|DELETE|INSERT|SET|COMMIT|ROLLBACK|CREATE|DROP|ALTER|CALL" 
# | sed 's%\(.*\)\([.]\{4\}\)\(.*\)%\3%'
# tcpflow -c -p -i any dst port 3306|sed 's%\(.*\)\([.]\{4\}\)\(.*\)%\3%'

#parse 8507/4444 as mysql protocol, default only parse 3306 as mysql.
# tshark -i eth0 -d tcp.port==8507,mysql -T fields -e mysql.query 'port 8507'
# tshark -i any -c 50 -d tcp.port==4444,mysql -Y " ((tcp.port eq 4444 )  )" -o tcp.calculate_timestamps:true -T fields -e frame.number -e frame.time_epoch  -e frame.time_delta_displayed  -e ip.src -e tcp.srcport -e tcp.dstport -e ip.dst -e tcp.time_delta -e tcp.stream -e tcp.len -e mysql.query



#query time
# tshark -i eth0 -Y " ((tcp.port eq 3306 ) and tcp.len>0 )" -o tcp.calculate_timestamps:true -T fields -e frame.number -e frame.time_epoch  -e frame.time_delta_displayed  -e ip.src -e tcp.srcport -e tcp.dstport -e ip.dst -e tcp.time_delta -e tcp.stream -e tcp.len -e mysql.query

tshark -i any -d tcp.port==3306,mysql -T fields -e tcp.analysis.ack_rtt -e ip.src -e tcp.srcport -e tcp.dstport -e ip.dst  -e mysql.query

# tshark -i any -t ad  -d tcp.port==3306,mysql -o tcp.calculate_timestamps:true -T fields -e frame.number -e frame.t
# ime_epoch  -e frame.time_delta_displayed  -e ip.src -e tcp.srcport -e tcp.dstport -e ip.dst -e tcp.time_delta -e tcp.s
# tream -e tcp.len -e tcp.analysis.ack_rtt -e mysql.query

# tcpdump -i any -s 0 -l -w - dst port 3306 | strings | perl -e '
# while(<>) { chomp; next if /^[^ ]+[ ]*$/;
#     if(/^(SELECT|UPDATE|DELETE|INSERT|SET|COMMIT|ROLLBACK|CREATE|DROP|ALTER|CALL)/i)
#     {
#         if (defined $q) { print "$q\n"; }
#         $q=$_;
#     } else {
#         $_ =~ s/^[ \t]+//; $q.=" $_";
#     }
# }'