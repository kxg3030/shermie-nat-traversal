#### 介绍
Go语言实现的内网穿透工具，使用tcp协议代理请求

#### 使用

- 服务器上执行
```go
go run Server.go --pass=123456 --port=8888
```
参数：

--pass：指定内网客户端连接你的服务器时的密码

--port：指定外网能够访问的端口

- 内网主机执行
```go
go run Client.go --bind=127.0.0.1:80 --pass=123456 --port=8888 --host=111.111.111.111
```
参数：

--pass：指定内网客户端连接你的服务器时的密码

--bind：内网客户端提供服务的端口

--port：公网服务器的外网port

--host：公网服务器的外网host

- 配置nginx
> 使用nginx将公网地址请求转发到穿透服务的端口
```go
# 内网穿透
location /nat {
    proxy_pass http://127.0.0.1:8888;		
}
```
