package main

import (
	"net/http"
	"strconv"
	"strings"
	"sync"

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
	Nin        []int  `form:"nin"`
	Username   string `form:"username"`
	ByUsername bool   `form:"byUsername"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		return origin == CLIENT_DOMAIN
	},
}

// var clients = make(map[int]*websocket.Conn)
var clients sync.Map

func Chat(w http.ResponseWriter, r *http.Request, c *gin.Context) {
	idCookie, _ := c.Cookie("id")
	id, _ := strconv.Atoi(idCookie)

	con, _ := upgrader.Upgrade(w, r, nil)
	// defer con.Close() // Закрываем соединение
	// if _, exist := clients[id]; exist {
	// 	clients[id].Close()
	// 	delete(clients, id)
	// }
	// clients[id] = con
	// defer quit(id)

	defer con.Close() // Закрываем соединение
	if v, exist := clients.Load(id); exist {
		c := v.(*websocket.Conn)
		clients.Delete(id)
		c.Close()
	}
	clients.Store(id, con)
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

		if v, exist := clients.Load(msg.Target); exist {
			c := v.(*websocket.Conn)
			c.WriteJSON(msg)
		}

		switch msg.Api {
		case "message":
			writeMessage(msg.Text, msg.User, msg.Target, msg.Created_at)
		case "view":
			viewMessages(msg.User, msg.Target)
		case "access":
			if msg.Access {
				AddAccess(msg.Target, c)
			} else {
				DeleteAccess(msg.Target, c)
			}
		}
	}
}

func GetRooms(data rooms, c *gin.Context) (map[string][]bson.M, []int) {
	var res = make(map[string][]bson.M)
	id, _ := c.Cookie("id")
	idInt, _ := strconv.Atoi(id)
	var rooms []bson.M
	var ids = []int{}
	var users []bson.M

	groupFilter := bson.D{{Key: "$group",
		Value: bson.D{
			{Key: "_id", Value: "$roommates"},
			{Key: "user", Value: bson.D{{Key: "$first", Value: "$user"}}},
			{Key: "target", Value: bson.D{{Key: "$first", Value: "$target"}}},
			{Key: "text", Value: bson.D{{Key: "$first", Value: "$text"}}},
			{Key: "created_at", Value: bson.D{{Key: "$first", Value: "$created_at"}}},
			{Key: "viewed", Value: bson.D{{Key: "$first", Value: "$viewed"}}},
			{Key: "news", Value: bson.D{{Key: "$sum", Value: bson.D{{
				Key: "$cond", Value: bson.D{
					{Key: "if", Value: bson.D{
						{Key: "$and", Value: []bson.D{
							{
								{Key: "$eq", Value: []interface{}{"$viewed", false}},
							},
							{
								{Key: "$eq", Value: []interface{}{"$target", idInt}},
							},
						}},
					}},
					{Key: "then", Value: 1},
					{Key: "else", Value: 0},
				},
			}}}}},
		},
	}}

	sortFilter := bson.D{{Key: "$sort", Value: bson.D{{Key: "created_at", Value: -1}}}}
	limitFilter := bson.D{{Key: "$limit", Value: 25}}

	if !data.ByUsername {
		matchFilter := bson.D{{Key: "$match",
			Value: bson.D{{
				Key: "roommates", Value: bson.D{
					{Key: "$in", Value: []int{idInt}},
					{Key: "$nin", Value: data.Nin},
				},
			}},
		}}

		cursor, _ := DB["messages"].Aggregate(ctx, mongo.Pipeline{matchFilter, sortFilter, groupFilter, sortFilter, limitFilter})
		cursor.All(ctx, &rooms)

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
	} else {
		nin := data.Nin
		nin = append(nin, idInt)
		opts := options.Find().SetProjection(bson.M{"username": 1, "avatar": 1, "online": 1})
		cur, _ := DB["users"].Find(ctx, bson.M{"username": bson.M{"$regex": data.Username, "$options": "gi"}, "_id": bson.M{"$nin": nin}, "status": true}, opts)
		cur.All(ctx, &users)

		for _, v := range users {
			ids = append(ids, int(v["_id"].(int32)))
		}

		if len(ids) > 0 {
			freshIds := ids
			freshIds = append(freshIds, idInt)

			matchFilter := bson.D{{Key: "$match",
				Value: bson.D{
					{Key: "user", Value: bson.D{{Key: "$in", Value: freshIds}}},
					{Key: "target", Value: bson.D{{Key: "$in", Value: freshIds}}},
				},
			}}

			cursor, _ := DB["messages"].Aggregate(ctx, mongo.Pipeline{matchFilter, sortFilter, groupFilter, sortFilter, limitFilter})
			cursor.All(ctx, &rooms)
		}
	}

	res["rooms"] = rooms
	res["users"] = users

	return res, ids
}

func GetMessages(data messages, c *gin.Context) map[string][]bson.M {
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

func CheckOnlineInChat(data rooms) []bson.M {
	var users []bson.M
	opts := options.Find().SetProjection(bson.M{"online": 1})
	cur, _ := DB["users"].Find(ctx, bson.M{"_id": bson.M{"$in": data.Nin}, "status": true}, opts)
	cur.All(ctx, &users)
	return users
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
	clients.Delete(id)
}
