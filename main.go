package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"
)

type todoitem struct {
	Value     string
	CreatedAt time.Time
}

func NewTodoItem(value string) todoitem {
	return todoitem{Value: value, CreatedAt: time.Now()}
}

type PrincipalContextType struct{}

var (
	principalKey = &PrincipalContextType{}
	todos        = make(map[string][]todoitem)
)

type NewTodoRequest struct {
	Value string `json:"value"`
}

func handleTodoStuff(w http.ResponseWriter, r *http.Request) {
	method := r.Method
	principal := r.Context().Value(principalKey).(string)
	log.Printf("%s logged in", principal)

	switch method {
	case "POST":
		var ri NewTodoRequest
		if err := json.NewDecoder(r.Body).Decode(&ri); err != nil {
			log.Fatal(err)
		}
		todos[principal] = append(todos[principal], NewTodoItem(ri.Value))
	case "GET":
		if err := json.NewEncoder(w).Encode(todos); err != nil {
			log.Fatal(err)
		}
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func main() {

	api := http.NewServeMux()
	api.HandleFunc("/", handleTodoStuff)

	root := http.NewServeMux()

	root.Handle("/", http.FileServer(http.Dir("public")))
	root.Handle("/todo/", http.StripPrefix("/todo", middlewareAuthBasic("Todo Realm", api)))

	if err := http.ListenAndServe(":8080", root); err != nil {
		log.Fatal(err)
	}
}

func middlewareAuthBasic(realm string, next http.Handler) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		authorizationHeader := r.Header.Get("authorization")
		if authorizationHeader == "" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Header().Set("www-authenticate", "Basic Realm")
			return
		}

		log.Println(authorizationHeader)
		principal, err := base64.StdEncoding.DecodeString(authorizationHeader[len("Basic "):])
		if err != nil {
			log.Println(err)
		}
		foo := strings.Split(string(principal), ":")
		// Skipping the user validation
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), principalKey, foo[0])))

	}

	return fn
}
