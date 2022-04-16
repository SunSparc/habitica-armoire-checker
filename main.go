package main

import (
	"context"
	"fmt"
	"os/exec"
)

var APIClient string

func main() {
	//fmt.Print("\x1bc") // clear screen, works on macOS
	exec.Command("cmd", "/c", "cls").Run() // clear screen, works on macOS

	fmt.Println(welcomeText)
	ctx, cancel := context.WithCancel(context.Background())
	// todo: do config first, then start escape mode, then start checker ??
	NewArmoireChecker(NewEscapeMode(ctx, cancel)).Run(ctx) // todo: go
	//NewEscapeMode(ctx, cancel).Run() // todo: run this while the checker is going

	// todo: get the name of the user for personalization
}

const welcomeText = `
Habitica Armoire Checker
=========================
`

////////////////////////////////////
