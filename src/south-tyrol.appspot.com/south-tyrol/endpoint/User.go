package endpoint

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"issue"
	"encoding/json"
	"appengine"
	"model/user"

	"errors"
	"strconv"
)

type User struct {
	Router *httprouter.Router
}

func (*User) Get(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	userID := r.Header.Get("Authorization")

	if userID == "" {
		issue.Handle(w, errors.New("Unauthorized"), http.StatusUnauthorized)
		return
	}

	entity, err := user.Get(appengine.NewContext(r), userID)
	issue.Handle(w, err, http.StatusBadRequest)

	data, err := json.Marshal(entity)
	issue.Handle(w, err, http.StatusInternalServerError)

	w.Write(data)
}

func (*User) Add(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	userID := r.Header.Get("Authorization")

	if userID == "" {
		issue.Handle(w, errors.New("Unauthorized"), http.StatusUnauthorized)
		return
	}

	entity, err := user.New(appengine.NewContext(r), r.Body, userID)
	issue.Handle(w, err, http.StatusBadRequest)

	data, err := json.Marshal(entity)
	issue.Handle(w, err, http.StatusInternalServerError)

	w.Write(data)
}

func (*User) AddPoints(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	userId := r.Header.Get("Authorization")

	if userId == "" {
		issue.Handle(w, errors.New("Unauthorized"), http.StatusUnauthorized)
		return
	}

	points, err := strconv.ParseInt(p.ByName("points"), 10, 64)

	if err != nil {
		issue.Handle(w, err, http.StatusBadRequest)
		return
	}

	if points < 0 {
		issue.Handle(w, errors.New("Points cannot be negative"), http.StatusBadRequest)
	}

	entity, err := user.AddPoints(appengine.NewContext(r), userId, points)
	issue.Handle(w, err, http.StatusBadRequest)

	data, err := json.Marshal(entity)
	issue.Handle(w, err, http.StatusInternalServerError)

	w.Write(data)
}

func (*User) Update(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	userID := r.Header.Get("Authorization")

	if userID == "" {
		issue.Handle(w, errors.New("Unauthorized"), http.StatusUnauthorized)
		return
	}

	if _, err := user.Update(appengine.NewContext(r), userID, r.Body); err != nil {
		issue.Handle(w, err, http.StatusBadRequest)
	}

	w.WriteHeader(http.StatusNoContent)
}

func (*User) Delete(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	userID := r.Header.Get("Authorization")

	if userID == "" {
		issue.Handle(w, errors.New("Unauthorized"), http.StatusUnauthorized)
		return
	}

	if _, err := user.Delete(appengine.NewContext(r), userID); err != nil {
		issue.Handle(w, err, http.StatusBadRequest)
	}

	w.WriteHeader(http.StatusNoContent)
}