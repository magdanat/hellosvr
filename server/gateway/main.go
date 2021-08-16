package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

//MethodMux sends the request to the function
//associated with the HTTP request method
type MethodMux struct {
	//use a map where the  key is a string (method name)
	//and the value is the associated handler function
	HandlerFuncs map[string]func(http.ResponseWriter, *http.Request)
}

//ServeHTTP sends the request to the appropriate handler based on
//the HTTP method in the request
func (mm *MethodMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//r.Method will be the method used in the request (GET, PUT, PATCH, POST, etc.)
	fn := mm.HandlerFuncs[r.Method]
	if fn != nil {
		fn(w, r)
	} else {
		http.Error(w, "that method is not allowed", http.StatusMethodNotAllowed)
	}
}

//NewMethodMux constructs a new MethodMux
func NewMethodMux() *MethodMux {
	return &MethodMux{
		HandlerFuncs: map[string]func(http.ResponseWriter, *http.Request){},
	}
}

//HelloHandler handles requests for the `/hello` resource
func HelloHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, Web!\n"))
}

func main() {
	//get the value of the ADDR environment variable
	addr := os.Getenv("ADDR")

	//if it's blank, default to ":80", which means
	//listen port 80 for requests addressed to any host
	if len(addr) == 0 {
		addr = ":80"
	}

	//get the TLS key and cert paths from environment variables
	//this allows us to use a self-signed cert/key during development
	//and the Let's Encrypt cert/key in production
	tlsKeyPath := os.Getenv("TLSKEY")
	tlsCertPath := os.Getenv("TLSCERT")

	// if len(tlsKeyPath) < 0 || len(tlsCertPath) < 0 {
	// 	log.Fatal("No envrionment variable found for either TLSKEY or TLSCERT")
	// }

	//create a new mux (router)
	//the mux calls different functions for
	//different resource paths
	mux := http.NewServeMux()

	methmux := NewMethodMux()
	methmux.HandlerFuncs["GET"] = HelloHandler

	mux.Handle("/hello", methmux)

	//start the web server using the mux as the root handler,
	//and report any errors that occur.
	//the ListenAndServe() function will block so
	//this program will continue to run until killed
	//start the server
	// fmt.Printf(addr)
	fmt.Printf("listening on %s...\n", addr)
	log.Fatal(http.ListenAndServeTLS(addr, tlsCertPath, tlsKeyPath, mux))
}
