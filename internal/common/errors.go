package common

import (
	"log"
)

// CheckErr reusable error check
func CheckErr(e error) {
	if e != nil {
		log.Fatal(e)
	}
}
