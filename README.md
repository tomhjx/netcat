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

* ? [Percona Toolkit](https://www.percona.com/software/database-tools/percona-toolkit)


## 参考

* MySQL Query
    * https://www.shuzhiduo.com/A/nAJvvNXoJr/