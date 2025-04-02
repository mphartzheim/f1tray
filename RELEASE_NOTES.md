## [v0.4.0] - 2025-04-01

### 🚀 Major Overhaul
- **🔨 Total Rebuild**: The F1Tray app was completely rewritten from scratch with modular architecture, cleaner logic, and maintainable code. Every feature was thoughtfully restructured for performance, readability, and scalability.
- **🧠 Smarter Structure**: Modular packages for `schedule`, `results`, `standings`, `upcoming`, and `layout` keep logic clean and easy to work with.

### ⚡ Performance Upgrades
- **📡 Concurrent Data Loading**: Tab data is now fetched in parallel using Go channels — dramatically improving performance at app startup and when switching seasons.
- **🧠 Lazy Loading Tabs**: Content is only fetched when tabs are opened, saving bandwidth and load time.
- **🕹️ Debug Mode**: Toggle debug logs with a command-line flag. See exactly what’s loading and when.

### 🧭 Navigation Enhancements
- **🗓️ Year Selector**: Quickly switch between seasons with an intuitive dropdown.
- **📊 Standings Tabs**: View driver and constructor standings from any season, instantly.
- **📅 Schedule + Upcoming**: Check session dates, circuits, and countdowns to the next event — all fetched live.

### 🔧 Stability & UX
- **💥 Tab Isolation**: Crashes in one tab won't affect others. Each tab manages its own state and failures.
- **⏱️ Countdown Timer**: Shows time until the next session, updated in real-time.
- **📁 Cleaner Error Handling**: UI provides graceful fallbacks for failed API calls.

### 🐞 Known Issues
- This is still a WIP — bugs may exist and polish is ongoing. But the core is **solid**.
