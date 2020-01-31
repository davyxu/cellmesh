# cellmesh
基于Service Mesh的分布式高性能游戏服务器框架

# 概念及关系

## 端(Peer)
link/peermgr.go 中管理Peer的生命期

### 侦听端(Acceptor Peer)
侦听在本地端口的Peer, 等待其他进程的Peer连接过来

### 连接端(Connector Peer)
连接端分为两种端实现:

* 服务间通信Peer

    使用cellnet的peer/tcp实现
 
* 连接已有的第三方服务, 例如: Redis, MySQL等

    使用cellnet的peer/redix, peer/mysql等实现
    
## Peer相关API

* link.GetPeerDesc(cellnet.Peer) *discovery.ServiceDesc

    获取Peer关联的服务发现信息

## 链接(Link)
当进程与另外一个进程建立服务间通信Socket连接(非第三方服务)时, 将创建一条连接两端Peer的链接(Link)

链接使用cellnet的Session实现

### 实现

当link.LinkService被调用后, 将从服务发现同步现有的服务,并自动连接.当Link断开后, 连接会自动尝试重连

link在linkmgr.go中被管理

## Link相关API

* link.GetLink(string) cellnet.Session

    从SvcID获取链接的Session

* link.GetLinkSvcID(cellnet.Session) string

    从链接的Session获得链接指向的SvcID

* link.GetLinkSvcName(cellnet.Session) string

    从链接的Session获得链接指向的SvcName

# TODO

- 服务间通信hub

- 带上下文的日志及相关的编写习惯, 改进日志分析难度

- rpc基本通信及接口,支持自动Relay

- 以Notify发送通知消息

- 网关

- 根据消息类型自动确定发送目的Peer, 全局只需一个Send实现


