package endpoint

import (
    "github.com/julienschmidt/httprouter"
    "net/http"
    "model/car"
    "issue"
    "encoding/json"
    "appengine"
)

type Car struct {
    Router *httprouter.Router
}

func (*Car) GetAll(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
    w.Header().Set("Content-Type", "application/json")

    cars, err := car.GetAll(appengine.NewContext(r))
    issue.Handle(w, err, http.StatusBadRequest)

    data, err := json.Marshal(cars)
    issue.Handle(w, err, http.StatusInternalServerError)

    w.Write(data)
}

func (*Car) Add(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
    w.Header().Set("Content-Type", "application/json")

    entity, err := car.New(appengine.NewContext(r), r.Body)
    issue.Handle(w, err, http.StatusBadRequest)

    data, err := json.Marshal(entity)
    issue.Handle(w, err, http.StatusInternalServerError)

    w.Write(data)
}

func (*Car) GetOne(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
    w.Header().Set("Content-Type", "application/json")

    entity, err := car.GetOne(appengine.NewContext(r), p.ByName("car_id"))
    issue.Handle(w, err, http.StatusBadRequest)

    data, err := json.Marshal(entity)
    issue.Handle(w, err, http.StatusInternalServerError)

    w.Write(data)
}

func (*Car) Update(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
    w.Header().Set("Content-Type", "application/json")

    if _, err := car.Update(appengine.NewContext(r), p.ByName("car_id"), r.Body); err != nil {
        issue.Handle(w, err, http.StatusBadRequest)
    }

    w.WriteHeader(http.StatusNoContent)
}

func (*Car) Delete(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
    w.Header().Set("Content-Type", "application/json")

    if _, err := car.Delete(appengine.NewContext(r), p.ByName("car_id")); err != nil {
        issue.Handle(w, err, http.StatusBadRequest)
    }

    w.WriteHeader(http.StatusNoContent)
}
