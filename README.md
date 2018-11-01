# cellmesh
基于cellnet的游戏服务框架

# 特点

## Based On Service Discovery(基于服务发现)

   通过服务发现自动实现服务互联,探测,挂接.无需配置服务器间的端口.

## Zero Config File(零配置文件)

   任何服务器配置均通过Consul的KV存取,无任何配置文件.

## Code Generation(代码生成)

   基于github.com/davyxu/protoplus的代码生成技术,迅速粘合逻辑与底层,代码好看易懂且高效.

   使用更简单更强大的schema编写协议, 并自动生成protobuf schema.

## Transport based on cellnet(基于cellnet的网络传输)

   提供强大的扩展及适配能力.

# 运行Demo

## 下载cellmesh源码

```
    go get -v github.com/davyxu/cellmesh
```

## 准备第三方包

```
    # 转到cellmesh的shell目录
    cd github.com/davyxu/cellmesh/shell
    sh ./DownloadThirdParty.sh
```

## 准备consul服务发现

```
    # 转到cellmesh的shell目录
    cd github.com/davyxu/cellmesh/shell

    # 这里会自动编译和运行consul
    sh ./StartConsul.sh
```

在浏览器里访问http://localhost:8500 地址可查看当前Consul状态

Mac下可直接运行下面的命令
```
    cd github.com/davyxu/cellmesh/shell
    sh ./ShowUI.sh
```

## 更新协议及上传路由规则

执行下面指令更新路由规则到Consul
```
    cd github.com/davyxu/cellmesh/demo/proto
    sh ./Gen.sh
```

## 启动demo服务

按照下面shell分别启动login, agent, game, client
```
    cd github.com/davyxu/cellmesh/shell

    sh ./RunDemoSvc.sh login
```

启动client后,可在命令行中输入文字作为聊天内容发送

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

  后台服务通过BindBackendACK消息,将后台连接与客户端绑定,客户端固定将对应消息发送到绑定的后台服务器.

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


# 开发进度
- [x] 基于Consul的服务发现及Watch机制
- [x] 网关基本逻辑
- [x] 带服务发现的连接器,侦听器
- [x] 网关会话绑定
- [x] 服务发现Consul的KV封装
- [x] 网关路由规则生成,上传和下载
- [x] 系统消息响应入口
- [x] 网关广播
- [x] 网关心跳处理
- [x] 频道订阅发布hub
- [ ] 服务在线人数更新及连接选择
- [ ] 登录服务器，JWT验证
- [ ] 玩家数据读取
- [ ] 游戏服务器，成长逻辑，花钱升级等级
- [ ] 社交服务器，聊天逻辑
- [ ] 机器人
