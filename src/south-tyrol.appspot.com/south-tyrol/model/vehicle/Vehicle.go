package vehicle

import (
    "io"

    "appengine"
    "appengine/datastore"
    "encoding/json"
    "id"
)

type Vehicle struct {
    Address      string `json:"address"`
    Availability string `json:"availability"`
    Available    bool `json:"available"`
    Currency     string `json:"currency"`
    Description  string `json:"description"`
    VehicleId    string `json:"id" datastore:"-"`
    Image        string `json:"image"`
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

func GetMy(c appengine.Context, userId string) ([]Vehicle, error) {
    q := datastore.NewQuery("Vehicle").Filter("UserId =", userId)

    var vehicles []Vehicle

    keys, err := q.GetAll(c, &vehicles)

    if err != nil {
        return nil, err
    }

    for i := 0; i < len(vehicles); i++ {
        vehicles[i].VehicleId = keys[i].StringID()
    }

    return vehicles, nil
}

func GetOne(c appengine.Context, vehicleId string) (*Vehicle, error) {
    vehicle := Vehicle{VehicleId: vehicleId}

    k := vehicle.key(c)
    err := datastore.Get(c, k, &vehicle)

    if err != nil {
        return nil, err
    }

    vehicle.VehicleId = k.StringID()

    return &vehicle, nil
}

func GetAll(c appengine.Context) ([]Vehicle, error) {
    q := datastore.NewQuery("Vehicle").Filter("Available =", true)

    var vehicles []Vehicle

    keys, err := q.GetAll(c, &vehicles)

    if err != nil {
        return nil, err
    }

    for i := 0; i < len(vehicles); i++ {
        vehicles[i].VehicleId = keys[i].StringID()
    }

    return vehicles, nil
}

func New(c appengine.Context, r io.ReadCloser, userId string) (*Vehicle, error) {
    vehicle := new(Vehicle)

    if err := json.NewDecoder(r).Decode(&vehicle); err != nil {
        return nil, err
    }

    vehicle.Owner = userId
    vehicle.VehicleId = id.Alphanumeric()
    vehicle.QrCode = ""

    if err := vehicle.Save(c); err != nil {
        return nil, err
    }

    return vehicle, nil
}

func Update(c appengine.Context, vehicleId string, r io.ReadCloser) (*Vehicle, error) {
    vehicle, err := GetOne(c, vehicleId)

    if err != nil {
        return nil, err
    }

    if err := json.NewDecoder(r).Decode(&vehicle); err != nil {
        return nil, err
    }

    if err := vehicle.Save(c); err != nil {
        return nil, err
    }

    return vehicle, nil
}

func Delete(c appengine.Context, id string) (*Vehicle, error) {
    vehicle, err := GetOne(c, id)

    if err != nil {
        return nil, err
    }

    err = datastore.Delete(c, vehicle.key(c))

    if err != nil {
        return nil, err
    }

    return vehicle, nil
}