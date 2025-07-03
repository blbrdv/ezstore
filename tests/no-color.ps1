Set-StrictMode -Version 3.0;
$ErrorActionPreference = "Stop";
trap { Write-Error $_ -ErrorAction Continue; exit 1 };

. ".\tests\_core.ps1";

$Env:NO_COLOR = "1"
$Output = Install "9mvsm3j7zj7c" "1.1.0.0";

if ( $Output -match $ColorRegexp ) {
    throw 'Output has colors.';
}

Write-Output "Test passed!";
