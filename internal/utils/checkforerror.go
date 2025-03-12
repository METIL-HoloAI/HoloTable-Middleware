package utils

import "log"

func CheckForError(err error) {
	if err != nil {
		log.Fatal("Error occured, view error message:\n", err)
	}
}
