package helper

const (
	ApiPath string = "/api"

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
	ResourcePath  string = "../uploads"
	SettingsPath  string = "../settings.json"
	BackendPort   uint64 = 8080
	FrontendPort  uint64 = 5173
	MaxUploadSize int64  = TranslateSize("5GB")

	defaultSettings = map[string]any{
		"AllowOtherIPs": false,
	}
	CurrentSettings map[string]any
)

type FileInfo struct {
	Name string `json:"name"`
	Size string `json:"size"`
	Mime string `json:"mime"`
}
