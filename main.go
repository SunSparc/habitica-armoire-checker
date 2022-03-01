package main

import (
	"fmt"
	"log"
	"time"
)

var APIClient string

func main() {
	NewArmoireChecker().run()

	// todo: how do we want users to stop the application? watch for keypress? listen for signal?
	// todo: what if the application is canceled(signaled)?
	// todo: do a partial report on each check so that we are not collecting system memory endlessly
}

func NewArmoireChecker() *ArmoireChecker {
	//config := NewConfig(APIClient)
	return &ArmoireChecker{
		//Config:   config,
		Requester: NewRequester(),
		DropsMap:  map[string][]Armoire{},
	}
}
func (this *ArmoireChecker) run() {
	//os.Exit(0)
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
	err := this.getGoldAmount()
	if err != nil {
		log.Println("[ERROR] getGoldAmount:", err)
		return false
	}
	if this.InitialGold <= 0 {
		this.InitialGold = this.User.Data.Stats.Gold
	}
	//log.Println("gold:", gold)
	if this.User.Data.Stats.Gold < 55000 { // TODO: accept minimum gold level as parameter
		fmt.Println("No more gold, go earn some more :)")
		return false
	}
	err = this.checkArmoire()
	if err != nil {
		return false
	}
	this.recordResponse()
	return true
}
func (this *ArmoireChecker) recordResponse() {
	if !this.User.Data.Flags.ArmoireOpened {
		log.Println("[WARN] Armoire is not opened. How does that work?")
	}
	if !this.User.Data.Flags.ArmoireEnabled {
		log.Println("[WARN] Armoire is not enabled. Why?")
	}
	if !this.User.Data.Flags.ArmoireEmpty {
		//log.Println("[WARN] Armoire is not empty. That means we get some new gear. :)")
		// todo-maybe: make a toggle that announces when the empty status of the Armoire changes.
	}

	//log.Printf("response.Data.Armoire: %v\n", response.Data.Armoire)
	fmt.Print("*")

	this.DropsMap[this.User.Data.Armoire.Type] = append(this.DropsMap[this.User.Data.Armoire.Type], this.User.Data.Armoire)
	this.DropsCount = this.DropsCount + 1
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
