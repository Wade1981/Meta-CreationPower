#!/usr/bin/env pwsh
# ELR Container Token Manager Script

param(
    [string]$Action = "help",
    [string]$Token = "",
    [string]$Description = "ELR Container Token"
)

$TokenFile = "$PSScriptRoot\elr_token.json"

function Test-TokenFile {
    if (-not (Test-Path $TokenFile)) {
        $initialTokenData = @{
            tokens = @()
            last_updated = [int][double]::Parse((Get-Date -UFormat %s))
        }
        $initialTokenData | ConvertTo-Json -Depth 3 | Set-Content -Path $TokenFile -Encoding UTF8
    }
}

function New-ELRToken {
    param(
        [string]$Description = "ELR Container Token"
    )
    
    Test-TokenFile
    
    $token = [System.Guid]::NewGuid().ToString("N") + [System.Guid]::NewGuid().ToString("N").Substring(0, 8)
    $tokenId = [System.Guid]::NewGuid().ToString("N").Substring(0, 8)
    $created = [int][double]::Parse((Get-Date -UFormat %s))
    $expires = $created + (7 * 24 * 60 * 60)
    
    $tokenData = Get-Content $TokenFile | ConvertFrom-Json
    $newToken = @{
        id = $tokenId
        token = $token
        description = $Description
        created = $created
        expires = $expires
        revoked = $false
    }
    
    $tokenData.tokens += $newToken
    $tokenData.last_updated = $created
    $tokenData | ConvertTo-Json -Depth 3 | Set-Content $TokenFile
    
    Write-Host "===================================="
    Write-Host "New Token Generated:"
    Write-Host $token
    Write-Host "===================================="
    Write-Host "Token ID: $tokenId"
    Write-Host "Valid for: 7 days"
    Write-Host "Please save this token securely"
    return $token
}

function Test-ELRToken {
    param(
        [string]$Token
    )
    
    if ([string]::IsNullOrEmpty($Token)) {
        Write-Host "Error: Token cannot be empty"
        return $false
    }
    
    Test-TokenFile
    
    $tokenData = Get-Content $TokenFile | ConvertFrom-Json
    $found = $false
    $valid = $false
    $message = ""
    
    foreach ($t in $tokenData.tokens) {
        if ($t.token -eq $Token) {
            $found = $true
            if ($t.revoked) {
                $message = "Token has been revoked"
            } elseif ([int][double]::Parse((Get-Date -UFormat %s)) -gt $t.expires) {
                $message = "Token has expired"
            } else {
                $valid = $true
                $message = "Token is valid"
            }
            break
        }
    }
    
    if (-not $found) {
        $message = "Token not found"
    }
    
    Write-Host "===================================="
    Write-Host "Token Validation Result:"
    if ($valid) {
        Write-Host "Status: Valid"
    } else {
        Write-Host "Status: Invalid"
    }
    Write-Host "Message: $message"
    Write-Host "===================================="
    return $valid
}

function Update-ELRToken {
    param(
        [string]$OldToken,
        [string]$Description = "Refreshed ELR Container Token"
    )
    
    if ([string]::IsNullOrEmpty($OldToken)) {
        Write-Host "Error: Old token cannot be empty"
        return
    }
    
    Test-TokenFile
    
    $tokenData = Get-Content $TokenFile | ConvertFrom-Json
    $found = $false
    $newToken = ""
    
    for ($i = 0; $i -lt $tokenData.tokens.Count; $i++) {
        if ($tokenData.tokens[$i].token -eq $OldToken) {
            $found = $true
            if ($tokenData.tokens[$i].revoked) {
                Write-Host "===================================="
                Write-Host "Token refresh failed:"
                Write-Host "Message: Token has been revoked"
                Write-Host "===================================="
                return
            }
            
            $newToken = [System.Guid]::NewGuid().ToString("N") + [System.Guid]::NewGuid().ToString("N").Substring(0, 8)
            $created = [int][double]::Parse((Get-Date -UFormat %s))
            $expires = $created + (7 * 24 * 60 * 60)
            
            $tokenData.tokens[$i].token = $newToken
            $tokenData.tokens[$i].description = $Description
            $tokenData.tokens[$i].created = $created
            $tokenData.tokens[$i].expires = $expires
            break
        }
    }
    
    if (-not $found) {
        Write-Host "===================================="
        Write-Host "Token refresh failed:"
        Write-Host "Message: Token not found"
        Write-Host "===================================="
        return
    }
    
    $tokenData.last_updated = [int][double]::Parse((Get-Date -UFormat %s))
    $tokenData | ConvertTo-Json -Depth 3 | Set-Content $TokenFile
    
    Write-Host "===================================="
    Write-Host "Token refreshed successfully:"
    Write-Host $newToken
    Write-Host "===================================="
    Write-Host "Valid for: 7 days"
    Write-Host "Please save this token securely"
    return $newToken
}

function Get-ELRTokenList {
    Test-TokenFile
    
    $tokenData = Get-Content $TokenFile | ConvertFrom-Json
    
    Write-Host "===================================="
    Write-Host "Token List:"
    Write-Host "===================================="
    Write-Host "ID       | Description         | Status  | Created"
    Write-Host "-------- | ------------------- | ------- | --------"
    
    $currentTime = [int][double]::Parse((Get-Date -UFormat %s))
    
    foreach ($t in $tokenData.tokens) {
        if ($t.revoked) {
            $status = "Revoked"
        } elseif ($currentTime -gt $t.expires) {
            $status = "Expired"
        } else {
            $status = "Valid"
        }
        $created = [DateTimeOffset]::FromUnixTimeSeconds($t.created).LocalDateTime.ToString("yyyy-MM-dd")
        $desc = $t.description
        if ($desc.Length -gt 19) {
            $desc = $desc.Substring(0, 19)
        }
        Write-Host "$($t.id) | $($desc.PadRight(19)) | $($status.PadRight(7)) | $created"
    }
    
    Write-Host "===================================="
}

function Revoke-ELRToken {
    param(
        [string]$TokenId
    )
    
    if ([string]::IsNullOrEmpty($TokenId)) {
        Write-Host "Error: Token ID cannot be empty"
        return
    }
    
    Test-TokenFile
    
    $tokenData = Get-Content $TokenFile | ConvertFrom-Json
    $found = $false
    
    for ($i = 0; $i -lt $tokenData.tokens.Count; $i++) {
        if ($tokenData.tokens[$i].id -eq $TokenId) {
            $found = $true
            $tokenData.tokens[$i].revoked = $true
            break
        }
    }
    
    if (-not $found) {
        Write-Host "===================================="
        Write-Host "Token revocation result:"
        Write-Host "Status: Failed"
        Write-Host "Message: Token ID not found"
        Write-Host "===================================="
        return
    }
    
    $tokenData.last_updated = [int][double]::Parse((Get-Date -UFormat %s))
    $tokenData | ConvertTo-Json -Depth 3 | Set-Content $TokenFile
    
    Write-Host "===================================="
    Write-Host "Token revocation result:"
    Write-Host "Status: Success"
    Write-Host "Message: Token has been revoked"
    Write-Host "===================================="
}

function Get-NetworkStatus {
    Write-Host "===================================="
    Write-Host "ELR Container Network Status"
    Write-Host "===================================="
    
    $services = @(
        @{ Name = "Desktop API"; Port = 8081; Path = "/api/desktop/health" },
        @{ Name = "Public API"; Port = 8080; Path = "/health" },
        @{ Name = "Model Service"; Port = 8082; Path = "/health" },
        @{ Name = "Micro Model Server"; Port = 8083; Path = "/health" }
    )
    
    foreach ($service in $services) {
        try {
            $url = "http://localhost:$($service.Port)$($service.Path)"
            $request = [System.Net.WebRequest]::Create($url)
            $request.Timeout = 5000
            $request.Method = "GET"
            $response = $request.GetResponse()
            $response.Close()
            Write-Host "$($service.Name): Running"
            Write-Host "  Address: http://localhost:$($service.Port)"
        } catch {
            Write-Host "$($service.Name): Not running"
        }
    }
    
    Write-Host "===================================="
}

function Show-Help {
    Write-Host "===================================="
    Write-Host "ELR Container Token Manager"
    Write-Host "===================================="
    Write-Host "Usage: .\token_manager.ps1 [Action] [Parameters]"
    Write-Host ""
    Write-Host "Actions:"
    Write-Host "  help              Show this help"
    Write-Host "  create            Create a new token"
    Write-Host "  validate          Validate a token"
    Write-Host "  refresh           Refresh a token"
    Write-Host "  list              List all tokens"
    Write-Host "  revoke            Revoke a token"
    Write-Host "  network-status    Check network status"
    Write-Host ""
    Write-Host "Parameters:"
    Write-Host "  -Token            Token value (for validate and refresh)"
    Write-Host "  -Description      Token description (for create and refresh)"
    Write-Host ""
    Write-Host "Examples:"
    Write-Host "  .\token_manager.ps1 create -Description 'Admin Token'"
    Write-Host "  .\token_manager.ps1 validate -Token 'token-value'"
    Write-Host "  .\token_manager.ps1 refresh -Token 'old-token' -Description 'Refreshed Token'"
    Write-Host "  .\token_manager.ps1 list"
    Write-Host "  .\token_manager.ps1 network-status"
    Write-Host "===================================="
}

switch ($Action) {
    "help" {
        Show-Help
    }
    "create" {
        New-ELRToken -Description $Description
    }
    "validate" {
        Test-ELRToken -Token $Token
    }
    "refresh" {
        Update-ELRToken -OldToken $Token -Description $Description
    }
    "list" {
        Get-ELRTokenList
    }
    "revoke" {
        Revoke-ELRToken -TokenId $Token
    }
    "network-status" {
        Get-NetworkStatus
    }
    default {
        Show-Help
    }
}