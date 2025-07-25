package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/complex-syndrome/file-server/backend/handlers"
	"github.com/complex-syndrome/file-server/backend/helper"
)

func main() {
	// Envs
	helper.ImportEnvs()
	helper.GenerateCleanedPaths()

	// Resources
	helper.TryMkdir(helper.ResourcePath)
	helper.RefreshSettings()

	// Log
	log.Println("Upload Path: " + helper.ResourcePath)
	log.Println("Settings Path: " + helper.SettingsPath)
	log.Println("Max upload size: " + helper.CalculateSize(helper.MaxUploadSize))

	// Handlers
	http.HandleFunc(helper.LoginURL, handlers.LoginHandler)
	http.HandleFunc(helper.ListFilesURL, handlers.ListFilesHandler)
	http.HandleFunc(helper.UploadFileURL, handlers.UploadFileHandler)
	http.HandleFunc(helper.DownloadFileURL, handlers.DownloadFileHandler)
	http.HandleFunc(helper.DeleteFileURL, handlers.DeleteFileHandler)
	http.HandleFunc(helper.ListSettingsURL, handlers.ListSettingsHandler)
	http.HandleFunc(helper.UpdateSettingsURL, handlers.UpdateSettingsHandler)

	// Fsnotify + Websocket
	fanOut := &helper.FanOut{}
	go handlers.WatchFiles(fanOut.Publish, helper.ResourcePath, helper.FSLabel)       // Watch folder change (pub)
	go handlers.WatchFiles(fanOut.Publish, helper.SettingsPath, helper.SettingsLabel) // Watch file change (pub)

	go handlers.RefreshSettingsOnChange(fanOut.Subscribe(), helper.SettingsLabel) // Refresh settings on change (sub)
	go handlers.Broadcaster(fanOut.Subscribe())                                   // Broadcast change on change (sub)

	http.HandleFunc(helper.WebSocketURL, // Websocket handler
		func(w http.ResponseWriter, r *http.Request) { handlers.FSChangeWebsocket(fanOut.Subscribe(), w, r) })

	// Log
	log.Printf("Backend server started at http://%s:%d%s (Served for localhost only)\n", helper.GetMyIP().String(), helper.BackendPort, helper.ApiPath)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("127.0.0.1:%d", helper.BackendPort), nil))
	<-make(chan struct{})
}
