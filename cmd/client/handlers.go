package main

import (
	"fmt"

	"github.com/danmiranda227/learn-pub-sub-starter/internal/gamelogic"
	"github.com/danmiranda227/learn-pub-sub-starter/internal/routing"
)

func handlerPause(gs *gamelogic.GameState) func(routing.PlayingState) {
	return func(ps routing.PlayingState) {
		defer fmt.Print("> ")
		gs.HandlePause(ps)
	}
}
