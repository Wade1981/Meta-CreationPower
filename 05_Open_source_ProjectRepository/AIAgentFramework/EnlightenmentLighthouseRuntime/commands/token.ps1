# Token management command module
function Manage-Token {
    if (-not $global:RUNTIME_STARTED) {
        Write-Host "Error: ELR runtime is not running"
        return
    }
    $action = "help"
    $token = ""
    $description = "ELR Container Token"
    for ($i = 1; $i -lt $args.Length; $i++) {
        if ($args[$i] -eq "--action" -and $i + 1 -lt $args.Length) {
            $action = $args[$i + 1]
            $i++
        } elseif ($args[$i] -eq "--token" -and $i + 1 -lt $args.Length) {
            $token = $args[$i + 1]
            $i++
        } elseif ($args[$i] -eq "--description" -and $i + 1 -lt $args.Length) {
            $description = $args[$i + 1]
            $i++
        }
    }
    $tokenManagerScript = "$PSScriptRoot\..\elr\token_manager.ps1"
    if (Test-Path $tokenManagerScript) {
        try {
            & $tokenManagerScript -Action $action -Token $token -Description $description
        } catch {
            Write-Host "Error managing token: $($_.Exception.Message)"
        }
    } else {
        Write-Host "Error: Token manager script not found"
        Write-Host "Please ensure the script exists at: $tokenManagerScript"
    }
}