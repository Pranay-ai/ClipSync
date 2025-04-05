

# ClipSync Backend Server

The ClipSync backend server is a real-time communication system that enables cross-device clipboard synchronization using WebSockets, Redis, and a local relational database for user authentication.

## Features

- User registration and login via email/password
- OTP-based password reset flow
- JWT token generation and validation for secure authentication
- WebSocket endpoint for device-to-device clipboard sync
- Redis Pub/Sub to sync messages across multiple server instances
- Local database for managing user accounts

## Tech Stack

- Go (Golang)
- Redis (Pub/Sub)
- PostgreSQL (or any SQL-compatible DB via GORM)
- Gorilla WebSocket
- JWT (github.com/golang-jwt/jwt)

---

## Architecture Overview

1. **User Auth System**
   - Users register and log in using email/password.
   - Passwords are hashed using bcrypt.
   - On login or registration, a JWT is generated and sent to the client.
   - Clients use this JWT to authenticate WebSocket connections.

2. **JWT Authentication**
   - The JWT includes `user_id` and is passed as a query parameter when initiating a WebSocket connection.
   - The backend validates the token before accepting the connection.

3. **WebSocket Communication**
   - Each connected client is registered to the server with a unique `user_id` and `device_id`.
   - Clipboard content sent from one device is relayed to other devices under the same user, excluding the sender.

4. **Redis Pub/Sub**
   - When a client sends a message over WebSocket, the message is published to a Redis channel based on the user's ID.
   - All server instances subscribed to that channel receive the message and broadcast it to their connected clients.
   - This allows WebSocket connections to scale across multiple servers without direct communication between them.

5. **Database**
   - GORM is used to connect to a local SQL database.
   - User information, including credentials and metadata, is stored and managed securely.

---

## Key Endpoints

### POST `/register`
Registers a new user and returns a JWT.

### POST `/login`
Authenticates a user and returns a JWT.

### POST `/forgot-password`
Sends an OTP to the user's email (mocked in logs).

### POST `/reset-password`
Validates the OTP and allows the user to reset their password.

### GET `/login-client`
Serves an HTML page for browser-based login (used by desktop tray apps).

### GET `/ws`
WebSocket upgrade endpoint used by clients to send and receive clipboard sync messages.

---

## How Redis is Used

- The server publishes clipboard messages to a Redis channel named:
  ```
  clipboard_sync:user:<user_id>
  ```
- All server instances subscribe to that channel for their respective users.
- When a message is published, it is re-broadcasted to all connected WebSocket clients (except the one that originated it).

This allows seamless communication across distributed server instances.

---

## Setup

- Requires Redis and a SQL database running locally or in your environment.
- Configure database connection via GORM settings.
- Run the server and connect your client (tray app or any WebSocket-capable app) with a valid JWT.

---

Let me know if you’d like to include `.env` configuration instructions, database schema, or Docker support.