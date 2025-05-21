Set-StrictMode -Version 3.0;
$ErrorActionPreference = "Stop";
trap { Write-Error $_ -ErrorAction Continue; exit 1 }

$ScriptPath = Split-Path -parent $MyInvocation.MyCommand.Definition;
$ExePath = Join-Path (Join-Path $ScriptPath "bin") "ezstore.exe";
$TempPath = [Environment]::GetFolderPath('LocalApplicationData');
$InstallerPath = Join-Path (Join-Path $TempPath "ezstore") "ezsetup.exe";

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

$CurrentVersion = [Version](Get-Item $ExePath).VersionInfo.ProductVersion;

Write-Host "Current version: $CurrentVersion";

$Response = (Invoke-WebRequest -Uri "https://api.github.com/repos/blbrdv/ezstore/releases/latest").Content | ConvertFrom-Json;
$LastVersion = [Version]$Response.tag_name.Substring(1);

Write-Host "Last version: $LastVersion";

if ($LastVersion -gt $CurrentVersion) {
    Write-Host "Update needed!";
    Invoke-WebRequest -Uri "https://github.com/blbrdv/ezstore/releases/download/${Response.tag_name}/ezsetup.exe" -OutFile $InstallerPath;
    Start-Process $InstallerPath -Wait;
} else {
    Write-Host "No update needed.";
}

Wait;
exit 0;
