$ScriptPath = Split-Path -parent $MyInvocation.MyCommand.Definition;
$ExePath = Join-Path $ScriptPath "ezstore.exe";
$TempPath = [Environment]::GetFolderPath('LocalApplicationData');
$InstallerPath = Join-Path $TempPath "ezstore" "ezsetup.exe";

$global:Wait = $true;

if ($myInvocation.line -notmatch "ExecutionPolicy") {
    $global:Wait = $false;
}

function Wait {
    if ( $global:Wait ) {
        Pause;
    }
}

if ( -not (Test-Path $ExePath) ) {
    Write-Host "Error: $ExePath not found";
    Wait;
    exit 1;
}

$CurrentVersion = [Version](Get-Item $Path).VersionInfo.ProductVersion;

Write-Host "Current version: $CurrentVersion";

$Tag = (
    (Invoke-WebRequest -Uri "https://api.github.com/repos/blbrdv/ezstore/releases/latest").Content
    | ConvertFrom-Json
).tag_name;

$LastVersion = [Version]$Tag.Substring(1);

Write-Host "Last version: $LastVersion";

if ($LastVersion > $CurrentVersion) {
    Write-Host "Update needed!";
    Invoke-WebRequest -Uri "https://github.com/blbrdv/ezstore/releases/download/$Tag/ezsetup.exe" -OutFile $InstallerPath;
    Start-Process $InstallerPath -Wait;
} else {
    Write-Host "No update needed.";
}

Wait;
exit 0;
