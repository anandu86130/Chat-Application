package main

import (
	"chat-application/internal/chat"
	"fmt"
	"log"
	"net/http"
)

func main() {
	// Register HTTP handlers
	http.HandleFunc("/join", chat.JoinHandler)
	http.HandleFunc("/leave", chat.LeaveHandler)
	http.HandleFunc("/send", chat.SendHandler)
	http.HandleFunc("/messages", chat.MessageHandler)
	// Start HTTP server
	port := ":8080"
	fmt.Println("chat server started...")
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
