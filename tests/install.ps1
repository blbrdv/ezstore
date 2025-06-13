Set-StrictMode -Version 3.0;
$ErrorActionPreference = "Stop";
trap { Write-Error $_ -ErrorAction Continue; exit 1 };

. ".\tests\_core.ps1";

Install $Data.Id $Data.Version;

Import-Module -Name Appx -UseWindowsPowerShell -WarningAction SilentlyContinue;

$Package = Get-AppxPackage -Name $Data.Name;

if ( $null -eq $Package ) {
    throw "Package $($Data.Name) not installed";
}

if ( $Package.Version -ne $Data.Version ) {
    throw "Wrong version installed. Expected $($Data.Version), actual: $($Package.Version)."
}

Write-Host "Test passed!";
