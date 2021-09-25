# netcat
集成网络抓包工具及提供便捷的使用方式，便于快速实现网络问题的监控、分析


## 使用

```

docker build ./docker/build -t tomhjx/netcat

```


* 截取某个容器发起的请求

```
docker run -it --rm --net container:目标容器名称 tomhjx/netcat what-mysql

```

## 依赖

* [TCPFLOW](https://github.com/simsong/tcpflow)