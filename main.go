package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func myFunc(w http.ResponseWriter, r *http.Request) {
	hostname := r.Host
	w.Write([]byte("hello func, executed on " + hostname))
}

func productsKey(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	w.Write([]byte("you came looking for " + key))
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", myFunc)
	router.HandleFunc("/products/{key}", productsKey)
	log.Fatal(http.ListenAndServe(":8080", router))

}
