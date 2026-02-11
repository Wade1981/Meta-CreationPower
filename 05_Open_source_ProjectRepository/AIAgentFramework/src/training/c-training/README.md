# AIAgentFramework 跨平台C语言训练模块

本目录包含AIAgentFramework的跨平台C语言训练模块实现，支持Windows、Linux和macOS等操作系统。该模块提供了高性能的AI模型训练功能，支持TensorFlow、PyTorch等主流AI框架。

## 目录结构

```
c-training/
├── training_module.h    # 跨平台训练模块头文件
├── training_module.c    # 跨平台训练模块实现
└── README.md           # 本指南
```

## 功能特性

- **跨平台支持**：支持Windows、Linux、macOS等操作系统
- **AI框架集成**：支持TensorFlow、PyTorch、TensorFlow Lite、ONNX Runtime
- **训练作业管理**：创建、启动、停止和监控训练作业
- **模型管理**：注册、部署、卸载模型
- **资源管理**：监控和限制CPU、内存、GPU使用
- **训练调度**：支持并发训练作业
- **自动部署**：训练完成后自动部署模型
- **详细日志**：训练过程和结果的详细日志
- **算法管理**：注册、运行和管理高性能算法
- **AdamW优化器**：内置高性能AdamW优化器算法
- **算法优化**：使用算法优化训练过程
- **性能评估**：算法执行性能评估和评分

## 构建指南

### 前提条件

- **编译器**：支持C99标准的编译器
  - Windows: Visual Studio, MinGW
  - Linux: GCC
  - macOS: Clang
- **Python**：版本 3.6 或更高（用于运行AI框架）
- **AI框架**（可选）：
  - TensorFlow
  - PyTorch
  - TensorFlow Lite
  - ONNX Runtime

### 编译命令

#### Windows (Visual Studio)

```bash
# 使用Developer Command Prompt for Visual Studio
cl /c training_module.c /I.
link training_module.obj /OUT:training_module.lib
```

#### Windows (MinGW)

```bash
gcc -c training_module.c -I.
gcc -shared -o training_module.dll training_module.o
```

#### Linux

```bash
gcc -c training_module.c -I.
gcc -shared -o libtraining_module.so training_module.o -lpthread
```

#### macOS

```bash
clang -c training_module.c -I.
clang -shared -o libtraining_module.dylib training_module.o -lpthread
```

### 静态库构建

#### Windows

```bash
# Visual Studio
lib training_module.obj /OUT:training_module.lib

# MinGW
gcc -c training_module.c -I.
ar rcs libtraining_module.a training_module.o
```

#### Linux/macOS

```bash
gcc -c training_module.c -I.
ar rcs libtraining_module.a training_module.o
```

## 使用指南

### 基本用法

```c
#include "training_module.h"

int main() {
    // 初始化训练模块
    TrainingModuleConfig config;
    strcpy(config.models_directory, "./models");
    strcpy(config.datasets_directory, "./datasets");
    strcpy(config.logs_directory, "./logs");
    config.max_concurrent_jobs = 5;
    config.max_memory_mb = 4096;
    config.use_gpu = 1;
    config.gpu_memory_mb = 2048;
    strcpy(config.python_executable, "python3");
    
    if (training_module_init(&config) != 0) {
        printf("Failed to initialize training module\n");
        return 1;
    }
    
    // 创建训练配置
    TrainingConfig training_config;
    strcpy(training_config.model_name, "my_model");
    strcpy(training_config.dataset_id, "dataset_1");
    strcpy(training_config.output_dir, "./models/my_model");
    training_config.epochs = 10;
    training_config.batch_size = 32;
    training_config.learning_rate = 0.001;
    strcpy(training_config.optimizer, "adam");
    strcpy(training_config.loss_function, "categorical_crossentropy");
    training_config.framework = AI_FRAMEWORK_TENSORFLOW;
    training_config.auto_deploy = 1;
    strcpy(training_config.framework_version, "2.10.0");
    strcpy(training_config.additional_params, "");
    
    // 开始训练
    const char* job_id = "train_job_1";
    if (start_training(job_id, &training_config) != 0) {
        printf("Failed to start training\n");
        training_module_cleanup();
        return 1;
    }
    
    printf("Training started with job ID: %s\n", job_id);
    
    // 等待训练完成（实际应用中可能需要异步处理）
    Sleep(60000); // 等待60秒
    
    // 获取训练状态
    TrainingJob* job = get_training_job(job_id);
    if (job) {
        printf("Training status: %d\n", job->status);
        printf("Accuracy: %.4f\n", job->accuracy);
        printf("Loss: %.4f\n", job->loss);
        printf("Model path: %s\n", job->model_path);
    }
    
    // 清理
    training_module_cleanup();
    
    return 0;
}
```

### 模型评估示例

```c
#include "training_module.h"

int main() {
    // 初始化训练模块
    training_module_init(NULL);
    
    // 评估模型
    const char* model_id = "model_1";
    const char* dataset_id = "test_dataset";
    float accuracy, loss;
    
    if (evaluate_model(model_id, dataset_id, &accuracy, &loss) == 0) {
        printf("Model evaluation completed\n");
        printf("Accuracy: %.4f\n", accuracy);
        printf("Loss: %.4f\n", loss);
    } else {
        printf("Failed to evaluate model\n");
    }
    
    // 清理
    training_module_cleanup();
    
    return 0;
}
```

### 资源监控示例

```c
#include "training_module.h"

int main() {
    // 初始化训练模块
    training_module_init(NULL);
    
    // 获取资源使用情况
    float cpu, memory, gpu;
    if (get_resource_usage(&cpu, &memory, &gpu) == 0) {
        printf("Resource usage:\n");
        printf("CPU: %.2f%%\n", cpu);
        printf("Memory: %.2f%%\n", memory);
        printf("GPU: %.2f%%\n", gpu);
    }
    
    // 设置资源限制
    if (set_resource_limits(8192, 1, 4096) == 0) {
        printf("Resource limits updated\n");
    }
    
    // 清理
    training_module_cleanup();
    
    return 0;
}
```

### 算法管理示例

```c
#include "training_module.h"

int main() {
    // 初始化训练模块
    training_module_init(NULL);
    
    // 注册自定义算法
    register_algorithm("custom_optimizer", "Custom Optimizer", ALGORITHM_TYPE_OPTIMIZER, "learning_rate=0.0005,beta1=0.9,beta2=0.999");
    
    // 列出所有算法
    Algorithm** algorithms;
    int count;
    if (list_algorithms(&algorithms, &count) == 0) {
        printf("Available algorithms: %d\n", count);
        for (int i = 0; i < count; i++) {
            printf("- %s: %s (Type: %d)\n", algorithms[i]->id, algorithms[i]->name, algorithms[i]->type);
        }
        free(algorithms);
    }
    
    // 运行算法
    float execution_time;
    if (run_algorithm("adamw_optimizer", NULL, NULL, &execution_time) == 0) {
        printf("AdamW optimizer executed in %.4f seconds\n", execution_time);
    }
    
    // 清理
    training_module_cleanup();
    
    return 0;
}
```

### AdamW优化器使用示例

```c
#include "training_module.h"

int main() {
    // 初始化训练模块
    training_module_init(NULL);
    
    // 初始化AdamW参数
    AdamWParams params;
    params.learning_rate = 0.001;
    params.beta1 = 0.9;
    params.beta2 = 0.999;
    params.epsilon = 1e-8;
    params.weight_decay = 0.01;
    params.use_amsgrad = 0;
    
    adamw_optimizer_init(&params);
    
    // 模拟权重和梯度
    int weight_count = 1000;
    float* weights = (float*)malloc(sizeof(float) * weight_count);
    float* gradients = (float*)malloc(sizeof(float) * weight_count);
    
    // 初始化权重和梯度
    for (int i = 0; i < weight_count; i++) {
        weights[i] = (float)rand() / RAND_MAX;
        gradients[i] = (float)rand() / RAND_MAX - 0.5;
    }
    
    // 执行优化器更新
    for (int step = 1; step <= 10; step++) {
        adamw_optimizer_update(weights, gradients, weight_count, &params, step);
        printf("Step %d completed\n", step);
    }
    
    // 释放内存
    free(weights);
    free(gradients);
    
    // 清理
    training_module_cleanup();
    
    return 0;
}
```

### 使用算法优化训练示例

```c
#include "training_module.h"

int main() {
    // 初始化训练模块
    TrainingModuleConfig config;
    strcpy(config.models_directory, "./models");
    strcpy(config.datasets_directory, "./datasets");
    strcpy(config.logs_directory, "./logs");
    config.max_concurrent_jobs = 5;
    config.max_memory_mb = 4096;
    config.use_gpu = 1;
    config.gpu_memory_mb = 2048;
    strcpy(config.python_executable, "python3");
    
    training_module_init(&config);
    
    // 创建训练配置
    TrainingConfig training_config;
    strcpy(training_config.model_name, "my_model");
    strcpy(training_config.dataset_id, "dataset_1");
    strcpy(training_config.output_dir, "./models/my_model");
    training_config.epochs = 10;
    training_config.batch_size = 32;
    training_config.learning_rate = 0.001;
    strcpy(training_config.optimizer, "adam");
    strcpy(training_config.loss_function, "categorical_crossentropy");
    training_config.framework = AI_FRAMEWORK_TENSORFLOW;
    training_config.auto_deploy = 1;
    strcpy(training_config.framework_version, "2.10.0");
    strcpy(training_config.additional_params, "");
    
    // 开始训练
    const char* job_id = "train_job_1";
    start_training(job_id, &training_config);
    
    // 使用AdamW优化器算法优化训练
    optimize_training_with_algorithm(job_id, "adamw_optimizer");
    
    printf("Training started with job ID: %s\n", job_id);
    printf("Optimized with AdamW optimizer\n");
    
    // 等待训练完成
    Sleep(60000);
    
    // 清理
    training_module_cleanup();
    
    return 0;
}
```

## 部署指南

### Windows

1. **复制库文件**：
   - 将编译生成的 `training_module.dll` 或 `training_module.lib` 复制到应用程序目录

2. **运行时依赖**：
   - Windows XP及以上：无需额外依赖
   - 确保已安装Python 3.6+
   - 可选：安装所需的AI框架

### Linux

1. **复制库文件**：
   ```bash
   # 系统-wide安装
   sudo cp libtraining_module.so /usr/lib/
   sudo ldconfig
   
   # 或本地安装
   cp libtraining_module.so /path/to/application/
   ```

2. **运行时设置**：
   ```bash
   # 设置库路径
   export LD_LIBRARY_PATH=/path/to/application:$LD_LIBRARY_PATH
   
   # 运行应用程序
   ./your_application
   ```

3. **安装依赖**：
   ```bash
   # 安装Python
   sudo apt-get install python3 python3-pip
   
   # 安装AI框架（可选）
   pip3 install tensorflow
   pip3 install torch torchvision
   ```

### macOS

1. **复制库文件**：
   ```bash
   # 系统-wide安装
   sudo cp libtraining_module.dylib /usr/local/lib/
   
   # 或本地安装
   cp libtraining_module.dylib /path/to/application/
   ```

2. **运行时设置**：
   ```bash
   # 设置库路径
   export DYLD_LIBRARY_PATH=/path/to/application:$DYLD_LIBRARY_PATH
   
   # 运行应用程序
   ./your_application
   ```

3. **安装依赖**：
   ```bash
   # 安装Python
   brew install python
   
   # 安装AI框架（可选）
   pip3 install tensorflow
   pip3 install torch torchvision
   ```

## 配置选项

### 训练模块配置

```c
TrainingModuleConfig config;
// 模型存储目录
strcpy(config.models_directory, "./models");
// 数据集存储目录
strcpy(config.datasets_directory, "./datasets");
// 日志存储目录
strcpy(config.logs_directory, "./logs");
// 最大并发训练作业数
config.max_concurrent_jobs = 5;
// 最大内存使用（MB）
config.max_memory_mb = 4096;
// 是否使用GPU
config.use_gpu = 1;
// GPU内存限制（MB）
config.gpu_memory_mb = 2048;
// Python可执行文件路径
strcpy(config.python_executable, "python3");
```

### 训练配置

```c
TrainingConfig training_config;
// 模型名称
strcpy(training_config.model_name, "my_model");
// 数据集ID
strcpy(training_config.dataset_id, "dataset_1");
// 输出目录
strcpy(training_config.output_dir, "./models/my_model");
// 训练轮数
training_config.epochs = 10;
// 批次大小
training_config.batch_size = 32;
// 学习率
training_config.learning_rate = 0.001;
// 优化器
strcpy(training_config.optimizer, "adam");
// 损失函数
strcpy(training_config.loss_function, "categorical_crossentropy");
// AI框架
training_config.framework = AI_FRAMEWORK_TENSORFLOW;
// 训练完成后自动部署
training_config.auto_deploy = 1;
// 框架版本
strcpy(training_config.framework_version, "2.10.0");
// 额外参数
strcpy(training_config.additional_params, "");
```

## 故障排除

### 常见问题

1. **Python 未找到**
   - 错误：`python3: command not found`
   - 解决：确保Python已安装并添加到系统路径

2. **AI框架未安装**
   - 错误：`ModuleNotFoundError: No module named 'tensorflow'`
   - 解决：使用pip安装所需的AI框架

3. **内存不足**
   - 错误：`MemoryError`
   - 解决：增加max_memory_mb配置或减少batch_size

4. **GPU不可用**
   - 错误：`CUDA out of memory`
   - 解决：减少gpu_memory_mb配置或禁用GPU（use_gpu=0）

5. **权限错误**
   - 错误：`Permission denied`
   - 解决：确保有写入模型和日志目录的权限

### 日志文件

训练模块会生成以下日志文件：

- **training_module.log**：模块核心日志
- **训练作业日志**：存储在logs_directory目录中

日志文件包含训练过程、错误信息和性能指标，有助于故障排除。

## 性能优化

### 训练性能优化

1. **使用GPU**：设置use_gpu=1以利用GPU加速
2. **调整批次大小**：根据可用内存调整batch_size
3. **优化学习率**：根据模型类型选择合适的learning_rate
4. **使用合适的优化器**：根据任务选择合适的optimizer

### 内存优化

1. **限制并发作业数**：根据系统资源调整max_concurrent_jobs
2. **设置内存限制**：合理设置max_memory_mb
3. **清理不再使用的模型**：及时undeploy不需要的模型

### 存储优化

1. **清理临时文件**：定期清理训练过程中产生的临时文件
2. **压缩模型**：对于部署的模型，考虑使用模型压缩技术
3. **合理组织目录**：按模型类型和版本组织存储目录

## 安全考虑

1. **输入验证**：验证所有训练配置参数
2. **数据安全**：确保训练数据的安全性和隐私保护
3. **模型安全**：保护模型文件不被未授权访问
4. **依赖安全**：定期更新AI框架和依赖库
5. **网络安全**：如果模型部署在网络上，确保网络安全

## 版本历史

- **v1.0.0**：初始版本
  - 跨平台支持（Windows、Linux、macOS）
  - TensorFlow、PyTorch集成
  - 训练作业管理
  - 模型管理和部署
  - 资源监控和限制

## 许可证

本模块采用MIT许可证，详见项目根目录的LICENSE文件。
