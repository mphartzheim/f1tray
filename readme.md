# F1 Tray Application

A cross-platform desktop tray application built in Go using the Fyne GUI framework. It provides quick access to real-time Formula 1 session results, schedules, and standings using the Jolpica Ergast API.

---

## 🚀 Features

- View full race schedule
- Check live or recent session results:
  - Practice 1, 2, 3
  - Qualifying
  - Sprint Qualifying
  - Sprint
  - Race Results
- Scrollable, compact GUI with back-navigation support
- Debug mode with extra UI elements

---

## 🛠 Tech Stack

- **Go (Golang)**
- **Fyne** for GUI
- **Jolpica Ergast API** for F1 data

---

## 📁 Folder Structure

```
/f1tray
├── assets/                  # Tray icon and other images
├── cmd/
│   └── f1tray/
│       └── main.go          # Main application entry point
├── internal/
│   ├── api/
│   │   ├── jolpica.go       # API logic
│   │   ├── models.go        # JSON models for Ergast data
│   │   └── jolpica_test.go  # API unit tests
│   ├── config/
│   │   └── config.go        # App-wide constants
│   ├── gui/
│   │   ├── assets.go        # Icon loading
│   │   ├── menu.go          # Main menu UI
│   │   ├── navigation.go    # Shared navigation helpers
│   │   ├── schedule.go      # Race schedule view
│   │   └── session_view.go  # Session results view
│   └── tray/                # [Planned] Tray integration logic
├── go.mod
├── go.sum
└── README.md
```

---

## 🧪 Testing

Unit tests for the API module live in:
```
internal/api/jolpica_test.go
```
Run tests with:
```bash
go test ./...
```

---

## 🧩 Planned Features

- Session reminders via system notifications
- User preferences system
- Customizable refresh intervals
- Highlighting next/upcoming race in the schedule
- Native system tray integration

---

## ✅ How to Run

1. Clone the repository
2. Install Go and run:
```bash
cd cmd/f1tray
go run .
```

---

## ✨ License
MIT

