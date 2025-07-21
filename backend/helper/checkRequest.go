package helper

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"path/filepath"
	"strings"
)

func FromInvalidIPs(addr string, localConnectionsOnly bool) bool {
	if !strings.Contains(addr, ":") { addr += ":1234"; } // add dummy port if not in format x.x.x.x:x
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		fmt.Println("log: ERROR", err)
		return true
	}
	
	if net.ParseIP(host).IsLoopback() || host == GetMyIP().String() {
		return false
	}
	
	if localConnectionsOnly { // Enforce restriction to localhost only (e.g., for settings endpoints)
		return true
	}
	
	val, ok := GetCurrentSettings("AllowOtherIPs")
	if ok {
		if AllowOtherIPs, ok := val.(bool); ok {
			return !AllowOtherIPs
			
		} else {
			log.Println("AllowOtherIPs is not a boolean.")
			return true
		}
		
	} else {
		log.Println("Unable to get current settings.")
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
		val, ok := GetCurrentSettings("AllowOtherIPs")
		
		if !ok {
			log.Println("AllowOtherIPs is not a boolean.")
			return
		}

		AllowOtherIPs, ok := val.(bool)
		if !ok {
			log.Println("Unable to get current settings.")	
			return
		}

		if !FromInvalidIPs(r.RemoteAddr, false) || AllowOtherIPs { // If localhost, proxy / webUI
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