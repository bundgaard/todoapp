package main

import "net/http"

type PrincipalContextType struct{}

func principal(r *http.Request) string {
	user := r.Context().Value(principalKey)
	if user == nil {
		return "anonymous"
	}
	return user.(string)
}
