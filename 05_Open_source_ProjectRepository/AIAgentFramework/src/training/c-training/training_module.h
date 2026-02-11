// training_module.h - AIAgent训练模块头文件 (跨平台版本)

#ifndef TRAINING_MODULE_H
#define TRAINING_MODULE_H

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <time.h>
#include <stdarg.h>

// 跨平台头文件
#ifdef _WIN32
    #include <windows.h>
    #include <process.h>
    #define THREAD_CREATE(thread, start, arg) _beginthreadex(NULL, 0, start, arg, 0, thread)
    #define THREAD_JOIN(thread) WaitForSingleObject((HANDLE)thread, INFINITE)
    #define THREAD_T unsigned int
#else
    #include <pthread.h>
    #include <unistd.h>
    #include <sys/types.h>
    #include <sys/wait.h>
    #define THREAD_CREATE(thread, start, arg) pthread_create(thread, NULL, start, arg)
    #define THREAD_JOIN(thread) pthread_join(*thread, NULL)
    #define THREAD_T pthread_t*
#endif

// 训练状态枚举
typedef enum {
    TRAINING_STATUS_PENDING,
    TRAINING_STATUS_RUNNING,
    TRAINING_STATUS_COMPLETED,
    TRAINING_STATUS_FAILED,
    TRAINING_STATUS_STOPPED
} TrainingStatus;

// AI框架类型
typedef enum {
    AI_FRAMEWORK_TENSORFLOW,
    AI_FRAMEWORK_PYTORCH,
    AI_FRAMEWORK_TENSORFLOW_LITE,
    AI_FRAMEWORK_ONNX_RUNTIME
} AIFrameworkType;

// 训练配置结构体
typedef struct {
    char model_name[256];
    char dataset_id[256];
    char output_dir[512];
    int epochs;
    int batch_size;
    float learning_rate;
    char optimizer[128];
    char loss_function[128];
    AIFrameworkType framework;
    int auto_deploy;
    char framework_version[64];
    char additional_params[1024];
} TrainingConfig;

// 训练作业结构体
typedef struct {
    char id[256];
    TrainingConfig config;
    TrainingStatus status;
    time_t start_time;
    time_t end_time;
    float accuracy;
    float loss;
    float precision;
    float recall;
    float f1_score;
    char model_path[512];
    char logs[4096];
    int log_size;
    THREAD_T thread;
    int running;
} TrainingJob;

// 模型结构体
typedef struct {
    char id[256];
    char name[256];
    char path[512];
    char framework[128];
    char version[64];
    time_t created_at;
    time_t deployed_at;
    char endpoint[512];
    int is_deployed;
} Model;

// 训练模块配置
typedef struct {
    char models_directory[512];
    char datasets_directory[512];
    char logs_directory[512];
    int max_concurrent_jobs;
    int max_memory_mb;
    int use_gpu;
    int gpu_memory_mb;
    char python_executable[256];
} TrainingModuleConfig;

// 初始化和清理函数
int training_module_init(TrainingModuleConfig* config);
int training_module_cleanup();

// 训练作业管理函数
int start_training(const char* job_id, TrainingConfig* config);
int stop_training(const char* job_id);
TrainingJob* get_training_job(const char* job_id);
int list_training_jobs(TrainingJob*** jobs, int* count);

// 模型管理函数
int register_model(const char* model_id, const char* model_path, const char* framework);
int deploy_model(const char* model_id, const char* endpoint);
int undeploy_model(const char* model_id);
Model* get_model(const char* model_id);
int list_models(Model*** models, int* count);

// 评估函数
int evaluate_model(const char* model_id, const char* dataset_id, float* accuracy, float* loss);

// 资源管理函数
int get_resource_usage(float* cpu, float* memory, float* gpu);
int set_resource_limits(int max_memory_mb, int use_gpu, int gpu_memory_mb);

// 日志函数
int log_training_event(const char* format, ...);
int get_training_logs(const char* job_id, char** logs, int* log_size);

// 框架管理函数
int check_framework_availability(AIFrameworkType framework, int* available, char* version);
int install_framework(AIFrameworkType framework, const char* version);

// 算法相关结构和函数
typedef enum {
    ALGORITHM_TYPE_OPTIMIZER,
    ALGORITHM_TYPE_LOSS,
    ALGORITHM_TYPE_COMPRESSION,
    ALGORITHM_TYPE_DATA_AUGMENTATION,
    ALGORITHM_TYPE_OTHER
} AlgorithmType;

typedef struct {
    char id[256];
    char name[256];
    AlgorithmType type;
    char parameters[2048];
    float performance_score;
    time_t created_at;
} Algorithm;

// 高性能算法函数
int register_algorithm(const char* algorithm_id, const char* name, AlgorithmType type, const char* parameters);
int run_algorithm(const char* algorithm_id, void* input_data, void* output_data, float* execution_time);
int optimize_training_with_algorithm(const char* job_id, const char* algorithm_id);
Algorithm* get_algorithm(const char* algorithm_id);
int list_algorithms(Algorithm*** algorithms, int* count);

// AdamW优化器参数
typedef struct {
    float learning_rate;
    float beta1;
    float beta2;
    float epsilon;
    float weight_decay;
    int use_amsgrad;
} AdamWParams;

// AdamW优化器函数
int adamw_optimizer_init(AdamWParams* params);
int adamw_optimizer_update(float* weights, float* gradients, int weight_count, AdamWParams* params, int step);

// 常量定义
#define MAX_TRAINING_JOBS 100
#define MAX_MODELS 500
#define MAX_LOG_SIZE 4096
#define MAX_COMMAND_LENGTH 2048

#endif // TRAINING_MODULE_H
