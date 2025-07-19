package helper

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/joho/godotenv"
)

func ImportEnvs() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Println("Error loading .env file. Setting all configs to default.")
	}

	newBPort, err := strconv.ParseUint(os.Getenv("BACKEND_PORT"), 10, 64)
	if err != nil {
		log.Printf("BACKEND_PORT env value not found, using default port %d.", BackendPort)
		log.Println(err)
	} else {
		BackendPort = newBPort
	}

	newFPort, err := strconv.ParseUint(os.Getenv("FRONTEND_PORT"), 10, 64)
	if err != nil {
		log.Printf("FRONTEND_PORT env value not found, using default port %d.", FrontendPort)
		log.Println(err)
	} else {
		FrontendPort = newFPort
	}

	newResourcePath := os.Getenv("UPLOADS_FOLDER")
	if newResourcePath == "" {
		log.Printf("UPLOADS_FOLDER env value not found, using default path %s.", ResourcePath)
		log.Println(err)
	} else {
		if strings.HasPrefix(newResourcePath, "./") || strings.HasPrefix(newResourcePath, "../") {
			newResourcePath = filepath.Join("..", newResourcePath)
		}
		ResourcePath = newResourcePath
	}

	if newMaxUploadSize := TranslateSize(os.Getenv("MAX_UPLOAD_SIZE")); newMaxUploadSize == -1 {
		log.Printf("MAX_UPLOAD_SIZE env value not found, using default value %s.", CalculateSize(MaxUploadSize))
		log.Println(err)
	} else {
		MaxUploadSize = newMaxUploadSize
	}
}

func GenerateCleanedPaths() {
	ResourcePath = CleanPath(ResourcePath)
	SettingsPath = CleanPath(SettingsPath)
}

func IndexOf(slice []string, find string) int {
	for k, v := range slice {
		if find == v {
			return k
		}
	}
	return -1
}

func GetMyIP() net.IP {
	conn, err := net.Dial("udp", "1.1.1.1:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	return conn.LocalAddr().(*net.UDPAddr).IP
}

func TryMkdir(path string) {
	log.Println("Path:", path)
	if err := os.MkdirAll(path, 0777); err != nil {
		log.Fatalf("Could not create directory: %v", err)
	}
}

func CleanPath(path string) string {
	cPath, err := filepath.Abs(path)
	if err != nil {
		log.Fatal("Path cleaning error:", err)
	}
	return cPath
}

func CalculateSize(size int64) string {
	if size < 1024 {
		return fmt.Sprintf("%d B", size)
	}
	units := []string{"KB", "MB", "GB"}
	s := float64(size)
	for i, unit := range units {
		s /= 1024
		if s < 1024 || i == len(units)-1 {
			return fmt.Sprintf("%.2f %s", s, unit)
		}
	}
	return fmt.Sprintf("%d B", size)
}

func TranslateSize(size string) int64 {
	size = strings.ToUpper(strings.TrimSpace(size))
	units := []string{"B", "KB", "MB", "GB"}
	re := regexp.MustCompile(`(?i)^([\d.]+)\s*(B|KB|MB|GB)$`)
	matches := re.FindStringSubmatch(size)
	if len(matches) != 3 {
		log.Printf("Invalid size to translate: %s\n", size)
		return -1
	}

	value, err := strconv.ParseFloat(matches[1], 64)
	if err != nil {
		log.Printf("Invalid size to translate: %s\n", size)
		return -1
	}
	return int64(value * math.Pow(1024, float64(IndexOf(units, matches[2]))))
}

func ReplyJSON(w http.ResponseWriter, json_obj any) {
	w.Header().Set("Content-type", "application/json")
	if err := json.NewEncoder(w).Encode(json_obj); err != nil {
		http.Error(w, err.Error(), 500)
		log.Println("JSON encode error:", err)
	}
}

func ReadSettingsJSON() map[string]any {
	data, err := os.ReadFile(SettingsPath)
	if err != nil {
		log.Printf("Unable to read %s. Generating default settings...\n", SettingsPath)
		WriteSettings(defaultSettings)
	}

	var settings map[string]any
	if err = json.Unmarshal(data, &settings); err != nil {
		log.Printf("Error during decoding %s: %v\n", SettingsPath, err)
		return defaultSettings
	}
	return settings
}

var settingsMutex sync.RWMutex

func WriteSettings(newSettings map[string]any) {
	settingsMutex.Lock()
	defer settingsMutex.Unlock()

	data, err := json.MarshalIndent(newSettings, "", "  ")
	if err != nil {
		log.Println("Unable to encode new setings into JSON format")
		return
	}

	if err := os.WriteFile(SettingsPath, data, 0644); err != nil {
		log.Fatalf("Error writing settings to %s: %v", SettingsPath, err)
	}
}

func RefreshSettings() {
	settingsMutex.Lock()
	defer settingsMutex.Unlock()

	CurrentSettings = ReadSettingsJSON()

	log.Println()
	log.Println("Current Settings:")
	for k, v := range CurrentSettings {
		log.Printf("%s = %v\n", k, v)
	}
	log.Println()
}

func GetCurrentSettings(key string) (any, bool) {
	settingsMutex.RLock()
	defer settingsMutex.RUnlock()
	val, ok := CurrentSettings[key]
	return val, ok
}
