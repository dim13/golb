package route

import (
	"log"
	"net/http"
	"regexp"
)

type HandlerFunc http.HandlerFunc

type SelectHandler interface {
	Select([]string) bool
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type route struct {
	re      *regexp.Regexp
	handler SelectHandler
}

type ReHandler struct {
	routes   []*route
	notFound HandlerFunc
}

func New() *ReHandler {
	return &ReHandler{notFound: http.NotFound}
}

func (f HandlerFunc) Select(_ []string) bool {
	return true
}

func (f HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f(w, r)
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

func (h *ReHandler) NotFound(handler HandlerFunc) {
	h.notFound = handler
}

func (h *ReHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, r.URL)
	for _, route := range h.routes {
		matches := route.re.FindStringSubmatch(r.URL.Path)
		if matches != nil {
			log.Println("Match", matches, r.URL)
			if !route.handler.Select(matches[1:]) {
				log.Println(route.re, "NotFound")
				h.notFound(w, r)
				return
			}
			r.ParseForm()
			route.handler.ServeHTTP(w, r)
			if r.Method == "POST" {
				http.Redirect(w, r, r.URL.Path, http.StatusFound)
			}
			return
		}
	}
	h.notFound(w, r)
}
