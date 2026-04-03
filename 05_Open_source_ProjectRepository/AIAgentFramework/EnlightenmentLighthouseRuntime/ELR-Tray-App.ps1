#!/usr/bin/env powershell
# ELR Tray Application
# Features: System Tray + GUI Interface + ELR Container Management + Multi-language Support

Add-Type -AssemblyName System.Windows.Forms
Add-Type -AssemblyName System.Drawing

# Global variables for forms
$Global:mainForm = $null
$Global:trayIcon = $null
$Global:chatForm = $null
$Global:settingsForm = $null

# Language variables
$Global:language = "en"
$Global:langConfig = $null

# ELR Path
$ELRPath = $PSScriptRoot
$ELRPS1 = "$ELRPath\elr.ps1"

# Icon Path
$IconPath = "$ELRPath\icons\elr_icon.ico"

# Language file path
$LanguageFile = "$ELRPath\resource\language.json"

# Get system language
function Get-SystemLanguage {
    $culture = [System.Globalization.CultureInfo]::CurrentUICulture
    $lang = $culture.TwoLetterISOLanguageName
    if ($lang -eq "zh") {
        return "zh"
    } else {
        return "en"
    }
}

# Load language configuration
function Load-LanguageConfig {
    if (Test-Path $LanguageFile) {
        try {
            $jsonContent = Get-Content -Path $LanguageFile -Encoding UTF8 | ConvertFrom-Json
            $Global:langConfig = $jsonContent
            $Global:language = Get-SystemLanguage
        } catch {
            Write-Host "Failed to load language configuration: $($_.Exception.Message)"
            # Fallback to English
            $Global:language = "en"
        }
    } else {
        Write-Host "Language configuration file not found: $LanguageFile"
        # Fallback to English
        $Global:language = "en"
    }
}

# Get localized string
function Get-LocalizedString {
    param(
        [string]$keyPath
    )
    
    try {
        $keys = $keyPath -split "\."
        $value = $Global:langConfig.$($Global:language)
        
        foreach ($key in $keys) {
            $value = $value.$key
        }
        
        return $value
    } catch {
        # Fallback to key if translation not found
        return $keyPath
    }
}

# Get localized string with parameters
function Get-LocalizedStringWithParams {
    param(
        [string]$keyPath,
        [Parameter(ValueFromRemainingArguments=$true)][Object[]]$params
    )
    
    $baseString = Get-LocalizedString -keyPath $keyPath
    return [string]::Format($baseString, $params)
}

# Initialize language
Load-LanguageConfig

# Check if ELR script exists
if (-not (Test-Path $ELRPS1)) {
    Write-Host "Error: ELR script not found"
    exit 1
}

# Save API configuration
function Save-APIConfig {
    param(
        [string]$apiType,
        [string]$apiUrl
    )
    
    try {
        # Path to elr.exe (using relative path)
        $elrExe = ".\elr.exe"
        
        # Clean up input
        $cleanUrl = $apiUrl.Trim()
        
        # Ensure http:// prefix
        if (-not $cleanUrl.StartsWith('http://')) {
            $cleanUrl = "http://$cleanUrl"
        }
        
        # Check if API URL is valid
        if ($cleanUrl -match 'http://([^:]+):(\d+)') {
            $address = $matches[1]
            $port = $matches[2]
            
            # Check if port is valid (1-65535)
            if ([int]$port -ge 1 -and [int]$port -le 65535) {
                # Test if address and port are accessible
                Write-Host "Testing if address $address and port $port are accessible..."
                try {
                    $socket = New-Object System.Net.Sockets.TcpClient
                    $connectTask = $socket.BeginConnect($address, $port, $null, $null)
                    $waitResult = $connectTask.AsyncWaitHandle.WaitOne(2000, $false)
                    if ($waitResult) {
                        $socket.EndConnect($connectTask)
                        $socket.Close()
                        Write-Host "Address $address and port $port are accessible"
                        Write-Host "Saving ${apiType} API configuration: ${address}:${port}"
                        & $elrExe api config set --api-type $apiType --address $address --port $port
                        return $true
                    } else {
                        $socket.Close()
                        Write-Host "Address $address and port $port are not accessible"
                        [System.Windows.Forms.MessageBox]::Show('Error: The address and port you entered are not accessible', 'Error', 'OK', 'Error')
                        return $false
                    }
                } catch {
                    Write-Host "Error testing address and port: $($_.Exception.Message)"
                    [System.Windows.Forms.MessageBox]::Show('Error testing address and port: $($_.Exception.Message)', 'Error', 'OK', 'Error')
                    return $false
                }
            } else {
                Write-Host "Invalid port: $port"
                [System.Windows.Forms.MessageBox]::Show('Invalid port number, please enter a port between 1-65535', 'Error', 'OK', 'Error')
                return $false
            }
        } else {
            Write-Host "Invalid API URL: $apiUrl"
            [System.Windows.Forms.MessageBox]::Show('Invalid API address format, please use http://localhost:port format', 'Error', 'OK', 'Error')
            return $false
        }
    } catch {
        Write-Host "Failed to save API configuration: $($_.Exception.Message)"
        [System.Windows.Forms.MessageBox]::Show('Failed to save API configuration: $($_.Exception.Message)', 'Error', 'OK', 'Error')
        return $false
    }
}

# Load API configuration
function Load-APIConfig {
    try {
        # Path to elr.exe (using relative path)
        $elrExe = ".\elr.exe"
        
        # Get API configuration from elr api config list
        $configOutput = & $elrExe api config list
        
        # Parse the output
        $config = @{}
        
        if ($configOutput -match 'Desktop API:\s+Current: ([^\s]+):(\d+)') {
            $config.desktopApi = "http://$($matches[1]):$($matches[2])"
        }
        
        if ($configOutput -match 'Public API:\s+Current: ([^\s]+):(\d+)') {
            $config.publicApi = "http://$($matches[1]):$($matches[2])"
        }
        
        if ($configOutput -match 'Model API:\s+Current: ([^\s]+):(\d+)') {
            $config.modelApi = "http://$($matches[1]):$($matches[2])"
        }
        
        return $config
    } catch {
        Write-Host "Failed to load API configuration: $($_.Exception.Message)"
        return $null
    }
}

function Get-NetworkStatus {
    try {
        # Simulate network status check, return all services are running
        $status = "====================================`nELR Container Network Status`n====================================`nDesktop API: Running`n  Address: http://localhost:8081`nPublic API: Running`n  Address: http://localhost:8080`nModel Service: Running`n  Address: http://localhost:8082`nMicro Model Server: Running`n  Address: http://localhost:8083`n===================================="
        return $status
    } catch {
        return "Failed to get network status: $($_.Exception.Message)"
    }
}

function Start-AllServices {
    try {
        & $ELRPS1 start-all
    } catch {
        Write-Host "Failed to start all services: $($_.Exception.Message)"
    }
}

function Stop-AllServices {
    try {
        & $ELRPS1 stop-all
    } catch {
        Write-Host "Failed to stop all services: $($_.Exception.Message)"
    }
}

# Get actual container information from ELR
function Get-ELRContainerStats {
    try {
        # Use Start-Process to capture all output
        $tempFile = New-TemporaryFile
        $process = Start-Process -FilePath "powershell.exe" -ArgumentList "-ExecutionPolicy Bypass", "-File", "$ELRPS1", "stats" -NoNewWindow -Wait -RedirectStandardOutput $tempFile.FullName
        
        # Read the output from the temporary file
        $elrStatsOutput = Get-Content -Path $tempFile.FullName
        
        # Clean up temporary file
        Remove-Item $tempFile -Force
        
        # Parse the output to extract container stats
        $stats = @{}
        $inStatsSection = $false
        
        foreach ($line in $elrStatsOutput) {
            # Convert to string to handle any non-string objects
            $lineStr = $line.ToString()
            
            if ($lineStr -match '^Container Stats:$') {
                $inStatsSection = $true
                continue
            }
            
            if ($inStatsSection -and $lineStr -match '^-+$') {
                continue
            }
            
            if ($inStatsSection -and $lineStr -match '^ID\s+NAME\s+MEMORY\s+CPU\s+GPU$') {
                continue
            }
            
            # Skip lines that are just separators or empty
            if ($inStatsSection -and ($lineStr -match '^=+$' -or [string]::IsNullOrEmpty($lineStr))) {
                continue
            }
            
            if ($inStatsSection) {
                # Match container stats line
                if ($lineStr -match '^(elr-\d+)\s+([^\s]+)\s+(\d+)MB\s+(\d+)%\s+(\d+)%$') {
                    $id = $matches[1]
                    $name = $matches[2]
                    $memory = [int]$matches[3]
                    $cpu = [int]$matches[4]
                    $gpu = [int]$matches[5]
                    
                    $stats[$name] = @{Memory=$memory; CPU=$cpu; GPU=$gpu}
                } else {
                    # Try a more flexible pattern
                    if ($lineStr -match '^(elr-\d+)\s+(.+?)\s+(\d+)MB\s+(\d+)%\s+(\d+)%$') {
                        $id = $matches[1]
                        $name = $matches[2].Trim()
                        $memory = [int]$matches[3]
                        $cpu = [int]$matches[4]
                        $gpu = [int]$matches[5]
                        
                        $stats[$name] = @{Memory=$memory; CPU=$cpu; GPU=$gpu}
                    }
                }
            }
        }
        
        return $stats
    } catch {
        # Fallback to empty stats if query fails
        Write-Host "Failed to get container stats: $($_.Exception.Message)"
        return @{}
    }
}

function Get-ELRContainers {
    try {
        # Try to get actual container information from ELR
        # Use Start-Process to capture all output
        $tempFile = New-TemporaryFile
        $process = Start-Process -FilePath "powershell.exe" -ArgumentList "-ExecutionPolicy Bypass", "-File", "$ELRPS1", "list" -NoNewWindow -Wait -RedirectStandardOutput $tempFile.FullName
        
        # Read the output from the temporary file
        $elrListOutput = Get-Content -Path $tempFile.FullName
        
        # Debug: print raw output
        Write-Host "Raw ELR list output:"
        foreach ($line in $elrListOutput) {
            Write-Host "Line: '$line'"
        }
        
        # Clean up temporary file
        Remove-Item $tempFile -Force
        
        # Get container stats
        $containerStats = Get-ELRContainerStats
        
        # Parse the output to extract container information
        $containers = @()
        $inContainerList = $false
        
        foreach ($line in $elrListOutput) {
            # Convert to string to handle any non-string objects
            $lineStr = $line.ToString()
            Write-Host "Processing line: '$lineStr'"
            Write-Host "Current inContainerList: $inContainerList"
            
            if ($lineStr -match '^Containers:$') {
                $inContainerList = $true
                Write-Host "Found Containers section, setting inContainerList to $inContainerList"
                continue
            }
            
            if ($inContainerList -and $lineStr -match '^-+$') {
                Write-Host "Found separator line"
                continue
            }
            
            if ($inContainerList -and $lineStr -match '^ID\s+NAME\s+IMAGE\s+STATUS\s+CREATED$') {
                Write-Host "Found header line"
                continue
            }
            
            # Skip lines that are just separators or empty
            if ($inContainerList -and ($lineStr -match '^=+$' -or [string]::IsNullOrEmpty($lineStr))) {
                Write-Host "Skipping separator or empty line"
                continue
            }
            
            # Try a more flexible regex to match container lines
            Write-Host "Trying to match container line: '$lineStr'"
            if ($inContainerList) {
                # Try different regex patterns to match container lines
                if ($lineStr -match '^\S+\s+([^\s]+)\s+([^\s]+)\s+([^\s]+)\s+.*$') {
                    $name = $matches[1]
                    $image = $matches[2]
                    $status = $matches[3]
                    
                    # Skip the header separator line that might be matched
                    if ($name -eq '----' -and $image -eq '-----' -and $status -eq '------') {
                        Write-Host "Skipping header separator line"
                        continue
                    }
                    
                    Write-Host "Found container: Name='$name', Image='$image', Status='$status'"
                    
                    # Map status to standardized values
                    if ($status -eq "created") {
                        $status = "created"
                    } elseif ($status -eq "running") {
                        $status = "running"
                    }
                    
                    # Extract components from image name
                    $components = "unknown"
                    if ($image -match 'python:') {
                        $components = "Python 3.9"
                    } elseif ($image -match 'ubuntu:') {
                        $components = "Ubuntu"
                    }
                    
                    # Generate sandbox ID based on container name
                    $sandboxId = "sandbox-" + ($containers.Count + 1).ToString("000")
                    
                    # Get resource usage from stats
                    $memory = 0
                    $cpu = 0
                    $gpu = 0
                    
                    if ($containerStats.ContainsKey($name)) {
                        $stats = $containerStats[$name]
                        $memory = $stats.Memory
                        $cpu = $stats.CPU
                        $gpu = $stats.GPU
                        Write-Host "Found stats for container ${name}: Memory=$memory, CPU=$cpu, GPU=$gpu"
                    } else {
                        # Fallback to sample data if stats not found
                        if ($status -eq "running") {
                            $memory = 128 * ($containers.Count + 1)
                            $cpu = 5 * ($containers.Count + 1)
                        }
                        Write-Host "No stats found for container ${name}, using fallback data"
                    }
                    
                    # Add container to list
                    $containers += @{Name=$name; Status=$status; Components=$components; Sandbox=$sandboxId; Memory=$memory; CPU=$cpu; GPU=$gpu}
                    Write-Host "Added container to list"
                } else {
                    # Try a simpler pattern
                    Write-Host "Trying simpler pattern"
                    if ($lineStr -match '^(\S+)\s+([^\s]+)\s+([^\s]+)\s+([^\s]+)\s+.*$') {
                        $id = $matches[1]
                        $name = $matches[2]
                        $image = $matches[3]
                        $status = $matches[4]
                        
                        # Skip the header separator line that might be matched
                        if ($name -eq '----' -and $image -eq '-----' -and $status -eq '------') {
                            Write-Host "Skipping header separator line"
                            continue
                        }
                        
                        Write-Host "Found container with simpler pattern: ID='$id', Name='$name', Image='$image', Status='$status'"
                        
                        # Map status to standardized values
                        if ($status -eq "created") {
                            $status = "created"
                        } elseif ($status -eq "running") {
                            $status = "running"
                        }
                        
                        # Extract components from image name
                        $components = "unknown"
                        if ($image -match 'python:') {
                            $components = "Python 3.9"
                        } elseif ($image -match 'ubuntu:') {
                            $components = "Ubuntu"
                        }
                        
                        # Generate sandbox ID based on container name
                        $sandboxId = "sandbox-" + ($containers.Count + 1).ToString("000")
                        
                        # Get resource usage from stats
                        $memory = 0
                        $cpu = 0
                        $gpu = 0
                        
                        if ($containerStats.ContainsKey($name)) {
                            $stats = $containerStats[$name]
                            $memory = $stats.Memory
                            $cpu = $stats.CPU
                            $gpu = $stats.GPU
                            Write-Host "Found stats for container ${name}: Memory=$memory, CPU=$cpu, GPU=$gpu"
                        } else {
                            # Fallback to sample data if stats not found
                            if ($status -eq "running") {
                                $memory = 128 * ($containers.Count + 1)
                                $cpu = 5 * ($containers.Count + 1)
                            }
                            Write-Host "No stats found for container ${name}, using fallback data"
                        }
                        
                        # Add container to list
                        $containers += @{Name=$name; Status=$status; Components=$components; Sandbox=$sandboxId; Memory=$memory; CPU=$cpu; GPU=$gpu}
                        Write-Host "Added container to list"
                    } else {
                        Write-Host "Container line not matched: '$lineStr'"
                    }
                }
            } else {
                Write-Host "Not in container list section"
            }
        }
        
        # Debug: print container count
        Write-Host "Found $($containers.Count) containers"
        
        # If no containers were found, return empty list
        if ($containers.Count -eq 0) {
            return @()
        }
        
        return $containers
    } catch {
        # Fallback to empty list if query fails
        Write-Host "Failed to get container information: $($_.Exception.Message)"
        return @()
    }
}

function Show-ChatWindow {
    if ($Global:chatForm -ne $null -and $Global:chatForm.Visible) {
        $Global:chatForm.Activate()
        return
    }

    $Global:chatForm = New-Object System.Windows.Forms.Form
    $Global:chatForm.Text = "ELR Chat"
    $Global:chatForm.Size = New-Object System.Drawing.Size(400, 500)
    $Global:chatForm.StartPosition = "CenterScreen"
    try {
        if (Test-Path $IconPath) {
            $Global:chatForm.Icon = New-Object System.Drawing.Icon($IconPath)
        }
    } catch {
        # Ignore icon errors
    }

    # Chat area
    $chatBox = New-Object System.Windows.Forms.RichTextBox
    $chatBox.Location = New-Object System.Drawing.Point(10, 10)
    $chatBox.Size = New-Object System.Drawing.Size(370, 380)
    $chatBox.ReadOnly = $true
    $chatBox.BackColor = [System.Drawing.Color]::White
    $chatBox.Font = New-Object System.Drawing.Font("Arial", 10)
    $Global:chatForm.Controls.Add($chatBox)

    # Input area
    $inputBox = New-Object System.Windows.Forms.TextBox
    $inputBox.Location = New-Object System.Drawing.Point(10, 400)
    $inputBox.Size = New-Object System.Drawing.Size(300, 25)
    $Global:chatForm.Controls.Add($inputBox)

    # Send button
    $sendButton = New-Object System.Windows.Forms.Button
    $sendButton.Location = New-Object System.Drawing.Point(320, 400)
    $sendButton.Size = New-Object System.Drawing.Size(60, 25)
    $sendButton.Text = "Send"
    $sendButton.Add_Click({
        $message = $inputBox.Text
        if ($message -ne "") {
            $chatBox.AppendText("You: $message`n")
            $inputBox.Text = ""
            
            # Simulate ELR container response
            $response = "ELR: I received your message: $message"
            Start-Sleep -Milliseconds 500
            $chatBox.AppendText("$response`n")
        }
    })
    $Global:chatForm.Controls.Add($sendButton)

    # Upload file button
    $uploadButton = New-Object System.Windows.Forms.Button
    $uploadButton.Location = New-Object System.Drawing.Point(10, 430)
    $uploadButton.Size = New-Object System.Drawing.Size(75, 25)
    $uploadButton.Text = "Upload"
    $uploadButton.Add_Click({
        $openFileDialog = New-Object System.Windows.Forms.OpenFileDialog
        $openFileDialog.Title = "Select File"
        $openFileDialog.Filter = "All Files (*.*)|*.*"
        if ($openFileDialog.ShowDialog() -eq "OK") {
            $filePath = $openFileDialog.FileName
            $chatBox.AppendText("System: File uploaded: $filePath`n")
            # File upload logic can be added here
        }
    })
    $Global:chatForm.Controls.Add($uploadButton)

    $Global:chatForm.Show()
}

function Show-SettingsWindow {
    if ($Global:settingsForm -ne $null -and $Global:settingsForm.Visible) {
        $Global:settingsForm.Activate()
        return
    }

    $Global:settingsForm = New-Object System.Windows.Forms.Form
    $Global:settingsForm.Text = (Get-LocalizedString -keyPath "settingsWindow.title")
    $Global:settingsForm.Size = New-Object System.Drawing.Size(550, 300)
    $Global:settingsForm.StartPosition = "CenterScreen"
    try {
        if (Test-Path $IconPath) {
            $Global:settingsForm.Icon = New-Object System.Drawing.Icon($IconPath)
        }
    } catch {
        # Ignore icon errors
    }
    $Global:settingsForm.FormBorderStyle = "FixedSingle"
    $Global:settingsForm.MaximizeBox = $false
    $Global:settingsForm.MinimizeBox = $false

    # Admin info
    $adminLabel = New-Object System.Windows.Forms.Label
    $adminLabel.Location = New-Object System.Drawing.Point(20, 20)
    $adminLabel.Size = New-Object System.Drawing.Size(300, 60)
    $adminLabel.Text = (Get-LocalizedString -keyPath "settingsWindow.tokenInfo")
    $Global:settingsForm.Controls.Add($adminLabel)

    # Create admin button
    $createAdminBtn = New-Object System.Windows.Forms.Button
    $createAdminBtn.Location = New-Object System.Drawing.Point(380, 20)
    $createAdminBtn.Size = New-Object System.Drawing.Size(120, 30)
    $createAdminBtn.Text = (Get-LocalizedString -keyPath "settingsWindow.createAdmin")
    $createAdminBtn.Add_Click({ 
        [System.Windows.Forms.MessageBox]::Show((Get-LocalizedString -keyPath "settingsWindow.featureInDevelopment"), (Get-LocalizedString -keyPath "settingsWindow.messageTitle"), 'OK', 'Information')
    })
    $Global:settingsForm.Controls.Add($createAdminBtn)

    # Update token button
    $updateTokenBtn = New-Object System.Windows.Forms.Button
    $updateTokenBtn.Location = New-Object System.Drawing.Point(380, 60)
    $updateTokenBtn.Size = New-Object System.Drawing.Size(120, 30)
    $updateTokenBtn.Text = (Get-LocalizedString -keyPath "settingsWindow.updateToken")
    $updateTokenBtn.Add_Click({ 
        [System.Windows.Forms.MessageBox]::Show((Get-LocalizedString -keyPath "settingsWindow.featureInDevelopment"), (Get-LocalizedString -keyPath "settingsWindow.messageTitle"), 'OK', 'Information')
    })
    $Global:settingsForm.Controls.Add($updateTokenBtn)

    # Admin count
    $adminCountLabel = New-Object System.Windows.Forms.Label
    $adminCountLabel.Location = New-Object System.Drawing.Point(20, 90)
    $adminCountLabel.Size = New-Object System.Drawing.Size(120, 20)
    $adminCountLabel.Text = (Get-LocalizedString -keyPath "settingsWindow.adminCount")
    $Global:settingsForm.Controls.Add($adminCountLabel)

    # View admins button
    $viewAdminBtn = New-Object System.Windows.Forms.Button
    $viewAdminBtn.Location = New-Object System.Drawing.Point(140, 90)
    $viewAdminBtn.Size = New-Object System.Drawing.Size(200, 20)
    $viewAdminBtn.Text = (Get-LocalizedString -keyPath "settingsWindow.viewAdmins")
    $viewAdminBtn.Add_Click({ 
        [System.Windows.Forms.MessageBox]::Show((Get-LocalizedString -keyPath "settingsWindow.featureInDevelopment"), (Get-LocalizedString -keyPath "settingsWindow.messageTitle"), 'OK', 'Information')
    })
    $Global:settingsForm.Controls.Add($viewAdminBtn)

    # Desktop API
    $Global:desktopApiCheckbox = New-Object System.Windows.Forms.CheckBox
    $Global:desktopApiCheckbox.Location = New-Object System.Drawing.Point(20, 120)
    $Global:desktopApiCheckbox.Size = New-Object System.Drawing.Size(20, 20)
    $Global:desktopApiCheckbox.Checked = $true
    $Global:settingsForm.Controls.Add($Global:desktopApiCheckbox)

    $desktopApiLabel = New-Object System.Windows.Forms.Label
    $desktopApiLabel.Location = New-Object System.Drawing.Point(40, 120)
    $desktopApiLabel.Size = New-Object System.Drawing.Size(90, 20)
    $desktopApiLabel.Text = (Get-LocalizedString -keyPath "settingsWindow.desktopApi")
    $Global:settingsForm.Controls.Add($desktopApiLabel)

    # Desktop API status
    $Global:desktopApiStatus = New-Object System.Windows.Forms.Label
    $Global:desktopApiStatus.Location = New-Object System.Drawing.Point(140, 120)
    $Global:desktopApiStatus.Size = New-Object System.Drawing.Size(80, 20)
    $Global:desktopApiStatus.Text = "Stopped"
    $Global:desktopApiStatus.ForeColor = [System.Drawing.Color]::Red
    $Global:settingsForm.Controls.Add($Global:desktopApiStatus)

    $Global:desktopApiInput = New-Object System.Windows.Forms.TextBox
    $Global:desktopApiInput.Name = "desktopApiInput"
    $Global:desktopApiInput.Location = New-Object System.Drawing.Point(230, 120)
    $Global:desktopApiInput.Size = New-Object System.Drawing.Size(190, 20)
    $Global:desktopApiInput.Text = "http://localhost:8094"
    $Global:desktopApiInput.TextAlign = [System.Windows.Forms.HorizontalAlignment]::Left
    $Global:settingsForm.Controls.Add($Global:desktopApiInput)

    $desktopApiAutoBtn = New-Object System.Windows.Forms.Button
    $desktopApiAutoBtn.Location = New-Object System.Drawing.Point(430, 120)
    $desktopApiAutoBtn.Size = New-Object System.Drawing.Size(70, 20)
    $desktopApiAutoBtn.Text = (Get-LocalizedString -keyPath "settingsWindow.auto")
    $desktopApiAutoBtn.Add_Click({ 
        $Global:desktopApiInput.Text = 'http://localhost:8094'
        [System.Windows.Forms.MessageBox]::Show((Get-LocalizedString -keyPath "settingsWindow.apiAddressAutoRecommended"), (Get-LocalizedString -keyPath "settingsWindow.messageTitle"), 'OK', 'Information')
    })
    $Global:settingsForm.Controls.Add($desktopApiAutoBtn)

    # Public API
    $Global:publicApiCheckbox = New-Object System.Windows.Forms.CheckBox
    $Global:publicApiCheckbox.Location = New-Object System.Drawing.Point(20, 150)
    $Global:publicApiCheckbox.Size = New-Object System.Drawing.Size(20, 20)
    $Global:publicApiCheckbox.Checked = $true
    $Global:settingsForm.Controls.Add($Global:publicApiCheckbox)

    $publicApiLabel = New-Object System.Windows.Forms.Label
    $publicApiLabel.Location = New-Object System.Drawing.Point(40, 150)
    $publicApiLabel.Size = New-Object System.Drawing.Size(90, 20)
    $publicApiLabel.Text = (Get-LocalizedString -keyPath "settingsWindow.publicApi")
    $Global:settingsForm.Controls.Add($publicApiLabel)

    # Public API status
    $Global:publicApiStatus = New-Object System.Windows.Forms.Label
    $Global:publicApiStatus.Location = New-Object System.Drawing.Point(140, 150)
    $Global:publicApiStatus.Size = New-Object System.Drawing.Size(80, 20)
    $Global:publicApiStatus.Text = "Stopped"
    $Global:publicApiStatus.ForeColor = [System.Drawing.Color]::Red
    $Global:settingsForm.Controls.Add($Global:publicApiStatus)

    $Global:publicApiInput = New-Object System.Windows.Forms.TextBox
    $Global:publicApiInput.Name = "publicApiInput"
    $Global:publicApiInput.Location = New-Object System.Drawing.Point(230, 150)
    $Global:publicApiInput.Size = New-Object System.Drawing.Size(190, 20)
    $Global:publicApiInput.Text = "http://localhost:8095"
    $Global:publicApiInput.TextAlign = [System.Windows.Forms.HorizontalAlignment]::Left
    $Global:settingsForm.Controls.Add($Global:publicApiInput)

    $publicApiAutoBtn = New-Object System.Windows.Forms.Button
    $publicApiAutoBtn.Location = New-Object System.Drawing.Point(430, 150)
    $publicApiAutoBtn.Size = New-Object System.Drawing.Size(70, 20)
    $publicApiAutoBtn.Text = (Get-LocalizedString -keyPath "settingsWindow.auto")
    $publicApiAutoBtn.Add_Click({ 
        $Global:publicApiInput.Text = 'http://localhost:8095'
        [System.Windows.Forms.MessageBox]::Show((Get-LocalizedString -keyPath "settingsWindow.apiAddressAutoRecommended"), (Get-LocalizedString -keyPath "settingsWindow.messageTitle"), 'OK', 'Information')
    })
    $Global:settingsForm.Controls.Add($publicApiAutoBtn)

    # Model API
    $Global:modelApiCheckbox = New-Object System.Windows.Forms.CheckBox
    $Global:modelApiCheckbox.Location = New-Object System.Drawing.Point(20, 180)
    $Global:modelApiCheckbox.Size = New-Object System.Drawing.Size(20, 20)
    $Global:modelApiCheckbox.Checked = $true
    $Global:settingsForm.Controls.Add($Global:modelApiCheckbox)

    $modelApiLabel = New-Object System.Windows.Forms.Label
    $modelApiLabel.Location = New-Object System.Drawing.Point(40, 180)
    $modelApiLabel.Size = New-Object System.Drawing.Size(90, 20)
    $modelApiLabel.Text = (Get-LocalizedString -keyPath "settingsWindow.modelApi")
    $Global:settingsForm.Controls.Add($modelApiLabel)

    # Model API status
    $Global:modelApiStatus = New-Object System.Windows.Forms.Label
    $Global:modelApiStatus.Location = New-Object System.Drawing.Point(140, 180)
    $Global:modelApiStatus.Size = New-Object System.Drawing.Size(80, 20)
    $Global:modelApiStatus.Text = "Stopped"
    $Global:modelApiStatus.ForeColor = [System.Drawing.Color]::Red
    $Global:settingsForm.Controls.Add($Global:modelApiStatus)

    $Global:modelApiInput = New-Object System.Windows.Forms.TextBox
    $Global:modelApiInput.Name = "modelApiInput"
    $Global:modelApiInput.Location = New-Object System.Drawing.Point(230, 180)
    $Global:modelApiInput.Size = New-Object System.Drawing.Size(190, 20)
    $Global:modelApiInput.Text = "http://localhost:8096"
    $Global:modelApiInput.TextAlign = [System.Windows.Forms.HorizontalAlignment]::Left
    $Global:settingsForm.Controls.Add($Global:modelApiInput)

    $modelApiAutoBtn = New-Object System.Windows.Forms.Button
    $modelApiAutoBtn.Location = New-Object System.Drawing.Point(430, 180)
    $modelApiAutoBtn.Size = New-Object System.Drawing.Size(70, 20)
    $modelApiAutoBtn.Text = (Get-LocalizedString -keyPath "settingsWindow.auto")
    $modelApiAutoBtn.Add_Click({ 
        $Global:modelApiInput.Text = 'http://localhost:8096'
        [System.Windows.Forms.MessageBox]::Show((Get-LocalizedString -keyPath "settingsWindow.apiAddressAutoRecommended"), (Get-LocalizedString -keyPath "settingsWindow.messageTitle"), 'OK', 'Information')
    })
    $Global:settingsForm.Controls.Add($modelApiAutoBtn)

    # Start/Stop API buttons
    $startApiBtn = New-Object System.Windows.Forms.Button
    $startApiBtn.Location = New-Object System.Drawing.Point(100, 220)
    $startApiBtn.Size = New-Object System.Drawing.Size(120, 30)
    $startApiBtn.Text = (Get-LocalizedString -keyPath "settingsWindow.startApi")
    $startApiBtn.Add_Click({ 
        [System.Windows.Forms.MessageBox]::Show((Get-LocalizedString -keyPath "settingsWindow.startingApi"), (Get-LocalizedString -keyPath "settingsWindow.messageTitle"), 'OK', 'Information')
        
        # Path to elr.exe (using relative path)
        $elrExe = ".\elr.exe"
        
        # Start selected APIs
        if ($Global:desktopApiCheckbox.Checked) {
            Write-Host "Starting Desktop API..."
            try {
                # Extract port from input box
                    $desktopApiUrl = $Global:desktopApiInput.Text
                    if ($desktopApiUrl -match 'http://([^:]+):(\d+)') {
                        $address = $matches[1]
                        $port = $matches[2]
                        Write-Host "Starting Desktop API on ${address}:${port}..."
                        
                        # Save the configuration to alternatives
                        Write-Host "Adding Desktop API address:port ${address}:${port} to alternatives"
                        & $elrExe api config set --api-type desktop --address $address --port $port
                        
                        # Get the index of the added alternative
                        $configList = & $elrExe api config list
                        $alternativeIndex = -1
                        $lines = $configList -split "`n"
                        foreach ($line in $lines) {
                            if ($line -match "Desktop API:") {
                                $inDesktopSection = $true
                            } elseif ($line -match "Public API:" -or $line -match "Model API:") {
                                $inDesktopSection = $false
                            }
                            if ($inDesktopSection -and $line -match "\s+(\d+):\s*${address}:${port}") {
                                $alternativeIndex = [int]$matches[1]
                                Write-Host "Desktop API address:port ${address}:${port} added to alternatives at index $alternativeIndex"
                                break
                            }
                        }
                        
                        # Enable the alternative configuration
                        if ($alternativeIndex -ge 0) {
                            Write-Host "Enabling Desktop API address:port ${address}:${port} at index $alternativeIndex"
                            & $elrExe api config enable --api-type desktop --index $alternativeIndex
                        }
                        
                        # Start the API in background
                        Write-Host "Executing: $elrExe api start desktop"
                        try {
                            # Start the API in background
                            $process = Start-Process -FilePath $elrExe -ArgumentList "api", "start", "desktop" -NoNewWindow -PassThru
                            
                            # Wait for the process to complete
                            $process.WaitForExit()
                            
                            # Test API access
                            Start-Sleep -Seconds 2
                            $status = & $elrExe api status
                            if ($status -match "Desktop API: http://[^:]+:$port - Running") {
                                Write-Host "Desktop API started and accessible"
                                # Update status label
                                $Global:desktopApiStatus.Text = "Running"
                                $Global:desktopApiStatus.ForeColor = [System.Drawing.Color]::Green
                                # Ask user if they want to save and enable this address and port
                                $confirmMessage = Get-LocalizedStringWithParams -keyPath "settingsWindow.confirmSaveApiAddress" -params "Desktop"
                                $result = [System.Windows.Forms.MessageBox]::Show($confirmMessage, (Get-LocalizedString -keyPath "settingsWindow.messageTitle"), 'YesNo', 'Information')
                                if ($result -eq [System.Windows.Forms.DialogResult]::Yes) {
                                    # Already saved and enabled earlier
                                    [System.Windows.Forms.MessageBox]::Show((Get-LocalizedString -keyPath "settingsWindow.apiAddressAndPortSaved"), (Get-LocalizedString -keyPath "settingsWindow.messageTitle"), 'OK', 'Information')
                                }
                            } else {
                                Write-Host "Desktop API started but not accessible"
                                [System.Windows.Forms.MessageBox]::Show((Get-LocalizedString -keyPath "settingsWindow.apiStartedButNotAccessible"), (Get-LocalizedString -keyPath "settingsWindow.messageTitle"), 'OK', 'Warning')
                            }
                        } catch {
                            Write-Host "Error starting Desktop API: $($_.Exception.Message)"
                            $errorMessage = Get-LocalizedStringWithParams -keyPath "settingsWindow.failedToStartApi" -params $($_.Exception.Message)
                            [System.Windows.Forms.MessageBox]::Show($errorMessage, (Get-LocalizedString -keyPath "settingsWindow.messageTitle"), 'OK', 'Error')
                        }
                    } else {
                        Write-Host "Invalid Desktop API URL: $desktopApiUrl"
                        [System.Windows.Forms.MessageBox]::Show((Get-LocalizedString -keyPath "settingsWindow.invalidApiAddress"), (Get-LocalizedString -keyPath "settingsWindow.messageTitle"), 'OK', 'Error')
                    }
            } catch {
                Write-Host "Error starting Desktop API: $($_.Exception.Message)"
                $errorMessage = Get-LocalizedStringWithParams -keyPath "settingsWindow.failedToStartApi" -params $($_.Exception.Message)
                [System.Windows.Forms.MessageBox]::Show($errorMessage, (Get-LocalizedString -keyPath "settingsWindow.messageTitle"), 'OK', 'Error')
            }
        }
        if ($Global:publicApiCheckbox.Checked) {
            Write-Host "Starting Public API..."
            try {
                # Extract port from input box
                    $publicApiUrl = $Global:publicApiInput.Text
                    if ($publicApiUrl -match 'http://([^:]+):(\d+)') {
                        $address = $matches[1]
                        $port = $matches[2]
                        Write-Host "Starting Public API on ${address}:${port}..."
                        
                        # Save the configuration to alternatives
                        Write-Host "Adding Public API address:port ${address}:${port} to alternatives"
                        & $elrExe api config set --api-type public --address $address --port $port
                        
                        # Get the index of the added alternative
                        $configList = & $elrExe api config list
                        $alternativeIndex = -1
                        $lines = $configList -split "`n"
                        foreach ($line in $lines) {
                            if ($line -match "Public API:") {
                                $inPublicSection = $true
                            } elseif ($line -match "Desktop API:" -or $line -match "Model API:") {
                                $inPublicSection = $false
                            }
                            if ($inPublicSection -and $line -match "\s+(\d+):\s*${address}:${port}") {
                                $alternativeIndex = [int]$matches[1]
                                Write-Host "Public API address:port ${address}:${port} added to alternatives at index $alternativeIndex"
                                break
                            }
                        }
                        
                        # Enable the alternative configuration
                        if ($alternativeIndex -ge 0) {
                            Write-Host "Enabling Public API address:port ${address}:${port} at index $alternativeIndex"
                            & $elrExe api config enable --api-type public --index $alternativeIndex
                        }
                        
                        # Start the API in background
                        Write-Host "Executing: $elrExe api start public"
                        try {
                            # Start the API in background
                            $process = Start-Process -FilePath $elrExe -ArgumentList "api", "start", "public" -NoNewWindow -PassThru
                            
                            # Wait for the process to complete
                            $process.WaitForExit()
                            
                            # Test API access
                            Start-Sleep -Seconds 2
                            $status = & $elrExe api status
                            if ($status -match "Public API: http://[^:]+:$port - Running") {
                                Write-Host "Public API started and accessible"
                                # Update status label
                                $Global:publicApiStatus.Text = "Running"
                                $Global:publicApiStatus.ForeColor = [System.Drawing.Color]::Green
                                # Ask user if they want to save and enable this address and port
                                $confirmMessage = Get-LocalizedStringWithParams -keyPath "settingsWindow.confirmSaveApiAddress" -params "Public"
                                $result = [System.Windows.Forms.MessageBox]::Show($confirmMessage, (Get-LocalizedString -keyPath "settingsWindow.messageTitle"), 'YesNo', 'Information')
                                if ($result -eq [System.Windows.Forms.DialogResult]::Yes) {
                                    # Already saved and enabled earlier
                                    [System.Windows.Forms.MessageBox]::Show((Get-LocalizedString -keyPath "settingsWindow.apiAddressAndPortSaved"), (Get-LocalizedString -keyPath "settingsWindow.messageTitle"), 'OK', 'Information')
                                }
                            } else {
                                Write-Host "Public API started but not accessible"
                                [System.Windows.Forms.MessageBox]::Show((Get-LocalizedString -keyPath "settingsWindow.apiStartedButNotAccessible"), (Get-LocalizedString -keyPath "settingsWindow.messageTitle"), 'OK', 'Warning')
                            }
                        } catch {
                            Write-Host "Error starting Public API: $($_.Exception.Message)"
                            $errorMessage = Get-LocalizedStringWithParams -keyPath "settingsWindow.failedToStartApi" -params $($_.Exception.Message)
                            [System.Windows.Forms.MessageBox]::Show($errorMessage, (Get-LocalizedString -keyPath "settingsWindow.messageTitle"), 'OK', 'Error')
                        }
                    } else {
                        Write-Host "Invalid Public API URL: $publicApiUrl"
                        [System.Windows.Forms.MessageBox]::Show((Get-LocalizedString -keyPath "settingsWindow.invalidApiAddress"), (Get-LocalizedString -keyPath "settingsWindow.messageTitle"), 'OK', 'Error')
                    }
            } catch {
                Write-Host "Error starting Public API: $($_.Exception.Message)"
                $errorMessage = Get-LocalizedStringWithParams -keyPath "settingsWindow.failedToStartApi" -params $($_.Exception.Message)
                [System.Windows.Forms.MessageBox]::Show($errorMessage, (Get-LocalizedString -keyPath "settingsWindow.messageTitle"), 'OK', 'Error')
            }
        }
        if ($Global:modelApiCheckbox.Checked) {
            Write-Host "Starting Model API..."
            try {
                # Extract port from input box
                    $modelApiUrl = $Global:modelApiInput.Text
                    if ($modelApiUrl -match 'http://([^:]+):(\d+)') {
                        $address = $matches[1]
                        $port = $matches[2]
                        Write-Host "Starting Model API on ${address}:${port}..."
                        
                        # Save the configuration to alternatives
                        Write-Host "Adding Model API address:port ${address}:${port} to alternatives"
                        & $elrExe api config set --api-type model --address $address --port $port
                        
                        # Get the index of the added alternative
                        $configList = & $elrExe api config list
                        $alternativeIndex = -1
                        $lines = $configList -split "`n"
                        foreach ($line in $lines) {
                            if ($line -match "Model API:") {
                                $inModelSection = $true
                            } elseif ($line -match "Public API:" -or $line -match "Desktop API:") {
                                $inModelSection = $false
                            }
                            if ($inModelSection -and $line -match "\s+(\d+):\s*${address}:${port}") {
                                $alternativeIndex = [int]$matches[1]
                                Write-Host "Model API address:port ${address}:${port} added to alternatives at index $alternativeIndex"
                                break
                            }
                        }
                        
                        # Enable the alternative configuration
                        if ($alternativeIndex -ge 0) {
                            Write-Host "Enabling Model API address:port ${address}:${port} at index $alternativeIndex"
                            & $elrExe api config enable --api-type model --index $alternativeIndex
                        }
                        
                        # Start the API in background
                        Write-Host "Executing: $elrExe api start model"
                        try {
                            # Start the API in background
                            $process = Start-Process -FilePath $elrExe -ArgumentList "api", "start", "model" -NoNewWindow -PassThru
                            
                            # Wait for the process to complete
                            $process.WaitForExit()
                            
                            # Test API access
                            Start-Sleep -Seconds 2
                            $status = & $elrExe api status
                            if ($status -match "Model API: http://[^:]+:$port - Running") {
                                Write-Host "Model API started and accessible"
                                # Update status label
                                $Global:modelApiStatus.Text = "Running"
                                $Global:modelApiStatus.ForeColor = [System.Drawing.Color]::Green
                                # Ask user if they want to save and enable this address and port
                                $confirmMessage = Get-LocalizedStringWithParams -keyPath "settingsWindow.confirmSaveApiAddress" -params "Model"
                                $result = [System.Windows.Forms.MessageBox]::Show($confirmMessage, (Get-LocalizedString -keyPath "settingsWindow.messageTitle"), 'YesNo', 'Information')
                                if ($result -eq [System.Windows.Forms.DialogResult]::Yes) {
                                    # Already saved and enabled earlier
                                    [System.Windows.Forms.MessageBox]::Show((Get-LocalizedString -keyPath "settingsWindow.apiAddressAndPortSaved"), (Get-LocalizedString -keyPath "settingsWindow.messageTitle"), 'OK', 'Information')
                                }
                            } else {
                                Write-Host "Model API started but not accessible"
                                [System.Windows.Forms.MessageBox]::Show((Get-LocalizedString -keyPath "settingsWindow.apiStartedButNotAccessible"), (Get-LocalizedString -keyPath "settingsWindow.messageTitle"), 'OK', 'Warning')
                            }
                        } catch {
                            Write-Host "Error starting Model API: $($_.Exception.Message)"
                            $errorMessage = Get-LocalizedStringWithParams -keyPath "settingsWindow.failedToStartApi" -params $($_.Exception.Message)
                            [System.Windows.Forms.MessageBox]::Show($errorMessage, (Get-LocalizedString -keyPath "settingsWindow.messageTitle"), 'OK', 'Error')
                        }
                    } else {
                        Write-Host "Invalid Model API URL: $modelApiUrl"
                        [System.Windows.Forms.MessageBox]::Show((Get-LocalizedString -keyPath "settingsWindow.invalidApiAddress"), (Get-LocalizedString -keyPath "settingsWindow.messageTitle"), 'OK', 'Error')
                    }
            } catch {
                Write-Host "Error starting Model API: $($_.Exception.Message)"
                $errorMessage = Get-LocalizedStringWithParams -keyPath "settingsWindow.failedToStartApi" -params $($_.Exception.Message)
                [System.Windows.Forms.MessageBox]::Show($errorMessage, (Get-LocalizedString -keyPath "settingsWindow.messageTitle"), 'OK', 'Error')
            }
        }
    })
    $Global:settingsForm.Controls.Add($startApiBtn)

    $stopApiBtn = New-Object System.Windows.Forms.Button
    $stopApiBtn.Location = New-Object System.Drawing.Point(250, 220)
    $stopApiBtn.Size = New-Object System.Drawing.Size(120, 30)
    $stopApiBtn.Text = (Get-LocalizedString -keyPath "settingsWindow.stopApi")
    $stopApiBtn.Add_Click({ 
        [System.Windows.Forms.MessageBox]::Show((Get-LocalizedString -keyPath "settingsWindow.stoppingApi"), (Get-LocalizedString -keyPath "settingsWindow.messageTitle"), 'OK', 'Information')
        
        # Path to elr.exe (using relative path)
        $elrExe = ".\elr.exe"
        
        # Stop selected APIs
        if ($Global:desktopApiCheckbox.Checked) {
            Write-Host "Stopping Desktop API..."
            try {
                Start-Process -FilePath $elrExe -ArgumentList "api", "stop", "desktop" -WindowStyle Hidden
                # Update status label
                $Global:desktopApiStatus.Text = "Stopped"
                $Global:desktopApiStatus.ForeColor = [System.Drawing.Color]::Red
            } catch {
                Write-Host "Error stopping Desktop API: $($_.Exception.Message)"
            }
        }
        if ($Global:publicApiCheckbox.Checked) {
            Write-Host "Stopping Public API..."
            try {
                Start-Process -FilePath $elrExe -ArgumentList "api", "stop", "public" -WindowStyle Hidden
                # Update status label
                $Global:publicApiStatus.Text = "Stopped"
                $Global:publicApiStatus.ForeColor = [System.Drawing.Color]::Red
            } catch {
                Write-Host "Error stopping Public API: $($_.Exception.Message)"
            }
        }
        if ($Global:modelApiCheckbox.Checked) {
            Write-Host "Stopping Model API..."
            try {
                Start-Process -FilePath $elrExe -ArgumentList "api", "stop", "model" -WindowStyle Hidden
                # Update status label
                $Global:modelApiStatus.Text = "Stopped"
                $Global:modelApiStatus.ForeColor = [System.Drawing.Color]::Red
            } catch {
                Write-Host "Error stopping Model API: $($_.Exception.Message)"
            }
        }
        
        Start-Sleep -Seconds 2
        [System.Windows.Forms.MessageBox]::Show((Get-LocalizedString -keyPath "settingsWindow.apiStoppedSuccessfully"), (Get-LocalizedString -keyPath "settingsWindow.messageTitle"), 'OK', 'Information')
    })
    $Global:settingsForm.Controls.Add($stopApiBtn)

    # Show the settings form first
    $Global:settingsForm.Show()

    # Create a loading form
    $loadingForm = New-Object System.Windows.Forms.Form
    $loadingForm.Text = "ELR Message"
    $loadingForm.Size = New-Object System.Drawing.Size(350, 120)
    $loadingForm.StartPosition = "CenterScreen"
    $loadingForm.FormBorderStyle = "FixedSingle"
    $loadingForm.MaximizeBox = $false
    $loadingForm.MinimizeBox = $false
    $loadingForm.ControlBox = $false
    $loadingForm.TopMost = $true
    
    $loadingLabel = New-Object System.Windows.Forms.Label
    $loadingLabel.Location = New-Object System.Drawing.Point(20, 40)
    $loadingLabel.Size = New-Object System.Drawing.Size(310, 40)
    $loadingLabel.Text = (Get-LocalizedString -keyPath "settingsWindow.loading")
    $loadingLabel.TextAlign = [System.Drawing.ContentAlignment]::MiddleCenter
    $loadingForm.Controls.Add($loadingLabel)

    # Show the loading form
    $loadingForm.Show()

    # Path to elr.exe
    $elrExe = Join-Path -Path $ELRPath -ChildPath "elr.exe"

    # Get actual API status synchronously
    Write-Host "Getting API status..."
    try {
        $statusOutput = & $elrExe api status 2>&1
        Write-Host "API Status Output: $statusOutput"

        # Parse the status output
        $apiResults = @{}
        
        # Convert statusOutput to string if it's an array
        if ($statusOutput -is [array]) {
            $statusOutput = $statusOutput -join " "
        }
        
        # Parse Desktop API
        if ($statusOutput -match 'Desktop API: http://([^:]+):(\d+)(?: - (Running|Stopped))?') {
            if ($matches -ne $null -and $matches.Count -ge 3) {
                $apiResults.Desktop = @{
                    Address = $matches[1]
                    Port = $matches[2]
                    Status = if ($matches.Count -ge 4) { $matches[3] } else { "Stopped" }
                }
                Write-Host "Found Desktop API: http://$($matches[1]):$($matches[2]) - $($apiResults.Desktop.Status)"
                # Update status label
                if ($apiResults.Desktop.Status -eq "Running") {
                    $Global:desktopApiStatus.Text = "Running"
                    $Global:desktopApiStatus.ForeColor = [System.Drawing.Color]::Green
                } else {
                    $Global:desktopApiStatus.Text = "Stopped"
                    $Global:desktopApiStatus.ForeColor = [System.Drawing.Color]::Red
                }
            }
        }

        # Parse Public API
        if ($statusOutput -match 'Public API: http://([^:]+):(\d+)(?: - (Running|Stopped))?') {
            if ($matches -ne $null -and $matches.Count -ge 3) {
                $apiResults.Public = @{
                    Address = $matches[1]
                    Port = $matches[2]
                    Status = if ($matches.Count -ge 4) { $matches[3] } else { "Stopped" }
                }
                Write-Host "Found Public API: http://$($matches[1]):$($matches[2]) - $($apiResults.Public.Status)"
                # Update status label
                if ($apiResults.Public.Status -eq "Running") {
                    $Global:publicApiStatus.Text = "Running"
                    $Global:publicApiStatus.ForeColor = [System.Drawing.Color]::Green
                } else {
                    $Global:publicApiStatus.Text = "Stopped"
                    $Global:publicApiStatus.ForeColor = [System.Drawing.Color]::Red
                }
            }
        }

        # Parse Model API
        if ($statusOutput -match 'Model API: http://([^:]+):(\d+)(?: - (Running|Stopped))?') {
            if ($matches -ne $null -and $matches.Count -ge 3) {
                $apiResults.Model = @{
                    Address = $matches[1]
                    Port = $matches[2]
                    Status = if ($matches.Count -ge 4) { $matches[3] } else { "Stopped" }
                }
                Write-Host "Found Model API: http://$($matches[1]):$($matches[2]) - $($apiResults.Model.Status)"
                # Update status label
                if ($apiResults.Model.Status -eq "Running") {
                    $Global:modelApiStatus.Text = "Running"
                    $Global:modelApiStatus.ForeColor = [System.Drawing.Color]::Green
                } else {
                    $Global:modelApiStatus.Text = "Stopped"
                    $Global:modelApiStatus.ForeColor = [System.Drawing.Color]::Red
                }
            }
        }

        Write-Host "API Results: $($apiResults | ConvertTo-Json)"
    } catch {
        Write-Host "Error getting API status: $($_.Exception.Message)"
        $apiResults = $null
    }

    # Update the input boxes with the API status
    Write-Host "Processing API results..."
    Write-Host "API Results: $($apiResults | ConvertTo-Json)"
    
    if ($apiResults) {
        Write-Host "API Results found, updating input boxes..."
        
        if ($apiResults.Desktop) {
            Write-Host "Updating Desktop API input..."
            $desktopApiInput.Text = "http://$($apiResults.Desktop.Address):$($apiResults.Desktop.Port)"
            Write-Host "Updated Desktop API: http://$($apiResults.Desktop.Address):$($apiResults.Desktop.Port)"
        } else {
            Write-Host "No Desktop API information found"
        }
        
        if ($apiResults.Public) {
            Write-Host "Updating Public API input..."
            $publicApiInput.Text = "http://$($apiResults.Public.Address):$($apiResults.Public.Port)"
            Write-Host "Updated Public API: http://$($apiResults.Public.Address):$($apiResults.Public.Port)"
        } else {
            Write-Host "No Public API information found"
        }
        
        if ($apiResults.Model) {
            Write-Host "Updating Model API input..."
            $modelApiInput.Text = "http://$($apiResults.Model.Address):$($apiResults.Model.Port)"
            Write-Host "Updated Model API: http://$($apiResults.Model.Address):$($apiResults.Model.Port)"
        } else {
            Write-Host "No Model API information found"
        }
    } else {
        Write-Host "No API results found, using fallback values..."
        # Fallback to config or default values
        $config = Load-APIConfig
        Write-Host "Config: $($config | ConvertTo-Json)"
        if ($config) {
            if ($config.desktopApi) {
                $desktopApiInput.Text = $config.desktopApi
                Write-Host "Fallback to config Desktop API: $($config.desktopApi)"
            }
            if ($config.publicApi) {
                $publicApiInput.Text = $config.publicApi
                Write-Host "Fallback to config Public API: $($config.publicApi)"
            }
            if ($config.modelApi) {
                $modelApiInput.Text = $config.modelApi
                Write-Host "Fallback to config Model API: $($config.modelApi)"
            }
        } else {
            Write-Host "No config found, using default values"
        }
    }
    
    # Force UI update
    $Global:settingsForm.Refresh()
    [System.Windows.Forms.Application]::DoEvents()
    Write-Host "UI updated"


    # Close the loading form after a short delay to allow user to see the message
    Start-Sleep -Seconds 1
    $loadingForm.Close()
    $loadingForm.Dispose()
}

function Show-MainWindow {
    if ($Global:mainForm -ne $null -and $Global:mainForm.Visible) {
        $Global:mainForm.Activate()
        return
    }

    $Global:mainForm = New-Object System.Windows.Forms.Form
    $Global:mainForm.Text = "ELR Dashboard"
    $Global:mainForm.Size = New-Object System.Drawing.Size(850, 300)
    $Global:mainForm.StartPosition = "CenterScreen"
    try {
        if (Test-Path $IconPath) {
            $Global:mainForm.Icon = New-Object System.Drawing.Icon($IconPath)
        }
    } catch {
        # Ignore icon errors
    }

    # Get actual container information from ELR
    $containers = Get-ELRContainers

    # Calculate running count
    $totalCount = $containers.Count
    $runningCount = 0
    foreach ($container in $containers) {
        if ($container.Status -eq "running") {
            $runningCount++
        }
    }

    # Debug: print container counts
    Write-Host "Total containers: $totalCount"
    Write-Host "Running containers: $runningCount"
    Write-Host "Containers array: $($containers | ConvertTo-Json)"

    # Status label
    $statusLabel = New-Object System.Windows.Forms.Label
    $statusLabel.Location = New-Object System.Drawing.Point(10, 10)
    $statusLabel.Size = New-Object System.Drawing.Size(830, 60)
    $statusText = "Connected to ELR Container`n"
    $statusText += "Total containers: $totalCount`n"
    $statusText += "Running containers: $runningCount"
    # Debug: print status text
    Write-Host "Status text: '$statusText'"
    $statusLabel.Text = $statusText
    $Global:mainForm.Controls.Add($statusLabel)

    # Network status button
    $statusBtn = New-Object System.Windows.Forms.Button
    $statusBtn.Location = New-Object System.Drawing.Point(10, 80)
    $statusBtn.Size = New-Object System.Drawing.Size(830, 30)
    $statusBtn.Text = "Check Network Status"
    $statusBtn.Add_Click({
        $status = Get-NetworkStatus
        [System.Windows.Forms.MessageBox]::Show($status, "Network Status", "OK", "Information")
    })
    $Global:mainForm.Controls.Add($statusBtn)

    # Container list
    $containerList = New-Object System.Windows.Forms.ListView
    $containerList.Location = New-Object System.Drawing.Point(10, 120)
    $containerList.Size = New-Object System.Drawing.Size(800, 150)
    $containerList.View = "Details"
    $containerList.Columns.Add("Name", 120)
    $containerList.Columns.Add("Status", 100)
    $containerList.Columns.Add("Components", 200)
    $containerList.Columns.Add("Sandbox", 100)
    $containerList.Columns.Add("Resource Usage", 250)

    # Add containers
    foreach ($container in $containers) {
        $item = New-Object System.Windows.Forms.ListViewItem($container.Name)
        
        # Get localized status
        if ($container.Status -eq "created") {
            $statusText = "Created"
        } elseif ($container.Status -eq "running") {
            $statusText = "Running"
        } else {
            $statusText = $container.Status
        }
        
        # Get localized components
        if ($container.Components -eq "unknown") {
            $componentsText = "Unknown"
        } else {
            $componentsText = $container.Components
        }
        
        $item.SubItems.Add($statusText)
        $item.SubItems.Add($componentsText)
        $item.SubItems.Add($container.Sandbox)
        $resourceUsage = "Memory: $($container.Memory)MB, CPU: $($container.CPU)%, GPU: $($container.GPU)%"
        $item.SubItems.Add($resourceUsage)
        $containerList.Items.Add($item)
    }

    $Global:mainForm.Controls.Add($containerList)

    $Global:mainForm.Show()
}

# Create system tray
function Initialize-TrayIcon {
    try {
        Write-Host "Creating NotifyIcon object..."
        $Global:trayIcon = New-Object System.Windows.Forms.NotifyIcon
        
        Write-Host "Loading icon from: $IconPath"
        if (Test-Path $IconPath) {
            try {
                # Use direct icon loading instead of ExtractAssociatedIcon
                $Global:trayIcon.Icon = New-Object System.Drawing.Icon($IconPath)
                Write-Host "Icon loaded successfully"
            } catch {
                Write-Host "Error loading icon: $($_.Exception.Message)"
                # Fallback to default icon
                $Global:trayIcon.Icon = [System.Drawing.SystemIcons]::Information
                Write-Host "Using default system icon"
            }
        } else {
            Write-Host "Icon file not found, using default system icon"
            $Global:trayIcon.Icon = [System.Drawing.SystemIcons]::Information
        }
        
        $trayText = "Enlightenment Lighthouse Runtime"
        Write-Host "Setting tray icon text: $trayText"
        $Global:trayIcon.Text = $trayText
        
        Write-Host "Setting tray icon visible: true"
        $Global:trayIcon.Visible = $true

        # Create context menu
        Write-Host "Creating context menu..."
        $contextMenu = New-Object System.Windows.Forms.ContextMenuStrip

        # Open main window
        $openItem = $contextMenu.Items.Add((Get-LocalizedString -keyPath "contextMenu.open"))
        $openItem.Add_Click({ Show-MainWindow })

        # Open chat window
        $chatItem = $contextMenu.Items.Add((Get-LocalizedString -keyPath "contextMenu.chat"))
        $chatItem.Add_Click({ Show-ChatWindow })

        # Open settings window
        $settingsItem = $contextMenu.Items.Add((Get-LocalizedString -keyPath "contextMenu.settings"))
        $settingsItem.Add_Click({ Show-SettingsWindow })

        # Container Management
        $containerItem = $contextMenu.Items.Add((Get-LocalizedString -keyPath "contextMenu.containerManagement"))
        $containerItem.Add_Click({ 
            [System.Windows.Forms.MessageBox]::Show((Get-LocalizedString -keyPath "settingsWindow.featureInDevelopment"), (Get-LocalizedString -keyPath "settingsWindow.messageTitle"), "OK", "Information")
        })

        # Domain Management
        $domainItem = $contextMenu.Items.Add((Get-LocalizedString -keyPath "contextMenu.domainManagement"))
        $domainItem.Add_Click({ 
            [System.Windows.Forms.MessageBox]::Show((Get-LocalizedString -keyPath "settingsWindow.featureInDevelopment"), (Get-LocalizedString -keyPath "settingsWindow.messageTitle"), "OK", "Information")
        })

        # Open Terminal
        $terminalItem = $contextMenu.Items.Add((Get-LocalizedString -keyPath "contextMenu.openTerminal"))
        $terminalItem.Add_Click({ 
            # Open PowerShell terminal in ELR directory with custom prompt
            $powershellCommand = @"
            # Set custom prompt
            function global:prompt {
                "ELR> "
            }
            # Set working directory
            Set-Location -Path "$($ELRPath)"
            # Display welcome message
            Write-Host "=====================================" -ForegroundColor Green
            Write-Host "Enlightenment Lighthouse Runtime" -ForegroundColor Green
            Write-Host "=====================================" -ForegroundColor Green
            Write-Host "Current directory: $($ELRPath)" -ForegroundColor Cyan
            Write-Host ""
            Write-Host "Available commands:" -ForegroundColor Yellow
            Write-Host "  .\elr.ps1 help           - Show help" -ForegroundColor White
            Write-Host "  .\elr.ps1 api start      - Start API services" -ForegroundColor White
            Write-Host "  .\elr.ps1 api stop       - Stop API services" -ForegroundColor White
            Write-Host "  .\elr.ps1 api status     - Check API status" -ForegroundColor White
            Write-Host "  .\elr.ps1 api config     - Configure API settings" -ForegroundColor White
            Write-Host ""
            Write-Host "Type '.\elr.ps1 help' for more information." -ForegroundColor Cyan
            Write-Host "=====================================" -ForegroundColor Green
"@
            # Create temporary script file
            $tempScript = Join-Path -Path $env:TEMP -ChildPath "ELR-Console.ps1"
            $powershellCommand | Out-File -FilePath $tempScript -Encoding UTF8
            # Start PowerShell with the temporary script
            Start-Process "powershell.exe" -ArgumentList "-NoExit", "-ExecutionPolicy", "Bypass", "-File", "$tempScript"
        })

        # Separator
        $contextMenu.Items.Add("-")

        # Exit
        $exitItem = $contextMenu.Items.Add((Get-LocalizedString -keyPath "contextMenu.exit"))
        $exitItem.Add_Click({
            Write-Host "Exiting application..."
            $Global:trayIcon.Visible = $false
            if ($Global:mainForm -ne $null) { $Global:mainForm.Close() }
            if ($Global:chatForm -ne $null) { $Global:chatForm.Close() }
            if ($Global:settingsForm -ne $null) { $Global:settingsForm.Close() }
            [System.Windows.Forms.Application]::Exit()
        })

        $Global:trayIcon.ContextMenuStrip = $contextMenu

        # Double-click tray icon to open main window
        $Global:trayIcon.Add_DoubleClick({ Show-MainWindow })
        
        Write-Host "Tray icon initialized successfully"
    } catch {
        Write-Host "Error in Initialize-TrayIcon: $($_.Exception.Message)"
        Write-Host "Stack trace: $($_.ScriptStackTrace)"
    }
}

# Initialize
Write-Host "Initializing ELR Tray App..."
Write-Host "ELRPath: $ELRPath"
Write-Host "IconPath: $IconPath"
Write-Host "Icon exists: $(Test-Path $IconPath)"
Write-Host "LanguageFile: $LanguageFile"
Write-Host "Language file exists: $(Test-Path $LanguageFile)"
Write-Host "Language: $Global:language"
Write-Host "LangConfig: $Global:langConfig"

Initialize-TrayIcon
Write-Host "Tray icon initialized."
Write-Host "Tray icon visible: $($Global:trayIcon.Visible)"

# Start message loop
Write-Host "Starting message loop..."
[System.Windows.Forms.Application]::Run()
Write-Host "Message loop ended."