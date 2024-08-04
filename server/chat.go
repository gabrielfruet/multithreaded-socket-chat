package main

import (
	"errors"
	"sync"
)

type Chat struct {
	clients       map[string]*Client
	clients_mutex sync.Mutex
}

func (c *Chat) RemoveClient(username string) {
	c.clients_mutex.Lock()
	defer c.clients_mutex.Unlock()
	delete(c.clients, username)
}

func (c *Chat) AddClient(client *Client) error {
	c.clients_mutex.Lock()
	defer c.clients_mutex.Unlock()

	_, exists := c.clients[client.username]

	if exists {
		return errors.New("Client already existed.")
	}

	c.clients[client.username] = client

	return nil
}

func (c *Chat) SendToClients(msg Message) {
	c.clients_mutex.Lock()
	defer c.clients_mutex.Unlock()
	for username, client := range c.clients {
		if msg.username != username {
			client.send <- msg
		}
	}
}

func createChat() Chat {
	clients := make(map[string]*Client)
	return Chat{clients, sync.Mutex{}}
}
