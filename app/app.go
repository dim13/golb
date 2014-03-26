// Test APP
package main

import (
	"github.com/dim13/golb"
	"log"
)

func main() {
	d := golb.Open("test.json")
	if err := d.Write(); err != nil {
		log.Fatal(err)
	}
	if err := d.Read(); err != nil {
		log.Fatal(err)
	}
}
