package main

import (
    "database/sql"
    "fmt"
    "log"

    "github.com/streadway/amqp"
    _ "github.com/lib/pq" // Import for PostgreSQL driver
)

func failOnError(err error, msg string) {
    if err != nil {
        log.Fatalf("%s: %s", msg, err)
        panic(fmt.Sprintf("%s: %s", msg, err))
    }
}

func main() {
    // Connect to RabbitMQ
    conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
    failOnError(err, "Failed to connect to RabbitMQ")
    defer conn.Close()

    ch, err := conn.Channel()
    failOnError(err, "Failed to open a channel")
    defer ch.Close()

    // Connect to PostgreSQL (Replace with your credentials)
    dsn := "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"

    db, err := sql.Open("postgres", dsn)
    failOnError(err, "Failed to connect to PostgreSQL")
    defer db.Close()

    // Declare queue
    q, err := ch.QueueDeclare(
        "hello-queue", // name
        false,   // durable
        false,   // delete when unused
        false,   // exclusive
        false,   // no-wait
        nil,     // arguments
    )
    failOnError(err, "Failed to declare a queue")

    // Consume messages
    msgs, err := ch.Consume(
        q.Name, // queue
        "",     // consumer
        true,   // auto-ack
        false,  // exclusive
        false,  // no-local
        false,  // no-wait
        nil,    // args
    )
    failOnError(err, "Failed to register a consumer")

    forever := make(chan bool)

    go func() {
        for d := range msgs {
            message := string(d.Body)
            log.Printf("Received a message: %s", message)

            // Prepare and execute SQL statement
            stmt, err := db.Prepare("INSERT INTO messages (message) VALUES ($1)")
            failOnError(err, "Failed to prepare statement")

            _, err = stmt.Exec(message)
            failOnError(err, "Failed to save message")

        }
    }()

    log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
    <-forever
}
