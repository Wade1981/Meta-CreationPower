#ifndef ELR_CLIENT_H
#define ELR_CLIENT_H

#include <string>
#include <vector>
#include <map>

namespace elr_test {

// ELR 客户端类
class ELRClient {
public:
    ELRClient(const std::string& base_url);
    ~ELRClient() = default;

    // 发送 HTTP 请求
    std::string send_request(const std::string& method, const std::string& path, const std::string& body = "");

    // ELR 服务管理
    bool start_service();
    bool stop_service();
    bool restart_service();
    std::string get_service_status();

    // ELR 容器管理
    std::string create_container(const std::string& name, const std::string& image, 
                                const std::map<std::string, std::string>& environment = {});
    bool start_container(const std::string& container_id);
    bool stop_container(const std::string& container_id);
    bool delete_container(const std::string& container_id);
    std::string get_container_status(const std::string& container_id);
    std::vector<std::map<std::string, std::string>> list_containers();

    // ELR API 调用
    std::string call_api(const std::string& endpoint, const std::string& method = "GET", 
                        const std::string& body = "");

private:
    std::string base_url_;

    // 辅助方法
    std::string build_url(const std::string& path);
    std::string handle_response(const void* response);
};

} // namespace elr_test

#endif // ELR_CLIENT_H