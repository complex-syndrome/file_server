package helper

import (
	"log"
	"net"
	"net/http"
	"path/filepath"
	"strings"
)

func FromInvalidIPs(addr string, important bool) bool {
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return true
	}
	
	if net.ParseIP(host).IsLoopback() || host == GetMyIP().String() {
		return false
	}
	
	if important { // Stuff like settings should be only editable by the host machine or webui
		return true
	}
	
	val, ok := GetCurrentSettings("AllowOtherIPs")
	if ok {
		if AllowOtherIPs, ok := val.(bool); ok {
			return !AllowOtherIPs

		} else {
			log.Println("Unable to get current settings.")
			return true
		}
		
	} else {
		log.Println("AllowOtherIPs is not a boolean.")
		return true
	}
}


func IsInvalidFileName(fileName string, safeFileName string) bool {
	if safeFileName == "" {
		safeFileName = filepath.Base(fileName)
	}
	return fileName != safeFileName ||
		safeFileName == "" ||
		strings.Contains(safeFileName, "..") ||
		strings.Contains(safeFileName, "/") ||
		strings.Contains(safeFileName, "\\")
}

func WithCORS(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !FromInvalidIPs(r.RemoteAddr, true) {
			w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
		}

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		handler.ServeHTTP(w, r)
	})
}