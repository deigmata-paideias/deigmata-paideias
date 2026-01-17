# Select 

使用 select 实现的简单 TCP 服务器 demo，可以处理多个客户端连接

```shell
gcc select_server.c -o select_server
./select_server

# Poll
gcc poll_server.c -o poll_server
./poll_server

# Epoll linux
gcc epoll_server_linux.c -o epoll_server_linux
./epoll_server_linux

# Epoll mac
gcc epoll_server_mac.c -o epoll_server_mac
./epoll_server_mac

telnet localhost 8888
```

```shell
$ ./select_server
Server listening on port 8888
New connection, socket fd: 4, IP: 127.0.0.1, Port: 60858
New connection, socket fd: 5, IP: 127.0.0.1, Port: 60902
New connection, socket fd: 6, IP: 127.0.0.1, Port: 60929
New connection, socket fd: 7, IP: 127.0.0.1, Port: 60943
New connection, socket fd: 8, IP: 127.0.0.1, Port: 60955
Received message from client 8: 

Client disconnected, socket fd: 6
Client disconnected, socket fd: 4
Client disconnected, socket fd: 5
Client disconnected, socket fd: 8
Client disconnected, socket fd: 7
```
