// main.go
package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/rs/cors"
	_ "modernc.org/sqlite"
)

// Message structure
type Message struct {
	ID      int    `json:"id"`
	Content string `json:"content"`
}

var (
	db        *sql.DB
	clients   = make(map[*websocket.Conn]bool)
	broadcast = make(chan Message)
	upgrader  = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	clientsMutex sync.Mutex
)

func main() {
	// Initialize SQLite database
	initDB()

	// Set up router
	router := mux.NewRouter()

	// Define API routes
	router.HandleFunc("/api/messages", getMessages).Methods("GET")
	router.HandleFunc("/api/messages", addMessage).Methods("POST")
	router.HandleFunc("/ws", handleWebSocket)

	// Set up CORS middleware
	corsHandler := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000"},
		AllowedMethods: []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type"},
	}).Handler(router)

	// Start server
	go handleMessages()
	log.Fatal(http.ListenAndServe(":8080", corsHandler))
}

// Initialize SQLite database
func initDB() {
	var err error
	db, err = sql.Open("sqlite", "./messages.db")
	if err != nil {
		log.Fatal(err)
	}

	// Create messages table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS messages (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			content TEXT
		)`)
	if err != nil {
		log.Fatal(err)
	}
}

// Retrieve messages from the database
func getMessages(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, content FROM messages")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var result []Message
	for rows.Next() {
		var message Message
		err := rows.Scan(&message.ID, &message.Content)
		if err != nil {
			log.Fatal(err)
		}
		result = append(result, message)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// Add a new message to the database
func addMessage(w http.ResponseWriter, r *http.Request) {
	var newMessage Message
	json.NewDecoder(r.Body).Decode(&newMessage)

	_, err := db.Exec("INSERT INTO messages (content) VALUES (?)", newMessage.Content)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newMessage)

	// Broadcast the new message to all connected clients
	broadcast <- newMessage
}

// WebSocket handler
func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	// Register client
	clientsMutex.Lock()
	clients[conn] = true
	clientsMutex.Unlock()

	// Listen for close events
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			break
		}
	}

	// Unregister client on close
	clientsMutex.Lock()
	delete(clients, conn)
	clientsMutex.Unlock()
}

// Broadcast new messages to all connected clients
func handleMessages() {
	for {
		// Wait for a new message to broadcast
		message := <-broadcast

		// Send the message to all clients
		clientsMutex.Lock()
		for client := range clients {
			err := client.WriteJSON(message)
			if err != nil {
				log.Println(err)
				client.Close()
				delete(clients, client)
			}
		}
		clientsMutex.Unlock()
	}
}
