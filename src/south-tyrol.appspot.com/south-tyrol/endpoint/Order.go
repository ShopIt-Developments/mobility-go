package endpoint

import (
    "github.com/julienschmidt/httprouter"
    "net/http"
    "model/order"
    "issue"
    "encoding/json"
    "appengine"
    "id"
    "network"
)

type Order struct {
    Router *httprouter.Router
}

func (*Order) GetOne(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
    w.Header().Set("Content-Type", "application/json")

    entity, err := order.GetOne(appengine.NewContext(r), p.ByName("order_id"))
    issue.Handle(w, err, http.StatusBadRequest)

    data, err := json.Marshal(entity)
    issue.Handle(w, err, http.StatusInternalServerError)

    w.Write(data)
}

func (*Order) New(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
    w.Header().Set("Content-Type", "application/json")

    entity, err := order.New(appengine.NewContext(r), r, p.ByName("vehicle_id"), network.Authorization(w, r))
    issue.Handle(w, err, http.StatusBadRequest)

    data, err := json.Marshal(id.Id{Id: entity.OrderId})
    issue.Handle(w, err, http.StatusInternalServerError)

    w.Write(data)
}

func (*Order) Delete(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
    w.Header().Set("Content-Type", "application/json")

    if _, err := order.Delete(r, p.ByName("order_id")); err != nil {
        issue.Handle(w, err, http.StatusBadRequest)
    }

    w.WriteHeader(http.StatusNoContent)
}
