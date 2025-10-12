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
	fmt.Println("Starting Peril client...")

	conn, err := amqp.Dial(connectionString)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}

	defer conn.Close()

	username, err := gamelogic.ClientWelcome()
	if err != nil {
		log.Fatalf("Error during client welcome: %v", err)
	}

	_, _, err = pubsub.DeclareAndBind(conn, routing.ExchangePerilDirect, routing.PauseKey+"."+username, routing.PauseKey, pubsub.TransientQueue)
	if err != nil {
		log.Fatalf("Failed to declare and bind pause queue: %v", err)
	}

	game_state := gamelogic.NewGameState(username)

	for {
		words := gamelogic.GetInput()
		if len(words) == 0 {
			continue
		}

		switch words[0] {
		case "spawn":
			fmt.Println("spawning...")
			err := game_state.CommandSpawn(words)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				continue
			}
		case "move":
			fmt.Println("moving...")
			_, err := game_state.CommandMove(words)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				continue
			}
		case "status":
			game_state.CommandStatus()
		case "help":
			gamelogic.PrintClientHelp()
		case "spam":
			fmt.Println("Spamming not allowed yet!")
		case "quit":
			gamelogic.PrintQuit()
			return
		default:
			fmt.Printf("unknown command: %s\n", words[0])
		}
	}
}
