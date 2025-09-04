package manager

import (
    "golang.org/x/sys/unix"
    "sync"
)

type Client struct {
    FD int // file descriptor from Accept
}

type ClientManager struct {
    clients map[int]*Client
    mu      sync.Mutex
}

// Create new manager
func NewClientManager() *ClientManager {
    return &ClientManager{
        clients: make(map[int]*Client),
    }
}

// Add client
func (cm *ClientManager) Add(fd int) {
    cm.mu.Lock()
    defer cm.mu.Unlock()
    cm.clients[fd] = &Client{FD: fd}
}

// Remove client
func (cm *ClientManager) Remove(fd int) {
    cm.mu.Lock()
    defer cm.mu.Unlock()
    delete(cm.clients, fd)
    unix.Close(fd) // close socket
}

// Broadcast message to all clients
func (cm *ClientManager) Broadcast(data []byte) {
    cm.mu.Lock()
    defer cm.mu.Unlock()
    for _, c := range cm.clients {
        unix.Write(c.FD, data)
    }
}

// List connected clients
func (cm *ClientManager) List() []*Client {
    cm.mu.Lock()
    defer cm.mu.Unlock()
    out := make([]*Client, 0, len(cm.clients))
    for _, c := range cm.clients {
        out = append(out, c)
    }
    return out
}