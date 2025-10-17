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

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	username, err := gamelogic.ClientWelcome()
	if err != nil {
		log.Fatalf("Error during client welcome: %v", err)
	}

	gs := gamelogic.NewGameState(username)

	err = pubsub.SubscribeJSON(
		conn,
		routing.ExchangePerilDirect,
		routing.PauseKey+"."+gs.GetUsername(),
		routing.PauseKey,
		pubsub.TransientQueue,
		handlerPause(gs),
	)
	if err != nil {
		log.Fatalf("Failed to subscribe to pause messages: %v", err)
	}

	err = pubsub.SubscribeJSON(
		conn,
		routing.ExchangePerilTopic,
		routing.ArmyMovesPrefix+"."+gs.GetUsername(),
		routing.ArmyMovesPrefix+".*",
		pubsub.TransientQueue,
		handlerMove(gs),
	)
	if err != nil {
		log.Fatalf("Failed to subscribe to army move messages: %v", err)
	}

	for {
		words := gamelogic.GetInput()
		if len(words) == 0 {
			continue
		}

		switch words[0] {
		case "spawn":
			fmt.Println("spawning...")
			err := gs.CommandSpawn(words)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				continue
			}
		case "move":
			fmt.Println("moving...")
			_, err := gs.CommandMove(words)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				continue
			}
			mv, err := gs.CommandMove(words)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				continue
			}

			err = pubsub.PublishJSON(
				ch,
				routing.ExchangePerilTopic,
				routing.ArmyMovesPrefix+"."+gs.GetUsername(),
				mv,
			)
			if err != nil {
				fmt.Printf("Failed to publish move: %v\n", err)
			}

			fmt.Println("Move command processed.")
		case "status":
			gs.CommandStatus()
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

func handlerPause(gs *gamelogic.GameState) func(routing.PlayingState) pubsub.Acktype {
	return func(ps routing.PlayingState) pubsub.Acktype {
		defer fmt.Print("> ")
		gs.HandlePause(ps)
		return pubsub.Ack
	}
}

func handlerMove(gs *gamelogic.GameState) func(gamelogic.ArmyMove) pubsub.Acktype {
	return func(am gamelogic.ArmyMove) pubsub.Acktype {
		defer fmt.Print("> ")
		mv_outcome := gs.HandleMove(am)
		switch mv_outcome {
		case gamelogic.MoveOutComeSafe:
			return pubsub.Ack
		case gamelogic.MoveOutcomeMakeWar:
			return pubsub.Ack
		case gamelogic.MoveOutcomeSamePlayer:
			return pubsub.NackDiscard
		default:
			return pubsub.NackDiscard
		}
	}
}
