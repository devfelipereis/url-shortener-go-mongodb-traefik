package routes

import (
	"context"
	"devfelipereis/urlShortener/database"
	"devfelipereis/urlShortener/env"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	COLLECTION = "links"
)

type LinkSchema struct {
	ID        primitive.ObjectID `bson:"_id"`
	Code      string             `bson:"code"`
	Url       string             `bson:"url"`
	CreatedAt time.Time          `bson:"createdAt"`
}

var linksCollection = database.Connect().Database("mydatabase").Collection(COLLECTION) // get collection "links" from mongo.Connect() which returns *mongo.Client

func Home(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Welcome to the api"})
}

func Redirect(c *gin.Context) {
	code := c.Param("code")

	var doc LinkSchema

	err := linksCollection.FindOne(context.TODO(), bson.M{"code": code}).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		log.Fatal(err)
	}

	c.Redirect(http.StatusMovedPermanently, doc.Url)
}

func Generate(c *gin.Context) {
	type Url struct {
		Url string `json:"url" binding:"required"`
	}

	var json Url
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	uri, err := url.ParseRequestURI(json.Url)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	code := randomCode()
	link := &LinkSchema{
		ID:        primitive.NewObjectID(),
		Code:      code,
		Url:       uri.String(),
		CreatedAt: time.Now(),
	}

	if _, err := linksCollection.InsertOne(context.TODO(), link); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"url": getBaseUrl() + code,
	})
}

func GetOne(c *gin.Context) {
	code := c.Param("code")

	var doc bson.M

	err := linksCollection.FindOne(context.TODO(), bson.M{"code": code}).Decode(&doc)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H(doc))
}

func randomCode() string {
	return randSeq(7)
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func getBaseUrl() string {
	appEnv := env.Get()

	apiPort := appEnv.ApiPort
	port := ""

	if apiPort != ":80" && appEnv.Env == "development" {
		port = apiPort
	}

	return "http://" + appEnv.ApiDomain + port + "/"
}
