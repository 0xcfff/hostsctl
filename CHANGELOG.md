# Changelog

All notable changes to this project are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.0.2] - 2026-08-18

### Added

- Homebrew installation instructions to the README (`brew tap 0xcfff/tap` / `brew install hostsctl`).

### Changed

- Modernized the Go toolchain: bumped the `go` directive from 1.20 to 1.26.
- Replaced `golang.org/x/exp` with the standard library (`slices`, `maps`).
- Removed the unused `sirupsen/logrus` dependency.
- Updated dependencies: `spf13/cobra` 1.5.0 -> 1.10.2, `spf13/afero` 1.9.2 -> 1.15.0, `stretchr/testify` 1.8.0 -> 1.12.0.
- Modernized CI: Go 1.21 -> 1.26, `actions/setup-go` v4 -> v5, `actions/upload-artifact` and `actions/download-artifact` v3 -> v4.

### Fixed

- Preserve the trailing newline when reading `/etc/hosts` on Go 1.22+ (changed `bufio.Scanner` behavior).
- Repaired a malformed struct tag in `commands/block/model.go` that newer `go vet` rejects.

### Security

- Bumped `golang.org/x/text` 0.3.7 -> 0.41.0, clearing known CVEs.

## [0.0.1] - 2023-11-01

### Added

- Initial release: manage individual and block DNS alias records in `/etc/hosts`, format the file, and back up/restore the configuration.

[0.0.2]: https://github.com/0xcfff/hostsctl/compare/v0.0.1...v0.0.2
[0.0.1]: https://github.com/0xcfff/hostsctl/releases/tag/v0.0.1
