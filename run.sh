#!/bin/bash
set -euo pipefail

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null && pwd )"

cd "$DIR/.build"
go mod download -x
cd "$DIR/.build/golangci-lint"
go mod download -x
cd "$DIR"
go mod download -x

cd "$DIR/.build"
go run -trimpath=1 . "$@"
