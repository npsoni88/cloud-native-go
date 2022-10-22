package main

import (
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Global variable for key value store
var store = make(map[string]string)
var errNoSuchKey = errors.New("no such key")

func myFunc(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome to my simple API "))
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", myFunc)
	router.HandleFunc("/v1/{key}", keyValuePutHandler).Methods("PUT")
	router.HandleFunc("/v1/{key}", keyValueGetHandler).Methods("GET")
	log.Fatal(http.ListenAndServe(":8080", router))

}

/* Put function and handler of the API that must
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

/* GET function and handler that must
- only match GET at /v1/key/{key}
- must call the get function from the API
- must respond with a 404 when key doesn't exist
- must respond with the value and a status 200 if the key exists
- must respond with 500 in case of unexpected errors
*/

func Get(key string) (string, error) {
	value, ok := store[key]
	if !ok {
		return "", errNoSuchKey
	}
	return value, nil

}

func keyValueGetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	value, err := Get(key)
	if errors.Is(err, errNoSuchKey) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte(value))
}
