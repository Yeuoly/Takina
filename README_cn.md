# Takina
Takina是由Yeuoly开发的一款用于做Docker集群内网穿透的代理工具，名字取自《Lycoris Recoil》中的泷奈<br>

[README](README.md) | [中文文档](README_cn.md)

我们知道，当Docker网络处于Internal模式并且网络模式为overlay的时候，PublishPort功能将失效，这意味着如果我们想要通过传统的主机端口映射的方式实现内网穿透将变得困难，但是上帝给你关上了一扇门，总是会给你开一扇窗容器是可以连接多个网络的，因此，我们可以在网络中接入一个用于做端口转发的容器（反向代理），它负责将外网流量转发到网络中的各个节点 <br>

Takina拥有几个基础组件

- TakinaServer
- TakinaClientDaemon
- TakinaClientCli

<br>
你可以在`/cmd/*`中看到各个组件的入口
<br>

Takina被设计为一个分布式的工具，客户端运行在容器中，这个容器需要连接到所有需要被穿透的网络中，同时连接到Docker的默认网络中，当然，如果你喜欢，也可以自建一个网络，只不过这个网络需要能够连接互联网以便Takina可以接收来自外部的流量<br>

当有流量到达TakinaServer的时候，它会将流量转发到对应的TakinaClientDaemon，TakinaClientDaemon会紧接着将流量转发到对应的网络，因此，TakinaClientDaemon需要连接到所有需要被穿透的网络中，这样可以确保容器网络处于隔离状态的情况下，也能够实现内网穿透，并且无法被外部感知<br>

如果您想更便利地使用Takina，可以使用Kisara项目，它封装了Takina在内，并且Kisara本身就是一个分布式Docker调度工具，提供了便捷的API接口<br>

# TakinaServer
Takina有一个重要的任务就是负载均衡，我们有时候有3台公网机器，它们作为流量入口负责将流量转发到容器内，Server可以有多个，但需要在Client中做配置

# TakinaClientDaemon
最终Takina是需要被部署到容器内部的，ClientDaemon就是负责流量转发的客户端，它将长期被挂起在后台，并接收TakinaClientCli的命令

# TakinaClientCli
因为当容器启动后，到控制端的网络已经全面被隔离，通信方式只有命令执行或虚拟文件等，Takina采用的是命令执行的访问，由控制端发出指令，而ClientCli就负责接收来自控制端的指令，并且将指令转发到ClientDaemon

# Third-party libraries Thanks

- [github.com/gin-gonic/gin](https://github.com/gin-gonic/gin)
- [github.com/fatedier/frp](https://github.com/fatedier/frp)
