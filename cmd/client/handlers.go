package main

import (
	"fmt"
	"time"

	"github.com/danmiranda227/learn-pub-sub-starter/internal/gamelogic"
	"github.com/danmiranda227/learn-pub-sub-starter/internal/pubsub"
	"github.com/danmiranda227/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func handlerPause(gs *gamelogic.GameState) func(routing.PlayingState) pubsub.Actype {
	return func(ps routing.PlayingState) pubsub.Actype {
		defer fmt.Print("> ")
		gs.HandlePause(ps)
		return pubsub.Ack
	}
}

func handlerMove(gs *gamelogic.GameState, ch *amqp.Channel) func(gamelogic.ArmyMove) pubsub.Actype {
	return func(move gamelogic.ArmyMove) pubsub.Actype {
		defer fmt.Print("> ")
		moveOutcome := gs.HandleMove(move)
		switch moveOutcome {
		case gamelogic.MoveOutcomeSamePlayer:
			return pubsub.NackDiscard
		case gamelogic.MoveOutComeSafe:
			return pubsub.Ack
		case gamelogic.MoveOutcomeMakeWar:
			err := pubsub.PublishJSON(
				ch,
				routing.ExchangePerilTopic,
				routing.WarRecognitionsPrefix+"."+gs.GetUsername(),
				gamelogic.RecognitionOfWar{
					Attacker: move.Player,
					Defender: gs.GetPlayerSnap(),
				},
			)
			if err != nil {
				fmt.Printf("error publishing war recognition: %v\n", err)
				return pubsub.NackRequeue
			}
			return pubsub.Ack
		}
		fmt.Println("error: unknown move outcome")
		return pubsub.NackDiscard
	}
}

func handlerWar(gs *gamelogic.GameState, ch *amqp.Channel) func(gamelogic.RecognitionOfWar) pubsub.Actype {
	return func(rw gamelogic.RecognitionOfWar) pubsub.Actype {
		defer fmt.Print("> ")
		outcome, winner, loser := gs.HandleWar(rw)

		var message string
		switch outcome {
		case gamelogic.WarOutcomeNotInvolved:
			return pubsub.NackRequeue
		case gamelogic.WarOutcomeNoUnits:
			return pubsub.NackDiscard
		case gamelogic.WarOutcomeOpponentWon:
			message = fmt.Sprintf("%s won a war against %s", winner, loser)
		case gamelogic.WarOutcomeYouWon:
			message = fmt.Sprintf("%s won a war against %s", winner, loser)
		case gamelogic.WarOutcomeDraw:
			message = fmt.Sprintf("A war between %s and %s resulted in a draw", winner, loser)
		default:
			fmt.Println("error: unknown war outcome")
			return pubsub.NackDiscard
		}

		err := publishGameLog(ch, rw.Attacker.Username, message)
		if err != nil {
			fmt.Printf("error publishing game log: %v\n", err)
			return pubsub.NackRequeue
		}
		return pubsub.Ack
	}
}

func publishGameLog(ch *amqp.Channel, username, message string) error {
	return pubsub.PublishGob(
		ch,
		routing.ExchangePerilTopic,
		routing.GameLogSlug+"."+username,
		routing.GameLog{
			CurrentTime: time.Now(),
			Message:     message,
			Username:    username,
		},
	)
}
