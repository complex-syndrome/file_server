package handlers

import (
	"log"
	"sync"
	"net/http"

	"github.com/fsnotify/fsnotify"
	"github.com/gorilla/websocket"

	"github.com/complex-syndrome/file-server/backend/helper"
)

var (
	upgrader websocket.Upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	clients   = make(map[*websocket.Conn]bool)
	mutex     sync.Mutex
)


func FSChangeWebsocket(nchan <-chan string, w http.ResponseWriter, r *http.Request) {
	if helper.FromInvalidIPs(r.RemoteAddr, true) {
		http.Error(w, "Access Denied: Local Connections Only", http.StatusForbidden)
		log.Printf("Websocket: Failed attempt to access by address: %s\n", r.RemoteAddr)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade Error:", err)
		return
	}
	defer conn.Close()

	// register client
	mutex.Lock()
	clients[conn] = true
	mutex.Unlock()

	// Remove client
	defer func() {
		mutex.Lock()
		delete(clients, conn)
		mutex.Unlock()
		conn.Close()
	}()

	// Keeps conn alive
	for {if _, _, err := conn.ReadMessage(); err != nil { break }}
}

func Broadcaster(nchan <-chan string) {
	for {
		msg := <-nchan // When have msg, broadcast
		log.Println("Websocket: ", msg)
		
		mutex.Lock()
		for conn := range clients {
			if err := conn.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
				conn.Close()
				delete(clients, conn)	
			}
		}
		mutex.Unlock()
	}
}


func WatchFiles(nchan chan<- string, path string) { // chan<- only receive
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	err = watcher.Add(path)
	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case event, ok := <-watcher.Events: // <-Events only send
			if !ok {
				return
			}
			log.Println("Watch Event:", event)
			nchan <- "(FS) " + event.String()

		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Println("Watch Error:", err)
		}
	}
}