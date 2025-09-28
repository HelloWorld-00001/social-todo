package common

import (
	"log"
)

// Recovery is a reusable panic recovery helper.
// Call it inside defer to catch panics and log them.
func Recovery() {
	if r := recover(); r != nil {
		log.Printf("Recovered from panic: %v", r)
	}
}
