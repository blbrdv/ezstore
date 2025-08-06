Set-StrictMode -Version 3.0;
$ErrorActionPreference = "Stop";
trap { Write-Error $_ -ErrorAction Continue; exit 1 };

$CurrentLocation = Get-Location;

Set-Location "$PSScriptRoot\.build";
& go mod download -x
Set-Location "$PSScriptRoot\.build\golangci-lint";
& go mod download -x
Set-Location "$PSScriptRoot";
& go mod download -x

Set-Location "$PSScriptRoot\.build";
& go run -trimpath=1 . $args;
Set-Location "$CurrentLocation";
exit $global:LASTEXITCODE;
