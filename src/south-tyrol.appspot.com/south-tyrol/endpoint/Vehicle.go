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
    "id"
    "network"
)

type Vehicle struct {
    Router *httprouter.Router
}

func (*Vehicle) GetAvailable(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
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

    vehicles, err := vehicle.GetAll(appengine.NewContext(r))
    issue.Handle(w, err, http.StatusBadRequest)

    data, _ := json.Marshal(vehicle.Vehicles{
        Vehicles: vehicles,
        Buses: buses,
    })

    w.Write(data)
}

func (*Vehicle) GetMy(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
    w.Header().Set("Content-Type", "application/json")

    vehicles, err := vehicle.GetMy(appengine.NewContext(r), r.Header.Get("Authorization"))
    issue.Handle(w, err, http.StatusBadRequest)

    data, err := json.Marshal(vehicles)
    issue.Handle(w, err, http.StatusInternalServerError)

    w.Write(data)
}

func (*Vehicle) GetAll(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
    w.Header().Set("Content-Type", "application/json")

    vehicles, err := vehicle.GetAll(appengine.NewContext(r))
    issue.Handle(w, err, http.StatusBadRequest)

    data, err := json.Marshal(vehicles)
    issue.Handle(w, err, http.StatusInternalServerError)

    w.Write(data)
}

func (*Vehicle) New(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
    w.Header().Set("Content-Type", "application/json")

    entity, err := vehicle.New(appengine.NewContext(r), r.Body, network.Authorization(w, r))
    issue.Handle(w, err, http.StatusBadRequest)

    data, err := json.Marshal(id.Id{Id: entity.VehicleId})
    issue.Handle(w, err, http.StatusInternalServerError)

    w.Write(data)
}

func (*Vehicle) GetOne(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
    w.Header().Set("Content-Type", "application/json")

    entity, err := vehicle.GetOne(appengine.NewContext(r), p.ByName("vehicle_id"))
    issue.Handle(w, err, http.StatusBadRequest)

    data, err := json.Marshal(entity)
    issue.Handle(w, err, http.StatusInternalServerError)

    w.Write(data)
}

func (*Vehicle) Update(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
    w.Header().Set("Content-Type", "application/json")

    if _, err := vehicle.Update(appengine.NewContext(r), p.ByName("vehicle_id"), r.Body); err != nil {
        issue.Handle(w, err, http.StatusBadRequest)
    }

    w.WriteHeader(http.StatusNoContent)
}

func (*Vehicle) Delete(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
    w.Header().Set("Content-Type", "application/json")

    if _, err := vehicle.Delete(appengine.NewContext(r), p.ByName("vehicle_id")); err != nil {
        issue.Handle(w, err, http.StatusBadRequest)
    }

    w.WriteHeader(http.StatusNoContent)
}
