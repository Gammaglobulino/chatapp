package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type room struct {
	forward chan []byte  // message queue
	join    chan *client // queue containing the wishing to join clients
	leave   chan *client // queue containing the wishing to leave clients
	clients map[*client]bool
}

func (r *room) run() {
	for {
		select {
		case client := <-r.join: //request to join message
			r.clients[client] = true
		case client := <-r.leave: // request to leave
			delete(r.clients, client) // delete client
			close(client.send)        // close the send channel
		case msg := <-r.forward:
			for client := range r.clients {
				select {
				case client.send <- msg:
				//send the message
				default:
					//failed to send
					delete(r.clients, client)
					close(client.send)
				}

			}
		}

	}
}

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{

	ReadBufferSize:  socketBufferSize,
	WriteBufferSize: socketBufferSize,
}

func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("ServeHTTP:", err)
		return
	}
	client := &client{
		socket: socket,
		send:   make(chan []byte, messageBufferSize),
		room:   r,
	}
	r.join <- client
	defer func() { r.leave <- client }()
	go client.write()
	client.read()
}

func newRoom() *room {
	return &room{
		forward: make(chan []byte),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
	}
}
