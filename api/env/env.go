package env

import (
	"log"
	"os"

	"github.com/Netflix/go-env"
	"github.com/joho/godotenv"
)

type Environment struct {
	Env           string `env:"ENV,default=development"`
	ApiDomain     string `env:"API_DOMAIN,default=foo.bar.localhost"`
	ApiPort       string `env:"API_PORT,default=:3000"`
	MongoHost     string `env:"MONGO_HOST,default=mongodb://127.0.0.1:27017"`
	MongoUsername string `env:"MONGO_USERNAME,default=:root"`
	MongoPassword string `env:"MONGO_PASSWORD,default=:root"`
}

var environment Environment

func init() {
	currentEnv := os.Getenv("ENV")
	if currentEnv == "development" {
		// loads values from .env into the system
		if err := godotenv.Load(); err != nil {
			log.Fatal("No .env file found")
		}
	}

	_, err := env.UnmarshalFromEnviron(&environment)
	if err != nil {
		log.Fatal(err)
	}
}

func Get() Environment {
	return environment
}
