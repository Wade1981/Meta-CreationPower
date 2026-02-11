// network_module.h - AIAgent网络通信模块头文件 (跨平台版本)

#ifndef NETWORK_MODULE_H
#define NETWORK_MODULE_H

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <time.h>
#include <stdarg.h>

// 跨平台网络头文件
#ifdef _WIN32
    #include <winsock2.h>
    #include <ws2tcpip.h>
    #include <windows.h>
    #define SOCKET_ERROR_RETURN -1
    #define CLOSE_SOCKET(s) closesocket(s)
    #define SOCKET_T SOCKET
    #define INVALID_SOCKET_VALUE INVALID_SOCKET
#else
    #include <sys/socket.h>
    #include <netinet/in.h>
    #include <arpa/inet.h>
    #include <unistd.h>
    #include <fcntl.h>
    #define SOCKET_ERROR_RETURN -1
    #define CLOSE_SOCKET(s) close(s)
    #define SOCKET_T int
    #define INVALID_SOCKET_VALUE -1
    #define SOCKET_ERROR -1
    #define WSAGetLastError() errno
    #define WSACleanup() (void)0
    #define WSAStartup(a, b) 0
#endif

// 网络节点结构体
typedef struct {
    char id[64];
    char hostname[256];
    int port;
    SOCKET_T socket;
    struct sockaddr_in address;
    int is_connected;
    char status[32];
} NetworkNode;

// 消息结构体
typedef struct {
    char id[64];
    char type[64];
    char source[64];
    char destination[64];
    int payload_size;
    char* payload;
    long timestamp;
} NetworkMessage;

// 网络模块配置
typedef struct {
    int max_connections;
    int buffer_size;
    int timeout;
    int keepalive_interval;
    char log_file[256];
} NetworkConfig;

// 网络模块状态
typedef struct {
    int initialized;
    int running;
    int active_connections;
    int total_messages_sent;
    int total_messages_received;
    long uptime;
} NetworkStatus;

// 协议类型
typedef enum {
    PROTOCOL_HTTP,
    PROTOCOL_TCP,
    PROTOCOL_UDP,
    PROTOCOL_WEBSOCKET
} ProtocolType;

// 初始化和清理函数
int network_module_init(NetworkConfig* config);
int network_module_cleanup();

// 节点管理函数
NetworkNode* create_network_node(const char* id, const char* hostname, int port);
int destroy_network_node(NetworkNode* node);
int connect_node(NetworkNode* node);
int disconnect_node(NetworkNode* node);

// 消息处理函数
NetworkMessage* create_message(const char* type, const char* source, const char* destination, const char* payload);
int destroy_message(NetworkMessage* message);
int send_message(NetworkNode* node, NetworkMessage* message);
NetworkMessage* receive_message(NetworkNode* node);

// 服务器函数
SOCKET_T create_server(int port, int backlog);
NetworkNode* accept_connection(SOCKET_T server_socket);
int close_server(SOCKET_T server_socket);

// 工具函数
int serialize_message(NetworkMessage* message, char** buffer, int* buffer_size);
NetworkMessage* deserialize_message(const char* buffer, int buffer_size);
int get_network_status(NetworkStatus* status);
int log_network_event(const char* format, ...);

// UDP函数
int create_udp_socket();
int send_udp_message(int socket, const char* host, int port, const char* message, int message_len);
int receive_udp_message(int socket, char* buffer, int buffer_size, struct sockaddr_in* client_addr, socklen_t* addr_len);

// 常量定义
#define MAX_BUFFER_SIZE 4096
#define MAX_NODE_ID_LENGTH 64
#define MAX_HOSTNAME_LENGTH 256
#define MAX_MESSAGE_TYPE_LENGTH 64
#define MAX_PAYLOAD_SIZE 1048576 // 1MB

#endif // NETWORK_MODULE_H
