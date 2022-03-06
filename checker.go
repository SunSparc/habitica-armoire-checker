package main

import (
	"fmt"
	"log"
	"time"
)

func NewArmoireChecker() *ArmoireChecker {
	config := NewConfig(APIClient)
	return &ArmoireChecker{
		InitialGold:   -1,
		SpendingLimit: config.SpendingLimit,
		Requester:     NewRequester(*config),
		DropsCount:    0,
		DropsMap:      map[string][]Armoire{},
	}
}

func (this *ArmoireChecker) run() {
	fmt.Println("Checking your Enchanted Armoire")
	defer this.report()
	// todo: how do we want users to stop the application? watch for keypress? listen for signal? Escape key at any time?
	//       what if the application is canceled(signaled)?
	//       create ticker/timer to watch channels with termination signals

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
	if !this.goldReservesAreAdequate() {
		return false
	}
	err := this.checkArmoire()
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

func (this *ArmoireChecker) goldReservesAreAdequate() bool {
	err := this.getGoldAmount()
	if err != nil {
		log.Println("[ERROR] getGoldAmount:", err)
		return false
	}
	if this.InitialGold <= -1 {
		this.InitialGold = this.User.Data.Stats.Gold
	}
	if this.User.Data.Stats.Gold < this.SpendingLimit {
		fmt.Println("No more gold. Go earn some more. =)")
		return false
	}
	return true
}
