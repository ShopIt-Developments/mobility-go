package car

import (
    "io"

    "appengine"
    "appengine/datastore"
    "encoding/json"
)

type Car struct {
    LicensePlate string `json:"license_plate" datastore:"-"`
    Description  string `json:"description"`
    PricePerHour float64 `json:"price_per_hour"`
    Availability Availability `json:"availability"`
}

func (car *Car) key(c appengine.Context) *datastore.Key {
    return datastore.NewKey(c, "Car", car.LicensePlate, 0, nil)
}

func (car *Car) save(c appengine.Context) error {
    k, err := datastore.Put(c, car.key(c), car)

    if err != nil {
        return err
    }

    car.LicensePlate = k.StringID()

    return nil
}

func GetOne(c appengine.Context, carId string) (*Car, error) {
    car := Car{LicensePlate: carId}

    k := car.key(c)
    err := datastore.Get(c, k, &car)

    if err != nil {
        return nil, err
    }

    car.LicensePlate = k.StringID()

    return &car, nil
}

func GetAll(c appengine.Context) ([]Car, error) {
    q := datastore.NewQuery("Car")

    var cars []Car

    keys, err := q.GetAll(c, &cars)

    if err != nil {
        return nil, err
    }

    for i := 0; i < len(cars); i++ {
        cars[i].LicensePlate = keys[i].StringID()
    }

    return cars, nil
}

func New(c appengine.Context, r io.ReadCloser) (*Car, error) {
    car := new(Car)

    if err := json.NewDecoder(r).Decode(&car); err != nil {
        return nil, err
    }

    if err := car.save(c); err != nil {
        return nil, err
    }

    return car, nil
}

func Update(c appengine.Context, carId string, r io.ReadCloser) (*Car, error) {
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

func Delete(c appengine.Context, id string) (*Car, error) {
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