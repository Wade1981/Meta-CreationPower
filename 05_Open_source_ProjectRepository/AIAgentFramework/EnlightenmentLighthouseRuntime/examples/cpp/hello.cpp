// Hello World example for Enlightenment Lighthouse Runtime (ELR)
// This example demonstrates how to run a C++ application in ELR

#include <iostream>
#include <string>
#include <vector>

int main() {
    // Print welcome message
    std::cout << "====================================" << std::endl;
    std::cout << "Hello from Enlightenment Lighthouse Runtime!" << std::endl;
    std::cout << "====================================" << std::endl;
    std::cout << "Language: C++" << std::endl;
    std::cout << "Runtime: ELR" << std::endl;
    std::cout << "====================================" << std::endl;
    
    // Demonstrate basic C++ features
    std::cout << "\nBasic C++ Features:" << std::endl;
    
    // Variables and output
    int number = 42;
    std::string message = "The answer to life, the universe, and everything";
    std::cout << "Integer variable: " << number << std::endl;
    std::cout << "String variable: " << message << std::endl;
    
    // Loop
    std::cout << "\nLoop demonstration:" << std::endl;
    for (int i = 1; i <= 5; ++i) {
        std::cout << "Iteration " << i << std::endl;
    }
    
    // Vector
    std::cout << "\nVector demonstration:" << std::endl;
    std::vector<std::string> languages = {"C++", "Python", "Java", "JavaScript", "Go"};
    for (const auto& lang : languages) {
        std::cout << "Supported language: " << lang << std::endl;
    }
    
    // Function
    std::cout << "\nFunction demonstration:" << std::endl;
    
    // Lambda function to calculate square
    auto square = [](int x) {
        return x * x;
    };
    
    for (int i = 1; i <= 3; ++i) {
        std::cout << "Square of " << i << " is " << square(i) << std::endl;
    }
    
    // Environment variables
    std::cout << "\nEnvironment variables:" << std::endl;
    if (const char* env_p = std::getenv("ELR_CONTAINER_ID")) {
        std::cout << "ELR_CONTAINER_ID: " << env_p << std::endl;
    } else {
        std::cout << "ELR_CONTAINER_ID: Not set (running outside ELR)" << std::endl;
    }
    
    // End message
    std::cout << "\n====================================" << std::endl;
    std::cout << "C++ example completed successfully!" << std::endl;
    std::cout << "====================================" << std::endl;
    
    return 0;
}
