package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"io/ioutil"
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

	todoApi := &todo{Items: make(map[string][]todoitem)}
	basicAuthMiddleware1 := NewBasicAuth("Todo Test Realm", basicAuthVerifier, todoApi.ServeHTTP)

	server := httptest.NewServer(basicAuthMiddleware1)
	defer server.Close()

	authRequest := basicAuthRequest("user1", "test")

	resp, err := authRequest("GET", server.URL+"/", nil)
	if err != nil {
		t.Error(err)
	}

	content, _ := ioutil.ReadAll(resp.Body)
	t.Logf("%s", content)

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected %d, got %d", http.StatusOK, resp.StatusCode)
	}

}

func TestPostTodo(t *testing.T) {
	todoApi := &todo{Items: make(map[string][]todoitem)}
	basicAuthMiddleware1 := NewBasicAuth("Todo Test Realm", basicAuthVerifier, todoApi.ServeHTTP)

	server := httptest.NewServer(basicAuthMiddleware1)
	defer server.Close()

	authRequest := basicAuthRequest("user1", "test")
	var body bytes.Buffer
	if err := json.NewEncoder(&body).Encode(requests.NewTodoRequest{Value: "Hello, WOrld"}); err != nil {
		t.Error(err)
	}
	resp, err := authRequest("POST", server.URL+"/", &body)
	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("expected %d. got %d", http.StatusCreated, resp.StatusCode)
	}
}
