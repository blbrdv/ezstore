#Requires -Version 5.0

[System.Diagnostics.CodeAnalysis.SuppressMessageAttribute(
    "PSReviewUnusedParameter", "",
    Justification="Param SkipInstallTests used in Install tests explicitly."
)]
param (
    [Parameter(Mandatory=$true)]
    [string] $Path,
    [Parameter(Mandatory=$true)]
    [string] $Archs,
    [string[]] $Tags,
    [string[]] $ExcludeTags
)

. "$PSScriptRoot\Utils.ps1"

Set-StrictMode -Version 3.0;
$ProgressPreference = "SilentlyContinue";
$ErrorActionPreference = "Stop";
trap { Write-Output $_; exit 1 };

$Targets = $Archs -split "," | ForEach-Object {
    @{
        Arch = $_
        Path = [IO.Path]::Combine($Path, $_, "bin")
    }
};
$TestsList = Get-ChildItem -Path $PSScriptRoot -Filter "*.Tests.ps1" -ErrorAction 'SilentlyContinue';
$TargetsList = $Targets | ForEach-Object { $_ | ForEach-Object { "{ Arch: '$($_.Arch)'; Path: '$($_.Path)' }" } };

Write-Output "Targets:`t[ $($TargetsList -join ", ") ]";
Write-Output "Test files:`t[ $($TestsList -join ", ") ]";
if ( ($Tags | Measure-Object).Count -gt 0 ) {
    Write-Output "Tags:`t[ $($Tags -join ", ") ]";
}
if ( ($ExcludeTags | Measure-Object).Count -gt 0 ) {
    Write-Output "ExcludeTags:`t[ $($ExcludeTags -join ", ") ]";
}

if ( $null -eq $TestsList ) {
    Write-Error "No test files found.";
    exit 1;
}

Install-ModuleSafe -Name "Pester" -Version "5.7.1";
Import-ModuleSafe -Name "Pester" -Version "5.7.1";

. "$PSScriptRoot\Asserts.ps1";

$Config = [PesterConfiguration]::Default;
$Config.Should.ErrorAction = "Continue";
$Config.Output.Verbosity = "Detailed";
$Config.Run.Container = $TestsList | ForEach-Object {
        New-PesterContainer -Path $_;
    };
$Config.Filter.Tag = $Tags;
$Config.Filter.ExcludeTag = $ExcludeTags;

Invoke-Pester -Configuration $config;
