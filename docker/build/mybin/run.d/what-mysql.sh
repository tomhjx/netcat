#!/bin/sh

# tcpflow -c -p -i any port 3306| grep -i -E "TRANSACTION|REPLACE|SELECT|UPDATE|DELETE|INSERT|SET|COMMIT|ROLLBACK|CREATE|DROP|ALTER|CALL" 
# | sed 's%\(.*\)\([.]\{4\}\)\(.*\)%\3%'
tcpflow -C -p -i any dst port 3306|grep '[^\.\S]'