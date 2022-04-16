package service

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	f "github.com/Faizal-Asep/crud-app/function"
	m "github.com/Faizal-Asep/crud-app/model"
	"github.com/gorilla/mux"
)

func (a *App) listNews(w http.ResponseWriter, r *http.Request) {

	var filter m.Newsfilter
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&filter)

	defer r.Body.Close()

	cacheKey := f.CacheNewsKey(filter)
	exist, result, _ := f.GetCache(a.Redis, cacheKey)
	if exist {
		f.ByteRespond(w, http.StatusOK, []byte(result))
		return
	}

	news, err := m.ListNews(a.DB, filter)
	if err != nil {
		f.ErrorRespond(w, http.StatusInternalServerError, err.Error())
		return
	}

	f.JSONRespond(w, http.StatusOK, news)

	f.SetCache(a.Redis, cacheKey, news)
}

func (a *App) getNews(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		f.ErrorRespond(w, http.StatusBadRequest, "Invalid news ID")
		return
	}

	n := m.News{ID: int64(id)}
	if err := n.GetNews(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			f.ErrorRespond(w, http.StatusNotFound, "news not found")
		default:

			f.ErrorRespond(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	f.JSONRespond(w, http.StatusOK, n)
}

func (a *App) createNews(w http.ResponseWriter, r *http.Request) {

	var n m.News
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&n); err != nil {
		f.ErrorRespond(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	if err := n.CreateNews(a.DB); err != nil {
		f.ErrorRespond(w, http.StatusInternalServerError, err.Error())
		return
	}

	f.JSONRespond(w, http.StatusCreated, n)
}

func (a *App) deleteNews(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		f.ErrorRespond(w, http.StatusBadRequest, "Invalid news ID")
		return
	}

	n := m.News{ID: int64(id)}
	if err := n.DeleteNews(a.DB); err != nil {
		f.ErrorRespond(w, http.StatusInternalServerError, err.Error())
		return
	}

	f.JSONRespond(w, http.StatusOK, map[string]string{"result": "success"})
}

func (a *App) updateNews(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		f.ErrorRespond(w, http.StatusBadRequest, "Invalid news ID")
		return
	}

	var n m.News
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&n); err != nil {
		f.ErrorRespond(w, http.StatusBadRequest, "Invalid resquest payload")
		return
	}

	defer r.Body.Close()
	n.ID = int64(id)

	if err := n.UpdateNews(a.DB); err != nil {
		f.ErrorRespond(w, http.StatusInternalServerError, err.Error())
		return
	}

	f.JSONRespond(w, http.StatusOK, n)
}
