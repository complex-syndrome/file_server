package handlers

// Handling websocket connections

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/gorilla/websocket"
)

var (
	upgrader websocket.Upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	clients = make(map[*websocket.Conn]bool)
	mutex   sync.Mutex
)

func FSChangeWebsocket(nchan <-chan string, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WS Upgrade Error:", err)
		return
	}
	defer conn.Close()

	// register client
	mutex.Lock()
	log.Println("WS: Client", r.RemoteAddr, "connected.")
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
	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			break
		}
	}
}

func Broadcaster(nchan <-chan string) {
	for {
		msg := <-nchan // When have msg, broadcast

		mutex.Lock()
		for conn := range clients {
			conn.SetWriteDeadline(time.Now().Add(2 * time.Second))
			if err := conn.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
				log.Printf("WebSocket write failed: %v", err)
				conn.Close()
				delete(clients, conn)
			} else {
				log.Printf("WS sent message to client: %s", msg)
			}
		}
		mutex.Unlock()
	}
}

func WatchFiles(publish func(msg string), path string, label string) { // chan<- only receive
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
			log.Println("FS Watch Event:", event)
			publish(label + ": " + event.String())

		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Println("FS Watch Error:", err)
		}
	}
}
