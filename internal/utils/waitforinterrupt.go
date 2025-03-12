package utils

import "fmt"

func WaitForInterrupt() {
	var input string
	for input != "quit" {
		_, err := fmt.Scanln(&input)
		CheckForError(err)
	}
}
