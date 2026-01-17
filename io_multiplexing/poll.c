#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <sys/types.h>
#include <sys/socket.h>
#include <netinet/in.h>
#include <arpa/inet.h>
#include <poll.h>
#include <errno.h>

#define MAX_CLIENTS 10
#define BUFFER_SIZE 1024
#define PORT 8888

int main() {
    int server_fd;
    struct sockaddr_in server_addr, client_addr;
    char buffer[BUFFER_SIZE];
    struct pollfd fds[MAX_CLIENTS + 1];  // +1 for server socket
    int nfds = 1, current_size = 0;
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

    // 初始化 poll 结构
    memset(fds, 0, sizeof(fds));
    fds[0].fd = server_fd;
    fds[0].events = POLLIN;

    printf("Server listening on port %d\n", PORT);

    while (1) {
        // 等待事件发生
        int ret = poll(fds, nfds, -1);
        if (ret < 0) {
            perror("Poll failed");
            break;
        }

        // 检查是否有新连接
        if (fds[0].revents & POLLIN) {
            int new_socket = accept(server_fd, (struct sockaddr *)&client_addr, 
                                  (socklen_t*)&addrlen);
            if (new_socket < 0) {
                perror("Accept failed");
                continue;
            }

            printf("New connection, socket fd: %d, IP: %s, Port: %d\n",
                   new_socket, inet_ntoa(client_addr.sin_addr), 
                   ntohs(client_addr.sin_port));

            // 添加新连接到 poll 集合
            if (nfds < MAX_CLIENTS + 1) {
                fds[nfds].fd = new_socket;
                fds[nfds].events = POLLIN;
                nfds++;
            } else {
                printf("Too many connections\n");
                close(new_socket);
            }
        }

        // 处理客户端数据
        for (int i = 1; i < nfds; i++) {
            if (fds[i].revents & POLLIN) {
                int valread = read(fds[i].fd, buffer, BUFFER_SIZE);
                if (valread <= 0) {
                    printf("Client disconnected, socket fd: %d\n", fds[i].fd);
                    close(fds[i].fd);
                    // 移除断开的客户端
                    for (int j = i; j < nfds - 1; j++) {
                        fds[j] = fds[j + 1];
                    }
                    nfds--;
                    i--;
                } else {
                    buffer[valread] = '\0';
                    printf("Received message from client %d: %s\n", fds[i].fd, buffer);
                    send(fds[i].fd, buffer, strlen(buffer), 0);
                }
            }
        }
    }

    return 0;
}
