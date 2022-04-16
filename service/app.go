package service

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	f "github.com/Faizal-Asep/crud-app/function"
	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
	Redis  *redis.Client
}

func (a *App) Initialize() {

	a.Redis = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("APP_REDIS_HOST"),
		Password: os.Getenv("APP_REDIS_PASSWORD"),
		DB:       0,
	})
	key := "app/start"
	data := "Hello Redis"
	ttl := time.Duration(3) * time.Second
	op1 := a.Redis.Set(context.Background(), key, data, ttl)
	if err := op1.Err(); err != nil {
		log.Fatal(err)
	}

	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", os.Getenv("APP_DB_USERNAME"), os.Getenv("APP_DB_PASSWORD"), os.Getenv("APP_DB_HOST"), os.Getenv("APP_DB_PORT"), os.Getenv("APP_DB_NAME"))
	var err error
	a.DB, err = sql.Open("mysql", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	a.Router = mux.NewRouter()

	a.initializeRoutes()
}

func (a *App) Run(addr string) {

	Logfile := f.NewLoger("./logs/service_.log", (24 * time.Hour), (30 * 24 * time.Hour))
	router := handlers.CombinedLoggingHandler(Logfile, a.Router)
	log.SetOutput(Logfile)

	srv := &http.Server{
		Addr:         "0.0.0.0:" + addr,
		WriteTimeout: time.Minute * time.Duration(30),
		ReadTimeout:  time.Minute * 30,
		IdleTimeout:  time.Minute * 60,
		Handler:      router,
	}

	go func() {
		log.Println("starting")
		err := srv.ListenAndServe()
		if err != nil {
			log.Fatal(err)
		}

	}()

	var wait time.Duration
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c

	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()

	srv.Shutdown(ctx)
	log.Println("shutting down")
}

func (a *App) initializeRoutes() {
	a.Router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Handle 404
		f.ErrorRespond(w, http.StatusNotFound, "404 page not found")

	})
	a.Router.MethodNotAllowedHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Handle 405
		f.ErrorRespond(w, http.StatusMethodNotAllowed, "405 method not allowed")
	})
	//tag route
	a.Router.HandleFunc("/tags", a.getTags).Methods("GET")
	a.Router.HandleFunc("/tag", a.createTag).Methods("POST")
	a.Router.HandleFunc("/tag/{id:[0-9]+}", a.getTag).Methods("GET")
	a.Router.HandleFunc("/tag/{id:[0-9]+}", a.updateTag).Methods("PUT")
	a.Router.HandleFunc("/tag/{id:[0-9]+}", a.deleteTag).Methods("DELETE")

	//news route
	a.Router.HandleFunc("/news", a.listNews).Methods("GET")
	a.Router.HandleFunc("/news", a.createNews).Methods("POST")
	a.Router.HandleFunc("/news/{id:[0-9]+}", a.getNews).Methods("GET")
	a.Router.HandleFunc("/news/{id:[0-9]+}", a.updateNews).Methods("PUT")
	a.Router.HandleFunc("/news/{id:[0-9]+}", a.deleteNews).Methods("DELETE")
}
