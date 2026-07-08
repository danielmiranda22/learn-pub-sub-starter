package main

import (
	"fmt"
	"log"

	"github.com/danmiranda227/learn-pub-sub-starter/internal/gamelogic"
	"github.com/danmiranda227/learn-pub-sub-starter/internal/pubsub"
	"github.com/danmiranda227/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril client...")

	rabbitConnStr := "amqp://guest:guest@localhost:5672/"
	conn, err := amqp.Dial(rabbitConnStr)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()
	fmt.Println("Peril game client connected to RabbitMQ server")

	username, err := gamelogic.ClientWelcome()
	if err != nil {
		log.Fatalf("Failed to welcome client: %v", err)
	}
	fmt.Printf("Welcome, %s!\n", username)

	_, _, err = pubsub.DeclareAndBind(
		conn,
		routing.ExchangePerilDirect,
		routing.PauseKey+"."+username,
		routing.PauseKey,
		pubsub.QueueTransient,
	)
	if err != nil {
		log.Fatalf("Failed to declare and bind queue: %v", err)
	}

	gs := gamelogic.NewGameState(username)

	for {
		words := gamelogic.GetInput()
		if len(words) == 0 {
			continue
		}
		switch words[0] {
		case "move":
			_, err := gs.CommandMove(words)
			if err != nil {
				fmt.Println(err)
				continue
			}
			// TODO: publish the move
		case "spawn":
			err = gs.CommandSpawn(words)
			if err != nil {
				fmt.Println(err)
				continue
			}
		case "status":
			gs.CommandStatus()
		case "help":
			gamelogic.PrintClientHelp()
		case "spam":
			// TODO: publish n malicious logs
			fmt.Println("Spamming not allowed yet!")
		case "quit":
			gamelogic.PrintQuit()
			return
		default:
			fmt.Println("unknown command")
		}
	}
}
