# cellmesh
基于cellnet的游戏服务框架

# 基础包依赖

```
    go get -v github.com/hashicorp/consul
```

# 启动服务发现

```
    cd shell
    ./StartConsul.sh
```


# 本工程包含什么？

- 客户端直连服务器例子(demo/login)

  使用基于cellnet封装网络通信+Consul服务发现的例子

- 网关例子(demo/agent)

- 基于Consul的服务发现封装(discovery)

- 协议代码生成工具

- 网关路由规则自动上传Consul

# 概念

## Service（服务）

一个Service为一套连接器或侦听器，挂接消息处理的派发器Dispatcher

- 侦听端口自动分配

  Service默认启动时以地址:0启动，网络底层自动分配端口，由cellmesh将服务信息报告到服务发现

  其他Service发现新的服务进入网络时，根据需要自动连接服务


## Agent（网关）

frontend(前端)

与客户端连接的前端通信的侦听器

- 使用心跳+底层断开通知确认User断开

- 客户端断开时通知后台服务器(ClientClosedACK消息)

backend(后端)

与后台服务器通信的侦听器

- 后台认证

  后台服务通过RouterBindUserACK消息,将后台连接与客户端绑定,客户端固定将对应消息发送到绑定的后台服务器.

- 后台断线重连

  后台服务断开重连时，自动维护连接，保证客户端正常收发后台消息

routerule(路由规则)

在proto文件中,消息的RouteRule属性描述如何路由消息到指定的后台服务器

- 阻断(不填写RouteRule)

   消息被路由阻断,无法发送到后台服务器

- 通透(RouteRule=pass)

  消息始终被路由到后台服务器

- 后台认证(RouteRule=auth)

  消息需要后台认证后才可被路由



## Connection Management（连接维护）

从服务发现的服务信息，创建到不同服务间的长连接。同时在这些连接断开时维护连接

逻辑中根据策略从已有连接及信息选择出连接与目标通信，例如：选择负载最低的一台游戏服务器


# TODO
- [x] 基于Consul的服务发现及Watch机制
- [x] 网关基本逻辑
- [x] 带服务发现的连接器,侦听器
- [x] 网关会话绑定
- [x] 服务发现Consul的KV封装
- [x] 网关路由规则生成,上传和下载
- [x] 系统消息响应入口
- [ ] 网关广播
- [ ] 服务在线人数更新及连接选择
- [ ] 分布式hub
- [ ] Docker部署
- [ ] 网关心跳处理
- [ ] 登录服务器，JWT验证
- [ ] 玩家数据读取（mysql）,goorm
- [ ] 游戏服务器，成长逻辑，花钱升级等级
- [ ] 社交服务器，聊天逻辑（redis)
- [ ] 机器人
