package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port       string
	MongoDBUri string
	Database   string
}

var config *Config

func EnvSetup() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	config = &Config{
		Port:       getEnv("PORT"),
		MongoDBUri: getEnv("MONGODB_URI"),
		Database:   getEnv("DB_NAME"),
	}
}

func GetConfig() Config {
	if config == nil {
		EnvSetup()
	}
	return *config
}

func getEnv(v string) string {
	v = os.Getenv(v)
	if v == "" {
		log.Fatalf("%s not found !!", v)
	}

	return v
}
