package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type basicAuthMiddleware struct {
	Realm    string                    // Just a Realm
	Verifier func(string, string) bool // Verifier func is the custom function to verify incoming users, here we can setup LDAP / Password file or whatever, read .htaccess
	Next     http.HandlerFunc          // next is the handlerfunc we want to protect
}

func NewBasicAuth(realm string, verifier func(string, string) bool, next http.HandlerFunc) *basicAuthMiddleware {
	return &basicAuthMiddleware{
		Realm:    realm,
		Next:     next,
		Verifier: verifier}
}
func (bam *basicAuthMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	authorizationHeader := r.Header.Get("authorization")
	if authorizationHeader == "" {
		w.Header().Set("www-authenticate", fmt.Sprintf("Basic realm=%q", bam.Realm))
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	details, err := base64.StdEncoding.DecodeString(authorizationHeader[len("Basic "):])
	if err != nil {
		log.Println("base64 decode", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	usernameAndPassword := strings.Split(string(details), ":")

	if !bam.Verifier(usernameAndPassword[0], usernameAndPassword[1]) {
		log.Println("verifier: password or username was incorrect")
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	bam.Next(w, r.WithContext(context.WithValue(r.Context(), principalKey, usernameAndPassword[0])))

}
