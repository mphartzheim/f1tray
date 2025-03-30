# Changelog

<!-- Only the latest version entry will appear in the GitHub Release -->

## [v0.1.2-dev.2] - 2025-03-29

### Added
- GitHub Actions workflow to automatically generate GitHub Releases
- Uploads Linux `.AppImage` and Windows `.zip` from versioned build folders
- Supports wildcard filenames and manual triggers

### Tested
- First automated release run using tag `v0.1.2-dev.2`

---

## [v0.1.2-dev.1] - 2025-03-29

### Added
- New "Standings" tab with inner tabs for "Drivers" and "Constructors"
- Parsers for driver and constructor standings using the Jolpica API
- Asterisk notification logic for Standings tab, matching existing Results tab behavior
- Standings content now updates dynamically when the selected year changes
- Preserves selected inner tab (Drivers/Constructors) when refreshing Standings

### Fixed
- Typo in `ConstructorsStandingsURL` constant
- Results and Standings tabs now properly mark themselves with an asterisk on background refresh or year change

### Internal
- Refactored `parsers.go` to support modular parsing of various data types
- Organized models to separate JSON responses and reusable types
