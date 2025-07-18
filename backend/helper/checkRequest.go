package helper

import (
	"log"
	"net"
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
