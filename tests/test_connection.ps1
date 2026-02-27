$ip = "192.168.1.100"
$port = 8766

try {
    $tcpClient = New-Object System.Net.Sockets.TcpClient
    $tcpClient.Connect($ip, $port)
    Write-Host "✅ 成功连接到 $ip`:$port"
    $tcpClient.Close()
} catch {
    Write-Host "❌ 无法连接到 $ip`:$port"
    Write-Host "错误信息: " $_.Exception.Message
}
