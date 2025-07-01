Set-StrictMode -Version 3.0;
$ErrorActionPreference = "Stop";
trap { Write-Error $_ -ErrorAction Continue; exit 1 };

Set-Location ".mage"
& go mod download -x
Set-Location "golangci-lint"
& go mod download -x
Set-Location ..
Set-Location ..

& go mod download -x

& go tool -modfile='.mage\go.mod' mage -d .mage -w . $args
exit $global:LASTEXITCODE
