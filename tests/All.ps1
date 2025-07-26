#Requires -Version 5.0

param (
    [Parameter(Mandatory=$true)]
    [string] $Path,
    [Parameter(Mandatory=$true)]
    [string] $Archs,
    [AllowEmptyString()]
    [string] $Version
)

. "$PSScriptRoot\Utils.ps1"

Set-StrictMode -Version 3.0;
$ErrorActionPreference = "Stop";
trap { Write-Output $_; exit 1 };

$Targets = $Archs -split "," | ForEach-Object {
    @{
        Arch = $_
        Path = [IO.Path]::Combine($Path, $_, "bin")
    }
};

Import-ModuleSafe -Name "Pester" -Version "5.7.1";
Import-ModuleSafe-Name "Pester" -Version "5.7.1";

# For some reason default glob search in New-PesterContainer Path parameter didn't work for me so using Get-ChildItem
$Containers = Get-ChildItem -Path $PSScriptRoot -Filter "*.Tests.ps1" -ErrorAction 'SilentlyContinue'
    | ForEach-Object {
        New-PesterContainer -Path $_;
    };

Invoke-Pester -Container $Containers -Output Detailed;
