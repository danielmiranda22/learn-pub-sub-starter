package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/danmiranda227/learn-pub-sub-starter/internal/gamelogic"
	"github.com/danmiranda227/learn-pub-sub-starter/internal/pubsub"
	"github.com/danmiranda227/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril server...")

	rabbitConnStr := "amqp://guest:guest@localhost:5672/"

	conn, err := amqp.Dial(rabbitConnStr)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()
	fmt.Println("Peril game server connected to RabbitMQ server")

	// wait for a signal to exit control-c
	// wait for ctrl+c
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
	fmt.Println("Peril game server is shutting down.")

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	err = pubsub.PublishJSON(
		ch,
		routing.ExchangePerilDirect,
		routing.PauseKey,
		routing.PlayingState{IsPaused: true},
	)
	if err != nil {
		log.Fatalf("Failed to publish JSON message: %v", err)
	}

	gamelogic.PrintServerHelp()

	for {
		input := gamelogic.GetInput()
		if len(input) == 0 {
			continue
		}
		command := input[0]
		switch command {
		case "pause":
			fmt.Println("Sending a pause message...")
			err = pubsub.PublishJSON(
				ch,
				routing.ExchangePerilDirect,
				routing.PauseKey,
				routing.PlayingState{IsPaused: true},
			)
			if err != nil {
				log.Fatalf("Failed to publish JSON message: %v", err)
			}
		case "resume":
			fmt.Println("Sending a resume message...")
			err = pubsub.PublishJSON(
				ch,
				routing.ExchangePerilDirect,
				routing.PauseKey,
				routing.PlayingState{IsPaused: false},
			)
			if err != nil {
				log.Fatalf("Failed to publish JSON message: %v", err)
			}
		case "quit":
			fmt.Println("Quitting the server...")
			return
		case "help":
			gamelogic.PrintServerHelp()
		default:
			fmt.Println("Unknown command. Type 'help' for a list of commands.")
		}
	}
}
