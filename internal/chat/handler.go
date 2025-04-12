package chat

import (
	"fmt"
	"net/http"
	"time"
)

var ChatRoomInstance = NewChatRoom()

func JoinHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing id", http.StatusBadRequest)
		return
	}

	ChatRoomInstance.mutex.Lock()
	_, exists := ChatRoomInstance.clients[id]
	ChatRoomInstance.mutex.Unlock()

	if exists {
		http.Error(w, "Client already joined", http.StatusConflict)
		return
	}

	client := &Client{ID: id, MsgChan: make(chan string, 10)}
	ChatRoomInstance.joinChan <- client
	fmt.Fprintf(w, "Client %s joined\n", id)
	// w.Write([]byte("joined chat room"))
}

func LeaveHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	ChatRoomInstance.leaveChan <- id
	exists := <-ChatRoomInstance.leaveAskChan

	if exists {
		fmt.Fprintf(w, "Client %s left\n", id)
	} else {
		http.Error(w, "Client not found or already left", http.StatusNotFound)
	}
	// w.Write([]byte("chat room left"))
}

func SendHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	msg := r.URL.Query().Get("message")
	if id == "" || msg == "" {
		http.Error(w, "missing id or message", http.StatusBadRequest)
		return
	}
	// Check if the client is still connected
	ChatRoomInstance.mutex.Lock()
	_, ok := ChatRoomInstance.clients[id]
	ChatRoomInstance.mutex.Unlock()

	if !ok {
		http.Error(w, "Client not found. Please join first", http.StatusForbidden)
		return
	}

	formattedMsg := fmt.Sprintf("[%s]:  %s", id, msg)
	ChatRoomInstance.broadcastChan <- formattedMsg
	// w.Write([]byte("message sent"))
	fmt.Fprintf(w, "Message sent: %s\n", formattedMsg)
}

func MessageHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing id", http.StatusBadRequest)
		return
	}
	ChatRoomInstance.mutex.Lock()
	client, ok := ChatRoomInstance.clients[id]
	ChatRoomInstance.mutex.Unlock()

	if !ok {
		http.Error(w, "Client not found", http.StatusNotFound)
		return
	}

	select {
	case msg := <-client.MsgChan:
		fmt.Fprint(w, msg)
	case <-time.After(10 * time.Second): // timeout to avoid indefinite block
		http.Error(w, "No new messages", http.StatusRequestTimeout)
	}
}
