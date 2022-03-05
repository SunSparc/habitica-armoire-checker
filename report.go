package main

import (
	"fmt"
	"log"
)

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
