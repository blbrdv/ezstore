Set-StrictMode -Version 3.0;
$ErrorActionPreference = "Stop";
trap { Write-Error $_ -ErrorAction Continue; exit 1 };

$Id = "9mvsm3j7zj7c";
$Name = "PeterEtelej.TreeCLI";
$Version = "1.1.0.0";

.\output\bin\ezstore.exe install $Id --ver $Version --verbosity m;

Import-Module -Name Appx -UseWindowsPowerShell -WarningAction SilentlyContinue;

$Package = Get-AppxPackage -Name $Name;

if ( $null -eq $Package ) {
    throw "Package $Name not installed";
}

if ( $Package.Version -ne $Version ) {
    throw "Wrong version installed. Expected $Version, actual: $($Package.Version).";
}

Write-Host "Package '$Name' ($Id) $Version successfully installed!";
