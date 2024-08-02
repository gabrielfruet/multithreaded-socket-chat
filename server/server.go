package main

import (
	"errors"
	"fmt"
	"net"
	"sync"
)

const (
    CONNECTION_SUCCESFULL= "OK"
    CONNECTION_UNSUCCESFULL="ERR"
)

type Message struct {
    username, text string
}

func (m Message) String() string {
    return fmt.Sprintf("%s: %s", m.username, m.text)
}

type Chat struct {
    clients map[string]Client
    clients_mutex sync.Mutex
}

func (c *Chat) RemoveClient(username string) {
    c.clients_mutex.Lock()
    delete(c.clients, username)
    c.clients_mutex.Unlock()
}

func (c *Chat) AddClient(client Client) error {
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
            conn.Close()
            continue
        }


        for {

            username := make([]byte, 128)

            _, err = conn.Read(username)

            if err != nil {
                fmt.Println(err)
                conn.Close()
                break
            }

            client := Client {
                &mainchat,
                conn,
                string(username),
                make(chan Message),
            }

            err = client.chat.AddClient(client)

            if err != nil {
                conn.Write([]byte(CONNECTION_UNSUCCESFULL))

                fmt.Println(err)
            } else {
                conn.Write([]byte(CONNECTION_SUCCESFULL))

                fmt.Printf("%s logged in\n", username)

                go client.ReceiveMsgFromClient()
                go client.SendMsgToClient()

                break
            }
        }



    }
}

func (c Client) SendMsgToClient() {
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

func (c Client) Disconnect() {
    c.conn.Close()
    c.chat.RemoveClient(c.username)
}

func (c Client) ReceiveMsgFromClient() {

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
