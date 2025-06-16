Set-StrictMode -Version 3.0;
$ErrorActionPreference = "Stop";
trap { Write-Error $_ -ErrorAction Continue; exit 1 };

. ".\tests\_core.ps1";

$Env:NO_COLOR = "1"
$Output = Install $Data.Id $Data.Version;

if ( $Output -match $ColorRegexp ) {
    throw 'Output has colors.';
}

Write-Output "Test passed!";

