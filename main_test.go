package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"todoapp/internal/requests"
)

func basicAuthRequest(username, password string) func(method, url string, body io.Reader) (*http.Response, error) {
	return func(method, url string, body io.Reader) (*http.Response, error) {
		req, err := http.NewRequest(method, url, body)
		if err != nil {
			return nil, err
		}

		req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(username+":"+password)))

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, err
		}

		return resp, nil
	}

}
func TestGetTodos(t *testing.T) {

	server := httptest.NewServer(middlewareAuthBasic("Todo Realm", &todo{}))
	defer server.Close()

	authRequest := basicAuthRequest("user1", "test")

	resp, err := authRequest("GET", server.URL+"/", nil)
	if err != nil {
		t.Fail()
	}

	if resp.StatusCode != http.StatusOK {
		t.Fail()
	}

}

func TestPostTodo(t *testing.T) {
	server := httptest.NewServer(middlewareAuthBasic("Todo Realm", &todo{Items: make(map[string][]todoitem)}))
	defer server.Close()

	authRequest := basicAuthRequest("user1", "test")
	var body bytes.Buffer
	json.NewEncoder(&body).Encode(requests.NewTodoRequest{Value: "Hello, WOrld"})
	resp, err := authRequest("POST", server.URL+"/", &body)
	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != http.StatusCreated {
		t.Fail()
	}
}
