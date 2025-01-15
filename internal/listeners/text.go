package listeners

import "fmt"

func StartTextListener() {
	for {
		fmt.Println("Please type your prompt or one of the following options:")
		fmt.Println("(q)uit, (r)eload yaml")

		var input string
		fmt.Scanln(&input)

		if input == "q" {
			return
		}

		if input == "r" {
			fmt.Println("TODO")
		}

		fmt.Print("\033[H\033[2J")
	}
}
