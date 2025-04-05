package ws

import (
	"log"
	"net/http"

	"clipsync.com/m/utils" // Assuming your GenerateJWT is here
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins — lock this down in production
	},
}

func ServeWS(server *Server, w http.ResponseWriter, r *http.Request) {
	// 🔐 Step 1: Get token and device_id from query params
	tokenStr := r.URL.Query().Get("token")
	deviceID := r.URL.Query().Get("device_id")

	if tokenStr == "" || deviceID == "" {
		http.Error(w, "Missing token or device_id", http.StatusBadRequest)
		return
	}

	// 🔐 Step 2: Validate JWT and extract user_id
	userID, err := utils.ValidateJWT(tokenStr)
	if err != nil {
		http.Error(w, "Invalid or expired token: "+err.Error(), http.StatusUnauthorized)
		return
	}

	// 🔗 Step 3: Upgrade to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}

	// ✅ Step 4: Create and register client
	client := &Client{
		UserID:   userID,
		DeviceID: deviceID,
		Conn:     conn,
		SendChan: make(chan []byte, 256),
		Server:   server,
	}

	client.Server.register <- client

	go client.WritePump()
	go client.ReadPump()
}
