#!/bin/sh

echo "https"


# tcpdump -i any dst host api.juejin.cn

ssldump -k /res/rsa_private_key.pem -Ad  -i any host api.juejin.cn
# ssldump -k /ssldump.pem -Ad  -i any host api.juejin.cn