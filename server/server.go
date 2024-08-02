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
