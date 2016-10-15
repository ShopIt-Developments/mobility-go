package order

import (
	"io"

	"appengine"
	"appengine/datastore"
	"errors"
	"encoding/json"
	"util"
	"time"
)

type Order struct {
	OrderId      string `json:"order_id" datastore:"-"`
	OrderDate    time.Time `json:"order_date"`
	DeliveryDate time.Time `json:"delivery_date"`
	Location     Location `json:"location"`
	Products     []Product `json:"products"`
	Currency     string `json:"currency"`
	UserId       string `json:"user_id"`
	Accepted     bool `json:"-"`
}

func (order *Order) key(c appengine.Context) *datastore.Key {
	if order.OrderId == "" {
		return datastore.NewKey(c, "Order", util.Alphanumeric(), 0, nil)
	}

	return datastore.NewKey(c, "Order", order.OrderId, 0, nil)
}

func (order *Order) save(c appengine.Context) error {
	k, err := datastore.Put(c, order.key(c), order)

	if err != nil {
		return err
	}

	order.OrderId = k.StringID()

	return nil
}

func GetOne(c appengine.Context, orderId string) (*Order, error) {
	order := Order{OrderId: orderId}

	k := order.key(c)
	err := datastore.Get(c, k, &order)

	if err != nil {
		return nil, err
	}

	order.OrderId = k.StringID()

	return &order, nil
}

func GetAll(c appengine.Context) ([]Order, error) {
	q := datastore.NewQuery("Order").Filter("Accepted =", false)

	var orders []Order

	keys, err := q.GetAll(c, &orders)

	if err != nil {
		return nil, err
	}

	for i := 0; i < len(orders); i++ {
		orders[i].OrderId = keys[i].StringID()
	}

	return orders, nil
}

func New(c appengine.Context, r io.ReadCloser) (*Order, error) {
	order := new(Order)

	if err := json.NewDecoder(r).Decode(&order); err != nil {
		return nil, err
	}

	if err := order.assertDeliveryDateInPast(); err != nil {
		return nil, err
	}

	order.OrderDate = time.Now()
	order.Accepted = false

	if err := order.save(c); err != nil {
		return nil, err
	}

	return order, nil
}

func Update(c appengine.Context, orderId string, r io.ReadCloser) (*Order, error) {
	order, err := GetOne(c, orderId)

	if order.Accepted {
		return nil, errors.New("An accepted order cannot be updated.")
	}

	if err != nil {
		return nil, err
	}

	if err := json.NewDecoder(r).Decode(&order); err != nil {
		return nil, err
	}

	if err := order.assertDeliveryDateInPast(); err != nil {
		return nil, err
	}

	order.OrderId = orderId
	order.OrderDate = time.Now()

	if err := order.save(c); err != nil {
		return nil, err
	}

	return order, nil
}

func Delete(c appengine.Context, id string) (*Order, error) {
	order, err := GetOne(c, id)

	if err != nil {
		return nil, err
	}

	err = datastore.Delete(c, order.key(c))

	if err != nil {
		return nil, err
	}

	return order, nil
}

func (order *Order) assertDeliveryDateInPast() error {
	if order.DeliveryDate.Before(time.Now()) {
		return errors.New("The delivery date cannot be past.")
	}

	return nil
}

func (order *Order) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		OrderId      string `json:"order_id" datastore:"-"`
		OrderDate    string `json:"order_date"`
		DeliveryDate string `json:"delivery_date"`
		Currency     string `json:"currency"`
		Location     Location `json:"location"`
		Products     []Product `json:"products"`
		UserId       string `json:"user_id"`
	}{
		OrderId: order.OrderId,
		OrderDate: order.OrderDate.Format(time.RFC3339),
		DeliveryDate: order.DeliveryDate.Format(time.RFC3339),
		Location: order.Location,
		Products: order.Products,
		Currency: order.Currency,
		UserId: order.UserId,
	})
}