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
	loginCommand    = "/login"

	WsNotifyCommand = "/ws"

	// WS labels
	FSLabel       = "Resource Folder"
	SettingsLabel = "Settings"

	// URLS
	LoginURL = ApiPath + loginCommand

	ListFilesURL    = ApiPath + listCommand
	UploadFileURL   = ApiPath + uploadCommand
	DownloadFileURL = ApiPath + downloadCommand
	DeleteFileURL   = ApiPath + deleteCommand

	ListSettingsURL   = ApiPath + settingsCommand + listCommand
	UpdateSettingsURL = ApiPath + settingsCommand + updateCommand

	WebSocketURL = WsNotifyCommand
)

var (
	ResourcePath  string = "../uploads"
	SettingsPath  string = "../settings.json"
	BackendPort   uint64 = 8080
	FrontendPort  uint64 = 4173
	MaxUploadSize int64  = TranslateSize("5GB")
	Password      []byte

	defaultSettings = map[string]any{
		"AllowOtherIPs": true,
	}
	CurrentSettings map[string]any
)

type FileInfo struct {
	Name string `json:"name"`
	Size string `json:"size"`
	Mime string `json:"mime"`
}


