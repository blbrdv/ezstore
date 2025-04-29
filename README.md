# ![icon](/icons/icon16.png) ezstore
![GitHub License](https://img.shields.io/github/license/blbrdv/ezstore)
![GitHub Release](https://img.shields.io/github/v/release/blbrdv/ezstore)
![GitHub Downloads (all assets, all releases)](https://img.shields.io/github/downloads/blbrdv/ezstore/total)

Easy install apps from MS Store on Long-term Servicing Windows.

## Installation

Download installer or portable version from
[release page](https://github.com/blbrdv/ezstore/releases).

## Usage

```
Easy install apps from MS Store

Usage:
    ezstore install <id> [--version=<xyz>] [--locale=<cc-CC>] [--debug]
    ezstore --help

Options:
    -h --help      Show this screen.
    -v --version   Sets the version of the product [default: latest].
    -l --locale    Sets the locale name of the product [default: current value in the OS or en_US].
    -d --debug     Sets the debug mode [default: false].

Examples:
    ezstore install 9nh2gph4jzs4
    ezstore install 9nh2gph4jzs4 -v 1.0.3.0 --locale cs_CZ --debug
```

## Development

### Requirements

[Golang](https://go.dev/dl/) version `1.24` or later must be installed

#### Lint

[staticcheck](https://staticcheck.dev/) must be installed

#### Build 

1. [go-winres](https://github.com/tc-hib/go-winres) must be installed
2. [7-Zip](https://7-zip.org/) must be installed and put in $PATH
3. [Inno Setup 6](https://jrsoftware.org/isinfo.php) must be installed and put in $PATH

### Automation

Use `.\run.ps1 <task>`.

Available tasks:
 - `clean` - removes build directories and files
 - `format` - fix codestyle
 - `lint` - run analysis for code
 - `test` - run unit tests
 - `build` - build cli and compile installer

## License

This project licensed under [MIT](https://opensource.org/license/mit/) license.

See [LICENSE](LICENSE) file for more info.
