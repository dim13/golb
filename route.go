package main

import (
	"log"
	"net/http"
	"regexp"
)

type HandlerFunc http.HandlerFunc

type SelectHandler interface {
	http.Handler
	Select([]string)
	Store(*http.Request)
}

type route struct {
	re      *regexp.Regexp
	handler SelectHandler
}

type ReHandler struct {
	routes []*route
}

func (f HandlerFunc) Select(s []string)                                {}
func (f HandlerFunc) Store(r *http.Request)                            {}
func (f HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) { f(w, r) }

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
	for _, route := range h.routes {
		matches := route.re.FindStringSubmatch(r.URL.Path)
		if matches != nil {
			route.handler.Select(matches[1:])
			if r.Method == "POST" {
				route.handler.Store(r)
				//http.Redirect(w, r, r.URL.Path, http.StatusFound)
			}
			route.handler.ServeHTTP(w, r)
			return
		}
	}
	http.NotFound(w, r)
}
