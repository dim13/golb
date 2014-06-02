package main

import (
	"fmt"
	"net/http"
)

func robotsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "User-agent: *")
	fmt.Fprintln(w, "Sitemap:", "http://"+r.Host+"/sitemap.xml")
}
