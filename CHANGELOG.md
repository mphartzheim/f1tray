# Changelog

All notable changes to this project are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),  
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

---

## [v0.3.1] - 2025-xx-xx

### Notes
- Store list of all constructors available via Jolpica so we can reference their URL.
- Favorite constructors are now implemented app wide.
- Constructor team links are now available anywhere you see them.
- Constructor list only loads if the file is not present.

---

## [v0.3.0] - 2025-03-31

### Added
- **Favorite Constructor:** Introduced the initial phase of the Favorite Constructor feature for a personalized user experience.
- **Interactive Labels:** Made all cells clickable, enhancing navigation and interactivity.

### Fixed
- **Constructor Standings:** Addressed issues in the Constructor standings tab caused by the Favorite Drivers implementation.
- **F1TV Button Border:** Corrected the border color of the "Watch on F1TV" button for consistent design.

### Changed
- **Driver Name Code:** Refactored driver name code to improve reusability and maintainability.
- **UI Cleanup:** Removed unnecessary visibility from the "Watch on F1TV" button to enhance interface consistency.
- **Window Resizing:** Disabled window resizing to ensure a uniform user experience across platforms.

### Automation
- **Build & Release:** Implemented automation for building and releasing, facilitating smoother and more reliable deployments.

---

## [v0.2.3] - 2025-xx-xx

### Notes
- *Details for this release are forthcoming.*

---

## [0.2.2-test] - 2025-03-30

### Notes
- **Test Release for Automated Builds:**  
  This test release focuses on automating the build process, refining workflows from PR creation to building.

---

## [0.2.2] - 2025-03-30

### Added
- **"Watch on F1TV" Button:** Enhanced prominence for easier access to live action.
- **Favorite Drivers Highlight:** Enabled users to spotlight their top drivers with a simple click.
- **Global Notification System:** Implemented alerts to inform users when selecting too many favorites.

---

## [0.2.1] - 2025-03-30

### Added
- **Clickable Driver Names:** Driver names in standings and results now open bio pages, with a fallback to the API URL.
- **Helper Functions:** Added functions to reduce parser repetition and a theme lookup helper using `themes.AvailableThemes`.

### Fixed
- **Tray Menu Options:** Ensured correct display and responsiveness across themes.
- **Team Themes:** Corrected highlight colors application on initial load.

---

## [0.2.0] - 2025-03-30

### Added
- **"Standings" Tab:** Displays live Driver and Constructor championship standings.
- **Session Notifications:** Configurable notifications with custom timing.
- **Theme Support:** Added light, dark, and team-based themes for all 10 constructors.
- **Screenshots:** Included a screenshots section in the README and `/screenshots/` directory.

### Changed
- **UI Layout:** Refactored for better consistency and theming.
- **Theme Selection Dropdown:** Prioritized "System", "Dark", and "Light" themes, followed by teams alphabetically.
- **Data Refresh:** Implemented lazy-loading to optimize startup performance.
- **README Reorganization:** Moved screenshots below Features and updated image layout.

### Fixed
- **Session Time Display:** Resolved issues in certain time zones.
- **Notification Delivery:** Improved reliability on Linux and Windows.
- **API Response Handling:** Fixed minor bugs with Jolpica API responses.

---

## [0.1.2] - 2025-03-17

### Added
- **Upcoming Tab:** Displays live session information.
- **Schedule and Results Views:** Provides comprehensive race schedules and results.
- **Jolpica API Integration:** Early support for enhanced data retrieval.
