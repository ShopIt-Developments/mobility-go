package endpoint

import (
    "github.com/julienschmidt/httprouter"
    "net/http"
    "model/vehicle"
    "issue"
    "encoding/json"
    "appengine"
    "model/sasa"
    "strconv"
    "strings"
)

type Car struct {
    Router *httprouter.Router
}

func (*Car) GetAvailable(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
    w.Header().Set("Content-Type", "application/json")

    response := sasa.ReadJsonFromUrl(w, r, "http://realtimetest.opensasa.info/positions")

    projection := new(sasa.Buses)
    json.NewDecoder(response.Body).Decode(projection)

    features := projection.Features
    buses := make([]sasa.RealtimeBus, len(features))

    for key, bus := range features {
        p := bus.Properties;

        vehicleId, _ := strconv.Atoi(strings.Split(p.VehicleId, " ")[0]);
        latitude, longitude := sasa.ToWgs84(bus.Geometry.Coordinates, 32)

        buses[key] = sasa.RealtimeBus{
            LineName: p.LineName,
            BusStop: p.BusStopName,
            HydrogenBus: vehicleId >= 428 && vehicleId <= 432,
            TripId: p.TripId,
            Latitude: latitude,
            Longitude: longitude,
        }
    }

    vehicles, _ := vehicle.GetAll(appengine.NewContext(r))

    data, _ := json.Marshal(vehicle.Vehicles{
        Vehicles: vehicles,
        Buses: buses,
    })

    w.Write(data)
}

func (*Car) GetMy(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
    w.Header().Set("Content-Type", "application/json")

    cars, err := vehicle.GetMy(appengine.NewContext(r), r.Header.Get("Authorization"))
    issue.Handle(w, err, http.StatusBadRequest)

    data, err := json.Marshal(cars)
    issue.Handle(w, err, http.StatusInternalServerError)

    w.Write(data)
}

func (*Car) GetAll(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
    w.Header().Set("Content-Type", "application/json")

    cars, err := vehicle.GetAll(appengine.NewContext(r))
    issue.Handle(w, err, http.StatusBadRequest)

    data, err := json.Marshal(cars)
    issue.Handle(w, err, http.StatusInternalServerError)

    w.Write(data)
}

func (*Car) Add(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
    w.Header().Set("Content-Type", "application/json")

    entity, err := vehicle.New(appengine.NewContext(r), r.Body)
    issue.Handle(w, err, http.StatusBadRequest)

    data, err := json.Marshal(entity)
    issue.Handle(w, err, http.StatusInternalServerError)

    w.Write(data)
}

func (*Car) GetOne(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
    w.Header().Set("Content-Type", "application/json")

    entity, err := vehicle.GetOne(appengine.NewContext(r), p.ByName("car_id"))
    issue.Handle(w, err, http.StatusBadRequest)

    data, err := json.Marshal(entity)
    issue.Handle(w, err, http.StatusInternalServerError)

    w.Write(data)
}

func (*Car) Update(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
    w.Header().Set("Content-Type", "application/json")

    if _, err := vehicle.Update(appengine.NewContext(r), p.ByName("car_id"), r.Body); err != nil {
        issue.Handle(w, err, http.StatusBadRequest)
    }

    w.WriteHeader(http.StatusNoContent)
}

func (*Car) Delete(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
    w.Header().Set("Content-Type", "application/json")

    if _, err := vehicle.Delete(appengine.NewContext(r), p.ByName("car_id")); err != nil {
        issue.Handle(w, err, http.StatusBadRequest)
    }

    w.WriteHeader(http.StatusNoContent)
}
