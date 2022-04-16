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

func (a *App) getTags(w http.ResponseWriter, r *http.Request) {

	count, _ := strconv.Atoi(r.FormValue("count"))
	start, _ := strconv.Atoi(r.FormValue("start"))
	if count > 10 || count < 1 {
		count = 10
	}
	if start < 0 {
		start = 0
	}

	tags, err := m.ListTags(a.DB, start, count)
	if err != nil {
		f.ErrorRespond(w, http.StatusInternalServerError, err.Error())
		return
	}

	f.JSONRespond(w, http.StatusOK, tags)
}

func (a *App) getTag(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		f.ErrorRespond(w, http.StatusBadRequest, "Invalid tag ID")
		return
	}

	t := m.Tag{ID: int64(id)}
	if err := t.GetTag(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			f.ErrorRespond(w, http.StatusNotFound, "tag not found")
		default:

			f.ErrorRespond(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	f.JSONRespond(w, http.StatusOK, t)
}

func (a *App) createTag(w http.ResponseWriter, r *http.Request) {

	var t m.Tag
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&t); err != nil {
		f.ErrorRespond(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	if err := t.CreateTag(a.DB); err != nil {
		f.ErrorRespond(w, http.StatusInternalServerError, err.Error())
		return
	}

	f.JSONRespond(w, http.StatusCreated, t)
}

func (a *App) updateTag(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		f.ErrorRespond(w, http.StatusBadRequest, "Invalid tag ID")
		return
	}

	var t m.Tag
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&t); err != nil {
		f.ErrorRespond(w, http.StatusBadRequest, "Invalid resquest payload")
		return
	}
	defer r.Body.Close()
	t.ID = int64(id)

	if err := t.UpdateTag(a.DB); err != nil {
		f.ErrorRespond(w, http.StatusInternalServerError, err.Error())
		return
	}

	f.JSONRespond(w, http.StatusOK, t)
}

func (a *App) deleteTag(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		f.ErrorRespond(w, http.StatusBadRequest, "Invalid Tag ID")
		return
	}

	t := m.Tag{ID: int64(id)}
	if err := t.DeleteTag(a.DB); err != nil {
		f.ErrorRespond(w, http.StatusInternalServerError, err.Error())
		return
	}

	f.JSONRespond(w, http.StatusOK, map[string]string{"result": "success"})
}
