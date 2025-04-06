package main

import (
	"fmt"
	"net/http"

	"clipsync.com/m/db"
	"clipsync.com/m/handlers"
	"clipsync.com/m/ws"
)

func main() {
	db.ConnectDB()
	db.ConnectRedis()
	server := ws.NewServer()
	go server.Run()

	http.HandleFunc("/register", handlers.RegisterHandler)
	http.HandleFunc("/login", handlers.LoginHandler)
	http.HandleFunc("/update-password", handlers.UpdatePasswordHandler)
	http.HandleFunc("/forgot-password", handlers.ForgotPasswordHandler)
	http.HandleFunc("/reset-password", handlers.ResetPasswordHandler)
	http.HandleFunc("/login-client", handlers.LoginClientPage)
	http.HandleFunc("/register-client", handlers.RegisterClientPage)
	http.HandleFunc("/forgot-password-client", handlers.ForgotPasswordClientPage)
	http.HandleFunc("/reset-password-client", handlers.ResetPasswordClientPage)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ws.ServeWS(server, w, r)
	})
	fmt.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
