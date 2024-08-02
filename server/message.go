package main

import "fmt"

type Message struct {
    username, text string
}

func (m Message) String() string {
    return fmt.Sprintf("%s: %s", m.username, m.text)
}
