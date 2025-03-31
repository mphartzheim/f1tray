# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),  
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

---

## [02.2.] - 2025-xx-xx

### Added
- Watch on F1TV button is now more prominent.
- Users can now highlight their favorite drivers!
- Global notification system (eg. warn the user if selecting too many favorites) implemented.

---

## [0.2.1] - 2025-03-30

### Added
- Driver names in standings and results are now clickable, opening bio pages with a fallback to the API URL.
- Added helper functions to reduce parser repetition.
- Added theme lookup helper using themes.AvailableThemes.

### Fixed
- Tray menu options now display and respond correctly across themes.
- Team themes now correctly apply highlight colors on initial load.

---

## [0.2.0] - 2025-03-30

### Added
- New "Standings" tab displaying live Driver and Constructor championship standings.
- Configurable session notifications with custom timing.
- Light and dark theme support.
- Team-based themes for all 10 constructors!
- Screenshots section added to the README and `/screenshots/` directory.

### Changed
- Refactored UI layout for better consistency and theming.
- Theme selection dropdown now prioritizes "System", "Dark", and "Light", followed by teams alphabetically.
- Lazy-loaded data refresh to optimize performance on startup.
- README reorganized: Screenshots moved below Features; image layout updated.

### Fixed
- Session time display issues in certain time zones.
- Notification delivery reliability on Linux and Windows.
- Minor Jolpica API response handling bugs.

---

## [0.1.2] - 2025-03-17

### Added
- Initial Upcoming tab with live session info.
- Schedule and Results views.
- Early support for Jolpica API integration.
