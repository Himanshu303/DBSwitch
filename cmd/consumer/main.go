package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/Himanshu303/mysql-mongo-migration/config"
	"github.com/Himanshu303/mysql-mongo-migration/models"
	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	cfg := config.LoadConfig()

	conn, err := amqp.Dial(cfg.RabbitMQUrl)

	if err != nil {
		log.Fatal("Failed to connect to rabbitmq", err)
	}

	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal("Failed to open a channel:", err)
	}
	defer ch.Close()

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(cfg.MongoURI))

	if err != nil {
		log.Fatal("Failed to connect to mongodb:", err)
	}

	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	coll := client.Database(cfg.DBName).Collection("students")

	msgs, err := ch.Consume("mysql-mongo_migration_queue", "", true, false, false, false, nil)

	if err != nil {
		log.Fatal("Failed to register a consumer:", err)
	}

	forever := make(chan bool)

	go func() {
		for msg := range msgs {
			var student models.Student
			if err := json.Unmarshal(msg.Body, &student); err != nil {
				log.Println("Failed to unmarshal message:", err)
				continue
			}

			_, err := coll.InsertOne(context.TODO(), student)
			if err != nil {
				log.Println("Failed to insert record into MongoDB:", err)
				continue
			}

			fmt.Println("Migrated Record:", student)
		}
	}()

	fmt.Println("Waiting for messages. Press Ctrl+C to exit.")
	<-forever
}
