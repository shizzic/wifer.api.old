package main

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type messages struct {
	Target int   `form:"target"`
	Skip   int64 `form:"skip"`
	Access bool  `form:"access"`
}

type rooms struct {
	Nin []int `form:"nin"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var clients = make(map[int]*websocket.Conn)

func Chat(w http.ResponseWriter, r *http.Request, c gin.Context) {
	idCookie, _ := c.Cookie("id")
	id, _ := strconv.Atoi(idCookie)

	con, _ := upgrader.Upgrade(w, r, nil)
	defer con.Close() // Закрываем соединение
	if _, exist := clients[id]; exist {
		clients[id].Close()
		delete(clients, id)
	}
	clients[id] = con
	defer quit(id)

	for {
		var msg struct {
			Api        string `json:"api"`
			Text       string `json:"text"`
			Username   string `json:"username"`
			User       int    `json:"user"`
			Target     int    `json:"target"`
			Access     bool   `json:"access"`
			Typing     bool   `json:"typing"`
			Avatar     bool   `json:"avatar"`
			Created_at int64  `json:"created_at"`
		}

		err := con.ReadJSON(&msg)
		if err != nil {
			break
		}

		if _, exist := clients[msg.Target]; exist {
			clients[msg.Target].WriteJSON(msg)
		}

		switch msg.Api {
		case "message":
			writeMessage(msg.Text, msg.User, msg.Target, msg.Created_at)
			break
		case "view":
			viewMessages(msg.User, msg.Target)
			break
		case "access":
			if msg.Access {
				AddAccess(msg.Target, c)
			} else {
				DeleteAccess(msg.Target, c)
			}
			break
		}
	}
}

func GetRooms(data rooms, c gin.Context) (map[string][]bson.M, []int) {
	var res = make(map[string][]bson.M)
	id, _ := c.Cookie("id")
	idInt, _ := strconv.Atoi(id)
	var rooms []bson.M
	var ids = []int{}
	var users []bson.M

	matchFilter := bson.D{{Key: "$match",
		Value: bson.D{{
			Key: "roommates", Value: bson.D{
				{Key: "$in", Value: []int{idInt}},
				{Key: "$nin", Value: data.Nin},
			},
		}},
	}}

	groupFilter := bson.D{{Key: "$group",
		Value: bson.D{
			{Key: "_id", Value: "$roommates"},
			{Key: "user", Value: bson.D{{Key: "$first", Value: "$user"}}},
			{Key: "target", Value: bson.D{{Key: "$first", Value: "$target"}}},
			{Key: "text", Value: bson.D{{Key: "$first", Value: "$text"}}},
			{Key: "created_at", Value: bson.D{{Key: "$first", Value: "$created_at"}}},
			{Key: "viewed", Value: bson.D{{Key: "$first", Value: "$viewed"}}},
		},
	}}

	sortFilter := bson.D{{Key: "$sort", Value: bson.D{{Key: "created_at", Value: -1}}}}
	limitFilter := bson.D{{Key: "$limit", Value: 25}}

	cursor, _ := DB["messages"].Aggregate(ctx, mongo.Pipeline{matchFilter, sortFilter, groupFilter, sortFilter, limitFilter})
	cursor.All(ctx, &rooms)
	res["rooms"] = rooms

	for _, v := range rooms {
		user := int(v["user"].(int32))

		if user != idInt {
			ids = append(ids, user)
			continue
		}

		ids = append(ids, int(v["target"].(int32)))
	}

	if len(ids) > 0 {
		opts := options.Find().SetProjection(bson.M{"username": 1, "avatar": 1, "online": 1})
		cur, _ := DB["users"].Find(ctx, bson.M{"_id": bson.M{"$in": ids}, "status": true}, opts)
		cur.All(ctx, &users)
	}

	res["users"] = users

	return res, ids
}

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

func writeMessage(text string, user, target int, created_at int64) {
	trimmed := strings.TrimSpace(text)
	var roommates []int

	if user > target {
		roommates = []int{target, user}
	} else {
		roommates = []int{user, target}
	}

	DB["messages"].InsertOne(ctx, bson.D{
		{Key: "roommates", Value: roommates},
		{Key: "user", Value: user},
		{Key: "target", Value: target},
		{Key: "viewed", Value: false},
		{Key: "text", Value: trimmed},
		{Key: "created_at", Value: created_at},
	})
}

func viewMessages(user, target int) {
	DB["messages"].UpdateMany(ctx, bson.M{"user": target, "target": user}, bson.D{
		{Key: "$set", Value: bson.D{{Key: "viewed", Value: true}}},
	})
}

// User left socket for whatever reason
func quit(id int) {
	delete(clients, id) // Удаляем соединение
}
