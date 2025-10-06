package main

import (
	"fmt"
	"github.com/Mr-Rafael/gator/internal/config"
)

func main() {
	setErr := config.SetUser("Mr-Rafael")
	if setErr != nil {
		fmt.Printf("\nThere was an error setting the user: %v\n", setErr)
	}
	config.PrintCurrentConfig()
}