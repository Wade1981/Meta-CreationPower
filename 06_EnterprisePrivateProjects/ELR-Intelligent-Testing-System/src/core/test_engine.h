#ifndef TEST_ENGINE_H
#define TEST_ENGINE_H

#include <string>
#include <vector>
#include <map>
#include <memory>

namespace elr_test {

// 测试状态枚举
enum class TestStatus {
    PENDING,
    RUNNING,
    PASSED,
    FAILED,
    SKIPPED
};

// 测试用例基类
class TestCase {
public:
    TestCase(const std::string& name, const std::string& description);
    virtual ~TestCase() = default;

    // 执行测试
    virtual TestStatus execute() = 0;

    // 获取测试结果
    const std::string& get_name() const { return name_; }
    const std::string& get_description() const { return description_; }
    TestStatus get_status() const { return status_; }
    const std::string& get_message() const { return message_; }
    double get_duration() const { return duration_; }

protected:
    std::string name_;
    std::string description_;
    TestStatus status_;
    std::string message_;
    double duration_; // 测试执行时间（秒）
};

// 测试套件类
class TestSuite {
public:
    TestSuite(const std::string& name);

    // 添加测试用例
    void add_test_case(std::unique_ptr<TestCase> test_case);

    // 执行所有测试
    void execute();

    // 获取测试结果
    const std::string& get_name() const { return name_; }
    const std::vector<std::unique_ptr<TestCase>>& get_test_cases() const { return test_cases_; }
    int get_total_tests() const { return test_cases_.size(); }
    int get_passed_tests() const; 
    int get_failed_tests() const;
    int get_skipped_tests() const;
    double get_total_duration() const;

private:
    std::string name_;
    std::vector<std::unique_ptr<TestCase>> test_cases_;
};

// 测试引擎类
class TestEngine {
public:
    TestEngine();

    // 添加测试套件
    void add_test_suite(std::unique_ptr<TestSuite> test_suite);

    // 执行所有测试
    void run();

    // 获取测试结果
    const std::vector<std::unique_ptr<TestSuite>>& get_test_suites() const { return test_suites_; }
    int get_total_tests() const;
    int get_passed_tests() const;
    int get_failed_tests() const;
    int get_skipped_tests() const;
    double get_total_duration() const;

private:
    std::vector<std::unique_ptr<TestSuite>> test_suites_;
};

} // namespace elr_test

#endif // TEST_ENGINE_H