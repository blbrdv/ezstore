Set-StrictMode -Version 3.0;
$ErrorActionPreference = "Stop";
trap { Write-Error $_ -ErrorAction Continue; exit 1 };

. ".\tests\_core.ps1";

$Key = $($Targets.Keys)[2];

$Env:NO_COLOR = "1"
$Output = Install $Key $Targets[$Key].Version;

if ( $Output -match $ColorRegexp ) {
    throw 'Output has colors.';
}

Write-Output "Test passed!";
