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
    "network"
    "storage"
)

type Vehicle struct {
    Router *httprouter.Router
}

func (*Vehicle) GetBooked(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
    w.Header().Set("Content-Type", "application/json")

    vehicles, err := vehicle.GetBooked(r, network.Authorization(w, r))
    issue.Handle(w, err, http.StatusBadRequest)

    data, err := json.Marshal(&vehicles)
    issue.Handle(w, err, http.StatusInternalServerError)

    w.Write(data)
}

func (*Vehicle) GetBus(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    w.Header().Set("Content-Type", "application/json")

    response := sasa.ReadJsonFromUrl(w, r, "http://realtimetest.opensasa.info/positions")

    projection := new(sasa.Buses)
    json.NewDecoder(response.Body).Decode(projection)

    for _, bus := range projection.Features {
        p := bus.Properties;

        vehicleId, _ := strconv.Atoi(strings.Split(p.VehicleId, " ")[0]);

        if strconv.Itoa(vehicleId) == ps.ByName("id") {
            latitude, longitude := sasa.ToWgs84(bus.Geometry.Coordinates, 32)
            variant, _ := strconv.Atoi(strings.Split(p.Variant, " ")[0])

            hydrogen := vehicleId >= 428 && vehicleId <= 432
            lineName := p.LineName

            if hydrogen {
                lineName = "Hydrogen bus: " + lineName
            }

            bus := sasa.RealtimeBus{
                LineName: lineName,
                LineId: p.LineId,
                Variant: variant,
                BusStop: p.BusStopName,
                HydrogenBus: hydrogen,
                TripId: strconv.Itoa(p.TripId),
                Latitude: latitude,
                Longitude: longitude,
            }

            data, _ := json.Marshal(bus)
            w.Write(data)
            return
        }
    }

    w.WriteHeader(http.StatusBadRequest)
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
        variant, _ := strconv.Atoi(strings.Split(p.Variant, " ")[0])

        buses[key] = sasa.RealtimeBus{
            LineName: p.LineName,
            LineId: p.LineId,
            Variant: variant,
            BusStop: p.BusStopName,
            HydrogenBus: vehicleId >= 428 && vehicleId <= 432,
            TripId: strconv.Itoa(p.TripId),
            Latitude: latitude,
            Longitude: longitude,
        }
    }

    vehicles, err := vehicle.GetAll(r)
    issue.Handle(w, err, http.StatusBadRequest)

    data, _ := json.Marshal(vehicle.Vehicles{
        Vehicles: vehicles,
        Buses: buses,
    })

    w.Write(data)
}

func (*Vehicle) GetMy(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
    w.Header().Set("Content-Type", "application/json")

    vehicles, err := vehicle.GetMy(r, network.Authorization(w, r))
    issue.Handle(w, err, http.StatusBadRequest)

    data, err := json.Marshal(vehicles)
    issue.Handle(w, err, http.StatusInternalServerError)

    w.Write(data)
}

func (*Vehicle) GetAll(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
    w.Header().Set("Content-Type", "application/json")

    vehicles, err := vehicle.GetAll(r)
    issue.Handle(w, err, http.StatusBadRequest)

    data, err := json.Marshal(vehicles)
    issue.Handle(w, err, http.StatusInternalServerError)

    w.Write(data)
}

func (*Vehicle) New(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
    w.Header().Set("Content-Type", "application/json")

    entity, err := vehicle.New(appengine.NewContext(r), r, network.Authorization(w, r))
    issue.Handle(w, err, http.StatusBadRequest)

    data, err := json.Marshal(entity)
    issue.Handle(w, err, http.StatusInternalServerError)

    w.Write(data)
}

func (*Vehicle) GetOne(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
    w.Header().Set("Content-Type", "application/json")

    entity, err := vehicle.GetOne(appengine.NewContext(r), r, p.ByName("vehicle_id"))
    issue.Handle(w, err, http.StatusBadRequest)

    data, err := json.Marshal(entity)
    issue.Handle(w, err, http.StatusInternalServerError)

    w.Write(data)
}

func (*Vehicle) Update(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
    w.Header().Set("Content-Type", "application/json")

    if _, err := vehicle.Update(appengine.NewContext(r), p.ByName("vehicle_id"), r); err != nil {
        issue.Handle(w, err, http.StatusBadRequest)
    }

    w.WriteHeader(http.StatusNoContent)
}

func (*Vehicle) Delete(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
    w.Header().Set("Content-Type", "application/json")

    storage.DeleteFile(r, "images/vehicles/" + p.ByName("vehicle_id") + ".txt")

    if _, err := vehicle.Delete(appengine.NewContext(r), r, p.ByName("vehicle_id")); err != nil {
        issue.Handle(w, err, http.StatusBadRequest)
    }

    w.WriteHeader(http.StatusNoContent)
}
