package main

import (
	"context"
	"fmt"
	"time"
)

var APIClient string

func main() {
	clearScreen()
	fmt.Println(welcomeText)
	ctx, cancel := context.WithCancel(context.Background())
	// todo: do config first, then start escape mode, then start checker ??
	NewArmoireChecker(NewEscapeMode(ctx, cancel)).Run(ctx) // todo: go
	//NewEscapeMode(ctx, cancel).Run() // todo: run this while the checker is going

	// todo: get the name of the user for personalization

	time.Sleep(time.Second * 60) // todo: press any key to exit
}

func clearScreen() {
	fmt.Print("\x1bc") // clear screen, works on macOS
	//exec.Command("cmd", "/c", "cls").Run() // clear screen, works on macOS
}

const welcomeText = `
Habitica Armoire Checker
=========================
`

////////////////////////////////////
