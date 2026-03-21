#include "test_engine.h"
#include <chrono>
#include <iostream>

namespace elr_test {

// TestCase 实现
TestCase::TestCase(const std::string& name, const std::string& description)
    : name_(name), description_(description), status_(TestStatus::PENDING), duration_(0.0) {}

// TestSuite 实现
TestSuite::TestSuite(const std::string& name) : name_(name) {}

void TestSuite::add_test_case(std::unique_ptr<TestCase> test_case) {
    test_cases_.push_back(std::move(test_case));
}

void TestSuite::execute() {
    std::cout << "Running test suite: " << name_ << std::endl;
    for (auto& test_case : test_cases_) {
        std::cout << "  Running test: " << test_case->get_name() << std::endl;
        test_case->execute();
        
        switch (test_case->get_status()) {
        case TestStatus::PASSED:
            std::cout << "  ✓ PASSED" << std::endl;
            break;
        case TestStatus::FAILED:
            std::cout << "  ✗ FAILED: " << test_case->get_message() << std::endl;
            break;
        case TestStatus::SKIPPED:
            std::cout << "  ⚠ SKIPPED" << std::endl;
            break;
        default:
            break;
        }
    }
    std::cout << "Test suite " << name_ << " completed" << std::endl;
}

int TestSuite::get_passed_tests() const {
    int count = 0;
    for (const auto& test_case : test_cases_) {
        if (test_case->get_status() == TestStatus::PASSED) {
            count++;
        }
    }
    return count;
}

int TestSuite::get_failed_tests() const {
    int count = 0;
    for (const auto& test_case : test_cases_) {
        if (test_case->get_status() == TestStatus::FAILED) {
            count++;
        }
    }
    return count;
}

int TestSuite::get_skipped_tests() const {
    int count = 0;
    for (const auto& test_case : test_cases_) {
        if (test_case->get_status() == TestStatus::SKIPPED) {
            count++;
        }
    }
    return count;
}

double TestSuite::get_total_duration() const {
    double total = 0.0;
    for (const auto& test_case : test_cases_) {
        total += test_case->get_duration();
    }
    return total;
}

// TestEngine 实现
TestEngine::TestEngine() {}

void TestEngine::add_test_suite(std::unique_ptr<TestSuite> test_suite) {
    test_suites_.push_back(std::move(test_suite));
}

void TestEngine::run() {
    std::cout << "Starting test engine..." << std::endl;
    auto start_time = std::chrono::high_resolution_clock::now();
    
    for (auto& test_suite : test_suites_) {
        test_suite->execute();
    }
    
    auto end_time = std::chrono::high_resolution_clock::now();
    double total_duration = std::chrono::duration<double>(end_time - start_time).count();
    
    std::cout << "\nTest engine completed" << std::endl;
    std::cout << "Total tests: " << get_total_tests() << std::endl;
    std::cout << "Passed: " << get_passed_tests() << std::endl;
    std::cout << "Failed: " << get_failed_tests() << std::endl;
    std::cout << "Skipped: " << get_skipped_tests() << std::endl;
    std::cout << "Total duration: " << total_duration << " seconds" << std::endl;
}

int TestEngine::get_total_tests() const {
    int count = 0;
    for (const auto& test_suite : test_suites_) {
        count += test_suite->get_total_tests();
    }
    return count;
}

int TestEngine::get_passed_tests() const {
    int count = 0;
    for (const auto& test_suite : test_suites_) {
        count += test_suite->get_passed_tests();
    }
    return count;
}

int TestEngine::get_failed_tests() const {
    int count = 0;
    for (const auto& test_suite : test_suites_) {
        count += test_suite->get_failed_tests();
    }
    return count;
}

int TestEngine::get_skipped_tests() const {
    int count = 0;
    for (const auto& test_suite : test_suites_) {
        count += test_suite->get_skipped_tests();
    }
    return count;
}

double TestEngine::get_total_duration() const {
    double total = 0.0;
    for (const auto& test_suite : test_suites_) {
        total += test_suite->get_total_duration();
    }
    return total;
}

} // namespace elr_test