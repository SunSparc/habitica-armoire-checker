package main

import (
	"fmt"
	"log"
	"strings"
)

// todo: do a partial report on each check so that we are not collecting system memory endlessly

// todo: take all the report text and send it to a formatter to make the output look nice

func (this *ArmoireChecker) report() {
	// todo: if we did no work, report that nothing was done, instead of saying "you started with no gold, etc..."
	clearScreen()
	fmt.Println("==================================================")
	fmt.Println("                     REPORT")
	fmt.Println("==================================================")
	fmt.Printf("You started with %s gold.\n", showMeTheMoney(this.InitialGold))
	fmt.Printf("There were %d drops from your Enchanted Armoire.\n", this.DropsCount)
	//log.Printf("Armoire dropsmap: %#v\n", this.DropsMap)
	for dropType, drops := range this.DropsMap {
		fmt.Printf("\n%s: %d drops\n", strings.ToTitle(dropType), len(drops))
		fmt.Printf("--------------------------------------------------\n")
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
				fmt.Printf("- %s x %d\n", food, count)
			}
		case "gear":
			//fmt.Printf("gear drops: %#v\n", drops)
			for _, drop := range drops {
				fmt.Printf("- %s\n", drop.DropText)
			}
		default:
			log.Println("unknown dropType:", dropType, drops)
		}
	}
	fmt.Println("==================================================")
}

func showMeTheMoney(gold int64) string {
	if gold <= 0 {
		return "no"
	}
	return fmt.Sprintf("%.0d", gold)
}
