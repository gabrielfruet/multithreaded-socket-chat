package main

import (
    "fmt"
    "net"
    "sync"
)

type Message struct {
    username, text string
}

func (m Message) String() string {
    return fmt.Sprintf("%s: %s", m.username, m.text)
}

type Chat struct {
    clients map[string]Client
    sends_mutex sync.Mutex
}

func (c *Chat) AddClient(client Client) {
    c.sends_mutex.Lock()
    defer c.sends_mutex.Unlock()
    c.clients[client.username] = client
}

func (c *Chat) SendToClients(msg Message) {
    c.sends_mutex.Lock()
    defer c.sends_mutex.Unlock()
    for _, client := range c.clients {
        client.send <- msg
    }
}

func createChat() Chat {
    clients := make(map[string]Client)
    return Chat { clients, sync.Mutex{} }
}


type Client struct {
    chat *Chat
    conn net.Conn
    username string
    send chan Message
}


func main() {
    ln, err := net.Listen("tcp",":8080")

    if err != nil {
        fmt.Println(err)
        return
    }

    mainchat := createChat()

    for {
        conn, err := ln.Accept()

        if err != nil {
            fmt.Println(err)
            continue
        }


        username := make([]byte, 128)

        _, err = conn.Read(username)

        if err != nil {
            fmt.Println(err)
            return
        }

        send := make(chan Message)

        client := Client {&mainchat, conn, string(username), send}

        client.chat.AddClient(client)

        fmt.Printf("%s logged in\n", username)

        go client.ReceiveMsgFromClient()
        go client.SendMsgToClient()
    }
}

func (c Client) SendMsgToClient() {
    for {
        msg := <-c.send
        _, err := c.conn.Write([]byte(msg.String()))

        if err != nil {
            fmt.Println(err)
            return
        }
    } 
}

func (c Client) ReceiveMsgFromClient() {

    for {

        buf := make([]byte, 1024)

        _, err := c.conn.Read(buf)
        if err != nil {
            fmt.Println(err)
            return
        }

        msg := Message{c.username, string(buf)}

        fmt.Printf("Received: %s\n", buf)
        go c.chat.SendToClients(msg)
    }
}
