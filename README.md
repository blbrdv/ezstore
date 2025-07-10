# ![icon](/icons/icon16.png) ezstore
![GitHub License](https://img.shields.io/github/license/blbrdv/ezstore)
![GitHub Release](https://img.shields.io/github/v/release/blbrdv/ezstore)
[![E2e tests](https://github.com/blbrdv/ezstore/actions/workflows/e2e-tests.yaml/badge.svg)](https://github.com/blbrdv/ezstore/actions/workflows/e2e-tests.yaml)
![GitHub Downloads (all assets, all releases)](https://img.shields.io/github/downloads/blbrdv/ezstore/total)

Easy install apps from MS Store.

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

#### Requirements

1. Windows OS 8 or higher
2. Powershell 5.1 or higher

#### Execution

Run `.\tests\<scenario>.ps1`.

## License

This project licensed under [MIT](https://opensource.org/license/mit/) license.

See [LICENSE](LICENSE) file for more info.
