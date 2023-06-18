package handler

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"
)

const (
	adminName     string = "ADMIN_NAME"
	adminPassword string = "ADMIN_PASSWORD"
)

type adminHandlers struct{}

func NewAdminHandler() *adminHandlers {
	h := &adminHandlers{}
	adminRoutes(h)
	return h
}

func adminRoutes(h *adminHandlers) {
	addRoutes([]route{
		newRouter("GET", "/admin/?", h.getInfo),
	})
}

func (h *adminHandlers) Handlers(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s", r.Method, r.URL)

	username, password, ok := r.BasicAuth()
	if !ok {
		writeDefaultStr(w, http.StatusUnauthorized, "basic authentication is required")
		return
	}
	envName := os.Getenv(adminName)
	envPassword := os.Getenv(adminPassword)
	if envName != username || envPassword != password {
		writeDefaultStr(w, http.StatusUnauthorized, "username or password is invalid")
		return
	}

	var allow []string
	for _, route := range routes {
		matches := route.regex.FindStringSubmatch(r.URL.Path)
		log.Printf("Admin Route: {%v, %v}, Matches: %v", route.method, route.regex.String(), matches)
		if len(matches) > 0 {
			if r.Method != route.method {
				allow = append(allow, route.method)
				continue
			}
			ctx := context.WithValue(r.Context(), ctxKey{}, matches[1:])
			route.handler(w, r.WithContext(ctx))
			return
		}
	}
	if len(allow) > 0 {
		w.Header().Add("Allow", strings.Join(allow, ", "))
		writeDefaultStr(w, http.StatusMethodNotAllowed, "405 method not allowed")
		return
	}
	http.NotFound(w, r)
}

func (h *adminHandlers) getInfo(w http.ResponseWriter, r *http.Request) {
	writeDefaultStr(w, http.StatusOK, "Hello, 世界")
}
