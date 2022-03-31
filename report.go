package main

import (
	"fmt"
	"log"
)

// todo: do a partial report on each check so that we are not collecting system memory endlessly

func (this *ArmoireChecker) report() {
	// todo: if we did no work, report that nothing was done, instead of saying "you started with no gold, etc..."
	fmt.Printf("\nREPORT\n")
	fmt.Println("-------")
	fmt.Printf("You started with %s gold.\n", showMeTheMoney(this.InitialGold))
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
	fmt.Println("-------")
}

func showMeTheMoney(gold float64) string {
	if gold <= 0 {
		return "no"
	}
	return fmt.Sprintf("%.0f", gold)
}
