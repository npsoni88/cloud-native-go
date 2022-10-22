package main

import (
	"errors"
	"io"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
)

// Global variable for key value store
var store = struct {
	sync.RWMutex
	m map[string]string
}{m: make(map[string]string)}

var errNoSuchKey = errors.New("no such key")

func myFunc(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome to my simple API "))
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", myFunc)
	router.HandleFunc("/v1/{key}", keyValuePutHandler).Methods("PUT")
	router.HandleFunc("/v1/{key}", keyValueGetHandler).Methods("GET")
	router.HandleFunc("/v1/{key}", keyValueDeleteHandler).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8080", router))

}

/* Put function and handler of the API that must
- only match PUT requests for /v1/{key} path
- must respond with 201 (created) when a k-v is actually created
- must respond with 500 (internal error) in case of unexpected errors
*/

func Put(key, value string) error {
	store.Lock()
	store.m[key] = value
	store.Unlock()
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
- only match GET at /v1/{key}
- must call the get function from the API
- must respond with a 404 when key doesn't exist
- must respond with the value and a status 200 if the key exists
- must respond with 500 in case of unexpected errors
*/

func Get(key string) (string, error) {
	store.RLock()
	value, ok := store.m[key]
	store.RUnlock()
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

/* Delete function and handler that must
- only match DELETE at /v1/{key}
- must call the function to Delete
- respond with 200 status code if deleted
- respond with internal server error in case of unexpected error
*/

func Delete(key string) error {
	store.Lock()
	delete(store.m, key)
	store.Unlock()
	return nil
}

func keyValueDeleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	err := Delete(key)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
