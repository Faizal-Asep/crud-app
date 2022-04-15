package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/Faizal-Asep/crud-app/service"
	"github.com/joho/godotenv"
)

var a service.App

const newsCreationQuery = `
CREATE TABLE IF NOT EXISTS news (
	id int(11) NOT NULL AUTO_INCREMENT,
	topic varchar(125) DEFAULT NULL,
	title varchar(255) DEFAULT NULL,
	content text,
	status enum('draf','publish','deleted') DEFAULT 'draf',
	PRIMARY KEY (id),
	UNIQUE KEY title (title)
);`

const tagsCreationQuery = `
CREATE TABLE IF NOT EXISTS tags (
	id int NOT NULL AUTO_INCREMENT,
	name varchar(75) NULL,
	status enum('publish', 'deleted') NOT NULL DEFAULT 'publish' ,
	PRIMARY KEY (id),
	UNIQUE INDEX(name)
);`

const newstagsCreationQuery = `
CREATE TABLE IF NOT EXISTS news_tags (
	news_id int(11) NOT NULL,
	tags_id int(11) NOT NULL,
	PRIMARY KEY (news_id,tags_id),
	CONSTRAINT news_tags_ibfk_1 FOREIGN KEY (news_id) REFERENCES news (id) ON DELETE CASCADE,
	CONSTRAINT news_tags_ibfk_2 FOREIGN KEY (tags_id) REFERENCES tags (id) ON DELETE CASCADE
);`

func TestMain(m *testing.M) {

	err := godotenv.Load("./config/.env")
	if err != nil {
		log.Fatal(err)
	}

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	a.Initialize()
	ensureTableExists()
	code := m.Run()
	// clearTable()
	os.Exit(code)
}

func ensureTableExists() {
	if _, err := a.DB.Exec(newsCreationQuery); err != nil {
		log.Fatal(err)
	}
	if _, err := a.DB.Exec(tagsCreationQuery); err != nil {
		log.Fatal(err)
	}
	if _, err := a.DB.Exec(newstagsCreationQuery); err != nil {
		log.Fatal(err)
	}
	clearTable()
}

func clearTable() {
	clearTags()
	clearNews()
}

func clearTags() {
	a.DB.Exec("DELETE FROM tags")
	a.DB.Exec("ALTER TABLE tags AUTO_INCREMENT=1;")
}

func clearNews() {
	a.DB.Exec("DELETE FROM news")
	a.DB.Exec("ALTER TABLE news AUTO_INCREMENT=1;")
	a.Redis.FlushAll(context.Background())
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func addtags(count int) {
	if count < 1 {
		count = 1
	}

	for i := 0; i < count; i++ {
		_, err := a.DB.Exec("INSERT INTO tags(name) VALUES(?)", "tag "+strconv.Itoa(i))
		if err != nil {
			fmt.Println(err)
		}
	}
}
func addnews(count int, topic string, status string) {
	if count < 1 {
		count = 1
	}

	for i := 0; i < count; i++ {
		_, err := a.DB.Exec("INSERT INTO news(topic,title,content,status) VALUES(?,?,?,?)", topic, "title "+strconv.Itoa(i), "content "+strconv.Itoa(i), status)
		if err != nil {
			fmt.Println(err)
		}
	}
}

//tags test
func TestEmptyTagsTable(t *testing.T) {
	clearTags()
	req, _ := http.NewRequest("GET", "/tags", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body != "[]" {
		t.Errorf("Expected an empty array. Got %s", body)
	}
}

func TestGetNonExistenttag(t *testing.T) {
	clearTags()
	req, _ := http.NewRequest("GET", "/tag/11", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "tag not found" {
		t.Errorf("Expected the 'error' key of the response to be set to 'tag not found'. Got '%s'", m["error"])
	}
}

func TestCreatetag(t *testing.T) {
	clearTags()
	var jsonStr = []byte(`{"name":"test tag"}`)
	req, _ := http.NewRequest("POST", "/tag", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)
	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["name"] != "test tag" {
		t.Errorf("Expected tag name to be 'test tag'. Got '%v'", m["name"])
	}

	if m["id"] != 1.0 {
		t.Errorf("Expected tag ID to be '1'. Got '%v'", m["id"])
	}
}

func TestGettag(t *testing.T) {
	clearTags()
	addtags(1)

	req, _ := http.NewRequest("GET", "/tag/1", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestUpdatetag(t *testing.T) {
	clearTags()
	addtags(1)

	req, _ := http.NewRequest("GET", "/tag/1", nil)
	response := executeRequest(req)
	var originaltag map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &originaltag)

	var jsonStr = []byte(`{"name":"test tag - updated name"}`)
	req, _ = http.NewRequest("PUT", "/tag/1", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response = executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["id"] != originaltag["id"] {
		t.Errorf("Expected the id to remain the same (%v). Got %v", originaltag["id"], m["id"])
	}

	if m["name"] == originaltag["name"] {
		t.Errorf("Expected the name to change from '%v' to '%v'. Got '%v'", originaltag["name"], m["name"], m["name"])
	}
}

func TestDeletetag(t *testing.T) {
	clearTags()
	addtags(1)

	req, _ := http.NewRequest("GET", "/tag/1", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("DELETE", "/tag/1", nil)
	response = executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("GET", "/tag/1", nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)
}

//news test
func TestEmptyNewsTable(t *testing.T) {
	clearNews()
	var jsonStr = []byte(``)
	req, _ := http.NewRequest("GET", "/news", bytes.NewBuffer(jsonStr))
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body != "[]" {
		t.Errorf("Expected an empty array. Got %s", body)
	}
}

func TestGetNonExistentNews(t *testing.T) {
	clearNews()
	req, _ := http.NewRequest("GET", "/news/11", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "news not found" {
		t.Errorf("Expected the 'error' key of the response to be set to 'tag not found'. Got '%s'", m["error"])
	}
}

func TestCreateNews(t *testing.T) {
	clearTags()
	addtags(2)
	clearNews()

	var jsonStr = []byte(`{"topic": "investment","title": "how to start investment","content": "news content","status": "draf","tags": [{"name": "tag 0"},{"name": "tag 1"}]}`)
	req, _ := http.NewRequest("POST", "/news", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)
	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["title"] != "how to start investment" {
		t.Errorf("Expected news name to be 'how to start investment'. Got '%v'", m["title"])
	}

	if m["id"] != 1.0 {
		t.Errorf("Expected tag ID to be '1'. Got '%v'", m["id"])
	}
}

func TestGetnews(t *testing.T) {
	clearNews()
	addnews(1, "investment", "draf")

	req, _ := http.NewRequest("GET", "/news/1", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestGetFilterTopicNews(t *testing.T) {
	clearNews()
	addnews(1, "investment", "draf")
	var jsonStr = []byte(`{ "topic": "investment"}`)
	req, _ := http.NewRequest("GET", "/news", bytes.NewBuffer(jsonStr))
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body == "[]" {
		t.Errorf("Expected not empty array. Got %s", body)
	}

	var arrm []map[string]string
	json.Unmarshal(response.Body.Bytes(), &arrm)
	for _, m := range arrm {
		if m["topic"] != "investment" {
			t.Errorf("Expected the 'topic' value 'investment'. Got '%s'", m["topic"])
		}
	}
}

func TestGetFilterStatusNews(t *testing.T) {
	clearNews()
	addnews(1, "investment", "publish")
	var jsonStr = []byte(`{ "status": "publish"}`)
	req, _ := http.NewRequest("GET", "/news", bytes.NewBuffer(jsonStr))
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body == "[]" {
		t.Errorf("Expected not empty array. Got %s", body)
	}

	var arrm []map[string]string
	json.Unmarshal(response.Body.Bytes(), &arrm)
	for _, m := range arrm {
		if m["status"] != "publish" {
			t.Errorf("Expected the 'status' value 'publish'. Got '%s'", m["status"])
		}
	}
}

func TestUpdateNews(t *testing.T) {
	clearNews()
	addnews(1, "investment", "draf")

	req, _ := http.NewRequest("GET", "/news/1", nil)
	response := executeRequest(req)
	var originalnews map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &originalnews)

	var jsonStr = []byte(`{"title":"test news - updated title"}`)
	req, _ = http.NewRequest("PUT", "/news/1", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response = executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["id"] != originalnews["id"] {
		t.Errorf("Expected the id to remain the same (%v). Got %v", originalnews["id"], m["id"])
	}

	if m["title"] == originalnews["title"] {
		t.Errorf("Expected the name to change from '%v' to '%v'. Got '%v'", originalnews["name"], m["name"], m["name"])
	}
}

func TestDeleteNews(t *testing.T) {
	clearNews()
	addnews(1, "investment", "draf")
	req, _ := http.NewRequest("GET", "/news/1", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("DELETE", "/news/1", nil)
	response = executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("GET", "/news/1", nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)
}
