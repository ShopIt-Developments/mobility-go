package purchase

import (
	"net/http"
	"github.com/julienschmidt/httprouter"

	"endpoint"
    "encoding/json"
    "strconv"
    "strings"
    "model/sasa"
    "model/vehicle"
    "appengine"
)

func init() {
	router := httprouter.New()

    router.GET("/mobility/vehicles/available", rt)

    c := endpoint.Car{Router: router}
    c.Router.GET("/mobility/vehicles/my", c.GetOwn)
    c.Router.GET("/mobility/vehicles", c.GetAll)
    c.Router.GET("/mobility/vehicle/:vehicle_id", c.GetOne)
    c.Router.POST("/mobility/vehicle", c.Add)
    c.Router.PUT("/mobility/vehicles/:vehicle_id", c.Update)
    c.Router.DELETE("/mobility/vehicles/:vehicle_id", c.Delete)

	u := endpoint.User{Router: router}
	u.Router.GET("/user", u.Get)
	u.Router.POST("/user", u.Add)
	u.Router.PUT("/user", u.Update)
	u.Router.DELETE("/user", u.Delete)

	http.Handle("/", router)
}

func rt(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
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