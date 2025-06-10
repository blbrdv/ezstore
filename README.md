# ![icon](/icons/icon16.png) ezstore
![GitHub License](https://img.shields.io/github/license/blbrdv/ezstore)
![GitHub Release](https://img.shields.io/github/v/release/blbrdv/ezstore)
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

1. Clone project
2. Install [Golang](https://go.dev/dl/) version `1.24` or later
3. Run `go get -t .\...` to install dependencies and tools for project

### Automation

Use `.\run.ps1 -l` to see available tasks.

Use `.\run.ps1 -h <task>` to see task description.

Use `.\run.ps1 [-v] <task>` to run task. Flag `-v` sets verbose output. 

## License

This project licensed under [MIT](https://opensource.org/license/mit/) license.

See [LICENSE](LICENSE) file for more info.
