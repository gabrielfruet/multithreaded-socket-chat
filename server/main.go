package main

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

const (
	CONNECTION_SUCCESFULL   = "OK"
	CONNECTION_UNSUCCESFULL = "ERR"
)

func main() {
	ln, err := net.Listen("tcp", ":5555")
	fmt.Println("Listening on port 5555...")

	if err != nil {
		fmt.Println(err)
		return
	}

	var chatList [101]Chat
	for i := 0; i < 101; i++ {
		chatList[i] = createChat()
	}

	for {
		conn, err := ln.Accept()

		if err != nil {
			fmt.Println(err)
			conn.Close()
			continue
		}

        go handleConnection(conn, &chatList)
	}
}

func handleConnection(conn net.Conn, chatList *[101]Chat) {
    for i:=0; i<5; i++{
        var chatNumberInt int
        for i:=0; i<5 ;i++{
            chatNumber := make([]byte, 128)

            _, err := conn.Read(chatNumber)
            if err != nil {
                fmt.Println(err)
                conn.Close()
                break
            }
            chatNumberStr := string(chatNumber)
            chatNumberStr = strings.TrimRight(chatNumberStr, "\x00")
            fmt.Println(chatNumberStr)
            v, err := strconv.Atoi(chatNumberStr)
            if err != nil {
                conn.Write([]byte(CONNECTION_UNSUCCESFULL))

                fmt.Println(err)
            } else if v > 100 || v < 0 {
                conn.Write([]byte(CONNECTION_UNSUCCESFULL))

                fmt.Println("Room number not in range.")
            } else {
                chatNumberInt = v
                break
            }
        }
        conn.Write([]byte(CONNECTION_SUCCESFULL))

        username := make([]byte, 128)

        _, err := conn.Read(username)

        if err != nil {
            fmt.Println(err)
            conn.Close()
            break
        }

        client := createClient(&chatList[chatNumberInt],conn,string(username))

        err = client.chat.AddClient(&client)

        if err != nil {
            conn.Write([]byte(CONNECTION_UNSUCCESFULL))

            fmt.Println(err)
        } else {
            conn.Write([]byte(CONNECTION_SUCCESFULL))

            fmt.Printf("%s logged in room %d\n", username, chatNumberInt)

            go client.ReceiveMsgFromClient()
            go client.SendMsgToClient()
            return
        }
    }
}
