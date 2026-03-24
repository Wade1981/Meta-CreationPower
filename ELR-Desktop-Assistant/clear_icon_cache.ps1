# 清除Windows图标缓存的PowerShell脚本

Write-Host "正在清除Windows图标缓存..."

# 停止Windows资源管理器
Stop-Process -Name explorer -Force

# 等待资源管理器停止
Start-Sleep -Seconds 2

# 删除图标缓存文件
$iconCacheFiles = @(
    "$env:LOCALAPPDATA\Microsoft\Windows\Explorer\iconcache_*.db",
    "$env:LOCALAPPDATA\Microsoft\Windows\Explorer\thumbcache_*.db"
)

foreach ($file in $iconCacheFiles) {
    Remove-Item -Path $file -Force -ErrorAction SilentlyContinue
}

# 重启Windows资源管理器
Start-Process explorer

Write-Host "图标缓存已清除，资源管理器已重启。"
Write-Host "请检查ELRDesktopAssistant.exe的图标是否已更新。"
