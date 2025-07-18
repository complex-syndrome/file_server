package helper


const (
	ResourcePath string = "../uploads"
	ApiPath      string = "/api"
	Port         uint64 = 8080
	SettingsJSON string = "./settings.json"

	HelpCommand     = "/help"
	ListCommand     = "/list"
	UploadCommand   = "/upload"
	DownloadCommand = "/download"
	DeleteCommand   = "/delete"

	SettingsCommand = "/settings"
	WsNotifyCommand = "/ws"

	FSLabel       = "Resource Folder"
	SettingsLabel = "Settings"
)

var (
	defaultSettings = map[string]any{
		"AllowOtherIPs": false,
	}
	CurrentSettings map[string]any

	MaxUploadSize int64 = TranslateSize("5GB")

	CleanedResourcePath string = CleanPath(ResourcePath)
	CleanedSettingsPath string = CleanPath(SettingsJSON)
)

type FileInfo struct {
	Name string `json:"name"`
	Size string `json:"size"`
	Mime string `json:"mime"`
}
