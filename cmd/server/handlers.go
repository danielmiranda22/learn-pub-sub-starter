package main

import (
	"fmt"

	"github.com/danmiranda227/learn-pub-sub-starter/internal/gamelogic"
	"github.com/danmiranda227/learn-pub-sub-starter/internal/pubsub"
	"github.com/danmiranda227/learn-pub-sub-starter/internal/routing"
)

func handlerGameLog() func(routing.GameLog) pubsub.Actype {
	return func(gl routing.GameLog) pubsub.Actype {
		defer fmt.Print("> ")
		err := gamelogic.WriteLog(gl)
		if err != nil {
			fmt.Printf("error writing log: %v\n", err)
			return pubsub.NackRequeue
		}
		return pubsub.Ack
	}
}
