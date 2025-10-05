package main

import (
	"fmt"
	"github.com/Mr-Rafael/gator/internal/config"
)

func main() {
	gatorConf, err := config.Read()
	if err != nil {
		fmt.Printf("\nThere was an error reading the configuration: %v\n", err)
	}
	fmt.Printf("\nCurrent Gator configuration: %v\n", gatorConf)
}