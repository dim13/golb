package gold

import (
	"net/http"
	"regexp"
)

type HandlerFunc http.HandlerFunc

type HandlerMatch interface {
	http.Handler
	StoreMatch([]string)
}

type route struct {
	re      *regexp.Regexp
	handler HandlerMatch
}

type ReHandler struct {
	routes []*route
}

func (f HandlerFunc) StoreMatch(s []string)                            {}
func (f HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) { f(w, r) }

func (h *ReHandler) Handle(re string, handler HandlerMatch) {
	r := &route{regexp.MustCompile(re), handler}
	h.routes = append(h.routes, r)
}

func (h *ReHandler) HandleFunc(re string, handler func(http.ResponseWriter, *http.Request)) {
	r := &route{
		re:      regexp.MustCompile(re),
		handler: HandlerFunc(handler),
	}
	h.routes = append(h.routes, r)
}

func (h *ReHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, route := range h.routes {
		matches := route.re.FindStringSubmatch(r.URL.Path)
		if matches != nil {
			route.handler.StoreMatch(matches[1:])
			route.handler.ServeHTTP(w, r)
			return
		}
	}
	http.NotFound(w, r)
}
