Set-StrictMode -Version 3.0;
$ErrorActionPreference = "Stop";
trap { Write-Error $_ -ErrorAction Continue; exit 1 };

cd "magefiles"
& go mod download -x
cd ..

& go mod download -x

& go tool -modfile='magefiles\go.mod' mage $args
exit $global:LASTEXITCODE
