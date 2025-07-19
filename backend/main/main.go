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
	fmt.Println("Upload Path: " + helper.ResourcePath)
	fmt.Println("Settings Path: " + helper.SettingsPath)
	fmt.Println("Max upload size: " + helper.CalculateSize(helper.MaxUploadSize))

	// Handlers
	http.HandleFunc(helper.ApiPath+helper.ListCommand, handlers.ListFilesHandler)
	http.HandleFunc(helper.ApiPath+helper.UploadCommand, handlers.UploadFileHandler)
	http.HandleFunc(helper.ApiPath+helper.DownloadCommand, handlers.DownloadFileHandler)
	http.HandleFunc(helper.ApiPath+helper.DeleteCommand, handlers.DeleteFileHandler)
	http.HandleFunc(helper.ApiPath+helper.SettingsCommand, handlers.SettingsHandler)

	// Fsnotify + Websocket
	fanOut := &helper.FanOut{}
	go handlers.WatchFiles(fanOut.Publish, helper.ResourcePath, helper.FSLabel)       // Watch folder change (pub)
	go handlers.WatchFiles(fanOut.Publish, helper.SettingsPath, helper.SettingsLabel) // Watch file change (pub)

	go handlers.RefreshSettingsOnChange(fanOut.Subscribe(), helper.SettingsLabel) // Refresh settings on change (sub)
	go handlers.Broadcaster(fanOut.Subscribe())                                   // Broadcast change on change (sub)

	http.HandleFunc(helper.ApiPath+helper.WsNotifyCommand, // Websocket handler
		func(w http.ResponseWriter, r *http.Request) { handlers.FSChangeWebsocket(fanOut.Subscribe(), w, r) })

	// Log again
	log.Printf("Backend server started at http://%s:%d%s\n", helper.GetMyIP().String(), helper.BackendPort, helper.ApiPath)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", helper.BackendPort), helper.WithCORS(http.DefaultServeMux)))

	<-make(chan struct{})
}
