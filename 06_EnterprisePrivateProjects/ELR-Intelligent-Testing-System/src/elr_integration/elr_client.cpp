#include "elr_client.h"
#include <curl/curl.h>
#include <sstream>
#include <iostream>

namespace elr_test {

// 回调函数用于获取 HTTP 响应
static size_t write_callback(void* contents, size_t size, size_t nmemb, void* userp) {
    ((std::string*)userp)->append((char*)contents, size * nmemb);
    return size * nmemb;
}

// ELRClient 实现
ELRClient::ELRClient(const std::string& base_url) : base_url_(base_url) {
    curl_global_init(CURL_GLOBAL_DEFAULT);
}

std::string ELRClient::build_url(const std::string& path) {
    std::string url = base_url_;
    if (url.back() != '/') {
        url += '/';
    }
    if (!path.empty() && path.front() == '/') {
        url += path.substr(1);
    } else {
        url += path;
    }
    return url;
}

std::string ELRClient::send_request(const std::string& method, const std::string& path, const std::string& body) {
    CURL* curl = curl_easy_init();
    if (!curl) {
        return "Error: Failed to initialize curl";
    }

    std::string url = build_url(path);
    std::string response_string;

    curl_easy_setopt(curl, CURLOPT_URL, url.c_str());
    curl_easy_setopt(curl, CURLOPT_WRITEFUNCTION, write_callback);
    curl_easy_setopt(curl, CURLOPT_WRITEDATA, &response_string);
    curl_easy_setopt(curl, CURLOPT_FOLLOWLOCATION, 1L);

    // 设置请求方法
    if (method == "POST") {
        curl_easy_setopt(curl, CURLOPT_POST, 1L);
        curl_easy_setopt(curl, CURLOPT_POSTFIELDS, body.c_str());
        curl_easy_setopt(curl, CURLOPT_POSTFIELDSIZE, body.length());
    } else if (method == "PUT") {
        curl_easy_setopt(curl, CURLOPT_CUSTOMREQUEST, "PUT");
        curl_easy_setopt(curl, CURLOPT_POSTFIELDS, body.c_str());
        curl_easy_setopt(curl, CURLOPT_POSTFIELDSIZE, body.length());
    } else if (method == "DELETE") {
        curl_easy_setopt(curl, CURLOPT_CUSTOMREQUEST, "DELETE");
    }

    // 执行请求
    CURLcode res = curl_easy_perform(curl);
    if (res != CURLE_OK) {
        std::string error = "Error: " + std::string(curl_easy_strerror(res));
        curl_easy_cleanup(curl);
        return error;
    }

    curl_easy_cleanup(curl);
    return response_string;
}

bool ELRClient::start_service() {
    std::string response = send_request("POST", "service/start");
    return response.find("success") != std::string::npos;
}

bool ELRClient::stop_service() {
    std::string response = send_request("POST", "service/stop");
    return response.find("success") != std::string::npos;
}

bool ELRClient::restart_service() {
    std::string response = send_request("POST", "service/restart");
    return response.find("success") != std::string::npos;
}

std::string ELRClient::get_service_status() {
    return send_request("GET", "service/status");
}

std::string ELRClient::create_container(const std::string& name, const std::string& image, 
                                       const std::map<std::string, std::string>& environment) {
    std::stringstream body;
    body << "{\"name\": \"" << name << "\", \"image\": \"" << image << "\", \"environment\": {";
    
    bool first = true;
    for (const auto& env : environment) {
        if (!first) {
            body << ",";
        }
        body << "\"" << env.first << "\": \"" << env.second << "\"";
        first = false;
    }
    
    body << "}}";
    
    return send_request("POST", "containers", body.str());
}

bool ELRClient::start_container(const std::string& container_id) {
    std::string response = send_request("POST", "containers/" + container_id + "/start");
    return response.find("success") != std::string::npos;
}

bool ELRClient::stop_container(const std::string& container_id) {
    std::string response = send_request("POST", "containers/" + container_id + "/stop");
    return response.find("success") != std::string::npos;
}

bool ELRClient::delete_container(const std::string& container_id) {
    std::string response = send_request("DELETE", "containers/" + container_id);
    return response.find("success") != std::string::npos;
}

std::string ELRClient::get_container_status(const std::string& container_id) {
    return send_request("GET", "containers/" + container_id + "/status");
}

std::vector<std::map<std::string, std::string>> ELRClient::list_containers() {
    // 简化实现，实际应该解析 JSON 响应
    std::vector<std::map<std::string, std::string>> containers;
    std::string response = send_request("GET", "containers");
    // 这里应该解析 JSON 响应并填充 containers 向量
    return containers;
}

std::string ELRClient::call_api(const std::string& endpoint, const std::string& method, const std::string& body) {
    return send_request(method, "api/" + endpoint, body);
}

std::string ELRClient::handle_response(const void* response) {
    // 简化实现，实际应该根据响应类型进行处理
    return ""; 
}

} // namespace elr_test