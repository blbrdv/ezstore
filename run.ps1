Set-StrictMode -Version 3.0;
$ErrorActionPreference = "Stop";
trap { Write-Error $_ -ErrorAction Continue; exit 1 };

Set-Location "magefiles"
& go mod download -x
Set-Location ..

& go mod download -x

& go tool -modfile='magefiles\go.mod' mage $args
exit $global:LASTEXITCODE
