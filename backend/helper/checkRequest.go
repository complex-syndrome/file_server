package helper

import (
	"net"
	"strings"
	"path/filepath"
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
	
	return !AllowOtherIPs
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
