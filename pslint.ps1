Set-StrictMode -Version 3.0;
$ErrorActionPreference = "Stop";
trap { Write-Error $_ -ErrorAction Continue; exit 1 };

if ( Test-Path env:GITHUB_ACTIONS ) {
    Install-Module -Name PSScriptAnalyzer -Force;
} else {
    Install-Module -Name PSScriptAnalyzer;
}

$Settings = @{
    ExcludeRules = @(
        'PSAvoidGlobalVars',
        'PSUseDeclaredVarsMoreThanAssignments',
        'PSAvoidUsingInvokeExpression'
    )
}

$Files = Get-ChildItem -Path .\*.ps1 -Recurse;
$Found = $false;

foreach ($File in $Files) {
    $Output = Invoke-ScriptAnalyzer -Path $File -Settings $Settings ;
    foreach ($Data in $Output) {
        $Found = $true;
        Write-Output "$($Data.ScriptName):$($Data.Line) $($Data.Message)";
    }
}

if ( $Found ) {
    exit 1;
}
