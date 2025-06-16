#!/bin/bash
set -euo pipefail

cd magefiles
go mod download -x
cd ..

go mod download -x

go tool -modfile='magefiles/go.mod' mage "$@"
