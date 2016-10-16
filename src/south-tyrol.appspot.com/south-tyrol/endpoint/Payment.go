package endpoint

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"issue"
	"encoding/json"
	"model/payment"
	"appengine"
	"model/order"
	"appengine/datastore"
)

type Payment struct {
	Router *httprouter.Router
}

func (*Payment) Scan(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	newPayment, err := payment.New(r, p.ByName("vehicle_id"))
	issue.Handle(w, err, http.StatusBadRequest)

	data, err := json.Marshal(newPayment)
	issue.Handle(w, err, http.StatusInternalServerError)

	w.Write(data)
}

func (*Payment) Accept(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	if err := payment.Accept(r, p.ByName("vehicle_id")); err != nil {
		issue.Handle(w, err, http.StatusBadRequest)
		return
	}

	q := datastore.NewQuery("Order").Filter("VehicleId =", p.ByName("vehicle_id"))

	orders := []order.Order{}
	keys, err := q.GetAll(appengine.NewContext(r), &orders)

	if err != nil {
		issue.Handle(w, err, http.StatusInternalServerError)
	}

	for i := 0; i < len(orders); i++ {
		orders[i].OrderId = keys[i].StringID()
	}

	if _, err := order.Delete(r, orders[0].OrderId); err != nil {
		issue.Handle(w, err, http.StatusInternalServerError)
	}
}

func (*Payment) Notify(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	payment.Notify(appengine.NewContext(r), r, p.ByName("vehicle_id"))
}