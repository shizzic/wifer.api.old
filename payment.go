package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	net "net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PayPal struct {
	Token       string `json:"access_token"`
	Last_Update string `json:"update_time"`
	Status      string `json:"status"`
}

func CheckPayment(orderID string, c gin.Context) (int64, error) {
	_, err := DB["payments"].InsertOne(ctx, bson.D{{Key: "_id", Value: orderID}})

	if err != nil {
		return 0, errors.New("0")
	}

	token, err := getPayPalToken()

	if err != nil {
		return 0, errors.New("0")
	}

	isOrderCompleted := checkPayPalOrder(orderID, token)

	if !isOrderCompleted {
		return 0, errors.New("0")
	}

	premium := updatePremium(c)

	return premium, nil
}

func updatePremium(c gin.Context) int64 {
	id, _ := c.Cookie("id")
	idInt, _ := strconv.Atoi(id)

	var user bson.M
	opts := options.FindOne().SetProjection(bson.M{"_id": 0, "premium": 1})
	DB["users"].FindOne(ctx, bson.M{"_id": idInt}, opts).Decode(&user)

	premium := user["premium"].(int64)
	now := time.Now().Unix()

	if premium > 0 && premium > now {
		expires := int64(premium + int64(1*60*60*24*30))

		DB["users"].UpdateOne(ctx, bson.M{"_id": idInt}, bson.D{
			{Key: "$set", Value: bson.D{{Key: "premium", Value: expires}}},
		})

		c.SetSameSite(net.SameSiteNoneMode)
		c.SetCookie("premium", "premium", int(premium-now+(1*60*60*24*30)), "/", "."+domainBack, true, true)

		return expires
	} else {
		expires := int64(now + int64(1*60*60*24*30))

		DB["users"].UpdateOne(ctx, bson.M{"_id": idInt}, bson.D{
			{Key: "$set", Value: bson.D{{Key: "premium", Value: expires}}},
		})

		c.SetSameSite(net.SameSiteNoneMode)
		c.SetCookie("premium", "premium", 1*60*60*24*30, "/", "."+domainBack, true, true)

		return expires
	}
}

func getPayPalToken() (string, error) {
	client := &http.Client{}
	data := strings.NewReader(`grant_type=client_credentials`)
	req, err := http.NewRequest("POST", "https://api-m.paypal.com/v1/oauth2/token", data)

	if err != nil {
		return "", errors.New("0")
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth("AQWOBzexECxYSIUxzjRH1YplNI_bxuc6Vzxa2uRWltHZzgOAjR23hTbEPXVm82UkF8Ebvj-QSQHo3pbH", "EGqQToUYR77NM_fZQaNPKTFJ8YDtFCsgraxE9GlOFlf3YI653tUsvSJ27yZGPlGXFgEqnAfyjiM34QlM")
	resp, err := client.Do(req)

	if err != nil {
		return "", errors.New("0")
	}

	defer resp.Body.Close()
	bodyText, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return "", errors.New("0")
	}

	paypal := PayPal{}
	json.Unmarshal([]byte(bodyText), &paypal)

	return paypal.Token, nil
}

func checkPayPalOrder(orderID, token string) bool {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.paypal.com/v2/checkout/orders/"+orderID, nil)

	if err != nil {
		return false
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := client.Do(req)

	if err != nil {
		return false
	}

	defer resp.Body.Close()
	bodyText, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return false
	}

	paypal := PayPal{}
	json.Unmarshal([]byte(bodyText), &paypal)

	if paypal.Status != "COMPLETED" {
		return false
	}

	if !checkPayPalOrderTime(paypal.Last_Update) {
		return false
	}

	return true
}

func checkPayPalOrderTime(last_update string) bool {
	date, _ := time.Parse(time.RFC3339, last_update)
	update := date.Unix()
	minus := time.Now().Unix() - update

	if minus > 3600 {
		return false
	}

	return true
}
