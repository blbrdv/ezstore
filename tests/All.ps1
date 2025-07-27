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

$SkipInstallTests = $false;

# skipping 386 and ARM architectures on Windows 11 ARM64 on Github VMs due to this issue
# https://learn.microsoft.com/en-us/windows/release-health/status-windows-11-21h2#2819msgdesc
if ( $null -ne $Env:GITHUB_ACTION ) {
    $OSArch = [System.Runtime.InteropServices.RuntimeInformation,mscorlib]::OSArchitecture.ToString().ToLower();
    $OSBuild = [Environment]::OSVersion.Version.BuildNum;

    $Is32bitApp = ( $Arch -eq "386" ) -or ( $Arch -eq "arm" );
    $IsArm64Win11 = ( $OSBuild -ge 22000 ) -and ( $OSArch -eq "arm64" )

    if ( $IsArm64Win11 -and $Is32bitApp ) {
        $SkipInstallTests = $True;
    }
}

Import-ModuleSafe -Name "Pester" -Version "5.7.1";
Import-ModuleSafe -Name "Pester" -Version "5.7.1";

$Config = [PesterConfiguration]::Default;
$Config.Should.ErrorAction = "Continue";
$Config.Output.Verbosity = "Detailed";
$Config.Run.Container = Get-ChildItem -Path $PSScriptRoot -Filter "*.Tests.ps1" -ErrorAction 'SilentlyContinue'
    | ForEach-Object {
        New-PesterContainer -Path $_;
    };

Invoke-Pester -Configuration $config;
