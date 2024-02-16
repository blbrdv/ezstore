# ![icon](/dist/icon16.png) ezstore
Easy install apps from MS Store on Long-term Servicing Windows.

## Installation

Download installer or portable version from
[release page](https://github.com/blbrdv/ezstore/releases).

## Usage

`.\ezstore.exe [OPTIONS] install <id>`

where `id` is product identifier form MS store, e.g.
[9NH2GPH4JZS4](https://apps.microsoft.com/store/detail/tiktok/9NH2GPH4JZS4).

### Options
  - `--version <value>`, `-v <value>` for product version (default: "latest")
  - `--locale <value>`, `-l <value>` for product locale
  - `--debug`, `-d` for debug output (default: false)
  - `--help`, `-h` - show help

## Development

### Run tests

`task test`

### Build

Requirements:

 1. [go-task](https://github.com/go-task/task) installed;
 2. [staticcheck](https://staticcheck.dev/) installed;
 3. [Inno Setup 6](https://jrsoftware.org/isinfo.php) installed and put in $PATH.

`task build`

## License

This project licensed under [MIT](https://opensource.org/license/mit/) license.

See [LICENSE](LICENSE) file for more info.
