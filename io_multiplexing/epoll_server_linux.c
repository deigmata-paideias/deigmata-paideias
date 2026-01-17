#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <sys/types.h>
#include <sys/socket.h>
#include <netinet/in.h>
#include <arpa/inet.h>
#include <sys/epoll.h>
#include <errno.h>
#include <fcntl.h>

#define MAX_EVENTS 10
#define BUFFER_SIZE 1024
#define PORT 8888

// 设置socket为非阻塞
void setnonblocking(int sock) {
    int opts = fcntl(sock, F_GETFL);
    if (opts < 0) {
        perror("fcntl(F_GETFL)");
        exit(EXIT_FAILURE);
    }
    opts = (opts | O_NONBLOCK);
    if (fcntl(sock, F_SETFL, opts) < 0) {
        perror("fcntl(F_SETFL)");
        exit(EXIT_FAILURE);
    }
}

int main() {
    int server_fd;
    struct sockaddr_in server_addr, client_addr;
    char buffer[BUFFER_SIZE];
    int addrlen = sizeof(client_addr);

    // 创建服务器socket
    if ((server_fd = socket(AF_INET, SOCK_STREAM, 0)) == 0) {
        perror("Socket creation failed");
        exit(EXIT_FAILURE);
    }

    // 设置socket选项
    int opt = 1;
    if (setsockopt(server_fd, SOL_SOCKET, SO_REUSEADDR, &opt, sizeof(opt))) {
        perror("Setsockopt failed");
        exit(EXIT_FAILURE);
    }

    // 配置服务器地址
    server_addr.sin_family = AF_INET;
    server_addr.sin_addr.s_addr = INADDR_ANY;
    server_addr.sin_port = htons(PORT);

    // 绑定socket
    if (bind(server_fd, (struct sockaddr *)&server_addr, sizeof(server_addr)) < 0) {
        perror("Bind failed");
        exit(EXIT_FAILURE);
    }

    // 监听连接
    if (listen(server_fd, 3) < 0) {
        perror("Listen failed");
        exit(EXIT_FAILURE);
    }

    // 创建epoll实例
    int epoll_fd = epoll_create1(0);
    if (epoll_fd == -1) {
        perror("Epoll create failed");
        exit(EXIT_FAILURE);
    }

    // 添加服务器socket到epoll
    struct epoll_event ev;
    ev.events = EPOLLIN;
    ev.data.fd = server_fd;
    if (epoll_ctl(epoll_fd, EPOLL_CTL_ADD, server_fd, &ev) == -1) {
        perror("Epoll_ctl: server_fd");
        exit(EXIT_FAILURE);
    }

    // 事件数组
    struct epoll_event events[MAX_EVENTS];
    
    printf("Server listening on port %d\n", PORT);

    while (1) {
        int nfds = epoll_wait(epoll_fd, events, MAX_EVENTS, -1);
        if (nfds == -1) {
            perror("Epoll_wait");
            break;
        }

        for (int n = 0; n < nfds; ++n) {
            if (events[n].data.fd == server_fd) {
                // 处理新连接
                int new_socket = accept(server_fd, (struct sockaddr *)&client_addr, 
                                     (socklen_t*)&addrlen);
                if (new_socket == -1) {
                    perror("Accept");
                    continue;
                }

                setnonblocking(new_socket);
                printf("New connection, socket fd: %d, IP: %s, Port: %d\n",
                       new_socket, inet_ntoa(client_addr.sin_addr), 
                       ntohs(client_addr.sin_port));

                // 添加新客户端到epoll
                ev.events = EPOLLIN | EPOLLET;  // 边缘触发
                ev.data.fd = new_socket;
                if (epoll_ctl(epoll_fd, EPOLL_CTL_ADD, new_socket, &ev) == -1) {
                    perror("Epoll_ctl: new_socket");
                    close(new_socket);
                }
            } else {
                // 处理客户端数据
                int fd = events[n].data.fd;
                int valread = read(fd, buffer, BUFFER_SIZE);
                
                if (valread <= 0) {
                    printf("Client disconnected, socket fd: %d\n", fd);
                    epoll_ctl(epoll_fd, EPOLL_CTL_DEL, fd, NULL);
                    close(fd);
                } else {
                    buffer[valread] = '\0';
                    printf("Received message from client %d: %s\n", fd, buffer);
                    send(fd, buffer, strlen(buffer), 0);
                }
            }
        }
    }

    close(epoll_fd);
    return 0;
}
