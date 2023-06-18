// Launch and serve customer Rest API
package main

import (
	"github.com/clebersonp/customer-go-rest-api/handler"
	"log"
	"net/http"
)

const (
	host = "localhost:8080"
)

func main() {
	// sets the log pattern
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)

	// customer handlers
	customerHandlers := handler.NewCustomerHandler().Handlers
	http.HandleFunc("/customers", customerHandlers)
	http.HandleFunc("/customers/", customerHandlers)

	// admin handlers
	adminHandlers := handler.NewAdminHandler().Handlers
	http.HandleFunc("/admin", adminHandlers)
	http.HandleFunc("/admin/", adminHandlers)

	// msg for listening at
	log.Printf("The customer Rest API has just been started on '%s'\n", host)

	// starts the server and listen to
	log.Fatal(http.ListenAndServe(host, nil))

	// TODO documentation
}
