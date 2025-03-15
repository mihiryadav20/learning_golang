package chat

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/gofiber/websocket/v2"
)

// ChatHub manages WebSocket connections and message broadcasting
type ChatHub struct {
	// Registered clients
	clients map[*websocket.Conn]int64 // Connection -> UserID

	// Mutex for thread-safe operations on the clients map
	clientsMu sync.RWMutex

	// Channels for communication between goroutines
	broadcast  chan *ChatMessage
	register   chan *ClientRegistration
	unregister chan *websocket.Conn

	// User repository for database operations
	userRepo UserRepository
}

// ChatMessage represents a message sent in the chat
type ChatMessage struct {
	Type      string    `json:"type"` // "message", "user_joined", "user_left"
	UserID    int64     `json:"user_id"`
	Username  string    `json:"username"`
	Content   string    `json:"content,omitempty"` // Optional for system messages
	Timestamp time.Time `json:"timestamp"`
}

// ClientRegistration contains data for registering a new client
type ClientRegistration struct {
	Conn   *websocket.Conn
	UserID int64
}

// NewChatHub creates a new chat hub
func NewChatHub(userRepo UserRepository) *ChatHub {
	return &ChatHub{
		clients:    make(map[*websocket.Conn]int64),
		broadcast:  make(chan *ChatMessage, 256), // Buffered channel
		register:   make(chan *ClientRegistration, 10),
		unregister: make(chan *websocket.Conn, 10),
		userRepo:   userRepo,
	}
}

// Run starts the ChatHub's main loop in a goroutine
func (h *ChatHub) Run() {
	// This is the main goroutine that handles all operations on the shared state
	go func() {
		for {
			select {
			case registration := <-h.register:
				h.registerClient(registration)

			case conn := <-h.unregister:
				h.unregisterClient(conn)

			case message := <-h.broadcast:
				h.broadcastMessage(message)
			}
		}
	}()
}

// registerClient adds a new client to the hub
func (h *ChatHub) registerClient(registration *ClientRegistration) {
	// Get user information from database
	user, err := h.userRepo.GetUserByID(registration.UserID)
	if err != nil {
		log.Printf("Error getting user %d: %v", registration.UserID, err)
		return
	}

	// Update user status to online
	if err := h.userRepo.UpdateUserStatus(user.ID, "online"); err != nil {
		log.Printf("Error updating user status: %v", err)
	}

	// Register the connection
	h.clientsMu.Lock()
	h.clients[registration.Conn] = registration.UserID
	h.clientsMu.Unlock()

	// Broadcast user joined message
	h.broadcast <- &ChatMessage{
		Type:      "user_joined",
		UserID:    user.ID,
		Username:  user.Username,
		Timestamp: time.Now(),
		Content:   "joined the chat",
	}

	// Send current online users to the new client
	h.sendOnlineUsers(registration.Conn)
}

// unregisterClient removes a client from the hub
func (h *ChatHub) unregisterClient(conn *websocket.Conn) {
	h.clientsMu.Lock()
	userID, exists := h.clients[conn]
	if exists {
		delete(h.clients, conn)
	}
	h.clientsMu.Unlock()

	if !exists {
		return
	}

	// Get user information
	user, err := h.userRepo.GetUserByID(userID)
	if err != nil {
		log.Printf("Error getting user %d: %v", userID, err)
		return
	}

	// Update user status to offline
	if err := h.userRepo.UpdateUserStatus(user.ID, "offline"); err != nil {
		log.Printf("Error updating user status: %v", err)
	}

	// Broadcast user left message
	h.broadcast <- &ChatMessage{
		Type:      "user_left",
		UserID:    user.ID,
		Username:  user.Username,
		Timestamp: time.Now(),
		Content:   "left the chat",
	}
}

// broadcastMessage sends a message to all connected clients
func (h *ChatHub) broadcastMessage(message *ChatMessage) {
	// Marshal the message to JSON
	jsonMessage, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling message: %v", err)
		return
	}

	// Make a copy of the clients map to avoid holding the lock while sending messages
	h.clientsMu.RLock()
	clients := make(map[*websocket.Conn]int64, len(h.clients))
	for conn, userID := range h.clients {
		clients[conn] = userID
	}
	h.clientsMu.RUnlock()

	// Send message to all clients
	for conn := range clients {
		// Use a goroutine to send each message concurrently
		go func(c *websocket.Conn) {
			if err := c.WriteMessage(websocket.TextMessage, jsonMessage); err != nil {
				log.Printf("Error sending message to client: %v", err)
				h.unregister <- c
			}
		}(conn)
	}
}

// sendOnlineUsers sends a list of currently online users to a specific client
func (h *ChatHub) sendOnlineUsers(conn *websocket.Conn) {
	// Get all user IDs currently connected
	h.clientsMu.RLock()
	userIDs := make([]int64, 0, len(h.clients))
	for _, userID := range h.clients {
		// Avoid duplicates
		found := false
		for _, id := range userIDs {
			if id == userID {
				found = true
				break
			}
		}
		if !found {
			userIDs = append(userIDs, userID)
		}
	}
	h.clientsMu.RUnlock()

	// Get user details from the database
	type OnlineUser struct {
		ID       int64  `json:"id"`
		Username string `json:"username"`
	}

	onlineUsers := make([]OnlineUser, 0, len(userIDs))
	for _, userID := range userIDs {
		user, err := h.userRepo.GetUserByID(userID)
		if err != nil {
			log.Printf("Error getting user %d: %v", userID, err)
			continue
		}
		onlineUsers = append(onlineUsers, OnlineUser{
			ID:       user.ID,
			Username: user.Username,
		})
	}

	// Create and send online users message
	message := struct {
		Type        string       `json:"type"`
		OnlineUsers []OnlineUser `json:"online_users"`
		Timestamp   time.Time    `json:"timestamp"`
	}{
		Type:        "online_users",
		OnlineUsers: onlineUsers,
		Timestamp:   time.Now(),
	}

	jsonMessage, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling online users message: %v", err)
		return
	}

	if err := conn.WriteMessage(websocket.TextMessage, jsonMessage); err != nil {
		log.Printf("Error sending online users to client: %v", err)
		h.unregister <- conn
	}
}

// HandleWebSocket handles a WebSocket connection
func (h *ChatHub) HandleWebSocket(c *websocket.Conn, userID int64) {
	// Register the client
	h.register <- &ClientRegistration{
		Conn:   c,
		UserID: userID,
	}

	// Unregister client when the function returns
	defer func() {
		h.unregister <- c
		c.Close()
	}()

	// Handle incoming messages
	for {
		messageType, data, err := c.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v", err)
			break
		}

		// Handle text messages
		if messageType == websocket.TextMessage {
			var message struct {
				Content string `json:"content"`
			}

			if err := json.Unmarshal(data, &message); err != nil {
				log.Printf("Error unmarshaling message: %v", err)
				continue
			}

			// Get user information
			user, err := h.userRepo.GetUserByID(userID)
			if err != nil {
				log.Printf("Error getting user %d: %v", userID, err)
				continue
			}

			// Broadcast the message
			h.broadcast <- &ChatMessage{
				Type:      "message",
				UserID:    userID,
				Username:  user.Username,
				Content:   message.Content,
				Timestamp: time.Now(),
			}
		}
	}
}
