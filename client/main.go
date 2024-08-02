package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

const (
	CONNECTION_SUCCESFULL   = "OK"
	CONNECTION_UNSUCCESFULL = "ERR"
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
			os.Exit(0)
		}

		_, err := conn.Write([]byte(text))

		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

func main() {
	err := godotenv.Load()

	if err != nil {
		fmt.Println("Error loading .env file", err)
	}

	SERVER_IPADDR := os.Getenv("SERVER_IPADDR")

	conn, err := net.Dial("tcp", SERVER_IPADDR+":8080")
	defer conn.Close()

	if err != nil {
		fmt.Println(err)
		return
	}

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Enter your room number (1-100)")
	for scanner.Scan() {
		chatNumber := scanner.Text()
		v, err := strconv.Atoi(chatNumber)
		if err != nil {
			fmt.Println("Invalid room number, enter again: ")
			continue
		}
		if v > 100 || v < 0 {
			fmt.Println("Invalid room number, enter again: ")
			continue
		}
		_, err = conn.Write([]byte(chatNumber))

		if err != nil {
			fmt.Println(err)
			return
		}

		buf := make([]byte, 4)
		n, err := conn.Read(buf)

		if err != nil {
			fmt.Println(err)
			return
		}

		if string(buf[:n]) == CONNECTION_SUCCESFULL {
			break
		} else {
			fmt.Println("Invalid room number, enter again: ")
		}
	}

	fmt.Println("Enter your username")

	for scanner.Scan() {
		username := scanner.Text()
		_, err = conn.Write([]byte(username))

		if err != nil {
			fmt.Println(err)
			return
		}

		buf := make([]byte, 4)
		n, err := conn.Read(buf)

		if err != nil {
			fmt.Println(err)
			return
		}

		if string(buf[:n]) == CONNECTION_SUCCESFULL {
			break
		} else {
			fmt.Println("Username already existed, enter again: ")
		}
	}

	go msgSender(conn)
	msgReceiver(conn)
}
