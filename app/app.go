// Test APP
package main

import (
	"fmt"
	"github.com/dim13/golb"
	"log"
	"net/http"
)

const listen = ":8000"

func root(w http.ResponseWriter, r *http.Request, s []string) {
	fmt.Fprint(w, s)
}

func main() {
	d := golb.Open("test.json")
	if err := d.Write(); err != nil {
		log.Fatal(err)
	}
	if err := d.Read(); err != nil {
		log.Fatal(err)
	}
	re := new(golb.ReHandler)
	re.AddRoute("^/(\\d+)/(.*)$", root)
	log.Fatal(http.ListenAndServe(listen, re))
}
