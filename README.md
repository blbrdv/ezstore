# ![icon](/icons/icon16.png) ezstore
![GitHub License](https://img.shields.io/github/license/blbrdv/ezstore)
![GitHub Release](https://img.shields.io/github/v/release/blbrdv/ezstore)
[![E2e tests](https://github.com/blbrdv/ezstore/actions/workflows/e2e-tests.yaml/badge.svg)](https://github.com/blbrdv/ezstore/actions/workflows/e2e-tests.yaml)
![GitHub Downloads (all assets, all releases)](https://img.shields.io/github/downloads/blbrdv/ezstore/total)

Easy install apps from MS Store.

## Requirements

1. Windows OS 8 or higher
2. Powershell 5.1 or higher

App tested on following Windows versions:

| Version | Buid  | OS arch | App arch   |
|---------|-------|---------|------------|
| 8       | 9200  | x86     | x86        |
| 8.1     | 9600  | x86     | x86        |
| 10 22H2 | 19042 | x64     | x64        |
| 10 21H2 | 20348 | x64     | x64, x86   |
| 10 24H2 | 26100 | x64     | x64, x86   |
| 11 21H2 | 22000 | x64     | x64        |
| 11 24H2 | 26100 | arm64   | arm64, x64 |

## Installation

Download installer or portable version from
[release page](https://github.com/blbrdv/ezstore/releases).

## Usage

```
Easy install apps from MS Store

Usage:
    ezstore install <id> [options]
    ezstore --help
    ezstore --version

Options:
    -h --help      Print this text.
    -v --version   Print app version.
    --ver          Sets the version of the product [default: latest].
    --locale       Sets the locale name of the product [default: current value in the OS or en_US].
    --verbosity    Sets verbosity level [default: n].
                   Available log levels:
                     * q - quiet, no output at all
                     * m - minimal, only SUCCESS and ERROR logs
                     * n - normal, same as minimal plus INFO and WARNING logs
                     * d - detailed, same as normal plus DEBUG logs and tracing net errors to log file

Examples:
    ezstore install 9nh2gph4jzs4
    ezstore install 9nh2gph4jzs4 -v 1.0.3.0 --locale cs_CZ --verbosity d
```

## Development

Project need [Golang](https://go.dev/dl/) version `1.24` or later.

### Automation

Use `.\run.ps1 [flags] [tasks]`.

Use `.\run.ps1 help` to see available flags and tasks.

### End-to-end tests

Run `.\tests\All.ps1 -Path <path> -Archs <archs> -Tags <tags> -ExcludeTags <exclude-tags>`.

 - `<path>` - path to directory with app binaries, e.g. `.\output`.
 - `<archs>` - list of app architectures to run separated with comma, e.g. `amd64,386`.
Allowed values: `amd64`, `386`, `arm64`, `arm`.
 - `<tags>` - list of test tags to include in test run. Optional parameter.
 - `<exclude-tags>` - list of test tags to exclude from test run. Optional parameter.

## License

This project licensed under [MIT](https://opensource.org/license/mit/) license.

See [LICENSE](LICENSE) file for more info.
