param (
    [Parameter(Mandatory=$true,Position=0)]
    [string]$JSON
)

Set-StrictMode -Version 3.0;
$ErrorActionPreference = "Stop";
trap { Write-Error $_ -ErrorAction Continue; exit 1 };

. ".\tests\_core.ps1";

$Target = $JSON | ConvertFrom-Json;

Install $Target.id $Target.version > $null;

Import-Module -Name Appx -UseWindowsPowerShell -WarningAction SilentlyContinue;

$Package = Get-AppxPackage -Name $Target.name;

if ( $null -eq $Package ) {
    throw "Package $($Target.name) not installed";
}

if ( $Package.Version -ne $Target.version ) {
    throw "Wrong version installed. Expected $($Target.version), actual: $($Package.Version)."
}

Write-Output "Test passed!";
