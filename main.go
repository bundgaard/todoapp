package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"todoapp/internal/requests"
)

type todo struct {
	mu    sync.RWMutex
	Items map[string][]todoitem
}

func (t *todo) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	method := r.Method
	path := r.URL.Path

	user := principal(r)

	info("%s %s %s\n", user, method, path)

	switch method {
	case "POST":
		t.postNewTodo(user, path)(w, r)
	case "GET":
		t.getAllTodo(user)(w, r)
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}

}

func (t *todo) postNewTodo(user, path string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var ri requests.NewTodoRequest
		if err := json.NewDecoder(r.Body).Decode(&ri); err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		t.mu.Lock()
		t.Items[user] = append(t.Items[user], NewTodoItem(ri.Value))
		t.mu.Unlock()

	}
}

func (t *todo) getAllTodo(user string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		t.mu.RLock()
		items, ok := t.Items[user]
		if !ok {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		if err := json.NewEncoder(w).Encode(items); err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		t.mu.RUnlock()
	}
}

func (t *todo) getSpecificTodo(user, value string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

var (
	principalKey = &PrincipalContextType{}
)

func main() {

	root := http.NewServeMux()

	root.Handle("/", http.FileServer(http.Dir("public")))
	root.Handle("/todo/", http.StripPrefix("/todo", middlewareAuthBasic("Todo Realm", &todo{Items: make(map[string][]todoitem)})))

	if err := http.ListenAndServe(":8080", root); err != nil {
		log.Fatal(err)
	}
}
