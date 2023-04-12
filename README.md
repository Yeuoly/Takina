# Takina

Takina is a proxy tool developed by Yeuoly for doing Docker cluster intranet penetration. The name comes from Takina, a character in \"Lycoris Recoil\".

[README](README.md) | [中文文档](README_cn.md)

When the Docker network is in Internal mode, the PublishPort function will be disabled. This means that if we want to achieve intranet penetration through the traditional host port mapping method, it will become difficult. However, when one door closes, another opens. Containers can access the internet, so we can add a container for port forwarding (reverse proxy) in the network. It is responsible for forwarding external traffic to various nodes in the network.

Takina has several basic components:

- TakinaServer
- TakinaClientDaemon
- TakinaClientCli

You can find the entry points for each component in `/cmd/*`.

If you want to use Takina more conveniently, you can use the Kisara project. It encapsulates Takina and is itself a distributed Docker scheduling tool that provides convenient API interfaces.

## TakinaServer

One important task of Takina is load balancing. Sometimes we have three public machines that serve as traffic entry points and are responsible for forwarding traffic to containers. There can be multiple servers, but they need to be configured in the client.

## TakinaClientDaemon

Ultimately, Takina needs to be deployed inside the container. ClientDaemon is the client responsible for traffic forwarding. It will be suspended in the background for a long time and receive commands from TakinaClientCli.

## TakinaClientCli

After the container starts, the network to the control side is fully isolated, and the only communication methods are command execution or virtual files. Takina uses command execution for access. The control side issues commands, and ClientCli is responsible for receiving the commands from the control side and forwarding them to ClientDaemon.

## Third-party libraries Thanks

- [github.com/aceld/zinx](https://github.com/aceld/zinx)
- [github.com/fatedier/frp](https://github.com/fatedier/frp)
