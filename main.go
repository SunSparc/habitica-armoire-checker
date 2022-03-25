package main

import (
	"context"
	"fmt"
)

var APIClient string

func main() {
	fmt.Println(welcomeText)
	ctx, cancel := context.WithCancel(context.Background())
	NewArmoireChecker(ctx).Run()     // todo: go
	NewEscapeMode(ctx, cancel).Run() // todo: run this while the checker is going

	// todo: get the name of the user for personalization
}

const welcomeText = `
Habitica Armoire Checker
=========================
`

////////////////////////////////////
