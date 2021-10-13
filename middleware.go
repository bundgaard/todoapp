package main

import (
	"context"
	"encoding/base64"
	"log"
	"net/http"
	"strings"
)

func middlewareAuthBasic(realm string, next http.Handler) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		authorizationHeader := r.Header.Get("authorization")
		if authorizationHeader == "" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Header().Set("www-authenticate", "Basic Realm")
			return
		}

		log.Println("header", authorizationHeader)
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
