# netcat
集成网络抓包工具及提供便捷的使用方式，便于快速实现网络问题的监控、分析


## 使用

```

docker build ./docker/build -t tomhjx/netcat

```


* debug

```

docker run -it --rm -v /Users/tom/Repos/github.com/tomhjx/netcat/docker/build/mybin:/mybin  -v /Users/tom/Repos/github.com/tomhjx/lab/php-framework-thinkphp5/src:/res  --net container:php-framework-thinkphp5_thinkphp5-fpm_1 tomhjx/netcat /bin/sh

```



* 截取某个容器发起的请求

```
docker run -it --rm --net container:目标容器名称 tomhjx/netcat what-mysql

```





## 依赖

* [TCPFLOW](https://github.com/simsong/tcpflow)

* [tshark](https://www.wireshark.org/docs/man-pages/tshark.html)

* ? [Percona Toolkit](https://www.percona.com/software/database-tools/percona-toolkit)


## 参考

* MySQL Query
    * https://www.shuzhiduo.com/A/nAJvvNXoJr/
    * https://www.wireshark.org/docs/dfref/m/mysql.html
    * https://plantegg.github.io/2019/06/21/%E5%B0%B1%E6%98%AF%E8%A6%81%E4%BD%A0%E6%87%82%E6%8A%93%E5%8C%85--WireShark%E4%B9%8B%E5%91%BD%E4%BB%A4%E8%A1%8C%E7%89%88tshark/