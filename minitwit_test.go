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

func logout() *httptest.ResponseRecorder {
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

func TestRegister(t *testing.T) {
	w1 := register("user1", "default", "", "")
	assert.Contains(t, w1.Body.String(), "You were successfully registered and can login now")

	w2 := register("user1", "default", "", "")
	assert.Contains(t, w2.Body.String(), "The username is already taken")

	w3 := register("", "default", "", "")
	assert.Contains(t, w3.Body.String(), "You have to enter a username")

	w4 := register("meh", "", "", "")
	assert.Contains(t, w4.Body.String(), "You have to enter a password")

	w5 := register("meh", "x", "y", "")
	assert.Contains(t, w5.Body.String(), "The two passwords do not match")

	w6 := register("meh", "default", "", "broken")
	assert.Contains(t, w6.Body.String(), "You have to enter a valid email address")
}

func TestLoginLogout(t *testing.T) {
	w1 := register_and_login("user1", "default")
	assert.Contains(t, w1.Body.String(), "You were logged in")

	w2 := logout()
	assert.Contains(t, w2.Body.String(), "You were logged out")

	w3 := login("user1", "wrongpassword")
	assert.Contains(t, w3.Body.String, "Invalid password")

	w4 := login("user2", "wrongpassword")
	assert.Contains(t, w4.Body.String, "Invalid username")
}

func TestMessageRecording(t *testing.T) {
	register_and_login("foo", "default")

	w1 := add_message("test message 1")
	assert.Contains(t, w1.Body.String, "Your message was recorded")

	w2 := add_message("<test message 2>")
	assert.Contains(t, w2.Body.String, "Your message was recorded")

	w3 := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	router.ServeHTTP(w3, req)
	assert.Contains(t, w3.Body.String, "test message 1")
	assert.Contains(t, w3.Body.String, "<test message 2>") // OR: "&lt;test message 2&gt"
}

func TestTimeline(t *testing.T) {
	register_and_login("foo", "default")
	add_message("the message by foo")
	logout()
	register_and_login("bar", "default")
	add_message("the message by bar")

	w1 := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/public", nil)
	router.ServeHTTP(w1, req)
	assert.Contains(t, w1.Body.String, "the message by foo")
	assert.Contains(t, w1.Body.String, "the message by bar")

	w2 := httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/", nil)
	router.ServeHTTP(w2, req)
	assert.NotContains(t, w2.Body.String, "the message by foo")
	assert.Contains(t, w2.Body.String, "the message by bar")

	w3 := httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/foo/follow", nil)
	router.ServeHTTP(w3, req)
	assert.Contains(t, w3.Body.String, "You are now following \"foo\"") // OR: "You are now following &#34;foo&#34;"

	w4 := httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/", nil)
	router.ServeHTTP(w4, req)
	assert.Contains(t, w4.Body.String, "the message by foo")
	assert.Contains(t, w4.Body.String, "the message by bar")

	w5 := httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/bar", nil)
	router.ServeHTTP(w5, req)
	assert.NotContains(t, w5.Body.String, "the message by foo")
	assert.Contains(t, w5.Body.String, "the message by bar")

	w6 := httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/fpo", nil)
	router.ServeHTTP(w6, req)
	assert.Contains(t, w6.Body.String, "the message by foo")
	assert.NotContains(t, w6.Body.String, "the message by bar")

	w7 := httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/foo/funollow", nil)
	router.ServeHTTP(w7, req)
	assert.Contains(t, w7.Body.String, "You are no longer following \"foo\"") // OR: "You are no longer following &#34;foo&#34;"

	w8 := httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/", nil)
	router.ServeHTTP(w8, req)
	assert.NotContains(t, w8.Body.String, "the message by foo")
	assert.Contains(t, w8.Body.String, "the message by bar")
}
