package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

var APIClient string

func main() {
	NewArmoireChecker(NewConfig(APIClient)).run()

	// todo: how do we want users to stop the application? watch for keypress? listen for signal?
	// todo: what if the application is canceled(signaled)?
}

func NewArmoireChecker(config *Config) *ArmoireChecker {
	return &ArmoireChecker{
		Config:   config,
		DropsMap: map[string][]Armoire{},
	}
}
func (this *ArmoireChecker) run() {
	os.Exit(0)
	fmt.Println("Checking your Enchanted Armoire")
	defer this.report()

	for x := 0; x <= 5; x++ {
		if !this.check() {
			//log.Println("[runner] done")
			break
		}
		for t := 0; t < 30; t++ {
			fmt.Print(".")
			time.Sleep(time.Second * 2) // no faster than 1 request every 30 seconds
		}
	}
}
func (this *ArmoireChecker) check() bool {
	gold, err := getGoldAmount("user?userFields=stats.gp", this.Config)
	if this.InitialGold <= 0 {
		this.InitialGold = gold
	}
	if err != nil {
		log.Println("main.getGoldAmount error:", err)
		return false
	}
	//log.Println("gold:", gold)
	if gold < 55000 { // TODO: accept minimum gold level as parameter
		fmt.Println("No more gold, go earn some more :)")
		return false
	}
	response := doArmoireRequest(this.Config)
	if !responseIsFavorable(response) {
		return false
	}
	this.recordResponse(response)
	return true
}
func (this *ArmoireChecker) recordResponse(response User) {
	if !response.Data.Flags.ArmoireOpened {
		log.Println("[WARN] Armoire is not opened. How does that work?")
	}
	if !response.Data.Flags.ArmoireEnabled {
		log.Println("[WARN] Armoire is not enabled. Why?")
	}
	if !response.Data.Flags.ArmoireEmpty {
		//log.Println("[WARN] Armoire is not empty. That means we get some new gear. :)")
		// todo-maybe: make a toggle that announces when the empty status of the Armoire changes.
	}

	//log.Printf("response.Data.Armoire: %v\n", response.Data.Armoire)
	fmt.Print("*")

	this.DropsMap[response.Data.Armoire.Type] = append(this.DropsMap[response.Data.Armoire.Type], response.Data.Armoire)
	this.DropsCount = this.DropsCount + 1
}
func responseIsFavorable(response User) bool {
	if response.StatusCode == 200 &&
		response.Success &&
		response.Data.Flags.ArmoireEnabled {
		return true
	}
	log.Println("response was not favorable:")
	log.Println("response.StatusCode:", response.StatusCode)
	log.Println("response.Success:", response.Success)
	log.Println("response.Data.Flags.ArmoireEnabled:", response.Data.Flags.ArmoireEnabled)
	log.Println("response.Data.Flags.ArmoireEmpty:", response.Data.Flags.ArmoireEmpty)
	return false
}
func (this *ArmoireChecker) report() {
	fmt.Printf("\nYou started with %.0f gold.\n", this.InitialGold)
	fmt.Printf("There were %d drops from your Enchanted Armoire.\n", this.DropsCount)
	//log.Printf("Armoire dropsmap: %#v\n", this.DropsMap)
	for dropType, drops := range this.DropsMap {
		fmt.Printf("\n%s: %d drops\n", dropType, len(drops))
		fmt.Printf("-------------------\n")
		switch dropType {
		case "experience":
			xpTotal := 0
			for _, drop := range drops {
				xpTotal = xpTotal + drop.Value
			}
			fmt.Println("Total experience dropped:", xpTotal)
		case "food":
			//log.Printf("food drops: %#v", drops)
			foodMap := map[string]int{}
			for _, drop := range drops {
				foodMap[drop.DropKey] = foodMap[drop.DropKey] + 1
			}
			//fmt.Printf("food map: %#v", foodMap)
			for food, count := range foodMap {
				fmt.Printf("- %s x %d", food, count)
			}
		case "gear":
			fmt.Printf("gear drops: %#v", drops)
		default:
			log.Println("unknown dropType:", dropType, drops)
		}
	}
}
