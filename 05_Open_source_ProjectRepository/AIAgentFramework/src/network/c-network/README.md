# AIAgentFramework 跨平台C语言网络模块

## 概述

本目录包含AIAgentFramework的跨平台C语言网络模块实现，支持Windows、Linux和macOS等操作系统。该模块提供了高性能的网络通信功能，包括TCP/UDP通信、消息序列化/反序列化、服务器/客户端实现等。

## 目录结构

```
c-network/
├── network_module.h    # 跨平台网络模块头文件
├── network_module.c    # 跨平台网络模块实现
└── README.md           # 本指南
```

## 功能特性

- **跨平台支持**：支持Windows、Linux、macOS等操作系统
- **TCP通信**：可靠的面向连接通信
- **UDP通信**：无连接的数据包通信
- **消息序列化**：支持结构化消息的序列化和反序列化
- **服务器功能**：支持创建和管理TCP服务器
- **客户端功能**：支持创建和管理TCP客户端
- **错误处理**：全面的错误检查和报告
- **日志系统**：详细的事件日志记录

## 构建指南

### 前提条件

- **编译器**：支持C99标准的编译器
  - Windows: Visual Studio, MinGW
  - Linux: GCC
  - macOS: Clang
- **网络库**：
  - Windows: Winsock2 (内置)
  - Linux/macOS: POSIX sockets (内置)

### 编译命令

#### Windows (Visual Studio)

```bash
# 使用Developer Command Prompt for Visual Studio
cl /c network_module.c /I.
link network_module.obj /OUT:network_module.lib
```

#### Windows (MinGW)

```bash
gcc -c network_module.c -I.
gcc -shared -o network_module.dll network_module.o -lws2_32
```

#### Linux

```bash
gcc -c network_module.c -I.
gcc -shared -o libnetwork_module.so network_module.o
```

#### macOS

```bash
clang -c network_module.c -I.
clang -shared -o libnetwork_module.dylib network_module.o
```

### 静态库构建

#### Windows

```bash
# Visual Studio
lib network_module.obj /OUT:network_module.lib

# MinGW
gcc -c network_module.c -I.
ar rcs libnetwork_module.a network_module.o
```

#### Linux/macOS

```bash
gcc -c network_module.c -I.
ar rcs libnetwork_module.a network_module.o
```

## 使用指南

### 基本用法

```c
#include "network_module.h"

int main() {
    // 初始化网络模块
    NetworkConfig config;
    config.max_connections = 100;
    config.buffer_size = MAX_BUFFER_SIZE;
    config.timeout = 30000;
    config.keepalive_interval = 60000;
    strcpy(config.log_file, "network.log");
    
    if (network_module_init(&config) != 0) {
        printf("Failed to initialize network module\n");
        return 1;
    }
    
    // 创建服务器
    int server_port = 8080;
    SOCKET_T server_socket = create_server(server_port, 5);
    if (server_socket == INVALID_SOCKET_VALUE) {
        printf("Failed to create server\n");
        network_module_cleanup();
        return 1;
    }
    
    printf("Server started on port %d\n", server_port);
    
    // 接受连接
    while (1) {
        NetworkNode* client_node = accept_connection(server_socket);
        if (client_node) {
            printf("Accepted connection from %s:%d\n", 
                   client_node->hostname, client_node->port);
            
            // 处理客户端消息
            NetworkMessage* message = receive_message(client_node);
            if (message) {
                printf("Received message: %s\n", message->payload);
                
                // 发送响应
                NetworkMessage* response = create_message(
                    "response", 
                    "server", 
                    message->source, 
                    "Hello from server"
                );
                if (response) {
                    send_message(client_node, response);
                    destroy_message(response);
                }
                
                destroy_message(message);
            }
            
            destroy_network_node(client_node);
        }
    }
    
    // 清理
    close_server(server_socket);
    network_module_cleanup();
    
    return 0;
}
```

### 客户端示例

```c
#include "network_module.h"

int main() {
    // 初始化网络模块
    if (network_module_init(NULL) != 0) {
        printf("Failed to initialize network module\n");
        return 1;
    }
    
    // 创建客户端节点
    NetworkNode* client_node = create_network_node(
        "client1", 
        "localhost", 
        8080
    );
    
    if (!client_node) {
        printf("Failed to create client node\n");
        network_module_cleanup();
        return 1;
    }
    
    // 连接到服务器
    if (connect_node(client_node) != 0) {
        printf("Failed to connect to server\n");
        destroy_network_node(client_node);
        network_module_cleanup();
        return 1;
    }
    
    // 发送消息
    NetworkMessage* message = create_message(
        "request", 
        "client1", 
        "server", 
        "Hello from client"
    );
    
    if (message) {
        if (send_message(client_node, message) == 0) {
            printf("Message sent successfully\n");
            
            // 接收响应
            NetworkMessage* response = receive_message(client_node);
            if (response) {
                printf("Received response: %s\n", response->payload);
                destroy_message(response);
            }
        }
        destroy_message(message);
    }
    
    // 清理
    disconnect_node(client_node);
    destroy_network_node(client_node);
    network_module_cleanup();
    
    return 0;
}
```

### UDP示例

```c
#include "network_module.h"

int main() {
    // 初始化网络模块
    if (network_module_init(NULL) != 0) {
        printf("Failed to initialize network module\n");
        return 1;
    }
    
    // 创建UDP套接字
    int udp_socket = create_udp_socket();
    if (udp_socket == INVALID_SOCKET_VALUE) {
        printf("Failed to create UDP socket\n");
        network_module_cleanup();
        return 1;
    }
    
    // 发送UDP消息
    const char* message = "Hello UDP";
    int bytes_sent = send_udp_message(
        udp_socket, 
        "localhost", 
        8081, 
        message, 
        strlen(message)
    );
    
    if (bytes_sent > 0) {
        printf("Sent %d bytes via UDP\n", bytes_sent);
    }
    
    // 接收UDP消息
    char buffer[1024];
    struct sockaddr_in client_addr;
    socklen_t addr_len = sizeof(client_addr);
    
    int bytes_received = receive_udp_message(
        udp_socket, 
        buffer, 
        sizeof(buffer), 
        &client_addr, 
        &addr_len
    );
    
    if (bytes_received > 0) {
        buffer[bytes_received] = '\0';
        char client_ip[INET_ADDRSTRLEN];
        inet_ntop(AF_INET, &(client_addr.sin_addr), client_ip, INET_ADDRSTRLEN);
        printf("Received %d bytes from %s:%d: %s\n", 
               bytes_received, client_ip, ntohs(client_addr.sin_port), buffer);
    }
    
    // 清理
    CLOSE_SOCKET(udp_socket);
    network_module_cleanup();
    
    return 0;
}
```

## 部署指南

### Windows

1. **复制库文件**：
   - 将编译生成的 `network_module.dll` 或 `network_module.lib` 复制到应用程序目录

2. **运行时依赖**：
   - Windows XP及以上：无需额外依赖
   - 确保应用程序以管理员权限运行（如果需要绑定低端口）

### Linux

1. **复制库文件**：
   ```bash
   # 系统-wide安装
   sudo cp libnetwork_module.so /usr/lib/
   sudo ldconfig
   
   # 或本地安装
   cp libnetwork_module.so /path/to/application/
   ```

2. **运行时设置**：
   ```bash
   # 设置库路径
   export LD_LIBRARY_PATH=/path/to/application:$LD_LIBRARY_PATH
   
   # 运行应用程序
   ./your_application
   ```

### macOS

1. **复制库文件**：
   ```bash
   # 系统-wide安装
   sudo cp libnetwork_module.dylib /usr/local/lib/
   
   # 或本地安装
   cp libnetwork_module.dylib /path/to/application/
   ```

2. **运行时设置**：
   ```bash
   # 设置库路径
   export DYLD_LIBRARY_PATH=/path/to/application:$DYLD_LIBRARY_PATH
   
   # 运行应用程序
   ./your_application
   ```

## 故障排除

### 常见问题

1. **Windows编译错误**：
   - 错误：`undefined reference to Winsock functions`
   - 解决：链接时添加 `-lws2_32` 库

2. **Linux连接错误**：
   - 错误：`Permission denied` 当绑定低端口时
   - 解决：使用sudo运行或绑定高于1024的端口

3. **macOS动态库加载错误**：
   - 错误：`dyld: Library not loaded`
   - 解决：设置正确的 `DYLD_LIBRARY_PATH`

4. **网络连接失败**：
   - 检查防火墙设置
   - 确保目标主机和端口可访问
   - 检查网络配置

### 日志文件

网络模块会生成详细的日志文件，默认名为 `network_module.log`，可以在配置中指定自定义路径。日志文件包含以下信息：

- 模块初始化和清理事件
- 连接建立和断开事件
- 消息发送和接收事件
- 错误和异常情况

## API参考

### 初始化和清理

- `int network_module_init(NetworkConfig* config)` - 初始化网络模块
- `int network_module_cleanup()` - 清理网络模块

### 节点管理

- `NetworkNode* create_network_node(const char* id, const char* hostname, int port)` - 创建网络节点
- `int destroy_network_node(NetworkNode* node)` - 销毁网络节点
- `int connect_node(NetworkNode* node)` - 连接节点
- `int disconnect_node(NetworkNode* node)` - 断开节点连接

### 消息处理

- `NetworkMessage* create_message(const char* type, const char* source, const char* destination, const char* payload)` - 创建消息
- `int destroy_message(NetworkMessage* message)` - 销毁消息
- `int send_message(NetworkNode* node, NetworkMessage* message)` - 发送消息
- `NetworkMessage* receive_message(NetworkNode* node)` - 接收消息

### 服务器功能

- `SOCKET_T create_server(int port, int backlog)` - 创建服务器
- `NetworkNode* accept_connection(SOCKET_T server_socket)` - 接受连接
- `int close_server(SOCKET_T server_socket)` - 关闭服务器

### UDP功能

- `int create_udp_socket()` - 创建UDP套接字
- `int send_udp_message(int socket, const char* host, int port, const char* message, int message_len)` - 发送UDP消息
- `int receive_udp_message(int socket, char* buffer, int buffer_size, struct sockaddr_in* client_addr, socklen_t* addr_len)` - 接收UDP消息

### 工具函数

- `int serialize_message(NetworkMessage* message, char** buffer, int* buffer_size)` - 序列化消息
- `NetworkMessage* deserialize_message(const char* buffer, int buffer_size)` - 反序列化消息
- `int get_network_status(NetworkStatus* status)` - 获取网络状态
- `int log_network_event(const char* format, ...)` - 记录网络事件

## 性能优化

1. **连接池**：对于频繁的客户端连接，考虑实现连接池
2. **非阻塞I/O**：对于高并发场景，使用非阻塞I/O
3. **线程池**：对于多客户端处理，使用线程池
4. **缓冲区大小**：根据实际需求调整缓冲区大小
5. **心跳机制**：实现心跳机制检测连接状态

## 安全考虑

1. **输入验证**：验证所有网络输入，防止缓冲区溢出
2. **加密通信**：对于敏感数据，考虑使用TLS/SSL
3. **认证机制**：实现客户端认证，防止未授权访问
4. **速率限制**：实现速率限制，防止DoS攻击
5. **防火墙**：配置适当的防火墙规则

## 版本历史

- **v1.0.0**：初始版本
  - 跨平台TCP/UDP支持
  - 消息序列化/反序列化
  - 服务器/客户端功能
  - 错误处理和日志

## 许可证

本模块采用MIT许可证，详见项目根目录的LICENSE文件。
