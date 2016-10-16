package payment

import (
	"appengine"
	"appengine/datastore"
	"encoding/json"
	"model/vehicle"
	"model/order"
	"errors"
	"time"
	"net/http"
	"model/user"
	"appengine/urlfetch"
    "strconv"
)

type Payment struct {
	QrCode      string `json:"qr_code" datastore:"-"`
	OrderId     string `json:"order_id"`
	Price       float64 `json:"price"`
	//CreditCard CreditCard `json:"credit_card"`
	PaymentType string `json:"payment_type"`
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

func New(r *http.Request, vehicleId string) (*Payment, error) {
	c := appengine.NewContext(r)
	payment := new(Payment)

	if err := json.NewDecoder(r.Body).Decode(&payment); err != nil {
		return nil, err
	}

	orders := []order.Order{}
	keys, err := datastore.NewQuery("Order").Filter("VehicleId =", vehicleId).GetAll(c, &orders)

	if err != nil {
		return nil, err
	}

	Order := orders[0]
	Order.OrderId = keys[0].StringID()
	v, err := vehicle.GetOne(c, r, Order.VehicleId)

	if err != nil {
		return nil, err
	}

	if v.QrCode != payment.QrCode {
		return nil, errors.New("QR codes do not match")
	}

	payment.Price = v.PricePerHour * (float64(time.Now().Hour()) - float64(Order.OrderDate.Hour() + 1))

	if err := payment.save(c); err != nil {
		return nil, err
	}

	return payment, nil
}

func Accept(r *http.Request, vehicleId string) error {
	c := appengine.NewContext(r)
    payment := new(Payment)

    if err := json.NewDecoder(r.Body).Decode(&payment); err != nil {
        return err
    }

    orders := []order.Order{}
    keys, err := datastore.NewQuery("Order").Filter("VehicleId =", vehicleId).GetAll(c, &orders)

    if err != nil {
        return err
    }

    orders[0].OrderId = keys[0].StringID()
    v, err := vehicle.GetOne(c, r, orders[0].VehicleId)

    if err != nil {
        return err
    }

    if v.QrCode != payment.QrCode {
        return nil, errors.New("QR codes do not match")
    }

	payment.Price = v.PricePerHour * (float64(time.Now().Hour()) - float64(orders[0].OrderDate.Hour()) + 1)

	Vehicle, err := vehicle.GetOne(c, r, vehicleId)

	if err != nil {
		return err
	}

	owner, err := user.Get(c, Vehicle.Owner)

	if err != nil {
		return err
	}

	_, e := urlfetch.Client(c).Get("https://sasa-bus.appspot.com/accept/" + owner.Token + "/" + strconv.FormatFloat(payment.Price, 'f', 6, 64))

	return e
}

func Notify(c appengine.Context, r *http.Request, vehicleId string) error {
	v, _ := vehicle.GetOne(c, r, vehicleId)
	u, _ := user.Get(c, v.Borrower)

	owner, err := user.Get(c, v.Owner)

	if err != nil {
		return err
	}

	_, e := urlfetch.Client(c).Get("https://sasa-bus.appspot.com/notify/" + owner.Token + "/" + u.Name + "/" + v.QrCode)

	return e
}
