# Release Notes – v0.1.2-dev.1

**Release Date:** March 29, 2025  
**Branch:** `dev`  
**Tag:** `v0.1.2-dev.1`

---

## 🚀 Highlights

### 📊 New Standings Tab!
A brand new **"Standings"** tab has been added to the app, providing quick access to:
- 🧑‍💼 **Driver Standings**
- 🏎️ **Constructor Standings**

Both sub-tabs are dynamically populated using the Jolpica API and reflect the currently selected year. Just like the "Results" tab, data is displayed in a clean table format with labeled positions, points, teams, and nationalities.

---

## 🧠 Smart Year-Aware Updates

- Changing the year in the dropdown now **automatically refreshes the Standings tab** to reflect the newly selected season.
- Your current view (Drivers or Constructors) is **preserved** when the tab refreshes.

---

## ✳️ Notification Logic Improvements

- The **"Standings" tab now supports the asterisk indicator** (`*`) if its data becomes outdated or is updated in the background.
- This visual cue helps you know which tabs have unseen updates — consistent with how the "Results" tab already behaves.

---

## 🔧 Fixes & Tweaks

- Fixed a typo in `ConstructorsStandingsURL` constant.
- Improved consistency across parsing logic by centralizing and expanding response structs.
- Cleaned up redundant comments and clarified refresh responsibilities.

---

## 📁 Files of Note

- `models/standings.go` → New response structs for driver and constructor standings
- `processes/parsers.go` → Added `ParseDriverStandings` and `ParseConstructorStandings`
- `standings_tab.go` / `standings_main.go` → New UI structure for nested tabs

---

## 🧪 Dev Notes

This is a `-dev.1` pre-release as the new Standings tab is still being tested and refined. Future iterations may:
- Add sorting options
- Persist tab state across sessions
- Include additional driver/construction metadata

Feedback welcome!

---

