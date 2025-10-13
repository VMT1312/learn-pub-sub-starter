package main

import (
	"fmt"
	"log"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

const connectionString = "amqp://guest:guest@localhost:5672/"

func main() {
	fmt.Println("Starting Peril server...")

	conn, err := amqp.Dial(connectionString)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}

	defer conn.Close()

	fmt.Println("Connected to RabbitMQ successfully")

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	_, _, err = pubsub.DeclareAndBind(
		conn, routing.ExchangePerilTopic,
		"game_logs", routing.GameLogSlug+".*",
		pubsub.DurableQueue,
	)
	if err != nil {
		log.Fatalf("Failed to declare and bind game logs queue: %v", err)
	}

	for {
		words := gamelogic.GetInput()
		if len(words) == 0 {
			continue
		}

		switch words[0] {
		case "pause":
			fmt.Println("sending pause")
			pubsub.PublishJSON(
				ch, routing.ExchangePerilDirect,
				routing.PauseKey, routing.PlayingState{IsPaused: true},
			)
		case "resume":
			fmt.Println("sending resume")
			pubsub.PublishJSON(
				ch, routing.ExchangePerilDirect,
				routing.PauseKey,
				routing.PlayingState{IsPaused: false},
			)
		case "quit":
			fmt.Println("exiting...")
			return
		default:
			fmt.Printf("unknown command: %s\n", words[0])
		}
	}
}
