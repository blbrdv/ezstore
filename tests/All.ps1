#Requires -Version 5.0

param (
    [Parameter(Mandatory=$true)]
    [string] $Path,
    [Parameter(Mandatory=$true)]
    [string] $Archs
)

. "$PSScriptRoot\Utils.ps1"

Set-StrictMode -Version 3.0;
$ErrorActionPreference = "Stop";
trap { Write-Output $_; exit 1 };

[System.Diagnostics.CodeAnalysis.SuppressMessageAttribute(
    "PSUseDeclaredVarsMoreThanAssignments", "",
    Justification="This variable used in tests explicitly."
)]
$Targets = $Archs -split "," | ForEach-Object {
    @{
        Arch = $_
        Path = [IO.Path]::Combine($Path, $_, "bin")
    }
};

Import-ModuleSafe -Name "Pester" -Version "5.7.1";
Import-ModuleSafe -Name "Pester" -Version "5.7.1";

$Config = [PesterConfiguration]::Default;
$Config.Should.ErrorAction = "Continue";
$Config.Output = "Detailed";
$Config.Run.Container = Get-ChildItem -Path $PSScriptRoot -Filter "*.Tests.ps1" -ErrorAction 'SilentlyContinue'
    | ForEach-Object {
        New-PesterContainer -Path $_;
    };

Invoke-Pester -Configuration $config;
