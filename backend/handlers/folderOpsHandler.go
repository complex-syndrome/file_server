package handlers

// Handling api calls for folder operations

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gabriel-vasile/mimetype"

	"github.com/complex-syndrome/file-server/backend/helper"
)

func ListFilesHandler(w http.ResponseWriter, r *http.Request) {
	// Ensure local access only (or setted otherwise)
	if !helper.ValidRequest(r, false) {
		http.Error(w, "Access Denied: Local Connections Only", http.StatusForbidden)
		log.Printf("Folder Ops: Failed attempt to access by address: %s\n", r.RemoteAddr)
		return
	}

	// GET
	if r.Method != http.MethodGet {
		http.Error(w, "Method not Allowed", http.StatusMethodNotAllowed)
		return
	}
	log.Printf("List file request from %s\n", r.RemoteAddr)

	// Read folder
	files, err := os.ReadDir(helper.ResourcePath)
	if err != nil {
		http.Error(w, "Unable to read saved files", http.StatusInternalServerError)
		log.Printf("ReadDir error (%s): %v\n", helper.ResourcePath, err)
		return
	}

	var fileInfos []helper.FileInfo
	for _, file := range files {
		// File issue
		info, err := file.Info()
		if err != nil {
			log.Println("Skipping file:", file.Name(), err)
			continue
		}

		// Neglect Directory
		if info.IsDir() {
			log.Println("Skipping dir:", file.Name())
			continue
		}

		// Mime type
		filePath := filepath.Join(helper.ResourcePath, info.Name())
		mime, err := mimetype.DetectFile(filePath)
		if err != nil {
			log.Printf("Mime-type detection error (%s): %v\n", info.Name(), err)
			continue
		}

		ext := mime.Extension()
		if ext == "" {
			ext = ".*"
		}

		// As json
		fileInfos = append(fileInfos, helper.FileInfo{
			Name: info.Name(),
			Size: helper.CalculateSize(info.Size()),
			Mime: fmt.Sprintf("%s (%s)", mime.String(), ext),
		})
	}
	helper.ReplyJSON(w, fileInfos)
}

func UploadFileHandler(w http.ResponseWriter, r *http.Request) {
	// Ensure local access only (or setted otherwise)
	if !helper.ValidRequest(r, false) {
		http.Error(w, "Access Denied: Local Connections Only", http.StatusForbidden)
		log.Printf("Folder Ops: Failed attempt to access by address: %s\n", r.RemoteAddr)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not Allowed", http.StatusMethodNotAllowed)
		return
	}
	log.Printf("Upload request from: %s\n", r.RemoteAddr)

	// File Size
	r.Body = http.MaxBytesReader(w, r.Body, helper.MaxUploadSize)
	err := r.ParseMultipartForm(helper.MaxUploadSize)
	if err != nil {
		http.Error(w, "Error uploading file", http.StatusRequestEntityTooLarge)
		log.Println("Upload file error:", err)
		return
	}

	// Form file
	src, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error uploading file", http.StatusInternalServerError)
		log.Println("Error forming file:", err)
		return
	}
	defer src.Close()

	// Traversal
	safeFileName := filepath.Base(handler.Filename)
	if helper.IsInvalidFileName(handler.Filename, safeFileName) {
		http.Error(w, "Invalid file name: "+handler.Filename, http.StatusBadRequest)
		log.Println("Invalid file name:", handler.Filename)
		return
	}

	// Handle overwrite
	dstPath := filepath.Join(helper.ResourcePath, safeFileName)
	i := 1
	for {
		if _, err := os.Stat(dstPath); err != nil {
			break
		}
		dstPath = filepath.Join(helper.ResourcePath, fmt.Sprintf("%s.%d", safeFileName, i))
		i++
	}

	dst, err := os.Create(dstPath)
	if err != nil {
		http.Error(w, "Error saving file: "+handler.Filename, http.StatusInternalServerError)
		log.Println("Save error:", err)
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		http.Error(w, "Failed to save file: "+handler.Filename, http.StatusInternalServerError)
		log.Println("Copy error:", err)
		return
	}

	fmt.Fprintf(w, "Upload successful: %s\n", handler.Filename)
	log.Printf("Successful upload request from: %s (%s)\n", r.RemoteAddr, handler.Filename)
}

func DownloadFileHandler(w http.ResponseWriter, r *http.Request) {
	// Ensure local access only (or setted otherwise)
	if !helper.ValidRequest(r, false) {
		http.Error(w, "Access Denied: Local Connections Only", http.StatusForbidden)
		log.Printf("Folder Ops: Failed attempt to access by address: %s\n", r.RemoteAddr)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "Method not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Empty
	fileName := r.URL.Query().Get("file")
	if fileName == "" {
		http.Error(w, "Url missing 'file' parameter", http.StatusBadRequest)
		return
	}

	log.Printf("Download request from: %s\n", r.RemoteAddr)
	// Safe
	safeFileName := filepath.Base(fileName)
	if fileName != safeFileName || helper.IsInvalidFileName(fileName, safeFileName) {
		http.Error(w, "Invalid File Name: "+fileName, http.StatusBadRequest)
		log.Printf("Download blocked for filename: %s\n", fileName)
		return
	}

	// Not found
	filePath := filepath.Join(helper.ResourcePath, safeFileName)
	if _, err := os.Stat(filePath); err != nil {
		http.Error(w, "File not found: "+fileName, http.StatusNotFound)
		log.Printf("Download Error: %v\n", err)
		return
	}

	http.ServeFile(w, r, fmt.Sprintf("%s/%s", helper.ResourcePath, fileName))
	fmt.Fprintf(w, "Download successful: %s\n", fileName)
	log.Printf("Successful download request from: %s (%s)\n", fileName, r.RemoteAddr)
}

func DeleteFileHandler(w http.ResponseWriter, r *http.Request) {
	// Ensure local access only, or require password from webui
	if !helper.ValidRequest(r, true) {
		http.Error(w, "Access Denied: Local Connections Only", http.StatusForbidden)
		log.Printf("Folder Ops: Failed attempt to access by address: %s\n", r.RemoteAddr)
		return
	}

	// DELETE
	if r.Method != http.MethodDelete && r.Method != http.MethodPost {
		http.Error(w, "Method not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Empty
	fileName := r.URL.Query().Get("file")
	if fileName == "" {
		http.Error(w, "Url missing 'file' parameter", http.StatusBadRequest)
		return
	}

	log.Printf("Delete request from: %s\n", r.RemoteAddr)
	// Safe
	safeFileName := filepath.Base(fileName)
	if helper.IsInvalidFileName(fileName, safeFileName) {
		http.Error(w, "Invalid file name: "+fileName, http.StatusBadRequest)
		log.Printf("Delete blocked for filename: %s\n", fileName)
		return
	}

	// Not found
	filePath := filepath.Join(helper.ResourcePath, safeFileName)
	if _, err := os.Stat(filePath); err != nil {
		http.Error(w, "File not found: "+fileName, http.StatusNotFound)
		log.Printf("Delete Error: %v\n", err)
		return
	}

	// Fail delete
	if err := os.Remove(filePath); err != nil {
		http.Error(w, "Failed to delete file: "+fileName, http.StatusNotFound)
		log.Printf("Deletion Failed Error: %v\n", err)
		return
	}

	fmt.Fprintf(w, "Deleted successfully: %s\n", fileName)
	log.Printf("Successful delete request from: %s (%s)\n", fileName, r.RemoteAddr)
}
