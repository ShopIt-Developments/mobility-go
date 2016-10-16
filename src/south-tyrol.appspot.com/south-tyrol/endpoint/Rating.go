package endpoint

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"issue"
	"encoding/json"
	"appengine"
	"model/rating"
	"strconv"
	"network"
)

type Rating struct {
	Router *httprouter.Router
}

func (*Rating) Add(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	entity, err := rating.New(appengine.NewContext(r), r.Body, network.Authorization(w, r))
	issue.Handle(w, err, http.StatusBadRequest)

	data, err := json.Marshal(entity)
	issue.Handle(w, err, http.StatusInternalServerError)

	w.Write(data)
}

func (*Rating) GetOne(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	ratingId, err := strconv.ParseInt(p.ByName("rating_id"), 10, 64)

	if err != nil {
		issue.Handle(w, err, http.StatusNotFound)
		return
	}

	entity, err := rating.GetOne(appengine.NewContext(r), ratingId)
	issue.Handle(w, err, http.StatusBadRequest)

	data, err := json.Marshal(entity)
	issue.Handle(w, err, http.StatusInternalServerError)

	w.Write(data)
}

func (*Rating) GetRatings(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	ratings, err := rating.GetRatings(appengine.NewContext(r), network.Authorization(w, r))
	issue.Handle(w, err, http.StatusBadRequest)

	data, err := json.Marshal(ratings)
	issue.Handle(w, err, http.StatusInternalServerError)

	w.Write(data)
}

func (*Rating) GetRated(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	ratings, err := rating.GetRated(appengine.NewContext(r), network.Authorization(w, r))
	issue.Handle(w, err, http.StatusBadRequest)

	data, err := json.Marshal(ratings)
	issue.Handle(w, err, http.StatusInternalServerError)

	w.Write(data)
}

func (*Rating) Update(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	ratingId, err := strconv.ParseInt(p.ByName("rating_id"), 10, 64)

	if err != nil {
		issue.Handle(w, err, http.StatusNotFound)
	}

	if _, err := rating.Update(appengine.NewContext(r), ratingId, r.Body); err != nil {
		issue.Handle(w, err, http.StatusBadRequest)
	}

	w.WriteHeader(http.StatusNoContent)
}

func (*Rating) Delete(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	ratingId, err := strconv.ParseInt(p.ByName("rating_id"), 10, 64)

	if err != nil {
		issue.Handle(w, err, http.StatusNotFound)
	}

	if _, err := rating.Delete(appengine.NewContext(r), ratingId); err != nil {
		issue.Handle(w, err, http.StatusBadRequest)
	}

	w.WriteHeader(http.StatusNoContent)
}