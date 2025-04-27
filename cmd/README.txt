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
