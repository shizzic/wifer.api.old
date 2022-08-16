package main

import (
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var clients = make(map[int]*websocket.Conn)
var rooms = make(map[int]map[int]bool)

func Chat(w http.ResponseWriter, r *http.Request, id int) {
	con, _ := upgrader.Upgrade(w, r, nil)
	defer con.Close() // Закрываем соединение
	clients[id] = con
	enter(id)
	defer quit(id)

	for {
		var msg struct {
			Message string `json:"message"`
			Api     string `json:"api"`
		}

		err := con.ReadJSON(&msg)
		if err != nil {
			break
		}

		con.WriteJSON(msg)
	}
}

// User enters socket
func enter(id int) {
	rooms[id] = make(map[int]bool)
}

// User left socket for whatever reason
func quit(id int) {
	for client := range rooms[id] {
		if _, exist := rooms[client]; exist {
			delete(rooms[client], id)
		}
	}

	delete(clients, id) // Удаляем соединение
}

func GetRooms(id int) {

}

func GetMessages(id int) {

}
