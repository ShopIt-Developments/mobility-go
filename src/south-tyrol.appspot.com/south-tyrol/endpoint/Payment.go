package endpoint

import (
    "github.com/julienschmidt/httprouter"
    "net/http"
    "issue"
    "encoding/json"
    "model/payment"
    "appengine"
    "model/order"
    "network"
)

type Payment struct {
    Router *httprouter.Router
}

func (*Payment) Scan(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
    w.Header().Set("Content-Type", "application/json")

    newPayment, err := payment.New(appengine.NewContext(r), r, p.ByName("order_id"))
    issue.Handle(w, err, http.StatusBadRequest)

    data, err := json.Marshal(newPayment)
    issue.Handle(w, err, http.StatusInternalServerError)

    w.Write(data)
}

func (*Payment) Accept(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
    w.Header().Set("Content-Type", "application/json")

    if err := payment.Accept(r, p.ByName("order_id"), network.Authorization(w, r)); err != nil {
        issue.Handle(w, err, http.StatusBadRequest)
        return
    }

    if _, err := order.Delete(r, p.ByName("order_id")); err != nil {
        issue.Handle(w, err, http.StatusInternalServerError)
    }
}

func (*Payment) Notify(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
    w.Header().Set("Content-Type", "application/json")

    err := payment.Notify(appengine.NewContext(r), r, p.ByName("order_id"))
    issue.Handle(w, err, http.StatusBadRequest)
}