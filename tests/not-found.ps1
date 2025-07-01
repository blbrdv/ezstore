Set-StrictMode -Version 3.0;
$ErrorActionPreference = "Stop";
trap { Write-Error $_ -ErrorAction Continue; exit 1 };

. ".\tests\_core.ps1";

$Id = "f1o2o3b4a5r6";
$Expected = '^\[ERR\] Finished with error: can not fetch product info: product with id "' + $Id + '" and locale "en-US" not found$';
$Actual = Install $Id "1.0.0.0" 1;
$Actual = $Actual -replace $ColorRegexp;

if ( $Actual -match $Expected ) {
    Write-Output "Test passed!";
    exit 0;
}

Write-Host "Expected (regexp): '$Expected'";
throw 'Incorrect app output.';
