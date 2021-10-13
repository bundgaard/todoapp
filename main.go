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

type PasswordEncoder interface {
	Encode([]byte) ([]byte, error)
}

type NoOpPasswordEncoder struct{}

func (n *NoOpPasswordEncoder) Encode(s1 []byte) ([]byte, error) {
	return s1, nil
}

type BCryptPasswordEncoder struct {
}

func (b *BCryptPasswordEncoder) Encode(s1 string) ([]byte, error) {
	// bcrypt.GenerateFromPassword()
	return nil, nil
}
func (t *todo) postNewTodo(user, path string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var ri requests.NewTodoRequest
		if err := json.NewDecoder(r.Body).Decode(&ri); err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		t.mu.Lock()
		defer t.mu.Unlock()
		t.Items[user] = append(t.Items[user], NewTodoItem(ri.Value))

		log.Println("Unlocked and sending back 201")
		w.WriteHeader(http.StatusCreated)
	}
}

func (t *todo) getAllTodo(user string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		t.mu.RLock()
		defer t.mu.RUnlock()
		items, ok := t.Items[user]

		if !ok {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		if err := json.NewEncoder(w).Encode(items); err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

	}
}

func (t *todo) getSpecificTodo(user, value string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

var (
	principalKey   = &PrincipalContextType{}
	basicAuthUsers = map[string]string{"user1": "test", "user2": "test"}
)

func basicAuthVerifier(username, password string) bool {
	pw, ok := basicAuthUsers[username]
	return ok && pw == password
}

func basicAuthPasswordEncoder(s1 []byte) ([]byte, error) {
	return nil, nil
}
func main() {

	root := http.NewServeMux()

	root.Handle("/", http.FileServer(http.Dir("public")))
	todoApi := &todo{Items: make(map[string][]todoitem)}
	root.Handle("/todo/", http.StripPrefix("/todo", NewBasicAuth("Todo Realm", basicAuthVerifier, todoApi.ServeHTTP)))
	log.Println("startin server on 8080")
	if err := http.ListenAndServe(":8080", root); err != nil {
		log.Fatal(err)
	}
}
