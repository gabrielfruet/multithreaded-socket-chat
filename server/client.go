package main

import (
	"fmt"
	"net"
    "sync"
)

type Client struct {
	chat     *Chat
	conn     net.Conn
	username string
	send     chan Message
    disconnect_mutex sync.Mutex
    connected bool
}

func createClient(chat *Chat, conn net.Conn, username string) Client {
    return Client{
        chat,
        conn,
        string(username),
        make(chan Message),
        sync.Mutex{},
        true,
    }
}

func (c *Client) SendMsgToClient() {
	for {
		msg := <-c.send
		_, err := c.conn.Write([]byte(msg.String()))

		if err != nil {
			fmt.Println(err)
			c.Disconnect()
			return
		}
	}
}

func (c *Client) Disconnect() {
    c.disconnect_mutex.Lock()
    defer c.disconnect_mutex.Unlock()

    if !c.connected { return }

    c.conn.Close()
    c.connected = false
    c.chat.SendToClients(Message{c.username, "Disconnected..."})
    fmt.Printf("%s disconnected.", c.username)
    c.chat.RemoveClient(c.username)
}

func (c *Client) ReceiveMsgFromClient() {
	for {

		buf := make([]byte, 1024)

		_, err := c.conn.Read(buf)

		if err != nil {
			fmt.Println(err)
			c.Disconnect()
			return
		}

		msg := Message{c.username, string(buf)}

		fmt.Printf("Received: %s\n", buf)
		go c.chat.SendToClients(msg)
	}
}
