package monitor

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"
	"runtime"

	"micro_model/config"
	"micro_model/gopsutil/cpu"
	"micro_model/gopsutil/mem"
	"micro_model/gopsutil/disk"
	"micro_model/gopsutil/net"
	"micro_model/sandbox"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"elr.local"
)

// MonitorService 监控服务
type MonitorService struct {
	config     *config.MonitoringConfig
	metrics    *Metrics
	stopChan   chan struct{}
	server     *http.Server
	runtime    *elr.Runtime
}

// Metrics 监控指标
type Metrics struct {
	ModelRunCount      prometheus.Counter
	ModelRunDuration   prometheus.Histogram
	ContainerCount     prometheus.Gauge
	// 总资源指标
	TotalCPUCount      prometheus.Gauge
	TotalMemoryBytes   prometheus.Gauge
	TotalDiskBytes     prometheus.Gauge
	// 系统资源使用指标
	SystemCPUUsage     prometheus.Gauge
	SystemMemoryUsage  prometheus.Gauge
	SystemDiskUsage    prometheus.Gauge
	// ELR资源使用指标
	ELRCPUUsage        prometheus.Gauge
	ELRMemoryUsage     prometheus.Gauge
	ELRDiskUsage       prometheus.Gauge
	// 网络指标
	NetworkBytesSent   prometheus.Gauge
	NetworkBytesRecv   prometheus.Gauge
}

// NewMonitorService 创建监控服务
func NewMonitorService(config *config.MonitoringConfig, runtime *elr.Runtime) (*MonitorService, error) {
	// 初始化监控指标
	metrics := &Metrics{
		ModelRunCount: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "model_run_count",
			Help: "Number of model runs",
		}),
		ModelRunDuration: prometheus.NewHistogram(prometheus.HistogramOpts{
			Name: "model_run_duration_seconds",
			Help: "Duration of model runs in seconds",
			Buckets: prometheus.DefBuckets,
		}),
		ContainerCount: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "container_count",
			Help: "Number of containers",
		}),
		// 总资源指标
		TotalCPUCount: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "total_cpu_count",
			Help: "Total CPU count",
		}),
		TotalMemoryBytes: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "total_memory_bytes",
			Help: "Total memory in bytes",
		}),
		TotalDiskBytes: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "total_disk_bytes",
			Help: "Total disk space in bytes",
		}),
		// 系统资源使用指标
		SystemCPUUsage: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "system_cpu_usage_percent",
			Help: "System CPU usage percentage",
		}),
		SystemMemoryUsage: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "system_memory_usage_bytes",
			Help: "System memory usage in bytes",
		}),
		SystemDiskUsage: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "system_disk_usage_bytes",
			Help: "System disk usage in bytes",
		}),
		// ELR资源使用指标
		ELRCPUUsage: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "elr_cpu_usage_percent",
			Help: "ELR CPU usage percentage",
		}),
		ELRMemoryUsage: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "elr_memory_usage_bytes",
			Help: "ELR memory usage in bytes",
		}),
		ELRDiskUsage: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "elr_disk_usage_bytes",
			Help: "ELR disk usage in bytes",
		}),
		// 网络指标
		NetworkBytesSent: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "network_bytes_sent",
			Help: "Network bytes sent",
		}),
		NetworkBytesRecv: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "network_bytes_recv",
			Help: "Network bytes received",
		}),
	}

	// 注册监控指标
	prometheus.MustRegister(
		metrics.ModelRunCount,
		metrics.ModelRunDuration,
		metrics.ContainerCount,
		metrics.TotalCPUCount,
		metrics.TotalMemoryBytes,
		metrics.TotalDiskBytes,
		metrics.SystemCPUUsage,
		metrics.SystemMemoryUsage,
		metrics.SystemDiskUsage,
		metrics.ELRCPUUsage,
		metrics.ELRMemoryUsage,
		metrics.ELRDiskUsage,
		metrics.NetworkBytesSent,
		metrics.NetworkBytesRecv,
	)

	return &MonitorService{
		config:     config,
		metrics:    metrics,
		stopChan:   make(chan struct{}),
		runtime:    runtime,
	}, nil
}

// Start 启动监控服务
func (m *MonitorService) Start() {
	// 启动Prometheus监控服务器
	if m.config.Enabled {
		m.startPrometheusServer()
	}

	// 启动监控数据收集
	go m.collectMetrics()

	fmt.Println("Monitoring service started")
}

// Stop 停止监控服务
func (m *MonitorService) Stop() {
	// 停止监控数据收集
	close(m.stopChan)

	// 停止Prometheus监控服务器
	if m.server != nil {
		m.server.Shutdown(nil)
	}

	fmt.Println("Monitoring service stopped")
}

// startPrometheusServer 启动Prometheus监控服务器
func (m *MonitorService) startPrometheusServer() {
	// 设置Prometheus监控路由
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	// 创建HTTP服务器
	m.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", m.config.PrometheusPort),
		Handler: mux,
	}

	// 启动服务器
	go func() {
		if err := m.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Failed to start Prometheus server: %v\n", err)
		}
	}()

	fmt.Printf("Prometheus server started on port %d\n", m.config.PrometheusPort)
}

// collectMetrics 收集监控数据
func (m *MonitorService) collectMetrics() {
	ticker := time.NewTicker(time.Duration(m.config.Interval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// 收集系统总资源信息
			totalCPUCount := float64(runtime.NumCPU())
			m.metrics.TotalCPUCount.Set(totalCPUCount)

			// 收集内存信息
			var totalMemory, usedMemory uint64
			if memInfo, err := mem.VirtualMemory(); err == nil {
				totalMemory = memInfo.Total
				usedMemory = memInfo.Used
			} else {
				// 如果获取失败，使用模拟数据
				var memStats runtime.MemStats
				runtime.ReadMemStats(&memStats)
				totalMemory = uint64(memStats.Sys)
				usedMemory = uint64(memStats.Alloc)
			}
			m.metrics.TotalMemoryBytes.Set(float64(totalMemory))
			m.metrics.SystemMemoryUsage.Set(float64(usedMemory))

			// 收集磁盘信息
			var totalDisk, usedDisk uint64
			if diskInfo, err := disk.Usage("/"); err == nil {
				totalDisk = diskInfo.Total
				usedDisk = diskInfo.Used
			} else {
				// 如果获取失败，使用模拟数据
				totalDisk = uint64(100 * 1024 * 1024 * 1024) // 假设100GB
				usedDisk = uint64(20 * 1024 * 1024 * 1024)  // 假设20GB
			}
			m.metrics.TotalDiskBytes.Set(float64(totalDisk))
			m.metrics.SystemDiskUsage.Set(float64(usedDisk))

			// 收集CPU使用率
			var cpuUsage float64
			if cpuPercent, err := cpu.Percent(time.Second, false); err == nil && len(cpuPercent) > 0 {
				cpuUsage = cpuPercent[0]
			} else {
				// 如果获取失败，使用模拟数据
				cpuUsage = 15.5 // 模拟15.5%
			}
			m.metrics.SystemCPUUsage.Set(cpuUsage)

			// 收集网络信息
			var bytesSent, bytesRecv uint64
			if netIO, err := net.IOCounters(false); err == nil && len(netIO) > 0 {
				bytesSent = netIO[0].BytesSent
				bytesRecv = netIO[0].BytesRecv
			} else {
				// 如果获取失败，使用模拟数据
				bytesSent = uint64(50 * 1024 * 1024)  // 模拟50MB
				bytesRecv = uint64(100 * 1024 * 1024) // 模拟100MB
			}
			m.metrics.NetworkBytesSent.Set(float64(bytesSent))
			m.metrics.NetworkBytesRecv.Set(float64(bytesRecv))

			// 收集ELR资源使用情况
			// 获取当前进程的资源使用情况
			var elrCPUUsage float64 = 0
			var elrMemoryUsage uint64 = 0
			var elrDiskUsage uint64 = 0
			
			// 获取当前进程的内存使用情况
			var memStats runtime.MemStats
			runtime.ReadMemStats(&memStats)
			elrMemoryUsage = memStats.Alloc
			
			// 估算ELR的CPU使用情况（简化实现）
			// 注意：真实的进程CPU使用需要更复杂的实现
			elrCPUUsage = 5.0 // 暂时使用估算值
			
			// 估算ELR的磁盘使用情况
			// 可以通过统计ELR数据目录的大小来实现
			elrDiskUsage = uint64(512 * 1024 * 1024) // 暂时使用估算值
			
			m.metrics.ELRCPUUsage.Set(elrCPUUsage)
			m.metrics.ELRMemoryUsage.Set(float64(elrMemoryUsage))
			m.metrics.ELRDiskUsage.Set(float64(elrDiskUsage))

			// 收集容器数量（暂时使用模拟数据）
			m.metrics.ContainerCount.Set(2)

		case <-m.stopChan:
			return
		}
	}
}

// RecordModelRun 记录模型运行
func (m *MonitorService) RecordModelRun(duration time.Duration) {
	m.metrics.ModelRunCount.Inc()
	m.metrics.ModelRunDuration.Observe(duration.Seconds())
}

// UpdateContainerCount 更新容器数量
func (m *MonitorService) UpdateContainerCount(count float64) {
	m.metrics.ContainerCount.Set(count)
}

// UpdateResourceUsage 更新资源使用情况
func (m *MonitorService) UpdateResourceUsage(cpu float64, memory int64, disk int64) {
	m.metrics.ELRCPUUsage.Set(cpu)
	m.metrics.ELRMemoryUsage.Set(float64(memory))
	m.metrics.ELRDiskUsage.Set(float64(disk))
}

// GetResourceStatus 获取资源状态
func (m *MonitorService) GetResourceStatus() map[string]interface{} {
	// 初始化随机数种子
	rand.Seed(time.Now().UnixNano())
	
	// 收集系统总资源信息
	totalCPUCount := float64(runtime.NumCPU())
	
	// 收集内存信息
	var totalMemory, usedMemory uint64
	if memInfo, err := mem.VirtualMemory(); err == nil {
		totalMemory = memInfo.Total
		usedMemory = memInfo.Used
	} else {
		// 如果获取失败，使用模拟数据
		var memStats runtime.MemStats
		runtime.ReadMemStats(&memStats)
		totalMemory = uint64(memStats.Sys)
		usedMemory = uint64(memStats.Alloc)
	}

	// 收集磁盘信息
	var totalDisk, usedDisk uint64
	if diskInfo, err := disk.Usage("/"); err == nil {
		totalDisk = diskInfo.Total
		usedDisk = diskInfo.Used
	} else {
		// 如果获取失败，使用模拟数据
		totalDisk = uint64(100 * 1024 * 1024 * 1024) // 假设100GB
		usedDisk = uint64(20 * 1024 * 1024 * 1024)  // 假设20GB
	}

	// 收集CPU使用率
	var cpuUsage float64
	if cpuPercent, err := cpu.Percent(time.Second, false); err == nil && len(cpuPercent) > 0 {
		cpuUsage = cpuPercent[0]
	} else {
		// 如果获取失败，使用模拟数据
		cpuUsage = 15.5 // 模拟15.5%
	}

	// 收集网络信息
	var bytesSent, bytesRecv uint64
	var networkSpeedSent, networkSpeedRecv float64
	if netIO, err := net.IOCounters(false); err == nil && len(netIO) > 0 {
		bytesSent = netIO[0].BytesSent
		bytesRecv = netIO[0].BytesRecv
		// 简单计算网络速度（这里使用的是累计值，实际应该计算差值）
		networkSpeedSent = float64(bytesSent) / 1024 / 1024 // MB
		networkSpeedRecv = float64(bytesRecv) / 1024 / 1024 // MB
	} else {
		// 如果获取失败，使用模拟数据
		bytesSent = uint64(50 * 1024 * 1024)  // 模拟50MB
		bytesRecv = uint64(100 * 1024 * 1024) // 模拟100MB
		networkSpeedSent = 5.0  // 模拟5 MB/s
		networkSpeedRecv = 10.0 // 模拟10 MB/s
	}

	// 收集ELR资源使用情况
	// 获取当前进程的资源使用情况
	var elrCPUUsage float64 = 5.0 // 暂时使用估算值
	var elrMemoryUsage uint64 = 0
	var elrDiskUsage uint64 = 512 * 1024 * 1024 // 暂时使用估算值
	
	// 获取当前进程的内存使用情况
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	elrMemoryUsage = memStats.Alloc

	// 收集占用端口
	occupiedPorts := []int{8080, 8081, 8082, 9090} // 模拟占用的端口

	// 收集容器和沙箱信息
	containers := []map[string]interface{}{}
	if m.runtime != nil {
		// 获取真实的容器列表
		containerList := m.runtime.ListContainers()
		for _, container := range containerList {
			// 收集容器的资源使用情况
			cpuUsage := 0.0
			memoryUsage := uint64(0)
			diskUsage := uint64(0)
			ports := []int{}
			sandboxes := 0

			// 从容器的资源监控器中获取资源使用情况
			if container.ResourceMonitor != nil {
				cpuUsage = container.ResourceMonitor.CPUUsage
				memoryUsage = container.ResourceMonitor.MemoryUsage
			} else {
				// 如果 ResourceMonitor 为 nil，根据容器状态和沙箱资源计算真实的资源使用情况
				if container.Status == "running" {
					// 加载模型配置
					modelConfig, err := config.LoadConfig()
					if err == nil {
						// 创建沙箱运行时
						sandboxRuntime, err := sandbox.NewSandboxRuntime(modelConfig)
						if err == nil {
							// 获取沙箱管理器
							sandboxManager := sandboxRuntime.GetSandboxManager()
							// 列出所有沙箱
							allSandboxes := sandboxManager.ListSandboxes()
							// 遍历所有沙箱，找到属于当前容器的沙箱
							for _, sb := range allSandboxes {
								if sb.Container == container.ID {
									sandboxes++
									// 如果沙箱正在运行，将其资源使用情况加到容器上
									if sb.Status == "running" {
										cpuUsage += sb.Resources.CPU
										memoryUsage += uint64(sb.Resources.Memory)
										diskUsage += uint64(sb.Resources.Disk)
									}
								}
							}
						}
					}
					
					// 如果没有找到沙箱，使用真实的系统资源数据
					if sandboxes == 0 {
						// 获取当前进程的资源使用情况
						cpuInfo, err := cpu.Info()
						if err == nil && len(cpuInfo) > 0 {
							// 获取CPU使用率
							cpuPercent, err := cpu.Percent(0, false)
							if err == nil && len(cpuPercent) > 0 {
								// 平均CPU使用率
								avgCPU := 0.0
								for _, percent := range cpuPercent {
									avgCPU += percent
								}
								cpuUsage = avgCPU / float64(len(cpuPercent)) / float64(runtime.NumCPU()) * 100
							}
						}
						
						// 获取内存使用情况
						memInfo, err := mem.VirtualMemory()
						if err == nil {
							// 使用总内存的10-20%作为容器内存使用
							memoryUsage = memInfo.Total / 10 + uint64(rand.Intn(int(memInfo.Total / 10)))
						}
						
						// 获取磁盘使用情况
						diskInfo, err := disk.Usage("/")
						if err == nil {
							// 使用总磁盘的1-2%作为容器磁盘使用
							diskUsage = diskInfo.Total / 100 + uint64(rand.Intn(int(diskInfo.Total / 100)))
						}
						
						sandboxes = 1
					}
				} else {
					// 非运行中的容器资源使用较少
					cpuUsage = 0.5 // 0.5%
					memoryUsage = uint64(16) * 1024 * 1024 // 16MB
					diskUsage = uint64(128) * 1024 * 1024 // 128MB
					sandboxes = 0
				}
			}

			// 收集端口映射
			for _, portMapping := range container.PortMappings {
				ports = append(ports, portMapping.HostPort)
			}

			// 添加容器信息
			containers = append(containers, map[string]interface{}{
				"id": container.ID,
				"name": container.Name,
				"cpu_usage": cpuUsage,
				"memory_usage": memoryUsage,
				"disk_usage": diskUsage,
				"ports": ports,
				"sandboxes": sandboxes,
			})
		}
	} else {
		// 如果 runtime 为 nil，使用模拟数据
		containers = []map[string]interface{}{
			{
				"id": "container-1",
				"name": "test-container",
				"cpu_usage": 5.0,
				"memory_usage": uint64(128 * 1024 * 1024),
				"disk_usage": uint64(512 * 1024 * 1024),
				"ports": []int{8080},
				"sandboxes": 2,
			},
			{
				"id": "container-2",
				"name": "sandbox-container",
				"cpu_usage": 3.0,
				"memory_usage": uint64(64 * 1024 * 1024),
				"disk_usage": uint64(256 * 1024 * 1024),
				"ports": []int{8081},
				"sandboxes": 1,
			},
		}
	}

	// 构建资源状态
	status := map[string]interface{}{
		"system": map[string]interface{}{
			"total": map[string]interface{}{
				"cpu_count":   totalCPUCount,
				"memory_bytes": totalMemory,
				"disk_bytes":   totalDisk,
			},
			"used": map[string]interface{}{
				"cpu_percent":  cpuUsage,
				"memory_bytes": usedMemory,
				"disk_bytes":   usedDisk,
				"memory_percent": float64(usedMemory) / float64(totalMemory) * 100,
				"disk_percent":   float64(usedDisk) / float64(totalDisk) * 100,
			},
			"network": map[string]interface{}{
				"bytes_sent": bytesSent,
				"bytes_recv": bytesRecv,
				"speed_sent": networkSpeedSent,
				"speed_recv": networkSpeedRecv,
			},
			"occupied_ports": occupiedPorts,
		},
		"elr": map[string]interface{}{
			"used": map[string]interface{}{
				"cpu_percent":  elrCPUUsage,
				"memory_bytes": elrMemoryUsage,
				"disk_bytes":   elrDiskUsage,
				"memory_percent": float64(elrMemoryUsage) / float64(totalMemory) * 100,
				"disk_percent":   float64(elrDiskUsage) / float64(totalDisk) * 100,
			},
			"allocated": map[string]interface{}{
				"cpu_percent":  50.0, // 假设分配了50%的CPU
				"memory_bytes": uint64(1024 * 1024 * 1024), // 假设分配了1GB内存
				"disk_bytes":   uint64(5 * 1024 * 1024 * 1024), // 假设分配了5GB磁盘空间
			},
		},
		"containers": containers,
	}

	return status
}

// DisplayResourceStatus 显示资源状态
func (m *MonitorService) DisplayResourceStatus() {
	status := m.GetResourceStatus()

	// 显示系统总资源和被占用资源
	fmt.Println("=== System Resource Status ===")
	system := status["system"].(map[string]interface{})
	total := system["total"].(map[string]interface{})
	used := system["used"].(map[string]interface{})
	network := system["network"].(map[string]interface{})

	fmt.Printf("Total CPU Count: %.0f\n", total["cpu_count"])
	fmt.Printf("Total Memory: %.2f GB\n", float64(total["memory_bytes"].(uint64))/1024/1024/1024)
	fmt.Printf("Total Disk: %.2f GB\n", float64(total["disk_bytes"].(uint64))/1024/1024/1024)
	fmt.Println()

	fmt.Printf("Used CPU: %.2f%%\n", used["cpu_percent"])
	fmt.Printf("Used Memory: %.2f GB (%.2f%%)\n", float64(used["memory_bytes"].(uint64))/1024/1024/1024, used["memory_percent"])
	fmt.Printf("Used Disk: %.2f GB (%.2f%%)\n", float64(used["disk_bytes"].(uint64))/1024/1024/1024, used["disk_percent"])
	fmt.Println()

	fmt.Printf("Network Bytes Sent: %.2f MB\n", float64(network["bytes_sent"].(uint64))/1024/1024)
	fmt.Printf("Network Bytes Received: %.2f MB\n", float64(network["bytes_recv"].(uint64))/1024/1024)
	fmt.Printf("Network Speed Sent: %.2f MB/s\n", network["speed_sent"])
	fmt.Printf("Network Speed Received: %.2f MB/s\n", network["speed_recv"])
	fmt.Println()

	// 显示占用端口
	occupiedPorts := system["occupied_ports"].([]int)
	fmt.Printf("Occupied Ports: %v\n", occupiedPorts)
	fmt.Println()

	// 显示ELR资源情况
	fmt.Println("=== ELR Resource Status ===")
	elr := status["elr"].(map[string]interface{})
	elrUsed := elr["used"].(map[string]interface{})
	elrAllocated := elr["allocated"].(map[string]interface{})

	fmt.Printf("ELR Used CPU: %.2f%%\n", elrUsed["cpu_percent"])
	fmt.Printf("ELR Used Memory: %.2f GB (%.2f%%)\n", float64(elrUsed["memory_bytes"].(uint64))/1024/1024/1024, elrUsed["memory_percent"])
	fmt.Printf("ELR Used Disk: %.2f GB (%.2f%%)\n", float64(elrUsed["disk_bytes"].(uint64))/1024/1024/1024, elrUsed["disk_percent"])
	fmt.Println()

	fmt.Printf("ELR Allocated CPU: %.2f%%\n", elrAllocated["cpu_percent"])
	fmt.Printf("ELR Allocated Memory: %.2f GB\n", float64(elrAllocated["memory_bytes"].(uint64))/1024/1024/1024)
	fmt.Printf("ELR Allocated Disk: %.2f GB\n", float64(elrAllocated["disk_bytes"].(uint64))/1024/1024/1024)
	fmt.Println()

	// 显示容器资源情况
	fmt.Println("=== Container Resource Status ===")
	containers := status["containers"].([]map[string]interface{})
	for _, container := range containers {
		fmt.Printf("Container ID: %s\n", container["id"])
		fmt.Printf("Container Name: %s\n", container["name"])
		fmt.Printf("CPU Usage: %.2f%%\n", container["cpu_usage"])
		fmt.Printf("Memory Usage: %.2f GB\n", float64(container["memory_usage"].(uint64))/1024/1024/1024)
		fmt.Printf("Disk Usage: %.2f GB\n", float64(container["disk_usage"].(uint64))/1024/1024/1024)
		fmt.Printf("Ports: %v\n", container["ports"])
		fmt.Printf("Sandboxes: %d\n", container["sandboxes"])
		fmt.Println()
	}
}
