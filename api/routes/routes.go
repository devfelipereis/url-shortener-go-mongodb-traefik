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
	"go.mongodb.org/mongo-driver/mongo/options"
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

var linksCollection *mongo.Collection

func init() {
	linksCollection = database.Connect().Database("mydatabase").Collection(COLLECTION)
	indexModel := mongo.IndexModel{
		Keys: bson.M{
			"code": 1, // index in ascending order
		}, Options: options.Index().SetUnique(true),
	}

	if _, err := linksCollection.Indexes().CreateOne(context.TODO(), indexModel); err != nil {
		log.Fatal(err)
	}
}

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
	maxTries := 1000
	found := false

	for {
		if maxTries == 0 {
			break
		}
		var doc LinkSchema
		err := linksCollection.FindOne(context.TODO(), bson.M{"code": code}).Decode(&doc)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				found = true
				break
			}
			log.Fatal(err)
		}
		maxTries--
		code = randomCode()
	}

	if !found {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Impossible to create more links"})
		return
	}

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
