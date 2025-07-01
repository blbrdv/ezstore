#!/bin/bash
set -euo pipefail

cd .mage
go mod download -x
cd ..

go mod download -x

go tool -modfile='.mage/go.mod' mage -d .mage -w . "$@"
