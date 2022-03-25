package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/term"
)

type EscapeMode struct {
	signalChannel chan os.Signal
	escapeChannel chan struct{}
}

func NewEscapeMode(ctx context.Context, cancel context.CancelFunc) *EscapeMode {
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)
	return &EscapeMode{
		signalChannel: signalChannel,
		escapeChannel: make(chan struct{}),
	}
}
func (this *EscapeMode) Run() {
	go this.Listen()
	for {
		select {
		case receivedSignal := <-this.signalChannel:
			log.Println("[DEV] receivedSignal:", receivedSignal)
			os.Exit(1)
		case <-this.escapeChannel:
			log.Println("[DEV] received escape")
			os.Exit(27)
		}
	}
}

func (this *EscapeMode) Listen() {

	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	oneChar := []byte{1}
	//reader := io.LimitReader(os.Stdin, 1)
	//num, err := reader.Read(oneChar)
	num, err := os.Stdin.Read(oneChar)
	if err != nil {
		log.Println("[DEV] reader.read error:", err)
	}
	log.Printf("[DEV] reader.read num and oneChar: %#v, %#v", num, oneChar)
	if string(oneChar) == "\x1B" {
		log.Printf("[DEV] oneChar equals \\x1B: %s", oneChar)
	} else {
		log.Printf("[DEV] oneChar does not equal \\x1B: %s", oneChar)
		return
	}
	//utf8 escape is 1B (hex), 27 (dec)
	// U+241B
	// 2022/03/23 10:00:27 [DEV] limitreader.read num and oneChar: 1, []byte{0x1b}
	// 2022/03/23 10:00:27 [DEV] oneChar equals  [27]

	this.escapeChannel <- struct{}{}
}
