package main

type Rooms struct {
    rooms []Chat
}

func createRooms(n int) Rooms {
    rooms := make([]Chat, n)
    for i := 0; i < n; i++ {
        rooms[i] = createChat()
    }
    return Rooms{rooms}
}

func (r Rooms) At(i int) *Chat {
    return &(r.rooms[i])
}
