package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"

	"fyne.io/fyne/v2/app"
)

var (
	configLock sync.Mutex
	config     Preferences
	loaded     = false
)

func Get() Preferences {
	configLock.Lock()
	defer configLock.Unlock()

	if !loaded {
		loadConfig()
	}
	return config
}

func Set(p Preferences) {
	configLock.Lock()
	config = p
	configLock.Unlock()
	Save()
}

func Save() {
	configLock.Lock()
	defer configLock.Unlock()

	path := getConfigPath()
	file, err := os.Create(path)
	if err != nil {
		return
	}
	defer file.Close()

	_ = json.NewEncoder(file).Encode(config)
}

func loadConfig() {
	path := getConfigPath()
	file, err := os.Open(path)
	if err != nil {
		config = defaultPreferences()
		loaded = true
		return
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		config = defaultPreferences()
	} else {
		config = SanitizePreferences(config)
		loaded = true
	}
}

func getConfigPath() string {
	storage := app.NewWithID("com.f1tray.app").Storage()
	dir := storage.RootURI().Path()
	return filepath.Join(dir, "preferences.json")
}
