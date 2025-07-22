package handlers

import (
	"crypto/sha512"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"github.com/complex-syndrome/file-server/backend/helper"
)

type PasswordJSON struct {
	Password string `json:"password"`
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if !helper.ValidRequest(r, true) { // For local or from webui only
		http.Error(w, "Access Denied: Local Connections Only", http.StatusForbidden)
		log.Printf("Login: Failed attempt to access by address: %s\n", r.RemoteAddr)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Method not Allowed", http.StatusMethodNotAllowed)
		return
	}
	log.Printf("Setting verify request from: %s\n", r.RemoteAddr)

	var req PasswordJSON
	var decoder = json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		http.Error(w, "An error occured while decoding JSON", http.StatusBadRequest)
		log.Println("JSON Decoding Error: ", err)
		return
	}

	sum := sha512.Sum512([]byte(req.Password))
	if err := bcrypt.CompareHashAndPassword(helper.Password, sum[:]); err != nil {
		log.Println("Rejected invalid password.")
		http.Error(w, "Invalid password", http.StatusForbidden)
	} else {
		log.Println("Allowed valid password.")
		fmt.Fprintln(w, "Valid password")
	}
}
