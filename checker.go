package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func NewArmoireChecker(escapeMode *EscapeMode) *ArmoireChecker {
	return &ArmoireChecker{
		//Requester:  NewRequester(NewConfig(APIClient)),
		DropsCount: 0,
		DropsMap:   map[string][]Armoire{},
		EscapeMode: escapeMode,
	}
}

func (this *ArmoireChecker) Run(ctx context.Context) {
	//go this.EscapeMode.Run()
	reset := false
	for {
		this.Requester = NewRequester(NewConfig(APIClient, reset))
		if this.getInitialGold() {
			break
		}
		reset = true
	}
	this.getSpendLimit()
	this.manageChecker(ctx)
}

func (this *ArmoireChecker) manageChecker(ctx context.Context) {
	// todo: use the context

	fmt.Println("\n-------------------------------")
	fmt.Println("Checking your Enchanted Armoire")
	//go this.EscapeMode.Run()

	defer this.report()

	// todo: what if the application is canceled(signaled)?
	//       create ticker/timer to watch channels with termination signals

	for {
		if !this.check() {
			break
		}
		for t := 0; t < 30; t++ {
			fmt.Print(".")
			time.Sleep(time.Second * 1) // note: no faster than 1 request every 30 seconds
		}
	}
}

func (this *ArmoireChecker) getInitialGold() bool {
	err := this.getGoldAmount()
	if err != nil {
		log.Println("[ERROR] getGoldAmount:", err)
		return false
	}
	this.InitialGold = int64(this.User.Data.Stats.Gold)
	fmt.Println("\nSuccess! We are connected to Habitica.")
	fmt.Println("----------------------------------------")
	fmt.Println()
	return true
}
func (this *ArmoireChecker) getSpendLimit() {
	fmt.Printf("The Enchanted Armoire requires 100 gold\n  each time it is opened.\n")
	fmt.Printf("\nYou currently have %d gold.\n", this.InitialGold)

	fmt.Println("You can spend it all! Or set a limit.")
	fmt.Println("  Examples:")
	fmt.Println("  - 0 (no limit, Spend it all!!)")
	fmt.Println("  - 1000 (only spend 1000 gold)")
	fmt.Println()
	fmt.Print("Spending Limit: ")
	reader := bufio.NewReader(os.Stdin)
	spendingLimit, err := reader.ReadString('\n')
	if err != nil {
		log.Println("[ERROR] reading spending limit:", err)
	}
	this.SpendingLimit, err = strconv.ParseInt(strings.TrimSpace(spendingLimit), 10, 64)
	if err != nil {
		log.Println("[ERROR] parsing spending limit:", err)
	}
	// todo: what if the user is silly and inputs text or other nonsense other than numbers?

	if this.SpendingLimit > this.InitialGold {
		log.Printf("Sweet! We are just going to blow through\n  the whole pile of gold!\n\n")
		this.SpendingLimit = 0
	} else {
		this.SpendingLimit = this.InitialGold - this.SpendingLimit
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
		fmt.Println("\nSpending limit reached. =)")
		return false
	}
	if (this.User.Data.Stats.Gold - 100) < 0 {
		fmt.Println("\nInsufficient funds to continue.")
		return false
	}
	return true
}
