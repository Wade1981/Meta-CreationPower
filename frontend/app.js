// 前端应用核心逻辑

// 页面导航功能
function setupNavigation() {
    const navLinks = document.querySelectorAll('.nav-link');
    const pages = document.querySelectorAll('.page');
    
    navLinks.forEach(link => {
        link.addEventListener('click', function(e) {
            e.preventDefault();
            
            // 移除所有导航链接的活跃状态
            navLinks.forEach(navLink => navLink.classList.remove('active'));
            // 添加当前导航链接的活跃状态
            this.classList.add('active');
            
            // 隐藏所有页面
            pages.forEach(page => page.classList.remove('active'));
            // 显示对应页面
            const targetId = this.getAttribute('href').substring(1);
            const targetPage = document.getElementById(targetId);
            if (targetPage) {
                targetPage.classList.add('active');
                
                // 根据页面ID加载对应数据
                if (targetId === 'dashboard') {
                    loadDashboardData();
                } else if (targetId === 'assets') {
                    loadAssets();
                } else if (targetId === 'system') {
                    loadSystemStatus();
                } else if (targetId === 'workflow') {
                    loadWorkflows();
                } else if (targetId === 'elr-models') {
                    loadElrModels();
                } else if (targetId === 'elr-containers') {
                    loadElrContainers();
                } else if (targetId === 'elr-sandbox') {
                    // 沙箱页面无需初始加载数据
                } else if (targetId === 'elr-monitor') {
                    loadElrMonitor();
                }
            }
        });
    });
}

// API基础URL
const API_BASE_URL = 'http://localhost:5000'; // 后端API运行在5000端口
// ELR API基础URL
const ELR_API_BASE_URL = 'http://localhost:5000/api/elr'; // ELR API通过后端代理

// API请求函数 - 使用模拟数据
async function apiRequest(endpoint, method = 'GET', data = null) {
    // 直接返回模拟数据，无需实际API调用
    console.log(`模拟API请求: ${method} ${endpoint}`, data);
    return getMockData(endpoint, method, data);
}

// 模拟数据，用于API请求失败时的 fallback
function getMockData(endpoint, method, data) {
    const mockData = {
        dashboard: {
            total_assets: 15,
            active_processes: 3,
            total_transactions: 45,
            system_health: 'healthy'
        },
        assets: [
            {
                asset_id: 'asset-001',
                asset_type: 'general',
                status: 'active',
                created_at: '2026-02-20T10:00:00Z'
            },
            {
                asset_id: 'asset-002',
                asset_type: 'algorithm',
                status: 'active',
                created_at: '2026-02-21T11:30:00Z'
            },
            {
                asset_id: 'asset-003',
                asset_type: 'feature',
                status: 'inactive',
                created_at: '2026-02-22T14:20:00Z'
            }
        ],
        systemStatus: {
            dashboard: {
                total_assets: 15,
                active_processes: 3,
                total_transactions: 45,
                system_health: 'healthy'
            },
            modules: {
                raw_material_engine: 'active',
                feature_asset_engine: 'active',
                algorithm_asset_engine: 'active',
                asset_packaging_engine: 'active',
                trading_engine: 'active',
                elr_integration: 'active',
                compression_utils: 'active',
                encryption_utils: 'active',
                file_transfer_utils: 'active',
                collaborative_network: 'active',
                ide_integration: 'active'
            },
            containers: [
                {
                    container_id: 'container-001',
                    status: 'running',
                    asset_id: 'asset-001'
                },
                {
                    container_id: 'container-002',
                    status: 'running',
                    asset_id: 'asset-002'
                }
            ]
        },
        workflows: [
            {
                execution_id: 'exec-001',
                workflow_id: 'workflow-001',
                status: 'completed',
                started_at: '2026-02-23T09:00:00Z',
                completed_at: '2026-02-23T09:05:00Z'
            },
            {
                execution_id: 'exec-002',
                workflow_id: 'workflow-002',
                status: 'running',
                started_at: '2026-02-24T10:30:00Z',
                completed_at: null
            }
        ]
    };
    
    // 根据端点返回相应的模拟数据
    if (endpoint.includes('/dashboard')) {
        return { success: true, data: mockData.dashboard };
    } else if (endpoint.includes('/assets') && method === 'GET') {
        return { success: true, data: mockData.assets, total: mockData.assets.length };
    } else if (endpoint.includes('/system/status')) {
        return { success: true, data: mockData.systemStatus };
    } else if (endpoint.includes('/workflows')) {
        if (method === 'GET') {
            return { success: true, data: mockData.workflows };
        } else if (method === 'POST') {
            return { success: true, data: {
                execution_id: 'exec-' + Math.random().toString(36).substr(2, 8),
                workflow_id: data.workflow_id,
                status: 'completed',
                started_at: new Date().toISOString(),
                completed_at: new Date().toISOString()
            }};
        }
    } else if (endpoint.includes('/assets') && method === 'POST') {
        return { success: true, data: {
            asset_id: 'asset-' + Math.random().toString(36).substr(2, 8),
            asset_type: data.type || 'general',
            status: 'active',
            created_at: new Date().toISOString(),
            last_modified: new Date().toISOString()
        }, asset_id: 'asset-' + Math.random().toString(36).substr(2, 8) };
    }
    
    return { success: false, message: 'No mock data available' };
}

// 加载仪表盘数据
async function loadDashboardData() {
    // 显示加载状态
    const totalAssetsEl = document.getElementById('total-assets');
    const activeProcessesEl = document.getElementById('active-processes');
    const totalTransactionsEl = document.getElementById('total-transactions');
    const systemHealthEl = document.getElementById('system-health');
    
    const originalValues = {
        total: totalAssetsEl.textContent,
        active: activeProcessesEl.textContent,
        transactions: totalTransactionsEl.textContent,
        health: systemHealthEl.textContent
    };
    
    // 设置加载状态
    totalAssetsEl.textContent = '<loading>';
    activeProcessesEl.textContent = '<loading>';
    totalTransactionsEl.textContent = '<loading>';
    systemHealthEl.textContent = '<loading>';
    
    try {
        // 调用API获取仪表盘数据
        const response = await apiRequest('/api/dashboard');
        if (response.success) {
            const dashboardData = response.data;
            
            // 更新仪表盘数据
            document.getElementById('total-assets').textContent = dashboardData.total_assets;
            document.getElementById('active-processes').textContent = dashboardData.active_processes;
            document.getElementById('total-transactions').textContent = dashboardData.total_transactions;
            document.getElementById('system-health').textContent = dashboardData.system_health;
            
            // 初始化资产分布图表
            initAssetDistributionChart();
        }
    } catch (error) {
        console.error('Error loading dashboard data:', error);
        // 恢复原始值
        totalAssetsEl.textContent = originalValues.total;
        activeProcessesEl.textContent = originalValues.active;
        totalTransactionsEl.textContent = originalValues.transactions;
        systemHealthEl.textContent = originalValues.health;
    }
}

// 初始化资产分布图表
function initAssetDistributionChart() {
    const ctx = document.getElementById('asset-distribution-chart').getContext('2d');
    
    // 销毁现有图表
    if (window.assetChart) {
        window.assetChart.destroy();
    }
    
    // 创建新图表
    window.assetChart = new Chart(ctx, {
        type: 'bar',
        data: {
            labels: ['通用资产', '算法资产', '特征资产', '打包资产'],
            datasets: [{
                label: '资产数量',
                data: [8, 3, 2, 2],
                backgroundColor: [
                    'rgba(54, 162, 235, 0.6)',
                    'rgba(75, 192, 192, 0.6)',
                    'rgba(153, 102, 255, 0.6)',
                    'rgba(255, 159, 64, 0.6)'
                ],
                borderColor: [
                    'rgba(54, 162, 235, 1)',
                    'rgba(75, 192, 192, 1)',
                    'rgba(153, 102, 255, 1)',
                    'rgba(255, 159, 64, 1)'
                ],
                borderWidth: 1
            }]
        },
        options: {
            responsive: true,
            maintainAspectRatio: false,
            scales: {
                y: {
                    beginAtZero: true
                }
            }
        }
    });
}

// 加载资产列表
async function loadAssets() {
    const tableBody = document.getElementById('asset-table-body');
    
    // 显示加载状态
    tableBody.innerHTML = '<tr><td colspan="5" class="text-center"><div class="loading"></div> 加载中...</td></tr>';
    
    try {
        // 调用API获取资产列表
        const response = await apiRequest('/api/assets');
        if (response.success) {
            const assets = response.data;
            
            // 清空表格
            tableBody.innerHTML = '';
            
            // 填充表格数据
            assets.forEach(asset => {
                const row = document.createElement('tr');
                
                // 格式化创建时间
                const createdDate = new Date(asset.created_at).toLocaleString();
                
                row.innerHTML = `
                    <td>${asset.asset_id}</td>
                    <td>${asset.asset_type}</td>
                    <td>
                        <span class="badge ${asset.status === 'active' ? 'bg-success' : 'bg-warning'}">
                            ${asset.status}
                        </span>
                    </td>
                    <td>${createdDate}</td>
                    <td>
                        <button class="btn btn-sm btn-info" onclick="viewAsset('${asset.asset_id}')">查看</button>
                        <button class="btn btn-sm btn-primary" onclick="editAsset('${asset.asset_id}')">编辑</button>
                        <button class="btn btn-sm btn-danger" onclick="deleteAsset('${asset.asset_id}')">删除</button>
                    </td>
                `;
                
                tableBody.appendChild(row);
            });
        } else {
            tableBody.innerHTML = '<tr><td colspan="5" class="text-center text-danger">加载资产失败</td></tr>';
        }
    } catch (error) {
        console.error('Error loading assets:', error);
        tableBody.innerHTML = '<tr><td colspan="5" class="text-center text-danger">加载资产失败</td></tr>';
    }
}

// 查看资产详情
function viewAsset(assetId) {
    alert(`查看资产详情: ${assetId}`);
    // 这里可以实现查看资产详情的逻辑
}

// 编辑资产
function editAsset(assetId) {
    alert(`编辑资产: ${assetId}`);
    // 这里可以实现编辑资产的逻辑
}

// 删除资产
async function deleteAsset(assetId) {
    if (confirm(`确定要删除资产 ${assetId} 吗？`)) {
        try {
            // 调用API删除资产
            const response = await apiRequest(`/api/assets/${assetId}`, 'DELETE');
            if (response.success) {
                alert(`资产 ${assetId} 已删除`);
                loadAssets(); // 重新加载资产列表
                loadDashboardData(); // 重新加载仪表盘数据
            } else {
                alert(`删除资产失败: ${response.message || '未知错误'}`);
            }
        } catch (error) {
            console.error('Error deleting asset:', error);
            alert('删除资产失败，请稍后重试');
        }
    }
}

// 加载系统状态
async function loadSystemStatus() {
    const moduleStatuses = document.getElementById('module-statuses');
    const containerTableBody = document.getElementById('container-table-body');
    
    // 显示加载状态
    moduleStatuses.innerHTML = '<div class="col-md-12 text-center"><div class="loading"></div> 加载中...</div>';
    containerTableBody.innerHTML = '<tr><td colspan="3" class="text-center"><div class="loading"></div> 加载中...</td></tr>';
    
    try {
        // 调用API获取系统状态
        const response = await apiRequest('/api/system/status');
        if (response.success) {
            const systemData = response.data;
            
            // 清空模块状态
            moduleStatuses.innerHTML = '';
            
            // 填充模块状态
            Object.entries(systemData.modules).forEach(([moduleName, status]) => {
                const moduleDiv = document.createElement('div');
                moduleDiv.className = `col-md-4 module-status ${status === 'active' ? 'active' : 'inactive'}`;
                moduleDiv.innerHTML = `
                    <strong>${moduleName.replace('_', ' ')}</strong>
                    <span class="float-end">${status}</span>
                `;
                moduleStatuses.appendChild(moduleDiv);
            });
            
            // 清空容器表格
            containerTableBody.innerHTML = '';
            
            // 填充容器表格
            systemData.containers.forEach(container => {
                const row = document.createElement('tr');
                row.innerHTML = `
                    <td>${container.container_id}</td>
                    <td>
                        <span class="badge ${container.status === 'running' ? 'bg-success' : 'bg-warning'}">
                            ${container.status}
                        </span>
                    </td>
                    <td>${container.asset_id}</td>
                `;
                containerTableBody.appendChild(row);
            });
        } else {
            moduleStatuses.innerHTML = '<div class="col-md-12 text-center text-danger">加载模块状态失败</div>';
            containerTableBody.innerHTML = '<tr><td colspan="3" class="text-center text-danger">加载容器状态失败</td></tr>';
        }
    } catch (error) {
        console.error('Error loading system status:', error);
        moduleStatuses.innerHTML = '<div class="col-md-12 text-center text-danger">加载模块状态失败</div>';
        containerTableBody.innerHTML = '<tr><td colspan="3" class="text-center text-danger">加载容器状态失败</td></tr>';
    }
}

// 加载工作流数据
async function loadWorkflows() {
    const tableBody = document.getElementById('workflow-table-body');
    
    // 显示加载状态
    tableBody.innerHTML = '<tr><td colspan="5" class="text-center"><div class="loading"></div> 加载中...</td></tr>';
    
    try {
        // 调用API获取工作流列表
        const response = await apiRequest('/api/workflows');
        if (response.success) {
            const workflows = response.data;
            
            // 清空表格
            tableBody.innerHTML = '';
            
            // 填充表格数据
            workflows.forEach(workflow => {
                const row = document.createElement('tr');
                
                // 格式化时间
                const startedDate = workflow.started_at ? new Date(workflow.started_at).toLocaleString() : '-';
                const completedDate = workflow.completed_at ? new Date(workflow.completed_at).toLocaleString() : '-';
                
                row.innerHTML = `
                    <td>${workflow.execution_id}</td>
                    <td>${workflow.workflow_id}</td>
                    <td>
                        <span class="badge ${workflow.status === 'completed' ? 'bg-success' : workflow.status === 'running' ? 'bg-info' : 'bg-warning'}">
                            ${workflow.status}
                        </span>
                    </td>
                    <td>${startedDate}</td>
                    <td>${completedDate}</td>
                `;
                
                tableBody.appendChild(row);
            });
        } else {
            tableBody.innerHTML = '<tr><td colspan="5" class="text-center text-danger">加载工作流失败</td></tr>';
        }
    } catch (error) {
        console.error('Error loading workflows:', error);
        tableBody.innerHTML = '<tr><td colspan="5" class="text-center text-danger">加载工作流失败</td></tr>';
    }
}

// 设置表单提交事件
function setupFormSubmission() {
    // 创建资产表单提交
    document.getElementById('submit-asset').addEventListener('click', async function() {
        const assetType = document.getElementById('asset-type').value;
        const assetContent = document.getElementById('asset-content').value;
        const dataType = document.getElementById('data-type').value;
        const structureType = document.getElementById('structure-type').value;
        
        if (!assetContent) {
            alert('请输入资产内容');
            return;
        }
        
        // 显示加载状态
        this.disabled = true;
        this.innerHTML = '<div class="loading"></div> 创建中...';
        
        try {
            // 调用API创建资产
            const response = await apiRequest('/api/assets', 'POST', {
                type: assetType,
                content: assetContent,
                data_type: dataType,
                structure_type: structureType
            });
            
            if (response.success) {
                alert('资产创建成功！');
                // 关闭模态框
                const modal = bootstrap.Modal.getInstance(document.getElementById('createAssetModal'));
                modal.hide();
                // 清空表单
                document.getElementById('create-asset-form').reset();
                // 重新加载资产列表和仪表盘数据
                loadAssets();
                loadDashboardData();
            } else {
                alert(`创建资产失败: ${response.message || '未知错误'}`);
            }
        } catch (error) {
            console.error('Error creating asset:', error);
            alert('创建资产失败，请稍后重试');
        } finally {
            // 恢复按钮状态
            this.disabled = false;
            this.innerHTML = '创建';
        }
    });
    
    // 执行工作流表单提交
    document.getElementById('submit-workflow').addEventListener('click', async function() {
        const workflowId = document.getElementById('workflow-id').value;
        const workflowData = document.getElementById('workflow-data').value;
        
        if (!workflowId) {
            alert('请输入工作流ID');
            return;
        }
        
        let workflowDataObj = {};
        try {
            if (workflowData) {
                workflowDataObj = JSON.parse(workflowData);
            }
        } catch (e) {
            alert('工作流数据格式错误，请输入有效的JSON');
            return;
        }
        
        // 显示加载状态
        this.disabled = true;
        this.innerHTML = '<div class="loading"></div> 执行中...';
        
        try {
            // 调用API执行工作流
            const response = await apiRequest('/api/workflows', 'POST', {
                workflow_id: workflowId,
                ...workflowDataObj
            });
            
            if (response.success) {
                alert('工作流执行成功！');
                // 关闭模态框
                const modal = bootstrap.Modal.getInstance(document.getElementById('createWorkflowModal'));
                modal.hide();
                // 清空表单
                document.getElementById('create-workflow-form').reset();
                // 重新加载工作流列表
                loadWorkflows();
            } else {
                alert(`执行工作流失败: ${response.message || '未知错误'}`);
            }
        } catch (error) {
            console.error('Error executing workflow:', error);
            alert('执行工作流失败，请稍后重试');
        } finally {
            // 恢复按钮状态
            this.disabled = false;
            this.innerHTML = '执行';
        }
    });
}

// 设置搜索功能
function setupSearch() {
    document.getElementById('search-btn').addEventListener('click', async function() {
        const searchTerm = document.getElementById('asset-search').value;
        const tableBody = document.getElementById('asset-table-body');
        
        // 显示加载状态
        tableBody.innerHTML = '<tr><td colspan="5" class="text-center"><div class="loading"></div> 搜索中...</td></tr>';
        
        try {
            // 调用API搜索资产
            // 这里简化处理，实际项目中可能需要后端支持搜索功能
            const response = await apiRequest('/api/assets');
            if (response.success) {
                let assets = response.data;
                
                // 前端过滤搜索结果
                if (searchTerm) {
                    assets = assets.filter(asset => 
                        asset.asset_id.includes(searchTerm) ||
                        asset.asset_type.includes(searchTerm) ||
                        asset.status.includes(searchTerm)
                    );
                }
                
                // 清空表格
                tableBody.innerHTML = '';
                
                // 填充表格数据
                if (assets.length === 0) {
                    tableBody.innerHTML = '<tr><td colspan="5" class="text-center">未找到匹配的资产</td></tr>';
                } else {
                    assets.forEach(asset => {
                        const row = document.createElement('tr');
                        
                        // 格式化创建时间
                        const createdDate = new Date(asset.created_at).toLocaleString();
                        
                        row.innerHTML = `
                            <td>${asset.asset_id}</td>
                            <td>${asset.asset_type}</td>
                            <td>
                                <span class="badge ${asset.status === 'active' ? 'bg-success' : 'bg-warning'}">
                                    ${asset.status}
                                </span>
                            </td>
                            <td>${createdDate}</td>
                            <td>
                                <button class="btn btn-sm btn-info" onclick="viewAsset('${asset.asset_id}')">查看</button>
                                <button class="btn btn-sm btn-primary" onclick="editAsset('${asset.asset_id}')">编辑</button>
                                <button class="btn btn-sm btn-danger" onclick="deleteAsset('${asset.asset_id}')">删除</button>
                            </td>
                        `;
                        
                        tableBody.appendChild(row);
                    });
                }
            } else {
                tableBody.innerHTML = '<tr><td colspan="5" class="text-center text-danger">搜索失败</td></tr>';
            }
        } catch (error) {
            console.error('Error searching assets:', error);
            tableBody.innerHTML = '<tr><td colspan="5" class="text-center text-danger">搜索失败，请稍后重试</td></tr>';
        }
    });
}

// 初始化应用
async function initApp() {
    setupNavigation();
    setupFormSubmission();
    setupSearch();
    
    // 初始加载仪表盘数据
    await loadDashboardData();
}

// 页面加载完成后初始化应用
document.addEventListener('DOMContentLoaded', initApp);

// 模拟后端API服务
function startMockApiServer() {
    // 这里可以添加模拟API服务器的逻辑
    console.log('Mock API server started');
}

// ELR API请求函数
async function elrApiRequest(endpoint, method = 'GET', data = null) {
    const url = `${ELR_API_BASE_URL}${endpoint}`;
    try {
        const options = {
            method,
            headers: {
                'Content-Type': 'application/json'
            }
        };
        if (data) {
            options.body = JSON.stringify(data);
        }
        const response = await fetch(url, options);
        const result = await response.json();
        return result;
    } catch (error) {
        console.error('ELR API request error:', error);
        return { success: false, message: 'API request failed' };
    }
}

// 加载ELR模型列表
async function loadElrModels() {
    const tableBody = document.getElementById('elr-models-table-body');
    tableBody.innerHTML = '<tr><td colspan="5" class="text-center"><div class="loading"></div> 加载中...</td></tr>';
    
    try {
        const response = await elrApiRequest('/models');
        if (response.success) {
            const models = response.data;
            tableBody.innerHTML = '';
            models.forEach(model => {
                const row = document.createElement('tr');
                row.innerHTML = `
                    <td>${model.id || model.model_id}</td>
                    <td>${model.type || model.model_type}</td>
                    <td>${model.status || 'unknown'}</td>
                    <td>${model.created_at || 'unknown'}</td>
                    <td>
                        <button class="btn btn-sm btn-info" onclick="viewElrModel('${model.id || model.model_id}')">查看</button>
                        <button class="btn btn-sm btn-primary" onclick="runElrModel('${model.id || model.model_id}')">运行</button>
                        <button class="btn btn-sm btn-danger" onclick="deleteElrModel('${model.id || model.model_id}')">删除</button>
                    </td>
                `;
                tableBody.appendChild(row);
            });
        } else {
            tableBody.innerHTML = '<tr><td colspan="5" class="text-center text-danger">加载模型失败</td></tr>';
        }
    } catch (error) {
        console.error('Error loading ELR models:', error);
        tableBody.innerHTML = '<tr><td colspan="5" class="text-center text-danger">加载模型失败</td></tr>';
    }
}

// 加载ELR容器列表
async function loadElrContainers() {
    const tableBody = document.getElementById('elr-containers-table-body');
    tableBody.innerHTML = '<tr><td colspan="4" class="text-center"><div class="loading"></div> 加载中...</td></tr>';
    
    try {
        const response = await elrApiRequest('/containers');
        if (response.success) {
            const containers = response.data;
            tableBody.innerHTML = '';
            containers.forEach(container => {
                const row = document.createElement('tr');
                row.innerHTML = `
                    <td>${container.name || container.container_name}</td>
                    <td>${container.status || 'unknown'}</td>
                    <td>${container.model_id || 'unknown'}</td>
                    <td>
                        <button class="btn btn-sm btn-info" onclick="viewElrContainer('${container.name || container.container_name}')">查看</button>
                        <button class="btn btn-sm btn-success" onclick="startElrContainer('${container.name || container.container_name}')">启动</button>
                        <button class="btn btn-sm btn-warning" onclick="stopElrContainer('${container.name || container.container_name}')">停止</button>
                        <button class="btn btn-sm btn-danger" onclick="deleteElrContainer('${container.name || container.container_name}')">删除</button>
                    </td>
                `;
                tableBody.appendChild(row);
            });
        } else {
            tableBody.innerHTML = '<tr><td colspan="4" class="text-center text-danger">加载容器失败</td></tr>';
        }
    } catch (error) {
        console.error('Error loading ELR containers:', error);
        tableBody.innerHTML = '<tr><td colspan="4" class="text-center text-danger">加载容器失败</td></tr>';
    }
}

// 加载ELR监控数据
async function loadElrMonitor() {
    const monitorContent = document.getElementById('elr-monitor-content');
    monitorContent.innerHTML = '<div class="text-center"><div class="loading"></div> 加载中...</div>';
    
    try {
        const response = await elrApiRequest('/monitor/metrics');
        if (response.success) {
            const metrics = response.data;
            monitorContent.innerHTML = `
                <div class="card">
                    <div class="card-body">
                        <h5 class="card-title">ELR 监控指标</h5>
                        <pre>${JSON.stringify(metrics, null, 2)}</pre>
                    </div>
                </div>
            `;
        } else {
            monitorContent.innerHTML = '<div class="text-center text-danger">加载监控数据失败</div>';
        }
    } catch (error) {
        console.error('Error loading ELR monitor:', error);
        monitorContent.innerHTML = '<div class="text-center text-danger">加载监控数据失败</div>';
    }
}

// ELR模型操作函数
async function viewElrModel(modelId) {
    try {
        const response = await elrApiRequest(`/models/${modelId}`);
        if (response.success) {
            alert(JSON.stringify(response.data, null, 2));
        } else {
            alert('查看模型失败');
        }
    } catch (error) {
        console.error('Error viewing ELR model:', error);
        alert('查看模型失败');
    }
}

async function runElrModel(modelId) {
    const containerName = prompt('请输入容器名称：');
    const input = prompt('请输入模型输入：');
    if (containerName && input) {
        try {
            const response = await elrApiRequest('/models/run', 'POST', {
                container_name: containerName,
                model_id: modelId,
                input: input
            });
            if (response.success) {
                alert(`模型运行结果：${response.data.output}`);
            } else {
                alert('运行模型失败');
            }
        } catch (error) {
            console.error('Error running ELR model:', error);
            alert('运行模型失败');
        }
    }
}

async function deleteElrModel(modelId) {
    if (confirm(`确定要删除模型 ${modelId} 吗？`)) {
        try {
            const response = await elrApiRequest(`/models/${modelId}`, 'DELETE');
            if (response.success) {
                alert('模型删除成功');
                loadElrModels();
            } else {
                alert('删除模型失败');
            }
        } catch (error) {
            console.error('Error deleting ELR model:', error);
            alert('删除模型失败');
        }
    }
}

// ELR容器操作函数
async function viewElrContainer(containerName) {
    try {
        const response = await elrApiRequest(`/containers/${containerName}`);
        if (response.success) {
            alert(JSON.stringify(response.data, null, 2));
        } else {
            alert('查看容器失败');
        }
    } catch (error) {
        console.error('Error viewing ELR container:', error);
        alert('查看容器失败');
    }
}

async function startElrContainer(containerName) {
    try {
        const response = await elrApiRequest('/containers/start', 'POST', {
            container_name: containerName
        });
        if (response.success) {
            alert('容器启动成功');
            loadElrContainers();
        } else {
            alert('启动容器失败');
        }
    } catch (error) {
        console.error('Error starting ELR container:', error);
        alert('启动容器失败');
    }
}

async function stopElrContainer(containerName) {
    try {
        const response = await elrApiRequest('/containers/stop', 'POST', {
            container_name: containerName
        });
        if (response.success) {
            alert('容器停止成功');
            loadElrContainers();
        } else {
            alert('停止容器失败');
        }
    } catch (error) {
        console.error('Error stopping ELR container:', error);
        alert('停止容器失败');
    }
}

async function deleteElrContainer(containerName) {
    if (confirm(`确定要删除容器 ${containerName} 吗？`)) {
        try {
            const response = await elrApiRequest(`/containers/${containerName}`, 'DELETE');
            if (response.success) {
                alert('容器删除成功');
                loadElrContainers();
            } else {
                alert('删除容器失败');
            }
        } catch (error) {
            console.error('Error deleting ELR container:', error);
            alert('删除容器失败');
        }
    }
}

// 创建ELR容器
async function createElrContainer() {
    const containerName = document.getElementById('container-name').value;
    const modelId = document.getElementById('container-model').value;
    if (containerName && modelId) {
        try {
            const response = await elrApiRequest('/containers', 'POST', {
                container_name: containerName,
                model_id: modelId,
                resources: {}
            });
            if (response.success) {
                alert('容器创建成功');
                // 关闭模态框
                const modal = bootstrap.Modal.getInstance(document.getElementById('createContainerModal'));
                modal.hide();
                // 清空表单
                document.getElementById('create-container-form').reset();
                // 重新加载容器列表
                loadElrContainers();
            } else {
                alert('创建容器失败');
            }
        } catch (error) {
            console.error('Error creating ELR container:', error);
            alert('创建容器失败');
        }
    } else {
        alert('请输入容器名称和模型ID');
    }
}

// 下载ELR模型
async function downloadElrModel() {
    const modelId = document.getElementById('model-id').value;
    const modelType = document.getElementById('model-type').value;
    const downloadUrl = document.getElementById('model-url').value;
    if (modelId && downloadUrl) {
        try {
            const response = await elrApiRequest('/models', 'POST', {
                model_id: modelId,
                model_type: modelType,
                download_url: downloadUrl
            });
            if (response.success) {
                alert('模型下载成功');
                // 关闭模态框
                const modal = bootstrap.Modal.getInstance(document.getElementById('downloadModelModal'));
                modal.hide();
                // 清空表单
                document.getElementById('download-model-form').reset();
                // 重新加载模型列表
                loadElrModels();
            } else {
                alert('下载模型失败');
            }
        } catch (error) {
            console.error('Error downloading ELR model:', error);
            alert('下载模型失败');
        }
    } else {
        alert('请输入模型ID和下载URL');
    }
}

// 执行沙箱命令
async function executeSandboxCommand() {
    const containerName = document.getElementById('sandbox-container').value;
    const command = document.getElementById('sandbox-command').value;
    if (containerName && command) {
        try {
            const response = await elrApiRequest('/sandbox/execute', 'POST', {
                container_name: containerName,
                command: command.split(' ')
            });
            if (response.success) {
                document.getElementById('sandbox-output').value = response.data.output;
            } else {
                alert('执行命令失败');
            }
        } catch (error) {
            console.error('Error executing sandbox command:', error);
            alert('执行命令失败');
        }
    } else {
        alert('请输入容器名称和命令');
    }
}

// 启动模拟API服务器
startMockApiServer();