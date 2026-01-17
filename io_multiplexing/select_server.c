#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <sys/types.h>
#include <sys/socket.h>
#include <netinet/in.h>
#include <arpa/inet.h>
#include <errno.h>
#include <fcntl.h>

#define MAX_CLIENTS 10
#define BUFFER_SIZE 1024
#define PORT 8888

int main() {
    int server_fd, client_fds[MAX_CLIENTS];
    fd_set read_fds;
    struct sockaddr_in server_addr, client_addr;
    char buffer[BUFFER_SIZE];
    int max_fd, activity, i, valread, sd;
    int addrlen = sizeof(client_addr);

    // 初始化客户端socket数组
    for (i = 0; i < MAX_CLIENTS; i++) {
        client_fds[i] = 0;
    }

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

    printf("Server listening on port %d\n", PORT);

    while (1) {
        // 清空文件描述符集合
        FD_ZERO(&read_fds);

        // 添加服务器socket到集合
        FD_SET(server_fd, &read_fds);
        max_fd = server_fd;

        // 添加客户端socket到集合
        for (i = 0; i < MAX_CLIENTS; i++) {
            sd = client_fds[i];
            if (sd > 0) {
                FD_SET(sd, &read_fds);
            }
            if (sd > max_fd) {
                max_fd = sd;
            }
        }

        // 等待活动发生
        activity = select(max_fd + 1, &read_fds, NULL, NULL, NULL);
        if ((activity < 0) && (errno != EINTR)) {
            perror("Select error");
        }

        // 处理新连接
        if (FD_ISSET(server_fd, &read_fds)) {
            int new_socket;
            if ((new_socket = accept(server_fd, (struct sockaddr *)&client_addr, 
                                   (socklen_t*)&addrlen)) < 0) {
                perror("Accept failed");
                exit(EXIT_FAILURE);
            }

            printf("New connection, socket fd: %d, IP: %s, Port: %d\n",
                   new_socket, inet_ntoa(client_addr.sin_addr), 
                   ntohs(client_addr.sin_port));

            // 添加新客户端到数组
            for (i = 0; i < MAX_CLIENTS; i++) {
                if (client_fds[i] == 0) {
                    client_fds[i] = new_socket;
                    break;
                }
            }
        }

        // 处理客户端数据
        for (i = 0; i < MAX_CLIENTS; i++) {
            sd = client_fds[i];

            if (FD_ISSET(sd, &read_fds)) {
                // 检查是否是断开连接
                if ((valread = read(sd, buffer, BUFFER_SIZE)) == 0) {
                    printf("Client disconnected, socket fd: %d\n", sd);
                    close(sd);
                    client_fds[i] = 0;
                }
                // 处理客户端消息
                else {
                    buffer[valread] = '\0';
                    printf("Received message from client %d: %s\n", sd, buffer);
                    // 回显消息给客户端
                    send(sd, buffer, strlen(buffer), 0);
                }
            }
        }
    }

    return 0;
}
