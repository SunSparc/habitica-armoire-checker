package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

func NewArmoireChecker() *ArmoireChecker {
	return &ArmoireChecker{
		Requester:  NewRequester(NewConfig(APIClient)),
		DropsCount: 0,
		DropsMap:   map[string][]Armoire{},
	}
}

func (this *ArmoireChecker) Run(ctx context.Context) {
	this.getInitialGold()
	this.getSpendLimit()
	// todo: use the context

	fmt.Println("Checking your Enchanted Armoire")
	defer this.report()
	// todo: what if the application is canceled(signaled)?
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

func (this *ArmoireChecker) getInitialGold() {
	err := this.getGoldAmount()
	if err != nil {
		log.Println("[ERROR] getGoldAmount:", err)
		return
	}
	this.InitialGold = int64(this.User.Data.Stats.Gold)
}
func (this *ArmoireChecker) getSpendLimit() {
	fmt.Printf("The Enchanted Armoire requires 100 gold each time it is opened.")
	fmt.Printf("You currently have %d gold.\n", this.InitialGold)

	fmt.Println("You can spend it all! Or set a limit.")
	fmt.Println("- 0 (no limit)")
	fmt.Println("- 1000 (example: only spend 1000 gold)")
	fmt.Print("Limit: ")
	reader := bufio.NewReader(os.Stdin)
	spendingLimit, err := reader.ReadString('\n')
	if err != nil {
		log.Println("[ERROR] reading spending limit:", err)
	}
	this.SpendingLimit, err = strconv.ParseInt(spendingLimit, 10, 64)
	if err != nil {
		log.Println("[ERROR]", err)
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
	if (int64(this.User.Data.Stats.Gold) - 100) < this.SpendingLimit {
		fmt.Println("Spending limit reached. Change the limit or go earn some more gold. =)")
		return false
	}
	if (this.User.Data.Stats.Gold - 100) < 0 {
		fmt.Println("Insufficient funds.")
		return false
	}
	return true
}
