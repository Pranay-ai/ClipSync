package ws

import (
	"log"
)

type Message struct {
	UserID     string
	FromDevice string
	Payload    []byte
}

type Server struct {
	clients    map[string]map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan Message
}

func NewServer() *Server {
	return &Server{
		clients:    make(map[string]map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan Message),
	}
}

func (s *Server) Run() {
	for {
		select {

		case client := <-s.register:
			// First client for this user? Initialize map and subscribe to Redis
			if s.clients[client.UserID] == nil {
				s.clients[client.UserID] = make(map[*Client]bool)
				go SubscribeToUserChannel(client.UserID, s)
			}

			s.clients[client.UserID][client] = true
			log.Printf("Client registered: user %s (%s)", client.UserID, client.DeviceID)

		case client := <-s.unregister:
			if clients, ok := s.clients[client.UserID]; ok {
				if _, exists := clients[client]; exists {
					delete(clients, client)
					close(client.SendChan)
					log.Printf("Client disconnected: user %s (%s)", client.UserID, client.DeviceID)

					if len(clients) == 0 {
						delete(s.clients, client.UserID)
						log.Printf("No more clients for user %s", client.UserID)
					}
				}
			}

		case msg := <-s.broadcast:
			if clients, ok := s.clients[msg.UserID]; ok {
				for c := range clients {
					if c.DeviceID != msg.FromDevice {
						select {
						case c.SendChan <- msg.Payload:
						default:
							log.Println("Send buffer full, closing connection")
							close(c.SendChan)
							delete(clients, c)
						}
					}
				}
			}
		}
	}
}
