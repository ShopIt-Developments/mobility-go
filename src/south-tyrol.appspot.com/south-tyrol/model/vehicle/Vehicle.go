package vehicle

import (
    "io"

    "appengine"
    "appengine/datastore"
    "encoding/json"
    "util"
)

type Vehicle struct {
    VehicleId    string `json:"id" datastore:"-"`
    Name         string `json:"name"`
    LicencePlate string `json:"licence_plate,omitempty"`
    Type         string `json:"type"`
    Description  string `json:"description"`
    PricePerHour float64 `json:"price_per_hour"`
    Currency     string `json:"currency"`
    Availability string `json:"availability"`
    Latitude     float64 `json:"lat"`
    Longitude    float64 `json:"lng"`
    UserId       string `json:"user_id"`
    Image        string `json:"image"`
}

func (car *Vehicle) key(c appengine.Context) *datastore.Key {
    if car.VehicleId == "" {
        return datastore.NewKey(c, "Vehicle", util.Alphanumeric(), 0, nil)
    }

    return datastore.NewKey(c, "Vehicle", car.VehicleId, 0, nil)
}

func (car *Vehicle) save(c appengine.Context) error {
    k, err := datastore.Put(c, car.key(c), car)

    if err != nil {
        return err
    }

    car.VehicleId = k.StringID()

    return nil
}

func GetMy(c appengine.Context, userId string) ([]Vehicle, error) {
    q := datastore.NewQuery("Vehicle").Filter("UserId =", userId)

    var cars []Vehicle

    keys, err := q.GetAll(c, &cars)

    if err != nil {
        return nil, err
    }

    for i := 0; i < len(cars); i++ {
        cars[i].VehicleId = keys[i].StringID()
    }

    return cars, nil
}

func GetOne(c appengine.Context, carId string) (*Vehicle, error) {
    car := Vehicle{VehicleId: carId}

    k := car.key(c)
    err := datastore.Get(c, k, &car)

    if err != nil {
        return nil, err
    }

    car.VehicleId = k.StringID()

    return &car, nil
}

func GetAll(c appengine.Context) ([]Vehicle, error) {
    q := datastore.NewQuery("Vehicle")

    var cars []Vehicle

    keys, err := q.GetAll(c, &cars)

    if err != nil {
        return nil, err
    }

    for i := 0; i < len(cars); i++ {
        cars[i].VehicleId = keys[i].StringID()
    }

    return cars, nil
}

func New(c appengine.Context, r io.ReadCloser) (*Vehicle, error) {
    car := new(Vehicle)

    if err := json.NewDecoder(r).Decode(&car); err != nil {
        return nil, err
    }

    if err := car.save(c); err != nil {
        return nil, err
    }

    return car, nil
}

func Update(c appengine.Context, carId string, r io.ReadCloser) (*Vehicle, error) {
    car, err := GetOne(c, carId)

    if err != nil {
        return nil, err
    }

    if err := json.NewDecoder(r).Decode(&car); err != nil {
        return nil, err
    }

    if err := car.save(c); err != nil {
        return nil, err
    }

    return car, nil
}

func Delete(c appengine.Context, id string) (*Vehicle, error) {
    car, err := GetOne(c, id)

    if err != nil {
        return nil, err
    }

    err = datastore.Delete(c, car.key(c))

    if err != nil {
        return nil, err
    }

    return car, nil
}