//+build wireinject

package main

import "github.com/google/wire"

func initializeEvent(msg string) event {
	wire.Build(newMessage, newGreeter, newEvent)
	return event{}
}
