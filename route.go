package golb

import (
	"net/http"
	"regexp"
)

type RouteHandler func(http.ResponseWriter, *http.Request, []string)

type route struct {
	re      *regexp.Regexp
	handler RouteHandler
}

type ReHandler struct {
	routes []*route
}

func (h *ReHandler) AddRoute(re string, handler RouteHandler) {
	r := &route{regexp.MustCompile(re), handler}
	h.routes = append(h.routes, r)
}

func (h *ReHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, route := range h.routes {
		matches := route.re.FindStringSubmatch(r.URL.Path)
		if matches != nil {
			route.handler(w, r, matches)
			return
		}
	}
	http.NotFound(w, r)
}