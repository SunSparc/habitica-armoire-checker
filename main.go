package main

import "fmt"

var APIClient string

func main() {
	fmt.Println(welcomeText)
	NewArmoireChecker().run()

	// todo: get the name of the user for personalization
}

const welcomeText = `
Welcome to the Habitica Armoire Checker
========================================
`
