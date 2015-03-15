package ssignal

import (
	"os"
	"fmt"
	"os/signal"

	"meta"
	"slog"
)

var ch chan os.Signal

func Init() {
	ch = make(chan os.Signal, 1)
	go sigHanderTask()
	
	signal.Notify(ch, os.Interrupt, os.Kill)
}

func sigHanderTask() {
	slog.Sig = fmt.Sprintf("%v", <-ch)

	meta.ToClose<- true
	meta.ToClose<- true

	<-meta.GotIt
	<-meta.GotIt

	slog.ToClose<- true
	<-slog.GotIt

	os.Exit(0)
}

