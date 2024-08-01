package main

import (
    "os"
    "bufio"
    "fmt"
    "net"
)

func msgReceiver(conn net.Conn) {
    for {
        buf := make([]byte, 1024)

        _, err := conn.Read(buf)
        if err != nil {
            fmt.Println(err)
            return
        }

        fmt.Printf("%s\n", buf)
    } 
}

func msgSender(conn net.Conn) {
    scanner := bufio.NewScanner(os.Stdin)
    for scanner.Scan() {
        text := scanner.Text()

        if text == "/disconnect" {
            fmt.Println("Disconnecting.")
            return
        }

        _, err := conn.Write([]byte(text))

        if err != nil {
            fmt.Println(err)
            return
        }
    }
}

func main() {
    conn, err := net.Dial("tcp", "localhost:8080")
    defer conn.Close()

    if err != nil {
        fmt.Println(err)
        return
    }

    scanner := bufio.NewScanner(os.Stdin)

    fmt.Println("Enter your username")

    scanner.Scan()
    username := scanner.Text()
    _, err = conn.Write([]byte(username))

    if err != nil {
        fmt.Println(err)
        return
    }

    go msgSender(conn)
    msgReceiver(conn)
}
