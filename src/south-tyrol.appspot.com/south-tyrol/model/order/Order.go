package order

import (
    "time"
    "id"
    "appengine/datastore"
    "model/vehicle"
    "appengine"
    "net/http"
)

type Order struct {
    UserId     string `json:"user_id"`
    VehicleId  string `json:"user_id"`
    OrderId    string `json:"order_id" datastore:"-"`
    OrderDate  time.Time `json:"order_date"`
    BilledDate time.Time `json:"billed_at"`
}

func (order *Order) key(c appengine.Context) *datastore.Key {
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

/*func GetMulti(c appengine.Context, userId string) (*[]Order, error) {
    orders := []Order{}
    _, err := datastore.NewQuery("Order").Filter("UserId =", userId).GetAll(c, &orders)

    if err != nil {
        return nil, err
    }

    return orders, nil
}*/

func New(c appengine.Context, r *http.Request, vehicleId string, userId string) (*Order, error) {
    v, err := vehicle.GetOne(c, r, vehicleId)

    if err != nil {
        return nil, err
    }

    v.Available = false
    v.QrCode = id.Alphanumeric()
    v.Save(c)

    order := Order{
        OrderId: id.Alphanumeric(),
        OrderDate: time.Now(),
        UserId: userId,
        VehicleId: vehicleId,
    }

    if err := order.save(c); err != nil {
        return nil, err
    }

    return &order, nil
}

func Delete(r *http.Request, orderId string) (*Order, error) {
    c := appengine.NewContext(r)
    order, err := GetOne(c, orderId)

    if err != nil {
        return nil, err
    }

    v, err := vehicle.GetOne(c, r, order.VehicleId)

    if err != nil {
        return nil, err
    }

    v.Available = true
    v.QrCode = ""
    v.Save(c)

    err = datastore.Delete(c, order.key(c))

    if err != nil {
        return nil, err
    }

    return order, nil
}
