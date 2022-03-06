package main

var APIClient string

func main() {
	NewArmoireChecker().run()
	// todo: before we start running the checker, see how much gold is available

	// todo: how do we want users to stop the application? watch for keypress? listen for signal?
	// todo: get the name of the user for personalization
	// todo: how do we want users to stop the application? watch for keypress? listen for signal? Escape key at any time?
	// todo: what if the application is canceled(signaled)?
	// todo: do a partial report on each check so that we are not collecting system memory endlessly
}
