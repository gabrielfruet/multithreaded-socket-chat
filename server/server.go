package main

import (
    "fmt"
    "net"
    "sync"
)

type Chat struct {
    clients map[string]net.Conn
    sends []chan string
    sends_mutex sync.Mutex
}

func (c *Chat) AddClient(username string, conn net.Conn, send chan string) {
    c.sends_mutex.Lock()
    defer c.sends_mutex.Unlock()
    c.clients[username] = conn
    //c.sends[username] = send
    c.sends = append(c.sends, send)
}

func (c Chat) SendToClients(msg string) {
    c.sends_mutex.Lock()
    defer c.sends_mutex.Unlock()
    for _, sendChan := range c.sends {
        sendChan <- msg
    }
}

func createChat() Chat {
    clients := make(map[string]net.Conn)
    sends := make([]chan string, 0)
    return Chat { clients, sends, sync.Mutex{} }
}


type Client struct {
    chat *Chat
    conn net.Conn
}


func main() {
    ln, err := net.Listen("tcp",":8080")

    if err != nil {
        fmt.Println(err)
        return
    }

    mainchat := createChat()

    for {
        fmt.Println(mainchat.sends)
        conn, err := ln.Accept()

        if err != nil {
            fmt.Println(err)
            continue
        }

        client := Client {&mainchat, conn}

        go client.HandleConnection()
    }
}

func (c Client) SendMsg(recv chan string) {
    for {
        msg := <-recv
        _, err := c.conn.Write([]byte(msg))

        if err != nil {
            fmt.Println(err)
            return
        }
    } 
}

func (c Client) ReceiveMsg() {

    for {

        buf := make([]byte, 1024)

        _, err := c.conn.Read(buf)
        if err != nil {
            fmt.Println(err)
            return
        }

        fmt.Printf("Received: %s\n", buf)
        go c.chat.SendToClients(string(buf))
    }
    
}

func (c Client) HandleConnection() {
    conn := c.conn

    defer conn.Close()

    send := make(chan string)

    username := make([]byte, 128)

    _, err := conn.Read(username)

    if err != nil {
        fmt.Println(err)
        return
    }

    c.chat.AddClient(string(username), conn, send)
    fmt.Printf("%s logged in\n", username)

    go c.ReceiveMsg()
    c.SendMsg(send)
}

