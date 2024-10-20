package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Himanshu303/mysql-mongo-migration/config"
	"github.com/Himanshu303/mysql-mongo-migration/models"
	_ "github.com/go-sql-driver/mysql"
	"github.com/streadway/amqp"
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

	_, err = ch.QueueDeclare("mysql-mongo_migration_queue", true, false, false, false, nil)

	if err != nil {
		log.Fatal("Failed to declare a queue:", err)
	}

	db, err := sql.Open("mysql", cfg.SqlDbURI)

	if err != nil {
		log.Fatal("Failed to connect to MYSQL", err)

	}

	defer db.Close()

	rows, err := db.Query("SELECT ID, name, bdate, marks, gpa  from students")

	if err != nil {
		log.Fatal("Failed to query data:", err)

	}

	defer rows.Close()

	for rows.Next() {
		var student models.Student
		var bdate []uint8

		if err := rows.Scan(&student.ID, &student.Name, &bdate, &student.Marks, &student.Gpa); err != nil {
			log.Fatal("Failed to scan row:", err)
		}

		if bdate != nil {
			t, err := time.Parse("2006-01-02", string(bdate))
			if err != nil {
				log.Fatal("failed to parse bdate: %w", err)
			}
			student.Bdate = &t
		} else {
			student.Bdate = nil
		}

		body, _ := json.Marshal(student)

		err = ch.Publish("", "mysql-mongo_migration_queue", false, false, amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})

		if err != nil {
			log.Fatal("Failed to publish message:", err)
		}
		fmt.Println("Published:", student)
	}

}
