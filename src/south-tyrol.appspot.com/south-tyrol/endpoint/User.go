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
	"network"
)

type User struct {
	Router *httprouter.Router
}

func (*User) Get(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	entity, err := user.Get(appengine.NewContext(r), network.Authorization(w, r))
	issue.Handle(w, err, http.StatusBadRequest)

	data, err := json.Marshal(entity)
	issue.Handle(w, err, http.StatusInternalServerError)

	w.Write(data)
}

func (*User) Add(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	entity, err := user.New(appengine.NewContext(r), r.Body, network.Authorization(w, r))
	issue.Handle(w, err, http.StatusBadRequest)

	data, err := json.Marshal(entity)
	issue.Handle(w, err, http.StatusInternalServerError)

	w.Write(data)
}

func (*User) AddPoints(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	points, err := strconv.ParseInt(p.ByName("points"), 10, 64)

	if err != nil {
		issue.Handle(w, err, http.StatusBadRequest)
		return
	}

	if points < 0 {
		issue.Handle(w, errors.New("Points cannot be negative"), http.StatusBadRequest)
		return
	}

	entity, err := user.AddPoints(appengine.NewContext(r), network.Authorization(w, r), points)
	issue.Handle(w, err, http.StatusBadRequest)

	data, err := json.Marshal(entity)
	issue.Handle(w, err, http.StatusInternalServerError)

	w.Write(data)
}

func (*User) GetPoints(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	entity, err := user.GetPoints(appengine.NewContext(r), network.Authorization(w, r))
	issue.Handle(w, err, http.StatusBadRequest)

	data, err := json.Marshal(entity)
	issue.Handle(w, err, http.StatusInternalServerError)

	w.Write(data)
}

func (*User) Update(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	if _, err := user.Update(appengine.NewContext(r), network.Authorization(w, r), r.Body); err != nil {
		issue.Handle(w, err, http.StatusBadRequest)
	}

	w.WriteHeader(http.StatusNoContent)
}

func (*User) Delete(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	if _, err := user.Delete(appengine.NewContext(r), network.Authorization(w, r)); err != nil {
		issue.Handle(w, err, http.StatusBadRequest)
	}

	w.WriteHeader(http.StatusNoContent)
}

func (*User) SetToken(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	if err := user.SetToken(appengine.NewContext(r), network.Authorization(w, r), p.ByName("token")); err != nil {
		issue.Handle(w, err, http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}