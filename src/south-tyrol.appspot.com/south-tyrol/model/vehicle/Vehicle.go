package vehicle

import (
    "appengine"
    "appengine/datastore"
    "encoding/json"
    "id"
    "storage"
    "net/http"
)

type Vehicle struct {
    Address      string `json:"address"`
    Availability string `json:"availability"`
    Available    bool `json:"available"`
    Currency     string `json:"currency"`
    Description  string `json:"description"`
    VehicleId    string `json:"id" datastore:"-"`
    Image        string `json:"image" datastore:"-"`
    Latitude     float64 `json:"lat"`
    LicencePlate string `json:"licence_plate,omitempty"`
    Longitude    float64 `json:"lng"`
    Name         string `json:"name"`
    PricePerHour float64 `json:"price_per_hour"`
    QrCode       string `json:"qr_code,omitempty"`
    Type         string `json:"type"`
    Owner        string `json:"owner"`
}

func (vehicle *Vehicle) key(c appengine.Context) *datastore.Key {
    return datastore.NewKey(c, "Vehicle", vehicle.VehicleId, 0, nil)
}

func (vehicle *Vehicle) Save(c appengine.Context) error {
    k, err := datastore.Put(c, vehicle.key(c), vehicle)

    if err != nil {
        return err
    }

    vehicle.VehicleId = k.StringID()

    return nil
}

func GetBooked(r *http.Request, userId string) ([]Vehicle, error) {
    c := appengine.NewContext(r)

    orders := []*order.Order{}
    _, err := datastore.NewQuery("Order").Filter("UserId =", userId).GetAll(c, &orders)

    if err != nil {
        return nil, err
    }

    keys := make([]*datastore.Key, len(orders))

    for i := 0; i < len(orders); i++ {
        keys = datastore.NewKey(c, "Vehicle", orders[i].VehicleId, 0, nil);
    }

    vehicles := []*Vehicle{}

    if err := datastore.GetMulti(c, &keys, &vehicles); err != nil {
        return nil, err
    }

    return vehicles, nil
}

func GetMy(c appengine.Context, userId string) ([]Vehicle, error) {
    q := datastore.NewQuery("Vehicle").Filter("Owner =", userId)

    vehicles := []Vehicle{}
    keys, err := q.GetAll(c, &vehicles)

    if err != nil {
        return nil, err
    }

    for i := 0; i < len(vehicles); i++ {
        vehicles[i].VehicleId = keys[i].StringID()
    }

    return vehicles, nil
}

func GetOne(c appengine.Context, r *http.Request, vehicleId string) (*Vehicle, error) {
    vehicle := Vehicle{VehicleId: vehicleId}

    k := vehicle.key(c)
    err := datastore.Get(c, k, &vehicle)

    if err != nil {
        return nil, err
    }

    file, err := storage.ReadFile(r, "images/vehicles/" + vehicleId + ".txt")

    if err != nil {
        return nil, err
    }

    vehicle.VehicleId = k.StringID()
    vehicle.Image = string(file)

    return &vehicle, nil
}

func GetAll(r *http.Request) ([]Vehicle, error) {
    q := datastore.NewQuery("Vehicle").Filter("Available =", true)

    vehicles := []Vehicle{}
    keys, err := q.GetAll(appengine.NewContext(r), &vehicles)

    if err != nil {
        return nil, err
    }

    for i := 0; i < len(vehicles); i++ {
        file, err := storage.ReadFile(r, "images/vehicles/" + vehicles[i].VehicleId + ".txt")

        if err != nil {
            return nil, err
        }

        vehicles[i].VehicleId = keys[i].StringID()
        vehicles[i].Image = string(file)
    }

    return vehicles, nil
}

func New(c appengine.Context, r *http.Request, userId string) (*Vehicle, error) {
    vehicle := new(Vehicle)

    if err := json.NewDecoder(r.Body).Decode(&vehicle); err != nil {
        return nil, err
    }

    vehicle.Available = true
    vehicle.Owner = userId
    vehicle.VehicleId = id.Alphanumeric()
    vehicle.QrCode = ""

    storage.WriteFile(r, "images/vehicles/" + vehicle.VehicleId + ".txt", vehicle.Image)

    if err := vehicle.Save(c); err != nil {
        return nil, err
    }

    return vehicle, nil
}

func Update(c appengine.Context, vehicleId string, r *http.Request) (*Vehicle, error) {
    vehicle, err := GetOne(c, r, vehicleId)

    if err != nil {
        return nil, err
    }

    if err := json.NewDecoder(r.Body).Decode(&vehicle); err != nil {
        return nil, err
    }

    if err := vehicle.Save(c); err != nil {
        return nil, err
    }

    return vehicle, nil
}

func Delete(c appengine.Context, r *http.Request, id string) (*Vehicle, error) {
    vehicle, err := GetOne(c, r, id)

    if err != nil {
        return nil, err
    }

    err = datastore.Delete(c, vehicle.key(c))

    if err != nil {
        return nil, err
    }

    return vehicle, nil
}