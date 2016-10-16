package endpoint

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"issue"
	"encoding/json"
	"model/payment"
	"appengine"
	"model/order"
)

type Payments struct {
	Router *httprouter.Router
}

func (*Payments) Scan(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	newPayment, err := payment.New(appengine.NewContext(r), r.Body, p.ByName("order_id"))
	issue.Handle(w, err, http.StatusBadRequest)

	data, err := json.Marshal(newPayment)
	issue.Handle(w, err, http.StatusInternalServerError)

	w.Write(data)
}

func (*Payments) Accept(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	if err := payment.Accept(r.Body); err != nil {
		issue.Handle(w, err, http.StatusBadRequest)
		return
	}

	if _, err := order.Delete(r, p.ByName("order_id")); err != nil {
		issue.Handle(w, err, http.StatusInternalServerError)
	}
}