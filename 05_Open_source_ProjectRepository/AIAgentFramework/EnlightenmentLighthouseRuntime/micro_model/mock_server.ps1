# Mock ELR API server in PowerShell

$port = 8082
$listener = New-Object System.Net.HttpListener
$listener.Prefixes.Add("http://*:$port/")
$listener.Start()

Write-Host "Mock ELR API server started on port $port"
Write-Host "API endpoints available at: http://localhost:$port/api"
Write-Host "Press Ctrl+C to stop the server"

try {
    while ($listener.IsListening) {
        $context = $listener.GetContext()
        $request = $context.Request
        $response = $context.Response
        
        # Set response headers
        $response.ContentType = "application/json"
        $response.Headers.Add("Access-Control-Allow-Origin", "*")
        
        # Handle different endpoints
        switch ($request.Url.LocalPath) {
            "/api" {
                $jsonResponse = @{
                    status = "ok"
                    message = "ELR API is running"
                    version = "1.0.0"
                    endpoints = @(
                        "/api/models",
                        "/api/containers",
                        "/api/sandbox"
                    )
                } | ConvertTo-Json
            }
            "/api/models" {
                $jsonResponse = @{
                    status = "ok"
                    models = @(
                        @{
                            name = "fish-speech"
                            version = "1.0.0"
                            status = "loaded"
                            path = "model/models/fish-speech"
                        }
                    )
                } | ConvertTo-Json
            }
            "/api/containers" {
                $jsonResponse = @{
                    status = "ok"
                    containers = @(
                        @{
                            id = "elr-1234567890"
                            name = "test-container"
                            status = "running"
                        }
                    )
                } | ConvertTo-Json
            }
            "/api/sandbox" {
                $jsonResponse = @{
                    status = "ok"
                    sandbox = @{
                        status = "running"
                        models = @("fish-speech")
                    }
                } | ConvertTo-Json
            }
            default {
                $response.StatusCode = 404
                $jsonResponse = @{
                    status = "error"
                    message = "Endpoint not found"
                } | ConvertTo-Json
            }
        }
        
        # Write response
        $buffer = [System.Text.Encoding]::UTF8.GetBytes($jsonResponse)
        $response.ContentLength64 = $buffer.Length
        $response.OutputStream.Write($buffer, 0, $buffer.Length)
        $response.OutputStream.Close()
        
        # Log request
        Write-Host "[$(Get-Date -Format 'yyyy-MM-dd HH:mm:ss')] $($request.HttpMethod) $($request.Url.LocalPath)"
    }
} catch {
    Write-Host "Error: $($_.Exception.Message)"
} finally {
    $listener.Stop()
    Write-Host "Server stopped"
}
