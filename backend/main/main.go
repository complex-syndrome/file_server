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

	// File Ops
	http.HandleFunc(helper.ApiPath+helper.ListCommand, handlers.ListFilesHandler)
	http.HandleFunc(helper.ApiPath+helper.UploadCommand, handlers.UploadFileHandler)
	http.HandleFunc(helper.ApiPath+helper.DownloadCommand, handlers.DownloadFileHandler)
	http.HandleFunc(helper.ApiPath+helper.DeleteCommand, handlers.DeleteFileHandler)

	// Settings
	http.HandleFunc(helper.ApiPath+helper.SettingsCommand, handlers.SettingsHandler)

	// Notify when file changes
	broadcast := make(chan string)
	go handlers.WatchFiles(broadcast, helper.CleanedResourcePath) // Watch folder change

	// Broadcast to all connections
	go handlers.Broadcaster(broadcast) 
	http.HandleFunc(helper.ApiPath+helper.WsNotifyCommand, 
		func(w http.ResponseWriter, r *http.Request) { handlers.FSChangeWebsocket(broadcast, w, r) })

	log.Printf("Server started at http://%s:%d%s\n", helper.GetMyIP().String(), helper.Port, helper.ApiPath)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", helper.Port), nil))

	<-make(chan struct{})
}
