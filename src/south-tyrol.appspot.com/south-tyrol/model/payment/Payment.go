package payment

import (
    "appengine"
    "appengine/datastore"
    "encoding/json"
    "model/vehicle"
    "model/order"
    "errors"
    "time"
    "github.com/google/go-gcm"
    "net/http"
    "model/user"
)

const FCM_KEY = "AIzaSyALYozs9Jc2TqRUVW7uecvstoBiR6PXfbs"

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

func New(c appengine.Context, r *http.Request, vehicleId string) (*Payment, error) {
    payment := new(Payment)

    if err := json.NewDecoder(r.Body).Decode(&payment); err != nil {
        return nil, err
    }

    orders := []order.Order{}
    keys, err := datastore.NewQuery("Order").Filter("UserId =", vehicleId).GetAll(c, &orders)

    if err != nil {
        return nil, err
    }

    orders[0].OrderId = keys[0].StringID()
    v, err := vehicle.GetOne(c, r, orders[0].VehicleId)

    if err != nil {
        return nil, err
    }

    if v.QrCode != payment.QrCode {
        return nil, errors.New("QR codes do not match")
    }

    payment.Price = v.PricePerHour * (float64(time.Now().Unix()) - float64(orders[0].OrderDate.Unix()) / float64(3600))

    if err := payment.save(c); err != nil {
        return nil, err
    }

    return payment, nil
}

func Accept(r *http.Request, vehicleId string, userId string) (error) {
    c := appengine.NewContext(r)

    Vehicle, err := vehicle.GetOne(c, r, vehicleId)

    if err != nil {
        return err
    }

    owner, err := user.Get(c, Vehicle.Owner)

    if err != nil {
        return err
    }

    User, err := user.Get(c, userId)

    if err != nil {
        return err
    }

    d := gcm.Data{"data": map[string]string{"receiver": "action", "action": "payment_successful"}}

    _, err = gcm.SendHttp(FCM_KEY, gcm.HttpMessage{To: owner.Token, Data:d})
    _, err = gcm.SendHttp(FCM_KEY, gcm.HttpMessage{To: User.Token, Data:d})

    return err
}

func Notify(c appengine.Context, r *http.Request, vehicleId string) error {
    var tokens []Token

    if err := json.NewDecoder(r.Body).Decode(&tokens); err != nil {
        return err
    }

    v, _ := vehicle.GetOne(c, r, vehicleId)
    u, _ := user.Get(c, v.Borrower)
    d := gcm.Data{"data": map[string]string{"receiver": "action", "action": "payment_initiate", "username": u.Name, "qr_code": v.QrCode}}

    var err error

    for i := 0; i < len(tokens) && err == nil; i++ {
        _, err = gcm.SendHttp(FCM_KEY, gcm.HttpMessage{To: tokens[i].Token, Data: d})
    }

    return err
}