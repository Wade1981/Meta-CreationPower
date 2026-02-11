// network_module.c - AIAgent网络通信模块实现 (跨平台版本)

#include "network_module.h"

// 全局变量
static NetworkConfig g_config;
static NetworkStatus g_status;
static FILE* g_log_file = NULL;
static long g_start_time = 0;
static int g_initialized = 0;

// 初始化和清理函数
int network_module_init(NetworkConfig* config) {
    #ifdef _WIN32
    WSADATA wsaData;
    int result;
    #endif

    // 检查是否已经初始化
    if (g_initialized) {
        log_network_event("Network module already initialized\n");
        return 0;
    }

    // 复制配置
    if (config) {
        memcpy(&g_config, config, sizeof(NetworkConfig));
    } else {
        // 默认配置
        g_config.max_connections = 100;
        g_config.buffer_size = MAX_BUFFER_SIZE;
        g_config.timeout = 30000; // 30秒
        g_config.keepalive_interval = 60000; // 60秒
        strcpy(g_config.log_file, "network_module.log");
    }

    // 初始化Winsock (Windows only)
    #ifdef _WIN32
    result = WSAStartup(MAKEWORD(2, 2), &wsaData);
    if (result != 0) {
        log_network_event("WSAStartup failed: %d\n", result);
        return -1;
    }
    #endif

    // 打开日志文件
    g_log_file = fopen(g_config.log_file, "a");
    if (!g_log_file) {
        log_network_event("Failed to open log file: %s\n", g_config.log_file);
        #ifdef _WIN32
        WSACleanup();
        #endif
        return -1;
    }

    // 初始化状态
    memset(&g_status, 0, sizeof(NetworkStatus));
    g_status.initialized = 1;
    g_status.running = 1;
    g_start_time = time(NULL);
    g_status.uptime = 0;

    log_network_event("Network module initialized successfully\n");
    g_initialized = 1;
    return 0;
}

int network_module_cleanup() {
    if (!g_initialized) {
        return 0;
    }

    // 关闭日志文件
    if (g_log_file) {
        fclose(g_log_file);
        g_log_file = NULL;
    }

    // 清理Winsock (Windows only)
    #ifdef _WIN32
    WSACleanup();
    #endif

    g_status.initialized = 0;
    g_status.running = 0;
    g_initialized = 0;

    log_network_event("Network module cleaned up successfully\n");
    return 0;
}

// 节点管理函数
NetworkNode* create_network_node(const char* id, const char* hostname, int port) {
    NetworkNode* node = (NetworkNode*)malloc(sizeof(NetworkNode));
    if (!node) {
        log_network_event("Failed to allocate memory for network node\n");
        return NULL;
    }

    // 初始化节点
    strcpy(node->id, id);
    strcpy(node->hostname, hostname);
    node->port = port;
    node->socket = INVALID_SOCKET_VALUE;
    node->is_connected = 0;
    strcpy(node->status, "created");

    // 设置地址
    memset(&node->address, 0, sizeof(node->address));
    node->address.sin_family = AF_INET;
    node->address.sin_port = htons(port);

    // 解析主机名
    #ifdef _WIN32
    struct hostent* host = gethostbyname(hostname);
    if (host) {
        memcpy(&node->address.sin_addr.s_addr, host->h_addr, host->h_length);
    } else {
        node->address.sin_addr.s_addr = inet_addr(hostname);
    }
    #else
    struct addrinfo hints, *res;
    memset(&hints, 0, sizeof(hints));
    hints.ai_family = AF_INET;
    hints.ai_socktype = SOCK_STREAM;
    
    int status = getaddrinfo(hostname, NULL, &hints, &res);
    if (status == 0) {
        struct sockaddr_in* addr = (struct sockaddr_in*)res->ai_addr;
        node->address.sin_addr = addr->sin_addr;
        freeaddrinfo(res);
    } else {
        node->address.sin_addr.s_addr = inet_addr(hostname);
    }
    #endif

    log_network_event("Created network node: %s at %s:%d\n", id, hostname, port);
    return node;
}

int destroy_network_node(NetworkNode* node) {
    if (!node) {
        return -1;
    }

    // 断开连接
    if (node->is_connected) {
        disconnect_node(node);
    }

    // 关闭套接字
    if (node->socket != INVALID_SOCKET_VALUE) {
        CLOSE_SOCKET(node->socket);
    }

    free(node);
    log_network_event("Destroyed network node\n");
    return 0;
}

int connect_node(NetworkNode* node) {
    if (!node) {
        return -1;
    }

    // 创建套接字
    node->socket = socket(AF_INET, SOCK_STREAM, IPPROTO_TCP);
    if (node->socket == INVALID_SOCKET_VALUE) {
        log_network_event("Failed to create socket: %d\n", WSAGetLastError());
        return -1;
    }

    // 连接到目标
    int result = connect(node->socket, (struct sockaddr*)&node->address, sizeof(node->address));
    if (result == SOCKET_ERROR) {
        log_network_event("Failed to connect: %d\n", WSAGetLastError());
        CLOSE_SOCKET(node->socket);
        node->socket = INVALID_SOCKET_VALUE;
        return -1;
    }

    node->is_connected = 1;
    strcpy(node->status, "connected");
    g_status.active_connections++;

    log_network_event("Connected to node %s at %s:%d\n", node->id, node->hostname, node->port);
    return 0;
}

int disconnect_node(NetworkNode* node) {
    if (!node || !node->is_connected) {
        return -1;
    }

    // 关闭套接字
    if (node->socket != INVALID_SOCKET_VALUE) {
        CLOSE_SOCKET(node->socket);
        node->socket = INVALID_SOCKET_VALUE;
    }

    node->is_connected = 0;
    strcpy(node->status, "disconnected");
    g_status.active_connections--;

    log_network_event("Disconnected from node %s\n", node->id);
    return 0;
}

// 消息处理函数
NetworkMessage* create_message(const char* type, const char* source, const char* destination, const char* payload) {
    NetworkMessage* message = (NetworkMessage*)malloc(sizeof(NetworkMessage));
    if (!message) {
        log_network_event("Failed to allocate memory for message\n");
        return NULL;
    }

    // 生成消息ID
    sprintf(message->id, "msg_%ld_%d", time(NULL), rand());
    strcpy(message->type, type);
    strcpy(message->source, source);
    strcpy(message->destination, destination);
    message->timestamp = time(NULL);

    // 设置负载
    if (payload) {
        message->payload_size = strlen(payload);
        message->payload = (char*)malloc(message->payload_size + 1);
        if (!message->payload) {
            free(message);
            log_network_event("Failed to allocate memory for message payload\n");
            return NULL;
        }
        strcpy(message->payload, payload);
    } else {
        message->payload_size = 0;
        message->payload = NULL;
    }

    return message;
}

int destroy_message(NetworkMessage* message) {
    if (!message) {
        return -1;
    }

    if (message->payload) {
        free(message->payload);
    }
    free(message);
    return 0;
}

int send_message(NetworkNode* node, NetworkMessage* message) {
    if (!node || !message || !node->is_connected) {
        return -1;
    }

    // 序列化消息
    char* buffer = NULL;
    int buffer_size = 0;
    if (serialize_message(message, &buffer, &buffer_size) != 0) {
        return -1;
    }

    // 发送消息
    int bytes_sent = send(node->socket, buffer, buffer_size, 0);
    if (bytes_sent == SOCKET_ERROR) {
        log_network_event("Failed to send message: %d\n", WSAGetLastError());
        free(buffer);
        return -1;
    }

    free(buffer);
    g_status.total_messages_sent++;
    log_network_event("Sent message %s to node %s\n", message->id, node->id);
    return 0;
}

NetworkMessage* receive_message(NetworkNode* node) {
    if (!node || !node->is_connected) {
        return NULL;
    }

    // 接收消息长度
    int buffer_size = 0;
    int bytes_received = recv(node->socket, (char*)&buffer_size, sizeof(int), 0);
    if (bytes_received == SOCKET_ERROR) {
        log_network_event("Failed to receive message length: %d\n", WSAGetLastError());
        return NULL;
    }

    if (bytes_received == 0) {
        // 连接关闭
        node->is_connected = 0;
        strcpy(node->status, "disconnected");
        g_status.active_connections--;
        return NULL;
    }

    // 接收消息内容
    char* buffer = (char*)malloc(buffer_size);
    if (!buffer) {
        log_network_event("Failed to allocate memory for message buffer\n");
        return NULL;
    }

    bytes_received = recv(node->socket, buffer, buffer_size, 0);
    if (bytes_received == SOCKET_ERROR) {
        log_network_event("Failed to receive message: %d\n", WSAGetLastError());
        free(buffer);
        return NULL;
    }

    // 反序列化消息
    NetworkMessage* message = deserialize_message(buffer, buffer_size);
    free(buffer);

    if (message) {
        g_status.total_messages_received++;
        log_network_event("Received message %s from node %s\n", message->id, node->id);
    }

    return message;
}

// 服务器函数
SOCKET_T create_server(int port, int backlog) {
    SOCKET_T server_socket = socket(AF_INET, SOCK_STREAM, IPPROTO_TCP);
    if (server_socket == INVALID_SOCKET_VALUE) {
        log_network_event("Failed to create server socket: %d\n", WSAGetLastError());
        return INVALID_SOCKET_VALUE;
    }

    // 设置地址重用
    int opt = 1;
    if (setsockopt(server_socket, SOL_SOCKET, SO_REUSEADDR, (char*)&opt, sizeof(opt)) < 0) {
        log_network_event("Failed to set socket options: %d\n", WSAGetLastError());
        CLOSE_SOCKET(server_socket);
        return INVALID_SOCKET_VALUE;
    }

    // 设置地址
    struct sockaddr_in server_addr;
    memset(&server_addr, 0, sizeof(server_addr));
    server_addr.sin_family = AF_INET;
    server_addr.sin_addr.s_addr = INADDR_ANY;
    server_addr.sin_port = htons(port);

    // 绑定地址
    int result = bind(server_socket, (struct sockaddr*)&server_addr, sizeof(server_addr));
    if (result == SOCKET_ERROR) {
        log_network_event("Failed to bind server socket: %d\n", WSAGetLastError());
        CLOSE_SOCKET(server_socket);
        return INVALID_SOCKET_VALUE;
    }

    // 开始监听
    result = listen(server_socket, backlog);
    if (result == SOCKET_ERROR) {
        log_network_event("Failed to listen on server socket: %d\n", WSAGetLastError());
        CLOSE_SOCKET(server_socket);
        return INVALID_SOCKET_VALUE;
    }

    log_network_event("Server created on port %d\n", port);
    return server_socket;
}

NetworkNode* accept_connection(SOCKET_T server_socket) {
    SOCKET_T client_socket;
    struct sockaddr_in client_addr;
    socklen_t client_addr_len = sizeof(client_addr);

    // 接受连接
    client_socket = accept(server_socket, (struct sockaddr*)&client_addr, &client_addr_len);
    if (client_socket == INVALID_SOCKET_VALUE) {
        log_network_event("Failed to accept connection: %d\n", WSAGetLastError());
        return NULL;
    }

    // 创建节点
    char client_ip[INET_ADDRSTRLEN];
    inet_ntop(AF_INET, &(client_addr.sin_addr), client_ip, INET_ADDRSTRLEN);

    char node_id[64];
    sprintf(node_id, "client_%ld", time(NULL));

    NetworkNode* node = create_network_node(node_id, client_ip, ntohs(client_addr.sin_port));
    if (!node) {
        CLOSE_SOCKET(client_socket);
        return NULL;
    }

    node->socket = client_socket;
    memcpy(&node->address, &client_addr, sizeof(client_addr));
    node->is_connected = 1;
    strcpy(node->status, "connected");
    g_status.active_connections++;

    log_network_event("Accepted connection from %s:%d\n", client_ip, ntohs(client_addr.sin_port));
    return node;
}

int close_server(SOCKET_T server_socket) {
    if (server_socket == INVALID_SOCKET_VALUE) {
        return -1;
    }

    CLOSE_SOCKET(server_socket);
    log_network_event("Server closed\n");
    return 0;
}

// UDP函数
int create_udp_socket() {
    int socket_fd = socket(AF_INET, SOCK_DGRAM, 0);
    if (socket_fd == INVALID_SOCKET_VALUE) {
        log_network_event("Failed to create UDP socket: %d\n", WSAGetLastError());
        return INVALID_SOCKET_VALUE;
    }
    return socket_fd;
}

int send_udp_message(int socket, const char* host, int port, const char* message, int message_len) {
    struct sockaddr_in addr;
    memset(&addr, 0, sizeof(addr));
    addr.sin_family = AF_INET;
    addr.sin_port = htons(port);
    
    #ifdef _WIN32
    struct hostent* hostent = gethostbyname(host);
    if (hostent) {
        memcpy(&addr.sin_addr.s_addr, hostent->h_addr, hostent->h_length);
    } else {
        addr.sin_addr.s_addr = inet_addr(host);
    }
    #else
    struct addrinfo hints, *res;
    memset(&hints, 0, sizeof(hints));
    hints.ai_family = AF_INET;
    hints.ai_socktype = SOCK_DGRAM;
    
    int status = getaddrinfo(host, NULL, &hints, &res);
    if (status == 0) {
        struct sockaddr_in* in_addr = (struct sockaddr_in*)res->ai_addr;
        addr.sin_addr = in_addr->sin_addr;
        freeaddrinfo(res);
    } else {
        addr.sin_addr.s_addr = inet_addr(host);
    }
    #endif
    
    int bytes_sent = sendto(socket, message, message_len, 0, (struct sockaddr*)&addr, sizeof(addr));
    if (bytes_sent == SOCKET_ERROR) {
        log_network_event("Failed to send UDP message: %d\n", WSAGetLastError());
        return -1;
    }
    
    return bytes_sent;
}

int receive_udp_message(int socket, char* buffer, int buffer_size, struct sockaddr_in* client_addr, socklen_t* addr_len) {
    int bytes_received = recvfrom(socket, buffer, buffer_size, 0, (struct sockaddr*)client_addr, addr_len);
    if (bytes_received == SOCKET_ERROR) {
        log_network_event("Failed to receive UDP message: %d\n", WSAGetLastError());
        return -1;
    }
    
    return bytes_received;
}

// 工具函数
int serialize_message(NetworkMessage* message, char** buffer, int* buffer_size) {
    if (!message || !buffer || !buffer_size) {
        return -1;
    }

    // 计算缓冲区大小
    int size = sizeof(int) + // message ID length
               strlen(message->id) +
               sizeof(int) + // message type length
               strlen(message->type) +
               sizeof(int) + // source length
               strlen(message->source) +
               sizeof(int) + // destination length
               strlen(message->destination) +
               sizeof(int) + // payload size
               message->payload_size +
               sizeof(long); // timestamp

    // 分配缓冲区
    *buffer = (char*)malloc(size);
    if (!*buffer) {
        return -1;
    }

    // 序列化数据
    char* ptr = *buffer;
    int len;

    // 消息ID
    len = strlen(message->id);
    memcpy(ptr, &len, sizeof(int));
    ptr += sizeof(int);
    memcpy(ptr, message->id, len);
    ptr += len;

    // 消息类型
    len = strlen(message->type);
    memcpy(ptr, &len, sizeof(int));
    ptr += sizeof(int);
    memcpy(ptr, message->type, len);
    ptr += len;

    // 源
    len = strlen(message->source);
    memcpy(ptr, &len, sizeof(int));
    ptr += sizeof(int);
    memcpy(ptr, message->source, len);
    ptr += len;

    // 目标
    len = strlen(message->destination);
    memcpy(ptr, &len, sizeof(int));
    ptr += sizeof(int);
    memcpy(ptr, message->destination, len);
    ptr += len;

    // 负载大小
    memcpy(ptr, &message->payload_size, sizeof(int));
    ptr += sizeof(int);

    // 负载
    if (message->payload_size > 0 && message->payload) {
        memcpy(ptr, message->payload, message->payload_size);
        ptr += message->payload_size;
    }

    // 时间戳
    memcpy(ptr, &message->timestamp, sizeof(long));

    *buffer_size = size;
    return 0;
}

NetworkMessage* deserialize_message(const char* buffer, int buffer_size) {
    if (!buffer || buffer_size <= 0) {
        return NULL;
    }

    NetworkMessage* message = (NetworkMessage*)malloc(sizeof(NetworkMessage));
    if (!message) {
        return NULL;
    }

    // 反序列化数据
    const char* ptr = buffer;
    int len;

    // 消息ID
    memcpy(&len, ptr, sizeof(int));
    ptr += sizeof(int);
    if (len > 0 && len < 64) {
        memcpy(message->id, ptr, len);
        message->id[len] = '\0';
    } else {
        strcpy(message->id, "unknown");
    }
    ptr += len;

    // 消息类型
    memcpy(&len, ptr, sizeof(int));
    ptr += sizeof(int);
    if (len > 0 && len < 64) {
        memcpy(message->type, ptr, len);
        message->type[len] = '\0';
    } else {
        strcpy(message->type, "unknown");
    }
    ptr += len;

    // 源
    memcpy(&len, ptr, sizeof(int));
    ptr += sizeof(int);
    if (len > 0 && len < 64) {
        memcpy(message->source, ptr, len);
        message->source[len] = '\0';
    } else {
        strcpy(message->source, "unknown");
    }
    ptr += len;

    // 目标
    memcpy(&len, ptr, sizeof(int));
    ptr += sizeof(int);
    if (len > 0 && len < 64) {
        memcpy(message->destination, ptr, len);
        message->destination[len] = '\0';
    } else {
        strcpy(message->destination, "unknown");
    }
    ptr += len;

    // 负载大小
    memcpy(&message->payload_size, ptr, sizeof(int));
    ptr += sizeof(int);

    // 负载
    if (message->payload_size > 0) {
        message->payload = (char*)malloc(message->payload_size + 1);
        if (message->payload) {
            memcpy(message->payload, ptr, message->payload_size);
            message->payload[message->payload_size] = '\0';
        } else {
            message->payload_size = 0;
        }
    } else {
        message->payload = NULL;
    }
    ptr += message->payload_size;

    // 时间戳
    memcpy(&message->timestamp, ptr, sizeof(long));

    return message;
}

int get_network_status(NetworkStatus* status) {
    if (!status) {
        return -1;
    }

    // 更新运行时间
    if (g_initialized) {
        g_status.uptime = time(NULL) - g_start_time;
    }

    memcpy(status, &g_status, sizeof(NetworkStatus));
    return 0;
}

int log_network_event(const char* format, ...) {
    va_list args;
    va_start(args, format);

    // 打印到控制台
    vprintf(format, args);

    // 写入日志文件
    if (g_log_file) {
        time_t now = time(NULL);
        struct tm* tm_info = localtime(&now);
        char timestamp[64];
        strftime(timestamp, sizeof(timestamp), "%Y-%m-%d %H:%M:%S", tm_info);

        fprintf(g_log_file, "[%s] ", timestamp);
        vfprintf(g_log_file, format, args);
        fflush(g_log_file);
    }

    va_end(args);
    return 0;
}
