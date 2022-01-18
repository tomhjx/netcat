# netcat
集成网络抓包工具及提供便捷的使用方式，便于快速实现网络问题的监控、分析


# 支持

* [*] MySQL
* [*] Redis
* [*] AMQP (eg. Rabbit MQ)
* [ ] Mongodb
* [ ] Kafka

## 使用

```

docker build ./docker/build -t tomhjx/netcat

```


* debug

```

docker run -it --rm -v /Users/tom/Repos/github.com/tomhjx/netcat/docker/build/mybin:/mybin  -v /Users/tom/Repos/github.com/tomhjx/lab/php-framework-thinkphp5/src:/res  --net container:php-framework-thinkphp5_thinkphp5-fpm_1 tomhjx/netcat /bin/sh

```


* pcap

```
docker run -it --rm --cap-add=ALL -v /Users/tom/Work/project/github.com/tomhjx/netcat/resources:/data/resources  --net container:lab_app_console tomhjx/netcat:0.1.0-alpine-3.14.2 /bin/sh -c "tcpdump -i eth0 -s 0 -w /data/resources/mysql.pcap"

docker run -it --rm --cap-add=ALL -v /Users/tom/Work/project/github.com/tomhjx/netcat/resources:/data/resources  --net container:lab_app_console tomhjx/netcat:0.1.0-alpine-3.14.2 /bin/sh -c "tcpdump -i eth0 -s 0 -w /data/resources/rabbit.pcap"


```



* 截取某个容器发起的请求

```
docker run -it --rm --net container:目标容器名称 tomhjx/netcat what-mysql

```



## 嗅探流程设计

* 处理器
    * 定义`输出器`信道，开启`输出器`协程
        * `输出器`
            * 对接输出设备 
        * 从`输出器`信道读取内容，作为`输出器`入参
        * 由`输出器`实现执行细节
    * 定义`解析器`信道，开启`解析器`协程
        * `解析器`
            * 解析内容，转换为结构化对象
        * 从`解析器`信道读取内容，作为`解析器`入参
        * 由`解析器`实现执行细节
        * 将`解析器`执行结果写入`输出器`信道
    * 启动`输入器`
        * `输入器`
            * 读取输入源（文件、流量）
            * 解包 
            * 包体结构化
        * 将`输入器`执行结果写入`解析器`信道 


## 依赖

* [TCPFLOW](https://github.com/simsong/tcpflow)

* [tshark](https://www.wireshark.org/docs/man-pages/tshark.html)

* ? [Percona Toolkit](https://www.percona.com/software/database-tools/percona-toolkit)


## 参考

* MySQL Query
    * https://www.shuzhiduo.com/A/nAJvvNXoJr/
    * https://www.wireshark.org/docs/dfref/m/mysql.html
    * https://plantegg.github.io/2019/06/21/%E5%B0%B1%E6%98%AF%E8%A6%81%E4%BD%A0%E6%87%82%E6%8A%93%E5%8C%85--WireShark%E4%B9%8B%E5%91%BD%E4%BB%A4%E8%A1%8C%E7%89%88tshark/
    * https://github.com/40t/go-sniffer