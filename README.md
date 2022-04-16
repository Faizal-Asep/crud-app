# crud-app
simple crud using gorilla/mux golang

## API Documentation
Config
  - file -> config/.env
  - content :
```
APP_DB_HOST = 127.0.0.1
APP_DB_PORT = 3306
APP_DB_USERNAME = user_bareksa
APP_DB_PASSWORD = password_bareksa
APP_DB_NAME = bareksa

APP_REDIS_HOST = 127.0.0.1:6379
APP_REDIS_PASSWORD = eYVX7EwVmmxKPCDmwMtyKVge8oLd2t81
```
Unit test
  - file -> main_test.go
  - run test -> go test -v
  - result :
```
=== RUN   TestEmptyTagsTable
--- PASS: TestEmptyTagsTable (0.16s)
=== RUN   TestGetNonExistenttag
--- PASS: TestGetNonExistenttag (0.15s)
=== RUN   TestCreatetag
--- PASS: TestCreatetag (0.14s)
=== RUN   TestGettag
--- PASS: TestGettag (0.14s)
=== RUN   TestUpdatetag
--- PASS: TestUpdatetag (0.56s)
=== RUN   TestDeletetag
--- PASS: TestDeletetag (1.04s)
=== RUN   TestEmptyNewsTable
--- PASS: TestEmptyNewsTable (0.21s)
=== RUN   TestGetNonExistentNews
--- PASS: TestGetNonExistentNews (0.51s)
=== RUN   TestCreateNews
--- PASS: TestCreateNews (0.34s)
=== RUN   TestGetnews
--- PASS: TestGetnews (0.32s)
=== RUN   TestGetFilterTopicNews
--- PASS: TestGetFilterTopicNews (0.23s)
=== RUN   TestGetFilterStatusNews
--- PASS: TestGetFilterStatusNews (0.56s)
=== RUN   TestUpdateNews
--- PASS: TestUpdateNews (0.33s)
=== RUN   TestDeleteNews
--- PASS: TestDeleteNews (0.39s)
PASS
ok      github.com/Faizal-Asep/crud-app 5.554s
```

### API
  - Schema
```
type Tag struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type News struct {
	ID      int64  `json:"id"`
	Topic   string `json:"topic"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Status  string `json:"status"`
	Tags    []Tag  `json:"tags"`
}

type Newsfilter struct {
	Status string `json:"status,omitempty"`
	Topic  string `json:"topic,omitempty"`
}
```
  - End Point
```
// CRUD Tags
GET     http://localhost:8081/tags 
POST    http://localhost:8081/tag
GET     http://localhost:8081/tag/{id}
PUT     http://localhost:8081/tag/{id}
DELETE  http://localhost:8081/tag/{id}

// CRUD News
GET     http://localhost:8081/news 
POST    http://localhost:8081/news
GET     http://localhost:8081/news/{id}
PUT     http://localhost:8081/news/{id}
DELETE  http://localhost:8081/news/{id}
```
