
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <sys/types.h>
#include <sys/socket.h>
#include <netinet/in.h>
#include <arpa/inet.h>
#include <sys/event.h>
#include <errno.h>
#include <fcntl.h>

#define MAX_EVENTS 10
#define BUFFER_SIZE 1024
#define PORT 8888

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

    // 创建kqueue实例
    int kq = kqueue();
    if (kq == -1) {
        perror("Kqueue creation failed");
        exit(EXIT_FAILURE);
    }

    // 设置服务器socket的事件
    struct kevent change_event;
    EV_SET(&change_event, server_fd, EVFILT_READ, EV_ADD, 0, 0, NULL);
    if (kevent(kq, &change_event, 1, NULL, 0, NULL) == -1) {
        perror("Kevent register failed");
        exit(EXIT_FAILURE);
    }

    // 事件数组
    struct kevent events[MAX_EVENTS];
    
    printf("Server listening on port %d\n", PORT);

    while (1) {
        int new_events = kevent(kq, NULL, 0, events, MAX_EVENTS, NULL);
        if (new_events == -1) {
            perror("Kevent wait failed");
            exit(EXIT_FAILURE);
        }

        for (int i = 0; i < new_events; i++) {
            int event_fd = events[i].ident;

            if (event_fd == server_fd) {
                // 处理新连接
                int new_socket = accept(server_fd, (struct sockaddr *)&client_addr, 
                                     (socklen_t*)&addrlen);
                if (new_socket == -1) {
                    perror("Accept failed");
                    continue;
                }

                setnonblocking(new_socket);
                printf("New connection, socket fd: %d, IP: %s, Port: %d\n",
                       new_socket, inet_ntoa(client_addr.sin_addr), 
                       ntohs(client_addr.sin_port));

                // 注册新socket的读事件
                struct kevent new_event;
                EV_SET(&new_event, new_socket, EVFILT_READ, EV_ADD, 0, 0, NULL);
                if (kevent(kq, &new_event, 1, NULL, 0, NULL) == -1) {
                    perror("Kevent register failed");
                    close(new_socket);
                }
            } else {
                // 处理客户端数据
                int valread = read(event_fd, buffer, BUFFER_SIZE);
                
                if (valread <= 0) {
                    printf("Client disconnected, socket fd: %d\n", event_fd);
                    // 从kqueue中删除socket
                    struct kevent del_event;
                    EV_SET(&del_event, event_fd, EVFILT_READ, EV_DELETE, 0, 0, NULL);
                    kevent(kq, &del_event, 1, NULL, 0, NULL);
                    close(event_fd);
                } else {
                    buffer[valread] = '\0';
                    printf("Received message from client %d: %s\n", event_fd, buffer);
                    send(event_fd, buffer, strlen(buffer), 0);
                }
            }
        }
    }

    close(kq);
    return 0;
}
