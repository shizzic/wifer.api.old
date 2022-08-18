package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type messages struct {
	Target int   `form:"target"`
	Skip   int64 `form:"skip"`
	Access bool  `form:"access"`
}

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
	if _, exist := clients[id]; exist {
		clients[id].Close()
		delete(clients, id)
	}
	clients[id] = con
	enter(id)
	defer quit(id)

	for {
		var msg struct {
			Message string `json:"message"`
			Api     string `json:"api"`
			Access  bool   `json:"access"`
			Target  int    `json:"target"`
			User    int    `json:"user"`
		}

		err := con.ReadJSON(&msg)
		if err != nil {
			break
		}

		if _, exist := clients[msg.Target]; exist {
			switch msg.Api {
			case "access":
				clients[msg.Target].WriteJSON(msg)
			}
		}
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

// func GetRooms(c gin.Context) {
// 	id, _ := c.Cookie("id")
// 	idInt, _ := strconv.Atoi(id)
// }

func GetMessages(data messages, c gin.Context) map[string][]bson.M {
	var res = make(map[string][]bson.M)
	id, _ := c.Cookie("id")
	idInt, _ := strconv.Atoi(id)
	access := true
	filter := bson.M{"user": bson.M{"$in": []int{idInt, data.Target}}, "target": bson.M{"$in": []int{idInt, data.Target}}}

	if data.Access {
		if _, err := c.Cookie("premium"); err != nil {
			access = false
		}

		accesses := CheckRoomAccess(idInt, data.Target, filter)
		res["accesses"] = accesses

		if !access {
			for _, v := range accesses {
				if v["target"].(int32) == int32(idInt) {
					access = true
				}
			}
		}
	}

	if access {
		var messages []bson.M
		opts := options.Find().SetProjection(bson.M{
			"_id":        0,
			"user":       1,
			"text":       1,
			"created_at": 1,
			"viewed":     1,
		}).
			SetSort(bson.D{{Key: "created_at", Value: -1}}).SetLimit(25).SetSkip(data.Skip)

		cursor, _ := DB["messages"].Find(ctx, filter, opts)
		cursor.All(ctx, &messages)

		res["messages"] = messages
	}

	return res
}

func CheckRoomAccess(id, target int, filter primitive.M) []bson.M {
	var accesses []bson.M
	opts := options.Find().SetProjection(bson.M{"_id": 0, "user": 1, "target": 1})
	cursor, _ := DB["access"].Find(ctx, filter, opts)
	cursor.All(ctx, &accesses)
	return accesses
}
