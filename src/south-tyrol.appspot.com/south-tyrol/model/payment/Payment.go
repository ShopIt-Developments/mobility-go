package payment

import (
	"appengine"
	"appengine/datastore"
	"io"
	"encoding/json"
	"model/vehicle"
	"model/order"
	"errors"
	"time"
	"github.com/google/go-gcm"
	"net/http"
)

type Payment struct {
	QrCode     string `json:"qr_code" datastore:"-"`
	OrderId    string `json:"order_id"`
	Price      float64 `json:"price"`
	CreditCard CreditCard `json:"credit_card"`
}

func (payment *Payment) key(c appengine.Context) *datastore.Key {
	return datastore.NewKey(c, "Payment", payment.QrCode, 0, nil)
}

func (payment *Payment) save(c appengine.Context) error {
	k, err := datastore.Put(c, payment.key(c), payment)

	if err != nil {
		return err
	}

	payment.QrCode = k.StringID()

	return nil
}

func New(c appengine.Context, r *http.Request, orderId string) (*Payment, error) {
	payment := new(Payment)

	if err := json.NewDecoder(r.Body).Decode(&payment); err != nil {
		return nil, err
	}

	o, err := order.GetOne(c, orderId)

	if err != nil {
		return nil, err
	}

	v, err := vehicle.GetOne(c, r, o.VehicleId)

	if err != nil {
		return nil, err
	}

	if v.QrCode != payment.QrCode {
		return nil, errors.New("Qr-Codes do not match!")
	}

	payment.Price = v.PricePerHour * (float64(time.Now().Unix()) - float64(o.BilledDate.Unix()) / float64(3600))

	if err := payment.save(c); err != nil {
		return nil, err
	}

	return payment, nil
}

func Accept(r io.ReadCloser) (error) {
	var tokens []Token

	if err := json.NewDecoder(r).Decode(&tokens); err != nil {
		return err
	}

	d := gcm.Data{"action": "successful"}

	var err error

	for i := 0; i < len(tokens) && err == nil; i++ {
		_, err = gcm.SendHttp("AIzaSyALYozs9Jc2TqRUVW7uecvstoBiR6PXfbs", gcm.HttpMessage{To: tokens[i].Token, Data:d})
	}
	return err
}
