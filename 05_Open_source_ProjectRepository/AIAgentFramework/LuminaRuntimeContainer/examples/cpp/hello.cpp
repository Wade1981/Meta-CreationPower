// Lumina Runtime Container - C++示例

#include <iostream>
#include <string>
#include <vector>

// 主函数
int main() {
    // 打印欢迎信息
    std::cout << "====================================" << std::endl;
    std::cout << "Hello from Lumina Runtime Container!" << std::endl;
    std::cout << "====================================" << std::endl;
    std::cout << "Language: C++" << std::endl;
    std::cout << "Container Version: 1.0.0" << std::endl;
    std::cout << "====================================" << std::endl;
    
    // 演示基本功能
    std::cout << "\nBasic C++ Features:" << std::endl;
    
    // 变量和输出
    int number = 42;
    std::string message = "The answer to life, the universe, and everything";
    std::cout << "Integer variable: " << number << std::endl;
    std::cout << "String variable: " << message << std::endl;
    
    // 循环
    std::cout << "\nLoop demonstration:" << std::endl;
    for (int i = 1; i <= 5; ++i) {
        std::cout << "Iteration " << i << std::endl;
    }
    
    // 向量
    std::cout << "\nVector demonstration:" << std::endl;
    std::vector<std::string> languages = {"C++", "Python", "Java", "JavaScript", "Go"};
    for (const auto& lang : languages) {
        std::cout << "Supported language: " << lang << std::endl;
    }
    
    // 函数调用
    std::cout << "\nFunction demonstration:" << std::endl;
    
    // 计算平方的函数
    auto square = [](int x) {
        return x * x;
    };
    
    for (int i = 1; i <= 3; ++i) {
        std::cout << "Square of " << i << " is " << square(i) << std::endl;
    }
    
    // 结束信息
    std::cout << "====================================" << std::endl;
    std::cout << "C++ example completed successfully!" << std::endl;
    std::cout << "====================================" << std::endl;
    
    return 0;
}
