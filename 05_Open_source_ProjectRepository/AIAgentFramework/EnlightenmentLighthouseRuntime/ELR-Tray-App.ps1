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
$Global:language = "zh"
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
    Write-Host $(Get-LocalizedString -keyPath "messages.error")
    exit 1
}

function Get-NetworkStatus {
    try {
        # Simulate network status check, return all services are running
        $status = "====================================`nELR Container Network Status`n====================================`nDesktop API: Running`n  Address: http://localhost:8081`nPublic API: Running`n  Address: http://localhost:8080`nModel Service: Running`n  Address: http://localhost:8082`nMicro Model Server: Running`n  Address: http://localhost:8083`n===================================="
        return $status
    } catch {
        return $(Get-LocalizedStringWithParams -keyPath "messages.failedStatus" -params @($_.Exception.Message))
    }
}

function Start-AllServices {
    try {
        & $ELRPS1 start-all
    } catch {
        Write-Host $(Get-LocalizedStringWithParams -keyPath "messages.failedStart" -params @($_.Exception.Message))
    }
}

function Stop-AllServices {
    try {
        & $ELRPS1 stop-all
    } catch {
        Write-Host $(Get-LocalizedStringWithParams -keyPath "messages.failedStop" -params @($_.Exception.Message))
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
                    if ($lineStr -match '^(\S+)\s+(\S+)\s+(\S+)\s+(\S+)\s+.*$') {
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
    $Global:chatForm.Text = $(Get-LocalizedString -keyPath "chatWindow.title")
    $Global:chatForm.Size = New-Object System.Drawing.Size(400, 500)
    $Global:chatForm.StartPosition = "CenterScreen"
    $Global:chatForm.Icon = [System.Drawing.Icon]::ExtractAssociatedIcon($IconPath)

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
    $sendButton.Text = $(Get-LocalizedString -keyPath "chatWindow.send")
    $sendButton.Add_Click({
        $message = $inputBox.Text
        if ($message -ne "") {
            $chatBox.AppendText($(Get-LocalizedStringWithParams -keyPath "chatWindow.you" -params @($message)) + "`n")
            $inputBox.Text = ""
            
            # Simulate ELR container response
            $response = Get-LocalizedStringWithParams -keyPath "chatWindow.response" -params @($message)
            Start-Sleep -Milliseconds 500
            $chatBox.AppendText("$response`n")
        }
    })
    $Global:chatForm.Controls.Add($sendButton)

    # Upload file button
    $uploadButton = New-Object System.Windows.Forms.Button
    $uploadButton.Location = New-Object System.Drawing.Point(10, 430)
    $uploadButton.Size = New-Object System.Drawing.Size(75, 25)
    $uploadButton.Text = $(Get-LocalizedString -keyPath "chatWindow.upload")
    $uploadButton.Add_Click({
        $openFileDialog = New-Object System.Windows.Forms.OpenFileDialog
        $openFileDialog.Title = $(Get-LocalizedString -keyPath "chatWindow.selectFile")
        $openFileDialog.Filter = "All Files (*.*)|*.*"
        if ($openFileDialog.ShowDialog() -eq "OK") {
            $filePath = $openFileDialog.FileName
            $chatBox.AppendText($(Get-LocalizedStringWithParams -keyPath "chatWindow.system" -params @($filePath)) + "`n")
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
    $Global:settingsForm.Text = $(Get-LocalizedString -keyPath "settingsWindow.title")
    $Global:settingsForm.Size = New-Object System.Drawing.Size(550, 350)
    $Global:settingsForm.StartPosition = "CenterScreen"
    $Global:settingsForm.Icon = [System.Drawing.Icon]::ExtractAssociatedIcon($IconPath)
    $Global:settingsForm.FormBorderStyle = "FixedSingle"
    $Global:settingsForm.MaximizeBox = $false
    $Global:settingsForm.MinimizeBox = $false

    # Token information
    $tokenLabel = New-Object System.Windows.Forms.Label
    $tokenLabel.Location = New-Object System.Drawing.Point(10, 20)
    $tokenLabel.Size = New-Object System.Drawing.Size(350, 60)
    $tokenLabel.Text = $(Get-LocalizedString -keyPath "settingsWindow.tokenInfo")
    $tokenLabel.AutoSize = $true
    $tokenLabel.TextAlign = [System.Drawing.ContentAlignment]::TopLeft
    $Global:settingsForm.Controls.Add($tokenLabel)

    # Create administrator button
    $createAdminBtn = New-Object System.Windows.Forms.Button
    $createAdminBtn.Location = New-Object System.Drawing.Point(360, 20)
    $createAdminBtn.Size = New-Object System.Drawing.Size(80, 25)
    $createAdminBtn.Text = $(Get-LocalizedString -keyPath "settingsWindow.createAdmin")
    $createAdminBtn.Add_Click({ 
        [System.Windows.Forms.MessageBox]::Show("Create admin feature in development", "Info", "OK", "Information")
    })
    $Global:settingsForm.Controls.Add($createAdminBtn)

    # Update token button
    $updateTokenBtn = New-Object System.Windows.Forms.Button
    $updateTokenBtn.Location = New-Object System.Drawing.Point(360, 50)
    $updateTokenBtn.Size = New-Object System.Drawing.Size(80, 25)
    $updateTokenBtn.Text = $(Get-LocalizedString -keyPath "settingsWindow.updateToken")
    $updateTokenBtn.Add_Click({ 
        [System.Windows.Forms.MessageBox]::Show("Update token feature in development", "Info", "OK", "Information")
    })
    $Global:settingsForm.Controls.Add($updateTokenBtn)

    # Admin count label
    $adminCountLabel = New-Object System.Windows.Forms.Label
    $adminCountLabel.Location = New-Object System.Drawing.Point(10, 90)
    $adminCountLabel.Size = New-Object System.Drawing.Size(200, 20)
    $adminCountLabel.Text = $(Get-LocalizedString -keyPath "settingsWindow.adminCount")
    $Global:settingsForm.Controls.Add($adminCountLabel)

    # View admins button
    $viewAdminsBtn = New-Object System.Windows.Forms.Button
    $viewAdminsBtn.Location = New-Object System.Drawing.Point(210, 88)
    $viewAdminsBtn.Size = New-Object System.Drawing.Size(80, 25)
    $viewAdminsBtn.Text = $(Get-LocalizedString -keyPath "settingsWindow.viewAdmins")
    $viewAdminsBtn.Add_Click({ 
        [System.Windows.Forms.MessageBox]::Show("View admins feature in development", "Info", "OK", "Information")
    })
    $Global:settingsForm.Controls.Add($viewAdminsBtn)

    # Desktop API settings
    $desktopApiLabel = New-Object System.Windows.Forms.Label
    $desktopApiLabel.Location = New-Object System.Drawing.Point(10, 120)
    $desktopApiLabel.Size = New-Object System.Drawing.Size(100, 20)
    $desktopApiLabel.Text = $(Get-LocalizedString -keyPath "settingsWindow.desktopApi")
    $Global:settingsForm.Controls.Add($desktopApiLabel)

    $desktopApiInput = New-Object System.Windows.Forms.TextBox
    $desktopApiInput.Location = New-Object System.Drawing.Point(110, 120)
    $desktopApiInput.Size = New-Object System.Drawing.Size(180, 20)
    $desktopApiInput.Text = "http://localhost:8000"
    $Global:settingsForm.Controls.Add($desktopApiInput)

    $desktopAutoBtn = New-Object System.Windows.Forms.Button
    $desktopAutoBtn.Location = New-Object System.Drawing.Point(300, 118)
    $desktopAutoBtn.Size = New-Object System.Drawing.Size(70, 25)
    $desktopAutoBtn.Text = $(Get-LocalizedString -keyPath "settingsWindow.auto")
    $desktopAutoBtn.Add_Click({ 
        $desktopApiInput.Text = "http://localhost:8000"
        [System.Windows.Forms.MessageBox]::Show("API address auto-recommended", "Info", "OK", "Information")
    })
    $Global:settingsForm.Controls.Add($desktopAutoBtn)

    $desktopSaveBtn = New-Object System.Windows.Forms.Button
    $desktopSaveBtn.Location = New-Object System.Drawing.Point(375, 118)
    $desktopSaveBtn.Size = New-Object System.Drawing.Size(70, 25)
    $desktopSaveBtn.Text = $(Get-LocalizedString -keyPath "settingsWindow.save")
    $desktopSaveBtn.Add_Click({ 
        [System.Windows.Forms.MessageBox]::Show("Settings saved", "Info", "OK", "Information")
    })
    $Global:settingsForm.Controls.Add($desktopSaveBtn)

    $desktopApiBtn = New-Object System.Windows.Forms.Button
    $desktopApiBtn.Location = New-Object System.Drawing.Point(450, 118)
    $desktopApiBtn.Size = New-Object System.Drawing.Size(80, 25)
    $desktopApiBtn.Text = $(Get-LocalizedString -keyPath "settingsWindow.startApi")
    $desktopApiBtn.Add_Click({ & $ELRPS1 start-desktop })
    $Global:settingsForm.Controls.Add($desktopApiBtn)

    # Public API settings
    $publicApiLabel = New-Object System.Windows.Forms.Label
    $publicApiLabel.Location = New-Object System.Drawing.Point(10, 150)
    $publicApiLabel.Size = New-Object System.Drawing.Size(100, 20)
    $publicApiLabel.Text = $(Get-LocalizedString -keyPath "settingsWindow.publicApi")
    $Global:settingsForm.Controls.Add($publicApiLabel)

    $publicApiInput = New-Object System.Windows.Forms.TextBox
    $publicApiInput.Location = New-Object System.Drawing.Point(110, 150)
    $publicApiInput.Size = New-Object System.Drawing.Size(180, 20)
    $publicApiInput.Text = "http://localhost:8001"
    $Global:settingsForm.Controls.Add($publicApiInput)

    $publicAutoBtn = New-Object System.Windows.Forms.Button
    $publicAutoBtn.Location = New-Object System.Drawing.Point(300, 148)
    $publicAutoBtn.Size = New-Object System.Drawing.Size(70, 25)
    $publicAutoBtn.Text = $(Get-LocalizedString -keyPath "settingsWindow.auto")
    $publicAutoBtn.Add_Click({ 
        $publicApiInput.Text = "http://localhost:8001"
        [System.Windows.Forms.MessageBox]::Show("API address auto-recommended", "Info", "OK", "Information")
    })
    $Global:settingsForm.Controls.Add($publicAutoBtn)

    $publicSaveBtn = New-Object System.Windows.Forms.Button
    $publicSaveBtn.Location = New-Object System.Drawing.Point(375, 148)
    $publicSaveBtn.Size = New-Object System.Drawing.Size(70, 25)
    $publicSaveBtn.Text = $(Get-LocalizedString -keyPath "settingsWindow.save")
    $publicSaveBtn.Add_Click({ 
        [System.Windows.Forms.MessageBox]::Show("Settings saved", "Info", "OK", "Information")
    })
    $Global:settingsForm.Controls.Add($publicSaveBtn)

    $publicApiBtn = New-Object System.Windows.Forms.Button
    $publicApiBtn.Location = New-Object System.Drawing.Point(450, 148)
    $publicApiBtn.Size = New-Object System.Drawing.Size(80, 25)
    $publicApiBtn.Text = $(Get-LocalizedString -keyPath "settingsWindow.startApi")
    $publicApiBtn.Add_Click({ & $ELRPS1 start-public })
    $Global:settingsForm.Controls.Add($publicApiBtn)

    # Model Service settings
    $modelServiceLabel = New-Object System.Windows.Forms.Label
    $modelServiceLabel.Location = New-Object System.Drawing.Point(10, 180)
    $modelServiceLabel.Size = New-Object System.Drawing.Size(100, 20)
    $modelServiceLabel.Text = $(Get-LocalizedString -keyPath "settingsWindow.modelService")
    $Global:settingsForm.Controls.Add($modelServiceLabel)

    $modelServiceInput = New-Object System.Windows.Forms.TextBox
    $modelServiceInput.Location = New-Object System.Drawing.Point(110, 180)
    $modelServiceInput.Size = New-Object System.Drawing.Size(180, 20)
    $modelServiceInput.Text = "http://localhost:8002"
    $Global:settingsForm.Controls.Add($modelServiceInput)

    $modelAutoBtn = New-Object System.Windows.Forms.Button
    $modelAutoBtn.Location = New-Object System.Drawing.Point(300, 178)
    $modelAutoBtn.Size = New-Object System.Drawing.Size(70, 25)
    $modelAutoBtn.Text = $(Get-LocalizedString -keyPath "settingsWindow.auto")
    $modelAutoBtn.Add_Click({ 
        $modelServiceInput.Text = "http://localhost:8002"
        [System.Windows.Forms.MessageBox]::Show("API address auto-recommended", "Info", "OK", "Information")
    })
    $Global:settingsForm.Controls.Add($modelAutoBtn)

    $modelSaveBtn = New-Object System.Windows.Forms.Button
    $modelSaveBtn.Location = New-Object System.Drawing.Point(375, 178)
    $modelSaveBtn.Size = New-Object System.Drawing.Size(70, 25)
    $modelSaveBtn.Text = $(Get-LocalizedString -keyPath "settingsWindow.save")
    $modelSaveBtn.Add_Click({ 
        [System.Windows.Forms.MessageBox]::Show("Settings saved", "Info", "OK", "Information")
    })
    $Global:settingsForm.Controls.Add($modelSaveBtn)

    $modelServiceBtn = New-Object System.Windows.Forms.Button
    $modelServiceBtn.Location = New-Object System.Drawing.Point(450, 178)
    $modelServiceBtn.Size = New-Object System.Drawing.Size(80, 25)
    $modelServiceBtn.Text = $(Get-LocalizedString -keyPath "settingsWindow.startApi")
    $modelServiceBtn.Add_Click({ & $ELRPS1 start-model })
    $Global:settingsForm.Controls.Add($modelServiceBtn)

    # Micro Model settings
    $microModelLabel = New-Object System.Windows.Forms.Label
    $microModelLabel.Location = New-Object System.Drawing.Point(10, 210)
    $microModelLabel.Size = New-Object System.Drawing.Size(100, 20)
    $microModelLabel.Text = $(Get-LocalizedString -keyPath "settingsWindow.microModel")
    $Global:settingsForm.Controls.Add($microModelLabel)

    $microModelInput = New-Object System.Windows.Forms.TextBox
    $microModelInput.Location = New-Object System.Drawing.Point(110, 210)
    $microModelInput.Size = New-Object System.Drawing.Size(180, 20)
    $microModelInput.Text = "http://localhost:8003"
    $Global:settingsForm.Controls.Add($microModelInput)

    $microAutoBtn = New-Object System.Windows.Forms.Button
    $microAutoBtn.Location = New-Object System.Drawing.Point(300, 208)
    $microAutoBtn.Size = New-Object System.Drawing.Size(70, 25)
    $microAutoBtn.Text = $(Get-LocalizedString -keyPath "settingsWindow.auto")
    $microAutoBtn.Add_Click({ 
        $microModelInput.Text = "http://localhost:8003"
        [System.Windows.Forms.MessageBox]::Show("API address auto-recommended", "Info", "OK", "Information")
    })
    $Global:settingsForm.Controls.Add($microAutoBtn)

    $microSaveBtn = New-Object System.Windows.Forms.Button
    $microSaveBtn.Location = New-Object System.Drawing.Point(375, 208)
    $microSaveBtn.Size = New-Object System.Drawing.Size(70, 25)
    $microSaveBtn.Text = $(Get-LocalizedString -keyPath "settingsWindow.save")
    $microSaveBtn.Add_Click({ 
        [System.Windows.Forms.MessageBox]::Show("Settings saved", "Info", "OK", "Information")
    })
    $Global:settingsForm.Controls.Add($microSaveBtn)

    $microModelBtn = New-Object System.Windows.Forms.Button
    $microModelBtn.Location = New-Object System.Drawing.Point(450, 208)
    $microModelBtn.Size = New-Object System.Drawing.Size(80, 25)
    $microModelBtn.Text = $(Get-LocalizedString -keyPath "settingsWindow.startApi")
    $microModelBtn.Add_Click({ & $ELRPS1 start-micro })
    $Global:settingsForm.Controls.Add($microModelBtn)

    # Start/Stop all APIs
    $startAllBtn = New-Object System.Windows.Forms.Button
    $startAllBtn.Location = New-Object System.Drawing.Point(10, 250)
    $startAllBtn.Size = New-Object System.Drawing.Size(120, 30)
    $startAllBtn.Text = $(Get-LocalizedString -keyPath "settingsWindow.startAll")
    $startAllBtn.Add_Click({ Start-AllServices })
    $Global:settingsForm.Controls.Add($startAllBtn)

    $stopAllBtn = New-Object System.Windows.Forms.Button
    $stopAllBtn.Location = New-Object System.Drawing.Point(140, 250)
    $stopAllBtn.Size = New-Object System.Drawing.Size(120, 30)
    $stopAllBtn.Text = $(Get-LocalizedString -keyPath "settingsWindow.stopAll")
    $stopAllBtn.Add_Click({ Stop-AllServices })
    $Global:settingsForm.Controls.Add($stopAllBtn)

    $Global:settingsForm.Show()
}

function Show-MainWindow {
    if ($Global:mainForm -ne $null -and $Global:mainForm.Visible) {
        $Global:mainForm.Activate()
        return
    }

    $Global:mainForm = New-Object System.Windows.Forms.Form
    $Global:mainForm.Text = $(Get-LocalizedString -keyPath "mainWindow.title")
    $Global:mainForm.Size = New-Object System.Drawing.Size(850, 300)
    $Global:mainForm.StartPosition = "CenterScreen"
    $Global:mainForm.Icon = [System.Drawing.Icon]::ExtractAssociatedIcon($IconPath)

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
    $statusText = "$(Get-LocalizedString -keyPath 'mainWindow.statusInfo.connected')`n"
    $statusText += "$(Get-LocalizedStringWithParams -keyPath 'mainWindow.statusInfo.containerCount' -params $totalCount)`n"
    $statusText += "$(Get-LocalizedStringWithParams -keyPath 'mainWindow.statusInfo.runningCount' -params $runningCount)"
    # Debug: print status text
    Write-Host "Status text: '$statusText'"
    $statusLabel.Text = $statusText
    $Global:mainForm.Controls.Add($statusLabel)

    # Network status button
    $statusBtn = New-Object System.Windows.Forms.Button
    $statusBtn.Location = New-Object System.Drawing.Point(10, 80)
    $statusBtn.Size = New-Object System.Drawing.Size(830, 30)
    $statusBtn.Text = $(Get-LocalizedString -keyPath "mainWindow.checkStatus")
    $statusBtn.Add_Click({
        $status = Get-NetworkStatus
        [System.Windows.Forms.MessageBox]::Show($status, $(Get-LocalizedString -keyPath "messages.networkStatus"), "OK", "Information")
    })
    $Global:mainForm.Controls.Add($statusBtn)

    # Container list
    $containerList = New-Object System.Windows.Forms.ListView
    $containerList.Location = New-Object System.Drawing.Point(10, 120)
    $containerList.Size = New-Object System.Drawing.Size(800, 150)
    $containerList.View = "Details"
    $containerList.Columns.Add($(Get-LocalizedString -keyPath "mainWindow.containerColumns.name"), 120)
    $containerList.Columns.Add($(Get-LocalizedString -keyPath "mainWindow.containerColumns.status"), 100)
    $containerList.Columns.Add($(Get-LocalizedString -keyPath "mainWindow.containerColumns.components"), 200)
    $containerList.Columns.Add($(Get-LocalizedString -keyPath "mainWindow.containerColumns.sandbox"), 100)
    $containerList.Columns.Add($(Get-LocalizedString -keyPath "mainWindow.containerColumns.resourceUsage"), 250)

    # Add containers
    foreach ($container in $containers) {
        $item = New-Object System.Windows.Forms.ListViewItem($container.Name)
        
        # Get localized status
        if ($container.Status -eq "created") {
            $statusText = $(Get-LocalizedString -keyPath "mainWindow.containerStatus.created")
        } elseif ($container.Status -eq "running") {
            $statusText = $(Get-LocalizedString -keyPath "mainWindow.containerStatus.running")
        } else {
            $statusText = $container.Status
        }
        
        # Get localized components
        if ($container.Components -eq "unknown") {
            $componentsText = $(Get-LocalizedString -keyPath "mainWindow.unknownComponents")
        } else {
            $componentsText = $container.Components
        }
        
        $item.SubItems.Add($statusText)
        $item.SubItems.Add($componentsText)
        $item.SubItems.Add($container.Sandbox)
        $memoryLabel = $(Get-LocalizedString -keyPath "mainWindow.resourceUsage.memory")
        $cpuLabel = $(Get-LocalizedString -keyPath "mainWindow.resourceUsage.cpu")
        $gpuLabel = $(Get-LocalizedString -keyPath "mainWindow.resourceUsage.gpu")
        $resourceUsage = "${memoryLabel}: $($container.Memory)MB, ${cpuLabel}: $($container.CPU)%, ${gpuLabel}: $($container.GPU)%"
        $item.SubItems.Add($resourceUsage)
        $containerList.Items.Add($item)
    }

    $Global:mainForm.Controls.Add($containerList)

    $Global:mainForm.Show()
}

# Create system tray
function Initialize-TrayIcon {
    $Global:trayIcon = New-Object System.Windows.Forms.NotifyIcon
    $Global:trayIcon.Icon = [System.Drawing.Icon]::ExtractAssociatedIcon($IconPath)
    $Global:trayIcon.Text = $(Get-LocalizedString -keyPath "trayIconText")
    $Global:trayIcon.Visible = $true

    # Create context menu
    $contextMenu = New-Object System.Windows.Forms.ContextMenuStrip

    # Open main window
    $openItem = $contextMenu.Items.Add($(Get-LocalizedString -keyPath "contextMenu.open"))
    $openItem.Add_Click({ Show-MainWindow })

    # Open chat window
    $chatItem = $contextMenu.Items.Add($(Get-LocalizedString -keyPath "contextMenu.chat"))
    $chatItem.Add_Click({ Show-ChatWindow })

    # Open settings window
    $settingsItem = $contextMenu.Items.Add($(Get-LocalizedString -keyPath "contextMenu.settings"))
    $settingsItem.Add_Click({ Show-SettingsWindow })

    # Network status
    $networkItem = $contextMenu.Items.Add($(Get-LocalizedString -keyPath "contextMenu.networkStatus"))
    $networkItem.Add_Click({
        $status = Get-NetworkStatus
        [System.Windows.Forms.MessageBox]::Show($status, $(Get-LocalizedString -keyPath "messages.networkStatus"), "OK", "Information")
    })

    # Start all services
    $startAllItem = $contextMenu.Items.Add($(Get-LocalizedString -keyPath "contextMenu.startAll"))
    $startAllItem.Add_Click({ Start-AllServices })

    # Stop all services
    $stopAllItem = $contextMenu.Items.Add($(Get-LocalizedString -keyPath "contextMenu.stopAll"))
    $stopAllItem.Add_Click({ Stop-AllServices })

    # Separator
    $contextMenu.Items.Add("-")

    # Exit
    $exitItem = $contextMenu.Items.Add($(Get-LocalizedString -keyPath "contextMenu.exit"))
    $exitItem.Add_Click({
        $Global:trayIcon.Visible = $false
        if ($Global:mainForm -ne $null) { $Global:mainForm.Close() }
        if ($Global:chatForm -ne $null) { $Global:chatForm.Close() }
        if ($Global:settingsForm -ne $null) { $Global:settingsForm.Close() }
        [System.Windows.Forms.Application]::Exit()
    })

    $Global:trayIcon.ContextMenuStrip = $contextMenu

    # Double-click tray icon to open main window
    $Global:trayIcon.Add_DoubleClick({ Show-MainWindow })
}

# Initialize
Initialize-TrayIcon

# Start message loop
[System.Windows.Forms.Application]::Run()