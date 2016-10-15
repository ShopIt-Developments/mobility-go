package endpoint

import (
    "github.com/julienschmidt/httprouter"
    "net/http"
    "model/order"
    "issue"
    "encoding/json"
    "appengine"
)

type Order struct {
    Router *httprouter.Router
}

//noinspection GoUnusedParameter
func (*Order) GetAll(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
    w.Header().Set("Content-Type", "application/json")

    orders, err := order.GetAll(appengine.NewContext(r))
    issue.Handle(w, err, http.StatusBadRequest)

    data, err := json.Marshal(orders)
    issue.Handle(w, err, http.StatusInternalServerError)

    w.Write(data)
}

//noinspection GoUnusedParameter
func (*Order) Add(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
    w.Header().Set("Content-Type", "application/json")

    entity, err := order.New(appengine.NewContext(r), r.Body)
    issue.Handle(w, err, http.StatusBadRequest)

    data, err := json.Marshal(entity)
    issue.Handle(w, err, http.StatusInternalServerError)

    w.Write(data)
}

func (*Order) GetOne(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
    w.Header().Set("Content-Type", "application/json")

    entity, err := order.GetOne(appengine.NewContext(r), p.ByName("order_id"))
    issue.Handle(w, err, http.StatusBadRequest)

    data, err := json.Marshal(entity)
    issue.Handle(w, err, http.StatusInternalServerError)

    w.Write(data)
}

func (*Order) Update(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
    w.Header().Set("Content-Type", "application/json")

    if _, err := order.Update(appengine.NewContext(r), p.ByName("order_id"), r.Body); err != nil {
        issue.Handle(w, err, http.StatusBadRequest)
    }

    w.WriteHeader(http.StatusNoContent)
}

func (*Order) Delete(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
    w.Header().Set("Content-Type", "application/json")

    if _, err := order.Delete(appengine.NewContext(r), p.ByName("order_id")); err != nil {
        issue.Handle(w, err, http.StatusBadRequest)
    }

    w.WriteHeader(http.StatusNoContent)
}

func (*Order) Accept(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
    w.Header().Set("Content-Type", "application/json")

    entity, err := order.GetOne(appengine.NewContext(r), p.ByName("order_id"))

    if err != nil {
        issue.Handle(w, err, http.StatusBadRequest)
    }

    entity.Accepted = true

    w.WriteHeader(http.StatusNoContent)
}