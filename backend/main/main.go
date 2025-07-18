package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/complex-syndrome/file-server/backend/handlers"
	"github.com/complex-syndrome/file-server/backend/helper"
)

func main() {
	helper.TryMkdir(helper.CleanedResourcePath)
	fmt.Println("Upload Path: " + helper.CleanedResourcePath)
	helper.RefreshSettings()

	fanOut := &helper.FanOut{}

	go handlers.WatchFiles(fanOut.Publish, helper.CleanedResourcePath, helper.FSLabel)       // Watch folder change
	go handlers.WatchFiles(fanOut.Publish, helper.CleanedSettingsPath, helper.SettingsLabel) // Watch file change

	// File Ops
	http.HandleFunc(helper.ApiPath+helper.ListCommand, handlers.ListFilesHandler)
	http.HandleFunc(helper.ApiPath+helper.UploadCommand, handlers.UploadFileHandler)
	http.HandleFunc(helper.ApiPath+helper.DownloadCommand, handlers.DownloadFileHandler)
	http.HandleFunc(helper.ApiPath+helper.DeleteCommand, handlers.DeleteFileHandler)

	// Settings
	http.HandleFunc(helper.ApiPath+helper.SettingsCommand, handlers.SettingsHandler)
	go handlers.RefreshSettingsOnChange(fanOut.Subscribe(), helper.SettingsLabel) // Refresh on change

	// Broadcast to all connections when file / folder changes
	go handlers.Broadcaster(fanOut.Subscribe())
	http.HandleFunc(helper.ApiPath+helper.WsNotifyCommand,
		func(w http.ResponseWriter, r *http.Request) { handlers.FSChangeWebsocket(fanOut.Subscribe(), w, r) })

	log.Printf("Server started at http://%s:%d%s\n", helper.GetMyIP().String(), helper.Port, helper.ApiPath)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", helper.Port), nil))

	<-make(chan struct{})
}
