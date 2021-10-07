#!/bin/sh

# tcpflow -c -p -i any port 3306| grep -i -E "TRANSACTION|REPLACE|SELECT|UPDATE|DELETE|INSERT|SET|COMMIT|ROLLBACK|CREATE|DROP|ALTER|CALL" 
# | sed 's%\(.*\)\([.]\{4\}\)\(.*\)%\3%'
tcpflow -c -p -i any dst port 3306|sed 's%\(.*\)\([.]\{4\}\)\(.*\)%\3%'

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