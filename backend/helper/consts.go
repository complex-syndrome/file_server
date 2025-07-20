package helper

const (
	// Root for api
	ApiPath string = "/api"

	// Commands
	// HelpCommand     = "/help"
	listCommand     = "/list"
	uploadCommand   = "/upload"
	downloadCommand = "/download"
	deleteCommand   = "/delete"

	settingsCommand = "/settings"
	updateCommand   = "/update"
	allowCommand    = "/allow"

	WsNotifyCommand = "/ws"

	// WS labels
	FSLabel       = "Resource Folder"
	SettingsLabel = "Settings"

	// URLS
	ListFilesURL    = ApiPath + listCommand
	UploadFileURL   = ApiPath + uploadCommand
	DownloadFileURL = ApiPath + downloadCommand
	DeleteFileURL   = ApiPath + deleteCommand

	ListSettingsURL    = ApiPath + settingsCommand + listCommand
	UpdateSettingsURL  = ApiPath + settingsCommand + updateCommand
	AllowIPSettingsURL = ApiPath + settingsCommand + allowCommand

	WebSocketURL = ApiPath + WsNotifyCommand
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

type IpJSON struct {
	IP string `json:"ip"`
}
