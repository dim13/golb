package main

import (
	"log"
	"net/http"
	"regexp"
)

type HandlerFunc http.HandlerFunc

type SelectHandler interface {
	Select([]string)
	Get(w http.ResponseWriter, r *http.Request)
	Post(w http.ResponseWriter, r *http.Request)
}

type route struct {
	re      *regexp.Regexp
	handler SelectHandler
}

type ReHandler struct {
	routes []*route
}

func (f HandlerFunc) Select(_ []string) {}
func (f HandlerFunc) Get(w http.ResponseWriter, r *http.Request) {
	f(w, r)
}
func (f HandlerFunc) Post(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, r.URL.Path, http.StatusFound)
}

func (h *ReHandler) Handle(re string, handler SelectHandler) {
	log.Println("SelectHandler", re)
	r := &route{
		re:      regexp.MustCompile(re),
		handler: handler,
	}
	h.routes = append(h.routes, r)
}

func (h *ReHandler) HandleFunc(re string, handler HandlerFunc) {
	log.Println("HandlerFunc", re)
	r := &route{
		re:      regexp.MustCompile(re),
		handler: handler,
	}
	h.routes = append(h.routes, r)
}

func (h *ReHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, r.URL)
	for _, route := range h.routes {
		matches := route.re.FindStringSubmatch(r.URL.Path)
		if matches != nil {
			log.Println("Match", matches, r.URL)
			route.handler.Select(matches[1:])
			switch r.Method {
			case "POST":
				r.ParseForm()
				route.handler.Post(w, r)
			case "GET":
				route.handler.Get(w, r)
			}
			return
		}
	}
	http.NotFound(w, r)
}
