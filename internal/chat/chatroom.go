package chat

import (
	"fmt"
	"sync"
)

type ChatRoom struct {
	clients       map[string]*Client
	joinChan      chan *Client
	leaveChan     chan string
	leaveAskChan  chan bool
	broadcastChan chan string
	mutex         sync.Mutex
}

// Chat room initialization
func NewChatRoom() *ChatRoom {
	cr := &ChatRoom{
		clients:       make(map[string]*Client),
		joinChan:      make(chan *Client),
		leaveChan:     make(chan string),
		leaveAskChan:  make(chan bool),
		broadcastChan: make(chan string),
	}

	go cr.Run()
	return cr
}

// Goroutine loop
func (cr *ChatRoom) Run() {
	for {
		select {
		case client := <-cr.joinChan:
			cr.mutex.Lock()
			cr.clients[client.ID] = client
			cr.mutex.Unlock()
			fmt.Println("Client joined:", client.ID)

		case id := <-cr.leaveChan:
			cr.mutex.Lock()
			client, exists := cr.clients[id]
			if exists {
				delete(cr.clients, id)
				close(client.MsgChan) // closing channel
			}
			cr.mutex.Unlock()

			cr.leaveAskChan <- exists

			if exists {
				fmt.Println("Client left:", id)
			}

		case msg := <-cr.broadcastChan:
			cr.mutex.Lock()
			for _, client := range cr.clients {
				client.MsgChan <- msg
			}
			cr.mutex.Unlock()
		}
	}
}
