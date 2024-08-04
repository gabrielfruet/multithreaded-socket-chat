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
	ln, err := net.Listen("tcp", ":8080")
	fmt.Println("Listening on port 8080...")

	if err != nil {
		fmt.Println(err)
		return
	}

    rooms := createRooms(101)

	for {
		conn, err := ln.Accept()

		if err != nil {
			fmt.Println(err)
			conn.Close()
			continue
		}
        go handleConnection(conn, &rooms)
	}
}

func getRoomNumber(conn net.Conn) int {
    for {
        var chatNumberInt int
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

            fmt.Println("Was not a number... " + err.Error())
        } else if v > 100 || v < 0 {
            conn.Write([]byte(CONNECTION_UNSUCCESFULL))

            fmt.Println("Room number not in range.")
        } else {
            chatNumberInt = v
            conn.Write([]byte(CONNECTION_SUCCESFULL))
            return chatNumberInt
        }
    }
    return  -1
}

func handleConnection(conn net.Conn, rooms *Rooms) {
    roomNumber := getRoomNumber(conn)

    for {
        username := make([]byte, 128)

        _, err := conn.Read(username)

        if err != nil {
            fmt.Println(err)
            conn.Close()
            break
        }

        client := createClient(
            rooms.At(roomNumber),
            conn,
            string(username))

        err = client.chat.AddClient(&client)

        if err != nil {
            conn.Write([]byte(CONNECTION_UNSUCCESFULL))

            fmt.Println(err)
        } else {
            conn.Write([]byte(CONNECTION_SUCCESFULL))

            fmt.Printf("%s logged in room %d\n", username, roomNumber)

            go client.ReceiveMsgFromClient()
            go client.SendMsgToClient()

            break
        }
    }
}
