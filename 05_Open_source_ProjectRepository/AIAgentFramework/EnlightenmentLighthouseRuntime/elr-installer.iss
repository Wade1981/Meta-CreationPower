; ELR Container Installer Script for Inno Setup
; This script creates a binary installer for ELR

[Setup]
AppName=Enlightenment Lighthouse Runtime
AppVersion=1.0.0
AppPublisher=Enlightenment Lighthouse Origin Team
AppPublisherURL=https://github.com/enlightenment-lighthouse
AppSupportURL=https://github.com/enlightenment-lighthouse/support
AppUpdatesURL=https://github.com/enlightenment-lighthouse/updates
DefaultDirName={userprofile}\ELR
DefaultGroupName=Enlightenment Lighthouse Runtime
OutputDir=output
OutputBaseFilename=elr-installer
SetupIconFile=icons\elr_icon.ico
Compression=lzma
SolidCompression=yes
PrivilegesRequired=lowest

[Languages]
Name: "english"; MessagesFile: "compiler:Default.isl"
Name: "chinese"; MessagesFile: "compiler:Languages\ChineseSimplified.isl"

[Tasks]
Name: "desktopicon"; Description: "{cm:CreateDesktopIcon}"; GroupDescription: "{cm:AdditionalIcons}"; Flags: unchecked
Name: "quicklaunchicon"; Description: "{cm:CreateQuickLaunchIcon}"; GroupDescription: "{cm:AdditionalIcons}"; Flags: unchecked; OnlyBelowVersion: 0,6.1

[Files]
; Main ELR files
Source: "elr.ps1"; DestDir: "{app}\bin";
Source: "elr.bat"; DestDir: "{app}\bin";
Source: "elr_api_server.py"; DestDir: "{app}\bin";

; Micro model directory
Source: "micro_model\*"; DestDir: "{app}\lib\micro_model"; Flags: recursesubdirs

; Models directory
Source: "models\*"; DestDir: "{app}\models"; Flags: recursesubdirs

; API server dependencies
Source: "python-portable\*"; DestDir: "{app}\python-portable"; Flags: recursesubdirs

[Icons]
Name: "{group}\Enlightenment Lighthouse Runtime"; Filename: "{app}\bin\elr.cmd"
Name: "{group}\Uninstall ELR"; Filename: "{uninstallexe}"
Name: "{commondesktop}\Enlightenment Lighthouse Runtime"; Filename: "{app}\bin\elr.cmd"; Tasks: desktopicon

[Run]
Filename: "{app}\bin\elr.cmd"; Parameters: "version"; Description: "Check ELR version"; Flags: postinstall nowait skipifsilent

[UninstallRun]
Filename: "{app}\bin\elr.cmd"; Parameters: "stop"; Flags: nowait

[Registry]
; Add ELR bin directory to PATH
Root: HKCU; Subkey: "Environment"; ValueType: expandsz; ValueName: "PATH"; ValueData: "{olddata};{app}\bin"; Flags: preservestringtype

[Code]
function InitializeSetup(): Boolean;
begin
  Result := True;
  // Check if PowerShell is available
  if not RegKeyExists(HKLM, 'SOFTWARE\Microsoft\PowerShell\1') and 
     not RegKeyExists(HKLM, 'SOFTWARE\Microsoft\PowerShell\3') and
     not RegKeyExists(HKLM, 'SOFTWARE\Microsoft\PowerShell\5') then
  begin
    MsgBox('PowerShell is required for ELR. Please install PowerShell 5.1 or later.', mbError, MB_OK);
    Result := False;
  end;
end;
