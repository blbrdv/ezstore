Set-StrictMode -Version 3.0;
$ErrorActionPreference = "Stop";
trap { Write-Error $_ -ErrorAction Continue; exit 1 };

if ( Test-Path env:GITHUB_ACTIONS ) {
    Install-Module -Name PSScriptAnalyzer -Force;
} else {
    Install-Module -Name PSScriptAnalyzer;
}

$Files = Get-ChildItem -Path .\*.ps1 -Recurse;
$Problems = [string[]]@();

foreach ($File in $Files) {
    $Output = Invoke-ScriptAnalyzer -Path $File;
    Write-Output "Analyzing '$File'";
    foreach ($Data in $Output) {
        $Problems += "$($Data.ScriptName):$($Data.Line)`t[$($Data.RuleName)]`t$($Data.Message)"
    }
}

if ( $Problems.Count -ne 0 ) {
    Write-Output "Problems found:";
    foreach ( $Problem in $Problems ) {
        Write-Output $Problem;
    }
    exit 1;
}

Write-Output "No problems found."
