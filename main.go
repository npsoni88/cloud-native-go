package main

import (
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Global variable for key value store
var store = make(map[string]string)

func myFunc(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome to my simple API "))
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", myFunc)
	router.HandleFunc("/v1/{key}", keyValuePutHandler).Methods("PUT")
	log.Fatal(http.ListenAndServe(":8080", router))

}

/* Put function of the API that must
- only match PUT requests for /v1/key/{key} path
- must respond with 201 (created) when a k-v is actually created
- must respond with 500 (internal error) in case of unexpected errors
*/

func Put(key, value string) error {
	store[key] = value
	return nil
}

func keyValuePutHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = Put(key, string(body))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusCreated)

}
