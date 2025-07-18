package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"net/http"

	"github.com/complex-syndrome/file-server/backend/helper"
)


func SettingsHandler(w http.ResponseWriter, r *http.Request) {	
	if helper.FromInvalidIPs(r.RemoteAddr, true) {
		http.Error(w, "Access Denied: Local Connections Only", http.StatusForbidden)
		log.Printf("Settings: Failed attempt to access by address: %s\n", r.RemoteAddr)
		return
	}

	if r.Method == http.MethodPost {
		log.Printf("Setting change request from: %s\n", r.RemoteAddr)
		editSettingsJSON(w, r)

	} else {
		log.Printf("Setting list request from: %s\n", r.RemoteAddr)
		helper.ReplyJSON(w, helper.CurrentSettings)
	}
}

func editSettingsJSON(w http.ResponseWriter, r *http.Request) {
	var newSettings map[string]any

	if err := json.NewDecoder(r.Body).Decode(&newSettings); err != nil {
		http.Error(w, "An error occured while changing settings", http.StatusBadRequest)
		log.Println("JSON Decoding Error: ", err)
	}

	updated := false
	for k, v := range newSettings {
		switch k {
		case "AllowOtherIPs": // Allow other IPs to directly access the API
			if bV, ok := v.(bool); ok {
				helper.CurrentSettings["AllowOtherIPs"] = bV
				updated = true
			}
		default:
			fmt.Fprintf(w, "Unknown setting: %s\n", k)
		}
	}

	if updated {
		helper.WriteSettings(newSettings)
		helper.RefreshSettings()
		
		fmt.Fprintln(w, "Settings successfully updated.")
		log.Printf("Settings successfully changed by %s.\n", r.RemoteAddr)

	} else {
		http.Error(w, "Failed to update settings.", http.StatusNotModified)
		log.Printf("Failed to change settings by %s.\n", r.RemoteAddr)
	}

}

func RefreshSettingsOnChange(nchan <-chan string, settingsLabel string) {
	for {
		msg := <-nchan
		if strings.HasPrefix(msg, settingsLabel) {
			helper.RefreshSettings()
		}
	}
}