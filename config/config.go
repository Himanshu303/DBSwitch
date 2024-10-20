package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	RabbitMQUrl string
	SqlDbURI    string
	MongoURI    string
	DBName      string
}

func init() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")

	}
}

func LoadConfig() *Config {
	return &Config{
		RabbitMQUrl: os.Getenv("RABBITMQ_URL"),
		SqlDbURI:    os.Getenv("MYSQL_DSN"),
		MongoURI:    os.Getenv("MONGO_URI"),
		DBName:      os.Getenv("DB_NAME"),
	}
}
