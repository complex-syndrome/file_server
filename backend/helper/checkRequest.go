package helper

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func IsFrontendRequest(r *http.Request) bool {
	return r.Header.Get("X-From-Frontend") == os.Getenv("VITE_CUSTOM_VALUE")
}

func FromLocalhost(host string) bool {
	return net.ParseIP(host).IsLoopback() || host == GetMyIP().String()
}

func ValidRequest(r *http.Request, important bool) bool {
	fmt.Println("AllowOtherIPs:", GetCurrentSettings("AllowOtherIPs"))
	fmt.Println("IsFrontendRequest:", IsFrontendRequest(r))

	head := r.Header.Get("X-Forwarded-For")
	fmt.Println("X-Forwarded-For:", head)
	fmt.Println("FromLocalhost:", FromLocalhost(head))

	AllowOtherIPs, ok := GetCurrentSettings("AllowOtherIPs").(bool)
	if !ok {
		log.Println("Unable to get current settings.")
		return false
	}

	return FromLocalhost(r.Header.Get("X-Forwarded-For")) || IsFrontendRequest(r) || (!important && AllowOtherIPs)
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
