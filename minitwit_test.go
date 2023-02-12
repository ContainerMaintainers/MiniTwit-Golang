package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ContainerMaintainers/MiniTwit-Golang/database"
	"github.com/ContainerMaintainers/MiniTwit-Golang/entities"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var router *gin.Engine

// ------ Initialization ------ //

func init() {
	router = setupRouter()
	database.ConnectToTestDatabase()
	database.MigrateEntities()
}

// ------ Helper Functions ------ //

func login(username string, password string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	body := map[string]interface{}{
		"username": username,
		"password": password,
	}
	json_body, err := json.Marshal(body)

	if err != nil {
		log.Fatalln(err)
	}

	req, _ := http.NewRequest("POST", "/login", bytes.NewReader(json_body))
	router.ServeHTTP(w, req)
	return w
}

func register(username string, password string, password2 string, email string) *httptest.ResponseRecorder {

	if password2 == "" {
		password2 = password
	}
	if email == "" {
		email = fmt.Sprintf("%s@example.com", username)
	}

	w := httptest.NewRecorder()
	body := map[string]interface{}{
		"username":  username,
		"password":  password,
		"password2": password2,
		"email":     email,
	}
	json_body, err := json.Marshal(body)

	if err != nil {
		log.Fatalln(err)
	}

	req, _ := http.NewRequest("POST", "/register", bytes.NewReader(json_body))
	router.ServeHTTP(w, req)
	return w
}

func register_and_login(username string, password string) *httptest.ResponseRecorder {
	register(username, password, "", "")
	return login(username, password)
}

func logout(username string, password string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/logout", nil)
	router.ServeHTTP(w, req)
	return w
}

func add_message(text string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	body := map[string]interface{}{
		"text": text,
	}
	json_body, err := json.Marshal(body)

	if err != nil {
		log.Fatalln(err)
	}

	req, _ := http.NewRequest("POST", "/add_message", bytes.NewReader(json_body))
	router.ServeHTTP(w, req)
	return w
}

// ------ Tests ------ //

func TestPingRoute(t *testing.T) {

	user := entities.User{User_ID: 1, Username: "name", Email: "email", PW_Hash: "hash"}
	database.DB.Create(&user)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	router.ServeHTTP(w, req)

	var response map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		log.Fatalln(err)
	}

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "pong", response["message"])
}
