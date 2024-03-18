package main

import (
    "fmt"
    "log"
    "os"

    "github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
    if err != nil {
        log.Fatalf("%s: %s", msg, err)
        panic(fmt.Sprintf("%s: %s", msg, err))
    }
}

func main() {
    message := os.Args[1]  // Get message from command-line arguments

    conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
    failOnError(err, "Failed to connect to RabbitMQ")
    defer conn.Close()

    ch, err := conn.Channel()
    failOnError(err, "Failed to open a channel")
    defer ch.Close()

    q, err := ch.QueueDeclare(
        "hello-queue", // name
        false,         // durable
        false,         // delete when unused
        false,         // exclusive
        false,         // no-wait
        nil,           // arguments
    )
    failOnError(err, "Failed to declare a queue")

    err = ch.Publish(
        "",         // exchange
        q.Name,     // routing key
        false,      // mandatory
        false,      // immediate
        amqp.Publishing{
            ContentType: "application/json",
            Body:        []byte(message),  // Use the command-line argument as message body
        })
    log.Printf(" [x] Sent %s", message)
    failOnError(err, "Failed to publish a message")
}