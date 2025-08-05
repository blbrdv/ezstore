# Changelog

## [1.4.1] - 2025-08-25

### Changed

- `github.com/imroc/req/v3` bump from `v3.53.0` to `v3.54.0`
- `golang.org/x/net` bump from `v0.41.0` to `v0.42.0`

### Fixed

- Fixed panic and error on old Windows versions [#42](https://github.com/blbrdv/ezstore/issues/42)

## [1.4.0] - 2025-07-02

### Added

- Add check for min/max version support for framework dependencies [#32](https://github.com/blbrdv/ezstore/issues/32)

### Fixed

- Fixed fetching and installing bundle dependency with different architecture that required by app
- Fixed incorrect bundle installation order

### Changed

- `github.com/imroc/req` bump from `v3.52.1` to `v3.53.0`

## [1.3.3] - 2025-06-13

### Fixed

- Files location in portable version [#25](https://github.com/blbrdv/ezstore/issues/25).

### Changed

- `golang.org/x/net` bump from `v0.40.0` to `v0.41.0`

## [1.3.2] - 2025-05-21

### Fixed

- Update script

## [1.3.1] - 2025-05-20

### Fixed

- Go dependency vulnerabilities [#18](https://github.com/blbrdv/ezstore/pull/18)
- Release CI/CD

## [1.3.0] - 2025-05-20

### Added

- Proper logging based on log level
- Script for automatic update check and download
- Version flag for printing app version

### Fixed

- Go dependency vulnerabilities [#14](https://github.com/blbrdv/ezstore/pull/14)
- Bug ["this package is not compatible with the device"](https://github.com/blbrdv/ezstore/issues/12)
- Downloading app with provided version, locale and architecture
- Inconsistent file fetching from MS Store API
- Unclear error messages
- CLI help output

## [1.2.0] - 2024-02-16

### Fixed

- Go dependency vulnerabilities
- CI/CD
- Lower version app dependency installation

## [1.1.1] - 2023-10-26

### Fixed

- Dependency installing order

## [1.1.0] - 2023-10-24

### Added

- Added dependency downloading

## [1.0.3] - 2023-10-11

### Added

- Installer

## [1.0.2] - 2023-09-27

### Fixed

- Runtime errors

## [1.0.1] - 2023-09-26

### Fixed

- Runtime errors

## [1.0.0]

Initial release.

[1.4.1]: https://github.com/blbrdv/ezstore/releases/tag/v1.4.1
[1.4.0]: https://github.com/blbrdv/ezstore/releases/tag/v1.4.0
[1.3.3]: https://github.com/blbrdv/ezstore/releases/tag/v1.3.3
[1.3.2]: https://github.com/blbrdv/ezstore/releases/tag/v1.3.2
[1.3.1]: https://github.com/blbrdv/ezstore/releases/tag/v1.3.1
[1.3.0]: https://github.com/blbrdv/ezstore/releases/tag/v1.3.0
[1.2.0]: https://github.com/blbrdv/ezstore/releases/tag/v1.2.0
[1.1.1]: https://github.com/blbrdv/ezstore/releases/tag/v1.1.1
[1.1.0]: https://github.com/blbrdv/ezstore/releases/tag/v1.1.0
[1.0.3]: https://github.com/blbrdv/ezstore/releases/tag/v1.0.3
[1.0.2]: https://github.com/blbrdv/ezstore/releases/tag/v1.0.2
[1.0.1]: https://github.com/blbrdv/ezstore/releases/tag/v1.0.1
[1.0.0]: https://github.com/blbrdv/ezstore/releases/tag/v1.0.0
