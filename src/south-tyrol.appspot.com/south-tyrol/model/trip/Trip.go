package trip

import (
    "time"
    "io"
    "appengine/datastore"
    "appengine"
    "encoding/json"
    "id"
)

type Trip struct {
    TripId    string `json:"-" datastore:"-"`
    UserId    string `json:"-"`
    Departure time.Time `json:"departure"`
    Arrival   time.Time `json:"arrival"`
}

func (trip *Trip) key(c appengine.Context) *datastore.Key {
    return datastore.NewKey(c, "Trip", trip.TripId, 0, nil)
}

func (order *Trip) save(c appengine.Context) error {
    k, err := datastore.Put(c, order.key(c), order)

    if err != nil {
        return err
    }

    order.TripId = k.StringID()

    return nil
}

func New(c appengine.Context, r io.ReadCloser, userId string) (*Trip, error) {
    trip := new(Trip)

    if err := json.NewDecoder(r).Decode(&trip); err != nil {
        return nil, err
    }

    trip.TripId = id.Alphanumeric()
    trip.UserId = userId

    if err := trip.save(c); err != nil {
        return nil, err
    }

    return trip, nil
}