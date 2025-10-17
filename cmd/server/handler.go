package main

import (
	"fmt"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
)

func handlerWriteLog(gamelog routing.GameLog) pubsub.Acktype {
	defer fmt.Println("> ")
	err := gamelogic.WriteLog(gamelog)
	if err != nil {
		return pubsub.NackDiscard
	}
	return pubsub.Ack
}
